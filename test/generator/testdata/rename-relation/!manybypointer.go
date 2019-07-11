package object

// ERROR = can't merge binding model information: uid annotation value must not be empty (model relation UID = 8514850266767180993) on relation Groups, entity NegTaskRelManyPtr

type NegTaskRelManyPtr struct {
	Id     uint64
	Groups []*Group `objectbox:"uid"`
}
