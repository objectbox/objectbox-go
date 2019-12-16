package object

/* ERROR:
can't merge binding model information: uid annotation value must not be empty on property Group, entity NegTaskRelId:
    [rename] apply the current UID 6745438398739480977
    [change/reset] apply a new UID 3959279844101328186
*/

type NegTaskRelId struct {
	Id    uint64
	Group uint64 `objectbox:"link:Group uid"`
}
