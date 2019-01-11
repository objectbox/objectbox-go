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
	"unsafe"

	"github.com/google/flatbuffers/go"
)

// Internal: won't be publicly exposed in a future version!
type cursor struct {
	cursor  *C.OBX_cursor
	binding ObjectBinding
	fbb     *flatbuffers.Builder
}

func (cursor *cursor) Close() error {
	rc := C.obx_cursor_close(cursor.cursor)
	cursor.cursor = nil
	if rc != 0 {
		return createError()
	}
	return nil
}

func (cursor *cursor) Get(id uint64) (object interface{}, err error) {
	bytes, err := cursor.getBytes(id)
	if bytes == nil || err != nil {
		return
	}
	return cursor.binding.ToObject(bytes), nil
}

func (cursor *cursor) GetAll() (slice interface{}, err error) {
	cBytesArray := C.obx_cursor_get_all(cursor.cursor)
	if cBytesArray == nil {
		return nil, createError()
	}
	return cursor.cBytesArrayToObjects(cBytesArray), nil
}

func (cursor *cursor) getBytes(id uint64) (bytes []byte, err error) {
	var data *C.void
	var dataSize C.size_t
	dataPtr := unsafe.Pointer(data) // Need ptr to an unsafe ptr here
	rc := C.obx_cursor_get(cursor.cursor, C.obx_id(id), &dataPtr, &dataSize)
	if rc != 0 {
		if rc != C.OBX_NOT_FOUND {
			err = createError()
		}
		return
	}
	bytes = C.GoBytes(dataPtr, C.int(dataSize))
	return
}

func (cursor *cursor) First() (bytes []byte, err error) {
	var data *C.void
	var dataSize C.size_t
	dataPtr := unsafe.Pointer(data) // Need ptr to an unsafe ptr here
	rc := C.obx_cursor_first(cursor.cursor, &dataPtr, &dataSize)
	if rc != 0 {
		if rc != C.OBX_NOT_FOUND {
			err = createError()
		}
		return
	}
	bytes = C.GoBytes(dataPtr, C.int(dataSize))
	return
}

func (cursor *cursor) Next() (bytes []byte, err error) {
	var data *C.void
	var dataSize C.size_t
	dataPtr := unsafe.Pointer(data) // Need ptr to an unsafe ptr here
	rc := C.obx_cursor_next(cursor.cursor, &dataPtr, &dataSize)
	if rc != 0 {
		if rc != C.OBX_NOT_FOUND {
			err = createError()
		}
		return
	}
	bytes = C.GoBytes(dataPtr, C.int(dataSize))
	return
}

func (cursor *cursor) Count() (count uint64, err error) {
	var cCount C.uint64_t
	rc := C.obx_cursor_count(cursor.cursor, &cCount)
	if rc != 0 {
		err = createError()
		return
	}
	return uint64(cCount), nil
}

func (cursor *cursor) CountMax(max uint64) (count uint64, err error) {
	var cCount C.uint64_t
	rc := C.obx_cursor_count_max(cursor.cursor, C.uint64_t(max), &cCount)
	if rc != 0 {
		err = createError()
		return
	}
	return uint64(cCount), nil
}

func (cursor *cursor) IsEmpty() (result bool, err error) {
	var cResult C.bool
	rc := C.obx_cursor_is_empty(cursor.cursor, &cResult)
	if rc != 0 {
		err = createError()
		return
	}
	return bool(cResult), nil
}
func (cursor *cursor) Put(object interface{}) (id uint64, err error) {
	idFromObject, err := cursor.binding.GetId(object)
	if err != nil {
		return
	}

	id, err = cursor.IdForPut(idFromObject)
	if err != nil {
		return
	}

	cursor.binding.Flatten(object, cursor.fbb, id)

	checkForPreviousValue := idFromObject != 0
	if err = cursor.finishInternalFbbAndPut(id, checkForPreviousValue); err != nil {
		return 0, err
	}

	// update the id on the object
	if idFromObject != id {
		// TODO SetId never errs
		if err = cursor.binding.SetId(object, id); err != nil {
			return 0, err
		}
	}

	return id, nil
}

func (cursor *cursor) finishInternalFbbAndPut(id uint64, checkForPreviousObject bool) (err error) {
	fbb := cursor.fbb
	fbb.Finish(fbb.EndObject())
	bytes := fbb.FinishedBytes()

	rc := C.obx_cursor_put(cursor.cursor,
		C.obx_id(id), unsafe.Pointer(&bytes[0]), C.size_t(len(bytes)), C.bool(checkForPreviousObject))
	if rc != 0 {
		err = createError()
	}

	// Reset to have a clear state for the next caller
	fbb.Reset()

	return
}

func (cursor *cursor) IdForPut(idCandidate uint64) (id uint64, err error) {
	id = uint64(C.obx_cursor_id_for_put(cursor.cursor, C.obx_id(idCandidate)))
	if id == 0 {
		err = createError()
	}
	return
}

func (cursor *cursor) Remove(id uint64) (err error) {
	rc := C.obx_cursor_remove(cursor.cursor, C.obx_id(id))
	if rc != 0 {
		err = createError()
	}
	return
}

func (cursor *cursor) RemoveAll() (err error) {
	rc := C.obx_cursor_remove_all(cursor.cursor)
	if rc != 0 {
		err = createError()
	}
	return
}

func (cursor *cursor) cBytesArrayToObjects(cBytesArray *C.OBX_bytes_array) (slice interface{}) {
	bytesArray := cBytesArrayToGo(cBytesArray)
	defer bytesArray.free()
	return cursor.bytesArrayToObjects(bytesArray)
}

func (cursor *cursor) bytesArrayToObjects(bytesArray *bytesArray) (slice interface{}) {
	slice = cursor.binding.MakeSlice(len(bytesArray.BytesArray))
	for _, bytesData := range bytesArray.BytesArray {
		object := cursor.binding.ToObject(bytesData)
		slice = cursor.binding.AppendToSlice(slice, object)
	}
	return
}
