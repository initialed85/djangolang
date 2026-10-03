package openapi

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	_pluralize "github.com/gertd/go-pluralize"

	"github.com/chanced/caps"
	"github.com/initialed85/djangolang/pkg/config"
	"github.com/initialed85/djangolang/pkg/query"
	"github.com/initialed85/djangolang/pkg/server"
	"github.com/initialed85/djangolang/pkg/types"
	"github.com/initialed85/structmeta/pkg/introspect"
)

// this file contains contains a fairly limited implementation of the OpenAPI v3 structure; limited as in just enough to convey the information
// needed to generate clients for a Djangolang API

const (
	contentTypeApplicationJSON = "application/json"
)

const (
	statusCodeDefault = "default"
)

var (
	matchers = []string{
		"eq",
		"ne",
		"gt",
		"gte",
		"lt",
		"lte",
		"in",
		"notin",
		"contains",
		"notcontains",
		"isnull",
		"isnotnull",
		"isfalse",
		"istrue",
		"like",
		"notlike",
		"ilike",
		"notilike",
		"desc",
		"asc",
	}

	descriptionByMatcher = map[string]string{
		"eq":          "SQL = comparison",
		"ne":          "SQL != comparison",
		"gt":          "SQL > comparison, may not work with all column types",
		"gte":         "SQL >= comparison, may not work with all column types",
		"lt":          "SQL < comparison, may not work with all column types",
		"lte":         "SQL <= comparison, may not work with all column types",
		"in":          "SQL IN comparison, permits comma-separated values",
		"contains":    "SQL @> comparison",
		"notcontains": "SQL NOT @> comparison",
		"notin":       "SQL NOT IN comparison, permits comma-separated values",
		"isnull":      "SQL IS null comparison, value is ignored (presence of key is sufficient)",
		"isnotnull":   "SQL IS NOT null comparison, value is ignored (presence of key is sufficient)",
		"isfalse":     "SQL IS false comparison, value is ignored (presence of key is sufficient)",
		"istrue":      "SQL IS true comparison, value is ignored (presence of key is sufficient)",
		"like":        "SQL LIKE comparison, value is implicitly prefixed and suffixed with %",
		"notlike":     "SQL NOT LIKE comparison, value is implicitly prefixed and suffixed with %",
		"ilike":       "SQL ILIKE comparison, value is implicitly prefixed and suffixed with %",
		"notilike":    "SQL NOT ILIKE comparison, value is implicitly prefixed and suffixed with %",
		"desc":        "SQL ORDER BY _ DESC clause, value is ignored (presence of key is sufficient)",
		"asc":         "SQL ORDER BY _ ASC clause, value is ignored (presence of key is sufficient)",
	}

	ignoredValueByMatcher = map[string]struct{}{
		"isnull":    {},
		"isnotnull": {},
		"isfalse":   {},
		"istrue":    {},
		"desc":      {},
		"asc":       {},
	}
)

var (
	pluralize = _pluralize.NewClient()
)

func init() {
	converter, ok := caps.DefaultConverter.(caps.StdConverter)
	if !ok {
		panic(fmt.Sprintf("failed to cast %#+v to caps.StdConverter", caps.DefaultConverter))
	}

	converter.Set("Dob", "DOB")
	converter.Set("Cpo", "CPO")
	converter.Set("Mwh", "MWH")
	converter.Set("Kwh", "KWH")
	converter.Set("Wh", "WH")
	converter.Set("Json", "JSON")
	converter.Set("Jsonb", "JSONB")
	converter.Set("Mac", "MAC")
	converter.Set("Ip", "IP")
}

var genericPattern = regexp.MustCompile(`.*\.(.*)\[.*\.(.*)\]`)
var packageStructPattern = regexp.MustCompile(`.*\.(.*)`)

func getTypeName(object *introspect.Object) string {
	if object == nil {
		return "Unknown"
	}

	typeName := object.Name
	if object.Type != nil {
		typeName = object.Type.String()

		if strings.Contains(typeName, "[") {
			if object.StructFields != nil {
				rawTypeName := typeName
				matches := genericPattern.FindAllStringSubmatch(rawTypeName, -1)
				if len(matches) == 1 && len(matches[0]) == 3 {
					typeName = fmt.Sprintf("%s_with_generic_of_%s", matches[0][1], matches[0][2])
				}
			}
		} else if strings.Contains(typeName, ".") {
			rawTypeName := typeName
			matches := packageStructPattern.FindAllStringSubmatch(rawTypeName, -1)
			if len(matches) == 1 && len(matches[0]) == 2 {
				typeName = matches[0][1]
			}
		}
	}

	return typeName
}

func getRefName(object *introspect.Object) string {
	if object == nil {
		return "Unknown"
	}

	typeName := getTypeName(object)

	refName := ""

	if object.PointerValue != nil {
		refName += "nullable_" + getRefName(object.PointerValue)
	} else if object.SliceValue != nil {
		refName += "array_of_" + getRefName(object.SliceValue)
	} else if object.MapKey != nil && object.MapValue != nil {
		refName += "map_of_" + getRefName(object.MapKey) + "_" + getRefName(object.MapValue)
	} else if object.StructFields != nil {
		refName = typeName
	} else {
		refName = typeName
	}

	camelCaseRefName := caps.ToCamel(refName)

	return camelCaseRefName
}

func getSchemaRef(object *introspect.Object) string {
	refName := getRefName(object)

	return fmt.Sprintf("#/components/schemas/%v", refName)
}

func isPrimitive(schema *types.Schema) bool {
	if schema == nil {
		return false
	}

	switch schema.Type {
	case types.TypeOfBoolean, types.TypeOfString, types.TypeOfNumber, types.TypeOfInteger:
		return true
	}

	return false
}

