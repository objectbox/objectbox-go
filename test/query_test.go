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
	"reflect"
	"regexp"
	"runtime"
	"strings"
	"testing"

	"github.com/objectbox/objectbox-go/objectbox"
	"github.com/objectbox/objectbox-go/test/assert"
	"github.com/objectbox/objectbox-go/test/model"
)

// Following methods use many test-cases defined as a list of queryTestCase and run all Query.* methods on each test case
// see queryTestCase, queryTestOptions and especially testQueries for more details
func TestQueries(t *testing.T) {
	env := model.NewTestEnv(t)
	defer env.Close()

	var box = env.Box

	// let's alias the entity to make the test cases easier to read
	var E = model.Entity_

	// Use this special entity for testing descriptions
	var e = model.Entity47()

	testQueries(t, env, queryTestOptions{baseCount: 1000}, []queryTestCase{
		{1, s{`String == "Val-1"`}, box.Query(E.String.Equals(e.String, true)), nil},
		{2, s{`String ==(i) "Val-1"`}, box.Query(E.String.Equals(e.String, false)), nil},
		{999, s{`String != "Val-1"`}, box.Query(E.String.NotEquals(e.String, true)), nil},
		{998, s{`String !=(i) "Val-1"`}, box.Query(E.String.NotEquals(e.String, false)), nil},
		{64, s{`String contains "Val-1"`}, box.Query(E.String.Contains(e.String, true)), nil},
		{131, s{`String contains(i) "Val-1"`}, box.Query(E.String.Contains(e.String, false)), nil},
		{64, s{`String starts with "Val-1"`}, box.Query(E.String.HasPrefix(e.String, true)), nil},
		{131, s{`String starts with(i) "Val-1"`}, box.Query(E.String.HasPrefix(e.String, false)), nil},
		{1, s{`String ends with "Val-1"`}, box.Query(E.String.HasSuffix(e.String, true)), nil},
		{2, s{`String ends with(i) "Val-1"`}, box.Query(E.String.HasSuffix(e.String, false)), nil},
		{998, s{`String > "Val-1"`}, box.Query(E.String.GreaterThan(e.String, true)), nil},
		{498, s{`String >(i) "Val-1"`}, box.Query(E.String.GreaterThan(e.String, false)), nil},
		{999, s{`String >= "Val-1"`}, box.Query(E.String.GreaterOrEqual(e.String, true)), nil},
		{500, s{`String >=(i) "Val-1"`}, box.Query(E.String.GreaterOrEqual(e.String, false)), nil},
		{1, s{`String < "Val-1"`}, box.Query(E.String.LessThan(e.String, true)), nil},
		{500, s{`String <(i) "Val-1"`}, box.Query(E.String.LessThan(e.String, false)), nil},
		{2, s{`String <= "Val-1"`}, box.Query(E.String.LessOrEqual(e.String, true)), nil},
		{502, s{`String <=(i) "Val-1"`}, box.Query(E.String.LessOrEqual(e.String, false)), nil},
		{2, s{`String in ["VAL-1", "val-860714888"]`, `String in ["val-860714888", "VAL-1"]`}, box.Query(E.String.In(true, "VAL-1", "val-860714888")), nil},
		{3, s{`String in(i) ["val-1", "val-860714888"]`, `String in(i) ["val-860714888", "val-1"]`}, box.Query(E.String.In(false, "VAL-1", "val-860714888")), nil},

		{2, s{`StringVector contains "first-1"`}, box.Query(E.StringVector.Contains("first-1", true)), nil},
		{2, s{`StringVector contains(i) "FIRST-1"`}, box.Query(E.StringVector.Contains("FIRST-1", false)), nil},

		{1, s{`Int64 == 0`}, box.Query(E.Int64.Equals(0)), nil},
		{2, s{`Int64 == 47`}, box.Query(E.Int64.Equals(e.Int64)), nil},
		{998, s{`Int64 != 47`}, box.Query(E.Int64.NotEquals(e.Int64)), nil},
		{498, s{`Int64 > 47`}, box.Query(E.Int64.GreaterThan(e.Int64)), nil},
		{500, s{`Int64 < 47`}, box.Query(E.Int64.LessThan(e.Int64)), nil},
		{1, s{`Int64 between -1 and 1`}, box.Query(E.Int64.Between(-1, 1)), nil},
		{2, s{`Int64 between 47 and 94`}, box.Query(E.Int64.Between(e.Int64, e.Int64*2)), nil},
		{2, s{`Int64 in [94|47]`, `Int64 in [47|94]`}, box.Query(E.Int64.In(e.Int64, e.Int64*2)), nil},
		{998, s{`Int64 not in [94|47]`, `Int64 not in [47|94]`}, box.Query(E.Int64.NotIn(e.Int64, e.Int64*2)), nil},

		{1, s{`Uint64 == 0`}, box.Query(E.Uint64.Equals(0)), nil},
		{2, s{`Uint64 == 47`}, box.Query(E.Uint64.Equals(e.Uint64)), nil},
		{998, s{`Uint64 != 47`}, box.Query(E.Uint64.NotEquals(e.Uint64)), nil},
		{997, s{`Uint64 > 47`}, box.Query(E.Uint64.GreaterThan(e.Uint64)), nil},
		{1, s{`Uint64 < 47`}, box.Query(E.Uint64.LessThan(e.Uint64)), nil},
		{1, s{`Uint64 between 0 and 1`}, box.Query(E.Uint64.Between(0, 1)), nil},
		{2, s{`Uint64 between 47 and 94`}, box.Query(E.Uint64.Between(e.Uint64, e.Uint64*2)), nil},
		{2, s{`Uint64 in [94|47]`, `Uint64 in [47|94]`}, box.Query(E.Uint64.In(e.Uint64, e.Uint64*2)), nil},
		{998, s{`Uint64 not in [94|47]`, `Uint64 not in [47|94]`}, box.Query(E.Uint64.NotIn(e.Uint64, e.Uint64*2)), nil},

		{1, s{`Int == 0`}, box.Query(E.Int.Equals(0)), nil},
		{3, s{`Int == 47`}, box.Query(E.Int.Equals(e.Int)), nil},
		{997, s{`Int != 47`}, box.Query(E.Int.NotEquals(e.Int)), nil},
		{498, s{`Int > 47`}, box.Query(E.Int.GreaterThan(e.Int)), nil},
		{499, s{`Int < 47`}, box.Query(E.Int.LessThan(e.Int)), nil},
		{1, s{`Int between -1 and 1`}, box.Query(E.Int.Between(-1, 1)), nil},
		{3, s{`Int between 47 and 94`}, box.Query(E.Int.Between(e.Int, e.Int*2)), nil},
		{3, s{`Int in [94|47]`, `Int in [47|94]`}, box.Query(E.Int.In(e.Int, e.Int*2)), nil},
		{997, s{`Int not in [94|47]`, `Int not in [47|94]`}, box.Query(E.Int.NotIn(e.Int, e.Int*2)), nil},

		{1, s{`Uint == 0`}, box.Query(E.Uint.Equals(0)), nil},
		{3, s{`Uint == 47`}, box.Query(E.Uint.Equals(e.Uint)), nil},
		{997, s{`Uint != 47`}, box.Query(E.Uint.NotEquals(e.Uint)), nil},
		{996, s{`Uint > 47`}, box.Query(E.Uint.GreaterThan(e.Uint)), nil},
		{1, s{`Uint < 47`}, box.Query(E.Uint.LessThan(e.Uint)), nil},
		{1, s{`Uint between 0 and 1`}, box.Query(E.Uint.Between(0, 1)), nil},
		{3, s{`Uint between 47 and 94`}, box.Query(E.Uint.Between(e.Uint, e.Uint*2)), nil},
		{3, s{`Uint in [94|47]`, `Uint in [47|94]`}, box.Query(E.Uint.In(e.Uint, e.Uint*2)), nil},
		{997, s{`Uint not in [94|47]`, `Uint not in [47|94]`}, box.Query(E.Uint.NotIn(e.Uint, e.Uint*2)), nil},

		{1, s{`Rune == 0`}, box.Query(E.Rune.Equals(0)), nil},
		{3, s{`Rune == 47`}, box.Query(E.Rune.Equals(e.Rune)), nil},
		{997, s{`Rune != 47`}, box.Query(E.Rune.NotEquals(e.Rune)), nil},
		{498, s{`Rune > 47`}, box.Query(E.Rune.GreaterThan(e.Rune)), nil},
		{499, s{`Rune < 47`}, box.Query(E.Rune.LessThan(e.Rune)), nil},
		{1, s{`Rune between -1 and 1`}, box.Query(E.Rune.Between(-1, 1)), nil},
		{3, s{`Rune between 47 and 94`}, box.Query(E.Rune.Between(e.Rune, e.Rune*2)), nil},
		{3, s{`Rune in [94|47]`, `Rune in [47|94]`}, box.Query(E.Rune.In(e.Rune, e.Rune*2)), nil},
		{997, s{`Rune not in [94|47]`, `Rune not in [47|94]`}, box.Query(E.Rune.NotIn(e.Rune, e.Rune*2)), nil},

		{1, s{`Int32 == 0`}, box.Query(E.Int32.Equals(0)), nil},
		{3, s{`Int32 == 47`}, box.Query(E.Int32.Equals(e.Int32)), nil},
		{997, s{`Int32 != 47`}, box.Query(E.Int32.NotEquals(e.Int32)), nil},
		{498, s{`Int32 > 47`}, box.Query(E.Int32.GreaterThan(e.Int32)), nil},
		{499, s{`Int32 < 47`}, box.Query(E.Int32.LessThan(e.Int32)), nil},
		{1, s{`Int32 between -1 and 1`}, box.Query(E.Int32.Between(-1, 1)), nil},
		{3, s{`Int32 between 47 and 94`}, box.Query(E.Int32.Between(e.Int32, e.Int32*2)), nil},
		{3, s{`Int32 in [94|47]`, `Int32 in [47|94]`}, box.Query(E.Int32.In(e.Int32, e.Int32*2)), nil},
		{997, s{`Int32 not in [94|47]`, `Int32 not in [47|94]`}, box.Query(E.Int32.NotIn(e.Int32, e.Int32*2)), nil},

		{1, s{`Uint32 == 0`}, box.Query(E.Uint32.Equals(0)), nil},
		{3, s{`Uint32 == 47`}, box.Query(E.Uint32.Equals(e.Uint32)), nil},
		{997, s{`Uint32 != 47`}, box.Query(E.Uint32.NotEquals(e.Uint32)), nil},
		{996, s{`Uint32 > 47`}, box.Query(E.Uint32.GreaterThan(e.Uint32)), nil},
		{1, s{`Uint32 < 47`}, box.Query(E.Uint32.LessThan(e.Uint32)), nil},
		{1, s{`Uint32 between 0 and 1`}, box.Query(E.Uint32.Between(0, 1)), nil},
		{3, s{`Uint32 between 47 and 94`}, box.Query(E.Uint32.Between(e.Uint32, e.Uint32*2)), nil},
		{3, s{`Uint32 in [94|47]`, `Uint32 in [47|94]`}, box.Query(E.Uint32.In(e.Uint32, e.Uint32*2)), nil},
		{997, s{`Uint32 not in [94|47]`, `Uint32 not in [47|94]`}, box.Query(E.Uint32.NotIn(e.Uint32, e.Uint32*2)), nil},

		{1, s{`Int16 == 0`}, box.Query(E.Int16.Equals(0)), nil},
		{3, s{`Int16 == 47`}, box.Query(E.Int16.Equals(e.Int16)), nil},
		{997, s{`Int16 != 47`}, box.Query(E.Int16.NotEquals(e.Int16)), nil},
		{498, s{`Int16 > 47`}, box.Query(E.Int16.GreaterThan(e.Int16)), nil},
		{499, s{`Int16 < 47`}, box.Query(E.Int16.LessThan(e.Int16)), nil},
		{1, s{`Int16 between -1 and 1`}, box.Query(E.Int16.Between(-1, 1)), nil},
		{4, s{`Int16 between 47 and 94`}, box.Query(E.Int16.Between(e.Int16, e.Int16*2)), nil},
		//{0, s{`Int16 in [94|47]`, `Int16 in [47|94]`}, box.Query(E.Int16.In(e.Int16, e.Int16*2)), nil},
		//{0, s{`Int16 in [94|47]`, `Int16 in [47|94]`}, box.Query(E.Int16.NotIn(e.Int16, e.Int16*2)), nil},

		{1, s{`Uint16 == 0`}, box.Query(E.Uint16.Equals(0)), nil},
		{3, s{`Uint16 == 47`}, box.Query(E.Uint16.Equals(e.Uint16)), nil},
		{997, s{`Uint16 != 47`}, box.Query(E.Uint16.NotEquals(e.Uint16)), nil},
		{996, s{`Uint16 > 47`}, box.Query(E.Uint16.GreaterThan(e.Uint16)), nil},
		{1, s{`Uint16 < 47`}, box.Query(E.Uint16.LessThan(e.Uint16)), nil},
		{1, s{`Uint16 between 0 and 1`}, box.Query(E.Uint16.Between(0, 1)), nil},
		{4, s{`Uint16 between 47 and 94`}, box.Query(E.Uint16.Between(e.Uint16, e.Uint16*2)), nil},
		//{0, s{`Uint16 in [94|47]`, `Uint16 in [47|94]`}, box.Query(E.Uint16.In(e.Uint16, e.Uint16*2)), nil},
		//{0, s{`Uint16 in [94|47]`, `Uint16 in [47|94]`}, box.Query(E.Uint16.NotIn(e.Uint16, e.Uint16*2)), nil},

		{5, s{`Int8 == 0`}, box.Query(E.Int8.Equals(0)), nil},
		{6, s{`Int8 == 47`}, box.Query(E.Int8.Equals(e.Int8)), nil},
		{994, s{`Int8 != 47`}, box.Query(E.Int8.NotEquals(e.Int8)), nil},
		{308, s{`Int8 > 47`}, box.Query(E.Int8.GreaterThan(e.Int8)), nil},
		{686, s{`Int8 < 47`}, box.Query(E.Int8.LessThan(e.Int8)), nil},
		{11, s{`Int8 between -1 and 1`}, box.Query(E.Int8.Between(-1, 1)), nil},
		{179, s{`Int8 between 47 and 94`}, box.Query(E.Int8.Between(e.Int8, e.Int8*2)), nil},
		//{0, s{`Int8 in [94|47]`, `Int8 in [47|94]`}, box.Query(E.Int8.In(e.Int8, e.Int8*2)), nil},
		//{0, s{`Int8 in [94|47]`, `Int8 in [47|94]`}, box.Query(E.Int8.NotIn(e.Int8, e.Int8*2)), nil},

		{5, s{`Uint8 == 0`}, box.Query(E.Uint8.Equals(0)), nil},
		{6, s{`Uint8 == 47`}, box.Query(E.Uint8.Equals(e.Uint8)), nil},
		{994, s{`Uint8 != 47`}, box.Query(E.Uint8.NotEquals(e.Uint8)), nil},
		{806, s{`Uint8 > 47`}, box.Query(E.Uint8.GreaterThan(e.Uint8)), nil},
		{188, s{`Uint8 < 47`}, box.Query(E.Uint8.LessThan(e.Uint8)), nil},
		{8, s{`Uint8 between 0 and 1`}, box.Query(E.Uint8.Between(0, 1)), nil},
		{179, s{`Uint8 between 47 and 94`}, box.Query(E.Uint8.Between(e.Uint8, e.Uint8*2)), nil},
		//{0, s{`Uint8 in [94|47]`, `Uint8 in [47|94]`}, box.Query(E.Uint8.In(e.Uint8, e.Uint8*2)), nil},
		//{0, s{`Uint8 in [94|47]`, `Uint8 in [47|94]`}, box.Query(E.Uint8.NotIn(e.Uint8, e.Uint8*2)), nil},

		{5, s{`Byte == 0`}, box.Query(E.Byte.Equals(0)), nil},
		{6, s{`Byte == 47`}, box.Query(E.Byte.Equals(e.Byte)), nil},
		{994, s{`Byte != 47`}, box.Query(E.Byte.NotEquals(e.Byte)), nil},
		{308, s{`Byte > 47`}, box.Query(E.Byte.GreaterThan(e.Byte)), nil},
		{686, s{`Byte < 47`}, box.Query(E.Byte.LessThan(e.Byte)), nil},
		{8, s{`Byte between 0 and 1`}, box.Query(E.Byte.Between(0, 1)), nil},
		{179, s{`Byte between 47 and 94`}, box.Query(E.Byte.Between(e.Byte, e.Byte*2)), nil},
		//{0, {`Byte in [94|47]`, `Byte in [47|94]`}, box.Query(E.Byte.In(e.Byte, e.Byte*2)), nil},
		//{0, {`Byte in [94|47]`, `Byte in [47|94]`}, box.Query(E.Byte.NotIn(e.Byte, e.Byte*2)), nil},

		{2, s{`Float64 between 47.739999 and 47.740001`}, box.Query(E.Float64.Between(e.Float64-0.000001, e.Float64+0.000001)), nil},
		{498, s{`Float64 > 47.740000`}, box.Query(E.Float64.GreaterThan(e.Float64)), nil},
		{500, s{`Float64 < 47.740000`}, box.Query(E.Float64.LessThan(e.Float64)), nil},
		{1, s{`Float64 between -1.000000 and 1.000000`}, box.Query(E.Float64.Between(-1, 1)), nil},
		{2, s{`Float64 between 47.740000 and 95.480000`}, box.Query(E.Float64.Between(e.Float64, e.Float64*2)), nil},

		{2, s{`Float32 between 47.739990 and 47.740013`}, box.Query(E.Float32.Between(e.Float32-0.00001, e.Float32+0.00001)), nil},
		{498, s{`Float32 > 47.740002`}, box.Query(E.Float32.GreaterThan(e.Float32)), nil},
		{500, s{`Float32 < 47.740002`}, box.Query(E.Float32.LessThan(e.Float32)), nil},
		{1, s{`Float32 between -1.000000 and 1.000000`}, box.Query(E.Float32.Between(-1, 1)), nil},
		{2, s{`Float32 between 47.740002 and 95.480003`}, box.Query(E.Float32.Between(e.Float32, e.Float32*2)), nil},

		{6, s{`ByteVector == byte[5]{0x01020305 08}`}, box.Query(E.ByteVector.Equals(e.ByteVector)), nil},
		//{994, s{`ByteVector != byte[5]{0x01020305 08`}, box.Query(E.ByteVector.NotEquals(e.ByteVector)), nil},
		{989, s{`ByteVector > byte[5]{0x01020305 08}`}, box.Query(E.ByteVector.GreaterThan(e.ByteVector)), nil},
		{995, s{`ByteVector >= byte[5]{0x01020305 08}`}, box.Query(E.ByteVector.GreaterOrEqual(e.ByteVector)), nil},
		{5, s{`ByteVector < byte[5]{0x01020305 08}`}, box.Query(E.ByteVector.LessThan(e.ByteVector)), nil},
		{11, s{`ByteVector <= byte[5]{0x01020305 08}`}, box.Query(E.ByteVector.LessOrEqual(e.ByteVector)), nil},

		{256, s{`Bool == 1`}, box.Query(E.Bool.Equals(true)), nil},
		{744, s{`Bool == 0`}, box.Query(E.Bool.Equals(false)), nil},

		{3, s{`(Bool == 1 AND Byte == 1)`}, box.Query(E.Bool.Equals(true), E.Byte.Equals(1)), nil},
	})

}

