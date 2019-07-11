package object

// ERROR = can't merge binding model information: uid annotation value must not be empty (model entity UID = 8717895732742165505) on entity A

// will fail as uid-request (print UID from model)
// `objectbox:"uid"`
type A struct {
	Id uint64
}
