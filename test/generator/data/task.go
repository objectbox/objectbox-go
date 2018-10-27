package data

//go:generate objectbox-bindings

type Task struct {
	Id   uint64
	Date uint64 `t:date`
	Text string
	//Data []byte
}
