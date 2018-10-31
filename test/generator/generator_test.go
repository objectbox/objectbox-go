package generator

import (
	"os"
	"testing"

	"github.com/objectbox/objectbox-go/internal/generator"
	"github.com/objectbox/objectbox-go/test/assert"
)

// TODO implement test similar to gofmt
// i. e. GLOB("data/*.input"), process & compare with "data/*.expected" files

func processAndTest(t *testing.T, sourceFile, bindingFile string) {
	var err error

	err = generator.Process(sourceFile)
	assert.NoErr(t, err)

	infoBinding, err := os.Stat(bindingFile)
	assert.NoErr(t, err)
	assert.Eq(t, int64(3079), infoBinding.Size())

	// check the permissions
	infoSource, err := os.Stat(sourceFile)
	assert.NoErr(t, err)
	assert.Eq(t, infoSource.Mode(), infoBinding.Mode())
}

func TestGeneratorSimple(t *testing.T) {
	var sourceFile = "data/task.go"
	var bindingFile = "data/taskbinding.go"

	// test when there's no binding file before
	os.Remove(bindingFile)
	processAndTest(t, sourceFile, bindingFile)

	// test when the binding file already exists
	processAndTest(t, sourceFile, bindingFile)
}
