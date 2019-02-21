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

func TestEntityByValue(t *testing.T) {
	env := model.NewTestEnv(t)
	defer env.Close()

	box := model.BoxForEntityByValue(env.ObjectBox)

	var object = &model.EntityByValue{}

	id, err := box.Put(object)
	assert.NoErr(t, err)
	assert.Eq(t, uint64(1), id)

	fetched, err := box.Get(id)
	assert.NoErr(t, err)
	assert.Eq(t, id, fetched.Id)

	var object2 = &model.EntityByValue{}
	id2, err := box.PutAsync(object2)
	assert.NoErr(t, err)
	assert.Eq(t, uint64(2), id2)

	env.ObjectBox.AwaitAsyncCompletion()

	var objects []model.EntityByValue

	objects, err = box.GetAll()
	assert.NoErr(t, err)
	assert.Eq(t, 2, len(objects))
	assert.Eq(t, id, objects[0].Id)
	assert.Eq(t, id2, objects[1].Id)

	err = box.Remove(object2)
	assert.NoErr(t, err)

	count, err := box.Count()
	assert.NoErr(t, err)
	assert.Eq(t, uint64(1), count)

	objects = []model.EntityByValue{{}, {}}
	ids, err := box.PutAll(objects)
	assert.NoErr(t, err)
	assert.Eq(t, len(objects), len(ids))
	assert.Eq(t, uint64(0), objects[0].Id)
	assert.Eq(t, uint64(0), objects[1].Id)
	assert.Eq(t, uint64(3), ids[0])
	assert.Eq(t, uint64(4), ids[1])

	count, err = box.Count()
	assert.NoErr(t, err)
	assert.Eq(t, uint64(3), count)

	objects, err = box.Query().Find()
	assert.NoErr(t, err)
	assert.Eq(t, 3, len(objects))
	assert.Eq(t, uint64(1), objects[0].Id)
	assert.Eq(t, uint64(3), objects[1].Id)
	assert.Eq(t, uint64(4), objects[2].Id)
}
