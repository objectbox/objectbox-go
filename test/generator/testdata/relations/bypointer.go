package object

type TaskRelPtr struct {
	Id    uint64
	Group *Group `objectbox:"link"`
}
