package object

type Empty struct {
}

type Id struct {
	Id uint64
}

type Float64Value struct {
	Value float64 `unique`
}

type BytesValue struct {
	Value []byte
}

type IdAndFloat64Value struct {
	Id    uint64
	Value float64
}

type Combined struct {
	Text string
	Empty
	Id
	Float64Value
}

type val uint64
