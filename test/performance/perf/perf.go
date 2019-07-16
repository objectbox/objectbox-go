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

package perf

import (
	"fmt"
	"path/filepath"
	"runtime"
	"time"

	"github.com/objectbox/objectbox-go/objectbox"
)

type Executor struct {
	ob    *objectbox.ObjectBox
	box   *EntityBox
	times map[string][]time.Duration // arrays of runtimes indexed by function name
}

func CreateExecutor(dbName string) *Executor {
	result := &Executor{
		times: map[string][]time.Duration{},
	}
	result.initObjectBox(dbName)
	return result
}

func (perf *Executor) initObjectBox(dbName string) {
	defer perf.trackTime(time.Now())

	objectBox, err := objectbox.NewBuilder().Directory(dbName).Model(ObjectBoxModel()).Build()
	if err != nil {
		panic(err)
	}

	perf.ob = objectBox
	perf.box = BoxForEntity(objectBox)
}

func (perf *Executor) Close() {
	defer perf.trackTime(time.Now())

	perf.ob.Close()
}

func (perf *Executor) trackTime(start time.Time) {
	elapsed := time.Since(start)

	pc, _, _, _ := runtime.Caller(1)
	fun := filepath.Ext(runtime.FuncForPC(pc).Name())[1:]
	perf.times[fun] = append(perf.times[fun], elapsed)
}

func (perf *Executor) PrintTimes(functions []string) {
	// print the whole data as a table
	fmt.Println("Function\tRuns\tAverage ms\tAll times")

	if len(functions) == 0 {
		for fun := range perf.times {
			functions = append(functions, fun)
		}
	}

	for _, fun := range functions {
		times := perf.times[fun]

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

func (perf *Executor) RemoveAll() {
	defer perf.trackTime(time.Now())
	err := perf.box.RemoveAll()
	if err != nil {
		panic(err)
	}
}

func (perf *Executor) PrepareData(count int) []*Entity {
	defer perf.trackTime(time.Now())

	var result = make([]*Entity, count)
	for i := 0; i < count; i++ {
		result[i] = &Entity{
			String:  fmt.Sprintf("Entity no. %d", i),
			Float64: float64(i),
			Int32:   int32(i),
			Int64:   int64(i),
		}
	}

	return result
}

func (perf *Executor) PutAsync(items []*Entity) {
	defer perf.trackTime(time.Now())

	const retries = 20

	var putErr error
	for _, item := range items {
		for i := 0; i < retries; i++ {
			if _, putErr = perf.box.PutAsync(item); putErr != nil {
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

	if err := perf.ob.AwaitAsyncCompletion(); err != nil {
		panic(err)
	}

	// if retrying failed
	if putErr != nil {
		panic(putErr)
	}
}

func (perf *Executor) PutMany(items []*Entity) {
	defer perf.trackTime(time.Now())

	if _, err := perf.box.PutMany(items); err != nil {
		panic(err)
	}
}

func (perf *Executor) ReadAll(count int) []*Entity {
	defer perf.trackTime(time.Now())

	if items, err := perf.box.GetAll(); err != nil {
		panic(err)
	} else if len(items) != count {
		panic("invalid number of objects read")
	} else {
		return items
	}
}

func (perf *Executor) ChangeValues(items []*Entity) {
	defer perf.trackTime(time.Now())

	count := len(items)
	for i := 0; i < count; i++ {
		items[i].Int64 = items[i].Int64 * 2
	}
}

func (perf *Executor) UpdateAll(items []*Entity) {
	defer perf.trackTime(time.Now())

	if _, err := perf.box.PutMany(items); err != nil {
		panic(err)
	}
}
