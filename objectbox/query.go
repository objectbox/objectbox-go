package objectbox

/*
#cgo LDFLAGS: -lobjectbox
#include <stdlib.h>
#include "objectbox.h"
*/
import "C"
import "unsafe"

type Query struct {
	cquery    *C.OBX_query
	typeId    TypeId
	objectBox *ObjectBox
}

func (query *Query) Close() (err error) {
	rc := C.obx_query_close(query.cquery)
	query.cquery = nil
	if rc != 0 {
		err = createError()
	}
	return
}

func (query *Query) Find() (slice interface{}, err error) {
	err = query.objectBox.runWithCursor(query.typeId, true, func(cursor *Cursor) error {
		var errInner error
		slice, errInner = query.find(cursor)
		return errInner
	})
	return
}

func (query *Query) find(cursor *Cursor) (slice interface{}, err error) {
	bytesArray, err := query.findBytes(cursor)
	if err != nil {
		return
	}
	defer bytesArray.Free()
	return cursor.bytesArrayToObjects(bytesArray), nil
}

// Won't be public in the future
func (query *Query) FindBytes() (bytesArray *BytesArray, err error) {
	err = query.objectBox.runWithCursor(query.typeId, true, func(cursor *Cursor) error {
		var errInner error
		bytesArray, errInner = query.findBytes(cursor)
		return errInner
	})
	return
}

func (query *Query) findBytes(cursor *Cursor) (bytesArray *BytesArray, err error) {
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
	rc := C.obx_query_string_param(query.cquery, C.obx_schema_id(propertyId), cvalue)
	if rc != 0 {
		return createError()
	}
	return
}

func (query *Query) SetParamInt(propertyId TypeId, value int64) (err error) {
	rc := C.obx_query_int_param(query.cquery, C.obx_schema_id(propertyId), C.int64_t(value))
	if rc != 0 {
		return createError()
	}
	return
}