func NewFromIntrospectedSchema(httpHandlerSummaries []server.HTTPHandlerSummary) (*types.OpenAPI, error) {
	apiRootForOpenAPI := config.APIRootForOpenAPI()

	endpointPrefix := strings.TrimRight(apiRootForOpenAPI, "/")

	_ = endpointPrefix

	o := types.OpenAPI{
		OpenAPI: "3.1.1",
		Info: &types.Info{
			Title:   "Djangolang",
			Version: "1.0",
		},
		Paths: make(map[string]*types.Path),
		Components: &types.Components{
			Schemas:   make(map[string]*types.Schema),
			Responses: make(map[string]*types.Response),
		},
	}

	existingSchemaByIntrospectIntrospectedObjectTypeName := make(map[string]*types.Schema)

	var getSchema func(*introspect.Object) (*types.Schema, error)

	getSchema = func(thisObject *introspect.Object) (*types.Schema, error) {
		var schema *types.Schema

		existingSchema, ok := existingSchemaByIntrospectIntrospectedObjectTypeName[thisObject.Name]
		if ok {
			return existingSchema, nil
		}

		defer func() {
			if schema != nil {
				if !isPrimitive(schema) {
					o.Components.Schemas[getRefName(thisObject)] = schema
				}
			}
		}()

		typeTemplate := thisObject.Name
		if thisObject.Type != nil {
			typeTemplate = thisObject.Type.String()
		}

		typeMeta, err := types.GetTypeMetaForTypeTemplate(typeTemplate)
		if err == nil {
			schema = typeMeta.GetOpenAPISchema()
			if schema != nil {
				existingSchemaByIntrospectIntrospectedObjectTypeName[thisObject.Name] = schema
				return schema, nil
			}
		}

		if schema == nil {
			if thisObject.PointerValue != nil {
				pointerValueSchema, err := getSchema(thisObject.PointerValue)
				if err != nil {
					return nil, err
				}

				if isPrimitive(pointerValueSchema) {
					schema = &types.Schema{
						Type:     pointerValueSchema.Type,
						Format:   pointerValueSchema.Format,
						Nullable: true,
					}
				} else {
					refName := getSchemaRef(thisObject.PointerValue)

					schema = &types.Schema{
						Ref:      refName,
						Nullable: true,
					}
				}
			} else if thisObject.SliceValue != nil {
				sliceValueSchema, err := getSchema(thisObject.SliceValue)
				if err != nil {
					return nil, err
				}
				_ = sliceValueSchema

				schema = &types.Schema{
					Type:     types.TypeOfArray,
					Nullable: true,
					Items:    sliceValueSchema,
				}
			} else if thisObject.MapKey != nil && thisObject.MapValue != nil {
				mapValueSchema, err := getSchema(thisObject.MapValue)
				if err != nil {
					return nil, err
				}
				_ = mapValueSchema

				schema = &types.Schema{
					Type:                 types.TypeOfObject,
					Nullable:             true,
					AdditionalProperties: mapValueSchema,
				}
			} else if thisObject.StructFields != nil {
				schema = &types.Schema{
					Type:       types.TypeOfObject,
					Properties: make(map[string]*types.Schema),
				}

				existingSchemaByIntrospectIntrospectedObjectTypeName[thisObject.Name] = schema

				for _, introspectedStructFieldObject := range thisObject.StructFields {
					structFieldSchema, err := getSchema(introspectedStructFieldObject.Object)
					if err != nil {
						return nil, err
					}
					_ = structFieldSchema

					field := introspectedStructFieldObject.Tag.Get("json")
					if field == "" {
						field = introspectedStructFieldObject.Field
					}

					if strings.Contains(field, ",") {
						field = strings.Split(field, ",")[0]
					}

					schema.Properties[field] = structFieldSchema
				}
			} else {
				switch thisObject.Zero().(type) {
				case string:
					schema = &types.Schema{
						Type: types.TypeOfString,
					}
				case uint8:
					schema = &types.Schema{
						Type:   types.TypeOfNumber,
						Format: types.FormatOfUInt8,
					}
				case uint16:
					schema = &types.Schema{
						Type:   types.TypeOfNumber,
						Format: types.FormatOfUInt16,
					}
				case uint, uint32:
					schema = &types.Schema{
						Type:   types.TypeOfNumber,
						Format: types.FormatOfUInt32,
					}
				case uint64:
					schema = &types.Schema{
						Type:   types.TypeOfNumber,
						Format: types.FormatOfUInt64,
					}
				case int8:
					schema = &types.Schema{
						Type:   types.TypeOfNumber,
						Format: types.FormatOfInt8,
					}
				case int16:
					schema = &types.Schema{
						Type:   types.TypeOfNumber,
						Format: types.FormatOfInt16,
					}
				case int, int32:
					schema = &types.Schema{
						Type:   types.TypeOfNumber,
						Format: types.FormatOfInt32,
					}
				case int64:
					schema = &types.Schema{
						Type:   types.TypeOfNumber,
						Format: types.FormatOfInt64,
					}
				case float32:
					schema = &types.Schema{
						Type:   types.TypeOfNumber,
						Format: types.FormatOfFloat,
					}
				case float64:
					schema = &types.Schema{
						Type:   types.TypeOfNumber,
						Format: types.FormatOfDouble,
					}
				}
			}
		}

		if schema == nil || (schema.Type == "" && schema.Ref == "") {
			return nil, fmt.Errorf("failed to work out schema for %#+v", thisObject)
		}

		if schema != nil {
			existingSchemaByIntrospectIntrospectedObjectTypeName[thisObject.Name] = schema
		}

		return schema, nil
	}

	introspectedObjects := make([]*introspect.Object, 0)

	for _, httpHandlerSummary := range httpHandlerSummaries {
		for _, introspectStructFieldObject := range httpHandlerSummary.PathParamsIntrospectedStructFieldObjectByKey {
			if introspectStructFieldObject.Object != nil {
				introspectedObjects = append(introspectedObjects, introspectStructFieldObject.Object)
			}
		}

		for _, introspectStructFieldObject := range httpHandlerSummary.QueryParamsIntrospectedStructFieldObjectByKey {
			if introspectStructFieldObject.Object != nil {
				introspectedObjects = append(introspectedObjects, introspectStructFieldObject.Object)
			}
		}

		if httpHandlerSummary.RequestIntrospectedObject != nil {
			introspectedObjects = append(introspectedObjects, httpHandlerSummary.RequestIntrospectedObject)
		}

		if httpHandlerSummary.ResponseIntrospectedObject != nil {
			introspectedObjects = append(introspectedObjects, httpHandlerSummary.ResponseIntrospectedObject)
		}
	}

	for _, introspectedObject := range introspectedObjects {
		_ = introspectedObject

		schema, err := getSchema(introspectedObject)
		if err != nil {
			return nil, err
		}

		if isPrimitive(schema) {
			continue
		}

		o.Components.Schemas[getRefName(introspectedObject)] = schema
	}

	for _, httpHandlerSummary := range httpHandlerSummaries {
		endpointName := caps.ToCamel(strings.ReplaceAll(strings.Trim(httpHandlerSummary.FullPath, "/"), "/", "_"))

		endpointTag := pluralize.Singular(caps.ToCamel(strings.Split(strings.ReplaceAll(strings.Trim(httpHandlerSummary.FullPath, "/"), "/", "_"), "_")[0]))

		if httpHandlerSummary.Builtin {
			if strings.Contains(httpHandlerSummary.FullPath, "/{primaryKey}") {
				endpointTag = strings.ReplaceAll(endpointTag, "PrimaryKey", "")
				endpointName = strings.ReplaceAll(endpointName, "PrimaryKey", "")
				endpointName = pluralize.Singular(endpointName)
			} else {
				endpointName = pluralize.Plural(endpointName)
			}
		}

		parameters := []*types.Parameter{}

		for key, structFieldObject := range httpHandlerSummary.PathParamsIntrospectedStructFieldObjectByKey {
			object := structFieldObject.Object
			if object.PointerValue != nil {
				object = structFieldObject.PointerValue
			}

			structFieldObjectSchema, err := getSchema(object)
			if err != nil {
				return nil, err
			}

			tag := key

			required := false
			for possibleKey := range httpHandlerSummary.RequiredPathParamKeys {
				if possibleKey == key {
					required = true
					break
				}
			}

			parameters = append(parameters, &types.Parameter{
				Name:        tag,
				In:          types.InPath,
				Required:    required,
				Schema:      structFieldObjectSchema,
				Description: fmt.Sprintf("Path parameter %s", tag),
			})
		}

		for key, structFieldObject := range httpHandlerSummary.QueryParamsIntrospectedStructFieldObjectByKey {
			object := structFieldObject.Object
			if object.PointerValue != nil {
				object = structFieldObject.PointerValue
			}

			structFieldObjectSchema, err := getSchema(object)
			if err != nil {
				return nil, err
			}

			tag := key

			required := false
			for possibleKey := range httpHandlerSummary.RequiredQueryParamKeys {
				if possibleKey == key {
					required = true
					break
				}
			}

			parameters = append(parameters, &types.Parameter{
				Name:        tag,
				In:          types.InQuery,
				Required:    required,
				Schema:      structFieldObjectSchema,
				Description: fmt.Sprintf("Query parameter %s", tag),
			})
		}

		if httpHandlerSummary.Builtin && len(httpHandlerSummary.QueryParamsIntrospectedStructFieldObjectByKey) == 0 {
			parameters = append(parameters, &types.Parameter{
				Name:     "limit",
				In:       types.InQuery,
				Required: false,
				Schema: &types.Schema{
					Type:   types.TypeOfInteger,
					Format: types.FormatOfInt32,
				},
				Description: "SQL LIMIT operator",
			})

			parameters = append(parameters, &types.Parameter{
				Name:     "offset",
				In:       types.InQuery,
				Required: false,
				Schema: &types.Schema{
					Type:   types.TypeOfInteger,
					Format: types.FormatOfInt32,
				},
				Description: "SQL OFFSET operator",
			})

			parameters = append(parameters, &types.Parameter{
				Name:     "depth",
				In:       types.InQuery,
				Required: false,
				Schema: &types.Schema{
					Type: types.TypeOfInteger,
				},
				Description: "Max recursion depth for loading foreign objects; default = 1\n\n(0 = recurse until graph cycle detected, 1 = this object only, 2 = this object + neighbours, 3 = this object + neighbours + their neighbours... etc)",
			})

			if httpHandlerSummary.BuiltinTable != nil {
				tableNames := make([]string, 0)
				handledTableNames := make(map[string]struct{})

				for _, column := range httpHandlerSummary.BuiltinTable.Columns {
					if column.ForeignColumn == nil {
						continue
					}

					tableName := column.ForeignColumn.TableName

					_, ok := handledTableNames[tableName]
					if ok {
						continue
					}

					tableNames = append(tableNames, tableName)
					handledTableNames[tableName] = struct{}{}

				}

				for _, tableName := range tableNames {
					parameters = append(parameters, &types.Parameter{
						Name:     fmt.Sprintf("%v__%v", tableName, "load"),
						In:       types.InQuery,
						Required: false,
						Schema: &types.Schema{
							Type: types.TypeOfString,
						},
						Description: "load the given directly related object, value is ignored (presence of key is sufficient)",
					})
				}

				tableNames = make([]string, 0)
				handledTableNames = make(map[string]struct{})

				for _, column := range httpHandlerSummary.BuiltinTable.ReferencedByColumns {
					if column.ForeignColumn == nil {
						continue
					}

					tableName := column.TableName

					_, ok := handledTableNames[tableName]
					if ok {
						continue
					}

					tableNames = append(tableNames, tableName)
					handledTableNames[tableName] = struct{}{}
				}

				for _, tableName := range tableNames {
					parameters = append(parameters, &types.Parameter{
						Name:     fmt.Sprintf("referenced_by_%v__%v", tableName, "load"),
						In:       types.InQuery,
						Required: false,
						Schema: &types.Schema{
							Type: types.TypeOfString,
						},
						Description: "load the given indirectly related objects, value is ignored (presence of key is sufficient)",
					})
				}
			}

			if httpHandlerSummary.BuiltinIntrospectedModelObject != nil {
				for _, introspectedStructFieldObject := range httpHandlerSummary.BuiltinIntrospectedModelObject.StructFields {
					for _, matcher := range matchers {
						schema, err := getSchema(introspectedStructFieldObject.Object)
						if err != nil {
							return nil, err
						}

						_, ignored := ignoredValueByMatcher[matcher]
						if ignored {
							schema = &types.Schema{
								Type: types.TypeOfString,
							}
						}

						if schema.Type == types.TypeOfObject || schema.Type == types.TypeOfArray || schema.Type == types.TypeOfAny {
							if !(matcher == "contains" || matcher == "notcontains") {
								continue
							}

							schema = types.GetOpenAPISchemaJSON()
						}

						if matcher == "isfalse" || matcher == "istrue" {
							if schema.Type != types.TypeOfBoolean {
								continue
							}
						}

						if matcher == "isnull" || matcher == "isnotnull" {
							if !schema.Nullable {
								continue
							}
						}

						if matcher == "like" || matcher == "notlike" || matcher == "ilike" || matcher == "notilike" {
							if schema.Type != types.TypeOfString {
								continue
							}
						}

						tag := introspectedStructFieldObject.Tag.Get("json")
						if tag == "" {
							tag = introspectedStructFieldObject.Field
						}

						if strings.Contains(tag, ",") {
							tag = strings.Split(tag, ",")[0]
						}

						parameters = append(parameters, &types.Parameter{
							Name:        fmt.Sprintf("%v__%v", tag, matcher),
							In:          types.InQuery,
							Required:    false,
							Schema:      schema,
							Description: descriptionByMatcher[matcher],
						})
					}
				}
			}
		}

		getRequest := func(method string) *types.RequestBody {
			if !(method == http.MethodPost || method == http.MethodPut || method == http.MethodPatch) {
				return nil
			}

			if httpHandlerSummary.RequestIsEmpty || httpHandlerSummary.RequestIntrospectedObject == nil {
				return nil
			}

			schemaRef := getSchemaRef(httpHandlerSummary.RequestIntrospectedObject)

			return &types.RequestBody{
				Content: map[string]*types.MediaType{
					contentTypeApplicationJSON: {
						Schema: &types.Schema{
							Ref: schemaRef,
						},
					},
				},
				Required: true,
			}
		}

		getSuccessResponse := func(method string) *types.Response {
			description := fmt.Sprintf("%v%v success", caps.ToCamel(method), endpointName)
			if httpHandlerSummary.ResponseIsEmpty || httpHandlerSummary.ResponseIntrospectedObject == nil {
				return &types.Response{
					Description: description,
				}
			}

			schemaRef := getSchemaRef(httpHandlerSummary.ResponseIntrospectedObject)

			return &types.Response{
				Description: description,
				Content: map[string]*types.MediaType{
					contentTypeApplicationJSON: {
						Schema: &types.Schema{
							Ref: schemaRef,
						},
					},
				},
			}
		}

		getErrorResponse := func(method string) *types.Response {
			description := fmt.Sprintf("%v%v failure", caps.ToCamel(method), endpointName)
			if httpHandlerSummary.ResponseIsEmpty {
				return &types.Response{
					Description: description,
				}
			}

			return &types.Response{
				Description: description,
				Content: map[string]*types.MediaType{
					contentTypeApplicationJSON: {
						Schema: &types.Schema{
							Type:     types.TypeOfObject,
							Nullable: false,
							Properties: map[string]*types.Schema{
								"status": {
									Type:   types.TypeOfInteger,
									Format: types.FormatOfInt32,
								},
								"success": {
									Type: types.TypeOfBoolean,
								},
								"error": {
									Type: types.TypeOfArray,
									Items: &types.Schema{
										Type: types.TypeOfString,
									},
								},
							},
							Required: []string{"status", "success", "error"},
						},
					},
				},
			}
		}

		getOperation := func() *types.Operation {
			return &types.Operation{
				Tags:        []string{endpointTag},
				OperationID: fmt.Sprintf("%v%v", caps.ToCamel(httpHandlerSummary.Method), endpointName),
				Parameters:  parameters,
				RequestBody: getRequest(httpHandlerSummary.Method),
				Responses: map[string]*types.Response{
					fmt.Sprintf("%v", httpHandlerSummary.Status): getSuccessResponse(httpHandlerSummary.Method),
					statusCodeDefault: getErrorResponse(httpHandlerSummary.Method),
				},
			}
		}

		fullPath := fmt.Sprintf("%s/%s", endpointPrefix, strings.TrimLeft(httpHandlerSummary.FullPath, "/"))
		if o.Paths[fullPath] == nil {
			o.Paths[fullPath] = &types.Path{}
		}

		switch httpHandlerSummary.Method {
		case http.MethodGet:
			o.Paths[fullPath].Get = getOperation()
		case http.MethodPost:
			o.Paths[fullPath].Post = getOperation()
		case http.MethodPut:
			o.Paths[fullPath].Put = getOperation()
		case http.MethodPatch:
			o.Paths[fullPath].Patch = getOperation()
		case http.MethodDelete:
			o.Paths[fullPath].Delete = getOperation()
		default:
			return nil, fmt.Errorf("unsupported method %s for %s", httpHandlerSummary.Method, httpHandlerSummary.FullPath)
		}

	}

	return &o, nil
}

