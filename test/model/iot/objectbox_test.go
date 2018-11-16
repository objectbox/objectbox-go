package iot

import (
	"testing"

	"github.com/objectbox/objectbox-go/test/assert"
)

func TestObjectBoxEvents(t *testing.T) {
	objectBox := CreateObjectBox()

	box := BoxForEvent(objectBox)
	assert.NoErr(t, box.RemoveAll())

	event := Event{
		Device: "my device",
	}
	objectId, err := box.Put(&event)
	assert.NoErr(t, err)
	t.Logf("Added object ID %v", objectId)

	event.Device = "2nd device"
	objectId, err = box.Put(&event)
	assert.NoErr(t, err)
	t.Logf("Added 2nd object ID %v", objectId)

	eventRead, err := box.Get(objectId)
	assert.NoErr(t, err)
	if objectId != eventRead.Id || event.Device != eventRead.Device {
		t.Fatalf("Event data error: %v vs. %v", event, eventRead)
	}

	all, err := box.GetAll()
	assert.NoErr(t, err)
	assert.EqInt(t, 2, len(all))
}
