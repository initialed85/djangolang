# djangolang

# status: under development, sort of works

An opinionated framework that's trying to make it easy to turn any (with various caveats) Postgres database into a HTTP RESTful API server and a WebSocket CDC server,
using Redis for caching, supporting pluggable middleware (for things like authentication / authorization), pluggable post-mutation actions and custom endpoints.

## Tasks

- [TODO] Fix up inserts of points / polygons
- [IN PROGRESS] Support foreign key children in endpoints; thoughts:
  - Special query parameter to include an array of the reverse-relationship children
  - Gonna be a recursive mess
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
- [TODO] Better support for recursive schemas (in the case that they cause a graph cycle)
- [TODO] Support for views
- [TODO] Fix up the various templating shortcuts I've taken that cause `staticcheck` warnings (e.g. `unnecessary use of fmt.Sprintf`)
- [TODO] Think about how to do hot-reloading on schema changes (is this mostly an infra problem? Not sure)
- [TODO] Document all the features

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

### Workflow

```shell
# shell 1 - spin up the dependencies
./run-env.sh

# shell 2 - run the generated HTTP API and WebSocket CDC server
find ./pkg/model_generated -type f -name '*.go' | entr -n -r -cc -s "DJANGOLANG_SET_REPLICA_IDENTITY=full PORT=7070 REDIS_URL=redis://default:some-password@localhost:6379 POSTGRES_DB=some_db POSTGRES_PASSWORD=some-password go run ./pkg/model_generated/cmd/ serve"

# shell 3 - consume from the WebSocket CDC stream
find ./pkg/model_generated -type f -name '*.go' | entr -n -r -cc -s "while true; do unbuffer websocat ws://localhost:7070/__stream | jq; done"

# shell 4 - run the templating tests and then the integration tests for the generated code (causes some changes to be seen at the WebSocket CDC stream)
find . -type f -name '*.*' | grep -v '/model_generated/' | entr -n -r -cc -s "DJANGOLANG_DEBUG=1 REDIS_URL=redis://default:some-password@localhost:6379 POSTGRES_DB=some_db POSTGRES_PASSWORD=some-password go test  -failfast -count=1 ./pkg/types ./pkg/query ./pkg/template && DJANGOLANG_DEBUG=1 DJANGOLANG_SET_REPLICA_IDENTITY=full REDIS_URL=redis://default:some-password@localhost:6379 POSTGRES_DB=some_db POSTGRES_PASSWORD=some-password go test -v -failfast -count=1 ./pkg/model_generated_test"
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

## Notes

### Templating

I still haven't found an intuitive way to do templating with Go (I'm sure it's out there, I just haven't researched hard enough) so I've taken what I'm calling a "reference package"
approach; i.e. there's a package that contains (more or less) the desired state for the generated package (with the addition of a few meta comments used to tag blocks of code).

I've found this easier to work with than just string or string-file templating approach (i.e. something with the Go `{{ .Thing }}` markup) because the reference package contains
actual sane Go code (so you get all the benefits of linting, language server, being able to actually run the code etc).

The basic workflow to add a new feature that needs some of the meta comments is something like:

- Have the resultant `pkg/model_generated/logical_things.go` open
- Add the new feature in `pkg/model_reference/logical_things.go`
- Comment it with a new unique HTML-like tag as required
- Add a new parse task in `pkg/template/parse.go`
- Add templating for the new parse task in `pkg/template/template.go`
- Check that you get the desired result in `pkg/model_generated/logical_things.go`
- Cross-reference with `pkg/model_generated/location_history.go` or `pkg/model_generated/physical_things.go` to ensure your changes aren't merely accidentally correct
  - i.e. a field name isn't being correctly replaced (but just happens to be correct because both the reference and the generated are for the `logical_things` table)
- Add some test coverage as applicable

### Caching

The cache approach is as follows:

- Cache key of `(table name):(optional foreign key):(filter hash)`
- GET requests try to return a cached item but fall back to making a database query (and setting the result as the cached item)
- POST, PUT and PATCH requests will always return the result of the database query
- DELETE requests will always return nothing
- Between introspection and generation, record which tables reference which other tables
- Perform cache invalidation at the WebSocket CDC server using the table reference information

### OpenAPI generation and client generation

At the time of writing, OpenAPI generation appears to work and I was able to manually generate a TypeScript / SWR client; I think I'll flesh this out further as part of the example repo
rather than make it specifically a Djangolang thing.
