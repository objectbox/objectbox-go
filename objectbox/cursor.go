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
	"reflect"
	"unsafe"

	"github.com/google/flatbuffers/go"
)

// Internal: won't be publicly exposed in a future version!
type Cursor struct {
	txn    *Transaction
	cursor *C.OBX_cursor
	entity *entity
	fbb    *flatbuffers.Builder
}

const defaultSliceCapacity = 16

func (cursor *Cursor) Close() error {
	rc := C.obx_cursor_close(cursor.cursor)
	cursor.cursor = nil
	if rc != 0 {
		return createError()
	}
	return nil
}

func (cursor *Cursor) Get(id uint64) (object interface{}, err error) {
	var data *C.void
	var dataSize C.size_t
	var dataPtr = unsafe.Pointer(data)

	var rc = C.obx_cursor_get(cursor.cursor, C.obx_id(id), &dataPtr, &dataSize)
	if rc == 0 {
		var bytes []byte
		cVoidPtrToByteSlice(dataPtr, int(dataSize), &bytes)
		return cursor.entity.binding.Load(cursor.txn, bytes)
	} else if rc == C.OBX_NOT_FOUND {
		return nil, nil
	} else {
		return nil, createError()
	}
}

func (cursor *Cursor) GetAll() (slice interface{}, err error) {
	if supportsBytesArray {
		cBytesArray := C.obx_cursor_get_all(cursor.cursor)
		if cBytesArray == nil {
			return nil, createError()
		}
		return cursor.cBytesArrayToObjects(cBytesArray)
	} else {
		return cursor.getAllSequential()
	}
}

func (cursor *Cursor) getAllSequential() (slice interface{}, err error) {
	var binding = cursor.entity.binding
	slice = binding.MakeSlice(defaultSliceCapacity)

	var data *C.void
	var dataSize C.size_t
	var dataPtr = unsafe.Pointer(data)

	var bytes []byte
	var rc C.obx_err
	for rc = C.obx_cursor_first(cursor.cursor, &dataPtr, &dataSize); rc == 0; rc = C.obx_cursor_next(cursor.cursor, &dataPtr, &dataSize) {
		cVoidPtrToByteSlice(dataPtr, int(dataSize), &bytes)
		if object, err := binding.Load(cursor.txn, bytes); err != nil {
			return nil, err
		} else {
			slice = binding.AppendToSlice(slice, object)
		}
	}

	// if there was an error
	if rc != 0 && rc != C.OBX_NOT_FOUND {
		return nil, createError()
	}

	return slice, nil
}

func (cursor *Cursor) Count() (count uint64, err error) {
	var cCount C.uint64_t
	rc := C.obx_cursor_count(cursor.cursor, &cCount)
	if rc != 0 {
		err = createError()
		return
	}
	return uint64(cCount), nil
}

func (cursor *Cursor) CountMax(max uint64) (count uint64, err error) {
	var cCount C.uint64_t
	rc := C.obx_cursor_count_max(cursor.cursor, C.uint64_t(max), &cCount)
	if rc != 0 {
		err = createError()
		return
	}
	return uint64(cCount), nil
}

func (cursor *Cursor) IsEmpty() (result bool, err error) {
	var cResult C.bool
	rc := C.obx_cursor_is_empty(cursor.cursor, &cResult)
	if rc != 0 {
		err = createError()
		return
	}
	return bool(cResult), nil
}

// seek moves Cursor to the object with the given ID (if it exists)
func (cursor *Cursor) seek(id uint64) (bool, error) {
	rc := C.obx_cursor_seek(cursor.cursor, C.obx_id(id))
	if rc == 0 {
		return true, nil
	} else if rc == C.OBX_NOT_FOUND {
		return false, nil
	} else {
		return false, createError()
	}
}

func (cursor *Cursor) Put(object interface{}) (id uint64, err error) {
	var binding = cursor.entity.binding
	idFromObject, err := binding.GetId(object)
	if err != nil {
		return 0, err
	}

	if id, err = cursor.IdForPut(idFromObject); err != nil {
		return 0, err
	}

	if cursor.entity.hasRelations {
		if err = binding.PutRelated(cursor.txn, object, id); err != nil {
			return 0, err
		}
	}

	if err = binding.Flatten(object, cursor.fbb, id); err != nil {
		return 0, err
	}

	checkForPreviousValue := idFromObject != 0
	if err = cursor.finishInternalFbbAndPut(id, checkForPreviousValue); err != nil {
		return 0, err
	}

	// update the id on the object
	if idFromObject != id {
		binding.SetId(object, id)
	}

	return id, nil
}

func (cursor *Cursor) finishInternalFbbAndPut(id uint64, checkForPreviousObject bool) (err error) {
	fbb := cursor.fbb
	fbb.Finish(fbb.EndObject())
	bytes := fbb.FinishedBytes()

	rc := C.obx_cursor_put(cursor.cursor,
		C.obx_id(id), unsafe.Pointer(&bytes[0]), C.size_t(len(bytes)), C.bool(checkForPreviousObject))
	if rc != 0 {
		err = createError()
	}

	// Reset to have a clear state for the next caller
	fbb.Reset()

	return
}

func (cursor *Cursor) IdForPut(idCandidate uint64) (id uint64, err error) {
	id = uint64(C.obx_cursor_id_for_put(cursor.cursor, C.obx_id(idCandidate)))
	if id == 0 {
		err = createError()
	}
	return
}

