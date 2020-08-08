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

/*
Generates objectbox related code for ObjectBox entities (Go structs)

It can be used by adding `//go:generate go run github.com/objectbox/objectbox-go/cmd/objectbox-gogen` comment inside a .go file
containing the struct that you want to persist and executing `go generate` in the module


Alternatively, you can run the command manually:

	objectbox-gogen [flags] {source-file}
		to generate the binding code

or

	objectbox-gogen clean {path}
		to remove the generated files instead of creating them - this removes *.obx.go and objectbox-model.go but keeps objectbox-model.json

path:
  * a source file path or a valid path pattern as accepted by the go tool (e.g. ./...)
  * if not given, the generator expects GOFILE environment variable to be set

Available flags:
  -byValue
    	getters should return a struct value (a copy) instead of a struct pointer
  -help
    	print this help
  -out string
    	output path for generated source files
  -persist string
    	path to the model information persistence file (JSON)
  -version
    	print the generator version info


To learn more about different configuration and annotations for entities, see docs at https://golang.objectbox.io/
*/
package main

import (
	"github.com/objectbox/objectbox-generator/cmd/objectbox-gogen"
)

func main() {
	gogen.Main()
}
