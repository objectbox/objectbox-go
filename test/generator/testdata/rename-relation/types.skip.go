package object

type WithGroup struct {
	Group     *Group   `objectbox:"link"`
	GroupsNew []*Group `objectbox:"uid:6392442863481646880"`
}

type WithGroupUidRequest struct {
	Group  *Group   `objectbox:"link"`
	Groups []*Group `objectbox:"uid"`
}
