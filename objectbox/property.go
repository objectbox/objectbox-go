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

// TODO consider not using closures but defining conditions for each operation
// test performance to make an informed decision as that approach requires much more code and is not so clean

type PropertyString struct {
	*BaseProperty
}

func (property PropertyString) Equals(text string, caseSensitive bool) Condition {
	return &conditionClosure{
		func(qb *QueryBuilder) (ConditionId, error) {
			return qb.StringEquals(property.BaseProperty, text, caseSensitive)
		},
	}
}

func (property PropertyString) NotEquals(text string, caseSensitive bool) Condition {
	return &conditionClosure{
		func(qb *QueryBuilder) (ConditionId, error) {
			return qb.StringNotEquals(property.BaseProperty, text, caseSensitive)
		},
	}
}

func (property PropertyString) Contains(text string, caseSensitive bool) Condition {
	return &conditionClosure{
		func(qb *QueryBuilder) (ConditionId, error) {
			return qb.StringContains(property.BaseProperty, text, caseSensitive)
		},
	}
}

func (property PropertyString) HasPrefix(text string, caseSensitive bool) Condition {
	return &conditionClosure{
		func(qb *QueryBuilder) (ConditionId, error) {
			return qb.StringHasPrefix(property.BaseProperty, text, caseSensitive)
		},
	}
}

func (property PropertyString) HasSuffix(text string, caseSensitive bool) Condition {
	return &conditionClosure{
		func(qb *QueryBuilder) (ConditionId, error) {
			return qb.StringHasSuffix(property.BaseProperty, text, caseSensitive)
		},
	}
}

func (property PropertyString) GreaterThan(text string, caseSensitive bool) Condition {
	return &conditionClosure{
		func(qb *QueryBuilder) (ConditionId, error) {
			return qb.StringGreater(property.BaseProperty, text, caseSensitive, false)
		},
	}
}

func (property PropertyString) GreaterOrEqual(text string, caseSensitive bool) Condition {
	return &conditionClosure{
		func(qb *QueryBuilder) (ConditionId, error) {
			return qb.StringGreater(property.BaseProperty, text, caseSensitive, true)
		},
	}
}

func (property PropertyString) LessThan(text string, caseSensitive bool) Condition {
	return &conditionClosure{
		func(qb *QueryBuilder) (ConditionId, error) {
			return qb.StringLess(property.BaseProperty, text, caseSensitive, false)
		},
	}
}

func (property PropertyString) LessOrEqual(text string, caseSensitive bool) Condition {
	return &conditionClosure{
		func(qb *QueryBuilder) (ConditionId, error) {
			return qb.StringLess(property.BaseProperty, text, caseSensitive, true)
		},
	}
}

func (property PropertyString) In(caseSensitive bool, texts ...string) Condition {
	return &conditionClosure{
		func(qb *QueryBuilder) (ConditionId, error) {
			return qb.StringIn(property.BaseProperty, texts, caseSensitive)
		},
	}
}

type PropertyStringVector struct {
	*BaseProperty
}

func (property PropertyStringVector) Contains(text string, caseSensitive bool) Condition {
	return &conditionClosure{
		func(qb *QueryBuilder) (ConditionId, error) {
			return qb.StringVectorContains(property.BaseProperty, text, caseSensitive)
		},
	}
}

type PropertyInt64 struct {
	*BaseProperty
}

func (property PropertyInt64) Equals(value int64) Condition {
	return &conditionClosure{
		func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntEqual(property.BaseProperty, value)
		},
	}
}

func (property PropertyInt64) NotEquals(value int64) Condition {
	return &conditionClosure{
		func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntNotEqual(property.BaseProperty, value)
		},
	}
}

func (property PropertyInt64) GreaterThan(value int64) Condition {
	return &conditionClosure{
		func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntGreater(property.BaseProperty, value)
		},
	}
}

func (property PropertyInt64) LessThan(value int64) Condition {
	return &conditionClosure{
		func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntLess(property.BaseProperty, value)
		},
	}
}

func (property PropertyInt64) Between(a, b int64) Condition {
	return &conditionClosure{
		func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntBetween(property.BaseProperty, a, b)
		},
	}
}

