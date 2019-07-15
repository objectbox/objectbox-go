# Quick setup for ObjectBox Go

This readme is based on our corresponding blog article, read it [here](https://objectbox.io/how-to-set-up-objectbox-go-on-raspberry-pi/)!

## Installation

The scripts and sources in this directory are needed to setup ObjectBox Go as easily as possible. 
To get started right away, just execute `./setup.sh` in your shell and you're basically done. 
This command creates the following subdirectories in your home directory:

- `goroot` with the binaries of Go 1.12.7 (only if Go wasn't installed before or the installed version was <1.12)
- `objectbox` with the shell script `update-objectbox.sh` you can execute to easily update ObjectBox upon a new release
- `projects/objectbox-go-test` mainly with the file `main.go` (also part of this directory) which contains a tiny demo application, based on the Tasks example found in ObjectBox GitHub repo 

## Working with ObjectBox Go

For the following commands to work, your current working directory needs to be `~/projects/objectbox-go-test`.

When you have changed your database model, execute `go generate ./...`

Whenever you'd like to rebuild and run your entire Go program, run the following two commands:

    go build
    ./objectbox-go-test
