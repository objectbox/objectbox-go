package object

// ERROR = can't prepare bindings for testdata/cycles/embedding.fail.go: embedded struct cycle detected: EmbeddingChainA.EmbeddingChainB.EmbeddingChainC on property EmbeddingChainA found in EmbeddingChainA.EmbeddingChainB.EmbeddingChainC

type EmbeddingChainA struct {
	Id uint64
	*EmbeddingChainB
}

type EmbeddingChainB struct {
	*EmbeddingChainC
}

type EmbeddingChainC struct {
	*EmbeddingChainA
}
