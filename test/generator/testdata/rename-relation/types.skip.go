package object

type WithGroup struct {
	Group     *Group   `link`
	GroupsNew []*Group `uid:"6392442863481646880"`
}

type WithGroupUidRequest struct {
	Group  *Group   `link`
	Groups []*Group `uid`
}
