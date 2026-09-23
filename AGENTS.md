# AGENTS.md

This repository is **Djangolang**, a Go/Postgres framework that generates a REST/WebSocket API from an introspected database schema. It is not a conventional hand-written Go service: database introspection, reference templates, generated source, and integration tests are tightly coupled.

## Project map

- `pkg/query/` — generic PostgreSQL query helpers (`Select`, `Insert`, `BulkInsert`, `Update`, etc.).
- `pkg/types/` — Parse/Format functions and type metadata used by introspection and generated code.
- `pkg/introspect/` — reads PostgreSQL metadata and maps columns to Go/query/stream types.
- `pkg/model_reference/` — canonical reference model source used as the starting point for generation. This is the important source for model method structure.
- `pkg/template/` — transforms reference source and introspected tables into generated Go source.
- `pkg/model_generated/` — generated models for the normal schema.
- `pkg/model_generated_from_schema/` — generated models for the `test` schema used by schema-driven tests.
- `pkg/model_generated_test/` — database/API/CDC integration tests against generated models.
- `pkg/schema/` — schema YAML parsing/Dump and schema-generation tests.
- `database/migrations/` — PostgreSQL migrations used by Docker Compose.

## Generated-code rules

1. **Do not hand-edit files under `pkg/model_generated/` or `pkg/model_generated_from_schema/`.** Change the reference/template/introspection source instead.
2. `pkg/template.TestTemplate` writes generated files into `pkg/model_generated/`.
3. `pkg/schema.TestSchema` writes generated files into `pkg/model_generated_from_schema/`.
4. Therefore, normal test runs intentionally dirty tracked generated files. Before reviewing or committing, inspect `git status` and reset generated output unless generated snapshots are explicitly part of the change.
5. Preserve unrelated local edits, especially `Dockerfile.test`; do not use a broad `git restore .`.
6. Generated methods must compile for plural table names and non-`id` primary keys. Derive identifiers from the singular object name and `table.PrimaryKeyColumn`; never assume every model has `ID`/`TableIDColumn`.
7. Field-update generation should use the canonical formatter from `types.GetTypeMetaForTypeTemplate(...).FormatFuncTemplate`, not a partial hand-maintained type switch. Nested `fmt.Sprintf` templates are easy to break: `%sTableIDColumn` in the outer format string produces `CameraTableIDColumn`; `%s + "TableIDColumn"` emits invalid generated Go such as `Camera + "TableIDColumn"`.

## Development environment and lifecycle

Prerequisites include Docker Compose, Go, `entr`, `golang-migrate`, and the command-line tools listed in the README.

The normal project pattern is: bring up one repository's containerized environment once, then run native Go/scripts against those services. Do not assume that whatever is listening on localhost belongs to this repository.

Start this repository's environment in the background and return to the shell:

```sh
./run.sh env up
```

`./run.sh env` without `up` also follows `docker compose logs -f`; it is a long-running log-viewing command and may tear the environment down when interrupted. Use it only when that behavior is intentional.

Before running tests, confirm the current Compose project is the djangolang one and that migrations completed:

```sh
docker compose ps -a
docker ps --format '{{.Names}}\t{{.Ports}}\t{{.Status}}'
```

The `post-migrate` service should have exited successfully, while PostgreSQL, Redis, and (when needed) the test service should be running/healthy. Native tests normally connect to `localhost` with `POSTGRES_DB=some_db`, `POSTGRES_PASSWORD=some-password`, and the local Redis URL.

Stop and remove this repository's environment (including volumes, so test data is discarded):

```sh
./run.sh env down
```

**Cross-project port collision:** djangolang and Camry currently use overlapping localhost ports (notably PostgreSQL 5432, Redis 6379, and Swagger 7071). They cannot be run simultaneously on the same machine. Before bringing one up, stop the other project's environment; before taking anything down, verify the Compose project/container names so you do not destroy the other project. Never run djangolang tests against Camry's PostgreSQL just because port 5432 is open.

`Dockerfile.test` controls the test image. Treat local edits to it as intentional unless confirmed otherwise.

## Tests

### Preferred headless/full workflow

After `./run.sh env` has completed migrations, use:

```sh
./run.sh test-ci
```

This is the non-interactive CI-style workflow. It runs race-enabled, verbose, fail-fast tests in the Docker `test` container: template tests first, then all test-package directories. It waits for the `post-migrate` service to have exited successfully.

