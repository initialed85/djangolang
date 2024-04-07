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
-   Intelligent Redis cache layer
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
-   Create
    -   Done-ish
        -   Select SQL helper function generation works
            -   No cascading foreign object upsert (not sure if TODO, it's a bit ORM-y)
        -   POST endpoint generation works
-   Read
    -   Done-ish
        -   Select SQL helper function generation works
            -   Child-to-parent foreign object loading works
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
-   Postgres logical replication events
    -   In progress
        -   Doesn't decode all the data types it needs to
        -   Doesn't do anything other than print to the console
-   Redis Caching layer
    -   Not yet started
-   OpenAPI spec generation
    -   Not yet started
-   Client generation
    -   Not yet started

### TODO

-   Fix up handling for recursive schemas
    -   Right now we use r.Context() for loop prevention based on object; but it needs to be (object, field) to ensure that we don't exit recursion too early (e.g in the case that a table has multiple foreign keys to the same table)
-   Constraints / Upsert
-   Refactor the cumbersome string-based templating into proper struct-based templating, or something better anyway
-   Write more tests

## Usage

```shell
# generate the stubs for your database
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

If you want to build on top of this, there are some useful functions that are generated alongside your stubs:

-   `GetTableByName() (map[string]*introspect.Table, error)`
    -   Dump out structs that describe the schema
-   `GetRawTableByName() []byte`
    -   Dump out the raw JSON that describes the schema

## Examples

Check out the tests for some example usage as code:

-   [./test/sql_helpers_test.go](test/sql_helpers_test.go)
-   [./test/endpoints_test.go](test/endpoints_test.go)
