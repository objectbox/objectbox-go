package objectbox_test

import (
	"testing"

	"github.com/objectbox/objectbox-go/test/assert"

	"github.com/objectbox/objectbox-go/objectbox"
	"github.com/objectbox/objectbox-go/test/model/iot"
)

func TestTransactionInsert(t *testing.T) {
	ob := iot.CreateObjectBox()

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
