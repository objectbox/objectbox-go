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

package modelinfo

import (
	"fmt"
	"strings"
)

// Entity represents a DB entity
type Entity struct {
	Id             IdUid                 `json:"id"`
	LastPropertyId IdUid                 `json:"lastPropertyId"`
	Name           string                `json:"name"`
	Properties     []*Property           `json:"properties"`
	Relations      []*StandaloneRelation `json:"relations,omitempty"`

	model *ModelInfo
}

// CreateEntity constructs an Entity
func CreateEntity(model *ModelInfo, id Id, uid Uid) *Entity {
	return &Entity{
		model:      model,
		Id:         CreateIdUid(id, uid),
		Properties: make([]*Property, 0),
	}
}

// Validate performs initial validation of loaded data so that it doesn't have to be checked in each function
func (entity *Entity) Validate() (err error) {
	if entity.model == nil {
		return fmt.Errorf("undefined parent model")
	}

	if err = entity.Id.Validate(); err != nil {
		return err
	}

	if len(entity.Name) == 0 {
		return fmt.Errorf("name is undefined")
	}

	if len(entity.Properties) > 0 {
		if err = entity.LastPropertyId.Validate(); err != nil {
			return fmt.Errorf("lastPropertyId: %s", err)
		}

		var lastID = entity.LastPropertyId.getIdSafe()
		var lastUID = entity.LastPropertyId.getUidSafe()

		var found = false
		for _, property := range entity.Properties {
			if property.entity == nil {
				property.entity = entity
			} else if property.entity != entity {
				return fmt.Errorf("relation %s %s has incorrect parent entity reference",
					property.Name, property.Id)
			}

			if lastID == property.Id.getIdSafe() {
				if lastUID != property.Id.getUidSafe() {
					return fmt.Errorf("lastPropertyId %s doesn't match relation %s %s",
						entity.LastPropertyId, property.Name, property.Id)
				}
				found = true
			} else if lastID < property.Id.getIdSafe() {
				return fmt.Errorf("lastPropertyId %s is lower than relation %s %s",
					entity.LastPropertyId, property.Name, property.Id)
			}
		}

		if !found && !searchSliceUID(entity.model.RetiredPropertyUids, lastUID) {
			return fmt.Errorf("lastPropertyId %s doesn't match any relation", entity.LastPropertyId)
		}
	}

	if entity.Properties == nil {
		return fmt.Errorf("properties are not defined or not an array")
	}

	for _, property := range entity.Properties {
		err = property.Validate()
		if err != nil {
			return fmt.Errorf("property %s %s is invalid: %s", property.Name, string(property.Id), err)
		}
	}

	for _, relation := range entity.Relations {
		if relation.entity == nil {
			relation.entity = entity
		} else if relation.entity != entity {
			return fmt.Errorf("relation %s %s has incorrect parent model reference", relation.Name, relation.Id)
		}

		err = relation.Validate()
		if err != nil {
			return fmt.Errorf("relation %s %s is invalid: %s", relation.Name, string(relation.Id), err)
		}
	}

	return nil
}

// FindPropertyByUID finds a property by UID
func (entity *Entity) FindPropertyByUID(uid Uid) (*Property, error) {
	for _, property := range entity.Properties {
		propertyUID, _ := property.Id.GetUid()
		if propertyUID == uid {
			return property, nil
		}
	}

	return nil, fmt.Errorf("property with Uid %d not found", uid)
}

//FindPropertyByName finds a property by name
func (entity *Entity) FindPropertyByName(name string) (*Property, error) {
	for _, property := range entity.Properties {
		if strings.ToLower(property.Name) == strings.ToLower(name) {
			return property, nil
		}
	}

	return nil, fmt.Errorf("property with Name %s not found", name)
}

// CreateProperty creates a property
func (entity *Entity) CreateProperty() (*Property, error) {
	var id Id = 1
	if len(entity.Properties) > 0 {
		id = entity.LastPropertyId.getIdSafe() + 1
	}

	uniqueUID, err := entity.model.generateUID()

	if err != nil {
		return nil, err
	}

	var property = CreateProperty(entity, id, uniqueUID)

	entity.Properties = append(entity.Properties, property)
	entity.LastPropertyId = property.Id

	return property, nil
}

// RemoveProperty removes a property
func (entity *Entity) RemoveProperty(property *Property) error {
	var indexToRemove = -1
	for index, prop := range entity.Properties {
		if prop == property {
			indexToRemove = index
		}
	}

	if indexToRemove < 0 {
		return fmt.Errorf("can't remove property %s %s - not found", property.Name, property.Id)
	}

	// remove index from the property
	if property.IndexId != nil {
		if err := property.RemoveIndex(); err != nil {
			return err
		}
	}

	// remove from list
	entity.Properties = append(entity.Properties[:indexToRemove], entity.Properties[indexToRemove+1:]...)

	// store the UID in the "retired" list so that it's not reused in the future
	entity.model.RetiredPropertyUids = append(entity.model.RetiredPropertyUids, property.Id.getUidSafe())

	return nil
}

// FindRelationByUID Finds relation by UID
func (entity *Entity) FindRelationByUID(uid Uid) (*StandaloneRelation, error) {
	for _, relation := range entity.Relations {
		relationUID, _ := relation.Id.GetUid()
		if relationUID == uid {
			return relation, nil
		}
	}

	return nil, fmt.Errorf("relation with Uid %d not found", uid)
}

// FindRelationByName finds relation by name
func (entity *Entity) FindRelationByName(name string) (*StandaloneRelation, error) {
	for _, relation := range entity.Relations {
		if strings.ToLower(relation.Name) == strings.ToLower(name) {
			return relation, nil
		}
	}

	return nil, fmt.Errorf("relation with Name %s not found", name)
}

// CreateRelation creates relation
func (entity *Entity) CreateRelation() (*StandaloneRelation, error) {
	id, err := entity.model.createRelationID()
	if err != nil {
		return nil, err
	}

	var relation = CreateStandaloneRelation(entity, id)
	entity.Relations = append(entity.Relations, relation)
	return relation, nil
}

// RemoveRelation removes relation
func (entity *Entity) RemoveRelation(relation *StandaloneRelation) error {
	var indexToRemove = -1
	for index, rel := range entity.Relations {
		if rel == relation {
			indexToRemove = index
		}
	}

	if indexToRemove < 0 {
		return fmt.Errorf("can't remove relation %s %s - not found", relation.Name, relation.Id)
	}

	// remove from list
	entity.Relations = append(entity.Relations[:indexToRemove], entity.Relations[indexToRemove+1:]...)

	// store the UID in the "retired" list so that it's not reused in the future
	entity.model.RetiredRelationUids = append(entity.model.RetiredRelationUids, relation.Id.getUidSafe())

	return nil
}

// containsUID recursively checks whether given UID is present in the model
func (entity *Entity) containsUID(searched Uid) bool {
	if entity.Id.getUidSafe() == searched {
		return true
	}

	if entity.LastPropertyId.getUidSafe() == searched {
		return true
	}

	for _, property := range entity.Properties {
		if property.containsUID(searched) {
			return true
		}
	}

	for _, relation := range entity.Relations {
		if relation.Id.getUidSafe() == searched {
			return true
		}
	}

	return false
}
