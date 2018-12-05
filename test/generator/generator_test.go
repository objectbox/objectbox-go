/*
 * Copyright 2018 ObjectBox Ltd. All rights reserved.
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

package generator

import (
	"flag"
	"testing"
)

// used during development of generator to overwrite the "golden" files
var overwriteExpected = flag.Bool("update", false,
	"Update all '.expected' files with the generated content. "+
		"It's up to the developer to actually check before committing whether the newly generated files are correct")

// used during development of generator to test a single directory
var target = flag.String("target", "", "Specify target subdirectory of testdata to generate")

func TestGenerator(t *testing.T) {
	if *target == "" {
		generateAllDirs(t, *overwriteExpected)
	} else {
		generateOneDir(t, *overwriteExpected, "testdata/"+*target)
	}
}
