#!/usr/bin/env bash
set -eu

# macOS does not have realpath and readlink does not have -f option, so do this instead:
script_dir=$( cd "$(dirname "$0")" ; pwd -P )
cd "${script_dir}/.." # move to project root dir

go_version=$(go version)
echo $go_version
if [[ $go_version == *go1.17.* ]]; then  # Keep in sync with our CI gatekeeper job; TODO update to latest
  # gofmt is version specific, so only run this for our reference version for development
  echo "Reference Go version found for gofmt; checking source format..."
  echo "******** Testing: gofmt ********"
  unformatted_files=$(gofmt -l .)
  if [[ ${unformatted_files} ]]; then
      echo "Some files are not formatted properly. You can use \`gofmt -l -w .\` to fix them:"
      printf "%s\n" "${unformatted_files}"
      exit 1
  fi
else
  echo "The found Go version is not our reference for gofmt; skipping source format check"
fi

echo "******** Testing: go vet ********"
go vet ./...

echo "******** Testing: go test ********"
go test "$@" ./...
