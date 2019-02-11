package object

type NegTaskRelId struct {
	Id    uint64
	Group uint64 `link:"Group" uid`
}
