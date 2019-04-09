#!/usr/bin/env bash
set -eu

# use GNU grep when available (usually MacOS)
if [[ -x "$(command -v ggrep)" ]]; then
    grep=ggrep
else
    grep=grep
fi

# TODO vet
# go fmt $(go list ./... | grep -v /vendor/)
# go vet $(go list ./... | grep -v /vendor/)

# find a list of all tests in the project (go list output ending with "test]")
tests=$(go list -test ./... | $grep -oP '\[\K.*(?=.test\])')

go test "$@" $tests
