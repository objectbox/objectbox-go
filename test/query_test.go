package objectbox_test

import (
	"github.com/objectbox/objectbox-go/objectbox"
	"github.com/objectbox/objectbox-go/test/assert"
	"github.com/objectbox/objectbox-go/test/model/iot"
	"github.com/objectbox/objectbox-go/test/model/iot/object"
	"testing"
)

func TestQueryBuilder(t *testing.T) {
	objectBox := iot.CreateObjectBox()
	box := objectBox.Box(1)
	box.RemoveAll()

	qb, err := objectBox.Query(1)
	if err != nil {
		t.Fatalf("Could not build QB: %v", err)
	}
	query, err := qb.BuildAndDestroy()
	if err != nil {
		t.Fatalf("Could not build query: %v", err)
	}
	defer query.Destroy()

	objectBox.RunWithCursor(1, true, func(cursor *objectbox.Cursor) (err error) {
		bytesArray, err := query.FindBytes(cursor)
		if err != nil {
			t.Fatalf("Could not run query: %v", err)
		}
		assert.EqInt(0, len(bytesArray.BytesArray))

		slice, err := query.Find(cursor)
		if err != nil {
			t.Fatalf("Could not run query: %v", err)
		}
		// TODO should be empty slice instead of nil
		if slice != nil && len(slice.([]object.Event)) != 0 {
			t.Fatalf("unexpected size")
		}
		return
	})

	event := object.Event{
		Device: "dev1",
	}
	id1, err := box.Put(&event)
	if err != nil {
		t.Fatalf("Could not put event: %v", err)
	}

	event.Device = "dev2"
	id2, err := box.Put(&event)
	if err != nil {
		t.Fatalf("Could not add 2nd event: %v", err)
	}

	objectBox.RunWithCursor(1, true, func(cursor *objectbox.Cursor) (err error) {
		bytesArray, err := query.FindBytes(cursor)
		if err != nil {
			t.Fatalf("Could not run query: %v", err)
		}
		assert.EqInt(2, len(bytesArray.BytesArray))

		slice, err := query.Find(cursor)
		if err != nil {
			t.Fatalf("Could not run query: %v", err)
		}
		events := slice.([]object.Event)
		if len(events) != 2 {
			t.Fatalf("unexpected size")
		}

		assert.EqInt(int(id1), int(events[0].Id))
		assert.EqString("dev1", events[0].Device)

		assert.EqInt(int(id2), int(events[1].Id))
		assert.EqString("dev2", events[1].Device)

		return
	})

}
