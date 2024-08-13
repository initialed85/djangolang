package openapi

import (
	"fmt"
	"net/http"
	"strings"

	_pluralize "github.com/gertd/go-pluralize"

	"github.com/chanced/caps"
	"github.com/initialed85/djangolang/pkg/helpers"
	"github.com/initialed85/djangolang/pkg/types"
	"github.com/initialed85/structmeta/pkg/introspect"
)

// this file contains contains a fairly limited implementation of the OpenAPI v3 structure; limited as in just enough to convey the information
// needed to generate clients for a Djangolang API

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
		"nin",
		"notin",
		"isnull",
		"nisnull",
		"isnotnull",
		"l",
		"like",
		"nl",
		"nlike",
		"notlike",
		"il",
		"ilike",
		"nil",
		"nilike",
		"notilike",
		"desc",
		"asc",
	}

	descriptionByMatcher = map[string]string{
		"eq":        "SQL = operator",
		"ne":        "SQL != operator",
		"gt":        "SQL > operator, may not work with all column types",
		"gte":       "SQL >= operator, may not work with all column types",
		"lt":        "SQL < operator, may not work with all column types",
		"lte":       "SQL <= operator, may not work with all column types",
		"in":        "SQL IN operator, permits comma-separated values",
		"nin":       "SQL NOT IN operator, permits comma-separated values",
		"notin":     "SQL NOT IN operator, permits comma-separated values",
		"isnull":    "SQL IS NULL operator, value is ignored (presence of key is sufficient)",
		"nisnull":   "SQL IS NOT NULL operator, value is ignored (presence of key is sufficient)",
		"isnotnull": "SQL IS NOT NULL operator, value is ignored (presence of key is sufficient)",
		"l":         "SQL LIKE operator, value is implicitly prefixed and suffixed with %",
		"like":      "SQL LIKE operator, value is implicitly prefixed and suffixed with %",
		"nl":        "SQL NOT LIKE operator, value is implicitly prefixed and suffixed with %",
		"nlike":     "SQL NOT LIKE operator, value is implicitly prefixed and suffixed with %",
		"notlike":   "SQL NOT LIKE operator, value is implicitly prefixed and suffixed with %",
		"il":        "SQL ILIKE operator, value is implicitly prefixed and suffixed with %",
		"ilike":     "SQL ILIKE operator, value is implicitly prefixed and suffixed with %",
		"nil":       "SQL NOT ILIKE operator, value is implicitly prefixed and suffixed with %",
		"nilike":    "SQL NOT ILIKE operator, value is implicitly prefixed and suffixed with %",
		"notilike":  "SQL NOT ILIKE operator, value is implicitly prefixed and suffixed with %",
		"desc":      "SQL ORDER BY _ DESC operator, value is ignored (presence of key is sufficient)",
		"asc":       "SQL ORDER BY _ ASC operator, value is ignored (presence of key is sufficient)",
	}

	ignoredValueByMatcher = map[string]struct{}{
		"isnull":    {},
		"nisnull":   {},
		"isnotnull": {},
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

func NewFromIntrospectedSchema(inputObjects []any) (*types.OpenAPI, error) {
	apiRoot := helpers.GetEnvironmentVariableOrDefault("DJANGOLANG_API_ROOT", "/")
	apiRoot = helpers.GetEnvironmentVariableOrDefault("DJANGOLANG_API_ROOT_FOR_OPENAPI", apiRoot)

	endpointPrefix := strings.TrimRight(apiRoot, "/")

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

	getTypeName := func(object *introspect.Object) string {
		typeName := object.Name
		if object.Type != nil {
			typeName = object.Type.String()
		}

		return typeName
	}

	for i, inputObject := range inputObjects {
		introspectedObject, err := introspect.Introspect(inputObject)
		if err != nil {
			return nil, fmt.Errorf("failed to introspect %#+v: %v", inputObject, err)
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

	getRefName := func(object *introspect.Object) string {
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

	getSchemaRef := func(object *introspect.Object) string {
		refName := getRefName(object)

		return fmt.Sprintf("#/components/schemas/%v", refName)
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
				if schema.Type == types.TypeOfObject {
					o.Components.Schemas[getRefName(thisObject)] = schema

					if isDjangolangObject {
						o.Components.Schemas["Nullable"+getRefName(thisObject)] = &types.Schema{
							Ref:      getSchemaRef(thisObject),
							Nullable: true,
						}

						o.Components.Schemas["NullableArrayOf"+getRefName(thisObject)] = &types.Schema{
							Type:     types.TypeOfArray,
							Nullable: true,
							Items: &types.Schema{
								Ref: getSchemaRef(thisObject),
							},
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

				schema.Properties[structFieldObject.Tag.Get("json")] = structFieldSchema
			}

			return schema, nil
		}

		typeTemplate := thisObject.Name
		if thisObject.Type != nil {
			typeTemplate = thisObject.Type.String()
		}

		typeMeta, err := types.GetTypeMetaForTypeTemplate(typeTemplate)
		if err != nil {
			return nil, err
		}

		schema = typeMeta.GetOpenAPISchema()

		return schema, nil
	}

	for _, introspectedObject := range introspectedObjects {
		schema, err := getSchema(introspectedObject)
		if err != nil {
			return nil, fmt.Errorf(
				"no schema for root: %v: %v",
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
									Type: types.TypeOfString,
								},
								"objects": {
									Type: types.TypeOfArray,
									Items: &types.Schema{
										Ref: ref,
									},
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
									Type: types.TypeOfString,
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
				Parameters:  make([]*types.Parameter, 0),
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

	return &o, nil
}
