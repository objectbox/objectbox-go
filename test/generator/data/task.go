package object

//go:generate objectbox-bindings

type Event struct {
	// TODO check tag conventions if no value is used
	// https://golang.org/pkg/reflect/#StructTag
	Id       uint64 `id`
	Uid      string `index`
	Device   string `nameInDb:"device"`
	Date     uint64 `date`
	tempInfo string `transient`
	//Data []byte
}
