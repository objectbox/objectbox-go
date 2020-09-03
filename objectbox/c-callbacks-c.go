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
func cVoidCallbackDispatch(callbackIdPtr unsafe.Pointer) {
	var callbackId = *(*cCallbackId)(callbackIdPtr)
	var callback = cCallbackLookup(callbackId)
	if callback != nil {
		callback.callVoid()
	}
}

//export cVoidUint64CallbackDispatch
func cVoidUint64CallbackDispatch(callbackIdPtr unsafe.Pointer, arg uint64) {
	var callbackId = *(*cCallbackId)(callbackIdPtr)
	var callback = cCallbackLookup(callbackId)
	if callback != nil {
		callback.callVoidUint64(arg)
	}
}

//export cVoidConstVoidCallbackDispatch
func cVoidConstVoidCallbackDispatch(callbackIdPtr unsafe.Pointer, arg unsafe.Pointer) {
	var callbackId = *(*cCallbackId)(callbackIdPtr)
	var callback = cCallbackLookup(callbackId)
	if callback != nil {
		callback.callVoidConstVoid(arg)
	}
}
