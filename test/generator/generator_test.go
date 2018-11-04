package generator

import (
	"os"
	"testing"

	"github.com/objectbox/objectbox-go/internal/generator"
	"github.com/objectbox/objectbox-go/test/assert"
)

// TODO implement test similar to gofmt
// i. e. GLOB("data/*.input"), process & compare with "data/*.expected" files

func processAndTest(t *testing.T, sourceFile, bindingFile string, expectedSize int64) {
	var err error

	err = generator.Process(sourceFile)
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
	var expectedSize = int64(3518)

	// test when there's no binding file before
	os.Remove(bindingFile)
	processAndTest(t, sourceFile, bindingFile, expectedSize)

	// test when the binding file already exists
	processAndTest(t, sourceFile, bindingFile, expectedSize)
}

func TestTypeful(t *testing.T) {
	var sourceFile = "data/typeful.go"
	var bindingFile = "data/typefulbinding.go"
	var expectedSize = int64(4919)

	// test when there's no binding file before
	os.Remove(bindingFile)
	processAndTest(t, sourceFile, bindingFile, expectedSize)

	// test when the binding file already exists
	processAndTest(t, sourceFile, bindingFile, expectedSize)
}
