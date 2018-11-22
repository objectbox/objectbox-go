package objectbox

/*
#cgo LDFLAGS: -lobjectbox
#include <stdlib.h>
#include "objectbox.h"
*/
import "C"

import (
	"strconv"
)

type ObjectBoxBuilder struct {
	name          string
	model         *Model
	Err           error
	lastEntityId  TypeId
	lastEntityUid uint64

	bindingsById   map[TypeId]ObjectBinding
	bindingsByName map[string]ObjectBinding
}

func NewObjectBoxBuilder() (builder *ObjectBoxBuilder) {
	if !C.obx_version_is_at_least(0, 3, 0) {
		var version string
		msg := C.obx_version_string()
		if msg == nil {
			version = "unknown"
		} else {
			version = C.GoString(msg)
		}
		panic("Minimum libobjectbox version 0.3.0 required, but found " + version +
			". Check https://github.com/objectbox/objectbox-c for updates.")
	}
	model, err := NewModel()
	if err != nil {
		panic("Could not create model: " + err.Error())
	}
	builder = &ObjectBoxBuilder{}
	builder.model = model
	builder.bindingsById = make(map[TypeId]ObjectBinding)
	builder.bindingsByName = make(map[string]ObjectBinding)
	return
}

func (builder *ObjectBoxBuilder) Name(name string) *ObjectBoxBuilder {
	builder.name = name
	return builder
}

func (builder *ObjectBoxBuilder) RegisterBinding(binding ObjectBinding) {
	binding.AddToModel(builder.model)
	id := builder.model.lastEntityId
	name := builder.model.lastEntityName
	if id == 0 {
		panic("No type ID; did you forget to add an entity to the model?")
	}
	if name == "" {
		panic("No type name")
	}
	existingBinding := builder.bindingsById[id]
	if existingBinding != nil {
		panic("Already registered a binding for ID " + strconv.Itoa(int(id)))
	}
	existingBinding = builder.bindingsByName[name]
	if existingBinding != nil {
		panic("Already registered a binding for name " + name)
	}
	builder.bindingsById[id] = binding
	builder.bindingsByName[name] = binding
}

func (builder *ObjectBoxBuilder) LastEntityId(id TypeId, uid uint64) *ObjectBoxBuilder {
	builder.lastEntityId = id
	builder.lastEntityUid = uid
	return builder
}

func (builder *ObjectBoxBuilder) Build() (objectBox *ObjectBox, err error) {
	if builder.model.Err != nil {
		err = builder.model.Err
		return
	}
	if builder.Err != nil {
		err = builder.Err
		return
	}
	if builder.lastEntityId == 0 || builder.lastEntityUid == 0 {
		panic("Configuration error: last entity ID/UID must be set")
	}
	builder.model.LastEntityId(builder.lastEntityId, builder.lastEntityUid)

	// TODO implement or remove
	//fmt.Println("Ignoring DB name: " + builder.name)
	//cname := C.CString(builder.name)
	//defer C.free(unsafe.Pointer(cname))

	objectBox = &ObjectBox{}
	objectBox.store = C.obx_store_open(builder.model.model, nil)
	if objectBox.store == nil {
		objectBox = nil
		err = createError()
	}
	if err == nil {
		objectBox.bindingsById = builder.bindingsById
		objectBox.bindingsByName = builder.bindingsByName
	}
	return
}
