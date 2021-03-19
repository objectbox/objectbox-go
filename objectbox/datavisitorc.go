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

// This file implements externs defined in datavisitor.go.
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

//export dataVisitorDispatch
// This function finds the data visitor (based on the pointer to the visitorId) and calls it with the given data
// NOTE: don't change ptr contents, it's `const void*` in C but go doesn't support const pointers
func dataVisitorDispatch(visitorIdPtr unsafe.Pointer, data unsafe.Pointer, size C.size_t) C.bool {
	var visitorId = *(*uint32)(visitorIdPtr)

	// create an empty byte slice and map the C data to it, no copy required
	var bytes []byte
	cVoidPtrToByteSlice(data, int(size), &bytes)

	var fn = dataVisitorLookup(uint32(visitorId))
	return C.bool(fn(bytes))
}
