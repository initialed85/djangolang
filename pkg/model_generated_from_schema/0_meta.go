package model_generated_from_schema

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gomodule/redigo/redis"
	"github.com/initialed85/djangolang/pkg/config"
	"github.com/initialed85/djangolang/pkg/introspect"
	"github.com/initialed85/djangolang/pkg/openapi"
	"github.com/initialed85/djangolang/pkg/server"
	"github.com/initialed85/djangolang/pkg/types"
	"github.com/jackc/pgx/v5/pgxpool"
	"gopkg.in/yaml.v2"

	"net/http/pprof"
)

var mu = new(sync.Mutex)
var newFromItemFnByTableName = make(map[string]func(map[string]any) (any, error))
var getRouterFnByPattern = make(map[string]server.GetRouterFn)
var allObjects = make([]any, 0)
var openApi *types.OpenAPI
var profile = config.Profile()

var httpHandlerSummaries []server.HTTPHandlerSummary = make([]server.HTTPHandlerSummary, 0)

func isRequired(columns map[string]*introspect.Column, columnName string) bool {
	column := columns[columnName]
	if column == nil {
		return false
	}

	return column.NotNull && !column.HasDefault
}

func register(
	tableName string,
	object any,
	newFromItem func(map[string]any) (any, error),
	pattern string,
	getRouterFn server.GetRouterFn,
) {
	allObjects = append(allObjects, object)
	newFromItemFnByTableName[tableName] = newFromItem
	getRouterFnByPattern[pattern] = getRouterFn
}

func GetOpenAPI() (*types.OpenAPI, error) {
	mu.Lock()
	defer mu.Unlock()

	if openApi != nil {
		return openApi, nil
	}

	var err error
	openApi, err = openapi.NewFromIntrospectedSchema(httpHandlerSummaries)
	if err != nil {
		return nil, err
	}

	return openApi, nil
}

func NewFromItem(tableName string, item map[string]any) (any, error) {
	if item == nil {
		return nil, nil
	}

	mu.Lock()
	newFromItemFn, ok := newFromItemFnByTableName[tableName]
	mu.Unlock()

	if !ok {
		return nil, fmt.Errorf("table name %v not known", tableName)
	}

	return newFromItemFn(item)
}

func GetRouter(db *pgxpool.Pool, redisPool *redis.Pool, httpMiddlewares []server.HTTPMiddleware, objectMiddlewares []server.ObjectMiddleware, waitForChange server.WaitForChange) chi.Router {
	r := chi.NewRouter()

	mu.Lock()
	getRouterFnByPattern := getRouterFnByPattern
	mu.Unlock()

	for pattern, getRouterFn := range getRouterFnByPattern {
		r.Mount(pattern, getRouterFn(db, redisPool, httpMiddlewares, objectMiddlewares, waitForChange))
	}

	healthzMu := new(sync.Mutex)
	healthzExpiresAt := time.Now().Add(-time.Second * 5)
	var lastHealthz error

	healthz := func(ctx context.Context) error {
		healthzMu.Lock()
		defer healthzMu.Unlock()

		if time.Now().Before(healthzExpiresAt) {
			return lastHealthz
		}

		lastHealthz = func() error {
			err := db.Ping(ctx)
			if err != nil {
				return fmt.Errorf("db ping failed; %v", err)
			}

			redisConn, err := redisPool.GetContext(ctx)
			if err != nil {
				return fmt.Errorf("redis pool get failed; %v", err)
			}

			defer func() {
				_ = redisConn.Close()
			}()

			_, err = redisConn.Do("PING")
			if err != nil {
				return fmt.Errorf("redis ping failed; %v", err)
			}

			return nil
		}()

		healthzExpiresAt = time.Now().Add(time.Second * 5)

		return lastHealthz
	}

	if profile {
		r.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		r.HandleFunc("/debug/pprof/profile", pprof.Profile)
		r.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		r.HandleFunc("/debug/pprof/trace", pprof.Trace)

		r.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
		r.Handle("/debug/pprof/heap", pprof.Handler("heap"))
		r.Handle("/debug/pprof/allocs", pprof.Handler("allocs"))
		r.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
		r.Handle("/debug/pprof/block", pprof.Handler("block"))
		r.Handle("/debug/pprof/mutex", pprof.Handler("mutex"))

		r.HandleFunc("/debug/pprof/", pprof.Index)
	}

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
		defer cancel()

		err := healthz(ctx)
		if err != nil {
			server.HandleErrorResponse(w, http.StatusInternalServerError, err)
			return
		}

		server.HandleObjectsResponse(w, http.StatusOK, nil)
	})

	r.Get("/openapi.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-type", "application/json")

		openApi, err := GetOpenAPI()
		if err != nil {
			server.HandleErrorResponse(w, http.StatusInternalServerError, fmt.Errorf("failed to get OpenAPI schema; %v", err))
			return
		}

		b, err := json.MarshalIndent(openApi, "", "  ")
		if err != nil {
			server.HandleErrorResponse(w, http.StatusInternalServerError, fmt.Errorf("failed to get OpenAPI schema; %v", err))
			return
		}

		server.WriteResponse(w, http.StatusOK, b)
	})

	r.Get("/openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-type", "application/yaml")

		openApi, err := GetOpenAPI()
		if err != nil {
			server.HandleErrorResponse(w, http.StatusInternalServerError, fmt.Errorf("failed to get OpenAPI schema; %v", err))
			return
		}

		b, err := yaml.Marshal(openApi)
		if err != nil {
			server.HandleErrorResponse(w, http.StatusInternalServerError, fmt.Errorf("failed to get OpenAPI schema; %v", err))
			return
		}

		server.WriteResponse(w, http.StatusOK, b)
	})

	return r
}

