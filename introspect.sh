#!/bin/bash

set -e

DJANGOLANG_DEBUG=${DJANGOLANG_DEBUG:-0} POSTGRES_DB=some_db POSTGRES_PASSWORD=some-password go run ./cmd introspect
