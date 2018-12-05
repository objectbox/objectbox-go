package object

// negative test, tag `uid` on an unknown property
type C struct {
	Id  uint64
	New string `uid`
}
