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
	"testing"

	"github.com/objectbox/objectbox-go/test/assert"

	"github.com/objectbox/objectbox-go/test/model"
)

func TestRelations(t *testing.T) {
	// run once for Put & once for PutAsync
	for i := 0; i <= 1; i++ {
		// starting the environment inside the for loop ensures the database is empty & IDs start from 0
		var env = model.NewTestEnv(t)
		var relBox = model.BoxForTestEntityRelated(env.ObjectBox)

		var id uint64
		var err error

		var entity = &model.Entity{RelatedPtr: &model.TestEntityRelated{Name: "Ptr"}, Related: model.TestEntityRelated{Name: "Val"}}

		if i == 0 {
			id, err = env.Box.Put(entity)
		} else {
			id, err = env.Box.PutAsync(entity)
			env.ObjectBox.AwaitAsyncCompletion()
		}

		assert.NoErr(t, err)
		assert.Eq(t, id, uint64(1))
		assert.Eq(t, entity.Related.Id, uint64(1))
		assert.Eq(t, entity.RelatedPtr.Id, uint64(2))

		// check that the relations are inserted and are the same as the ones in the entity
		rels, err := relBox.GetAll()
		assert.NoErr(t, err)
		assert.Eq(t, *rels[0], entity.Related)
		assert.Eq(t, *rels[1], *entity.RelatedPtr)

		// try to read the entity and validate it's read correctly with relations assigned
		entityRead, err := env.Box.Get(id)
		assert.NoErr(t, err)
		assert.Eq(t, entity.Related, entityRead.Related)
		assert.Eq(t, entity.RelatedPtr, entityRead.RelatedPtr)
		assert.Eq(t, entity.RelatedPtr2, entityRead.RelatedPtr2)

		// remove one relation and try to read the entity
		assert.NoErr(t, relBox.Remove(&entity.Related))
		assert.NoErr(t, relBox.Remove(entity.RelatedPtr))
		entityRead, err = env.Box.Get(id)
		assert.Eq(t, entityRead.Related.Id, uint64(0))
		assert.Eq(t, entityRead.RelatedPtr == nil, true)
	}
}
