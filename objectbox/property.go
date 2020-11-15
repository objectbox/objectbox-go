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

// BaseProperty serves as a common base for all the property types
type BaseProperty struct {
	Id     TypeId
	Entity *Entity
}

func (property BaseProperty) propertyId() TypeId {
	return property.Id
}

func (property BaseProperty) entityId() TypeId {
	return property.Entity.Id
}

// implementing propertyOrAlias
func (property BaseProperty) alias() *string {
	return nil
}

// IsNil finds entities with the stored property value nil
func (property BaseProperty) IsNil() Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IsNil(&property)
		},
	}
}

// IsNotNil finds entities with the stored property value not nil
func (property BaseProperty) IsNotNil() Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IsNotNil(&property)
		},
	}
}

// Note - the following order* methods are not public because they're only applied selectively
//  e.g. StringVector doesn't support ordering at all

// orderAsc sets ascending order based on this property
func (property BaseProperty) orderAsc() Condition {
	return &orderClosure{
		apply: func(qb *QueryBuilder) error {
			return qb.orderAsc(&property)
		},
	}
}

// orderDesc sets descending order based on this property
func (property BaseProperty) orderDesc() Condition {
	return &orderClosure{
		apply: func(qb *QueryBuilder) error {
			return qb.orderDesc(&property)
		},
	}
}

// orderNilLast puts objects with nil value of the property at the end of the result set
func (property BaseProperty) orderNilLast() Condition {
	return &orderClosure{
		apply: func(qb *QueryBuilder) error {
			return qb.orderNilLast(&property)
		},
	}
}

// orderNilAsZero treats the nil value of the property the same as if it was 0
func (property BaseProperty) orderNilAsZero() Condition {
	return &orderClosure{
		apply: func(qb *QueryBuilder) error {
			return qb.orderNilAsZero(&property)
		},
	}
}

// PropertyString holds information about a property and provides query building methods
type PropertyString struct {
	*BaseProperty
}

// Equals finds entities with the stored property value equal to the given value
func (property PropertyString) Equals(text string, caseSensitive bool) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.StringEquals(property.BaseProperty, text, caseSensitive)
		},
	}
}

// NotEquals finds entities with the stored property value different than the given value
func (property PropertyString) NotEquals(text string, caseSensitive bool) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.StringNotEquals(property.BaseProperty, text, caseSensitive)
		},
	}
}

// Contains finds entities with the stored property value contains the given text
func (property PropertyString) Contains(text string, caseSensitive bool) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.StringContains(property.BaseProperty, text, caseSensitive)
		},
	}
}

// HasPrefix finds entities with the stored property value starts with the given text
func (property PropertyString) HasPrefix(text string, caseSensitive bool) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.StringHasPrefix(property.BaseProperty, text, caseSensitive)
		},
	}
}

// HasSuffix finds entities with the stored property value ends with the given text
func (property PropertyString) HasSuffix(text string, caseSensitive bool) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.StringHasSuffix(property.BaseProperty, text, caseSensitive)
		},
	}
}

// GreaterThan finds entities with the stored property value greater than the given value
func (property PropertyString) GreaterThan(text string, caseSensitive bool) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.StringGreater(property.BaseProperty, text, caseSensitive, false)
		},
	}
}

// GreaterOrEqual finds entities with the stored property value greater than the given value or they're equal
func (property PropertyString) GreaterOrEqual(text string, caseSensitive bool) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.StringGreater(property.BaseProperty, text, caseSensitive, true)
		},
	}
}

// LessThan finds entities with the stored property value less than the given value
func (property PropertyString) LessThan(text string, caseSensitive bool) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.StringLess(property.BaseProperty, text, caseSensitive, false)
		},
	}
}

// LessOrEqual finds entities with the stored property value less than the given value or they're equal
func (property PropertyString) LessOrEqual(text string, caseSensitive bool) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.StringLess(property.BaseProperty, text, caseSensitive, true)
		},
	}
}

// In finds entities with the stored property value equal to any of the given values
// In finds entities with the stored property value equal to any of the given values
func (property PropertyString) In(caseSensitive bool, texts ...string) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.StringIn(property.BaseProperty, texts, caseSensitive)
		},
	}
}

// OrderAsc sets ascending order based on this property
func (property PropertyString) OrderAsc(caseSensitive bool) Condition {
	return &orderClosure{
		apply: func(qb *QueryBuilder) error {
			if err := qb.orderAsc(property.BaseProperty); err != nil {
				return err
			}
			return qb.orderCaseSensitive(property.BaseProperty, caseSensitive)
		},
	}
}

// OrderDesc sets descending order based on this property
func (property PropertyString) OrderDesc(caseSensitive bool) Condition {
	return &orderClosure{
		apply: func(qb *QueryBuilder) error {
			if err := qb.orderDesc(property.BaseProperty); err != nil {
				return err
			}
			return qb.orderCaseSensitive(property.BaseProperty, caseSensitive)
		},
	}
}

