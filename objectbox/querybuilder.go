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
#cgo LDFLAGS: -lobjectbox
#include <stdlib.h>
#include "objectbox.h"
*/
import "C"
import (
	"errors"
	"fmt"
	"strings"
	"unsafe"
)

// Internal class; use Box.Query instead.
// Allows construction of queries; just check queryBuilder.Error or err from Build()
type QueryBuilder struct {
	objectBox     *ObjectBox
	cqb           *C.OBX_query_builder
	typeId        TypeId
	innerBuilders []*QueryBuilder

	// The first error that occurred during a any of the calls on the query builder
	Err error
}

func newQueryBuilder(ob *ObjectBox, typeId TypeId) *QueryBuilder {
	cqb := C.obx_qb_create(ob.store, C.obx_schema_id(typeId))
	var err error = nil
	if cqb == nil {
		err = createError()
	}

	var qb = &QueryBuilder{
		objectBox: ob,
		cqb:       cqb,
		typeId:    typeId,
		Err:       err,
	}

	// log.Printf("QB %p created for entity %d\n", qb, typeId)
	return qb
}

func (qb *QueryBuilder) newInnerBuilder(typeId TypeId, cqb *C.OBX_query_builder) *QueryBuilder {
	if cqb == nil {
		qb.Err = createError()
		return nil
	}

	var iqb = &QueryBuilder{
		objectBox: qb.objectBox,
		cqb:       cqb,
		typeId:    typeId,
	}

	qb.innerBuilders = append(qb.innerBuilders, iqb)

	// log.Printf("QB %p attached inner QB %p, entity %d\n", qb, iqb, iqb.typeId)
	return iqb
}

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
		rc := C.obx_qb_close(toClose)
		if rc != 0 {
			errs = append(errs, createError().Error())
		}
	}

	if errs != nil {
		return errors.New(strings.Join(errs, "; "))
	}
	return nil
}

func (qb *QueryBuilder) Build(box *Box) (*Query, error) {
	qb.checkForCError() // TODO why is this called here? It could lead to incorrect error messages in a parallel app
	if qb.Err != nil {
		return nil, qb.Err
	}

	// log.Printf("Building %p\n", qb)

	cQuery := C.obx_query_create(qb.cqb)
	if cQuery == nil {
		qb.Err = createError()
		return nil, qb.Err
	}

	query := &Query{
		objectBox: qb.objectBox,
		cQuery:    cQuery,
		box:       box,
		entity:    qb.objectBox.getEntityById(qb.typeId),
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

func (qb *QueryBuilder) LinkOneToMany(relation *RelationToOne, conditions []Condition) error {
	if qb.Err != nil {
		return qb.Err
	}

	// create a new "inner" query builder
	var iqb *QueryBuilder

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

func (qb *QueryBuilder) LinkManyToMany(relation *RelationToMany, conditions []Condition) error {
	if qb.Err != nil {
		return qb.Err
	}

	// create a new "inner" query builder
	var iqb *QueryBuilder

	// recognize whether it's a link or a backlink
	if relation.Source.Id == qb.typeId && relation.Target.Id != qb.typeId {
		//log.Printf("QB %p creating link to entity %d over relation %d", qb, relation.Target.Id, relation.Id)
		iqb = qb.newInnerBuilder(relation.Target.Id, C.obx_qb_link_standalone(qb.cqb, C.obx_schema_id(relation.Id)))
	} else if relation.Source.Id != qb.typeId && relation.Target.Id == qb.typeId {
		//log.Printf("QB %p creating backlink from entity %d over relation %d", qb, relation.Source.Id, relation.Id)
		iqb = qb.newInnerBuilder(relation.Source.Id, C.obx_qb_backlink_standalone(qb.cqb, C.obx_schema_id(relation.Id)))
	} else {
		return errors.New("relation not recognized as either link or backlink")
	}
	if iqb == nil {
		return qb.Err // this has been set by newInnerBuilder()
	}

	return iqb.applyConditions(conditions)
}

func (qb *QueryBuilder) checkForCError() {
	if qb.Err != nil { // TODO why if err != nil, doesn't make sense at a first glance
		errCode := C.obx_qb_error_code(qb.cqb)
		if errCode != 0 {
			msg := C.obx_qb_error_message(qb.cqb)
			if msg == nil {
				qb.Err = fmt.Errorf("could not create query builder (code %v)", int(errCode))
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

func (qb *QueryBuilder) Any(ids []ConditionId) (ConditionId, error) {
	var cid ConditionId

	if qb.Err == nil {
		cid = qb.getConditionId(C.obx_qb_any(qb.cqb, (*C.obx_qb_cond)(unsafe.Pointer(&ids[0])), C.int(len(ids))))
	}

	return cid, qb.Err
}

func (qb *QueryBuilder) All(ids []ConditionId) (ConditionId, error) {
	var cid ConditionId

	if qb.Err == nil {
		cid = qb.getConditionId(C.obx_qb_all(qb.cqb, (*C.obx_qb_cond)(unsafe.Pointer(&ids[0])), C.int(len(ids))))
	}

	return cid, qb.Err
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
		if len(values) > 0 {
			cStringArray := goStringArrayToC(values)
			defer cStringArray.free()
			cid = qb.getConditionId(C.obx_qb_string_in(qb.cqb, C.obx_schema_id(property.Id), cStringArray.cArray, C.int(cStringArray.size), C.bool(caseSensitive)))
		} else {
			cid = qb.getConditionId(C.obx_qb_string_in(qb.cqb, C.obx_schema_id(property.Id), nil, 0, C.bool(caseSensitive)))
		}
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

func (qb *QueryBuilder) StringVectorContains(property *BaseProperty, value string, caseSensitive bool) (ConditionId, error) {
	var cid ConditionId

	if qb.Err == nil && qb.checkEntityId(property.Entity.Id) {
		cvalue := C.CString(value)
		defer C.free(unsafe.Pointer(cvalue))
		cid = qb.getConditionId(C.obx_qb_strings_contain(qb.cqb, C.obx_schema_id(property.Id), cvalue, C.bool(caseSensitive)))
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
		cid = qb.getConditionId(C.obx_qb_int64_in(qb.cqb, C.obx_schema_id(property.Id), goInt64ArrayToC(values), C.int(len(values))))
	}

	return cid, qb.Err
}

func (qb *QueryBuilder) Int64NotIn(property *BaseProperty, values []int64) (ConditionId, error) {
	var cid ConditionId

	if qb.Err == nil && qb.checkEntityId(property.Entity.Id) {
		cid = qb.getConditionId(C.obx_qb_int64_not_in(qb.cqb, C.obx_schema_id(property.Id), goInt64ArrayToC(values), C.int(len(values))))
	}

	return cid, qb.Err
}

func (qb *QueryBuilder) Int32In(property *BaseProperty, values []int32) (ConditionId, error) {
	var cid ConditionId

	if qb.Err == nil && qb.checkEntityId(property.Entity.Id) {
		cid = qb.getConditionId(C.obx_qb_int32_in(qb.cqb, C.obx_schema_id(property.Id), goInt32ArrayToC(values), C.int(len(values))))
	}

	return cid, qb.Err
}

func (qb *QueryBuilder) Int32NotIn(property *BaseProperty, values []int32) (ConditionId, error) {
	var cid ConditionId

	if qb.Err == nil && qb.checkEntityId(property.Entity.Id) {
		cid = qb.getConditionId(C.obx_qb_int32_not_in(qb.cqb, C.obx_schema_id(property.Id), goInt32ArrayToC(values), C.int(len(values))))
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
