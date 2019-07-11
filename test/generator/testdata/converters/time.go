package object

import "time"

type TimeEntity struct {
	Id   uint64    `objectbox:"id"`
	Time time.Time `objectbox:"date converter:timeInt64 type:int64"`
}
