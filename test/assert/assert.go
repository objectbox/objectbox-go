package assert

import "strconv"

func EqString(expected string, actual string) {
	if expected != actual {
		panic("Expected \"" + expected + "\", but got \"" + actual + "\"")
	}
}

func EqInt(expected int, actual int) {
	if expected != actual {
		panic("Expected \"" + strconv.Itoa(expected) + "\", but got \"" + strconv.Itoa(actual) + "\"")
	}
}
