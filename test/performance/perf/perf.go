package perf

import (
	"log"
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
	for fun, times := range perf.times {
		sum := int64(0)
		for _, duration := range times {
			sum += duration.Nanoseconds()
		}

		avgDuration := time.Duration(sum / int64(len(times)))
		log.Printf("%s has been executed %d times and took %s on average", fun, len(times), avgDuration)
	}
}

func (perf *Executor) RemoveAll() {
	defer perf.trackTime(time.Now())
	err := perf.box.RemoveAll()
	if err != nil {
		panic(err)
	}
}

func (perf *Executor) PutAsync(count int) {
	defer perf.trackTime(time.Now())

	var proto = &Entity{}
	for i := 0; i < count; i++ {
		if _, err := perf.box.PutAsync(proto); err != nil {
			panic(err)
		}
	}

	perf.ob.AwaitAsyncCompletion()
}

func (perf *Executor) PutAll(count int) {
	defer perf.trackTime(time.Now())

	items := make([]*Entity, count)

	var proto = &Entity{}
	for i := 0; i < count; i++ {
		items[i] = proto
	}

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

	// NOTE this takes about 20% of the function time
	var newValue = uint32(1)
	for _, item := range items {
		item.Value = newValue
	}

	if _, err := perf.box.PutAll(items); err != nil {
		panic(err)
	}
}
