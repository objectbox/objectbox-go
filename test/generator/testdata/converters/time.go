package object

import "time"

type TimeEntity struct {
	Id   uint64    `id`
	Time time.Time `date converter:"timeInt64" type:"int64"`
}
