/*
 * Copyright 2019 ObjectBox Ltd. All rights reserved.
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
#include <stdlib.h>
#include "objectbox.h"
*/
import "C"

import (
	"fmt"
	"runtime"
	"unsafe"
)

// Builder provides tools to fully configure and construct ObjectBox
type Builder struct {
	model *Model
	Error error

	// these options are used when creating the underlying store using the C-api calls
	// pointers are used to distinguish whether a value is present or not
	directory   *string
	maxSizeInKb *uint64
	maxReaders  *uint

	// these options are passed-through to the created ObjectBox struct
	options
}

// NewBuilder creates a new ObjectBox instance builder object
func NewBuilder() *Builder {
	// these constants are based on the objectbox.h file, not on the loaded library
	var expectedVersion = Version{C.OBX_VERSION_MAJOR, C.OBX_VERSION_MINOR, C.OBX_VERSION_PATCH, ""}

	// NOTE, this is currently not used as the C-API currently still breaks compatibility before v1.0.
	//if !C.obx_version_is_at_least(C.int(expectedVersion.Major), C.int(expectedVersion.Minor), C.int(expectedVersion.Patch)) {
	//	var version string
	//	msg := C.obx_version_string()
	//	if msg == nil {
	//		version = "unknown"
	//	} else {
	//		version = C.GoString(msg)
	//	}
	//	panic("Minimum libobjectbox version " + expectedVersion.String() + " required, but found " + version + "\n" +
	//		"Please run install.sh for a full upgrade\n" +
	//		"Or check https://github.com/objectbox/objectbox-c for info about the required library.")
	//}

	// Therefore, we're checking the exact version match
	var version = VersionLib()
	// patch version can be higher but only if it's not a "dev" version (100+)
	var patchMatches = version.Patch >= expectedVersion.Patch && (version.Patch < 100 || version.Patch == expectedVersion.Patch)
	if version.Major != expectedVersion.Major || version.Minor != expectedVersion.Minor || !patchMatches {
		panic("libobjectbox version " + expectedVersion.String() + " required, but found " + version.String() + ".\n" +
			"Please run install.sh for a full upgrade.\n" +
			"Or check https://github.com/objectbox/objectbox-c for info about the required library.")
	}

	return &Builder{
		options: options{
			// defaults
			asyncTimeout: 1000, // 1s ; TODO make this 0 to use core default?
		},
	}
}

// Directory configures the path where the database is stored
func (builder *Builder) Directory(path string) *Builder {
	builder.directory = &path
	return builder
}

// MaxSizeInKb defines maximum size the database can take on disk (default: 1 GByte).
func (builder *Builder) MaxSizeInKb(maxSizeInKb uint64) *Builder {
	builder.maxSizeInKb = &maxSizeInKb
	return builder
}

// MaxReaders defines maximum concurrent readers (default: 126).
// Increase only if you are getting errors (highly concurrent scenarios).
func (builder *Builder) MaxReaders(maxReaders uint) *Builder {
	builder.maxReaders = &maxReaders
	return builder
}

// asyncTimeoutTBD configures the default enqueue timeout for async operations (default is 1 second).
// See Box.PutAsync method doc for more information.
// TODO: implement this option in core and use it
func (builder *Builder) asyncTimeoutTBD(milliseconds uint) *Builder {
	builder.asyncTimeout = milliseconds
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

// BuildOrError validates the configuration and tries to init the ObjectBox.
func (builder *Builder) BuildOrError() (*ObjectBox, error) {
	if builder.Error != nil {
		return nil, builder.Error
	}

	if builder.model == nil {
		return nil, fmt.Errorf("model is not defined")
	}

	// for native calls/createError()
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	cOptions := C.obx_opt()
	if cOptions == nil {
		return nil, createError()
	}

	if builder.directory != nil {
		cDir := C.CString(*builder.directory)
		defer C.free(unsafe.Pointer(cDir))
		if 0 != C.obx_opt_directory(cOptions, cDir) {
			C.obx_opt_free(cOptions)
			return nil, createError()
		}
	}

	if builder.maxSizeInKb != nil {
		C.obx_opt_max_db_size_in_kb(cOptions, C.size_t(*builder.maxSizeInKb))
	}

	if builder.maxReaders != nil {
		C.obx_opt_max_readers(cOptions, C.uint(*builder.maxReaders))
	}

	C.obx_opt_model(cOptions, builder.model.cModel)

	// cOptions is consumed by obx_store_open() so no need to free it
	cStore := C.obx_store_open(cOptions)
	if cStore == nil {
		return nil, createError()
	}

	ob := &ObjectBox{
		store:          cStore,
		entitiesById:   builder.model.entitiesById,
		entitiesByName: builder.model.entitiesByName,
		boxes:          make(map[TypeId]*Box, len(builder.model.entitiesById)),
		options:        builder.options,
	}

	for _, entity := range builder.model.entitiesById {
		entity.objectBox = ob
	}
	return ob, nil
}
