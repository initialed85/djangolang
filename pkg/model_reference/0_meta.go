package model_reference

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/gomodule/redigo/redis"
	"github.com/initialed85/djangolang/pkg/server"
	"github.com/jmoiron/sqlx"
)

var mu = new(sync.Mutex)
var newFromItemFnByTableName = make(map[string]func(map[string]any) (any, error))
var getRouterFnByPattern = make(map[string]func(*sqlx.DB, redis.Conn, ...func(http.Handler) http.Handler) chi.Router)

func register(
	tableName string,
	newFromItem func(map[string]any) (any, error),
	pattern string,
	getRouterFn func(*sqlx.DB, redis.Conn, ...func(http.Handler) http.Handler) chi.Router,
) {
	newFromItemFnByTableName[tableName] = newFromItem
	getRouterFnByPattern[pattern] = getRouterFn
}

func NewFromItem(tableName string, item map[string]any) (any, error) {
	mu.Lock()
	newFromItemFn, ok := newFromItemFnByTableName[tableName]
	mu.Unlock()

	if !ok {
		return nil, fmt.Errorf("table name %v not known", tableName)
	}

	return newFromItemFn(item)
}

func GetRouter(db *sqlx.DB, redisConn redis.Conn, middlewares ...func(http.Handler) http.Handler) chi.Router {
	r := chi.NewRouter()

	for _, m := range middlewares {
		r.Use(m)
	}

	mu.Lock()
	for pattern, getRouterFn := range getRouterFnByPattern {
		r.Mount(pattern, getRouterFn(db, redisConn))
	}
	mu.Unlock()

	return r
}

func GetHandlerFunc(db *sqlx.DB, redisConn redis.Conn, middlewares ...func(http.Handler) http.Handler) http.HandlerFunc {
	r := chi.NewRouter()

	for _, m := range middlewares {
		r.Use(m)
	}

	mu.Lock()
	for pattern, getRouterFn := range getRouterFnByPattern {
		r.Mount(pattern, getRouterFn(db, redisConn))
	}
	mu.Unlock()

	return r.ServeHTTP
}

func RunServer(ctx context.Context, changes chan server.Change, addr string, db *sqlx.DB, redisConn redis.Conn, middlewares ...func(http.Handler) http.Handler) error {
	return server.RunServer(ctx, changes, addr, NewFromItem, GetRouter, db, redisConn, middlewares...)
}
