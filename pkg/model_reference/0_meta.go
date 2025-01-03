package model_reference

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

type patternAndMutateRouterFn struct {
	pattern        string
	mutateRouterFn server.MutateRouterFn
}

var mu = new(sync.Mutex)
var newFromItemFnByTableName = make(map[string]func(map[string]any) (any, error))
var patternsAndMutateRouterFns = make([]patternAndMutateRouterFn, 0)
var allObjects = make([]any, 0)
var openApi *types.OpenAPI
var profile = config.Profile()
var schema = "public"

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
	getRouterFn server.MutateRouterFn,
) {
	allObjects = append(allObjects, object)
	newFromItemFnByTableName[tableName] = newFromItem
	patternsAndMutateRouterFns = append(patternsAndMutateRouterFns, patternAndMutateRouterFn{
		pattern:        pattern,
		mutateRouterFn: getRouterFn,
	})
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

func MutateRouter(r chi.Router, db *pgxpool.Pool, redisPool *redis.Pool, objectMiddlewares []server.ObjectMiddleware, waitForChange server.WaitForChange) {
	mu.Lock()
	patternsAndGetRouterFns := patternsAndMutateRouterFns
	mu.Unlock()

	for _, thisPatternAndGetRouterFn := range patternsAndGetRouterFns {
		thisPatternAndGetRouterFn.mutateRouterFn(r, db, redisPool, objectMiddlewares, waitForChange)
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
}

func getHTTPHandler[T any, S any, Q any, R any](method string, path string, status int, handle func(context.Context, T, S, Q, any) (R, error), modelObject any, table *introspect.Table) (*server.HTTPHandler[T, S, Q, R], error) {
	customHTTPHandler, err := server.GetHTTPHandler(method, path, status, handle)
	if err != nil {
		return nil, err
	}

	customHTTPHandler.Builtin = true
	customHTTPHandler.BuiltinModelObject = modelObject
	customHTTPHandler.BuiltinTable = table

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

var tableByNameAsJSON = []byte(`{}`)

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
	nodeNames ...string,
) error {
	mu.Lock()
	thisTableByName := tableByName
	mu.Unlock()

	return server.RunServer(ctx, changes, addr, NewFromItem, MutateRouter, db, redisPool, httpMiddlewares, objectMiddlewares, addCustomHandlers, thisTableByName, nodeNames...)
}
