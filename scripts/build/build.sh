#!/usr/bin/env bash

source ./VERSION

export CGO_ENABLED=0
export GOOS=linux
export GOARCH=amd64

echo "Building go-tracker server binary..."
go build -v -o tracker-server \
    -ldflags "-X github.com/tywkeene/go-tracker/cmd/server/version.Version=$VERSION -X github.com/tywkeene/go-tracker/cmd/server/version.CommitHash=$COMMIT" \
    github.com/tywkeene/go-tracker/cmd/server

echo "Building go-tracker client binary..."
go build -v -o tracker-client \
    -ldflags "-X github.com/tywkeene/go-tracker/cmd/client/version.Version=$VERSION -X github.com/tywkeene/go-tracker/cmd/client/version.CommitHash=$COMMIT" \
    github.com/tywkeene/go-tracker/cmd/client 
