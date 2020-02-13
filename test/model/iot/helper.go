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

package iot

import (
	"github.com/objectbox/objectbox-go/objectbox"
	"github.com/objectbox/objectbox-go/test/assert"
	"io/ioutil"
	"os"
	"strconv"
)

type TestEnv struct {
	*objectbox.ObjectBox
	dir string
}

// Close closes ObjectBox and removes the database
func (env *TestEnv) Close() {
	env.ObjectBox.Close()
	os.RemoveAll(env.dir)
}

// NewTestEnv creates an empty ObjectBox instance
func NewTestEnv() *TestEnv {

	// Test in a temporary directory - if tested by an end user, the repo is read-only.
	tempDir, err := ioutil.TempDir("", "objectbox-test")
	if err != nil {
		panic(err)
	}

	objectBox, err := objectbox.NewBuilder().Directory(tempDir).Model(ObjectBoxModel()).Build()
	if err != nil {
		panic(err)
	}

	return &TestEnv{
		ObjectBox: objectBox,
		dir:       tempDir,
	}
}

// PutEvent creates an event
func PutEvent(ob *objectbox.ObjectBox, device string, date int64) *Event {
	event := Event{Device: device, Date: date}
	id, err := BoxForEvent(ob).Put(&event)
	assert.NoErr(nil, err)
	event.Id = id
	return &event
}

// PutReading creates a reading
func PutReading(ob *objectbox.ObjectBox, name string, ValueString string, ValueInteger int64, ValueFloating float64, ValueInt32 int32, ValueFloating32 float32) *Reading {
	event := Reading{ValueName: name, ValueString: ValueString, ValueInteger: ValueInteger, ValueFloating: ValueFloating, ValueInt32: ValueInt32, ValueFloating32: ValueFloating32}
	id, err := BoxForReading(ob).Put(&event)
	assert.NoErr(nil, err)
	event.Id = id
	return &event
}

// PutEvents creates multiple events
func PutEvents(ob *objectbox.ObjectBox, count int) []*Event {
	// TODO TX
	events := make([]*Event, 0, count)
	for i := 1; i <= count; i++ {
		event := PutEvent(ob, "device "+strconv.Itoa(i), int64(10000+i))
		events = append(events, event)
	}
	return events
}

// PutReadings creates multiple readings
func PutReadings(ob *objectbox.ObjectBox, count int) []*Reading {
	// TODO TX
	readings := make([]*Reading, 0, count)
	for i := 1; i <= count; i++ {
		reading := PutReading(ob, "reading"+strconv.Itoa(i), "string"+strconv.Itoa(i), int64(10000+i), float64(10000+i), int32(10000+i), float32(10000+i))
		readings = append(readings, reading)
	}
	return readings
}
