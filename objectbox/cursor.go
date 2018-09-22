package objectbox

/*
#cgo LDFLAGS: -lobjectbox
#include <stdlib.h>
#include "objectbox.h"
*/
import "C"

import (
	"github.com/google/flatbuffers/go"
	"unsafe"
)

type Cursor struct {
	cursor  *C.OBX_cursor
	binding ObjectBinding
	fbb     *flatbuffers.Builder
}

func (cursor *Cursor) Destroy() (err error) {
	rc := C.obx_cursor_destroy(cursor.cursor)
	cursor.cursor = nil
	if rc != 0 {
		err = createError()
	}
	return
}

func (cursor *Cursor) Get(id uint64) (object interface{}, err error) {
	bytes, err := cursor.GetBytes(id)
	if bytes == nil || err != nil {
		return
	}
	return cursor.binding.ToObject(bytes), nil
}

func (cursor *Cursor) GetAll() (slice interface{}, err error) {
	var bytes []byte
	binding := cursor.binding
	slice = cursor.binding.MakeSlice(16)

	for bytes, err = cursor.First(); bytes != nil; bytes, err = cursor.Next() {
		if err != nil || bytes == nil {
			slice = nil
			return
		}
		object := binding.ToObject(bytes)
		slice = binding.AppendToSlice(slice, object)
	}
	return slice, nil
}

func (cursor *Cursor) GetBytes(id uint64) (bytes []byte, err error) {
	var data *C.void
	var dataSize C.size_t
	dataPtr := unsafe.Pointer(data) // Need ptr to an unsafe ptr here
	rc := C.obx_cursor_get(cursor.cursor, C.uint64_t(id), &dataPtr, &dataSize)
	if rc != 0 {
		if rc != 404 {
			err = createError()
		}
		return
	}
	bytes = C.GoBytes(dataPtr, C.int(dataSize))
	return
}

func (cursor *Cursor) First() (bytes []byte, err error) {
	var data *C.void
	var dataSize C.size_t
	dataPtr := unsafe.Pointer(data) // Need ptr to an unsafe ptr here
	rc := C.obx_cursor_first(cursor.cursor, &dataPtr, &dataSize)
	if rc != 0 {
		if rc != 404 {
			err = createError()
		}
		return
	}
	bytes = C.GoBytes(dataPtr, C.int(dataSize))
	return
}

func (cursor *Cursor) Next() (bytes []byte, err error) {
	var data *C.void
	var dataSize C.size_t
	dataPtr := unsafe.Pointer(data) // Need ptr to an unsafe ptr here
	rc := C.obx_cursor_next(cursor.cursor, &dataPtr, &dataSize)
	if rc != 0 {
		if rc != 404 {
			err = createError()
		}
		return
	}
	bytes = C.GoBytes(dataPtr, C.int(dataSize))
	return
}

func (cursor *Cursor) Count() (count uint64, err error) {
	var cCount C.uint64_t
	rc := C.obx_cursor_count(cursor.cursor, &cCount)
	if rc != 0 {
		err = createError()
		return
	}
	return uint64(cCount), nil
}

func (cursor *Cursor) Put(object interface{}) (id uint64, err error) {
	idFromObject, err := cursor.binding.GetId(object)
	if err != nil {
		return
	}
	checkForPreviousValue := idFromObject != 0
	id, err = cursor.IdForPut(idFromObject)
	if err != nil {
		return
	}
	cursor.binding.Flatten(object, cursor.fbb, id)
	return id, cursor.finishInternalFbbAndPut(id, checkForPreviousValue)
}

func (cursor *Cursor) finishInternalFbbAndPut(id uint64, checkForPreviousObject bool) (err error) {
	fbb := cursor.fbb
	fbb.Finish(fbb.EndObject())
	bytes := fbb.FinishedBytes()

	cCheckPrevious := 0
	if checkForPreviousObject {
		cCheckPrevious = 1
	}
	rc := C.obx_cursor_put(cursor.cursor, C.uint64_t(id), unsafe.Pointer(&bytes[0]), C.size_t(len(bytes)),
		C.int(cCheckPrevious))
	if rc != 0 {
		err = createError()
	}

	// Reset to have a clear state for the next caller
	fbb.Reset()

	return
}

func (cursor *Cursor) IdForPut(idCandidate uint64) (id uint64, err error) {
	id = uint64(C.obx_cursor_id_for_put(cursor.cursor, C.uint64_t(idCandidate)))
	if id == 0 {
		err = createError()
	}
	return
}

func (cursor *Cursor) RemoveAll() (err error) {
	rc := C.obx_cursor_remove_all(cursor.cursor)
	if rc != 0 {
		err = createError()
	}
	return
}

func (cursor *Cursor) FindByString(propertyId uint, value string) (bytesArray *BytesArray, err error) {
	cvalue := C.CString(value)
	defer C.free(unsafe.Pointer(cvalue))

	cBytesArray := C.obx_query_by_string(cursor.cursor, C.uint32_t(propertyId), cvalue)
	if cBytesArray == nil {
		err = createError()
		return
	}
	size := int(cBytesArray.size)
	plainBytesArray := make([][]byte, size)
	if size > 0 {
		goBytesArray := (*[1 << 30]C.OBX_bytes)(unsafe.Pointer(cBytesArray.bytes))[:size:size]
		for i := 0; i < size; i++ {
			cBytes := goBytesArray[i]
			dataBytes := C.GoBytes(cBytes.data, C.int(cBytes.size))
			plainBytesArray[i] = dataBytes
		}
	}

	return &BytesArray{plainBytesArray, cBytesArray}, nil
}

func (cursor *Cursor) cBytesArrayToObjects(cBytesArray *C.OBX_bytes_array) (slice interface{}) {
	bytesArray := cBytesArrayToGo(cBytesArray)
	defer bytesArray.Destroy()
	return cursor.bytesArrayToObjects(bytesArray)
}

func (cursor *Cursor) bytesArrayToObjects(bytesArray *BytesArray) (slice interface{}) {
	slice = cursor.binding.MakeSlice(len(bytesArray.BytesArray))
	for _, bytesData := range bytesArray.BytesArray {
		object := cursor.binding.ToObject(bytesData)
		slice = cursor.binding.AppendToSlice(slice, object)
	}
	return
}

func cBytesArrayToGo(cBytesArray *C.OBX_bytes_array) (bytesArray *BytesArray) {
	size := int(cBytesArray.size)
	plainBytesArray := make([][]byte, size)
	if size > 0 {
		goBytesArray := (*[1 << 30]C.OBX_bytes)(unsafe.Pointer(cBytesArray.bytes))[:size:size]
		for i := 0; i < size; i++ {
			cBytes := goBytesArray[i]
			dataBytes := C.GoBytes(cBytes.data, C.int(cBytes.size))
			plainBytesArray[i] = dataBytes
		}
	}
	return &BytesArray{plainBytesArray, cBytesArray}
}
