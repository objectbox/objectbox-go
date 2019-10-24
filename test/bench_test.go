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
	"flag"
	"fmt"
	"github.com/objectbox/objectbox-go/objectbox"
	"github.com/objectbox/objectbox-go/test/performance/perf"
	"os"
	"testing"
)

// Implements simple benchmarks as an alternative to the "test/performance". However, it doesn't achieve the optimal
// performance as the standalone one so the following benchmarks are only for quick regression testing.

var bulkCount = 10000

func init() {
	// need to do the following two manually in init() function order to have access to testing.Short()
	testing.Init()
	flag.Parse()

	if testing.Short() {
		bulkCount = 100
	}
}

func prepareBenchData(b *testing.B, count int) []*perf.Entity {
	b.StopTimer()

	var result = make([]*perf.Entity, count)
	for i := 0; i < count; i++ {
		result[i] = &perf.Entity{
			String:  fmt.Sprintf("Entity no. %d", i),
			Float64: float64(i),
			Int32:   int32(i),
			Int64:   int64(i),
		}
	}

	b.StartTimer()
	return result
}

type benchmarkEnv struct {
	dbName string
	ob     *objectbox.ObjectBox
	box    *perf.EntityBox
	b      *testing.B
}

func newBenchEnv(b *testing.B) *benchmarkEnv {
	b.StopTimer()

	var env = &benchmarkEnv{
		dbName: "testdata",
		b:      b,
	}

	var err error
	env.ob, err = objectbox.NewBuilder().Directory(env.dbName).Model(perf.ObjectBoxModel()).Build()
	env.check(err)

	env.box = perf.BoxForEntity(env.ob)

	b.StartTimer()
	b.ReportAllocs()
	return env
}

func (env *benchmarkEnv) close() {
	env.b.StopTimer()
	env.ob.Close()
	os.RemoveAll(env.dbName)
	env.b.StartTimer()
}

func (env *benchmarkEnv) check(err error) {
	if err != nil {
		env.b.Error(err)
	}
}

func BenchmarkPutMany(b *testing.B) {
	var env = newBenchEnv(b)
	defer env.close()
	var inserts = prepareBenchData(b, bulkCount)

	b.Run(fmt.Sprintf("count=%v", bulkCount), func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			_, err := env.box.PutMany(inserts)
			env.check(err)

			b.StopTimer()
			env.box.RemoveAll()
			b.StartTimer()
		}
	})
}

func BenchmarkGetAll(b *testing.B) {
	var env = newBenchEnv(b)
	defer env.close()
	var inserts = prepareBenchData(b, bulkCount)

	b.StopTimer()
	_, err := env.box.PutMany(inserts)
	env.check(err)
	b.StartTimer()

	b.Run("GetAll", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			objects, err := env.box.GetAll()
			if err != nil {
				b.Error(err)
			} else if len(objects) != bulkCount {
				b.Errorf("invalid number of objects received: %v instead of %v", len(objects), bulkCount)
			}
		}
	})
}

func BenchmarkRemoveAll(b *testing.B) {
	var env = newBenchEnv(b)
	defer env.close()
	var inserts = prepareBenchData(b, bulkCount)

	b.Run("count=%v", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			b.StopTimer()
			_, err := env.box.PutMany(inserts)
			env.check(err)
			b.StartTimer()
			err = env.box.RemoveAll()
			env.check(err)
		}
	})
}

// BenchmarkTxPut executes many individual puts in a single transaction.
func BenchmarkTxPut(b *testing.B) {
	var env = newBenchEnv(b)
	defer env.close()

	// prepare the data first
	var inserts = prepareBenchData(b, b.N)

	// execute in a single transaction
	env.check(env.ob.RunInWriteTx(func() error {
		for n := 0; n < b.N; n++ {
			_, err := env.box.Put(inserts[n])
			env.check(err)
		}
		return nil
	}))
}

// BenchmarkSlowPut executes many individual puts, each in its own transaction (internally).
func BenchmarkSlowPut(b *testing.B) {
	var env = newBenchEnv(b)
	defer env.close()

	// prepare the data first
	var inserts = prepareBenchData(b, b.N)

	for n := 0; n < b.N; n++ {
		_, err := env.box.Put(inserts[n])
		env.check(err)
	}
}

// BenchmarkTxGet reads a single object from DB many times, all in a single transaction
func BenchmarkTxGet(b *testing.B) {
	var env = newBenchEnv(b)
	defer env.close()

	// prepare the data first
	var inserts = prepareBenchData(b, 1)
	b.StopTimer()
	_, err := env.box.Put(inserts[0])
	env.check(err)
	b.StartTimer()

	env.check(env.ob.RunInReadTx(func() error {
		for n := 0; n < b.N; n++ {
			_, err := env.box.Get(inserts[0].Id)
			env.check(err)
		}
		return nil
	}))
}

// BenchmarkGet reads a single object from DB many times, each in its own transaction (internally)
func BenchmarkSlowGet(b *testing.B) {
	var env = newBenchEnv(b)
	defer env.close()

	// prepare the data first
	var inserts = prepareBenchData(b, 1)
	b.StopTimer()
	_, err := env.box.Put(inserts[0])
	env.check(err)
	b.StartTimer()

	for n := 0; n < b.N; n++ {
		_, err := env.box.Get(inserts[0].Id)
		env.check(err)
	}
}
