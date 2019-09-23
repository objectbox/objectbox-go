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

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/objectbox/objectbox-go/internal/generator"
)

func main() {
	path, clean, options := getArgs()

	var err error
	if clean {
		fmt.Printf("Removing ObjectBox bindings for %s\n", path)
		err = generator.Clean(path)
	} else {
		fmt.Printf("Generating ObjectBox bindings for %s\n", path)
		err = generator.Process(path, options)
	}

	stopOnError(err)
}

func stopOnError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
}

func showUsage() {
	fmt.Fprint(flag.CommandLine.Output(), `Usage:
	objectbox-gogen [flags] [path-pattern]
		to generate the binding code

or

	objectbox-gogen clean [path-pattern]
		to remove the generated files instead of creating them - this removes *.obx.go and objectbox-model.go but keeps objectbox-model.json

path-pattern:
  * a path or a valid path pattern as accepted by the go tool (e.g. ./...)
  * if not given, the generator expects GOFILE environment variable to be set

Available flags:
`)
	flag.PrintDefaults()
}

func showUsageAndExit() {
	showUsage()
	os.Exit(1)
}

func getArgs() (path string, clean bool, options generator.Options) {
	var printVersion bool
	var printHelp bool
	flag.Usage = showUsage
	flag.StringVar(&path, "source", "", "@deprecated, equivalent to passing the given source file path as as the path-pattern argument")
	flag.StringVar(&options.ModelInfoFile, "persist", "", "path to the model information persistence file")
	flag.BoolVar(&options.ByValue, "byValue", false, "getters should return a struct value (a copy) instead of a struct pointer")
	flag.BoolVar(&printVersion, "version", false, "print the generator version info")
	flag.BoolVar(&printHelp, "help", false, "print this help")
	flag.Parse()

	if printHelp {
		showUsage()
		os.Exit(0)
	}

	if printVersion {
		fmt.Println(fmt.Sprintf("ObjectBox Go binding code generator version: %d", generator.Version))
		os.Exit(0)
	}

	if len(path) != 0 {
		fmt.Println("'source' flag is deprecated and will be removed in the future - use a standard positional argument instead. See command help for more information.")
	}

	var argPath string

	if flag.NArg() == 2 {
		clean = true
		if flag.Arg(0) != "clean" {
			fmt.Printf("Unknown argument %s", flag.Arg(0))
			showUsageAndExit()
		}

		argPath = flag.Arg(1)

	} else if flag.NArg() == 1 {
		argPath = flag.Arg(0)
	} else if flag.NArg() != 0 {
		showUsageAndExit()
	}

	// if the path-pattern positional argument was given
	if len(argPath) > 0 {
		if len(path) == 0 {
			path = argPath
		} else if argPath != path {
			fmt.Printf("Path argument mismatch - given 'source' flag '%s' and the positional path argument '%s'\n", path, argPath)
			showUsageAndExit()
		}
	}

	if len(path) == 0 {
		// if the command is run by go:generate some environment variables are set
		// https://golang.org/pkg/cmd/go/internal/generate/
		if gofile, exists := os.LookupEnv("GOFILE"); exists {
			path = gofile
		}

		if len(path) == 0 {
			showUsageAndExit()
		}
	}

	return
}
