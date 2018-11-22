ObjectBox Go API
================
ObjectBox is a superfast database for objects.
Using this Golang API, you can us ObjectBox as an embedded database in your Go application.
In this embedded mode, it runs within your application process.

Some features
-------------
* Object storage based on [FlatBuffers](https://google.github.io/flatbuffers/)
* Secondary indexes based on object properties
* Simple CRUD API
* Asynchronous puts
* Automatic model migration (no schema upgrade scripts etc.) 
* (Coming soon: Powerful queries) 
* (Coming soon: Relations to other objects) 

Installation
------------
To get started with ObjectBox you can get the repository code as usual with go get 
and install the two prerequisites - pre-compiled library and a bindings generator.

```bash
go get github.com/objectbox/objectbox-go

wget https://raw.githubusercontent.com/objectbox/objectbox-c/master/download.sh
bash ./download.sh

go install $GOPATH/src/github.com/objectbox/objectbox-go/cmd/objectbox-bindings/
```

See [installation docs](https://golang.objectbox.io/install) for more details and further instructions.

Docs
----
Documentation is still on-going work.
To get started, please have a look at the [examples](examples) directory and [golang.objectbox.io](https://golang.objectbox.io).

Current state
-------------
As this is still an early version of the Go APIs, they are not as convenient as the [Java/Kotlin APIs](https://docs.objectbox.io/),
which deeply integrate into the language using e.g. [@Entity annotations](https://docs.objectbox.io/entity-annotations).

A better language integration could be build based on reflection, code generation, or a combination of both.

Building
------------
```sh
./build/build.sh
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

