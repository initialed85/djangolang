#!/bin/bash

set -e

target=${1:-./pkg/some_db}

rm -frv "${target}-temp" >/dev/null 2>&1 || true

# shellcheck disable=SC2068
DJANGOLANG_DEBUG=${DJANGOLANG_DEBUG:-0} POSTGRES_DB=${POSTGRES_DB:-some_db} POSTGRES_PASSWORD=${POSTGRES_PASSWORD:-some-password} go run ./cmd template "${target}"
