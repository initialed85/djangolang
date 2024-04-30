package types

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/google/uuid"
	"github.com/initialed85/djangolang/pkg/helpers"
	_pgtype "github.com/jackc/pgtype"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/paulmach/orb/geojson"

	"github.com/lib/pq"
	"github.com/lib/pq/hstore"
)

type Type[T any] struct {
	DataType           string `json:"datatype"`
	ZeroType           any    `json:"zero_type"`
	QueryTypeTemplate  string `json:"query_type_template"`
	StreamTypeTemplate string `json:"stream_type_template"`
	TypeTemplate       string `json:"type_template"`
	ParseFunc          func(any) (T, error)
	ParseFuncTemplate  string `json:"parse_func_template"`
}

var theTypes = make([]*Type[any], 0)
var typeByDataType = make(map[string]*Type[any])
var typeByTypeTemplate = make(map[string]*Type[any])

func init() {
	// TODO: not an exhaustive source type list
	dataTypes := []string{
		"timestamp without time zone[]",
		"timestamp with time zone[]",
		"timestamp without time zone",
		"timestamp with time zone",
		"interval[]",
		"interval",
		"json[]",
		"jsonb[]",
		"json",
		"jsonb",
		"character varying[]",
		"text[]",
		"character varying",
		"text",
		"smallint[]",
		"integer[]",
		"bigint[]",
		"smallint",
		"integer",
		"bigint",
		"real[]",
		"float[]",
		"numeric[]",
		"double precision[]",
		"float",
		"real",
		"numeric",
		"double precision",
		"boolean[]",
		"boolean",
		"tsvector[]",
		"tsvector",
		"uuid[]",
		"uuid",
		"point[]",
		"polygon[]",
		"collection[]",
		"geometry[]",
		"hstore[]",
		"hstore",
		"point",
		"polygon",
		"collection",
		"geometry",
	}

	for _, dataType := range dataTypes {
		var zeroType any
		queryTypeTemplate := ""
		streamTypeTemplate := ""
		typeTemplate := ""
		var parseFunc func(any) (any, error)
		parseFuncTemplate := ""

		// TODO: not an exhaustive suite of implementations for the source type list
		switch dataType {

		//
		// slices
		//

		case "timestamp without time zone[]":
			fallthrough
		case "timestamp with time zone[]":
			zeroType = make([]time.Time, 0)
			queryTypeTemplate = "[]time.Time"
			// TODO: parseFunc
		case "timestamp without time zone":
			fallthrough
		case "timestamp with time zone":
			zeroType = time.Time{}
			queryTypeTemplate = "time.Time"
			parseFunc = ParseTime
			parseFuncTemplate = "types.ParseTime(v)"

		case "interval[]":
			zeroType = make([]time.Duration, 0)
			queryTypeTemplate = "[]time.Duration"
			// TODO: parseFunc
		case "interval":
			zeroType = helpers.Deref(new(time.Duration))
			queryTypeTemplate = "time.Duration"
			// TODO: parseFunc

		case "json[]":
			fallthrough
		case "jsonb[]":
			zeroType = make([]any, 0)
			queryTypeTemplate = "[]any"
			// TODO: parseFunc
		case "json":
			fallthrough
		case "jsonb":
			zeroType = nil
			queryTypeTemplate = "any"
			parseFunc = ParseJSON
			parseFuncTemplate = "types.ParseJSON(v)"

		case "character varying[]":
			fallthrough
		case "text[]":
			zeroType = make(pq.StringArray, 0)
			queryTypeTemplate = "pq.StringArray"
			typeTemplate = "[]string"
			parseFunc = ParseStringArray
			parseFuncTemplate = "types.ParseStringArray(v)"
		case "character varying":
			fallthrough
		case "text":
			zeroType = helpers.Deref(new(string))
			queryTypeTemplate = "string"
			parseFunc = ParseString
			parseFuncTemplate = "types.ParseString(v)"

		case "smallint[]":
			fallthrough
		case "integer[]":
			fallthrough
		case "bigint[]":
			zeroType = make(pq.Int64Array, 0)
			queryTypeTemplate = "pq.Int64Array"
			// TODO: parseFunc
		case "smallint":
			fallthrough
		case "integer":
			fallthrough
		case "bigint":
			zeroType = helpers.Deref(new(int64))
			queryTypeTemplate = "int64"
			parseFunc = ParseInt
			parseFuncTemplate = "types.ParseInt(v)"

		case "real[]":
			fallthrough
		case "float[]":
			fallthrough
		case "numeric[]":
			fallthrough
		case "double precision[]":
			zeroType = make(pq.Float64Array, 0)
			queryTypeTemplate = "pq.Float64Array"
			// TODO: parseFunc
		case "float":
			fallthrough
		case "real":
			fallthrough
		case "numeric":
			fallthrough
		case "double precision":
			zeroType = helpers.Deref(new(float64))
			queryTypeTemplate = "float64"
			parseFunc = ParseFloat
			parseFuncTemplate = "types.ParseFloat(v)"

		case "boolean[]":
			zeroType = make(pq.BoolArray, 0)
			queryTypeTemplate = "pq.BoolArray"
			// TODO: parseFunc
		case "boolean":
			zeroType = helpers.Deref(new(bool))
			queryTypeTemplate = "bool"
			parseFunc = ParseBool
			parseFuncTemplate = "types.ParseBool(v)"

		case "tsvector[]":
			zeroType = make([]any, 0)
			queryTypeTemplate = "[]any"
			// TODO: parseFunc
		case "tsvector":
			zeroType = nil
			queryTypeTemplate = "any"
			// TODO: parseFunc

		case "uuid[]":
			zeroType = make([]uuid.UUID, 0)
			queryTypeTemplate = "[]uuid.UUID"
			streamTypeTemplate = "[][16]uint8"
			// TODO: parseFunc
		case "uuid":
			zeroType = uuid.UUID{}
			queryTypeTemplate = "uuid.UUID"
			streamTypeTemplate = "[16]uint8"
			parseFunc = ParseUUID
			parseFuncTemplate = "types.ParseUUID(v)"

		case "point[]":
			zeroType = make([]pgtype.Point, 0)
			queryTypeTemplate = "[]pgtype.Point"
			// TODO: parseFunc
		case "polygon[]":
			zeroType = make([]pgtype.Polygon, 0)
			queryTypeTemplate = "[]pgtype.Polygon"
			// TODO: parseFunc
		case "collection[]":
			zeroType = make([]geojson.Feature, 0)
			queryTypeTemplate = "[]pgtype.Feature"
			// TODO: parseFunc
		case "geometry[]":
			zeroType = make([]geojson.Geometry, 0)
			queryTypeTemplate = "[]pgtype.Geometry"
			// TODO: parseFunc

		case "hstore[]":
			zeroType = make([]hstore.Hstore, 0)
			queryTypeTemplate = "[]hstore.Hstore"
			typeTemplate = "[]map[string]*string"
			// TODO: parseFunc
		case "hstore":
			zeroType = hstore.Hstore{}
			queryTypeTemplate = "hstore.Hstore"
			typeTemplate = "map[string]*string"
			parseFunc = ParseHstore
			parseFuncTemplate = "types.ParseHstore(v)"

		case "point":
			zeroType = pgtype.Point{}
			queryTypeTemplate = "pgtype.Point"
			parseFunc = ParsePoint
			parseFuncTemplate = "types.ParsePoint(v)"
		case "polygon":
			zeroType = pgtype.Polygon{}
			queryTypeTemplate = "pgtype.Polygon"
			parseFunc = ParsePolygon
			parseFuncTemplate = "types.ParsePolygon(v)"
		case "collection":
			zeroType = geojson.Feature{}
			queryTypeTemplate = "geojson.Feature"
			// TODO: parseFunc
		case "geometry":
			zeroType = geojson.Geometry{}
			queryTypeTemplate = "geojson.Geometry"
			// TODO: parseFunc

		default:
			panic(
				fmt.Sprintf("failed to work out Go type details for Postgres type %#+v",
					dataType,
				),
			)
		}

		if streamTypeTemplate == "" {
			streamTypeTemplate = queryTypeTemplate
		}

		if typeTemplate == "" {
			typeTemplate = queryTypeTemplate
		}

		theType := Type[any]{
			DataType:           dataType,
			ZeroType:           zeroType,
			QueryTypeTemplate:  queryTypeTemplate,
			StreamTypeTemplate: streamTypeTemplate,
			TypeTemplate:       typeTemplate,
			ParseFuncTemplate:  parseFuncTemplate,
		}

		if parseFunc != nil {
			theType.ParseFunc = parseFunc
		}

		theTypes = append(theTypes, &theType)
	}

	for _, theType := range theTypes {
		typeByDataType[theType.DataType] = theType
		typeByTypeTemplate[theType.TypeTemplate] = theType
	}
}

