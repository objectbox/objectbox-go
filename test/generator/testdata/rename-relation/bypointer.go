package object

type TaskRelPtr struct {
	Id       uint64
	GroupNew *Group `objectbox:"link uid:1774932891286980153"`
}
