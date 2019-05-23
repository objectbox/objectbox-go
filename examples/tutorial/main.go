package main

import "fmt"
import "github.com/objectbox/objectbox-go/objectbox"

type Task struct {
	Id   uint64
	Text string
}

func main() {
	obx, err := objectbox.NewBuilder().Model(ObjectBoxModel()).Build()
	if err != nil {
		panic(err)
	}
	box := BoxForTask(obx)
	id, _ := box.Put(&Task{Text: "Buy milk"})
	fmt.Printf("new entry got id %d\n", id)
}
