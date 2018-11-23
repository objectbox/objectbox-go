// Package objectbox provides a super-fast, light-weight object persistence framework.
package objectbox

/*
#cgo LDFLAGS: -lobjectbox
#include <stdlib.h>
#include "objectbox.h"
*/
import "C"

import (
	"errors"
	"fmt"
	"runtime"
	"strconv"

	"github.com/google/flatbuffers/go"
)

const Unavailable = flatbuffers.UOffsetT(0)

//noinspection GoUnusedConst
const (
	DebugFlags_LOG_TRANSACTIONS_READ  = 1
	DebugFlags_LOG_TRANSACTIONS_WRITE = 2
	DebugFlags_LOG_QUERIES            = 4
	DebugFlags_LOG_QUERY_PARAMETERS   = 8
	DebugFlags_LOG_ASYNC_QUEUE        = 16
)

type TypeId uint32

// An ObjectBinding provides an interface for various object types to be included in the model
type ObjectBinding interface {
	AddToModel(model *Model)
	GetId(object interface{}) (id uint64, err error)
	Flatten(object interface{}, fbb *flatbuffers.Builder, id uint64)
	ToObject(bytes []byte) interface{}
	MakeSlice(capacity int) interface{}
	AppendToSlice(slice interface{}, object interface{}) (sliceNew interface{})
}

type ObjectBox struct {
	store          *C.OBX_store
	bindingsById   map[TypeId]ObjectBinding
	bindingsByName map[string]ObjectBinding
}

type BytesArray struct {
	BytesArray  [][]byte
	cBytesArray *C.OBX_bytes_array
}

type TxnFun func(transaction *Transaction) (err error)
type CursorFun func(cursor *Cursor) (err error)

func (ob *ObjectBox) Close() {
	storeToClose := ob.store
	ob.store = nil
	if storeToClose != nil {
		C.obx_store_close(storeToClose)
	}
}

func (ob *ObjectBox) BeginTxn() (txn *Transaction, err error) {
	var ctxn = C.obx_txn_begin(ob.store)
	if ctxn == nil {
		return nil, createError()
	}
	return &Transaction{ctxn, ob}, nil
}

func (ob *ObjectBox) BeginTxnRead() (txn *Transaction, err error) {
	var ctxn = C.obx_txn_begin_read(ob.store)
	if ctxn == nil {
		return nil, createError()
	}
	return &Transaction{ctxn, ob}, nil
}

func (ob *ObjectBox) RunInTxn(readOnly bool, txnFun TxnFun) (err error) {
	runtime.LockOSThread()
	var txn *Transaction
	if readOnly {
		txn, err = ob.BeginTxnRead()
	} else {
		txn, err = ob.BeginTxn()
	}
	if err != nil {
		runtime.UnlockOSThread()
		return
	}

	//fmt.Println(">>> START TX")
	//os.Stdout.Sync()

	err = txnFun(txn)

	//fmt.Println("<<< END TX")
	//os.Stdout.Sync()

	if !readOnly && err == nil {
		err = txn.Commit()
	}
	err2 := txn.Close()
	if err == nil {
		err = err2
	}
	runtime.UnlockOSThread()

	//fmt.Println("<<< END TX Close")
	//os.Stdout.Sync()

	return
}

func (ob ObjectBox) getBindingById(typeId TypeId) ObjectBinding {
	binding := ob.bindingsById[typeId]
	if binding == nil {
		// Configuration error by the dev, OK to panic
		panic("Configuration error; no binding registered for type ID " + strconv.Itoa(int(typeId)))
	}
	return binding
}

func (ob ObjectBox) getBindingByName(typeName string) ObjectBinding {
	binding := ob.bindingsByName[typeName]
	if binding == nil {
		// Configuration error by the dev, OK to panic
		panic("Configuration error; no binding registered for type name " + typeName)
	}
	return binding
}

func (ob *ObjectBox) runWithCursor(typeId TypeId, readOnly bool, cursorFun CursorFun) (err error) {
	binding := ob.getBindingById(typeId)
	return ob.RunInTxn(readOnly, func(txn *Transaction) (err error) {
		cursor, err := txn.createCursor(typeId, binding)
		if err != nil {
			return
		}
		//fmt.Println(">>> START C")
		//os.Stdout.Sync()

		err = cursorFun(cursor)

		//fmt.Println("<<< END C")
		//os.Stdout.Sync()

		err2 := cursor.Close()
		if err == nil {
			err = err2
		}
		return
	})
}

func (ob *ObjectBox) SetDebugFlags(flags uint) (err error) {
	rc := C.obx_store_debug_flags(ob.store, C.OBDebugFlags(flags))
	if rc != 0 {
		err = createError()
	}
	return
}

// Returns a Box, panics on error (see BoxOrError)
func (ob *ObjectBox) Box(typeId TypeId) *Box {
	box, err := ob.BoxOrError(typeId)
	if err != nil {
		panic("Could not create box for type ID " + strconv.Itoa(int(typeId)) + ": " + err.Error())
	}
	return box
}

func (ob *ObjectBox) BoxOrError(typeId TypeId) (*Box, error) {
	binding := ob.getBindingById(typeId)
	cbox := C.obx_box_create(ob.store, C.obx_schema_id(typeId))
	if cbox == nil {
		return nil, createError()
	}
	return &Box{
		objectBox: ob,
		box:       cbox,
		typeId:    typeId,
		binding:   binding,
		fbb:       flatbuffers.NewBuilder(512),
	}, nil
}

func (ob *ObjectBox) AwaitAsyncCompletion() *ObjectBox {
	if C.obx_store_await_async_completion(ob.store) != 0 {
		fmt.Println(createError())
	}
	return ob
}

func (ob *ObjectBox) Query(typeId TypeId) *QueryBuilder {
	qb := C.obx_qb_create(ob.store, C.obx_schema_id(typeId))
	var err error = nil
	if qb == nil {
		err = createError()
	}
	return &QueryBuilder{
		typeId:    typeId,
		objectBox: ob,
		cqb:       qb,
		Err:       err,
	}
}

func (bytesArray *BytesArray) Free() {
	cBytesArray := bytesArray.cBytesArray
	if cBytesArray != nil {
		bytesArray.cBytesArray = nil
		C.obx_bytes_array_free(cBytesArray)
	}
	bytesArray.BytesArray = nil
}

func createError() error {
	msg := C.obx_last_error_message()
	if msg == nil {
		return errors.New("no error info available; please report")
	} else {
		return errors.New(C.GoString(msg))
	}
}
