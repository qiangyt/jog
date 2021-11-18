#!/bin/sh

set -e

go test ./... -v -covermode=count -coverprofile=coverage.out
go tool cover -html=./coverage.out -o ./coverage.html
