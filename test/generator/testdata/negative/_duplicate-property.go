package negative

type DuplicateProperty struct {
	Id   uint64 `id`
	Text string `nameInDb:"text"`
	text string
}
