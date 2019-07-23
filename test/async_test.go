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

	for testCase := 0; testCase <= 1; testCase++ {
		var object = &model.TestEntityInline{BaseWithValue: &model.BaseWithValue{}}
		id, err := async.Put(object)
		assert.NoErr(t, err)
		assert.Eq(t, id, object.Id)

		assert.NoErr(t, ob.AwaitAsyncCompletion())

		count, err := box.Count()
		assert.NoErr(t, err)
		assert.True(t, count == 1)

		objectRead, err := box.Get(id)
		assert.NoErr(t, err)
		assert.Eq(t, *object, *objectRead)

		if testCase == 0 {
			err = async.Remove(object)
		} else {
			err = async.RemoveId(id)
		}
		assert.NoErr(t, err)

		assert.NoErr(t, ob.AwaitAsyncCompletion())
		count, err = box.Count()
		assert.NoErr(t, err)
		assert.True(t, count == 0)
	}
}
