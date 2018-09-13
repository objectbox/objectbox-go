package iot

import (
	"github.com/objectbox/objectbox-go/test/assert"
	. "github.com/objectbox/objectbox-go/test/model/iot/object"
	"testing"
)

func TestObjectBoxEvents(t *testing.T) {
	objectBox := CreateObjectBox()
	box := objectBox.Box(1)
	box.RemoveAll()
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

	objectRead, err := box.Get(objectId)
	assert.NoErr(t, err)
	eventRead := objectRead.(*Event)
	if objectId != eventRead.Id || event.Device != eventRead.Device {
		t.Fatalf("Event data error: %v vs. %v", event, eventRead)
	}

	all, err := box.GetAll()
	assert.NoErr(t, err)
	assert.EqInt(t, 2, len(all.([]Event)))
}
