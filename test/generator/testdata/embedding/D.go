package object

type D struct {
	*IdAndFloat64Value `objectbox:"inline"` // pointer
}
