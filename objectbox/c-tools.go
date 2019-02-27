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
	"errors"
	"sync"
)

// provides wrappers for objectbox C-api calls, making sure the returned error belongs to this call
// TODO migrate all native C.obx_* calls so that they use these wrappers

// TODO instead of this mutex, we could make c-api block before an error is overwritten (so it has to be "acknowledged" first)
var cCallMutex sync.Mutex

func cMaybeErr(fn func() C.obx_err) error {
	cCallMutex.Lock()
	defer cCallMutex.Unlock()

	if rc := fn(); rc != 0 {
		return createError()
	} else {
		return nil
	}
}

func cGetIds(fn func() *C.OBX_id_array) ([]uint64, error) {
	cCallMutex.Lock()
	defer cCallMutex.Unlock()

	var cArray = fn()
	if cArray == nil {
		return nil, createError()
	}

	ids := cIdsArrayToGo(cArray)
	defer ids.free()

	return ids.ids, nil
}

func cGetBytesArray(fn func() *C.OBX_bytes_array) ([][]byte, error) {
	cCallMutex.Lock()
	defer cCallMutex.Unlock()

	var cArray = fn()
	if cArray == nil {
		return nil, createError()
	}

	bytes := cBytesArrayToGo(cArray)
	defer bytes.free()

	return bytes.array, nil
}

func createError() error {
	msg := C.obx_last_error_message()
	if msg == nil {
		return errors.New("no error info available; please report")
	} else {
		return errors.New(C.GoString(msg))
	}
}

