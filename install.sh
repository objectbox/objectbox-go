#!/usr/bin/env bash

mkdir objectboxlib && cd objectboxlib
curl https://raw.githubusercontent.com/objectbox/objectbox-c/master/download.sh > download.sh
bash download.sh

go get github.com/objectbox/objectbox-go/...
go install github.com/objectbox/objectbox-go/cmd/objectbox-gogen/

echo "Done"