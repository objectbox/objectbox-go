package object

type TaskIndexed struct {
	Id       uint64 `objectbox:"id"`
	Uid      string `objectbox:"unique"`
	Name     string `objectbox:"index"` // uses HASH as default
	Priority int    `objectbox:"index"` // uses VALUE as default
	Group    string `objectbox:"index:value"`
	Place    string `objectbox:"index:hash"`
	Source   string `objectbox:"index:hash64"`
}
