package object

type Task struct {
	Id       uint64 `objectbox:"id"`
	Uid      string `objectbox:"unique"`
	Text     string `objectbox:"name:text"`
	Date     uint64 `objectbox:"date" json:"date"`
	tempInfo string `objectbox:"transient"`
	GroupId  uint64 `objectbox:"link:Group"`
}

type Group struct {
	Id uint64 `objectbox:"id"`
}
