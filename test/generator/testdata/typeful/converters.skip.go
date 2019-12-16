package object

import (
	"time"
)

// converts Unix timestamp in milliseconds (ObjectBox date field) to time.Time
func timeInt64PtrToEntityProperty(dbValue *int64) (*time.Time, error) {
	if dbValue == nil {
		return nil, nil
	}
	var result = time.Unix(*dbValue/1000, *dbValue%1000*1000000).UTC()
	return &result, nil
}

// converts time.Time to Unix timestamp in milliseconds (internal format expected by ObjectBox on a date field)
func timeInt64PtrToDatabaseValue(goValue *time.Time) (*int64, error) {
	if goValue == nil {
		return nil, nil
	}

	var ms = int64(goValue.Nanosecond()) / 1000000
	ms = goValue.Unix()*1000 + ms
	return &ms, nil
}
