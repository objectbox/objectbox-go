package object

type TaskRelId struct {
	Id      uint64
	GroupId uint64 `link:"Group"`
}
