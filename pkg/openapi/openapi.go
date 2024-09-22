package openapi

import (
	"fmt"
	"log"
	"regexp"
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

var genericPattern = regexp.MustCompile(`.*\.(.*)\[.*\.(.*)\]`)
var packageStructPattern = regexp.MustCompile(`.*\.(.*)`)

func getTypeName(object *introspect.Object) string {
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

var refNameByIntrospectedObjectTypeName = make(map[string]string)

func getRefName(object *introspect.Object) string {
	typeName := getTypeName(object)

	refName, ok := refNameByIntrospectedObjectTypeName[typeName]
	if ok {
		return refName
	}

	if object.PointerValue != nil {
		refName += "nullable_" + getRefName(object.PointerValue)
	} else if object.SliceValue != nil {
		refName += "array_of_" + getRefName(object.SliceValue)
	} else if object.MapKey != nil && object.MapValue != nil {
		refName += "map_of_" + getRefName(object.MapKey) + "_" + getRefName(object.MapValue)
	} else if object.StructFields != nil {
		refName = typeName
	} else {
		refName = "primitive_" + typeName
	}

	camelCaseRefName := caps.ToCamel(refName)

	refNameByIntrospectedObjectTypeName[typeName] = camelCaseRefName
	log.Printf("%s | %s", typeName, camelCaseRefName)

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

	if strings.Contains(schema.Ref, "schemas/Primitive") || strings.Contains(schema.Ref, "schemas/NullablePrimitive") {
		return true
	}

	return false
}

func NewFromIntrospectedSchema(httpHandlerSummaries []server.HTTPHandlerSummary) (*types.OpenAPI, error) {
	apiRootForOpenAPI := config.APIRootForOpenAPI()

	endpointPrefix := strings.TrimRight(apiRootForOpenAPI, "/")

	_ = endpointPrefix

	o := types.OpenAPI{
		OpenAPI: "3.0.0",
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
				o.Components.Schemas[getTypeName(thisObject)] = schema
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
				_ = pointerValueSchema

				refName := getSchemaRef(thisObject.PointerValue)

				schema = &types.Schema{
					Ref:      refName,
					Nullable: true,
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

		// if schema == nil || (schema.Type == "" && schema.Ref == "") {
		// 	return nil, fmt.Errorf("failed to work out schema for %#+v", thisObject)
		// }

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

		// if !isPrimitive(schema) {
		// 	typeName := getRefName(introspectedObject)
		// 	o.Components.Schemas[typeName] = schema
		// }

		typeName := getRefName(introspectedObject)
		o.Components.Schemas[typeName] = schema
	}

	return &o, nil
}
