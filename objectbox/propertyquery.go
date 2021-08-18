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
#include <stdlib.h>
#include "objectbox.h"
*/
import "C"
import (
	"fmt"
	"runtime"
	"sync"
	"unsafe"
)

// PropertyQuery provides access to values or aggregate functions over a single property (entity field).
type PropertyQuery struct {
	cPropQuery *C.OBX_query_prop
	closeMutex sync.Mutex
	query      *Query
}

func newPropertyQuery(query *Query, propertyId TypeId) (*PropertyQuery, error) {
	var pq = &PropertyQuery{query: query}

	if err := cCallBool(func() bool {
		pq.cPropQuery = C.obx_query_prop(query.cQuery, C.obx_schema_id(propertyId))
		return pq.cPropQuery != nil
	}); err != nil {
		return nil, err
	}

	runtime.SetFinalizer(pq, propQueryFinalizer)
	return pq, nil
}

// Close frees (native) resources held by this PropertyQuery.
// While SetFinalizer() is used to close automatically after GC, it's usually still preferable to close() manually
// after you don't need the object anymore.
func (pq *PropertyQuery) Close() error {
	pq.closeMutex.Lock()
	defer pq.closeMutex.Unlock()

	if pq.cPropQuery != nil {
		return cCall(func() C.obx_err {
			var err = C.obx_query_prop_close(pq.cPropQuery)
			pq.cPropQuery = nil
			runtime.SetFinalizer(pq, nil) // remove
			return err
		})
	}

	return nil
}

func propQueryFinalizer(pq *PropertyQuery) {
	err := pq.Close()
	if err != nil {
		fmt.Printf("Error in PropertyQuery finalizer: %s", err)
	}
}

// Distinct configures the property query to work only on distinct values.
// Note: not all methods support distinct, those that don't will return an error.
func (pq *PropertyQuery) Distinct(value bool) error {
	return cCall(func() C.obx_err {
		return C.obx_query_prop_distinct(pq.cPropQuery, C.bool(value))
	})
}

// DistinctString configures the property query to work only on distinct values.
// Note: not all methods support distinct, those that don't will return an error.
func (pq *PropertyQuery) DistinctString(value, caseSensitive bool) error {
	return cCall(func() C.obx_err {
		return C.obx_query_prop_distinct_case(pq.cPropQuery, C.bool(value), C.bool(caseSensitive))
	})
}

// Count returns a number of non-NULL values of the given property across all objects matching the query.
func (pq *PropertyQuery) Count() (uint64, error) {
	var cResult C.uint64_t
	if err := cCall(func() C.obx_err { return C.obx_query_prop_count(pq.cPropQuery, &cResult) }); err != nil {
		return 0, err
	}
	return uint64(cResult), nil
}

// Average returns an average value for the given numeric property across all objects matching the query.
func (pq *PropertyQuery) Average() (float64, error) {
	var cResult C.double
	var cCount C.int64_t
	if err := cCall(func() C.obx_err { return C.obx_query_prop_avg(pq.cPropQuery, &cResult, &cCount) }); err != nil {
		return 0, err
	}
	return float64(cResult), nil
}

// MinFloat64 finds the minimum value of the given floating-point property across all objects matching the query.
func (pq *PropertyQuery) MinFloat64() (float64, error) {
	var cResult C.double
	if err := cCall(func() C.obx_err { return C.obx_query_prop_min(pq.cPropQuery, &cResult, nil) }); err != nil {
		return 0, err
	}
	return float64(cResult), nil
}

// MaxFloat64 finds the maximum value of the given floating-point property across all objects matching the query.
func (pq *PropertyQuery) MaxFloat64() (float64, error) {
	var cResult C.double
	if err := cCall(func() C.obx_err { return C.obx_query_prop_max(pq.cPropQuery, &cResult, nil) }); err != nil {
		return 0, err
	}
	return float64(cResult), nil
}

// SumFloat64 calculates the sum of the given floating-point property across all objects matching the query.
func (pq *PropertyQuery) SumFloat64() (float64, error) {
	var cResult C.double
	if err := cCall(func() C.obx_err { return C.obx_query_prop_sum(pq.cPropQuery, &cResult, nil) }); err != nil {
		return 0, err
	}
	return float64(cResult), nil
}

// Min finds the minimum value of the given property across all objects matching the query.
func (pq *PropertyQuery) Min() (int64, error) {
	var cResult C.int64_t
	if err := cCall(func() C.obx_err { return C.obx_query_prop_min_int(pq.cPropQuery, &cResult, nil) }); err != nil {
		return 0, err
	}
	return int64(cResult), nil
}

