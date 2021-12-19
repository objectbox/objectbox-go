#!/usr/bin/env bash
set -eu

unformatted_files=$(gofmt -l .)
if [[ ${unformatted_files} ]]; then
    echo "Some files are not formatted properly. You can use \`gofmt -l -w .\` to fix them:"
    printf "%s\n" "${unformatted_files}"
    exit 1
fi

go vet ./...

go test "$@" ./...
