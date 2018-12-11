package object

//go:generate objectbox-gogen -byValue

type TaskByValue struct {
	Id uint64
}

type TaskStringByValue struct {
	Id string
}
