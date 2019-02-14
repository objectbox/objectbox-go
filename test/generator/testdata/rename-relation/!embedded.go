package object

// ERROR = can't merge binding model information: uid annotation value must not be empty (model relation UID = 7845762441295307478) on relation Groups, entity NegTaskRelEmbedded

type NegTaskRelEmbedded struct {
	Id uint64
	WithGroupUidRequest
}
