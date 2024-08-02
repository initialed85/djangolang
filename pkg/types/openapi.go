package types

import (
	"fmt"

	"gopkg.in/yaml.v2"
)

type In string

const (
	InQuery In = "query"
	InPath  In = "path"
)

type Type string

const (
	TypeOfAny     Type = ""
	TypeOfBoolean Type = "boolean"
	TypeOfString  Type = "string"
	TypeOfNumber  Type = "number"
	TypeOfInteger Type = "integer"
	TypeOfArray   Type = "array"
	TypeOfObject  Type = "object"
)

type Format string

const (
	// as per OAS 3.0
	FormatOfInt32    Format = "int32" // signed 32 bits
	FormatOfInt64    Format = "int64" // signed 64 bits (a.k.a long)
	FormatOfFloat    Format = "float"
	FormatOfDouble   Format = "double"
	FormatOfByte     Format = "byte"      // base64 encoded characters
	FormatOfBinary   Format = "binary"    // any sequence of octets
	FormatOfDate     Format = "date"      // As defined by full-date - RFC3339
	FormatOfDateTime Format = "date-time" // As defined by date-time - RFC3339
	FormatOfPassword Format = "password"  // A hint to UIs to obscure input.

	// not per-spec but apparently anything is permitted / welcomed
	FormatOfUUID  Format = "uuid"
	FormatOfEmail Format = "email"
	FormatOfIPv4  Format = "ipv4"
)

type Info struct {
	Title   string `yaml:"title" json:"title"`
	Version string `yaml:"version" json:"version"`
}

type Schema struct {
	Ref                  string             `yaml:"$ref,omitempty" json:"$ref,omitempty"`
	Type                 Type               `yaml:"type,omitempty" json:"type,omitempty"`
	Format               Format             `yaml:"format,omitempty" json:"format,omitempty"`
	Nullable             bool               `yaml:"nullable,omitempty" json:"nullable,omitempty"`
	Properties           map[string]*Schema `yaml:"properties,omitempty" json:"properties,omitempty"`
	AdditionalProperties *Schema            `yaml:"additionalProperties,omitempty" json:"additionalProperties,omitempty"`
	Required             []string           `yaml:"required,omitempty" json:"required,omitempty"`
	Items                *Schema            `yaml:"items,omitempty" json:"items,omitempty"`
}

type MediaType struct {
	Schema *Schema `yaml:"schema" json:"schema"`
}

type Parameter struct {
	Name     string  `yaml:"name" json:"name"`
	In       In      `yaml:"in" json:"in"`
	Required bool    `yaml:"required" json:"required"`
	Schema   *Schema `yaml:"schema" json:"schema"`
}

type RequestBody struct {
	Ref      string                `yaml:"$ref,omitempty" json:"$ref,omitempty"`
	Content  map[string]*MediaType `yaml:"content" json:"content"`
	Required bool                  `yaml:"required,omitempty" json:"required,omitempty"`
}

type Response struct {
	Ref         string                `yaml:"$ref,omitempty" json:"$ref,omitempty"`
	Description string                `yaml:"description" json:"description"`
	Content     map[string]*MediaType `yaml:"content,omitempty" json:"content,omitempty"`
}

type Operation struct {
	Tags        []string             `yaml:"tags,omitempty" json:"tags,omitempty"`
	OperationID string               `yaml:"operationId,omitempty" json:"operationId,omitempty"`
	Parameters  []*Parameter         `yaml:"parameters,omitempty" json:"parameters,omitempty"`
	RequestBody *RequestBody         `yaml:"requestBody,omitempty" json:"requestBody,omitempty"`
	Responses   map[string]*Response `yaml:"responses" json:"responses"`
}

type Path struct {
	Get    *Operation `yaml:"get,omitempty" json:"get,omitempty"`
	Post   *Operation `yaml:"post,omitempty" json:"post,omitempty"`
	Put    *Operation `yaml:"put,omitempty" json:"put,omitempty"`
	Patch  *Operation `yaml:"patch,omitempty" json:"patch,omitempty"`
	Delete *Operation `yaml:"delete,omitempty" json:"delete,omitempty"`
}

type Components struct {
	Schemas   map[string]*Schema   `yaml:"schemas,omitempty" json:"schemas,omitempty"`
	Responses map[string]*Response `yaml:"responses,omitempty" json:"responses,omitempty"`
}

type OpenAPI struct {
	OpenAPI    string           `yaml:"openapi" json:"openapi"`
	Info       *Info            `yaml:"info" json:"info"`
	Paths      map[string]*Path `yaml:"paths" json:"paths"`
	Components *Components      `yaml:"components" json:"components"`
}

func (o *OpenAPI) String() string {
	b, err := yaml.Marshal(o)
	if err != nil {
		panic(fmt.Errorf("failed to marshal %#+v to JSON: %v", o, err))
	}

	return string(b)
}
