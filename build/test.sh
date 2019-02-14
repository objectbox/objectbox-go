#!/usr/bin/env bash
set -eu

# TODO vet
# go fmt $(go list ./... | grep -v /vendor/)
# go vet $(go list ./... | grep -v /vendor/)

# this is a list of all tests in the project
tests=$(go list -test ./... | grep -oP '\[\K.*(?=.test\])')

go test "$@" $tests
