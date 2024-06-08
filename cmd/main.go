package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/initialed85/djangolang/pkg/helpers"
	"github.com/initialed85/djangolang/pkg/introspect"
	"github.com/initialed85/djangolang/pkg/model_generated"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if len(os.Args) < 2 {
		log.Fatal("first argument must be command")
	}

	command := strings.TrimSpace(strings.ToLower(os.Args[1]))

	var err error

	switch command {

	case "introspect":
		err = introspect.Run(ctx)
		if err != nil {
			log.Fatalf("%v failed; err: %v", command, err)
		}
		return

	case "serve":
		port, err := helpers.GetPort()
		if err != nil {
			log.Fatalf("%v failed; err: %v", command, err)
		}

		db, err := helpers.GetDBFromEnvironment(ctx)
		if err != nil {
			log.Fatalf("%v failed; err: %v", command, err)
		}
		defer func() {
			_ = db.Close()
		}()

		go func() {
			helpers.WaitForCtrlC(ctx)
			cancel()
		}()

		err = model_generated.RunServer(ctx, nil, fmt.Sprintf("0.0.0.0:%v", port), db)
		if err != nil {
			log.Fatalf("%v failed; err: %v", command, err)
		}

	default:
	}
}
