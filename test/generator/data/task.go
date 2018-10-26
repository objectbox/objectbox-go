package data

//go:generate objectbox-bindings

type Task struct {
	Id   uint64
	Date int64
	Text string
}
