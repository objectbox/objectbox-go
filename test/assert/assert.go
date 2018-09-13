package assert

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

func EqString(t *testing.T, expected string, actual string) {
	if expected != actual {
		Failf(t, "Expected \""+expected+"\", but got \""+actual+"\"")
	}
}

func EqInt(t *testing.T, expected int, actual int) {
	if expected != actual {
		Failf(t, "Expected %v, but got %v", expected, actual)
	}
}

// Uses reflect.DeepEqual to test for equality
func Eq(t *testing.T, expected interface{}, actual interface{}) {
	if !reflect.DeepEqual(expected, actual) {
		Failf(t, "Expected %v, but got %v", expected, actual)
	}
}

func NoErr(t *testing.T, err error) {
	if err != nil {
		Failf(t, "Unexpected error occurred: %v", err)
	}
}

func Failf(t *testing.T, format string, args ...interface{}) {
	Fail(t, fmt.Sprintf(format, args...))
}

func Fail(t *testing.T, text string) {
	stackString := "Call stack:\n"
	for idx := 1; ; idx++ {
		_, file, line, ok := runtime.Caller(idx)
		if !ok {
			break
		}
		_, filename := filepath.Split(file)
		if filename == "assert.go" {
			continue
		}
		if filename == "testing.go" {
			break
		}
		stackString += fmt.Sprintf("%v:%v\n", filename, line)
	}
	if t != nil {
		t.Fatal(text, "\n", stackString)
	} else {
		fmt.Print(text, "\n", stackString)
	}
}
