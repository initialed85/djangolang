package openapi

import (
	"fmt"
	"net/http"
	"strings"

	_pluralize "github.com/gertd/go-pluralize"

	"github.com/chanced/caps"
	"github.com/initialed85/djangolang/pkg/config"
	"github.com/initialed85/djangolang/pkg/server"
	"github.com/initialed85/djangolang/pkg/types"
	"github.com/initialed85/structmeta/pkg/introspect"
)

// this file contains contains a fairly limited implementation of the OpenAPI v3 structure; limited as in just enough to convey the information
// needed to generate clients for a Djangolang API

type CustomHTTPHandlerSummary struct {
	PathParams  any
	QueryParams any
	Request     any
	Response    any
	Method      string
	Path        string
	Status      int
}

type EndpointVariant string

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
		"eq":        "SQL = comparison",
		"ne":        "SQL != comparison",
		"gt":        "SQL > comparison, may not work with all column types",
		"gte":       "SQL >= comparison, may not work with all column types",
		"lt":        "SQL < comparison, may not work with all column types",
		"lte":       "SQL <= comparison, may not work with all column types",
		"in":        "SQL IN comparison, permits comma-separated values",
		"notin":     "SQL NOT IN comparison, permits comma-separated values",
		"isnull":    "SQL IS null comparison, value is ignored (presence of key is sufficient)",
		"isnotnull": "SQL IS NOT null comparison, value is ignored (presence of key is sufficient)",
		"isfalse":   "SQL IS false comparison, value is ignored (presence of key is sufficient)",
		"istrue":    "SQL IS true comparison, value is ignored (presence of key is sufficient)",
		"like":      "SQL LIKE comparison, value is implicitly prefixed and suffixed with %",
		"notlike":   "SQL NOT LIKE comparison, value is implicitly prefixed and suffixed with %",
		"ilike":     "SQL ILIKE comparison, value is implicitly prefixed and suffixed with %",
		"notilike":  "SQL NOT ILIKE comparison, value is implicitly prefixed and suffixed with %",
		"desc":      "SQL ORDER BY _ DESC clause, value is ignored (presence of key is sufficient)",
		"asc":       "SQL ORDER BY _ ASC clause, value is ignored (presence of key is sufficient)",
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

func getTypeName(object *introspect.Object) string {
	typeName := object.Name
	if object.Type != nil {
		typeName = object.Type.String()
	}

	return typeName
}

func getRefName(object *introspect.Object) string {
	refName := getTypeName(object)

	isArray := false
	if strings.HasPrefix(refName, "[]") || strings.HasPrefix(refName, "*[]") {
		isArray = true
		refName = strings.ReplaceAll(refName, "[]", "")
	}

	isNullable := false
	if strings.HasPrefix(refName, "*") {
		isNullable = true
	}

	parts := strings.Split(refName, ".")
	refName = parts[len(parts)-1]

	if isArray {
		refName = "array_of_" + refName
	}

	if isNullable {
		refName = "nullable_" + refName
	}

	refName = strings.ReplaceAll(refName, "interface {}", "any")

	if object.MapKey != nil && object.MapValue != nil {
		refName = strings.ReplaceAll(refName, "map[", "map_of_")
		refName = strings.ReplaceAll(refName, "]", "_")
	}

	refName = caps.ToCamel(refName)

	return refName
}

func getSchemaRef(object *introspect.Object) string {
	refName := getRefName(object)

	return fmt.Sprintf("#/components/schemas/%v", refName)
}

func isPrimitive(schema *types.Schema) bool {
	switch schema.Type {
	case types.TypeOfBoolean, types.TypeOfString, types.TypeOfNumber, types.TypeOfInteger:
		return true
	}

	return false
}

