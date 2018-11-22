package generator

import (
	"testing"
)

// NOTE overwriteExpected is used during development to update all ".expected" files with the generated content
// it's up to the developer to actually check whether the newly generated files are correct before commit
// NOTE - never commit this file with `overwriteExpected = true` as it means nothing is actually tested
var overwriteExpected = false

func TestAll(t *testing.T) {
	generateAllDirs(t, overwriteExpected)
}
