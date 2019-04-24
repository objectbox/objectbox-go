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
#include "datavisitor.h"
*/
import "C"

import (
	"errors"
	"fmt"
	"reflect"
	"unsafe"

	"github.com/google/flatbuffers/go"
)

// Box provides CRUD access to objects of a common type
type Box struct {
	objectBox *ObjectBox
	box       *C.OBX_box
	entity    *entity
}

const defaultSliceCapacity = 16

// close fully closes the Box connection and free's resources
func (box *Box) close() error {
	return cMaybeErr(func() C.obx_err { return C.obx_box_close(box.box) })
}

// Creates a query with the given conditions. Use generated properties to create conditions.
// Keep the Query object if you intend to execute it multiple times.
// Note: this function panics if you try to create illegal queries; e.g. use properties of an alien type.
// This is typically a programming error. Use QueryOrError instead if you want the explicit error check.
func (box *Box) Query(conditions ...Condition) *Query {
	query, err := box.QueryOrError(conditions...)
	if err != nil {
		panic(fmt.Sprintf("Could not create query - please check your query conditions: %s", err))
	}
	return query
}

// Like Query() but with error handling; e.g. when you build conditions dynamically that may fail.
func (box *Box) QueryOrError(conditions ...Condition) (query *Query, err error) {
	builder := newQueryBuilder(box.objectBox, box.entity.id)

	defer func() {
		err2 := builder.Close()
		if err == nil && err2 != nil {
			err = err2
			query = nil
		}
	}()

	if err = builder.applyConditions(conditions); err != nil {
		return nil, err
	}

	query, err = builder.Build(box)

	return // NOTE result might be overwritten by the deferred "closer" function
}

func (box *Box) idForPut(idCandidate uint64) (id uint64, err error) {
	id = uint64(C.obx_box_id_for_put(box.box, C.obx_id(idCandidate)))
	if id == 0 {
		err = createError()
	}
	return
}

func (box *Box) idsForPut(count int) ([]uint64, error) {
	if count == 0 {
		return nil, nil
	}

	ids, err := cGetIds(func() *C.OBX_id_array { return C.obx_box_ids_for_put(box.box, C.uint64_t(count)) })
	if err != nil {
		return nil, err
	} else if len(ids) != count {
		return nil, fmt.Errorf("invalid number of new IDs reserved: got %v instead of required %v", len(ids), count)
	} else {
		return ids, nil
	}
}

func (box *Box) put(object interface{}, async bool, timeoutMs uint) (id uint64, err error) {

	idFromObject, err := box.entity.binding.GetId(object)
	if err != nil {
		return 0, err
	}

	if async && box.entity.hasRelations {
		return 0, errors.New("PutAsync is currently not supported on entities that have relations")
	}

	// TODO in case of an error later during insert, this ID is not recovered (reused in the future)...
	//  we need to either move idForPut to C-api (preferrable) or use a transaction
	if id, err = box.idForPut(idFromObject); err != nil {
		return 0, err
	}

	//log.Printf("Put %v: %v, new=%v", box.entity.name, id, idFromObject==0)

	if box.entity.hasRelations {
		if err := box.entity.binding.PutRelated(box.objectBox, object, id); err != nil {
			return 0, err
		}
	}

	err = box.withObjectBytes(object, id, func(bytes []byte) error {
		return cMaybeErr(func() C.obx_err {
			if async {
				return C.obx_box_put_async(box.box, C.obx_id(id), unsafe.Pointer(&bytes[0]), C.size_t(len(bytes)),
					C.bool(idFromObject != 0), C.uint64_t(timeoutMs))
			} else {
				return C.obx_box_put(box.box, C.obx_id(id), unsafe.Pointer(&bytes[0]), C.size_t(len(bytes)),
					C.bool(idFromObject != 0))
			}
		})
	})

	if err != nil {
		return 0, err
	}

	// update the id on the object
	if idFromObject != id {
		box.entity.binding.SetId(object, id)
	}

	return id, nil
}

func (box *Box) withObjectBytes(object interface{}, id uint64, fn func([]byte) error) error {
	var fbb = fbbPool.Get().(*flatbuffers.Builder)
	fbb.Reset()

	if err := box.entity.binding.Flatten(object, fbb, id); err != nil {
		// put the fbb back to the pool for the others to use; don't use defer, it's slower
		fbbPool.Put(fbb)
		return err
	}

	fbb.Finish(fbb.EndObject())

	var result = fn(fbb.FinishedBytes())

	// put the fbb back to the pool for the others to use; don't use defer, it's slower
	fbbPool.Put(fbb)

	return result
}

