package object

type TaskRelValue struct {
	Id       uint64
	GroupNew Group `objectbox:"link uid:2661732831099943416"`
}
