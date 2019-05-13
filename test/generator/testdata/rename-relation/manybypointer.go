package object

type TaskRelManyPtr struct {
	Id        uint64
	GroupsNew []*Group `objectbox:"uid:3930927879439176946"`
}
