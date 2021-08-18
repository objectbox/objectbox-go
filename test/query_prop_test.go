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

package objectbox_test

import (
	"github.com/objectbox/objectbox-go/objectbox"
	"github.com/objectbox/objectbox-go/test/assert"
	"github.com/objectbox/objectbox-go/test/model"
	"math"
	"reflect"
	"regexp"
	"testing"
)

func TestPropQuerySetup(t *testing.T) {
	env := model.NewTestEnv(t)
	defer env.Close()

	{ // standard positive use case
		pq, err := env.Box.Query().PropertyOrError(model.Entity_.String)
		assert.NoErr(t, err)
		assert.True(t, pq != nil)

		assert.True(t, nil != env.Box.Query().Property(model.Entity_.String))
	}

	{ // invalid entity
		pq, err := env.Box.Query().PropertyOrError(model.TestEntityRelated_.Name)
		assert.Err(t, err)
		assert.True(t, pq == nil)

		func() {
			defer assert.MustPanic(t, regexp.MustCompile("property from a different entity"))
			env.Box.Query().Property(model.TestEntityRelated_.Name)
		}()
	}
}

func TestPropQueryFind(t *testing.T) {
	env := model.NewTestEnv(t)
	defer env.Close()

	const count = 10
	env.Populate(count / 2)

	// make sure there are duplicates
	objects, err := env.Box.GetAll()
	assert.NoErr(t, err)
	for i := range objects {
		objects[i].Id = 0
	}
	_, err = env.Box.PutMany(objects)
	assert.NoErr(t, err)

	// let's use aliases to make the test cases declarations shorter
	type I = interface{}
	type M = *model.Entity
	type PQ = *objectbox.PropertyQuery
	type e = error
	var E = model.Entity_

	testPropertyQueries(t, env, []propertyQueryTestCase{
		{E.Int, func(o M) I { return o.Int }, func(q PQ) (I, e) { return q.FindInts(nil) }},
		{E.Int8, func(o M) I { return o.Int8 }, func(q PQ) (I, e) { return q.FindInt8s(nil) }},
		{E.Int16, func(o M) I { return o.Int16 }, func(q PQ) (I, e) { return q.FindInt16s(nil) }},
		{E.Int32, func(o M) I { return o.Int32 }, func(q PQ) (I, e) { return q.FindInt32s(nil) }},
		{E.Int64, func(o M) I { return o.Int64 }, func(q PQ) (I, e) { return q.FindInt64s(nil) }},
		{E.Uint, func(o M) I { return o.Uint }, func(q PQ) (I, e) { return q.FindUints(nil) }},
		{E.Uint8, func(o M) I { return o.Uint8 }, func(q PQ) (I, e) { return q.FindUint8s(nil) }},
		{E.Uint16, func(o M) I { return o.Uint16 }, func(q PQ) (I, e) { return q.FindUint16s(nil) }},
		{E.Uint32, func(o M) I { return o.Uint32 }, func(q PQ) (I, e) { return q.FindUint32s(nil) }},
		{E.Uint64, func(o M) I { return o.Uint64 }, func(q PQ) (I, e) { return q.FindUint64s(nil) }},
		{E.Bool, func(o M) I { return o.Bool }, func(q PQ) (I, e) { return q.FindBools(nil) }},
		{E.String, func(o M) I { return o.String }, func(q PQ) (I, e) { return q.FindStrings(nil) }},
		{E.Date, func(o M) I { return o.Date.UnixNano() / 1000000 }, func(q PQ) (I, e) { return q.FindInt64s(nil) }},
		// TODO not available in C  {E.StringVector, func(o M) I { return o.StringVector }, func(q PQ) (I, e) { return q.FindStringSlices(nil) }},
		{E.Byte, func(o M) I { return o.Byte }, func(q PQ) (I, e) { return q.FindUint8s(nil) }},
		// TODO not available in C  {E.ByteVector, func(o M) I { return o.ByteVector }, func(q PQ) (I, e) { return q.FindByteVectors(nil) }},
		{E.Rune, func(o M) I { return o.Rune }, func(q PQ) (I, e) { return q.FindInt32s(nil) }},
		{E.Float32, func(o M) I { return o.Float32 }, func(q PQ) (I, e) { return q.FindFloat32s(nil) }},
		{E.Float64, func(o M) I { return o.Float64 }, func(q PQ) (I, e) { return q.FindFloat64s(nil) }},
	})

}

type propertyQueryTestCase struct {
	p   objectbox.Property
	exp func(o *model.Entity) interface{}                     // selector, returning the given property value for each object
	act func(q *objectbox.PropertyQuery) (interface{}, error) // execute the query, e.g. `return q.FindStrings()`
}

