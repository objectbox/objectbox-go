package object

// ERROR = can't prepare bindings for testdata/embedding/2.fail.go: duplicate name (note that property names are case insensitive) on property Id found in Negative2.IdAndFloat64Value

// duplicate field
type Negative2 struct {
	Id                `objectbox:"inline"`
	IdAndFloat64Value `objectbox:"inline"`
}
