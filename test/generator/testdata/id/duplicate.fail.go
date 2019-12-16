package object

// ERROR = can't prepare bindings for testdata/id/duplicate.fail.go: duplicate name (note that property names are case insensitive) on property id found in Duplicate

type Duplicate struct {
	Id uint64
	id uint64
}
