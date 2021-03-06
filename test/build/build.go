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

package build

import (
	"os/exec"
)

// Package builds a single Go package/directory, running `go build path`
func Package(path string) (stdOut []byte, stdErr []byte, err error) {
	var cmd = exec.Command("go", "build")
	cmd.Dir = path
	stdOut, err = cmd.Output()
	if ee, ok := err.(*exec.ExitError); ok {
		stdErr = ee.Stderr
	}
	return
}
