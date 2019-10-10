package negative

// ERROR = can't prepare bindings for testdata/negative/missing-id.fail.go: id field is missing on entity MissingId - either annotate a field with `objectbox:"id"` tag or use an (u)int64 field named 'Id/id/ID'

type MissingId struct {
	Text string
}
