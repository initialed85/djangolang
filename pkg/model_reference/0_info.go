package model_reference

import (
	_ "embed"
)

//go:embed 0_meta.go
var BaseFileData string

//go:embed 0_app.go
var AppFileData string

//go:embed cmd/main.go
var CmdMainFileData string

//go:embed logical_things.go
var ReferenceFileData string

var ReferencePackageName = "model_reference"
var ReferenceTableName = "logical_things"
var ReferenceObjectName = "LogicalThing"
var ReferenceEndpointName = "logical-things"
