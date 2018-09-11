package objectbox

/*
#cgo LDFLAGS: -L ${SRCDIR}/libs -lobjectboxc
#include <stdlib.h>
#include <string.h>
#include "objectbox.h"
*/
import "C"

import (
	"github.com/google/flatbuffers/go"
	"unsafe"
)

type Transaction struct {
	txn       *C.OB_txn
	objectBox *ObjectBox
}

func (txn *Transaction) Destroy() (err error) {
	rc := C.ob_txn_destroy(txn.txn)
	txn.txn = nil
	if rc != 0 {
		err = createError()
	}
	return
}

func (txn *Transaction) Abort() (err error) {
	rc := C.ob_txn_abort(txn.txn)
	if rc != 0 {
		err = createError()
	}
	return
}

func (txn *Transaction) Commit() (err error) {
	rc := C.ob_txn_commit(txn.txn)
	if rc != 0 {
		err = createError()
	}
	return
}

func (txn *Transaction) createCursor(typeId TypeId, binding ObjectBinding) (*Cursor, error) {
	ccursor := C.ob_cursor_create(txn.txn, C.uint(typeId))
	if ccursor == nil {
		return nil, createError()
	}
	return &Cursor{ccursor, binding, flatbuffers.NewBuilder(512)}, nil
}

func (txn *Transaction) CursorForName(entitySchemaName string) (*Cursor, error) {
	binding := txn.objectBox.getBindingByName(entitySchemaName)
	cname := C.CString(entitySchemaName)
	defer C.free(unsafe.Pointer(cname))

	ccursor := C.ob_cursor_create2(txn.txn, cname)
	if ccursor == nil {
		return nil, createError()
	}
	return &Cursor{ccursor, binding, flatbuffers.NewBuilder(512)}, nil
}
