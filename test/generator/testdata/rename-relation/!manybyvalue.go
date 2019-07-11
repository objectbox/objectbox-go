package object

// ERROR = can't merge binding model information: uid annotation value must not be empty (model relation UID = 4345851588384648695) on relation Groups, entity NegTaskRelManyValue

type NegTaskRelManyValue struct {
	Id     uint64
	Groups []Group `objectbox:"uid"`
}
