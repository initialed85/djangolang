#!/bin/bash

set -e

go vet ./...

staticcheck ./...

if test -e ./bin; then
    rm -frv ./bin
fi

mkdir -p ./bin

CGO_ENABLED=0 go build -o ./bin -trimpath ./cmd
