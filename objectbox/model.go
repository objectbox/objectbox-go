package objectbox

/*
#cgo LDFLAGS: -lobjectbox
#include <stdlib.h>
#include "objectbox.h"
*/
import "C"

import (
	"fmt"
	"strconv"
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

// Usually only used by the automatically generated artifacts of ObjectBox generator
type Model struct {
	model              *C.OBX_model
	previousEntityName string
	previousEntityId   TypeId
	Error              error

	bindingsById   map[TypeId]ObjectBinding
	bindingsByName map[string]ObjectBinding

	lastEntityId  TypeId
	lastEntityUid uint64

	lastIndexId  TypeId
	lastIndexUid uint64

	lastRelationId  TypeId
	lastRelationUid uint64
}

func NewModel() (*Model, error) {
	cModel := C.obx_model_create()
	if cModel == nil {
		return nil, createError()
	}
	return &Model{
		model:          cModel,
		bindingsById:   make(map[TypeId]ObjectBinding),
		bindingsByName: make(map[string]ObjectBinding),
	}, nil
}

func (model *Model) LastEntityId(id TypeId, uid uint64) {
	if model.Error != nil {
		return
	}
	model.lastEntityId = id
	model.lastEntityUid = uid
	C.obx_model_last_entity_id(model.model, C.obx_schema_id(id), C.obx_uid(uid))
}

func (model *Model) LastIndexId(id TypeId, uid uint64) {
	if model.Error != nil {
		return
	}
	model.lastIndexId = id
	model.lastIndexUid = uid
	C.obx_model_last_index_id(model.model, C.obx_schema_id(id), C.obx_uid(uid))
}

func (model *Model) LastRelationId(id TypeId, uid uint64) {
	if model.Error != nil {
		return
	}
	model.lastRelationId = id
	model.lastRelationUid = uid
	C.obx_model_last_relation_id(model.model, C.obx_schema_id(id), C.obx_uid(uid))
}

func (model *Model) Entity(name string, id TypeId, uid uint64) {
	if model.Error != nil {
		return
	}
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	rc := C.obx_model_entity(model.model, cname, C.obx_schema_id(id), C.obx_uid(uid))
	if rc != 0 {
		model.Error = createError()
		return
	}
	model.previousEntityName = name
	model.previousEntityId = id
	return
}

func (model *Model) EntityLastPropertyId(id TypeId, uid uint64) {
	if model.Error != nil {
		return
	}
	rc := C.obx_model_entity_last_property_id(model.model, C.obx_schema_id(id), C.obx_uid(uid))
	if rc != 0 {
		model.Error = createError()
	}
	return
}

func (model *Model) Property(name string, propertyType int, id TypeId, uid uint64) {
	if model.Error != nil {
		return
	}
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	rc := C.obx_model_property(model.model, cname, C.OBPropertyType(propertyType), C.obx_schema_id(id), C.obx_uid(uid))
	if rc != 0 {
		model.Error = createError()
	}
	return
}

func (model *Model) PropertyFlags(propertyFlags int) {
	if model.Error != nil {
		return
	}
	rc := C.obx_model_property_flags(model.model, C.OBPropertyFlags(propertyFlags))
	if rc != 0 {
		model.Error = createError()
	}
	return
}

func (model *Model) PropertyIndex(id TypeId, uid uint64) {
	if model.Error != nil {
		return
	}
	rc := C.obx_model_property_index_id(model.model, C.obx_schema_id(id), C.obx_uid(uid))
	if rc != 0 {
		model.Error = createError()
	}
	return
}

func (model *Model) PropertyRelation(targetEntityName string, indexId TypeId, indexUid uint64) {
	if model.Error != nil {
		return
	}
	cname := C.CString(targetEntityName)
	defer C.free(unsafe.Pointer(cname))
	rc := C.obx_model_property_relation(model.model, cname, C.obx_schema_id(indexId), C.obx_uid(indexUid))
	if rc != 0 {
		model.Error = createError()
	}
	return
}

func (model *Model) RegisterBinding(binding ObjectBinding) {
	binding.AddToModel(model)
	id := model.previousEntityId
	name := model.previousEntityName
	if id == 0 {
		panic("No type ID; did you forget to add an entity to the model?")
	}
	if name == "" {
		panic("No type name")
	}
	existingBinding := model.bindingsById[id]
	if existingBinding != nil {
		panic("Already registered a binding for ID " + strconv.Itoa(int(id)))
	}
	existingBinding = model.bindingsByName[name]
	if existingBinding != nil {
		panic("Already registered a binding for name " + name)
	}
	model.bindingsById[id] = binding
	model.bindingsByName[name] = binding
}

func (model *Model) validate() error {
	if model.Error != nil {
		return model.Error
	}

	if model.lastEntityId == 0 || model.lastEntityUid == 0 {
		return fmt.Errorf("last entity ID/UID is missing")
	}

	return nil
}
