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

func (property PropertyString) StartsWith(text string, caseSensitive bool) Condition {
	return &ConditionClosure{
		func(qb *queryBuilder) (conditionId, error) {
			return qb.StringStartsWith(property.Id, text, caseSensitive)
		},
	}
}

type PropertyUint64 struct {
	*Property
}

type PropertyInt64 struct {
	*Property
}

type PropertyUint32 struct {
	*Property
}

type PropertyInt32 struct {
	*Property
}

type PropertyUint16 struct {
	*Property
}

type PropertyInt16 struct {
	*Property
}

type PropertyUint8 struct {
	*Property
}

type PropertyInt8 struct {
	*Property
}

type PropertyUint struct {
	*Property
}

type PropertyInt struct {
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
