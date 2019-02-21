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
)

// implements "StringIdConvert" property value converter
func StringIdConvertToEntityProperty(dbValue uint64) (goValue string) {
	return strconv.FormatUint(dbValue, 10)
}

// implements "StringIdConvert" property value converter
func StringIdConvertToDatabaseValue(goValue string) uint64 {
	// in case the object was initialized by the user without setting the ID explicitly
	if goValue == "" {
		return 0
	}

	if id, err := strconv.ParseUint(goValue, 10, 64); err != nil {
		panic(fmt.Errorf("error parsing numeric ID represented as string: %s", err))
	} else {
		return uint64(id)
	}
}
