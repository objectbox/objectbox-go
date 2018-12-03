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
	"fmt"
	"runtime"
	"strings"
	"testing"

	"github.com/objectbox/objectbox-go/test/assert"
	"github.com/objectbox/objectbox-go/test/model"
)

// tests all queries using the Describe method which serializes the query to string
func TestQueries(t *testing.T) {
	env := model.NewTestEnv(t)
	defer env.Close()

	// the number of entities in the database when the queries are executed
	// if you want to change this, you need to update all the expected numbers in the test cases below
	const baseCount = 1000

	var resetDb = func() {
		assert.NoErr(t, env.Box.RemoveAll())

		// insert new entries
		env.Populate(baseCount)
	}

	var box = env.Box

	// let's alias the entity to make the test cases easier to read
	var E = model.Entity_

	// Use this special entity for testing descriptions
	var e = model.Entity47()

	testCases := []struct {
		// using short variable names because the IDE auto-fills (displays) them in the value initialization
		c int    // expected Query.Count()
		d string // expected Query.Describe()
		q *model.EntityQuery
	}{
		{1, `String == "Val-1"`, box.Query(E.String.Equals(e.String, true))},
		{2, `String == "Val-1"`, box.Query(E.String.Equals(e.String, false))},
		{999, `String != "Val-1"`, box.Query(E.String.NotEquals(e.String, true))},
		{998, `String != "Val-1"`, box.Query(E.String.NotEquals(e.String, false))},
		{64, `String contains "Val-1"`, box.Query(E.String.Contains(e.String, true))},
		{131, `String contains "Val-1"`, box.Query(E.String.Contains(e.String, false))},
		{64, `String starts with "Val-1"`, box.Query(E.String.HasPrefix(e.String, true))},
		{131, `String starts with "Val-1"`, box.Query(E.String.HasPrefix(e.String, false))},
		{1, `String ends with "Val-1"`, box.Query(E.String.HasSuffix(e.String, true))},
		{2, `String ends with "Val-1"`, box.Query(E.String.HasSuffix(e.String, false))},
		{998, `String > "Val-1"`, box.Query(E.String.GreaterThan(e.String, true))},
		{498, `String > "Val-1"`, box.Query(E.String.GreaterThan(e.String, false))},
		{999, `String >= "Val-1"`, box.Query(E.String.GreaterOrEqual(e.String, true))},
		{500, `String >= "Val-1"`, box.Query(E.String.GreaterOrEqual(e.String, false))},
		{1, `String < "Val-1"`, box.Query(E.String.LessThan(e.String, true))},
		{500, `String < "Val-1"`, box.Query(E.String.LessThan(e.String, false))},
		{2, `String <= "Val-1"`, box.Query(E.String.LessOrEqual(e.String, true))},
		{502, `String <= "Val-1"`, box.Query(E.String.LessOrEqual(e.String, false))},
		{2, `String in ["VAL-1", "val-860714888"]`, box.Query(E.String.In(true, "VAL-1", "val-860714888"))},
		{3, `String in ["val-1", "val-860714888"]`, box.Query(E.String.In(false, "VAL-1", "val-860714888"))},

		{1, `Int64 == 0`, box.Query(E.Int64.Equals(0))},
		{2, `Int64 == 47`, box.Query(E.Int64.Equals(e.Int64))},
		{998, `Int64 != 47`, box.Query(E.Int64.NotEquals(e.Int64))},
		{498, `Int64 > 47`, box.Query(E.Int64.GreaterThan(e.Int64))},
		{500, `Int64 < 47`, box.Query(E.Int64.LessThan(e.Int64))},
		{1, `Int64 between -1 and 1`, box.Query(E.Int64.Between(-1, 1))},
		{2, `Int64 between 47 and 94`, box.Query(E.Int64.Between(e.Int64, e.Int64*2))},
		{2, `Int64 in [94|47]`, box.Query(E.Int64.In(e.Int64, e.Int64*2))},
		{998, `Int64 in [94|47]`, box.Query(E.Int64.NotIn(e.Int64, e.Int64*2))},

		{1, `Uint64 == 0`, box.Query(E.Uint64.Equals(0))},
		{2, `Uint64 == 47`, box.Query(E.Uint64.Equals(e.Uint64))},
		{998, `Uint64 != 47`, box.Query(E.Uint64.NotEquals(e.Uint64))},
		//{498,`Uint64 > 47`, box.Query(E.Uint64.GreaterThan(e.Uint64))},
		//{500,`Uint64 < 47`, box.Query(E.Uint64.LessThan(e.Uint64))},
		//{1,`Uint64 between 0 and 1`, box.Query(E.Uint64.Between(0, 1))},
		//{2,`Uint64 between 47 and 94`, box.Query(E.Uint64.Between(e.Uint64, e.Uint64*2))},
		{2, `Uint64 in [94|47]`, box.Query(E.Uint64.In(e.Uint64, e.Uint64*2))},
		{998, `Uint64 in [94|47]`, box.Query(E.Uint64.NotIn(e.Uint64, e.Uint64*2))},

		{1, `Int == 0`, box.Query(E.Int.Equals(0))},
		{2, `Int == 47`, box.Query(E.Int.Equals(e.Int))},
		{998, `Int != 47`, box.Query(E.Int.NotEquals(e.Int))},
		{498, `Int > 47`, box.Query(E.Int.GreaterThan(e.Int))},
		{500, `Int < 47`, box.Query(E.Int.LessThan(e.Int))},
		{1, `Int between -1 and 1`, box.Query(E.Int.Between(-1, 1))},
		{2, `Int between 47 and 94`, box.Query(E.Int.Between(e.Int, e.Int*2))},
		{2, `Int in [94|47]`, box.Query(E.Int.In(e.Int, e.Int*2))},
		{998, `Int in [94|47]`, box.Query(E.Int.NotIn(e.Int, e.Int*2))},

		{1, `Uint == 0`, box.Query(E.Uint.Equals(0))},
		{2, `Uint == 47`, box.Query(E.Uint.Equals(e.Uint))},
		{998, `Uint != 47`, box.Query(E.Uint.NotEquals(e.Uint))},
		//{498,`Uint > 47`, box.Query(E.Uint.GreaterThan(e.Uint))},
		//{500,`Uint < 47`, box.Query(E.Uint.LessThan(e.Uint))},
		//{1,`Uint between 0 and 1`, box.Query(E.Uint.Between(0, 1))},
		//{2,`Uint between 47 and 94`, box.Query(E.Uint.Between(e.Uint, e.Uint*2))},
		{2, `Uint in [94|47]`, box.Query(E.Uint.In(e.Uint, e.Uint*2))},
		{998, `Uint in [94|47]`, box.Query(E.Uint.NotIn(e.Uint, e.Uint*2))},

		{1, `Rune == 0`, box.Query(E.Rune.Equals(0))},
		{3, `Rune == 47`, box.Query(E.Rune.Equals(e.Rune))},
		{997, `Rune != 47`, box.Query(E.Rune.NotEquals(e.Rune))},
		{498, `Rune > 47`, box.Query(E.Rune.GreaterThan(e.Rune))},
		{499, `Rune < 47`, box.Query(E.Rune.LessThan(e.Rune))},
		{1, `Rune between -1 and 1`, box.Query(E.Rune.Between(-1, 1))},
		{3, `Rune between 47 and 94`, box.Query(E.Rune.Between(e.Rune, e.Rune*2))},
		{3, `Rune in [94|47]`, box.Query(E.Rune.In(e.Rune, e.Rune*2))},
		{997, `Rune in [94|47]`, box.Query(E.Rune.NotIn(e.Rune, e.Rune*2))},

		{1, `Int32 == 0`, box.Query(E.Int32.Equals(0))},
		{3, `Int32 == 47`, box.Query(E.Int32.Equals(e.Int32))},
		{997, `Int32 != 47`, box.Query(E.Int32.NotEquals(e.Int32))},
		{498, `Int32 > 47`, box.Query(E.Int32.GreaterThan(e.Int32))},
		{499, `Int32 < 47`, box.Query(E.Int32.LessThan(e.Int32))},
		{1, `Int32 between -1 and 1`, box.Query(E.Int32.Between(-1, 1))},
		{3, `Int32 between 47 and 94`, box.Query(E.Int32.Between(e.Int32, e.Int32*2))},
		{3, `Int32 in [94|47]`, box.Query(E.Int32.In(e.Int32, e.Int32*2))},
		{997, `Int32 in [94|47]`, box.Query(E.Int32.NotIn(e.Int32, e.Int32*2))},

		{1, `Uint32 == 0`, box.Query(E.Uint32.Equals(0))},
		{3, `Uint32 == 47`, box.Query(E.Uint32.Equals(e.Uint32))},
		{997, `Uint32 != 47`, box.Query(E.Uint32.NotEquals(e.Uint32))},
		{498, `Uint32 > 47`, box.Query(E.Uint32.GreaterThan(e.Uint32))},
		{499, `Uint32 < 47`, box.Query(E.Uint32.LessThan(e.Uint32))},
		{1, `Uint32 between 0 and 1`, box.Query(E.Uint32.Between(0, 1))},
		{3, `Uint32 between 47 and 94`, box.Query(E.Uint32.Between(e.Uint32, e.Uint32*2))},
		{3, `Uint32 in [94|47]`, box.Query(E.Uint32.In(e.Uint32, e.Uint32*2))},
		{997, `Uint32 in [94|47]`, box.Query(E.Uint32.NotIn(e.Uint32, e.Uint32*2))},

		{1, `Int16 == 0`, box.Query(E.Int16.Equals(0))},
		{3, `Int16 == 47`, box.Query(E.Int16.Equals(e.Int16))},
		{997, `Int16 != 47`, box.Query(E.Int16.NotEquals(e.Int16))},
		{498, `Int16 > 47`, box.Query(E.Int16.GreaterThan(e.Int16))},
		{499, `Int16 < 47`, box.Query(E.Int16.LessThan(e.Int16))},
		{1, `Int16 between -1 and 1`, box.Query(E.Int16.Between(-1, 1))},
		{4, `Int16 between 47 and 94`, box.Query(E.Int16.Between(e.Int16, e.Int16*2))},
		//{0,`Int16 in [94|47]`, box.Query(E.Int16.In(e.Int16, e.Int16*2))},
		//{0,`Int16 in [94|47]`, box.Query(E.Int16.NotIn(e.Int16, e.Int16*2))},

		{1, `Uint16 == 0`, box.Query(E.Uint16.Equals(0))},
		{3, `Uint16 == 47`, box.Query(E.Uint16.Equals(e.Uint16))},
		{997, `Uint16 != 47`, box.Query(E.Uint16.NotEquals(e.Uint16))},
		{498, `Uint16 > 47`, box.Query(E.Uint16.GreaterThan(e.Uint16))},
		{499, `Uint16 < 47`, box.Query(E.Uint16.LessThan(e.Uint16))},
		{1, `Uint16 between 0 and 1`, box.Query(E.Uint16.Between(0, 1))},
		{4, `Uint16 between 47 and 94`, box.Query(E.Uint16.Between(e.Uint16, e.Uint16*2))},
		//{0,`Uint16 in [94|47]`, box.Query(E.Uint16.In(e.Uint16, e.Uint16*2))},
		//{0,`Uint16 in [94|47]`, box.Query(E.Uint16.NotIn(e.Uint16, e.Uint16*2))},

		{5, `Int8 == 0`, box.Query(E.Int8.Equals(0))},
		{6, `Int8 == 47`, box.Query(E.Int8.Equals(e.Int8))},
		{994, `Int8 != 47`, box.Query(E.Int8.NotEquals(e.Int8))},
		{308, `Int8 > 47`, box.Query(E.Int8.GreaterThan(e.Int8))},
		{686, `Int8 < 47`, box.Query(E.Int8.LessThan(e.Int8))},
		{11, `Int8 between -1 and 1`, box.Query(E.Int8.Between(-1, 1))},
		{179, `Int8 between 47 and 94`, box.Query(E.Int8.Between(e.Int8, e.Int8*2))},
		//{0,`Int8 in [94|47]`, box.Query(E.Int8.In(e.Int8, e.Int8*2))},
		//{0,`Int8 in [94|47]`, box.Query(E.Int8.NotIn(e.Int8, e.Int8*2))},

		{5, `Uint8 == 0`, box.Query(E.Uint8.Equals(0))},
		{6, `Uint8 == 47`, box.Query(E.Uint8.Equals(e.Uint8))},
		{994, `Uint8 != 47`, box.Query(E.Uint8.NotEquals(e.Uint8))},
		{308, `Uint8 > 47`, box.Query(E.Uint8.GreaterThan(e.Uint8))},
		{686, `Uint8 < 47`, box.Query(E.Uint8.LessThan(e.Uint8))},
		{8, `Uint8 between 0 and 1`, box.Query(E.Uint8.Between(0, 1))},
		{179, `Uint8 between 47 and 94`, box.Query(E.Uint8.Between(e.Uint8, e.Uint8*2))},
		//{0,`Uint8 in [94|47]`, box.Query(E.Uint8.In(e.Uint8, e.Uint8*2))},
		//{0,`Uint8 in [94|47]`, box.Query(E.Uint8.NotIn(e.Uint8, e.Uint8*2))},

		{5, `Byte == 0`, box.Query(E.Byte.Equals(0))},
		{6, `Byte == 47`, box.Query(E.Byte.Equals(e.Byte))},
		{994, `Byte != 47`, box.Query(E.Byte.NotEquals(e.Byte))},
		{308, `Byte > 47`, box.Query(E.Byte.GreaterThan(e.Byte))},
		{686, `Byte < 47`, box.Query(E.Byte.LessThan(e.Byte))},
		{8, `Byte between 0 and 1`, box.Query(E.Byte.Between(0, 1))},
		{179, `Byte between 47 and 94`, box.Query(E.Byte.Between(e.Byte, e.Byte*2))},
		//{0,`Byte in [94|47]`, box.Query(E.Byte.In(e.Byte, e.Byte*2))},
		//{0,`Byte in [94|47]`, box.Query(E.Byte.NotIn(e.Byte, e.Byte*2))},

		{2, `Float64 between 47.739999 and 47.740001`, box.Query(E.Float64.Between(e.Float64-0.000001, e.Float64+0.000001))},
		{498, `Float64 > 47.740000`, box.Query(E.Float64.GreaterThan(e.Float64))},
		{500, `Float64 < 47.740000`, box.Query(E.Float64.LessThan(e.Float64))},
		{1, `Float64 between -1.000000 and 1.000000`, box.Query(E.Float64.Between(-1, 1))},
		{2, `Float64 between 47.740000 and 95.480000`, box.Query(E.Float64.Between(e.Float64, e.Float64*2))},

		{2, `Float32 between 47.739990 and 47.740013`, box.Query(E.Float32.Between(e.Float32-0.00001, e.Float32+0.00001))},
		{498, `Float32 > 47.740002`, box.Query(E.Float32.GreaterThan(e.Float32))},
		{500, `Float32 < 47.740002`, box.Query(E.Float32.LessThan(e.Float32))},
		{1, `Float32 between -1.000000 and 1.000000`, box.Query(E.Float32.Between(-1, 1))},
		{2, `Float32 between 47.740002 and 95.480003`, box.Query(E.Float32.Between(e.Float32, e.Float32*2))},

		{6, `ByteVector == `, box.Query(E.ByteVector.Equals(e.ByteVector))},
		//{994,`ByteVector != `, box.Query(E.ByteVector.NotEquals(e.ByteVector))},
		{989, `ByteVector >= `, box.Query(E.ByteVector.GreaterThan(e.ByteVector))},
		{995, `ByteVector >= `, box.Query(E.ByteVector.GreaterOrEqual(e.ByteVector))},
		{5, `ByteVector <= `, box.Query(E.ByteVector.LessThan(e.ByteVector))},
		{11, `ByteVector <= `, box.Query(E.ByteVector.LessOrEqual(e.ByteVector))},

		{256, `Bool == 1`, box.Query(E.Bool.Equals(true))},
		{744, `Bool == 0`, box.Query(E.Bool.Equals(false))},

		{3, `(Bool == 1 AND Byte == 1)`, box.Query(E.Bool.Equals(true), E.Byte.Equals(1))},
	}

	t.Logf("Executing %d test cases", len(testCases))

	for i, tc := range testCases {
		// reset Db before each query, necessary due to query.Remove()
		// TODO we can replace this to make the test run faster by a managed transaction with rollback
		resetDb()

		// assign some readable variable names
		var count = tc.c
		var desc = tc.d
		var query = tc.q

		// Describe
		if actualDesc, err := query.Describe(); err != nil {
			assert.Failf(t, "case #%d {%s} - %s", i, desc, err)
		} else if actualDesc = strings.Replace(actualDesc, "\n", "", -1); desc != actualDesc {
			assert.Failf(t, "case #%d expected {%s}, but got {%s}", i, desc, actualDesc)
		}

		// Find
		var actualData []*model.Entity
		if data, err := query.Find(); err != nil {
			assert.Failf(t, "case #%d {%s} %s", i, desc, err)
		} else if data == nil {
			assert.Failf(t, "case #%d {%s} data is nil", i, desc)
		} else if len(data) != count {
			assert.Failf(t, "case #%d {%s} expected %d, but got %d len(Find())", i, desc, count, len(data))
		} else {
			actualData = data
		}

		// Count
		if actualCount, err := query.Count(); err != nil {
			assert.Failf(t, "case #%d {%s} %s", i, desc, err)
		} else if uint64(count) != actualCount {
			assert.Failf(t, "case #%d {%s} expected %d, but got %d Count()", i, desc, count, actualCount)
		}

		// FindIds
		if ids, err := query.FindIds(); err != nil {
			assert.Failf(t, "case #%d {%s} %s", i, desc, err)
		} else if err := matchAllEntityIds(ids, actualData); err != nil {
			assert.Failf(t, "case #%d {%s} %s", i, desc, err)
		}

		// Remove
		if removedCount, err := query.Remove(); err != nil {
			assert.Failf(t, "case #%d {%s} %s", i, desc, err)
		} else if uint64(count) != removedCount {
			assert.Failf(t, "case #%d {%s} expected %d, but got %d Remove()", i, desc, count, removedCount)
		} else if actualCount, err := env.Box.Count(); err != nil {
			assert.Failf(t, "case #%d {%s} %s", i, desc, err)
		} else if actualCount+removedCount != baseCount {
			assert.Failf(t, "case #%d {%s} expected %d, but got %d Box.Count() after remove",
				i, desc, baseCount-removedCount, actualCount)
		}

		if err := query.Close(); err != nil {
			assert.Failf(t, "case #%d {%s} %s", i, desc, err)
		}
	}
}

