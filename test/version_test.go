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

	versionGo := objectbox.VersionGo().String()
	if !format.MatchString(versionGo) {
		t.Errorf("ObjectBox-Go version %v doesn't match expected regexp %v", versionGo, format)
	}
	versionLib := objectbox.VersionLib().String()
	if !format.MatchString(versionGo) {
		t.Errorf("ObjectBox-C version %v doesn't match expected regexp %v", versionLib, format)
	}
	assert.Eq(t, true, strings.Contains(versionInfo, versionGo))
	assert.Eq(t, true, strings.Contains(versionInfo, versionLib))
}

func TestVersionLabel(t *testing.T) {
	var version = objectbox.Version{Major: 1, Minor: 2, Patch: 3, Label: "beta"}
	assert.Eq(t, version.String(), "1.2.3-beta")
	version.Label = ""
	assert.Eq(t, version.String(), "1.2.3")
}
