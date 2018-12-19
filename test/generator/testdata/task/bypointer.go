package object

type Task struct {
	Id       uint64 `id`
	Uid      string `unique`
	Text     string `nameInDb:"text"`
	Date     uint64 `date json:"date"`
	tempInfo string `transient`
	GroupId  uint64 `link:"Group"`
}

type Group struct {
	Id uint64 `id`
}
