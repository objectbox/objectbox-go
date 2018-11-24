package object

type C struct {
	Id uint64 `id`
	//Removed string `index` // removed in one generator run
	New int `index` // added in another generator run (actual test)
}
