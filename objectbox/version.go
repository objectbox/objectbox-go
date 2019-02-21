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
#cgo LDFLAGS: -lobjectbox
#include <stdlib.h>
#include "objectbox.h"
*/
import "C"
import "fmt"

// If you depend on a certain version of ObjectBox, you can check using this struct.
// See also VersionGo() and VersionLib().
type Version struct {
	Major int
	Minor int
	Patch int
}

func (v Version) String() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}

// ObjectBox Version Go
func VersionGo() Version {
	return Version{0, 9, 0}
}

// Version of the dynamic linked ObjectBox library
func VersionLib() Version {
	var major C.int
	var minor C.int
	var patch C.int
	C.obx_version(&major, &minor, &patch)
	return Version{int(major), int(minor), int(patch)}
}

// A printable version string
func VersionInfo() string {
	return "ObjectBox Go version " + VersionGo().String() + " using dynamic library version " + VersionLib().String()
}