func TestQueryOffsetLimit(t *testing.T) {
	env := model.NewTestEnv(t)
	defer env.Close()

	testQueries(t, env, queryTestOptions{baseCount: 10, skipCount: true, skipRemove: true}, []queryTestCase{
		{10, s{`TRUE`}, env.Box.Query(), nil},
		{5, s{`TRUE`}, env.Box.Query().Offset(5), nil},
		{3, s{`TRUE`}, env.Box.Query().Limit(3), nil},
		{3, s{`TRUE`}, env.Box.Query().Offset(3).Limit(3), nil},
		{1, s{`TRUE`}, env.Box.Query().Offset(9).Limit(3), nil},
	})
}

func TestQueryParams(t *testing.T) {
	env := model.NewTestEnv(t)
	defer env.Close()

	var box = env.Box

	// let's alias the entity to make the test cases easier to read
	var E = model.Entity_

	// Use this special entity for testing descriptions
	var e = model.Entity47()

	// and a shorter type name for the setup function argument
	type i = interface{}
	var eq = func(q interface{}) *objectbox.Query { return q.(*objectbox.Query) }

	testQueries(t, env, queryTestOptions{baseCount: 1000}, []queryTestCase{
		{1, s{`String == "Val-1"`}, box.Query(E.String.Equals("", true)),
			func(q i) error { return eq(q).SetStringParams(E.String, e.String) }},
		{1, s{`String in ["VAL-1"]`}, box.Query(E.String.In(true)),
			func(q i) error { return eq(q).SetStringParamsIn(E.String, "VAL-1") }},
		{2, s{`String in ["VAL-1", "val-860714888"]`, `String in ["val-860714888", "VAL-1"]`},
			box.Query(E.String.In(true)),
			func(q i) error { return eq(q).SetStringParamsIn(E.String, "val-860714888", "VAL-1") }},

		{2, s{`StringVector contains "first-1"`}, box.Query(E.StringVector.Contains("", true)),
			func(q i) error { return eq(q).SetStringParams(E.StringVector, "first-1") }},

		{2, s{`Int64 == 47`}, box.Query(E.Int64.Equals(0)),
			func(q i) error { return eq(q).SetInt64Params(E.Int64, e.Int64) }},
		{1, s{`Int64 between -1 and 1`}, box.Query(E.Int64.Between(0, 0)),
			func(q i) error { return eq(q).SetInt64Params(E.Int64, -1, 1) }},
		{2, s{`Int64 in [94|47]`, `Int64 in [47|94]`}, box.Query(E.Int64.In()),
			func(q i) error { return eq(q).SetInt64ParamsIn(E.Int64, e.Int64, e.Int64*2) }},

		{2, s{`Uint64 == 47`}, box.Query(E.Uint64.Equals(0)),
			func(q i) error { return eq(q).SetInt64Params(E.Uint64, int64(e.Uint64)) }},
		{2, s{`Uint64 in [94|47]`, `Uint64 in [47|94]`}, box.Query(E.Uint64.In()),
			func(q i) error { return eq(q).SetInt64ParamsIn(E.Uint64, int64(e.Uint64), int64(e.Uint64*2)) }},

		{3, s{`Int == 47`}, box.Query(E.Int.Equals(0)),
			func(q i) error { return eq(q).SetInt64Params(E.Int, int64(e.Int)) }},
		{1, s{`Int between -1 and 1`}, box.Query(E.Int.Between(0, 0)),
			func(q i) error { return eq(q).SetInt64Params(E.Int, -1, 1) }},
		{3, s{`Int in [94|47]`, `Int in [47|94]`}, box.Query(E.Int.In()),
			func(q i) error { return eq(q).SetInt64ParamsIn(E.Int, int64(e.Int), int64(e.Int*2)) }},

		{3, s{`Uint == 47`}, box.Query(E.Uint.Equals(0)),
			func(q i) error { return eq(q).SetInt64Params(E.Uint, int64(e.Uint)) }},
		{3, s{`Uint in [94|47]`, `Uint in [47|94]`}, box.Query(E.Uint.In()),
			func(q i) error { return eq(q).SetInt64ParamsIn(E.Uint, int64(e.Uint), int64(e.Uint*2)) }},

		{3, s{`Rune == 47`}, box.Query(E.Rune.Equals(0)),
			func(q i) error { return eq(q).SetInt64Params(E.Rune, int64(e.Rune)) }},
		{1, s{`Rune between -1 and 1`}, box.Query(E.Rune.Between(0, 0)),
			func(q i) error { return eq(q).SetInt64Params(E.Rune, -1, 1) }},
		{3, s{`Rune in [94|47]`, `Rune in [47|94]`}, box.Query(E.Rune.In()),
			func(q i) error { return eq(q).SetInt32ParamsIn(E.Rune, int32(e.Rune), int32(e.Rune*2)) }},

		{3, s{`Int32 == 47`}, box.Query(E.Int32.Equals(0)),
			func(q i) error { return eq(q).SetInt64Params(E.Int32, int64(e.Int32)) }},
		{1, s{`Int32 between -1 and 1`}, box.Query(E.Int32.Between(0, 0)),
			func(q i) error { return eq(q).SetInt64Params(E.Int32, -1, 1) }},
		{3, s{`Int32 in [94|47]`, `Int32 in [47|94]`}, box.Query(E.Int32.In()),
			func(q i) error { return eq(q).SetInt32ParamsIn(E.Int32, e.Int32, e.Int32*2) }},

		{3, s{`Uint32 == 47`}, box.Query(E.Uint32.Equals(0)),
			func(q i) error { return eq(q).SetInt64Params(E.Uint32, int64(e.Int32)) }},
		{1, s{`Uint32 between 0 and 1`}, box.Query(E.Uint32.Between(0, 0)),
			func(q i) error { return eq(q).SetInt64Params(E.Uint32, 0, 1) }},
		{3, s{`Uint32 in [94|47]`, `Uint32 in [47|94]`}, box.Query(E.Uint32.In()),
			func(q i) error { return eq(q).SetInt32ParamsIn(E.Uint32, int32(e.Uint32), int32(e.Uint32*2)) }},

		{3, s{`Int16 == 47`}, box.Query(E.Int16.Equals(0)),
			func(q i) error { return eq(q).SetInt64Params(E.Int16, int64(e.Int16)) }},
		{1, s{`Int16 between -1 and 1`}, box.Query(E.Int16.Between(0, 0)),
			func(q i) error { return eq(q).SetInt64Params(E.Int16, -1, 1) }},

		{3, s{`Uint16 == 47`}, box.Query(E.Uint16.Equals(0)),
			func(q i) error { return eq(q).SetInt64Params(E.Uint16, int64(e.Uint16)) }},
		{1, s{`Uint16 between 0 and 1`}, box.Query(E.Uint16.Between(0, 0)),
			func(q i) error { return eq(q).SetInt64Params(E.Uint16, 0, 1) }},

		{6, s{`Int8 == 47`}, box.Query(E.Int8.Equals(0)),
			func(q i) error { return eq(q).SetInt64Params(E.Int8, int64(e.Int8)) }},
		{11, s{`Int8 between -1 and 1`}, box.Query(E.Int8.Between(0, 0)),
			func(q i) error { return eq(q).SetInt64Params(E.Int8, -1, 1) }},

		{6, s{`Uint8 == 47`}, box.Query(E.Uint8.Equals(0)),
			func(q i) error { return eq(q).SetInt64Params(E.Uint8, int64(e.Uint8)) }},
		{8, s{`Uint8 between 0 and 1`}, box.Query(E.Uint8.Between(0, 0)),
			func(q i) error { return eq(q).SetInt64Params(E.Uint8, 0, 1) }},

		{6, s{`Byte == 47`}, box.Query(E.Byte.Equals(0)),
			func(q i) error { return eq(q).SetInt64Params(E.Byte, int64(e.Byte)) }},
		{8, s{`Byte between 0 and 1`}, box.Query(E.Byte.Between(0, 0)),
			func(q i) error { return eq(q).SetInt64Params(E.Byte, 0, 1) }},

		{2, s{`Float64 between 47.739999 and 47.740001`}, box.Query(E.Float64.Between(0, 0)),
			func(q i) error { return eq(q).SetFloat64Params(E.Float64, e.Float64-0.000001, e.Float64+0.000001) }},
		{498, s{`Float64 > 47.740000`}, box.Query(E.Float64.GreaterThan(0)),
			func(q i) error { return eq(q).SetFloat64Params(E.Float64, e.Float64) }},

		{2, s{`Float32 between 47.739990 and 47.740013`}, box.Query(E.Float32.Between(0, 0)),
			func(q i) error {
				return eq(q).SetFloat64Params(E.Float32, float64(e.Float32-0.00001), float64(e.Float32+0.00001))
			}},
		{498, s{`Float32 > 47.740002`}, box.Query(E.Float32.GreaterThan(e.Float32)),
			func(q i) error { return eq(q).SetFloat64Params(E.Float32, float64(e.Float32)) }},

		{6, s{`ByteVector == byte[5]{0x01020305 08}`}, box.Query(E.ByteVector.Equals(nil)),
			func(q i) error { return eq(q).SetBytesParams(E.ByteVector, e.ByteVector) }},
		{1000, s{`ByteVector > byte[0]""`}, box.Query(E.ByteVector.GreaterThan(nil)),
			func(q i) error { return eq(q).SetBytesParams(E.ByteVector, nil) }},
		{5, s{`ByteVector < byte[5]{0x01020305 08}`}, box.Query(E.ByteVector.LessThan(nil)),
			func(q i) error { return eq(q).SetBytesParams(E.ByteVector, e.ByteVector) }},
	})
}

