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

import "fmt"

type Property struct {
	Id      IdUid  `json:"id"`
	Name    string `json:"name"`
	IndexId *IdUid `json:"indexId,omitempty"`

	entity *Entity
}

func CreateProperty(entity *Entity, id Id, uid Uid) *Property {
	return &Property{
		entity: entity,
		Id:     CreateIdUid(id, uid),
	}
}

// performs initial validation of loaded data so that it doesn't have to be checked in each function
func (property *Property) Validate() error {
	if property.entity == nil {
		return fmt.Errorf("undefined parent entity")
	}

	if err := property.Id.Validate(); err != nil {
		return err
	}

	if property.IndexId != nil {
		if err := property.IndexId.Validate(); err != nil {
			return fmt.Errorf("indexId: %s", err)
		}
	}

	if len(property.Name) == 0 {
		return fmt.Errorf("name is undefined")
	}

	return nil
}

func (property *Property) CreateIndex() error {
	if property.IndexId != nil {
		return fmt.Errorf("can't create an index - it already exists")
	}

	if indexId, err := property.entity.model.createIndexId(); err != nil {
		return err
	} else {
		property.IndexId = &indexId
		return nil
	}
}

func (property *Property) RemoveIndex() error {
	if property.IndexId == nil {
		return fmt.Errorf("can't remove index - it's not defined")
	}

	property.entity.model.RetiredIndexUids = append(property.entity.model.RetiredIndexUids, property.IndexId.getUidSafe())

	property.IndexId = nil

	return nil
}

// recursively checks whether given UID is present in the model
func (property *Property) containsUid(searched Uid) bool {
	if property.Id.getUidSafe() == searched {
		return true
	}

	if property.IndexId != nil && property.IndexId.getUidSafe() == searched {
		return true
	}

	return false
}
