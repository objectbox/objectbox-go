package objectbox_test

import (
	"testing"

	"github.com/objectbox/objectbox-go/test/model/iot"
	. "github.com/objectbox/objectbox-go/test/model/iot/object"
)

func TestAsync(t *testing.T) {
	objectBox := iot.CreateObjectBox()
	box := objectBox.Box(1)
	event := Event{
		Device: "my device",
	}
	objectId, err := box.PutAsync(&event)
	if err != nil {
		t.Fatalf("Could not add event: %v", err)
	}
	objectBox.AwaitAsyncCompletion()

	objectRead, err := box.Get(objectId)
	if err != nil {
		t.Fatalf("Could not get back event by ID: %v", err)
	}
	eventRead := objectRead.(*Event)
	if objectId != eventRead.Id || event.Device != eventRead.Device {
		t.Fatalf("Event data error: %v vs. %v", event, eventRead)
	}
}