// ---------------------------------------------------------------------------
// builtinParameters — builds the builtin query parameters (limit, offset,
// depth, load markers, column filters) for a single builtin handler.
// The parameterSource callback determines how relationships and filters
// are discovered: either from introspected table data or from QueryModel.
// ---------------------------------------------------------------------------

type parameterSource func(summary server.HTTPHandlerSummary) ([]*types.Parameter, error)

func builtinParameters(summary server.HTTPHandlerSummary, source parameterSource) ([]*types.Parameter, error) {
	return source(summary)
}

// introspectedRelationshipParams — builds load params from BuiltinTable.
func introspectedRelationshipParams(summary server.HTTPHandlerSummary) ([]*types.Parameter, error) {
	if summary.BuiltinTable == nil {
		return nil, nil
	}

	var params []*types.Parameter

	// Outbound FK loads: {table}__load
	handledTables := make(map[string]struct{})
	for _, col := range summary.BuiltinTable.Columns {
		if col.ForeignColumn == nil {
			continue
		}
		if _, ok := handledTables[col.ForeignColumn.TableName]; ok {
			continue
		}
		handledTables[col.ForeignColumn.TableName] = struct{}{}
		params = append(params, &types.Parameter{
			Name:        fmt.Sprintf("%v__%v", col.ForeignColumn.TableName, "load"),
			In:          types.InQuery,
			Required:    false,
			Schema:      &types.Schema{Type: types.TypeOfString},
			Description: "load the given directly related object, value is ignored (presence of key is sufficient)",
		})
	}

	// Inbound FK loads: referenced_by_{table}__load
	handledTables = make(map[string]struct{})
	for _, col := range summary.BuiltinTable.ReferencedByColumns {
		if col.ForeignColumn == nil {
			continue
		}
		if _, ok := handledTables[col.TableName]; ok {
			continue
		}
		handledTables[col.TableName] = struct{}{}
		params = append(params, &types.Parameter{
			Name:        fmt.Sprintf("referenced_by_%v__%v", col.TableName, "load"),
			In:          types.InQuery,
			Required:    false,
			Schema:      &types.Schema{Type: types.TypeOfString},
			Description: "load the given indirectly related objects, value is ignored (presence of key is sufficient)",
		})
	}

	return params, nil
}

