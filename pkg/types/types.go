package types

import (
	"fmt"
	"net/netip"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/initialed85/djangolang/pkg/helpers"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/paulmach/orb/geojson"
	"golang.org/x/exp/maps"

	"github.com/cridenour/go-postgis"
	jsoniter "github.com/json-iterator/go"
)

var c = jsoniter.Config{
	EscapeHTML:              true,
	SortMapKeys:             false,
	MarshalFloatWith6Digits: true,
}.Froze()

type TypeMeta[T any] struct {
	DataType           string
	ZeroType           any
	QueryTypeTemplate  string
	StreamTypeTemplate string
	TypeTemplate       string
	GetOpenAPISchema   func() *Schema
	ParseFunc          func(any) (T, error)
	ParseFuncTemplate  string
	IsZeroFunc         func(any) bool
	IsZeroFuncTemplate string
	FormatFunc         func(T) (any, error)
	FormatFuncTemplate string
}

var typeMetas = make([]*TypeMeta[any], 0)
var typeMetaByDataType = make(map[string]*TypeMeta[any])
var typeMetaByTypeTemplate = make(map[string]*TypeMeta[any])

func init() {
	geojson.CustomJSONMarshaler = c
	geojson.CustomJSONUnmarshaler = c

	// TODO: not an exhaustive source type list
	dataTypes := []string{
		"bigint",
		"bigint[]",
		"boolean",
		"boolean[]",
		"bytea",
		"character",
		"character[]",
		"character varying",
		"character varying[]",
		"double precision",
		"double precision[]",
		"float",
		"float[]",
		"geometry",
		"geometry(PointZ)",
		"hstore",
		"inet",
		"integer",
		"integer[]",
		"interval",
		"json",
		"jsonb",
		"numeric",
		"numeric[]",
		"point",
		"polygon",
		"real",
		"real[]",
		"smallint",
		"smallint[]",
		"text",
		"text[]",
		"timestamp with time zone",
		"timestamp without time zone",
		"tsvector",
		"uuid",
	}

	for _, dataType := range dataTypes {
		var zeroType any
		queryTypeTemplate := ""
		streamTypeTemplate := ""
		typeTemplate := ""
		getOpenAPISchema := GetOpenAPISchemaNotImplemented
		parseFunc := ParseNotImplemented
		parseFuncTemplate := "types.ParseNotImplemented(v)"
		isZeroFunc := IsZeroNotImplemented
		isZeroFuncTemplate := "types.IsZeroNotImplemented"
		formatFunc := FormatNotImplemented
		formatFuncTemplate := "types.FormatNotImplemented"

		// TODO: not an exhaustive suite of implementations for the source type list
		switch dataType {

		case "timestamp without time zone":
			fallthrough
		case "timestamp with time zone":
			zeroType = time.Time{}
			queryTypeTemplate = "time.Time"
			getOpenAPISchema = GetOpenAPISchemaTime
			parseFunc = ParseTime
			parseFuncTemplate = "types.ParseTime(v)"
			isZeroFunc = IsZeroTime
			isZeroFuncTemplate = "types.IsZeroTime"
			formatFunc = FormatTime
			formatFuncTemplate = "types.FormatTime"

		case "interval":
			zeroType = helpers.Deref(new(time.Duration))
			queryTypeTemplate = "time.Duration"
			getOpenAPISchema = GetOpenAPISchemaDuration
			parseFunc = ParseDuration
			parseFuncTemplate = "types.ParseDuration(v)"
			isZeroFunc = IsZeroDuration
			isZeroFuncTemplate = "types.IsZeroDuration"
			formatFunc = FormatDuration
			formatFuncTemplate = "types.FormatDuration"

		case "json":
			fallthrough
		case "jsonb":
			zeroType = nil
			queryTypeTemplate = "any"
			getOpenAPISchema = GetOpenAPISchemaJSON
			parseFunc = ParseJSON
			parseFuncTemplate = "types.ParseJSON(v)"
			isZeroFunc = IsZeroJSON
			isZeroFuncTemplate = "types.IsZeroJSON"
			formatFunc = FormatJSON
			formatFuncTemplate = "types.FormatJSON"

		case "character[]":
			fallthrough
		case "character varying[]":
			fallthrough
		case "text[]":
			zeroType = make([]string, 0)
			queryTypeTemplate = "[]string"
			typeTemplate = "[]string"
			getOpenAPISchema = GetOpenAPISchemaStringArray
			parseFunc = ParseStringArray
			parseFuncTemplate = "types.ParseStringArray(v)"
			isZeroFunc = IsZeroStringArray
			isZeroFuncTemplate = "types.IsZeroStringArray"
			formatFunc = FormatStringArray
			formatFuncTemplate = "types.FormatStringArray"

		case "character":
			fallthrough
		case "character varying":
			fallthrough
		case "text":
			zeroType = helpers.Deref(new(string))
			queryTypeTemplate = "string"
			getOpenAPISchema = GetOpenAPISchemaString
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
			zeroType = make([]int64, 0)
			queryTypeTemplate = "[]int64"
			typeTemplate = "[]int64"
			getOpenAPISchema = GetOpenAPISchemaIntArray
			parseFunc = ParseIntArray
			parseFuncTemplate = "types.ParseIntArray(v)"
			isZeroFunc = IsZeroIntArray
			isZeroFuncTemplate = "types.IsZeroIntArray"
			formatFunc = FormatIntArray
			formatFuncTemplate = "types.FormatIntArray"

		case "smallint":
			fallthrough
		case "integer":
			fallthrough
		case "bigint":
			zeroType = helpers.Deref(new(int64))
			queryTypeTemplate = "int64"
			getOpenAPISchema = GetOpenAPISchemaInt
			parseFunc = ParseInt
			parseFuncTemplate = "types.ParseInt(v)"
			isZeroFunc = IsZeroInt
			isZeroFuncTemplate = "types.IsZeroInt"
			formatFunc = FormatInt
			formatFuncTemplate = "types.FormatInt"

		case "float[]":
			fallthrough
		case "real[]":
			fallthrough
		case "numeric[]":
			fallthrough
		case "double precision[]":
			zeroType = make([]float64, 0)
			queryTypeTemplate = "[]float64"
			typeTemplate = "[]float64"
			getOpenAPISchema = GetOpenAPISchemaFloatArray
			parseFunc = ParseFloatArray
			parseFuncTemplate = "types.ParseFloatArray(v)"
			isZeroFunc = IsZeroFloatArray
			isZeroFuncTemplate = "types.IsZeroFloatArray"
			formatFunc = FormatFloatArray
			formatFuncTemplate = "types.FormatFloatArray"

		case "float":
			fallthrough
		case "real":
			fallthrough
		case "numeric":
			fallthrough
		case "double precision":
			zeroType = helpers.Deref(new(float64))
			queryTypeTemplate = "float64"
			getOpenAPISchema = GetOpenAPISchemaFloat
			parseFunc = ParseFloat
			parseFuncTemplate = "types.ParseFloat(v)"
			isZeroFunc = IsZeroFloat
			isZeroFuncTemplate = "types.IsZeroFloat"
			formatFunc = FormatFloat
			formatFuncTemplate = "types.FormatFloat"

		case "boolean[]":
			zeroType = make([]bool, 0)
			queryTypeTemplate = "[]bool"
			typeTemplate = "[]bool"
			getOpenAPISchema = GetOpenAPISchemaBoolArray
			parseFunc = ParseBoolArray
			parseFuncTemplate = "types.ParseBoolArray(v)"
			isZeroFunc = IsZeroBoolArray
			isZeroFuncTemplate = "types.IsZeroBoolArray"
			formatFunc = FormatBoolArray
			formatFuncTemplate = "types.FormatBoolArray"

		case "boolean":
			zeroType = helpers.Deref(new(bool))
			queryTypeTemplate = "bool"
			getOpenAPISchema = GetOpenAPISchemaBool
			parseFunc = ParseBool
			parseFuncTemplate = "types.ParseBool(v)"
			isZeroFunc = IsZeroBool
			isZeroFuncTemplate = "types.IsZeroBool"
			formatFunc = FormatBool
			formatFuncTemplate = "types.FormatBool"

		case "tsvector":
			zeroType = map[string][]int{}
			queryTypeTemplate = "map[string][]int"
			getOpenAPISchema = GetOpenAPISchemaTSVector
			parseFunc = ParseTSVector
			parseFuncTemplate = "types.ParseTSVector(v)"
			isZeroFunc = IsZeroTSVector
			isZeroFuncTemplate = "types.IsZeroTSVector"
			formatFunc = FormatTSVector
			formatFuncTemplate = "types.FormatTSVector"

		case "uuid":
			zeroType = uuid.UUID{}
			queryTypeTemplate = "uuid.UUID"
			streamTypeTemplate = "[16]uint8"
			getOpenAPISchema = GetOpenAPISchemaUUID
			parseFunc = ParseUUID
			parseFuncTemplate = "types.ParseUUID(v)"
			isZeroFunc = IsZeroUUID
			isZeroFuncTemplate = "types.IsZeroUUID"
			formatFunc = FormatUUID
			formatFuncTemplate = "types.FormatUUID"
			formatFunc = FormatUUID
			formatFuncTemplate = "types.FormatUUID"

		case "hstore":
			zeroType = pgtype.Hstore{}
			queryTypeTemplate = "pgtype.Hstore"
			typeTemplate = "map[string]*string"
			getOpenAPISchema = GetOpenAPISchemaHstore
			parseFunc = ParseHstore
			parseFuncTemplate = "types.ParseHstore(v)"
			isZeroFunc = IsZeroHstore
			isZeroFuncTemplate = "types.IsZeroHstore"
			formatFunc = FormatHstore
			formatFuncTemplate = "types.FormatHstore"

		case "point":
			zeroType = pgtype.Vec2{}
			typeTemplate = "pgtype.Vec2"
			queryTypeTemplate = "pgtype.Point"
			getOpenAPISchema = GetOpenAPISchemaPoint
			parseFunc = ParsePoint
			parseFuncTemplate = "types.ParsePoint(v)"
			isZeroFunc = IsZeroPoint
			isZeroFuncTemplate = "types.IsZeroPoint"
			formatFunc = FormatPoint
			formatFuncTemplate = "types.FormatPoint"

		case "polygon":
			zeroType = []pgtype.Vec2{}
			typeTemplate = "[]pgtype.Vec2"
			queryTypeTemplate = "pgtype.Polygon"
			getOpenAPISchema = GetOpenAPISchemaPolygon
			parseFunc = ParsePolygon
			parseFuncTemplate = "types.ParsePolygon(v)"
			isZeroFunc = IsZeroPolygon
			isZeroFuncTemplate = "types.IsZeroPolygon"
			formatFunc = FormatPolygon
			formatFuncTemplate = "types.FormatPolygon"

		case "geometry", "geometry(PointZ)":
			zeroType = postgis.PointZ{}
			queryTypeTemplate = "postgis.PointZ"
			getOpenAPISchema = GetOpenAPISchemaGeometry
			parseFunc = ParseGeometry
			parseFuncTemplate = "types.ParseGeometry(v)"
			isZeroFunc = IsZeroGeometry
			isZeroFuncTemplate = "types.IsZeroGeometry"
			formatFunc = FormatGeometry
			formatFuncTemplate = "types.FormatGeometry"

		case "inet":
			zeroType = netip.Prefix{}
			queryTypeTemplate = "netip.Prefix"
			getOpenAPISchema = GetOpenAPISchemaInet
			parseFunc = ParseInet
			parseFuncTemplate = "types.ParseInet(v)"
			isZeroFunc = IsZeroInet
			isZeroFuncTemplate = "types.IsZeroInet"
			formatFunc = FormatInet
			formatFuncTemplate = "types.FormatInet"

		case "bytea":
			zeroType = []byte{}
			queryTypeTemplate = "[]byte"
			getOpenAPISchema = GetOpenAPISchemaBytes
			parseFunc = ParseBytes
			parseFuncTemplate = "types.ParseBytes(v)"
			isZeroFunc = IsZeroBytes
			isZeroFuncTemplate = "types.IsZeroBytes"
			formatFunc = FormatBytes
			formatFuncTemplate = "types.FormatBytes"

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

		theType := TypeMeta[any]{
			DataType:           dataType,
			ZeroType:           zeroType,
			QueryTypeTemplate:  queryTypeTemplate,
			StreamTypeTemplate: streamTypeTemplate,
			TypeTemplate:       typeTemplate,
			GetOpenAPISchema:   getOpenAPISchema,
			ParseFunc:          parseFunc,
			ParseFuncTemplate:  parseFuncTemplate,
			IsZeroFunc:         isZeroFunc,
			IsZeroFuncTemplate: isZeroFuncTemplate,
			FormatFunc:         formatFunc,
			FormatFuncTemplate: formatFuncTemplate,
		}

		typeMetas = append(typeMetas, &theType)
	}

	for _, theType := range typeMetas {
		typeMetaByDataType[theType.DataType] = theType
		typeMetaByTypeTemplate[theType.TypeTemplate] = theType
	}
}

func GetTypeMetaForDataType(dataType string) (*TypeMeta[any], error) {
	theType := typeMetaByDataType[dataType]
	if theType == nil {
		return nil, fmt.Errorf("unknown dataType %#+v (out of %v)", dataType, maps.Keys(typeMetaByDataType))
	}

	return theType, nil
}

func GetTypeMetaForTypeTemplate(typeTemplate string) (*TypeMeta[any], error) {
	// TODO: come up with something a bit less hacky, maybe a mapping?
	typeTemplate = strings.ReplaceAll(typeTemplate, "interface {}", "any")
	typeTemplate = strings.ReplaceAll(typeTemplate, "uint8", "byte")

	theType := typeMetaByTypeTemplate[typeTemplate]

	if theType == nil && strings.HasPrefix(typeTemplate, "*") {
		theType = typeMetaByTypeTemplate[typeTemplate[1:]]
	}

	if theType == nil && strings.HasSuffix(typeTemplate, "int") {
		typeTemplate += "64"
		theType = typeMetaByTypeTemplate[typeTemplate]
	}

	if theType == nil {
		return nil, fmt.Errorf("unknown typeTemplate %#+v", typeTemplate)
	}

	return theType, nil
}

func GetTypeMetas() []*TypeMeta[any] {
	return typeMetas
}
