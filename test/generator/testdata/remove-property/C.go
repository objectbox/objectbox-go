package object

type C struct {
	Id uint64 `objectbox:"id"`
	//Removed string // removed in one generator run
	New int // added in another generator run (actual test)
}
