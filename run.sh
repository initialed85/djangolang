#!/bin/bash

set -e

if [[ "${CI}" == "true" ]]; then
    POSTGRES_PORT="0"
    REDIS_PORT="0"
    SWAGGER_PORT="0"

    export POSTGRES_PORT
    export REDIS_PORT
    export SWAGGER_PORT
fi

PORT="${PORT:-7070}"
DJANGOLANG_DEBUG="${DJANGOLANG_DEBUG:-0}"
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
    shift

    ./env.sh "${@}"
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

    go test -race -v -failfast -count=1 ./pkg/template
    ;;

"test")
    while ! docker compose ps -a | grep post-migrate | grep 'Exited (0)' >/dev/null 2>&1; do
        sleep 0.1
    done

    shift

    find . -type f -name '*.*' | grep -v '/model_generated/' | grep -v '/model_generated_from_schema/' | entr -n -r -cc -s "DJANGOLANG_NODE_NAME=test-clean go test -race -v -failfast -count=1 ./pkg/template ./pkg/schema && DJANGOLANG_NODE_NAME=test-clean go test -race -v -failfast -count=1 ${*:-./pkg/model_generated_test}; echo -e '\n(done)'"
    ;;

"test-ci")
    while ! docker compose ps -a | grep post-migrate | grep 'Exited (0)' >/dev/null 2>&1; do
        sleep 0.1
    done

    shift

    docker compose exec -e DJANGOLANG_NODE_NAME=test-ci test go test -race -v -failfast -count=1 ./pkg/template
    docker compose exec -e DJANGOLANG_NODE_NAME=test-ci test go test -race -v -failfast -count=1 ./...

    echo -e '\n(done)'
    ;;

"test-specific")
    while ! docker compose ps -a | grep post-migrate | grep 'Exited (0)' >/dev/null 2>&1; do
        sleep 0.1
    done

    shift

    find . -type f -name '*.*' | grep -v '/model_generated/' | grep -v '/model_generated_from_schema/' | entr -n -r -cc -s "DJANGOLANG_NODE_NAME=test-specific go test -race -v -failfast -count=1 ${*}; echo -e '\n(done)'"
    ;;

"test-specific-no-clean")
    while ! docker compose ps -a | grep post-migrate | grep 'Exited (0)' >/dev/null 2>&1; do
        sleep 0.1
    done

    shift

    find . -type f -name '*.*' | grep -v '/model_generated/' | grep -v '/model_generated_from_schema/' | entr -n -r -cc -s "DJANGOLANG_NODE_NAME=test-specific-no-clean go test -race -v -failfast -count=1 ${*}; echo -e '\n(done)'"
    ;;

"serve")
    while ! docker compose ps -a | grep post-migrate | grep 'Exited (0)' >/dev/null 2>&1; do
        sleep 0.1
    done

    find ./pkg/model_generated* -type f -name '*.go' | entr -n -r -cc -s "DJANGOLANG_PROFILE=${DJANGOLANG_PROFILE:-0} DJANGOLANG_NODE_NAME=${DJANGOLANG_NODE_NAME:-serve} go run ./pkg/model_generated/cmd/ serve"
    ;;

"serve-schema")
    while ! docker compose ps -a | grep post-migrate | grep 'Exited (0)' >/dev/null 2>&1; do
        sleep 0.1
    done

    find ./pkg/model_generated_from_schema* -type f -name '*.go' | entr -n -r -cc -s "DJANGOLANG_PROFILE=${DJANGOLANG_PROFILE:-0} DJANGOLANG_NODE_NAME=${DJANGOLANG_NODE_NAME:-serve} POSTGRES_SCHEMA=test go run ./pkg/model_generated_from_schema/cmd/ serve"
    ;;

"stream")
    while ! docker compose ps -a | grep post-migrate | grep 'Exited (0)' >/dev/null 2>&1; do
        sleep 0.1
    done

    find ./pkg/model_generated* -type f -name '*.go' | entr -n -r -cc -s "while true; do unbuffer websocat ws://localhost:${PORT}/__stream | jq; done"
    ;;

"cli")
    shift

    DJANGOLANG_PROFILE=${DJANGOLANG_PROFILE:-0} DJANGOLANG_NODE_NAME=${DJANGOLANG_NODE_NAME:-serve} go run ./pkg/model_generated/cmd/ "${*}"
    ;;

*)
    echo "error: unrecognized command: ${1}"
    exit 1
    ;;
esac
