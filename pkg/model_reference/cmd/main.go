package main

import (
	"context"
	"fmt"
	"log"

	"github.com/initialed85/djangolang/pkg/helpers"
	"github.com/initialed85/djangolang/pkg/model_reference"
)

func main() {
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
		_ = db.Close()
	}()

	go func() {
		helpers.WaitForCtrlC(ctx)
		cancel()
	}()

	_ = port

	err = model_reference.RunServer(ctx, nil, fmt.Sprintf("0.0.0.0:%v", port), db)
	if err != nil {
		log.Fatalf("err: %v", err)
	}
}
