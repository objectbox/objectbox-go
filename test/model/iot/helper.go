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
	"os"
	"strconv"

	"github.com/objectbox/objectbox-go/objectbox"
	"github.com/objectbox/objectbox-go/test/assert"
)

func LoadEmptyTestObjectBox() *objectbox.ObjectBox {
	var dbName = "iot-test"

	// remove database beforehand - only used in tests
	os.RemoveAll(dbName)

	objectBox, err := objectbox.NewBuilder().Directory(dbName).Model(ObjectBoxModel()).Build()
	if err != nil {
		panic(err)
	}
	return objectBox
}

func PutEvent(ob *objectbox.ObjectBox, device string, date int64) *Event {
	event := Event{Device: device, Date: date}
	id, err := ob.Box(1).Put(&event)
	assert.NoErr(nil, err)
	event.Id = id
	return &event
}

func PutEvents(ob *objectbox.ObjectBox, count int) []*Event {
	// TODO TX
	events := make([]*Event, 0, count)
	for i := 1; i <= count; i++ {
		event := PutEvent(ob, "device "+strconv.Itoa(i), int64(10000+i))
		events = append(events, event)
	}
	return events
}
