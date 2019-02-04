package object

type NegTaskRelValue struct {
	Id    uint64
	Group Group `link uid`
}
