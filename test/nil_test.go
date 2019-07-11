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

	"github.com/objectbox/objectbox-go/test/assert"
	"github.com/objectbox/objectbox-go/test/model"
)

func TestNilPropertiesWhenNil(t *testing.T) {
	env := model.NewTestEnv(t)
	defer env.Close()

	box := model.BoxForEntity(env.ObjectBox)

	var object = &model.Entity{}
	id, err := box.Put(object)
	assert.NoErr(t, err)
	assert.Eq(t, uint64(1), id)

	object, err = box.Get(id)
	assert.NoErr(t, err)
	assert.Eq(t, id, object.Id)
	assert.True(t, nil == object.IntPtr)
	assert.True(t, nil == object.Int8Ptr)
	assert.True(t, nil == object.Int16Ptr)
	assert.True(t, nil == object.Int32Ptr)
	assert.True(t, nil == object.Int64Ptr)
	assert.True(t, nil == object.UintPtr)
	assert.True(t, nil == object.Uint8Ptr)
	assert.True(t, nil == object.Uint16Ptr)
	assert.True(t, nil == object.Uint32Ptr)
	assert.True(t, nil == object.Uint64Ptr)
	assert.True(t, nil == object.BoolPtr)
	assert.True(t, nil == object.StringPtr)
	assert.True(t, nil == object.StringVectorPtr)
	assert.True(t, nil == object.BytePtr)
	assert.True(t, nil == object.ByteVectorPtr)
	assert.True(t, nil == object.RunePtr)
	assert.True(t, nil == object.Float32Ptr)
	assert.True(t, nil == object.Float64Ptr)
}

func TestNilPropertiesWithValues(t *testing.T) {
	env := model.NewTestEnv(t)
	defer env.Close()

	box := model.BoxForEntity(env.ObjectBox)

	// source for the values
	var prototype = model.Entity47()

	var object = &model.Entity{}
	object.IntPtr = &prototype.Int
	object.Int8Ptr = &prototype.Int8
	object.Int16Ptr = &prototype.Int16
	object.Int32Ptr = &prototype.Int32
	object.Int64Ptr = &prototype.Int64
	object.UintPtr = &prototype.Uint
	object.Uint8Ptr = &prototype.Uint8
	object.Uint16Ptr = &prototype.Uint16
	object.Uint32Ptr = &prototype.Uint32
	object.Uint64Ptr = &prototype.Uint64
	object.BoolPtr = &prototype.Bool
	object.StringPtr = &prototype.String
	object.StringVectorPtr = &prototype.StringVector
	object.BytePtr = &prototype.Byte
	object.ByteVectorPtr = &prototype.ByteVector
	object.RunePtr = &prototype.Rune
	object.Float32Ptr = &prototype.Float32
	object.Float64Ptr = &prototype.Float64

	id, err := box.Put(object)
	assert.NoErr(t, err)
	assert.Eq(t, uint64(1), id)

	object, err = box.Get(id)
	assert.NoErr(t, err)
	assert.Eq(t, id, object.Id)

	assert.True(t, nil != object.IntPtr)
	assert.True(t, nil != object.Int8Ptr)
	assert.True(t, nil != object.Int16Ptr)
	assert.True(t, nil != object.Int32Ptr)
	assert.True(t, nil != object.Int64Ptr)
	assert.True(t, nil != object.UintPtr)
	assert.True(t, nil != object.Uint8Ptr)
	assert.True(t, nil != object.Uint16Ptr)
	assert.True(t, nil != object.Uint32Ptr)
	assert.True(t, nil != object.Uint64Ptr)
	assert.True(t, nil != object.BoolPtr)
	assert.True(t, nil != object.StringPtr)
	assert.True(t, nil != object.StringVectorPtr)
	assert.True(t, nil != object.BytePtr)
	assert.True(t, nil != object.ByteVectorPtr)
	assert.True(t, nil != object.RunePtr)
	assert.True(t, nil != object.Float32Ptr)
	assert.True(t, nil != object.Float64Ptr)

	assert.Eq(t, prototype.Int, *object.IntPtr)
	assert.Eq(t, prototype.Int8, *object.Int8Ptr)
	assert.Eq(t, prototype.Int16, *object.Int16Ptr)
	assert.Eq(t, prototype.Int32, *object.Int32Ptr)
	assert.Eq(t, prototype.Int64, *object.Int64Ptr)
	assert.Eq(t, prototype.Uint, *object.UintPtr)
	assert.Eq(t, prototype.Uint8, *object.Uint8Ptr)
	assert.Eq(t, prototype.Uint16, *object.Uint16Ptr)
	assert.Eq(t, prototype.Uint32, *object.Uint32Ptr)
	assert.Eq(t, prototype.Uint64, *object.Uint64Ptr)
	assert.Eq(t, prototype.Bool, *object.BoolPtr)
	assert.Eq(t, prototype.String, *object.StringPtr)
	assert.Eq(t, prototype.StringVector, *object.StringVectorPtr)
	assert.Eq(t, prototype.Byte, *object.BytePtr)
	assert.Eq(t, prototype.ByteVector, *object.ByteVectorPtr)
	assert.Eq(t, prototype.Rune, *object.RunePtr)
	assert.Eq(t, prototype.Float32, *object.Float32Ptr)
	assert.Eq(t, prototype.Float64, *object.Float64Ptr)
}

// This tests correct behaviour of both vectors that are nil and those that are empty
func TestNilPropertiesVectors(t *testing.T) {
	env := model.NewTestEnv(t)
	defer env.Close()

	box := model.BoxForEntity(env.ObjectBox)

	// empty vectors (not nil!)
	id, err := box.Put(&model.Entity{
		StringVector:    []string{},
		StringVectorPtr: &[]string{},
		ByteVector:      []byte{},
		ByteVectorPtr:   &[]byte{},
	})
	assert.NoErr(t, err)

	object, err := box.Get(id)
	assert.NoErr(t, err)
	assert.True(t, nil != object.StringVector)
	assert.True(t, nil != object.StringVectorPtr)
	assert.True(t, nil != object.ByteVector)
	assert.True(t, nil != object.ByteVectorPtr)

	// nil vectors
	id, err = box.Put(&model.Entity{
		StringVector:    nil,
		StringVectorPtr: nil,
		ByteVector:      nil,
		ByteVectorPtr:   nil,
	})
	assert.NoErr(t, err)

	object, err = box.Get(id)
	assert.NoErr(t, err)
	assert.True(t, nil == object.StringVector)
	assert.True(t, nil == object.StringVectorPtr)
	assert.True(t, nil == object.ByteVector)
	assert.True(t, nil == object.ByteVectorPtr)
}
