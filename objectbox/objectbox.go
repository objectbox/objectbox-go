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
	"errors"
	"fmt"
	"runtime"
	"strconv"

	"github.com/google/flatbuffers/go"
)

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

// Internal: Won't be public in the future
type BytesArray struct {
	BytesArray  [][]byte
	cBytesArray *C.OBX_bytes_array
}

type txnFun func(transaction *Transaction) error
type cursorFun func(cursor *cursor) error

// Close fully closes the database and free's resources
func (ob *ObjectBox) Close() {
	storeToClose := ob.store
	ob.store = nil
	if storeToClose != nil {
		C.obx_store_close(storeToClose)
	}
}

func (ob *ObjectBox) beginTxn() (*Transaction, error) {
	var ctxn = C.obx_txn_begin(ob.store)
	if ctxn == nil {
		return nil, createError()
	}
	return &Transaction{ctxn, ob}, nil
}

func (ob *ObjectBox) beginTxnRead() (*Transaction, error) {
	var ctxn = C.obx_txn_begin_read(ob.store)
	if ctxn == nil {
		return nil, createError()
	}
	return &Transaction{ctxn, ob}, nil
}

func (ob *ObjectBox) runInTxn(readOnly bool, txnFun txnFun) (err error) {
	runtime.LockOSThread()
	var txn *Transaction
	if readOnly {
		txn, err = ob.beginTxnRead()
	} else {
		txn, err = ob.beginTxn()
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

func (ob *ObjectBox) runWithCursor(typeId TypeId, readOnly bool, cursorFun cursorFun) error {
	binding := ob.getBindingById(typeId)
	return ob.runInTxn(readOnly, func(txn *Transaction) error {
		cursor, err := txn.createCursor(typeId, binding)
		if err != nil {
			return err
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
		return err
	})
}

// SetDebugFlags configures debug logging of the ObjectBox core
// see DebugFlags_* constants
func (ob *ObjectBox) SetDebugFlags(flags uint) error {
	rc := C.obx_store_debug_flags(ob.store, C.OBDebugFlags(flags))
	if rc != 0 {
		return createError()
	}
	return nil
}

// Box opens an Entity Box which provides CRUD access to objects
// panics on error (in case entity with the given ID doesn't exist)
func (ob *ObjectBox) Box(typeId TypeId) *Box {
	box, err := ob.BoxOrError(typeId)
	if err != nil {
		panic("Could not create box for type ID " + strconv.Itoa(int(typeId)) + ": " + err.Error())
	}
	return box
}

// BoxOrError opens an Entity Box which provides CRUD access to objects of the given type
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

// AwaitAsyncCompletion blocks until all PutAsync insert have been processed
func (ob *ObjectBox) AwaitAsyncCompletion() *ObjectBox {
	if C.obx_store_await_async_completion(ob.store) != 0 {
		fmt.Println(createError())
	}
	return ob
}

// Query starts to build a new Query
// Deprecated: this function is subject to change due to necessary typeId argument
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

func (bytesArray *BytesArray) free() {
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
