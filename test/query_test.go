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

	"github.com/objectbox/objectbox-go/test/model"

	"github.com/objectbox/objectbox-go/objectbox"
	"github.com/objectbox/objectbox-go/test/assert"
)

func TestQueryConditions(t *testing.T) {
	env := model.NewTestEnv(t)
	defer env.Close()

	var box = env.Box

	// let's alias the entity to make the test cases easier to read
	var E = model.Entity_

	// define some data for the tests
	var e = model.Entity{
		Int:        42,
		Int8:       42,
		Int16:      42,
		Int32:      42,
		Int64:      42,
		Uint:       42,
		Uint8:      42,
		Uint16:     42,
		Uint32:     42,
		Uint64:     42,
		Bool:       true,
		String:     "val",
		Byte:       42,
		ByteVector: []byte{1, 1, 2, 3, 5, 8, 13},
		Rune:       42,
		Float32:    42.24,
		Float64:    42.24,
	}

	testCases := []struct {
		expected string
		query    *objectbox.Query
	}{
		{`String == "val"`, box.Query(E.String.Equal(e.String, true))},
		{`String == "val"`, box.Query(E.String.Equal(e.String, false))},
		{`String != "val"`, box.Query(E.String.NotEqual(e.String, true))},
		{`String != "val"`, box.Query(E.String.NotEqual(e.String, false))},
		{`String contains "val"`, box.Query(E.String.Contains(e.String, true))},
		{`String contains "val"`, box.Query(E.String.Contains(e.String, false))},
		{`String starts with "val"`, box.Query(E.String.StartsWith(e.String, true))},
		{`String starts with "val"`, box.Query(E.String.StartsWith(e.String, false))},
		{`String ends with "val"`, box.Query(E.String.EndsWith(e.String, true))},
		{`String ends with "val"`, box.Query(E.String.EndsWith(e.String, false))},
		{`String > "val"`, box.Query(E.String.GreaterThan(e.String, true))},
		{`String > "val"`, box.Query(E.String.GreaterThan(e.String, false))},
		{`String >= "val"`, box.Query(E.String.GreaterOrEqual(e.String, true))},
		{`String >= "val"`, box.Query(E.String.GreaterOrEqual(e.String, false))},
		{`String < "val"`, box.Query(E.String.LessThan(e.String, true))},
		{`String < "val"`, box.Query(E.String.LessThan(e.String, false))},
		{`String <= "val"`, box.Query(E.String.LessOrEqual(e.String, true))},
		{`String <= "val"`, box.Query(E.String.LessOrEqual(e.String, false))},
		{`String in ["val1", "val2"]`, box.Query(E.String.In(true, "val1", "val2"))},
		{`String in ["val1", "val2"]`, box.Query(E.String.In(false, "val1", "val2"))},

		{`Int64 == 42`, box.Query(E.Int64.Equal(e.Int64))},
		{`Int64 != 42`, box.Query(E.Int64.NotEqual(e.Int64))},
		{`Int64 > 42`, box.Query(E.Int64.GreaterThan(e.Int64))},
		{`Int64 < 42`, box.Query(E.Int64.LessThan(e.Int64))},
		{`Int64 between 42 and 84`, box.Query(E.Int64.Between(e.Int64, e.Int64*2))},
		{`Int64 in [84|42]`, box.Query(E.Int64.In(e.Int64, e.Int64*2))},
		{`Int64 in [84|42]`, box.Query(E.Int64.NotIn(e.Int64, e.Int64*2))},

		{`Uint64 == 42`, box.Query(E.Uint64.Equal(e.Uint64))},
		{`Uint64 != 42`, box.Query(E.Uint64.NotEqual(e.Uint64))},
		//{`Uint64 > 42`, box.Query(E.Uint64.GreaterThan(e.Uint64))},
		//{`Uint64 < 42`, box.Query(E.Uint64.LessThan(e.Uint64))},
		//{`Uint64 between 42 and 84`, box.Query(E.Uint64.Between(e.Uint64, e.Uint64*2))},
		{`Uint64 in [84|42]`, box.Query(E.Uint64.In(e.Uint64, e.Uint64*2))},
		{`Uint64 in [84|42]`, box.Query(E.Uint64.NotIn(e.Uint64, e.Uint64*2))},

		{`Int == 42`, box.Query(E.Int.Equal(e.Int))},
		{`Int != 42`, box.Query(E.Int.NotEqual(e.Int))},
		{`Int > 42`, box.Query(E.Int.GreaterThan(e.Int))},
		{`Int < 42`, box.Query(E.Int.LessThan(e.Int))},
		{`Int between 42 and 84`, box.Query(E.Int.Between(e.Int, e.Int*2))},
		{`Int in [84|42]`, box.Query(E.Int.In(e.Int, e.Int*2))},
		{`Int in [84|42]`, box.Query(E.Int.NotIn(e.Int, e.Int*2))},

		{`Uint == 42`, box.Query(E.Uint.Equal(e.Uint))},
		{`Uint != 42`, box.Query(E.Uint.NotEqual(e.Uint))},
		//{`Uint > 42`, box.Query(E.Uint.GreaterThan(e.Uint))},
		//{`Uint < 42`, box.Query(E.Uint.LessThan(e.Uint))},
		//{`Uint between 42 and 84`, box.Query(E.Uint.Between(e.Uint, e.Uint*2))},
		{`Uint in [84|42]`, box.Query(E.Uint.In(e.Uint, e.Uint*2))},
		{`Uint in [84|42]`, box.Query(E.Uint.NotIn(e.Uint, e.Uint*2))},

		{`Rune == 42`, box.Query(E.Rune.Equal(e.Rune))},
		{`Rune != 42`, box.Query(E.Rune.NotEqual(e.Rune))},
		{`Rune > 42`, box.Query(E.Rune.GreaterThan(e.Rune))},
		{`Rune < 42`, box.Query(E.Rune.LessThan(e.Rune))},
		{`Rune between 42 and 84`, box.Query(E.Rune.Between(e.Rune, e.Rune*2))},
		{`Rune in [84|42]`, box.Query(E.Rune.In(e.Rune, e.Rune*2))},
		{`Rune in [84|42]`, box.Query(E.Rune.NotIn(e.Rune, e.Rune*2))},

		{`Int32 == 42`, box.Query(E.Int32.Equal(e.Int32))},
		{`Int32 != 42`, box.Query(E.Int32.NotEqual(e.Int32))},
		{`Int32 > 42`, box.Query(E.Int32.GreaterThan(e.Int32))},
		{`Int32 < 42`, box.Query(E.Int32.LessThan(e.Int32))},
		{`Int32 between 42 and 84`, box.Query(E.Int32.Between(e.Int32, e.Int32*2))},
		{`Int32 in [84|42]`, box.Query(E.Int32.In(e.Int32, e.Int32*2))},
		{`Int32 in [84|42]`, box.Query(E.Int32.NotIn(e.Int32, e.Int32*2))},

		{`Uint32 == 42`, box.Query(E.Uint32.Equal(e.Uint32))},
		{`Uint32 != 42`, box.Query(E.Uint32.NotEqual(e.Uint32))},
		{`Uint32 > 42`, box.Query(E.Uint32.GreaterThan(e.Uint32))},
		{`Uint32 < 42`, box.Query(E.Uint32.LessThan(e.Uint32))},
		{`Uint32 between 42 and 84`, box.Query(E.Uint32.Between(e.Uint32, e.Uint32*2))},
		{`Uint32 in [84|42]`, box.Query(E.Uint32.In(e.Uint32, e.Uint32*2))},
		{`Uint32 in [84|42]`, box.Query(E.Uint32.NotIn(e.Uint32, e.Uint32*2))},

		{`Int16 == 42`, box.Query(E.Int16.Equal(e.Int16))},
		{`Int16 != 42`, box.Query(E.Int16.NotEqual(e.Int16))},
		{`Int16 > 42`, box.Query(E.Int16.GreaterThan(e.Int16))},
		{`Int16 < 42`, box.Query(E.Int16.LessThan(e.Int16))},
		{`Int16 between 42 and 84`, box.Query(E.Int16.Between(e.Int16, e.Int16*2))},
		//{`Int16 in [84|42]`, box.Query(E.Int16.In(e.Int16, e.Int16*2))},
		//{`Int16 in [84|42]`, box.Query(E.Int16.NotIn(e.Int16, e.Int16*2))},

		{`Uint16 == 42`, box.Query(E.Uint16.Equal(e.Uint16))},
		{`Uint16 != 42`, box.Query(E.Uint16.NotEqual(e.Uint16))},
		{`Uint16 > 42`, box.Query(E.Uint16.GreaterThan(e.Uint16))},
		{`Uint16 < 42`, box.Query(E.Uint16.LessThan(e.Uint16))},
		{`Uint16 between 42 and 84`, box.Query(E.Uint16.Between(e.Uint16, e.Uint16*2))},
		//{`Uint16 in [84|42]`, box.Query(E.Uint16.In(e.Uint16, e.Uint16*2))},
		//{`Uint16 in [84|42]`, box.Query(E.Uint16.NotIn(e.Uint16, e.Uint16*2))},

		{`Int8 == 42`, box.Query(E.Int8.Equal(e.Int8))},
		{`Int8 != 42`, box.Query(E.Int8.NotEqual(e.Int8))},
		{`Int8 > 42`, box.Query(E.Int8.GreaterThan(e.Int8))},
		{`Int8 < 42`, box.Query(E.Int8.LessThan(e.Int8))},
		{`Int8 between 42 and 84`, box.Query(E.Int8.Between(e.Int8, e.Int8*2))},
		//{`Int8 in [84|42]`, box.Query(E.Int8.In(e.Int8, e.Int8*2))},
		//{`Int8 in [84|42]`, box.Query(E.Int8.NotIn(e.Int8, e.Int8*2))},

		{`Uint8 == 42`, box.Query(E.Uint8.Equal(e.Uint8))},
		{`Uint8 != 42`, box.Query(E.Uint8.NotEqual(e.Uint8))},
		{`Uint8 > 42`, box.Query(E.Uint8.GreaterThan(e.Uint8))},
		{`Uint8 < 42`, box.Query(E.Uint8.LessThan(e.Uint8))},
		{`Uint8 between 42 and 84`, box.Query(E.Uint8.Between(e.Uint8, e.Uint8*2))},
		//{`Uint8 in [84|42]`, box.Query(E.Uint8.In(e.Uint8, e.Uint8*2))},
		//{`Uint8 in [84|42]`, box.Query(E.Uint8.NotIn(e.Uint8, e.Uint8*2))},

		{`Byte == 42`, box.Query(E.Byte.Equal(e.Byte))},
		{`Byte != 42`, box.Query(E.Byte.NotEqual(e.Byte))},
		{`Byte > 42`, box.Query(E.Byte.GreaterThan(e.Byte))},
		{`Byte < 42`, box.Query(E.Byte.LessThan(e.Byte))},
		{`Byte between 42 and 84`, box.Query(E.Byte.Between(e.Byte, e.Byte*2))},
		//{`Byte in [84|42]`, box.Query(E.Byte.In(e.Byte, e.Byte*2))},
		//{`Byte in [84|42]`, box.Query(E.Byte.NotIn(e.Byte, e.Byte*2))},

		{`Float64 > 42.240000`, box.Query(E.Float64.GreaterThan(e.Float64))},
		{`Float64 < 42.240000`, box.Query(E.Float64.LessThan(e.Float64))},
		{`Float64 between 42.240000 and 84.480000`, box.Query(E.Float64.Between(e.Float64, e.Float64*2))},

		{`Float32 > 42.240002`, box.Query(E.Float32.GreaterThan(e.Float32))},
		{`Float32 < 42.240002`, box.Query(E.Float32.LessThan(e.Float32))},
		{`Float32 between 42.240002 and 84.480003`, box.Query(E.Float32.Between(e.Float32, e.Float32*2))},

		{`ByteVector == `, box.Query(E.ByteVector.Equal(e.ByteVector))},
		//{`ByteVector != `, box.Query(E.ByteVector.NotEqual(e.ByteVector))},
		{`ByteVector >= `, box.Query(E.ByteVector.GreaterThan(e.ByteVector))},
		{`ByteVector >= `, box.Query(E.ByteVector.GreaterOrEqual(e.ByteVector))},
		{`ByteVector <= `, box.Query(E.ByteVector.LessThan(e.ByteVector))},
		{`ByteVector <= `, box.Query(E.ByteVector.LessOrEqual(e.ByteVector))},

		{`Bool == 1`, box.Query(E.Bool.Equal(e.Bool))},
	}

	t.Logf("Executing %d test cases", len(testCases))

	for i, tc := range testCases {
		if desc, err := tc.query.Describe(); err != nil {
			assert.Failf(t, "error describing %d {%s} - %s", i, tc.expected, err)
		} else {
			assert.Eq(t, tc.expected, desc)
		}

		if data, err := tc.query.Find(); err != nil {
			assert.Failf(t, "error executing %d {%s} - %s", i, tc.expected, err)
		} else if data == nil {
			assert.Failf(t, "error executing %d {%s} - data is nil", i, tc.expected)
		}
	}
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
//	qb.StringEqual(iot.Event_.Device, "device 2", false)
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
//	qb.StringNotEqual(iot.Event_.Device, "device 3", false)
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