func (property PropertyInt64) In(values ...int64) Condition {
	return &conditionClosure{
		func(qb *QueryBuilder) (ConditionId, error) {
			return qb.Int64In(property.BaseProperty, values)
		},
	}
}

func (property PropertyInt64) NotIn(values ...int64) Condition {
	return &conditionClosure{
		func(qb *QueryBuilder) (ConditionId, error) {
			return qb.Int64NotIn(property.BaseProperty, values)
		},
	}
}

type PropertyInt struct {
	*BaseProperty
}

func (property PropertyInt) Equals(value int) Condition {
	return &conditionClosure{
		func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntEqual(property.BaseProperty, int64(value))
		},
	}
}

func (property PropertyInt) NotEquals(value int) Condition {
	return &conditionClosure{
		func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntNotEqual(property.BaseProperty, int64(value))
		},
	}
}

func (property PropertyInt) GreaterThan(value int) Condition {
	return &conditionClosure{
		func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntGreater(property.BaseProperty, int64(value))
		},
	}
}

func (property PropertyInt) LessThan(value int) Condition {
	return &conditionClosure{
		func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntLess(property.BaseProperty, int64(value))
		},
	}
}

func (property PropertyInt) Between(a, b int) Condition {
	return &conditionClosure{
		func(qb *QueryBuilder) (ConditionId, error) {
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

func (property PropertyInt) In(values ...int) Condition {
	return &conditionClosure{
		func(qb *QueryBuilder) (ConditionId, error) {
			return qb.Int64In(property.BaseProperty, property.int64Slice(values))
		},
	}
}

func (property PropertyInt) NotIn(values ...int) Condition {
	return &conditionClosure{
		func(qb *QueryBuilder) (ConditionId, error) {
			return qb.Int64NotIn(property.BaseProperty, property.int64Slice(values))
		},
	}
}

type PropertyUint64 struct {
	*BaseProperty
}

func (property PropertyUint64) Equals(value uint64) Condition {
	return &conditionClosure{
		func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntEqual(property.BaseProperty, int64(value))
		},
	}
}

func (property PropertyUint64) NotEquals(value uint64) Condition {
	return &conditionClosure{
		func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntNotEqual(property.BaseProperty, int64(value))
		},
	}
}

func (property PropertyUint64) GreaterThan(value uint64) Condition {
	return &conditionClosure{
		func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntGreater(property.BaseProperty, int64(value))
		},
	}
}

func (property PropertyUint64) LessThan(value uint64) Condition {
	return &conditionClosure{
		func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntLess(property.BaseProperty, int64(value))
		},
	}
}

func (property PropertyUint64) Between(a, b uint64) Condition {
	return &conditionClosure{
		func(qb *QueryBuilder) (ConditionId, error) {
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

func (property PropertyUint64) In(values ...uint64) Condition {
	return &conditionClosure{
		func(qb *QueryBuilder) (ConditionId, error) {
			return qb.Int64In(property.BaseProperty, property.int64Slice(values))
		},
	}
}

func (property PropertyUint64) NotIn(values ...uint64) Condition {
	return &conditionClosure{
		func(qb *QueryBuilder) (ConditionId, error) {
			return qb.Int64NotIn(property.BaseProperty, property.int64Slice(values))
		},
	}
}

type PropertyUint struct {
	*BaseProperty
}

func (property PropertyUint) Equals(value uint) Condition {
	return &conditionClosure{
		func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntEqual(property.BaseProperty, int64(value))
		},
	}
}

func (property PropertyUint) NotEquals(value uint) Condition {
	return &conditionClosure{
		func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntNotEqual(property.BaseProperty, int64(value))
		},
	}
}

func (property PropertyUint) GreaterThan(value uint) Condition {
	return &conditionClosure{
		func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntGreater(property.BaseProperty, int64(value))
		},
	}
}

func (property PropertyUint) LessThan(value uint) Condition {
	return &conditionClosure{
		func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntLess(property.BaseProperty, int64(value))
		},
	}
}

