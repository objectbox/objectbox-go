package objectbox

type Property struct {
	Id     TypeId
	Entity *Entity
}

// TODO consider not using closures but defining conditions for each operation
// test performance to make an informed decision as that approach requires much more code and is not so clean

type PropertyString struct {
	*Property
}

func (property PropertyString) Equal(text string, caseSensitive bool) Condition {
	return &ConditionClosure{
		func(qb *queryBuilder) (conditionId, error) {
			return qb.StringEqual(property.Id, text, caseSensitive)
		},
	}
}

func (property PropertyString) NotEqual(text string, caseSensitive bool) Condition {
	return &ConditionClosure{
		func(qb *queryBuilder) (conditionId, error) {
			return qb.StringNotEqual(property.Id, text, caseSensitive)
		},
	}
}

func (property PropertyString) Contains(text string, caseSensitive bool) Condition {
	return &ConditionClosure{
		func(qb *queryBuilder) (conditionId, error) {
			return qb.StringContains(property.Id, text, caseSensitive)
		},
	}
}

func (property PropertyString) StartsWith(text string, caseSensitive bool) Condition {
	return &ConditionClosure{
		func(qb *queryBuilder) (conditionId, error) {
			return qb.StringStartsWith(property.Id, text, caseSensitive)
		},
	}
}

func (property PropertyString) EndsWith(text string, caseSensitive bool) Condition {
	return &ConditionClosure{
		func(qb *queryBuilder) (conditionId, error) {
			return qb.StringEndsWith(property.Id, text, caseSensitive)
		},
	}
}

func (property PropertyString) GreaterThan(text string, caseSensitive bool) Condition {
	return &ConditionClosure{
		func(qb *queryBuilder) (conditionId, error) {
			return qb.StringGreater(property.Id, text, caseSensitive, false)
		},
	}
}

func (property PropertyString) GreaterOrEqual(text string, caseSensitive bool) Condition {
	return &ConditionClosure{
		func(qb *queryBuilder) (conditionId, error) {
			return qb.StringGreater(property.Id, text, caseSensitive, true)
		},
	}
}

func (property PropertyString) LessThan(text string, caseSensitive bool) Condition {
	return &ConditionClosure{
		func(qb *queryBuilder) (conditionId, error) {
			return qb.StringLess(property.Id, text, caseSensitive, false)
		},
	}
}

func (property PropertyString) LessOrEqual(text string, caseSensitive bool) Condition {
	return &ConditionClosure{
		func(qb *queryBuilder) (conditionId, error) {
			return qb.StringLess(property.Id, text, caseSensitive, true)
		},
	}
}

func (property PropertyString) In(caseSensitive bool, texts ...string) Condition {
	return &ConditionClosure{
		func(qb *queryBuilder) (conditionId, error) {
			return qb.StringIn(property.Id, texts, caseSensitive)
		},
	}
}

type PropertyInt64 struct {
	*Property
}

func (property PropertyInt64) Equal(value int64) Condition {
	return &ConditionClosure{
		func(qb *queryBuilder) (conditionId, error) {
			return qb.IntEqual(property.Id, value)
		},
	}
}

func (property PropertyInt64) NotEqual(value int64) Condition {
	return &ConditionClosure{
		func(qb *queryBuilder) (conditionId, error) {
			return qb.IntNotEqual(property.Id, value)
		},
	}
}

func (property PropertyInt64) GreaterThan(value int64) Condition {
	return &ConditionClosure{
		func(qb *queryBuilder) (conditionId, error) {
			return qb.IntGreater(property.Id, value)
		},
	}
}

func (property PropertyInt64) LessThan(value int64) Condition {
	return &ConditionClosure{
		func(qb *queryBuilder) (conditionId, error) {
			return qb.IntLess(property.Id, value)
		},
	}
}

func (property PropertyInt64) Between(a, b int64) Condition {
	return &ConditionClosure{
		func(qb *queryBuilder) (conditionId, error) {
			return qb.IntBetween(property.Id, a, b)
		},
	}
}

func (property PropertyInt64) In(values ...int64) Condition {
	return &ConditionClosure{
		func(qb *queryBuilder) (conditionId, error) {
			return qb.Int64In(property.Id, values)
		},
	}
}

func (property PropertyInt64) NotIn(values ...int64) Condition {
	return &ConditionClosure{
		func(qb *queryBuilder) (conditionId, error) {
			return qb.Int64NotIn(property.Id, values)
		},
	}
}

