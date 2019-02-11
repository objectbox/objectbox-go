package object

type TaskRelId struct {
	Id    uint64
	Group uint64 `link:"Group"`
}
