package object

type B struct {
	Id uint64 `id`
	//Removed string `index`
	New int `index` // added at the same generator run as the previous one have been removed
}
