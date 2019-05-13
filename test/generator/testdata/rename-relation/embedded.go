package object

type TaskRelEmbedded struct {
	Id        uint64
	WithGroup `objectbox:"inline"`
}
