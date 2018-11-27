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

package objectbox_test

import (
	"sync"
	"testing"

	"github.com/objectbox/objectbox-go/test/assert"
	"github.com/objectbox/objectbox-go/test/model/iot"
)

func TestConcurrentPut(t *testing.T) {
	concurrentInsert(t, 100, 20, false)
}

func TestConcurrentPutAsync(t *testing.T) {
	concurrentInsert(t, 100000, 20, true)
}

func concurrentInsert(t *testing.T, count, concurrency int, putAsync bool) {
	objectBox := iot.LoadEmptyTestObjectBox()
	defer objectBox.Close()
	box := iot.BoxForEvent(objectBox)
	defer box.Close()

	err := box.RemoveAll()
	assert.NoErr(t, err)

	event := iot.Event{
		Device: "my device",
	}

	var countPart = count / concurrency
	assert.EqInt(t, 0, count%concurrency)

	// prepare channels and launch the goroutines
	ids := make(chan uint64, count)
	errors := make(chan error, count)

	t.Logf("launching %d routines to insert %d objects each", concurrency, countPart)

	var wg sync.WaitGroup
	wg.Add(concurrency)
	for i := concurrency; i > 0; i-- {
		go func() {
			defer wg.Done()
			for i := countPart; i > 0; i-- {
				var id uint64
				var e error

				if putAsync {
					id, e = box.PutAsync(&event)
				} else {
					id, e = box.Put(&event)
				}

				if e != nil {
					errors <- e
				} else {
					ids <- id
				}
			}
		}()
	}

	// collect and check results after everything is done
	t.Log("waiting for all goroutines to finish")
	wg.Wait()

	if putAsync {
		t.Log("waiting until all PutAsync are processed")
		objectBox.AwaitAsyncCompletion()
	}

	t.Log("validating counts")
	assert.EqInt(t, 0, len(errors))
	assert.EqInt(t, count, len(ids))

	actualCount, err := box.Count()
	assert.NoErr(t, err)
	assert.Eq(t, uint64(count), actualCount)

	// check whether the IDs are unique
	t.Log("validating IDs")
	idsMap := make(map[uint64]bool)
	for i := count; i > 0; i-- {
		id := <-ids
		if idsMap[id] != false {
			assert.Failf(t, "duplicate ID %d", id)
		} else {
			idsMap[id] = true
		}
	}
}