// OrderNilLast puts objects with nil value of the property at the end of the result set
func (property PropertyString) OrderNilLast() Condition {
	return property.orderNilLast()
}

// PropertyStringVector holds information about a property and provides query building methods
type PropertyStringVector struct {
	*BaseProperty
}

// Contains finds entities with the stored property value contains the given text
func (property PropertyStringVector) Contains(text string, caseSensitive bool) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.StringVectorContains(property.BaseProperty, text, caseSensitive)
		},
	}
}

// PropertyInt64 holds information about a property and provides query building methods
type PropertyInt64 struct {
	*BaseProperty
}

// Equals finds entities with the stored property value equal to the given value
func (property PropertyInt64) Equals(value int64) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntEqual(property.BaseProperty, value)
		},
	}
}

// NotEquals finds entities with the stored property value different than the given value
func (property PropertyInt64) NotEquals(value int64) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntNotEqual(property.BaseProperty, value)
		},
	}
}

// GreaterThan finds entities with the stored property value greater than the given value
func (property PropertyInt64) GreaterThan(value int64) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntGreater(property.BaseProperty, value, false)
		},
	}
}

// GreaterOrEqual finds entities with the stored property value greater than the given value or they're equal
func (property PropertyInt64) GreaterOrEqual(value int64) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntGreater(property.BaseProperty, value, true)
		},
	}
}

// LessThan finds entities with the stored property value less than the given value
func (property PropertyInt64) LessThan(value int64) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntLess(property.BaseProperty, value, false)
		},
	}
}

// LessOrEqual finds entities with the stored property value less than the given value or they're equal
func (property PropertyInt64) LessOrEqual(value int64) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntLess(property.BaseProperty, value, true)
		},
	}
}

// Between finds entities with the stored property value between a and b (including a and b)
func (property PropertyInt64) Between(a, b int64) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntBetween(property.BaseProperty, a, b)
		},
	}
}

// In finds entities with the stored property value equal to any of the given values
func (property PropertyInt64) In(values ...int64) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.Int64In(property.BaseProperty, values)
		},
	}
}

// NotIn finds entities with the stored property value not equal to any of the given values
func (property PropertyInt64) NotIn(values ...int64) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.Int64NotIn(property.BaseProperty, values)
		},
	}
}

// OrderAsc sets ascending order based on this property
func (property PropertyInt64) OrderAsc() Condition {
	return property.orderAsc()
}

// OrderDesc sets descending order based on this property
func (property PropertyInt64) OrderDesc() Condition {
	return property.orderDesc()
}

// OrderNilLast puts objects with nil value of the property at the end of the result set
func (property PropertyInt64) OrderNilLast() Condition {
	return property.orderNilLast()
}

// OrderNilAsZero treats the nil value of the property the same as if it was 0
func (property PropertyInt64) OrderNilAsZero() Condition {
	return property.orderNilAsZero()
}

// PropertyInt holds information about a property and provides query building methods
type PropertyInt struct {
	*BaseProperty
}

// Equals finds entities with the stored property value equal to the given value
func (property PropertyInt) Equals(value int) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntEqual(property.BaseProperty, int64(value))
		},
	}
}

// NotEquals finds entities with the stored property value different than the given value
func (property PropertyInt) NotEquals(value int) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntNotEqual(property.BaseProperty, int64(value))
		},
	}
}

// GreaterThan finds entities with the stored property value greater than the given value
func (property PropertyInt) GreaterThan(value int) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntGreater(property.BaseProperty, int64(value), false)
		},
	}
}

// GreaterOrEqual finds entities with the stored property value greater than the given value or they're equal
func (property PropertyInt) GreaterOrEqual(value int) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntGreater(property.BaseProperty, int64(value), true)
		},
	}
}

// LessThan finds entities with the stored property value less than the given value
func (property PropertyInt) LessThan(value int) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntLess(property.BaseProperty, int64(value), false)
		},
	}
}

// LessOrEqual finds entities with the stored property value less than the given value or they're equal
func (property PropertyInt) LessOrEqual(value int) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntLess(property.BaseProperty, int64(value), true)
		},
	}
}

// Between finds entities with the stored property value between a and b (including a and b)
func (property PropertyInt) Between(a, b int) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntBetween(property.BaseProperty, int64(a), int64(b))
		},
	}
}

func (property PropertyInt) int64Slice(values []int) []int64 {
	result := make([]int64, len(values))

	for i, v := range values {
		result[i] = int64(v)
	}

	return result
}

// In finds entities with the stored property value equal to any of the given values
func (property PropertyInt) In(values ...int) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.Int64In(property.BaseProperty, property.int64Slice(values))
		},
	}
}

// NotIn finds entities with the stored property value not equal to any of the given values
func (property PropertyInt) NotIn(values ...int) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.Int64NotIn(property.BaseProperty, property.int64Slice(values))
		},
	}
}

