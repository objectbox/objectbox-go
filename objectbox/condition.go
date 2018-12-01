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

package objectbox

type Condition interface {
	build(qb *queryBuilder) (conditionId, error)
}

type conditionId = int

type conditionClosure struct {
	buildFun func(qb *queryBuilder) (conditionId, error)
}

func (condition *conditionClosure) build(qb *queryBuilder) (conditionId, error) {
	return condition.buildFun(qb)
}

// Combines multiple conditions with an operator
type conditionCombination struct {
	or         bool // AND by default
	conditions []Condition
}

func (condition *conditionCombination) build(qb *queryBuilder) (conditionId, error) {
	ids := make([]conditionId, len(condition.conditions))

	for _, sub := range condition.conditions {
		if id, err := sub.build(qb); err != nil {
			return 0, err
		} else {
			ids = append(ids, id)
		}
	}

	// TODO
	//if condition.or {
	//	return qb.Any(ids)
	//} else {
	//	return qb.All(ids)
	//}
	return 0, nil
}
