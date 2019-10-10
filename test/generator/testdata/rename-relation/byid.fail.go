package object

// ERROR = can't merge binding model information: uid annotation value must not be empty (model property UID = 6745438398739480977) on property Group, entity NegTaskRelId

type NegTaskRelId struct {
	Id    uint64
	Group uint64 `objectbox:"link:Group uid"`
}