// OrderAsc sets ascending order based on this property
func (property PropertyInt) OrderAsc() Condition {
	return property.orderAsc()
}

// OrderDesc sets descending order based on this property
func (property PropertyInt) OrderDesc() Condition {
	return property.orderDesc()
}

// OrderNilLast puts objects with nil value of the property at the end of the result set
func (property PropertyInt) OrderNilLast() Condition {
	return property.orderNilLast()
}

// OrderNilAsZero treats the nil value of the property the same as if it was 0
func (property PropertyInt) OrderNilAsZero() Condition {
	return property.orderNilAsZero()
}

// PropertyUint64 holds information about a property and provides query building methods
type PropertyUint64 struct {
	*BaseProperty
}

// Equals finds entities with the stored property value equal to the given value
func (property PropertyUint64) Equals(value uint64) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntEqual(property.BaseProperty, int64(value))
		},
	}
}

// NotEquals finds entities with the stored property value different than the given value
func (property PropertyUint64) NotEquals(value uint64) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntNotEqual(property.BaseProperty, int64(value))
		},
	}
}

// GreaterThan finds entities with the stored property value greater than the given value
func (property PropertyUint64) GreaterThan(value uint64) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntGreater(property.BaseProperty, int64(value), false)
		},
	}
}

// GreaterOrEqual finds entities with the stored property value greater than the given value or they're equal
func (property PropertyUint64) GreaterOrEqual(value uint64) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntGreater(property.BaseProperty, int64(value), true)
		},
	}
}

// LessThan finds entities with the stored property value less than the given value
func (property PropertyUint64) LessThan(value uint64) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntLess(property.BaseProperty, int64(value), false)
		},
	}
}

// LessOrEqual finds entities with the stored property value less than the given value or they're equal
func (property PropertyUint64) LessOrEqual(value uint64) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntLess(property.BaseProperty, int64(value), true)
		},
	}
}

// Between finds entities with the stored property value between a and b (including a and b)
func (property PropertyUint64) Between(a, b uint64) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntBetween(property.BaseProperty, int64(a), int64(b))
		},
	}
}

func (property PropertyUint64) int64Slice(values []uint64) []int64 {
	result := make([]int64, len(values))

	for i, v := range values {
		result[i] = int64(v)
	}

	return result
}

// In finds entities with the stored property value equal to any of the given values
func (property PropertyUint64) In(values ...uint64) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.Int64In(property.BaseProperty, property.int64Slice(values))
		},
	}
}

// NotIn finds entities with the stored property value not equal to any of the given values
func (property PropertyUint64) NotIn(values ...uint64) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.Int64NotIn(property.BaseProperty, property.int64Slice(values))
		},
	}
}

// OrderAsc sets ascending order based on this property
func (property PropertyUint64) OrderAsc() Condition {
	return property.orderAsc()
}

// OrderDesc sets descending order based on this property
func (property PropertyUint64) OrderDesc() Condition {
	return property.orderDesc()
}

// OrderNilLast puts objects with nil value of the property at the end of the result set
func (property PropertyUint64) OrderNilLast() Condition {
	return property.orderNilLast()
}

// OrderNilAsZero treats the nil value of the property the same as if it was 0
func (property PropertyUint64) OrderNilAsZero() Condition {
	return property.orderNilAsZero()
}

// PropertyUint holds information about a property and provides query building methods
type PropertyUint struct {
	*BaseProperty
}

// Equals finds entities with the stored property value equal to the given value
func (property PropertyUint) Equals(value uint) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntEqual(property.BaseProperty, int64(value))
		},
	}
}

// NotEquals finds entities with the stored property value different than the given value
func (property PropertyUint) NotEquals(value uint) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntNotEqual(property.BaseProperty, int64(value))
		},
	}
}

// GreaterThan finds entities with the stored property value greater than the given value
func (property PropertyUint) GreaterThan(value uint) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntGreater(property.BaseProperty, int64(value), false)
		},
	}
}

// GreaterOrEqual finds entities with the stored property value greater than the given value
func (property PropertyUint) GreaterOrEqual(value uint) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntGreater(property.BaseProperty, int64(value), true)
		},
	}
}

// LessThan finds entities with the stored property value less than the given value
func (property PropertyUint) LessThan(value uint) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntLess(property.BaseProperty, int64(value), false)
		},
	}
}

// LessOrEqual finds entities with the stored property value less than the given value
func (property PropertyUint) LessOrEqual(value uint) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntLess(property.BaseProperty, int64(value), true)
		},
	}
}

// Between finds entities with the stored property value between a and b (including a and b)
func (property PropertyUint) Between(a, b uint) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntBetween(property.BaseProperty, int64(a), int64(b))
		},
	}
}

func (property PropertyUint) int64Slice(values []uint) []int64 {
	result := make([]int64, len(values))

	for i, v := range values {
		result[i] = int64(v)
	}

	return result
}

