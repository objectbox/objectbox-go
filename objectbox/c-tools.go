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
#include <stdlib.h>
#include "objectbox.h"
*/
import "C"
import (
	"errors"
	"runtime"
)

// provides wrappers for objectbox C-api calls, making sure the returned error belongs to this call.

func cCall(fn func() C.obx_err) (err error) {
	runtime.LockOSThread()

	if rc := fn(); rc != 0 {
		err = createError()
	}

	runtime.UnlockOSThread()
	return err
}

func cCallBool(fn func() bool) (err error) {
	runtime.LockOSThread()

	if successful := fn(); !successful {
		err = createError()
	}

	runtime.UnlockOSThread()
	return err
}

func cGetIds(fn func() *C.OBX_id_array) (ids []uint64, err error) {
	runtime.LockOSThread()

	var cArray = fn()
	if cArray == nil {
		err = createError()
	} else {
		ids = cIdsArrayToGo(cArray)
		C.obx_id_array_free(cArray)
	}

	runtime.UnlockOSThread()
	return ids, err
}

func cGetBytesArray(fn func() *C.OBX_bytes_array) (array [][]byte, err error) {
	runtime.LockOSThread()

	var cArray = fn()
	if cArray == nil {
		err = createError()
	} else {
		array = cBytesArrayToGo(cArray)
		C.obx_bytes_array_free(cArray)
	}

	runtime.UnlockOSThread()
	return array, err
}

// createError fetches the latest error that happened in the c-api on a current-thread.
// The c-api uses thread-local storage for the latest error so we need to lock the current goroutine to a thread.
// Must only be called when runtime.LockOSThread() is active. Either use one of the above cCall-style functions or a TX.
func createError() error {
	msg := C.obx_last_error_message()
	if msg == nil {
		return errors.New("no error info available; please report")
	}
	return errors.New(C.GoString(msg))
}
