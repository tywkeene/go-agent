#!/usr/bin/env bash

source ./VERSION

export CGO_ENABLED=0
export GOOS=linux
export GOARCH=amd64

echo "Building go-tracker server binary..."

go build -v \
    -ldflags "-X github.com/tywkeene/go-tracker/version.Version=$VERSION -X github.com/tywkeene/go-tracker/version.CommitHash=$COMMIT" \
    github.com/tywkeene/go-tracker
