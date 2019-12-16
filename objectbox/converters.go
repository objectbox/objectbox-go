/*
 * Copyright 2019 ObjectBox Ltd. All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package objectbox

import (
	"fmt"
	"strconv"
	"time"
)

// StringIdConvertToEntityProperty implements "StringIdConvert" property value converter
func StringIdConvertToEntityProperty(dbValue uint64) (string, error) {
	return strconv.FormatUint(dbValue, 10), nil
}

// StringIdConvertToDatabaseValue implements "StringIdConvert" property value converter
func StringIdConvertToDatabaseValue(goValue string) (uint64, error) {
	// in case the object was initialized by the user without setting the ID explicitly
	if goValue == "" {
		return 0, nil
	}
	return strconv.ParseUint(goValue, 10, 64)
}

// TimeInt64ConvertToEntityProperty converts Unix timestamp in milliseconds (ObjectBox date field) to time.Time
// NOTE - you lose precision - anything smaller then milliseconds is dropped
func TimeInt64ConvertToEntityProperty(dbValue int64) (time.Time, error) {
	return time.Unix(dbValue/1000, dbValue%1000*1000000).UTC(), nil
}

// TimeInt64ConvertToDatabaseValue converts time.Time to Unix timestamp in milliseconds (internal format expected by ObjectBox on a date field)
// NOTE - you lose precision - anything smaller then milliseconds is dropped
func TimeInt64ConvertToDatabaseValue(goValue time.Time) (int64, error) {
	var ms = int64(goValue.Nanosecond()) / 1000000
	return goValue.Unix()*1000 + ms, nil
}

// TimeTextConvertToEntityProperty uses time.Time.UnmarshalText() to decode RFC 3339 formatted string to time.Time.
func TimeTextConvertToEntityProperty(dbValue string) (goValue time.Time, err error) {
	err = goValue.UnmarshalText([]byte(dbValue))
	if err != nil {
		err = fmt.Errorf("error unmarshalling time %v: %v", dbValue, err)
	}
	return goValue, err
}

// TimeTextConvertToDatabaseValue uses time.Time.MarshalText() to encode time.Time into RFC 3339 formatted string.
func TimeTextConvertToDatabaseValue(goValue time.Time) (string, error) {
	bytes, err := goValue.MarshalText()
	if err != nil {
		err = fmt.Errorf("error marshalling time %v: %v", goValue, err)
	}
	return string(bytes), err
}

// TimeBinaryConvertToEntityProperty uses time.Time.UnmarshalBinary() to decode time.Time.
func TimeBinaryConvertToEntityProperty(dbValue []byte) (goValue time.Time, err error) {
	err = goValue.UnmarshalBinary(dbValue)
	if err != nil {
		err = fmt.Errorf("error unmarshalling time %v: %v", dbValue, err)
	}
	return goValue, err
}

// TimeBinaryConvertToDatabaseValue uses time.Time.MarshalBinary() to encode time.Time.
func TimeBinaryConvertToDatabaseValue(goValue time.Time) ([]byte, error) {
	bytes, err := goValue.MarshalBinary()
	if err != nil {
		err = fmt.Errorf("error marshalling time %v: %v", goValue, err)
	}
	return bytes, err
}
