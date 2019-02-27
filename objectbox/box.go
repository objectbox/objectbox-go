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
	"fmt"
	"reflect"
	"sync/atomic"
	"unsafe"

	"github.com/google/flatbuffers/go"
)

// TODO remove this temporary transaction variable from the methods that use it
// transaction is handled by Box implicitly
var todoTemporaryTxn = &Transaction{}

// Box provides CRUD access to objects of a common type
type Box struct {
	objectBox *ObjectBox
	box       *C.OBX_box
	entity    *entity

	// Must be used in combination with fbbInUseAtomic
	fbb *flatbuffers.Builder

	// Values 0 (fbb available) or 1 (fbb in use); use only with CompareAndSwapInt32
	fbbInUseAtomic uint32
}

// close fully closes the Box connection and free's resources
func (box *Box) close() (err error) {
	rc := C.obx_box_close(box.box)
	box.box = nil
	if rc != 0 {
		err = createError()
	}
	return
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

	query, err = builder.Build()

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

	cIdsArray := C.obx_box_bulk_ids_for_put(box.box, C.uint64_t(count))
	if cIdsArray == nil {
		return nil, createError()
	}

	idsArray := cIdsArrayToGo(cIdsArray)
	defer idsArray.free()

	if len(idsArray.ids) != count {
		return nil, fmt.Errorf("invalid number of new IDs reserved: got %v instead of required %v",
			len(idsArray.ids), count)
	}

	return idsArray.ids, nil
}

