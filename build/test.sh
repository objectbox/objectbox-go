#!/usr/bin/env bash
set -eu

go fmt $(go list ./... | grep -v /vendor/)
go vet $(go list ./... | grep -v /vendor/)
go test -race $(go list ./... | grep -v /vendor/)
