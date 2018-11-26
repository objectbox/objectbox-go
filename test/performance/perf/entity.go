package perf

//go:generate objectbox-gogen

type Entity struct {
	Id      uint64
	Int32   int32
	Int64   int64
	String  string
	Float64 float64
}
