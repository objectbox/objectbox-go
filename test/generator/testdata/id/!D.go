package object

// ERROR = can't prepare bindings for testdata/id/!D.go: struct D has multiple ID properties - Id and id

type D struct {
	Id uint64 `objectbox:"id"`
	id uint64 `objectbox:"id"`
}
