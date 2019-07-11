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

package generator

import (
	"bytes"
	"errors"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/objectbox/objectbox-go/internal/generator"
	"github.com/objectbox/objectbox-go/test/assert"
)

// generateAllDirs walks through the "data" and generates bindings for each subdirectory
// set overwriteExpected to TRUE to update all ".expected" files with the generated content
func generateAllDirs(t *testing.T, overwriteExpected bool) {
	var datadir = "testdata"
	folders, err := ioutil.ReadDir(datadir)
	assert.NoErr(t, err)

	for _, folder := range folders {
		if !folder.IsDir() {
			continue
		}

		var dir = filepath.Join(datadir, folder.Name())
		t.Run(folder.Name(), func(t *testing.T) {
			t.Parallel()
			generateOneDir(t, overwriteExpected, dir)
		})
	}
}

func generateOneDir(t *testing.T, overwriteExpected bool, dir string) {
	modelInfoFile := generator.ModelInfoFile(dir)
	modelInfoExpectedFile := modelInfoFile + ".expected"
	modelInfoInitialFile := modelInfoFile + ".initial"

	modelFile := generator.ModelFile(modelInfoFile)
	modelExpectedFile := modelFile + ".expected"

	// run the generation twice, first time with deleting old modelInfo
	for i := 0; i <= 1; i++ {
		if i == 0 {
			t.Logf("Testing %s without model info JSON", filepath.Base(dir))
			os.Remove(modelInfoFile)
		} else {
			t.Logf("Testing %s with previous model info JSON", filepath.Base(dir))
		}

		if fileExists(modelInfoInitialFile) {
			assert.NoErr(t, copyFile(modelInfoInitialFile, modelInfoFile))
		}

		generateAllFiles(t, overwriteExpected, dir, modelInfoFile)

		assertSameFile(t, modelInfoFile, modelInfoExpectedFile, overwriteExpected)
		assertSameFile(t, modelFile, modelExpectedFile, overwriteExpected)
	}
}

func assertSameFile(t *testing.T, file string, expectedFile string, overwriteExpected bool) {
	// if no file is expected
	if !fileExists(expectedFile) {
		// there can be no source file either
		if fileExists(file) {
			assert.Failf(t, "%s is missing but %s exists", expectedFile, file)
		}
		return
	}

	content, err := ioutil.ReadFile(file)
	assert.NoErr(t, err)

	if overwriteExpected {
		assert.NoErr(t, copyFile(file, expectedFile))
	}

	contentExpected, err := ioutil.ReadFile(expectedFile)
	assert.NoErr(t, err)

	if 0 != bytes.Compare(content, contentExpected) {
		assert.Failf(t, "generated file %s is not the same as %s", file, expectedFile)
	}
}

func generateAllFiles(t *testing.T, overwriteExpected bool, dir string, modelInfoFile string) {
	var modelFile = generator.ModelFile(modelInfoFile)

	// remove generated files during development (they might be syntactically wrong)
	if overwriteExpected {
		inputFiles, err := filepath.Glob(filepath.Join(dir, "*.obx.go"))
		assert.NoErr(t, err)

		for _, sourceFile := range inputFiles {
			assert.NoErr(t, os.Remove(sourceFile))
		}
	}

	// process all *.go files in the directory
	inputFiles, err := filepath.Glob(filepath.Join(dir, "*.go"))
	assert.NoErr(t, err)
	for _, sourceFile := range inputFiles {
		// skip generated files & "expected results" files
		if strings.HasSuffix(sourceFile, ".obx.go") ||
			strings.HasSuffix(sourceFile, ".skip.go") ||
			strings.HasSuffix(sourceFile, "expected") ||
			strings.HasSuffix(sourceFile, "initial") ||
			sourceFile == modelFile {
			continue
		}

		t.Logf("  %s", filepath.Base(sourceFile))

		err = generator.Process(sourceFile, getOptions(t, sourceFile, modelInfoFile))

		// handle negative test
		var shouldFail = strings.HasPrefix(filepath.Base(sourceFile), "!")
		if shouldFail {
			if err == nil {
				assert.Failf(t, "Unexpected PASS on a negative test %s", sourceFile)
			} else {
				var errPlatformIndependent = strings.Replace(err.Error(), "\\", "/", -1)
				assert.Eq(t, getExpectedError(t, sourceFile).Error(), errPlatformIndependent)
				continue
			}
		}

		assert.NoErr(t, err)

		var bindingFile = generator.BindingFile(sourceFile)
		var expectedFile = bindingFile + ".expected"
		assertSameFile(t, bindingFile, expectedFile, overwriteExpected)
	}
}

var generatorArgsRegexp = regexp.MustCompile("//go:generate objectbox-gogen (.+)[\n|\r]")

func getOptions(t *testing.T, sourceFile, modelInfoFile string) generator.Options {
	var options = generator.Options{
		ModelInfoFile: modelInfoFile,
		// NOTE zero seed for test-only - avoid changes caused by random numbers by fixing them to the same seed
		Rand: rand.New(rand.NewSource(0)),
	}

	source, err := ioutil.ReadFile(sourceFile)
	assert.NoErr(t, err)

	var match = generatorArgsRegexp.FindSubmatch(source)
	if len(match) > 1 {
		var args = argsToMap(string(match[1]))

		setArgs(t, args, &options)
	}

	return options
}

var expectedErrorRegexp = regexp.MustCompile("// *ERROR *=(.+)[\n|\r]")

func getExpectedError(t *testing.T, sourceFile string) error {
	source, err := ioutil.ReadFile(sourceFile)
	assert.NoErr(t, err)

	var match = expectedErrorRegexp.FindSubmatch(source)
	if len(match) > 1 {
		return errors.New(strings.TrimSpace(string(match[1])))
	}

	return nil
}

func setArgs(t *testing.T, args map[string]string, options *generator.Options) {
	for name, value := range args {
		_ = value // get rid of the compiler warning until we start using some options with values

		switch name {
		case "byValue":
			options.ByValue = true
		default:
			t.Fatalf("unknown option '%s'", name)
		}
	}
}

func argsToMap(args string) map[string]string {
	var result = map[string]string{}

	for _, arg := range strings.Split(strings.TrimSpace(args), "-") {
		arg = strings.TrimSpace(arg)

		if len(arg) == 0 {
			continue
		}

		var pair = strings.Split(arg, " ")
		if len(pair) == 1 {
			result[pair[0]] = ""
		} else {
			result[pair[0]] = pair[1]
		}
	}

	return result
}
