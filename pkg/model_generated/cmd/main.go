package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gomodule/redigo/redis"
	"gopkg.in/yaml.v2"

	"github.com/initialed85/djangolang/pkg/helpers"
	"github.com/initialed85/djangolang/pkg/model_generated"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if len(os.Args) < 2 {
		log.Fatal("first argument must be command (one of 'serve', 'dump-openapi-json', 'dump-openapi-yaml')")
	}

	command := strings.TrimSpace(strings.ToLower(os.Args[1]))

	switch command {

	case "serve":
		port, err := helpers.GetPort()
		if err != nil {
			log.Fatalf("err: %v", err)
		}

		db, err := helpers.GetDBFromEnvironment(ctx)
		if err != nil {
			log.Fatalf("err: %v", err)
		}
		defer func() {
			_ = db.Close()
		}()

		go func() {
			helpers.WaitForCtrlC(ctx)
			cancel()
		}()

		redisURL := helpers.GetRedisURL()
		var redisConn redis.Conn
		if redisURL != "" {
			redisConn, err = redis.DialURLContext(ctx, redisURL)
			if err != nil {
				log.Fatalf("err: %v", err)
			}
			defer func() {
				_ = redisConn.Close()
			}()
		}

		err = model_generated.RunServer(ctx, nil, fmt.Sprintf("0.0.0.0:%v", port), db, redisConn, nil, nil)
		if err != nil {
			log.Fatalf("err: %v", err)
		}

	case "dump-openapi-json":
		openApi, err := model_generated.GetOpenAPI()
		if err != nil {
			log.Fatalf("err: %v", err)
		}

		b, err := json.MarshalIndent(openApi, "", "  ")
		if err != nil {
			log.Fatalf("err: %v", err)
		}

		fmt.Printf("%v", string(b))
	case "dump-openapi-yaml":
		openApi, err := model_generated.GetOpenAPI()
		if err != nil {
			log.Fatalf("err: %v", err)
		}

		b, err := yaml.Marshal(openApi)
		if err != nil {
			log.Fatalf("err: %v", err)
		}

		fmt.Printf("%v", string(b))
	}
}
