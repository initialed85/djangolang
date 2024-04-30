#!/bin/bash

set -e

function teardown() {
    docker compose down --remove-orphans --volumes
}

trap teardown exit

docker compose pull

docker compose build

if ! docker compose up -d; then
    docker compose logs -t
    echo "error: docker compose up failed; scroll up for logs"
    exit 1
fi

docker compose exec -it postgres psql -U postgres -c 'ALTER SYSTEM SET wal_level = logical;'

docker compose restart postgres

docker compose up -d

docker compose logs -f -t