func NewFromIntrospectedSchema(inputObjects []any, customHTTPHandlerSummaries []CustomHTTPHandlerSummary) (*types.OpenAPI, error) {
	apiRootForOpenAPI := config.APIRootForOpenAPI()

	endpointPrefix := strings.TrimRight(apiRootForOpenAPI, "/")

	o := types.OpenAPI{
		OpenAPI: "3.0.0",
		Info: &types.Info{
			Title:   "Djangolang",
			Version: "1.0",
		},
		// TODO
		// Servers: []types.Server{
		// 	{
		// 		URL: fmt.Sprintf("/%s", strings.Trim(apiRoot, "/")),
		// 	},
		// 	{
		// 		URL: fmt.Sprintf("http://localhost:7070/%s", strings.Trim(apiRoot, "/")),
		// 	},
		// },
		Paths: map[string]*types.Path{},
		Components: &types.Components{
			Schemas:   map[string]*types.Schema{},
			Responses: map[string]*types.Response{},
		},
	}

	introspectedObjectByName := make(map[string]*introspect.Object)
	introspectedObjects := make([]*introspect.Object, 0)

	for i, inputObject := range inputObjects {
		introspectedObject, err := introspect.Introspect(inputObject)
		if err != nil {
			return nil, fmt.Errorf("failed to introspect %#+v; %v", inputObject, err)
		}

		if len(introspectedObject.StructFields) == 0 {
			return nil, fmt.Errorf(
				"%v: %#+v has no struct fields- that doesn't seem right (it should be a Djangolang object representing a database table)",
				i, inputObject,
			)
		}

		typeName := strings.TrimLeft(getTypeName(introspectedObject), "*")

		introspectedObjectByName[typeName] = introspectedObject
		introspectedObjectByName["*"+typeName] = introspectedObject
		introspectedObjectByName["[]*"+typeName] = introspectedObject
		introspectedObjects = append(introspectedObjects, introspectedObject)
	}

	o.Components = &types.Components{
		Schemas: make(map[string]*types.Schema),
	}

	var getSchema func(*introspect.Object) (*types.Schema, error)

	getSchema = func(thisObject *introspect.Object) (*types.Schema, error) {
		var schema *types.Schema
		var err error

		isBehindPointer := thisObject.PointerValue != nil
		typeName := getTypeName(thisObject)
		_, isDjangolangObject := introspectedObjectByName[typeName]

		defer func() {
			if schema != nil {
				if !schema.Nullable {
					if schema.Type == types.TypeOfObject {
						o.Components.Schemas[getRefName(thisObject)] = schema

						if !schema.Nullable && isDjangolangObject {
							o.Components.Schemas["Nullable"+getRefName(thisObject)] = &types.Schema{
								Ref:      getSchemaRef(thisObject),
								Nullable: true,
							}

							o.Components.Schemas["ArrayOf"+getRefName(thisObject)] = &types.Schema{
								Type: types.TypeOfArray,
								Items: &types.Schema{
									Ref: getSchemaRef(thisObject),
								},
							}

							o.Components.Schemas["NullableArrayOf"+getRefName(thisObject)] = &types.Schema{
								Nullable: true,
								Ref:      strings.ReplaceAll(getSchemaRef(thisObject), "schemas/", "schemas/ArrayOf"),
							}
						}
					}
				} else {
					if schema.Type == types.TypeOfObject {
						o.Components.Schemas[getRefName(thisObject)] = schema

						if schema.Nullable && isDjangolangObject {
							o.Components.Schemas[getRefName(thisObject.PointerValue)] = &types.Schema{
								Ref:      getSchemaRef(thisObject.PointerValue),
								Nullable: true,
							}

							o.Components.Schemas["ArrayOf"+getRefName(thisObject.PointerValue)] = &types.Schema{
								Type:     types.TypeOfArray,
								Nullable: true,
								Items: &types.Schema{
									Ref: getSchemaRef(thisObject.PointerValue),
								},
							}
						}
					}
				}
			}
		}()

		if isBehindPointer {
			actualSchema, err := getSchema(thisObject.PointerValue)
			if err != nil {
				return nil, err
			}

			if actualSchema != nil {
				isNullable := true
				if actualSchema.Format == "" {
					isNullable = false
				}

				schema = &types.Schema{
					Ref:                  actualSchema.Ref,
					Type:                 actualSchema.Type,
					Format:               actualSchema.Format,
					Nullable:             isNullable,
					Properties:           actualSchema.Properties,
					AdditionalProperties: actualSchema.AdditionalProperties,
					Required:             actualSchema.Required,
					Items:                actualSchema.Items,
				}

				return schema, nil
			}
		}

		if isDjangolangObject {
			schema = &types.Schema{
				Type:       types.TypeOfObject,
				Nullable:   isBehindPointer,
				Properties: make(map[string]*types.Schema),
			}

			for _, structFieldObject := range thisObject.StructFields {
				var structFieldSchema *types.Schema

				structFieldTypeName := getTypeName(structFieldObject.Object)

				_, structFieldIsDjangolangObject := introspectedObjectByName[structFieldTypeName]
				if structFieldIsDjangolangObject {
					if structFieldObject.SliceValue != nil {
						structFieldSchema = &types.Schema{
							Ref: getSchemaRef(structFieldObject.Object),
						}
					} else {
						structFieldSchema = &types.Schema{
							Ref: getSchemaRef(structFieldObject.Object),
						}
					}
				} else {
					structFieldSchema, err = getSchema(structFieldObject.Object)
					if err != nil {
						return nil, err
					}
				}

				tag := structFieldObject.Tag.Get("json")
				if tag == "" {
					tag = structFieldObject.Field
				}

				schema.Properties[tag] = structFieldSchema
			}

			return schema, nil
		}

		typeTemplate := thisObject.Name
		if thisObject.Type != nil {
			typeTemplate = thisObject.Type.String()
		}

		typeMeta, err := types.GetTypeMetaForTypeTemplate(typeTemplate)
		if err != nil {
			var schema *types.Schema

			if thisObject.StructFields != nil {
				schema = &types.Schema{
					Type:       types.TypeOfObject,
					Nullable:   isBehindPointer,
					Properties: make(map[string]*types.Schema),
					Required:   make([]string, 0),
				}

				for _, structFieldObject := range thisObject.StructFields {
					var structFieldSchema *types.Schema

					structFieldSchema, err = getSchema(structFieldObject.Object)
					if err != nil {
						return nil, err
					}

					tag := structFieldObject.Tag.Get("json")
					if tag == "" {
						tag = structFieldObject.Field
					}

					refName := getRefName(structFieldObject.Object)
					schemaRef := getSchemaRef(structFieldObject.Object)
					o.Components.Schemas[refName] = structFieldSchema
					schema.Properties[tag] = &types.Schema{
						Ref: schemaRef,
					}

					if structFieldObject.PointerValue == nil {
						schema.Required = append(schema.Required, tag)
					}
				}
			}

			if thisObject.SliceValue != nil {
				itemsSchema, err := getSchema(thisObject.SliceValue)
				if err != nil {
					return nil, err
				}

				schema = &types.Schema{
					Type:     types.TypeOfArray,
					Nullable: true,
				}

				refName := getRefName(thisObject.SliceValue)
				schemaRef := getSchemaRef(thisObject.SliceValue)
				o.Components.Schemas[refName] = itemsSchema

				schema.Items = &types.Schema{
					Ref: schemaRef,
				}
			}

			if thisObject.MapKey != nil && thisObject.MapValue != nil {
				additionalPropertiesSchema, err := getSchema(thisObject.MapValue)
				if err != nil {
					return nil, err
				}

				schema = &types.Schema{
					Type:     types.TypeOfObject,
					Nullable: true,
				}

				refName := getRefName(thisObject.MapValue)
				schemaRef := getSchemaRef(thisObject.MapValue)
				o.Components.Schemas[refName] = additionalPropertiesSchema

				schema.AdditionalProperties = &types.Schema{
					Ref: schemaRef,
				}
			}

			if thisObject.PointerValue != nil {
				pointerValueSchema, err := getSchema(thisObject.PointerValue)
				if err != nil {
					return nil, err
				}

				refName := getRefName(thisObject.PointerValue)
				schemaRef := getSchemaRef(thisObject.PointerValue)
				o.Components.Schemas[refName] = pointerValueSchema

				schema = &types.Schema{
					Ref:      schemaRef,
					Nullable: true,
				}
			}

			if schema != nil {
				return schema, nil
			}

			return nil, err
		}

		schema = typeMeta.GetOpenAPISchema()

		return schema, nil
	}

	for _, introspectedObject := range introspectedObjects {
		schema, err := getSchema(introspectedObject)
		if err != nil {
			return nil, fmt.Errorf(
				"no schema for root: %v; %v",
				introspectedObject.Name, err,
			)
		}

		_ = schema
	}

	o.Paths = make(map[string]*types.Path)

	for _, introspectedObject := range introspectedObjects {
		ref := getSchemaRef(introspectedObject)

		objectNamePlural := pluralize.Plural(caps.ToCamel(introspectedObject.Name))
		endpointNamePlural := pluralize.Plural(caps.ToKebab(introspectedObject.Name))

		objectNameSingular := pluralize.Singular(caps.ToCamel(introspectedObject.Name))

		get200Response := func(description string) *types.Response {
			return &types.Response{
				Description: fmt.Sprintf("Successful %v", description),
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
								"objects": {
									Type: types.TypeOfArray,
									Items: &types.Schema{
										Ref: ref,
									},
								},
								"count": {
									Type:   types.TypeOfInteger,
									Format: types.FormatOfInt64,
								},
								"total_count": {
									Type:   types.TypeOfInteger,
									Format: types.FormatOfInt64,
								},
								"limit": {
									Type:   types.TypeOfInteger,
									Format: types.FormatOfInt64,
								},
								"offset": {
									Type:   types.TypeOfInteger,
									Format: types.FormatOfInt64,
								},
							},
							Required: []string{"status", "success"},
						},
					},
				},
			}
		}

		get204Response := func(description string) *types.Response {
			return &types.Response{
				Description: fmt.Sprintf("Successful %v", description),
			}
		}

		getErrorResponse := func(description string) *types.Response {
			return &types.Response{
				Description: fmt.Sprintf("Failed %v", description),
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
							Required: []string{"status", "success"},
						},
					},
				},
			}
		}

		listParameters := make([]*types.Parameter, 0)

		listParameters = append(listParameters, &types.Parameter{
			Name:     "limit",
			In:       types.InQuery,
			Required: false,
			Schema: &types.Schema{
				Type:   types.TypeOfInteger,
				Format: types.FormatOfInt32,
			},
			Description: "SQL LIMIT operator",
		})

		listParameters = append(listParameters, &types.Parameter{
			Name:     "offset",
			In:       types.InQuery,
			Required: false,
			Schema: &types.Schema{
				Type:   types.TypeOfInteger,
				Format: types.FormatOfInt32,
			},
			Description: "SQL OFFSET operator",
		})

		listParameters = append(listParameters, &types.Parameter{
			Name:     "depth",
			In:       types.InQuery,
			Required: false,
			Schema: &types.Schema{
				Type: types.TypeOfInteger,
			},
			Description: "Max recursion depth for loading foreign objects; default = 1\n\n(0 = recurse until graph cycle detected, 1 = this object only, 2 = this object + neighbours, 3 = this object + neighbours + their neighbours... etc)",
		})

		for _, structFieldObject := range introspectedObject.StructFields {
			structFieldTypeName := getTypeName(structFieldObject.Object)

			_, structFieldIsDjangolangObject := introspectedObjectByName[structFieldTypeName]
			if structFieldIsDjangolangObject {
				continue
			}

			typeTemplate := structFieldObject.Name
			if structFieldObject.Type != nil {
				typeTemplate = structFieldObject.Type.String()
			}

			typeMeta, err := types.GetTypeMetaForTypeTemplate(typeTemplate)
			if err != nil {
				return nil, err
			}

			schema := typeMeta.GetOpenAPISchema()

			// skip objects / arrays / unsupported
			if schema.Type == types.TypeOfObject ||
				schema.Type == types.TypeOfArray ||
				schema.Type == "" {
				continue
			}

			for _, matcher := range matchers {
				_, ignored := ignoredValueByMatcher[matcher]
				if ignored {
					schema = &types.Schema{
						Type: types.TypeOfString,
					}
				}

				listParameters = append(listParameters, &types.Parameter{
					Name:        fmt.Sprintf("%v__%v", structFieldObject.Tag.Get("json"), matcher),
					In:          types.InQuery,
					Required:    false,
					Schema:      schema,
					Description: descriptionByMatcher[matcher],
				})
			}
		}

		listRequestBody := &types.RequestBody{
			Content: map[string]*types.MediaType{
				contentTypeApplicationJSON: {
					Schema: &types.Schema{
						Type:     types.TypeOfArray,
						Nullable: false,
						Items: &types.Schema{
							Ref: getSchemaRef(introspectedObject),
						},
					},
				},
			},
			Required: true,
		}

		o.Paths[fmt.Sprintf("%v/%v", endpointPrefix, endpointNamePlural)] = &types.Path{
			Get: &types.Operation{
				Tags:        []string{introspectedObject.Name},
				OperationID: fmt.Sprintf("%v%v", caps.ToCamel(http.MethodGet), objectNamePlural),
				Parameters:  listParameters,
				Responses: map[string]*types.Response{
					fmt.Sprintf("%v", http.StatusOK): get200Response(fmt.Sprintf("List Fetch for %v", objectNamePlural)),
					statusCodeDefault:                getErrorResponse(fmt.Sprintf("List Fetch for %v", objectNamePlural)),
				},
			},
			Post: &types.Operation{
				Tags:        []string{introspectedObject.Name},
				OperationID: fmt.Sprintf("%v%v", caps.ToCamel(http.MethodPost), objectNamePlural),
				Parameters: []*types.Parameter{
					{
						Name:     "depth",
						In:       types.InQuery,
						Required: false,
						Schema: &types.Schema{
							Type: types.TypeOfInteger,
						},
						Description: "Max recursion depth for loading foreign objects; default = 1\n\n(0 = recurse until graph cycle detected, 1 = this object only, 2 = this object + neighbours, 3 = this object + neighbours + their neighbours... etc)",
					},
				},
				RequestBody: listRequestBody,
				Responses: map[string]*types.Response{
					fmt.Sprintf("%v", http.StatusOK): get200Response(fmt.Sprintf("List Create for %v", objectNamePlural)),
					statusCodeDefault:                getErrorResponse(fmt.Sprintf("List Create for %v", objectNamePlural)),
				},
			},
		}

		itemParameters := []*types.Parameter{
			{
				Name:        "primaryKey",
				In:          types.InPath,
				Required:    true,
				Schema:      &types.Schema{},
				Description: fmt.Sprintf("Primary key for %v", introspectedObject.Name),
			},
			{
				Name:     "depth",
				In:       types.InQuery,
				Required: false,
				Schema: &types.Schema{
					Type: types.TypeOfInteger,
				},
				Description: "Max recursion depth for loading foreign objects; default = 1\n\n(0 = recurse until graph cycle detected, 1 = this object only, 2 = this object + neighbours, 3 = this object + neighbours + their neighbours... etc)",
			},
		}

		itemRequestBody := &types.RequestBody{
			Content: map[string]*types.MediaType{
				contentTypeApplicationJSON: {
					Schema: &types.Schema{
						Ref:      getSchemaRef(introspectedObject),
						Nullable: false,
					},
				},
			},
			Required: true,
		}

		o.Paths[fmt.Sprintf("%v/%v/{primaryKey}", endpointPrefix, endpointNamePlural)] = &types.Path{
			Get: &types.Operation{
				Tags:        []string{introspectedObject.Name},
				OperationID: fmt.Sprintf("%v%v", caps.ToCamel(http.MethodGet), objectNameSingular),
				Parameters:  itemParameters,
				Responses: map[string]*types.Response{
					fmt.Sprintf("%v", http.StatusOK): get200Response(fmt.Sprintf("Item Fetch for %v", objectNamePlural)),
					statusCodeDefault:                getErrorResponse(fmt.Sprintf("Item Fetch for %v", objectNamePlural)),
				},
			},
			Put: &types.Operation{
				Tags:        []string{introspectedObject.Name},
				OperationID: fmt.Sprintf("%v%v", caps.ToCamel(http.MethodPut), objectNameSingular),
				Parameters:  itemParameters,
				RequestBody: itemRequestBody,
				Responses: map[string]*types.Response{
					fmt.Sprintf("%v", http.StatusOK): get200Response(fmt.Sprintf("Item Replace for %v", objectNamePlural)),
					statusCodeDefault:                getErrorResponse(fmt.Sprintf("Item Replace for %v", objectNamePlural)),
				},
			},
			Patch: &types.Operation{
				Tags:        []string{introspectedObject.Name},
				OperationID: fmt.Sprintf("%v%v", caps.ToCamel(http.MethodPatch), objectNameSingular),
				Parameters:  itemParameters,
				RequestBody: itemRequestBody,
				Responses: map[string]*types.Response{
					fmt.Sprintf("%v", http.StatusOK): get200Response(fmt.Sprintf("Item Update for %v", objectNamePlural)),
					statusCodeDefault:                getErrorResponse(fmt.Sprintf("Item Update for %v", objectNamePlural)),
				},
			},
			Delete: &types.Operation{
				Tags:        []string{introspectedObject.Name},
				OperationID: fmt.Sprintf("%v%v", caps.ToCamel(http.MethodDelete), objectNameSingular),
				Parameters:  itemParameters,
				Responses: map[string]*types.Response{
					fmt.Sprintf("%v", http.StatusNoContent): get204Response(fmt.Sprintf("Item Delete for %v", objectNamePlural)),
					statusCodeDefault:                       getErrorResponse(fmt.Sprintf("Item Delete for %v", objectNamePlural)),
				},
			},
		}
	}

	introspectedEmptyRequest, err := introspect.Introspect(server.EmptyRequest{})
	if err != nil {
		return nil, err
	}

	introspectedEmptyResponse, err := introspect.Introspect(server.EmptyResponse{})
	if err != nil {
		return nil, err
	}

	for _, customHTTPHandlerSummary := range customHTTPHandlerSummaries {
		pathParamsIntrospectedObject, err := introspect.Introspect(customHTTPHandlerSummary.PathParams)
		if err != nil {
			return nil, err
		}

		queryParamsIntrospectedObject, err := introspect.Introspect(customHTTPHandlerSummary.QueryParams)
		if err != nil {
			return nil, err
		}

		requestIntrospectedObject, err := introspect.Introspect(customHTTPHandlerSummary.Request)
		if err != nil {
			return nil, err
		}

		requestIsEmpty := requestIntrospectedObject == introspectedEmptyRequest

		responseIntrospectedObject, err := introspect.Introspect(customHTTPHandlerSummary.Response)
		if err != nil {
			return nil, err
		}

		responseIsEmpty := responseIntrospectedObject == introspectedEmptyResponse

		fullPath := fmt.Sprintf("%s/%s%s", endpointPrefix, "custom", customHTTPHandlerSummary.Path)

		_, ok := o.Paths[fullPath]
		if !ok {
			o.Paths[fullPath] = &types.Path{}
		}

		endpointTag := "Custom"
		endpointName := caps.ToCamel(strings.ReplaceAll(strings.Trim(customHTTPHandlerSummary.Path, "/"), "/", "_"))

		parameters := []*types.Parameter{}

		for _, structFieldObject := range pathParamsIntrospectedObject.StructFields {
			structFieldObjectSchema, err := getSchema(structFieldObject.Object)
			if err != nil {
				return nil, err
			}

			tag := structFieldObject.Tag.Get("json")
			if tag == "" {
				tag = structFieldObject.Field
			}

			schemaRef := getSchemaRef(structFieldObject.Object)

			refName := getRefName(structFieldObject.Object)
			if !isPrimitive(structFieldObjectSchema) {
				o.Components.Schemas[refName] = structFieldObjectSchema
			}

			parameters = append(parameters, &types.Parameter{
				Name:     tag,
				In:       types.InPath,
				Required: true,
				Schema: &types.Schema{
					Ref: schemaRef,
				},
				Description: fmt.Sprintf("Path parameter %s", tag),
			})
		}

		for _, structFieldObject := range queryParamsIntrospectedObject.StructFields {
			structFieldObjectSchema, err := getSchema(structFieldObject.Object)
			if err != nil {
				return nil, err
			}

			tag := structFieldObject.Tag.Get("json")
			if tag == "" {
				tag = structFieldObject.Field
			}

			schemaRef := getSchemaRef(structFieldObject.Object)

			refName := getRefName(structFieldObject.Object)
			if !isPrimitive(structFieldObjectSchema) {
				o.Components.Schemas[refName] = structFieldObjectSchema
			}

			parameters = append(parameters, &types.Parameter{
				Name:     tag,
				In:       types.InQuery,
				Required: true,
				Schema: &types.Schema{
					Ref: schemaRef,
				},
				Description: fmt.Sprintf("Query parameter %s", tag),
			})
		}

		requestBodySchema, err := getSchema(requestIntrospectedObject)
		if err != nil {
			return nil, err
		}

		getRequest := func(method string) *types.RequestBody {
			if !(method == http.MethodPost || method == http.MethodPut || method == http.MethodPatch) {
				return nil
			}

			if requestIsEmpty {
				return nil
			}

			refName := getRefName(requestIntrospectedObject)
			schemaRef := getSchemaRef(requestIntrospectedObject)

			if !isPrimitive(requestBodySchema) {
				o.Components.Schemas[refName] = requestBodySchema
			}

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

		responseSchema, err := getSchema(responseIntrospectedObject)
		if err != nil {
			return nil, err
		}

		getSuccessResponse := func(method string) *types.Response {
			description := fmt.Sprintf("%v%v success", caps.ToCamel(method), endpointName)
			if responseIsEmpty {
				return &types.Response{
					Description: description,
				}
			}

			refName := getRefName(responseIntrospectedObject)
			schemaRef := getSchemaRef(responseIntrospectedObject)

			if !isPrimitive(responseSchema) {
				o.Components.Schemas[refName] = responseSchema
			}

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
			if responseIsEmpty {
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

		operation := &types.Operation{
			Tags:        []string{endpointTag},
			OperationID: fmt.Sprintf("%v%v", caps.ToCamel(customHTTPHandlerSummary.Method), endpointName),
			Parameters:  parameters,
			RequestBody: getRequest(customHTTPHandlerSummary.Method),
			Responses: map[string]*types.Response{
				fmt.Sprintf("%v", customHTTPHandlerSummary.Status): getSuccessResponse(customHTTPHandlerSummary.Method),
				statusCodeDefault: getErrorResponse(customHTTPHandlerSummary.Method),
			},
		}

		switch customHTTPHandlerSummary.Method {
		case http.MethodGet:
			o.Paths[fullPath].Get = operation
		case http.MethodPost:
			o.Paths[fullPath].Post = operation
		case http.MethodPut:
			o.Paths[fullPath].Put = operation
		case http.MethodPatch:
			o.Paths[fullPath].Patch = operation
		case http.MethodDelete:
			o.Paths[fullPath].Delete = operation
		default:
			return nil, fmt.Errorf("unsupported method %s for %s", customHTTPHandlerSummary.Method, customHTTPHandlerSummary.Path)
		}
	}

	return &o, nil
}
