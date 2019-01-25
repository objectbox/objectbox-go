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

package modelinfo

import "fmt"

type Relation struct {
	Id   IdUid  `json:"id"`
	Name string `json:"name"`
}

// performs initial validation of loaded data so that it doesn't have to be checked in each function
func (relation *Relation) Validate() error {
	if err := relation.Id.Validate(); err != nil {
		return err
	}

	if len(relation.Name) == 0 {
		return fmt.Errorf("name is undefined")
	}

	return nil
}

type StandaloneRelation struct {
	Relation
	TargetEntityId IdUid `json:"-"` // currently not in the json file
}

func CreateStandaloneRelation(id IdUid) *StandaloneRelation {
	return &StandaloneRelation{
		Relation: Relation{Id: id},
	}
}

// performs initial validation of loaded data so that it doesn't have to be checked in each function
func (relation *StandaloneRelation) Validate() error {
	if err := relation.Relation.Validate(); err != nil {
		return err
	}

	if err := relation.TargetEntityId.Validate(); err != nil {
		return err
	}

	return nil
}