// In finds entities with the stored property value equal to any of the given values
func (property PropertyUint) In(values ...uint) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.Int64In(property.BaseProperty, property.int64Slice(values))
		},
	}
}

// NotIn finds entities with the stored property value not equal to any of the given values
func (property PropertyUint) NotIn(values ...uint) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.Int64NotIn(property.BaseProperty, property.int64Slice(values))
		},
	}
}

// OrderAsc sets ascending order based on this property
func (property PropertyUint) OrderAsc() Condition {
	return property.orderAsc()
}

// OrderDesc sets descending order based on this property
func (property PropertyUint) OrderDesc() Condition {
	return property.orderDesc()
}

// OrderNilLast puts objects with nil value of the property at the end of the result set
func (property PropertyUint) OrderNilLast() Condition {
	return property.orderNilLast()
}

// OrderNilAsZero treats the nil value of the property the same as if it was 0
func (property PropertyUint) OrderNilAsZero() Condition {
	return property.orderNilAsZero()
}

// PropertyRune holds information about a property and provides query building methods
type PropertyRune struct {
	*BaseProperty
}

// Equals finds entities with the stored property value equal to the given value
func (property PropertyRune) Equals(value rune) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntEqual(property.BaseProperty, int64(value))
		},
	}
}

// NotEquals finds entities with the stored property value different than the given value
func (property PropertyRune) NotEquals(value rune) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntNotEqual(property.BaseProperty, int64(value))
		},
	}
}

// GreaterThan finds entities with the stored property value greater than the given value
func (property PropertyRune) GreaterThan(value rune) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntGreater(property.BaseProperty, int64(value), false)
		},
	}
}

// GreaterOrEqual finds entities with the stored property value greater than the given value
func (property PropertyRune) GreaterOrEqual(value rune) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntGreater(property.BaseProperty, int64(value), true)
		},
	}
}

// LessThan finds entities with the stored property value less than the given value
func (property PropertyRune) LessThan(value rune) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntLess(property.BaseProperty, int64(value), false)
		},
	}
}

// LessOrEqual finds entities with the stored property value less than the given value
func (property PropertyRune) LessOrEqual(value rune) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntLess(property.BaseProperty, int64(value), true)
		},
	}
}

// Between finds entities with the stored property value between a and b (including a and b)
func (property PropertyRune) Between(a, b rune) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntBetween(property.BaseProperty, int64(a), int64(b))
		},
	}
}

func (property PropertyRune) int32Slice(values []rune) []int32 {
	result := make([]int32, len(values))

	for i, v := range values {
		result[i] = int32(v)
	}

	return result
}

// In finds entities with the stored property value equal to any of the given values
func (property PropertyRune) In(values ...rune) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.Int32In(property.BaseProperty, property.int32Slice(values))
		},
	}
}

// NotIn finds entities with the stored property value not equal to any of the given values
func (property PropertyRune) NotIn(values ...rune) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.Int32NotIn(property.BaseProperty, property.int32Slice(values))
		},
	}
}

// OrderAsc sets ascending order based on this property
func (property PropertyRune) OrderAsc() Condition {
	return property.orderAsc()
}

// OrderDesc sets descending order based on this property
func (property PropertyRune) OrderDesc() Condition {
	return property.orderDesc()
}

// OrderNilLast puts objects with nil value of the property at the end of the result set
func (property PropertyRune) OrderNilLast() Condition {
	return property.orderNilLast()
}

// OrderNilAsZero treats the nil value of the property the same as if it was 0
func (property PropertyRune) OrderNilAsZero() Condition {
	return property.orderNilAsZero()
}

// PropertyInt32 holds information about a property and provides query building methods
type PropertyInt32 struct {
	*BaseProperty
}

// Equals finds entities with the stored property value equal to the given value
func (property PropertyInt32) Equals(value int32) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntEqual(property.BaseProperty, int64(value))
		},
	}
}

// NotEquals finds entities with the stored property value different than the given value
func (property PropertyInt32) NotEquals(value int32) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntNotEqual(property.BaseProperty, int64(value))
		},
	}
}

// GreaterThan finds entities with the stored property value greater than the given value
func (property PropertyInt32) GreaterThan(value int32) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntGreater(property.BaseProperty, int64(value), false)
		},
	}
}

// GreaterOrEqual finds entities with the stored property value greater than the given value
func (property PropertyInt32) GreaterOrEqual(value int32) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntGreater(property.BaseProperty, int64(value), true)
		},
	}
}

// LessThan finds entities with the stored property value less than the given value
func (property PropertyInt32) LessThan(value int32) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntLess(property.BaseProperty, int64(value), false)
		},
	}
}

// LessOrEqual finds entities with the stored property value less than the given value
func (property PropertyInt32) LessOrEqual(value int32) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntLess(property.BaseProperty, int64(value), true)
		},
	}
}

// Between finds entities with the stored property value between a and b (including a and b)
func (property PropertyInt32) Between(a, b int32) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntBetween(property.BaseProperty, int64(a), int64(b))
		},
	}
}

