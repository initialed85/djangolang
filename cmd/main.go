package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/initialed85/djangolang/pkg/cameranator"
	"github.com/initialed85/djangolang/pkg/helpers"
	"github.com/initialed85/djangolang/pkg/introspect"
	"github.com/initialed85/djangolang/pkg/template"
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
		db, err := helpers.GetDBFromEnvironment(ctx)
		if err != nil {
			log.Fatalf("%v failed; err: %v", command, err)
		}
		defer func() {
			_ = db.Close()
		}()

		schema := helpers.GetSchema()

		tableByName, err := introspect.Introspect(ctx, db, schema)
		if err != nil {
			log.Fatalf("%v failed; err: %v", command, err)
		}

		packageName := helpers.GetPackageName()

		templateDataByFileName, err := template.Template(
			tableByName,
			packageName,
		)
		if err != nil {
			log.Fatalf("%v failed; err: %v", command, err)
		}

		wd, err := os.Getwd()
		if err != nil {
			log.Fatalf("%v failed; err: %v", command, err)
		}

		_ = os.MkdirAll(path.Join(wd, "pkg", packageName), 0o777)

		for fileName, templateData := range templateDataByFileName {
			err = os.WriteFile(path.Join(wd, "pkg", packageName, fileName), []byte(templateData), 0o777)
			if err != nil {
				log.Fatalf("%v failed; err: %v", command, err)
			}
		}

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

		_ = port

		err = cameranator.RunServer(ctx, nil, fmt.Sprintf("0.0.0.0:%v", port), db)
		if err != nil {
			log.Fatalf("%v failed; err: %v", command, err)
		}

	default:
	}
}
