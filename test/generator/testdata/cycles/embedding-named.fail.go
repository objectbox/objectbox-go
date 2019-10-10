package object

// ERROR = can't prepare bindings for testdata/cycles/embedding-named.fail.go: embedded struct cycle detected: EmbeddingNamedChainA.BPtr.CPtr on property APtr found in EmbeddingNamedChainA.BPtr.CPtr

type EmbeddingNamedChainA struct {
	Id   uint64
	BPtr *EmbeddingNamedChainB
}

type EmbeddingNamedChainB struct {
	CPtr *EmbeddingNamedChainC
}

type EmbeddingNamedChainC struct {
	APtr *EmbeddingNamedChainA
}
