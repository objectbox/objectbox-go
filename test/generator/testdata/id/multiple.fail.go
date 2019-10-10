package object

// ERROR = can't prepare bindings for testdata/id/multiple.fail.go: struct Multiple has multiple ID properties - Id and id

type Multiple struct {
	Id uint64 `objectbox:"id"`
	id uint64 `objectbox:"id"`
}
