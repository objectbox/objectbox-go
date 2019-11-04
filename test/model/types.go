package model

import (
	"bytes"
	"encoding/gob"
	"time"
)

// BaseWithDate model
type BaseWithDate struct {
	Date int64 `objectbox:"date"`
}

// BaseWithValue model
type BaseWithValue struct {
	Value float64
}

// converts Unix timestamp in milliseconds (ObjectBox date field) to time.Time
func timeInt64ToEntityProperty(dbValue int64) (time.Time, error) {
	return time.Unix(dbValue/1000, dbValue%1000*1000000).UTC(), nil
}

// converts time.Time to Unix timestamp in milliseconds (internal format expected by ObjectBox on a date field)
func timeInt64ToDatabaseValue(goValue time.Time) (int64, error) {
	var ms = int64(goValue.Nanosecond()) / 1000000
	return goValue.Unix()*1000 + ms, nil
}

// decodes the given byte slice as a complex number
func complex128BytesToEntityProperty(dbValue []byte) (complex128, error) {
	// NOTE that constructing the decoder each time is inefficient and only serves as an example for the property converters
	var b = bytes.NewBuffer(dbValue)
	var decoder = gob.NewDecoder(b)

	var value complex128
	err := decoder.Decode(&value)
	return value, err
}

// encodes the given complex number as a byte slice
func complex128BytesToDatabaseValue(goValue complex128) ([]byte, error) {
	// NOTE that constructing the encoder each time is inefficient and only serves as an example for the property converters
	var b bytes.Buffer
	var encoder = gob.NewEncoder(&b)
	err := encoder.Encode(goValue)
	return b.Bytes(), err
}
