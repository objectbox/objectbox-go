/*
 * Copyright 2018-2021 ObjectBox Ltd. All rights reserved.
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
	"fmt"
	"github.com/objectbox/objectbox-go/objectbox"
	"github.com/objectbox/objectbox-go/test/performance/perf"
	"path/filepath"
	"runtime"
	"time"
)

type executor struct {
	ob    *objectbox.ObjectBox
	box   *perf.EntityBox
	times map[string][]time.Duration // arrays of runtimes indexed by function name
}

func createExecutor(dbName string) *executor {
	result := &executor{
		times: map[string][]time.Duration{},
	}
	result.initObjectBox(dbName)
	return result
}

func (exec *executor) initObjectBox(dbName string) {
	defer exec.trackTime(time.Now())

	objectBox, err := objectbox.NewBuilder().Directory(dbName).Model(perf.ObjectBoxModel()).Build()
	if err != nil {
		panic(err)
	}

	exec.ob = objectBox
	exec.box = perf.BoxForEntity(objectBox)
}

func (exec *executor) close() {
	defer exec.trackTime(time.Now())

	exec.ob.Close()
}

func (exec *executor) trackTime(start time.Time) {
	elapsed := time.Since(start)

	pc, _, _, _ := runtime.Caller(1)
	fun := filepath.Ext(runtime.FuncForPC(pc).Name())[1:]
	exec.times[fun] = append(exec.times[fun], elapsed)
}

func (exec *executor) printTimes(functions []string) {
	// print the whole data as a table
	fmt.Println("Function\tRuns\tAverage ms\tAll times")

	if len(functions) == 0 {
		for fun := range exec.times {
			functions = append(functions, fun)
		}
	}

	for _, fun := range functions {
		times := exec.times[fun]

		sum := int64(0)
		for _, duration := range times {
			sum += duration.Nanoseconds()
		}
		fmt.Printf("%s\t%d\t%f", fun, len(times), float64(sum/int64(len(times)))/1000000)

		for _, duration := range times {
			fmt.Printf("\t%f", float64(duration.Nanoseconds())/1000000)
		}
		fmt.Println()
	}
}

func (exec *executor) removeAll() {
	defer exec.trackTime(time.Now())
	err := exec.box.RemoveAll()
	if err != nil {
		panic(err)
	}
}

func (exec *executor) prepareData(count int) []*perf.Entity {
	defer exec.trackTime(time.Now())

	var result = make([]*perf.Entity, count)
	for i := 0; i < count; i++ {
		result[i] = &perf.Entity{
			String:  fmt.Sprintf("Entity no. %d", i),
			Float64: float64(i),
			Int32:   int32(i),
			Int64:   int64(i),
		}
	}

	return result
}

func (exec *executor) putAsync(items []*perf.Entity) {
	defer exec.trackTime(time.Now())

	const retries = 20

	var putErr error
	for _, item := range items {
		for i := 0; i < retries; i++ {
			if _, putErr = exec.box.PutAsync(item); putErr != nil {
				// before each retry we sleep for a little more
				time.Sleep(time.Duration(i+1) * time.Second)
			} else {
				break
			}
		}

		// if retrying failed, stop completely
		if putErr != nil {
			break
		}
	}

	if err := exec.ob.AwaitAsyncCompletion(); err != nil {
		panic(err)
	}

	// if retrying failed
	if putErr != nil {
		panic(putErr)
	}
}

func (exec *executor) putMany(items []*perf.Entity) {
	defer exec.trackTime(time.Now())

	if _, err := exec.box.PutMany(items); err != nil {
		panic(err)
	}
}

func (exec *executor) readAll(count int) []*perf.Entity {
	defer exec.trackTime(time.Now())

	if items, err := exec.box.GetAll(); err != nil {
		panic(err)
	} else if len(items) != count {
		panic("invalid number of objects read")
	} else {
		return items
	}
}

func (exec *executor) changeValues(items []*perf.Entity) {
	defer exec.trackTime(time.Now())

	count := len(items)
	for i := 0; i < count; i++ {
		items[i].Int64 = items[i].Int64 * 2
	}
}

func (exec *executor) updateAll(items []*perf.Entity) {
	defer exec.trackTime(time.Now())

	if _, err := exec.box.PutMany(items); err != nil {
		panic(err)
	}
}