func (box *Box) put(object interface{}, async bool, timeoutMs uint) (id uint64, err error) {
	idFromObject, err := box.entity.binding.GetId(object)
	if err != nil {
		return 0, err
	}

	// TODO in case of an error later during insert, this ID is not recovered (reused in the future)...
	//  we need to either move idForPut to C-api (preferrable) or use a transaction
	if id, err = box.idForPut(idFromObject); err != nil {
		return 0, err
	}

	if box.entity.hasRelations {
		err = box.entity.binding.PutRelated(todoTemporaryTxn, object, id)
		if err != nil {
			return 0, err
		}
	}

	err = box.withObjectBytes(object, id, func(bytes []byte) error {
		var rc C.obx_err
		if async {
			box.entity.markOutOfSync()
			rc = C.obx_box_put_async(box.box, C.obx_id(id), unsafe.Pointer(&bytes[0]), C.size_t(len(bytes)),
				C.bool(idFromObject != 0), C.uint64_t(timeoutMs))
		} else {
			rc = C.obx_box_put(box.box, C.obx_id(id), unsafe.Pointer(&bytes[0]), C.size_t(len(bytes)),
				C.bool(idFromObject != 0))
		}

		if rc != 0 {
			return createError()
		}
		return nil
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
	var fbb *flatbuffers.Builder
	if atomic.CompareAndSwapUint32(&box.fbbInUseAtomic, 0, 1) {
		defer func() {
			atomic.StoreUint32(&box.fbbInUseAtomic, 0)
			fbb.Reset()
		}()
		fbb = box.fbb
	} else {
		fbb = flatbuffers.NewBuilder(256)
	}

	if err := box.entity.binding.Flatten(object, fbb, id); err != nil {
		return err
	}

	fbb.Finish(fbb.EndObject())

	return fn(fbb.FinishedBytes())
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

	// flatten all the objects
	var allBytes = make([][]byte, count)
	for i := 0; i < count; i++ {
		var object = slice.Index(i).Interface()
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

	var rc = C.obx_box_bulk_put(box.box, bytesArray.cBytesArray, goUint64ArrayToCObxId(ids), goBoolArrayToC(isUpdate))
	if rc != 0 {
		return nil, createError()
	}

	// restore update IDs on the new objects
	for index := range indexesNewObjects {
		binding.SetId(slice.Index(index).Interface(), ids[index])
	}

	return
}

// Remove deletes a single object
func (box *Box) Remove(id uint64) (err error) {
	return box.objectBox.runWithCursor(box.entity, false, func(cursor *Cursor) error {
		return cursor.Remove(id)
	})
}

// RemoveAll removes all stored objects.
// This is much faster than removing objects one by one in a loop.
func (box *Box) RemoveAll() (err error) {
	return box.objectBox.runWithCursor(box.entity, false, func(cursor *Cursor) error {
		return cursor.RemoveAll()
	})
}

// Count returns a number of objects stored
func (box *Box) Count() (count uint64, err error) {
	err = box.objectBox.runWithCursor(box.entity, true, func(cursor *Cursor) error {
		var errInner error
		count, errInner = cursor.Count()
		return errInner
	})
	return
}

// CountMax returns a number of objects stored (up to a given maximum)
func (box *Box) CountMax(max uint64) (count uint64, err error) {
	err = box.objectBox.runWithCursor(box.entity, true, func(cursor *Cursor) error {
		var errInner error
		count, errInner = cursor.CountMax(max)
		return errInner
	})
	return
}

// IsEmpty checks whether the box contains any objects
func (box *Box) IsEmpty() (result bool, err error) {
	err = box.objectBox.runWithCursor(box.entity, true, func(cursor *Cursor) error {
		var errInner error
		result, errInner = cursor.IsEmpty()
		return errInner
	})
	return
}

// Get reads a single object.
//
// Returns an interface that should be cast to the appropriate type.
// Returns nil in case the object with the given ID doesn't exist.
// The cast is done automatically when using the generated BoxFor* code
func (box *Box) Get(id uint64) (object interface{}, err error) {
	var data *C.void
	var dataSize C.size_t
	var dataPtr = unsafe.Pointer(data)

	// TODO use runInTxn so that the related entities are fetched within the same transaction in binding.Load()
	var rc = C.obx_box_get(box.box, C.obx_id(id), &dataPtr, &dataSize)
	if rc == 0 {
		var bytes []byte
		cVoidPtrToByteSlice(dataPtr, int(dataSize), &bytes)
		return box.entity.binding.Load(todoTemporaryTxn, bytes)
	} else if rc == C.OBX_NOT_FOUND {
		return nil, nil
	} else {
		return nil, createError()
	}
}

// GetAll reads all stored objects
//
// Returns a slice of objects that should be cast to the appropriate type.
// The cast is done automatically when using the generated BoxFor* code
func (box *Box) GetAll() (slice interface{}, err error) {
	if supportsBytesArray {
		cBytesArray := C.obx_box_get_all(box.box)
		if cBytesArray == nil {
			return nil, createError()
		}
		bytesArray := cBytesArrayToGo(cBytesArray)
		defer bytesArray.free()
		return box.bytesArrayToObjects(bytesArray)
	} else {
		return box.getAllSequential()
	}
}

func (box *Box) getAllSequential() (slice interface{}, err error) {
	var cCall = func(visitorArg unsafe.Pointer) C.obx_err {
		return C.obx_box_visit(box.box, C.data_visitor, visitorArg)
	}

	return runDataVisitor(box.entity.binding, defaultSliceCapacity, cCall)
}

func (box *Box) bytesArrayToObjects(bytesArray *bytesArray) (slice interface{}, err error) {
	var binding = box.entity.binding
	slice = binding.MakeSlice(len(bytesArray.BytesArray))
	for _, bytesData := range bytesArray.BytesArray {
		if object, err := binding.Load(todoTemporaryTxn, bytesData); err != nil {
			return nil, err
		} else {
			slice = binding.AppendToSlice(slice, object)
		}
	}
	return slice, nil
}

// Contains checks whether an object with the given ID is stored.
func (box *Box) Contains(id uint64) (bool, error) {
	var found = false
	var err = box.objectBox.runWithCursor(box.entity, true, func(cursor *Cursor) error {
		var errInner error
		found, errInner = cursor.seek(id)
		return errInner
	})
	return found, err
}
