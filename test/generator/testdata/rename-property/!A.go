package object

// ERROR = can't merge binding model information: uid annotation value must not be empty (model property UID = 3390393562759376202) on property Old, entity A

// negative test, tag `objectbox:"uid"` will cause the build tool to print the UID of the property and fail
type A struct {
	Id  uint64
	Old string `objectbox:"uid"`
}
