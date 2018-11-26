package perf

import (
	"fmt"
	"path"
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

	objectBox, err := objectbox.NewObjectBoxBuilder().Name(dbName).Model(CreateObjectBoxModel()).Build()
	if err != nil {
		panic(err)
	}

	perf.ob = objectBox
	perf.box = BoxForEntity(objectBox)
}

func (perf *Executor) Close() {
	defer perf.trackTime(time.Now())

	perf.box.Close()
	perf.ob.Close()
}

func (perf *Executor) trackTime(start time.Time) {
	elapsed := time.Since(start)

	pc, _, _, _ := runtime.Caller(1)
	fun := path.Ext(runtime.FuncForPC(pc).Name())
	//if perf.logTimes {
	//	//log.Printf("%s took %s", f.Name()), elapsed)
	//} else {
	//	// TODO
	//}

	perf.times[fun] = append(perf.times[fun], elapsed)
}

func (perf *Executor) PrintTimes() {
	// print the whole data as a table
	fmt.Println("Function\tRuns\tAverage ms\tAll times")
	for fun, times := range perf.times {
		sum := int64(0)
		for _, duration := range times {
			sum += duration.Nanoseconds()
		}
		fmt.Printf("%s\t%d\t%f", fun[1:], len(times), float64(sum/int64(len(times)))/1000000)

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

	for _, item := range items {
		if _, err := perf.box.PutAsync(item); err != nil {
			panic(err)
		}
	}

	perf.ob.AwaitAsyncCompletion()
}

func (perf *Executor) PutAll(items []*Entity) {
	defer perf.trackTime(time.Now())

	if _, err := perf.box.PutAll(items); err != nil {
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

func (perf *Executor) UpdateAll(items []*Entity) {
	defer perf.trackTime(time.Now())

	for _, item := range items {
		item.Int64 = item.Int64 * 2
	}

	if _, err := perf.box.PutAll(items); err != nil {
		panic(err)
	}
}
