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
	cquery    *C.OBX_query
	typeId    TypeId
	objectBox *ObjectBox
	condition Condition
}

func (query *Query) Close() (err error) {
	if query.cquery != nil {
		rc := C.obx_query_close(query.cquery)
		query.cquery = nil
		if rc != 0 {
			err = createError()
		}
	}
	return
}

// builds query JiT when it's needed for execution
func (query *Query) build() error {
	qb := query.objectBox.QueryBuilder(query.typeId)
	defer qb.Close()

	var err error
	if _, err = query.condition.build(qb); err != nil {
		return err
	}

	query.cquery, err = qb.build()
	return err
}

func (query *Query) Find() (objects interface{}, err error) {
	defer query.Close()
	if err = query.build(); err != nil {
		return nil, err
	}

	err = query.objectBox.runWithCursor(query.typeId, true, func(cursor *cursor) error {
		var errInner error
		objects, errInner = query.find(cursor)
		return errInner
	})

	return
}

func (query *Query) find(cursor *cursor) (slice interface{}, err error) {
	bytesArray, err := query.findBytes(cursor)
	if err != nil {
		return
	}
	defer bytesArray.free()
	return cursor.bytesArrayToObjects(bytesArray), nil
}

// Deprecated: Won't be public in the future
func (query *Query) FindBytes() (bytesArray *BytesArray, err error) {
	defer query.Close()
	if err = query.build(); err != nil {
		return nil, err
	}

	err = query.objectBox.runWithCursor(query.typeId, true, func(cursor *cursor) error {
		var errInner error
		bytesArray, errInner = query.findBytes(cursor)
		return errInner
	})
	return
}

func (query *Query) findBytes(cursor *cursor) (*BytesArray, error) {
	cBytesArray := C.obx_query_find(query.cquery, cursor.cursor)
	if cBytesArray == nil {
		return nil, createError()
	}
	return cBytesArrayToGo(cBytesArray), nil
}
