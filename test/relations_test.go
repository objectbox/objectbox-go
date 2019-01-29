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
	"fmt"
	"testing"

	"github.com/objectbox/objectbox-go/test/assert"

	"github.com/objectbox/objectbox-go/test/model"
)

func TestRelationsInsert(t *testing.T) {
	// run once for Put & once for PutAsync
	for i := 0; i <= 1; i++ {
		// starting the environment inside the for loop ensures the database is empty & IDs start from 0
		var env = model.NewTestEnv(t)
		var relBox = model.BoxForTestEntityRelated(env.ObjectBox)
		var relValueBox = model.BoxForEntityByValue(env.ObjectBox)

		var id uint64
		var err error

		// this object is used in many-to-many (slice) & one-to-many relations and is inserted only once
		var relReused = &model.TestEntityRelated{Name: "Ptr"}
		var entity = &model.Entity{
			Related:         model.TestEntityRelated{Name: "Val"},
			RelatedPtr:      relReused,
			RelatedSlice:    []model.EntityByValue{{}, {}},
			RelatedPtrSlice: []*model.TestEntityRelated{relReused, {Name: "New"}},
		}

		if i == 0 {
			id, err = env.Box.Put(entity)
		} else {
			id, err = env.Box.PutAsync(entity)
			env.ObjectBox.AwaitAsyncCompletion()
		}

		assert.NoErr(t, err)
		assert.Eq(t, uint64(1), id)
		assert.Eq(t, uint64(1), entity.Related.Id)
		assert.Eq(t, uint64(2), entity.RelatedPtr.Id)
		assert.Eq(t, uint64(2), entity.RelatedPtr.Id)
		assert.Eq(t, uint64(1), entity.RelatedSlice[0].Id)
		assert.Eq(t, uint64(2), entity.RelatedSlice[1].Id)
		assert.Eq(t, uint64(2), entity.RelatedPtrSlice[0].Id)
		assert.Eq(t, uint64(3), entity.RelatedPtrSlice[1].Id)

		// check that the relations are inserted and are the same as the ones in the entity
		rels, err := relBox.GetAll()
		assert.NoErr(t, err)
		assert.Eq(t, 3, len(rels))
		assert.Eq(t, *rels[0], entity.Related)
		assert.Eq(t, *rels[1], *entity.RelatedPtr)
		assert.Eq(t, *rels[1], *entity.RelatedPtrSlice[0])
		assert.Eq(t, *rels[2], *entity.RelatedPtrSlice[1])

		relsV, err := relValueBox.GetAll()
		assert.NoErr(t, err)
		assert.Eq(t, 2, len(relsV))
		assert.Eq(t, relsV[0], entity.RelatedSlice[0])
		assert.Eq(t, relsV[1], entity.RelatedSlice[1])

		// try to read the entity and validate it's read correctly with relations assigned
		entityRead, err := env.Box.Get(id)
		assert.NoErr(t, err)
		assert.Eq(t, entity.Related, entityRead.Related)
		assert.Eq(t, entity.RelatedPtr, entityRead.RelatedPtr)
		assert.Eq(t, entity.RelatedPtr2, entityRead.RelatedPtr2)
		assert.Eq(t, entity.RelatedSlice, entityRead.RelatedSlice)
		assert.Eq(t, entity.RelatedPtrSlice, entityRead.RelatedPtrSlice)

		// remove one target of each relation, read the entity and check everything looks as expected (relations are removed)
		assert.NoErr(t, relBox.Remove(&entity.Related))
		assert.NoErr(t, relBox.Remove(entity.RelatedPtr))
		assert.NoErr(t, relValueBox.Remove(&entity.RelatedSlice[0]))
		entityRead, err = env.Box.Get(id)

		fmt.Println(rels[0])
		assert.Eq(t, uint64(0), entityRead.Related.Id)
		assert.Eq(t, true, entityRead.RelatedPtr == nil)
		assert.Eq(t, 1, len(entityRead.RelatedSlice))
		assert.Eq(t, relsV[1], entityRead.RelatedSlice[0])
		assert.Eq(t, 1, len(entityRead.RelatedPtrSlice))
		assert.Eq(t, rels[2], entityRead.RelatedPtrSlice[0])
	}
}

// TODO test update of a source entity while inserting/deleting relations
