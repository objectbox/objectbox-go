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
This file implements obx_txn_callable_* forwarding to Go callbacks

Overview:
	* register a txnCallable callback
	* pass the registered callback Id together with a generic C.txn_callable_read|write to C.obx_store_exec_read|write()
	* when ObjectBox calls the C.txn_callable_read|write, the call is forwarded to Go txnCallableDispatch
	* txnCallableDispatch finds the callback registered under that ID and calls it
	* after there can be no more callbacks, the visitor must be unregistered

Code example:
	TODO update
*/

/*
#include "objectbox.h"

// this is a Go function defined bellow and called from C
extern bool txnCallableDispatch(uint32_t id, OBX_txn* txn);
*/
import "C"
import (
	"fmt"
	"sync"
	"unsafe"
)

type txnCallable = func(tx *Transaction) bool

var txnCallableRead = (*C.obx_txn_callable_read)(unsafe.Pointer(C.txnCallableDispatch))
var txnCallableWrite = (*C.obx_txn_callable_write)(unsafe.Pointer(C.txnCallableDispatch))
var txnCallableId uint32
var txnCallableMutex sync.Mutex
var txnCallables = make(map[uint32]txnCallable)

func txnCallableRegister(fn txnCallable) (uint32, error) {
	txnCallableMutex.Lock()
	defer txnCallableMutex.Unlock()

	// cycle through ids until we find an empty slot
	txnCallableId++
	var initialId = txnCallableId
	for txnCallables[txnCallableId] != nil {
		txnCallableId++

		if initialId == txnCallableId {
			return 0, fmt.Errorf("full queue of txn-callables - can't allocate another")
		}
	}

	txnCallables[txnCallableId] = fn
	return txnCallableId, nil
}

func txnCallableLookup(id uint32) txnCallable {
	txnCallableMutex.Lock()
	defer txnCallableMutex.Unlock()

	return txnCallables[id]
}

func txnCallableUnregister(id uint32) {
	txnCallableMutex.Lock()
	defer txnCallableMutex.Unlock()

	delete(txnCallables, id)
}
