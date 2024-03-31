#!/bin/bash

set -e

function cleanup() {
    if [[ "${SKIP_CLEANUP}" == "1" ]]; then
        exit 0
    fi

    docker compose down --remove-orphans --volumes || true
}
trap cleanup exit

if docker compose ps | grep postgres | grep Up | grep healthy >/dev/null 2>&1; then
    SKIP_CLEANUP=1
else
    docker compose up -d
fi

if [[ "${SKIP_TEMPLATE}" != "1" ]]; then
    DJANGOLANG_DEBUG=${DJANGOLANG_DEBUG:-1} POSTGRES_DB=some_db POSTGRES_PASSWORD=some-password ./template.sh
fi

PORT=8000 DJANGOLANG_DEBUG=${DJANGOLANG_DEBUG:-1} POSTGRES_DB=some_db POSTGRES_PASSWORD=some-password go test -v -count=1 -failfast ./...
