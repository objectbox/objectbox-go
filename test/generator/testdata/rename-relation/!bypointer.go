package object

// ERROR = can't merge binding model information: uid annotation value must not be empty (model property UID = 5974317550424871033) on property Group, entity NegTaskRelPtr

type NegTaskRelPtr struct {
	Id    uint64
	Group *Group `link uid`
}
