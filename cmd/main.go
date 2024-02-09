package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/initialed85/djangolang/pkg/generate"
	"github.com/initialed85/djangolang/pkg/introspect"
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
	case "generate":
		err = generate.Run(ctx)
	default:
		err = fmt.Errorf("unrecognized command: %v", command)
	}

	if err != nil {
		log.Fatal(err)
	}
}
