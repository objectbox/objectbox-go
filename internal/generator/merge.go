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

package generator

import (
	"fmt"

	"github.com/objectbox/objectbox-go/internal/generator/modelinfo"
)

func mergeBindingWithModelInfo(binding *Binding, modelInfo *modelinfo.ModelInfo) error {
	for _, bindingEntity := range binding.Entities {
		if modelEntity, err := getModelEntity(bindingEntity, modelInfo); err != nil {
			return err
		} else if err := mergeModelEntity(bindingEntity, modelEntity); err != nil {
			return err
		}
	}

	// NOTE this is not ideal as there could be models across multiple packages
	modelInfo.Package = binding.Package

	return nil
}

func getModelEntity(bindingEntity *Entity, modelInfo *modelinfo.ModelInfo) (*modelinfo.Entity, error) {
	if bindingEntity.Uid != 0 {
		return modelInfo.FindEntityByUid(bindingEntity.Uid)
	} else if entity, err := modelInfo.FindEntityByName(bindingEntity.Name); entity != nil {
		return entity, err
	} else {
		return modelInfo.CreateEntity()
	}
}

func mergeModelEntity(bindingEntity *Entity, modelEntity *modelinfo.Entity) (err error) {
	modelEntity.Name = bindingEntity.Name

	if bindingEntity.Id, bindingEntity.Uid, err = modelEntity.Id.Get(); err != nil {
		return err
	}

	// add all properties from the bindings to the model and update/rename the changed ones
	for _, bindingProperty := range bindingEntity.Properties {
		if modelProperty, err := getModelProperty(bindingProperty, modelEntity); err != nil {
			return err
		} else if err := mergeModelProperty(bindingProperty, modelProperty); err != nil {
			return err
		}
	}

	// remove the missing (removed) properties
	removedProperties := make([]*modelinfo.Property, 0)
	for _, modelProperty := range modelEntity.Properties {
		if !bindingPropertyExists(modelProperty, bindingEntity) {
			removedProperties = append(removedProperties, modelProperty)
		}
	}

	for _, property := range removedProperties {
		if err := modelEntity.RemoveProperty(property); err != nil {
			return err
		}
	}

	bindingEntity.LastPropertyId = modelEntity.LastPropertyId

	return nil
}

func getModelProperty(bindingProperty *Property, modelEntity *modelinfo.Entity) (*modelinfo.Property, error) {
	if bindingProperty.Uid != 0 {
		return modelEntity.FindPropertyByUid(bindingProperty.Uid)
	}

	property, err := modelEntity.FindPropertyByName(bindingProperty.Name)

	// handle uid request
	if bindingProperty.uidRequest {
		var errInfo string
		if property != nil {
			if uid, err := property.Id.GetUid(); err != nil {
				return nil, err
			} else {
				errInfo = fmt.Sprintf("model property UID = %d", uid)
			}
		} else {
			errInfo = "property not found in the model"
		}
		return nil, fmt.Errorf("uid annotation value must not be empty (%s) on property %s, entity %s",
			errInfo, bindingProperty.Name, bindingProperty.entity.Name)
	}

	if property != nil {
		return property, err
	} else {
		return modelEntity.CreateProperty()
	}
}

func mergeModelProperty(bindingProperty *Property, modelProperty *modelinfo.Property) (err error) {
	modelProperty.Name = bindingProperty.Name

	if bindingProperty.Id, bindingProperty.Uid, err = modelProperty.Id.Get(); err != nil {
		return err
	}

	if bindingProperty.Index == nil {
		// if there shouldn't be an index
		if modelProperty.IndexId != nil {
			// if there originally was an index, remove it
			if err = modelProperty.RemoveIndex(); err != nil {
				return err
			}
		}
	} else {
		// if there should be an index, create it (or reuse an existing one)
		if modelProperty.IndexId == nil {
			if err = modelProperty.CreateIndex(); err != nil {
				return err
			}
		}

		if bindingProperty.Index.Id, bindingProperty.Index.Uid, err = modelProperty.IndexId.Get(); err != nil {
			return err
		}
	}

	return nil
}

func bindingPropertyExists(modelProperty *modelinfo.Property, bindingEntity *Entity) bool {
	for _, bindingProperty := range bindingEntity.Properties {
		if bindingProperty.Name == modelProperty.Name {
			return true
		}
	}

	return false
}
