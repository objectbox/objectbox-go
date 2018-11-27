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

package main

import (
	"flag"
	"log"
	"os"
	"runtime"
	"runtime/debug"

	"github.com/objectbox/objectbox-go/test/performance/perf"
)

func main() {
	var dbName = flag.String("db", "db", "database directory")
	var count = flag.Int("count", 100000, "number of objects")
	var runs = flag.Int("runs", 30, "number of times the tests should be executed")
	flag.Parse()

	log.Printf("running the test %d times with %d objects", *runs, *count)

	// remove old database in case it already exists (and remove it after the test as well)
	os.RemoveAll(*dbName)
	defer os.RemoveAll(*dbName)

	// disable automatic garbage collector
	debug.SetGCPercent(-1)

	executor := perf.CreateExecutor(*dbName)
	defer executor.Close()

	inserts := executor.PrepareData(*count)

	for i := 0; i < *runs; i++ {
		executor.PutAll(inserts)
		items := executor.ReadAll(*count)
		executor.UpdateAll(items)
		executor.RemoveAll()

		log.Printf("%d/%d finished", i+1, *runs)

		// manually invoke GC out of benchmarked time
		runtime.GC()
		log.Printf("%d/%d garbage-collector executed", i+1, *runs)
	}

	executor.PrintTimes([]string{
		"PutAll",
		"ReadAll",
		"UpdateAll",
		"RemoveAll",
	})
}
