package object

// ERROR = can't prepare bindings for testdata/id/type-int.fail.go: id field 'Id' has unsupported type 'int' on entity TypeInt - must be one of [int64, uint64, string]

type TypeInt struct {
	Id int `objectbox:"id"`
}
