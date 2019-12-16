package object

// let's create a dependency cycle by using the type that's not yet created
type _ = a_EntityInfo

type samePackageAlias = int
type samePackageNamed int
