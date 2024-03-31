#!/bin/bash

set -e

# shellcheck disable=SC2068
PORT=${PORT:-8000} DJANGOLANG_DEBUG=${DJANGOLANG_DEBUG:-0} POSTGRES_DB=some_db POSTGRES_PASSWORD=some-password go run ./cmd serve
