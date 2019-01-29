package object

import "github.com/objectbox/objectbox-go/test/generator/testdata/embedding/other"

type E struct {
	other.Trackable
	id uint64
	other.ForeignAlias
	other.ForeignNamed
}
