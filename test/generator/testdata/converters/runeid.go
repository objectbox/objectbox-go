package object

type RuneIdEntity struct {
	Id rune `objectbox:"id type:uint64 converter:runeId"`
}
