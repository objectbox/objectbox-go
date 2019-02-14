package object

// ERROR = can't merge binding model information: uid annotation value must not be empty (entity not found in the model) on entity C

// negative test, tag `uid` on an unknown (new) entity
// `uid`
type C struct {
	Id uint64
}
