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

import (
	"errors"
	"fmt"
)

// Condition is used by Query to limit object selection or specify their order
type Condition interface {
	applyTo(qb *QueryBuilder, isRoot bool) (ConditionId, error)

	// Alias sets a string alias for the given condition. It can later be used in Query.Set*Params() methods.
	Alias(alias string) Condition

	// As sets a string alias for the given condition. It can later be used in Query.Set*Params() methods.
	As(alias *alias) Condition
}

// ConditionId is a condition identifier type, used when building queries
type ConditionId = int32

const conditionIdFakeOrder = -1
const conditionIdFakeLink = -2

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

// As sets an alias for the given condition. It can later be used in Query.Set*Params() methods.
func (condition *conditionClosure) As(alias *alias) Condition {
	condition.alias = alias.alias()
	return condition
}

// Combines multiple conditions with an operator
type conditionCombination struct {
	or         bool // AND by default
	conditions []Condition
	alias      *string // this is only used to report an error
}

// assertNoLinks makes sure there are no links (0 condition IDs) among given conditions
func (*conditionCombination) assertNoLinks(conditionIds []ConditionId) error {
	for _, cid := range conditionIds {
		if cid == conditionIdFakeLink {
			return errors.New("using Link inside Any/All is not supported")
		}
	}
	return nil
}

func (condition *conditionCombination) applyTo(qb *QueryBuilder, isRoot bool) (ConditionId, error) {
	if condition.alias != nil {
		return 0, fmt.Errorf("using Alias/As(\"%s\") on a combination of conditions is not supported", *condition.alias)
	}

	if len(condition.conditions) == 0 {
		return 0, nil
	} else if len(condition.conditions) == 1 {
		return condition.conditions[0].applyTo(qb, isRoot)
	}

	ids := make([]ConditionId, 0, len(condition.conditions))
	for _, sub := range condition.conditions {
		cid, err := sub.applyTo(qb, false)
		if err != nil {
			return 0, err
		}

		// Skip order pseudo conditions.
		// Note: conditionIdFakeLink is allowed here and is caught below if used in non-root or in an "ALL" combination.
		if cid != conditionIdFakeOrder {
			ids = append(ids, cid)
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

// Alias sets a string alias for the given condition. It can later be used in Query.Set*Params() methods.
// This is an invalid call on a combination of conditions and will result in an error.
func (condition *conditionCombination) Alias(alias string) Condition {
	condition.alias = &alias
	return condition
}

// As sets an alias for the given condition. It can later be used in Query.Set*Params() methods.
// This is an invalid call on a combination of conditions and will result in an error.
func (condition *conditionCombination) As(alias *alias) Condition {
	condition.alias = alias.alias()
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

// implements propertyOrAlias
type alias struct {
	string
}

func (alias) propertyId() TypeId {
	return 0
}

func (alias) entityId() TypeId {
	return 0
}

func (as *alias) alias() *string {
	return &as.string
}

// Alias wraps a string as an identifier usable for Query.Set*Params*() methods.
func Alias(value string) *alias {
	return &alias{value}
}

type orderClosure struct {
	apply func(qb *QueryBuilder) error
	alias *string // this is only used to report an error
}

func (order *orderClosure) applyTo(qb *QueryBuilder, isRoot bool) (ConditionId, error) {
	if order.alias != nil {
		return 0, fmt.Errorf("using Alias/As(\"%s\") on Order*() is not supported", *order.alias)
	}

	return conditionIdFakeOrder, order.apply(qb)
}

// Alias sets a string alias for the given condition. It can later be used in Query.Set*Params() methods.
// This is an invalid call on order definition and will result in an error.
func (order *orderClosure) Alias(alias string) Condition {
	order.alias = &alias
	return order
}

// As sets an alias for the given condition. It can later be used in Query.Set*Params() methods.
// This is an invalid call on order definition and will result in an error.
func (order *orderClosure) As(alias *alias) Condition {
	order.alias = alias.alias()
	return order
}
