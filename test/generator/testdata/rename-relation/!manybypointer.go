package object

// ERROR = can't merge binding model information: uid annotation value must not be empty (relation not found in the model) on relation Groups, entity NegTaskRelManyPtr

type NegTaskRelManyPtr struct {
	Id     uint64
	Groups []*Group `uid`
}
