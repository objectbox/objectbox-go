package objectbox

/*
#cgo LDFLAGS: -lobjectbox
#include <stdlib.h>
#include "objectbox.h"
*/
import "C"

import (
	"fmt"
	"unsafe"
)

type Builder struct {
	model *Model
	Error error

	name        string
	maxSizeInKb uint64
	maxReaders  uint
}

func NewObjectBoxBuilder() *Builder {
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
	return &Builder{}
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

func (builder *Builder) Model(model *Model, err error) *Builder {
	if err == nil {
		err = model.validate()
	}

	if err != nil {
		builder.model = nil
		builder.Error = err
	} else {
		builder.model = model
	}
	return builder
}

func (builder *Builder) Build() (*ObjectBox, error) {
	if builder.Error != nil {
		return nil, builder.Error
	}

	if builder.model == nil {
		return nil, fmt.Errorf("model is not defined")
	}

	cOptions := C.struct_OBX_store_options{}

	if builder.name != "" {
		cname := C.CString(builder.name)
		defer C.free(unsafe.Pointer(cname))
		cOptions.directory = cname
	}

	cOptions.maxReaders = C.uint(builder.maxReaders)            // Zero is the default on both sides
	cOptions.maxDbSizeInKByte = C.uint64_t(builder.maxSizeInKb) // Zero is the default on both sides

	cStore := C.obx_store_open(builder.model.model, &cOptions)
	if cStore == nil {
		return nil, createError()
	}

	return &ObjectBox{
		store:          cStore,
		bindingsById:   builder.model.bindingsById,
		bindingsByName: builder.model.bindingsByName,
	}, nil
}
