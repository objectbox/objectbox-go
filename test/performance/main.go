package main

import (
	"log"
	"os"

	"github.com/objectbox/objectbox-go/test/performance/perf"
)

func main() {
	var dbName = "db"
	var count = 100000
	var runs = 30

	log.Printf("running the test %d times with %d objects", runs, count)

	// remove old database in case it already exists (and remove it after the test as well)
	os.RemoveAll(dbName)
	defer os.RemoveAll(dbName)

	executor := perf.CreateExecutor(dbName)
	defer executor.Close()

	for i := 0; i < runs; i++ {
		executor.PutAll(count)
		items := executor.ReadAll(count)
		executor.UpdateAll(items)
		executor.RemoveAll()
	}

	executor.PrintTimes()
}
