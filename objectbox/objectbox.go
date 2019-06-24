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
*/
import "C"

import (
	"fmt"
	"runtime"
	"strconv"
	"sync"
)

const (
	DebugflagsLogTransactionsRead  = 1
	DebugflagsLogTransactionsWrite = 2
	DebugflagsLogQueries           = 4
	DebugflagsLogQueryParameters   = 8
	DebugflagsLogAsyncQueue        = 16
)

const (
	// Standard put ("insert or update")
	cPutModePut = 1

	// Put succeeds only if the entity does not exist yet.
	cPutModeInsert = 2

	// Put succeeds only if the entity already exist.
	cPutModeUpdate = 3

	// Not used yet (does not make sense for asnyc puts)
	// The given ID (non-zero) is guaranteed to be new; don't use unless you know exactly what you are doing!
	// This is primarily used internally. Wrong usage leads to inconsistent data (e.g. index data not updated)!
	cPutModePutIdGuaranteedToBeNew = 4
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
	boxesMutex     sync.Mutex
	options        options
}

type options struct {
	putAsyncTimeout uint
}

// constant during runtime so no need to call this each time it's necessary
var supportsBytesArray = bool(C.obx_supports_bytes_array())

// Close fully closes the database and free's resources
func (ob *ObjectBox) Close() {
	storeToClose := ob.store
	ob.store = nil
	if storeToClose != nil {
		C.obx_store_close(storeToClose)
	}
}

// RunInReadTx executes the given function inside a read transaction.
// Note, you must not launch Go routines inside this function - the call must be sequential.
// The error returned by your callback is passed-through as the output error
func (ob *ObjectBox) RunInReadTx(fn func() error) error {
	return ob.runInTxn(true, fn)
}

// RunInWriteTx executes the given function inside a write transaction.
// Note, you must not launch Go routines inside this function - the call must be sequential.
// The error returned by your callback is passed-through as the output error.
// If the resulting error is not nil, the transaction is aborted (rolled-back)
func (ob *ObjectBox) RunInWriteTx(fn func() error) error {
	return ob.runInTxn(false, fn)
}

func (ob *ObjectBox) beginTxn() (*transaction, error) {
	var ctxn = C.obx_txn_begin(ob.store)
	if ctxn == nil {
		return nil, createError()
	}
	return &transaction{ctxn, ob}, nil
}

func (ob *ObjectBox) beginTxnRead() (*transaction, error) {
	var ctxn = C.obx_txn_begin_read(ob.store)
	if ctxn == nil {
		return nil, createError()
	}
	return &transaction{ctxn, ob}, nil
}

func (ob *ObjectBox) runInTxn(readOnly bool, fn func() error) (err error) {
	runtime.LockOSThread()
	var txn *transaction
	if readOnly {
		txn, err = ob.beginTxnRead()
	} else {
		txn, err = ob.beginTxn()
	}

	if err != nil {
		runtime.UnlockOSThread()
		return err
	}

	// Defer to ensure a TX is ALWAYS closed, even in a panic
	defer func() {
		err2 := txn.Close()
		if err == nil {
			err = err2
		}
		runtime.UnlockOSThread()
	}()

	err = fn()

	if !readOnly && err == nil {
		err = txn.Commit()
	}

	return err
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

// SetDebugFlags configures debug logging of the ObjectBox core.
// See DebugFlags* constants
func (ob *ObjectBox) SetDebugFlags(flags uint) error {
	rc := C.obx_store_debug_flags(ob.store, C.OBXDebugFlags(flags))
	if rc != 0 {
		return createError()
	}
	return nil
}

// InternalBox returns an Entity Box or panics on error (in case entity with the given ID doesn't exist)
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
	cbox := C.obx_box(ob.store, C.obx_schema_id(typeId))
	if cbox == nil {
		return nil, createError()
	}

	box = &Box{
		objectBox: ob,
		cBox:      cbox,
		entity:    entity,
	}
	ob.boxes[typeId] = box
	return box, nil
}

// AwaitAsyncCompletion blocks until all PutAsync insert have been processed
func (ob *ObjectBox) AwaitAsyncCompletion() *ObjectBox {
	if C.obx_store_await_async_completion(ob.store) != true {
		fmt.Println(createError())
	}
	return ob
}
