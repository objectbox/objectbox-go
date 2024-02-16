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
set +e
# ignore '# <package>' comments to stablize checking false positive below; 
# newer go version outputs a second-line with bracket ('# [<package]').
go_vet_result=$(go 2>&1 vet ./... | grep -v ^#)
go_vet_rc=$?
set -e

echo "$go_vet_result"
go_vet_result_lines=$(echo "$go_vet_result" | wc -l)
if [ $go_vet_rc -ne 0 ]; then
  if [[ $go_vet_result_lines -eq 1 && $go_vet_result == *objectbox[/\\]c-callbacks.go*possible\ misuse\ of\ unsafe.Pointer* ]]; then
    echo "Ignoring known false positive of go vet"
    go_vet_rc=0
  else
    echo "go vet failed ($go_vet_rc)"
    # Fail later because we want to run tests for now too; was: exit $go_vet_rc
  fi
fi

echo "******** Testing: go test ********"
go test "$@" ./...

if [ $go_vet_rc -ne 0 ]; then
  echo "go vet failed ($go_vet_rc)"
  exit $go_vet_rc
fi