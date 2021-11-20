#!/bin/sh

set -e

rm -f coverage.out coverage.html
go test ./... -v  -covermode=count -coverprofile=coverage.out gcflags=all=-l
go tool cover -html=./coverage.out -o ./coverage.html