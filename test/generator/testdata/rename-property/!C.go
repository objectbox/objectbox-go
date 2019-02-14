package object

// ERROR = can't merge binding model information: uid annotation value must not be empty (property not found in the model) on property New, entity C

// negative test, tag `uid` on an unknown property
type C struct {
	Id  uint64
	New string `uid`
}