// introspectedFilterParams — builds column filter params from BuiltinIntrospectedModelObject.
func introspectedFilterParams(summary server.HTTPHandlerSummary, getSchema func(*introspect.Object) (*types.Schema, error)) ([]*types.Parameter, error) {
	if summary.BuiltinIntrospectedModelObject == nil {
		return nil, nil
	}

	var params []*types.Parameter

	for _, sf := range summary.BuiltinIntrospectedModelObject.StructFields {
		obj := sf.Object
		if obj.PointerValue != nil {
			obj = obj.PointerValue
		}

		schema, err := getSchema(obj)
		if err != nil {
			return nil, err
		}

		tag := sf.Tag.Get("json")
		if tag == "" {
			tag = sf.Field
		}
		if strings.Contains(tag, ",") {
			tag = strings.Split(tag, ",")[0]
		}

		for _, matcher := range matchers {
			_, ignored := ignoredValueByMatcher[matcher]
			if ignored {
				schema = &types.Schema{Type: types.TypeOfString}
			}

			if schema.Type == types.TypeOfObject || schema.Type == types.TypeOfArray || schema.Type == types.TypeOfAny {
				if !(matcher == "contains" || matcher == "notcontains") {
					continue
				}
				schema = types.GetOpenAPISchemaJSON()
			}

			if matcher == "isfalse" || matcher == "istrue" {
				if schema.Type != types.TypeOfBoolean {
					continue
				}
			}
			if matcher == "isnull" || matcher == "isnotnull" {
				if !schema.Nullable {
					continue
				}
			}
			if matcher == "like" || matcher == "notlike" || matcher == "ilike" || matcher == "notilike" {
				if schema.Type != types.TypeOfString {
					continue
				}
			}

			params = append(params, &types.Parameter{
				Name:        fmt.Sprintf("%v__%v", tag, matcher),
				In:          types.InQuery,
				Required:    false,
				Schema:      schema,
				Description: descriptionByMatcher[matcher],
			})
		}
	}

	return params, nil
}

