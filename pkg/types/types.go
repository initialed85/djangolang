package types

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/google/uuid"
	"github.com/initialed85/djangolang/pkg/helpers"
	_pgtype "github.com/jackc/pgtype"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"

	"github.com/lib/pq"
	"github.com/lib/pq/hstore"

	"github.com/cridenour/go-postgis"
	jsoniter "github.com/json-iterator/go"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/ewkb"
)

var c = jsoniter.Config{
	EscapeHTML:              true,
	SortMapKeys:             false,
	MarshalFloatWith6Digits: true,
}.Froze()

type Type[T any] struct {
	DataType           string `json:"datatype"`
	ZeroType           any    `json:"zero_type"`
	QueryTypeTemplate  string `json:"query_type_template"`
	StreamTypeTemplate string `json:"stream_type_template"`
	TypeTemplate       string `json:"type_template"`
	ParseFunc          func(any) (T, error)
	ParseFuncTemplate  string `json:"parse_func_template"`
	IsZeroFunc         func(any) bool
	IsZeroFuncTemplate string `json:"is_zero_func_template"`
	FormatFunc         func(T) (any, error)
	FormatFuncTemplate string `json:"format_func_template"`
}

var theTypes = make([]*Type[any], 0)
var typeByDataType = make(map[string]*Type[any])
var typeByTypeTemplate = make(map[string]*Type[any])

