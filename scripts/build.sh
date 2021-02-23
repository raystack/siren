#!/bin/bash
NAME="github.com/odpf/siren"
VERSION=$(git describe --always --tags 2>/dev/null)

SYS=("linux" "darwin")
ARCH=("386" "amd64")
BUILD_DIR="dist"

build() {
    EXECUTABLE_NAME=$1
    LD_FLAGS=$2
    for os in ${SYS[*]}; do
        for arch in ${ARCH[*]}; do

            # create a folder named via the combination of os and arch
            TARGET="./$BUILD_DIR/${os}-${arch}"
            mkdir -p $TARGET

            # place the executable within that folder
            executable="${TARGET}/$EXECUTABLE_NAME"
            echo $executable
            GOOS=$os GOARCH=$arch go build -ldflags "$LD_FLAGS" -o $executable $NAME
        done
    done
}

build siren "-X main.Version=${VERSION}"
