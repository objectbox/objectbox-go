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
	"regexp"
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
	assert.Eq(t, "0", objectbox.StringIdConvertToEntityProperty(0))
	assert.Eq(t, uint64(0), objectbox.StringIdConvertToDatabaseValue("0"))
	assert.Eq(t, "10", objectbox.StringIdConvertToEntityProperty(10))
	assert.Eq(t, uint64(10), objectbox.StringIdConvertToDatabaseValue("10"))

	func() {
		defer assert.MustPanic(t, regexp.MustCompile("error parsing numeric ID represented as string"))
		objectbox.StringIdConvertToDatabaseValue("invalid")
	}()
}

func TestTimeInt64Converter(t *testing.T) {
	assert.Eq(t, "1970-01-01 00:00:00 +0000 UTC", objectbox.TimeInt64ConvertToEntityProperty(0).String())
	assert.Eq(t, "1970-01-01 00:00:01.234 +0000 UTC", objectbox.TimeInt64ConvertToEntityProperty(1234).String())
	assert.Eq(t, "1969-12-31 23:59:54.322 +0000 UTC", objectbox.TimeInt64ConvertToEntityProperty(-5678).String())
	var date = time.Now()
	assert.Eq(t, date.UnixNano()/1000000, objectbox.TimeInt64ConvertToDatabaseValue(date))
}

func TestTimeTextConverter(t *testing.T) {
	date := time.Unix(time.Now().Unix(), int64(time.Now().Nanosecond())) // get date without monotonic clock reading
	bytes, err := date.MarshalText()
	assert.NoErr(t, err)
	assert.Eq(t, string(bytes), objectbox.TimeTextConvertToDatabaseValue(date))
	assert.Eq(t, date.UnixNano(), objectbox.TimeTextConvertToEntityProperty(string(bytes)).UnixNano())
}

func TestTimeBinaryConverter(t *testing.T) {
	date := time.Unix(time.Now().Unix(), int64(time.Now().Nanosecond())) // get date without monotonic clock reading
	bytes, err := date.MarshalBinary()
	assert.NoErr(t, err)
	assert.Eq(t, bytes, objectbox.TimeBinaryConvertToDatabaseValue(date))
	assert.Eq(t, date, objectbox.TimeBinaryConvertToEntityProperty(bytes))
}
