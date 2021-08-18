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
#include "objectbox-sync.h"

char** newCharArray(int size) {
	return malloc(sizeof(char*) * size);
}

void setArrayString(const char** array, size_t index, const char* value) {
    array[index] = value;
}

void freeCharArray(char** a, int size) {
	// old compiler errors out if part of the for loop: 'for' loop initial declarations are only allowed in C99 mode
	size_t i;
    for (i = 0; i < size; i++) {
    	free(a[i]);
    }
    free(a);
}
*/
import "C"
import (
	"reflect"
	"runtime"
	"unsafe"
)

type bytesArray struct {
	array       [][]byte
	cBytesArray *C.OBX_bytes_array
}

func (bytesArray *bytesArray) free() {
	cBytesArray := bytesArray.cBytesArray
	if cBytesArray != nil {
		bytesArray.cBytesArray = nil
		C.obx_bytes_array_free(cBytesArray)
	}
	bytesArray.array = nil
}

func goBytesArrayToC(goArray [][]byte) (*bytesArray, error) {
	// for native calls/createError()
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	var cArray = C.obx_bytes_array(C.size_t(len(goArray)))
	if cArray == nil {
		return nil, createError()
	}

	for i, bytes := range goArray {
		rc := C.obx_bytes_array_set(cArray, C.size_t(i), cBytesPtr(bytes), C.size_t(len(bytes)))
		if rc != 0 {
			C.obx_bytes_array_free(cArray)
			return nil, createError()
		}
	}

	return &bytesArray{goArray, cArray}, nil
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

func goIdsArrayToC(ids []uint64) (*idsArray, error) {
	// for native calls/createError()
	runtime.LockOSThread()

	var err error
	var cArray = C.obx_id_array(goUint64ArrayToCObxId(ids), C.size_t(len(ids)))
	if cArray == nil {
		err = createError()
	}

	runtime.UnlockOSThread()
	return &idsArray{ids, cArray}, err
}

type charPtrsArray struct {
	cArray **C.char
	size   int
}

func (array *charPtrsArray) free() {
	if array.cArray != nil {
		C.freeCharArray(array.cArray, C.int(array.size))
		array.cArray = nil
	}
}

func goStringArrayToC(values []string) *charPtrsArray {
	result := &charPtrsArray{
		cArray: C.newCharArray(C.int(len(values))),
		size:   len(values),
	}
	for i, s := range values {
		C.setArrayString(result.cArray, C.size_t(i), C.CString(s))
	}
	return result
}

func goInt64ArrayToC(values []int64) *C.int64_t {
	if len(values) == 0 {
		return nil
	}
	return (*C.int64_t)(unsafe.Pointer(&values[0]))
}

func goInt32ArrayToC(values []int32) *C.int32_t {
	if len(values) == 0 {
		return nil
	}
	return (*C.int32_t)(unsafe.Pointer(&values[0]))
}

func goUint64ArrayToCObxId(values []uint64) *C.obx_id {
	if len(values) == 0 {
		return nil
	}
	return (*C.obx_id)(unsafe.Pointer(&values[0]))
}

func cBytesPtr(value []byte) unsafe.Pointer {
	if len(value) == 0 {
		return nil
	}
	return unsafe.Pointer(&value[0])
}

/**
Functions for reading OBX_*_array to Go. The passed C array still needs to be freed manually.
Generics please!!!
*/

// WARN: this function doesn't create a copy of each byte-vector (item) but references it in the C source instead
// Therefore, it's only supposed to be used intra-library (usually inside a read transaction)
func cBytesArrayToGo(cArray *C.OBX_bytes_array) [][]byte {
	size := int(cArray.count)
	plainBytesArray := make([][]byte, size)

	if size > 0 {
		var sliceOfCBytes []C.OBX_bytes
		// see cVoidPtrToByteSlice for documentation of the following approach in general
		header := (*reflect.SliceHeader)(unsafe.Pointer(&sliceOfCBytes))
		header.Data = uintptr(unsafe.Pointer(cArray.bytes))
		header.Len = size
		header.Cap = size

		for i := 0; i < size; i++ {
			cBytes := sliceOfCBytes[i]
			cVoidPtrToByteSlice(unsafe.Pointer(cBytes.data), int(cBytes.size), &(plainBytesArray[i])) // reference
		}
	}

	return plainBytesArray
}

func cStringsArrayToGo(cArray *C.OBX_string_array) []string {
	var size = uint(cArray.count)
	var items = make([]string, size)
	if size > 0 {
		var cArrayStart = unsafe.Pointer(cArray.items)
		var cItemSize = unsafe.Sizeof(*cArray.items)
		for i := uint(0); i < size; i++ {
			var cCharPtr = (**C.char)(unsafe.Pointer(uintptr(cArrayStart) + uintptr(i)*cItemSize))
			items[i] = C.GoString(*cCharPtr) // make a copy
		}
	}
	return items
}

func cIdsArrayToGo(cArray *C.OBX_id_array) []uint64 {
	var size = uint(cArray.count)
	var items = make([]uint64, size)
	if size > 0 {
		var cArrayStart = unsafe.Pointer(cArray.ids)
		var cItemSize = unsafe.Sizeof(*cArray.ids)
		for i := uint(0); i < size; i++ {
			items[i] = *(*uint64)(unsafe.Pointer(uintptr(cArrayStart) + uintptr(i)*cItemSize)) // make a copy
		}
	}
	return items
}

func cIntsArrayToGo(cArray *C.OBX_int64_array) []int {
	var size = uint(cArray.count)
	var result = make([]int, size)
	if size > 0 {
		var cArrayStart = unsafe.Pointer(cArray.items)
		var cItemSize = unsafe.Sizeof(*cArray.items)
		for i := uint(0); i < size; i++ {
			result[i] = int(*(*int64)(unsafe.Pointer(uintptr(cArrayStart) + uintptr(i)*cItemSize))) // make a copy
		}
	}
	return result
}

func cUintsArrayToGo(cArray *C.OBX_int64_array) []uint {
	var size = uint(cArray.count)
	var result = make([]uint, size)
	if size > 0 {
		var cArrayStart = unsafe.Pointer(cArray.items)
		var cItemSize = unsafe.Sizeof(*cArray.items)
		for i := uint(0); i < size; i++ {
			result[i] = uint(*(*uint64)(unsafe.Pointer(uintptr(cArrayStart) + uintptr(i)*cItemSize))) // make a copy
		}
	}
	return result
}

func cInt64sArrayToGo(cArray *C.OBX_int64_array) []int64 {
	var size = uint(cArray.count)
	var result = make([]int64, size)
	if size > 0 {
		var cArrayStart = unsafe.Pointer(cArray.items)
		var cItemSize = unsafe.Sizeof(*cArray.items)
		for i := uint(0); i < size; i++ {
			result[i] = *(*int64)(unsafe.Pointer(uintptr(cArrayStart) + uintptr(i)*cItemSize)) // make a copy
		}
	}
	return result
}

func cUint64sArrayToGo(cArray *C.OBX_int64_array) []uint64 {
	var size = uint(cArray.count)
	var result = make([]uint64, size)
	if size > 0 {
		var cArrayStart = unsafe.Pointer(cArray.items)
		var cItemSize = unsafe.Sizeof(*cArray.items)
		for i := uint(0); i < size; i++ {
			result[i] = *(*uint64)(unsafe.Pointer(uintptr(cArrayStart) + uintptr(i)*cItemSize)) // make a copy
		}
	}
	return result
}

func cInt32sArrayToGo(cArray *C.OBX_int32_array) []int32 {
	var size = uint(cArray.count)
	var result = make([]int32, size)
	if size > 0 {
		var cArrayStart = unsafe.Pointer(cArray.items)
		var cItemSize = unsafe.Sizeof(*cArray.items)
		for i := uint(0); i < size; i++ {
			result[i] = *(*int32)(unsafe.Pointer(uintptr(cArrayStart) + uintptr(i)*cItemSize)) // make a copy
		}
	}
	return result
}

func cUint32sArrayToGo(cArray *C.OBX_int32_array) []uint32 {
	var size = uint(cArray.count)
	var result = make([]uint32, size)
	if size > 0 {
		var cArrayStart = unsafe.Pointer(cArray.items)
		var cItemSize = unsafe.Sizeof(*cArray.items)
		for i := uint(0); i < size; i++ {
			result[i] = *(*uint32)(unsafe.Pointer(uintptr(cArrayStart) + uintptr(i)*cItemSize)) // make a copy
		}
	}
	return result
}

func cInt16sArrayToGo(cArray *C.OBX_int16_array) []int16 {
	var size = uint(cArray.count)
	var result = make([]int16, size)
	if size > 0 {
		var cArrayStart = unsafe.Pointer(cArray.items)
		var cItemSize = unsafe.Sizeof(*cArray.items)
		for i := uint(0); i < size; i++ {
			result[i] = *(*int16)(unsafe.Pointer(uintptr(cArrayStart) + uintptr(i)*cItemSize)) // make a copy
		}
	}
	return result
}

func cUint16sArrayToGo(cArray *C.OBX_int16_array) []uint16 {
	var size = uint(cArray.count)
	var result = make([]uint16, size)
	if size > 0 {
		var cArrayStart = unsafe.Pointer(cArray.items)
		var cItemSize = unsafe.Sizeof(*cArray.items)
		for i := uint(0); i < size; i++ {
			result[i] = *(*uint16)(unsafe.Pointer(uintptr(cArrayStart) + uintptr(i)*cItemSize)) // make a copy
		}
	}
	return result
}

func cInt8sArrayToGo(cArray *C.OBX_int8_array) []int8 {
	var size = uint(cArray.count)
	var result = make([]int8, size)
	if size > 0 {
		var cArrayStart = unsafe.Pointer(cArray.items)
		var cItemSize = unsafe.Sizeof(*cArray.items)
		for i := uint(0); i < size; i++ {
			result[i] = *(*int8)(unsafe.Pointer(uintptr(cArrayStart) + uintptr(i)*cItemSize)) // make a copy
		}
	}
	return result
}

func cUint8sArrayToGo(cArray *C.OBX_int8_array) []uint8 {
	var size = uint(cArray.count)
	var result = make([]uint8, size)
	if size > 0 {
		var cArrayStart = unsafe.Pointer(cArray.items)
		var cItemSize = unsafe.Sizeof(*cArray.items)
		for i := uint(0); i < size; i++ {
			result[i] = *(*uint8)(unsafe.Pointer(uintptr(cArrayStart) + uintptr(i)*cItemSize)) // make a copy
		}
	}
	return result
}

func cBoolsArrayToGo(cArray *C.OBX_int8_array) []bool {
	var size = uint(cArray.count)
	var result = make([]bool, size)
	if size > 0 {
		var cArrayStart = unsafe.Pointer(cArray.items)
		var cItemSize = unsafe.Sizeof(*cArray.items)
		for i := uint(0); i < size; i++ {
			result[i] = 1 == *(*uint8)(unsafe.Pointer(uintptr(cArrayStart) + uintptr(i)*cItemSize)) // make a copy
		}
	}
	return result
}

func cFloat64sArrayToGo(cArray *C.OBX_double_array) []float64 {
	var size = uint(cArray.count)
	var result = make([]float64, size)
	if size > 0 {
		var cArrayStart = unsafe.Pointer(cArray.items)
		var cItemSize = unsafe.Sizeof(*cArray.items)
		for i := uint(0); i < size; i++ {
			result[i] = *(*float64)(unsafe.Pointer(uintptr(cArrayStart) + uintptr(i)*cItemSize)) // make a copy
		}
	}
	return result
}

func cFloat32sArrayToGo(cArray *C.OBX_float_array) []float32 {
	var size = uint(cArray.count)
	var result = make([]float32, size)
	if size > 0 {
		var cArrayStart = unsafe.Pointer(cArray.items)
		var cItemSize = unsafe.Sizeof(*cArray.items)
		for i := uint(0); i < size; i++ {
			result[i] = *(*float32)(unsafe.Pointer(uintptr(cArrayStart) + uintptr(i)*cItemSize)) // make a copy
		}
	}
	return result
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
