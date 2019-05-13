package object

// rename existing property
type B struct {
	Id  uint64 `objectbox:"id"`
	New string `objectbox:"uid:2669985732393126063"`
}
