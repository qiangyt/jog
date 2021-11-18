#!/bin/sh

set -e

go test ./... -v gcflags=all=-l -covermode=count -coverprofile=coverage.out
go tool cover -html=./coverage.out -o ./coverage.html
