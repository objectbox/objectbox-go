package object

type TaskIndexed struct {
	Id       uint64 `id`
	Uid      string `unique`
	Name     string `index` // uses HASH as default
	Priority int    `index` // uses VALUE as default
	Group    string `index:"value"`
	Place    string `index:"hash"`
	Source   string `index:"hash64"`
}
