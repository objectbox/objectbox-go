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
	"regexp"
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

	var box = env.Box

	// let's alias the entity to make the test cases easier to read
	var E = model.Entity_

	// Use this special entity for testing descriptions
	var e = model.Entity47()

	testQueries(t, env, queryTestOptions{baseCount: 1000}, []queryTestCase{
		{1, s{`String == "Val-1"`}, box.Query(E.String.Equals(e.String, true))},
		{2, s{`String ==(i) "Val-1"`}, box.Query(E.String.Equals(e.String, false))},
		{999, s{`String != "Val-1"`}, box.Query(E.String.NotEquals(e.String, true))},
		{998, s{`String !=(i) "Val-1"`}, box.Query(E.String.NotEquals(e.String, false))},
		{64, s{`String contains "Val-1"`}, box.Query(E.String.Contains(e.String, true))},
		{131, s{`String contains(i) "Val-1"`}, box.Query(E.String.Contains(e.String, false))},
		{64, s{`String starts with "Val-1"`}, box.Query(E.String.HasPrefix(e.String, true))},
		{131, s{`String starts with(i) "Val-1"`}, box.Query(E.String.HasPrefix(e.String, false))},
		{1, s{`String ends with "Val-1"`}, box.Query(E.String.HasSuffix(e.String, true))},
		{2, s{`String ends with(i) "Val-1"`}, box.Query(E.String.HasSuffix(e.String, false))},
		{998, s{`String > "Val-1"`}, box.Query(E.String.GreaterThan(e.String, true))},
		{498, s{`String >(i) "Val-1"`}, box.Query(E.String.GreaterThan(e.String, false))},
		{999, s{`String >= "Val-1"`}, box.Query(E.String.GreaterOrEqual(e.String, true))},
		{500, s{`String >=(i) "Val-1"`}, box.Query(E.String.GreaterOrEqual(e.String, false))},
		{1, s{`String < "Val-1"`}, box.Query(E.String.LessThan(e.String, true))},
		{500, s{`String <(i) "Val-1"`}, box.Query(E.String.LessThan(e.String, false))},
		{2, s{`String <= "Val-1"`}, box.Query(E.String.LessOrEqual(e.String, true))},
		{502, s{`String <=(i) "Val-1"`}, box.Query(E.String.LessOrEqual(e.String, false))},
		{2, s{`String in ["VAL-1", "val-860714888"]`, `String in ["val-860714888", "VAL-1"]`}, box.Query(E.String.In(true, "VAL-1", "val-860714888"))},
		{3, s{`String in(i) ["val-1", "val-860714888"]`, `String in(i) ["val-860714888", "val-1"]`}, box.Query(E.String.In(false, "VAL-1", "val-860714888"))},

		{1, s{`Int64 == 0`}, box.Query(E.Int64.Equals(0))},
		{2, s{`Int64 == 47`}, box.Query(E.Int64.Equals(e.Int64))},
		{998, s{`Int64 != 47`}, box.Query(E.Int64.NotEquals(e.Int64))},
		{498, s{`Int64 > 47`}, box.Query(E.Int64.GreaterThan(e.Int64))},
		{500, s{`Int64 < 47`}, box.Query(E.Int64.LessThan(e.Int64))},
		{1, s{`Int64 between -1 and 1`}, box.Query(E.Int64.Between(-1, 1))},
		{2, s{`Int64 between 47 and 94`}, box.Query(E.Int64.Between(e.Int64, e.Int64*2))},
		{2, s{`Int64 in [94|47]`, `Int64 in [47|94]`}, box.Query(E.Int64.In(e.Int64, e.Int64*2))},
		{998, s{`Int64 not in [94|47]`, `Int64 not in [47|94]`}, box.Query(E.Int64.NotIn(e.Int64, e.Int64*2))},

		{1, s{`Uint64 == 0`}, box.Query(E.Uint64.Equals(0))},
		{2, s{`Uint64 == 47`}, box.Query(E.Uint64.Equals(e.Uint64))},
		{998, s{`Uint64 != 47`}, box.Query(E.Uint64.NotEquals(e.Uint64))},
		//{498,`Uint64 > 47`}, box.Query(E.Uint64.GreaterThan(e.Uint64))},
		//{500,`Uint64 < 47`}, box.Query(E.Uint64.LessThan(e.Uint64))},
		//{1,`Uint64 between 0 and 1`}, box.Query(E.Uint64.Between(0, 1))},
		//{2,`Uint64 between 47 and 94`}, box.Query(E.Uint64.Between(e.Uint64, e.Uint64*2))},
		{2, s{`Uint64 in [94|47]`, `Uint64 in [47|94]`}, box.Query(E.Uint64.In(e.Uint64, e.Uint64*2))},
		{998, s{`Uint64 not in [94|47]`, `Uint64 not in [47|94]`}, box.Query(E.Uint64.NotIn(e.Uint64, e.Uint64*2))},

		{1, s{`Int == 0`}, box.Query(E.Int.Equals(0))},
		{2, s{`Int == 47`}, box.Query(E.Int.Equals(e.Int))},
		{998, s{`Int != 47`}, box.Query(E.Int.NotEquals(e.Int))},
		{498, s{`Int > 47`}, box.Query(E.Int.GreaterThan(e.Int))},
		{500, s{`Int < 47`}, box.Query(E.Int.LessThan(e.Int))},
		{1, s{`Int between -1 and 1`}, box.Query(E.Int.Between(-1, 1))},
		{2, s{`Int between 47 and 94`}, box.Query(E.Int.Between(e.Int, e.Int*2))},
		{2, s{`Int in [94|47]`, `Int in [47|94]`}, box.Query(E.Int.In(e.Int, e.Int*2))},
		{998, s{`Int not in [94|47]`, `Int not in [47|94]`}, box.Query(E.Int.NotIn(e.Int, e.Int*2))},

		{1, s{`Uint == 0`}, box.Query(E.Uint.Equals(0))},
		{2, s{`Uint == 47`}, box.Query(E.Uint.Equals(e.Uint))},
		{998, s{`Uint != 47`}, box.Query(E.Uint.NotEquals(e.Uint))},
		//{498,`Uint > 47`}, box.Query(E.Uint.GreaterThan(e.Uint))},
		//{500,`Uint < 47`}, box.Query(E.Uint.LessThan(e.Uint))},
		//{1,`Uint between 0 and 1`}, box.Query(E.Uint.Between(0, 1))},
		//{2,`Uint between 47 and 94`}, box.Query(E.Uint.Between(e.Uint, e.Uint*2))},
		{2, s{`Uint in [94|47]`, `Uint in [47|94]`}, box.Query(E.Uint.In(e.Uint, e.Uint*2))},
		{998, s{`Uint not in [94|47]`, `Uint not in [47|94]`}, box.Query(E.Uint.NotIn(e.Uint, e.Uint*2))},

		{1, s{`Rune == 0`}, box.Query(E.Rune.Equals(0))},
		{3, s{`Rune == 47`}, box.Query(E.Rune.Equals(e.Rune))},
		{997, s{`Rune != 47`}, box.Query(E.Rune.NotEquals(e.Rune))},
		{498, s{`Rune > 47`}, box.Query(E.Rune.GreaterThan(e.Rune))},
		{499, s{`Rune < 47`}, box.Query(E.Rune.LessThan(e.Rune))},
		{1, s{`Rune between -1 and 1`}, box.Query(E.Rune.Between(-1, 1))},
		{3, s{`Rune between 47 and 94`}, box.Query(E.Rune.Between(e.Rune, e.Rune*2))},
		{3, s{`Rune in [94|47]`, `Rune in [47|94]`}, box.Query(E.Rune.In(e.Rune, e.Rune*2))},
		{997, s{`Rune not in [94|47]`, `Rune not in [47|94]`}, box.Query(E.Rune.NotIn(e.Rune, e.Rune*2))},

		{1, s{`Int32 == 0`}, box.Query(E.Int32.Equals(0))},
		{3, s{`Int32 == 47`}, box.Query(E.Int32.Equals(e.Int32))},
		{997, s{`Int32 != 47`}, box.Query(E.Int32.NotEquals(e.Int32))},
		{498, s{`Int32 > 47`}, box.Query(E.Int32.GreaterThan(e.Int32))},
		{499, s{`Int32 < 47`}, box.Query(E.Int32.LessThan(e.Int32))},
		{1, s{`Int32 between -1 and 1`}, box.Query(E.Int32.Between(-1, 1))},
		{3, s{`Int32 between 47 and 94`}, box.Query(E.Int32.Between(e.Int32, e.Int32*2))},
		{3, s{`Int32 in [94|47]`, `Int32 in [47|94]`}, box.Query(E.Int32.In(e.Int32, e.Int32*2))},
		{997, s{`Int32 not in [94|47]`, `Int32 not in [47|94]`}, box.Query(E.Int32.NotIn(e.Int32, e.Int32*2))},

		{1, s{`Uint32 == 0`}, box.Query(E.Uint32.Equals(0))},
		{3, s{`Uint32 == 47`}, box.Query(E.Uint32.Equals(e.Uint32))},
		{997, s{`Uint32 != 47`}, box.Query(E.Uint32.NotEquals(e.Uint32))},
		{498, s{`Uint32 > 47`}, box.Query(E.Uint32.GreaterThan(e.Uint32))},
		{499, s{`Uint32 < 47`}, box.Query(E.Uint32.LessThan(e.Uint32))},
		{1, s{`Uint32 between 0 and 1`}, box.Query(E.Uint32.Between(0, 1))},
		{3, s{`Uint32 between 47 and 94`}, box.Query(E.Uint32.Between(e.Uint32, e.Uint32*2))},
		{3, s{`Uint32 in [94|47]`, `Uint32 in [47|94]`}, box.Query(E.Uint32.In(e.Uint32, e.Uint32*2))},
		{997, s{`Uint32 not in [94|47]`, `Uint32 not in [47|94]`}, box.Query(E.Uint32.NotIn(e.Uint32, e.Uint32*2))},

		{1, s{`Int16 == 0`}, box.Query(E.Int16.Equals(0))},
		{3, s{`Int16 == 47`}, box.Query(E.Int16.Equals(e.Int16))},
		{997, s{`Int16 != 47`}, box.Query(E.Int16.NotEquals(e.Int16))},
		{498, s{`Int16 > 47`}, box.Query(E.Int16.GreaterThan(e.Int16))},
		{499, s{`Int16 < 47`}, box.Query(E.Int16.LessThan(e.Int16))},
		{1, s{`Int16 between -1 and 1`}, box.Query(E.Int16.Between(-1, 1))},
		{4, s{`Int16 between 47 and 94`}, box.Query(E.Int16.Between(e.Int16, e.Int16*2))},
		//{0, s{`Int16 in [94|47]`, `Int16 in [47|94]`}, box.Query(E.Int16.In(e.Int16, e.Int16*2))},
		//{0, s{`Int16 in [94|47]`, `Int16 in [47|94]`}, box.Query(E.Int16.NotIn(e.Int16, e.Int16*2))},

		{1, s{`Uint16 == 0`}, box.Query(E.Uint16.Equals(0))},
		{3, s{`Uint16 == 47`}, box.Query(E.Uint16.Equals(e.Uint16))},
		{997, s{`Uint16 != 47`}, box.Query(E.Uint16.NotEquals(e.Uint16))},
		{498, s{`Uint16 > 47`}, box.Query(E.Uint16.GreaterThan(e.Uint16))},
		{499, s{`Uint16 < 47`}, box.Query(E.Uint16.LessThan(e.Uint16))},
		{1, s{`Uint16 between 0 and 1`}, box.Query(E.Uint16.Between(0, 1))},
		{4, s{`Uint16 between 47 and 94`}, box.Query(E.Uint16.Between(e.Uint16, e.Uint16*2))},
		//{0, {`Uint16 in [94|47]`, `Uint16 in [47|94]`}, box.Query(E.Uint16.In(e.Uint16, e.Uint16*2))},
		//{0, {`Uint16 in [94|47]`, `Uint16 in [47|94]`}, box.Query(E.Uint16.NotIn(e.Uint16, e.Uint16*2))},

		{5, s{`Int8 == 0`}, box.Query(E.Int8.Equals(0))},
		{6, s{`Int8 == 47`}, box.Query(E.Int8.Equals(e.Int8))},
		{994, s{`Int8 != 47`}, box.Query(E.Int8.NotEquals(e.Int8))},
		{308, s{`Int8 > 47`}, box.Query(E.Int8.GreaterThan(e.Int8))},
		{686, s{`Int8 < 47`}, box.Query(E.Int8.LessThan(e.Int8))},
		{11, s{`Int8 between -1 and 1`}, box.Query(E.Int8.Between(-1, 1))},
		{179, s{`Int8 between 47 and 94`}, box.Query(E.Int8.Between(e.Int8, e.Int8*2))},
		//{0, {`Int8 in [94|47]`, `Int8 in [47|94]`}, box.Query(E.Int8.In(e.Int8, e.Int8*2))},
		//{0, {`Int8 in [94|47]`, `Int8 in [47|94]`}, box.Query(E.Int8.NotIn(e.Int8, e.Int8*2))},

		{5, s{`Uint8 == 0`}, box.Query(E.Uint8.Equals(0))},
		{6, s{`Uint8 == 47`}, box.Query(E.Uint8.Equals(e.Uint8))},
		{994, s{`Uint8 != 47`}, box.Query(E.Uint8.NotEquals(e.Uint8))},
		{308, s{`Uint8 > 47`}, box.Query(E.Uint8.GreaterThan(e.Uint8))},
		{686, s{`Uint8 < 47`}, box.Query(E.Uint8.LessThan(e.Uint8))},
		{8, s{`Uint8 between 0 and 1`}, box.Query(E.Uint8.Between(0, 1))},
		{179, s{`Uint8 between 47 and 94`}, box.Query(E.Uint8.Between(e.Uint8, e.Uint8*2))},
		//{0, {`Uint8 in [94|47]`, `Uint8 in [47|94]`}, box.Query(E.Uint8.In(e.Uint8, e.Uint8*2))},
		//{0, {`Uint8 in [94|47]`, `Uint8 in [47|94]`}, box.Query(E.Uint8.NotIn(e.Uint8, e.Uint8*2))},

		{5, s{`Byte == 0`}, box.Query(E.Byte.Equals(0))},
		{6, s{`Byte == 47`}, box.Query(E.Byte.Equals(e.Byte))},
		{994, s{`Byte != 47`}, box.Query(E.Byte.NotEquals(e.Byte))},
		{308, s{`Byte > 47`}, box.Query(E.Byte.GreaterThan(e.Byte))},
		{686, s{`Byte < 47`}, box.Query(E.Byte.LessThan(e.Byte))},
		{8, s{`Byte between 0 and 1`}, box.Query(E.Byte.Between(0, 1))},
		{179, s{`Byte between 47 and 94`}, box.Query(E.Byte.Between(e.Byte, e.Byte*2))},
		//{0, {`Byte in [94|47]`, `Byte in [47|94]`}, box.Query(E.Byte.In(e.Byte, e.Byte*2))},
		//{0, {`Byte in [94|47]`, `Byte in [47|94]`}, box.Query(E.Byte.NotIn(e.Byte, e.Byte*2))},

		{2, s{`Float64 between 47.739999 and 47.740001`}, box.Query(E.Float64.Between(e.Float64-0.000001, e.Float64+0.000001))},
		{498, s{`Float64 > 47.740000`}, box.Query(E.Float64.GreaterThan(e.Float64))},
		{500, s{`Float64 < 47.740000`}, box.Query(E.Float64.LessThan(e.Float64))},
		{1, s{`Float64 between -1.000000 and 1.000000`}, box.Query(E.Float64.Between(-1, 1))},
		{2, s{`Float64 between 47.740000 and 95.480000`}, box.Query(E.Float64.Between(e.Float64, e.Float64*2))},

		{2, s{`Float32 between 47.739990 and 47.740013`}, box.Query(E.Float32.Between(e.Float32-0.00001, e.Float32+0.00001))},
		{498, s{`Float32 > 47.740002`}, box.Query(E.Float32.GreaterThan(e.Float32))},
		{500, s{`Float32 < 47.740002`}, box.Query(E.Float32.LessThan(e.Float32))},
		{1, s{`Float32 between -1.000000 and 1.000000`}, box.Query(E.Float32.Between(-1, 1))},
		{2, s{`Float32 between 47.740002 and 95.480003`}, box.Query(E.Float32.Between(e.Float32, e.Float32*2))},

		{6, s{`ByteVector == byte[5]{0x01020305 08}`}, box.Query(E.ByteVector.Equals(e.ByteVector))},
		//{994,`ByteVector != byte[5]{0x01020305 08`}, box.Query(E.ByteVector.NotEquals(e.ByteVector))},
		{989, s{`ByteVector > byte[5]{0x01020305 08}`}, box.Query(E.ByteVector.GreaterThan(e.ByteVector))},
		{995, s{`ByteVector >= byte[5]{0x01020305 08}`}, box.Query(E.ByteVector.GreaterOrEqual(e.ByteVector))},
		{5, s{`ByteVector < byte[5]{0x01020305 08}`}, box.Query(E.ByteVector.LessThan(e.ByteVector))},
		{11, s{`ByteVector <= byte[5]{0x01020305 08}`}, box.Query(E.ByteVector.LessOrEqual(e.ByteVector))},

		{256, s{`Bool == 1`}, box.Query(E.Bool.Equals(true))},
		{744, s{`Bool == 0`}, box.Query(E.Bool.Equals(false))},

		{3, s{`(Bool == 1 AND Byte == 1)`}, box.Query(E.Bool.Equals(true), E.Byte.Equals(1))},
	})

}

