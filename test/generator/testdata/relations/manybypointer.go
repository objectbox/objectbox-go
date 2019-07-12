package object

type TaskRelManyPtr struct {
	Id     uint64
	Groups []*Group `objectbox:"lazy"`
}
