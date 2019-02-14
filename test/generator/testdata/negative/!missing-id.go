package negative

// ERROR = can't prepare bindings for testdata/negative/!missing-id.go: id field is missing on entity MissingId - either annotate a field with `id` tag or use an uint64 field named 'Id/id/ID'

type MissingId struct {
	Text string
}
