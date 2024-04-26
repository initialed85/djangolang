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
    -   Tooling for HTTP and WebSocket server for CRUD endpoints / change event stream
    -   Hooks for Authentication and Authorization
-   Client generation
    -   Generated TypeScript interface types / clients for the CRUD endpoints
    -   Generated WebSocket client for change event stream
    -   Maybe generated Go HTTP CRUD client because you love microservices and don't want to connect to your database twice or something
-   Intelligent cache layer
    -   Use change event stream to set / reset cache items
    -   Provides an interface over the cached data using the generated Go structs
-   OpenAPI spec generation (or some kind of IR)

The SQL helpers aren't really an ORM per se, but I guess they're in that direction.

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
        -   TODO: GET item endpoint generation
-   Update
    -   Done-ish
        -   Select SQL helper function generation works
            -   No cascading foreign object upsert (not sure if TODO, it's a bit ORM-y)
        -   PUT endpoint generation works
        -   TODO: PATCH endpoint generation
-   Delete
    -   Done-ish
        -   Delete SQL helper function generation works
            -   No cascading foreign object delete (not sure if TODO, it's a bit ORM-y)
        -   DELETE endpoint generation works
-   Stream (i.e. Postgres logical replication events)
    -   In progress
        -   Consumes a Postgres logical replication slot
        -   Attempts to select (with child-to-parent foreign object loading) the inserted / updated row
            -   TODO: Provide the actual shallow changed object as well as a nested refetched one
        -   Probably doesn't decode all the data types it needs to
            -   TODO: Need a fairly detailed test schema and associated test to hit everything here I think
            -   TODO: Generics or code generation?
        -   Produces a "Change" struct that gets sent to a channel
        -   Provides a WebSocket server that publishes "Change" structs as JSON to subscribers
            -   TODO: Honour table name filtering
-   Client generation
    -   TODO
        -   TypeScript interfaces types
        -   SWR tooling / RTK Query tooling
        -   WebSocket client for the change event stream
        -   Probably some SWR tooling and RTK Query tooling for TypeScript
-   Authentication layer
    -   TODO
        -   Pattern for injecting an `Authenticate(token string) error` function
        -   Wire it throughout the generated endpoints and at WebSocket upgrade for the change event stream
        -   `/__authenticate` endpoint that extracts a token from `Authorization: Bearer (.*)` and feeds it to the above for authentication checks
-   Authorization layer
    -   TODO
        -   Pattern for injecting a `GetAuthorizedFilterJoins[T any](token string, tableName string) (string, error)` function
            -   Wire it through the generated endpoints
        -   Pattern for injecting a `IsAuthorized[T any](token string, tableName string, foreignKeyColumn string, foreignKeyValue T)` function
            -   Wire it in as a filter for the produced change event stream messages
-   Caching layer
    -   TODO
-   OpenAPI spec generation
    -   TODO

## What to do next

-   Completely refactor the templating stuff so that I can move faster
-   Go back to brass tacks with database serialization / deserialization
    -   Maybe an intermediate object? It's getting messy.
-   Get closer to a drop-in-replacement for DRF use cases (at least interface / behaviourally):
    -   Fix the hacked-in recursion protection
    -   Single-item (e.g. `GET {table}/{primaryKeyValue}`) endpoint
    -   Partial update (e.g. `PATCH {table}/{primaryKeyValue}`) endpoint
        -   Need to generate an `OptionalTable` struct where every field is a pointer so it's clear which fields to change
-   Get some client generation working
    -   The idea being that by this point, I can start refactoring some of my personal projects to use `djangolang` to learn more
-   Do a few things to be at all feasible for any sort of production system:
    -   At least Authentication
    -   Probably also Authorization (though internal / non-user-facing services might not need it)
-   Tackle the scale problem
    -   Ensure everything shards nicely
    -   Caching layer
-   Support future external integrations
    -   OpenAPI spec generation

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
