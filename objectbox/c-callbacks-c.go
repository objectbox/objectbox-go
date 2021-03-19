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

// This file implements externs defined in c-callbacks.go.
// It needs to be separate or it would cause duplicate symbol errors during linking.
// See https://golang.org/cmd/cgo/#hdr-C_references_to_Go for more details.

/*
#include <stdbool.h>
#include <stdint.h>
*/
import "C"
import (
	"unsafe"
)

// These functions find the callback based on the pointer to the callbackId and call it.

//export cVoidCallbackDispatch
func cVoidCallbackDispatch(callbackIdPtr C.uintptr_t) {
	var callback = cCallbackLookup(callbackIdPtr)
	if callback != nil {
		callback.callVoid()
	}
}

//export cVoidUint64CallbackDispatch
func cVoidUint64CallbackDispatch(callbackIdPtr C.uintptr_t, arg uint64) {
	var callback = cCallbackLookup(callbackIdPtr)
	if callback != nil {
		callback.callVoidUint64(arg)
	}
}

//export cVoidInt64CallbackDispatch
func cVoidInt64CallbackDispatch(callbackIdPtr C.uintptr_t, arg int64) {
	var callback = cCallbackLookup(callbackIdPtr)
	if callback != nil {
		callback.callVoidInt64(arg)
	}
}

//export cVoidConstVoidCallbackDispatch
func cVoidConstVoidCallbackDispatch(callbackIdPtr C.uintptr_t, arg unsafe.Pointer) {
	var callback = cCallbackLookup(callbackIdPtr)
	if callback != nil {
		callback.callVoidConstVoid(arg)
	}
}
