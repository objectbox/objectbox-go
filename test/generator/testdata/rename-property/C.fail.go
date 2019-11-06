package object

// ERROR = can't merge binding model information: uid annotation value must not be empty on an unknown property New, entity C

// negative test, tag `objectbox:"uid"` on an unknown property
type C struct {
	Id  uint64
	New string `objectbox:"uid"`
}
