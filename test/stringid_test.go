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
	"strconv"
	"testing"

	"github.com/objectbox/objectbox-go/test/assert"

	"github.com/objectbox/objectbox-go/test/model"
)

func TestStringIdSingleOps(t *testing.T) {
	env := model.NewTestEnv(t)
	box := model.BoxForTestStringIdEntity(env.ObjectBox)

	var object = &model.TestStringIdEntity{}

	id, err := box.Put(object)
	assert.NoErr(t, err)
	assert.Eq(t, uint64(1), id)
	assert.Eq(t, strconv.FormatUint(id, 10), object.Id)

	fetched, err := box.Get(id)
	assert.NoErr(t, err)
	assert.Eq(t, strconv.FormatUint(id, 10), fetched.Id)

	var object2 = &model.TestStringIdEntity{}
	id2, err := box.PutAsync(object2)
	assert.NoErr(t, err)
	assert.Eq(t, uint64(2), id2)
	assert.Eq(t, strconv.FormatUint(id2, 10), object2.Id)

	env.ObjectBox.AwaitAsyncCompletion()

	all, err := box.GetAll()
	assert.NoErr(t, err)
	assert.Eq(t, 2, len(all))
	assert.Eq(t, strconv.FormatUint(id, 10), all[0].Id)
	assert.Eq(t, strconv.FormatUint(id2, 10), all[1].Id)

	err = box.Remove(object2)
	assert.NoErr(t, err)

	count, err := box.Count()
	assert.NoErr(t, err)
	assert.Eq(t, uint64(1), count)
}

func TestStringIdMultiOps(t *testing.T) {
	env := model.NewTestEnv(t)
	box := model.BoxForTestStringIdEntity(env.ObjectBox)

	objects := []*model.TestStringIdEntity{{}, {}}

	ids, err := box.PutAll(objects)
	assert.NoErr(t, err)
	assert.Eq(t, len(objects), len(ids))
	assert.Eq(t, "1", objects[0].Id)
	assert.Eq(t, "2", objects[1].Id)
	assert.Eq(t, uint64(1), ids[0])
	assert.Eq(t, uint64(2), ids[1])

	count, err := box.Count()
	assert.NoErr(t, err)
	assert.Eq(t, uint64(2), count)

	query := box.Query(model.TestStringIdEntity_.Id.Equals(2))
	found, err := query.Find()
	assert.NoErr(t, err)
	assert.Eq(t, 1, len(found))
	assert.Eq(t, "2", found[0].Id)

	err = box.RemoveAll()
	assert.NoErr(t, err)

	count, err = box.Count()
	assert.NoErr(t, err)
	assert.Eq(t, uint64(0), count)
}
