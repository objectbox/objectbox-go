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

package objectbox_test

import (
	"errors"
	"fmt"
	"runtime"
	"testing"
)

// 500000000	         3.93 ns/op
func BenchmarkLockOsThread(b *testing.B) {
	for n := 0; n < b.N; n++ {
		runtime.LockOSThread()
		runtime.UnlockOSThread()
	}
}

// prevent compiler optimizing out an unused return value
var globalErr error

// 100000000	        27.1 ns/op
func BenchmarkErrorReturn(b *testing.B) {
	var recovery = func() error {
		return errors.New("oh")
	}
	for n := 0; n < b.N; n++ {
		globalErr = recovery()
	}
}

// Panic/Failing         	20000000	        88.2 ns/op
// Panic/Successful      	50000000	        36.1 ns/op
// Error/Failing         	50000000	        31.0 ns/op
// Error/Successful      	1000000000	         2.14 ns/op
func BenchmarkErrorVsPanicRecover(b *testing.B) {
	// Using a single function with boolean switches doesn't hurt the benchmarks.
	// Originally I had multiple functions and the results and the results are about the same
	var withPanic = func(shouldFail bool) (err error) {
		defer func() {
			if r := recover(); r != nil {
				switch x := r.(type) {
				case string:
					err = errors.New(x)
				case error:
					err = x
				default:
					err = fmt.Errorf("%v", r)
				}
			}
		}()

		if shouldFail {
			panic("oh")
		}

		return nil
	}
	var withError = func(shouldFail bool) error {
		if shouldFail {
			return errors.New("oh")
		}
		return nil
	}

	b.Run("Panic/Failing", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			globalErr = withPanic(true)
		}
	})

	b.Run("Panic/Successful", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			globalErr = withPanic(false)
		}
	})

	b.Run("Error/Failing", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			globalErr = withError(true)
		}
	})

	b.Run("Error/Successful", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			globalErr = withError(false)
		}
	})
}
