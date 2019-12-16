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

/*
#include <stdlib.h>
#include "objectbox.h"
*/
import "C"
import "fmt"

// Version represents a semantic-version If you depend on a certain version of ObjectBox, you can check using this struct.
// See also VersionGo() and VersionLib().
type Version struct {
	Major int
	Minor int
	Patch int
	Label string
}

func (v Version) String() string {
	versionString := fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
	if len(v.Label) > 0 {
		versionString += "-" + v.Label
	}
	return versionString
}

// VersionGo returns the Version of the ObjectBox-Go binding
func VersionGo() Version {
	// for label, use `beta.0` format, increasing the counter for each subsequent release
	return Version{1, 1, 0, ""}
}

// VersionLib returns the Version of the dynamic linked ObjectBox library
func VersionLib() Version {
	var major C.int
	var minor C.int
	var patch C.int
	C.obx_version(&major, &minor, &patch)
	return Version{int(major), int(minor), int(patch), ""}
}

// VersionInfo returns a printable version string
func VersionInfo() string {
	return "ObjectBox Go version " + VersionGo().String() + " using dynamic library version " + VersionLib().String()
}
