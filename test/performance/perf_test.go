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
	"log"
	"os"
	"testing"

	"github.com/objectbox/objectbox-go/test/performance/perf"
)

func TestPerformanceSimple(t *testing.T) {
	var dbName = "db"
	var count = 100000

	if testing.Short() {
		count = 1000
	}

	log.Printf("running the test with %d objects", count)

	// remove old database in case it already exists (and remove it after the test as well)
	os.RemoveAll(dbName)
	defer os.RemoveAll(dbName)

	executor := perf.CreateExecutor(dbName)
	defer executor.Close()

	inserts := executor.PrepareData(count)

	executor.PutAsync(inserts)
	executor.RemoveAll()

	executor.PutMany(inserts)

	items := executor.ReadAll(count)
	executor.ChangeValues(items)
	executor.UpdateAll(items)

	executor.PrintTimes([]string{})
}
