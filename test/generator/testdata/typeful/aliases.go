package object

// Tests type aliases

type sameFileAlias = string

type Aliases struct {
	Id          uint64
	SameFile    sameFileAlias
	SamePackage samePackageAlias
}
