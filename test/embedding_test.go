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

// NOTE nil pointers are currently not supported
//func TestStructEmbeddingNilPtr(t *testing.T) {
//	var env = model.NewTestEnv(t)
//	defer env.Close()
//
//	box := model.BoxForTestEntityInline(env.ObjectBox)
//
//	entity := &model.TestEntityInline{
//		BaseWithValue: nil,
//	}
//	id, err := box.Put(entity)
//	assert.NoErr(t, err)
//	assert.Eq(t, id, entity.Id)
//}
