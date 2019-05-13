package object

type TaskRelValue struct {
	Id    uint64
	Group Group `objectbox:"link"`
}
