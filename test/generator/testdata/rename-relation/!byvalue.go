package object

// ERROR = can't merge binding model information: uid annotation value must not be empty (model property UID = 167566062957544642) on property Group, entity NegTaskRelValue

type NegTaskRelValue struct {
	Id    uint64
	Group Group `objectbox:"link uid"`
}
