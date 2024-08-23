package model_generated

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/gomodule/redigo/redis"
	"github.com/initialed85/djangolang/pkg/helpers"
	"github.com/jackc/pgx/v5/pgxpool"
	"gopkg.in/yaml.v2"
)

func RunDumpOpenAPIJSON() {
	openApi, err := GetOpenAPI()
	if err != nil {
		log.Fatalf("err: %v", err)
	}

	b, err := json.MarshalIndent(openApi, "", "  ")
	if err != nil {
		log.Fatalf("err: %v", err)
	}

	fmt.Printf("%v", string(b))
}

func RunDumpOpenAPIYAML() {
	openApi, err := GetOpenAPI()
	if err != nil {
		log.Fatalf("err: %v", err)
	}

	b, err := yaml.Marshal(openApi)
	if err != nil {
		log.Fatalf("err: %v", err)
	}

	fmt.Printf("%v", string(b))
}

func RunServeWithArguments(
	ctx context.Context,
	cancel context.CancelFunc,
	port uint16,
	db *pgxpool.Pool,
	redisPool *redis.Pool,
) {
	defer cancel()

	go func() {
		helpers.WaitForCtrlC(ctx)
		cancel()
	}()

	err := RunServer(ctx, nil, fmt.Sprintf("0.0.0.0:%v", port), db, redisPool, nil, nil)
	if err != nil {
		log.Fatalf("err: %v", err)
	}
}

func RunServeWithEnvironment() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	port, err := helpers.GetPort()
	if err != nil {
		log.Fatalf("err: %v", err)
	}

	db, err := helpers.GetDBFromEnvironment(ctx)
	if err != nil {
		log.Fatalf("err: %v", err)
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
		log.Fatalf("err: %v", err)
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

	RunServeWithArguments(ctx, cancel, port, db, redisPool)
}
