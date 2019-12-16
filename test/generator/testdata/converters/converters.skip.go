package object

import (
	"fmt"
	"time"
)

// converts Unix timestamp in milliseconds (ObjectBox date field) to time.Time
func timeInt64ToEntityProperty(dbValue int64) (time.Time, error) {
	return time.Unix(dbValue/1000, dbValue%1000*1000000).UTC(), nil
}

// converts time.Time to Unix timestamp in milliseconds (internal format expected by ObjectBox on a date field)
func timeInt64ToDatabaseValue(goValue time.Time) (int64, error) {
	var ms = int64(goValue.Nanosecond()) / 1000000
	return goValue.Unix()*1000 + ms, nil
}

func runeIdToEntityProperty(dbValue uint64) (rune, error) {
	if uint64(rune(dbValue)) != dbValue {
		return 0, fmt.Errorf("ID %d out of range for the used type (rune)", dbValue)
	}
	return rune(dbValue), nil
}

func runeIdToDatabaseValue(goValue rune) (uint64, error) {
	return uint64(goValue), nil
}
