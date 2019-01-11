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
#include "datavisitor.h"

extern bool dataVisitorDispatch(uint32_t id, void* data, size_t size);
*/
import "C"
import (
	"fmt"
	"reflect"
	"sync"
	"unsafe"
)

// this file implements obx_data_visitor forwarding to Go callbacks
// see https://github.com/golang/go/wiki/cgo#function-variables for introduction

type dataVisitorCallback = func([]byte) bool

var dataVisitorIndex uint32
var dataVisitorMutex sync.Mutex
var dataVisitorCallbacks = make(map[uint32]dataVisitorCallback)

func dataVisitorRegister(fn dataVisitorCallback) (uint32, error) {
	dataVisitorMutex.Lock()
	defer dataVisitorMutex.Unlock()

	// cycle through indexes until we find an empty slot
	dataVisitorIndex++
	var initialIndex = dataVisitorIndex
	for dataVisitorCallbacks[dataVisitorIndex] != nil {
		dataVisitorIndex++

		if initialIndex == dataVisitorIndex {
			return 0, fmt.Errorf("full queue of data-visitor callbacks - can't allocate another")
		}
	}

	dataVisitorCallbacks[dataVisitorIndex] = fn
	return dataVisitorIndex, nil
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

//export dataVisitorDispatch
// NOTE: don't change ptr, it's `const void*` in C
func dataVisitorDispatch(id C.uint, ptr unsafe.Pointer, size C.size_t) C.bool {
	// create an empty byte slice and map the C data to it, no copy required
	var bytes []byte
	header := (*reflect.SliceHeader)(unsafe.Pointer(&bytes))
	*header = reflect.SliceHeader{Data: uintptr(ptr), Len: int(size), Cap: int(size)}

	var fn = dataVisitorLookup(uint32(id))
	return C.bool(fn(bytes))
}
