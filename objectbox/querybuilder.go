package objectbox

/*
#cgo LDFLAGS: -lobjectbox
#include <stdlib.h>
#include "objectbox.h"
*/
import "C"
import "unsafe"

type QueryBuilder struct {
	objectBox *ObjectBox
	cqb       *C.OBX_query_builder
}

func (qb *QueryBuilder) Close() (err error) {
	if qb.cqb != nil {
		rc := C.obx_qb_close(qb.cqb)
		qb.cqb = nil
		if rc != 0 {
			err = createError()
		}
	}
	return
}

func (qb *QueryBuilder) StringEq(propertyId TypeId, value string, caseSensitive bool) (err error) {
	cvalue := C.CString(value)
	defer C.free(unsafe.Pointer(cvalue))
	rc := C.obx_qb_string_equal(qb.cqb, C.uint32_t(propertyId), cvalue, C.bool(caseSensitive))
	if rc != 0 {
		err = createError()
	}
	return
}

func (qb *QueryBuilder) IntBetween(propertyId TypeId, value1 int64, value2 int64) (err error) {
	rc := C.obx_qb_int_between(qb.cqb, C.uint32_t(propertyId), C.int64_t(value1), C.int64_t(value2))
	if rc != 0 {
		err = createError()
	}
	return
}

func (qb *QueryBuilder) Build() (query *Query, err error) {
	cquery, err := C.obx_query_create(qb.cqb)
	if err != nil {
		return nil, err
	}
	return &Query{cquery: cquery}, nil
}

func (qb *QueryBuilder) BuildAndClose() (query *Query, err error) {
	query, err = qb.Build()
	err2 := qb.Close()
	if err == nil && err2 != nil {
		query.Close()
		return nil, err2
	}
	return
}
