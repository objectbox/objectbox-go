package object

// Tests all available GO & ObjectBox types
type Typeful struct {
	Id           uint64 `objectbox:"id"` // NOTE ID is currently required
	Int          int
	Int8         int8
	Int16        int16
	Int32        int32
	Int64        int64
	Uint         uint
	Uint8        uint8
	Uint16       uint16
	Uint32       uint32
	Uint64       uint64
	Bool         bool
	String       string
	StringVector []string
	Byte         byte
	ByteVector   []byte
	Rune         rune
	Float32      float32
	Float64      float64
	Date         int64 `objectbox:"date"`
}
