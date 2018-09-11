package iot

import (
	. "github.com/objectbox/objectbox-go/objectbox"
	"github.com/objectbox/objectbox-go/test/model/iot/binding"
)

func CreateObjectBox() *ObjectBox {
	builder := NewObjectBoxBuilder().Name("iot-test").LastEntityId(2, 10002)
	//objectBox.SetDebugFlags(DebugFlags_LOG_ASYNC_QUEUE)
	builder.RegisterBinding(binding.EventBinding{})
	builder.RegisterBinding(binding.ReadingBinding{})
	objectBox, err := builder.Build()
	if err != nil {
		panic(err)
	}
	return objectBox
}
