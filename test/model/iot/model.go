package iot

//go:generate objectbox-bindings

type Event struct {
	Id     uint64 `id`
	Device string
	Date   int64 `date`
}

type Reading struct {
	Id   uint64 `id`
	Date int64  `date`

	/// to-one relation
	EventId uint64

	ValueName string

	/// Device sensor data value
	ValueString string

	/// Device sensor data value
	ValueInteger int64

	/// Device sensor data value
	ValueFloating float64
}
