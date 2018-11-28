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

	"github.com/objectbox/objectbox-go/objectbox"
	"github.com/objectbox/objectbox-go/test/assert"
	"github.com/objectbox/objectbox-go/test/model/iot"
)

func TestQueryBuilder(t *testing.T) {
	objectBox := iot.LoadEmptyTestObjectBox()
	defer objectBox.Close()
	box := iot.BoxForEvent(objectBox)
	defer box.Close()
	box.RemoveAll()

	qb := objectBox.Query(1)
	assert.NoErr(t, qb.Err)
	query, err := qb.BuildAndClose()
	assert.NoErr(t, err)
	defer query.Close()

	bytesArray, err := query.FindBytes()
	assert.NoErr(t, err)
	assert.EqInt(t, 0, len(bytesArray.BytesArray))

	slice, err := query.Find()
	assert.NoErr(t, err)
	assert.EqInt(t, 0, len(slice.([]*iot.Event)))

	event := iot.Event{
		Device: "dev1",
	}
	id1, err := box.Put(&event)
	assert.NoErr(t, err)

	event.Device = "dev2"
	id2, err := box.Put(&event)
	assert.NoErr(t, err)

	bytesArray, err = query.FindBytes()
	assert.NoErr(t, err)
	assert.EqInt(t, 2, len(bytesArray.BytesArray))

	slice, err = query.Find()
	assert.NoErr(t, err)
	events := slice.([]*iot.Event)
	if len(events) != 2 {
		t.Fatalf("unexpected size")
	}

	assert.Eq(t, id1, events[0].Id)
	assert.EqString(t, "dev1", events[0].Device)

	assert.Eq(t, id2, events[1].Id)
	assert.EqString(t, "dev2", events[1].Device)

	return
}

func TestQueryBuilder_StringEq(t *testing.T) {
	objectBox := iot.LoadEmptyTestObjectBox()
	defer objectBox.Close()
	box := iot.BoxForEvent(objectBox)
	defer box.Close()
	box.RemoveAll()

	iot.PutEvents(objectBox, 3)

	qb := objectBox.Query(1)
	assert.NoErr(t, qb.Err)
	defer qb.Close()
	qb.StringEq(2, "device 2", false)
	query, err := qb.Build()
	assert.NoErr(t, err)
	defer query.Close()

	slice, err := query.Find()
	assert.NoErr(t, err)
	events := slice.([]*iot.Event)
	assert.EqInt(t, 1, len(events))
	assert.EqString(t, "device 2", events[0].Device)

	query.SetParamString(2, "device 1")
	slice, err = query.Find()
	assert.NoErr(t, err)
	events = slice.([]*iot.Event)
	assert.EqInt(t, 1, len(events))
	assert.EqString(t, "device 1", events[0].Device)
}

func TestQueryBuilder_StringNotEq(t *testing.T) {
	objectBox := iot.LoadEmptyTestObjectBox()
	defer objectBox.Close()
	box := iot.BoxForEvent(objectBox)
	defer box.Close()
	box.RemoveAll()

	iot.PutEvents(objectBox, 3)

	qb := objectBox.Query(1)
	assert.NoErr(t, qb.Err)
	defer qb.Close()
	qb.StringNotEq(2, "device 2", false)
	query, err := qb.Build()
	assert.NoErr(t, err)
	defer query.Close()

	slice, err := query.Find()
	assert.NoErr(t, err)
	events := slice.([]*iot.Event)
	assert.EqInt(t, 2, len(events))

}

func TestQueryBuilder_IntBetween(t *testing.T) {
	objectBox := iot.LoadEmptyTestObjectBox()
	defer objectBox.Close()
	box := iot.BoxForEvent(objectBox)
	defer box.Close()
	box.RemoveAll()

	objectBox.SetDebugFlags(objectbox.DebugFlags_LOG_QUERIES | objectbox.DebugFlags_LOG_QUERY_PARAMETERS)

	events := iot.PutEvents(objectBox, 6)

	qb := objectBox.Query(1)
	assert.NoErr(t, qb.Err)
	defer qb.Close()
	start := events[2].Date
	end := events[4].Date
	qb.IntBetween(3, start, end)
	query, err := qb.Build()
	assert.NoErr(t, err)
	defer query.Close()

	slice, err := query.Find()
	assert.NoErr(t, err)
	events = slice.([]*iot.Event)
	assert.EqInt(t, 3, len(events))
	assert.Eq(t, start, events[0].Date)
	assert.Eq(t, start+1, events[1].Date)
	assert.Eq(t, end, events[2].Date)
}