// modelQueryParams — builds load and filter params from QueryModel metadata.
func modelQueryParams(summary server.HTTPHandlerSummary, getSchema func(*introspect.Object) (*types.Schema, error)) ([]*types.Parameter, error) {
	qm, ok := summary.BuiltinModelObject.(query.QueryModel)
	if !ok {
		return nil, nil
	}

	info := qm.GetQueryInfo()

	var params []*types.Parameter

	// Load params from relationships
	for _, rel := range info.Relationships {
		var loadName string
		if rel.Direction == query.RelationDirectionInbound {
			loadName = rel.Name
		} else {
			loadName = rel.TargetTable
		}
		params = append(params, &types.Parameter{
			Name:        fmt.Sprintf("%v__%v", loadName, "load"),
			In:          types.InQuery,
			Required:    false,
			Schema:      &types.Schema{Type: types.TypeOfString},
			Description: "load the given " + string(rel.Direction) + " relationship, value is ignored (presence of key is sufficient)",
		})
	}

	// Filter params from column definitions
	for _, f := range info.Filters {
		// Look up the column's introspected type to get the schema
		// We need to find the corresponding struct field to get its Go type
		// For now, map data types to OpenAPI schemas directly
		schema := dataTypeToSchema(f.DataType, f.Nullable)

		for _, op := range f.Operators {
			matcher := filterOpToMatcher(op)
			if matcher == "" {
				continue
			}

			_, ignored := ignoredValueByMatcher[matcher]
			if ignored {
				schema = &types.Schema{Type: types.TypeOfString}
			}

			if schema.Type == types.TypeOfObject || schema.Type == types.TypeOfArray || schema.Type == types.TypeOfAny {
				if !(matcher == "contains" || matcher == "notcontains") {
					continue
				}
				schema = types.GetOpenAPISchemaJSON()
			}

			params = append(params, &types.Parameter{
				Name:        fmt.Sprintf("%v__%v", f.Name, matcher),
				In:          types.InQuery,
				Required:    false,
				Schema:      schema,
				Description: descriptionByMatcher[matcher],
			})
		}
	}

	return params, nil
}

