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

char** newCharArray(int size) {
	return calloc(sizeof(char*), size);
}

void setArrayString(const char **array, size_t index, const char *value) {
    array[index] = value;
}

void freeCharArray(char **a, int size) {
    for (size_t i = 0; i < size; i++) {
    	free(a[i]);
    }
    free(a);
}
*/
import "C"
import (
	"reflect"
	"unsafe"
)

type bytesArray struct {
	BytesArray  [][]byte
	cBytesArray *C.OBX_bytes_array
}

func (bytesArray *bytesArray) free() {
	cBytesArray := bytesArray.cBytesArray
	if cBytesArray != nil {
		bytesArray.cBytesArray = nil
		C.obx_bytes_array_free(cBytesArray)
	}
	bytesArray.BytesArray = nil
}

func cBytesArrayToGo(cBytesArray *C.OBX_bytes_array) *bytesArray {
	size := int(cBytesArray.count)
	plainBytesArray := make([][]byte, size)

	if size > 0 {
		var sliceOfCBytes []C.OBX_bytes
		// see cVoidPtrToByteSlice for documentation of the following approach in general
		header := (*reflect.SliceHeader)(unsafe.Pointer(&sliceOfCBytes))
		header.Data = uintptr(unsafe.Pointer(cBytesArray.bytes))
		header.Len = size
		header.Cap = size

		for i := 0; i < size; i++ {
			cBytes := sliceOfCBytes[i]
			cVoidPtrToByteSlice(unsafe.Pointer(cBytes.data), int(cBytes.size), &(plainBytesArray[i]))
		}
	}

	return &bytesArray{plainBytesArray, cBytesArray}
}

type idsArray struct {
	ids    []uint64
	cArray *C.OBX_id_array
}

func (array *idsArray) free() {
	if array.cArray != nil {
		C.obx_id_array_free(array.cArray)
		array.cArray = nil
	}
	array.ids = nil
}

func cIdsArrayToGo(cArray *C.OBX_id_array) *idsArray {
	var size = uint(cArray.count)
	var ids = make([]uint64, size)
	if size > 0 {
		var cArrayStart = uintptr(unsafe.Pointer(cArray.ids))
		var cIdSize = unsafe.Sizeof(*cArray.ids)
		for i := uint(0); i < size; i++ {
			ids[i] = *(*uint64)(unsafe.Pointer(cArrayStart + uintptr(i)*cIdSize))
		}
	}
	return &idsArray{ids, cArray}
}

type stringArray struct {
	cArray **C.char
	size   int
}

func (array *stringArray) free() {
	if array.cArray != nil {
		C.freeCharArray(array.cArray, C.int(array.size))
		array.cArray = nil
	}
}

func goStringArrayToC(values []string) *stringArray {
	result := &stringArray{
		cArray: C.newCharArray(C.int(len(values))),
		size:   len(values),
	}
	for i, s := range values {
		C.setArrayString(result.cArray, C.size_t(i), C.CString(s))
	}
	return result
}

func goInt64ArrayToC(values []int64) *C.int64_t {
	if len(values) > 0 {
		return (*C.int64_t)(unsafe.Pointer(&values[0]))
	} else {
		return nil
	}
}

func goInt32ArrayToC(values []int32) *C.int32_t {
	if len(values) > 0 {
		return (*C.int32_t)(unsafe.Pointer(&values[0]))
	} else {
		return nil
	}
}

func cBytesPtr(value []byte) unsafe.Pointer {
	if len(value) >= 1 {
		return unsafe.Pointer(&value[0])
	} else {
		return nil
	}
}

// Maps a C void* to the given byte-slice. The void* is not garbage collected and must be managed outside.
//
// Previous alternative without reflect https://github.com/golang/go/wiki/cgo#turning-c-arrays-into-go-slices
// was broken on 32-bit platforms, see https://github.com/golang/go/issues/13656#issuecomment-303246650
// thus we have chosen a solution mapping the C pointers to a Go slice.
// Performance-wise, there's no noticeable difference and the current solution is more "obvious"
//
// NOTE watch https://github.com/golang/go/issues/19367 for possible changes & new solutions.
// As both unsafe package as well as reflect.SliceHeader might change in the future,
// the above-mentioned issue might describe an alternative solution
func cVoidPtrToByteSlice(data unsafe.Pointer, size int, bytes *[]byte) {
	header := (*reflect.SliceHeader)(unsafe.Pointer(bytes))
	header.Data = uintptr(data)
	header.Len = size
	header.Cap = size
}
