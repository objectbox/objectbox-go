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
	"fmt"
	"testing"

	"github.com/objectbox/objectbox-go/test/assert"

	"github.com/objectbox/objectbox-go/test/model"
)

func TestRelationsInsert(t *testing.T) {
	// run once for Put & once for PutAsync
	for i := 0; i <= 1; i++ {
		if i == 1 {
			// TODO Box.PutAsync currently doesn't support relations
			continue
		}

		// starting the environment inside the for loop ensures the database is empty & IDs start from 0
		var env = model.NewTestEnv(t)
		var relBox = model.BoxForTestEntityRelated(env.ObjectBox)
		var relValueBox = model.BoxForEntityByValue(env.ObjectBox)

		var id uint64
		var err error

		// this object is used in many-to-many (slice) & one-to-many relations and is inserted only once
		var relReused = &model.TestEntityRelated{Name: "Ptr", NextSlice: []model.EntityByValue{}}
		var object = &model.Entity{
			Related:         model.TestEntityRelated{Name: "Val", NextSlice: []model.EntityByValue{}},
			RelatedPtr:      relReused,
			RelatedSlice:    []model.EntityByValue{{}, {}},
			RelatedPtrSlice: []*model.TestEntityRelated{relReused, {Name: "New", NextSlice: []model.EntityByValue{}}},
		}

		if i == 0 {
			id, err = env.Box.Put(object)
		} else {
			id, err = env.Box.PutAsync(object)
			env.ObjectBox.AwaitAsyncCompletion()
		}

		assert.NoErr(t, err)
		assert.Eq(t, uint64(1), id)
		assert.Eq(t, uint64(1), object.Related.Id)
		assert.Eq(t, uint64(2), object.RelatedPtr.Id)
		assert.Eq(t, uint64(2), object.RelatedPtr.Id)
		assert.Eq(t, uint64(1), object.RelatedSlice[0].Id)
		assert.Eq(t, uint64(2), object.RelatedSlice[1].Id)
		assert.Eq(t, uint64(2), object.RelatedPtrSlice[0].Id)
		assert.Eq(t, uint64(3), object.RelatedPtrSlice[1].Id)

		// check that the relations are inserted and are the same as the ones in the object
		rels, err := relBox.GetAll()
		assert.NoErr(t, err)
		assert.Eq(t, 3, len(rels))
		assert.Eq(t, *rels[0], object.Related)
		assert.Eq(t, *rels[1], *object.RelatedPtr)
		assert.Eq(t, *rels[1], *object.RelatedPtrSlice[0])
		assert.Eq(t, *rels[2], *object.RelatedPtrSlice[1])

		relsV, err := relValueBox.GetAll()
		assert.NoErr(t, err)
		assert.Eq(t, 2, len(relsV))
		assert.Eq(t, relsV[0], object.RelatedSlice[0])
		assert.Eq(t, relsV[1], object.RelatedSlice[1])

		// try to read the object and validate it's read correctly with relations assigned
		objectRead, err := env.Box.Get(id)
		assert.NoErr(t, err)
		assert.Eq(t, object.Related, objectRead.Related)
		assert.Eq(t, object.RelatedPtr, objectRead.RelatedPtr)
		assert.Eq(t, object.RelatedPtr2, objectRead.RelatedPtr2) // this one is empty
		assert.Eq(t, object.RelatedSlice, objectRead.RelatedSlice)
		assert.Eq(t, object.RelatedPtrSlice, objectRead.RelatedPtrSlice)

		// remove one target of each relation, read the object and check everything looks as expected (relations are removed)
		assert.NoErr(t, relBox.Remove(&object.Related))
		assert.NoErr(t, relBox.Remove(object.RelatedPtr))
		assert.NoErr(t, relValueBox.Remove(&object.RelatedSlice[0]))
		objectRead, err = env.Box.Get(id)

		fmt.Println(rels[0])
		assert.Eq(t, uint64(0), objectRead.Related.Id)
		assert.Eq(t, true, objectRead.RelatedPtr == nil)
		assert.Eq(t, 1, len(objectRead.RelatedSlice))
		assert.Eq(t, relsV[1], objectRead.RelatedSlice[0])
		assert.Eq(t, 1, len(objectRead.RelatedPtrSlice))
		assert.Eq(t, rels[2], objectRead.RelatedPtrSlice[0])
	}
}