// this function executes tests for all query methods on each test case
func testPropertyQueries(t *testing.T, env *model.TestEnv, testCases []propertyQueryTestCase) {
	var query = env.Box.Query()

	objects, err := env.Box.GetAll()
	assert.NoErr(t, err)
	count := len(objects)
	assert.True(t, count > 0)

	// utilities to read a single property from all objects = simulating what Property Query does
	var readProperty = func(fn func(o *model.Entity) interface{}) interface{} {
		var elType = reflect.TypeOf(fn(objects[0]))
		var sliceType = reflect.SliceOf(elType)

		var slice = reflect.MakeSlice(sliceType, 0, 0)

		for _, object := range objects {
			slice = reflect.Append(slice, reflect.ValueOf(fn(object)))
		}
		return slice.Interface()
	}

	// NOTE: this will probably break on vectors
	var distinct = func(list interface{}) interface{} {
		var inputSlice = reflect.ValueOf(list)
		var elType = inputSlice.Index(0).Type()
		var mapType = reflect.MapOf(elType, elType)
		var rmap = reflect.MakeMap(mapType)

		for i := 0; i < inputSlice.Len(); i++ {
			var el = inputSlice.Index(i)
			rmap.SetMapIndex(el, el)
		}

		var keys = rmap.MapKeys()
		var result = reflect.MakeSlice(inputSlice.Type(), 0, 0)
		for _, key := range keys {
			result = reflect.Append(result, key)
		}
		return result.Interface()
	}

	for k, tc := range testCases {
		t.Logf("TestCase %v - %v", k+1, reflect.TypeOf(tc.p).String())

		var pq = query.Property(tc.p)

		// distinct = false
		var expected = readProperty(tc.exp)
		actual, err := tc.act(pq)
		assert.NoErr(t, err)
		assert.EqItems(t, expected, actual)

		// distinct = true
		if _, ok := tc.p.(*objectbox.PropertyString); ok {
			assert.NoErr(t, pq.DistinctString(true, true))
		} else {
			assert.NoErr(t, pq.Distinct(true))
		}

		expected = distinct(expected)
		var expectedLen = reflect.ValueOf(expected).Len()
		t.Logf("  count all = %v, count distinct = %v", count, expectedLen)
		assert.True(t, expectedLen > 0 && expectedLen < count)
		actual, err = tc.act(pq)
		assert.NoErr(t, err)
		assert.EqItems(t, expected, actual)
	}
}

func TestPropQueryAggregate(t *testing.T) {
	env := model.NewTestEnv(t)
	defer env.Close()

	// let's alias the entity to make the test cases easier to read
	var E = model.Entity_

	// matcher for invalid function variant called - e.g. double() variant for int property
	var na = regexp.MustCompile("This operation is not supported for")
	var overflow = regexp.MustCompile("Numeric overflow")

	var maxInsertedUint64 = uint64(18446744023243685935) // can't be inlined

	type i = int64
	type f = float64

	// first test when the database is empty
	testPropertyQueriesAggregate(t, env, []propertyQueryAggregateTestCase{
		// Test case # 1
		{E.Int, math.NaN(), i(0), i(0), i(0), na, na, f(0)},
		{E.Int8, math.NaN(), i(0), i(0), i(0), na, na, f(0)},
		{E.Int16, math.NaN(), i(0), i(0), i(0), na, na, f(0)},
		{E.Int32, math.NaN(), i(0), i(0), i(0), na, na, f(0)},
		{E.Int64, math.NaN(), i(0), i(0), i(0), na, na, f(0)},
		// Test case # 6
		{E.Uint, math.NaN(), i(0), i(0), i(0), na, na, f(0)},
		{E.Uint8, math.NaN(), i(0), i(0), i(0), na, na, f(0)},
		{E.Uint16, math.NaN(), i(0), i(0), i(0), na, na, f(0)},
		{E.Uint32, math.NaN(), i(0), i(0), i(0), na, na, f(0)},
		{E.Uint64, math.NaN(), i(0), i(0), i(0), na, na, f(0)},
		// Test case # 11
		{E.Bool, math.NaN(), na, na, i(0), na, na, f(0)},
		{E.String, na, na, na, na, na, na, na},
		{E.StringVector, na, na, na, na, na, na, na},
		{E.Byte, math.NaN(), i(0), i(0), i(0), na, na, f(0)},
		{E.ByteVector, na, na, na, na, na, na, na},
		// Test case # 16
		{E.Rune, math.NaN(), i(0), i(0), i(0), na, na, f(0)},
		{E.Float32, math.NaN(), na, na, na, math.NaN(), math.NaN(), f(0)},
		{E.Float64, math.NaN(), na, na, na, math.NaN(), math.NaN(), f(0)},
	})

	// Next insert some data and test again
	env.Populate(10)

	// DEV:
	// objects, _ := env.Box.GetAll()
	// for i := range objects {
	//	fmt.Println(objects[i].Float32, " ")
	// }

	// We choose to hard-code the numbers below (calculated "manually") instead of calculating them to notice potential
	// issues and make sure they're the same across platforms.
	testPropertyQueriesAggregate(t, env, []propertyQueryAggregateTestCase{
		// Test case # 1
		{E.Int, 14.1, i(-2147483601), i(2147483601), i(141), na, na, f(141)},
		{E.Int8, 14.1, i(-47), i(47), i(141), na, na, f(141)},
		{E.Int16, 14.1, i(-47), i(47), i(141), na, na, f(141)},
		{E.Int32, 14.1, i(-2147483601), i(2147483601), i(141), na, na, f(141)},
		{E.Int64, -20186346277.1, i(-201863462865), i(151397597137), i(-201863462771), na, na, f(-201863462771)},
		// Test case # 6
		{E.Uint, 1288490202.9, i(0), i(3221225519), i(12884902029), na, na, f(12884902029)},
		{E.Uint8, 90.9, i(0), i(209), i(909), na, na, f(909)},
		{E.Uint16, 19674.9, i(0), i(65489), i(196749), na, na, f(196749)},
		{E.Uint32, 1288490202.9, i(0), i(3221225519), i(12884902029), na, na, f(12884902029)},
		{E.Uint64, 7378697609297474369.3, i(0), i(maxInsertedUint64), overflow, na, na, f(73786976092974743693)},
		// Test case # 11
		{E.Bool, 0.5, na, na, i(5), na, na, f(5)},
		{E.String, na, na, na, na, na, na, na},
		{E.StringVector, na, na, na, na, na, na, na},
		{E.Byte, 90.9, i(0), i(209), i(909), na, na, f(909)},
		{E.ByteVector, na, na, na, na, na, na, na},
		// Test case # 16
		{E.Rune, 14.1, i(-2147483601), i(2147483601), i(141), na, na, f(141)},
		{E.Float32, -20504174582.452, na, na, na, -205041745920.0, 153781305344.0, -205041745824.52},
		{E.Float64, -20504173856.782, na, na, na, -205041738663.30002, 153781303985.54, -205041738567.82},
	})
}

