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
	"time"

	"github.com/objectbox/objectbox-go/internal/generator"
)

func main() {
	sourceFile, options := getArgs()

	fmt.Printf("Generating ObjectBox bindings for %s", sourceFile)
	fmt.Println()

	// we need to do random seeding here instead of the internal/generator so that it can be easily testable
	rand.Seed(time.Now().UTC().UnixNano())

	err := generator.Process(sourceFile, options)
	stopOnError(err)
}

func stopOnError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
}

func showUsageAndExit() {
	flag.Usage()
	os.Exit(1)
}

func getArgs() (file string, options generator.Options) {
	flag.StringVar(&file, "source", "", "path to the source file containing structs to process")
	flag.StringVar(&options.ModelInfoFile, "persist", "", "path to the model information persistence file")
	flag.BoolVar(&options.ByValue, "byValue", false, "getters should return the struct value (copy) instead of struct pointer")
	flag.Parse()

	if len(file) == 0 {
		// if the command is run by go:generate some environment variables are set
		// https://golang.org/pkg/cmd/go/internal/generate/
		if gofile, exists := os.LookupEnv("GOFILE"); exists {
			file = gofile
		}

		if len(file) == 0 {
			showUsageAndExit()
		}
	}

	return
}
