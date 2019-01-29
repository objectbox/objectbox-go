package object

import (
	"fmt"
	"time"
)

// converts Unix timestamp in milliseconds (ObjectBox date field) to time.Time
func timeInt64ToEntityProperty(dbValue int64) (goValue time.Time) {
	return time.Unix(dbValue/1000, dbValue%1000*1000000)
}

// converts time.Time to Unix timestamp in milliseconds (internal format expected by ObjectBox on a date field)
func timeInt64ToDatabaseValue(goValue time.Time) int64 {
	var ms = int64(goValue.Nanosecond()) / 1000000
	return goValue.Unix()*1000 + ms
}

func runeIdToEntityProperty(dbValue uint64) (goValue rune) {
	if uint64(rune(dbValue)) != dbValue {
		panic(fmt.Errorf("ID %d out of range for the used type (rune)", dbValue))
	}
	return rune(dbValue)
}

func runeIdToDatabaseValue(goValue rune) uint64 {
	return uint64(goValue)
}
