package object

type C struct {
	int64 // embed a simple type
	val   // embed a named type
	Id    `objectbox:"inline"`
}
