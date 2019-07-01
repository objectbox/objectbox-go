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
This file implements obx_data_visitor forwarding to Go callbacks

Overview:
	* Register a dataVisitor callback, getting a visitor ID.
	* Pass the registered visitor ID together with a generic dataVisitor (C.dataVisitorDispatch) to a C.obx_* function.
	* When ObjectBox calls dataVisitorDispatch, it finds the callback registered under that ID and calls it.
	* After there can be no more callbacks, the visitor must be unregistered.

Code example:
	var visitorId uint32
	visitorId, err = dataVisitorRegister(func(bytes []byte) bool {
		// do your thing with the data
		object := Cursor.binding.Load(bytes)
		return true // this return value is passed back to the ObjectBox, usually used to break the traversal
	})

	if err != nil {
		return err
	}

	// don't forget to unregister the visitor after it's no longer going to be called or you would fill the queue up quickly
	defer dataVisitorUnregister(visitorId)

	rc := C.obx_query_visit(cQuery, cCursor, dataVisitor, unsafe.Pointer(&visitorId), C.uint64_t(offset), C.uint64_t(limit))
*/

/*
#include "objectbox.h"

// this implements the obx_data_visitor forwarding, it's called from ObjectBox C-api (see `dataVisitor` go var)
extern bool dataVisitorDispatch(void* visitorId, void* data, size_t size);
*/
import "C"
import (
	"fmt"
	"sync"
	"unsafe"
)

type dataVisitorCallback = func([]byte) bool

var dataVisitor = (*C.obx_data_visitor)(unsafe.Pointer(C.dataVisitorDispatch))
var dataVisitorId uint32
var dataVisitorMutex sync.Mutex
var dataVisitorCallbacks = make(map[uint32]dataVisitorCallback)

func dataVisitorRegister(fn dataVisitorCallback) (uint32, error) {
	dataVisitorMutex.Lock()
	defer dataVisitorMutex.Unlock()

	// cycle through ids until we find an empty slot
	dataVisitorId++
	var initialId = dataVisitorId
	for dataVisitorCallbacks[dataVisitorId] != nil {
		dataVisitorId++

		if initialId == dataVisitorId {
			return 0, fmt.Errorf("full queue of data-visitor callbacks - can't allocate another")
		}
	}

	dataVisitorCallbacks[dataVisitorId] = fn
	return dataVisitorId, nil
}

func dataVisitorLookup(id uint32) dataVisitorCallback {
	dataVisitorMutex.Lock()
	defer dataVisitorMutex.Unlock()

	return dataVisitorCallbacks[id]
}

func dataVisitorUnregister(id uint32) {
	dataVisitorMutex.Lock()
	defer dataVisitorMutex.Unlock()

	delete(dataVisitorCallbacks, id)
}

