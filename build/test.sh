#!/usr/bin/env bash
set -eu

# TODO vet
# go fmt $(go list ./... | grep -v /vendor/)
# go vet $(go list ./... | grep -v /vendor/)
# go test -race $(go list ./... | grep -v /vendor/)

# on amd64, we run extended tests (memory sanitizer & race checks)
if [[ $(go env GOARCH) == "amd64" ]]; then
    # exclude the generator
    extendedTests=$(go list ./... | grep -v test/generator)
    go test ./test/generator/...
    go test -race $extendedTests
    go test -msan $extendedTests

else
    go test ./...
fi
