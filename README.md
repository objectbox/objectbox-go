Do you ♥️ using ObjectBox?
==========================
We want to [hear about your app](https://docs.google.com/forms/d/e/1FAIpQLScIYiOIThcq-AnDVoCvnZOMgxO4S-fBtDSFPQfWldJnhi2c7Q/viewform)!
It will - literally - take just a minute, but help us a lot. Thank you!​ 🙏​

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

Latest release: [v1.2.0 (2020-08-25)](https://golang.objectbox.io/)

Some features
-------------
* [Powerful queries](https://golang.objectbox.io/queries)
* [Relations (to-one, to-many)](https://golang.objectbox.io/relations)
* Secondary indexes based on object properties
* Asynchronous puts
* Automatic model migration (no schema upgrade scripts etc.)
* inline/prefix complex embedded structs (or you can use custom converters)

Getting started
---------------
To install ObjectBox, execute the following command in your project directory. 
You can have a look at [installation docs](https://golang.objectbox.io/install) for more details and further instructions. 
```bash
bash <(curl -s https://raw.githubusercontent.com/objectbox/objectbox-go/main/install.sh)
```

You can run tests to validate your installation
```bash
go test github.com/objectbox/objectbox-go/...
```

With the dependencies installed, you can start adding entities to your project:
```go
//go:generate go run github.com/objectbox/objectbox-go/cmd/objectbox-gogen
​
type Task struct {
	Id   uint64
	Text string
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
* ObjectBox library (libobjectbox.so|dylib; objectbox.dll)
* ObjectBox code generator

This is important as diverging versions of any component might result in errors.
  
The `install.sh` script can also be used for upgrading:
 ```bash
bash <(curl -s https://raw.githubusercontent.com/objectbox/objectbox-go/main/install.sh)
 ```
 
Afterwards, don't forget to re-run the code generation on your project
```bash
go generate ./...
```

License
-------
    Copyright 2018-2019 ObjectBox Ltd. All rights reserved.
    
    Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at
    
        http://www.apache.org/licenses/LICENSE-2.0
    
    Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.

