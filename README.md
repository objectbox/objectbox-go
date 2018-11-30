ObjectBox Go API
================
ObjectBox is a superfast database for objects.
Using this Golang API, you can use ObjectBox as an embedded database in your Go application.
In this embedded mode, it runs within your application process.

Some features
-------------
* Object storage: put and get native Go structs
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

Docs
----
Documentation is still on-going work.
To get started, please have a look at the [examples](examples) directory and [golang.objectbox.io](https://golang.objectbox.io).

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

