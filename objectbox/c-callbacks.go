/*
 * Copyright 2018-2021 ObjectBox Ltd. All rights reserved.
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
This file implements some universal formats of C callbacks forwarding to Go callbacks

Overview:
	* Register a callback, getting a callback ID.
	* Pass the registered callback ID together with a generic C callback (e.g. C.cVoidCallbackDispatch) to a C.obx_* function.
	* When ObjectBox calls C.cVoidCallbackDispatch, it finds the callback registered under that ID and calls it.
	* After there can be no more callbacks, the callback must be unregistered.

Code example:
	callbackId, err := cCallbackRegister(cVoidCallback(func() { // cVoidCallback() is a type-cast
		// do your thing here
	}))
	if err != nil {
		return err
	}

	// don't forget to unregister the callback after it's no longer going to be called or you would fill the queue up quickly
	defer cCallbackUnregister(callbackId)

	rc := C.obx_callback_taking_method(cStruct, (*C.actual_fn_type_required_by_c_method)(cVoidCallbackDispatch), unsafe.Pointer(&callbackId))
*/

/*
#include "objectbox.h"

// following functions implement forwarding and are passed to the c-api

// void return, no arguments
extern void cVoidCallbackDispatch(uintptr_t callbackId);
typedef void cVoidCallback(uintptr_t callbackId);

// void return, uint64 argument
extern void cVoidUint64CallbackDispatch(uintptr_t callbackId);
typedef void cVoidUint64Callback(uintptr_t callbackId, uint64_t arg);

// void return, int64 argument
extern void cVoidInt64CallbackDispatch(uintptr_t callbackId);
typedef void cVoidInt64Callback(uintptr_t callbackId, int64_t arg);

// void return, const uintptr_t argument
extern void cVoidConstVoidCallbackDispatch(uintptr_t callbackId);
typedef void cVoidConstVoidCallback(uintptr_t callbackId, const void* arg);
*/
import "C"
import (
	"fmt"
	"sync"
	"unsafe"
)

// cCallable allows us to avoid defining below register/unregister/lookup methods for all possible function signatures.
// This interface is implemented by callback functions and each function signature should only implement a single method
// and panic in all others. Method name format used below: "call<ReturnType><...ArgNType>()"
type cCallable interface {
	callVoid()
	callVoidUint64(uint64)
	callVoidInt64(int64)
	callVoidConstVoid(unsafe.Pointer)
}

// programming error - using an incorrect `cCallable` (arguments and return-type combination)
const cCallablePanicMsg = "invalid callback signature"

type cVoidCallback func()

func (fn cVoidCallback) callVoid()                        { fn() }
func (fn cVoidCallback) callVoidUint64(uint64)            { panic(cCallablePanicMsg) }
func (fn cVoidCallback) callVoidInt64(int64)              { panic(cCallablePanicMsg) }
func (fn cVoidCallback) callVoidConstVoid(unsafe.Pointer) { panic(cCallablePanicMsg) }

var cVoidCallbackDispatchPtr = (*C.cVoidCallback)(unsafe.Pointer(C.cVoidCallbackDispatch))

type cVoidUint64Callback func(uint64)

func (fn cVoidUint64Callback) callVoid()                        { panic(cCallablePanicMsg) }
func (fn cVoidUint64Callback) callVoidUint64(arg uint64)        { fn(arg) }
func (fn cVoidUint64Callback) callVoidInt64(int64)              { panic(cCallablePanicMsg) }
func (fn cVoidUint64Callback) callVoidConstVoid(unsafe.Pointer) { panic(cCallablePanicMsg) }

var cVoidUint64CallbackDispatchPtr = (*C.cVoidUint64Callback)(unsafe.Pointer(C.cVoidUint64CallbackDispatch))

type cVoidInt64Callback func(int64)

func (fn cVoidInt64Callback) callVoid()                        { panic(cCallablePanicMsg) }
func (fn cVoidInt64Callback) callVoidUint64(uint64)            { panic(cCallablePanicMsg) }
func (fn cVoidInt64Callback) callVoidInt64(arg int64)          { fn(arg) }
func (fn cVoidInt64Callback) callVoidConstVoid(unsafe.Pointer) { panic(cCallablePanicMsg) }

var cVoidInt64CallbackDispatchPtr = (*C.cVoidInt64Callback)(unsafe.Pointer(C.cVoidInt64CallbackDispatch))

type cVoidConstVoidCallback func(unsafe.Pointer)

func (fn cVoidConstVoidCallback) callVoid()                            { panic(cCallablePanicMsg) }
func (fn cVoidConstVoidCallback) callVoidUint64(uint64)                { panic(cCallablePanicMsg) }
func (fn cVoidConstVoidCallback) callVoidInt64(int64)                  { panic(cCallablePanicMsg) }
func (fn cVoidConstVoidCallback) callVoidConstVoid(arg unsafe.Pointer) { fn(arg) }

var cVoidConstVoidCallbackDispatchPtr = (*C.cVoidConstVoidCallback)(unsafe.Pointer(C.cVoidConstVoidCallbackDispatch))

type cCallbackId uint32

var cCallbackLastId cCallbackId
var cCallbackMutex sync.Mutex
var cCallbackMap = make(map[cCallbackId]cCallable)

// The result is actually not a memory pointer, just a number. That's also how it's used in cCallbackLookup().
func (cbId cCallbackId) cPtr() unsafe.Pointer {
	return unsafe.Pointer(uintptr(cbId))
}

// Returns the next cCallbackId in a sequence (NOT checking its availability), skipping zero.
func cCallbackNextId() cCallbackId {
	cCallbackLastId++
	if cCallbackLastId == 0 {
		cCallbackLastId++
	}
	return cCallbackLastId
}

func cCallbackRegister(fn cCallable) (cCallbackId, error) {
	cCallbackMutex.Lock()
	defer cCallbackMutex.Unlock()

	// cycle through ids until we find an empty slot
	var initialId = cCallbackNextId()
	for cCallbackMap[cCallbackLastId] != nil {
		cCallbackNextId()

		if initialId == cCallbackLastId {
			return 0, fmt.Errorf("full queue of data-callback callbacks - can't allocate another")
		}
	}

	cCallbackMap[cCallbackLastId] = fn
	return cCallbackLastId, nil
}

func cCallbackLookup(id C.uintptr_t) cCallable {
	cCallbackMutex.Lock()
	defer cCallbackMutex.Unlock()

	fn, found := cCallbackMap[cCallbackId(id)]
	if !found {
		// this might happen in extraordinary circumstances, e.g. during shutdown if there are still some sync listeners
		fmt.Println(fmt.Errorf("invalid C-API callback ID %d", id))
		return nil
	}

	return fn
}

func cCallbackUnregister(id cCallbackId) {
	// special value - not registered
	if id == 0 {
		return
	}

	cCallbackMutex.Lock()
	defer cCallbackMutex.Unlock()

	delete(cCallbackMap, id)
}
