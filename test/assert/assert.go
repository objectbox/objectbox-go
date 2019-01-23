/*
 * Copyright 2018 ObjectBox Ltd. All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package assert

import (
	"fmt"
	"path/filepath"
	"reflect"
	"regexp"
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
	if expected == nil && actual == nil {
		return
	}
	if !reflect.DeepEqual(expected, actual) {
		Failf(t, "Expected %v, but got %v", expected, actual)
	}
}

// Uses reflect.DeepEqual to test for equality
func NotEq(t *testing.T, notThisValue interface{}, actual interface{}) {
	if reflect.DeepEqual(notThisValue, actual) {
		Failf(t, "Expected a value other than %v", notThisValue)
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

// mustPanic ensures that the caller's context will panic and that the panic will match the given regular expression
//   func() {
//   	defer mustPanic(t, regexp.MustCompile("+*"))
//		panic("some text")
//   }
func MustPanic(t *testing.T, match *regexp.Regexp) {
	if r := recover(); r != nil {
		// convert panic result to string
		var str string
		switch x := r.(type) {
		case string:
			str = x
		case error:
			str = x.Error()
		default:
			Failf(t, "unknown panic result '%v' for an expected panic: %s", r, match.String())
		}

		if !match.MatchString(str) {
			Failf(t, "expected panic '%s' but got '%s'", match.String(), str)
		}
	} else {
		Failf(t, "expected panic hasn't occurred: %s", match.String())
	}
}