func TestRelationsUpdate(t *testing.T) {
	// run once for Put & once for PutAsync
	for i := 0; i <= 1; i++ {
		if i == 1 {
			// TODO Box.PutAsync currently doesn't support relations
			continue
		}

		// starting the environment inside the for loop ensures the database is empty & IDs start from 0
		var env = model.NewTestEnv(t).SetOptions(model.TestEnvOptions{PopulateRelations: true})
		var relBox = model.BoxForTestEntityRelated(env.ObjectBox)
		var relValueBox = model.BoxForEntityByValue(env.ObjectBox)
		var err error

		// update the object
		var update = func(object *model.Entity) {
			if i == 0 {
				_, err = env.Box.Put(object)
			} else {
				_, err = env.Box.PutAsync(object)
				env.ObjectBox.AwaitAsyncCompletion()
			}

			assert.NoErr(t, err)
		}

		// add a few test entities (including relations)
		var baseCount uint = 10
		env.Populate(baseCount)

		// check that the relations are inserted correctly
		count, err := relBox.Count()
		assert.NoErr(t, err)
		assert.Eq(t, uint64(baseCount*3), count)

		count, err = relValueBox.Count()
		assert.NoErr(t, err)
		assert.Eq(t, uint64(baseCount*5), count)

		// get one of the entities
		var id uint64 = 2
		object, err := env.Box.Get(id)
		assert.NoErr(t, err)
		assert.Eq(t, uint64(4), object.Related.Id)
		assert.Eq(t, uint64(5), object.RelatedPtr.Id)
		assert.Eq(t, 1, len(object.RelatedSlice))
		assert.Eq(t, 1, len(object.RelatedPtrSlice))
		assert.Eq(t, uint64(8), object.RelatedSlice[0].Id)
		assert.Eq(t, uint64(6), object.RelatedPtrSlice[0].Id)

		// add new (non-existent) items to many-to-many relations
		object.RelatedSlice = append(object.RelatedSlice, model.EntityByValue{})
		object.RelatedPtrSlice = append(object.RelatedPtrSlice, &model.TestEntityRelated{})

		// add existing items to many-to-many relations
		relVal, err := relValueBox.Get(3)
		assert.NoErr(t, err)
		object.RelatedSlice = append(object.RelatedSlice, *relVal)
		object.RelatedPtrSlice = append(object.RelatedPtrSlice, object.RelatedPtr)
		update(object)

		// check if it was updated correctly
		object, err = env.Box.Get(id)
		assert.NoErr(t, err)
		assert.Eq(t, uint64(4), object.Related.Id)
		assert.Eq(t, uint64(5), object.RelatedPtr.Id)
		assert.Eq(t, 3, len(object.RelatedSlice))
		assert.Eq(t, 3, len(object.RelatedPtrSlice))
		assert.Eq(t, []uint64{3, 8, 51}, []uint64{object.RelatedSlice[0].Id, object.RelatedSlice[1].Id, object.RelatedSlice[2].Id})
		assert.EqItems(t, []uint64{5, 6, 31}, []uint64{object.RelatedPtrSlice[0].Id, object.RelatedPtrSlice[1].Id, object.RelatedPtrSlice[2].Id})

		// remove some relations
		object.RelatedPtr = nil
		object.RelatedSlice = object.RelatedSlice[1:]
		object.RelatedPtrSlice = object.RelatedPtrSlice[1:]
		update(object)

		// check if it was updated correctly
		object, err = env.Box.Get(id)
		assert.NoErr(t, err)
		assert.Eq(t, uint64(4), object.Related.Id)
		assert.Eq(t, true, nil == object.RelatedPtr)
		assert.Eq(t, 2, len(object.RelatedSlice))
		assert.Eq(t, 2, len(object.RelatedPtrSlice))
		assert.EqItems(t, []uint64{8, 51}, []uint64{object.RelatedSlice[0].Id, object.RelatedSlice[1].Id})
		assert.EqItems(t, []uint64{6, 31}, []uint64{object.RelatedPtrSlice[0].Id, object.RelatedPtrSlice[1].Id})
	}
}
