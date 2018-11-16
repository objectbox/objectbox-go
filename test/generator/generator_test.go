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

// TODO implement test similar to gofmt
// i. e. GLOB("data/*.input"), process & compare with "data/*.expected" files
//
func TestAll(t *testing.T) {
	var datadir = "data"
	folders, err := ioutil.ReadDir(datadir)
	assert.NoErr(t, err)

	for _, folder := range folders {
		if !folder.IsDir() {
			continue
		}

		fmt.Println("Testing " + folder.Name())

		// NOTE test-only - avoid changes caused by random numbers by fixing them to the same seed all the time
		rand.Seed(0)

		var dir = path.Join(datadir, folder.Name())

		modelInfoFile := path.Join(dir, "objectbox-model-info.json")
		modelInfoExpectedFile := modelInfoFile[0:len(modelInfoFile)-len(path.Ext(modelInfoFile))] + ".expected.json"

		// run the generation twice, first time with deleting old modelInfo
		for i := 0; i <= 1; i++ {
			if i == 0 {
				os.Remove(modelInfoFile)
			}

			generateAllFiles(t, dir, modelInfoFile)

			modelInfoFileContents, err := ioutil.ReadFile(modelInfoFile)
			assert.NoErr(t, err)

			modelInfoFileExpectedContents, err := ioutil.ReadFile(modelInfoExpectedFile)
			assert.NoErr(t, err)

			if 0 != bytes.Compare(modelInfoFileContents, modelInfoFileExpectedContents) {
				assert.Failf(t, "Generated model info file %s is not the same as %s",
					modelInfoFile, modelInfoExpectedFile)
			}
		}
	}

}

func generateAllFiles(t *testing.T, dir string, modelInfoFile string) {
	// process all *.go files in the directory
	inputFiles, err := filepath.Glob(path.Join(dir, "*.go"))
	assert.NoErr(t, err)
	for _, sourceFile := range inputFiles {
		// skip generated files & "expected results" files
		if strings.HasSuffix(sourceFile, "binding.go") || strings.HasSuffix(sourceFile, "expected.go") {
			continue
		}

		err = generator.Process(sourceFile, modelInfoFile)
		assert.NoErr(t, err)

		var bindingFile = generator.BindingFileName(sourceFile)
		var expectedFile = bindingFile[0:len(bindingFile)-3] + ".expected.go"

		bindingContents, err := ioutil.ReadFile(bindingFile)
		assert.NoErr(t, err)

		expectedContents, err := ioutil.ReadFile(expectedFile)
		assert.NoErr(t, err)

		if 0 != bytes.Compare(bindingContents, expectedContents) {
			assert.Failf(t, "Generated binding file %s is not the same as %s", bindingFile, expectedFile)
		}
	}
}
