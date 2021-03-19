/*
 * Copyright 2018-2021 ObjectBox Ltd. All rights reserved.
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
	"testing"
	"time"

	"github.com/objectbox/objectbox-go/test/assert"
	"github.com/objectbox/objectbox-go/test/model"
)

func TestTimeConverter(t *testing.T) {
	var env = model.NewTestEnv(t)
	defer env.Close()

	date, err := time.Parse(time.RFC3339, "2018-11-28T12:16:42.145+07:00")
	assert.NoErr(t, err)

	id := env.PutEntity(&model.Entity{Date: date})
	assert.Eq(t, uint64(1), id)

	read, err := env.Box.Get(id)
	assert.NoErr(t, err)
	assert.Eq(t, date.UnixNano(), read.Date.UnixNano())
}

func TestComplex128Converter(t *testing.T) {
	var env = model.NewTestEnv(t)
	defer env.Close()

	var value = complex(14, 125)

	id := env.PutEntity(&model.Entity{Complex128: value})
	assert.Eq(t, uint64(1), id)

	read, err := env.Box.Get(id)
	assert.NoErr(t, err)
	assert.Eq(t, value, read.Complex128)
}

func TestStringIdConverter(t *testing.T) {
	{
		value, err := objectbox.StringIdConvertToEntityProperty(0)
		assert.NoErr(t, err)
		assert.Eq(t, "0", value)
	}

	{
		value, err := objectbox.StringIdConvertToDatabaseValue("0")
		assert.NoErr(t, err)
		assert.Eq(t, uint64(0), value)
	}

	{
		value, err := objectbox.StringIdConvertToEntityProperty(10)
		assert.NoErr(t, err)
		assert.Eq(t, "10", value)
	}

	{
		value, err := objectbox.StringIdConvertToDatabaseValue("10")
		assert.NoErr(t, err)
		assert.Eq(t, uint64(10), value)
	}

	{
		value, err := objectbox.StringIdConvertToDatabaseValue("invalid")
		assert.Err(t, err)
		assert.Eq(t, uint64(0), value)
	}
}

func TestTimeInt64Converter(t *testing.T) {
	var test = func(expected string, timestamp int64) {
		value, err := objectbox.TimeInt64ConvertToEntityProperty(timestamp)
		assert.NoErr(t, err)
		assert.Eq(t, expected, value.String())
	}

	test("1970-01-01 00:00:00 +0000 UTC", 0)
	test("1970-01-01 00:00:01.234 +0000 UTC", 1234)
	test("1969-12-31 23:59:54.322 +0000 UTC", -5678)

	{
		var date = time.Now()
		value, err := objectbox.TimeInt64ConvertToDatabaseValue(date)
		assert.NoErr(t, err)
		assert.Eq(t, date.UnixNano()/1000000, value)
	}
}

func TestTimeTextConverter(t *testing.T) {
	date := time.Unix(time.Now().Unix(), int64(time.Now().Nanosecond())) // get date without monotonic clock reading
	bytes, err := date.MarshalText()
	assert.NoErr(t, err)

	{
		value, err := objectbox.TimeTextConvertToDatabaseValue(date)
		assert.NoErr(t, err)
		assert.Eq(t, string(bytes), value)
	}

	{
		value, err := objectbox.TimeTextConvertToEntityProperty(string(bytes))
		assert.NoErr(t, err)
		assert.Eq(t, date.UnixNano(), value.UnixNano())
	}
}

func TestTimeBinaryConverter(t *testing.T) {
	date := time.Unix(time.Now().Unix(), int64(time.Now().Nanosecond())) // get date without monotonic clock reading
	bytes, err := date.MarshalBinary()
	assert.NoErr(t, err)

	{
		value, err := objectbox.TimeBinaryConvertToDatabaseValue(date)
		assert.NoErr(t, err)
		assert.Eq(t, bytes, value)
	}

	{
		value, err := objectbox.TimeBinaryConvertToEntityProperty(bytes)
		assert.NoErr(t, err)
		assert.Eq(t, date, value)
	}
}
