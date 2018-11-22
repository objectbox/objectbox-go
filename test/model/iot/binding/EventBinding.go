package binding

import (
	"github.com/google/flatbuffers/go"
	. "github.com/objectbox/objectbox-go/objectbox"
	. "github.com/objectbox/objectbox-go/test/model/iot/flat"
	"github.com/objectbox/objectbox-go/test/model/iot/object"
)

// This file could be generated in the future
type EventBinding struct {
}

func (EventBinding) AddToModel(model *Model) {
	model.Entity("Event", 1, 10001)
	model.Property("id", PropertyType_Long, 1, 10001001)
	model.PropertyFlags(PropertyFlags_ID)
	model.Property("device", PropertyType_String, 2, 10001002)
	model.Property("date", PropertyType_Date, 3, 10001003)
	model.EntityLastPropertyId(3, 10001003)
}

func (EventBinding) GetId(entity interface{}) (id uint64, err error) {
	return entity.(*object.Event).Id, nil
}

func (EventBinding) Flatten(entity interface{}, fbb *flatbuffers.Builder, id uint64) {
	flattenModelEvent(entity.(*object.Event), fbb, id)
}

func flattenModelEvent(event *object.Event, fbb *flatbuffers.Builder, id uint64) {
	offsetDevice := Unavailable
	if event.Device != "" {
		offsetDevice = fbb.CreateString(event.Device)
	}

	EventStart(fbb)

	EventAddId(fbb, id)
	if offsetDevice != Unavailable {
		EventAddDevice(fbb, offsetDevice)
	}

	EventAddDate(fbb, event.Date)
}

func (EventBinding) ToObject(bytes []byte) interface{} {
	flatEvent := GetRootAsEvent(bytes, flatbuffers.UOffsetT(0))
	return toModelEvent(flatEvent)
}

func toModelEvent(src *Event) *object.Event {
	return &object.Event{
		Id:     src.Id(),
		Device: string(src.Device()),
		Date:   src.Date(),
	}
}

func (EventBinding) MakeSlice(capacity int) interface{} {
	return make([]object.Event, 0, 16)
}

func (EventBinding) AppendToSlice(slice interface{}, entity interface{}) interface{} {
	return append(slice.([]object.Event), *entity.(*object.Event))
}
