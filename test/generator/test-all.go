package generator

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/objectbox/objectbox-go/internal/generator"
	"github.com/objectbox/objectbox-go/test/assert"
)

// generateAllDirs walks through the "data" and generates bindings for each subdirectory
// set overwriteExpected to TRUE to update all ".expected" files with the generated content
func generateAllDirs(t *testing.T, overwriteExpected bool) {
	var datadir = "data"
	folders, err := ioutil.ReadDir(datadir)
	assert.NoErr(t, err)

	for _, folder := range folders {
		if !folder.IsDir() {
			continue
		}

		fmt.Println("Testing " + folder.Name())

		var dir = path.Join(datadir, folder.Name())

		modelInfoFile := path.Join(dir, "objectbox-model-info.json")
		modelInfoExpectedFile := modelInfoFile[0:len(modelInfoFile)-len(path.Ext(modelInfoFile))] + ".expected.json"
		modelInfoOriginalFile := modelInfoFile[0:len(modelInfoFile)-len(path.Ext(modelInfoFile))] + ".original.json"

		// run the generation twice, first time with deleting old modelInfo
		for i := 0; i <= 1; i++ {
			if i == 0 {
				os.Remove(modelInfoFile)
			}

			if fileExists(modelInfoOriginalFile) {
				assert.NoErr(t, copyFile(modelInfoOriginalFile, modelInfoFile))
			}

			generateAllFiles(t, overwriteExpected, dir, modelInfoFile)

			modelInfoFileContents, err := ioutil.ReadFile(modelInfoFile)
			assert.NoErr(t, err)

			if overwriteExpected {
				assert.NoErr(t, copyFile(modelInfoFile, modelInfoExpectedFile))
			}

			modelInfoFileExpectedContents, err := ioutil.ReadFile(modelInfoExpectedFile)
			assert.NoErr(t, err)

			if 0 != bytes.Compare(modelInfoFileContents, modelInfoFileExpectedContents) {
				assert.Failf(t, "Generated model info file %s is not the same as %s",
					modelInfoFile, modelInfoExpectedFile)
			}
		}
	}
}

func generateAllFiles(t *testing.T, overwriteExpected bool, dir string, modelInfoFile string) {
	// NOTE test-only - avoid changes caused by random numbers by fixing them to the same seed all the time
	rand.Seed(0)

	// process all *.go files in the directory
	inputFiles, err := filepath.Glob(path.Join(dir, "*.go"))
	assert.NoErr(t, err)
	for _, sourceFile := range inputFiles {
		// skip generated files & "expected results" files
		if strings.HasSuffix(sourceFile, "binding.go") || strings.HasSuffix(sourceFile, "expected.go") {
			continue
		}

		err = generator.Process(sourceFile, modelInfoFile)

		// handle negative test
		var shouldFail = strings.HasPrefix(path.Base(sourceFile), "!")
		if shouldFail {
			if err == nil {
				assert.Failf(t, "Unexpected PASS on a negative test %s", sourceFile)
			} else {
				return
			}
		}

		assert.NoErr(t, err)
		var bindingFile = generator.BindingFileName(sourceFile)
		var expectedFile = bindingFile[0:len(bindingFile)-3] + ".expected.go"

		bindingContents, err := ioutil.ReadFile(bindingFile)
		assert.NoErr(t, err)

		if overwriteExpected {
			assert.NoErr(t, copyFile(bindingFile, expectedFile))
		}

		expectedContents, err := ioutil.ReadFile(expectedFile)
		assert.NoErr(t, err)

		if 0 != bytes.Compare(bindingContents, expectedContents) {
			assert.Failf(t, "Generated binding file %s is not the same as %s", bindingFile, expectedFile)
		}
	}
}
