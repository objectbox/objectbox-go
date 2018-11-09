package generator

import (
	"os"
	"testing"

	"github.com/objectbox/objectbox-go/internal/generator"
	"github.com/objectbox/objectbox-go/test/assert"
)

// TODO implement test similar to gofmt
// i. e. GLOB("data/*.input"), process & compare with "data/*.expected" files

func processAndTest(t *testing.T, sourceFile, bindingFile, modelInfoFile string, expectedSize int64) {
	var err error

	err = generator.Process(sourceFile, modelInfoFile)
	assert.NoErr(t, err)

	infoBinding, err := os.Stat(bindingFile)
	assert.NoErr(t, err)
	assert.Eq(t, expectedSize, infoBinding.Size())

	// check the permissions
	infoSource, err := os.Stat(sourceFile)
	assert.NoErr(t, err)
	assert.Eq(t, infoSource.Mode(), infoBinding.Mode())
}

func TestTask(t *testing.T) {
	var sourceFile = "data/task.go"
	var bindingFile = "data/taskbinding.go"
	var modelInfoFile = "data/objectbox-model-info.js"
	var expectedSize = int64(3681)

	// test when there's no binding file before
	os.Remove(bindingFile)
	//os.Remove(modelInfoFile)
	processAndTest(t, sourceFile, bindingFile, modelInfoFile, expectedSize)

	// test when the binding file already exists
	processAndTest(t, sourceFile, bindingFile, modelInfoFile, expectedSize)
}

func TestTypeful(t *testing.T) {
	var sourceFile = "data/typeful.go"
	var bindingFile = "data/typefulbinding.go"
	var modelInfoFile = "data/objectbox-model-info.js"
	var expectedSize = int64(5331)

	// test when there's no binding file before
	os.Remove(bindingFile)
	os.Remove(modelInfoFile)
	processAndTest(t, sourceFile, bindingFile, modelInfoFile, expectedSize)

	// test when the binding file already exists
	processAndTest(t, sourceFile, bindingFile, modelInfoFile, expectedSize)
}
