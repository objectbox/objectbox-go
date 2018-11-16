package objectbox_test

import (
	"testing"

	"github.com/objectbox/objectbox-go/test/assert"
	"github.com/objectbox/objectbox-go/test/model/iot"
	. "github.com/objectbox/objectbox-go/test/model/iot/object"
)

func TestAsync(t *testing.T) {
	objectBox := iot.CreateObjectBox()
	box := objectBox.Box(1)

	err := box.RemoveAll()
	assert.NoErr(t, err)

	event := Event{
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

	objectRead, err := box.Get(objectId)
	if err != nil {
		t.Fatalf("Could not get back event by ID: %v", err)
	}
	eventRead := objectRead.(*Event)
	if objectId != eventRead.Id || event.Device != eventRead.Device {
		t.Fatalf("Event data error: %v vs. %v", event, eventRead)
	}

	err = box.Remove(objectId)
	assert.NoErr(t, err)

	objectRead, err = box.Get(objectId)
	assert.NoErr(t, err)
	if objectRead != nil {
		t.Fatalf("object hasn't been deleted by box.Remove()")
	}

	count, err = box.Count()
	assert.NoErr(t, err)
	assert.Eq(t, uint64(0), count)

}
