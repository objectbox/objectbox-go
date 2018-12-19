package model

type BaseWithDate struct {
	Id   uint64
	Date int64 `date`
}

type BaseWithValue struct {
	Id    uint64
	Value float64
}
