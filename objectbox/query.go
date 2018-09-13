package objectbox

/*
#cgo LDFLAGS: -lobjectboxc
#include "objectbox.h"
*/
import "C"
import "unsafe"

type Query struct {
	cquery *C.OBX_query
}

func (query *Query) Destroy() (err error) {
	rc := C.obx_query_close(query.cquery)
	query.cquery = nil
	if rc != 0 {
		err = createError()
	}
	return
}

func (query *Query) Find(cursor *Cursor) (slice interface{}, err error) {
	bytesArray, err := query.FindBytes(cursor)
	if err != nil {
		return
	}
	defer bytesArray.Destroy()

	slice = cursor.binding.MakeSlice(len(bytesArray.BytesArray))
	for _, bytesData := range bytesArray.BytesArray {
		object := cursor.binding.ToObject(bytesData)
		slice = cursor.binding.AppendToSlice(slice, object)
	}
	return
}

func (query *Query) FindBytes(cursor *Cursor) (bytesArray *BytesArray, err error) {
	cBytesArray := C.obx_query_find(query.cquery, cursor.cursor)
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
