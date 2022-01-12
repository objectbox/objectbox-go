#!/usr/bin/env bash
set -eu

# macOS does not have realpath and readlink does not have -f option, so do this instead:
script_dir=$( cd "$(dirname "$0")" ; pwd -P )
cd "${script_dir}/.." # move to project root dir

args="$@"
buildDir=${PWD}/build-artifacts

function preBuild {
    echo "******** Preparing build ********"
    echo "Creating build artifacts directory '$buildDir'"
    mkdir -p "$buildDir"
}

function build {
    echo "******** Building ********"
    for path in cmd/*; do
        echo "building ${path}"
        cd "${path}"
        cmd=$(basename "${path}")
        go build -o "${buildDir}/${cmd}"
        cd -
    done
}

function postBuild {
    echo "******** Collecting artifacts ********"

    echo "The $buildDir contains the following files: "
    ls -l "$buildDir"
}

function test {
    echo "******** Testing ********"

    # on amd64, we run extended tests (memory sanitizer & race checks)
    if [[ $(go env GOARCH) == "amd64" ]]; then
        ./build/test.sh $args -race
    else
        ./build/test.sh $args
    fi
}

function generate {
    echo "******** Generating ********"
    go generate ./...
}

go version

preBuild
build
generate
test
postBuild

