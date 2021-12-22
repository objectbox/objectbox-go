Go ObjectBox Database API
================
ObjectBox is a superfast Go database persisting objects; [check the performance benchmarks vs SQLite (GORM) & Storm](https://objectbox.io/go-1-0-release-and-performance-benchmarks/). Using this Golang API, you can use ObjectBox as an embedded database in your Go application.

ObjectBox persists your native Go structs using a simple CRUD API:

```go
id, err := box.Put(&Person{ FirstName: "Joe", LastName:  "Green" })
```

Want details? **[Read the docs](https://golang.objectbox.io/)** or
**[check out the API reference](https://godoc.org/github.com/objectbox/objectbox-go/objectbox)**.

Latest release: [v1.5.0 (2021-08-18)](https://golang.objectbox.io/)

High-performance Golang database
-------------
🏁 **High performance** on restricted devices, like IoT gateways, micro controllers, ECUs etc.\
🪂 **Resourceful** with minimal CPU, power and Memory usage for maximum flexibility and sustainability\
🔗 **Relations:** object links / relationships are built-in\
💻 **Multiplatform:** Linux, Windows, Android, iOS, macOS

🌱 **Scalable:** handling millions of objects resource-efficiently with ease\
💐 **Queries:** filter data as needed, even across relations\
🦮 **Statically typed:** compile time checks & optimizations\
📃 **Automatic schema migrations:** no update scripts needed

**And much more than just data persistence**\
👥 **[ObjectBox Sync](https://objectbox.io/sync/):** keeps data in sync between devices and servers\
🕒 **[ObjectBox TS](https://objectbox.io/time-series-database/):** time series extension for time based data

Enjoy ❤️


Getting started
---------------
To install ObjectBox, execute the following command in your project directory. 
You can have a look at [installation docs](https://golang.objectbox.io/install) for more details and further instructions. 
```bash
bash <(curl -s https://raw.githubusercontent.com/objectbox/objectbox-go/main/install.sh)
```

To install [ObjectBox Sync](https://objectbox.io/sync/) variant of the library, pass `--sync` argument to the command above:

```bash
bash <(curl -s https://raw.githubusercontent.com/objectbox/objectbox-go/main/install.sh) --sync
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

Other languages/bindings
------------------------
ObjectBox supports multiple platforms and languages:

* [ObjectBox Java/Kotlin](https://github.com/objectbox/objectbox-java): runs on Android, desktop, and servers
* [ObjectBox Swift](https://github.com/objectbox/objectbox-swift): build fast mobile apps for iOS (and macOS) 
* [ObjectBox Dart/Flutter](https://github.com/objectbox/objectbox-dart): cross-platform for mobile and desktop apps 
* [ObjectBox C and C++](https://github.com/objectbox/objectbox-c): native speed with zero copy access to FlatBuffer objects


How I help ObjectBox?
---------------------------
We believe, ObjectBox is super easy to use. We are on a mission to make developers’ lives better, by building developer tools that are intuitive and fun to code with.

To do that, we want your feedback: what do you love? What's amiss? Where do you struggle in everyday app development?

**We're looking forward to receiving your comments and requests:**

- Add [GitHub issues](https://github.com/ObjectBox/objectbox-go/issues) 
- Upvote issues you find important by hitting the 👍/+1 reaction button
- Drop us a line via [@ObjectBox_io](https://twitter.com/ObjectBox_io/)
- ⭐ us, if you like what you see 

Thank you! 🙏

Keep in touch: For general news on ObjectBox, [check our blog](https://objectbox.io/blog)!

License
-------
    Copyright 2018-2021 ObjectBox Ltd. All rights reserved.
    
    Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at
    
        http://www.apache.org/licenses/LICENSE-2.0
    
    Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.

