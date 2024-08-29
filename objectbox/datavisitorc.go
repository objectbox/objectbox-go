/*
 * Copyright 2018-2024 ObjectBox Ltd. All rights reserved.
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

// This file implements externs defined in datavisitor.go.
// It needs to be separate or it would cause duplicate symbol errors during linking.
// See https://golang.org/cmd/cgo/#hdr-C_references_to_Go for more details.

/*
#include <stdbool.h>
#include <stdint.h>
*/
import "C"
import (
	"fmt"
	"unsafe"
)

// "Implements" the C function obx_data_visitor and dispatches to the registered Go data visitor.
// This function finds the data visitor (based on the pointer to the visitorId) and calls it with the given data
// NOTE: don't change ptr contents, it's `const void*` in C but go doesn't support const pointers
//
//export dataVisitorDispatch
func dataVisitorDispatch(data unsafe.Pointer, size C.size_t, userData unsafe.Pointer) C.bool {
	if userData == nil {
		panic("Internal error: visitor ID pointer is nil")
	}
	var visitorId = *(*uint32)(userData)

	// create an empty byte slice and map the C data to it, no copy required
	var bytes []byte
	if data != nil {
		cVoidPtrToByteSlice(data, int(size), &bytes)
	}

	var fn = dataVisitorLookup(visitorId)
	if fn == nil {
		panic(fmt.Sprintf("Internal error: no data visitor found for ID %d", visitorId))
	}

	return C.bool(fn(bytes))
}