// PutAsync asynchronously inserts/updates a single object.
// When inserting, the ID property on the passed object will be assigned the new ID as well.
//
// It's executed on a separate internal thread for better performance.
//
// There are two main use cases:
//
// 1) "Put & Forget:" you gain faster puts as you don't have to wait for the transaction to finish.
//
// 2) Many small transactions: if your write load is typically a lot of individual puts that happen in parallel,
// this will merge small transactions into bigger ones. This results in a significant gain in overall throughput.
//
//
// In situations with (extremely) high async load, this method may be throttled (~1ms) or delayed
// up to options.putAsyncTimeout (10 seconds by default). In the unlikely event that the object could
// not be enqueued after delaying (because of a full queue), an error will be returned.
//
// Note that this method does not give you hard durability guarantees like the synchronous Put provides.
// There is a small time window in which the data may not have been committed durably yet.
func (box *Box) PutAsync(object interface{}) (id uint64, err error) {
	return box.put(object, true, box.objectBox.options.putAsyncTimeout)
}

// Same as PutAsync but with a custom enqueue timeout
func (box *Box) PutAsyncWithTimeout(object interface{}, timeoutMs uint) (id uint64, err error) {
	return box.put(object, true, timeoutMs)
}

// Put synchronously inserts/updates a single object.
// In case the ID is not specified, it would be assigned automatically (auto-increment).
// When inserting, the ID property on the passed object will be assigned the new ID as well.
func (box *Box) Put(object interface{}) (id uint64, err error) {
	return box.put(object, false, 0)
}

