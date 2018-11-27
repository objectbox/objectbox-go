/*
 * Copyright 2018 ObjectBox Ltd. All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package objectbox

/*
#cgo LDFLAGS: -lobjectbox
#include <stdlib.h>
#include "objectbox.h"
*/
import "C"

import (
	"reflect"
	"unsafe"

	"github.com/google/flatbuffers/go"
)

// Internal: won't be publicly exposed in a future version!
type cursor struct {
	cursor  *C.OBX_cursor
	binding ObjectBinding
	fbb     *flatbuffers.Builder
}

func (cursor *cursor) Close() error {
	rc := C.obx_cursor_close(cursor.cursor)
	cursor.cursor = nil
	if rc != 0 {
		return createError()
	}
	return nil
}

func (cursor *cursor) Get(id uint64) (object interface{}, err error) {
	bytes, err := cursor.getBytes(id)
	if bytes == nil || err != nil {
		return
	}
	return cursor.binding.ToObject(bytes), nil
}

func (cursor *cursor) GetAll() (slice interface{}, err error) {
	cBytesArray := C.obx_cursor_get_all(cursor.cursor)
	if cBytesArray == nil {
		return nil, createError()
	}
	return cursor.cBytesArrayToObjects(cBytesArray), nil
}

func (cursor *cursor) getBytes(id uint64) (bytes []byte, err error) {
	var data *C.void
	var dataSize C.size_t
	dataPtr := unsafe.Pointer(data) // Need ptr to an unsafe ptr here
	rc := C.obx_cursor_get(cursor.cursor, C.obx_id(id), &dataPtr, &dataSize)
	if rc != 0 {
		if rc != 404 {
			err = createError()
		}
		return
	}
	bytes = C.GoBytes(dataPtr, C.int(dataSize))
	return
}

func (cursor *cursor) First() (bytes []byte, err error) {
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

func (cursor *cursor) Next() (bytes []byte, err error) {
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

func (cursor *cursor) Count() (count uint64, err error) {
	var cCount C.uint64_t
	rc := C.obx_cursor_count(cursor.cursor, &cCount)
	if rc != 0 {
		err = createError()
		return
	}
	return uint64(cCount), nil
}

func (cursor *cursor) Put(object interface{}) (id uint64, err error) {
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

func (cursor *cursor) finishInternalFbbAndPut(id uint64, checkForPreviousObject bool) (err error) {
	fbb := cursor.fbb
	fbb.Finish(fbb.EndObject())
	bytes := fbb.FinishedBytes()

	rc := C.obx_cursor_put(cursor.cursor,
		C.obx_id(id), unsafe.Pointer(&bytes[0]), C.size_t(len(bytes)), C.bool(checkForPreviousObject))
	if rc != 0 {
		err = createError()
	}

	// Reset to have a clear state for the next caller
	fbb.Reset()

	return
}

func (cursor *cursor) IdForPut(idCandidate uint64) (id uint64, err error) {
	id = uint64(C.obx_cursor_id_for_put(cursor.cursor, C.obx_id(idCandidate)))
	if id == 0 {
		err = createError()
	}
	return
}

func (cursor *cursor) Remove(id uint64) (err error) {
	rc := C.obx_cursor_remove(cursor.cursor, C.obx_id(id))
	if rc != 0 {
		err = createError()
	}
	return
}

func (cursor *cursor) RemoveAll() (err error) {
	rc := C.obx_cursor_remove_all(cursor.cursor)
	if rc != 0 {
		err = createError()
	}
	return
}

func (cursor *cursor) cBytesArrayToObjects(cBytesArray *C.OBX_bytes_array) (slice interface{}) {
	bytesArray := cBytesArrayToGo(cBytesArray)
	defer bytesArray.free()
	return cursor.bytesArrayToObjects(bytesArray)
}

func (cursor *cursor) bytesArrayToObjects(bytesArray *BytesArray) (slice interface{}) {
	slice = cursor.binding.MakeSlice(len(bytesArray.BytesArray))
	for _, bytesData := range bytesArray.BytesArray {
		object := cursor.binding.ToObject(bytesData)
		slice = cursor.binding.AppendToSlice(slice, object)
	}
	return
}

func cBytesArrayToGo(cBytesArray *C.OBX_bytes_array) *BytesArray {
	size := int(cBytesArray.count)
	plainBytesArray := make([][]byte, size)
	if size > 0 {
		// Previous alternative without reflect:
		//   https://github.com/golang/go/wiki/cgo#turning-c-arrays-into-go-slices (2012)
		//   On a RPi 3, the size with 1<<30 did not work, but 1<<27 did
		// Raul measured both variants and did notice a visible perf impact (Go 1.11.2)
		var goBytesArray []C.OBX_bytes
		header := (*reflect.SliceHeader)(unsafe.Pointer(&goBytesArray))
		*header = reflect.SliceHeader{Data: uintptr(unsafe.Pointer(cBytesArray.bytes)), Len: size, Cap: size}
		for i := 0; i < size; i++ {
			cBytes := goBytesArray[i]
			dataBytes := C.GoBytes(cBytes.data, C.int(cBytes.size))
			plainBytesArray[i] = dataBytes
		}
	}
	return &BytesArray{plainBytesArray, cBytesArray}
}
