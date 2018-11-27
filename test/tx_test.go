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

	"github.com/objectbox/objectbox-go/objectbox"
	"github.com/objectbox/objectbox-go/test/model/iot"
)

func TestTransactionInsert(t *testing.T) {
	ob := iot.CreateObjectBox()
	defer ob.Close()

	assert.NoErr(t, ob.Box(1).RemoveAll())

	var insert = uint64(1000000)

	testObx := objectbox.InternalTestAccessObjectBox{ObjectBox: ob}
	assert.NoErr(t, testObx.RunInTxn(false, func(tx *objectbox.Transaction) (err error) {
		cursor, err := tx.CursorForName("Event")
		assert.NoErr(t, err)

		for i := insert; i > 0; i-- {
			_, err := cursor.Put(&iot.Event{})
			assert.NoErr(t, err)
		}
		return nil
	}))

	count, err := ob.Box(1).Count()
	assert.NoErr(t, err)

	assert.Eq(t, insert, count)
}
