#!/bin/bash

set -e

if ! command -v atlas >/dev/null 2>&1; then
    curl -sSf https://atlasgo.sh | sh
fi

command="${1}"

if [[ "${command}" == "" ]]; then
    echo "error: first argument must be command"
    exit 1
fi

atlas_dir="file://docker/postgres/migrations"
atlas_to="file://docker/postgres/schema.sql"
atlas_dev_image_1="djangolang-postgres:latest"
atlas_dev_image_2="postgres:16-bullseye-not-really"
atlas_dev_image_3="postgres/16-bullseye-not-really"
atlas_dev_url="docker://${atlas_dev_image_3}/djangolang?search_path=public"
atlas_format='{{ sql . "  " }}'
atlas_url="postgres://postgres:Password1@localhost:5432/djangolang?sslmode=disable"

case "${command}" in

"diff")
    docker image tag "${atlas_dev_image_1}" "${atlas_dev_image_2}"

    atlas migrate diff initial \
        --dir "${atlas_dir}" \
        --to "${atlas_to}" \
        --dev-url "${atlas_dev_url}" \
        --format "${atlas_format}"
    ;;

"create")
    atlas migrate new --dir "${atlas_dir}"
    exit 1
    ;;

"make")
    echo "error: TODO"
    exit 1
    ;;

"up")
    atlas migrate apply \
        --dir "${atlas_dir}" \
        --url "${atlas_url}"
    ;;

"down")
    echo "error: TODO"
    exit 1
    ;;

*)
    echo "error: unknown command ${command}"
    exit 1
    ;;
esac