func (property PropertyUint) Between(a, b uint) Condition {
	return &conditionClosure{
		func(qb *QueryBuilder) (ConditionId, error) {
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

func (property PropertyUint) In(values ...uint) Condition {
	return &conditionClosure{
		func(qb *QueryBuilder) (ConditionId, error) {
			return qb.Int64In(property.BaseProperty, property.int64Slice(values))
		},
	}
}

func (property PropertyUint) NotIn(values ...uint) Condition {
	return &conditionClosure{
		func(qb *QueryBuilder) (ConditionId, error) {
			return qb.Int64NotIn(property.BaseProperty, property.int64Slice(values))
		},
	}
}

type PropertyRune struct {
	*BaseProperty
}

func (property PropertyRune) Equals(value rune) Condition {
	return &conditionClosure{
		func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntEqual(property.BaseProperty, int64(value))
		},
	}
}

func (property PropertyRune) NotEquals(value rune) Condition {
	return &conditionClosure{
		func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntNotEqual(property.BaseProperty, int64(value))
		},
	}
}

func (property PropertyRune) GreaterThan(value rune) Condition {
	return &conditionClosure{
		func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntGreater(property.BaseProperty, int64(value))
		},
	}
}

func (property PropertyRune) LessThan(value rune) Condition {
	return &conditionClosure{
		func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntLess(property.BaseProperty, int64(value))
		},
	}
}

func (property PropertyRune) Between(a, b rune) Condition {
	return &conditionClosure{
		func(qb *QueryBuilder) (ConditionId, error) {
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

func (property PropertyRune) In(values ...rune) Condition {
	return &conditionClosure{
		func(qb *QueryBuilder) (ConditionId, error) {
			return qb.Int32In(property.BaseProperty, property.int32Slice(values))
		},
	}
}

func (property PropertyRune) NotIn(values ...rune) Condition {
	return &conditionClosure{
		func(qb *QueryBuilder) (ConditionId, error) {
			return qb.Int32NotIn(property.BaseProperty, property.int32Slice(values))
		},
	}
}

type PropertyInt32 struct {
	*BaseProperty
}

func (property PropertyInt32) Equals(value int32) Condition {
	return &conditionClosure{
		func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntEqual(property.BaseProperty, int64(value))
		},
	}
}

func (property PropertyInt32) NotEquals(value int32) Condition {
	return &conditionClosure{
		func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntNotEqual(property.BaseProperty, int64(value))
		},
	}
}

func (property PropertyInt32) GreaterThan(value int32) Condition {
	return &conditionClosure{
		func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntGreater(property.BaseProperty, int64(value))
		},
	}
}

func (property PropertyInt32) LessThan(value int32) Condition {
	return &conditionClosure{
		func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntLess(property.BaseProperty, int64(value))
		},
	}
}

func (property PropertyInt32) Between(a, b int32) Condition {
	return &conditionClosure{
		func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntBetween(property.BaseProperty, int64(a), int64(b))
		},
	}
}

func (property PropertyInt32) In(values ...int32) Condition {
	return &conditionClosure{
		func(qb *QueryBuilder) (ConditionId, error) {
			return qb.Int32In(property.BaseProperty, values)
		},
	}
}

func (property PropertyInt32) NotIn(values ...int32) Condition {
	return &conditionClosure{
		func(qb *QueryBuilder) (ConditionId, error) {
			return qb.Int32NotIn(property.BaseProperty, values)
		},
	}
}

type PropertyUint32 struct {
	*BaseProperty
}

func (property PropertyUint32) Equals(value uint32) Condition {
	return &conditionClosure{
		func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntEqual(property.BaseProperty, int64(value))
		},
	}
}

func (property PropertyUint32) NotEquals(value uint32) Condition {
	return &conditionClosure{
		func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntNotEqual(property.BaseProperty, int64(value))
		},
	}
}

func (property PropertyUint32) GreaterThan(value uint32) Condition {
	return &conditionClosure{
		func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntGreater(property.BaseProperty, int64(value))
		},
	}
}

func (property PropertyUint32) LessThan(value uint32) Condition {
	return &conditionClosure{
		func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntLess(property.BaseProperty, int64(value))
		},
	}
}

