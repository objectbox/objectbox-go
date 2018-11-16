package objectbox

/*
#cgo LDFLAGS: -lobjectbox
#include <stdlib.h>
#include "objectbox.h"
*/
import "C"

import (
	"sync/atomic"
	"unsafe"

	"github.com/google/flatbuffers/go"
)

type Box struct {
	objectBox *ObjectBox
	box       *C.OBX_box
	typeId    TypeId
	binding   ObjectBinding

	// Must be used in combination with fbbInUseAtomic
	fbb *flatbuffers.Builder

	// Values 0 (unused) or 1 (in use); use only with CompareAndSwapInt32
	fbbInUseAtomic uint32
}

func (box *Box) Destroy() (err error) {
	rc := C.obx_box_close(box.box)
	box.box = nil
	if rc != 0 {
		err = createError()
	}
	return
}

func (box *Box) IdForPut(idCandidate uint64) (id uint64, err error) {
	id = uint64(C.obx_box_id_for_put(box.box, C.uint64_t(idCandidate)))
	if id == 0 {
		err = createError()
	}
	return
}

func (box *Box) PutAsync(object interface{}) (id uint64, err error) {
	idFromObject, err := box.binding.GetId(object)
	if err != nil {
		return
	}
	checkForPreviousValue := idFromObject != 0
	id, err = box.IdForPut(idFromObject)
	if err != nil {
		return
	}

	var fbb *flatbuffers.Builder
	if atomic.CompareAndSwapUint32(&box.fbbInUseAtomic, 0, 1) {
		defer atomic.StoreUint32(&box.fbbInUseAtomic, 0)
		fbb = box.fbb
	} else {
		fbb = flatbuffers.NewBuilder(256)
	}
	box.binding.Flatten(object, fbb, id)
	return id, box.finishFbbAndPutAsync(fbb, id, checkForPreviousValue)
}

func (box *Box) finishFbbAndPutAsync(fbb *flatbuffers.Builder, id uint64, checkForPreviousObject bool) (err error) {
	fbb.Finish(fbb.EndObject())
	bytes := fbb.FinishedBytes()

	rc := C.obx_box_put_async(box.box,
		C.uint64_t(id), unsafe.Pointer(&bytes[0]), C.size_t(len(bytes)), C.bool(checkForPreviousObject))
	if rc != 0 {
		err = createError()
	}

	// Reset to have a clear state for the next caller
	fbb.Reset()

	return
}

func (box *Box) Put(object interface{}) (id uint64, err error) {
	err = box.objectBox.RunWithCursor(box.typeId, false, func(cursor *Cursor) error {
		var errInner error
		id, errInner = cursor.Put(object)
		return errInner
	})
	return
}

func (box *Box) Remove(id uint64) (err error) {
	return box.objectBox.RunWithCursor(box.typeId, false, func(cursor *Cursor) error {
		return cursor.Remove(id)
	})
}

func (box *Box) RemoveAll() (err error) {
	return box.objectBox.RunWithCursor(box.typeId, false, func(cursor *Cursor) error {
		return cursor.RemoveAll()
	})
}

func (box *Box) Count() (count uint64, err error) {
	err = box.objectBox.RunWithCursor(box.typeId, true, func(cursor *Cursor) error {
		var errInner error
		count, errInner = cursor.Count()
		return errInner
	})
	return
}

func (box *Box) Get(id uint64) (object interface{}, err error) {
	err = box.objectBox.RunWithCursor(box.typeId, true, func(cursor *Cursor) error {
		var errInner error
		object, errInner = cursor.Get(id)
		return errInner
	})
	return
}

func (box *Box) GetAll() (slice interface{}, err error) {
	err = box.objectBox.RunWithCursor(box.typeId, true, func(cursor *Cursor) error {
		var errInner error
		slice, errInner = cursor.GetAll()
		return errInner
	})
	return
}
