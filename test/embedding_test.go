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
	"testing"

	"github.com/objectbox/objectbox-go/test/assert"
	"github.com/objectbox/objectbox-go/test/model"
)

func TestStructEmbedding(t *testing.T) {
	var env = model.NewTestEnv(t)
	defer env.Close()

	box := model.BoxForTestEntityInline(env.ObjectBox)

	entity := &model.TestEntityInline{
		BaseWithDate: model.BaseWithDate{
			Date: 4,
		},
		BaseWithValue: &model.BaseWithValue{
			Value: 2,
		},
	}

	id, err := box.Put(entity)
	assert.NoErr(t, err)
	assert.Eq(t, id, entity.Id)

	read, err := box.Get(id)
	assert.NoErr(t, err)
	assert.Eq(t, entity.Id, read.Id)
	assert.Eq(t, entity.Date, read.Date)
	assert.Eq(t, entity.Value, read.Value)
}

func TestStructEmbeddingNilPtr(t *testing.T) {
	var env = model.NewTestEnv(t)
	defer env.Close()

	box := model.BoxForTestEntityInline(env.ObjectBox)

	entity := &model.TestEntityInline{
		BaseWithValue: nil,
	}
	id, err := box.Put(entity)
	assert.NoErr(t, err)
	assert.Eq(t, id, entity.Id)

	read, err := box.Get(id)
	assert.NoErr(t, err)
	assert.True(t, read != nil)
	// TODO invert the condition - should be nil
	// Currently, the generated Load() method creates the containing object regardless if the given slot is present in
	// FlatBuffers. That could be improved by constructing nil-able embedded structs similar to relations, setting to
	// nil if none of the slots was set.
	assert.True(t, read.BaseWithValue != nil)
}
