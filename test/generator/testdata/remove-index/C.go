package object

type C struct {
	Id uint64 `objectbox:"id"`
	//Removed string `objectbox:"index"` // removed in one generator run
	New int `objectbox:"index"` // added in another generator run (actual test)
}
