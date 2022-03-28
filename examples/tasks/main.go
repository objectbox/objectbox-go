/*
 * Copyright 2018-2022 ObjectBox Ltd. All rights reserved.
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

package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/objectbox/objectbox-go/examples/tasks/internal/model"
	"github.com/objectbox/objectbox-go/objectbox"
)

func main() {
	objectBox := initObjectBox()
	defer objectBox.Close()

	box := model.BoxForTask(objectBox)

	checkStartSyncClient(objectBox, box)

	runInteractiveShell(box)
}

func runInteractiveShell(box *model.TaskBox) {
	// our simple interactive shell
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Welcome to the ObjectBox tasks-list app example")
	printHelp()

	for {
		fmt.Print("$ ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		//input = strings.TrimSuffix(input, "\n")
		input = strings.TrimSpace(input)
		args := strings.Fields(input)

		switch strings.ToLower(args[0]) {
		case "new":
			createTask(box, strings.Join(args[1:], " "))
		case "done":
			if len(args) != 2 {
				fmt.Fprintf(os.Stderr, "wrong number of arguments, expecting exactly one\n")
			} else if id, err := strconv.ParseUint(args[1], 10, 64); err != nil {
				fmt.Fprintf(os.Stderr, "could not parse ID: %s\n", err)
			} else {
				setDone(box, id)
			}
		case "ls":
			if len(args) < 2 {
				printList(box, false)
			} else if args[1] == "-a" {
				printList(box, true)
			} else {
				fmt.Fprintf(os.Stderr, "unknown argument %s\n", args[1])
				fmt.Println()
			}
		case "exit":
			return
		case "help":
			printHelp()
		default:
			fmt.Fprintf(os.Stderr, "unknown command %s\n", input)
			printHelp()
		}
	}
}

func initObjectBox() *objectbox.ObjectBox {
	objectBox, err := objectbox.NewBuilder().Model(model.ObjectBoxModel()).Build()
	if err != nil {
		panic(err)
	}
	return objectBox
}

func printHelp() {
	fmt.Println("Available commands are: ")
	fmt.Println("    ls [-a]        list tasks - unfinished or all (-a flag)")
	fmt.Println("    new Task text  create a new task with the text 'Task text'")
	fmt.Println("    done ID        mark task with the given ID as done")
	fmt.Println("    exit           close the program")
	fmt.Println("    help           display this help")
}

func createTask(box *model.TaskBox, text string) {
	task := &model.Task{
		Text:         text,
		DateCreated:  time.Now(),
		DateFinished: time.Unix(0, 0), // use "epoch start" to unify values across platforms (e.g for Sync)
	}

	if id, err := box.Put(task); err != nil {
		fmt.Fprintf(os.Stderr, "could not create task: %s\n", err)
	} else {
		task.Id = id
		fmt.Printf("task ID %d successfully created\n", task.Id)
	}
}

func printList(box *model.TaskBox, all bool) {
	var list []*model.Task
	var err error

	if all { // load all tasks
		list, err = box.GetAll()
	} else { // load only unfinished tasks (value 0 is "epoch start")
		list, err = box.Query(model.Task_.DateFinished.Equals(0)).Find()
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "could not list tasks: %s\n", err)
	}

	fmt.Printf("%3s  %-23s  %-23s  %s\n", "ID", "Created", "Finished", "Text")
	for _, task := range list {
		fmt.Printf("%3d  %-23s  %-23s  %s\n",
			task.Id, task.DateCreated.Format("2006-01-02 15:04:05"), task.DateFinished.Format("2006-01-02 15:04:05"), task.Text)
	}
}

func setDone(box *model.TaskBox, id uint64) {
	if task, err := box.Get(id); err != nil {
		fmt.Fprintf(os.Stderr, "could not read task ID %d: %s\n", id, err)
	} else if task == nil {
		fmt.Fprintf(os.Stderr, "task ID %d doesn't exist\n", id)
	} else {
		task.DateFinished = time.Now()
		if _, err := box.Put(task); err != nil {
			fmt.Fprintf(os.Stderr, "could not update task ID %d: %s\n", id, err)
		} else {
			fmt.Printf("task ID %d completed at %s\n", id, task.DateFinished.String())
		}
	}
}

func checkStartSyncClient(ob *objectbox.ObjectBox, box *model.TaskBox) { // only if sync-enabled library is used

	if objectbox.SyncIsAvailable() {
		syncClient, err := objectbox.NewSyncClient(
			ob,
			"ws://127.0.0.1", // wss for SSL, ws for unencrypted traffic
			objectbox.SyncCredentialsNone())

		if err == nil {
			syncClient.Start() // Connect and start syncing.
			fmt.Println("Sync client started.")
			syncClient.SetChangeListener(func(changes []*objectbox.SyncChange) {

				fmt.Printf("received %d changes\n", len(changes))
				printList(box, true)
			})
		} else {
			fmt.Println("Could not start the sync client.")
		}
	} else {
		fmt.Println("Sync is not available. Please go to https://sync.objectbox.io/ for more information.")
	}
}
