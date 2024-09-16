package main

import (
	"os"
	"strings"

	"github.com/initialed85/djangolang/pkg/config"
	"github.com/initialed85/djangolang/pkg/model_reference"
)

var log = model_reference.ThisLogger()

func main() {
	if len(os.Args) < 2 {
		log.Fatal("first argument must be command (one of 'dump-config', 'dump-openapi-json', 'dump-openapi-yaml' or 'serve')")
	}

	command := strings.TrimSpace(strings.ToLower(os.Args[1]))

	switch command {

	case "dump-config":
		config.DumpConfig()

	case "dump-openapi-json":
		model_reference.RunDumpOpenAPIJSON()

	case "dump-openapi-yaml":
		model_reference.RunDumpOpenAPIYAML()

	case "serve":
		model_reference.RunServeWithEnvironment(nil, nil, nil)
	}
}
