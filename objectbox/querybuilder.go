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

// Allows construction of queries; just check QueryBuilder.Error or err from Build()
//
// WIP: Query interface is subject to change with full ObjectBox queries support
type QueryBuilder struct {
	objectBox *ObjectBox
	cqb       *C.OBX_query_builder
	typeId    TypeId

	// Currently unused
	cLastCondition C.obx_qb_cond

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

func (qb *QueryBuilder) StringEq(propertyId TypeId, value string, caseSensitive bool) {
	if qb.Err != nil {
		return
	}
	cvalue := C.CString(value)
	defer C.free(unsafe.Pointer(cvalue))
	qb.cLastCondition = C.obx_qb_string_equal(qb.cqb, C.obx_schema_id(propertyId), cvalue, C.bool(caseSensitive))
	qb.checkForCError() // Mirror C error early to Error

	// TBD: depending on Go's query API, return either *QueryBuilder or query condition
	return
}

func (qb *QueryBuilder) IntBetween(propertyId TypeId, value1 int64, value2 int64) {
	if qb.Err != nil {
		return
	}
	qb.cLastCondition = C.obx_qb_int_between(qb.cqb, C.obx_schema_id(propertyId), C.int64_t(value1), C.int64_t(value2))
	qb.checkForCError() // Mirror C error early to Error

	// TBD: depending on Go's query API, return either *QueryBuilder or query condition
	return
}

func (qb *QueryBuilder) Null(propertyId TypeId) {
	if qb.Err != nil {
		return
	}
	qb.cLastCondition = C.obx_qb_null(qb.cqb, C.obx_schema_id(propertyId))
	qb.checkForCError() // Mirror C error early to Error

	// TBD: depending on Go's query API, return either *QueryBuilder or query condition
	return
}

func (qb *QueryBuilder) StringNotEq(propertyId TypeId, value string, caseSensitive bool) {
	if qb.Err != nil {
		return
	}
	cvalue := C.CString(value)
	defer C.free(unsafe.Pointer(cvalue))
	qb.cLastCondition = C.obx_qb_string_not_equal(qb.cqb, C.obx_schema_id(propertyId), cvalue, C.bool(caseSensitive))
	qb.checkForCError() // Mirror C error early to Error

	// TBD: depending on Go's query API, return either *QueryBuilder or query condition
	return
}

func (qb *QueryBuilder) StringGreater(propertyId TypeId, value string, caseSensitive bool, withEqual bool) {
	if qb.Err != nil {
		return
	}
	cvalue := C.CString(value)
	defer C.free(unsafe.Pointer(cvalue))
	qb.cLastCondition = C.obx_qb_string_greater(qb.cqb, C.obx_schema_id(propertyId), cvalue, C.bool(caseSensitive), C.bool(withEqual))
	qb.checkForCError() // Mirror C error early to Error

	// TBD: depending on Go's query API, return either *QueryBuilder or query condition
	return
}

func (qb *QueryBuilder) NotNull(propertyId TypeId) {
	if qb.Err != nil {
		return
	}
	qb.cLastCondition = C.obx_qb_not_null(qb.cqb, C.obx_schema_id(propertyId))
	qb.checkForCError() // Mirror C error early to Error

	// TBD: depending on Go's query API, return either *QueryBuilder or query condition
	return
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


//// **************************


//// obx_qb_cond obx_qb_string_contains(OBX_query_builder* builder, obx_schema_id property_id, string value, bool case_sensitive);
//func (qb *QueryBuilder) StringContains(propertyId TypeId, value string, case_sensitive bool) {
//	if qb.Err != nil {
//		return
//	}
//	qb.cLastCondition = obx_qb_string_contains(qb.cqb, C.obx_schema_id(propertyId))
//	qb.checkForCError() // Mirror C error early to Error
//
//	// TBD: depending on Go's query API, return either *QueryBuilder or query condition
//	return
//}
//
////  const char* value, bool case_sensitive
//// obx_qb_cond obx_qb_string_starts_with(OBX_query_builder* builder, obx_schema_id property_id, const char* value, bool case_sensitive);
//func (qb *QueryBuilder) StringStartsWith(propertyId TypeId, value string, case_sensitive bool) {
//	if qb.Err != nil {
//		return
//	}
//	qb.cLastCondition = C.obx_qb_cond
//	obx_qb_string_starts_with(qb.cqb, C.obx_schema_id(propertyId))
//	qb.checkForCError() // Mirror C error early to Error
//
//	// TBD: depending on Go's query API, return either *QueryBuilder or query condition
//	return
//}
//
////  const char* value, bool case_sensitive
//// obx_qb_cond obx_qb_string_ends_with(OBX_query_builder* builder, obx_schema_id property_id, const char* value, bool case_sensitive);
//func (qb *QueryBuilder) StringEndsWith(propertyId TypeId, value string, case_sensitive bool) {
//	if qb.Err != nil {
//		return
//	}
//	qb.cLastCondition = C.obx_qb_cond
//	obx_qb_string_ends_with(qb.cqb, C.obx_schema_id(propertyId))
//	qb.checkForCError() // Mirror C error early to Error
//
//	// TBD: depending on Go's query API, return either *QueryBuilder or query condition
//	return
//}
//
//
////  const char* value, bool case_sensitive, bool with_equal
//// obx_qb_cond obx_qb_string_less(OBX_query_builder* builder, obx_schema_id property_id, const char* value, bool case_sensitive, bool with_equal);
//func (qb *QueryBuilder) StringLess(propertyId TypeId, value string, case_sensitive bool, with_equal bool) {
//	if qb.Err != nil {
//		return
//	}
//	qb.cLastCondition = C.obx_qb_cond
//	obx_qb_string_less(qb.cqb, C.obx_schema_id(propertyId))
//	qb.checkForCError() // Mirror C error early to Error
//
//	// TBD: depending on Go's query API, return either *QueryBuilder or query condition
//	return
//}
//
////  const char* values[], int count, bool case_sensitive
//// obx_qb_cond obx_qb_string_in(OBX_query_builder* builder, obx_schema_id property_id, const char* values[], int count, bool case_sensitive);
//func (qb *QueryBuilder) StringIn(propertyId TypeId, values [] string, count int, case_sensitive bool) {
//	if qb.Err != nil {
//		return
//	}
//	qb.cLastCondition = C.obx_qb_cond
//	obx_qb_string_in(qb.cqb, C.obx_schema_id(propertyId))
//	qb.checkForCError() // Mirror C error early to Error
//
//	// TBD: depending on Go's query API, return either *QueryBuilder or query condition
//	return
//}
//
//// obx_qb_cond obx_qb_int_equal(OBX_query_builder* builder, obx_schema_id property_id, int64_t value);
//func (qb *QueryBuilder) IntEqual(propertyId TypeId, value int64) {
//	if qb.Err != nil {
//		return
//	}
//	qb.cLastCondition = C.obx_qb_cond
//	obx_qb_int_equal(qb.cqb, C.obx_schema_id(propertyId), value)
//	qb.checkForCError() // Mirror C error early to Error
//
//	// TBD: depending on Go's query API, return either *QueryBuilder or query condition
//	return
//}
//
//// obx_qb_cond obx_qb_int_not_equal(OBX_query_builder* builder, obx_schema_id property_id, int64 value);
//func (qb *QueryBuilder) IntNotEqual(propertyId TypeId, value int64) {
//	if qb.Err != nil {
//		return
//	}
//	qb.cLastCondition = C.obx_qb_cond
//	obx_qb_int_not_equal(qb.cqb, C.obx_schema_id(propertyId))
//	qb.checkForCError() // Mirror C error early to Error
//
//	// TBD: depending on Go's query API, return either *QueryBuilder or query condition
//	return
//}
//
////  int64_t value
//// obx_qb_cond obx_qb_int_greater(OBX_query_builder* builder, obx_schema_id property_id, int64  value);
//func (qb *QueryBuilder) IntGreater(propertyId TypeId, value int64) {
//	if qb.Err != nil {
//		return
//	}
//	qb.cLastCondition = C.obx_qb_cond
//	obx_qb_int_greater(qb.cqb, C.obx_schema_id(propertyId))
//	qb.checkForCError() // Mirror C error early to Error
//
//	// TBD: depending on Go's query API, return either *QueryBuilder or query condition
//	return
//}
//
//// obx_qb_cond obx_qb_int_less(OBX_query_builder* builder, obx_schema_id property_id, int64  value);
//func (qb *QueryBuilder) IntLess(propertyId TypeId, value int64) {
//	if qb.Err != nil {
//		return
//	}
//	qb.cLastCondition = C.obx_qb_cond
//	obx_qb_int_less(qb.cqb, C.obx_schema_id(propertyId))
//	qb.checkForCError() // Mirror C error early to Error
//
//	// TBD: depending on Go's query API, return either *QueryBuilder or query condition
//	return
//}
//
//// obx_qb_cond obx_qb_int64_in(OBX_query_builder* builder, obx_schema_id property_id, const int64  values[], int count);
//func (qb *QueryBuilder) Int64In(propertyId TypeId, values [] int64, count int) {
//	if qb.Err != nil {
//		return
//	}
//	qb.cLastCondition = C.obx_qb_cond
//	obx_qb_int64_in(qb.cqb, C.obx_schema_id(propertyId))
//	qb.checkForCError() // Mirror C error early to Error
//
//	// TBD: depending on Go's query API, return either *QueryBuilder or query condition
//	return
//}
//
////  const int64_t values[], int count
//// obx_qb_cond obx_qb_int64_not_in(OBX_query_builder* builder, obx_schema_id property_id, const int64  values[], int count);
//func (qb *QueryBuilder) Int64NotIn(propertyId TypeId, values [] int64, count int) {
//	if qb.Err != nil {
//		return
//	}
//	qb.cLastCondition = C.obx_qb_cond
//	obx_qb_int64_not_in(qb.cqb, C.obx_schema_id(propertyId))
//	qb.checkForCError() // Mirror C error early to Error
//
//	// TBD: depending on Go's query API, return either *QueryBuilder or query condition
//	return
//}
//
//// obx_qb_cond obx_qb_int32_in(OBX_query_builder* builder, obx_schema_id property_id, const int32_t values[], int count);
//func (qb *QueryBuilder) Int32In(propertyId TypeId, values [] int32_t, count int) {
//	if qb.Err != nil {
//		return
//	}
//	qb.cLastCondition = C.obx_qb_cond
//	obx_qb_int32_in(qb.cqb, C.obx_schema_id(propertyId), values, count)
//	qb.checkForCError() // Mirror C error early to Error
//
//	// TBD: depending on Go's query API, return either *QueryBuilder or query condition
//	return
//}
//
////  const int32_t values[], int count
//// obx_qb_cond obx_qb_int32_not_in(OBX_query_builder* builder, obx_schema_id property_id, const int32_t values[], int count);
//func (qb *QueryBuilder) Int32NotIn(propertyId TypeId, values [] int32_t, count int) {
//	if qb.Err != nil {
//		return
//	}
//	qb.cLastCondition = C.obx_qb_cond
//	obx_qb_int32_not_in(qb.cqb, C.obx_schema_id(propertyId))
//	qb.checkForCError() // Mirror C error early to Error
//
//	// TBD: depending on Go's query API, return either *QueryBuilder or query condition
//	return
//}
//
////  double value
//// obx_qb_cond obx_qb_double_greater(OBX_query_builder* builder, obx_schema_id property_id, double value);
//func (qb *QueryBuilder) DoubleGreater(propertyId TypeId, value double) {
//	if qb.Err != nil {
//		return
//	}
//	qb.cLastCondition = C.obx_qb_cond
//	obx_qb_double_greater(qb.cqb, C.obx_schema_id(propertyId))
//	qb.checkForCError() // Mirror C error early to Error
//
//	// TBD: depending on Go's query API, return either *QueryBuilder or query condition
//	return
//}
//
////  double value
//// obx_qb_cond obx_qb_double_less(OBX_query_builder* builder, obx_schema_id property_id, double value);
//func (qb *QueryBuilder) DoubleLess(propertyId TypeId, value double) {
//	if qb.Err != nil {
//		return
//	}
//	qb.cLastCondition = C.obx_qb_cond
//	obx_qb_double_less(qb.cqb, C.obx_schema_id(propertyId))
//	qb.checkForCError() // Mirror C error early to Error
//
//	// TBD: depending on Go's query API, return either *QueryBuilder or query condition
//	return
//}
//
////  double value_a, double value_b
//// obx_qb_cond obx_qb_double_between(OBX_query_builder* builder, obx_schema_id property_id, double value_a, double value_b);
//func (qb *QueryBuilder) DoubleBetween(propertyId TypeId, value_a double, value_b double) {
//	if qb.Err != nil {
//		return
//	}
//	qb.cLastCondition = C.obx_qb_cond
//	obx_qb_double_between(qb.cqb, C.obx_schema_id(propertyId))
//	qb.checkForCError() // Mirror C error early to Error
//
//	// TBD: depending on Go's query API, return either *QueryBuilder or query condition
//	return
//}
//
//// obx_qb_cond obx_qb_bytes_equal(OBX_query_builder* builder, obx_schema_id property_id, const void* value, size_t size);
//func (qb *QueryBuilder) BytesEqual(propertyId TypeId, value []byte, size size_t) {
//	if qb.Err != nil {
//		return
//	}
//	qb.cLastCondition = C.obx_qb_cond
//	obx_qb_bytes_equal(qb.cqb, C.obx_schema_id(propertyId))
//	qb.checkForCError() // Mirror C error early to Error
//
//	// TBD: depending on Go's query API, return either *QueryBuilder or query condition
//	return
//}
//
////  const void* value, size_t size, bool with_equal
//// obx_qb_cond obx_qb_bytes_greater(OBX_query_builder* builder, obx_schema_id property_id, const void* value, size_t size, bool with_equal);
//func (qb *QueryBuilder) BytesGreater(propertyId TypeId, value []byte, size size_t, with_equal bool) {
//	if qb.Err != nil {
//		return
//	}
//	qb.cLastCondition = C.obx_qb_cond
//	obx_qb_bytes_greater(qb.cqb, C.obx_schema_id(propertyId))
//	qb.checkForCError() // Mirror C error early to Error
//
//	// TBD: depending on Go's query API, return either *QueryBuilder or query condition
//	return
//}
//
////  const void* value, size_t size, bool with_equal
//// obx_qb_cond obx_qb_bytes_less(OBX_query_builder* builder, obx_schema_id property_id, const void* value, size_t size, bool with_equal);
//func (qb *QueryBuilder) BytesLess(propertyId TypeId, value []byte, size size_t, with_equal bool) {
//	if qb.Err != nil {
//		return
//	}
//	qb.cLastCondition = C.obx_qb_cond
//	obx_qb_bytes_less(qb.cqb, C.obx_schema_id(propertyId))
//	qb.checkForCError() // Mirror C error early to Error
//
//	// TBD: depending on Go's query API, return either *QueryBuilder or query condition
//	return
//}
//
//////
//
//
//
