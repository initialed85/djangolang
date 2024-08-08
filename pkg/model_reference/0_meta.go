package model_reference

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/gomodule/redigo/redis"
	"github.com/initialed85/djangolang/pkg/helpers"
	"github.com/initialed85/djangolang/pkg/openapi"
	"github.com/initialed85/djangolang/pkg/server"
	"github.com/initialed85/djangolang/pkg/types"
	"github.com/jmoiron/sqlx"
	"gopkg.in/yaml.v2"
)

var mu = new(sync.Mutex)
var newFromItemFnByTableName = make(map[string]func(map[string]any) (any, error))
var getRouterFnByPattern = make(map[string]func(*sqlx.DB, *redis.Pool, []server.HTTPMiddleware, []server.ModelMiddleware) chi.Router)
var allObjects = make([]any, 0)
var openApi *types.OpenAPI

func register(
	tableName string,
	object any,
	newFromItem func(map[string]any) (any, error),
	pattern string,
	getRouterFn func(*sqlx.DB, *redis.Pool, []server.HTTPMiddleware, []server.ModelMiddleware) chi.Router,
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

func GetRouter(db *sqlx.DB, redisPool *redis.Pool, httpMiddlewares []server.HTTPMiddleware, modelMiddlewares []server.ModelMiddleware) chi.Router {
	r := chi.NewRouter()

	for _, m := range httpMiddlewares {
		r.Use(m)
	}

	mu.Lock()
	for pattern, getRouterFn := range getRouterFnByPattern {
		r.Mount(pattern, getRouterFn(db, redisPool, httpMiddlewares, modelMiddlewares))
	}
	mu.Unlock()

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
	db *sqlx.DB,
	redisPool *redis.Pool,
	httpMiddlewares []server.HTTPMiddleware,
	modelMiddlewares []server.ModelMiddleware,
) error {
	return server.RunServer(ctx, changes, addr, NewFromItem, GetRouter, db, redisPool, httpMiddlewares, modelMiddlewares)
}
