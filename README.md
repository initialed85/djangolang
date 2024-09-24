# djangolang

# status: under development, sort of works

An opinionated framework that's trying to make it easy to turn any (with various caveats) Postgres database into a HTTP RESTful API server and a WebSocket CDC server,
using Redis for caching and supporting pluggable middleware for things like authentication / authorization.

## Tasks

- Roadmap
  - [DONE] Introspection functionality
  - [DONE] Generated SQL helpers
  - [DONE] Generated endpoint
  - [DONE] Generic HTTP server
  - [DONE] Generic CDC server
  - [TODO] Cache presented via endpoints and invalidated via CDC
- Bits and pieces
  - [TODO] Support more Postgres data types as they come up
  - [TODO] Support recursive schemas
  - [TODO] Support views

## Usage

At the moment these are just instructions for myself while developing; this is my current workflow:

```shell
# shell 1 - spin up the dependencies
./run-env.sh

# shell 2 - run the generated HTTP API and WebSocket CDC server
find ./pkg/model_generated -type f -name '*.go' | entr -n -r -cc -s "PORT=7070 REDIS_URL=redis://default:some-password@localhost:6379 POSTGRES_DB=some_db POSTGRES_PASSWORD=some-password go run ./pkg/model_generated/cmd/"

# shell 3 - consume from the WebSocket CDC stream
find ./pkg/model_generated -type f -name '*.go' | entr -n -r -cc -s "while true; do unbuffer websocat ws://localhost:7070/__stream | jq; done"

# shell 4 - run the templating tests and then the integration tests for the generated code (causes some changes to be seen at the WebSocket CDC stream)
find . -type f -name '*.*' | grep -v '/model_generated/' | entr -n -r -cc -s "DJANGOLANG_DEBUG=1 REDIS_URL=redis://default:some-password@localhost:6379 POSTGRES_DB=some_db POSTGRES_PASSWORD=some-password go test -v -failfast -count=1 ./pkg/template && DJANGOLANG_DEBUG=1 REDIS_URL=redis://default:some-password@localhost:6379 POSTGRES_DB=some_db POSTGRES_PASSWORD=some-password go test -v -failfast -count=1 ./pkg/model_generated_test"
```

Everything should restart automatically when the code changes, testing first any templating aspects and then (if that works) testing the actual behaviours; shells 2 and 3 are
mostly just a smoke test.

## Notes

### Templating

I still haven't found an intuitive way to do templating with Go (I'm sure it's out there, I just haven't researched hard enough) so I've taken what I'm calling a "reference package"
approach; i.e. there's a package that contains (more or less) the desired state for the generated package (with the addition of a few meta comments used to tag blocks of code).

I've found this easier to work with than just string or string-file templating approach (i.e. something with the Go `{{ .Thing }}` markup) because the reference package contains
actual sane Go code (so you get all the benefits of linting, language server, being able to actually run the code etc).

The basic workflow to add a new feature that needs some of the meta comments is something like:

- Have the resultant `pkg/model_generated/logical_things.go` open
- Add the new feature in `pkg/model_reference/logical_thing.go`
- Comment it with a new unique HTML-like tag as required
- Add a new parse task in `pkg/template/parse.go`
- Add templating for the new parse task in `pkg/template/template.go`
- Check that you get the desired result in `pkg/model_generated/logical_things.go`
- Cross-reference with `pkg/model_generated/location_history.go` or `pkg/model_generated/physical_things.go` to ensure your changes aren't merely accidentally correct
  - i.e. a field name isn't being correctly replaced (but just happens to be correct because both the reference and the generated are for the `logical_things` table)
- Add some test coverage as applicable

### Caching

So the broad approach will be something like:

- Use a cache key of `(table name):(filter hash)` for many-object endpoints
- Use a cache key of `(table name):(foreign key hash)` for single-object endpoints
- GET requests will prefer to return a cached item but fall back to returning an uncached item (and caching it in the process)
- POST, PUT and PATCH requests will always return an uncached item
- DELETE requests will always return nothing
- Between introspection and generation, record which tables reference which other tables
- Perform cache invalidation at the WebSocket CDC server using the table reference information

Blah
