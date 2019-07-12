package object

//go:generate go run github.com/objectbox/objectbox-go/cmd/objectbox-gogen -byValue

type TaskByValue struct {
	Id   uint64
	Name string
}

type TaskStringByValue struct {
	Id   string
	Name string
}
