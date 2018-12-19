package object

type D struct {
	int64 // embed a simple type
	id    // embed a named type
	Id    `int64`
}