func (cursor *Cursor) Remove(id uint64) (err error) {
	rc := C.obx_cursor_remove(cursor.cursor, C.obx_id(id))
	if rc != 0 {
		err = createError()
	}
	return
}

func (cursor *Cursor) RemoveAll() (err error) {
	rc := C.obx_cursor_remove_all(cursor.cursor)
	if rc != 0 {
		err = createError()
	}
	return
}

func (cursor *Cursor) cBytesArrayToObjects(cBytesArray *C.OBX_bytes_array) (slice interface{}, err error) {
	bytesArray := cBytesArrayToGo(cBytesArray)
	defer bytesArray.free()
	return cursor.bytesArrayToObjects(bytesArray)
}

func (cursor *Cursor) bytesArrayToObjects(bytesArray *bytesArray) (slice interface{}, err error) {
	var binding = cursor.entity.binding
	slice = binding.MakeSlice(len(bytesArray.BytesArray))
	for _, bytesData := range bytesArray.BytesArray {
		if object, err := binding.Load(cursor.txn, bytesData); err != nil {
			return nil, err
		} else {
			slice = binding.AppendToSlice(slice, object)
		}
	}
	return slice, nil
}

// RelationReplace replaces all targets for a given source in a standalone many-to-many relation
// It also inserts new related objects (with a 0 ID).
// TODO don't require a targetEntityId, it can retrieved using a relationId
func (cursor *Cursor) RelationReplace(relationId TypeId, targetEntityId TypeId, sourceId uint64, sourceObject interface{}, targetObjects interface{}) (err error) {
	// get id from the object, if inserting, it would be 0 even if the argument id is already non-zero
	// this saves us an unnecessary request to RelationIds for new objects (there can't be any relations yet)
	objId, err := cursor.entity.binding.GetId(sourceObject)
	if err != nil {
		return err
	}

	// make a map of related target entity IDs, marking those that were originally related but should be removed
	var idsToRemove = make(map[uint64]bool)

	if objId != 0 {
		if oldRelIds, err := cursor.RelationIds(relationId, sourceId); err != nil {
			return err
		} else {
			for _, rId := range oldRelIds {
				idsToRemove[rId] = true
			}
		}
	}

	sliceValue := reflect.ValueOf(targetObjects)
	count := sliceValue.Len()

	if count > 0 {
		targetEntity := cursor.txn.objectBox.getEntityById(targetEntityId)

		if err := cursor.txn.runWithCursor(targetEntity, func(targetCursor *Cursor) error {
			// walk over the current related objects, mark those that still exist, add the new ones
			for i := 0; i < count; i++ {
				var reflObj = sliceValue.Index(i)
				var rel interface{}
				if reflObj.Kind() == reflect.Ptr {
					rel = reflObj.Interface()
				} else {
					rel = reflObj.Addr().Interface()
				}

				rId, err := targetEntity.binding.GetId(rel)
				if err != nil {
					return err
				} else if rId == 0 {
					if rId, err = targetCursor.Put(rel); err != nil {
						return err
					}
				}

				if idsToRemove[rId] {
					// old relation that still exists, keep it
					delete(idsToRemove, rId)
				} else {
					// new relation, add it
					if err := cursor.RelationPut(relationId, sourceId, rId); err != nil {
						return err
					}
				}
			}
			return nil
		}); err != nil {
			return err
		}
	}

	// remove those that were not found in the rSlice but were originally related to this entity
	for rId := range idsToRemove {
		if err := cursor.RelationRemove(relationId, sourceId, rId); err != nil {
			return err
		}
	}

	return nil
}

// Get all target objects from a standalone relation
// TODO don't require a targetEntityId, it can retrieved using a relationId
func (cursor *Cursor) RelationGetAll(relationId TypeId, targetEntityId TypeId, sourceId uint64) (slice interface{}, err error) {
	targetIds, err := cursor.RelationIds(relationId, sourceId)
	if err != nil {
		return nil, err
	}

	targetEntity := cursor.txn.objectBox.getEntityById(targetEntityId)
	slice = targetEntity.binding.MakeSlice(len(targetIds))

	err = cursor.txn.runWithCursor(targetEntity, func(targetCursor *Cursor) error {
		for _, id := range targetIds {
			if object, err := targetCursor.Get(id); err != nil {
				return err
			} else {
				slice = targetEntity.binding.AppendToSlice(slice, object)
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	} else {
		return slice, nil
	}
}

func (cursor *Cursor) RelationPut(relationId TypeId, sourceId, targetId uint64) error {
	rc := C.obx_cursor_rel_put(cursor.cursor, C.obx_schema_id(relationId), C.obx_id(sourceId), C.obx_id(targetId))
	if rc != 0 {
		return createError()
	}
	return nil
}

func (cursor *Cursor) RelationRemove(relationId TypeId, sourceId, targetId uint64) error {
	rc := C.obx_cursor_rel_remove(cursor.cursor, C.obx_schema_id(relationId), C.obx_id(sourceId), C.obx_id(targetId))
	if rc != 0 {
		return createError()
	}
	return nil
}

func (cursor *Cursor) RelationIds(relationId TypeId, sourceId uint64) ([]uint64, error) {
	cIdsArray := C.obx_cursor_rel_ids(cursor.cursor, C.obx_schema_id(relationId), C.obx_id(sourceId))
	if cIdsArray == nil {
		return nil, createError()
	}

	idsArray := cIdsArrayToGo(cIdsArray)
	defer idsArray.free()

	return idsArray.ids, nil
}
