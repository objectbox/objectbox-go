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
	applyTo(qb *QueryBuilder, isRoot bool) (ConditionId, error)
}

type ConditionId = int32

type conditionClosure struct {
	applyFun func(qb *QueryBuilder) (ConditionId, error)
}

func (condition *conditionClosure) applyTo(qb *QueryBuilder, isRoot bool) (ConditionId, error) {
	return condition.applyFun(qb)
}

// Combines multiple conditions with an operator
type conditionCombination struct {
	or         bool // AND by default
	conditions []Condition
}

func (condition *conditionCombination) applyTo(qb *QueryBuilder, isRoot bool) (ConditionId, error) {
	if len(condition.conditions) == 0 {
		return 0, nil
	} else if len(condition.conditions) == 1 {
		return condition.conditions[0].applyTo(qb, isRoot)
	}

	ids := make([]ConditionId, len(condition.conditions))

	var err error
	for k, sub := range condition.conditions {
		if ids[k], err = sub.applyTo(qb, false); err != nil {
			return 0, err
		}
	}

	if condition.or {
		return qb.Any(ids)
	} else if !isRoot {
		// only necessary to use AND if it's not a root condition group
		return qb.All(ids)
	} else {
		return 0, nil
	}
}

func Any(conditions ...Condition) Condition {
	return &conditionCombination{
		or:         true,
		conditions: conditions,
	}
}

func All(conditions ...Condition) Condition {
	return &conditionCombination{
		conditions: conditions,
	}
}
