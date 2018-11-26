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

	log.Printf("running the test with %d objects", count)

	// remove old database in case it already exists (and remove it after the test as well)
	os.RemoveAll(dbName)
	defer os.RemoveAll(dbName)

	executor := perf.CreateExecutor(dbName)
	defer executor.Close()

	inserts := executor.PrepareData(count)

	executor.PutAsync(inserts)
	executor.RemoveAll()

	executor.PutAll(inserts)

	items := executor.ReadAll(count)
	executor.ChangeValues(items)
	executor.UpdateAll(items)

	executor.PrintTimes([]string{})
}
