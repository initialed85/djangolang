package model_generated

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
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
var getRouterFnByPattern = make(map[string]server.GetRouterFn)
var allObjects = make([]any, 0)
var openApi *types.OpenAPI

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

func GetRouter(db *sqlx.DB, redisPool *redis.Pool, httpMiddlewares []server.HTTPMiddleware, objectMiddlewares []server.ObjectMiddleware, waitForChange server.WaitForChange) chi.Router {
	actualRouter := chi.NewRouter()

	for _, m := range httpMiddlewares {
		actualRouter.Use(m)
	}

	mu.Lock()
	for pattern, getRouterFn := range getRouterFnByPattern {
		actualRouter.Mount(pattern, getRouterFn(db, redisPool, httpMiddlewares, objectMiddlewares, waitForChange))
	}
	mu.Unlock()

	actualRouter.Get("/openapi.json", func(w http.ResponseWriter, r *http.Request) {
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

	actualRouter.Get("/openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
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

	finalRouter := chi.NewRouter()

	apiRoot := helpers.GetEnvironmentVariableOrDefault("DJANGOLANG_API_ROOT", "/")

	finalRouter.Mount(fmt.Sprintf("/%s", strings.Trim(apiRoot, "/")), actualRouter)

	return finalRouter
}

func RunServer(
	ctx context.Context,
	changes chan server.Change,
	addr string,
	db *sqlx.DB,
	redisPool *redis.Pool,
	httpMiddlewares []server.HTTPMiddleware,
	objectMiddlewares []server.ObjectMiddleware,
) error {
	return server.RunServer(ctx, changes, addr, NewFromItem, GetRouter, db, redisPool, httpMiddlewares, objectMiddlewares)
}
