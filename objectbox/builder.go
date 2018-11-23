package objectbox

/*
#cgo LDFLAGS: -lobjectbox
#include <stdlib.h>
#include "objectbox.h"
*/
import "C"

import (
	"strconv"
	"unsafe"
)

type Builder struct {
	model *Model
	Err   error

	name        string
	maxSizeInKb uint64
	maxReaders  uint

	lastEntityId  TypeId
	lastEntityUid uint64

	lastIndexId  TypeId
	lastIndexUid uint64

	lastRelationId  TypeId
	lastRelationUid uint64

	bindingsById   map[TypeId]ObjectBinding
	bindingsByName map[string]ObjectBinding
}

func NewObjectBoxBuilder() (builder *Builder) {
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
	builder = &Builder{}
	builder.model = model
	builder.bindingsById = make(map[TypeId]ObjectBinding)
	builder.bindingsByName = make(map[string]ObjectBinding)
	return
}

func (builder *Builder) Name(name string) *Builder {
	builder.name = name
	return builder
}

func (builder *Builder) MaxSizeInKb(maxSizeInKb uint64) *Builder {
	builder.maxSizeInKb = maxSizeInKb
	return builder
}

func (builder *Builder) MaxReaders(maxReaders uint) *Builder {
	builder.maxReaders = maxReaders
	return builder
}

func (builder *Builder) RegisterBinding(binding ObjectBinding) {
	binding.AddToModel(builder.model)
	id := builder.model.previousEntityId
	name := builder.model.previousEntityName
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

func (builder *Builder) LastEntityId(id TypeId, uid uint64) *Builder {
	builder.lastEntityId = id
	builder.lastEntityUid = uid
	return builder
}

func (builder *Builder) LastIndexId(id TypeId, uid uint64) *Builder {
	builder.lastIndexId = id
	builder.lastIndexUid = uid
	return builder
}

func (builder *Builder) LastRelationId(id TypeId, uid uint64) *Builder {
	builder.lastRelationId = id
	builder.lastRelationUid = uid
	return builder
}

func (builder *Builder) Build() (objectBox *ObjectBox, err error) {
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

	if builder.lastIndexId != 0 {
		builder.model.LastIndexId(builder.lastIndexId, builder.lastIndexUid)
	}
	if builder.lastRelationId != 0 {
		builder.model.LastRelationId(builder.lastRelationId, builder.lastRelationUid)
	}

	coptions := C.struct_OBX_store_options{}
	if builder.name != "" {
		cname := C.CString(builder.name)
		defer C.free(unsafe.Pointer(cname))
		coptions.directory = cname
	}
	coptions.maxReaders = C.uint(builder.maxReaders)            // Zero is the default on both sides
	coptions.maxDbSizeInKByte = C.uint64_t(builder.maxSizeInKb) // Zero is the default on both sides

	objectBox = &ObjectBox{}
	objectBox.store = C.obx_store_open(builder.model.model, &coptions)
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
