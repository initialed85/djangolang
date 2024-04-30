package model_reference

import (
	_ "embed"
)

//go:embed logical_thing.go
var ReferenceFileData string
var ReferencePackageName = "model_reference"
var ReferenceTableName = "logical_things"
var ReferenceObjectName = "LogicalThing"
