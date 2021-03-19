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
#include "objectbox.h"
*/
import "C"
import (
	"errors"
	"fmt"
	"runtime"
	"strings"
	"unsafe"
)

// QueryBuilder is an internal class; use Box.Query instead.
// Allows construction of queries; just check queryBuilder.Error or err from Build()
type QueryBuilder struct {
	objectBox     *ObjectBox
	cqb           *C.OBX_query_builder
	typeId        TypeId
	innerBuilders []*QueryBuilder
	orderFlags    map[TypeId]C.OBXOrderFlags

	// The first error that occurred during a any of the calls on the query builder
	Err error
}

func newQueryBuilder(ob *ObjectBox, typeId TypeId) *QueryBuilder {
	var qb = &QueryBuilder{
		objectBox:  ob,
		typeId:     typeId,
		orderFlags: make(map[TypeId]C.OBXOrderFlags),
	}

	qb.Err = cCallBool(func() bool {
		qb.cqb = C.obx_query_builder(ob.store, C.obx_schema_id(typeId))
		return qb.cqb != nil
	})

	// log.Printf("QB %p created for entity %d\n", qb, typeId)
	return qb
}

func (qb *QueryBuilder) newInnerBuilder(typeId TypeId, cqb *C.OBX_query_builder) *QueryBuilder {
	if cqb == nil {
		qb.Err = createError()
		return nil
	}

	var iqb = &QueryBuilder{
		objectBox:  qb.objectBox,
		cqb:        cqb,
		typeId:     typeId,
		orderFlags: make(map[TypeId]C.OBXOrderFlags),
	}

	qb.innerBuilders = append(qb.innerBuilders, iqb)

	// log.Printf("QB %p attached inner QB %p, entity %d\n", qb, iqb, iqb.typeId)
	return iqb
}

// Close is called internally
func (qb *QueryBuilder) Close() error {
	// close inner builders while collecting errors
	var errs []string
	for _, iqb := range qb.innerBuilders {
		if err := iqb.Close(); err != nil {
			// even though there's an error, try to close the rest
			errs = append(errs, err.Error())
		}
	}

	toClose := qb.cqb
	if toClose != nil {
		qb.cqb = nil
		if err := cCall(func() C.obx_err {
			return C.obx_qb_close(toClose)
		}); err != nil {
			errs = append(errs, err.Error())
		}
	}

	if errs != nil {
		return errors.New(strings.Join(errs, "; "))
	}
	return nil
}

// Build is called internally
func (qb *QueryBuilder) Build(box *Box) (*Query, error) {
	for propertyId, orderFlags := range qb.orderFlags {
		qb.order(C.obx_schema_id(propertyId), orderFlags)
	}

	if qb.Err != nil {
		return nil, qb.Err
	}

	query := &Query{
		objectBox: qb.objectBox,
		box:       box,
		entity:    qb.objectBox.getEntityById(qb.typeId),
	}

	if err := cCallBool(func() bool {
		query.cQuery = C.obx_query(qb.cqb)
		return query.cQuery != nil
	}); err != nil {
		qb.Err = err
		return nil, err
	}

	query.installFinalizer()

	// search all inner builders recursively and collect linked entity IDs
	qb.setQueryLinkedEntityIds(query)

	return query, nil
}

func (qb *QueryBuilder) setQueryLinkedEntityIds(query *Query) {
	for _, iqb := range qb.innerBuilders {
		query.linkedEntityIds = append(query.linkedEntityIds, iqb.typeId)
		iqb.setQueryLinkedEntityIds(query)
	}
}

func (qb *QueryBuilder) applyConditions(conditions []Condition) error {
	if qb.Err != nil {
		return qb.Err
	}

	if len(conditions) == 1 {
		_, qb.Err = conditions[0].applyTo(qb, true)
	} else if len(conditions) > 1 {
		_, qb.Err = (&conditionCombination{conditions: conditions}).applyTo(qb, true)
	}

	return qb.Err
}