type PropertyInt struct {
	*Property
}

func (property PropertyInt) Equal(value int) Condition {
	return &ConditionClosure{
		func(qb *queryBuilder) (conditionId, error) {
			return qb.IntEqual(property.Id, int64(value))
		},
	}
}

func (property PropertyInt) NotEqual(value int) Condition {
	return &ConditionClosure{
		func(qb *queryBuilder) (conditionId, error) {
			return qb.IntNotEqual(property.Id, int64(value))
		},
	}
}

func (property PropertyInt) GreaterThan(value int) Condition {
	return &ConditionClosure{
		func(qb *queryBuilder) (conditionId, error) {
			return qb.IntGreater(property.Id, int64(value))
		},
	}
}

func (property PropertyInt) LessThan(value int) Condition {
	return &ConditionClosure{
		func(qb *queryBuilder) (conditionId, error) {
			return qb.IntLess(property.Id, int64(value))
		},
	}
}

func (property PropertyInt) Between(a, b int) Condition {
	return &ConditionClosure{
		func(qb *queryBuilder) (conditionId, error) {
			return qb.IntBetween(property.Id, int64(a), int64(b))
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
	return &ConditionClosure{
		func(qb *queryBuilder) (conditionId, error) {
			return qb.Int64In(property.Id, property.int64Slice(values))
		},
	}
}

func (property PropertyInt) NotIn(values ...int) Condition {
	return &ConditionClosure{
		func(qb *queryBuilder) (conditionId, error) {
			return qb.Int64NotIn(property.Id, property.int64Slice(values))
		},
	}
}

type PropertyUint64 struct {
	*Property
}

func (property PropertyUint64) Equal(value uint64) Condition {
	return &ConditionClosure{
		func(qb *queryBuilder) (conditionId, error) {
			return qb.IntEqual(property.Id, int64(value))
		},
	}
}

func (property PropertyUint64) NotEqual(value uint64) Condition {
	return &ConditionClosure{
		func(qb *queryBuilder) (conditionId, error) {
			return qb.IntNotEqual(property.Id, int64(value))
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
	return &ConditionClosure{
		func(qb *queryBuilder) (conditionId, error) {
			return qb.Int64In(property.Id, property.int64Slice(values))
		},
	}
}

func (property PropertyUint64) NotIn(values ...uint64) Condition {
	return &ConditionClosure{
		func(qb *queryBuilder) (conditionId, error) {
			return qb.Int64NotIn(property.Id, property.int64Slice(values))
		},
	}
}

type PropertyUint struct {
	*Property
}

func (property PropertyUint) Equal(value uint) Condition {
	return &ConditionClosure{
		func(qb *queryBuilder) (conditionId, error) {
			return qb.IntEqual(property.Id, int64(value))
		},
	}
}

func (property PropertyUint) NotEqual(value uint) Condition {
	return &ConditionClosure{
		func(qb *queryBuilder) (conditionId, error) {
			return qb.IntNotEqual(property.Id, int64(value))
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
	return &ConditionClosure{
		func(qb *queryBuilder) (conditionId, error) {
			return qb.Int64In(property.Id, property.int64Slice(values))
		},
	}
}

func (property PropertyUint) NotIn(values ...uint) Condition {
	return &ConditionClosure{
		func(qb *queryBuilder) (conditionId, error) {
			return qb.Int64NotIn(property.Id, property.int64Slice(values))
		},
	}
}

type PropertyRune struct {
	*Property
}

func (property PropertyRune) Equal(value rune) Condition {
	return &ConditionClosure{
		func(qb *queryBuilder) (conditionId, error) {
			return qb.IntEqual(property.Id, int64(value))
		},
	}
}

func (property PropertyRune) NotEqual(value rune) Condition {
	return &ConditionClosure{
		func(qb *queryBuilder) (conditionId, error) {
			return qb.IntNotEqual(property.Id, int64(value))
		},
	}
}

func (property PropertyRune) GreaterThan(value rune) Condition {
	return &ConditionClosure{
		func(qb *queryBuilder) (conditionId, error) {
			return qb.IntGreater(property.Id, int64(value))
		},
	}
}

func (property PropertyRune) LessThan(value rune) Condition {
	return &ConditionClosure{
		func(qb *queryBuilder) (conditionId, error) {
			return qb.IntLess(property.Id, int64(value))
		},
	}
}

func (property PropertyRune) Between(a, b rune) Condition {
	return &ConditionClosure{
		func(qb *queryBuilder) (conditionId, error) {
			return qb.IntBetween(property.Id, int64(a), int64(b))
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
	return &ConditionClosure{
		func(qb *queryBuilder) (conditionId, error) {
			return qb.Int32In(property.Id, property.int32Slice(values))
		},
	}
}

func (property PropertyRune) NotIn(values ...rune) Condition {
	return &ConditionClosure{
		func(qb *queryBuilder) (conditionId, error) {
			return qb.Int32NotIn(property.Id, property.int32Slice(values))
		},
	}
}

type PropertyInt32 struct {
	*Property
}

func (property PropertyInt32) Equal(value int32) Condition {
	return &ConditionClosure{
		func(qb *queryBuilder) (conditionId, error) {
			return qb.IntEqual(property.Id, int64(value))
		},
	}
}

func (property PropertyInt32) NotEqual(value int32) Condition {
	return &ConditionClosure{
		func(qb *queryBuilder) (conditionId, error) {
			return qb.IntNotEqual(property.Id, int64(value))
		},
	}
}

func (property PropertyInt32) GreaterThan(value int32) Condition {
	return &ConditionClosure{
		func(qb *queryBuilder) (conditionId, error) {
			return qb.IntGreater(property.Id, int64(value))
		},
	}
}

func (property PropertyInt32) LessThan(value int32) Condition {
	return &ConditionClosure{
		func(qb *queryBuilder) (conditionId, error) {
			return qb.IntLess(property.Id, int64(value))
		},
	}
}

func (property PropertyInt32) Between(a, b int32) Condition {
	return &ConditionClosure{
		func(qb *queryBuilder) (conditionId, error) {
			return qb.IntBetween(property.Id, int64(a), int64(b))
		},
	}
}

func (property PropertyInt32) In(values ...int32) Condition {
	return &ConditionClosure{
		func(qb *queryBuilder) (conditionId, error) {
			return qb.Int32In(property.Id, values)
		},
	}
}

func (property PropertyInt32) NotIn(values ...int32) Condition {
	return &ConditionClosure{
		func(qb *queryBuilder) (conditionId, error) {
			return qb.Int32NotIn(property.Id, values)
		},
	}
}

type PropertyUint32 struct {
	*Property
}

func (property PropertyUint32) Equal(value uint32) Condition {
	return &ConditionClosure{
		func(qb *queryBuilder) (conditionId, error) {
			return qb.IntEqual(property.Id, int64(value))
		},
	}
}

func (property PropertyUint32) NotEqual(value uint32) Condition {
	return &ConditionClosure{
		func(qb *queryBuilder) (conditionId, error) {
			return qb.IntNotEqual(property.Id, int64(value))
		},
	}
}

func (property PropertyUint32) GreaterThan(value uint32) Condition {
	return &ConditionClosure{
		func(qb *queryBuilder) (conditionId, error) {
			return qb.IntGreater(property.Id, int64(value))
		},
	}
}

func (property PropertyUint32) LessThan(value uint32) Condition {
	return &ConditionClosure{
		func(qb *queryBuilder) (conditionId, error) {
			return qb.IntLess(property.Id, int64(value))
		},
	}
}

func (property PropertyUint32) Between(a, b uint32) Condition {
	return &ConditionClosure{
		func(qb *queryBuilder) (conditionId, error) {
			return qb.IntBetween(property.Id, int64(a), int64(b))
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
	return &ConditionClosure{
		func(qb *queryBuilder) (conditionId, error) {
			return qb.Int32In(property.Id, property.int32Slice(values))
		},
	}
}

func (property PropertyUint32) NotIn(values ...uint32) Condition {
	return &ConditionClosure{
		func(qb *queryBuilder) (conditionId, error) {
			return qb.Int32NotIn(property.Id, property.int32Slice(values))
		},
	}
}

type PropertyInt16 struct {
	*Property
}

func (property PropertyInt16) Equal(value int16) Condition {
	return &ConditionClosure{
		func(qb *queryBuilder) (conditionId, error) {
			return qb.IntEqual(property.Id, int64(value))
		},
	}
}

func (property PropertyInt16) NotEqual(value int16) Condition {
	return &ConditionClosure{
		func(qb *queryBuilder) (conditionId, error) {
			return qb.IntNotEqual(property.Id, int64(value))
		},
	}
}

func (property PropertyInt16) GreaterThan(value int16) Condition {
	return &ConditionClosure{
		func(qb *queryBuilder) (conditionId, error) {
			return qb.IntGreater(property.Id, int64(value))
		},
	}
}

func (property PropertyInt16) LessThan(value int16) Condition {
	return &ConditionClosure{
		func(qb *queryBuilder) (conditionId, error) {
			return qb.IntLess(property.Id, int64(value))
		},
	}
}

func (property PropertyInt16) Between(a, b int16) Condition {
	return &ConditionClosure{
		func(qb *queryBuilder) (conditionId, error) {
			return qb.IntBetween(property.Id, int64(a), int64(b))
		},
	}
}

type PropertyUint16 struct {
	*Property
}

func (property PropertyUint16) Equal(value uint16) Condition {
	return &ConditionClosure{
		func(qb *queryBuilder) (conditionId, error) {
			return qb.IntEqual(property.Id, int64(value))
		},
	}
}

func (property PropertyUint16) NotEqual(value uint16) Condition {
	return &ConditionClosure{
		func(qb *queryBuilder) (conditionId, error) {
			return qb.IntNotEqual(property.Id, int64(value))
		},
	}
}

func (property PropertyUint16) GreaterThan(value uint16) Condition {
	return &ConditionClosure{
		func(qb *queryBuilder) (conditionId, error) {
			return qb.IntGreater(property.Id, int64(value))
		},
	}
}

func (property PropertyUint16) LessThan(value uint16) Condition {
	return &ConditionClosure{
		func(qb *queryBuilder) (conditionId, error) {
			return qb.IntLess(property.Id, int64(value))
		},
	}
}

func (property PropertyUint16) Between(a, b uint16) Condition {
	return &ConditionClosure{
		func(qb *queryBuilder) (conditionId, error) {
			return qb.IntBetween(property.Id, int64(a), int64(b))
		},
	}
}

type PropertyInt8 struct {
	*Property
}

func (property PropertyInt8) Equal(value int8) Condition {
	return &ConditionClosure{
		func(qb *queryBuilder) (conditionId, error) {
			return qb.IntEqual(property.Id, int64(value))
		},
	}
}

func (property PropertyInt8) NotEqual(value int8) Condition {
	return &ConditionClosure{
		func(qb *queryBuilder) (conditionId, error) {
			return qb.IntNotEqual(property.Id, int64(value))
		},
	}
}

func (property PropertyInt8) GreaterThan(value int8) Condition {
	return &ConditionClosure{
		func(qb *queryBuilder) (conditionId, error) {
			return qb.IntGreater(property.Id, int64(value))
		},
	}
}

func (property PropertyInt8) LessThan(value int8) Condition {
	return &ConditionClosure{
		func(qb *queryBuilder) (conditionId, error) {
			return qb.IntLess(property.Id, int64(value))
		},
	}
}

func (property PropertyInt8) Between(a, b int8) Condition {
	return &ConditionClosure{
		func(qb *queryBuilder) (conditionId, error) {
			return qb.IntBetween(property.Id, int64(a), int64(b))
		},
	}
}

type PropertyUint8 struct {
	*Property
}

func (property PropertyUint8) Equal(value uint8) Condition {
	return &ConditionClosure{
		func(qb *queryBuilder) (conditionId, error) {
			return qb.IntEqual(property.Id, int64(value))
		},
	}
}

func (property PropertyUint8) NotEqual(value uint8) Condition {
	return &ConditionClosure{
		func(qb *queryBuilder) (conditionId, error) {
			return qb.IntNotEqual(property.Id, int64(value))
		},
	}
}

func (property PropertyUint8) GreaterThan(value uint8) Condition {
	return &ConditionClosure{
		func(qb *queryBuilder) (conditionId, error) {
			return qb.IntGreater(property.Id, int64(value))
		},
	}
}

func (property PropertyUint8) LessThan(value uint8) Condition {
	return &ConditionClosure{
		func(qb *queryBuilder) (conditionId, error) {
			return qb.IntLess(property.Id, int64(value))
		},
	}
}

func (property PropertyUint8) Between(a, b uint8) Condition {
	return &ConditionClosure{
		func(qb *queryBuilder) (conditionId, error) {
			return qb.IntBetween(property.Id, int64(a), int64(b))
		},
	}
}

type PropertyFloat64 struct {
	*Property
}

type PropertyFloat32 struct {
	*Property
}

type PropertyByte struct {
	*Property
}

type PropertyByteVector struct {
	*Property
}

type PropertyBool struct {
	*Property
}
