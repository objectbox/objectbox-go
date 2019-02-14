package object

// ERROR = can't merge binding model information: uid annotation value must not be empty (property not found in the model) on property Group, entity NegTaskRelId

type NegTaskRelId struct {
	Id    uint64
	Group uint64 `link:"Group" uid`
}
