package objectbox_test

import (
	"testing"

	"github.com/objectbox/objectbox-go/test/assert"
	"github.com/objectbox/objectbox-go/test/model/iot"
)

func TestAsync(t *testing.T) {
	objectBox := iot.CreateObjectBox()
	box := iot.BoxForEvent(objectBox)

	err := box.RemoveAll()
	assert.NoErr(t, err)

	event := iot.Event{
		Device: "my device",
	}
	objectId, err := box.PutAsync(&event)
	if err != nil {
		t.Fatalf("Could not add event: %v", err)
	}

	objectBox.AwaitAsyncCompletion()

	count, err := box.Count()
	assert.NoErr(t, err)
	assert.Eq(t, uint64(1), count)

	eventRead, err := box.Get(objectId)
	if err != nil {
		t.Fatalf("Could not get back event by ID: %v", err)
	}
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
