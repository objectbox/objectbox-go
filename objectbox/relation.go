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
	"fmt"
)

type conditionRelationOneToMany struct {
	relation   *RelationToOne
	conditions []Condition
	alias      *string // this is only used to report an error
}

func (condition *conditionRelationOneToMany) applyTo(qb *QueryBuilder, isRoot bool) (ConditionId, error) {
	if condition.alias != nil {
		return 0, fmt.Errorf("using Alias/As(\"%s\") on a OneToMany relation link is not supported", *condition.alias)
	}

	return conditionIdFakeLink, qb.LinkOneToMany(condition.relation, condition.conditions)
}

// Alias sets a string alias for the given condition. It can later be used in Query.Set*Params() methods.
// This is an invalid call on Relation links and will result in an error.
func (condition *conditionRelationOneToMany) Alias(alias string) Condition {
	condition.alias = &alias
	return condition
}

// As sets an alias for the given condition. It can later be used in Query.Set*Params() methods.
// This is an invalid call on Relation links and will result in an error.
func (condition *conditionRelationOneToMany) As(alias *alias) Condition {
	condition.alias = alias.alias()
	return condition
}

// RelationToOne holds information about a relation link on a property.
// It is used in generated entity code, providing a way to create a query across multiple related entities.
// Internally, the property value holds an ID of an object in the target entity.
type RelationToOne struct {
	Property *BaseProperty
	Target   *Entity
}

func (relation RelationToOne) entityId() TypeId {
	return relation.Property.Entity.Id
}

func (relation RelationToOne) propertyId() TypeId {
	return relation.Property.Id
}

// Link creates a connection and takes inner conditions to evaluate on the linked entity.
func (relation *RelationToOne) Link(conditions ...Condition) Condition {
	return &conditionRelationOneToMany{relation: relation, conditions: conditions}
}

// Equals finds entities with relation target ID equal to the given value
func (relation RelationToOne) Equals(value uint64) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntEqual(relation.Property, int64(value))
		},
	}
}

// NotEquals finds entities with relation target ID different than the given value
func (relation RelationToOne) NotEquals(value uint64) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntNotEqual(relation.Property, int64(value))
		},
	}
}

func (relation RelationToOne) int64Slice(values []uint64) []int64 {
	result := make([]int64, len(values))

	for i, v := range values {
		result[i] = int64(v)
	}

	return result
}

// In finds entities with relation target ID equal to any of the given values
func (relation RelationToOne) In(values ...uint64) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.Int64In(relation.Property, relation.int64Slice(values))
		},
	}
}

// NotIn finds entities with relation target ID not equal to any the given values
func (relation RelationToOne) NotIn(values ...uint64) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.Int64NotIn(relation.Property, relation.int64Slice(values))
		},
	}
}

type conditionRelationManyToMany struct {
	relation   *RelationToMany
	conditions []Condition
	alias      *string // this is only used to report an error
}

func (condition *conditionRelationManyToMany) applyTo(qb *QueryBuilder, isRoot bool) (ConditionId, error) {
	if condition.alias != nil {
		return 0, fmt.Errorf("using Alias/As(\"%s\") on a ManyToMany relation link is not supported", *condition.alias)
	}

	return conditionIdFakeLink, qb.LinkManyToMany(condition.relation, condition.conditions)
}

// Alias sets a string alias for the given condition. It can later be used in Query.Set*Params() methods.
// This is an invalid call on Relation links and will result in an error.
func (condition *conditionRelationManyToMany) Alias(alias string) Condition {
	condition.alias = &alias
	return condition
}

// As sets an alias for the given condition. It can later be used in Query.Set*Params() methods.
// This is an invalid call on Relation links and will result in an error.
func (condition *conditionRelationManyToMany) As(alias *alias) Condition {
	condition.alias = alias.alias()
	return condition
}

// RelationToMany holds information about a standalone relation link between two entities.
// It is used in generated entity code, providing a way to create a query across multiple related entities.
// Internally, the relation is stored separately, holding pairs of source & target object IDs.
type RelationToMany struct {
	Id     TypeId
	Source *Entity
	Target *Entity
}

// Link creates a connection and takes inner conditions to evaluate on the linked entity.
func (relation *RelationToMany) Link(conditions ...Condition) Condition {
	return &conditionRelationManyToMany{relation: relation, conditions: conditions}
}

// TODO contains() would make sense for many-to-many (slice)
