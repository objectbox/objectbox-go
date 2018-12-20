package model

import (
	"bytes"
	"encoding/gob"
	"time"
)

type BaseWithDate struct {
	Date int64 `date`
}

type BaseWithValue struct {
	Value float64
}

// converts Unix timestamp in microseconds (ObjectBox date field) to time.Time
func timeInt64ToEntityProperty(dbValue int64) time.Time {
	return time.Unix(dbValue/1000, dbValue%1000*1000)
}

// converts time.Time to Unix timestamp in microseconds (internal format expected by ObjectBox on a date field)
func timeInt64ToDatabaseValue(goValue time.Time) int64 {
	return goValue.Unix() * 1000
}

// decodes the given byte slice as a complex number
func complex128BytesToEntityProperty(dbValue []byte) complex128 {
	// NOTE that constructing the decoder each time is inefficient and only serves as an example for the property converters
	var b = bytes.NewBuffer(dbValue)
	var decoder = gob.NewDecoder(b)

	var value complex128
	if err := decoder.Decode(&value); err != nil {
		panic(err)
	}

	return value
}

// encodes the given complex number as a byte slice
func complex128BytesToDatabaseValue(goValue complex128) []byte {
	// NOTE that constructing the encoder each time is inefficient and only serves as an example for the property converters
	var b bytes.Buffer
	var encoder = gob.NewEncoder(&b)

	if err := encoder.Encode(goValue); err != nil {
		panic(err)
	}

	return b.Bytes()
}
