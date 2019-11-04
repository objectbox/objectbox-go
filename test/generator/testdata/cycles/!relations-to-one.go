package object

// ERROR = relation cycle detected: RelationToOneChainA.BPtr.CPtr.APtr (RelationToOneChainA)

type RelationToOneChainA struct {
	Id   uint64
	BPtr *RelationToOneChainB `objectbox:"link"`
}

type RelationToOneChainB struct {
	Id   uint64
	CPtr *RelationToOneChainC `objectbox:"link"`
}

type RelationToOneChainC struct {
	Id   uint64
	APtr *RelationToOneChainA `objectbox:"link"`
}
