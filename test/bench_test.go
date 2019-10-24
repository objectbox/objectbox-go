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
	"fmt"
	"github.com/objectbox/objectbox-go/objectbox"
	"github.com/objectbox/objectbox-go/test/performance/perf"
	"os"
	"testing"
)

// Implements simple benchmarks as an alternative to the "test/performance". However, it doesn't achieve the optimal
// performance as the standalone one so the following benchmarks are only for quick regression testing.

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
	if err != nil {
		b.Error(err)
		b.FailNow()
	}

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

func benchmarkBulk(b *testing.B, count int) {
	var env = newBenchEnv(b)
	defer env.close()

	// prepare the data first
	var inserts = prepareBenchData(b, count)

	b.Run("PutMany", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			_, err := env.box.PutMany(inserts)
			if err != nil {
				b.Error(err)
			}

			b.StopTimer()
			env.box.RemoveAll()
			b.StartTimer()
		}
	})

	b.Run("GetAll", func(b *testing.B) {
		b.ReportAllocs()
		b.StopTimer()
		_, err := env.box.PutMany(inserts)
		if err != nil {
			b.Error(err)
		}
		b.StartTimer()
		for n := 0; n < b.N; n++ {
			objects, err := env.box.GetAll()
			if err != nil {
				b.Error(err)
			} else if len(objects) != count {
				b.Errorf("invalid number of objects received: %v instead of %v", len(objects), count)
			}
		}
	})

	b.Run("RemoveAll", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			b.StopTimer()
			_, err := env.box.PutMany(inserts)
			if err != nil {
				b.Error(err)
			}
			b.StartTimer()
			err = env.box.RemoveAll()
			if err != nil {
				b.Error(err)
			}
		}
	})
}

// BenchmarkBulk tests bulk put, read & remove.
func BenchmarkBulk(b *testing.B) {
	b.Run("count=1000", func(b *testing.B) {
		benchmarkBulk(b, 1000)
	})

	if testing.Short() {
		return
	}

	b.Run("count=10000", func(b *testing.B) {
		benchmarkBulk(b, 10*1000)
	})

	b.Run("count=1000000", func(b *testing.B) {
		benchmarkBulk(b, 1000*1000)
	})
}

// BenchmarkTxPut executes many individual puts in a single transaction.
func BenchmarkTxPut(b *testing.B) {
	var env = newBenchEnv(b)
	defer env.close()

	// prepare the data first
	var inserts = prepareBenchData(b, b.N)

	// execute in a single transaction
	err := env.ob.RunInWriteTx(func() error {
		for n := 0; n < b.N; n++ {
			_, err := env.box.Put(inserts[n])
			if err != nil {
				b.Error(err)
			}
		}
		return nil
	})
	if err != nil {
		b.Error(err)
	}
}

// BenchmarkSlowPut executes many individual puts, each in its own transaction (internally).
func BenchmarkSlowPut(b *testing.B) {
	var env = newBenchEnv(b)
	defer env.close()

	// prepare the data first
	var inserts = prepareBenchData(b, b.N)

	for n := 0; n < b.N; n++ {
		_, err := env.box.Put(inserts[n])
		if err != nil {
			b.Error(err)
		}
	}
}

// BenchmarkGet reads a single object from DB many times
func BenchmarkGet(b *testing.B) {
	var env = newBenchEnv(b)
	defer env.close()

	// prepare the data first
	var inserts = prepareBenchData(b, 1)
	b.StopTimer()
	_, err := env.box.Put(inserts[0])
	if err != nil {
		b.Error(err)
	}
	b.StartTimer()

	for n := 0; n < b.N; n++ {
		_, err := env.box.Get(inserts[0].Id)
		if err != nil {
			b.Error(err)
		}
	}
}