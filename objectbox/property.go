package objectbox

type Property struct {
	Id TypeId
}

type PropertyString struct {
	*Property
}

func (property PropertyString) StartsWith(text string, caseSensitive bool) Condition {
	// TODO consider not using closures but defining conditions for each operation
	// test performance to make an informed decision as that approach requires much more code and is not so clear
	return &ConditionClosure{
		func(qb *QueryBuilder) (builtConditionId, error) {
			qb.StringStartsWith(property.Id, text, caseSensitive)
			return builtConditionId(qb.cLastCondition), qb.Err
		},
	}
}