func GetTypeForDataType(dataType string) (*Type[any], error) {
	theType := typeByDataType[dataType]
	if theType == nil {
		return nil, fmt.Errorf("unknown dataType %#+v", dataType)
	}

	return theType, nil
}

func GetTypeForTypeTemplate(typeTemplate string) (*Type[any], error) {
	theType := typeByTypeTemplate[typeTemplate]
	if theType == nil {
		return nil, fmt.Errorf("unknown typeTemplate %#+v", typeTemplate)
	}

	return theType, nil
}

func ParseUUID(v any) (any, error) {
	switch v1 := v.(type) {

	case []byte:
		v2, err := uuid.Parse(string(v1))
		if err != nil {
			return uuid.UUID{}, fmt.Errorf("%#+v (%v) could not be parsed with uuid.Parse for ParseUUID; err: %v", v, reflect.TypeOf(v).String(), err)
		}

		return v2, nil

	case [16]byte:
		v2, err := uuid.FromBytes(v1[:])
		if err != nil {
			return uuid.UUID{}, fmt.Errorf("%#+v (%v) could not be parsed with uuid.ParseBytes for ParseUUID; err: %v", v, reflect.TypeOf(v).String(), err)
		}

		return v2, nil
	}

	return uuid.UUID{}, fmt.Errorf("%#+v (%v) could not be identified for ParseTime", v, reflect.TypeOf(v).String())
}

