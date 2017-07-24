#!/usr/bin/env bash

set -e
set -u

gofmt -w -s .
golint -set_exit_status ./...
go test -v ./...