func (property PropertyUint32) Between(a, b uint32) Condition {
	return &conditionClosure{
		func(qb *QueryBuilder) (ConditionId, error) {
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

func (property PropertyUint32) In(values ...uint32) Condition {
	return &conditionClosure{
		func(qb *QueryBuilder) (ConditionId, error) {
			return qb.Int32In(property.BaseProperty, property.int32Slice(values))
		},
	}
}

func (property PropertyUint32) NotIn(values ...uint32) Condition {
	return &conditionClosure{
		func(qb *QueryBuilder) (ConditionId, error) {
			return qb.Int32NotIn(property.BaseProperty, property.int32Slice(values))
		},
	}
}

type PropertyInt16 struct {
	*BaseProperty
}

func (property PropertyInt16) Equals(value int16) Condition {
	return &conditionClosure{
		func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntEqual(property.BaseProperty, int64(value))
		},
	}
}

func (property PropertyInt16) NotEquals(value int16) Condition {
	return &conditionClosure{
		func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntNotEqual(property.BaseProperty, int64(value))
		},
	}
}

func (property PropertyInt16) GreaterThan(value int16) Condition {
	return &conditionClosure{
		func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntGreater(property.BaseProperty, int64(value))
		},
	}
}

func (property PropertyInt16) LessThan(value int16) Condition {
	return &conditionClosure{
		func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntLess(property.BaseProperty, int64(value))
		},
	}
}

func (property PropertyInt16) Between(a, b int16) Condition {
	return &conditionClosure{
		func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntBetween(property.BaseProperty, int64(a), int64(b))
		},
	}
}

type PropertyUint16 struct {
	*BaseProperty
}

func (property PropertyUint16) Equals(value uint16) Condition {
	return &conditionClosure{
		func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntEqual(property.BaseProperty, int64(value))
		},
	}
}

func (property PropertyUint16) NotEquals(value uint16) Condition {
	return &conditionClosure{
		func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntNotEqual(property.BaseProperty, int64(value))
		},
	}
}

func (property PropertyUint16) GreaterThan(value uint16) Condition {
	return &conditionClosure{
		func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntGreater(property.BaseProperty, int64(value))
		},
	}
}

func (property PropertyUint16) LessThan(value uint16) Condition {
	return &conditionClosure{
		func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntLess(property.BaseProperty, int64(value))
		},
	}
}

func (property PropertyUint16) Between(a, b uint16) Condition {
	return &conditionClosure{
		func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntBetween(property.BaseProperty, int64(a), int64(b))
		},
	}
}

type PropertyInt8 struct {
	*BaseProperty
}

func (property PropertyInt8) Equals(value int8) Condition {
	return &conditionClosure{
		func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntEqual(property.BaseProperty, int64(value))
		},
	}
}

func (property PropertyInt8) NotEquals(value int8) Condition {
	return &conditionClosure{
		func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntNotEqual(property.BaseProperty, int64(value))
		},
	}
}

func (property PropertyInt8) GreaterThan(value int8) Condition {
	return &conditionClosure{
		func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntGreater(property.BaseProperty, int64(value))
		},
	}
}

func (property PropertyInt8) LessThan(value int8) Condition {
	return &conditionClosure{
		func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntLess(property.BaseProperty, int64(value))
		},
	}
}

func (property PropertyInt8) Between(a, b int8) Condition {
	return &conditionClosure{
		func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntBetween(property.BaseProperty, int64(a), int64(b))
		},
	}
}

type PropertyUint8 struct {
	*BaseProperty
}

func (property PropertyUint8) Equals(value uint8) Condition {
	return &conditionClosure{
		func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntEqual(property.BaseProperty, int64(value))
		},
	}
}

func (property PropertyUint8) NotEquals(value uint8) Condition {
	return &conditionClosure{
		func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntNotEqual(property.BaseProperty, int64(value))
		},
	}
}

func (property PropertyUint8) GreaterThan(value uint8) Condition {
	return &conditionClosure{
		func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntGreater(property.BaseProperty, int64(value))
		},
	}
}

func (property PropertyUint8) LessThan(value uint8) Condition {
	return &conditionClosure{
		func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntLess(property.BaseProperty, int64(value))
		},
	}
}

func (property PropertyUint8) Between(a, b uint8) Condition {
	return &conditionClosure{
		func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntBetween(property.BaseProperty, int64(a), int64(b))
		},
	}
}