func ParseTime(v any) (any, error) {
	switch v1 := v.(type) {
	case time.Time:
		return v1, nil
	}

	return time.Time{}, fmt.Errorf("%#+v (%v) could not be identified ParseTime", v, reflect.TypeOf(v).String())
}

func ParseString(v any) (any, error) {
	switch v1 := v.(type) {
	case string:
		return v1, nil
	}

	return "", fmt.Errorf("%#+v (%v) could not be identified for ParseString", v, reflect.TypeOf(v).String())
}

func ParseStringArray(v any) (any, error) {
	switch v1 := v.(type) {
	case []byte:
		v2 := pq.StringArray{}
		err := v2.Scan(v1)
		if err != nil {
			return nil, fmt.Errorf("%#+v (%v) could not be parsed with pq.StringArray.Scan for ParseStringArray; err: %v", v1, reflect.TypeOf(v).String(), err)
		}

		v3 := []string(v2)
		return v3, nil
	case []any:
		temp2 := make([]string, 0)
		for _, v := range v1 {
			s, ok := v.(string)
			if !ok {
				return nil, fmt.Errorf("%#+v (%v) could not be cast to []string for ParseStringArray", v, reflect.TypeOf(v).String())
			}

			temp2 = append(temp2, s)
		}

		return temp2, nil
	}

	return nil, fmt.Errorf("%#+v (%v) could not be identified for ParseStringArray", v, reflect.TypeOf(v).String())
}