Capture output when the logs are large:

```sh
./run.sh test-ci > /tmp/djangolang-test-ci.log 2>&1
status=$?
tail -200 /tmp/djangolang-test-ci.log
exit "$status"
```

`./run.sh test` is an `entr` watcher intended for an interactive development shell. It runs a pass and then waits for file changes; it is expected not to exit on its own. Prefer `test-ci` in an agent/headless session.

`./run.sh test-ci` executes tests inside the Docker `test` container. The Compose file volume-mounts the repository at `/srv`, so it sees current source, but this is still different from native Go execution. When the requested workflow is native, run the same package list directly on the host after `./run.sh env up`; do not silently substitute the container workflow.

### Useful targeted commands

Compile all packages without running tests:

```sh
go test ./... -run '^$'
```

Run the query suite locally when PostgreSQL/Redis are exposed on localhost:

```sh
POSTGRES_DB=some_db \
POSTGRES_PASSWORD=some-password \
REDIS_URL=redis://default:some-password@localhost:6379 \
go test -race -v -failfast -count=1 ./pkg/query
```

The BulkInsert regression is in `pkg/query/query_test.go` and exercises more than 8,191 eight-parameter rows. It must remain a single caller-owned transaction, split below PostgreSQL's 65,535 bind-parameter limit, and concatenate `RETURNING` rows in input order.

Template/schema tests require the database environment and rewrite generated files. Their tests now also compile the generated packages so generator regressions fail at the source-generation stage rather than later in integration tests.

## Generated-file ownership recovery

The Docker test container may write generated files as `root`, especially after a failed test run or after the environment is stopped. If `git restore` reports permission errors, fix ownership only for generated directories, then reset them:

```sh
sudo chown -R "$(id -u):$(id -g)" \
  pkg/model_generated pkg/model_generated_from_schema

git restore pkg/model_generated pkg/model_generated_from_schema
```

Do not restore the whole worktree; that can erase a real local change.

## Follow established patterns first

Prefer extending an existing coherent path over inventing a parallel abstraction. In particular, `SelectLogicalThings` and the generated list handlers already define how typed filters, `OrderBy`, limits, offsets, and request arguments flow from an HTTP query into SQL. Claim endpoints should follow that same pattern: reuse `server.GetSelectManyArguments`, pass its validated `Where`, `OrderBy`, and `Values` through, and only apply a claim-specific default when `OrderBy` is nil. Do not add a separate `WithOrderBy` API or bypass the typed argument path unless the existing pattern genuinely cannot express the behavior.

A pattern can contain a bug, but first understand and preserve its surrounding conventions; coherence is a project feature and makes generated code/client behavior predictable.

## Query-layer cautions

- `BulkInsert` receives one flattened values slice; `len(columns)` is the number of values per row.
- PostgreSQL's extended protocol allows at most 65,535 bind parameters. `ON CONFLICT ... DO UPDATE` adds statement-level parameters and must be included in the budget.
- Chunking must use the transaction supplied by the caller. `BulkInsert` must not begin, commit, or rollback its own transaction.
- Preserve `RETURNING` behavior and append chunk results in input/chunk order.
- Run malformed bulk-shape validation before issuing the first chunk.

## Release checklist

1. Run `./run.sh test-ci` and record the package summaries/exit status.
2. Reset only test-generated output; retain intentional source/local changes.
3. Verify `git status`, `git diff --check`, and `git log`.
4. Commit source, tests, and intentional build/config changes.
5. Push `main` before publishing a release tag.
6. Tags follow the existing lightweight `v0.1.N` convention. Create and push the next tag explicitly:

```sh
git push origin main
git tag v0.1.N HEAD
git push origin refs/tags/v0.1.N
```

Never force-move a release tag after consumers may have fetched it without explicit approval.

## Common failure interpretation

- Compile errors mentioning `Camera + "TableIDColumn"`, plural names such as `LogicalThings`, or missing `ID` on a model with another primary key indicate generated field-update/template regressions, not query-layer failures.
- A dirty generated tree immediately after template/schema tests is expected; inspect the diff before deciding it is a source change.
- A watcher command that appears to hang after printing `(done)` is usually waiting for `entr` file changes, not stuck in Go tests.
