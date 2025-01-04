#!/bin/bash

set -e

go fmt ./...
goimports -w .
go vet ./...
staticcheck ./...

echo "(done)"
