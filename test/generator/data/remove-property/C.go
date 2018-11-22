package object

type C struct {
	Id uint64 `id`
	//Removed string // removed in one generator run
	New int // added in another generator run (actual test)
}
