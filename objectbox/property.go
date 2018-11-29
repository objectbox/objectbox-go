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