// PutAll inserts multiple objects in single transaction.
// The given argument must be a slice of the object type this Box represents (pointers to objects).
// In case IDs are not set on the objects, they would be assigned automatically (auto-increment).
//
// Returns: IDs of the put objects (in the same order).
//
// Note: In case an error occurs during the transaction, some of the objects may already have the ID assigned
// even though the transaction has been rolled back and the objects are not stored under those IDs.
//
// Note: The slice may be empty or even nil; in both cases, an empty IDs slice and no error is returned.
func (box *Box) PutAll(objects interface{}) (ids []uint64, err error) {
	// TODO we need a sequential version as well, starting a transaction and calling multiple puts()

	var binding = box.entity.binding
	var slice = reflect.ValueOf(objects)
	var count = slice.Len()

	// a little optimization for the edge case
	if count == 0 {
		return []uint64{}, nil
	}

	// find out ids of all the objects & whether they're new objects or updates
	ids = make([]uint64, count)
	var isUpdate = make([]bool, count)

	// indexes of new objects (zero IDs) in the `ids` slice
	var indexesNewObjects = make([]int, 0)

	for i := 0; i < count; i++ {
		var object = slice.Index(i).Interface()
		if id, err := binding.GetId(object); err != nil {
			return nil, err
		} else if id > 0 {
			ids[i] = id
			isUpdate[i] = true
		} else {
			indexesNewObjects = append(indexesNewObjects, i)
		}
	}

	// if there are any new objects, reserve IDs for them
	// TODO in case of an error later during insert, this ID is not recovered (reused in the future)...
	//  we need to either move idForPut to C-api (preferrable) or use a transaction
	if newIds, err := box.idsForPut(len(indexesNewObjects)); err != nil {
		return nil, err
	} else {
		for i, id := range newIds {
			ids[indexesNewObjects[i]] = id
		}
	}

	//log.Printf("PutAll %v: %v", box.entity.name, ids)

	// flatten all the objects
	var allBytes = make([][]byte, count)
	for i := 0; i < count; i++ {
		var object = slice.Index(i).Interface()

		if box.entity.hasRelations {
			if err := box.entity.binding.PutRelated(box.objectBox, object, ids[i]); err != nil {
				return nil, err
			}
		}

		err = box.withObjectBytes(object, ids[i], func(bytes []byte) error {
			allBytes[i] = make([]byte, len(bytes))
			copy(allBytes[i], bytes)
			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	// create a C representation of the objects array
	bytesArray, err := goBytesArrayToC(allBytes)
	if err != nil {
		return nil, err
	} else {
		defer bytesArray.free()
	}

	if err := cMaybeErr(func() C.obx_err {
		return C.obx_box_put_array(box.box, bytesArray.cBytesArray, goUint64ArrayToCObxId(ids), goBoolArrayToC(isUpdate))
	}); err != nil {
		return nil, err
	}

	// restore update IDs on the new objects
	for index := range indexesNewObjects {
		binding.SetId(slice.Index(index).Interface(), ids[index])
	}

	return
}

// Remove deletes a single object
func (box *Box) Remove(id uint64) error {
	return cMaybeErr(func() C.obx_err {
		return C.obx_box_remove(box.box, C.obx_id(id))
	})
}

// RemoveAll removes all stored objects.
// This is much faster than removing objects one by one in a loop.
func (box *Box) RemoveAll() error {
	return cMaybeErr(func() C.obx_err {
		return C.obx_box_remove_all(box.box, nil)
	})
}

// Count returns a number of objects stored
func (box *Box) Count() (uint64, error) {
	return box.CountMax(0)
}

// CountMax returns a number of objects stored (up to a given maximum)
// passing limit=0 is the same as calling Count() - counts all objects without a limit
func (box *Box) CountMax(limit uint64) (uint64, error) {
	var cResult C.uint64_t
	if err := cMaybeErr(func() C.obx_err { return C.obx_box_count(box.box, C.uint64_t(limit), &cResult) }); err != nil {
		return 0, err
	}
	return uint64(cResult), nil
}

// IsEmpty checks whether the box contains any objects
func (box *Box) IsEmpty() (bool, error) {
	var cResult C.bool
	if err := cMaybeErr(func() C.obx_err { return C.obx_box_is_empty(box.box, &cResult) }); err != nil {
		return false, err
	}
	return bool(cResult), nil
}

// Get reads a single object.
//
// Returns an interface that should be cast to the appropriate type.
// Returns nil in case the object with the given ID doesn't exist.
// The cast is done automatically when using the generated BoxFor* code.
func (box *Box) Get(id uint64) (object interface{}, err error) {
	var data *C.void
	var dataSize C.size_t
	var dataPtr = unsafe.Pointer(data)

	// TODO use runInTxn so that the related entities are fetched within the same transaction in binding.Load()
	var rc = C.obx_box_get(box.box, C.obx_id(id), &dataPtr, &dataSize)
	if rc == 0 {
		var bytes []byte
		cVoidPtrToByteSlice(dataPtr, int(dataSize), &bytes)
		return box.entity.binding.Load(box.objectBox, bytes)
	} else if rc == C.OBX_NOT_FOUND {
		return nil, nil
	} else {
		return nil, createError()
	}
}

// GetMany reads multiple objects at once.
//
// Returns a slice of objects that should be cast to the appropriate type.
// The cast is done automatically when using the generated BoxFor* code.
// If any of the objects doesn't exist, its position in the return slice
//  is nil or an empty object (depends on the binding)
func (box *Box) GetMany(ids ...uint64) (slice interface{}, err error) {
	if cIds, err := goIdsArrayToC(ids); err != nil {
		return nil, err
	} else if supportsBytesArray {
		data, err := cGetBytesArray(func() *C.OBX_bytes_array { return C.obx_box_get_ids(box.box, cIds.cArray) })
		if err != nil {
			return nil, err
		}
		return box.bytesArrayToObjects(data)

	} else {
		var cCall = func(visitorArg unsafe.Pointer) C.obx_err {
			return C.obx_box_visit_ids(box.box, cIds.cArray, C.data_visitor, visitorArg)
		}
		return readUsingVisitor(box.objectBox, box.entity.binding, defaultSliceCapacity, cCall)
	}
}

// GetAll reads all stored objects.
//
// Returns a slice of objects that should be cast to the appropriate type.
// The cast is done automatically when using the generated BoxFor* code.
func (box *Box) GetAll() (slice interface{}, err error) {
	if supportsBytesArray {
		data, err := cGetBytesArray(func() *C.OBX_bytes_array { return C.obx_box_get_all(box.box) })
		if err != nil {
			return nil, err
		}
		return box.bytesArrayToObjects(data)

	} else {
		var cCall = func(visitorArg unsafe.Pointer) C.obx_err {
			return C.obx_box_visit(box.box, C.data_visitor, visitorArg)
		}
		return readUsingVisitor(box.objectBox, box.entity.binding, defaultSliceCapacity, cCall)
	}
}

func (box *Box) bytesArrayToObjects(bytesArray [][]byte) (slice interface{}, err error) {
	var binding = box.entity.binding
	slice = binding.MakeSlice(len(bytesArray))
	for _, bytesData := range bytesArray {
		if object, err := binding.Load(box.objectBox, bytesData); err != nil {
			return nil, err
		} else {
			slice = binding.AppendToSlice(slice, object)
		}
	}
	return slice, nil
}

// Contains checks whether an object with the given ID is stored.
func (box *Box) Contains(id uint64) (bool, error) {
	var cResult C.bool
	if err := cMaybeErr(func() C.obx_err { return C.obx_box_contains(box.box, C.obx_id(id), &cResult) }); err != nil {
		return false, err
	}
	return bool(cResult), nil
}

// RelationIds returns IDs of all target objects related to the given source object ID
func (box *Box) RelationIds(relation *RelationToMany, sourceId uint64) ([]uint64, error) {
	return cGetIds(func() *C.OBX_id_array {
		return C.obx_box_rel_targets_ids(box.box, C.obx_schema_id(relation.Id), C.obx_id(sourceId))
	})
}

// RelationReplace replaces all targets for a given source in a standalone many-to-many relation
// It also inserts new related objects (with a 0 ID).
func (box *Box) RelationReplace(relation *RelationToMany, sourceId uint64, sourceObject interface{},
	targetObjects interface{}) (err error) {

	// TODO this whole func needs to be executed in a single transaction because it calls to the c-box multiple times

	// get id from the object, if inserting, it would be 0 even if the argument id is already non-zero
	// this saves us an unnecessary request to RelationIds for new objects (there can't be any relations yet)
	objId, err := box.entity.binding.GetId(sourceObject)
	if err != nil {
		return err
	}

	// make a map of related target entity IDs, marking those that were originally related but should be removed
	var idsToRemove = make(map[uint64]bool)

	if objId != 0 {
		if oldRelIds, err := box.RelationIds(relation, sourceId); err != nil {
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
		var targetBox = box.objectBox.InternalBox(relation.Target.Id)

		// walk over the current related objects, mark those that still exist, add the new ones
		for i := 0; i < count; i++ {
			var reflObj = sliceValue.Index(i)
			var rel interface{}
			if reflObj.Kind() == reflect.Ptr {
				rel = reflObj.Interface()
			} else {
				rel = reflObj.Addr().Interface()
			}

			rId, err := targetBox.entity.binding.GetId(rel)
			if err != nil {
				return err
			} else if rId == 0 {
				if rId, err = targetBox.Put(rel); err != nil {
					return err
				}
			}

			if idsToRemove[rId] {
				// old relation that still exists, keep it
				delete(idsToRemove, rId)
			} else {
				// new relation, add it
				if err := box.RelationPut(relation, sourceId, rId); err != nil {
					return err
				}
			}
		}
	}

	// remove those that were not found in the rSlice but were originally related to this entity
	for rId := range idsToRemove {
		if err := box.RelationRemove(relation, sourceId, rId); err != nil {
			return err
		}
	}

	return nil
}

// RelationPut creates a relation between the given source & target objects
func (box *Box) RelationPut(relation *RelationToMany, sourceId, targetId uint64) error {
	//log.Printf("RelationPut %v: %v (%s) -> %v (%s)", relation.Id,
	//	sourceId, box.objectBox.getEntityById(relation.Source.Id).name,
	//	targetId, box.objectBox.getEntityById(relation.Target.Id).name)
	return cMaybeErr(func() C.obx_err {
		return C.obx_box_rel_put(box.box, C.obx_schema_id(relation.Id), C.obx_id(sourceId), C.obx_id(targetId))
	})
}

// RelationRemove removes a relation between the given source & target objects
func (box *Box) RelationRemove(relation *RelationToMany, sourceId, targetId uint64) error {
	//log.Printf("RelationRemove %v: %v (%s) -> %v (%s)", relation.Id,
	//	sourceId, box.objectBox.getEntityById(relation.Source.Id).name,
	//	targetId, box.objectBox.getEntityById(relation.Target.Id).name)
	return cMaybeErr(func() C.obx_err {
		return C.obx_box_rel_remove(box.box, C.obx_schema_id(relation.Id), C.obx_id(sourceId), C.obx_id(targetId))
	})
}
