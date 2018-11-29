#!/usr/bin/env bash

go get github.com/objectbox/objectbox-go
go get github.com/google/flatbuffers/go
go install github.com/objectbox/objectbox-go/cmd/objectbox-gogen/

mkdir objectboxlib && cd objectboxlib
curl https://raw.githubusercontent.com/objectbox/objectbox-c/master/download.sh > download.sh
bash download.sh

echo "Done"