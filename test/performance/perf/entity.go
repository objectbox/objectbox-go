package perf

//go:generate objectbox-gogen

type Entity struct {
	Id    uint64
	Value uint32
}
