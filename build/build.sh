#!/usr/bin/env bash
set -eu

cwd=$PWD
buildDir=build-artifacts

function preBuild {
    echo "******** Preparing build ********"
    echo "Creating build artifacts directory '$buildDir'"
    mkdir -p $buildDir
}

function build {
    echo ""
    cd $buildDir

    echo "******** Building ********"
    for CMD in `ls $cwd/cmd`; do
        echo "building cmd/$CMD"
        go build $cwd/cmd/$CMD
    done

    echo ""
    cd $cwd
}

function postBuild {
    echo "******** Collecting artifacts ********"

    echo "The $buildDir contains the following files: "
    ls $buildDir -l
}

function test {
    echo "******** Testing ********"
    go test ./...

}

preBuild
build
test
postBuild

