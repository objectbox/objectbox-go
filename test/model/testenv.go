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

package model

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/objectbox/objectbox-go/test/assert"

	"github.com/objectbox/objectbox-go/objectbox"
)

type TestEnv struct {
	ObjectBox *objectbox.ObjectBox
	Box       *EntityBox

	t      *testing.T
	dbName string
}

func removeDb(name string) {
	os.Remove(filepath.Join(name, "data.mdb"))
	os.Remove(filepath.Join(name, "lock.mdb"))
}

func NewTestEnv(t *testing.T) *TestEnv {
	var dbName = "testdata"

	removeDb(dbName)

	ob, err := objectbox.NewBuilder().Directory(dbName).Model(ObjectBoxModel()).Build()
	if err != nil {
		t.Fatal(err)
	}
	return &TestEnv{
		ObjectBox: ob,
		Box:       BoxForEntity(ob),
		dbName:    dbName,
		t:         t,
	}
}

func (env *TestEnv) Close() {
	env.ObjectBox.Close()
	removeDb(env.dbName)
}

func (env *TestEnv) Populate(count int) {
	// the first one is always the special Entity47
	env.PutEntity(Entity47())

	// additionally an entity with upper case String
	var e = Entity47()
	e.String = strings.ToUpper(e.String)
	env.PutEntity(e)

	var toInsert = count - 2

	// insert some data - different values but dependable
	var limit = float64(4294967295) // uint max so that when we multiply with 42, we get some large int64 values
	var step = limit / float64(toInsert) * 2
	var entities = make([]*Entity, toInsert)
	var i = 0
	for coef := -limit; i < toInsert; coef += step {
		entities[i] = entity47(int(coef))
		i++
	}

	_, err := env.Box.PutAll(entities)
	assert.NoErr(env.t, err)

	c, err := env.Box.Count()
	assert.NoErr(env.t, err)
	assert.Eq(env.t, count, int(c))
}

func (env *TestEnv) PutEntity(entity *Entity) uint64 {
	id, err := env.Box.Put(entity)
	assert.NoErr(env.t, err)

	return id
}

func Entity47() *Entity {
	return entity47(1)
}

func entity47(coef int) *Entity {
	// NOTE, it doesn't really matter that we overflow the smaller types
	var Bool = coef%2 == 1

	var String string
	if Bool {
		String = fmt.Sprintf("Val-%d", coef)
	} else {
		String = fmt.Sprintf("val-%d", coef)
	}

	return &Entity{
		Int:        47 * coef,
		Int8:       47 * int8(coef),
		Int16:      47 * int16(coef),
		Int32:      47 * int32(coef),
		Int64:      47 * int64(coef),
		Uint:       47 * uint(coef),
		Uint8:      47 * uint8(coef),
		Uint16:     47 * uint16(coef),
		Uint32:     47 * uint32(coef),
		Uint64:     47 * uint64(coef),
		Bool:       Bool,
		String:     String,
		Byte:       47 * byte(coef),
		ByteVector: []byte{1 * byte(coef), 2 * byte(coef), 3 * byte(coef), 5 * byte(coef), 8 * byte(coef)},
		Rune:       47 * rune(coef),
		Float32:    47.74 * float32(coef),
		Float64:    47.74 * float64(coef),
	}
}
