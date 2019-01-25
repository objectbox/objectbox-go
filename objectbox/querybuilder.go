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
*/
import "C"
import (
	"errors"
	"fmt"
	"unsafe"
)

// Internal class; use Box.Query instead.
// Allows construction of queries; just check queryBuilder.Error or err from Build()
type QueryBuilder struct {
	objectBox *ObjectBox
	cqb       *C.OBX_query_builder
	typeId    TypeId

	// The first error that occurred during a any of the calls on the query builder
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

func (qb *QueryBuilder) setError(err error) {
	if qb.Err == nil {
		qb.Err = err
	}
}

func (qb *QueryBuilder) Build() (*Query, error) {
	qb.checkForCError()
	if qb.Err != nil {
		return nil, qb.Err
	}
	cQuery, err := C.obx_query_create(qb.cqb)
	if err != nil {
		return nil, err
	}
	query := &Query{
		objectBox: qb.objectBox,
		cQuery:    cQuery,
		entity:    qb.objectBox.entities[qb.typeId],
	}
	query.installFinalizer()
	return query, nil
}

func (qb *QueryBuilder) BuildWithConditions(conditions ...Condition) (*Query, error) {
	var condition Condition
	if len(conditions) == 1 {
		condition = conditions[0]
	} else {
		condition = &conditionCombination{
			conditions: conditions,
		}
	}

	var err error
	if _, err = condition.applyTo(qb); err != nil {
		return nil, err
	}
	return qb.Build()
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

func (qb *QueryBuilder) checkEntityId(entityId TypeId) bool {
	if qb.typeId == entityId {
		return true
	}

	if qb.Err == nil {
		qb.Err = fmt.Errorf("property from a different entity %d passed, expected %d", entityId, qb.typeId)
	}

	return false
}

func (qb *QueryBuilder) getConditionId(cid C.obx_qb_cond) ConditionId {
	if cid == 0 {
		// we only need to check & store the error if cid is 0, otherwise there can't be any error
		qb.checkForCError()
	}

	return ConditionId(cid)
}

func (qb *QueryBuilder) Null(property *BaseProperty) (ConditionId, error) {
	var cid ConditionId

	if qb.Err == nil && qb.checkEntityId(property.Entity.Id) {
		cid = qb.getConditionId(C.obx_qb_null(qb.cqb, C.obx_schema_id(property.Id)))
	}

	return cid, qb.Err
}

func (qb *QueryBuilder) NotNull(property *BaseProperty) (ConditionId, error) {
	var cid ConditionId

	if qb.Err == nil && qb.checkEntityId(property.Entity.Id) {
		cid = qb.getConditionId(C.obx_qb_not_null(qb.cqb, C.obx_schema_id(property.Id)))
	}

	return cid, qb.Err
}

func (qb *QueryBuilder) StringEquals(property *BaseProperty, value string, caseSensitive bool) (ConditionId, error) {
	var cid ConditionId

	if qb.Err == nil && qb.checkEntityId(property.Entity.Id) {
		cvalue := C.CString(value)
		defer C.free(unsafe.Pointer(cvalue))
		cid = qb.getConditionId(C.obx_qb_string_equal(qb.cqb, C.obx_schema_id(property.Id), cvalue, C.bool(caseSensitive)))
	}

	return cid, qb.Err
}

func (qb *QueryBuilder) StringIn(property *BaseProperty, values []string, caseSensitive bool) (ConditionId, error) {
	var cid ConditionId

	if qb.Err == nil && qb.checkEntityId(property.Entity.Id) {
		cStringArray := goStringArrayToC(values)
		defer cStringArray.free()
		cid = qb.getConditionId(C.obx_qb_string_in(qb.cqb, C.obx_schema_id(property.Id), cStringArray.cArray, C.int(cStringArray.size), C.bool(caseSensitive)))
	}

	return cid, qb.Err
}

func (qb *QueryBuilder) StringContains(property *BaseProperty, value string, caseSensitive bool) (ConditionId, error) {
	var cid ConditionId

	if qb.Err == nil && qb.checkEntityId(property.Entity.Id) {
		cvalue := C.CString(value)
		defer C.free(unsafe.Pointer(cvalue))
		cid = qb.getConditionId(C.obx_qb_string_contains(qb.cqb, C.obx_schema_id(property.Id), cvalue, C.bool(caseSensitive)))
	}

	return cid, qb.Err
}

func (qb *QueryBuilder) StringHasPrefix(property *BaseProperty, value string, caseSensitive bool) (ConditionId, error) {
	var cid ConditionId

	if qb.Err == nil && qb.checkEntityId(property.Entity.Id) {
		cvalue := C.CString(value)
		defer C.free(unsafe.Pointer(cvalue))
		cid = qb.getConditionId(C.obx_qb_string_starts_with(qb.cqb, C.obx_schema_id(property.Id), cvalue, C.bool(caseSensitive)))
	}

	return cid, qb.Err
}

func (qb *QueryBuilder) StringHasSuffix(property *BaseProperty, value string, caseSensitive bool) (ConditionId, error) {
	var cid ConditionId

	if qb.Err == nil && qb.checkEntityId(property.Entity.Id) {
		cvalue := C.CString(value)
		defer C.free(unsafe.Pointer(cvalue))
		cid = qb.getConditionId(C.obx_qb_string_ends_with(qb.cqb, C.obx_schema_id(property.Id), cvalue, C.bool(caseSensitive)))
	}

	return cid, qb.Err
}

func (qb *QueryBuilder) StringNotEquals(property *BaseProperty, value string, caseSensitive bool) (ConditionId, error) {
	var cid ConditionId

	if qb.Err == nil && qb.checkEntityId(property.Entity.Id) {
		cvalue := C.CString(value)
		defer C.free(unsafe.Pointer(cvalue))
		cid = qb.getConditionId(C.obx_qb_string_not_equal(qb.cqb, C.obx_schema_id(property.Id), cvalue, C.bool(caseSensitive)))
	}

	return cid, qb.Err
}

func (qb *QueryBuilder) StringGreater(property *BaseProperty, value string, caseSensitive bool, withEqual bool) (ConditionId, error) {
	var cid ConditionId

	if qb.Err == nil && qb.checkEntityId(property.Entity.Id) {
		cvalue := C.CString(value)
		defer C.free(unsafe.Pointer(cvalue))
		cid = qb.getConditionId(C.obx_qb_string_greater(qb.cqb, C.obx_schema_id(property.Id), cvalue, C.bool(caseSensitive), C.bool(withEqual)))
	}

	return cid, qb.Err
}

func (qb *QueryBuilder) StringLess(property *BaseProperty, value string, caseSensitive bool, withEqual bool) (ConditionId, error) {
	var cid ConditionId

	if qb.Err == nil && qb.checkEntityId(property.Entity.Id) {
		cvalue := C.CString(value)
		defer C.free(unsafe.Pointer(cvalue))
		cid = qb.getConditionId(C.obx_qb_string_less(qb.cqb, C.obx_schema_id(property.Id), cvalue, C.bool(caseSensitive), C.bool(withEqual)))
	}

	return cid, qb.Err
}

func (qb *QueryBuilder) IntBetween(property *BaseProperty, value1 int64, value2 int64) (ConditionId, error) {
	var cid ConditionId

	if qb.Err == nil && qb.checkEntityId(property.Entity.Id) {
		cid = qb.getConditionId(C.obx_qb_int_between(qb.cqb, C.obx_schema_id(property.Id), C.int64_t(value1), C.int64_t(value2)))
	}

	return cid, qb.Err
}

func (qb *QueryBuilder) IntEqual(property *BaseProperty, value int64) (ConditionId, error) {
	var cid ConditionId

	if qb.Err == nil && qb.checkEntityId(property.Entity.Id) {
		cid = qb.getConditionId(C.obx_qb_int_equal(qb.cqb, C.obx_schema_id(property.Id), C.int64_t(value)))
	}

	return cid, qb.Err
}

func (qb *QueryBuilder) IntNotEqual(property *BaseProperty, value int64) (ConditionId, error) {
	var cid ConditionId

	if qb.Err == nil && qb.checkEntityId(property.Entity.Id) {
		cid = qb.getConditionId(C.obx_qb_int_not_equal(qb.cqb, C.obx_schema_id(property.Id), C.int64_t(value)))
	}

	return cid, qb.Err
}

func (qb *QueryBuilder) IntGreater(property *BaseProperty, value int64) (ConditionId, error) {
	var cid ConditionId

	if qb.Err == nil && qb.checkEntityId(property.Entity.Id) {
		cid = qb.getConditionId(C.obx_qb_int_greater(qb.cqb, C.obx_schema_id(property.Id), C.int64_t(value)))
	}

	return cid, qb.Err
}

func (qb *QueryBuilder) IntLess(property *BaseProperty, value int64) (ConditionId, error) {
	var cid ConditionId

	if qb.Err == nil && qb.checkEntityId(property.Entity.Id) {
		cid = qb.getConditionId(C.obx_qb_int_less(qb.cqb, C.obx_schema_id(property.Id), C.int64_t(value)))
	}

	return cid, qb.Err
}

func (qb *QueryBuilder) Int64In(property *BaseProperty, values []int64) (ConditionId, error) {
	var cid ConditionId

	if qb.Err == nil && qb.checkEntityId(property.Entity.Id) {
		cid = qb.getConditionId(C.obx_qb_int64_in(qb.cqb, C.obx_schema_id(property.Id), (*C.int64_t)(unsafe.Pointer(&values[0])), C.int(len(values))))
	}

	return cid, qb.Err
}

func (qb *QueryBuilder) Int64NotIn(property *BaseProperty, values []int64) (ConditionId, error) {
	var cid ConditionId

	if qb.Err == nil && qb.checkEntityId(property.Entity.Id) {
		cid = qb.getConditionId(C.obx_qb_int64_not_in(qb.cqb, C.obx_schema_id(property.Id), (*C.int64_t)(unsafe.Pointer(&values[0])), C.int(len(values))))
	}

	return cid, qb.Err
}

func (qb *QueryBuilder) Int32In(property *BaseProperty, values []int32) (ConditionId, error) {
	var cid ConditionId

	if qb.Err == nil && qb.checkEntityId(property.Entity.Id) {
		cid = qb.getConditionId(C.obx_qb_int32_in(qb.cqb, C.obx_schema_id(property.Id), (*C.int32_t)(unsafe.Pointer(&values[0])), C.int(len(values))))
	}

	return cid, qb.Err
}

func (qb *QueryBuilder) Int32NotIn(property *BaseProperty, values []int32) (ConditionId, error) {
	var cid ConditionId

	if qb.Err == nil && qb.checkEntityId(property.Entity.Id) {
		cid = qb.getConditionId(C.obx_qb_int32_not_in(qb.cqb, C.obx_schema_id(property.Id), (*C.int32_t)(unsafe.Pointer(&values[0])), C.int(len(values))))
	}

	return cid, qb.Err
}

func (qb *QueryBuilder) DoubleGreater(property *BaseProperty, value float64) (ConditionId, error) {
	var cid ConditionId

	if qb.Err == nil && qb.checkEntityId(property.Entity.Id) {
		cid = qb.getConditionId(C.obx_qb_double_greater(qb.cqb, C.obx_schema_id(property.Id), C.double(value)))
	}

	return cid, qb.Err
}

func (qb *QueryBuilder) DoubleLess(property *BaseProperty, value float64) (ConditionId, error) {
	var cid ConditionId

	if qb.Err == nil && qb.checkEntityId(property.Entity.Id) {
		cid = qb.getConditionId(C.obx_qb_double_less(qb.cqb, C.obx_schema_id(property.Id), C.double(value)))
	}

	return cid, qb.Err
}

func (qb *QueryBuilder) DoubleBetween(property *BaseProperty, valueA float64, valueB float64) (ConditionId, error) {
	var cid ConditionId

	if qb.Err == nil && qb.checkEntityId(property.Entity.Id) {
		cid = qb.getConditionId(C.obx_qb_double_between(qb.cqb, C.obx_schema_id(property.Id), C.double(valueA), C.double(valueB)))
	}

	return cid, qb.Err
}

func (qb *QueryBuilder) BytesEqual(property *BaseProperty, value []byte) (ConditionId, error) {
	var cid ConditionId

	if qb.Err == nil && qb.checkEntityId(property.Entity.Id) {
		cid = qb.getConditionId(C.obx_qb_bytes_equal(qb.cqb, C.obx_schema_id(property.Id), cBytesPtr(value), C.size_t(len(value))))
	}

	return cid, qb.Err
}

func (qb *QueryBuilder) BytesGreater(property *BaseProperty, value []byte, withEqual bool) (ConditionId, error) {
	var cid ConditionId

	if qb.Err == nil && qb.checkEntityId(property.Entity.Id) {
		cid = qb.getConditionId(C.obx_qb_bytes_greater(qb.cqb, C.obx_schema_id(property.Id), cBytesPtr(value), C.size_t(len(value)), C.bool(withEqual)))
	}

	return cid, qb.Err
}

func (qb *QueryBuilder) BytesLess(property *BaseProperty, value []byte, withEqual bool) (ConditionId, error) {
	var cid ConditionId

	if qb.Err == nil && qb.checkEntityId(property.Entity.Id) {
		cid = qb.getConditionId(C.obx_qb_bytes_less(qb.cqb, C.obx_schema_id(property.Id), cBytesPtr(value), C.size_t(len(value)), C.bool(withEqual)))
	}

	return cid, qb.Err
}
