#!/bin/sh

set -e

PROJECT_DIR=$(cd "$(dirname $0)";pwd)
TARGET_DIR=${PROJECT_DIR}/target

rm -rf ${TARGET_DIR}
cd ${PROJECT_DIR}

go_build() {
    local _OS=$1
    local _PREFIX=$2
    local _OS_TARGET_DIR=${TARGET_DIR}/${_OS}

    mkdir -p ${_OS_TARGET_DIR}
    GOOS=${_OS} GOARCH=amd64 go build -trimpath -o ${_OS_TARGET_DIR}/jog${_PREFIX}
}

go generate

go test github.com/qiangyt/jog/util

go_build linux .linux
go_build darwin .darwin
go_build windows .exe
