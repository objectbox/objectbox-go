package generator

import (
	"testing"

	"github.com/objectbox/objectbox-go/internal/generator"
	"github.com/objectbox/objectbox-go/test/assert"
)

func TestGenerator(t *testing.T) {
	var err error

	err = generator.Process("data/task.go")
	assert.NoErr(t, err)
}
