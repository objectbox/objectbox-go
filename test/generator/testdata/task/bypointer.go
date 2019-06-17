package object

type Task struct {
	Id       uint64 `objectbox:"id"`
	Uid      string `objectbox:"UNIQUE"` // let's verify it's case insensitive as well
	Text     string `objectbox:"name:text"`
	Date     uint64 `objectbox:"date" json:"date"`
	tempInfo string `objectbox:"-"`
	GroupId  uint64 `objectbox:"link:Group"`
}

type Group struct {
	Id uint64 `objectbox:"id"`
}
