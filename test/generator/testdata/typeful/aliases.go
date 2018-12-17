package object

// Tests type aliases and definitions of named types

type sameFileAlias = string
type sameFileNamed string

type Aliases struct {
	Id           uint64
	SameFile     sameFileAlias
	SamePackage  samePackageAlias
	SameFile2    sameFileNamed
	SamePackage2 samePackageNamed
}
