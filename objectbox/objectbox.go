/*
 * Copyright 2019 ObjectBox Ltd. All rights reserved.
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
#include "txncallable.h"
*/
import "C"

import (
	"fmt"
	"runtime"
	"strconv"
	"sync"
	"unsafe"

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

// atomic boolean true & false
const aTrue = 1
const aFalse = 0

type TypeId uint32

type ObjectBox struct {
	store          *C.OBX_store
	entitiesById   map[TypeId]*entity
	entitiesByName map[string]*entity
	boxes          map[TypeId]*Box
	boxesMutex     *sync.Mutex
	options        options
}

type options struct {
	putAsyncTimeout  uint
}

type txnFun func(transaction *Transaction) error

// constant during runtime so no need to call this each time it's necessary
var supportsBytesArray = bool(C.obx_supports_bytes_array())

// Close fully closes the database and free's resources
func (ob *ObjectBox) Close() {
	storeToClose := ob.store
	ob.store = nil
	if storeToClose != nil {
		C.obx_store_close(storeToClose)
	}

	ob.boxesMutex.Lock()
	defer ob.boxesMutex.Unlock()
	for _, box := range ob.boxes {
		if err := box.close(); err != nil {
			fmt.Println(err)
		}
	}
	ob.boxes = nil
}

// View executes the given function inside a read transaction.
// Note, you must not launch Go routines inside this function - the call must be sequential.
// The error returned by your callback is passed-through as the output error
func (ob *ObjectBox) View(fn func() error) error {
	return ob.runInTxn(true, nil, fn)
}

// Update executes the given function inside a write transaction.
// Note, you must not launch Go routines inside this function - the call must be sequential.
// The error returned by your callback is passed-through as the output error.
// If the resulting error is not nil, the transaction is aborted (rolled-back)
func (ob *ObjectBox) Update(fn func() error) error {
	return ob.runInTxn(false, nil, fn)
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

// runInTxn executes one of the given function inside a transaction
// if the function returns an error, the transaction is rolled-back
func (ob *ObjectBox) runInTxn(readOnly bool, fnWithTxn txnFun, fnNoTxn func() error) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	var visitorId uint32
	var err error
	var pnc interface{}

	visitorId, err = txnCallableRegister(func(tx *Transaction) (result bool) {
		// this function must not panic or the call to C would not return and the transaction wouldn't finish
		defer func() {
			pnc = recover()
			if pnc != nil {
				fmt.Printf("panic during transaction callback: %v", pnc)
				err = fmt.Errorf("%v", pnc)
				result = false
			}
		}()

		if fnWithTxn != nil {
			if tx.objectBox == nil {
				tx.objectBox = ob
			}
			err = fnWithTxn(tx)
		} else {
			err = fnNoTxn()
		}
		return err == nil
	})
	if err != nil {
		return err
	}
	defer txnCallableUnregister(visitorId)

	var rc C.obx_err
	if readOnly {
		rc = C.obx_store_exec_read(ob.store, C.txn_callable_read, unsafe.Pointer(&visitorId))
	} else {
		rc = C.obx_store_exec_write(ob.store, C.txn_callable_write, unsafe.Pointer(&visitorId))
	}

	// handle errors in order of their priorities
	if pnc != nil {
		// propagate the panic
		panic(pnc)
	} else if err != nil {
		// err set by the visitor callback above
		return err
	} else if rc != 0 {
		return createError()
	} else {
		return nil
	}
}

func (ob ObjectBox) getEntityById(typeId TypeId) *entity {
	entity := ob.entitiesById[typeId]
	if entity == nil {
		// Configuration error by the dev, OK to panic
		panic("Configuration error; no entity registered for type ID " + strconv.Itoa(int(typeId)))
	}
	return entity
}

func (ob ObjectBox) getEntityByName(typeName string) *entity {
	entity := ob.entitiesByName[typeName]
	if entity == nil {
		// Configuration error by the dev, OK to panic
		panic("Configuration error; no entity registered for type name " + typeName)
	}
	return entity
}

func (ob *ObjectBox) runWithCursor(e *entity, readOnly bool, cursorFun cursorFun) error {
	if ob.options.alwaysAwaitAsync {
		e.awaitAsyncCompletion()
	}

	return ob.runInTxn(readOnly, func(txn *Transaction) error {
		return txn.runWithCursor(e, cursorFun)
	}, nil)
}

// SetDebugFlags configures debug logging of the ObjectBox core.
// See DebugFlags_* constants
func (ob *ObjectBox) SetDebugFlags(flags uint) error {
	rc := C.obx_store_debug_flags(ob.store, C.OBXDebugFlags(flags))
	if rc != 0 {
		return createError()
	}
	return nil
}

// panics on error (in case entity with the given ID doesn't exist)
func (ob *ObjectBox) InternalBox(typeId TypeId) *Box {
	box, err := ob.box(typeId)
	if err != nil {
		panic(fmt.Sprintf("Could not create box for type ID %d: %s", typeId, err))
	}
	return box
}

// Gets an Entity Box which provides CRUD access to objects of the given type
func (ob *ObjectBox) box(typeId TypeId) (*Box, error) {
	ob.boxesMutex.Lock()
	defer ob.boxesMutex.Unlock()

	box := ob.boxes[typeId]
	if box != nil {
		return box, nil
	}

	entity := ob.getEntityById(typeId)
	cbox := C.obx_box_create(ob.store, C.obx_schema_id(typeId))
	if cbox == nil {
		return nil, createError()
	}

	box = &Box{
		objectBox: ob,
		box:       cbox,
		entity:    entity,
	}
	ob.boxes[typeId] = box
	return box, nil
}

// AwaitAsyncCompletion blocks until all PutAsync insert have been processed
func (ob *ObjectBox) AwaitAsyncCompletion() *ObjectBox {
	if C.obx_store_await_async_completion(ob.store) != 0 {
		fmt.Println(createError())
	}
	return ob
}
