# djangolang

# status: under development, sort of works

An opinionated framework that's trying to make it easy to turn any (with various caveats) Postgres database into a HTTP RESTful API server and a WebSocket CDC server,
using Redis for caching, supporting pluggable middleware (for things like authentication / authorization), pluggable post-mutation actions and custom endpoints.

## Tasks

The TODOs aren't in any sensible order, but the DONEs / WIP are in the order completed

- [DONE] Database schema introspection
- [DONE] Code generation for SQL helpers
  - [DONE] Support loading of directly-related foreign key objects (i.e parents, more or less)
- [DONE] Code generation for CRUD endpoints
  - [DONE] GET list
  - [DONE] GET item
  - [DONE] PUT (replace)
  - [DONE] PATCH (update)
  - [DONE] DELETE
- [DONE] Consume logical replication stream and produce CDC events
- [DONE] Code generation for generic HTTP server
  - [DONE] Include WebSocket CDC server
- [DONE] Use Redis for a cache in generated CRUD endpoints
  - [DONE] Smart-ish (table-level) cache invalidation using the CDC events
- [DONE] Support soft-deletion
  - [TODO] Make soft-deletion less implicit and more configurable
- [DONE] OpenAPI JSON / YAML generation (to support generation of third party clients / servers)
- [DONE] Support loading of indirectly referenced-by foreign key objects (i.e. children, more or less)
  - [DONE] Make and then fix a big recursive mess (really make sure we can't go infinite)
  - [DONE] Fix up cache invalidation now that it has to flow both ways through the object graph
  - [DONE] Have an in-memory query cache (lifetime-scoped to a root query execution) now that we're touching so much of the object graph per root query
- [DONE] Some sort of query parameter to opt-out of the directly-related / referenced-by foreign object loading (`?depth=(int)`)
- [DONE] Figure out how to sensibly serialize / deserialize point / polygon / pointz
  - [TODO] Find a way to not rely on an `ST_PointZ` for `pointz` (I think nobody has a working Go lib that handles the binary type?)
- [TODO] Redo the type mapping stuff again
- [TODO] Fix all the tests I commented out (I was too lazy to fix the complicated structures)
- [TODO] Support create-or-update endpoints; thoughts:
  - Probably a special URL path
  - Would be nice to be able to create or update reverse-relationship children at the same time
  - Don't ask for an opinionated schema re: constraints- a specified primary key = update, no primary key = create
  - Think about some way to use the object content to potentially resolve a primary key
- [TODO] Middleware / hooks; all the loose thoughts:
  - Authentication HTTP middleware
  - Support for object-level middleware
  - Authorization object-level middleware
  - Before / after middleware pattern?
  - Before / after hooks? Triggered by stream?
- [TODO] Explicit tests for the cache
- [TODO] Mechanism to optionally partition cache by authentication (i.e. public vs private endpoint caching)
- [TODO] Better tests for OpenAPI schema generation
- [TODO] Ephemeral pub-sub using Redis w/ WebSocket client
- [TODO] Support for custom endpoints
- [TODO] Make the errors more readable (probably populate an `objects: []any` field)
- [TODO] Move some of the configuration injection further out (environment variables at the outside, but just function parameters further in)
- [TODO] Replace the various "this should probably be configurable" TODOs with mechanisms to configure
- [TODO] Support more Postgres data types as they come up, maybe redo the type-mapping stuff to be clearer in terms of to/from request / db / stream / intermediate
  - [TODO] Maybe just use `jackc/pgx` everywhere in the hope it'll save some type-juggling
- [TODO] Better support for recursive schemas (in the case that they cause a graph cycle)
- [TODO] Support for views
- [DONE] Fix up the various templating shortcuts I've taken that cause `staticcheck` warnings (e.g. `unnecessary use of fmt.Sprintf`)
- [TODO] Think about how to do hot-reloading on schema changes (is this mostly an infra problem? Not sure)
- [TODO] Document all the features
- [TODO] Cleaner handling for when Redis isn't available
- [TODO] Some sort of latest-by-parent endpoint variation that implicitly expects a `timestamp` or `time` column on the object in question

## Usage for prod

See [initialed85/camry](https://github.com/initialed85/camry) for a practical usage of Djangolang.

## Usage for dev

### Prerequisites

- Linux / MacOS (you can probably build and run this thing on Windows, but you're on your own)
- Docker and Docker Compose
- Go 1.21+
- curl
- entr
- websocat
- jq
- yq
- golang-migrate

### Workflow

```shell
# shell 1 - spin up the dependencies
./run.sh env

# shell 2 - run the generated HTTP API and WebSocket CDC server
./run.sh server

# shell 3 - consume from the WebSocket CDC stream
./run.sh stream

# shell 4 - run the templating tests and then the integration tests for the generated code (causes some changes to be seen at the WebSocket CDC stream)
./run.sh test-clean
```

Everything should restart automatically when the code changes, testing first any templating aspects and then (if that works) testing the actual behaviours; shells 2 and 3 are
mostly just a smoke test.

### Usage

- Interact with the generated API at [http://localhost:7070](http://localhost:7070)
  - e.g.:
    - `curl -X POST http://localhost:7070/logical-things -d '[{"name": "Name1", "type": "Type1"}, {"name": "Name2", "type": "Type2"}]' | jq`
    - `curl 'http://localhost:7070/logical-things?name__in=Name1,Name2' | jq`
    - `curl 'http://localhost:7070/logical-things?name__ilike=ame' | jq`
- Consume the WebSocket CDC stream at [http://localhost:7070/\_\_stream](http://localhost:7070/__stream)
  - e.g.: `websocat ws://localhost:7070/__stream`
- See the generated OpenAPI schema at [http://localhost:7070/openapi.json](http://localhost:7070/openapi.json) and [http://localhost:7070/openapi.yaml](http://localhost:7070/openapi.yaml)
  - `curl 'http://localhost:7070/openapi.json' | jq`
  - `curl 'http://localhost:7070/openapi.yaml' | yq`
- Explore generated OpenAPI schema using Swagger at [http://localhost:7071](http://localhost:7071)

### Notes

#### Database row representations

Database rows in Djangolang are represented as either **objects** or **items**.

Objects are the generated Go structs that represent database rows at a high level, with a strict type definition that matches the table and types that are intuitive to work with.

Items are simply `map[string]any` (where the key is the column name), but confusingly there are 3 patterns around value types, pending which part of Djangolang the items are pertinent to:

- Query helper functions
  - Types are determined by the produced output (and required input) for SQL queries as handled by:
    - `jmoiron/sqlx`
    - `lib/pq`
    - `database/sql`
- Logical replication stream
  - Types are determined by the product output for decoded logical replication messagesas handled by:
    - `jackc/pgx/v5/pgconn`
    - `jackc/pgx/v5/pgtype`
- Generated endpoint functions
  - Types are decided by what's most intuitive to work with (once it's marshaled to JSON)

So you can see there are 4 different representations (including the Djangolang object) of a database row each with slightly different types (for the same actual column).

The main challenges arise from the not-very-intuitive types that the two database-facing library suites unmarshal into (and indeed the fact that these don't entirely agree
between the `database/sql` query helper function side and the `jackc/pgx` logical replication stream side).

At some point I'll aim to change over to just using `jackx/pgx` for everything to cut down on at least some of this confusion, but we can't get away from needing a different
set of types for the Objects and a different set of types for the endpoint-flavoured items anyway.

##### Type conversions

The type conversion piece is a bit of a mess at the moment- it works largely as intended for the consumer, but the code is not anywhere near as clean as I'd like.

If we follow a row through from a Select into a GET response, it's something like this:

- A generated `SelectABC` helper function (where `ABC` is a generated Djangolang object struct) is called (expected to return `[]*ABC`)
- The `pkg/query/query.go::Select` SQL helper function returns items as `[]map[string]any`
  - These are in the native format from the `jmoiron/sqlx` / `lib/pq` / `database/sql` libs suite
- An appropriate empty Djangolang object struct is created for each item and `.FromItem([]map[string]any)` is called to set the struct's properties, using the various `types.ParseXYZ(any)` functions
  - An implication here is that the `types.ParseXYZ(any)` function must know how to convert from a `lib/pq` type to a Djangolang object-appropriate type
- The Djangolang object struct is marshaled with plain old `json.Marshal` and this is returned in the response
  - An implication here is the Djangolang object-appropriate types must also be Djangolang endpoint-appropriate

And if we do the same for a POST to an Insert, it's something like this:

- A response body is unmarshaled into `[]map[string]any`
  - These are expected to be in the Djangolang endpoint-appropriate types
- An appropriate empty Djangolang object struct is created for each item and `.FromItem([]map[string]any)` is called to set the struct's properties, using the various `types.ParseXYZ(any)` functions
  - An implication here is that the `types.ParseXYZ(any)` function must know how to convert from a Djangolang endpoint-appropriate type to a Djangolang object-appropriate type
- Now validated via the Djangolang object, `.Insert()` is called and a `map[string]any` is created using the various `types.FormatXYZ(any)` functions
  - An implication here is that the `types.FormatXYZ(any)` function must know how to convert from a Djangolang object-appropriate type to a `lib/pq` type
- The `pkg/query/query.go::Insert` SQL helper function columns as `[]string` and values `any` which is just our `map[string]any` split into two parameters

In the background, because we've made a mutation to the database, we've triggered a logical replication stream event; that goes something like this:

- `jackc/pgx/v5/pgconn` receives a logical replication stream message
  - It's decoded as an `INSERT`
- `jackc/pgx/v5/pgtype` turns that into a `map[string]any` item
  - The types are not deemed intuitive at this point
- An appropriate empty Djangolang object struct is created for the item and `.FromItem([]map[string]any)` is called to set the struct's properties, using the various `types.ParseXYZ(any)` functions
  - An implication here is that the `types.ParseXYZ(any)` function must know how to convert from a `jackc/pgx/v5/pgtype` type to a Djangolang object-appropriate type

So boiling down the above:

- `types.ParseXYZ(any) (any, error)`
  - Used for keys of query-flavoured `map[string]any` -> Djangolang object
  - Used for keys of stream-flavoured `map[string]any` -> Djangolang object
- `json.Marshal(any) (any, error)`
  - Used for Djangolang object -> endpoint-flavoured `map[string]any` as `[]byte`
- `json.Unmarshal([]byte, any) error`
  - Used for `[]byte` -> endpoint-flavoured `map[string]any`
- `(*SomeObject).FromItem(map[string]any) (*SomeObject, error)`
  - Used for endpoint-flavoured `map[string]any` -> Djangolang object
- `types.FormatXYZ(any) (any, error)`
  - Used for properties of Djangolang object -> query-flavoured `map[string]any`
