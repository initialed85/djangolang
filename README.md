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

## Non Goals

I'm not trying to write a rich ORM, I just want to make it easy to sling your Postgres DB around in Go- the best language for writing complicated SQL is SQL.

## Tasks

-   Foundation
    -   In progress
        -   Database introspection works (not all types catered for)
        -   Struct generation works
-   Create
    -   Not yet started
-   Read
    -   In progress
        -   Select SQL helper function generation works
            -   Child-to-parent foreign object loading works
        -   GET (list) endpoint generation works
            -   Supports a degree of filter, order, limit and offset
-   Update
    -   Not yet started
-   Delete
    -   Not yet started

### TODO

-   Insert / Update / Delete
-   Constraints
-   Refactor the cumbersome string-based templating into proper struct-based templating
-   Write literally any tests

## Usage

```shell
# generate the stubs for your database
POSTGRES_HOST=localhost POSTGRES_USER=postgres POSTGRES_PASSWORD=some-password POSTGRES_DB=some-db go run ./cmd/main.go template

# run the server
PORT=8080 POSTGRES_HOST=localhost POSTGRES_USER=postgres POSTGRES_PASSWORD=some-password POSTGRES_DB=some-db go run ./cmd/main.go server

# you can order and limit and offset (e.g. for pagination)
curl -s 'http://localhost:8080/event?order_by=start_timestamp&order=desc&limit=10&offset=10' | jq

# you can filter for an exact match (result is still a list of length 1)
curl -s 'http://localhost:8080/event?id=46709' | jq

# you can filter for a range
curl -s 'http://localhost:8080/detection?id__gte=14790680&id__lte=14790682' | jq

# you can filter for a list
curl -s 'http://localhost:8080/detection?id__in=14790680,14790681,14790682' | jq
```

## Advance Usage

If you want to build on top of this, there are some useful functions that are generated alongside your stubs:

-   `GetTableByName() (map[string]*introspect.Table, error)`
    -   Dump out structs that describe the schema
-   `GetRawTableByName() []byte`
    -   Dump out the raw JSON that describes the schema

## Example

Ref.: [./pkg/example/example.go](pkg/example/example.go)

```go
package example

import (
	"context"
	"fmt"
	"log"

	"github.com/initialed85/djangolang/pkg/helpers"
	"github.com/initialed85/djangolang/pkg/template/templates" // this is where my templates were generated- yours will be different
)

func Run(ctx context.Context) error {
	db, err := helpers.GetDBFromEnvironment(ctx) // e.g. POSTGRES_DB, POSTGRES_PASSWORD etc
	if err != nil {
		return err
	}
	defer func() {
		_ = db.Close()
	}()

	events, err := templates.SelectEvents(
		ctx,
		db,
		templates.EventTransformedColumns, // you'll mostly want this (all columns + any relevant transforms, e.g. PostGIS Point -> GeoJSON map[string]any)
		helpers.Ptr(fmt.Sprintf("%v DESC", templates.EventStartTimestampColumn)), // order by
		helpers.Ptr(10), // limit
		helpers.Ptr(10), // ofset
		fmt.Sprintf("%v > 1000", templates.EventIDColumn),                // implicitly AND'd together
		fmt.Sprintf("%v < 2000", templates.EventIDColumn),                // implicitly AND'd together
		fmt.Sprintf("%v IN (1, 2)", templates.EventSourceCameraIDColumn), // implicitly AND'd together
	)
	if err != nil {
		log.Fatal(err)
	}

	for _, event := range events {
		log.Printf("%#+v", event)
	}

	return nil
}
```