func init() {
	geojson.CustomJSONMarshaler = c
	geojson.CustomJSONUnmarshaler = c

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
		"geometry(PointZ)",
	}

	for _, dataType := range dataTypes {
		var zeroType any
		queryTypeTemplate := ""
		streamTypeTemplate := ""
		typeTemplate := ""
		parseFunc := ParseNotImplemented
		parseFuncTemplate := "types.ParseNotImplemented(v)"
		isZeroFunc := IsZeroNotImplemented
		isZeroFuncTemplate := "types.IsZeroNotImplemented"
		formatFunc := FormatNotImplemented
		formatFuncTemplate := "types.FormatNotImplemented"

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

		case "timestamp without time zone":
			fallthrough
		case "timestamp with time zone":
			zeroType = time.Time{}
			queryTypeTemplate = "time.Time"
			parseFunc = ParseTime
			parseFuncTemplate = "types.ParseTime(v)"
			isZeroFunc = IsZeroTime
			isZeroFuncTemplate = "types.IsZeroTime"
			formatFunc = FormatTime
			formatFuncTemplate = "types.FormatTime"

		case "interval[]":
			zeroType = make([]time.Duration, 0)
			queryTypeTemplate = "[]time.Duration"

		case "interval":
			zeroType = helpers.Deref(new(time.Duration))
			queryTypeTemplate = "time.Duration"
			parseFunc = ParseDuration
			parseFuncTemplate = "types.ParseDuration(v)"
			isZeroFunc = IsZeroDuration
			isZeroFuncTemplate = "types.IsZeroDuration"
			formatFunc = FormatDuration
			formatFuncTemplate = "types.FormatDuration"

		case "json[]":
			fallthrough
		case "jsonb[]":
			zeroType = make([]any, 0)
			queryTypeTemplate = "[]any"

		case "json":
			fallthrough
		case "jsonb":
			zeroType = nil
			queryTypeTemplate = "any"
			parseFunc = ParseJSON
			parseFuncTemplate = "types.ParseJSON(v)"
			isZeroFunc = IsZeroJSON
			isZeroFuncTemplate = "types.IsZeroJSON"
			formatFunc = FormatJSON
			formatFuncTemplate = "types.FormatJSON"

		case "character varying[]":
			fallthrough
		case "text[]":
			zeroType = make(pq.StringArray, 0)
			queryTypeTemplate = "pq.StringArray"
			typeTemplate = "[]string"
			parseFunc = ParseStringArray
			parseFuncTemplate = "types.ParseStringArray(v)"
			isZeroFunc = IsZeroStringArray
			isZeroFuncTemplate = "types.IsZeroStringArray"
			formatFunc = FormatStringArray
			formatFuncTemplate = "types.FormatStringArray"
		case "character varying":
			fallthrough
		case "text":
			zeroType = helpers.Deref(new(string))
			queryTypeTemplate = "string"
			parseFunc = ParseString
			parseFuncTemplate = "types.ParseString(v)"
			isZeroFunc = IsZeroString
			isZeroFuncTemplate = "types.IsZeroString"
			formatFunc = FormatString
			formatFuncTemplate = "types.FormatString"

		case "smallint[]":
			fallthrough
		case "integer[]":
			fallthrough
		case "bigint[]":
			zeroType = make(pq.Int64Array, 0)
			queryTypeTemplate = "pq.Int64Array"

		case "smallint":
			fallthrough
		case "integer":
			fallthrough
		case "bigint":
			zeroType = helpers.Deref(new(int64))
			queryTypeTemplate = "int64"
			parseFunc = ParseInt
			parseFuncTemplate = "types.ParseInt(v)"
			isZeroFunc = IsZeroInt
			isZeroFuncTemplate = "types.IsZeroInt"
			formatFunc = FormatInt
			formatFuncTemplate = "types.FormatInt"

		case "real[]":
			fallthrough
		case "float[]":
			fallthrough
		case "numeric[]":
			fallthrough
		case "double precision[]":
			zeroType = make(pq.Float64Array, 0)
			queryTypeTemplate = "pq.Float64Array"

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
			isZeroFunc = IsZeroFloat
			isZeroFuncTemplate = "types.IsZeroFloat"
			formatFunc = FormatFloat
			formatFuncTemplate = "types.FormatFloat"

		case "boolean[]":
			zeroType = make(pq.BoolArray, 0)
			queryTypeTemplate = "pq.BoolArray"

		case "boolean":
			zeroType = helpers.Deref(new(bool))
			queryTypeTemplate = "bool"
			parseFunc = ParseBool
			parseFuncTemplate = "types.ParseBool(v)"
			isZeroFunc = IsZeroBool
			isZeroFuncTemplate = "types.IsZeroBool"
			formatFunc = FormatBool
			formatFuncTemplate = "types.FormatBool"

		case "tsvector[]":
			zeroType = make([]any, 0)
			queryTypeTemplate = "[]any"

		case "tsvector":
			zeroType = nil
			queryTypeTemplate = "any"

		case "uuid[]":
			zeroType = make([]uuid.UUID, 0)
			queryTypeTemplate = "[]uuid.UUID"
			streamTypeTemplate = "[][16]uint8"

		case "uuid":
			zeroType = uuid.UUID{}
			queryTypeTemplate = "uuid.UUID"
			streamTypeTemplate = "[16]uint8"
			parseFunc = ParseUUID
			parseFuncTemplate = "types.ParseUUID(v)"
			isZeroFunc = IsZeroUUID
			isZeroFuncTemplate = "types.IsZeroUUID"
			formatFunc = FormatUUID
			formatFuncTemplate = "types.FormatUUID"
			formatFunc = FormatUUID
			formatFuncTemplate = "types.FormatUUID"

		case "point[]":
			zeroType = make([]pgtype.Point, 0)
			queryTypeTemplate = "[]pgtype.Point"

		case "polygon[]":
			zeroType = make([]pgtype.Polygon, 0)
			queryTypeTemplate = "[]pgtype.Polygon"

		case "collection[]":
			zeroType = make([]geojson.Feature, 0)
			queryTypeTemplate = "[]pgtype.Feature"

		case "geometry[]":
			zeroType = make([]orb.Collection, 0)
			queryTypeTemplate = "[]pgtype.Geometry"

		case "hstore[]":
			zeroType = make([]hstore.Hstore, 0)
			queryTypeTemplate = "[]hstore.Hstore"
			typeTemplate = "[]map[string]*string"

		case "hstore":
			zeroType = hstore.Hstore{}
			queryTypeTemplate = "hstore.Hstore"
			typeTemplate = "map[string]*string"
			parseFunc = ParseHstore
			parseFuncTemplate = "types.ParseHstore(v)"
			isZeroFunc = IsZeroHstore
			isZeroFuncTemplate = "types.IsZeroHstore"
			formatFunc = FormatHstore
			formatFuncTemplate = "types.FormatHstore"

		case "point":
			zeroType = pgtype.Vec2{}
			queryTypeTemplate = "pgtype.Vec2"
			parseFunc = ParsePoint
			parseFuncTemplate = "types.ParsePoint(v)"
			isZeroFunc = IsZeroPoint
			isZeroFuncTemplate = "types.IsZeroPoint"
			formatFunc = FormatPoint
			formatFuncTemplate = "types.FormatPoint"

		case "polygon":
			zeroType = []pgtype.Vec2{}
			queryTypeTemplate = "[]pgtype.Vec2"
			parseFunc = ParsePolygon
			parseFuncTemplate = "types.ParsePolygon(v)"
			isZeroFunc = IsZeroPolygon
			isZeroFuncTemplate = "types.IsZeroPolygon"
			formatFunc = FormatPolygon
			formatFuncTemplate = "types.FormatPolygon"

		case "collection":
			zeroType = geojson.Feature{}
			queryTypeTemplate = "geojson.Feature"

		case "geometry", "geometry(PointZ)":
			zeroType = postgis.PointZ{}
			queryTypeTemplate = "postgis.PointZ"
			parseFunc = ParseGeometry
			parseFuncTemplate = "types.ParseGeometry(v)"
			isZeroFunc = IsZeroGeometry
			isZeroFuncTemplate = "types.IsZeroGeometry"
			formatFunc = FormatGeometry
			formatFuncTemplate = "types.FormatGeometry"

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
			ParseFunc:          parseFunc,
			ParseFuncTemplate:  parseFuncTemplate,
			IsZeroFunc:         isZeroFunc,
			IsZeroFuncTemplate: isZeroFuncTemplate,
			FormatFunc:         formatFunc,
			FormatFuncTemplate: formatFuncTemplate,
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

func ParseNotImplemented(v any) (any, error) {
	return nil, fmt.Errorf("parse not implemented for %#+v", v)
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

func ParseDuration(v any) (any, error) {
	switch v1 := v.(type) {
	case time.Duration:
		return v1, nil
	case _pgtype.Interval:
		return time.Microsecond * time.Duration(v1.Microseconds), nil
	case pgtype.Interval:
		return time.Microsecond * time.Duration(v1.Microseconds), nil
	}

	return time.Duration(0), fmt.Errorf("%#+v (%v) could not be identified ParseDuration", v, reflect.TypeOf(v).String())
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
			return pgtype.Vec2{}, fmt.Errorf("%#+v (%v) could not be parsed with pgtype.Point.Scan for ParsePoint; err: %v", v1, reflect.TypeOf(v).String(), err)
		}

		return v2.P, nil
	case pgtype.Point:
		return v1.P, nil
	}

	return pgtype.Vec2{}, fmt.Errorf("%#+v (%v) could not be identified for ParsePoint", v, reflect.TypeOf(v).String())
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

		return v3.P, nil
	case pgtype.Polygon:
		return v1.P, nil
	}

	return nil, fmt.Errorf("%#+v (%v) could not be identified for ParsePolygon", v, reflect.TypeOf(v).String())
}

