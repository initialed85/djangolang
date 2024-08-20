#!/bin/bash

set -e

if [[ "${1}" == "" ]]; then
    echo "error: first argument must be command (one of 'env', 'test', 'serve' or 'stream')"
    exit 1
fi

PORT="${PORT:-7070}"
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

"migrate")
    while ! docker compose ps -a | grep post-migrate | grep 'Exited (0)' >/dev/null 2>&1; do
        sleep 0.1
    done

    shift

    # shellcheck disable=SC2068
    migrate -source file://./database/migrations -database postgres://postgres:some-password@localhost:5432/some_db?sslmode=disable ${@}
    ;;

"template")
    while ! docker compose ps -a | grep post-migrate | grep 'Exited (0)' >/dev/null 2>&1; do
        sleep 0.1
    done

    go test -failfast -count=1 ./pkg/template
    ;;

"test")
    while ! docker compose ps -a | grep post-migrate | grep 'Exited (0)' >/dev/null 2>&1; do
        sleep 0.1
    done

    shift

    find . -type f -name '*.*' | grep -v '/model_generated/' | entr -n -r -cc -s "DJANGOLANG_NODE_NAME=test go test -failfast -count=1 ./pkg/helpers ./pkg/types ./pkg/query ./pkg/template ./pkg/openapi && DJANGOLANG_NODE_NAME=test go test -v -failfast -count=1 ./pkg/model_generated_test ${*}"
    ;;

"test-clean")
    while ! docker compose ps -a | grep post-migrate | grep 'Exited (0)' >/dev/null 2>&1; do
        sleep 0.1
    done

    shift

    find . -type f -name '*.*' | grep -v '/model_generated/' | entr -n -r -cc -s "PAGER=cat PGPASSWORD=some-password psql -h localhost -p 5432 -U postgres some_db -c 'TRUNCATE TABLE physical_things CASCADE;' && DJANGOLANG_NODE_NAME=test-clean go test -failfast -count=1 ./pkg/helpers ./pkg/types ./pkg/query ./pkg/template ./pkg/openapi && DJANGOLANG_NODE_NAME=test go test -v -failfast -count=1 ./pkg/model_generated_test ${*}"
    ;;

"serve")
    while ! docker compose ps -a | grep post-migrate | grep 'Exited (0)' >/dev/null 2>&1; do
        sleep 0.1
    done

    find ./pkg/model_generated -type f -name '*.go' | entr -n -r -cc -s "DJANGOLANG_NODE_NAME=${DJANGOLANG_NODE_NAME:-serve} go run ./pkg/model_generated/cmd/ serve"
    ;;

"stream")
    while ! docker compose ps -a | grep post-migrate | grep 'Exited (0)' >/dev/null 2>&1; do
        sleep 0.1
    done

    find ./pkg/model_generated -type f -name '*.go' | entr -n -r -cc -s "while true; do unbuffer websocat ws://localhost:${PORT}/__stream | jq; done"
    ;;

*)
    echo "error: unrecognized command: ${1}"
    exit 1
    ;;
esac
