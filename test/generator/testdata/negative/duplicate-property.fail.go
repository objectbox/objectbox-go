package negative

// ERROR = can't prepare bindings for testdata/negative/duplicate-property.fail.go: duplicate name (note that property names are case insensitive) on property text found in DuplicateProperty

type DuplicateProperty struct {
	Id   uint64 `objectbox:"id"`
	Text string `objectbox:"name:text"`
	text string
}