func ParseGeometry(v any) (any, error) {
	switch v1 := v.(type) {
	case string:
		v2 := postgis.PointZ{}
		err := v2.Scan([]byte(v1))
		if err != nil {
			return nil, fmt.Errorf("%#+v (%v) could not be parsed with ewkbhex.Decode for ParseGeometry; err: %v", v1, reflect.TypeOf(v).String(), err)
		}

		return v2, nil
	}

	return nil, fmt.Errorf("%#+v (%v) could not be identified for ParseGeometry", v, reflect.TypeOf(v).String())
}

func IsZeroNotImplemented(v any) bool {
	return true
}

func IsZeroUUID(v any) bool {
	if v == nil {
		return true
	}

	v1, ok := v.(uuid.UUID)
	if !ok {
		return false
	}

	return v1 == uuid.Nil
}

func IsZeroTime(v any) bool {
	if v == nil {
		return true
	}

	v1, ok := v.(time.Time)
	if !ok {
		return false
	}

	return v1.IsZero()
}

func IsZeroDuration(v any) bool {
	if v == nil {
		return true
	}

	v1, ok := v.(time.Duration)
	if !ok {
		return false
	}

	return v1 == time.Duration(0)
}

func IsZeroString(v any) bool {
	if v == nil {
		return true
	}

	v1, ok := v.(string)
	if !ok {
		return false
	}

	return v1 == ""
}

