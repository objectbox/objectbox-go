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

func prepareBenchData(count int) []*perf.Entity {
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

func benchmark(b *testing.B, count int) {
	var dbName = "testdata"
	ob, err := objectbox.NewBuilder().Directory(dbName).Model(perf.ObjectBoxModel()).Build()
	if err != nil {
		b.Error(err)
		b.FailNow()
	}
	defer func() {
		ob.Close()
		os.RemoveAll(dbName)
	}()

	var box = perf.BoxForEntity(ob)

	// prepare the data first
	var inserts = prepareBenchData(count)

	b.Run("PutMany", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			_, err := box.PutMany(inserts)
			if err != nil {
				b.Error(err)
			}

			b.StopTimer()
			box.RemoveAll()
			b.StartTimer()
		}
	})

	b.Run("GetAll", func(b *testing.B) {
		b.StopTimer()
		_, err := box.PutMany(inserts)
		if err != nil {
			b.Error(err)
		}
		b.StartTimer()
		for n := 0; n < b.N; n++ {
			objects, err := box.GetAll()
			if err != nil {
				b.Error(err)
			} else if len(objects) != count {
				b.Errorf("invalid number of objects received: %v instead of %v", len(objects), count)
			}
		}
	})

	b.Run("RemoveAll", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			b.StopTimer()
			_, err := box.PutMany(inserts)
			if err != nil {
				b.Error(err)
			}
			b.StartTimer()
			err = box.RemoveAll()
			if err != nil {
				b.Error(err)
			}
		}
	})
}

// Implements simple benchmarks as an alternative to the "test/performance", though it doesn't achieve the optimal
// performance of that one.
func BenchmarkObjectBox(b *testing.B) {
	b.Run("count=1000", func(b *testing.B) {
		benchmark(b, 1000)
	})

	if testing.Short() {
		return
	}

	b.Run("count=10000", func(b *testing.B) {
		benchmark(b, 10*1000)
	})

	b.Run("count=1000000", func(b *testing.B) {
		benchmark(b, 1000*1000)
	})
}