func getHTTPHandler[T any, S any, Q any, R any](method string, path string, status int, handle func(context.Context, T, S, Q, any) (R, error), modelObject any) (*server.HTTPHandler[T, S, Q, R], error) {
	customHTTPHandler, err := server.GetHTTPHandler(method, path, status, handle)
	if err != nil {
		return nil, err
	}

	customHTTPHandler.Builtin = true
	customHTTPHandler.BuiltinModelObject = modelObject

	mu.Lock()
	httpHandlerSummaries = append(httpHandlerSummaries, customHTTPHandler.Summarize())
	mu.Unlock()

	return customHTTPHandler, nil
}

func GetHTTPHandler[T any, S any, Q any, R any](method string, path string, status int, handle func(context.Context, T, S, Q, any) (R, error)) (*server.HTTPHandler[T, S, Q, R], error) {
	customHTTPHandler, err := server.GetHTTPHandler(method, path, status, handle)
	if err != nil {
		return nil, err
	}

	mu.Lock()
	httpHandlerSummaries = append(httpHandlerSummaries, customHTTPHandler.Summarize())
	mu.Unlock()

	return customHTTPHandler, nil
}

func GetHTTPHandlerSummaries() ([]server.HTTPHandlerSummary, error) {
	mu.Lock()
	defer mu.Unlock()

	if len(httpHandlerSummaries) == 0 {
		return nil, fmt.Errorf("httpHandlerSummaries unexpectedly empty; you'll need to call GetRouter() so that this is populated")
	}

	return httpHandlerSummaries, nil
}

