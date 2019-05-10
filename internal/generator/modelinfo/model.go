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
	"math/rand"
	"os"
	"strings"
)

type Id = uint32
type Uid = uint64

const (
	ModelVersion = 5

	// modelVersion supported by this parser & generator
	minModelVersion = 4
	maxModelVersion = ModelVersion
)

type ModelInfo struct {
	// NOTE don't change order of these json exported properties because it will change users' model.json files
	Note1                string    `json:"_note1"`
	Note2                string    `json:"_note2"`
	Note3                string    `json:"_note3"`
	Entities             []*Entity `json:"entities"`
	LastEntityId         IdUid     `json:"lastEntityId"`
	LastIndexId          IdUid     `json:"lastIndexId"`
	LastRelationId       IdUid     `json:"lastRelationId"`
	ModelVersion         int       `json:"modelVersion"`
	MinimumParserVersion int       `json:"modelVersionParserMinimum"`
	RetiredEntityUids    []Uid     `json:"retiredEntityUids"`
	RetiredIndexUids     []Uid     `json:"retiredIndexUids"`
	RetiredPropertyUids  []Uid     `json:"retiredPropertyUids"`
	RetiredRelationUids  []Uid     `json:"retiredRelationUids"`
	Version              int       `json:"version"` // user specified version

	file *os.File // file handle, locked while the model is open

	// Model Template
	Package string `json:"-"`
}

var defaultModel = ModelInfo{
	Note1:                "KEEP THIS FILE! Check it into a version control system (VCS) like git.",
	Note2:                "ObjectBox manages crucial IDs for your object model. See docs for details.",
	Note3:                "If you have VCS merge conflicts, you must resolve them according to ObjectBox docs.",
	Entities:             make([]*Entity, 0),
	RetiredEntityUids:    make([]Uid, 0),
	RetiredIndexUids:     make([]Uid, 0),
	RetiredPropertyUids:  make([]Uid, 0),
	RetiredRelationUids:  make([]Uid, 0),
	ModelVersion:         maxModelVersion,
	MinimumParserVersion: maxModelVersion,
	Version:              1,
}

func createModelInfo() *ModelInfo {
	var model = defaultModel
	return &model
}

func (model *ModelInfo) fillMissing() {
	// just replace comments with the latest ones
	model.Note1 = defaultModel.Note1
	model.Note2 = defaultModel.Note2
	model.Note3 = defaultModel.Note3
}

// performs initial validation of loaded data so that it doesn't have to be checked in each function
func (model *ModelInfo) Validate() (err error) {
	if model.ModelVersion < minModelVersion {
		return fmt.Errorf("the loaded model is too old - version %d while the minimum supported is %d - "+
			" consider upgrading with an older generator or manually.", model.ModelVersion, minModelVersion)
	}

	if model.ModelVersion > maxModelVersion {
		if model.MinimumParserVersion == 0 || model.MinimumParserVersion > ModelVersion {
			return fmt.Errorf("the loaded model has been created with a newer generator version %d "+
				" while the maximimum supported version is %d. Please upgrade your toolchain/generator",
				model.ModelVersion, maxModelVersion)
		}
	}

	if model.Entities == nil {
		return fmt.Errorf("entities are not defined or not an array")
	}

	for _, entity := range model.Entities {
		if entity.model == nil {
			entity.model = model
		} else if entity.model != model {
			return fmt.Errorf("entity %s %s has incorrect parent model reference", entity.Name, entity.Id)
		}

		err = entity.Validate()
		if err != nil {
			return fmt.Errorf("entity %s %s is invalid: %s", entity.Name, entity.Id, err)
		}
	}

	if len(model.Entities) > 0 {
		if err = model.LastEntityId.Validate(); err != nil {
			return fmt.Errorf("lastEntityId: %s", err)
		}

		var lastId = model.LastEntityId.getIdSafe()
		var lastUid = model.LastEntityId.getUidSafe()

		var found = false
		for _, entity := range model.Entities {
			if lastId == entity.Id.getIdSafe() {
				if lastUid != entity.Id.getUidSafe() {
					return fmt.Errorf("lastEntityId %s doesn't match entity %s %s",
						model.LastEntityId, entity.Name, entity.Id)
				}
				found = true
			} else if lastId < entity.Id.getIdSafe() {
				return fmt.Errorf("lastEntityId %s is lower than entity %s %s",
					model.LastEntityId, entity.Name, entity.Id)
			}
		}

		if !found && !searchSliceUid(model.RetiredEntityUids, lastUid) {
			return fmt.Errorf("lastEntityId %s doesn't match any entity", model.LastEntityId)
		}
	}

	if len(model.LastIndexId) > 0 {
		if err = model.LastIndexId.Validate(); err != nil {
			return fmt.Errorf("lastIndexId: %s", err)
		}
	}

	if len(model.LastRelationId) > 0 || model.hasRelations() {
		if err = model.LastRelationId.Validate(); err != nil {
			return fmt.Errorf("lastRelationId: %s", err)
		}

		// find the last relation ID among entities' relations
		var lastId = model.LastRelationId.getIdSafe()
		var lastUid = model.LastRelationId.getUidSafe()
		var found = false

		for _, entity := range model.Entities {
			for _, relation := range entity.Relations {
				if relation.entity == nil {
					relation.entity = entity
				} else if relation.entity != entity {
					return fmt.Errorf("relation %s %s has incorrect parent entity reference",
						relation.Name, relation.Id)
				}

				if lastId == relation.Id.getIdSafe() {
					if lastUid != relation.Id.getUidSafe() {
						return fmt.Errorf("lastRelationId %s doesn't match relation %s %s",
							model.LastRelationId, relation.Name, relation.Id)
					}
					found = true
				} else if lastId < relation.Id.getIdSafe() {
					return fmt.Errorf("lastRelationId %s is lower than relation %s %s",
						model.LastRelationId, relation.Name, relation.Id)
				}
			}
		}

		if !found && !searchSliceUid(model.RetiredRelationUids, lastUid) {
			return fmt.Errorf("lastRelationId %s doesn't match any relation", model.LastRelationId)
		}
	}

	if model.RetiredEntityUids == nil {
		return fmt.Errorf("retiredEntityUids are not defined or not an array")
	}

	if model.RetiredIndexUids == nil {
		return fmt.Errorf("retiredIndexUids are not defined or not an array")
	}

	if model.RetiredPropertyUids == nil {
		return fmt.Errorf("retiredPropertyUids are not defined or not an array")
	}

	return nil
}

