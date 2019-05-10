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
	"errors"
	"fmt"
)

type StandaloneRelation struct {
	Id       IdUid   `json:"id"`
	Name     string  `json:"name"`
	Target   *Entity `json:"-"`
	TargetId IdUid   `json:"targetId"`

	entity *Entity
}

func CreateStandaloneRelation(entity *Entity, id IdUid) *StandaloneRelation {
	return &StandaloneRelation{entity: entity, Id: id}
}

// performs initial validation of loaded data so that it doesn't have to be checked in each function
func (relation *StandaloneRelation) Validate() error {
	if err := relation.Id.Validate(); err != nil {
		return err
	}

	if len(relation.Name) == 0 {
		return errors.New("name is undefined")
	}

	if len(relation.TargetId) > 0 {
		if err := relation.TargetId.Validate(); err != nil {
			return err
		}

		for _, entity := range relation.entity.model.Entities {
			if entity.Id == relation.TargetId {
				relation.Target = entity
			}
		}

		if relation.Target == nil {
			return fmt.Errorf("target entity ID %s not found", string(relation.TargetId))
		}
	}

	return nil
}

func (relation *StandaloneRelation) SetTarget(entity *Entity) {
	relation.Target = entity
	relation.TargetId = entity.Id
}