// Max finds the maximum value of the given property across all objects matching the query.
func (pq *PropertyQuery) Max() (int64, error) {
	var cResult C.int64_t
	if err := cCall(func() C.obx_err { return C.obx_query_prop_max_int(pq.cPropQuery, &cResult, nil) }); err != nil {
		return 0, err
	}
	return int64(cResult), nil
}

// Sum calculates the sum of the given property across all objects matching the query.
func (pq *PropertyQuery) Sum() (int64, error) {
	var cResult C.int64_t
	if err := cCall(func() C.obx_err { return C.obx_query_prop_sum_int(pq.cPropQuery, &cResult, nil) }); err != nil {
		return 0, err
	}
	return int64(cResult), nil
}

// FindInts returns an int slice composed of values of the given property across all objects matching the query.
// Parameter valueIfNil - value that should be returned instead of NULL values on object fields.
// If `valueIfNil = nil` is given, objects with NULL values of the specified field are skipped.
func (pq *PropertyQuery) FindInts(valueIfNil *int) ([]int, error) {
	return cGetInts(func() *C.OBX_int64_array {
		if valueIfNil == nil {
			return C.obx_query_prop_find_int64s(pq.cPropQuery, nil)
		} else {
			var cValueIfNil = C.int64_t(*valueIfNil)
			return C.obx_query_prop_find_int64s(pq.cPropQuery, &cValueIfNil)
		}
	})
}

// FindUints returns an uint slice composed of values of the given property across all objects matching the query.
// Parameter valueIfNil - value that should be returned instead of NULL values on object fields.
// If `valueIfNil = nil` is given, objects with NULL values of the specified field are skipped.
func (pq *PropertyQuery) FindUints(valueIfNil *uint) ([]uint, error) {
	return cGetUints(func() *C.OBX_int64_array {
		if valueIfNil == nil {
			return C.obx_query_prop_find_int64s(pq.cPropQuery, nil)
		} else {
			var cValueIfNil = C.int64_t(*valueIfNil)
			return C.obx_query_prop_find_int64s(pq.cPropQuery, &cValueIfNil)
		}
	})
}

// FindInt64s returns an int64 slice composed of values of the given property across all objects matching the query.
// Parameter valueIfNil - value that should be returned instead of NULL values on object fields.
// If `valueIfNil = nil` is given, objects with NULL values of the specified field are skipped.
func (pq *PropertyQuery) FindInt64s(valueIfNil *int64) ([]int64, error) {
	return cGetInt64s(func() *C.OBX_int64_array {
		return C.obx_query_prop_find_int64s(pq.cPropQuery, (*C.int64_t)(valueIfNil))
	})
}

// FindUint64s returns an uint64 slice composed of values of the given property across all objects matching the query.
// Parameter valueIfNil - value that should be returned instead of NULL values on object fields.
// If `valueIfNil = nil` is given, objects with NULL values of the specified field are skipped.
func (pq *PropertyQuery) FindUint64s(valueIfNil *uint64) ([]uint64, error) {
	return cGetUint64s(func() *C.OBX_int64_array {
		if valueIfNil == nil {
			return C.obx_query_prop_find_int64s(pq.cPropQuery, nil)
		} else {
			var cValueIfNil = C.int64_t(*valueIfNil)
			return C.obx_query_prop_find_int64s(pq.cPropQuery, &cValueIfNil)
		}
	})
}

// FindInt32s returns an int32 slice composed of values of the given property across all objects matching the query.
// Parameter valueIfNil - value that should be returned instead of NULL values on object fields.
// If `valueIfNil = nil` is given, objects with NULL values of the specified field are skipped.
func (pq *PropertyQuery) FindInt32s(valueIfNil *int32) ([]int32, error) {
	return cGetInt32s(func() *C.OBX_int32_array {
		return C.obx_query_prop_find_int32s(pq.cPropQuery, (*C.int32_t)(valueIfNil))
	})
}

// FindUint32s returns an uint32 slice composed of values of the given property across all objects matching the query.
// Parameter valueIfNil - value that should be returned instead of NULL values on object fields.
// If `valueIfNil = nil` is given, objects with NULL values of the specified field are skipped.
func (pq *PropertyQuery) FindUint32s(valueIfNil *uint32) ([]uint32, error) {
	return cGetUint32s(func() *C.OBX_int32_array {
		if valueIfNil == nil {
			return C.obx_query_prop_find_int32s(pq.cPropQuery, nil)
		} else {
			var cValueIfNil = C.int32_t(*valueIfNil)
			return C.obx_query_prop_find_int32s(pq.cPropQuery, &cValueIfNil)
		}
	})
}

// FindInt16s returns an int16 slice composed of values of the given property across all objects matching the query.
// Parameter valueIfNil - value that should be returned instead of NULL values on object fields.
// If `valueIfNil = nil` is given, objects with NULL values of the specified field are skipped.
func (pq *PropertyQuery) FindInt16s(valueIfNil *int16) ([]int16, error) {
	return cGetInt16s(func() *C.OBX_int16_array {
		return C.obx_query_prop_find_int16s(pq.cPropQuery, (*C.int16_t)(valueIfNil))
	})
}

