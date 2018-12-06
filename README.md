ObjectBox Go API
================
ObjectBox is a superfast database for objects.
Using this Golang API, you can use ObjectBox as an embedded database in your Go application.
A client/server mode will follow soon.

ObjectBox persists your native Go structs using a simple CRUD API:

```go
    id, err := box.Put(&Person{ FirstName: "Joe", LastName:  "Green" })
```

Want details? **[Read the docs](https://golang.objectbox.io/)** or
**[check out the API reference](https://godoc.org/github.com/objectbox/objectbox-go/objectbox)**.

Latest release: [v0.8.0 (2018-12-06)](https://golang.objectbox.io/)

Some features
-------------
* [Powerful queries](https://golang.objectbox.io/queries)
* Secondary indexes based on object properties
* Asynchronous puts
* Automatic model migration (no schema upgrade scripts etc.) 
* (Coming soon: Relations to other objects) 

Installation
------------
To get started with ObjectBox you can get the repository code as usual with go get 
and install the two prerequisites - pre-compiled library and a bindings generator.

```bash
go get github.com/objectbox/objectbox-go
go get github.com/google/flatbuffers/go
go install github.com/objectbox/objectbox-go/cmd/objectbox-gogen/

mkdir objectboxlib && cd objectboxlib
curl https://raw.githubusercontent.com/objectbox/objectbox-c/master/download.sh > download.sh
bash download.sh

```

See [installation docs](https://golang.objectbox.io/install) for more details and further instructions.

Additionally, you can run tests to validate your installation
```bash
go test github.com/objectbox/objectbox-go/...
```

Getting started
---------------
With the dependencies installed, you can start adding entities to your project:
```go
//go:generate objectbox-gogen
​
type Task struct {
	Id           uint64
	Text         string
}
```
And run code generation in your project dir
```bash
go generate ./...
```
This generates a few files in the same folder as the entity - remember to add those to version control (e. g. git).

Once code generation finished successfully, you can start using ObjectBox:
```go
   obx := objectbox.NewBuilder().Model(ObjectBoxModel()).Build()
   box := BoxForTask(obx) // Generated function to provide a Box for Task objects
   id, _ := box.Put(&Task{ Text: "Buy milk" })
```

See the [Getting started](https://golang.objectbox.io/getting-started) section of our docs for a more thorough intro. 

Also, please have a look at the [examples](examples) directory and for the API reference see 
[ObjectBox GoDocs](https://godoc.org/github.com/objectbox/objectbox-go/objectbox) - and the sources in this repo. 

Upgrading to a newer version
----------------------------
When you want to update, please re-run the entire installation process to ensure all components are updated:

* ObjectBox itself (objectbox/objectbox-go)
* Dependencies (flatbuffers)
* libobjectbox
* ObjectBox code generator

This is important as diverging versions of any component might result in errors.
  
This repository also come with a `install.sh` script that can be used for installation and upgrading:

 ```bash
~/go/src/github.com/objectbox/objectbox-go/install.sh
 ```
 
Afterwards, don't forget to re-run the code generation on your project
```bash
go generate ./...
```

License
-------
    Copyright 2018 ObjectBox Ltd. All rights reserved.
    
    Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at
    
        http://www.apache.org/licenses/LICENSE-2.0
    
    Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.

