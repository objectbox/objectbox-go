package objectbox

/*
#cgo LDFLAGS: -lobjectbox
#include <stdlib.h>
#include "objectbox.h"
*/
import "C"
import "unsafe"

type Query struct {
	cquery *C.OBX_query
}

func (query *Query) Close() (err error) {
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
	defer bytesArray.Free()
	return cursor.bytesArrayToObjects(bytesArray), nil
}

func (query *Query) FindBytes(cursor *Cursor) (bytesArray *BytesArray, err error) {
	cBytesArray := C.obx_query_find(query.cquery, cursor.cursor)
	if cBytesArray == nil {
		err = createError()
		return
	}
	return cBytesArrayToGo(cBytesArray), nil
}

func (query *Query) SetParamString(propertyId TypeId, value string) (err error) {
	cvalue := C.CString(value)
	defer C.free(unsafe.Pointer(cvalue))
	rc := C.obx_query_string_param(query.cquery, C.uint32_t(propertyId), cvalue)
	if rc != 0 {
		return createError()
	}
	return
}

func (query *Query) SetParamInt(propertyId TypeId, value int64) (err error) {
	rc := C.obx_query_int_param(query.cquery, C.uint32_t(propertyId), C.int64_t(value))
	if rc != 0 {
		return createError()
	}
	return
}
