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

// CheckRelationCycles finds relations cycles
func (model *ModelInfo) CheckRelationCycles() error {
	// DFS cycle check, storing relation path in the recursion stack
	var recursionStack = make(map[*Entity]bool)
	for _, entity := range model.Entities {
		if err := entity.checkRelationCycles(&recursionStack, entity.Name); err != nil {
			return err
		}
	}

	return nil
}

func (entity *Entity) checkRelationCycles(recursionStack *map[*Entity]bool, path string) error {
	(*recursionStack)[entity] = true

	// to-many relations
	for _, rel := range entity.Relations {
		if err := checkRelationCycle(recursionStack, path+"."+rel.Name, rel.Target); err != nil {
			return err
		}
	}

	// to-one relations
	for _, prop := range entity.Properties {
		if prop.RelationTarget == "" {
			continue
		}

		relTarget, _ := entity.model.FindEntityByName(prop.RelationTarget)

		if err := checkRelationCycle(recursionStack, path+"."+prop.Name, relTarget); err != nil {
			return err
		}
	}

	delete(*recursionStack, entity)
	return nil
}

func checkRelationCycle(recursionStack *map[*Entity]bool, path string, relTarget *Entity) error {
	// this happens if the entity containing this relation haven't been defined in this file
	if relTarget == nil {
		return nil
	}

	if (*recursionStack)[relTarget] {
		return fmt.Errorf("relation cycle detected: %s (%s)", path, relTarget.Name)
	}

	return relTarget.checkRelationCycles(recursionStack, path)
}