func ParseHstore(v any) (any, error) {
	v1, ok := v.([]byte)
	if !ok {
		temp, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("%#+v (%v) could not be cast to string for ParseHstore", v, reflect.TypeOf(v).String())
		}

		v1 = []byte(temp)
	}

	v2 := hstore.Hstore{}
	err := v2.Scan(v1)
	if err != nil {
		return nil, fmt.Errorf("%#+v (%v) could not be parsed with pq.StringArray.Scan for ParseHstore; err: %v", v1, reflect.TypeOf(v).String(), err)
	}

	v3 := make(map[string]*string)

	for k, v := range v2.Map {
		if v.Valid {
			v3[k] = helpers.Ptr(v.String)
		} else {
			v3[k] = helpers.Nil("")
		}
	}

	return v3, nil
}

func ParseJSON(v any) (any, error) {
	switch v := v.(type) {
	case []byte:
		var v1 any
		err := json.Unmarshal(v, &v1)
		if err != nil {
			return nil, fmt.Errorf("%#+v (%v) could not be parsed with json.Unmarshal for ParseJSON; err: %v", v1, reflect.TypeOf(v).String(), err)
		}

		return v1, nil
	case (map[string]any), []map[string]any, []any, any:
		return v, nil
	default:
	}

	return nil, fmt.Errorf("%#+v could not be identified ParseJSON", v)
}

func ParseInt(v any) (any, error) {
	switch v1 := v.(type) {
	case int64:
		return v1, nil
	}

	return 0, fmt.Errorf("%#+v (%v) could not be identified for ParseInt", v, reflect.TypeOf(v).String())
}

func ParseFloat(v any) (any, error) {
	switch v1 := v.(type) {
	case float64:
		return v1, nil
	}

	return 0, fmt.Errorf("%#+v (%v) could not be identified for ParseFloat", v, reflect.TypeOf(v).String())
}

func ParseBool(v any) (any, error) {
	switch v1 := v.(type) {
	case bool:
		return v1, nil
	}

	return 0, fmt.Errorf("%#+v (%v) could not be identified for ParseBool", v, reflect.TypeOf(v).String())
}

func ParsePoint(v any) (any, error) {
	switch v1 := v.(type) {
	case []byte:
		v2 := pgtype.Point{}
		err := v2.UnmarshalJSON(v1)
		if err != nil {
			return pgtype.Point{}, fmt.Errorf("%#+v (%v) could not be parsed with pgtype.Point.Scan for ParsePoint; err: %v", v1, reflect.TypeOf(v).String(), err)
		}

		return v2, nil
	case pgtype.Point:
		return v1, nil
	}

	return pgtype.Point{}, fmt.Errorf("%#+v (%v) could not be identified for ParsePoint", v, reflect.TypeOf(v).String())
}

func ParsePolygon(v any) (any, error) {
	switch v1 := v.(type) {
	case []byte:
		v2 := _pgtype.Polygon{}
		err := v2.Scan(v1)
		if err != nil {
			return pgtype.Polygon{}, fmt.Errorf("%#+v (%v) could not be parsed with _pgtype.Polygon.Scan for ParsePolygon; err: %v", v1, reflect.TypeOf(v).String(), err)
		}

		v3 := pgtype.Polygon{
			P:     make([]pgtype.Vec2, 0),
			Valid: v2.Status == _pgtype.Present,
		}

		for _, p := range v2.P {
			v3.P = append(v3.P, pgtype.Vec2{X: p.X, Y: p.Y})
		}

		return v3, nil
	case pgtype.Polygon:
		return v1, nil
	}

	return pgtype.Polygon{}, fmt.Errorf("%#+v (%v) could not be identified for ParsePolygon", v, reflect.TypeOf(v).String())
}
