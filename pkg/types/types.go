package types

import (
	"fmt"
	"net/netip"
	"time"

	"github.com/google/uuid"
	"github.com/initialed85/djangolang/pkg/helpers"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/paulmach/orb/geojson"

	"github.com/lib/pq"
	"github.com/lib/pq/hstore"

	"github.com/cridenour/go-postgis"
	jsoniter "github.com/json-iterator/go"
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
		"timestamp without time zone",
		"timestamp with time zone",
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
		"tsvector",
		"uuid",
		"hstore",
		"point",
		"polygon",
		"geometry",
		"geometry(PointZ)",
		"inet",
		"bytea",
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

		case "interval":
			zeroType = helpers.Deref(new(time.Duration))
			queryTypeTemplate = "time.Duration"
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

		case "tsvector":
			zeroType = map[string][]int{}
			queryTypeTemplate = "map[string][]int"
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
			parseFunc = ParseUUID
			parseFuncTemplate = "types.ParseUUID(v)"
			isZeroFunc = IsZeroUUID
			isZeroFuncTemplate = "types.IsZeroUUID"
			formatFunc = FormatUUID
			formatFuncTemplate = "types.FormatUUID"
			formatFunc = FormatUUID
			formatFuncTemplate = "types.FormatUUID"

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

		case "geometry", "geometry(PointZ)":
			zeroType = postgis.PointZ{}
			queryTypeTemplate = "postgis.PointZ"
			parseFunc = ParseGeometry
			parseFuncTemplate = "types.ParseGeometry(v)"
			isZeroFunc = IsZeroGeometry
			isZeroFuncTemplate = "types.IsZeroGeometry"
			formatFunc = FormatGeometry
			formatFuncTemplate = "types.FormatGeometry"

		case "inet":
			zeroType = netip.Prefix{}
			queryTypeTemplate = "netip.Prefix"
			parseFunc = ParseInet
			parseFuncTemplate = "types.ParseInet(v)"
			isZeroFunc = IsZeroInet
			isZeroFuncTemplate = "types.IsZeroInet"
			formatFunc = FormatInet
			formatFuncTemplate = "types.FormatInet"

		case "bytea":
			zeroType = []byte{}
			queryTypeTemplate = "[]byte"
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
