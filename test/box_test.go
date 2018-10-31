package objectbox_test

import (
	"fmt"
	"testing"

	"github.com/objectbox/objectbox-go/objectbox"
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

//go:generate objectbox-bindings
type Person struct {
	Id        uint64 `id`
	FirstName string
	LastName  string
}

func TestGeneratedBox(t *testing.T) {
	ob := iot.CreateObjectBox()

	box := BoxForPerson(ob)

	// Create
	id, _ := box.Put(&Person{
		FirstName: "Joe",
		LastName:  "Green",
	})

	person, _ := box.Get(id) // Read
	person.LastName = "Black"
	box.Put(person)    // Update
	box.Remove(person) // Delete
}

// TODO generated code follows
type PersonBox struct {
	objectbox.Box
	EntityId uint32
}

func BoxForPerson(ob *objectbox.ObjectBox) *PersonBox {
	return &PersonBox{
		EntityId: uint32(1),
	}
}

func (box *PersonBox) Get(id uint64) (person *Person, err error) {
	object, err := box.Box.Get(id)
	if err != nil {
		return nil, err
	}
	person, ok := object.(*Person)
	if !ok {
		return nil, fmt.Errorf("could not convert object to Person, invalid type")
	}
	return person, nil
}

func (box *PersonBox) Remove(person *Person) (err error) {
	return box.Box.Remove(person.Id)
}
