/*
 * Copyright 2018 ObjectBox Ltd. All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package objectbox

/*
#cgo LDFLAGS: -lobjectbox
#include <stdlib.h>
#include "objectbox.h"
*/
import "C"

import (
	"fmt"
	"sync"
	"unsafe"
)

// Builder provides tools to fully configure and construct ObjectBox
type Builder struct {
	model *Model
	Error error

	name        string
	maxSizeInKb uint64
	maxReaders  uint
}

func NewBuilder() *Builder {
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

// Directory configures the path where the database is stored
func (builder *Builder) Directory(name string) *Builder {
	builder.name = name
	return builder
}

// MaxSizeInKb defines maximum size the database can take on disk (default: 1 GByte).
func (builder *Builder) MaxSizeInKb(maxSizeInKb uint64) *Builder {
	builder.maxSizeInKb = maxSizeInKb
	return builder
}

// Maximum (concurrent) readers (default: 126). Increase only if you are getting errors (highly concurrent scenarios).
func (builder *Builder) MaxReaders(maxReaders uint) *Builder {
	builder.maxReaders = maxReaders
	return builder
}

// Model specifies schema for the database
//
// Pass the result of the generated function ObjectBoxModel as an argument: Model(ObjectBoxModel())
func (builder *Builder) Model(model *Model) *Builder {
	if builder.Error != nil {
		return builder
	}

	builder.Error = model.validate()
	if builder.Error != nil {
		builder.model = nil
	} else {
		builder.model = model
	}

	return builder
}

// Build validates the configuration and tries to init the ObjectBox
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

	// TODO move to objectbox.go?
	return &ObjectBox{
		store:          cStore,
		bindingsById:   builder.model.bindingsById,
		bindingsByName: builder.model.bindingsByName,
		boxes:          make(map[TypeId]*Box, len(builder.model.bindingsById)),
		boxesMutex:     &sync.RWMutex{},
	}, nil
}
