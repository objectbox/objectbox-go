package object

type WithGroup struct {
	Group  *Group `objectbox:"link"`
	Groups []*Group
}
