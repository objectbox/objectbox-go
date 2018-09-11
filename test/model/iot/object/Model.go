package object

type Event struct {
	Id     uint64
	Device string
	Date   int64
}

type Reading struct {
	Id   uint64
	Date int64

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
