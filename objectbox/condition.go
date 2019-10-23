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

import "errors"

// Condition is used by Query to limit object selection
type Condition interface {
	applyTo(qb *QueryBuilder, isRoot bool) (ConditionId, error)

	// Alias sets a string alias for the given condition. It can later be used in Query.Set*Params() methods
	Alias(alias string) Condition
}

type ConditionId = int32

type conditionClosure struct {
	apply func(qb *QueryBuilder) (ConditionId, error)
	alias *string
}

func (condition *conditionClosure) applyTo(qb *QueryBuilder, isRoot bool) (ConditionId, error) {
	cid, err := condition.apply(qb)
	if err != nil {
		return 0, err
	}

	if condition.alias != nil {
		err = qb.Alias(*condition.alias)
		if err != nil {
			return 0, err
		}
	}

	return cid, nil
}

// Alias sets a string alias for the given condition. It can later be used in Query.Set*Params() methods
func (condition *conditionClosure) Alias(alias string) Condition {
	condition.alias = &alias
	return condition
}

// Combines multiple conditions with an operator
type conditionCombination struct {
	or          bool // AND by default
	conditions  []Condition
	aliasCalled bool
}

// assertNoLinks makes sure there are no links (0 condition IDs) among given conditions
func (*conditionCombination) assertNoLinks(conditionIds []ConditionId) error {
	for _, cid := range conditionIds {
		if cid == 0 {
			return errors.New("using Link inside Any/All is not supported")
		}
	}
	return nil
}

func (condition *conditionCombination) applyTo(qb *QueryBuilder, isRoot bool) (ConditionId, error) {
	if condition.aliasCalled {
		return 0, errors.New("using Alias on a combination of conditions is not supported")
	}

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

	// root All (AND) is implicit so no need to actually combine the conditions
	if isRoot && !condition.or {
		return 0, nil
	}

	if err := condition.assertNoLinks(ids); err != nil {
		return 0, err
	}

	if condition.or {
		return qb.Any(ids)
	}

	return qb.All(ids)
}

// Alias sets a string alias for the given condition. It can later be used in Query.Set*Params() methods
func (condition *conditionCombination) Alias(alias string) Condition {
	condition.aliasCalled = true // this is invalid on condition-combinations
	return condition
}

// Any provides a way to combine multiple query conditions (equivalent to OR logical operator)
func Any(conditions ...Condition) Condition {
	return &conditionCombination{
		or:         true,
		conditions: conditions,
	}
}

// All provides a way to combine multiple query conditions (equivalent to AND logical operator)
func All(conditions ...Condition) Condition {
	return &conditionCombination{
		conditions: conditions,
	}
}
