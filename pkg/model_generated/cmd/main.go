package main

import (
	"log"
	"os"
	"strings"

	"github.com/initialed85/djangolang/pkg/model_generated"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("first argument must be command (one of 'serve', 'dump-openapi-json', 'dump-openapi-yaml')")
	}

	command := strings.TrimSpace(strings.ToLower(os.Args[1]))

	switch command {

	case "dump-openapi-json":
		model_generated.RunDumpOpenAPIJSON()

	case "dump-openapi-yaml":
		model_generated.RunDumpOpenAPIYAML()

	case "serve":
		model_generated.RunServeWithEnvironment(nil, nil, nil)
	}
}
