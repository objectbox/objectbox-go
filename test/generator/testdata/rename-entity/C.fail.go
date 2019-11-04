package object

// ERROR = can't merge binding model information: uid annotation value must not be empty (entity not found in the model) on entity C

// negative test, tag `objectbox:"uid"` on an unknown (new) entity
// `objectbox:"uid"`
type C struct {
	Id uint64
}
