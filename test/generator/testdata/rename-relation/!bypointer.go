package object

// ERROR = can't merge binding model information: uid annotation value must not be empty (property not found in the model) on property Group, entity NegTaskRelPtr

type NegTaskRelPtr struct {
	Id    uint64
	Group *Group `link uid`
}
