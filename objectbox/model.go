/*
 * Copyright 2018-2021 ObjectBox Ltd. All rights reserved.
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
	"unsafe"

	flatbuffers "github.com/google/flatbuffers/go"
	gogen "github.com/objectbox/objectbox-generator/v4/cmd/objectbox-gogen"
)

// ObjectBinding provides an interface for various object types to be included in the model
type ObjectBinding interface {
	// AddToModel adds the entity information, including properties, indexes, etc., to the model during construction.
	AddToModel(model *Model)

	// GetId reads the ID field of the given object.
	GetId(object interface{}) (id uint64, err error)

	// SetId sets the ID field on the given object.
	SetId(object interface{}, id uint64) error

	// PutRelated updates/inserts objects related to the given object, based on the available object data.
	PutRelated(ob *ObjectBox, object interface{}, id uint64) error

	// Flatten serializes the object to FlatBuffers. The given ID must be used instead of the object field.
	Flatten(object interface{}, fbb *flatbuffers.Builder, id uint64) error

	// Load constructs the object from serialized byte buffer. Also reads data for eagerly loaded related entities.
	Load(ob *ObjectBox, bytes []byte) (interface{}, error)

	// MakeSlice creates a slice of objects with the given capacity (0 length).
	MakeSlice(capacity int) interface{}

	// AppendToSlice adds the object at the end of the slice created by MakeSlice(). Returns the new slice.
	AppendToSlice(slice interface{}, object interface{}) (sliceNew interface{})

	// GeneratorVersion returns the version used to generate this binding - used to verify the compatibility.
	GeneratorVersion() int
}

// Model is used by the generated code to represent information about the ObjectBox database schema
type Model struct {
	cModel *C.OBX_model
	Error  error

	currentEntity  *entity
	entitiesById   map[TypeId]*entity
	entitiesByName map[string]*entity

	lastEntityId  TypeId
	lastEntityUid uint64

	lastIndexId  TypeId
	lastIndexUid uint64

	lastRelationId  TypeId
	lastRelationUid uint64

	generatorVersion int
}

// NewModel creates a model
func NewModel() *Model {
	var model = &Model{
		entitiesById:   make(map[TypeId]*entity),
		entitiesByName: make(map[string]*entity),
	}

	model.Error = cCallBool(func() bool {
		model.cModel = C.obx_model()
		return model.cModel != nil
	})

	return model
}

// GeneratorVersion configures version of the generator used to create this model
func (model *Model) GeneratorVersion(version int) {
	if model.Error != nil {
		return
	}

	model.generatorVersion = version
}

// LastEntityId declares an entity with the highest ID.
// Used as a compatibility check when opening DB with an older model version.
func (model *Model) LastEntityId(id TypeId, uid uint64) {
	if model.Error != nil {
		return
	}
	model.lastEntityId = id
	model.lastEntityUid = uid
	C.obx_model_last_entity_id(model.cModel, C.obx_schema_id(id), C.obx_uid(uid))
}

// LastIndexId declares an index with the highest ID.
// Used as a compatibility check when opening DB with an older model version.
func (model *Model) LastIndexId(id TypeId, uid uint64) {
	if model.Error != nil {
		return
	}
	model.lastIndexId = id
	model.lastIndexUid = uid
	C.obx_model_last_index_id(model.cModel, C.obx_schema_id(id), C.obx_uid(uid))
}

// LastRelationId declares a relation with the highest ID.
// Used as a compatibility check when opening DB with an older model version.
func (model *Model) LastRelationId(id TypeId, uid uint64) {
	if model.Error != nil {
		return
	}
	model.lastRelationId = id
	model.lastRelationUid = uid
	C.obx_model_last_relation_id(model.cModel, C.obx_schema_id(id), C.obx_uid(uid))
}

// Entity creates an entity in a model
func (model *Model) Entity(name string, id TypeId, uid uint64) {
	if model.Error != nil {
		return
	}
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	if err := cCall(func() C.obx_err {
		return C.obx_model_entity(model.cModel, cname, C.obx_schema_id(id), C.obx_uid(uid))
	}); err != nil {
		model.Error = err
		return
	}

	model.currentEntity = &entity{
		name: name,
		id:   id,
	}
}

// EntityFlags configures behavior of entities
func (model *Model) EntityFlags(entityFlags int) {
	if model.Error != nil {
		return
	}
	model.Error = cCall(func() C.obx_err {
		return C.obx_model_entity_flags(model.cModel, C.uint32_t(entityFlags))
	})
}

// TODO each Entity-related method (e.g. Property, Relation,...) should check whether currentEntity is not nil

// Relation adds a "standalone" many-to-many relation between the current entity and a target entity
func (model *Model) Relation(relationId TypeId, relationUid uint64, targetEntityId TypeId, targetEntityUid uint64) {
	if model.Error != nil {
		return
	}

	model.Error = cCall(func() C.obx_err {
		return C.obx_model_relation(model.cModel, C.obx_schema_id(relationId), C.obx_uid(relationUid),
			C.obx_schema_id(targetEntityId), C.obx_uid(targetEntityUid))
	})

	model.currentEntity.hasRelations = true
}

// EntityLastPropertyId declares a property with the highest ID.
// Used as a compatibility check when opening DB with an older model version.
func (model *Model) EntityLastPropertyId(id TypeId, uid uint64) {
	if model.Error != nil {
		return
	}

	model.Error = cCall(func() C.obx_err {
		return C.obx_model_entity_last_property_id(model.cModel, C.obx_schema_id(id), C.obx_uid(uid))
	})
}

// Property creates a property in an Entity
func (model *Model) Property(name string, propertyType int, id TypeId, uid uint64) {
	if model.Error != nil {
		return
	}
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	model.Error = cCall(func() C.obx_err {
		return C.obx_model_property(model.cModel, cname, C.OBXPropertyType(propertyType), C.obx_schema_id(id), C.obx_uid(uid))
	})
}

// PropertyFlags configures type and other information about the property
func (model *Model) PropertyFlags(propertyFlags int) {
	if model.Error != nil {
		return
	}
	model.Error = cCall(func() C.obx_err {
		return C.obx_model_property_flags(model.cModel, C.uint32_t(propertyFlags))
	})
}

// PropertyIndex creates a new index on the property
func (model *Model) PropertyIndex(id TypeId, uid uint64) {
	if model.Error != nil {
		return
	}
	model.Error = cCall(func() C.obx_err {
		return C.obx_model_property_index_id(model.cModel, C.obx_schema_id(id), C.obx_uid(uid))
	})
}

// PropertyRelation adds a property-based (i.e. to-one) relation
func (model *Model) PropertyRelation(targetEntityName string, indexId TypeId, indexUid uint64) {
	if model.Error != nil {
		return
	}

	cname := C.CString(targetEntityName)
	defer C.free(unsafe.Pointer(cname))

	model.Error = cCall(func() C.obx_err {
		return C.obx_model_property_relation(model.cModel, cname, C.obx_schema_id(indexId), C.obx_uid(indexUid))
	})

	model.currentEntity.hasRelations = true
}

// RegisterBinding attaches generated binding code to the model.
// The binding is used by ObjectBox for marshalling and other typed operations.
func (model *Model) RegisterBinding(binding ObjectBinding) {
	if model.Error != nil {
		return
	}

	model.currentEntity = nil

	binding.AddToModel(model)

	if model.currentEntity == nil {
		model.Error = fmt.Errorf("invalid binding - model.Entity() not called")
		return
	}

	id := model.currentEntity.id
	name := model.currentEntity.name

	if id == 0 {
		model.Error = fmt.Errorf("invalid binding - entity id is not set")
		return
	}

	if name == "" {
		model.Error = fmt.Errorf("invalid binding - entity name is not set")
		return
	}

	if model.entitiesById[id] != nil {
		model.Error = fmt.Errorf("duplicate binding - entity id %d is already registered", id)
		return
	}

	if model.entitiesByName[name] != nil {
		model.Error = fmt.Errorf("duplicate binding - entity name %s is already registered", name)
		return
	}

	var version = binding.GeneratorVersion()
	if version != gogen.VersionId {
		model.Error = fmt.Errorf("incompatible generator version %d used to generate the binding %s code "+
			"- please follow the upgrade procedure described in the README.md", version, name)
		return
	}

	model.currentEntity.binding = binding
	model.entitiesById[id] = model.currentEntity
	model.entitiesByName[name] = model.currentEntity

	model.currentEntity = nil
}

func (model *Model) validate() error {
	if model.Error != nil {
		return model.Error
	}

	if model.generatorVersion != gogen.VersionId {
		return fmt.Errorf("incompatible generator version %d used to generate the model code "+
			"- please follow the upgrade procedure described in the README.md", model.generatorVersion)
	}

	if model.lastEntityId == 0 || model.lastEntityUid == 0 {
		return fmt.Errorf("last entity ID/UID is missing")
	}

	return nil
}
