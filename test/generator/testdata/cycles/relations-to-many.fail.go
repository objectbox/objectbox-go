package object

// ERROR = relation cycle detected: RelationToManyChainA.BPtrSlice.CPtrSlice.APtrSlice (RelationToManyChainA)

type RelationToManyChainA struct {
	Id        uint64
	BPtrSlice []*RelationToManyChainB
}

type RelationToManyChainB struct {
	Id        uint64
	CPtrSlice []*RelationToManyChainC
}

type RelationToManyChainC struct {
	Id        uint64
	APtrSlice []*RelationToManyChainA
}
