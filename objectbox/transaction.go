package objectbox

/*
#cgo LDFLAGS: -lobjectbox
#include <stdlib.h>
#include "objectbox.h"
*/
import "C"

import (
	"unsafe"

	"github.com/google/flatbuffers/go"
)

type Transaction struct {
	txn       *C.OBX_txn
	objectBox *ObjectBox
}

func (txn *Transaction) Close() (err error) {
	rc := C.obx_txn_close(txn.txn)
	txn.txn = nil
	if rc != 0 {
		err = createError()
	}
	return
}

func (txn *Transaction) Abort() (err error) {
	rc := C.obx_txn_abort(txn.txn)
	if rc != 0 {
		err = createError()
	}
	return
}

func (txn *Transaction) Commit() (err error) {
	rc := C.obx_txn_commit(txn.txn)
	if rc != 0 {
		err = createError()
	}
	return
}

func (txn *Transaction) createCursor(typeId TypeId, binding ObjectBinding) (*Cursor, error) {
	ccursor := C.obx_cursor_create(txn.txn, C.obx_schema_id(typeId))
	if ccursor == nil {
		return nil, createError()
	}
	return &Cursor{ccursor, binding, flatbuffers.NewBuilder(512)}, nil
}

func (txn *Transaction) CursorForName(entitySchemaName string) (*Cursor, error) {
	binding := txn.objectBox.getBindingByName(entitySchemaName)
	cname := C.CString(entitySchemaName)
	defer C.free(unsafe.Pointer(cname))

	ccursor := C.obx_cursor_create2(txn.txn, cname)
	if ccursor == nil {
		return nil, createError()
	}
	return &Cursor{ccursor, binding, flatbuffers.NewBuilder(512)}, nil
}
