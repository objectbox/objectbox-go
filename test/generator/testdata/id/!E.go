package object

// ERROR = can't prepare bindings for testdata/id/!E.go: duplicate name (note that property names are case insensitive) on property id, entity E

type E struct {
	Id uint64
	id uint64
}
