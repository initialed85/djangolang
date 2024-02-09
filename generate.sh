#!/bin/bash

set -e

pushd "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)" >/dev/null 2>&1
function cleanup() {
    popd >/dev/null 2>&1 || true
}
trap cleanup EXIT

if ! command -v xo >/dev/null 2>&1; then
    go install github.com/xo/xo@v1.0.1
fi

POSTGRES_USER="${POSTGRES_USER:-postgres}"
POSTGRES_PASSWORD="${POSTGRES_PASSWORD?POSTGRES_PASSWORD env var unset}"
POSTGRES_HOST="${POSTGRES_HOST:-host.docker.internal}"
POSTGRES_PORT="${POSTGRES_PORT:-5432}"
POSTGRES_DB="${POSTGRES_DB?POSTGRES_DB env var unset}"

# rm -fr ./pkg/models
# mkdir -p ./pkg/models
# # shellcheck disable=SC2068
# xo schema --out ./pkg/models ${@} "postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable"

go mod tidy
go get ./...
go fmt ./...
go vet ./...
