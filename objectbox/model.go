package objectbox

/*
#cgo LDFLAGS: -L ${SRCDIR}/libs -lobjectboxc
#include <stdlib.h>
#include <string.h>
#include "objectbox.h"
*/
import "C"

import (
	"unsafe"
)

//noinspection GoUnusedConst
const (
	PropertyType_Bool       = 1
	PropertyType_Byte       = 2
	PropertyType_Short      = 3
	PropertyType_Char       = 4
	PropertyType_Int        = 5
	PropertyType_Long       = 6
	PropertyType_Float      = 7
	PropertyType_Double     = 8
	PropertyType_String     = 9
	PropertyType_Date       = 10
	PropertyType_Relation   = 11
	PropertyType_ByteVector = 23
)

//noinspection GoUnusedConst
const (
	/// One long property on an entity must be the ID
	PropertyFlags_ID = 1

	/// On languages like Java, a non-primitive type is used (aka wrapper types, allowing null)
	PropertyFlags_NON_PRIMITIVE_TYPE = 2

	/// Unused yet
	PropertyFlags_NOT_NULL = 4
	PropertyFlags_INDEXED  = 8
	PropertyFlags_RESERVED = 16
	/// Unused yet: Unique index
	PropertyFlags_UNIQUE = 32
	/// Unused yet: Use a persisted sequence to enforce ID to rise monotonic (no ID reuse)
	PropertyFlags_ID_MONOTONIC_SEQUENCE = 64
	/// Allow IDs to be assigned by the developer
	PropertyFlags_ID_SELF_ASSIGNABLE = 128
	/// Unused yet
	PropertyFlags_INDEX_PARTIAL_SKIP_NULL = 256
	/// Unused yet, used by References for 1) back-references and 2) to clear references to deleted objects (required for ID reuse)
	PropertyFlags_INDEX_PARTIAL_SKIP_ZERO = 512
	/// Virtual properties may not have a dedicated field in their entity class, e.g. target IDs of to-one relations
	PropertyFlags_VIRTUAL = 1024
	/// Index uses a 32 bit hash instead of the value
	/// (32 bits is shorter on disk, runs well on 32 bit systems, and should be OK even with a few collisions)

	PropertyFlags_INDEX_HASH = 2048
	/// Index uses a 64 bit hash instead of the value
	/// (recommended mostly for 64 bit machines with values longer >200 bytes; small values are faster with a 32 bit hash)
	PropertyFlags_INDEX_HASH64 = 4096
)

type Model struct {
	model          *C.OB_model
	lastEntityName string
	lastEntityId   TypeId
	Err            error
}

func NewModel() (model *Model, err error) {
	model = &Model{}
	model.model = C.ob_model_create()
	if model.model == nil {
		model = nil
		err = createError()
	}
	return
}

func (model *Model) LastEntityId(id TypeId, uid uint64) {
	if model.Err != nil {
		return
	}
	C.ob_model_last_entity_id(model.model, C.uint(id), C.ulong(uid))
}

func (model *Model) Entity(name string, id TypeId, uid uint64) (err error) {
	if model.Err != nil {
		return model.Err
	}
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	rc := C.ob_model_entity(model.model, cname, C.uint(id), C.ulong(uid))
	if rc != 0 {
		err = createError()
		model.Err = err
		return
	}
	model.lastEntityName = name
	model.lastEntityId = id
	return
}

func (model *Model) EntityLastPropertyId(id TypeId, uid uint64) (err error) {
	if model.Err != nil {
		return model.Err
	}
	rc := C.ob_model_entity_last_property_id(model.model, C.uint(id), C.ulong(uid))
	if rc != 0 {
		err = createError()
		model.Err = err
	}
	return
}

func (model *Model) Property(name string, propertyType int, id TypeId, uid uint64) (err error) {
	if model.Err != nil {
		return model.Err
	}
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	rc := C.ob_model_property(model.model, cname, C.OBPropertyType(propertyType), C.uint(id), C.ulong(uid))
	if rc != 0 {
		err = createError()
		model.Err = err
	}
	return
}

func (model *Model) PropertyFlags(propertyFlags int) (err error) {
	if model.Err != nil {
		return model.Err
	}
	rc := C.ob_model_property_flags(model.model, C.OBPropertyFlags(propertyFlags))
	if rc != 0 {
		err = createError()
		model.Err = err
	}
	return
}
