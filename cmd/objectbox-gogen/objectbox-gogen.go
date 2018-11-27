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

package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"path"
	"time"

	"github.com/objectbox/objectbox-go/internal/generator"
)

func main() {
	file, _, modelFile := getArgs()

	fmt.Printf("Generating ObjectBox bindings for %s", file)
	fmt.Println()

	// we need to do random seeding here instead of the internal/generator so that it can be easily testable
	rand.Seed(time.Now().UTC().UnixNano())

	err := generator.Process(file, modelFile)
	stopOnError(err)
}

func stopOnError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
}

func getArgs() (file string, line uint, modelInfoFile string) {
	var hasAll = true
	line = 0

	file = *flag.String("source", "", "path to the source file containing structs to process")

	if len(file) == 0 {
		// if the command is run by go:generate some environment variables are set
		// https://golang.org/pkg/cmd/go/internal/generate/
		if gofile, exists := os.LookupEnv("GOFILE"); exists {
			file = gofile
		}

		if len(file) == 0 {
			hasAll = false
		}
	}

	modelInfoFile = *flag.String("persist", "", "path to the model information persistence file")

	if len(modelInfoFile) == 0 {
		modelInfoFile = generator.ModelInfoFile(path.Dir(file))
	}

	if !hasAll {
		flag.Usage()
		os.Exit(1)
	}
	return
}
