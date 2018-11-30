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

type PropertyUint64 struct {
	*Property
}

type PropertyInt32 struct {
	*Property
}

type PropertyUint32 struct {
	*Property
}

type PropertyRune struct {
	*Property
}

type PropertyInt16 struct {
	*Property
}

type PropertyUint16 struct {
	*Property
}

type PropertyInt8 struct {
	*Property
}

type PropertyUint8 struct {
	*Property
}

type PropertyInt struct {
	*Property
}

type PropertyUint struct {
	*Property
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
