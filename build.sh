#!/bin/bash

set -e

pushd "$(dirname -- "${BASH_SOURCE[0]}")" >/dev/null 2>&1

function teardown() {
    popd >/dev/null 2>&1 || true
}
trap teardown exit

mkdir -p ./bin

if test -e ./bin/djangolang; then
    rm -f ./bin/djangolang
fi

go build -o ./bin/djangolang ./cmd
