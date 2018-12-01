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

// Query provides a way to search stored objects
type Query struct {
	typeId    TypeId
	objectBox *ObjectBox
	condition Condition

	cQuery *C.OBX_query
}

func (query *Query) Find() (objects interface{}, err error) {
	if err = query.cBuild(); err != nil {
		return nil, err
	} else {
		defer query.cFree()
	}

	err = query.objectBox.runWithCursor(query.typeId, true, func(cursor *cursor) error {
		var errInner error
		objects, errInner = query.find(cursor)
		return errInner
	})

	return
}

func (query *Query) FindIds() (ids []uint64, err error) {
	if err = query.cBuild(); err != nil {
		return nil, err
	} else {
		defer query.cFree()
	}

	err = query.objectBox.runWithCursor(query.typeId, true, func(cursor *cursor) error {
		var errInner error
		ids, errInner = query.findIds(cursor)
		return errInner
	})

	return
}

func (query *Query) Count() (count uint64, err error) {
	if err = query.cBuild(); err != nil {
		return 0, err
	} else {
		defer query.cFree()
	}

	err = query.objectBox.runWithCursor(query.typeId, true, func(cursor *cursor) error {
		var errInner error
		count, errInner = query.count(cursor)
		return errInner
	})

	return
}

func (query *Query) Remove() (count uint64, err error) {
	if err := query.cBuild(); err != nil {
		return 0, err
	} else {
		defer query.cFree()
	}

	err = query.objectBox.runWithCursor(query.typeId, false, func(cursor *cursor) error {
		var errInner error
		count, errInner = query.remove(cursor)
		return errInner
	})

	return
}

func (query *Query) Describe() (string, error) {
	if err := query.cBuild(); err != nil {
		return "", err
	} else {
		defer query.cFree()
	}

	// no need to free, it's handled by the cQuery internally
	cResult := C.obx_query_describe_parameters(query.cQuery)

	return C.GoString(cResult), nil
}

// builds query JiT when it's needed for execution
func (query *Query) cBuild() error {
	qb := query.objectBox.newQueryBuilder(query.typeId)
	defer qb.Close()

	var err error
	if _, err = query.condition.build(qb); err != nil {
		return err
	}

	query.cQuery, err = qb.build()
	return err
}

func (query *Query) cFree() (err error) {
	if query.cQuery != nil {
		rc := C.obx_query_close(query.cQuery)
		query.cQuery = nil
		if rc != 0 {
			err = createError()
		}
	}
	return
}

func (query *Query) count(cursor *cursor) (uint64, error) {
	var cCount C.uint64_t
	rc := C.obx_query_count(query.cQuery, cursor.cursor, &cCount)
	if rc != 0 {
		return 0, createError()
	}
	return uint64(cCount), nil
}

func (query *Query) remove(cursor *cursor) (uint64, error) {
	var cCount C.uint64_t
	rc := C.obx_query_remove(query.cQuery, cursor.cursor, &cCount)
	if rc != 0 {
		return 0, createError()
	}
	return uint64(cCount), nil
}

func (query *Query) findIds(cursor *cursor) (ids []uint64, err error) {
	cIdsArray := C.obx_query_find_ids(query.cQuery, cursor.cursor)
	if cIdsArray == nil {
		return nil, createError()
	}

	idsArray := cIdsArrayToGo(cIdsArray)
	defer idsArray.free()

	return idsArray.ids, nil
}

func (query *Query) find(cursor *cursor) (slice interface{}, err error) {
	cBytesArray := C.obx_query_find(query.cQuery, cursor.cursor)
	if cBytesArray == nil {
		return nil, createError()
	}

	return cursor.cBytesArrayToObjects(cBytesArray), nil
}
