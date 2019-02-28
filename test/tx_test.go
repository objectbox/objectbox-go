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

	"github.com/objectbox/objectbox-go/test/model/iot"
)

func TestTransactionInsert(t *testing.T) {
	ob := iot.LoadEmptyTestObjectBox()
	defer ob.Close()

	assert.NoErr(t, iot.BoxForEvent(ob).RemoveAll())

	// TODO use box with reentrant transactions
	//var insert = uint64(1000000)
	//
	//testObx := objectbox.InternalTestAccessObjectBox{ObjectBox: ob}
	//assert.NoErr(t, testObx.RunInTxn(false, func(tx *objectbox.Transaction) (err error) {
	//	return tx.RunWithCursor(iot.EventBinding.Id, func(cursor *objectbox.Cursor) error {
	//		for i := insert; i > 0; i-- {
	//			_, err := cursor.Put(&iot.Event{})
	//			assert.NoErr(t, err)
	//		}
	//		return nil
	//	})
	//}))
	//
	//count, err := iot.BoxForEvent(ob).Count()
	//assert.NoErr(t, err)
	//
	//assert.Eq(t, insert, count)
}