func (model *ModelInfo) hasRelations() bool {
	for _, entity := range model.Entities {
		if len(entity.Relations) > 0 {
			return true
		}
	}
	return false
}

func (model *ModelInfo) FindEntityByUid(uid Uid) (*Entity, error) {
	for _, entity := range model.Entities {
		entityUid, _ := entity.Id.GetUid()
		if entityUid == uid {
			return entity, nil
		}
	}

	return nil, fmt.Errorf("entity with uid %d was not found", uid)
}

func (model *ModelInfo) FindEntityByName(name string) (*Entity, error) {
	for _, entity := range model.Entities {
		if strings.ToLower(entity.Name) == strings.ToLower(name) {
			return entity, nil
		}
	}

	return nil, fmt.Errorf("entity named %s was not found", name)
}

func (model *ModelInfo) CreateEntity(name string) (*Entity, error) {
	var id Id = 1
	if len(model.Entities) > 0 {
		id = model.LastEntityId.getIdSafe() + 1
	}

	uniqueUid, err := model.generateUid()

	if err != nil {
		return nil, err
	}

	var entity = CreateEntity(model, id, uniqueUid)
	entity.Name = name

	model.Entities = append(model.Entities, entity)
	model.LastEntityId = entity.Id

	return entity, nil
}

func (model *ModelInfo) generateUid() (result Uid, err error) {
	result = 0

	for i := 0; i < 1000; i++ {
		t := Uid(rand.Int63())
		if !model.containsUid(t) {
			result = t
			break
		}
	}

	if result == 0 {
		err = fmt.Errorf("internal error = could not generate a unique UID")
	}

	return result, err
}

func (model *ModelInfo) createIndexId() (IdUid, error) {
	var id Id = 1
	if len(model.LastIndexId) > 0 {
		id = model.LastIndexId.getIdSafe() + 1
	}

	uniqueUid, err := model.generateUid()

	if err != nil {
		return "", err
	}

	model.LastIndexId = CreateIdUid(id, uniqueUid)
	return model.LastIndexId, nil
}

func (model *ModelInfo) createRelationId() (IdUid, error) {
	var id Id = 1
	if len(model.LastRelationId) > 0 {
		id = model.LastRelationId.getIdSafe() + 1
	}

	uniqueUid, err := model.generateUid()

	if err != nil {
		return "", err
	}

	model.LastRelationId = CreateIdUid(id, uniqueUid)
	return model.LastRelationId, nil
}

// recursively checks whether given UID is present in the model
func (model *ModelInfo) containsUid(searched Uid) bool {
	if model.LastEntityId.getUidSafe() == searched {
		return true
	}

	if model.LastIndexId.getUidSafe() == searched {
		return true
	}

	if model.LastRelationId.getUidSafe() == searched {
		return true
	}

	if searchSliceUid(model.RetiredEntityUids, searched) {
		return true
	}

	if searchSliceUid(model.RetiredIndexUids, searched) {
		return true
	}

	if searchSliceUid(model.RetiredPropertyUids, searched) {
		return true
	}

	for _, entity := range model.Entities {
		if entity.containsUid(searched) {
			return true
		}
	}

	return false
}

// the passed slices are not too large so let's just do linear search
func searchSliceUid(slice []Uid, searched Uid) bool {
	for _, i := range slice {
		if i == searched {
			return true
		}
	}

	return false
}
