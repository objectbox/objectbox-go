package object

type NegTaskRelPtr struct {
	Id    uint64
	Group *Group `link uid`
}
