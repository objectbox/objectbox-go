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

package model

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/objectbox/objectbox-go/objectbox"
	"github.com/objectbox/objectbox-go/test/assert"
)

// TestEnv provides environment for testing ObjectBox. It sets up the database and populates it with data.
type TestEnv struct {
	ObjectBox *objectbox.ObjectBox
	Box       *EntityBox
	Directory string

	t       *testing.T
	options TestEnvOptions
}

// TestEnvOptions configure the TestEnv
type TestEnvOptions struct {
	PopulateRelations bool
}

// NewTestEnv creates the test environment
func NewTestEnv(t *testing.T) *TestEnv {
	// Test in a temporary directory - if tested by an end user, the repo is read-only.
	tempDir, err := ioutil.TempDir("", "objectbox-test")
	if err != nil {
		t.Fatal(err)
	}

	var env = &TestEnv{
		Directory: tempDir,
		t:         t,
	}

	env.ObjectBox, err = objectbox.NewBuilder().Directory(env.Directory).Model(ObjectBoxModel()).Build()
	if err != nil {
		t.Fatal(err)
	}

	env.Box = BoxForEntity(env.ObjectBox)

	return env
}

// Close closes ObjectBox and removes the database
func (env *TestEnv) Close() {
	env.ObjectBox.Close()
	os.RemoveAll(env.Directory)
}

func (env *TestEnv) SyncClient(serverUri string) *objectbox.SyncClient {
	err, client := objectbox.NewSyncClient(env.ObjectBox, serverUri, objectbox.SyncCredentialsNone())
	assert.NoErr(env.t, err)
	assert.True(env.t, client != nil)
	return client
}

func removeFileIfExists(path string) error {
	if _, err := os.Stat(path); err == nil {
		return os.Remove(path)
	}
	return nil
}

// SetOptions configures options
func (env *TestEnv) SetOptions(options TestEnvOptions) *TestEnv {
	env.options = options
	return env
}

// Populate creates given number of entities in the database
func (env *TestEnv) Populate(count uint) {
	// the first one is always the special Entity47
	env.PutEntity(entity47(1, &env.options))

	if count > 1 {
		// additionally an entity with upper case String
		var e = entity47(1, &env.options)
		e.String = strings.ToUpper(e.String)
		env.PutEntity(e)
	}

	if count > 2 {
		var toInsert = count - 2

		// insert some data - different values but dependable
		var limit = float64(4294967295) // uint max so that when we multiply with 42, we get some large int64 values
		var step = limit / float64(toInsert) * 2
		var entities = make([]*Entity, toInsert)
		var i = uint(0)
		for coef := -limit; i < toInsert; coef += step {
			entities[i] = entity47(int64(coef), &env.options)
			i++
		}

		_, err := env.Box.PutMany(entities)
		assert.NoErr(env.t, err)
	}

	c, err := env.Box.Count()
	assert.NoErr(env.t, err)
	assert.Eq(env.t, uint64(count), c)
}

// PutEntity creates an entity
func (env *TestEnv) PutEntity(entity *Entity) uint64 {
	id, err := env.Box.Put(entity)
	assert.NoErr(env.t, err)

	return id
}

// Entity47 creates a test entity ("47" because int fields are multiples of 47)
func Entity47() *Entity {
	return entity47(1, nil)
}

// entity47 creates a test entity ("47" because int fields are multiples of 47)
func entity47(coef int64, options *TestEnvOptions) *Entity {
	// NOTE, it doesn't really matter that we overflow the smaller types
	var Bool = coef%2 == 1

	var String string
	if Bool {
		String = fmt.Sprintf("Val-%d", coef)
	} else {
		String = fmt.Sprintf("val-%d", coef)
	}

	var object = &Entity{
		Int:          int(int32(47 * coef)),
		Int8:         47 * int8(coef),
		Int16:        47 * int16(coef),
		Int32:        47 * int32(coef),
		Int64:        47 * int64(coef),
		Uint:         uint(uint32(47 * coef)),
		Uint8:        47 * uint8(coef),
		Uint16:       47 * uint16(coef),
		Uint32:       47 * uint32(coef),
		Uint64:       47 * uint64(coef),
		Bool:         Bool,
		String:       String,
		StringVector: []string{fmt.Sprintf("first-%d", coef), fmt.Sprintf("second-%d", coef), ""},
		Byte:         47 * byte(coef),
		ByteVector:   []byte{1 * byte(coef), 2 * byte(coef), 3 * byte(coef), 5 * byte(coef), 8 * byte(coef)},
		Rune:         47 * rune(coef),
		Float32:      47.74 * float32(coef),
		Float64:      47.74 * float64(coef),
	}
	var err error
	object.Date, err = objectbox.TimeInt64ConvertToEntityProperty(47 * int64(coef))
	if err != nil {
		panic(err)
	}

	if options != nil && options.PopulateRelations {
		object.Related = TestEntityRelated{Name: "rel-" + String}
		object.RelatedPtr = &TestEntityRelated{
			Name:      "relPtr-" + String,
			Next:      &EntityByValue{Text: "RelatedPtr-Next-" + String},
			NextSlice: []EntityByValue{{Text: "RelatedPtr-NextSlice-" + String}},
		}
		object.RelatedSlice = []EntityByValue{{Text: "relByValue-" + String}}
		object.RelatedPtrSlice = []*TestEntityRelated{{
			Name:      "relPtr-" + String,
			Next:      &EntityByValue{Text: "RelatedPtrSlice-Next-" + String},
			NextSlice: []EntityByValue{{Text: "RelatedPtrSlice-NextSlice-" + String}},
		}}
	} else {
		object.Related.NextSlice = []EntityByValue{}
		object.RelatedSlice = []EntityByValue{}
		object.RelatedPtrSlice = nil // lazy loaded so Get() would set nil here as well
	}

	return object
}