func TestQueryBuilder_Null(t *testing.T) {
	objectBox := iot.LoadEmptyTestObjectBox()
	defer objectBox.Close()
	box := iot.BoxForEvent(objectBox)
	defer box.Close()
	box.RemoveAll()

	objectBox.SetDebugFlags(objectbox.DebugFlags_LOG_QUERIES | objectbox.DebugFlags_LOG_QUERY_PARAMETERS)

	events := iot.PutEvents(objectBox, 6)

	qb := objectBox.Query(1)
	assert.NoErr(t, qb.Err)
	defer qb.Close()
	qb.Null(3)
	query, err := qb.Build()
	assert.NoErr(t, err)
	defer query.Close()

	slice, err := query.Find()
	assert.NoErr(t, err)
	events = slice.([]*iot.Event)
	assert.EqInt(t, 0, len(events))

}


func TestQueryBuilder_NotNull(t *testing.T) {
	objectBox := iot.LoadEmptyTestObjectBox()
	defer objectBox.Close()
	box := iot.BoxForEvent(objectBox)
	defer box.Close()
	box.RemoveAll()

	objectBox.SetDebugFlags(objectbox.DebugFlags_LOG_QUERIES | objectbox.DebugFlags_LOG_QUERY_PARAMETERS)

	events := iot.PutEvents(objectBox, 6)

	qb := objectBox.Query(1)
	assert.NoErr(t, qb.Err)
	defer qb.Close()
	qb.NotNull(3)
	query, err := qb.Build()
	assert.NoErr(t, err)
	defer query.Close()

	slice, err := query.Find()
	assert.NoErr(t, err)
	events = slice.([]*iot.Event)
	assert.EqInt(t, 6, len(events))

}


func TestQueryBuilder_StringGreater(t *testing.T) {
	objectBox := iot.LoadEmptyTestObjectBox()
	defer objectBox.Close()
	box := iot.BoxForEvent(objectBox)
	defer box.Close()
	box.RemoveAll()

	objectBox.SetDebugFlags(objectbox.DebugFlags_LOG_QUERIES | objectbox.DebugFlags_LOG_QUERY_PARAMETERS)

	events := iot.PutEvents(objectBox, 6)

	qb := objectBox.Query(1)
	assert.NoErr(t, qb.Err)
	defer qb.Close()
	qb.StringGreater(2, "device 2",  false, false)
	query, err := qb.Build()
	assert.NoErr(t, err)
	defer query.Close()

	slice, err := query.Find()
	assert.NoErr(t, err)
	events = slice.([]*iot.Event)
	assert.EqInt(t, 4, len(events))

}

func TestQueryBuilder_IntEqual(t *testing.T) {
	objectBox := iot.LoadEmptyTestObjectBox()
	defer objectBox.Close()
	box := iot.BoxForEvent(objectBox)
	defer box.Close()
	box.RemoveAll()

	objectBox.SetDebugFlags(objectbox.DebugFlags_LOG_QUERIES | objectbox.DebugFlags_LOG_QUERY_PARAMETERS)

	events := iot.PutEvents(objectBox, 6)

	qb := objectBox.Query(1)
	assert.NoErr(t, qb.Err)
	defer qb.Close()
	qb.IntEqual(1, 5)
	query, err := qb.Build()
	assert.NoErr(t, err)
	defer query.Close()

	slice, err := query.Find()
	assert.NoErr(t, err)
	events = slice.([]*iot.Event)
	assert.EqInt(t, 1, len(events))

}

func TestQueryBuilder_DoubleGreater(t *testing.T) {
	objectBox := iot.LoadEmptyTestObjectBox()
	defer objectBox.Close()
	box := iot.BoxForEvent(objectBox)
	defer box.Close()
	box.RemoveAll()

	objectBox.SetDebugFlags(objectbox.DebugFlags_LOG_QUERIES | objectbox.DebugFlags_LOG_QUERY_PARAMETERS)

	//events := iot.PutEvents(objectBox, 6)

	qb := objectBox.Query(1)
	assert.NoErr(t, qb.Err)
	defer qb.Close()

	//FIXME:
	qb.DoubleGreater(1, 2)
	query, err := qb.Build()
	assert.NoErr(t, err)
	defer query.Close()
	//
	//slice, err := query.Find()
	//assert.NoErr(t, err)
	//events = slice.([]*iot.Event)
	//assert.EqInt(t, 1, len(events))

}

func TestQueryBuilder_DoubleBetween(t *testing.T) {
	objectBox := iot.LoadEmptyTestObjectBox()
	defer objectBox.Close()
	box := iot.BoxForEvent(objectBox)
	defer box.Close()
	box.RemoveAll()

	objectBox.SetDebugFlags(objectbox.DebugFlags_LOG_QUERIES | objectbox.DebugFlags_LOG_QUERY_PARAMETERS)

	//events := iot.PutEvents(objectBox, 6)

	qb := objectBox.Query(1)
	assert.NoErr(t, qb.Err)
	defer qb.Close()

	//FIXME:
	qb.DoubleBetween(1, 2, 3)
	query, err := qb.Build()
	assert.NoErr(t, err)
	defer query.Close()
	//
	//slice, err := query.Find()
	//assert.NoErr(t, err)
	//events = slice.([]*iot.Event)
	//assert.EqInt(t, 1, len(events))

}