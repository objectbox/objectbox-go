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

// cGetIds converts the given C array to Go and frees the source
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

// cGetBytesArray converts the given C array to Go and frees the source
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

// cGetStrings converts the given C array to Go and frees the source
func cGetStrings(fn func() *C.OBX_string_array) (items []string, err error) {
	runtime.LockOSThread()

	var cArray = fn()
	if cArray == nil {
		err = createError()
	} else {
		items = cStringsArrayToGo(cArray)
		C.obx_string_array_free(cArray)
	}

	runtime.UnlockOSThread()
	return items, err
}

// cGetInts converts the given C array to Go and frees the source
func cGetInts(fn func() *C.OBX_int64_array) (items []int, err error) {
	runtime.LockOSThread()

	var cArray = fn()
	if cArray == nil {
		err = createError()
	} else {
		items = cIntsArrayToGo(cArray)
		C.obx_int64_array_free(cArray)
	}

	runtime.UnlockOSThread()
	return items, err
}

// cGetUint64s converts the given C array to Go and frees the source
func cGetUints(fn func() *C.OBX_int64_array) (items []uint, err error) {
	runtime.LockOSThread()

	var cArray = fn()
	if cArray == nil {
		err = createError()
	} else {
		items = cUintsArrayToGo(cArray)
		C.obx_int64_array_free(cArray)
	}

	runtime.UnlockOSThread()
	return items, err
}

// cGetInt64s converts the given C array to Go and frees the source
func cGetInt64s(fn func() *C.OBX_int64_array) (items []int64, err error) {
	runtime.LockOSThread()

	var cArray = fn()
	if cArray == nil {
		err = createError()
	} else {
		items = cInt64sArrayToGo(cArray)
		C.obx_int64_array_free(cArray)
	}

	runtime.UnlockOSThread()
	return items, err
}

// cGetUint64s converts the given C array to Go and frees the source
func cGetUint64s(fn func() *C.OBX_int64_array) (items []uint64, err error) {
	runtime.LockOSThread()

	var cArray = fn()
	if cArray == nil {
		err = createError()
	} else {
		items = cUint64sArrayToGo(cArray)
		C.obx_int64_array_free(cArray)
	}

	runtime.UnlockOSThread()
	return items, err
}

// cGetInt32s converts the given C array to Go and frees the source
func cGetInt32s(fn func() *C.OBX_int32_array) (items []int32, err error) {
	runtime.LockOSThread()

	var cArray = fn()
	if cArray == nil {
		err = createError()
	} else {
		items = cInt32sArrayToGo(cArray)
		C.obx_int32_array_free(cArray)
	}

	runtime.UnlockOSThread()
	return items, err
}

// cGetUint32s converts the given C array to Go and frees the source
func cGetUint32s(fn func() *C.OBX_int32_array) (items []uint32, err error) {
	runtime.LockOSThread()

	var cArray = fn()
	if cArray == nil {
		err = createError()
	} else {
		items = cUint32sArrayToGo(cArray)
		C.obx_int32_array_free(cArray)
	}

	runtime.UnlockOSThread()
	return items, err
}

// cGetInt16s converts the given C array to Go and frees the source
func cGetInt16s(fn func() *C.OBX_int16_array) (items []int16, err error) {
	runtime.LockOSThread()

	var cArray = fn()
	if cArray == nil {
		err = createError()
	} else {
		items = cInt16sArrayToGo(cArray)
		C.obx_int16_array_free(cArray)
	}

	runtime.UnlockOSThread()
	return items, err
}

// cGetUint16s converts the given C array to Go and frees the source
func cGetUint16s(fn func() *C.OBX_int16_array) (items []uint16, err error) {
	runtime.LockOSThread()

	var cArray = fn()
	if cArray == nil {
		err = createError()
	} else {
		items = cUint16sArrayToGo(cArray)
		C.obx_int16_array_free(cArray)
	}

	runtime.UnlockOSThread()
	return items, err
}

// cGetInt8s converts the given C array to Go and frees the source
func cGetInt8s(fn func() *C.OBX_int8_array) (items []int8, err error) {
	runtime.LockOSThread()

	var cArray = fn()
	if cArray == nil {
		err = createError()
	} else {
		items = cInt8sArrayToGo(cArray)
		C.obx_int8_array_free(cArray)
	}

	runtime.UnlockOSThread()
	return items, err
}

// cGetUint8s converts the given C array to Go and frees the source
func cGetUint8s(fn func() *C.OBX_int8_array) (items []uint8, err error) {
	runtime.LockOSThread()

	var cArray = fn()
	if cArray == nil {
		err = createError()
	} else {
		items = cUint8sArrayToGo(cArray)
		C.obx_int8_array_free(cArray)
	}

	runtime.UnlockOSThread()
	return items, err
}

// cGetBools converts the given C array to Go and frees the source
func cGetBools(fn func() *C.OBX_int8_array) (items []bool, err error) {
	runtime.LockOSThread()

	var cArray = fn()
	if cArray == nil {
		err = createError()
	} else {
		items = cBoolsArrayToGo(cArray)
		C.obx_int8_array_free(cArray)
	}

	runtime.UnlockOSThread()
	return items, err
}

// cGetFloat64s converts the given C array to Go and frees the source
func cGetFloat64s(fn func() *C.OBX_double_array) (items []float64, err error) {
	runtime.LockOSThread()

	var cArray = fn()
	if cArray == nil {
		err = createError()
	} else {
		items = cFloat64sArrayToGo(cArray)
		C.obx_double_array_free(cArray)
	}

	runtime.UnlockOSThread()
	return items, err
}

// cGetFloat32s converts the given C array to Go and frees the source
func cGetFloat32s(fn func() *C.OBX_float_array) (items []float32, err error) {
	runtime.LockOSThread()

	var cArray = fn()
	if cArray == nil {
		err = createError()
	} else {
		items = cFloat32sArrayToGo(cArray)
		C.obx_float_array_free(cArray)
	}

	runtime.UnlockOSThread()
	return items, err
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
