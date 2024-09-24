#!/bin/bash

set -e

if [[ "${1}" == "down" ]]; then
    docker compose down --remove-orphans --volumes
    exit 0
fi

if [[ "${1}" != "up" ]]; then
    function teardown() {
        docker compose down --remove-orphans --volumes
    }

    trap teardown exit
fi

docker compose pull

docker compose build

if ! docker compose up -d; then
    docker compose logs -t
    echo "error: docker compose up failed; scroll up for logs"
    exit 1
fi

docker compose exec -it postgres psql -U postgres -c 'ALTER SYSTEM SET wal_level = logical;'

docker compose exec -it postgres psql -U postgres -c "ALTER DATABASE some_db SET log_statement = 'all';"

docker compose restart postgres

docker compose up -d

if [[ "${1}" != "up" ]]; then
    docker compose logs -f -t
fi
