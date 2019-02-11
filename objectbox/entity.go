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

import (
	"sync"
	"sync/atomic"
)

// this is used publicly in the model/bindings
type Entity struct {
	Id TypeId
}

// this is used internally to automatically synchronize
// note that it must be a singleton
type entity struct {
	objectBox *ObjectBox
	id        TypeId
	name      string
	binding   ObjectBinding

	// whether this entity has any relations (standalone or property-rels) - configured during model creation
	hasRelations bool

	// whether there was an asynchronous operation recently
	isOutOfSync uint32

	// locked when currently waiting for async completion
	mutex sync.Mutex
}

func (e *entity) markOutOfSync() {
	atomic.StoreUint32(&e.isOutOfSync, aTrue)
}

func (e *entity) awaitAsyncCompletion() {
	// if this entity is currently in-sync, no need to do anything
	if aTrue != atomic.LoadUint32(&e.isOutOfSync) {
		return
	}

	e.mutex.Lock()
	defer e.mutex.Unlock()

	// check again after getting the mutex, it might have been cleared in the meantime
	if aTrue == atomic.LoadUint32(&e.isOutOfSync) {
		e.objectBox.AwaitAsyncCompletion()
		atomic.StoreUint32(&e.isOutOfSync, aFalse)
	}
}