// dataTypeToSchema maps a PostgreSQL data type to an OpenAPI schema.
func dataTypeToSchema(dataType string, nullable bool) *types.Schema {
	switch dataType {
	case "uuid":
		return &types.Schema{Type: types.TypeOfString, Format: types.FormatOfUUID, Nullable: nullable}
	case "text", "character varying", "character":
		return &types.Schema{Type: types.TypeOfString, Nullable: nullable}
	case "boolean":
		return &types.Schema{Type: types.TypeOfBoolean, Nullable: nullable}
	case "integer", "bigint", "smallint":
		return &types.Schema{Type: types.TypeOfInteger, Nullable: nullable}
	case "real", "double precision", "numeric", "float":
		return &types.Schema{Type: types.TypeOfNumber, Nullable: nullable}
	case "timestamp with time zone", "timestamp without time zone":
		return &types.Schema{Type: types.TypeOfString, Format: types.FormatOfDateTime, Nullable: nullable}
	case "interval":
		return &types.Schema{Type: types.TypeOfString, Nullable: nullable}
	case "json", "jsonb":
		return types.GetOpenAPISchemaJSON()
	case "hstore":
		return &types.Schema{Type: types.TypeOfObject, Nullable: nullable}
	case "bytea":
		return &types.Schema{Type: types.TypeOfString, Format: types.FormatOfByte, Nullable: nullable}
	default:
		return &types.Schema{Type: types.TypeOfString, Nullable: nullable}
	}
}

// filterOpToMatcher maps a query.FilterOperator to a matcher string.
func filterOpToMatcher(op query.FilterOperator) string {
	switch op {
	case query.FilterOpEq:
		return "eq"
	case query.FilterOpNe:
		return "ne"
	case query.FilterOpGt:
		return "gt"
	case query.FilterOpGte:
		return "gte"
	case query.FilterOpLt:
		return "lt"
	case query.FilterOpLte:
		return "lte"
	case query.FilterOpIn:
		return "in"
	case query.FilterOpNotIn:
		return "notin"
	case query.FilterOpContains:
		return "contains"
	case query.FilterOpNotContains:
		return "notcontains"
	case query.FilterOpIsNull:
		return "isnull"
	case query.FilterOpIsNotNull:
		return "isnotnull"
	case query.FilterOpIsFalse:
		return "isfalse"
	case query.FilterOpIsTrue:
		return "istrue"
	case query.FilterOpLike:
		return "like"
	case query.FilterOpNotLike:
		return "notlike"
	case query.FilterOpILike:
		return "ilike"
	case query.FilterOpNotILike:
		return "notilike"
	}
	return ""
}

