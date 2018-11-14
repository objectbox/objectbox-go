package main

import (
	"fmt"
	"time"

	"github.com/objectbox/objectbox-go/objectbox"

	"github.com/objectbox/objectbox-go/examples/tasks/internal/model"
)

func main() {
	//reader := bufio.NewReader(os.Stdin)

	fmt.Println("Welcome to the ObjectBox tasks-list example")
	fmt.Println()

	// load objectbox
	ob := initObjectBox()
	//defer ob.Destroy()
	box := model.BoxForTask(ob)
	//defer box.Destroy()

	createTask(box, "Text")
	printList(box)
	setDone(box, 1)
	printList(box)

	box.Destroy()
	ob.Destroy()
}

func initObjectBox() *objectbox.ObjectBox {
	builder := objectbox.NewObjectBoxBuilder().LastEntityId(1, 2351059844987470165)
	builder.RegisterBinding(model.TaskBinding{})
	objectBox, err := builder.Build()
	if err != nil {
		panic(err)
	}
	return objectBox
}

func createTask(box *model.TaskBox, text string) {
	task := &model.Task{
		Text:        text,
		DateCreated: obNow(),
	}

	if id, err := box.Put(task); err != nil {
		fmt.Printf("could not create task: %s\n", err)
	} else {
		task.Id = id
		fmt.Printf("task ID %d successfully created\n", task.Id)
	}
}

func printList(box *model.TaskBox) {
	if list, err := box.GetAll(); err != nil {
		fmt.Printf("could not list tasks: %s\n", err)
	} else {
		fmt.Printf("%3s  %-29s  %-29s  %s\n", "ID", "Created", "Finished", "Text")
		for _, task := range list {
			fmt.Printf("%3d  %-29s  %-29s  %s\n",
				task.Id, fmtTime(task.DateCreated), fmtTime(task.DateFinished), task.Text)
		}
	}
}

func setDone(box *model.TaskBox, id uint64) {
	if task, err := box.Get(id); err != nil {
		fmt.Printf("could not read task ID %d: %s\n", id, err)
	} else {
		task.DateFinished = obNow()
		if _, err := box.Put(task); err != nil {
			fmt.Printf("could not update task ID %d: %s\n", id, err)
		} else {
			fmt.Printf("task ID %d completed at %s\n", id, fmtTime(task.DateFinished))
		}
	}
}

func fmtTime(obTimestamp int64) string {
	if obTimestamp == 0 {
		return ""
	} else {
		return time.Unix(obTimestamp/1000, obTimestamp%1000*1000000).String()
	}
}

func obNow() int64 {
	return time.Now().Unix() * 1000
}
