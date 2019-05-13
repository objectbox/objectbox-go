package object

type TaskRelId struct {
	Id       uint64
	GroupNew uint64 `objectbox:"link:Group uid:6050128673802995827"`
}
