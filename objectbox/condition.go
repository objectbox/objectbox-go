package objectbox

type Condition interface {
	build(qb *QueryBuilder) (conditionId, error)
}

type conditionId = int

type ConditionClosure struct {
	buildFun func(qb *QueryBuilder) (conditionId, error)
}

func (condition *ConditionClosure) build(qb *QueryBuilder) (conditionId, error) {
	return condition.buildFun(qb)
}

// Combines multiple conditions with an operator
type ConditionCombination struct {
	or         bool // AND by default
	conditions []Condition
}

func (condition *ConditionCombination) build(qb *QueryBuilder) (conditionId, error) {
	ids := make([]conditionId, len(condition.conditions))

	for _, sub := range condition.conditions {
		if id, err := sub.build(qb); err != nil {
			return 0, err
		} else {
			ids = append(ids, id)
		}
	}

	// TODO
	//if condition.or {
	//	return qb.Any(ids)
	//} else {
	//	return qb.All(ids)
	//}
	return 0, nil
}
