package object

type C struct {
	Id         uint64 // fake
	identifier uint64 `objectbox:"id"` // real one
}
