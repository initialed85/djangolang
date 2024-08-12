# djangolang

# status: under development, sort of works

An opinionated framework that's trying to make it easy to turn any (with various caveats) Postgres database into a HTTP RESTful API server and a WebSocket CDC server,
using Redis for caching, supporting pluggable middleware (for things like authentication / authorization), pluggable post-mutation actions and custom endpoints.

## Tasks

- [TODO] Redo the type mapping stuff again
- [TODO] Fix all the tests I commented out (I was too lazy to fix the complicated structures)
- [WIP] Support foreign key children in endpoints; thoughts:
  - [DONE] Make a recursive mess
  - [TODO] Fix the recursive mess (really make sure we can't go infinite)
  - [TODO] Special query parameter to include an array of the reverse-relationship children
  - [TODO] Fix up cache invalidation now that it flows the other way too (conditionally)
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
./run.sh env

# shell 2 - run the generated HTTP API and WebSocket CDC server
./run.sh server

# shell 3 - consume from the WebSocket CDC stream
./run.sh stream

# shell 4 - run the templating tests and then the integration tests for the generated code (causes some changes to be seen at the WebSocket CDC stream)
./run.sh test
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
