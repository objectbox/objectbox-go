package assert

import (
	"bufio"
	"fmt"
	"os"
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
	firstLocation := ""
	stack_count := 0
	stack_string := ""
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
		location := fmt.Sprintf("%v:%v", filename, line)
		if firstLocation == "" {
			firstLocation = location
		}
		stack_string += fmt.Sprintf("    %v\n", location)
		stack_count++
	}
	if stack_count > 1 {
		fmt.Println("Full stack for [" + text + "]:\n" + stack_string)

	}

	bufio.NewWriter(os.Stdout).Flush()
	t.Fatal(text, "\nat", firstLocation)
}