// LinkOneToMany is called internally
func (qb *QueryBuilder) LinkOneToMany(relation *RelationToOne, conditions []Condition) error {
	if qb.Err != nil {
		return qb.Err
	}

	// create a new "inner" query builder
	var iqb *QueryBuilder

	// for native calls/createError() in newInnerBuilder
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	cRelPropertyId := C.obx_schema_id(relation.Property.Id)
	// recognize whether it's a link or a backlink
	if relation.Property.Entity.Id == qb.typeId && relation.Target.Id != qb.typeId {
		// if property belongs to the entity of the "main" query builder & target is another entity, it's a link
		// log.Printf("QB %p creating link to entity %d over property %d", qb, relation.Target.Id, relation.Property.Id)
		iqb = qb.newInnerBuilder(relation.Target.Id, C.obx_qb_link_property(qb.cqb, cRelPropertyId))
	} else if relation.Property.Entity.Id != qb.typeId && relation.Target.Id == qb.typeId {
		// if property is not from the same entity as this query builder but the target is, it's a backlink
		// log.Printf("QB %p creating backlink from entity %d over property %d", qb, relation.Property.Entity.Id, relation.Property.Id)
		cInnerQB := C.obx_qb_backlink_property(qb.cqb, C.obx_schema_id(relation.Property.Entity.Id), cRelPropertyId)
		iqb = qb.newInnerBuilder(relation.Property.Entity.Id, cInnerQB)
	} else {
		return errors.New("relation not recognized as either link or backlink")
	}

	if iqb == nil {
		return qb.Err // this has been set by newInnerBuilder()
	}

	return iqb.applyConditions(conditions)
}

// LinkManyToMany is called internally
func (qb *QueryBuilder) LinkManyToMany(relation *RelationToMany, conditions []Condition) error {
	if qb.Err != nil {
		return qb.Err
	}

	// create a new "inner" query builder
	var iqb *QueryBuilder

	// for native calls/createError() in newInnerBuilder
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	// recognize whether it's a link or a backlink
	if relation.Source.Id == qb.typeId && relation.Target.Id != qb.typeId {
		// log.Printf("QB %p creating link to entity %d over relation %d", qb, relation.Target.Id, relation.Id)
		iqb = qb.newInnerBuilder(relation.Target.Id, C.obx_qb_link_standalone(qb.cqb, C.obx_schema_id(relation.Id)))
	} else if relation.Source.Id != qb.typeId && relation.Target.Id == qb.typeId {
		// log.Printf("QB %p creating backlink from entity %d over relation %d", qb, relation.Source.Id, relation.Id)
		iqb = qb.newInnerBuilder(relation.Source.Id, C.obx_qb_backlink_standalone(qb.cqb, C.obx_schema_id(relation.Id)))
	} else {
		return errors.New("relation not recognized as either link or backlink")
	}
	if iqb == nil {
		return qb.Err // this has been set by newInnerBuilder()
	}

	return iqb.applyConditions(conditions)
}

func (qb *QueryBuilder) order(propertyId C.obx_schema_id, flags C.OBXOrderFlags) {
	if qb.Err == nil {
		qb.Err = cCall(func() C.obx_err {
			return C.obx_qb_order(qb.cqb, propertyId, flags)
		})
	}
}

// setOrderFlag stores the order flag to be applied later before building the query
// if value is true, the flag is set, otherwise the flag is cleared (unset)
func (qb *QueryBuilder) setOrderFlag(property *BaseProperty, flag C.OBXOrderFlags, value bool) error {
	if qb.Err == nil && qb.checkEntityId(property.Entity.Id) {
		if value {
			// set the flag
			qb.orderFlags[property.Id] = qb.orderFlags[property.Id] | flag
		} else {
			// clear the flag
			qb.orderFlags[property.Id] = qb.orderFlags[property.Id] &^ flag
		}
	}
	return qb.Err
}

func (qb *QueryBuilder) orderAsc(property *BaseProperty) error {
	return qb.setOrderFlag(property, C.OBXOrderFlags_DESCENDING, false)
}

func (qb *QueryBuilder) orderDesc(property *BaseProperty) error {
	return qb.setOrderFlag(property, C.OBXOrderFlags_DESCENDING, true)
}

func (qb *QueryBuilder) orderCaseSensitive(property *BaseProperty, value bool) error {
	return qb.setOrderFlag(property, C.OBXOrderFlags_CASE_SENSITIVE, value)
}

