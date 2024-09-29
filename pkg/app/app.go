package app

import (
	"context"
	"fmt"
	_log "log"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/initialed85/djangolang/pkg/config"
	"github.com/initialed85/djangolang/pkg/helpers"
	"github.com/initialed85/djangolang/pkg/introspect"
	"github.com/initialed85/djangolang/pkg/template"
)

var log = helpers.GetLogger("app")

func ThisLogger() *_log.Logger {
	return log
}

func Run() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if len(os.Args) < 2 {
		log.Fatal("first argument must be command (one of 'introspect' or 'template')")
	}

	command := strings.TrimSpace(strings.ToLower(os.Args[1]))

	var err error

	switch command {

	case "introspect":
		err = introspect.Run(ctx)
		if err != nil {
			log.Fatalf("%v failed; %v", command, err)
		}
		return

	case "template":
		db, err := config.GetDBFromEnvironment(ctx)
		if err != nil {
			log.Fatalf("%v failed; %v", command, err)
		}
		defer func() {
			db.Close()
		}()

		schema := config.GetSchema()

		tx, err := db.Begin(ctx)
		if err != nil {
			log.Fatalf("%v failed; %v", command, err)
		}

		defer func() {
			_ = tx.Rollback(ctx)
		}()

		tableByName, err := introspect.Introspect(ctx, tx, schema)
		if err != nil {
			log.Fatalf("%v failed; %v", command, err)
		}

		modulePath := config.ModulePath()

		packageName := config.PackageName()

		templateDataByFileName, err := template.Template(
			tableByName,
			modulePath,
			packageName,
		)
		if err != nil {
			log.Fatalf("%v failed; %v", command, err)
		}

		wd, err := os.Getwd()
		if err != nil {
			log.Fatalf("%v failed; %v", command, err)
		}

		fullPath := path.Join(wd, "pkg", packageName)

		_ = os.RemoveAll(fullPath)

		err = os.MkdirAll(path.Join(fullPath, "cmd"), 0o777)
		if err != nil {
			log.Fatalf("%v failed; %v", command, err)
		}

		for fileName, templateData := range templateDataByFileName {
			log.Printf("writing out %v...", path.Join(fullPath, fileName))

			err = os.WriteFile(path.Join(fullPath, fileName), []byte(templateData), 0o777)
			if err != nil {
				log.Fatalf("%v failed; %v", command, err)
			}
		}

		err = os.MkdirAll(path.Join(fullPath, "bin"), 0o777)
		if err != nil {
			log.Fatalf("%v failed; %v", command, err)
		}

		log.Printf("building %v into %v...", path.Join(fullPath, "cmd"), path.Join(fullPath, "bin", packageName))

		log.Printf("running go mod tidy...")
		cmd := exec.Command(
			"go",
			"mod",
			"tidy",
		)
		out, err := cmd.CombinedOutput()
		if err != nil {
			log.Printf("\n%v\n", string(out))
			log.Fatalf("%v failed; %v", command, err)
		}
		log.Printf("%v", string(out))

		log.Printf("running go get...")
		cmd = exec.Command(
			"go",
			"get",
			"./...",
		)
		out, err = cmd.CombinedOutput()
		if err != nil {
			log.Printf("\n%v\n", string(out))
			log.Fatalf("%v failed; %v", command, err)
		}
		log.Printf("%v", string(out))

		log.Printf("running go build...")
		cmd = exec.Command(
			"go",
			"build",
			"-o",
			path.Join(fullPath, "bin", packageName),
			path.Join(fullPath, "cmd"),
		)
		out, err = cmd.CombinedOutput()
		if err != nil {
			log.Printf("\n%v\n", string(out))
			log.Fatalf("%v failed; %v", command, err)
		}
		log.Printf("%v", string(out))

		log.Printf("done.")

		fmt.Printf("\n")
		fmt.Printf("For code-level integration, a Go source package for your Djangolang API has been generated:\n\n    %v\n\n", path.Join(fullPath))
		fmt.Printf("For simply using the API server, a default entrypoint binary has been built within that package:\n\n    %v\n\n", path.Join(fullPath, "bin", packageName))
		fmt.Printf("Here are some example usages of the binary:\n\n")

		fmt.Printf("    %v dump-config\n\n", path.Join(fullPath, "bin", packageName))
		fmt.Printf("    %v dump-openapi-json\n\n", path.Join(fullPath, "bin", packageName))
		fmt.Printf("    %v dump-openapi-yaml\n\n", path.Join(fullPath, "bin", packageName))

		envVars := ""

		envVars += fmt.Sprintf("POSTGRES_USER=%v ", config.PostgresUser())
		envVars += fmt.Sprintf("POSTGRES_PASSWORD=%v ", config.PostgresPassword())
		envVars += fmt.Sprintf("POSTGRES_HOST=%v ", config.PostgresHost())
		envVars += fmt.Sprintf("POSTGRES_PORT=%v ", config.PostgresPort())
		envVars += fmt.Sprintf("POSTGRES_DATABASe=%v ", config.PostgresDatabase())
		envVars += fmt.Sprintf("POSTGRES_SSLMODE=%v ", config.PostgresSSLMode())
		envVars += fmt.Sprintf("POSTGRES_SCHEMA=%v ", config.PostgresSchema())

		fmt.Printf("    %v%v serve\n\n", envVars, path.Join(fullPath, "bin", packageName))

		fmt.Printf("Invoke the binary with no arguments for information on usage.\n\n")
	default:
	}
}
