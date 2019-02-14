package object

// ERROR = can't prepare bindings for testdata/embedding/!2.go: duplicate name (note that property names are case insensitive) on property Id, entity Negative2.IdAndFloat64Value

// duplicate field
type Negative2 struct {
	Id
	IdAndFloat64Value
}
