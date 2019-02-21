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
