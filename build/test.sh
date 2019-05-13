#!/usr/bin/env bash
set -eu

# use GNU grep when available (usually MacOS)
if [[ -x "$(command -v ggrep)" ]]; then
    grep=ggrep
else
    grep=grep
fi

unformatted_files=$(gofmt -l .)
if [[ ${unformatted_files} ]]; then
    echo "Some files are not formatted properly. You can use \`gofmt -l -w .\` to fix them."
    printf "%s\n" ${unformatted_files}
fi

go vet ./...

# find a list of all tests in the project (go list output ending with "test]")
tests=$(go list -test ./... | $grep -oP '\[\K.*(?=.test\])')

go test "$@" $tests