func TestQueryOffsetLimit(t *testing.T) {
	env := model.NewTestEnv(t)
	defer env.Close()

	testQueries(t, env, queryTestOptions{baseCount: 10, skipCount: true, skipRemove: true}, []queryTestCase{
		{10, s{`TRUE`}, env.Box.Query()},
		{5, s{`TRUE`}, env.Box.Query().Offset(5)},
		{3, s{`TRUE`}, env.Box.Query().Limit(3)},
		{3, s{`TRUE`}, env.Box.Query().Offset(3).Limit(3)},
		{1, s{`TRUE`}, env.Box.Query().Offset(9).Limit(3)},
	})
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

func TestQueryWrongEntity(t *testing.T) {
	env := model.NewTestEnv(t)
	defer env.Close()

	// try to use condition on a different entity than the one in the box
	var box = model.BoxForEntity(env.ObjectBox)

	var setWrongCondition = func() {
		defer assert.MustPanic(t, regexp.MustCompile(fmt.Sprintf(
			"property from a different entity %d passed, expected %d", model.EntityByValueBinding.Id, model.EntityBinding.Id)))

		box.Query(model.EntityByValue_.Id.Equals(1))
	}

	setWrongCondition()
}

type queryTestCase struct {
	// using short variable names because the IDE auto-fills (displays) them in the value initialization
	c int      // expected Query.Count()
	d []string // expected Query.Describe()
	q *model.EntityQuery
}

type queryTestOptions struct {
	baseCount  uint
	skipCount  bool
	skipRemove bool
}

// to keep the test-case definitions short &readable
type s = []string

// this function executes tests for all query methods on the given test-cases
func testQueries(t *testing.T, env *model.TestEnv, options queryTestOptions, testCases []queryTestCase) {
	t.Logf("Executing %d test cases", len(testCases))

	var resetDb = func() {
		assert.NoErr(t, env.Box.RemoveAll())

		// insert new entries
		env.Populate(options.baseCount)
	}

	for i, tc := range testCases {
		// reset Db before each query, necessary due to query.Remove()
		// TODO we can replace this to make the test run faster by a managed transaction with rollback
		resetDb()

		// assign some readable variable names
		var count = tc.c
		var desc = tc.d
		var query = tc.q

		var isExpected = func(actualDesc string) bool {
			for _, expectedDescription := range desc {
				if expectedDescription == actualDesc {
					return true
				}
			}
			return false
		}

		// Describe
		if actualDesc, err := query.Describe(); err != nil {
			assert.Failf(t, "case #%d {%s} - %s", i, desc, err)
		} else if actualDesc = strings.Replace(actualDesc, "\n", "", -1); !isExpected(actualDesc) {
			assert.Failf(t, "case #%d expected one of %v, but got {%s}", i, desc, actualDesc)
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
		if !options.skipCount {
			if actualCount, err := query.Count(); err != nil {
				assert.Failf(t, "case #%d {%s} %s", i, desc, err)
			} else if uint64(count) != actualCount {
				assert.Failf(t, "case #%d {%s} expected %d, but got %d Count()", i, desc, count, actualCount)
			}
		}

		// FindIds
		if ids, err := query.FindIds(); err != nil {
			assert.Failf(t, "case #%d {%s} %s", i, desc, err)
		} else if err := matchAllEntityIds(ids, actualData); err != nil {
			assert.Failf(t, "case #%d {%s} %s", i, desc, err)
		}

		// Remove
		if !options.skipRemove {
			if removedCount, err := query.Remove(); err != nil {
				assert.Failf(t, "case #%d {%s} %s", i, desc, err)
			} else if uint64(count) != removedCount {
				assert.Failf(t, "case #%d {%s} expected %d, but got %d Remove()", i, desc, count, removedCount)
			} else if actualCount, err := env.Box.Count(); err != nil {
				assert.Failf(t, "case #%d {%s} %s", i, desc, err)
			} else if actualCount+removedCount != uint64(options.baseCount) {
				assert.Failf(t, "case #%d {%s} expected %d, but got %d Box.Count() after remove",
					i, desc, uint64(options.baseCount)-removedCount, actualCount)
			}
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
