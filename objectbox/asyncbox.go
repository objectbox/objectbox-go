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
#include <stdlib.h>
#include "objectbox.h"
*/
import "C"
import (
	"errors"
	"unsafe"
)

// AsyncBox provides asynchronous operations on objects of a common type.
//
// Asynchronous operations are executed on a separate internal thread for better performance.
//
// There are two main use cases:
//
// 1) "execute & forget:" you gain faster put/remove operations as you don't have to wait for the transaction to finish.
//
// 2) Many small transactions: if your write load is typically a lot of individual puts that happen in parallel,
// this will merge small transactions into bigger ones. This results in a significant gain in overall throughput.
//
// In situations with (extremely) high async load, an async method may be throttled (~1ms) or delayed up to 1 second.
// In the unlikely event that the object could still not be enqueued (full queue), an error will be returned.
//
// Note that async methods do not give you hard durability guarantees like the synchronous Box provides.
// There is a small time window in which the data may not have been committed durably yet.
type AsyncBox struct {
	box    *Box
	cAsync *C.OBX_async
	cOwned bool // whether the cAsync resource is owned by this struct
}

// NewAsyncBox creates a new async box with the given operation timeout in case an async queue is full.
// The returned struct must be freed explicitly using the Close() method.
// It's usually preferable to use Box::Async() which takes care of resource management and doesn't require closing.
func NewAsyncBox(ob *ObjectBox, entityId TypeId, timeoutMs uint64) (*AsyncBox, error) {
	var async = &AsyncBox{
		box:    ob.InternalBox(entityId),
		cOwned: true,
	}

	if err := cCallBool(func() bool {
		async.cAsync = C.obx_async_create(async.box.cBox, C.uint64_t(timeoutMs))
		return async.cAsync != nil
	}); err != nil {
		return nil, err
	}

	return async, nil
}

// Close frees resources of a customized AsyncBox (e.g. with a custom timeout).
// Not necessary for the standard (shared) instance from box.Async(); Close() can still be called for those:
// it just won't have any effect.
func (async *AsyncBox) Close() error {
	if !async.cOwned || async.cAsync == nil {
		return nil
	}
	var cAsync = async.cAsync
	async.cAsync = nil
	return cCall(func() C.obx_err {
		return C.obx_async_close(cAsync)
	})
}

func (async *AsyncBox) put(object interface{}, mode int) (uint64, error) {
	entity := async.box.entity
	idFromObject, err := entity.binding.GetId(object)
	if err != nil {
		return 0, err
	}

	if entity.hasRelations {
		return 0, errors.New("asynchronous Put/Insert/Update is currently not supported on entities that have" +
			" relations because it could result in partial inserts/broken relations")
	}

	id, err := async.box.idForPut(idFromObject)
	if err != nil {
		return 0, err
	}

	err = async.box.withObjectBytes(object, id, func(bytes []byte) error {
		return cCall(func() C.obx_err {
			return C.obx_async_put_mode(async.cAsync, C.obx_id(id), unsafe.Pointer(&bytes[0]), C.size_t(len(bytes)),
				C.OBXPutMode(mode))
		})
	})

	if err != nil {
		return 0, err
	}

	// update the id on the object
	if idFromObject != id {
		entity.binding.SetId(object, id)
	}

	return id, nil
}

// Put inserts/updates a single object asynchronously.
// When inserting a new object, the ID property on the passed object will be assigned a new ID the entity would hold
// if the insert is ultimately successful. The newly assigned ID may not become valid if the insert fails.
func (async *AsyncBox) Put(object interface{}) (id uint64, err error) {
	return async.put(object, cPutModePut)
}

// Insert a single object asynchronously.
// The ID property on the passed object will be assigned a new ID the entity would hold if the insert is ultimately
// successful. The newly assigned ID may not become valid if the insert fails.
// Fails silently if an object with the same ID already exists (this error is not returned).
func (async *AsyncBox) Insert(object interface{}) (id uint64, err error) {
	return async.put(object, cPutModeInsert)
}

// Update a single object asynchronously.
// The object must already exists or the update fails silently (without an error returned).
func (async *AsyncBox) Update(object interface{}) error {
	_, err := async.put(object, cPutModeUpdate)
	return err
}

// Remove deletes a single object asynchronously.
func (async *AsyncBox) Remove(object interface{}) error {
	id, err := async.box.entity.binding.GetId(object)
	if err != nil {
		return err
	}

	return async.RemoveId(id)
}

// RemoveId deletes a single object asynchronously.
func (async *AsyncBox) RemoveId(id uint64) error {
	return cCall(func() C.obx_err {
		return C.obx_async_remove(async.cAsync, C.obx_id(id))
	})
}

// Awaits for all (including future) async submissions to be completed (the async queue becomes idle for a moment).
// Currently this is not limited to the single entity this AsyncBox is working on but all entities in the store.
// Returns an error if shutting down or an error occurred
func (async *AsyncBox) AwaitCompletion() error {
	return cCallBool(func() bool {
		return bool(C.obx_store_await_async_completion(async.box.ObjectBox.store))
	})
}

// Awaits for previously submitted async operations to be completed (the async queue does not have to become idle).
// Currently this is not limited to the single entity this AsyncBox is working on but all entities in the store.
// Returns an error if shutting down or an error occurred
func (async *AsyncBox) AwaitSubmitted() error {
	return cCallBool(func() bool {
		return bool(C.obx_store_await_async_submitted(async.box.ObjectBox.store))
	})
}
