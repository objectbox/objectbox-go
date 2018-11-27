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

package iot

import (
	"testing"

	"github.com/objectbox/objectbox-go/test/assert"
)

func TestObjectBoxEvents(t *testing.T) {
	objectBox := LoadEmptyTestObjectBox()

	box := BoxForEvent(objectBox)
	assert.NoErr(t, box.RemoveAll())

	event := Event{
		Device: "my device",
	}
	objectId, err := box.Put(&event)
	assert.NoErr(t, err)
	t.Logf("Added object ID %v", objectId)

	event.Device = "2nd device"
	objectId, err = box.Put(&event)
	assert.NoErr(t, err)
	t.Logf("Added 2nd object ID %v", objectId)

	eventRead, err := box.Get(objectId)
	assert.NoErr(t, err)
	if objectId != eventRead.Id || event.Device != eventRead.Device {
		t.Fatalf("Event data error: %v vs. %v", event, eventRead)
	}

	all, err := box.GetAll()
	assert.NoErr(t, err)
	assert.EqInt(t, 2, len(all))
}

func TestObjectBoxReadings(t *testing.T) {
	objectBox := LoadEmptyTestObjectBox()
	box := BoxForReading(objectBox)
	assert.NoErr(t, box.RemoveAll())
	reading := Reading{
		ValueName:    "Temperature",
		ValueInteger: 77,
	}
	objectId, err := box.Put(&reading)
	assert.NoErr(t, err)
	t.Logf("Added object ID %v", objectId)

	reading.ValueInteger = 100
	objectId, err = box.Put(&reading)
	assert.NoErr(t, err)
	t.Logf("Added 2nd object ID %v", objectId)

	readingRead, err := box.Get(objectId)
	assert.NoErr(t, err)
	if objectId != readingRead.Id || reading.ValueInteger != readingRead.ValueInteger {
		t.Fatalf("Event data error: %v vs. %v", reading, readingRead)
	}

	all, err := box.GetAll()
	assert.NoErr(t, err)
	assert.EqInt(t, 2, len(all))
}
