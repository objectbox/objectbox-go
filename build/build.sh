#!/usr/bin/env bash
set -eu

buildDir=${PWD}/build-artifacts

PATH=$PATH:${buildDir}

function preBuild {
    echo "******** Preparing build ********"
    echo "Creating build artifacts directory '$buildDir'"
    mkdir -p $buildDir
}

function build {
    echo "******** Building ********"
    for CMD in `ls cmd`; do
        echo "building cmd/${CMD}"
        go build -o ${buildDir}/${CMD} ${PWD}/cmd/${CMD}
    done
}

function postBuild {
    echo "******** Collecting artifacts ********"

    echo "The $buildDir contains the following files: "
    ls $buildDir -l
}

function test {
    echo "******** Testing ********"
    ./build/test.sh
}

function generate {
    echo "******** Generating ********"
    go generate ./...
}

preBuild
build
generate
test
postBuild