func (qb *QueryBuilder) orderNilLast(property *BaseProperty) error {
	return qb.setOrderFlag(property, C.OBXOrderFlags_NULLS_LAST, true)
}

func (qb *QueryBuilder) orderNilAsZero(property *BaseProperty) error {
	return qb.setOrderFlag(property, C.OBXOrderFlags_NULLS_ZERO, true)
}

func (qb *QueryBuilder) checkForCError() {
	// if there's already an error logged, don't overwrite it
	if qb.Err != nil {
		return
	}

	code := C.obx_qb_error_code(qb.cqb)
	if code == 0 {
		return
	}

	msg := C.obx_qb_error_message(qb.cqb)
	if msg == nil {
		qb.Err = fmt.Errorf("unknown query builder error (code %v)", int(code))
	} else {
		qb.Err = errors.New(C.GoString(msg))
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

// Alias sets an alias for the last created condition
func (qb *QueryBuilder) Alias(alias string) error {
	if qb.Err == nil {
		qb.Err = cCall(func() C.obx_err {
			cvalue := C.CString(alias)
			defer C.free(unsafe.Pointer(cvalue))
			return C.obx_qb_param_alias(qb.cqb, cvalue)
		})
	}

	return qb.Err
}

// Any is called internally
func (qb *QueryBuilder) Any(ids []ConditionId) (ConditionId, error) {
	var cid ConditionId

	if qb.Err == nil {
		cid = qb.getConditionId(C.obx_qb_any(qb.cqb, (*C.obx_qb_cond)(unsafe.Pointer(&ids[0])), C.size_t(len(ids))))
	}

	return cid, qb.Err
}

// All is called internally
func (qb *QueryBuilder) All(ids []ConditionId) (ConditionId, error) {
	var cid ConditionId

	if qb.Err == nil {
		cid = qb.getConditionId(C.obx_qb_all(qb.cqb, (*C.obx_qb_cond)(unsafe.Pointer(&ids[0])), C.size_t(len(ids))))
	}

	return cid, qb.Err
}

// IsNil is called internally
func (qb *QueryBuilder) IsNil(property *BaseProperty) (ConditionId, error) {
	var cid ConditionId

	if qb.Err == nil && qb.checkEntityId(property.Entity.Id) {
		cid = qb.getConditionId(C.obx_qb_null(qb.cqb, C.obx_schema_id(property.Id)))
	}

	return cid, qb.Err
}

// IsNotNil is called internally
func (qb *QueryBuilder) IsNotNil(property *BaseProperty) (ConditionId, error) {
	var cid ConditionId

	if qb.Err == nil && qb.checkEntityId(property.Entity.Id) {
		cid = qb.getConditionId(C.obx_qb_not_null(qb.cqb, C.obx_schema_id(property.Id)))
	}

	return cid, qb.Err
}

// StringEquals is called internally
func (qb *QueryBuilder) StringEquals(property *BaseProperty, value string, caseSensitive bool) (ConditionId, error) {
	var cid ConditionId

	if qb.Err == nil && qb.checkEntityId(property.Entity.Id) {
		cvalue := C.CString(value)
		defer C.free(unsafe.Pointer(cvalue))
		cid = qb.getConditionId(C.obx_qb_equals_string(qb.cqb, C.obx_schema_id(property.Id), cvalue, C.bool(caseSensitive)))
	}

	return cid, qb.Err
}

// StringIn is called internally
func (qb *QueryBuilder) StringIn(property *BaseProperty, values []string, caseSensitive bool) (ConditionId, error) {
	var cid ConditionId

	if qb.Err == nil && qb.checkEntityId(property.Entity.Id) {
		if len(values) > 0 {
			cStringArray := goStringArrayToC(values)
			defer cStringArray.free()
			cid = qb.getConditionId(C.obx_qb_in_strings(qb.cqb, C.obx_schema_id(property.Id), cStringArray.cArray, C.size_t(cStringArray.size), C.bool(caseSensitive)))
		} else {
			cid = qb.getConditionId(C.obx_qb_in_strings(qb.cqb, C.obx_schema_id(property.Id), nil, 0, C.bool(caseSensitive)))
		}
	}

	return cid, qb.Err
}

// StringContains is called internally
func (qb *QueryBuilder) StringContains(property *BaseProperty, value string, caseSensitive bool) (ConditionId, error) {
	var cid ConditionId

	if qb.Err == nil && qb.checkEntityId(property.Entity.Id) {
		cvalue := C.CString(value)
		defer C.free(unsafe.Pointer(cvalue))
		cid = qb.getConditionId(C.obx_qb_contains_string(qb.cqb, C.obx_schema_id(property.Id), cvalue, C.bool(caseSensitive)))
	}

	return cid, qb.Err
}

// StringHasPrefix is called internally
func (qb *QueryBuilder) StringHasPrefix(property *BaseProperty, value string, caseSensitive bool) (ConditionId, error) {
	var cid ConditionId

	if qb.Err == nil && qb.checkEntityId(property.Entity.Id) {
		cvalue := C.CString(value)
		defer C.free(unsafe.Pointer(cvalue))
		cid = qb.getConditionId(C.obx_qb_starts_with_string(qb.cqb, C.obx_schema_id(property.Id), cvalue, C.bool(caseSensitive)))
	}

	return cid, qb.Err
}

// StringHasSuffix is called internally
func (qb *QueryBuilder) StringHasSuffix(property *BaseProperty, value string, caseSensitive bool) (ConditionId, error) {
	var cid ConditionId

	if qb.Err == nil && qb.checkEntityId(property.Entity.Id) {
		cvalue := C.CString(value)
		defer C.free(unsafe.Pointer(cvalue))
		cid = qb.getConditionId(C.obx_qb_ends_with_string(qb.cqb, C.obx_schema_id(property.Id), cvalue, C.bool(caseSensitive)))
	}

	return cid, qb.Err
}

// StringNotEquals is called internally
func (qb *QueryBuilder) StringNotEquals(property *BaseProperty, value string, caseSensitive bool) (ConditionId, error) {
	var cid ConditionId

	if qb.Err == nil && qb.checkEntityId(property.Entity.Id) {
		cvalue := C.CString(value)
		defer C.free(unsafe.Pointer(cvalue))
		cid = qb.getConditionId(C.obx_qb_not_equals_string(qb.cqb, C.obx_schema_id(property.Id), cvalue, C.bool(caseSensitive)))
	}

	return cid, qb.Err
}

// StringGreater is called internally
func (qb *QueryBuilder) StringGreater(property *BaseProperty, value string, caseSensitive bool, withEqual bool) (ConditionId, error) {
	var cid ConditionId

	if qb.Err == nil && qb.checkEntityId(property.Entity.Id) {
		cvalue := C.CString(value)
		defer C.free(unsafe.Pointer(cvalue))
		if withEqual {
			cid = qb.getConditionId(C.obx_qb_greater_or_equal_string(qb.cqb, C.obx_schema_id(property.Id), cvalue, C.bool(caseSensitive)))
		} else {
			cid = qb.getConditionId(C.obx_qb_greater_than_string(qb.cqb, C.obx_schema_id(property.Id), cvalue, C.bool(caseSensitive)))
		}
	}

	return cid, qb.Err
}

// StringLess is called internally
func (qb *QueryBuilder) StringLess(property *BaseProperty, value string, caseSensitive bool, withEqual bool) (ConditionId, error) {
	var cid ConditionId

	if qb.Err == nil && qb.checkEntityId(property.Entity.Id) {
		cvalue := C.CString(value)
		defer C.free(unsafe.Pointer(cvalue))
		if withEqual {
			cid = qb.getConditionId(C.obx_qb_less_or_equal_string(qb.cqb, C.obx_schema_id(property.Id), cvalue, C.bool(caseSensitive)))
		} else {
			cid = qb.getConditionId(C.obx_qb_less_than_string(qb.cqb, C.obx_schema_id(property.Id), cvalue, C.bool(caseSensitive)))
		}
	}

	return cid, qb.Err
}

// StringVectorContains is called internally
func (qb *QueryBuilder) StringVectorContains(property *BaseProperty, value string, caseSensitive bool) (ConditionId, error) {
	var cid ConditionId

	if qb.Err == nil && qb.checkEntityId(property.Entity.Id) {
		cvalue := C.CString(value)
		defer C.free(unsafe.Pointer(cvalue))
		cid = qb.getConditionId(C.obx_qb_any_equals_string(qb.cqb, C.obx_schema_id(property.Id), cvalue, C.bool(caseSensitive)))
	}

	return cid, qb.Err
}

// IntBetween is called internally
func (qb *QueryBuilder) IntBetween(property *BaseProperty, value1 int64, value2 int64) (ConditionId, error) {
	var cid ConditionId

	if qb.Err == nil && qb.checkEntityId(property.Entity.Id) {
		cid = qb.getConditionId(C.obx_qb_between_2ints(qb.cqb, C.obx_schema_id(property.Id), C.int64_t(value1), C.int64_t(value2)))
	}

	return cid, qb.Err
}

// IntEqual is called internally
func (qb *QueryBuilder) IntEqual(property *BaseProperty, value int64) (ConditionId, error) {
	var cid ConditionId

	if qb.Err == nil && qb.checkEntityId(property.Entity.Id) {
		cid = qb.getConditionId(C.obx_qb_equals_int(qb.cqb, C.obx_schema_id(property.Id), C.int64_t(value)))
	}

	return cid, qb.Err
}

// IntNotEqual is called internally
func (qb *QueryBuilder) IntNotEqual(property *BaseProperty, value int64) (ConditionId, error) {
	var cid ConditionId

	if qb.Err == nil && qb.checkEntityId(property.Entity.Id) {
		cid = qb.getConditionId(C.obx_qb_not_equals_int(qb.cqb, C.obx_schema_id(property.Id), C.int64_t(value)))
	}

	return cid, qb.Err
}

// IntGreater is called internally
func (qb *QueryBuilder) IntGreater(property *BaseProperty, value int64, withEqual bool) (ConditionId, error) {
	var cid ConditionId

	if qb.Err == nil && qb.checkEntityId(property.Entity.Id) {
		if withEqual {
			cid = qb.getConditionId(C.obx_qb_greater_or_equal_int(qb.cqb, C.obx_schema_id(property.Id), C.int64_t(value)))
		} else {
			cid = qb.getConditionId(C.obx_qb_greater_than_int(qb.cqb, C.obx_schema_id(property.Id), C.int64_t(value)))
		}
	}

	return cid, qb.Err
}

// IntLess is called internally
func (qb *QueryBuilder) IntLess(property *BaseProperty, value int64, withEqual bool) (ConditionId, error) {
	var cid ConditionId

	if qb.Err == nil && qb.checkEntityId(property.Entity.Id) {
		if withEqual {
			cid = qb.getConditionId(C.obx_qb_less_or_equal_int(qb.cqb, C.obx_schema_id(property.Id), C.int64_t(value)))
		} else {
			cid = qb.getConditionId(C.obx_qb_less_than_int(qb.cqb, C.obx_schema_id(property.Id), C.int64_t(value)))
		}
	}

	return cid, qb.Err
}

// Int64In is called internally
func (qb *QueryBuilder) Int64In(property *BaseProperty, values []int64) (ConditionId, error) {
	var cid ConditionId

	if qb.Err == nil && qb.checkEntityId(property.Entity.Id) {
		cid = qb.getConditionId(C.obx_qb_in_int64s(qb.cqb, C.obx_schema_id(property.Id), goInt64ArrayToC(values), C.size_t(len(values))))
	}

	return cid, qb.Err
}

// Int64NotIn is called internally
func (qb *QueryBuilder) Int64NotIn(property *BaseProperty, values []int64) (ConditionId, error) {
	var cid ConditionId

	if qb.Err == nil && qb.checkEntityId(property.Entity.Id) {
		cid = qb.getConditionId(C.obx_qb_not_in_int64s(qb.cqb, C.obx_schema_id(property.Id), goInt64ArrayToC(values), C.size_t(len(values))))
	}

	return cid, qb.Err
}

// Int32In is called internally
func (qb *QueryBuilder) Int32In(property *BaseProperty, values []int32) (ConditionId, error) {
	var cid ConditionId

	if qb.Err == nil && qb.checkEntityId(property.Entity.Id) {
		cid = qb.getConditionId(C.obx_qb_in_int32s(qb.cqb, C.obx_schema_id(property.Id), goInt32ArrayToC(values), C.size_t(len(values))))
	}

	return cid, qb.Err
}

// Int32NotIn is called internally
func (qb *QueryBuilder) Int32NotIn(property *BaseProperty, values []int32) (ConditionId, error) {
	var cid ConditionId

	if qb.Err == nil && qb.checkEntityId(property.Entity.Id) {
		cid = qb.getConditionId(C.obx_qb_not_in_int32s(qb.cqb, C.obx_schema_id(property.Id), goInt32ArrayToC(values), C.size_t(len(values))))
	}

	return cid, qb.Err
}

// DoubleGreater is called internally
func (qb *QueryBuilder) DoubleGreater(property *BaseProperty, value float64, withEqual bool) (ConditionId, error) {
	var cid ConditionId

	if qb.Err == nil && qb.checkEntityId(property.Entity.Id) {
		if withEqual {
			cid = qb.getConditionId(C.obx_qb_greater_or_equal_double(qb.cqb, C.obx_schema_id(property.Id), C.double(value)))
		} else {
			cid = qb.getConditionId(C.obx_qb_greater_than_double(qb.cqb, C.obx_schema_id(property.Id), C.double(value)))
		}
	}

	return cid, qb.Err
}

// DoubleLess is called internally
func (qb *QueryBuilder) DoubleLess(property *BaseProperty, value float64, withEqual bool) (ConditionId, error) {
	var cid ConditionId

	if qb.Err == nil && qb.checkEntityId(property.Entity.Id) {
		if withEqual {
			cid = qb.getConditionId(C.obx_qb_less_or_equal_double(qb.cqb, C.obx_schema_id(property.Id), C.double(value)))
		} else {
			cid = qb.getConditionId(C.obx_qb_less_than_double(qb.cqb, C.obx_schema_id(property.Id), C.double(value)))
		}
	}

	return cid, qb.Err
}

// DoubleBetween is called internally
func (qb *QueryBuilder) DoubleBetween(property *BaseProperty, valueA float64, valueB float64) (ConditionId, error) {
	var cid ConditionId

	if qb.Err == nil && qb.checkEntityId(property.Entity.Id) {
		cid = qb.getConditionId(C.obx_qb_between_2doubles(qb.cqb, C.obx_schema_id(property.Id), C.double(valueA), C.double(valueB)))
	}

	return cid, qb.Err
}

// BytesEqual is called internally
func (qb *QueryBuilder) BytesEqual(property *BaseProperty, value []byte) (ConditionId, error) {
	var cid ConditionId

	if qb.Err == nil && qb.checkEntityId(property.Entity.Id) {
		cid = qb.getConditionId(C.obx_qb_equals_bytes(qb.cqb, C.obx_schema_id(property.Id), cBytesPtr(value), C.size_t(len(value))))
	}

	return cid, qb.Err
}

// BytesGreater is called internally
func (qb *QueryBuilder) BytesGreater(property *BaseProperty, value []byte, withEqual bool) (ConditionId, error) {
	var cid ConditionId

	if qb.Err == nil && qb.checkEntityId(property.Entity.Id) {
		if withEqual {
			cid = qb.getConditionId(C.obx_qb_greater_or_equal_bytes(qb.cqb, C.obx_schema_id(property.Id), cBytesPtr(value), C.size_t(len(value))))
		} else {
			cid = qb.getConditionId(C.obx_qb_greater_than_bytes(qb.cqb, C.obx_schema_id(property.Id), cBytesPtr(value), C.size_t(len(value))))
		}
	}

	return cid, qb.Err
}

// BytesLess is called internally
func (qb *QueryBuilder) BytesLess(property *BaseProperty, value []byte, withEqual bool) (ConditionId, error) {
	var cid ConditionId

	if qb.Err == nil && qb.checkEntityId(property.Entity.Id) {
		if withEqual {
			cid = qb.getConditionId(C.obx_qb_less_or_equal_bytes(qb.cqb, C.obx_schema_id(property.Id), cBytesPtr(value), C.size_t(len(value))))
		} else {
			cid = qb.getConditionId(C.obx_qb_less_than_bytes(qb.cqb, C.obx_schema_id(property.Id), cBytesPtr(value), C.size_t(len(value))))
		}
	}

	return cid, qb.Err
}
