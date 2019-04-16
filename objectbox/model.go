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
#cgo LDFLAGS: -lobjectbox
#include <stdlib.h>
#include "objectbox.h"
*/
import "C"

import (
	"fmt"
	"unsafe"

	"github.com/google/flatbuffers/go"
	"github.com/objectbox/objectbox-go/internal/generator"
)

// An ObjectBinding provides an interface for various object types to be included in the model
type ObjectBinding interface {
	AddToModel(model *Model)
	GetId(object interface{}) (id uint64, err error)
	SetId(object interface{}, id uint64)
	PutRelated(txn *Transaction, object interface{}, id uint64) error
	Flatten(object interface{}, fbb *flatbuffers.Builder, id uint64) error
	Load(txn *Transaction, bytes []byte) (interface{}, error)
	MakeSlice(capacity int) interface{}
	AppendToSlice(slice interface{}, object interface{}) (sliceNew interface{})
	GeneratorVersion() int
}

// Model is used by the generated code to represent information about the ObjectBox database schema
type Model struct {
	model *C.OBX_model
	Error error

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

func NewModel() *Model {
	cModel := C.obx_model_create()
	var err error
	if cModel == nil {
		err = createError()
	}
	return &Model{
		model:          cModel,
		Error:          err,
		entitiesById:   make(map[TypeId]*entity),
		entitiesByName: make(map[string]*entity),
	}
}

func (model *Model) GeneratorVersion(version int) {
	if model.Error != nil {
		return
	}

	model.generatorVersion = version
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

	model.currentEntity = &entity{
		name: name,
		id:   id,
	}

	return
}

// TODO each Entity-related method (e.g. Property, Relation,...) should check whether currentEntity is not nil

// Relation adds a "standalone" many-to-many relation between the current entity and a target entity
func (model *Model) Relation(relationId TypeId, relationUid uint64, targetEntityId TypeId, targetEntityUid uint64) {
	if model.Error != nil {
		return
	}

	rc := C.obx_model_relation(model.model, C.obx_schema_id(relationId), C.obx_uid(relationUid),
		C.obx_schema_id(targetEntityId), C.obx_uid(targetEntityUid))
	if rc != 0 {
		model.Error = createError()
		return
	}

	model.currentEntity.hasRelations = true
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

	rc := C.obx_model_property(model.model, cname, C.OBXPropertyType(propertyType), C.obx_schema_id(id), C.obx_uid(uid))
	if rc != 0 {
		model.Error = createError()
	}
	return
}

func (model *Model) PropertyFlags(propertyFlags int) {
	if model.Error != nil {
		return
	}
	rc := C.obx_model_property_flags(model.model, C.OBXPropertyFlags(propertyFlags))
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

	model.currentEntity.hasRelations = true
	return
}

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
	if version != generator.Version {
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

	if model.generatorVersion != generator.Version {
		return fmt.Errorf("incompatible generator version %d used to generate the model code "+
			"- please follow the upgrade procedure described in the README.md", model.generatorVersion)
	}

	if model.lastEntityId == 0 || model.lastEntityUid == 0 {
		return fmt.Errorf("last entity ID/UID is missing")
	}

	return nil
}
