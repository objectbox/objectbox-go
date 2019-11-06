package object

// ERROR = can't merge binding model information: property with Uid 5617773211005988520 not found in 'D'; property named 'New' not found in 'D'

// completely new entity, already with an uid on a property
// this is quite unusual and indicates a migration from another DB
// or just a copy-paste, in any case, it needs to be handled gracefully
type D struct {
	Id  uint64 `objectbox:"id"`
	New string `objectbox:"uid:5617773211005988520"`
}
