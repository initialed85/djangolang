# djangolang

# status: under development, barely works

This started out as an attempt to write some Go that exposes a Postgres database in a HTTP API that feels like [Django Rest Framework](https://www.django-rest-framework.org/).

It's still gonna be the first part, it'll probably deviate from the DRF piece though.

## Goals

-   Point it at your DB
-   It introspects it and generates some code
-   Generated code includes:
    -   Structs for your tables
    -   Generated SQL helpers
    -   Generated CRUD endpoints that use the SQL helpers

The SQL helpers aren't really an ORM per se, but I guess they're in that direction.

## Stretch Goals

-   Publish events somewhere using Postgres logical replication
-   Intelligent cache layer
    -   Use Postgres logical replication to set / reset cache items
    -   Provides an interface over the cached data using the generated Go structs
-   OpenAPI spec generation
-   Client generation

## Non Goals

I'm not trying to write a rich ORM, I just want to make it easy to sling your Postgres DB around in Go- the best language for writing complicated SQL is SQL.

## Status

-   Foundation
    -   Done-ish
        -   Database introspection works (not all types catered for)
        -   Struct generation works
            -   TODO: Make the code generation smarter and less cumbersome (editing it is challenging)
-   Create
    -   Done-ish
        -   Select SQL helper function generation works
            -   No cascading foreign object upsert (not sure if TODO, it's a bit ORM-y)
        -   POST endpoint generation works
-   Read
    -   Done-ish
        -   Select SQL helper function generation works
            -   Child-to-parent foreign object loading works
                -   TODO: Fix up the recursion protection so that it doesn't bail out too early in the case of multiple foreign keys to the same table
        -   GET (list) endpoint generation works
            -   Supports a degree of filter, order, limit and offset
-   Update
    -   Done-ish
        -   Select SQL helper function generation works
            -   No cascading foreign object upsert (not sure if TODO, it's a bit ORM-y)
        -   PATCH endpoint generation works
-   Delete
    -   Done-ish
        -   Delete SQL helper function generation works
            -   No cascading foreign object delete (not sure if TODO, it's a bit ORM-y)
        -   DELETE endpoint generation works
-   Stream (i.e. Postgres logical replication events)
    -   In progress
        -   Consumes a Postgres logical replication slot
        -   Attempts to select (with child-to-parent foreign object loading) the inserted / updated row
            -   TODO: Populate the root object with the data from the change, not the selected data (that may be newer)
        -   Produces a "Change" struct (that is right now only printed)
        -   Doesn't decode all the data types it needs to
            -   TODO: Need a fairly detailed test schema and associated test to hit everything here I think
            -   TODO: Generics or code generation?
        -   Doesn't do anything other than print
            -   TODO: Spit the Change structs into a channel
-   Caching layer
    -   Not yet started
-   OpenAPI spec generation
    -   Not yet started
-   Client generation
    -   Not yet started

## Usage

```shell
# generate the code for your database
POSTGRES_HOST=localhost POSTGRES_USER=postgres POSTGRES_PASSWORD=some-password POSTGRES_DB=some_db go run ./cmd/main.go template

# run the server
PORT=8080 POSTGRES_HOST=localhost POSTGRES_USER=postgres POSTGRES_PASSWORD=some-password POSTGRES_DB=some_db go run ./cmd/main.go serve

# you can order and limit and offset (e.g. for pagination)
curl -s 'http://localhost:8080/event?order_by=start_timestamp&order=desc&limit=10&offset=10' | jq

# you can filter for an exact match (result is still a list of length 1)
curl -s 'http://localhost:8080/event?id=46709' | jq

# you can filter for a range
curl -s 'http://localhost:8080/detection?id__gte=14790680&id__lte=14790682' | jq

# you can filter for a list
curl -s 'http://localhost:8080/detection?id__in=14790680,14790681,14790682' | jq
```

## Advanced Usage

If you want to build on top of this, there are some useful functions that are generated alongside your generated code:

-   `GetTableByName() (map[string]*introspect.Table, error)`
    -   Dump out structs that describe the schema
-   `GetRawTableByName() []byte`
    -   Dump out the raw JSON that describes the schema

## Examples

Check out the tests for some example usage as code:

-   [./test/sql_helpers_test.go](test/sql_helpers_test.go)
-   [./test/endpoints_test.go](test/endpoints_test.go)

## Development

```shell
#
# shell 1
#

docker compose up -d && \
docker compose exec -it postgres psql -U postgres -c 'ALTER SYSTEM SET wal_level = logical;' && \
docker compose restart postgres && \
docker compose up -d && \
docker compose logs -f -t ; \
docker compose down --remove-orphans --volumes

#
# shell 2
#

# this just reads the database
./introspect.sh

# this reads the database and generates the code
./template.sh

# this runs a server on top of the generated code
./server.sh

# this cleans out the database, generates the code and executes some tests against both the generated code and the server
# note: ignore the warnings that the stream is spitting out about failure to select objects- that's ultimately just a
# byproduct of the stream piece being WIP
./test.sh
```
