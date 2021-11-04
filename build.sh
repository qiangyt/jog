#!/bin/sh

set -e

PROJECT_DIR=$(cd "$(dirname $0)";pwd)
TARGET_DIR=${PROJECT_DIR}/target

rm -rf ${TARGET_DIR}
cd ${PROJECT_DIR}

go_build() {
    local _OS=$1
    local _ARCH=$2
    local _PREFIX=$3
    local _OS_TARGET_DIR=${TARGET_DIR}/${_OS}

    mkdir -p ${_OS_TARGET_DIR}
    GOOS=${_OS} GOARCH=amd64 go build -trimpath -o ${_OS_TARGET_DIR}/jog${_PREFIX}
}

go generate
go fmt ./...
go test ./...

#go install github.com/gojp/goreportcard/cmd/goreportcard-cli@latest
#https://github.com/alecthomas/gometalinter/releases/tag/v3.0.0
#goreportcard-cli -v
goreportcard-cli

go_build linux amd64 .linux
go_build darwin amd64 .darwin_amd64
go_build darwin arm64 .darwin_arm64
go_build windows amd64 .exe
