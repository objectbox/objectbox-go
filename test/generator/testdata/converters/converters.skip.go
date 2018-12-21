package object

import "time"

// converts Unix timestamp in milliseconds (ObjectBox date field) to time.Time
func timeInt64ToEntityProperty(dbValue int64) (goValue time.Time) {
	return time.Unix(dbValue/1000, dbValue%1000*1000000)
}

// converts time.Time to Unix timestamp in milliseconds (internal format expected by ObjectBox on a date field)
func timeInt64ToDatabaseValue(goValue time.Time) int64 {
	var ms = int64(goValue.Nanosecond()) / 1000000
	return goValue.Unix()*1000 + ms
}
