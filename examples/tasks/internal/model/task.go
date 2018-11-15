package model

//go:generate objectbox-bindings

type Task struct {
	Id           uint64 `id`
	Text         string
	DateCreated  int64
	DateFinished int64
}
