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


static char**newCharArray(int size) {
        return calloc(sizeof(char*), size);
}

static void setArrayString(const char **array, size_t index, const char *value) {
        array[index] = value;
}

static void freeCharArray(char **a, int size) {
        for (size_t i = 0; i < size; i++)
                free(a[i]);
        free(a);
}
*/
import "C"
import (
	"errors"
	"fmt"
	"unsafe"
)

// Allows construction of queries; just check QueryBuilder.Error or err from Build()
//
// WIP: Query interface is subject to change with full ObjectBox queries support
type QueryBuilder struct {
	objectBox *ObjectBox
	cqb       *C.OBX_query_builder
	typeId    TypeId

	// Any error that occurred during a call to QueryBuilder or its construction
	Err error
}

func (qb *QueryBuilder) Close() error {
	toClose := qb.cqb
	if toClose != nil {
		qb.cqb = nil
		rc := C.obx_qb_close(toClose)
		if rc != 0 {
			return createError()
		}
	}
	return nil
}

func (qb *QueryBuilder) StringEq(propertyId TypeId, value string, caseSensitive bool) (conditionId, error) {
	if qb.Err != nil {
		return 0, qb.Err
	}
	cvalue := C.CString(value)
	defer C.free(unsafe.Pointer(cvalue))
	cid := C.obx_qb_string_equal(qb.cqb, C.obx_schema_id(propertyId), cvalue, C.bool(caseSensitive))
	qb.checkForCError() // Mirror C error early to Error

	return conditionId(cid), qb.Err
}

func (qb *QueryBuilder) StringIn(propertyId TypeId, values []string, caseSensitive bool) (conditionId, error) {
	if qb.Err != nil {
		return 0, qb.Err
	}

	cStringArray := C.newCharArray(C.int(len(values)))
	defer C.freeCharArray(cStringArray, C.int(len(values)))
	for i, s := range values {
		C.setArrayString(cStringArray, C.size_t(i), C.CString(s))
	}

	cid := C.obx_qb_string_in(qb.cqb, C.obx_schema_id(propertyId), cStringArray, C.int(len(values)), C.bool(caseSensitive))
	qb.checkForCError() // Mirror C error early to Error

	return conditionId(cid), qb.Err
}

func (qb *QueryBuilder) StringContains(propertyId TypeId, value string, caseSensitive bool) (conditionId, error) {
	if qb.Err != nil {
		return 0, qb.Err
	}
	cvalue := C.CString(value)
	defer C.free(unsafe.Pointer(cvalue))
	cid := C.obx_qb_string_contains(qb.cqb, C.obx_schema_id(propertyId), cvalue, C.bool(caseSensitive))
	qb.checkForCError() // Mirror C error early to Error

	return conditionId(cid), qb.Err
}

func (qb *QueryBuilder) StringStartsWith(propertyId TypeId, value string, caseSensitive bool) (conditionId, error) {
	if qb.Err != nil {
		return 0, qb.Err
	}
	cvalue := C.CString(value)
	defer C.free(unsafe.Pointer(cvalue))
	cid := C.obx_qb_string_starts_with(qb.cqb, C.obx_schema_id(propertyId), cvalue, C.bool(caseSensitive))
	qb.checkForCError() // Mirror C error early to Error

	return conditionId(cid), qb.Err
}

func (qb *QueryBuilder) StringEndsWith(propertyId TypeId, value string, caseSensitive bool) (conditionId, error) {
	if qb.Err != nil {
		return 0, qb.Err
	}
	cvalue := C.CString(value)
	defer C.free(unsafe.Pointer(cvalue))
	cid := C.obx_qb_string_ends_with(qb.cqb, C.obx_schema_id(propertyId), cvalue, C.bool(caseSensitive))
	qb.checkForCError() // Mirror C error early to Error

	return conditionId(cid), qb.Err
}

func (qb *QueryBuilder) IntBetween(propertyId TypeId, value1 int64, value2 int64) (conditionId, error) {
	if qb.Err != nil {
		return 0, qb.Err
	}
	cid := C.obx_qb_int_between(qb.cqb, C.obx_schema_id(propertyId), C.int64_t(value1), C.int64_t(value2))
	qb.checkForCError() // Mirror C error early to Error

	return conditionId(cid), qb.Err
}

func (qb *QueryBuilder) Null(propertyId TypeId) (conditionId, error) {
	if qb.Err != nil {
		return 0, qb.Err
	}
	cid := C.obx_qb_null(qb.cqb, C.obx_schema_id(propertyId))
	qb.checkForCError() // Mirror C error early to Error

	return conditionId(cid), qb.Err
}

