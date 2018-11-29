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

// WIP: Query interface is subject to change with full ObjectBox queries support
type Query struct {
	typeId    TypeId
	objectBox *ObjectBox
	condition Condition
}

// builds query JiT when it's needed for execution
func (query *Query) build() (*cQuery, error) {
	qb := query.objectBox.newQueryBuilder(query.typeId)
	defer qb.Close()

	var err error
	if _, err = query.condition.build(qb); err != nil {
		return nil, err
	}

	if cquery, err := qb.build(); err != nil {
		return nil, err
	} else {
		return &cQuery{cquery}, nil
	}
}

func (query *Query) Find() (objects interface{}, err error) {
	cQuery, err := query.build()
	if err != nil {
		return nil, err
	} else {
		defer cQuery.close()
	}

	err = query.objectBox.runWithCursor(query.typeId, true, func(cursor *cursor) error {
		var errInner error
		objects, errInner = cQuery.find(cursor)
		return errInner
	})

	return
}

// WIP: Query interface is subject to change with full ObjectBox queries support
type cQuery struct {
	cquery *C.OBX_query
}

func (query *cQuery) close() (err error) {
	if query.cquery != nil {
		rc := C.obx_query_close(query.cquery)
		query.cquery = nil
		if rc != 0 {
			err = createError()
		}
	}
	return
}

func (query *cQuery) find(cursor *cursor) (slice interface{}, err error) {
	bytesArray, err := query.findBytes(cursor)
	if err != nil {
		return
	}
	defer bytesArray.free()
	return cursor.bytesArrayToObjects(bytesArray), nil
}

func (query *cQuery) findBytes(cursor *cursor) (*BytesArray, error) {
	cBytesArray := C.obx_query_find(query.cquery, cursor.cursor)
	if cBytesArray == nil {
		return nil, createError()
	}
	return cBytesArrayToGo(cBytesArray), nil
}