func IsZeroStringArray(v any) bool {
	if v == nil {
		return true
	}

	v1, ok := v.([]string)
	if !ok {
		return false
	}

	return v1 == nil
}

func IsZeroHstore(v any) bool {
	if v == nil {
		return true
	}

	v1, ok := v.(map[string]*string)
	if !ok {
		return false
	}

	return v1 == nil
}

func IsZeroJSON(v any) bool {
	if v == nil {
		return true
	}

	return v == nil
}

func IsZeroInt(v any) bool {
	if v == nil {
		return true
	}

	v1, ok := v.(int64)
	if !ok {
		return false
	}

	return v1 == 0
}

func IsZeroFloat(v any) bool {
	if v == nil {
		return true
	}

	v1, ok := v.(float64)
	if !ok {
		return false
	}

	return v1 == 0.0
}

func IsZeroBool(v any) bool {
	if v == nil {
		return true
	}

	v1, ok := v.(bool)
	if !ok {
		return false
	}

	return !v1
}

func IsZeroPoint(v any) bool {
	if v == nil {
		return true
	}

	v1, ok := v.(pgtype.Point)
	if !ok {
		return false
	}

	return v1 == pgtype.Point{}
}

func IsZeroPolygon(v any) bool {
	if v == nil {
		return true
	}

	v1, ok := v.(pgtype.Polygon)
	if !ok {
		return false
	}

	return v1.P == nil
}

func IsZeroGeometry(v any) bool {
	if v == nil {
		return true
	}

	v1, ok := v.(geom.T)
	if !ok {
		return false
	}

	return v1 == nil || v1.Empty()
}

func FormatNotImplemented(v any) (any, error) {
	return nil, fmt.Errorf("format not implemented for %#+v", v)
}

func FormatUUID(v any) (any, error) {
	v1, ok := v.(*uuid.UUID)
	if ok {
		if v1 == nil {
			return nil, nil
		}

		return *v1, nil
	}

	v2, ok := v.(uuid.UUID)
	if !ok {
		return nil, fmt.Errorf("%#+v (%v) could not be cast to uuid.UUID for FormatUUID", v, reflect.TypeOf(v).String())
	}

	return v2, nil
}

func FormatTime(v any) (any, error) {
	v1, ok := v.(*time.Time)
	if ok {
		if v1 == nil {
			return nil, nil
		}

		return *v1, nil
	}

	v2, ok := v.(time.Time)
	if !ok {
		return nil, fmt.Errorf("%#+v (%v) could not be cast to time.Time for FormatTime", v, reflect.TypeOf(v).String())
	}

	return v2, nil
}

func FormatDuration(v any) (any, error) {
	v1, ok := v.(*time.Duration)
	if ok {
		if v1 == nil {
			return nil, nil
		}

		return *v1, nil
	}

	v2, ok := v.(time.Duration)
	if !ok {
		return nil, fmt.Errorf("%#+v (%v) could not be cast to time.Duration for FormatDuration", v, reflect.TypeOf(v).String())
	}

	return v2, nil
}

func FormatString(v any) (any, error) {
	v1, ok := v.(*string)
	if ok {
		if v1 == nil {
			return nil, nil
		}

		return *v1, nil
	}

	v2, ok := v.(string)
	if !ok {
		return nil, fmt.Errorf("%#+v (%v) could not be cast to string for FormatString", v, reflect.TypeOf(v).String())
	}

	return v2, nil
}

func FormatStringArray(v any) (any, error) {
	format := func(x []string) pq.StringArray {
		return pq.StringArray(x)
	}

	v1, ok := v.(*[]string)
	if ok {
		if v1 == nil {
			return nil, nil
		}

		return format(*v1), nil
	}

	v2, ok := v.([]string)
	if !ok {
		return nil, fmt.Errorf("%#+v (%v) could not be cast to []string for FormatStringArray", v, reflect.TypeOf(v).String())
	}

	return format(v2), nil
}

