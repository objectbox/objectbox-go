package object

type B struct {
	Id uint64 `objectbox:"id"`
	//Removed string `objectbox:"index"`
	New int `objectbox:"index"` // added at the same generator run as the previous one have been removed
}