func (qb *QueryBuilder) StringNotEq(propertyId TypeId, value string, caseSensitive bool) (conditionId, error) {
	if qb.Err != nil {
		return 0, qb.Err
	}
	cvalue := C.CString(value)
	defer C.free(unsafe.Pointer(cvalue))
	cid := C.obx_qb_string_not_equal(qb.cqb, C.obx_schema_id(propertyId), cvalue, C.bool(caseSensitive))
	qb.checkForCError() // Mirror C error early to Error

	return conditionId(cid), qb.Err
}

func (qb *QueryBuilder) StringGreater(propertyId TypeId, value string, caseSensitive bool, withEqual bool) (conditionId, error) {
	if qb.Err != nil {
		return 0, qb.Err
	}
	cvalue := C.CString(value)
	defer C.free(unsafe.Pointer(cvalue))
	cid := C.obx_qb_string_greater(qb.cqb, C.obx_schema_id(propertyId), cvalue, C.bool(caseSensitive), C.bool(withEqual))
	qb.checkForCError() // Mirror C error early to Error

	return conditionId(cid), qb.Err
}

func (qb *QueryBuilder) StringLess(propertyId TypeId, value string, caseSensitive bool, withEqual bool) (conditionId, error) {
	if qb.Err != nil {
		return 0, qb.Err
	}
	cvalue := C.CString(value)
	defer C.free(unsafe.Pointer(cvalue))
	cid := C.obx_qb_string_less(qb.cqb, C.obx_schema_id(propertyId), cvalue, C.bool(caseSensitive), C.bool(withEqual))
	qb.checkForCError() // Mirror C error early to Error

	return conditionId(cid), qb.Err
}

func (qb *QueryBuilder) IntEqual(propertyId TypeId, value int64) (conditionId, error) {
	if qb.Err != nil {
		return 0, qb.Err
	}
	cid := C.obx_qb_int_equal(qb.cqb, C.obx_schema_id(propertyId), C.int64_t(value))
	qb.checkForCError() // Mirror C error early to Error

	return conditionId(cid), qb.Err
}

func (qb *QueryBuilder) DoubleGreater(propertyId TypeId, value float64) (conditionId, error) {
	if qb.Err != nil {
		return 0, qb.Err
	}
	cid := C.obx_qb_double_greater(qb.cqb, C.obx_schema_id(propertyId), C.double(value))
	qb.checkForCError() // Mirror C error early to Error

	return conditionId(cid), qb.Err
}

func (qb *QueryBuilder) DoubleLess(propertyId TypeId, value float64) (conditionId, error) {
	if qb.Err != nil {
		return 0, qb.Err
	}
	cid := C.obx_qb_double_less(qb.cqb, C.obx_schema_id(propertyId), C.double(value))
	qb.checkForCError() // Mirror C error early to Error

	return conditionId(cid), qb.Err
}

func (qb *QueryBuilder) DoubleBetween(propertyId TypeId, valueA float64, valueB float64) (conditionId, error) {
	if qb.Err != nil {
		return 0, qb.Err
	}
	cid := C.obx_qb_double_between(qb.cqb, C.obx_schema_id(propertyId), C.double(valueA), C.double(valueB))
	qb.checkForCError() // Mirror C error early to Error

	return conditionId(cid), qb.Err
}

// obx_qb_cond obx_qb_int_not_equal(OBX_query_builder* builder, obx_schema_id property_id, int64 value);
func (qb *QueryBuilder) IntNotEqual(propertyId TypeId, value int64) (conditionId, error) {
	if qb.Err != nil {
		return 0, qb.Err
	}
	cid := C.obx_qb_int_not_equal(qb.cqb, C.obx_schema_id(propertyId), C.int64_t(value))
	qb.checkForCError() // Mirror C error early to Error

	return conditionId(cid), qb.Err
}

// obx_qb_cond obx_qb_int_greater(OBX_query_builder* builder, obx_schema_id property_id, int64  value);
func (qb *QueryBuilder) IntGreater(propertyId TypeId, value int64) (conditionId, error) {
	if qb.Err != nil {
		return 0, qb.Err
	}
	cid := C.obx_qb_int_greater(qb.cqb, C.obx_schema_id(propertyId), C.int64_t(value))
	qb.checkForCError() // Mirror C error early to Error

	return conditionId(cid), qb.Err
}

// obx_qb_cond obx_qb_int_less(OBX_query_builder* builder, obx_schema_id property_id, int64  value);
func (qb *QueryBuilder) IntLess(propertyId TypeId, value int64) (conditionId, error) {
	if qb.Err != nil {
		return 0, qb.Err
	}
	cid := C.obx_qb_int_less(qb.cqb, C.obx_schema_id(propertyId), C.int64_t(value))
	qb.checkForCError() // Mirror C error early to Error

	return conditionId(cid), qb.Err
}

func (qb *QueryBuilder) NotNull(propertyId TypeId) (conditionId, error) {
	if qb.Err != nil {
		return 0, qb.Err
	}
	cid := C.obx_qb_not_null(qb.cqb, C.obx_schema_id(propertyId))
	qb.checkForCError() // Mirror C error early to Error

	return conditionId(cid), qb.Err
}