var tableByNameAsJSON = []byte(`{
  "change": {
    "tablename": "change",
    "oid": "31499",
    "schema": "test",
    "reltuples": -1,
    "relkind": "r",
    "relam": "2",
    "relacl": null,
    "reltype": "31501",
    "relowner": "10",
    "relhasindex": true,
    "columns": [
      {
        "column": "id",
        "datatype": "uuid",
        "table": "change",
        "pos": 1,
        "typeid": "2950",
        "typelen": 16,
        "typemod": -1,
        "notnull": true,
        "hasdefault": true,
        "hasmissing": false,
        "ispkey": true,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "31499",
        "zero_type": "00000000-0000-0000-0000-000000000000",
        "query_type_template": "uuid.UUID",
        "stream_type_template": "[16]uint8",
        "type_template": "uuid.UUID"
      },
      {
        "column": "created_at",
        "datatype": "timestamp with time zone",
        "table": "change",
        "pos": 2,
        "typeid": "1184",
        "typelen": 8,
        "typemod": -1,
        "notnull": true,
        "hasdefault": true,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "31499",
        "zero_type": "0001-01-01T00:00:00Z",
        "query_type_template": "time.Time",
        "stream_type_template": "time.Time",
        "type_template": "time.Time"
      },
      {
        "column": "updated_at",
        "datatype": "timestamp with time zone",
        "table": "change",
        "pos": 3,
        "typeid": "1184",
        "typelen": 8,
        "typemod": -1,
        "notnull": true,
        "hasdefault": true,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "31499",
        "zero_type": "0001-01-01T00:00:00Z",
        "query_type_template": "time.Time",
        "stream_type_template": "time.Time",
        "type_template": "time.Time"
      },
      {
        "column": "deleted_at",
        "datatype": "timestamp with time zone",
        "table": "change",
        "pos": 4,
        "typeid": "1184",
        "typelen": 8,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "31499",
        "zero_type": "0001-01-01T00:00:00Z",
        "query_type_template": "time.Time",
        "stream_type_template": "time.Time",
        "type_template": "time.Time"
      },
      {
        "column": "commit_hash",
        "datatype": "text",
        "table": "change",
        "pos": 5,
        "typeid": "25",
        "typelen": -1,
        "typemod": -1,
        "notnull": true,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "31499",
        "zero_type": "",
        "query_type_template": "string",
        "stream_type_template": "string",
        "type_template": "string"
      },
      {
        "column": "authored_by",
        "datatype": "text",
        "table": "change",
        "pos": 6,
        "typeid": "25",
        "typelen": -1,
        "typemod": -1,
        "notnull": true,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "31499",
        "zero_type": "",
        "query_type_template": "string",
        "stream_type_template": "string",
        "type_template": "string"
      },
      {
        "column": "authored_at",
        "datatype": "timestamp with time zone",
        "table": "change",
        "pos": 7,
        "typeid": "1184",
        "typelen": 8,
        "typemod": -1,
        "notnull": true,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "31499",
        "zero_type": "0001-01-01T00:00:00Z",
        "query_type_template": "time.Time",
        "stream_type_template": "time.Time",
        "type_template": "time.Time"
      },
      {
        "column": "repository_id",
        "datatype": "uuid",
        "table": "change",
        "pos": 8,
        "typeid": "2950",
        "typelen": 16,
        "typemod": -1,
        "notnull": true,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": "repository",
        "fcolumn": "id",
        "parent_id": "31499",
        "zero_type": "00000000-0000-0000-0000-000000000000",
        "query_type_template": "uuid.UUID",
        "stream_type_template": "[16]uint8",
        "type_template": "uuid.UUID"
      }
    ]
  },
  "execution": {
    "tablename": "execution",
    "oid": "31512",
    "schema": "test",
    "reltuples": -1,
    "relkind": "r",
    "relam": "2",
    "relacl": null,
    "reltype": "31514",
    "relowner": "10",
    "relhasindex": true,
    "columns": [
      {
        "column": "id",
        "datatype": "uuid",
        "table": "execution",
        "pos": 1,
        "typeid": "2950",
        "typelen": 16,
        "typemod": -1,
        "notnull": true,
        "hasdefault": true,
        "hasmissing": false,
        "ispkey": true,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "31512",
        "zero_type": "00000000-0000-0000-0000-000000000000",
        "query_type_template": "uuid.UUID",
        "stream_type_template": "[16]uint8",
        "type_template": "uuid.UUID"
      },
      {
        "column": "created_at",
        "datatype": "timestamp with time zone",
        "table": "execution",
        "pos": 2,
        "typeid": "1184",
        "typelen": 8,
        "typemod": -1,
        "notnull": true,
        "hasdefault": true,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "31512",
        "zero_type": "0001-01-01T00:00:00Z",
        "query_type_template": "time.Time",
        "stream_type_template": "time.Time",
        "type_template": "time.Time"
      },
      {
        "column": "updated_at",
        "datatype": "timestamp with time zone",
        "table": "execution",
        "pos": 3,
        "typeid": "1184",
        "typelen": 8,
        "typemod": -1,
        "notnull": true,
        "hasdefault": true,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "31512",
        "zero_type": "0001-01-01T00:00:00Z",
        "query_type_template": "time.Time",
        "stream_type_template": "time.Time",
        "type_template": "time.Time"
      },
      {
        "column": "deleted_at",
        "datatype": "timestamp with time zone",
        "table": "execution",
        "pos": 4,
        "typeid": "1184",
        "typelen": 8,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "31512",
        "zero_type": "0001-01-01T00:00:00Z",
        "query_type_template": "time.Time",
        "stream_type_template": "time.Time",
        "type_template": "time.Time"
      },
      {
        "column": "status",
        "datatype": "text",
        "table": "execution",
        "pos": 5,
        "typeid": "25",
        "typelen": -1,
        "typemod": -1,
        "notnull": true,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "31512",
        "zero_type": "",
        "query_type_template": "string",
        "stream_type_template": "string",
        "type_template": "string"
      },
      {
        "column": "started_at",
        "datatype": "timestamp with time zone",
        "table": "execution",
        "pos": 6,
        "typeid": "1184",
        "typelen": 8,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "31512",
        "zero_type": "0001-01-01T00:00:00Z",
        "query_type_template": "time.Time",
        "stream_type_template": "time.Time",
        "type_template": "time.Time"
      },
      {
        "column": "ended_at",
        "datatype": "timestamp with time zone",
        "table": "execution",
        "pos": 7,
        "typeid": "1184",
        "typelen": 8,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "31512",
        "zero_type": "0001-01-01T00:00:00Z",
        "query_type_template": "time.Time",
        "stream_type_template": "time.Time",
        "type_template": "time.Time"
      },
      {
        "column": "task_id",
        "datatype": "uuid",
        "table": "execution",
        "pos": 8,
        "typeid": "2950",
        "typelen": 16,
        "typemod": -1,
        "notnull": true,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": "task",
        "fcolumn": "id",
        "parent_id": "31512",
        "zero_type": "00000000-0000-0000-0000-000000000000",
        "query_type_template": "uuid.UUID",
        "stream_type_template": "[16]uint8",
        "type_template": "uuid.UUID"
      },
      {
        "column": "m2m_rule_trigger_job_id",
        "datatype": "uuid",
        "table": "execution",
        "pos": 9,
        "typeid": "2950",
        "typelen": 16,
        "typemod": -1,
        "notnull": true,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": "m2m_rule_trigger_job",
        "fcolumn": "id",
        "parent_id": "31512",
        "zero_type": "00000000-0000-0000-0000-000000000000",
        "query_type_template": "uuid.UUID",
        "stream_type_template": "[16]uint8",
        "type_template": "uuid.UUID"
      }
    ]
  },
  "job": {
    "tablename": "job",
    "oid": "31470",
    "schema": "test",
    "reltuples": -1,
    "relkind": "r",
    "relam": "2",
    "relacl": null,
    "reltype": "31472",
    "relowner": "10",
    "relhasindex": true,
    "columns": [
      {
        "column": "id",
        "datatype": "uuid",
        "table": "job",
        "pos": 1,
        "typeid": "2950",
        "typelen": 16,
        "typemod": -1,
        "notnull": true,
        "hasdefault": true,
        "hasmissing": false,
        "ispkey": true,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "31470",
        "zero_type": "00000000-0000-0000-0000-000000000000",
        "query_type_template": "uuid.UUID",
        "stream_type_template": "[16]uint8",
        "type_template": "uuid.UUID"
      },
      {
        "column": "created_at",
        "datatype": "timestamp with time zone",
        "table": "job",
        "pos": 2,
        "typeid": "1184",
        "typelen": 8,
        "typemod": -1,
        "notnull": true,
        "hasdefault": true,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "31470",
        "zero_type": "0001-01-01T00:00:00Z",
        "query_type_template": "time.Time",
        "stream_type_template": "time.Time",
        "type_template": "time.Time"
      },
      {
        "column": "updated_at",
        "datatype": "timestamp with time zone",
        "table": "job",
        "pos": 3,
        "typeid": "1184",
        "typelen": 8,
        "typemod": -1,
        "notnull": true,
        "hasdefault": true,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "31470",
        "zero_type": "0001-01-01T00:00:00Z",
        "query_type_template": "time.Time",
        "stream_type_template": "time.Time",
        "type_template": "time.Time"
      },
      {
        "column": "deleted_at",
        "datatype": "timestamp with time zone",
        "table": "job",
        "pos": 4,
        "typeid": "1184",
        "typelen": 8,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "31470",
        "zero_type": "0001-01-01T00:00:00Z",
        "query_type_template": "time.Time",
        "stream_type_template": "time.Time",
        "type_template": "time.Time"
      },
      {
        "column": "name",
        "datatype": "text",
        "table": "job",
        "pos": 5,
        "typeid": "25",
        "typelen": -1,
        "typemod": -1,
        "notnull": true,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "31470",
        "zero_type": "",
        "query_type_template": "string",
        "stream_type_template": "string",
        "type_template": "string"
      },
      {
        "column": "rule_trigger_rule_id",
        "datatype": "uuid",
        "table": "job",
        "pos": 6,
        "typeid": "2950",
        "typelen": 16,
        "typemod": -1,
        "notnull": true,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": "rule",
        "fcolumn": "id",
        "parent_id": "31470",
        "zero_type": "00000000-0000-0000-0000-000000000000",
        "query_type_template": "uuid.UUID",
        "stream_type_template": "[16]uint8",
        "type_template": "uuid.UUID"
      }
    ]
  },
  "log": {
    "tablename": "log",
    "oid": "31538",
    "schema": "test",
    "reltuples": -1,
    "relkind": "r",
    "relam": "2",
    "relacl": null,
    "reltype": "31540",
    "relowner": "10",
    "relhasindex": true,
    "columns": [
      {
        "column": "id",
        "datatype": "uuid",
        "table": "log",
        "pos": 1,
        "typeid": "2950",
        "typelen": 16,
        "typemod": -1,
        "notnull": true,
        "hasdefault": true,
        "hasmissing": false,
        "ispkey": true,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "31538",
        "zero_type": "00000000-0000-0000-0000-000000000000",
        "query_type_template": "uuid.UUID",
        "stream_type_template": "[16]uint8",
        "type_template": "uuid.UUID"
      },
      {
        "column": "created_at",
        "datatype": "timestamp with time zone",
        "table": "log",
        "pos": 2,
        "typeid": "1184",
        "typelen": 8,
        "typemod": -1,
        "notnull": true,
        "hasdefault": true,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "31538",
        "zero_type": "0001-01-01T00:00:00Z",
        "query_type_template": "time.Time",
        "stream_type_template": "time.Time",
        "type_template": "time.Time"
      },
      {
        "column": "updated_at",
        "datatype": "timestamp with time zone",
        "table": "log",
        "pos": 3,
        "typeid": "1184",
        "typelen": 8,
        "typemod": -1,
        "notnull": true,
        "hasdefault": true,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "31538",
        "zero_type": "0001-01-01T00:00:00Z",
        "query_type_template": "time.Time",
        "stream_type_template": "time.Time",
        "type_template": "time.Time"
      },
      {
        "column": "deleted_at",
        "datatype": "timestamp with time zone",
        "table": "log",
        "pos": 4,
        "typeid": "1184",
        "typelen": 8,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "31538",
        "zero_type": "0001-01-01T00:00:00Z",
        "query_type_template": "time.Time",
        "stream_type_template": "time.Time",
        "type_template": "time.Time"
      },
      {
        "column": "buffer",
        "datatype": "bytea",
        "table": "log",
        "pos": 5,
        "typeid": "17",
        "typelen": -1,
        "typemod": -1,
        "notnull": true,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "31538",
        "zero_type": "",
        "query_type_template": "[]byte",
        "stream_type_template": "[]byte",
        "type_template": "[]byte"
      },
      {
        "column": "output_id",
        "datatype": "uuid",
        "table": "log",
        "pos": 6,
        "typeid": "2950",
        "typelen": 16,
        "typemod": -1,
        "notnull": true,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": "output",
        "fcolumn": "id",
        "parent_id": "31538",
        "zero_type": "00000000-0000-0000-0000-000000000000",
        "query_type_template": "uuid.UUID",
        "stream_type_template": "[16]uint8",
        "type_template": "uuid.UUID"
      }
    ]
  },
  "m2m_rule_trigger_job": {
    "tablename": "m2m_rule_trigger_job",
    "oid": "31561",
    "schema": "test",
    "reltuples": -1,
    "relkind": "r",
    "relam": "2",
    "relacl": null,
    "reltype": "31563",
    "relowner": "10",
    "relhasindex": true,
    "columns": [
      {
        "column": "id",
        "datatype": "uuid",
        "table": "m2m_rule_trigger_job",
        "pos": 1,
        "typeid": "2950",
        "typelen": 16,
        "typemod": -1,
        "notnull": true,
        "hasdefault": true,
        "hasmissing": false,
        "ispkey": true,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "31561",
        "zero_type": "00000000-0000-0000-0000-000000000000",
        "query_type_template": "uuid.UUID",
        "stream_type_template": "[16]uint8",
        "type_template": "uuid.UUID"
      },
      {
        "column": "created_at",
        "datatype": "timestamp with time zone",
        "table": "m2m_rule_trigger_job",
        "pos": 2,
        "typeid": "1184",
        "typelen": 8,
        "typemod": -1,
        "notnull": true,
        "hasdefault": true,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "31561",
        "zero_type": "0001-01-01T00:00:00Z",
        "query_type_template": "time.Time",
        "stream_type_template": "time.Time",
        "type_template": "time.Time"
      },
      {
        "column": "updated_at",
        "datatype": "timestamp with time zone",
        "table": "m2m_rule_trigger_job",
        "pos": 3,
        "typeid": "1184",
        "typelen": 8,
        "typemod": -1,
        "notnull": true,
        "hasdefault": true,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "31561",
        "zero_type": "0001-01-01T00:00:00Z",
        "query_type_template": "time.Time",
        "stream_type_template": "time.Time",
        "type_template": "time.Time"
      },
      {
        "column": "deleted_at",
        "datatype": "timestamp with time zone",
        "table": "m2m_rule_trigger_job",
        "pos": 4,
        "typeid": "1184",
        "typelen": 8,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "31561",
        "zero_type": "0001-01-01T00:00:00Z",
        "query_type_template": "time.Time",
        "stream_type_template": "time.Time",
        "type_template": "time.Time"
      }
    ]
  },
  "output": {
    "tablename": "output",
    "oid": "31525",
    "schema": "test",
    "reltuples": -1,
    "relkind": "r",
    "relam": "2",
    "relacl": null,
    "reltype": "31527",
    "relowner": "10",
    "relhasindex": true,
    "columns": [
      {
        "column": "id",
        "datatype": "uuid",
        "table": "output",
        "pos": 1,
        "typeid": "2950",
        "typelen": 16,
        "typemod": -1,
        "notnull": true,
        "hasdefault": true,
        "hasmissing": false,
        "ispkey": true,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "31525",
        "zero_type": "00000000-0000-0000-0000-000000000000",
        "query_type_template": "uuid.UUID",
        "stream_type_template": "[16]uint8",
        "type_template": "uuid.UUID"
      },
      {
        "column": "created_at",
        "datatype": "timestamp with time zone",
        "table": "output",
        "pos": 2,
        "typeid": "1184",
        "typelen": 8,
        "typemod": -1,
        "notnull": true,
        "hasdefault": true,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "31525",
        "zero_type": "0001-01-01T00:00:00Z",
        "query_type_template": "time.Time",
        "stream_type_template": "time.Time",
        "type_template": "time.Time"
      },
      {
        "column": "updated_at",
        "datatype": "timestamp with time zone",
        "table": "output",
        "pos": 3,
        "typeid": "1184",
        "typelen": 8,
        "typemod": -1,
        "notnull": true,
        "hasdefault": true,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "31525",
        "zero_type": "0001-01-01T00:00:00Z",
        "query_type_template": "time.Time",
        "stream_type_template": "time.Time",
        "type_template": "time.Time"
      },
      {
        "column": "deleted_at",
        "datatype": "timestamp with time zone",
        "table": "output",
        "pos": 4,
        "typeid": "1184",
        "typelen": 8,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "31525",
        "zero_type": "0001-01-01T00:00:00Z",
        "query_type_template": "time.Time",
        "stream_type_template": "time.Time",
        "type_template": "time.Time"
      },
      {
        "column": "status",
        "datatype": "text",
        "table": "output",
        "pos": 5,
        "typeid": "25",
        "typelen": -1,
        "typemod": -1,
        "notnull": true,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "31525",
        "zero_type": "",
        "query_type_template": "string",
        "stream_type_template": "string",
        "type_template": "string"
      },
      {
        "column": "started_at",
        "datatype": "timestamp with time zone",
        "table": "output",
        "pos": 6,
        "typeid": "1184",
        "typelen": 8,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "31525",
        "zero_type": "0001-01-01T00:00:00Z",
        "query_type_template": "time.Time",
        "stream_type_template": "time.Time",
        "type_template": "time.Time"
      },
      {
        "column": "ended_at",
        "datatype": "timestamp with time zone",
        "table": "output",
        "pos": 7,
        "typeid": "1184",
        "typelen": 8,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "31525",
        "zero_type": "0001-01-01T00:00:00Z",
        "query_type_template": "time.Time",
        "stream_type_template": "time.Time",
        "type_template": "time.Time"
      },
      {
        "column": "exit_status",
        "datatype": "integer",
        "table": "output",
        "pos": 8,
        "typeid": "23",
        "typelen": 4,
        "typemod": -1,
        "notnull": true,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "31525",
        "zero_type": 0,
        "query_type_template": "int64",
        "stream_type_template": "int64",
        "type_template": "int64"
      },
      {
        "column": "error",
        "datatype": "text",
        "table": "output",
        "pos": 9,
        "typeid": "25",
        "typelen": -1,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "31525",
        "zero_type": "",
        "query_type_template": "string",
        "stream_type_template": "string",
        "type_template": "string"
      },
      {
        "column": "task_id",
        "datatype": "uuid",
        "table": "output",
        "pos": 10,
        "typeid": "2950",
        "typelen": 16,
        "typemod": -1,
        "notnull": true,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": "task",
        "fcolumn": "id",
        "parent_id": "31525",
        "zero_type": "00000000-0000-0000-0000-000000000000",
        "query_type_template": "uuid.UUID",
        "stream_type_template": "[16]uint8",
        "type_template": "uuid.UUID"
      },
      {
        "column": "logid",
        "datatype": "uuid",
        "table": "output",
        "pos": 11,
        "typeid": "2950",
        "typelen": 16,
        "typemod": -1,
        "notnull": true,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": "log",
        "fcolumn": "id",
        "parent_id": "31525",
        "zero_type": "00000000-0000-0000-0000-000000000000",
        "query_type_template": "uuid.UUID",
        "stream_type_template": "[16]uint8",
        "type_template": "uuid.UUID"
      }
    ]
  },
  "repository": {
    "tablename": "repository",
    "oid": "31444",
    "schema": "test",
    "reltuples": -1,
    "relkind": "r",
    "relam": "2",
    "relacl": null,
    "reltype": "31446",
    "relowner": "10",
    "relhasindex": true,
    "columns": [
      {
        "column": "id",
        "datatype": "uuid",
        "table": "repository",
        "pos": 1,
        "typeid": "2950",
        "typelen": 16,
        "typemod": -1,
        "notnull": true,
        "hasdefault": true,
        "hasmissing": false,
        "ispkey": true,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "31444",
        "zero_type": "00000000-0000-0000-0000-000000000000",
        "query_type_template": "uuid.UUID",
        "stream_type_template": "[16]uint8",
        "type_template": "uuid.UUID"
      },
      {
        "column": "created_at",
        "datatype": "timestamp with time zone",
        "table": "repository",
        "pos": 2,
        "typeid": "1184",
        "typelen": 8,
        "typemod": -1,
        "notnull": true,
        "hasdefault": true,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "31444",
        "zero_type": "0001-01-01T00:00:00Z",
        "query_type_template": "time.Time",
        "stream_type_template": "time.Time",
        "type_template": "time.Time"
      },
      {
        "column": "updated_at",
        "datatype": "timestamp with time zone",
        "table": "repository",
        "pos": 3,
        "typeid": "1184",
        "typelen": 8,
        "typemod": -1,
        "notnull": true,
        "hasdefault": true,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "31444",
        "zero_type": "0001-01-01T00:00:00Z",
        "query_type_template": "time.Time",
        "stream_type_template": "time.Time",
        "type_template": "time.Time"
      },
      {
        "column": "deleted_at",
        "datatype": "timestamp with time zone",
        "table": "repository",
        "pos": 4,
        "typeid": "1184",
        "typelen": 8,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "31444",
        "zero_type": "0001-01-01T00:00:00Z",
        "query_type_template": "time.Time",
        "stream_type_template": "time.Time",
        "type_template": "time.Time"
      },
      {
        "column": "url",
        "datatype": "text",
        "table": "repository",
        "pos": 5,
        "typeid": "25",
        "typelen": -1,
        "typemod": -1,
        "notnull": true,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "31444",
        "zero_type": "",
        "query_type_template": "string",
        "stream_type_template": "string",
        "type_template": "string"
      }
    ]
  },
  "rule": {
    "tablename": "rule",
    "oid": "31457",
    "schema": "test",
    "reltuples": -1,
    "relkind": "r",
    "relam": "2",
    "relacl": null,
    "reltype": "31459",
    "relowner": "10",
    "relhasindex": true,
    "columns": [
      {
        "column": "id",
        "datatype": "uuid",
        "table": "rule",
        "pos": 1,
        "typeid": "2950",
        "typelen": 16,
        "typemod": -1,
        "notnull": true,
        "hasdefault": true,
        "hasmissing": false,
        "ispkey": true,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "31457",
        "zero_type": "00000000-0000-0000-0000-000000000000",
        "query_type_template": "uuid.UUID",
        "stream_type_template": "[16]uint8",
        "type_template": "uuid.UUID"
      },
      {
        "column": "created_at",
        "datatype": "timestamp with time zone",
        "table": "rule",
        "pos": 2,
        "typeid": "1184",
        "typelen": 8,
        "typemod": -1,
        "notnull": true,
        "hasdefault": true,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "31457",
        "zero_type": "0001-01-01T00:00:00Z",
        "query_type_template": "time.Time",
        "stream_type_template": "time.Time",
        "type_template": "time.Time"
      },
      {
        "column": "updated_at",
        "datatype": "timestamp with time zone",
        "table": "rule",
        "pos": 3,
        "typeid": "1184",
        "typelen": 8,
        "typemod": -1,
        "notnull": true,
        "hasdefault": true,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "31457",
        "zero_type": "0001-01-01T00:00:00Z",
        "query_type_template": "time.Time",
        "stream_type_template": "time.Time",
        "type_template": "time.Time"
      },
      {
        "column": "deleted_at",
        "datatype": "timestamp with time zone",
        "table": "rule",
        "pos": 4,
        "typeid": "1184",
        "typelen": 8,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "31457",
        "zero_type": "0001-01-01T00:00:00Z",
        "query_type_template": "time.Time",
        "stream_type_template": "time.Time",
        "type_template": "time.Time"
      },
      {
        "column": "branch",
        "datatype": "text",
        "table": "rule",
        "pos": 5,
        "typeid": "25",
        "typelen": -1,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "31457",
        "zero_type": "",
        "query_type_template": "string",
        "stream_type_template": "string",
        "type_template": "string"
      },
      {
        "column": "repository_id",
        "datatype": "uuid",
        "table": "rule",
        "pos": 6,
        "typeid": "2950",
        "typelen": 16,
        "typemod": -1,
        "notnull": true,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": "repository",
        "fcolumn": "id",
        "parent_id": "31457",
        "zero_type": "00000000-0000-0000-0000-000000000000",
        "query_type_template": "uuid.UUID",
        "stream_type_template": "[16]uint8",
        "type_template": "uuid.UUID"
      },
      {
        "column": "job_trigger_job_id",
        "datatype": "uuid",
        "table": "rule",
        "pos": 7,
        "typeid": "2950",
        "typelen": 16,
        "typemod": -1,
        "notnull": true,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": "job",
        "fcolumn": "id",
        "parent_id": "31457",
        "zero_type": "00000000-0000-0000-0000-000000000000",
        "query_type_template": "uuid.UUID",
        "stream_type_template": "[16]uint8",
        "type_template": "uuid.UUID"
      }
    ]
  },
  "task": {
    "tablename": "task",
    "oid": "31483",
    "schema": "test",
    "reltuples": -1,
    "relkind": "r",
    "relam": "2",
    "relacl": null,
    "reltype": "31485",
    "relowner": "10",
    "relhasindex": true,
    "columns": [
      {
        "column": "id",
        "datatype": "uuid",
        "table": "task",
        "pos": 1,
        "typeid": "2950",
        "typelen": 16,
        "typemod": -1,
        "notnull": true,
        "hasdefault": true,
        "hasmissing": false,
        "ispkey": true,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "31483",
        "zero_type": "00000000-0000-0000-0000-000000000000",
        "query_type_template": "uuid.UUID",
        "stream_type_template": "[16]uint8",
        "type_template": "uuid.UUID"
      },
      {
        "column": "created_at",
        "datatype": "timestamp with time zone",
        "table": "task",
        "pos": 2,
        "typeid": "1184",
        "typelen": 8,
        "typemod": -1,
        "notnull": true,
        "hasdefault": true,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "31483",
        "zero_type": "0001-01-01T00:00:00Z",
        "query_type_template": "time.Time",
        "stream_type_template": "time.Time",
        "type_template": "time.Time"
      },
      {
        "column": "updated_at",
        "datatype": "timestamp with time zone",
        "table": "task",
        "pos": 3,
        "typeid": "1184",
        "typelen": 8,
        "typemod": -1,
        "notnull": true,
        "hasdefault": true,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "31483",
        "zero_type": "0001-01-01T00:00:00Z",
        "query_type_template": "time.Time",
        "stream_type_template": "time.Time",
        "type_template": "time.Time"
      },
      {
        "column": "deleted_at",
        "datatype": "timestamp with time zone",
        "table": "task",
        "pos": 4,
        "typeid": "1184",
        "typelen": 8,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "31483",
        "zero_type": "0001-01-01T00:00:00Z",
        "query_type_template": "time.Time",
        "stream_type_template": "time.Time",
        "type_template": "time.Time"
      },
      {
        "column": "name",
        "datatype": "text",
        "table": "task",
        "pos": 5,
        "typeid": "25",
        "typelen": -1,
        "typemod": -1,
        "notnull": true,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "31483",
        "zero_type": "",
        "query_type_template": "string",
        "stream_type_template": "string",
        "type_template": "string"
      },
      {
        "column": "index",
        "datatype": "integer",
        "table": "task",
        "pos": 6,
        "typeid": "23",
        "typelen": 4,
        "typemod": -1,
        "notnull": true,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "31483",
        "zero_type": 0,
        "query_type_template": "int64",
        "stream_type_template": "int64",
        "type_template": "int64"
      },
      {
        "column": "platform",
        "datatype": "text",
        "table": "task",
        "pos": 7,
        "typeid": "25",
        "typelen": -1,
        "typemod": -1,
        "notnull": true,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "31483",
        "zero_type": "",
        "query_type_template": "string",
        "stream_type_template": "string",
        "type_template": "string"
      },
      {
        "column": "image",
        "datatype": "text",
        "table": "task",
        "pos": 8,
        "typeid": "25",
        "typelen": -1,
        "typemod": -1,
        "notnull": true,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "31483",
        "zero_type": "",
        "query_type_template": "string",
        "stream_type_template": "string",
        "type_template": "string"
      },
      {
        "column": "script",
        "datatype": "text",
        "table": "task",
        "pos": 9,
        "typeid": "25",
        "typelen": -1,
        "typemod": -1,
        "notnull": true,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "31483",
        "zero_type": "",
        "query_type_template": "string",
        "stream_type_template": "string",
        "type_template": "string"
      },
      {
        "column": "job_id",
        "datatype": "uuid",
        "table": "task",
        "pos": 10,
        "typeid": "2950",
        "typelen": 16,
        "typemod": -1,
        "notnull": true,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": "job",
        "fcolumn": "id",
        "parent_id": "31483",
        "zero_type": "00000000-0000-0000-0000-000000000000",
        "query_type_template": "uuid.UUID",
        "stream_type_template": "[16]uint8",
        "type_template": "uuid.UUID"
      }
    ]
  }
}`)

var tableByName introspect.TableByName

func init() {
	mu.Lock()
	defer mu.Unlock()

	err := json.Unmarshal(tableByNameAsJSON, &tableByName)
	if err != nil {
		panic(fmt.Errorf("failed to unmarshal tableByNameAsJSON into introspect.TableByName; %v", err))
	}
}

func RunServer(
	ctx context.Context,
	changes chan server.Change,
	addr string,
	db *pgxpool.Pool,
	redisPool *redis.Pool,
	httpMiddlewares []server.HTTPMiddleware,
	objectMiddlewares []server.ObjectMiddleware,
	addCustomHandlers func(chi.Router) error,
) error {
	mu.Lock()
	thisTableByName := tableByName
	mu.Unlock()

	return server.RunServer(ctx, changes, addr, NewFromItem, GetRouter, db, redisPool, httpMiddlewares, objectMiddlewares, addCustomHandlers, thisTableByName)
}
