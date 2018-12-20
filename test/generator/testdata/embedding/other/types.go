package other

// used to test embedding types from other packages

type Trackable struct {
	Location    string
	unavailable int // not available outside of this package
}
