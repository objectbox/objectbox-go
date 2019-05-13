package object

type A struct {
	Id   `objectbox:"inline"`
	Name string
}
