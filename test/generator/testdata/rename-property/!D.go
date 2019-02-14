package object

// ERROR = can't merge binding model information: property with Uid 5617773211005988520 not found

// completely new entity, already with an uid on a property
// this is quite unusual and indicates a migration from another DB
// or just a copy-paste, in any case, it needs to be handled gracefully
type D struct {
	Id  uint64 `id`
	New string `uid:"5617773211005988520"`
}
