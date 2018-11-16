package iot

import (
	"strconv"

	. "github.com/objectbox/objectbox-go/objectbox"
	"github.com/objectbox/objectbox-go/test/assert"
)

func CreateObjectBox() *ObjectBox {
	builder := NewObjectBoxBuilder().Name("iot-test").LastEntityId(2, 5284076134434938613)
	//objectBox.SetDebugFlags(DebugFlags_LOG_ASYNC_QUEUE)
	builder.RegisterBinding(EventBinding{})
	builder.RegisterBinding(ReadingBinding{})
	objectBox, err := builder.Build()
	if err != nil {
		panic(err)
	}
	return objectBox
}

func PutEvent(ob *ObjectBox, device string, date int64) *Event {
	event := Event{Device: device, Date: date}
	id, err := ob.Box(1).Put(&event)
	assert.NoErr(nil, err)
	event.Id = id
	return &event
}

func PutEvents(ob *ObjectBox, count int) []*Event {
	// TODO TX
	events := make([]*Event, 0, count)
	for i := 1; i <= count; i++ {
		event := PutEvent(ob, "device "+strconv.Itoa(i), int64(10000+i))
		events = append(events, event)
	}
	return events
}
