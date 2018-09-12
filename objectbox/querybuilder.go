package objectbox

/*
#cgo LDFLAGS: -lobjectboxc
#include "objectbox.h"
*/
import "C"

type QueryBuilder struct {
	objectBox *ObjectBox
	cqb       *C.OBX_query_builder
}

func (qb *QueryBuilder) Destroy() (err error) {
	if qb.cqb != nil {
		rc := C.obx_qb_close(qb.cqb)
		qb.cqb = nil
		if rc != 0 {
			err = createError()
		}
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

func (qb *QueryBuilder) BuildAndDestroy() (query *Query, err error) {
	query, err = qb.Build()
	err2 := qb.Destroy()
	if err == nil && err2 != nil {
		query.Destroy()
		return nil, err2
	}
	return
}
