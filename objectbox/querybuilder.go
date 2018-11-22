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

// Allows fluent construction of queries; just check QueryBuilder.Err or err from Build()
type QueryBuilder struct {
	objectBox *ObjectBox
	cqb       *C.OBX_query_builder

	// Currently unused
	cLastCondition C.obx_qb_cond

	// Any error that occurred during a call to QueryBuilder or its construction
	Err error
}

func (qb *QueryBuilder) Close() (err error) {
	toClose := qb.cqb
	if toClose != nil {
		qb.cqb = nil
		rc := C.obx_qb_close(toClose)
		if rc != 0 {
			err = createError()
		}
	}
	return
}

func (qb *QueryBuilder) StringEq(propertyId TypeId, value string, caseSensitive bool) {
	if qb.Err != nil {
		return
	}
	cvalue := C.CString(value)
	defer C.free(unsafe.Pointer(cvalue))
	qb.cLastCondition = C.obx_qb_string_equal(qb.cqb, C.obx_schema_id(propertyId), cvalue, C.bool(caseSensitive))
	qb.checkForCError() // Mirror C error early to Err

	// TBD: depending on Go's query API, return either *QueryBuilder or query condition
	return
}

func (qb *QueryBuilder) IntBetween(propertyId TypeId, value1 int64, value2 int64) {
	if qb.Err != nil {
		return
	}
	qb.cLastCondition = C.obx_qb_int_between(qb.cqb, C.obx_schema_id(propertyId), C.int64_t(value1), C.int64_t(value2))
	qb.checkForCError() // Mirror C error early to Err

	// TBD: depending on Go's query API, return either *QueryBuilder or query condition
	return
}

func (qb *QueryBuilder) Build() (query *Query, err error) {
	qb.checkForCError()
	if qb.Err != nil {
		return nil, qb.Err
	}
	cquery, err := C.obx_query_create(qb.cqb)
	if err != nil {
		return nil, err
	}
	return &Query{cquery: cquery}, nil
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

func (qb *QueryBuilder) BuildAndClose() (query *Query, err error) {
	err = qb.Err
	if err == nil {
		query, err = qb.Build()
	}
	err2 := qb.Close()
	if err == nil && err2 != nil {
		query.Close()
		return nil, err2
	}
	return
}