// In finds entities with the stored property value equal to any of the given values
func (property PropertyInt32) In(values ...int32) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.Int32In(property.BaseProperty, values)
		},
	}
}

// NotIn finds entities with the stored property value not equal to any of the given values
func (property PropertyInt32) NotIn(values ...int32) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.Int32NotIn(property.BaseProperty, values)
		},
	}
}

// OrderAsc sets ascending order based on this property
func (property PropertyInt32) OrderAsc() Condition {
	return property.orderAsc()
}

// OrderDesc sets descending order based on this property
func (property PropertyInt32) OrderDesc() Condition {
	return property.orderDesc()
}

// OrderNilLast puts objects with nil value of the property at the end of the result set
func (property PropertyInt32) OrderNilLast() Condition {
	return property.orderNilLast()
}

// OrderNilAsZero treats the nil value of the property the same as if it was 0
func (property PropertyInt32) OrderNilAsZero() Condition {
	return property.orderNilAsZero()
}

// PropertyUint32 holds information about a property and provides query building methods
type PropertyUint32 struct {
	*BaseProperty
}

// Equals finds entities with the stored property value equal to the given value
func (property PropertyUint32) Equals(value uint32) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntEqual(property.BaseProperty, int64(value))
		},
	}
}

// NotEquals finds entities with the stored property value different than the given value
func (property PropertyUint32) NotEquals(value uint32) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntNotEqual(property.BaseProperty, int64(value))
		},
	}
}

// GreaterThan finds entities with the stored property value greater than the given value
func (property PropertyUint32) GreaterThan(value uint32) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntGreater(property.BaseProperty, int64(value), false)
		},
	}
}

// GreaterOrEqual finds entities with the stored property value greater than the given value
func (property PropertyUint32) GreaterOrEqual(value uint32) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntGreater(property.BaseProperty, int64(value), true)
		},
	}
}

// LessThan finds entities with the stored property value less than the given value
func (property PropertyUint32) LessThan(value uint32) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntLess(property.BaseProperty, int64(value), false)
		},
	}
}

// LessOrEqual finds entities with the stored property value less than the given value
func (property PropertyUint32) LessOrEqual(value uint32) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntLess(property.BaseProperty, int64(value), true)
		},
	}
}

// Between finds entities with the stored property value between a and b (including a and b)
func (property PropertyUint32) Between(a, b uint32) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntBetween(property.BaseProperty, int64(a), int64(b))
		},
	}
}

func (property PropertyUint32) int32Slice(values []uint32) []int32 {
	result := make([]int32, len(values))

	for i, v := range values {
		result[i] = int32(v)
	}

	return result
}

// In finds entities with the stored property value equal to any of the given values
func (property PropertyUint32) In(values ...uint32) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.Int32In(property.BaseProperty, property.int32Slice(values))
		},
	}
}

// NotIn finds entities with the stored property value not equal to any of the given values
func (property PropertyUint32) NotIn(values ...uint32) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.Int32NotIn(property.BaseProperty, property.int32Slice(values))
		},
	}
}

// OrderAsc sets ascending order based on this property
func (property PropertyUint32) OrderAsc() Condition {
	return property.orderAsc()
}

// OrderDesc sets descending order based on this property
func (property PropertyUint32) OrderDesc() Condition {
	return property.orderDesc()
}

// OrderNilLast puts objects with nil value of the property at the end of the result set
func (property PropertyUint32) OrderNilLast() Condition {
	return property.orderNilLast()
}

// OrderNilAsZero treats the nil value of the property the same as if it was 0
func (property PropertyUint32) OrderNilAsZero() Condition {
	return property.orderNilAsZero()
}

// PropertyInt16 holds information about a property and provides query building methods
type PropertyInt16 struct {
	*BaseProperty
}

// Equals finds entities with the stored property value equal to the given value
func (property PropertyInt16) Equals(value int16) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntEqual(property.BaseProperty, int64(value))
		},
	}
}

// NotEquals finds entities with the stored property value different than the given value
func (property PropertyInt16) NotEquals(value int16) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntNotEqual(property.BaseProperty, int64(value))
		},
	}
}

// GreaterThan finds entities with the stored property value greater than the given value
func (property PropertyInt16) GreaterThan(value int16) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntGreater(property.BaseProperty, int64(value), false)
		},
	}
}

// GreaterOrEqual finds entities with the stored property value greater than the given value
func (property PropertyInt16) GreaterOrEqual(value int16) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntGreater(property.BaseProperty, int64(value), true)
		},
	}
}

// LessThan finds entities with the stored property value less than the given value
func (property PropertyInt16) LessThan(value int16) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntLess(property.BaseProperty, int64(value), false)
		},
	}
}

// LessOrEqual finds entities with the stored property value less than the given value
func (property PropertyInt16) LessOrEqual(value int16) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntLess(property.BaseProperty, int64(value), true)
		},
	}
}

