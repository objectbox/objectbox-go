package object

type F struct {
	id uint64 `objectbox:"id"`
	Combined
	BytesValue
	More Combined
}
