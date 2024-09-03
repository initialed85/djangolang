package model_reference

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/go-chi/chi/v5"
	"github.com/gomodule/redigo/redis"
	"github.com/initialed85/djangolang/pkg/helpers"
	"github.com/initialed85/djangolang/pkg/server"
	"github.com/jackc/pgx/v5/pgxpool"
	"gopkg.in/yaml.v2"
)

func RunDumpOpenAPIJSON() {
	openApi, err := GetOpenAPI()
	if err != nil {
		log.Fatalf("%v", err)
	}

	b, err := json.MarshalIndent(openApi, "", "  ")
	if err != nil {
		log.Fatalf("%v", err)
	}

	fmt.Printf("%v", string(b))
}

func RunDumpOpenAPIYAML() {
	openApi, err := GetOpenAPI()
	if err != nil {
		log.Fatalf("%v", err)
	}

	b, err := yaml.Marshal(openApi)
	if err != nil {
		log.Fatalf("%v", err)
	}

	fmt.Printf("%v", string(b))
}

func RunServeWithArguments(
	ctx context.Context,
	cancel context.CancelFunc,
	port uint16,
	db *pgxpool.Pool,
	redisPool *redis.Pool,
	httpMiddlewares []server.HTTPMiddleware,
	objectMiddlewares []server.ObjectMiddleware,
	addCustomHandlers func(chi.Router) error,
) {
	defer cancel()

	go func() {
		helpers.WaitForCtrlC(ctx)
		cancel()
	}()

	err := RunServer(ctx, nil, fmt.Sprintf("0.0.0.0:%v", port), db, redisPool, httpMiddlewares, objectMiddlewares, addCustomHandlers)
	if err != nil {
		log.Fatalf("%v", err)
	}
}

func RunServeWithEnvironment(
	httpMiddlewares []server.HTTPMiddleware,
	objectMiddlewares []server.ObjectMiddleware,
	addCustomHandlers func(chi.Router) error,
) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	port, err := helpers.GetPort()
	if err != nil {
		log.Fatalf("%v", err)
	}

	db, err := helpers.GetDBFromEnvironment(ctx)
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer func() {
		db.Close()
	}()

	go func() {
		helpers.WaitForCtrlC(ctx)
		cancel()
	}()

	redisURL, err := helpers.GetRedisURL()
	if err != nil {
		log.Fatalf("%v", err)
	}

	redisPool := &redis.Pool{
		DialContext: func(ctx context.Context) (redis.Conn, error) {
			return redis.DialURLContext(ctx, redisURL)
		},
		MaxIdle:         2,
		MaxActive:       100,
		IdleTimeout:     300,
		Wait:            false,
		MaxConnLifetime: 86400,
	}

	defer func() {
		_ = redisPool.Close()
	}()

	RunServeWithArguments(ctx, cancel, port, db, redisPool, nil, nil, nil)
}