// Between finds entities with the stored property value between a and b (including a and b)
func (property PropertyInt16) Between(a, b int16) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntBetween(property.BaseProperty, int64(a), int64(b))
		},
	}
}

// OrderAsc sets ascending order based on this property
func (property PropertyInt16) OrderAsc() Condition {
	return property.orderAsc()
}

// OrderDesc sets descending order based on this property
func (property PropertyInt16) OrderDesc() Condition {
	return property.orderDesc()
}

// OrderNilLast puts objects with nil value of the property at the end of the result set
func (property PropertyInt16) OrderNilLast() Condition {
	return property.orderNilLast()
}

// OrderNilAsZero treats the nil value of the property the same as if it was 0
func (property PropertyInt16) OrderNilAsZero() Condition {
	return property.orderNilAsZero()
}

// PropertyUint16 holds information about a property and provides query building methods
type PropertyUint16 struct {
	*BaseProperty
}

// Equals finds entities with the stored property value equal to the given value
func (property PropertyUint16) Equals(value uint16) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntEqual(property.BaseProperty, int64(value))
		},
	}
}

// NotEquals finds entities with the stored property value different than the given value
func (property PropertyUint16) NotEquals(value uint16) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntNotEqual(property.BaseProperty, int64(value))
		},
	}
}

// GreaterThan finds entities with the stored property value greater than the given value
func (property PropertyUint16) GreaterThan(value uint16) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntGreater(property.BaseProperty, int64(value), false)
		},
	}
}

// GreaterOrEqual finds entities with the stored property value greater than the given value
func (property PropertyUint16) GreaterOrEqual(value uint16) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntGreater(property.BaseProperty, int64(value), true)
		},
	}
}

// LessThan finds entities with the stored property value less than the given value
func (property PropertyUint16) LessThan(value uint16) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntLess(property.BaseProperty, int64(value), false)
		},
	}
}

// LessOrEqual finds entities with the stored property value less than the given value
func (property PropertyUint16) LessOrEqual(value uint16) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntLess(property.BaseProperty, int64(value), true)
		},
	}
}

// Between finds entities with the stored property value between a and b (including a and b)
func (property PropertyUint16) Between(a, b uint16) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntBetween(property.BaseProperty, int64(a), int64(b))
		},
	}
}

// OrderAsc sets ascending order based on this property
func (property PropertyUint16) OrderAsc() Condition {
	return property.orderAsc()
}

// OrderDesc sets descending order based on this property
func (property PropertyUint16) OrderDesc() Condition {
	return property.orderDesc()
}

// OrderNilLast puts objects with nil value of the property at the end of the result set
func (property PropertyUint16) OrderNilLast() Condition {
	return property.orderNilLast()
}

// OrderNilAsZero treats the nil value of the property the same as if it was 0
func (property PropertyUint16) OrderNilAsZero() Condition {
	return property.orderNilAsZero()
}

// PropertyInt8 holds information about a property and provides query building methods
type PropertyInt8 struct {
	*BaseProperty
}

// Equals finds entities with the stored property value equal to the given value
func (property PropertyInt8) Equals(value int8) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntEqual(property.BaseProperty, int64(value))
		},
	}
}

// NotEquals finds entities with the stored property value different than the given value
func (property PropertyInt8) NotEquals(value int8) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntNotEqual(property.BaseProperty, int64(value))
		},
	}
}

// GreaterThan finds entities with the stored property value greater than the given value
func (property PropertyInt8) GreaterThan(value int8) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntGreater(property.BaseProperty, int64(value), false)
		},
	}
}

// GreaterOrEqual finds entities with the stored property value greater than the given value
func (property PropertyInt8) GreaterOrEqual(value int8) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntGreater(property.BaseProperty, int64(value), true)
		},
	}
}

// LessThan finds entities with the stored property value less than the given value
func (property PropertyInt8) LessThan(value int8) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntLess(property.BaseProperty, int64(value), false)
		},
	}
}

// LessOrEqual finds entities with the stored property value less than the given value
func (property PropertyInt8) LessOrEqual(value int8) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntLess(property.BaseProperty, int64(value), true)
		},
	}
}

// Between finds entities with the stored property value between a and b (including a and b)
func (property PropertyInt8) Between(a, b int8) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntBetween(property.BaseProperty, int64(a), int64(b))
		},
	}
}

// OrderAsc sets ascending order based on this property
func (property PropertyInt8) OrderAsc() Condition {
	return property.orderAsc()
}

// OrderDesc sets descending order based on this property
func (property PropertyInt8) OrderDesc() Condition {
	return property.orderDesc()
}

// OrderNilLast puts objects with nil value of the property at the end of the result set
func (property PropertyInt8) OrderNilLast() Condition {
	return property.orderNilLast()
}

// OrderNilAsZero treats the nil value of the property the same as if it was 0
func (property PropertyInt8) OrderNilAsZero() Condition {
	return property.orderNilAsZero()
}

