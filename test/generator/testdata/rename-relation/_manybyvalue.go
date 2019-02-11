package object

type TaskRelManyValue struct {
	Id     uint64
	Groups []Group `uid`
}
