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
	"github.com/objectbox/objectbox-go/test/model"
	"testing"

	"github.com/objectbox/objectbox-go/test/assert"
)

// TestBoxAsync tests the implicit AsyncBox returned by Box.Async()
func TestBoxAsync(t *testing.T) {
	testAsync(t, func(box *model.TestEntityInlineBox) *model.TestEntityInlineAsyncBox {
		return box.Async()
	})
}

const timeoutMs = 100

// TestAsyncBox tests manually managed AsyncBox with custom timeout
func TestAsyncBox(t *testing.T) {
	testAsync(t, func(box *model.TestEntityInlineBox) *model.TestEntityInlineAsyncBox {
		asyncBox, err := objectbox.NewAsyncBox(box.ObjectBox, model.TestEntityInlineBinding.Id, timeoutMs)
		assert.NoErr(t, err)
		assert.True(t, asyncBox != nil)
		return &model.TestEntityInlineAsyncBox{AsyncBox: asyncBox}
	})
}

// TestAsyncBoxGenerated tests manually managed AsyncBox with custom timeout
func TestAsyncBoxGenerated(t *testing.T) {
	testAsync(t, func(box *model.TestEntityInlineBox) *model.TestEntityInlineAsyncBox {
		asyncBox := model.AsyncBoxForTestEntityInline(box.ObjectBox, timeoutMs)
		assert.True(t, asyncBox != nil)
		return asyncBox
	})
}

// testAsync tests all AsyncBox operations
func testAsync(t *testing.T, asyncF func(box *model.TestEntityInlineBox) *model.TestEntityInlineAsyncBox) {
	var env = model.NewTestEnv(t)
	defer env.Close()

	var ob = env.ObjectBox
	var box = model.BoxForTestEntityInline(ob)
	var async = asyncF(box)
	defer func() {
		assert.NoErr(t, async.Close())
	}()

	var waitAndCount = func(expected uint64) {
		assert.NoErr(t, ob.AwaitAsyncCompletion())
		count, err := box.Count()
		assert.NoErr(t, err)
		assert.Eq(t, expected, count)
	}

	var object = &model.TestEntityInline{BaseWithValue: &model.BaseWithValue{}}
	id, err := async.Put(object)
	assert.NoErr(t, err)
	assert.Eq(t, id, object.Id)
	waitAndCount(1)

	// check the inserted object
	objectRead, err := box.Get(id)
	assert.NoErr(t, err)
	assert.Eq(t, *object, *objectRead)

	// update the object
	object.Value = object.Value + 1
	assert.NoErr(t, async.Update(object))
	waitAndCount(1)

	// check the updated object
	objectRead, err = box.Get(id)
	assert.NoErr(t, err)
	assert.Eq(t, *object, *objectRead)

	err = async.Remove(object)
	assert.NoErr(t, err)
	waitAndCount(0)

	// while the update will ultimately fail because the object is removed, it can't know it in advance
	assert.NoErr(t, async.Update(object))
	waitAndCount(0)

	// insert with a custom ID will work just fine
	id, err = async.Insert(object)
	assert.NoErr(t, err)
	assert.Eq(t, object.Id, id)

	id, err = async.Insert(&model.TestEntityInline{BaseWithValue: &model.BaseWithValue{}})
	assert.NoErr(t, err)
	assert.Eq(t, object.Id+1, id)
	waitAndCount(2)

	objects, err := box.GetAll()
	assert.NoErr(t, err)
	assert.Eq(t, 2, len(objects))
	assert.EqItems(t, []uint64{object.Id, id}, []uint64{objects[0].Id, objects[1].Id})

	// insert with an existing ID will fail now (silently)
	var idBefore = object.Id
	id, err = async.Insert(object)
	assert.NoErr(t, err)
	assert.Eq(t, object.Id, idBefore)
	waitAndCount(2)

	assert.NoErr(t, async.RemoveId(object.Id))
	waitAndCount(1)
}
