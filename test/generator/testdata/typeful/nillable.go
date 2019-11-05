package object

import "time"

// Tests all available GO & ObjectBox types.
// We're using pointers to test the `nil` support
type Nillable struct {
	Id           uint64 `id`
	Int          *int
	Int8         *int8
	Int16        *int16
	Int32        *int32
	Int64        *int64
	Uint         *uint
	Uint8        *uint8
	Uint16       *uint16
	Uint32       *uint32
	Uint64       *uint64
	Bool         *bool
	String       *string
	StringVector *[]string
	Byte         *byte
	ByteVector   *[]byte
	Rune         *rune
	Float32      *float32
	Float64      *float64
	Date         *int64     `objectbox:"date"`
	Time         *time.Time `objectbox:"date converter:timeInt64Ptr type:*int64"`
}