type PropertyByte struct {
	*BaseProperty
}

func (property PropertyByte) Equals(value byte) Condition {
	return &conditionClosure{
		func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntEqual(property.BaseProperty, int64(value))
		},
	}
}

func (property PropertyByte) NotEquals(value byte) Condition {
	return &conditionClosure{
		func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntNotEqual(property.BaseProperty, int64(value))
		},
	}
}

func (property PropertyByte) GreaterThan(value byte) Condition {
	return &conditionClosure{
		func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntGreater(property.BaseProperty, int64(value))
		},
	}
}

func (property PropertyByte) LessThan(value byte) Condition {
	return &conditionClosure{
		func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntLess(property.BaseProperty, int64(value))
		},
	}
}

func (property PropertyByte) Between(a, b byte) Condition {
	return &conditionClosure{
		func(qb *QueryBuilder) (ConditionId, error) {
			return qb.IntBetween(property.BaseProperty, int64(a), int64(b))
		},
	}
}

type PropertyFloat64 struct {
	*BaseProperty
}

func (property PropertyFloat64) GreaterThan(value float64) Condition {
	return &conditionClosure{
		func(qb *QueryBuilder) (ConditionId, error) {
			return qb.DoubleGreater(property.BaseProperty, value)
		},
	}
}

func (property PropertyFloat64) LessThan(value float64) Condition {
	return &conditionClosure{
		func(qb *QueryBuilder) (ConditionId, error) {
			return qb.DoubleLess(property.BaseProperty, value)
		},
	}
}

func (property PropertyFloat64) Between(a, b float64) Condition {
	return &conditionClosure{
		func(qb *QueryBuilder) (ConditionId, error) {
			return qb.DoubleBetween(property.BaseProperty, a, b)
		},
	}
}

type PropertyFloat32 struct {
	*BaseProperty
}

func (property PropertyFloat32) GreaterThan(value float32) Condition {
	return &conditionClosure{
		func(qb *QueryBuilder) (ConditionId, error) {
			return qb.DoubleGreater(property.BaseProperty, float64(value))
		},
	}
}

func (property PropertyFloat32) LessThan(value float32) Condition {
	return &conditionClosure{
		func(qb *QueryBuilder) (ConditionId, error) {
			return qb.DoubleLess(property.BaseProperty, float64(value))
		},
	}
}

func (property PropertyFloat32) Between(a, b float32) Condition {
	return &conditionClosure{
		func(qb *QueryBuilder) (ConditionId, error) {
			return qb.DoubleBetween(property.BaseProperty, float64(a), float64(b))
		},
	}
}

type PropertyByteVector struct {
	*BaseProperty
}

func (property PropertyByteVector) Equals(value []byte) Condition {
	return &conditionClosure{
		func(qb *QueryBuilder) (ConditionId, error) {
			return qb.BytesEqual(property.BaseProperty, value)
		},
	}
}

func (property PropertyByteVector) GreaterThan(value []byte) Condition {
	return &conditionClosure{
		func(qb *QueryBuilder) (ConditionId, error) {
			return qb.BytesGreater(property.BaseProperty, value, false)
		},
	}
}

func (property PropertyByteVector) GreaterOrEqual(value []byte) Condition {
	return &conditionClosure{
		func(qb *QueryBuilder) (ConditionId, error) {
			return qb.BytesGreater(property.BaseProperty, value, true)
		},
	}
}

func (property PropertyByteVector) LessThan(value []byte) Condition {
	return &conditionClosure{
		func(qb *QueryBuilder) (ConditionId, error) {
			return qb.BytesLess(property.BaseProperty, value, false)
		},
	}
}

func (property PropertyByteVector) LessOrEqual(value []byte) Condition {
	return &conditionClosure{
		func(qb *QueryBuilder) (ConditionId, error) {
			return qb.BytesLess(property.BaseProperty, value, true)
		},
	}
}

type PropertyBool struct {
	*BaseProperty
}

func (property PropertyBool) Equals(value bool) Condition {
	return &conditionClosure{
		func(qb *QueryBuilder) (ConditionId, error) {
			if value {
				return qb.IntEqual(property.BaseProperty, 1)
			}
			return qb.IntEqual(property.BaseProperty, 0)
		},
	}
}