// PropertyUint8 holds information about a property and provides query building methods
type PropertyUint8 struct {
	*BaseProperty
}

// Equals finds entities with the stored property value equal to the given value
func (property PropertyUint8) Equals(value uint8) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntEqual(property.BaseProperty, int64(value))
		},
	}
}

// NotEquals finds entities with the stored property value different than the given value
func (property PropertyUint8) NotEquals(value uint8) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntNotEqual(property.BaseProperty, int64(value))
		},
	}
}

// GreaterThan finds entities with the stored property value greater than the given value
func (property PropertyUint8) GreaterThan(value uint8) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntGreater(property.BaseProperty, int64(value), false)
		},
	}
}

// GreaterOrEqual finds entities with the stored property value greater than the given value
func (property PropertyUint8) GreaterOrEqual(value uint8) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntGreater(property.BaseProperty, int64(value), true)
		},
	}
}

// LessThan finds entities with the stored property value less than the given value
func (property PropertyUint8) LessThan(value uint8) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntLess(property.BaseProperty, int64(value), false)
		},
	}
}

// LessOrEqual finds entities with the stored property value less than the given value
func (property PropertyUint8) LessOrEqual(value uint8) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntLess(property.BaseProperty, int64(value), true)
		},
	}
}

// Between finds entities with the stored property value between a and b (including a and b)
func (property PropertyUint8) Between(a, b uint8) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntBetween(property.BaseProperty, int64(a), int64(b))
		},
	}
}

// OrderAsc sets ascending order based on this property
func (property PropertyUint8) OrderAsc() Condition {
	return property.orderAsc()
}

// OrderDesc sets descending order based on this property
func (property PropertyUint8) OrderDesc() Condition {
	return property.orderDesc()
}

// OrderNilLast puts objects with nil value of the property at the end of the result set
func (property PropertyUint8) OrderNilLast() Condition {
	return property.orderNilLast()
}

// OrderNilAsZero treats the nil value of the property the same as if it was 0
func (property PropertyUint8) OrderNilAsZero() Condition {
	return property.orderNilAsZero()
}

// PropertyByte holds information about a property and provides query building methods
type PropertyByte struct {
	*BaseProperty
}

// Equals finds entities with the stored property value equal to the given value
func (property PropertyByte) Equals(value byte) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntEqual(property.BaseProperty, int64(value))
		},
	}
}

// NotEquals finds entities with the stored property value different than the given value
func (property PropertyByte) NotEquals(value byte) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntNotEqual(property.BaseProperty, int64(value))
		},
	}
}

// GreaterThan finds entities with the stored property value greater than the given value
func (property PropertyByte) GreaterThan(value byte) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntGreater(property.BaseProperty, int64(value), false)
		},
	}
}

// GreaterOrEqual finds entities with the stored property value greater than the given value
func (property PropertyByte) GreaterOrEqual(value byte) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntGreater(property.BaseProperty, int64(value), true)
		},
	}
}

// LessThan finds entities with the stored property value less than the given value
func (property PropertyByte) LessThan(value byte) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntLess(property.BaseProperty, int64(value), false)
		},
	}
}

// LessOrEqual finds entities with the stored property value less than the given value
func (property PropertyByte) LessOrEqual(value byte) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntLess(property.BaseProperty, int64(value), true)
		},
	}
}

// Between finds entities with the stored property value between a and b (including a and b)
func (property PropertyByte) Between(a, b byte) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntBetween(property.BaseProperty, int64(a), int64(b))
		},
	}
}

// OrderAsc sets ascending order based on this property
func (property PropertyByte) OrderAsc() Condition {
	return property.orderAsc()
}

// OrderDesc sets descending order based on this property
func (property PropertyByte) OrderDesc() Condition {
	return property.orderDesc()
}

// OrderNilLast puts objects with nil value of the property at the end of the result set
func (property PropertyByte) OrderNilLast() Condition {
	return property.orderNilLast()
}

// OrderNilAsZero treats the nil value of the property the same as if it was 0
func (property PropertyByte) OrderNilAsZero() Condition {
	return property.orderNilAsZero()
}

// PropertyFloat64 holds information about a property and provides query building methods
type PropertyFloat64 struct {
	*BaseProperty
}

// GreaterThan finds entities with the stored property value greater than the given value
func (property PropertyFloat64) GreaterThan(value float64) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.DoubleGreater(property.BaseProperty, value, false)
		},
	}
}

// GreaterOrEqual finds entities with the stored property value greater than the given value
func (property PropertyFloat64) GreaterOrEqual(value float64) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.DoubleGreater(property.BaseProperty, value, true)
		},
	}
}

// LessThan finds entities with the stored property value less than the given value
func (property PropertyFloat64) LessThan(value float64) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.DoubleLess(property.BaseProperty, value, false)
		},
	}
}

