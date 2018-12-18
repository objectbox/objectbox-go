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
		// Previous alternative without reflect:
		//   https://github.com/golang/go/wiki/cgo#turning-c-arrays-into-go-slices (2012)
		//   On a RPi 3, the size with 1<<30 did not work, but 1<<27 did
		// Raul measured both variants and did notice a visible perf impact (Go 1.11.2)
		var goBytesArray []C.OBX_bytes
		header := (*reflect.SliceHeader)(unsafe.Pointer(&goBytesArray))
		*header = reflect.SliceHeader{Data: uintptr(unsafe.Pointer(cBytesArray.bytes)), Len: size, Cap: size}
		for i := 0; i < size; i++ {
			cBytes := goBytesArray[i]
			dataBytes := C.GoBytes(cBytes.data, C.int(cBytes.size))
			plainBytesArray[i] = dataBytes
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
