package negative

// ERROR = can't prepare bindings for testdata/negative/!duplicate-property.go: duplicate name (note that property names are case insensitive) on property text, entity DuplicateProperty

type DuplicateProperty struct {
	Id   uint64 `id`
	Text string `nameInDb:"text"`
	text string
}
