#!/bin/bash

set -e

cd database

atlas migrate diff \
    --to file://schema.sql \
    --dev-url "docker://postgis/14-3.4/some_db?search_path=public" \
    --format '{{ sql . "  " }}'
