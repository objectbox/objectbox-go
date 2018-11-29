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

func TestNewQueries(t *testing.T) {
	// XXX TEMPORARY CODE
	var Event_Device = &objectbox.PropertyString{
		Property: &objectbox.Property{
			Id: iot.Event_.Device,
		},
	}

	objectBox := iot.LoadEmptyTestObjectBox()
	defer objectBox.Close()
	box := iot.BoxForEvent(objectBox)
	defer box.Close()

	box.RemoveAll()

	data, err := box.Query(Event_Device.StartsWith("dev", true)).Find()
	assert.NoErr(t, err)
	assert.Eq(t, 0, len(data.([]*iot.Event)))
}

//
//func TestQueryBuilder(t *testing.T) {
//	objectBox := iot.LoadEmptyTestObjectBox()
//	defer objectBox.Close()
//	box := iot.BoxForEvent(objectBox)
//	defer box.Close()
//	box.RemoveAll()
//
//	qb := objectBox.QueryBuilder(iot.EventBinding.Id)
//	assert.NoErr(t, qb.Err)
//	query, err := qb.BuildAndClose()
//	assert.NoErr(t, err)
//	defer query.Close()
//
//	bytesArray, err := query.FindBytes()
//	assert.NoErr(t, err)
//	assert.EqInt(t, 0, len(bytesArray.BytesArray))
//
//	slice, err := query.Find()
//	assert.NoErr(t, err)
//	assert.EqInt(t, 0, len(slice.([]*iot.Event)))
//
//	id1, err := box.Put(&iot.Event{
//		Device: "dev1",
//	})
//	assert.NoErr(t, err)
//
//	id2, err := box.Put(&iot.Event{
//		Device: "dev2",
//	})
//	assert.NoErr(t, err)
//
//	bytesArray, err = query.FindBytes()
//	assert.NoErr(t, err)
//	assert.EqInt(t, 2, len(bytesArray.BytesArray))
//
//	slice, err = query.Find()
//	assert.NoErr(t, err)
//	events := slice.([]*iot.Event)
//	if len(events) != 2 {
//		t.Fatalf("unexpected size")
//	}
//
//	assert.Eq(t, id1, events[0].Id)
//	assert.EqString(t, "dev1", events[0].Device)
//
//	assert.Eq(t, id2, events[1].Id)
//	assert.EqString(t, "dev2", events[1].Device)
//
//	return
//}
//
//// TODO refactor following conditions to table-like test
//
//func TestQueryBuilder_StringEq(t *testing.T) {
//	objectBox := iot.LoadEmptyTestObjectBox()
//	defer objectBox.Close()
//	box := iot.BoxForEvent(objectBox)
//	defer box.Close()
//	box.RemoveAll()
//
//	iot.PutEvents(objectBox, 3)
//
//	qb := objectBox.QueryBuilder(iot.EventBinding.Id)
//	assert.NoErr(t, qb.Err)
//	defer qb.Close()
//	qb.StringEq(iot.Event_.Device, "device 2", false)
//	query, err := qb.Build()
//	assert.NoErr(t, err)
//	defer query.Close()
//
//	slice, err := query.Find()
//	assert.NoErr(t, err)
//	events := slice.([]*iot.Event)
//	assert.EqInt(t, 1, len(events))
//	assert.EqString(t, "device 2", events[0].Device)
//}
//
//func TestQueryBuilder_StringIn(t *testing.T) {
//	objectBox := iot.LoadEmptyTestObjectBox()
//	defer objectBox.Close()
//	box := iot.BoxForEvent(objectBox)
//	defer box.Close()
//	box.RemoveAll()
//
//	iot.PutEvents(objectBox, 3)
//
//	qb := objectBox.QueryBuilder(iot.EventBinding.Id)
//	assert.NoErr(t, qb.Err)
//	defer qb.Close()
//	values := []string{"device 2", "device 3"}
//	qb.StringIn(iot.Event_.Device, values, false)
//	query, err := qb.Build()
//	assert.NoErr(t, err)
//	defer query.Close()
//
//	slice, err := query.Find()
//	assert.NoErr(t, err)
//	events := slice.([]*iot.Event)
//	assert.EqInt(t, len(values), len(events))
//}
//
//func TestQueryBuilder_StringContains(t *testing.T) {
//	objectBox := iot.LoadEmptyTestObjectBox()
//	defer objectBox.Close()
//	box := iot.BoxForEvent(objectBox)
//	defer box.Close()
//	box.RemoveAll()
//
//	iot.PutEvents(objectBox, 3)
//
//	qb := objectBox.QueryBuilder(iot.EventBinding.Id)
//	assert.NoErr(t, qb.Err)
//	defer qb.Close()
//	qb.StringContains(iot.Event_.Device, "device 2", false)
//	query, err := qb.Build()
//	assert.NoErr(t, err)
//	defer query.Close()
//
//	slice, err := query.Find()
//	assert.NoErr(t, err)
//	events := slice.([]*iot.Event)
//	assert.EqInt(t, 1, len(events))
//	assert.EqString(t, "device 2", events[0].Device)
//}
//
//func TestQueryBuilder_StringStartsWith(t *testing.T) {
//	objectBox := iot.LoadEmptyTestObjectBox()
//	defer objectBox.Close()
//	box := iot.BoxForEvent(objectBox)
//	defer box.Close()
//	box.RemoveAll()
//
//	iot.PutEvents(objectBox, 3)
//
//	qb := objectBox.QueryBuilder(iot.EventBinding.Id)
//	assert.NoErr(t, qb.Err)
//	defer qb.Close()
//	qb.StringStartsWith(iot.Event_.Device, "device 2", false)
//	query, err := qb.Build()
//	assert.NoErr(t, err)
//	defer query.Close()
//
//	slice, err := query.Find()
//	assert.NoErr(t, err)
//	events := slice.([]*iot.Event)
//	assert.EqInt(t, 1, len(events))
//	assert.EqString(t, "device 2", events[0].Device)
//}
//
//func TestQueryBuilder_StringEndsWith(t *testing.T) {
//	objectBox := iot.LoadEmptyTestObjectBox()
//	defer objectBox.Close()
//	box := iot.BoxForEvent(objectBox)
//	defer box.Close()
//	box.RemoveAll()
//
//	iot.PutEvents(objectBox, 3)
//
//	qb := objectBox.QueryBuilder(iot.EventBinding.Id)
//	assert.NoErr(t, qb.Err)
//	defer qb.Close()
//	qb.StringEndsWith(iot.Event_.Device, "device 2", false)
//	query, err := qb.Build()
//	assert.NoErr(t, err)
//	defer query.Close()
//
//	slice, err := query.Find()
//	assert.NoErr(t, err)
//	events := slice.([]*iot.Event)
//	assert.EqInt(t, 1, len(events))
//	assert.EqString(t, "device 2", events[0].Device)
//}
//
//func TestQueryBuilder_StringNotEq(t *testing.T) {
//	objectBox := iot.LoadEmptyTestObjectBox()
//	defer objectBox.Close()
//	box := iot.BoxForEvent(objectBox)
//	defer box.Close()
//	box.RemoveAll()
//
//	iot.PutEvents(objectBox, 3)
//
//	qb := objectBox.QueryBuilder(iot.EventBinding.Id)
//	assert.NoErr(t, qb.Err)
//	defer qb.Close()
//	qb.StringNotEq(iot.Event_.Device, "device 3", false)
//	query, err := qb.Build()
//	assert.NoErr(t, err)
//	defer query.Close()
//
//	slice, err := query.Find()
//	assert.NoErr(t, err)
//	events := slice.([]*iot.Event)
//	assert.EqInt(t, 2, len(events))
//
//}
//
//func TestQueryBuilder_StringLess(t *testing.T) {
//	objectBox := iot.LoadEmptyTestObjectBox()
//	defer objectBox.Close()
//	box := iot.BoxForEvent(objectBox)
//	defer box.Close()
//	box.RemoveAll()
//
//	iot.PutEvents(objectBox, 3)
//
//	qb := objectBox.QueryBuilder(iot.EventBinding.Id)
//	assert.NoErr(t, qb.Err)
//	defer qb.Close()
//	qb.StringLess(iot.Event_.Device, "device 3", false, false)
//	query, err := qb.Build()
//	assert.NoErr(t, err)
//	defer query.Close()
//
//	slice, err := query.Find()
//	assert.NoErr(t, err)
//	events := slice.([]*iot.Event)
//	assert.EqInt(t, 2, len(events))
//
//}
//
//func TestQueryBuilder_IntBetween(t *testing.T) {
//	objectBox := iot.LoadEmptyTestObjectBox()
//	defer objectBox.Close()
//	box := iot.BoxForEvent(objectBox)
//	defer box.Close()
//	box.RemoveAll()
//
//	objectBox.SetDebugFlags(objectbox.DebugFlags_LOG_QUERIES | objectbox.DebugFlags_LOG_QUERY_PARAMETERS)
//
//	events := iot.PutEvents(objectBox, 6)
//
//	qb := objectBox.QueryBuilder(iot.EventBinding.Id)
//	assert.NoErr(t, qb.Err)
//	defer qb.Close()
//	start := events[2].Date
//	end := events[4].Date
//	qb.IntBetween(iot.Event_.Date, start, end)
//	query, err := qb.Build()
//	assert.NoErr(t, err)
//	defer query.Close()
//
//	slice, err := query.Find()
//	assert.NoErr(t, err)
//	events = slice.([]*iot.Event)
//	assert.EqInt(t, 3, len(events))
//	assert.Eq(t, start, events[0].Date)
//	assert.Eq(t, start+1, events[1].Date)
//	assert.Eq(t, end, events[2].Date)
//}
//
//func TestQueryBuilder_Null(t *testing.T) {
//	objectBox := iot.LoadEmptyTestObjectBox()
//	defer objectBox.Close()
//	box := iot.BoxForEvent(objectBox)
//	defer box.Close()
//	box.RemoveAll()
//
//	objectBox.SetDebugFlags(objectbox.DebugFlags_LOG_QUERIES | objectbox.DebugFlags_LOG_QUERY_PARAMETERS)
//
//	events := iot.PutEvents(objectBox, 6)
//
//	qb := objectBox.QueryBuilder(iot.EventBinding.Id)
//	assert.NoErr(t, qb.Err)
//	defer qb.Close()
//	qb.Null(iot.Event_.Date)
//	query, err := qb.Build()
//	assert.NoErr(t, err)
//	defer query.Close()
//
//	slice, err := query.Find()
//	assert.NoErr(t, err)
//	events = slice.([]*iot.Event)
//	assert.EqInt(t, 0, len(events))
//
//}
//
//func TestQueryBuilder_NotNull(t *testing.T) {
//	objectBox := iot.LoadEmptyTestObjectBox()
//	defer objectBox.Close()
//	box := iot.BoxForEvent(objectBox)
//	defer box.Close()
//	box.RemoveAll()
//
//	objectBox.SetDebugFlags(objectbox.DebugFlags_LOG_QUERIES | objectbox.DebugFlags_LOG_QUERY_PARAMETERS)
//
//	events := iot.PutEvents(objectBox, 6)
//
//	qb := objectBox.QueryBuilder(iot.EventBinding.Id)
//	assert.NoErr(t, qb.Err)
//	defer qb.Close()
//	qb.NotNull(iot.Event_.Date)
//	query, err := qb.Build()
//	assert.NoErr(t, err)
//	defer query.Close()
//
//	slice, err := query.Find()
//	assert.NoErr(t, err)
//	events = slice.([]*iot.Event)
//	assert.EqInt(t, 6, len(events))
//
//}
//
//func TestQueryBuilder_StringGreater(t *testing.T) {
//	objectBox := iot.LoadEmptyTestObjectBox()
//	defer objectBox.Close()
//	box := iot.BoxForEvent(objectBox)
//	defer box.Close()
//	box.RemoveAll()
//
//	objectBox.SetDebugFlags(objectbox.DebugFlags_LOG_QUERIES | objectbox.DebugFlags_LOG_QUERY_PARAMETERS)
//
//	events := iot.PutEvents(objectBox, 6)
//
//	qb := objectBox.QueryBuilder(iot.EventBinding.Id)
//	assert.NoErr(t, qb.Err)
//	defer qb.Close()
//	qb.StringGreater(iot.Event_.Device, "device 2", false, false)
//	query, err := qb.Build()
//	assert.NoErr(t, err)
//	defer query.Close()
//
//	slice, err := query.Find()
//	assert.NoErr(t, err)
//	events = slice.([]*iot.Event)
//	assert.EqInt(t, 4, len(events))
//
//}
//
//func TestQueryBuilder_IntEqual(t *testing.T) {
//	objectBox := iot.LoadEmptyTestObjectBox()
//	defer objectBox.Close()
//	box := iot.BoxForEvent(objectBox)
//	defer box.Close()
//	box.RemoveAll()
//
//	objectBox.SetDebugFlags(objectbox.DebugFlags_LOG_QUERIES | objectbox.DebugFlags_LOG_QUERY_PARAMETERS)
//
//	events := iot.PutEvents(objectBox, 6)
//
//	qb := objectBox.QueryBuilder(iot.EventBinding.Id)
//	assert.NoErr(t, qb.Err)
//	defer qb.Close()
//	qb.IntEqual(iot.Event_.Id, 5)
//	query, err := qb.Build()
//	assert.NoErr(t, err)
//	defer query.Close()
//
//	slice, err := query.Find()
//	assert.NoErr(t, err)
//	events = slice.([]*iot.Event)
//	assert.EqInt(t, 1, len(events))
//
//}
//
//func TestQueryBuilder_IntNotEqual(t *testing.T) {
//	objectBox := iot.LoadEmptyTestObjectBox()
//	defer objectBox.Close()
//	box := iot.BoxForEvent(objectBox)
//	defer box.Close()
//	box.RemoveAll()
//
//	objectBox.SetDebugFlags(objectbox.DebugFlags_LOG_QUERIES | objectbox.DebugFlags_LOG_QUERY_PARAMETERS)
//
//	events := iot.PutEvents(objectBox, 6)
//
//	qb := objectBox.QueryBuilder(iot.EventBinding.Id)
//	assert.NoErr(t, qb.Err)
//	defer qb.Close()
//	qb.IntNotEqual(iot.Event_.Id, 5)
//	query, err := qb.Build()
//	assert.NoErr(t, err)
//	defer query.Close()
//
//	slice, err := query.Find()
//	assert.NoErr(t, err)
//	events = slice.([]*iot.Event)
//	assert.EqInt(t, 5, len(events))
//
//}
//
//func TestQueryBuilder_IntGreater(t *testing.T) {
//	objectBox := iot.LoadEmptyTestObjectBox()
//	defer objectBox.Close()
//	box := iot.BoxForEvent(objectBox)
//	defer box.Close()
//	box.RemoveAll()
//
//	objectBox.SetDebugFlags(objectbox.DebugFlags_LOG_QUERIES | objectbox.DebugFlags_LOG_QUERY_PARAMETERS)
//
//	events := iot.PutEvents(objectBox, 6)
//
//	qb := objectBox.QueryBuilder(iot.EventBinding.Id)
//	assert.NoErr(t, qb.Err)
//	defer qb.Close()
//	qb.IntGreater(iot.Event_.Id, 5)
//	query, err := qb.Build()
//	assert.NoErr(t, err)
//	defer query.Close()
//
//	slice, err := query.Find()
//	assert.NoErr(t, err)
//	events = slice.([]*iot.Event)
//	assert.EqInt(t, 1, len(events))
//
//}
//
//func TestQueryBuilder_IntLess(t *testing.T) {
//	objectBox := iot.LoadEmptyTestObjectBox()
//	defer objectBox.Close()
//	box := iot.BoxForEvent(objectBox)
//	defer box.Close()
//	box.RemoveAll()
//
//	objectBox.SetDebugFlags(objectbox.DebugFlags_LOG_QUERIES | objectbox.DebugFlags_LOG_QUERY_PARAMETERS)
//
//	events := iot.PutEvents(objectBox, 6)
//
//	qb := objectBox.QueryBuilder(iot.EventBinding.Id)
//	assert.NoErr(t, qb.Err)
//	defer qb.Close()
//	qb.IntLess(iot.Event_.Id, 5)
//	query, err := qb.Build()
//	assert.NoErr(t, err)
//	defer query.Close()
//
//	slice, err := query.Find()
//	assert.NoErr(t, err)
//	events = slice.([]*iot.Event)
//	assert.EqInt(t, 4, len(events))
//
//}
//
//func TestQueryBuilder_DoubleLess(t *testing.T) {
//	objectBox := iot.LoadEmptyTestObjectBox()
//	defer objectBox.Close()
//	box := iot.BoxForReading(objectBox)
//	defer box.Close()
//	box.RemoveAll()
//
//	objectBox.SetDebugFlags(objectbox.DebugFlags_LOG_QUERIES | objectbox.DebugFlags_LOG_QUERY_PARAMETERS)
//
//	readings := iot.PutReadings(objectBox, 6)
//
//	qb := objectBox.QueryBuilder(iot.ReadingBinding.Id)
//	assert.NoErr(t, qb.Err)
//	defer qb.Close()
//
//	qb.DoubleLess(iot.Reading_.ValueFloating32, 10003)
//	query, err := qb.Build()
//	assert.NoErr(t, err)
//	defer query.Close()
//
//	slice, err := query.Find()
//	assert.NoErr(t, err)
//	readings = slice.([]*iot.Reading)
//	assert.EqInt(t, 2, len(readings))
//}
//
//func TestQueryBuilder_DoubleGreater(t *testing.T) {
//	objectBox := iot.LoadEmptyTestObjectBox()
//	defer objectBox.Close()
//	box := iot.BoxForReading(objectBox)
//	defer box.Close()
//	box.RemoveAll()
//
//	objectBox.SetDebugFlags(objectbox.DebugFlags_LOG_QUERIES | objectbox.DebugFlags_LOG_QUERY_PARAMETERS)
//
//	readings := iot.PutReadings(objectBox, 6)
//
//	qb := objectBox.QueryBuilder(iot.ReadingBinding.Id)
//	assert.NoErr(t, qb.Err)
//	defer qb.Close()
//
//	qb.DoubleGreater(iot.Reading_.ValueFloating32, 10003)
//	query, err := qb.Build()
//	assert.NoErr(t, err)
//	defer query.Close()
//
//	slice, err := query.Find()
//	assert.NoErr(t, err)
//	readings = slice.([]*iot.Reading)
//	assert.EqInt(t, 3, len(readings))
//}
//
//func TestQueryBuilder_DoubleBetween(t *testing.T) {
//	objectBox := iot.LoadEmptyTestObjectBox()
//	defer objectBox.Close()
//	box := iot.BoxForReading(objectBox)
//	defer box.Close()
//	box.RemoveAll()
//
//	objectBox.SetDebugFlags(objectbox.DebugFlags_LOG_QUERIES | objectbox.DebugFlags_LOG_QUERY_PARAMETERS)
//
//	readings := iot.PutReadings(objectBox, 6)
//
//	qb := objectBox.QueryBuilder(iot.ReadingBinding.Id)
//	assert.NoErr(t, qb.Err)
//	defer qb.Close()
//
//	qb.DoubleBetween(iot.Reading_.ValueFloating32, 10003, 10005)
//	query, err := qb.Build()
//	assert.NoErr(t, err)
//	defer query.Close()
//
//	slice, err := query.Find()
//	assert.NoErr(t, err)
//	readings = slice.([]*iot.Reading)
//	assert.EqInt(t, 3, len(readings))
//}
//
//func TestQueryBuilder_BytesEqual(t *testing.T) {
//	objectBox := iot.LoadEmptyTestObjectBox()
//	defer objectBox.Close()
//	box := iot.BoxForEvent(objectBox)
//	defer box.Close()
//	box.RemoveAll()
//
//	objectBox.SetDebugFlags(objectbox.DebugFlags_LOG_QUERIES | objectbox.DebugFlags_LOG_QUERY_PARAMETERS)
//
//	events := iot.PutEvents(objectBox, 6)
//
//	qb := objectBox.QueryBuilder(iot.EventBinding.Id)
//	assert.NoErr(t, qb.Err)
//	defer qb.Close()
//
//	bytes := []byte{1, 2, 3}
//	qb.BytesEqual(iot.Event_.Picture, bytes)
//	query, err := qb.Build()
//	assert.NoErr(t, err)
//	defer query.Close()
//
//	slice, err := query.Find()
//	assert.NoErr(t, err)
//	events = slice.([]*iot.Event)
//	assert.EqInt(t, 0, len(events))
//
//}
//
//func TestQueryBuilder_BytesGreater(t *testing.T) {
//	objectBox := iot.LoadEmptyTestObjectBox()
//	defer objectBox.Close()
//	box := iot.BoxForEvent(objectBox)
//	defer box.Close()
//	box.RemoveAll()
//
//	objectBox.SetDebugFlags(objectbox.DebugFlags_LOG_QUERIES | objectbox.DebugFlags_LOG_QUERY_PARAMETERS)
//
//	events := iot.PutEvents(objectBox, 6)
//
//	qb := objectBox.QueryBuilder(iot.EventBinding.Id)
//	assert.NoErr(t, qb.Err)
//	defer qb.Close()
//
//	bytes := []byte{1, 2, 3}
//	qb.BytesGreater(iot.Event_.Picture, bytes, false)
//	query, err := qb.Build()
//	assert.NoErr(t, err)
//	defer query.Close()
//
//	slice, err := query.Find()
//	assert.NoErr(t, err)
//	events = slice.([]*iot.Event)
//	assert.EqInt(t, 0, len(events))
//
//}
//
//func TestQueryBuilder_BytesLess(t *testing.T) {
//	objectBox := iot.LoadEmptyTestObjectBox()
//	defer objectBox.Close()
//	box := iot.BoxForEvent(objectBox)
//	defer box.Close()
//	box.RemoveAll()
//
//	objectBox.SetDebugFlags(objectbox.DebugFlags_LOG_QUERIES | objectbox.DebugFlags_LOG_QUERY_PARAMETERS)
//
//	events := iot.PutEvents(objectBox, 6)
//
//	qb := objectBox.QueryBuilder(iot.EventBinding.Id)
//	assert.NoErr(t, qb.Err)
//	defer qb.Close()
//
//	bytes := []byte{1, 2, 3}
//	qb.BytesLess(iot.Event_.Picture, bytes, false)
//	query, err := qb.Build()
//	assert.NoErr(t, err)
//	defer query.Close()
//
//	slice, err := query.Find()
//	assert.NoErr(t, err)
//	events = slice.([]*iot.Event)
//	assert.EqInt(t, 0, len(events))
//
//}
//
//func TestQueryBuilder_Int64In(t *testing.T) {
//	objectBox := iot.LoadEmptyTestObjectBox()
//	defer objectBox.Close()
//	box := iot.BoxForReading(objectBox)
//	defer box.Close()
//	box.RemoveAll()
//
//	objectBox.SetDebugFlags(objectbox.DebugFlags_LOG_QUERIES | objectbox.DebugFlags_LOG_QUERY_PARAMETERS)
//
//	readings := iot.PutReadings(objectBox, 6)
//
//	qb := objectBox.QueryBuilder(iot.ReadingBinding.Id)
//	assert.NoErr(t, qb.Err)
//	defer qb.Close()
//
//	values := []int64{10002, 10003}
//	qb.Int64In(iot.Reading_.ValueInteger, values)
//	query, err := qb.Build()
//	assert.NoErr(t, err)
//	defer query.Close()
//
//	slice, err := query.Find()
//	assert.NoErr(t, err)
//	readings = slice.([]*iot.Reading)
//	assert.EqInt(t, 2, len(readings))
//
//}
//
//func TestQueryBuilder_Int64NotIn(t *testing.T) {
//	objectBox := iot.LoadEmptyTestObjectBox()
//	defer objectBox.Close()
//	box := iot.BoxForReading(objectBox)
//	defer box.Close()
//	box.RemoveAll()
//
//	objectBox.SetDebugFlags(objectbox.DebugFlags_LOG_QUERIES | objectbox.DebugFlags_LOG_QUERY_PARAMETERS)
//
//	readings :=
//		iot.PutReadings(objectBox, 6)
//
//	qb := objectBox.QueryBuilder(iot.ReadingBinding.Id)
//	assert.NoErr(t, qb.Err)
//	defer qb.Close()
//
//	values := []int64{10002, 10003}
//	qb.Int64NotIn(iot.Reading_.ValueInteger, values)
//	query, err := qb.Build()
//	assert.NoErr(t, err)
//	defer query.Close()
//
//	slice, err := query.Find()
//	assert.NoErr(t, err)
//	readings = slice.([]*iot.Reading)
//	assert.EqInt(t, 4, len(readings))
//
//}
//
//func TestQueryBuilder_Int32In(t *testing.T) {
//	objectBox := iot.LoadEmptyTestObjectBox()
//	defer objectBox.Close()
//	box := iot.BoxForReading(objectBox)
//	defer box.Close()
//	box.RemoveAll()
//
//	objectBox.SetDebugFlags(objectbox.DebugFlags_LOG_QUERIES | objectbox.DebugFlags_LOG_QUERY_PARAMETERS)
//
//	readings :=
//		iot.PutReadings(objectBox, 6)
//
//	qb := objectBox.QueryBuilder(iot.ReadingBinding.Id)
//	assert.NoErr(t, qb.Err)
//	defer qb.Close()
//
//	values := []int32{10002}
//	qb.Int32In(iot.Reading_.ValueInt32, values)
//	query, err := qb.Build()
//	assert.NoErr(t, err)
//	defer query.Close()
//
//	slice, err := query.Find()
//	assert.NoErr(t, err)
//	readings = slice.([]*iot.Reading)
//	assert.EqInt(t, 1, len(readings))
//}
//
//func TestQueryBuilder_Int32NotIn(t *testing.T) {
//	objectBox := iot.LoadEmptyTestObjectBox()
//	defer objectBox.Close()
//	box := iot.BoxForReading(objectBox)
//	defer box.Close()
//	box.RemoveAll()
//
//	objectBox.SetDebugFlags(objectbox.DebugFlags_LOG_QUERIES | objectbox.DebugFlags_LOG_QUERY_PARAMETERS)
//
//	readings :=
//		iot.PutReadings(objectBox, 6)
//
//	qb := objectBox.QueryBuilder(iot.ReadingBinding.Id)
//	assert.NoErr(t, qb.Err)
//	defer qb.Close()
//
//	values := []int32{10002}
//	qb.Int32NotIn(iot.Reading_.ValueInt32, values)
//	query, err := qb.Build()
//	assert.NoErr(t, err)
//	defer query.Close()
//
//	slice, err := query.Find()
//	assert.NoErr(t, err)
//	readings = slice.([]*iot.Reading)
//	assert.EqInt(t, 5, len(readings))
//}
