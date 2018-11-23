package model

//go:generate objectbox-gogen

type Task struct {
	Id           uint64 `id`
	Text         string
	DateCreated  int64
	DateFinished int64
}
