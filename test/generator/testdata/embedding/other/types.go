package other

// used to test embedding types from other packages

type Trackable struct {
	Location     string
	private      int
	_AlsoPrivate float64
	privateType
}

type privateType struct {
	Settings string
}

type ForeignAlias = string
type ForeignNamed string
