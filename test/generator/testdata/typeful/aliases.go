package object

import (
	"github.com/objectbox/objectbox-go/test/generator/testdata/embedding/other"
	ot "github.com/objectbox/objectbox-go/test/generator/testdata/embedding/other"
)

// Tests type aliases and definitions of named types

type sameFileAlias = string
type sameFileNamed string

type Aliases struct {
	Id            uint64
	SameFile      sameFileAlias
	SamePackage   samePackageAlias
	SameFile2     sameFileNamed
	SamePackage2  samePackageNamed
	OtherPackage  other.ForeignAlias
	OtherPackage2 ot.ForeignNamed
}
