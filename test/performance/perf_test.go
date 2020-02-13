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

package main

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestPerformanceSimple(t *testing.T) {
	var count = 100000

	if testing.Short() {
		count = 1000
	}

	log.Printf("running the test with %d objects", count)

	// Test in a temporary directory - if tested by an end user, the repo is read-only.
	tempDir, err := ioutil.TempDir("", "objectbox-test")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Fatal(err)
		}
	}()

	executor := createExecutor(tempDir)
	defer executor.close()

	inserts := executor.prepareData(count)

	executor.putAsync(inserts)
	executor.removeAll()

	executor.putMany(inserts)

	items := executor.readAll(count)
	executor.changeValues(items)
	executor.updateAll(items)

	executor.printTimes([]string{})
}
