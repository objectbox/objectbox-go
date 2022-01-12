#!/usr/bin/env bash
set -eu

# macOS does not have realpath and readlink does not have -f option, so do this instead:
script_dir=$( cd "$(dirname "$0")" ; pwd -P )
cd "${script_dir}/.." # move to project root dir

unformatted_files=$(gofmt -l .)
if [[ ${unformatted_files} ]]; then
    echo "Some files are not formatted properly. You can use \`gofmt -l -w .\` to fix them:"
    printf "%s\n" "${unformatted_files}"
    exit 1
fi

go vet ./...

go test "$@" ./...
