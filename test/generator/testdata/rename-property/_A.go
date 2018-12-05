package object

// negative test, tag `uid` will cause the build tool to print the UID of the property and fail
type A struct {
	Id  uint64
	Old string `uid`
}
