#!/bin/bash

set -e

if [[ "${1}" == "" ]]; then
    echo "error: first argument must be command (one of 'env', 'test', 'serve' or 'stream')"
    exit 1
fi

DJANGOLANG_DEBUG="${DJANGOLANG_DEBUG:-1}"
DJANGOLANG_SET_REPLICA_IDENTITY="${DJANGOLANG_SET_REPLICA_IDENTITY:-full}"
REDIS_URL="${REDIS_URL:-redis://default:some-password@localhost:6379}"
POSTGRES_DB="${POSTGRES_DB:-some_db}"
POSTGRES_PASSWORD="${POSTGRES_PASSWORD:-some-password}"

export DJANGOLANG_DEBUG
export DJANGOLANG_SET_REPLICA_IDENTITY
export REDIS_URL
export POSTGRES_DB
export POSTGRES_PASSWORD

case "${1}" in

"env")
    ./env.sh
    ;;

"test")
    find . -type f -name '*.*' | grep -v '/model_generated/' | entr -n -r -cc -s "go test -failfast -count=1 ./pkg/types ./pkg/query ./pkg/template ./pkg/openapi && go test -v -failfast -count=1 ./pkg/model_generated_test"
    ;;

"serve")
    find ./pkg/model_generated -type f -name '*.go' | entr -n -r -cc -s "go run ./pkg/model_generated/cmd/ serve"
    ;;

"stream")
    find ./pkg/model_generated -type f -name '*.go' | entr -n -r -cc -s "while true; do unbuffer websocat ws://localhost:7070/__stream | jq; done"
    ;;

*)
    echo "error: unrecognized command: ${1}"
    exit 1
    ;;
esac
