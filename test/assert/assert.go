/*
 * Copyright 2019 ObjectBox Ltd. All rights reserved.
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

// True asserts that the given value is a boolean true
func True(t *testing.T, actual interface{}) {
	if actual != true {
		Failf(t, "Expected: true\nReceived: %v", actual)
	}
}

// Eq asserts the values are equal. Uses reflect.DeepEqual to test for equality
func Eq(t *testing.T, expected interface{}, actual interface{}) {
	if expected == nil && actual == nil {
		return
	}
	if !reflect.DeepEqual(expected, actual) {
		Failf(t, "Values are not equal\nExpected: %v\nReceived: %v", expected, actual)
	}
}

// EqItems asserts the given slices have the same elements (regardless of their order)
func EqItems(t *testing.T, expected interface{}, actual interface{}) {
	var exp = reflect.ValueOf(expected)
	var act = reflect.ValueOf(actual)

	if exp.Type() != act.Type() {
		Failf(t, "Types are not equal\nExpected: %v\nReceived: %v", exp.Type(), act.Type())
	}

	if exp.Len() != act.Len() {
		Failf(t, "Lengths are not equal\nExpected: %v (%d elements)\nReceived: %v (%d elements)", exp, exp.Len(), act, act.Len())
	}

	if exp.Len() == 0 {
		return
	}

	// make a map[elem-type]int = number of occurrences of each element
	// we use reflection to create a dynamically typed map
	var keyType = exp.Index(0).Type()
	var valueType = reflect.TypeOf(int(0))
	var mapType = reflect.MapOf(keyType, valueType)
	merged := reflect.MakeMapWithSize(mapType, exp.Len())

	// count the number of expected occurrences
	for i := 0; i < exp.Len(); i++ {
		var existing = merged.MapIndex(exp.Index(i))
		if existing.IsValid() {
			merged.SetMapIndex(exp.Index(i), reflect.ValueOf(int(existing.Int())+1)) // increase by one
		} else {
			merged.SetMapIndex(exp.Index(i), reflect.ValueOf(int(1)))
		}
	}

	// count the number of actual occurrences
	for i := 0; i < act.Len(); i++ {
		var existing = merged.MapIndex(act.Index(i))
		if !existing.IsValid() {
			Failf(t, "Unexpected item %v found in %v, expecting %v", act.Index(i), act, exp)
		}

		merged.SetMapIndex(act.Index(i), reflect.ValueOf(int(existing.Int())-1)) // decrease by one
	}

	// check if all of the expected where actually found
	for _, k := range merged.MapKeys() {
		var existing = merged.MapIndex(k)
		if existing.Int() != 0 {
			Failf(t, "Expected %v more of item %v", existing.Int(), k)
		}
	}
}

// NotEq asserts the given values are not equal. Uses reflect.DeepEqual to test for equality
func NotEq(t *testing.T, notThisValue interface{}, actual interface{}) {
	if reflect.DeepEqual(notThisValue, actual) {
		Failf(t, "Expected a value other than %v", notThisValue)
	}
}

// NoErr asserts the error is nil
func NoErr(t *testing.T, err error) {
	if err != nil {
		Failf(t, "Unexpected error occurred: %v", err)
	}
}

// Failf fails immediately
func Failf(t *testing.T, format string, args ...interface{}) {
	Fail(t, fmt.Sprintf(format, args...))
}

// Fail fails immediately
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

// MustPanic ensures that the caller's context will panic and that the panic will match the given regular expression
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
			Failf(t, "Unknown panic result '%v' for an expected panic: %s", r, match.String())
		}

		if !match.MatchString(str) {
			Failf(t, "Errors are not equal\nExpected: panic '%s'\nReceived: '%s'", match.String(), str)
		}
	} else {
		Failf(t, "Expected panic hasn't occurred: %s", match.String())
	}
}