func matchAllEntityIds(ids []uint64, items []*model.Entity) error {
	if len(ids) != len(items) {
		return fmt.Errorf("count mismatch = ids=%d, items=%d", len(ids), len(items))
	}

	var merged = map[uint64]int{}

	for _, id := range ids {
		merged[id] = 1
	}

	for _, item := range items {
		if merged[item.Id] == 0 {
			return fmt.Errorf("item %d is missing in the IDs list %v", item.Id, ids)
		}
		merged[item.Id] = merged[item.Id] + 1
	}

	for id, count := range merged {
		if count != 2 {
			return fmt.Errorf("ID %d is missing in the items list", id)
		}
	}

	return nil
}

func TestQueryClose(t *testing.T) {
	env := model.NewTestEnv(t)
	defer env.Close()
	query := env.Box.Query()
	assert.NoErr(t, query.Close())

	_, err := query.Find()
	if err == nil {
		assert.Failf(t, "should fail after Close()")
	}

	// Double close
	assert.NoErr(t, query.Close())
}

// Forces the finalizer to run; not a "real" test with assertions
func TestQueryCloseFinalizer(t *testing.T) {
	env := model.NewTestEnv(t)
	defer env.Close()
	query := env.Box.Query()
	_, _ = query.Count()
	query = nil
	runtime.GC()
	runtime.GC() // 2nd GC allows to set a break point in the finalizer and actually stop there
}

func TestQueryCloseAfterObjectBox(t *testing.T) {
	env := model.NewTestEnv(t)
	query := env.Box.Query()
	queryFinalizer := env.Box.Query()
	_, _ = query.Count()
	_, _ = queryFinalizer.Count()
	env.Close()
	err := query.Close()
	assert.NoErr(t, err)
	queryFinalizer = nil
	runtime.GC()
	runtime.GC() // 2nd GC allows to set a break point in the finalizer and actually stop there
}