// FindUint16s returns an uint16 slice composed of values of the given property across all objects matching the query.
// Parameter valueIfNil - value that should be returned instead of NULL values on object fields.
// If `valueIfNil = nil` is given, objects with NULL values of the specified field are skipped.
func (pq *PropertyQuery) FindUint16s(valueIfNil *uint16) ([]uint16, error) {
	return cGetUint16s(func() *C.OBX_int16_array {
		if valueIfNil == nil {
			return C.obx_query_prop_find_int16s(pq.cPropQuery, nil)
		} else {
			var cValueIfNil = C.int16_t(*valueIfNil)
			return C.obx_query_prop_find_int16s(pq.cPropQuery, &cValueIfNil)
		}
	})
}

// FindInt8s returns an int8 slice composed of values of the given property across all objects matching the query.
// Parameter valueIfNil - value that should be returned instead of NULL values on object fields.
// If `valueIfNil = nil` is given, objects with NULL values of the specified field are skipped.
func (pq *PropertyQuery) FindInt8s(valueIfNil *int8) ([]int8, error) {
	return cGetInt8s(func() *C.OBX_int8_array {
		return C.obx_query_prop_find_int8s(pq.cPropQuery, (*C.int8_t)(valueIfNil))
	})
}

// FindUint8s returns an int8 slice composed of values of the given property across all objects matching the query.
// Parameter valueIfNil - value that should be returned instead of NULL values on object fields.
// If `valueIfNil = nil` is given, objects with NULL values of the specified field are skipped.
func (pq *PropertyQuery) FindUint8s(valueIfNil *uint8) ([]uint8, error) {
	return cGetUint8s(func() *C.OBX_int8_array {
		if valueIfNil == nil {
			return C.obx_query_prop_find_int8s(pq.cPropQuery, nil)
		} else {
			var cValueIfNil = C.int8_t(*valueIfNil)
			return C.obx_query_prop_find_int8s(pq.cPropQuery, &cValueIfNil)
		}
	})
}

// FindFloat64s returns a float64 slice composed of values of the given property across all objects matching the query.
// Parameter valueIfNil - value that should be returned instead of NULL values on object fields.
// If `valueIfNil = nil` is given, objects with NULL values of the specified field are skipped.
func (pq *PropertyQuery) FindFloat64s(valueIfNil *float64) ([]float64, error) {
	return cGetFloat64s(func() *C.OBX_double_array {
		return C.obx_query_prop_find_doubles(pq.cPropQuery, (*C.double)(valueIfNil))
	})
}

// FindFloat32s returns a float32 slice composed of values of the given property across all objects matching the query.
// Parameter valueIfNil - value that should be returned instead of NULL values on object fields.
// If `valueIfNil = nil` is given, objects with NULL values of the specified field are skipped.
func (pq *PropertyQuery) FindFloat32s(valueIfNil *float32) ([]float32, error) {
	return cGetFloat32s(func() *C.OBX_float_array {
		return C.obx_query_prop_find_floats(pq.cPropQuery, (*C.float)(valueIfNil))
	})
}

// FindBools returns an int8 slice composed of values of the given property across all objects matching the query.
// Parameter valueIfNil - value that should be returned instead of NULL values on object fields.
// If `valueIfNil = nil` is given, objects with NULL values of the specified field are skipped.
func (pq *PropertyQuery) FindBools(valueIfNil *bool) ([]bool, error) {
	return cGetBools(func() *C.OBX_int8_array {
		if valueIfNil == nil {
			return C.obx_query_prop_find_int8s(pq.cPropQuery, nil)
		} else {
			var cValueIfNil = C.int8_t(0)
			if *valueIfNil {
				cValueIfNil = 1
			}
			return C.obx_query_prop_find_int8s(pq.cPropQuery, &cValueIfNil)
		}
	})
}

// FindStrings returns a string slice composed of values of the given property across all objects matching the query.
// Parameter valueIfNil - value that should be returned instead of NULL values on object fields.
// If `valueIfNil = nil` is given, objects with NULL values of the specified field are skipped.
func (pq *PropertyQuery) FindStrings(valueIfNil *string) ([]string, error) {
	return cGetStrings(func() *C.OBX_string_array {
		if valueIfNil == nil {
			return C.obx_query_prop_find_strings(pq.cPropQuery, nil)
		} else {
			var cValueIfNil = C.CString(*valueIfNil)
			defer C.free(unsafe.Pointer(cValueIfNil))
			return C.obx_query_prop_find_strings(pq.cPropQuery, cValueIfNil)
		}
	})
}
