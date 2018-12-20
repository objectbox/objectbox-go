package object

import "time"

// converts Unix timestamp in microseconds (ObjectBox date field) to time.Time
func timeInt64ToEntityProperty(dbValue int64) (goValue time.Time) {
	return time.Unix(dbValue/1000, dbValue%1000*1000)
}

// converts time.Time to Unix timestamp in microseconds (internal format expected by ObjectBox on a date field)
func timeInt64ToDatabaseValue(goValue time.Time) (dbValue int64) {
	return goValue.Unix() * 1000
}
