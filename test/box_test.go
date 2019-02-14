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
	"testing"

	"github.com/objectbox/objectbox-go/test/assert"
	"github.com/objectbox/objectbox-go/test/model"
	"github.com/objectbox/objectbox-go/test/model/iot"
)

func TestBox(t *testing.T) {
	objectBox := iot.LoadEmptyTestObjectBox()
	defer objectBox.Close()
	box1 := iot.BoxForEvent(objectBox)
	box2 := iot.BoxForEvent(objectBox)

	assert.Eq(t, box1.Box, box2.Box)
}

func TestPutAsync(t *testing.T) {
	objectBox := iot.LoadEmptyTestObjectBox()
	defer objectBox.Close()
	box := iot.BoxForEvent(objectBox)
	err := box.RemoveAll()
	assert.NoErr(t, err)

	event := iot.Event{
		Device: "my device",
	}
	objectId, err := box.PutAsync(&event)
	assert.NoErr(t, err)
	assert.Eq(t, objectId, event.Id)

	objectBox.AwaitAsyncCompletion()

	count, err := box.Count()
	assert.NoErr(t, err)
	assert.Eq(t, uint64(1), count)

	eventRead, err := box.Get(objectId)
	assert.NoErr(t, err)
	if objectId != eventRead.Id || event.Device != eventRead.Device {
		t.Fatalf("Event data error: %v vs. %v", event, eventRead)
	}

	err = box.Remove(eventRead)
	assert.NoErr(t, err)

	eventRead, err = box.Get(objectId)
	assert.NoErr(t, err)
	if eventRead != nil {
		t.Fatalf("object hasn't been deleted by box.Remove()")
	}

	count, err = box.Count()
	assert.NoErr(t, err)
	assert.Eq(t, uint64(0), count)
}

func TestUnique(t *testing.T) {
	objectBox := iot.LoadEmptyTestObjectBox()
	defer objectBox.Close()
	box := iot.BoxForEvent(objectBox)

	err := box.RemoveAll()
	assert.NoErr(t, err)

	_, err = box.Put(&iot.Event{
		Device: "my device",
		Uid:    "duplicate-uid",
	})
	assert.NoErr(t, err)

	_, err = box.Put(&iot.Event{
		Device: "my device 2",
		Uid:    "duplicate-uid",
	})
	if err == nil {
		assert.Failf(t, "put() passed instead of an expected unique constraint violation")
	}

	count, err := box.Count()
	assert.NoErr(t, err)
	assert.Eq(t, uint64(1), count)
}

func TestPutAll(t *testing.T) {
	objectBox := iot.LoadEmptyTestObjectBox()
	defer objectBox.Close()
	box := iot.BoxForEvent(objectBox)

	err := box.RemoveAll()
	assert.NoErr(t, err)

	event1 := iot.Event{
		Device: "Pi 3B",
	}
	event2 := iot.Event{
		Device: "Pi Zero",
	}
	events := []*iot.Event{&event1, &event2}
	objectIds, err := box.PutAll(events)
	assert.NoErr(t, err)
	assert.Eq(t, uint64(1), objectIds[0])
	assert.Eq(t, objectIds[0], events[0].Id)
	assert.Eq(t, uint64(2), objectIds[1])
	assert.Eq(t, objectIds[1], events[1].Id)

	count, err := box.Count()
	assert.NoErr(t, err)
	assert.Eq(t, uint64(2), count)

	eventRead, err := box.Get(objectIds[0])
	assert.NoErr(t, err)
	assert.Eq(t, "Pi 3B", eventRead.Device)

	eventRead, err = box.Get(objectIds[1])
	assert.NoErr(t, err)
	assert.Eq(t, "Pi Zero", eventRead.Device)

	// And passing nil & empty slice
	objectIds, err = box.PutAll(nil)
	assert.NoErr(t, err)
	assert.Eq(t, len(objectIds), 0)
	//noinspection GoPreferNilSlice
	noEvents := []*iot.Event{}
	objectIds, err = box.PutAll(noEvents)
	assert.NoErr(t, err)
	assert.Eq(t, len(objectIds), 0)
}

