package object

// ERROR = can't prepare bindings for testdata/embedding/!1.go: duplicate name (note that property names are case insensitive) on property Value, entity Negative1.BytesValue

// both contain Value field but of two distinct types
type Negative1 struct {
	Float64Value `objectbox:"inline"`
	BytesValue `objectbox:"inline"`
}
