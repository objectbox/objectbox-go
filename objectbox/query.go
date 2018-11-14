package objectbox

/*
#cgo LDFLAGS: -lobjectbox
#include "objectbox.h"
*/
import "C"

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
	defer bytesArray.Close()
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
