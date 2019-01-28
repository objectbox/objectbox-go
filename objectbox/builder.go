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

	options
}

func NewBuilder() *Builder {
	// these constants are based on the objectbox.h file, not on the loaded library
	var obxMinVersion = Version{C.OBX_VERSION_MAJOR, C.OBX_VERSION_MINOR, C.OBX_VERSION_PATCH}

	if !C.obx_version_is_at_least(C.int(obxMinVersion.Major), C.int(obxMinVersion.Minor), C.int(obxMinVersion.Patch)) {
		var version string
		msg := C.obx_version_string()
		if msg == nil {
			version = "unknown"
		} else {
			version = C.GoString(msg)
		}
		panic("Minimum libobjectbox version " + obxMinVersion.String() + " required, but found " + version + ":\n." +
			">>> Please run install.sh for a full upgrade <<<\n" +
			"Or check https://github.com/objectbox/objectbox-c for info about the required library.")
	}
	return &Builder{
		options: options{
			// defaults
			putAsyncTimeout:  10000, // 10s
			alwaysAwaitAsync: false,
		},
	}
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

// Configures PutAsync enqueue timeout (default is 10 seconds). See Box.PutAsync method doc for more information.
func (builder *Builder) PutAsyncTimeout(milliseconds uint) *Builder {
	builder.putAsyncTimeout = milliseconds
	return builder
}

// Enables automatic waiting for async operations between executing a synchronous one.
// This can be replaced if you're using PutAsync in many places and need to make sure the operation has finished
// before your data you read/query/delete,... is executed. Calls ObjectBox.AwaitAsyncCompletion() internally.
func (builder *Builder) AlwaysAwaitAsync(value bool) *Builder {
	builder.alwaysAwaitAsync = value
	return builder
}

// Model specifies schema for the database.
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

// Build validates the configuration and tries to init the ObjectBox.
// This call panics on failures; if ObjectBox is optional for your app, consider BuildOrError().
func (builder *Builder) Build() (*ObjectBox, error) {
	objectBox, err := builder.BuildOrError()
	if err != nil {
		panic(fmt.Sprintf("Could not create ObjectBox - please check configuration: %s", err))
	}
	return objectBox, nil
}

// Build validates the configuration and tries to init the ObjectBox.
func (builder *Builder) BuildOrError() (*ObjectBox, error) {
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

	ob := &ObjectBox{
		store:          cStore,
		bindingsById:   builder.model.bindingsById,
		bindingsByName: builder.model.bindingsByName,
		boxes:          make(map[TypeId]*Box, len(builder.model.bindingsById)),
		boxesMutex:     &sync.Mutex{},
		entities:       make(map[TypeId]*entity),
		options:        builder.options,
	}

	for id := range builder.model.bindingsById {
		ob.entities[id] = &entity{id: id, objectBox: ob}
	}
	return ob, nil
}
