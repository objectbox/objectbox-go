package model

import (
	"bytes"
	"encoding/gob"
)

// BaseWithDate model
type BaseWithDate struct {
	Date int64 `objectbox:"date"`
}

// BaseWithValue model
type BaseWithValue struct {
	Value float64
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
