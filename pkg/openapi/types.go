package openapi

// this file contains contains a fairly limited implementation of the OpenAPI v3 structure; limited as in just enough to convey the information
// needed to generate clients for a Djangolang API

type Info struct {
	Title   string `json:"title"`
	Version string `json:"version"`
}

type In string

const (
	InQuery In = "query"
	InPath  In = "path"
)

type Parameter struct {
	Name     string `json:"name"`
	In       In     `json:"in"`
	Required bool   `json:"required"`
}

type Type string

const (
	TypeString  Type = "string"
	TypeNumber  Type = "number"
	TypeInteger Type = "integer"
	TypeObject  Type = "object"
	TypeArray   Type = "array"
	TypeNull    Type = "null"
	TypeBoolean Type = "boolean"
)

type Schema struct {
	Ref        string            `json:"$ref"`
	Type       Type              `json:"type"`
	Properties map[string]Schema `json:"properties"`
	Required   []string          `json:"required"`
	Items      *Schema           `json:"items"`
}

type MediaType struct {
	Schema Schema `json:"schema"`
}

type Response struct {
	Ref         string               `json:"$ref"`
	Description string               `json:"description"`
	Content     map[string]MediaType `json:"content"`
}

type Operation struct {
	Tags        []string         `json:"tags"`
	OperationID string           `json:"operationId"`
	Parameters  []Parameter      `json:"parameters"`
	Responses   map[int]Response `json:"responses"`
}

type Path struct {
	Get    Operation `json:"get"`
	Put    Operation `json:"put"`
	Post   Operation `json:"post"`
	Delete Operation `json:"delete"`
	Patch  Operation `json:"patch"`
}

type Components struct {
	Schemas   map[string]Schema   `json:"schemas"`
	Responses map[string]Response `json:"responses"`
}

type OpenAPI struct {
	OpenAPI    string          `json:"openapi"`
	Info       Info            `json:"info"`
	Paths      map[string]Path `json:"paths"`
	Components Components      `json:"components"`
}
