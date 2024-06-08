package model_reference

import (
	_ "embed"
)

//go:embed 0_meta.go
var BaseFileData string

//go:embed logical_thing.go
var ReferenceFileData string

var ReferencePackageName = "model_reference"
var ReferenceTableName = "logical_things"
var ReferenceObjectName = "LogicalThing"
var ReferenceEndpointName = "logical-things"