// NewFromQueryModels — builds the OpenAPI spec using QueryModel metadata
// for builtin parameter discovery instead of walking BuiltinTable
// and BuiltinIntrospectedModelObject directly. Falls back to the
// introspected schema approach for models that don't implement QueryModel.
func NewFromQueryModels(httpHandlerSummaries []server.HTTPHandlerSummary) (*types.OpenAPI, error) {
	o := types.OpenAPI{
		OpenAPI: "3.1.1",
		Info: &types.Info{
			Title:   "Djangolang",
			Version: "1.0",
		},
		Paths: make(map[string]*types.Path),
		Components: &types.Components{
			Schemas:   make(map[string]*types.Schema),
			Responses: make(map[string]*types.Response),
		},
	}

	// Reuse the schema-building logic from NewFromIntrospectedSchema
	// by temporarily collecting the introspected objects we need.
	// We build the schemas the same way, but use QueryModel for params.
	existingSchemaByIntrospectIntrospectedObjectTypeName := make(map[string]*types.Schema)

	var getSchema func(*introspect.Object) (*types.Schema, error)

	getSchema = func(thisObject *introspect.Object) (*types.Schema, error) {
		var schema *types.Schema

		existingSchema, ok := existingSchemaByIntrospectIntrospectedObjectTypeName[thisObject.Name]
		if ok {
			return existingSchema, nil
		}

		defer func() {
			if schema != nil {
				if !isPrimitive(schema) {
					o.Components.Schemas[getRefName(thisObject)] = schema
				}
			}
		}()

		typeTemplate := thisObject.Name
		if thisObject.Type != nil {
			typeTemplate = thisObject.Type.String()
		}

		typeMeta, err := types.GetTypeMetaForTypeTemplate(typeTemplate)
		if err == nil {
			schema = typeMeta.GetOpenAPISchema()
			if schema != nil {
				existingSchemaByIntrospectIntrospectedObjectTypeName[thisObject.Name] = schema
				return schema, nil
			}
		}

		if schema == nil {
			if thisObject.PointerValue != nil {
				pointerValueSchema, err := getSchema(thisObject.PointerValue)
				if err != nil {
					return nil, err
				}
				if isPrimitive(pointerValueSchema) {
					schema = &types.Schema{Type: pointerValueSchema.Type, Format: pointerValueSchema.Format, Nullable: true}
				} else {
					schema = &types.Schema{Ref: getSchemaRef(thisObject.PointerValue), Nullable: true}
				}
			} else if thisObject.SliceValue != nil {
				sliceValueSchema, err := getSchema(thisObject.SliceValue)
				if err != nil {
					return nil, err
				}
				schema = &types.Schema{Type: types.TypeOfArray, Nullable: true, Items: sliceValueSchema}
			} else if thisObject.MapKey != nil && thisObject.MapValue != nil {
				mapValueSchema, err := getSchema(thisObject.MapValue)
				if err != nil {
					return nil, err
				}
				schema = &types.Schema{Type: types.TypeOfObject, Nullable: true, AdditionalProperties: mapValueSchema}
			} else if thisObject.StructFields != nil {
				schema = &types.Schema{Type: types.TypeOfObject, Properties: make(map[string]*types.Schema)}
				existingSchemaByIntrospectIntrospectedObjectTypeName[thisObject.Name] = schema
				for _, sf := range thisObject.StructFields {
					sfSchema, err := getSchema(sf.Object)
					if err != nil {
						return nil, err
					}
					field := sf.Tag.Get("json")
					if field == "" {
						field = sf.Field
					}
					if strings.Contains(field, ",") {
						field = strings.Split(field, ",")[0]
					}
					schema.Properties[field] = sfSchema
				}
			} else {
				switch thisObject.Zero().(type) {
				case string:
					schema = &types.Schema{Type: types.TypeOfString}
				case uint8:
					schema = &types.Schema{Type: types.TypeOfNumber, Format: types.FormatOfUInt8}
				case uint16:
					schema = &types.Schema{Type: types.TypeOfNumber, Format: types.FormatOfUInt16}
				case uint, uint32:
					schema = &types.Schema{Type: types.TypeOfNumber, Format: types.FormatOfUInt32}
				case uint64:
					schema = &types.Schema{Type: types.TypeOfNumber, Format: types.FormatOfUInt64}
				case int8:
					schema = &types.Schema{Type: types.TypeOfNumber, Format: types.FormatOfInt8}
				case int16:
					schema = &types.Schema{Type: types.TypeOfNumber, Format: types.FormatOfInt16}
				case int, int32:
					schema = &types.Schema{Type: types.TypeOfNumber, Format: types.FormatOfInt32}
				case int64:
					schema = &types.Schema{Type: types.TypeOfNumber, Format: types.FormatOfInt64}
				case float32:
					schema = &types.Schema{Type: types.TypeOfNumber, Format: types.FormatOfFloat}
				case float64:
					schema = &types.Schema{Type: types.TypeOfNumber, Format: types.FormatOfDouble}
				}
			}
		}

		if schema == nil || (schema.Type == "" && schema.Ref == "") {
			return nil, fmt.Errorf("failed to work out schema for %#+v", thisObject)
		}

		if schema != nil {
			existingSchemaByIntrospectIntrospectedObjectTypeName[thisObject.Name] = schema
		}

		return schema, nil
	}

	// Collect all introspected objects for schema building
	introspectedObjects := make([]*introspect.Object, 0)
	for _, summary := range httpHandlerSummaries {
		for _, sf := range summary.PathParamsIntrospectedStructFieldObjectByKey {
			if sf.Object != nil {
				introspectedObjects = append(introspectedObjects, sf.Object)
			}
		}
		for _, sf := range summary.QueryParamsIntrospectedStructFieldObjectByKey {
			if sf.Object != nil {
				introspectedObjects = append(introspectedObjects, sf.Object)
			}
		}
		if summary.RequestIntrospectedObject != nil {
			introspectedObjects = append(introspectedObjects, summary.RequestIntrospectedObject)
		}
		if summary.ResponseIntrospectedObject != nil {
			introspectedObjects = append(introspectedObjects, summary.ResponseIntrospectedObject)
		}
	}

	for _, obj := range introspectedObjects {
		_, err := getSchema(obj)
		if err != nil {
			return nil, err
		}
	}

	// Build paths using QueryModel for builtin parameters
	for _, summary := range httpHandlerSummaries {
		endpointName := caps.ToCamel(strings.ReplaceAll(strings.Trim(summary.FullPath, "/"), "/", "_"))
		endpointTag := pluralize.Singular(caps.ToCamel(strings.Split(strings.ReplaceAll(strings.Trim(summary.FullPath, "/"), "/", "_"), "_")[0]))

		if summary.Builtin {
			if strings.Contains(summary.FullPath, "/{primaryKey}") {
				endpointName = strings.ReplaceAll(endpointName, "PrimaryKey", "")
				endpointName = pluralize.Singular(endpointName)
			} else {
				endpointName = pluralize.Plural(endpointName)
			}
		}

		parameters := []*types.Parameter{}

		// Path params
		for key, sf := range summary.PathParamsIntrospectedStructFieldObjectByKey {
			obj := sf.Object
			if obj.PointerValue != nil {
				obj = obj.PointerValue
			}
			sfSchema, err := getSchema(obj)
			if err != nil {
				return nil, err
			}
			tag := key
			required := false
			for pk := range summary.RequiredPathParamKeys {
				if pk == key {
					required = true
					break
				}
			}
			parameters = append(parameters, &types.Parameter{
				Name:        tag,
				In:          types.InPath,
				Required:    required,
				Schema:      sfSchema,
				Description: fmt.Sprintf("Path parameter %s", tag),
			})
		}

		// Query params (non-builtin)
		for key, sf := range summary.QueryParamsIntrospectedStructFieldObjectByKey {
			obj := sf.Object
			if obj.PointerValue != nil {
				obj = obj.PointerValue
			}
			sfSchema, err := getSchema(obj)
			if err != nil {
				return nil, err
			}
			tag := key
			required := false
			for pk := range summary.RequiredQueryParamKeys {
				if pk == key {
					required = true
					break
				}
			}
			parameters = append(parameters, &types.Parameter{
				Name:        tag,
				In:          types.InQuery,
				Required:    required,
				Schema:      sfSchema,
				Description: fmt.Sprintf("Query parameter %s", tag),
			})
		}

		// Builtin params: limit, offset, depth, and model-driven params
		if summary.Builtin && len(summary.QueryParamsIntrospectedStructFieldObjectByKey) == 0 {
			parameters = append(parameters, &types.Parameter{
				Name:     "limit",
				In:       types.InQuery,
				Required: false,
				Schema:   &types.Schema{Type: types.TypeOfInteger, Format: types.FormatOfInt32},
				Description: "SQL LIMIT operator",
			})
			parameters = append(parameters, &types.Parameter{
				Name:     "offset",
				In:       types.InQuery,
				Required: false,
				Schema:   &types.Schema{Type: types.TypeOfInteger, Format: types.FormatOfInt32},
				Description: "SQL OFFSET operator",
			})
			parameters = append(parameters, &types.Parameter{
				Name:     "depth",
				In:       types.InQuery,
				Required: false,
				Schema:   &types.Schema{Type: types.TypeOfInteger},
				Description: "Max recursion depth for loading foreign objects; default = 1\n\n(0 = recurse until graph cycle detected, 1 = this object only, 2 = this object + neighbours, 3 = this object + neighbours + their neighbours... etc)",
			})

			// Try QueryModel first, fall back to introspected table
			var builtinParams []*types.Parameter
			var err error

			if _, ok := summary.BuiltinModelObject.(query.QueryModel); ok {
				builtinParams, err = modelQueryParams(summary, getSchema)
				if err != nil {
					return nil, err
				}
			}

			if builtinParams == nil && summary.BuiltinTable != nil {
				builtinParams, err = introspectedRelationshipParams(summary)
				if err != nil {
					return nil, err
				}
				filterParams, err := introspectedFilterParams(summary, getSchema)
				if err != nil {
					return nil, err
				}
				builtinParams = append(builtinParams, filterParams...)
			}

			parameters = append(parameters, builtinParams...)
		}

		getRequest := func(method string) *types.RequestBody {
			if !(method == http.MethodPost || method == http.MethodPut || method == http.MethodPatch) {
				return nil
			}
			if summary.RequestIsEmpty || summary.RequestIntrospectedObject == nil {
				return nil
			}
			return &types.RequestBody{
				Content: map[string]*types.MediaType{
					contentTypeApplicationJSON: {Schema: &types.Schema{Ref: getSchemaRef(summary.RequestIntrospectedObject)}},
				},
				Required: true,
			}
		}

		getSuccessResponse := func(method string) *types.Response {
			desc := fmt.Sprintf("%v%v success", caps.ToCamel(method), endpointName)
			if summary.ResponseIsEmpty || summary.ResponseIntrospectedObject == nil {
				return &types.Response{Description: desc}
			}
			return &types.Response{
				Description: desc,
				Content: map[string]*types.MediaType{
					contentTypeApplicationJSON: {Schema: &types.Schema{Ref: getSchemaRef(summary.ResponseIntrospectedObject)}},
				},
			}
		}

		getErrorResponse := func(method string) *types.Response {
			desc := fmt.Sprintf("%v%v failure", caps.ToCamel(method), endpointName)
			if summary.ResponseIsEmpty {
				return &types.Response{Description: desc}
			}
			return &types.Response{
				Description: desc,
				Content: map[string]*types.MediaType{
					contentTypeApplicationJSON: {
						Schema: &types.Schema{
							Type:     types.TypeOfObject,
							Nullable: false,
							Properties: map[string]*types.Schema{
								"status":    {Type: types.TypeOfInteger, Format: types.FormatOfInt32},
								"success":   {Type: types.TypeOfBoolean},
								"error":     {Type: types.TypeOfArray, Items: &types.Schema{Type: types.TypeOfString}},
							},
							Required: []string{"status", "success", "error"},
						},
					},
				},
			}
		}

		getOperation := func() *types.Operation {
			return &types.Operation{
				Tags:        []string{endpointTag},
				OperationID: fmt.Sprintf("%v%v", caps.ToCamel(summary.Method), endpointName),
				Parameters:  parameters,
				RequestBody: getRequest(summary.Method),
				Responses: map[string]*types.Response{
					fmt.Sprintf("%v", summary.Status): getSuccessResponse(summary.Method),
					statusCodeDefault:                 getErrorResponse(summary.Method),
				},
			}
		}

		fullPath := fmt.Sprintf("%s/%s", strings.TrimRight(config.APIRootForOpenAPI(), "/"), strings.TrimLeft(summary.FullPath, "/"))
		if o.Paths[fullPath] == nil {
			o.Paths[fullPath] = &types.Path{}
		}

		switch summary.Method {
		case http.MethodGet:
			o.Paths[fullPath].Get = getOperation()
		case http.MethodPost:
			o.Paths[fullPath].Post = getOperation()
		case http.MethodPut:
			o.Paths[fullPath].Put = getOperation()
		case http.MethodPatch:
			o.Paths[fullPath].Patch = getOperation()
		case http.MethodDelete:
			o.Paths[fullPath].Delete = getOperation()
		default:
			return nil, fmt.Errorf("unsupported method %s for %s", summary.Method, summary.FullPath)
		}
	}

	return &o, nil
}
