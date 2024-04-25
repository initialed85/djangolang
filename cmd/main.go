package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/initialed85/djangolang/pkg/helpers"
	"github.com/initialed85/djangolang/pkg/introspect"
	"github.com/initialed85/djangolang/pkg/some_db"
	"github.com/initialed85/djangolang/pkg/template"
	"github.com/initialed85/djangolang/pkg/types"
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

	case "template":
		err = template.Run(ctx)
		if err != nil {
			log.Fatalf("%v failed; err: %v", command, err)
		}
		return

	default:
	}

	// a bit of a buffer to avoid blocking the stream while we're processing changes
	changes := make(chan types.Change, 1024)

	go func() {
		var err error

		switch command {

		case "serve":
			// TODO: this depends on some_db being templated correctly
			err = some_db.RunServer(ctx, changes)
			if err != nil {
				log.Fatalf("%v failed; err: %v", command, err)
			}

		default:
			err = fmt.Errorf("unrecognized command: %v", command)
			if err != nil {
				log.Fatalf("%v failed; err: %v", command, err)
			}
		}
	}()

	runtime.Gosched()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case change := <-changes: // it's important to drain this channel
				_ = change
			}
		}
	}()

	runtime.Gosched()

	time.Sleep(time.Second * 1)

	helpers.WaitForCtrlC(ctx)
	cancel()

	time.Sleep(time.Second * 1)
}