type propertyQueryAggregateTestCase struct {
	p objectbox.Property
	// each field represents the expected result (an error or a value) for the given function
	// common:
	avg interface{}
	// int64:
	min interface{}
	max interface{}
	sum interface{}
	// float64:
	minF interface{}
	maxF interface{}
	sumF interface{}
}

func propertyQueryAssertResult(t *testing.T, name string, expected interface{}, result interface{}, err error) {
	t.Logf("verifying %s()", name)

	if expected == nil {
		return
	}

	if f, ok := expected.(float64); ok && math.IsNaN(f) {
		assert.True(t, math.IsNaN(result.(float64)))
		return
	}

	if err == nil {
		assert.Eq(t, expected, result)
		return
	}

	if reg, ok := expected.(*regexp.Regexp); ok {
		assert.MustMatch(t, reg, err)
		return
	}

	assert.Eq(t, expected, err)
}

func propertyQueryAssertResultInt64(t *testing.T, name string, expected interface{}, fn func() (int64, error)) {
	result, err := fn()
	propertyQueryAssertResult(t, name, expected, result, err)
}

func propertyQueryAssertResultUint64(t *testing.T, name string, expected interface{}, fn func() (uint64, error)) {
	result, err := fn()
	propertyQueryAssertResult(t, name, expected, result, err)
}

func propertyQueryAssertResultFloat64(t *testing.T, name string, expected interface{}, fn func() (float64, error)) {
	result, err := fn()
	propertyQueryAssertResult(t, name, expected, result, err)
}

// this function executes tests for all query methods on each test case
func testPropertyQueriesAggregate(t *testing.T, env *model.TestEnv, testCases []propertyQueryAggregateTestCase) {
	var query = env.Box.Query()

	count, err := query.Count()
	assert.NoErr(t, err)

	for k, tc := range testCases {
		t.Logf("TestCase %v - %v", k+1, reflect.TypeOf(tc.p).String())

		var pq = query.Property(tc.p)

		propertyQueryAssertResultUint64(t, "count", count, pq.Count)
		propertyQueryAssertResultFloat64(t, "average", tc.avg, pq.Average)

		propertyQueryAssertResultInt64(t, "min", tc.min, pq.Min)
		propertyQueryAssertResultInt64(t, "max", tc.max, pq.Max)
		propertyQueryAssertResultInt64(t, "sum", tc.sum, pq.Sum)

		propertyQueryAssertResultFloat64(t, "minF", tc.minF, pq.MinFloat64)
		propertyQueryAssertResultFloat64(t, "maxF", tc.maxF, pq.MaxFloat64)
		propertyQueryAssertResultFloat64(t, "sumF", tc.sumF, pq.SumFloat64)
	}
}
