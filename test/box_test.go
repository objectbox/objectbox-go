package objectbox_test

import (
	"testing"

	"github.com/objectbox/objectbox-go/test/assert"
	"github.com/objectbox/objectbox-go/test/model/iot"
)

func TestAsync(t *testing.T) {
	objectBox := iot.CreateObjectBox()
	defer objectBox.Close()
	box := iot.BoxForEvent(objectBox)
	defer box.Close()
	err := box.RemoveAll()
	assert.NoErr(t, err)

	event := iot.Event{
		Device: "my device",
	}
	objectId, err := box.PutAsync(&event)
	assert.NoErr(t, err)

	objectBox.AwaitAsyncCompletion()

	count, err := box.Count()
	assert.NoErr(t, err)
	assert.Eq(t, uint64(1), count)

	eventRead, err := box.Get(objectId)
	assert.NoErr(t, err)
	if objectId != eventRead.Id || event.Device != eventRead.Device {
		t.Fatalf("Event data error: %v vs. %v", event, eventRead)
	}

	err = box.Remove(eventRead)
	assert.NoErr(t, err)

	eventRead, err = box.Get(objectId)
	assert.NoErr(t, err)
	if eventRead != nil {
		t.Fatalf("object hasn't been deleted by box.Remove()")
	}

	count, err = box.Count()
	assert.NoErr(t, err)
	assert.Eq(t, uint64(0), count)
}

func TestUnique(t *testing.T) {
	objectBox := iot.CreateObjectBox()
	defer objectBox.Close()
	box := iot.BoxForEvent(objectBox)
	defer box.Close()
	err := box.RemoveAll()
	assert.NoErr(t, err)

	event := iot.Event{
		Device: "my device",
		Uid:    "a",
	}
	_, err = box.Put(&event)
	assert.NoErr(t, err)

	_, err = box.Put(&event)
	if err == nil {
		assert.Failf(t, "put() passed instead of an expected unique constraint violation")
	}

	count, err := box.Count()
	assert.NoErr(t, err)
	assert.Eq(t, uint64(1), count)
}

func TestPutAll(t *testing.T) {
	objectBox := iot.CreateObjectBox()
	defer objectBox.Close()
	box := iot.BoxForEvent(objectBox)
	defer box.Close()
	err := box.RemoveAll()
	assert.NoErr(t, err)

	event1 := iot.Event{
		Device: "Pi 3B",
	}
	event2 := iot.Event{
		Device: "Pi Zero",
	}
	events := []*iot.Event{&event1, &event2}
	objectIds, err := box.PutAll(events)
	assert.NoErr(t, err)

	count, err := box.Count()
	assert.NoErr(t, err)
	assert.Eq(t, uint64(2), count)

	eventRead, err := box.Get(objectIds[0])
	assert.NoErr(t, err)
	assert.EqString(t, "Pi 3B", eventRead.Device)

	eventRead, err = box.Get(objectIds[1])
	assert.NoErr(t, err)
	assert.EqString(t, "Pi Zero", eventRead.Device)

	// And passing nil & empty slice
	objectIds, err = box.PutAll(nil)
	assert.NoErr(t, err)
	assert.EqInt(t, len(objectIds), 0)
	//noinspection GoPreferNilSlice
	noEvents := []*iot.Event{}
	objectIds, err = box.PutAll(noEvents)
	assert.NoErr(t, err)
	assert.EqInt(t, len(objectIds), 0)
}