func TestPut(t *testing.T) {
	objectBox := iot.LoadEmptyTestObjectBox()
	defer objectBox.Close()
	box := iot.BoxForEvent(objectBox)

	assert.NoErr(t, box.RemoveAll())

	event := iot.Event{
		Device: "my device",
	}
	objectId, err := box.Put(&event)
	assert.NoErr(t, err)
	assert.Eq(t, objectId, event.Id)
	t.Logf("Added object ID %v", objectId)

	event2 := iot.Event{
		Device: "2nd device",
	}
	objectId2, err := box.Put(&event2)
	assert.NoErr(t, err)
	t.Logf("Added 2nd object ID %v", objectId2)

	// read the previous object and compare
	eventRead, err := box.Get(objectId)
	assert.NoErr(t, err)
	assert.Eq(t, event, *eventRead)

	all, err := box.GetAll()
	assert.NoErr(t, err)
	assert.Eq(t, 2, len(all))
	assert.Eq(t, &event, all[0])
	assert.Eq(t, &event2, all[1])
}

func TestBoxCount(t *testing.T) {
	var env = model.NewTestEnv(t)
	defer env.Close()

	var c = uint64(10)
	env.Populate(uint(c))

	count, err := env.Box.Count()
	assert.NoErr(t, err)
	assert.Eq(t, c, count)

	count, err = env.Box.CountMax(c / 2)
	assert.NoErr(t, err)
	assert.Eq(t, c/2, count)
}

func TestBoxEmpty(t *testing.T) {
	var env = model.NewTestEnv(t)
	defer env.Close()

	isEmpty, err := env.Box.IsEmpty()
	assert.NoErr(t, err)
	assert.Eq(t, true, isEmpty)

	env.Populate(10)

	isEmpty, err = env.Box.IsEmpty()
	assert.NoErr(t, err)
	assert.Eq(t, false, isEmpty)

	env.Box.RemoveAll()

	isEmpty, err = env.Box.IsEmpty()
	assert.NoErr(t, err)
	assert.Eq(t, true, isEmpty)
}

func TestBoxContains(t *testing.T) {
	var env = model.NewTestEnv(t)
	defer env.Close()

	found, err := env.Box.Contains(1)
	assert.NoErr(t, err)
	assert.Eq(t, false, found)

	env.Populate(1)

	found, err = env.Box.Contains(1)
	assert.NoErr(t, err)
	assert.Eq(t, true, found)

	found, err = env.Box.Contains(2)
	assert.NoErr(t, err)
	assert.Eq(t, false, found)
}

// Includes testing the default string vector (containing 2 normal values and one "")
func TestBoxPutData(t *testing.T) {
	env := model.NewTestEnv(t)
	defer env.Close()

	var inserted = model.Entity47()

	id, err := env.Box.Put(inserted)
	assert.NoErr(t, err)

	read, err := env.Box.Get(id)
	assert.NoErr(t, err)
	assert.Eq(t, inserted, read)
}

func TestBoxPutAndGetStringVectorsEmptyAndNil(t *testing.T) {
	env := model.NewTestEnv(t)
	defer env.Close()

	var inserted = model.Entity47()

	// test empty vectors
	inserted.StringVector = []string{}
	inserted.ByteVector = []byte{}

	id, err := env.Box.Put(inserted)
	assert.NoErr(t, err)

	read, err := env.Box.Get(id)
	assert.NoErr(t, err)
	assert.Eq(t, *inserted, *read)

	// test nil vectors
	inserted.StringVector = nil
	inserted.ByteVector = nil

	id, err = env.Box.Put(inserted)
	assert.NoErr(t, err)

	read, err = env.Box.Get(id)
	assert.NoErr(t, err)
	assert.Eq(t, *inserted, *read)
}