func (qb *QueryBuilder) BytesEqual(propertyId TypeId, value []byte) (conditionId, error) {
	if qb.Err != nil {
		return 0, qb.Err
	}

	cid := C.obx_qb_bytes_equal(qb.cqb, C.obx_schema_id(propertyId), unsafe.Pointer(&value[0]), C.size_t(len(value)))
	qb.checkForCError() // Mirror C error early to Error

	return conditionId(cid), qb.Err
}

func (qb *QueryBuilder) BytesGreater(propertyId TypeId, value []byte, withEqual bool) (conditionId, error) {
	if qb.Err != nil {
		return 0, qb.Err
	}
	cid := C.obx_qb_bytes_greater(qb.cqb, C.obx_schema_id(propertyId), unsafe.Pointer(&value[0]), C.size_t(len(value)), C.bool(withEqual))
	qb.checkForCError() // Mirror C error early to Error

	return conditionId(cid), qb.Err
}

func (qb *QueryBuilder) BytesLess(propertyId TypeId, value []byte, withEqual bool) (conditionId, error) {
	if qb.Err != nil {
		return 0, qb.Err
	}
	cid := C.obx_qb_bytes_less(qb.cqb, C.obx_schema_id(propertyId), unsafe.Pointer(&value[0]), C.size_t(len(value)), C.bool(withEqual))
	qb.checkForCError() // Mirror C error early to Error

	return conditionId(cid), qb.Err
}

func (qb *QueryBuilder) Int64In(propertyId TypeId, values []int64) (conditionId, error) {
	if qb.Err != nil {
		return 0, qb.Err
	}
	cid := C.obx_qb_int64_in(qb.cqb, C.obx_schema_id(propertyId), (*C.int64_t)(unsafe.Pointer(&values[0])), C.int(len(values)))
	qb.checkForCError() // Mirror C error early to Error

	return conditionId(cid), qb.Err
}

func (qb *QueryBuilder) Int64NotIn(propertyId TypeId, values []int64) (conditionId, error) {
	if qb.Err != nil {
		return 0, qb.Err
	}
	cid := C.obx_qb_int64_not_in(qb.cqb, C.obx_schema_id(propertyId), (*C.int64_t)(unsafe.Pointer(&values[0])), C.int(len(values)))
	qb.checkForCError() // Mirror C error early to Error

	return conditionId(cid), qb.Err
}

func (qb *QueryBuilder) Int32In(propertyId TypeId, values []int32) (conditionId, error) {
	if qb.Err != nil {
		return 0, qb.Err
	}
	cid := C.obx_qb_int32_in(qb.cqb, C.obx_schema_id(propertyId), (*C.int32_t)(unsafe.Pointer(&values[0])), C.int(len(values)))
	qb.checkForCError() // Mirror C error early to Error

	return conditionId(cid), qb.Err
}

func (qb *QueryBuilder) Int32NotIn(propertyId TypeId, values []int32) (conditionId, error) {
	if qb.Err != nil {
		return 0, qb.Err
	}
	cid := C.obx_qb_int32_not_in(qb.cqb, C.obx_schema_id(propertyId), (*C.int32_t)(unsafe.Pointer(&values[0])), C.int(len(values)))
	qb.checkForCError() // Mirror C error early to Error

	return conditionId(cid), qb.Err
}

func (qb *QueryBuilder) build() (*C.OBX_query, error) {
	qb.checkForCError()
	if qb.Err != nil {
		return nil, qb.Err
	}
	cQuery, err := C.obx_query_create(qb.cqb)
	return cQuery, err
}

func (qb *QueryBuilder) Build() (*Query, error) {
	qb.checkForCError()
	if qb.Err != nil {
		return nil, qb.Err
	}
	cquery, err := C.obx_query_create(qb.cqb)
	if err != nil {
		return nil, err
	}
	return &Query{
		objectBox: qb.objectBox,
		cquery:    cquery,
		typeId:    qb.typeId,
	}, nil
}

func (qb *QueryBuilder) checkForCError() {
	if qb.Err != nil {
		errCode := C.obx_qb_error_code(qb.cqb)
		if errCode != 0 {
			msg := C.obx_qb_error_message(qb.cqb)
			if msg == nil {
				qb.Err = errors.New(fmt.Sprintf("Could not create query builder (code %v)", int(errCode)))
			} else {
				qb.Err = errors.New(C.GoString(msg))
			}
		}
	}
}

func (qb *QueryBuilder) BuildAndClose() (*Query, error) {
	var query *Query
	var err = qb.Err
	if err == nil {
		query, err = qb.Build()
	}
	err2 := qb.Close()
	if err == nil && err2 != nil {
		query.Close()
		return nil, err2
	}
	return query, err
}
