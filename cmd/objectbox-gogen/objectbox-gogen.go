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