func FormatHstore(v any) (any, error) {
	format := func(x map[string]*string) hstore.Hstore {
		y := hstore.Hstore{
			Map: make(map[string]sql.NullString),
		}

		for k, v := range x {
			z := sql.NullString{}

			if v != nil {
				z.String = *v
				z.Valid = true
			}

			y.Map[k] = z
		}

		return y
	}

	v1, ok := v.(*map[string]*string)
	if ok {
		if v1 == nil {
			return nil, nil
		}

		return format(*v1), nil
	}

	v2, ok := v.(map[string]*string)
	if !ok {
		return nil, fmt.Errorf("%#+v (%v) could not be cast to map[string]*string for FormatHstore", v, reflect.TypeOf(v).String())
	}

	return format(v2), nil
}

func FormatJSON(v any) (any, error) {
	v1, ok := v.(*any)
	if ok {
		if v1 == nil {
			return nil, nil
		}

		return json.Marshal(*v1)
	}

	v2, ok := v.(any)
	if !ok {
		return nil, fmt.Errorf("%#+v (%v) could not be cast to any for FormatJSON", v, reflect.TypeOf(v).String())
	}

	return json.Marshal(v2)
}

func FormatInt(v any) (any, error) {
	v1, ok := v.(*int64)
	if ok {
		if v1 == nil {
			return nil, nil
		}

		return *v1, nil
	}

	v2, ok := v.(int64)
	if !ok {
		return nil, fmt.Errorf("%#+v (%v) could not be cast to int64 for FormatInt", v, reflect.TypeOf(v).String())
	}

	return v2, nil
}

func FormatFloat(v any) (any, error) {
	v1, ok := v.(*float64)
	if ok {
		if v1 == nil {
			return nil, nil
		}

		return *v1, nil
	}

	v2, ok := v.(float64)
	if !ok {
		return nil, fmt.Errorf("%#+v (%v) could not be cast to float64 for FormatFloat", v, reflect.TypeOf(v).String())
	}

	return v2, nil
}

func FormatBool(v any) (any, error) {
	v1, ok := v.(*bool)
	if ok {
		if v1 == nil {
			return nil, nil
		}

		return *v1, nil
	}

	v2, ok := v.(bool)
	if !ok {
		return nil, fmt.Errorf("%#+v (%v) could not be cast to bool for FormatBool", v, reflect.TypeOf(v).String())
	}

	return v2, nil
}

func FormatPoint(v any) (any, error) {
	v1, ok := v.(*pgtype.Vec2)
	if ok {
		if v1 == nil {
			return nil, nil
		}

		return *v1, nil
	}

	v2, ok := v.(pgtype.Vec2)
	if !ok {
		return nil, fmt.Errorf("%#+v (%v) could not be cast to pgtype.Vec2 for FormatPoint", v, reflect.TypeOf(v).String())
	}

	return pgtype.Point{P: v2, Valid: true}, nil
}

func FormatPolygon(v any) (any, error) {
	v1, ok := v.(*[]pgtype.Vec2)
	if ok {
		if v1 == nil {
			return nil, nil
		}

		return *v1, nil
	}

	v2, ok := v.([]pgtype.Vec2)
	if !ok {
		return nil, fmt.Errorf("%#+v (%v) could not be cast to []pgtype.Vec2 for FormatPolygon", v, reflect.TypeOf(v).String())
	}

	return pgtype.Polygon{P: v2, Valid: true}, nil
}

func FormatGeometry(v any) (any, error) {
	v1, ok := v.(*ewkb.Point)
	if ok {
		if v1 == nil {
			return nil, nil
		}

		return *v1, nil
	}

	v2, ok := v.(ewkb.Point)
	if !ok {
		return nil, fmt.Errorf("%#+v (%v) could not be cast to ewkb.Point for FormatGeometry", v, reflect.TypeOf(v).String())
	}

	return v2, nil
}
