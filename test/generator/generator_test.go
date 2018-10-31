package generator

import (
	"os"
	"testing"

	"github.com/objectbox/objectbox-go/internal/generator"
	"github.com/objectbox/objectbox-go/test/assert"
)

func processAndTest(t *testing.T, sourceFile, bindingFile string) {
	var err error

	err = generator.Process(sourceFile)
	assert.NoErr(t, err)

	infoBinding, err := os.Stat(bindingFile)
	assert.NoErr(t, err)
	assert.Eq(t, infoBinding.Size(), int64(3079))

	// check the permissions
	infoSource, err := os.Stat(sourceFile)
	assert.NoErr(t, err)
	assert.Eq(t, infoBinding.Mode(), infoSource.Mode())
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