func TestQueryAndOr(t *testing.T) {
	env := model.NewTestEnv(t)
	defer env.Close()

	var box = env.Box

	// let's alias the entity to make the test cases easier to read
	var E = model.Entity_

	// Use this special entity for testing descriptions
	var e = model.Entity47()

	// and a shorter type name for the setup function argument
	type i = interface{}
	var eq = func(q interface{}) *objectbox.Query { return q.(*objectbox.Query) }

	// test standard queries
	testQueries(t, env, queryTestOptions{baseCount: 1000}, []queryTestCase{
		{1, s{`(Int == 0 AND Int32 == 0 AND Int64 == 0)`}, box.Query(E.Int.Equals(0), E.Int32.Equals(0), E.Int64.Equals(0)), nil},
		{1, s{`(Int == 0 AND Int32 == 0 AND Int64 == 0)`}, box.Query(objectbox.All(E.Int.Equals(0), E.Int32.Equals(0), E.Int64.Equals(0))), nil},
		{1, s{`(Int == 0 OR Int32 == 0 OR Int64 == 0)`}, box.Query(objectbox.Any(E.Int.Equals(0), E.Int32.Equals(0), E.Int64.Equals(0))), nil},
		{1, s{`(Int == 0 OR (Int == 0 AND Int32 == 0 AND Int64 == 0) OR Int64 == 0)`}, box.Query(objectbox.Any(E.Int.Equals(0), objectbox.All(E.Int.Equals(0), E.Int32.Equals(0), E.Int64.Equals(0)), E.Int64.Equals(0))), nil},
		{1, s{`(Int == 0 AND Int64 == 0)`}, box.Query(objectbox.Any(E.Int.Equals(0)), objectbox.All(E.Int64.Equals(0))), nil},
		{1000, s{`TRUE`}, box.Query(objectbox.Any(), objectbox.All()), nil},
		{1000, s{`TRUE`}, box.Query(), nil},
	})

	// test when using setParams
	testQueries(t, env, queryTestOptions{baseCount: 1000}, []queryTestCase{
		{0, s{`(Int == 0 AND Int32 == 0 AND Int64 == 47)`}, box.Query(E.Int.Equals(0), E.Int32.Equals(0), E.Int64.Equals(0)),
			func(q i) error { return eq(q).SetInt64Params(E.Int64, e.Int64) }},
		{3, s{`(Int == 0 OR Int32 == 0 OR Int64 == 47)`}, box.Query(objectbox.Any(E.Int.Equals(0), E.Int32.Equals(0), E.Int64.Equals(0))),
			func(q i) error { return eq(q).SetInt64Params(E.Int64, e.Int64) }},
		{1, s{`(Int == 0 OR (Int == 0 AND Int32 == 0 AND Int64 == 47) OR Int64 == 0)`}, box.Query(objectbox.Any(E.Int.Equals(0), objectbox.All(E.Int.Equals(0), E.Int32.Equals(0), E.Int64.Equals(0)), E.Int64.Equals(0))),
			func(q i) error { return eq(q).SetInt64Params(E.Int64, e.Int64) }},
		{0, s{`(Int == 0 AND Int64 == 47)`}, box.Query(objectbox.Any(E.Int.Equals(0)), objectbox.All(E.Int64.Equals(0))),
			func(q i) error { return eq(q).SetInt64Params(E.Int64, e.Int64) }},
	})
}