// LessOrEqual finds entities with the stored property value less than the given value
func (property PropertyFloat64) LessOrEqual(value float64) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.DoubleLess(property.BaseProperty, value, true)
		},
	}
}

// Between finds entities with the stored property value between a and b (including a and b)
func (property PropertyFloat64) Between(a, b float64) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.DoubleBetween(property.BaseProperty, a, b)
		},
	}
}

// OrderAsc sets ascending order based on this property
func (property PropertyFloat64) OrderAsc() Condition {
	return property.orderAsc()
}

// OrderDesc sets descending order based on this property
func (property PropertyFloat64) OrderDesc() Condition {
	return property.orderDesc()
}

// OrderNilLast puts objects with nil value of the property at the end of the result set
func (property PropertyFloat64) OrderNilLast() Condition {
	return property.orderNilLast()
}

// OrderNilAsZero treats the nil value of the property the same as if it was 0
func (property PropertyFloat64) OrderNilAsZero() Condition {
	return property.orderNilAsZero()
}

// PropertyFloat32 holds information about a property and provides query building methods
type PropertyFloat32 struct {
	*BaseProperty
}

// GreaterThan finds entities with the stored property value greater than the given value
func (property PropertyFloat32) GreaterThan(value float32) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.DoubleGreater(property.BaseProperty, float64(value), false)
		},
	}
}

// GreaterOrEqual finds entities with the stored property value greater than the given value
func (property PropertyFloat32) GreaterOrEqual(value float32) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.DoubleGreater(property.BaseProperty, float64(value), true)
		},
	}
}

// LessThan finds entities with the stored property value less than the given value
func (property PropertyFloat32) LessThan(value float32) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.DoubleLess(property.BaseProperty, float64(value), false)
		},
	}
}

// LessOrEqual finds entities with the stored property value less than the given value
func (property PropertyFloat32) LessOrEqual(value float32) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.DoubleLess(property.BaseProperty, float64(value), true)
		},
	}
}

// Between finds entities with the stored property value between a and b (including a and b)
func (property PropertyFloat32) Between(a, b float32) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.DoubleBetween(property.BaseProperty, float64(a), float64(b))
		},
	}
}

// OrderAsc sets ascending order based on this property
func (property PropertyFloat32) OrderAsc() Condition {
	return property.orderAsc()
}

// OrderDesc sets descending order based on this property
func (property PropertyFloat32) OrderDesc() Condition {
	return property.orderDesc()
}

// OrderNilLast puts objects with nil value of the property at the end of the result set
func (property PropertyFloat32) OrderNilLast() Condition {
	return property.orderNilLast()
}

// OrderNilAsZero treats the nil value of the property the same as if it was 0
func (property PropertyFloat32) OrderNilAsZero() Condition {
	return property.orderNilAsZero()
}

// PropertyByteVector holds information about a property and provides query building methods
type PropertyByteVector struct {
	*BaseProperty
}

// Equals finds entities with the stored property value equal to the given value
func (property PropertyByteVector) Equals(value []byte) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.BytesEqual(property.BaseProperty, value)
		},
	}
}

// GreaterThan finds entities with the stored property value greater than the given value
func (property PropertyByteVector) GreaterThan(value []byte) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.BytesGreater(property.BaseProperty, value, false)
		},
	}
}

// GreaterOrEqual finds entities with the stored property value greater than the given value or they're equal
func (property PropertyByteVector) GreaterOrEqual(value []byte) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.BytesGreater(property.BaseProperty, value, true)
		},
	}
}

// LessThan finds entities with the stored property value less than the given value
func (property PropertyByteVector) LessThan(value []byte) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.BytesLess(property.BaseProperty, value, false)
		},
	}
}

// LessOrEqual finds entities with the stored property value less than the given value or they're equal
func (property PropertyByteVector) LessOrEqual(value []byte) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			return qb.BytesLess(property.BaseProperty, value, true)
		},
	}
}

// PropertyBool holds information about a property and provides query building methods
type PropertyBool struct {
	*BaseProperty
}

// Equals finds entities with the stored property value equal to the given value
func (property PropertyBool) Equals(value bool) Condition {
	return &conditionClosure{
		apply: func(qb *QueryBuilder) (ConditionId, error) {
			if value {
				return qb.IntEqual(property.BaseProperty, 1)
			}
			return qb.IntEqual(property.BaseProperty, 0)
		},
	}
}

// OrderAsc sets ascending order based on this property
func (property PropertyBool) OrderAsc() Condition {
	return property.orderAsc()
}

// OrderDesc sets descending order based on this property
func (property PropertyBool) OrderDesc() Condition {
	return property.orderDesc()
}

// OrderNilLast puts objects with nil value of the property at the end of the result set
func (property PropertyBool) OrderNilLast() Condition {
	return property.orderNilLast()
}

// OrderNilAsFalse treats the nil value of the property the same as if it was 0
func (property PropertyBool) OrderNilAsFalse() Condition {
	return property.orderNilAsZero()
}
