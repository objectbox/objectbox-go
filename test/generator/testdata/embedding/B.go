package object

type B struct {
	Combined `objectbox:"inline"`
	Name     string
}