func TestQueryLinks(t *testing.T) {
	env := model.NewTestEnv(t).SetOptions(model.TestEnvOptions{PopulateRelations: true})
	defer env.Close()

	var box = env.Box
	var boxR = model.BoxForTestEntityRelated(env.ObjectBox)
	var boxV = model.BoxForEntityByValue(env.ObjectBox)

	// Use this special entity for testing descriptions
	var e = model.Entity47()

	// let's alias the entity to make the test cases easier to read
	var E = model.Entity_
	var R = model.TestEntityRelated_
	var V = model.EntityByValue_

	// and a shorter type name for the setup function argument
	type i = interface{}
	var eq = func(q interface{}) *objectbox.Query { return q.(*objectbox.Query) }

	testQueries(t, env, queryTestOptions{baseCount: 10}, []queryTestCase{
		// to-one link
		{2, s{`TRUE Link: Name == "rel-Val-1"`}, box.Query(E.Related.Link(R.Name.Equals("", true))),
			func(q i) error { return eq(q).SetStringParams(R.Name, "rel-Val-1") }},

		// to-one backlink = one-to-many
		{1, s{`TRUE Link: String == "Val-1"`}, boxR.Query(E.Related.Link(E.String.Equals("", true))),
			func(q i) error { return eq(q).SetStringParams(E.String, e.String) }},

		// to-one empty
		{10, s{`TRUE Link: TRUE`}, box.Query(E.Related.Link()), nil},
		{10, s{`TRUE Link: TRUE`}, boxR.Query(E.Related.Link()), nil},

		// to-many link
		{2, s{`TRUE Link: Name == "relPtr-Val-1"`}, box.Query(E.RelatedPtrSlice.Link(R.Name.Equals("", true))),
			func(q i) error { return eq(q).SetStringParams(R.Name, "relPtr-Val-1") }},

		// to-many backlink = many-to-many
		{1, s{`TRUE Link: String == "Val-1"`}, boxR.Query(E.RelatedPtrSlice.Link(E.String.Equals("", true))),
			func(q i) error { return eq(q).SetStringParams(E.String, e.String) }},

		// to-many empty
		{10, s{`TRUE Link: TRUE`}, box.Query(E.RelatedPtrSlice.Link()), nil},
		{10, s{`TRUE Link: TRUE`}, boxR.Query(E.RelatedPtrSlice.Link()), nil},

		// to-one three entities deep link
		{2, s{`TRUE Link: TRUE Link: Text == "RelatedPtr-Next-Val-1"`},
			box.Query(E.RelatedPtr.Link(R.Next.Link(V.Text.Equals("", true)))),
			func(q i) error { return eq(q).SetStringParams(V.Text, "RelatedPtr-Next-Val-1") }},

		// to-one three entities deep backlink
		{1, s{`TRUE Link: TRUE Link: String == "Val-1"`},
			boxV.Query(R.Next.Link(E.RelatedPtr.Link(E.String.Equals("", true)))),
			func(q i) error { return eq(q).SetStringParams(E.String, e.String) }},

		// to-many three entities deep link
		{2, s{`TRUE Link: TRUE Link: Text == "RelatedPtr-NextSlice-Val-1"`},
			box.Query(E.RelatedPtr.Link(R.NextSlice.Link(V.Text.Equals("", true)))),
			func(q i) error { return eq(q).SetStringParams(V.Text, "RelatedPtr-NextSlice-Val-1") }},

		// to-many three entities deep backlink
		{1, s{`TRUE Link: TRUE Link: String == "Val-1"`},
			boxV.Query(R.NextSlice.Link(E.RelatedPtr.Link(E.String.Equals("", true)))),
			func(q i) error { return eq(q).SetStringParams(E.String, e.String) }},
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

	// test Box.Query
	func() {
		defer assert.MustPanic(t, regexp.MustCompile(fmt.Sprintf(
			"property from a different entity %d passed, expected %d", model.EntityByValueBinding.Id, model.EntityBinding.Id)))

		box.Query(model.EntityByValue_.Id.Equals(1))
	}()

	// test Query.Set*Param
	{
		var expected = fmt.Errorf("property from a different entity %d passed, expected %d", model.EntityByValueBinding.Id, model.EntityBinding.Id)
		var query = box.Query(model.Entity_.Id.Equals(0))

		assert.Eq(t, expected, query.SetBytesParams(model.EntityByValue_.Id, []byte{}))
		assert.Eq(t, expected, query.SetFloat64Params(model.EntityByValue_.Id, 1))
		assert.Eq(t, expected, query.SetInt32ParamsIn(model.EntityByValue_.Id, 1))
		assert.Eq(t, expected, query.SetInt64ParamsIn(model.EntityByValue_.Id, 1))
		assert.Eq(t, expected, query.SetInt64Params(model.EntityByValue_.Id, 1))
		assert.Eq(t, expected, query.SetStringParamsIn(model.EntityByValue_.Id, ""))
		assert.Eq(t, expected, query.SetStringParams(model.EntityByValue_.Id, ""))
	}

}

func TestQueryEmptyString(t *testing.T) {
	env := model.NewTestEnv(t)
	defer env.Close()

	_, err := env.Box.Put(&model.Entity{})
	assert.NoErr(t, err)

	count, err := env.Box.Query(model.Entity_.String.Equals("", true)).Count()
	assert.NoErr(t, err)
	assert.Eq(t, uint64(1), count)
}

func TestQueryUint(t *testing.T) {
	env := model.NewTestEnv(t)
	defer env.Close()

	env.Box.Put(&model.Entity{})

	var e = model.Entity47()
	e.Uint = 9223372036854775807 + 10   // > int64 max
	e.Uint64 = 9223372036854775807 + 10 // > int64 max
	e.Uint32 = 2147483647 + 10          // > int32 max

	id, err := env.Box.Put(e)
	assert.NoErr(t, err)

	read, err := env.Box.Get(id)
	assert.NoErr(t, err)
	assert.Eq(t, *e, *read)

	count, err := env.Box.Query(model.Entity_.Uint64.GreaterThan(e.Uint64 - 1)).Count()
	fmt.Println(env.Box.Query(model.Entity_.Uint64.GreaterThan(e.Uint64 - 1)).DescribeParams())
	assert.NoErr(t, err)
	assert.Eq(t, uint64(1), count)
}

// define some type aliases to keep the test-case definitions short & readable
type s = []string

type queryTestCase struct {
	// using short variable names because the IDE auto-fills (displays) them in the value initialization
	c int                     // expected Query.Count()
	d s                       // expected Query.DescribeParams()
	q interface{}             // any of the model.*Query
	f func(interface{}) error // function to configure & check query before executing
}

type queryTestOptions struct {
	baseCount  uint
	skipCount  bool
	skipRemove bool
}

// this function executes tests for all query methods on the given test-cases
func testQueries(t *testing.T, env *model.TestEnv, options queryTestOptions, testCases []queryTestCase) {
	t.Logf("Executing %d test cases", len(testCases))

	var resetDb = func() {
		assert.NoErr(t, model.BoxForEntity(env.ObjectBox).RemoveAll())
		assert.NoErr(t, model.BoxForEntityByValue(env.ObjectBox).RemoveAll())
		assert.NoErr(t, model.BoxForTestEntityRelated(env.ObjectBox).RemoveAll())
		assert.NoErr(t, model.BoxForTestEntityInline(env.ObjectBox).RemoveAll())
		assert.NoErr(t, model.BoxForTestStringIdEntity(env.ObjectBox).RemoveAll())

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
		var setup = tc.f

		var query *objectbox.Query
		var box *objectbox.Box
		var baseCount uint64
		if q, valid := tc.q.(*model.EntityQuery); valid {
			query = q.Query
			box = model.BoxForEntity(env.ObjectBox).Box
			baseCount = uint64(options.baseCount)
		} else if q, valid := tc.q.(*model.TestEntityRelatedQuery); valid {
			query = q.Query
			box = model.BoxForTestEntityRelated(env.ObjectBox).Box
			baseCount = uint64(options.baseCount) * 3 // TestEnv::Populate() currently inserts 3 relations for each main entity
		} else if q, valid := tc.q.(*model.EntityByValueQuery); valid {
			query = q.Query
			box = model.BoxForEntityByValue(env.ObjectBox).Box
			baseCount = uint64(options.baseCount) * 5 // TestEnv::Populate() currently inserts 5 relations for each main entity
		} else {
			assert.Failf(t, "Query is not supported by the test executor: %v", tc.q)
		}

		// run query-setup function, if defined
		if setup != nil {
			if err := setup(query); err != nil {
				assert.Failf(t, "case #%d {%s} setup failed - %s", i, desc, err)
			}
		}

		var isExpected = func(actualDesc string) bool {
			for _, expectedDescription := range desc {
				if expectedDescription == actualDesc {
					return true
				}
			}
			return false
		}

		// DescribeParams
		var removeSpecialChars = strings.NewReplacer("\n", "", "\t", "")
		if actualDesc, err := query.DescribeParams(); err != nil {
			assert.Failf(t, "case #%d {%s} - %s", i, desc, err)
		} else if actualDesc = removeSpecialChars.Replace(actualDesc); !isExpected(actualDesc) {
			assert.Failf(t, "case #%d expected one of %v, but got {%s}", i, desc, actualDesc)
		}

		// Find
		var actualData interface{}
		if data, err := query.Find(); err != nil {
			assert.Failf(t, "case #%d {%s} %s", i, desc, err)
		} else if data == nil {
			assert.Failf(t, "case #%d {%s} data is nil", i, desc)
		} else if reflect.ValueOf(data).Len() != count {
			assert.Failf(t, "case #%d {%s} expected %d, but got %d len(Find())", i, desc, count, reflect.ValueOf(data).Len())
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
		} else {
			//t.Logf("case #%d {%s} - checking all IDs are present in the result", i, desc)
			matchAllEntityIds(t, ids, actualData)
		}

		// Remove
		if !options.skipRemove {
			if removedCount, err := query.Remove(); err != nil {
				assert.Failf(t, "case #%d {%s} %s", i, desc, err)
			} else if uint64(count) != removedCount {
				assert.Failf(t, "case #%d {%s} expected %d, but got %d Remove()", i, desc, count, removedCount)
			} else if actualCount, err := box.Count(); err != nil {
				assert.Failf(t, "case #%d {%s} %s", i, desc, err)
			} else if actualCount+removedCount != baseCount {
				assert.Failf(t, "case #%d {%s} expected %d, but got %d Box.Count() after remove",
					i, desc, baseCount-removedCount, actualCount)
			}
		}

		if err := query.Close(); err != nil {
			assert.Failf(t, "case #%d {%s} %s", i, desc, err)
		}
	}
}

// takes ids & items (slice of one of the model.*Entity) and makes sure all IDs are present
func matchAllEntityIds(t *testing.T, ids []uint64, items interface{}) {
	var actualIds []uint64

	var slice = reflect.ValueOf(items)
	for i := 0; i < slice.Len(); i++ {
		var item = slice.Index(i).Interface()

		if obj, valid := item.(*model.Entity); valid {
			actualIds = append(actualIds, obj.Id)
		} else if obj, valid := item.(*model.TestEntityRelated); valid {
			actualIds = append(actualIds, obj.Id)
		} else if obj, valid := item.(model.EntityByValue); valid {
			actualIds = append(actualIds, obj.Id)
		} else {
			t.Fatalf("type not supported: %v", slice.Type())
		}
	}

	assert.EqItems(t, ids, actualIds)
}
