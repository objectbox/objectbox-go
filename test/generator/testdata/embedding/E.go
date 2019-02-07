package object

import "github.com/objectbox/objectbox-go/test/generator/testdata/embedding/other"

type E struct {
	other.Trackable `inline`
	id              uint64
	other.ForeignAlias
	other.ForeignNamed
}
