package assert

import (
	"fmt"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"testing"
)

func EqString(t *testing.T, expected string, actual string) {
	if expected != actual {
		debug.PrintStack()
		Failf(t, "Expected \""+expected+"\", but got \""+actual+"\"")
	}
}

func EqInt(t *testing.T, expected int, actual int) {
	if expected != actual {
		debug.PrintStack()
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
	stack_string := "Call stack:\n"
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
		stack_string += fmt.Sprintf("%v:%v\n", filename, line)
	}
	t.Fatal(text, "\n", stack_string)
}
