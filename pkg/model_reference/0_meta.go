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
	"github.com/initialed85/djangolang/pkg/helpers"
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
var profile = helpers.GetEnvironmentVariable("DJANGOLANG_PROFILE") == "1"

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
	openApi, err = openapi.NewFromIntrospectedSchema(allObjects)
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

	for _, m := range httpMiddlewares {
		r.Use(m)
	}

	mu.Lock()
	for pattern, getRouterFn := range getRouterFnByPattern {
		r.Mount(pattern, getRouterFn(db, redisPool, httpMiddlewares, objectMiddlewares, waitForChange))
	}
	mu.Unlock()

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
				return fmt.Errorf("db ping failed: %v", err)
			}

			redisConn, err := redisPool.GetContext(ctx)
			if err != nil {
				return fmt.Errorf("redis pool get failed: %v", err)
			}

			defer func() {
				_ = redisConn.Close()
			}()

			_, err = redisConn.Do("PING")
			if err != nil {
				return fmt.Errorf("redis ping failed: %v", err)
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
			helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
			return
		}

		helpers.HandleObjectsResponse(w, http.StatusOK, nil)
	})

	r.Get("/openapi.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-type", "application/json")

		openApi, err := GetOpenAPI()
		if err != nil {
			helpers.HandleErrorResponse(w, http.StatusInternalServerError, fmt.Errorf("failed to get OpenAPI schema: %v", err))
			return
		}

		b, err := json.MarshalIndent(openApi, "", "  ")
		if err != nil {
			helpers.HandleErrorResponse(w, http.StatusInternalServerError, fmt.Errorf("failed to get OpenAPI schema: %v", err))
			return
		}

		helpers.WriteResponse(w, http.StatusOK, b)
	})

	r.Get("/openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-type", "application/yaml")

		openApi, err := GetOpenAPI()
		if err != nil {
			helpers.HandleErrorResponse(w, http.StatusInternalServerError, fmt.Errorf("failed to get OpenAPI schema: %v", err))
			return
		}

		b, err := yaml.Marshal(openApi)
		if err != nil {
			helpers.HandleErrorResponse(w, http.StatusInternalServerError, fmt.Errorf("failed to get OpenAPI schema: %v", err))
			return
		}

		helpers.WriteResponse(w, http.StatusOK, b)
	})

	return r
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
	return server.RunServer(ctx, changes, addr, NewFromItem, GetRouter, db, redisPool, httpMiddlewares, objectMiddlewares, addCustomHandlers)
}

var tableByNameAsJSON = []byte(`{}`)

var tableByName introspect.TableByName

func init() {
	err := json.Unmarshal(tableByNameAsJSON, &tableByName)
	if err != nil {
		panic(fmt.Errorf("failed to unmarshal tableByNameAsJSON into introspect.TableByName: %v", err))
	}
}
