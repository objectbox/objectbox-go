/*
 * Copyright 2018-2021 ObjectBox Ltd. All rights reserved.
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

package objectbox_test

import (
	"regexp"
	"strings"
	"testing"

	"github.com/objectbox/objectbox-go/objectbox"
	"github.com/objectbox/objectbox-go/test/assert"
)

func TestVersion(t *testing.T) {
	versionInfo := objectbox.VersionInfo()
	t.Log(versionInfo)

	var format = regexp.MustCompile(`^[0-9]+\.[0-9]+\.[0-9]+(-[a-z]+\.[0-9]+)?$`)

	versionGo := objectbox.VersionGo()
	versionGoString := versionGo.String()
	if !format.MatchString(versionGoString) {
		t.Errorf("ObjectBox-Go version %v doesn't match expected regexp %v", versionGoString, format)
	}
	versionGoInt := versionGo.Major*10000 + versionGo.Minor*100 + versionGo.Patch
	assert.True(t, versionGoInt >= 10600) // Update with new releases (won't fail if forgotten)
	assert.True(t, versionGoInt < 20000)  // Future next major release

	versionLib := objectbox.VersionLib()
	versionLibString := versionLib.String()
	if !format.MatchString(versionGoString) {
		t.Errorf("ObjectBox-C version %v doesn't match expected regexp %v", versionLibString, format)
	}
	versionLibInt := versionLib.Major*10000 + versionLib.Minor*100 + versionLib.Patch
	assert.True(t, versionLibInt >= 1500) // Update with new releases (won't fail if forgotten)
	assert.True(t, versionLibInt < 10000) // Future next major release

	assert.Eq(t, true, strings.Contains(versionInfo, versionGoString))
	assert.Eq(t, true, strings.Contains(versionInfo, versionLibString))
}

func TestVersionLabel(t *testing.T) {
	var version = objectbox.Version{Major: 1, Minor: 2, Patch: 3, Label: "beta"}
	assert.Eq(t, version.String(), "1.2.3-beta")
	version.Label = ""
	assert.Eq(t, version.String(), "1.2.3")
}
