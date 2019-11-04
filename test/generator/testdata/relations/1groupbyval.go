package object

//go:generate go run github.com/objectbox/objectbox-go/cmd/objectbox-gogen -byValue

type GroupByVal struct {
	Id   uint64
	Name string
}
