#!/usr/bin/env bash

source ./VERSION

export CGO_ENABLED=0
export GOOS=linux
export GOARCH=amd64


function build(){
    echo "Building go-tracker $1 binary..."
    rm -f tracker-$1
    go build -v -o tracker-$1 \
        -ldflags "-X github.com/tywkeene/go-tracker/version.Version=$VERSION \
        -X github.com/tywkeene/go-tracker/version.CommitHash=$COMMIT" \
        github.com/tywkeene/go-tracker/cmd/$1
}

function usage(){
    printf "Usage:\n$0 -c <build client binary>\n$0 -s <build server binary>\n$0 -a <build all binaries>\n"
}

if [ -z "$1" ]; then
    usage
    exit -1
fi

while getopts "cash" opt; do
    case "$opt" in
        h) usage
            ;;
        c) build "client"
            ;;
        s) build "server"
            ;;
        a) build "server"
            build "client"
            ;;
    esac
done
