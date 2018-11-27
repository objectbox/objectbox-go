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
	"testing"
)

// NOTE overwriteExpected is used during development to update all ".expected" files with the generated content
// it's up to the developer to actually check whether the newly generated files are correct before commit
// NOTE - never commit this file with `overwriteExpected = true` as it means nothing is actually tested
var overwriteExpected = true

func TestAll(t *testing.T) {
	generateAllDirs(t, overwriteExpected)
}
