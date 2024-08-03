package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/initialed85/djangolang/pkg/helpers"
	"github.com/initialed85/djangolang/pkg/introspect"
	"github.com/initialed85/djangolang/pkg/template"
)

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

		modulePath, packageName := helpers.GetModulePathAndPackageName()

		templateDataByFileName, err := template.Template(
			tableByName,
			modulePath,
			packageName,
		)
		if err != nil {
			log.Fatalf("%v failed; err: %v", command, err)
		}

		wd, err := os.Getwd()
		if err != nil {
			log.Fatalf("%v failed; err: %v", command, err)
		}

		fullPath := path.Join(wd, "pkg", packageName)

		_ = os.RemoveAll(fullPath)

		err = os.MkdirAll(path.Join(fullPath, "cmd"), 0o777)
		if err != nil {
			log.Fatalf("%v failed; err: %v", command, err)
		}

		for fileName, templateData := range templateDataByFileName {
			log.Printf("writing out %v...", path.Join(fullPath, fileName))

			err = os.WriteFile(path.Join(fullPath, fileName), []byte(templateData), 0o777)
			if err != nil {
				log.Fatalf("%v failed; err: %v", command, err)
			}
		}

		err = os.MkdirAll(path.Join(fullPath, "bin"), 0o777)
		if err != nil {
			log.Fatalf("%v failed; err: %v", command, err)
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
			log.Fatalf("%v failed; err: %v", command, err)
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
			log.Fatalf("%v failed; err: %v", command, err)
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
			log.Fatalf("%v failed; err: %v", command, err)
		}
		log.Printf("%v", string(out))

		log.Printf("done.")

		fmt.Printf("\n")
		fmt.Printf("For code-level integration, a Go source package for your Djangolang API has been generated:\n\n    %v\n\n", path.Join(fullPath))
		fmt.Printf("For simply using the API server, a default entrypoint binary has been built within that package:\n\n    %v\n\n", path.Join(fullPath, "bin", packageName))
		fmt.Printf("Here are some example usages of the binary:\n\n")

		fmt.Printf("    %v dump-openapi-json\n\n", path.Join(fullPath, "bin", packageName))
		fmt.Printf("    %v dump-openapi-yaml\n\n", path.Join(fullPath, "bin", packageName))

		envVars := ""
		for _, k := range []string{"POSTGRES_USER", "POSTGRES_PASSWORD", "POSTGRES_HOST", "POSTGRES_PORT", "POSTGRES_DB", "POSTGRES_SSLMODE", "POSTGRES_SCHEMA"} {
			v := strings.TrimSpace(os.Getenv(k))
			if v == "" {
				continue
			}

			envVars += fmt.Sprintf("%v=%v ", k, v)
		}

		fmt.Printf("    %v%v serve\n\n", envVars, path.Join(fullPath, "bin", packageName))

		fmt.Printf("Invoke the binary with no arguments for information on usage.\n\n")
	default:
	}
}
