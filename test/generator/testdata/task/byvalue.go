package object

//go:generate objectbox-gogen -byValue

type TaskByValue struct {
	Id uint64
	Name string
}

type TaskStringByValue struct {
	Id string
	Name string
}
