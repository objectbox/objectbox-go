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

package objectbox

import (
	"time"
)

// waitUntil waits up to `timeout` until `fn()` returns true. Fails fast if `fn()` returns an error
// Returns:
// 		(true, nil) in case of a time out
// 		(false, nil) in case the `fn()` returned true => overall success
// 		(false, error) if the `fn()` returned an error

func waitUntil(timeout time.Duration, step time.Duration, fn func() (bool, error)) (timedOut bool, err error) {
	var endTime = time.After(timeout)
	tick := time.Tick(step)

	// Keep trying until we're timed out or got a result or got an error
	for {
		select {
		// Got a timeout! fail with a timeout error
		case <-endTime:
			return true, nil
		// Got a tick, we should check on doSomething()
		case <-tick:
			if ok, err := fn(); err != nil {
				return false, err
			} else if ok {
				return false, nil
			}
			// fn() didn't work yet, but it didn't fail, so let's try again
			// this will exit up to the for loop
		}
	}
}

