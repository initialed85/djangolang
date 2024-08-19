package model_generated

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/netip"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/cridenour/go-postgis"
	"github.com/go-chi/chi/v5"
	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
	"github.com/initialed85/djangolang/pkg/helpers"
	"github.com/initialed85/djangolang/pkg/introspect"
	"github.com/initialed85/djangolang/pkg/query"
	"github.com/initialed85/djangolang/pkg/server"
	"github.com/initialed85/djangolang/pkg/stream"
	"github.com/initialed85/djangolang/pkg/types"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/lib/pq/hstore"
	"golang.org/x/exp/maps"
)

type NotNullFuzz struct {
	ID                        uuid.UUID          `json:"id"`
	SomeBigint                int64              `json:"some_bigint"`
	SomeBigintArray           []int64            `json:"some_bigint_array"`
	SomeBoolean               bool               `json:"some_boolean"`
	SomeBooleanArray          []bool             `json:"some_boolean_array"`
	SomeBytea                 []byte             `json:"some_bytea"`
	SomeCharacterVarying      string             `json:"some_character_varying"`
	SomeCharacterVaryingArray []string           `json:"some_character_varying_array"`
	SomeDoublePrecision       float64            `json:"some_double_precision"`
	SomeDoublePrecisionArray  []float64          `json:"some_double_precision_array"`
	SomeFloat                 float64            `json:"some_float"`
	SomeFloatArray            []float64          `json:"some_float_array"`
	SomeGeometryPointZ        postgis.PointZ     `json:"some_geometry_point_z"`
	SomeHstore                map[string]*string `json:"some_hstore"`
	SomeInet                  netip.Prefix       `json:"some_inet"`
	SomeInteger               int64              `json:"some_integer"`
	SomeIntegerArray          []int64            `json:"some_integer_array"`
	SomeInterval              time.Duration      `json:"some_interval"`
	SomeJSON                  any                `json:"some_json"`
	SomeJSONB                 any                `json:"some_jsonb"`
	SomeNumeric               float64            `json:"some_numeric"`
	SomeNumericArray          []float64          `json:"some_numeric_array"`
	SomePoint                 pgtype.Vec2        `json:"some_point"`
	SomePolygon               []pgtype.Vec2      `json:"some_polygon"`
	SomeReal                  float64            `json:"some_real"`
	SomeRealArray             []float64          `json:"some_real_array"`
	SomeSmallint              int64              `json:"some_smallint"`
	SomeSmallintArray         []int64            `json:"some_smallint_array"`
	SomeText                  string             `json:"some_text"`
	SomeTextArray             []string           `json:"some_text_array"`
	SomeTimestamptz           time.Time          `json:"some_timestamptz"`
	SomeTimestamp             time.Time          `json:"some_timestamp"`
	SomeTsvector              map[string][]int   `json:"some_tsvector"`
	SomeUUID                  uuid.UUID          `json:"some_uuid"`
}

var NotNullFuzzTable = "not_null_fuzz"

var (
	NotNullFuzzTableIDColumn                        = "id"
	NotNullFuzzTableSomeBigintColumn                = "some_bigint"
	NotNullFuzzTableSomeBigintArrayColumn           = "some_bigint_array"
	NotNullFuzzTableSomeBooleanColumn               = "some_boolean"
	NotNullFuzzTableSomeBooleanArrayColumn          = "some_boolean_array"
	NotNullFuzzTableSomeByteaColumn                 = "some_bytea"
	NotNullFuzzTableSomeCharacterVaryingColumn      = "some_character_varying"
	NotNullFuzzTableSomeCharacterVaryingArrayColumn = "some_character_varying_array"
	NotNullFuzzTableSomeDoublePrecisionColumn       = "some_double_precision"
	NotNullFuzzTableSomeDoublePrecisionArrayColumn  = "some_double_precision_array"
	NotNullFuzzTableSomeFloatColumn                 = "some_float"
	NotNullFuzzTableSomeFloatArrayColumn            = "some_float_array"
	NotNullFuzzTableSomeGeometryPointZColumn        = "some_geometry_point_z"
	NotNullFuzzTableSomeHstoreColumn                = "some_hstore"
	NotNullFuzzTableSomeInetColumn                  = "some_inet"
	NotNullFuzzTableSomeIntegerColumn               = "some_integer"
	NotNullFuzzTableSomeIntegerArrayColumn          = "some_integer_array"
	NotNullFuzzTableSomeIntervalColumn              = "some_interval"
	NotNullFuzzTableSomeJSONColumn                  = "some_json"
	NotNullFuzzTableSomeJSONBColumn                 = "some_jsonb"
	NotNullFuzzTableSomeNumericColumn               = "some_numeric"
	NotNullFuzzTableSomeNumericArrayColumn          = "some_numeric_array"
	NotNullFuzzTableSomePointColumn                 = "some_point"
	NotNullFuzzTableSomePolygonColumn               = "some_polygon"
	NotNullFuzzTableSomeRealColumn                  = "some_real"
	NotNullFuzzTableSomeRealArrayColumn             = "some_real_array"
	NotNullFuzzTableSomeSmallintColumn              = "some_smallint"
	NotNullFuzzTableSomeSmallintArrayColumn         = "some_smallint_array"
	NotNullFuzzTableSomeTextColumn                  = "some_text"
	NotNullFuzzTableSomeTextArrayColumn             = "some_text_array"
	NotNullFuzzTableSomeTimestamptzColumn           = "some_timestamptz"
	NotNullFuzzTableSomeTimestampColumn             = "some_timestamp"
	NotNullFuzzTableSomeTsvectorColumn              = "some_tsvector"
	NotNullFuzzTableSomeUUIDColumn                  = "some_uuid"
)

var (
	NotNullFuzzTableIDColumnWithTypeCast                        = fmt.Sprintf(`"id" AS id`)
	NotNullFuzzTableSomeBigintColumnWithTypeCast                = fmt.Sprintf(`"some_bigint" AS some_bigint`)
	NotNullFuzzTableSomeBigintArrayColumnWithTypeCast           = fmt.Sprintf(`"some_bigint_array" AS some_bigint_array`)
	NotNullFuzzTableSomeBooleanColumnWithTypeCast               = fmt.Sprintf(`"some_boolean" AS some_boolean`)
	NotNullFuzzTableSomeBooleanArrayColumnWithTypeCast          = fmt.Sprintf(`"some_boolean_array" AS some_boolean_array`)
	NotNullFuzzTableSomeByteaColumnWithTypeCast                 = fmt.Sprintf(`"some_bytea" AS some_bytea`)
	NotNullFuzzTableSomeCharacterVaryingColumnWithTypeCast      = fmt.Sprintf(`"some_character_varying" AS some_character_varying`)
	NotNullFuzzTableSomeCharacterVaryingArrayColumnWithTypeCast = fmt.Sprintf(`"some_character_varying_array" AS some_character_varying_array`)
	NotNullFuzzTableSomeDoublePrecisionColumnWithTypeCast       = fmt.Sprintf(`"some_double_precision" AS some_double_precision`)
	NotNullFuzzTableSomeDoublePrecisionArrayColumnWithTypeCast  = fmt.Sprintf(`"some_double_precision_array" AS some_double_precision_array`)
	NotNullFuzzTableSomeFloatColumnWithTypeCast                 = fmt.Sprintf(`"some_float" AS some_float`)
	NotNullFuzzTableSomeFloatArrayColumnWithTypeCast            = fmt.Sprintf(`"some_float_array" AS some_float_array`)
	NotNullFuzzTableSomeGeometryPointZColumnWithTypeCast        = fmt.Sprintf(`"some_geometry_point_z" AS some_geometry_point_z`)
	NotNullFuzzTableSomeHstoreColumnWithTypeCast                = fmt.Sprintf(`"some_hstore" AS some_hstore`)
	NotNullFuzzTableSomeInetColumnWithTypeCast                  = fmt.Sprintf(`"some_inet" AS some_inet`)
	NotNullFuzzTableSomeIntegerColumnWithTypeCast               = fmt.Sprintf(`"some_integer" AS some_integer`)
	NotNullFuzzTableSomeIntegerArrayColumnWithTypeCast          = fmt.Sprintf(`"some_integer_array" AS some_integer_array`)
	NotNullFuzzTableSomeIntervalColumnWithTypeCast              = fmt.Sprintf(`"some_interval" AS some_interval`)
	NotNullFuzzTableSomeJSONColumnWithTypeCast                  = fmt.Sprintf(`"some_json" AS some_json`)
	NotNullFuzzTableSomeJSONBColumnWithTypeCast                 = fmt.Sprintf(`"some_jsonb" AS some_jsonb`)
	NotNullFuzzTableSomeNumericColumnWithTypeCast               = fmt.Sprintf(`"some_numeric" AS some_numeric`)
	NotNullFuzzTableSomeNumericArrayColumnWithTypeCast          = fmt.Sprintf(`"some_numeric_array" AS some_numeric_array`)
	NotNullFuzzTableSomePointColumnWithTypeCast                 = fmt.Sprintf(`"some_point" AS some_point`)
	NotNullFuzzTableSomePolygonColumnWithTypeCast               = fmt.Sprintf(`"some_polygon" AS some_polygon`)
	NotNullFuzzTableSomeRealColumnWithTypeCast                  = fmt.Sprintf(`"some_real" AS some_real`)
	NotNullFuzzTableSomeRealArrayColumnWithTypeCast             = fmt.Sprintf(`"some_real_array" AS some_real_array`)
	NotNullFuzzTableSomeSmallintColumnWithTypeCast              = fmt.Sprintf(`"some_smallint" AS some_smallint`)
	NotNullFuzzTableSomeSmallintArrayColumnWithTypeCast         = fmt.Sprintf(`"some_smallint_array" AS some_smallint_array`)
	NotNullFuzzTableSomeTextColumnWithTypeCast                  = fmt.Sprintf(`"some_text" AS some_text`)
	NotNullFuzzTableSomeTextArrayColumnWithTypeCast             = fmt.Sprintf(`"some_text_array" AS some_text_array`)
	NotNullFuzzTableSomeTimestamptzColumnWithTypeCast           = fmt.Sprintf(`"some_timestamptz" AS some_timestamptz`)
	NotNullFuzzTableSomeTimestampColumnWithTypeCast             = fmt.Sprintf(`"some_timestamp" AS some_timestamp`)
	NotNullFuzzTableSomeTsvectorColumnWithTypeCast              = fmt.Sprintf(`"some_tsvector" AS some_tsvector`)
	NotNullFuzzTableSomeUUIDColumnWithTypeCast                  = fmt.Sprintf(`"some_uuid" AS some_uuid`)
)

var NotNullFuzzTableColumns = []string{
	NotNullFuzzTableIDColumn,
	NotNullFuzzTableSomeBigintColumn,
	NotNullFuzzTableSomeBigintArrayColumn,
	NotNullFuzzTableSomeBooleanColumn,
	NotNullFuzzTableSomeBooleanArrayColumn,
	NotNullFuzzTableSomeByteaColumn,
	NotNullFuzzTableSomeCharacterVaryingColumn,
	NotNullFuzzTableSomeCharacterVaryingArrayColumn,
	NotNullFuzzTableSomeDoublePrecisionColumn,
	NotNullFuzzTableSomeDoublePrecisionArrayColumn,
	NotNullFuzzTableSomeFloatColumn,
	NotNullFuzzTableSomeFloatArrayColumn,
	NotNullFuzzTableSomeGeometryPointZColumn,
	NotNullFuzzTableSomeHstoreColumn,
	NotNullFuzzTableSomeInetColumn,
	NotNullFuzzTableSomeIntegerColumn,
	NotNullFuzzTableSomeIntegerArrayColumn,
	NotNullFuzzTableSomeIntervalColumn,
	NotNullFuzzTableSomeJSONColumn,
	NotNullFuzzTableSomeJSONBColumn,
	NotNullFuzzTableSomeNumericColumn,
	NotNullFuzzTableSomeNumericArrayColumn,
	NotNullFuzzTableSomePointColumn,
	NotNullFuzzTableSomePolygonColumn,
	NotNullFuzzTableSomeRealColumn,
	NotNullFuzzTableSomeRealArrayColumn,
	NotNullFuzzTableSomeSmallintColumn,
	NotNullFuzzTableSomeSmallintArrayColumn,
	NotNullFuzzTableSomeTextColumn,
	NotNullFuzzTableSomeTextArrayColumn,
	NotNullFuzzTableSomeTimestamptzColumn,
	NotNullFuzzTableSomeTimestampColumn,
	NotNullFuzzTableSomeTsvectorColumn,
	NotNullFuzzTableSomeUUIDColumn,
}

var NotNullFuzzTableColumnsWithTypeCasts = []string{
	NotNullFuzzTableIDColumnWithTypeCast,
	NotNullFuzzTableSomeBigintColumnWithTypeCast,
	NotNullFuzzTableSomeBigintArrayColumnWithTypeCast,
	NotNullFuzzTableSomeBooleanColumnWithTypeCast,
	NotNullFuzzTableSomeBooleanArrayColumnWithTypeCast,
	NotNullFuzzTableSomeByteaColumnWithTypeCast,
	NotNullFuzzTableSomeCharacterVaryingColumnWithTypeCast,
	NotNullFuzzTableSomeCharacterVaryingArrayColumnWithTypeCast,
	NotNullFuzzTableSomeDoublePrecisionColumnWithTypeCast,
	NotNullFuzzTableSomeDoublePrecisionArrayColumnWithTypeCast,
	NotNullFuzzTableSomeFloatColumnWithTypeCast,
	NotNullFuzzTableSomeFloatArrayColumnWithTypeCast,
	NotNullFuzzTableSomeGeometryPointZColumnWithTypeCast,
	NotNullFuzzTableSomeHstoreColumnWithTypeCast,
	NotNullFuzzTableSomeInetColumnWithTypeCast,
	NotNullFuzzTableSomeIntegerColumnWithTypeCast,
	NotNullFuzzTableSomeIntegerArrayColumnWithTypeCast,
	NotNullFuzzTableSomeIntervalColumnWithTypeCast,
	NotNullFuzzTableSomeJSONColumnWithTypeCast,
	NotNullFuzzTableSomeJSONBColumnWithTypeCast,
	NotNullFuzzTableSomeNumericColumnWithTypeCast,
	NotNullFuzzTableSomeNumericArrayColumnWithTypeCast,
	NotNullFuzzTableSomePointColumnWithTypeCast,
	NotNullFuzzTableSomePolygonColumnWithTypeCast,
	NotNullFuzzTableSomeRealColumnWithTypeCast,
	NotNullFuzzTableSomeRealArrayColumnWithTypeCast,
	NotNullFuzzTableSomeSmallintColumnWithTypeCast,
	NotNullFuzzTableSomeSmallintArrayColumnWithTypeCast,
	NotNullFuzzTableSomeTextColumnWithTypeCast,
	NotNullFuzzTableSomeTextArrayColumnWithTypeCast,
	NotNullFuzzTableSomeTimestamptzColumnWithTypeCast,
	NotNullFuzzTableSomeTimestampColumnWithTypeCast,
	NotNullFuzzTableSomeTsvectorColumnWithTypeCast,
	NotNullFuzzTableSomeUUIDColumnWithTypeCast,
}

var NotNullFuzzTableColumnLookup = map[string]*introspect.Column{
	NotNullFuzzTableIDColumn:                        {Name: NotNullFuzzTableIDColumn, NotNull: true, HasDefault: true},
	NotNullFuzzTableSomeBigintColumn:                {Name: NotNullFuzzTableSomeBigintColumn, NotNull: true, HasDefault: true},
	NotNullFuzzTableSomeBigintArrayColumn:           {Name: NotNullFuzzTableSomeBigintArrayColumn, NotNull: true, HasDefault: true},
	NotNullFuzzTableSomeBooleanColumn:               {Name: NotNullFuzzTableSomeBooleanColumn, NotNull: true, HasDefault: true},
	NotNullFuzzTableSomeBooleanArrayColumn:          {Name: NotNullFuzzTableSomeBooleanArrayColumn, NotNull: true, HasDefault: true},
	NotNullFuzzTableSomeByteaColumn:                 {Name: NotNullFuzzTableSomeByteaColumn, NotNull: true, HasDefault: true},
	NotNullFuzzTableSomeCharacterVaryingColumn:      {Name: NotNullFuzzTableSomeCharacterVaryingColumn, NotNull: true, HasDefault: true},
	NotNullFuzzTableSomeCharacterVaryingArrayColumn: {Name: NotNullFuzzTableSomeCharacterVaryingArrayColumn, NotNull: true, HasDefault: true},
	NotNullFuzzTableSomeDoublePrecisionColumn:       {Name: NotNullFuzzTableSomeDoublePrecisionColumn, NotNull: true, HasDefault: true},
	NotNullFuzzTableSomeDoublePrecisionArrayColumn:  {Name: NotNullFuzzTableSomeDoublePrecisionArrayColumn, NotNull: true, HasDefault: true},
	NotNullFuzzTableSomeFloatColumn:                 {Name: NotNullFuzzTableSomeFloatColumn, NotNull: true, HasDefault: true},
	NotNullFuzzTableSomeFloatArrayColumn:            {Name: NotNullFuzzTableSomeFloatArrayColumn, NotNull: true, HasDefault: true},
	NotNullFuzzTableSomeGeometryPointZColumn:        {Name: NotNullFuzzTableSomeGeometryPointZColumn, NotNull: true, HasDefault: true},
	NotNullFuzzTableSomeHstoreColumn:                {Name: NotNullFuzzTableSomeHstoreColumn, NotNull: true, HasDefault: true},
	NotNullFuzzTableSomeInetColumn:                  {Name: NotNullFuzzTableSomeInetColumn, NotNull: true, HasDefault: true},
	NotNullFuzzTableSomeIntegerColumn:               {Name: NotNullFuzzTableSomeIntegerColumn, NotNull: true, HasDefault: true},
	NotNullFuzzTableSomeIntegerArrayColumn:          {Name: NotNullFuzzTableSomeIntegerArrayColumn, NotNull: true, HasDefault: true},
	NotNullFuzzTableSomeIntervalColumn:              {Name: NotNullFuzzTableSomeIntervalColumn, NotNull: true, HasDefault: true},
	NotNullFuzzTableSomeJSONColumn:                  {Name: NotNullFuzzTableSomeJSONColumn, NotNull: true, HasDefault: true},
	NotNullFuzzTableSomeJSONBColumn:                 {Name: NotNullFuzzTableSomeJSONBColumn, NotNull: true, HasDefault: true},
	NotNullFuzzTableSomeNumericColumn:               {Name: NotNullFuzzTableSomeNumericColumn, NotNull: true, HasDefault: true},
	NotNullFuzzTableSomeNumericArrayColumn:          {Name: NotNullFuzzTableSomeNumericArrayColumn, NotNull: true, HasDefault: true},
	NotNullFuzzTableSomePointColumn:                 {Name: NotNullFuzzTableSomePointColumn, NotNull: true, HasDefault: true},
	NotNullFuzzTableSomePolygonColumn:               {Name: NotNullFuzzTableSomePolygonColumn, NotNull: true, HasDefault: true},
	NotNullFuzzTableSomeRealColumn:                  {Name: NotNullFuzzTableSomeRealColumn, NotNull: true, HasDefault: true},
	NotNullFuzzTableSomeRealArrayColumn:             {Name: NotNullFuzzTableSomeRealArrayColumn, NotNull: true, HasDefault: true},
	NotNullFuzzTableSomeSmallintColumn:              {Name: NotNullFuzzTableSomeSmallintColumn, NotNull: true, HasDefault: true},
	NotNullFuzzTableSomeSmallintArrayColumn:         {Name: NotNullFuzzTableSomeSmallintArrayColumn, NotNull: true, HasDefault: true},
	NotNullFuzzTableSomeTextColumn:                  {Name: NotNullFuzzTableSomeTextColumn, NotNull: true, HasDefault: true},
	NotNullFuzzTableSomeTextArrayColumn:             {Name: NotNullFuzzTableSomeTextArrayColumn, NotNull: true, HasDefault: true},
	NotNullFuzzTableSomeTimestamptzColumn:           {Name: NotNullFuzzTableSomeTimestamptzColumn, NotNull: true, HasDefault: true},
	NotNullFuzzTableSomeTimestampColumn:             {Name: NotNullFuzzTableSomeTimestampColumn, NotNull: true, HasDefault: true},
	NotNullFuzzTableSomeTsvectorColumn:              {Name: NotNullFuzzTableSomeTsvectorColumn, NotNull: true, HasDefault: true},
	NotNullFuzzTableSomeUUIDColumn:                  {Name: NotNullFuzzTableSomeUUIDColumn, NotNull: true, HasDefault: true},
}

var (
	NotNullFuzzTablePrimaryKeyColumn = NotNullFuzzTableIDColumn
)
var _ = []any{
	time.Time{},
	time.Duration(0),
	nil,
	pq.StringArray{},
	string(""),
	pq.Int64Array{},
	int64(0),
	pq.Float64Array{},
	float64(0),
	pq.BoolArray{},
	bool(false),
	map[string][]int{},
	uuid.UUID{},
	hstore.Hstore{},
	pgtype.Point{},
	pgtype.Polygon{},
	postgis.PointZ{},
	netip.Prefix{},
	[]byte{},
	errors.Is,
}

func (m *NotNullFuzz) GetPrimaryKeyColumn() string {
	return NotNullFuzzTablePrimaryKeyColumn
}

func (m *NotNullFuzz) GetPrimaryKeyValue() any {
	return m.ID
}

func (m *NotNullFuzz) FromItem(item map[string]any) error {
	if item == nil {
		return fmt.Errorf(
			"item unexpectedly nil during NotNullFuzzFromItem",
		)
	}

	if len(item) == 0 {
		return fmt.Errorf(
			"item unexpectedly empty during NotNullFuzzFromItem",
		)
	}

	wrapError := func(k string, v any, err error) error {
		return fmt.Errorf("%v: %#+v; error: %v", k, v, err)
	}

	for k, v := range item {
		_, ok := NotNullFuzzTableColumnLookup[k]
		if !ok {
			return fmt.Errorf(
				"item contained unexpected key %#+v during NotNullFuzzFromItem; item: %#+v",
				k, item,
			)
		}

		switch k {
		case "id":
			if v == nil {
				continue
			}

			temp1, err := types.ParseUUID(v)
			if err != nil {
				return wrapError(k, v, err)
			}

			temp2, ok := temp1.(uuid.UUID)
			if !ok {
				if temp1 != nil {
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uuid.UUID", temp1))
				}
			}

			m.ID = temp2

		case "some_bigint":
			if v == nil {
				continue
			}

			temp1, err := types.ParseInt(v)
			if err != nil {
				return wrapError(k, v, err)
			}

			temp2, ok := temp1.(int64)
			if !ok {
				if temp1 != nil {
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to int64", temp1))
				}
			}

			m.SomeBigint = temp2

		case "some_bigint_array":
			if v == nil {
				continue
			}

			temp1, err := types.ParseIntArray(v)
			if err != nil {
				return wrapError(k, v, err)
			}

			temp2, ok := temp1.([]int64)
			if !ok {
				if temp1 != nil {
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to []int64", temp1))
				}
			}

			m.SomeBigintArray = temp2

		case "some_boolean":
			if v == nil {
				continue
			}

			temp1, err := types.ParseBool(v)
			if err != nil {
				return wrapError(k, v, err)
			}

			temp2, ok := temp1.(bool)
			if !ok {
				if temp1 != nil {
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to bool", temp1))
				}
			}

			m.SomeBoolean = temp2

		case "some_boolean_array":
			if v == nil {
				continue
			}

			temp1, err := types.ParseBoolArray(v)
			if err != nil {
				return wrapError(k, v, err)
			}

			temp2, ok := temp1.([]bool)
			if !ok {
				if temp1 != nil {
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to []bool", temp1))
				}
			}

			m.SomeBooleanArray = temp2

		case "some_bytea":
			if v == nil {
				continue
			}

			temp1, err := types.ParseBytes(v)
			if err != nil {
				return wrapError(k, v, err)
			}

			temp2, ok := temp1.([]byte)
			if !ok {
				if temp1 != nil {
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to []byte", temp1))
				}
			}

			m.SomeBytea = temp2

		case "some_character_varying":
			if v == nil {
				continue
			}

			temp1, err := types.ParseString(v)
			if err != nil {
				return wrapError(k, v, err)
			}

			temp2, ok := temp1.(string)
			if !ok {
				if temp1 != nil {
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to string", temp1))
				}
			}

			m.SomeCharacterVarying = temp2

		case "some_character_varying_array":
			if v == nil {
				continue
			}

			temp1, err := types.ParseStringArray(v)
			if err != nil {
				return wrapError(k, v, err)
			}

			temp2, ok := temp1.([]string)
			if !ok {
				if temp1 != nil {
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to []string", temp1))
				}
			}

			m.SomeCharacterVaryingArray = temp2

		case "some_double_precision":
			if v == nil {
				continue
			}

			temp1, err := types.ParseFloat(v)
			if err != nil {
				return wrapError(k, v, err)
			}

			temp2, ok := temp1.(float64)
			if !ok {
				if temp1 != nil {
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to float64", temp1))
				}
			}

			m.SomeDoublePrecision = temp2

		case "some_double_precision_array":
			if v == nil {
				continue
			}

			temp1, err := types.ParseFloatArray(v)
			if err != nil {
				return wrapError(k, v, err)
			}

			temp2, ok := temp1.([]float64)
			if !ok {
				if temp1 != nil {
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to []float64", temp1))
				}
			}

			m.SomeDoublePrecisionArray = temp2

		case "some_float":
			if v == nil {
				continue
			}

			temp1, err := types.ParseFloat(v)
			if err != nil {
				return wrapError(k, v, err)
			}

			temp2, ok := temp1.(float64)
			if !ok {
				if temp1 != nil {
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to float64", temp1))
				}
			}

			m.SomeFloat = temp2

		case "some_float_array":
			if v == nil {
				continue
			}

			temp1, err := types.ParseFloatArray(v)
			if err != nil {
				return wrapError(k, v, err)
			}

			temp2, ok := temp1.([]float64)
			if !ok {
				if temp1 != nil {
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to []float64", temp1))
				}
			}

			m.SomeFloatArray = temp2

		case "some_geometry_point_z":
			if v == nil {
				continue
			}

			temp1, err := types.ParseGeometry(v)
			if err != nil {
				return wrapError(k, v, err)
			}

			temp2, ok := temp1.(postgis.PointZ)
			if !ok {
				if temp1 != nil {
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to postgis.PointZ", temp1))
				}
			}

			m.SomeGeometryPointZ = temp2

		case "some_hstore":
			if v == nil {
				continue
			}

			temp1, err := types.ParseHstore(v)
			if err != nil {
				return wrapError(k, v, err)
			}

			temp2, ok := temp1.(map[string]*string)
			if !ok {
				if temp1 != nil {
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to map[string]*string", temp1))
				}
			}

			m.SomeHstore = temp2

		case "some_inet":
			if v == nil {
				continue
			}

			temp1, err := types.ParseInet(v)
			if err != nil {
				return wrapError(k, v, err)
			}

			temp2, ok := temp1.(netip.Prefix)
			if !ok {
				if temp1 != nil {
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to netip.Prefix", temp1))
				}
			}

			m.SomeInet = temp2

		case "some_integer":
			if v == nil {
				continue
			}

			temp1, err := types.ParseInt(v)
			if err != nil {
				return wrapError(k, v, err)
			}

			temp2, ok := temp1.(int64)
			if !ok {
				if temp1 != nil {
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to int64", temp1))
				}
			}

			m.SomeInteger = temp2

		case "some_integer_array":
			if v == nil {
				continue
			}

			temp1, err := types.ParseIntArray(v)
			if err != nil {
				return wrapError(k, v, err)
			}

			temp2, ok := temp1.([]int64)
			if !ok {
				if temp1 != nil {
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to []int64", temp1))
				}
			}

			m.SomeIntegerArray = temp2

		case "some_interval":
			if v == nil {
				continue
			}

			temp1, err := types.ParseDuration(v)
			if err != nil {
				return wrapError(k, v, err)
			}

			temp2, ok := temp1.(time.Duration)
			if !ok {
				if temp1 != nil {
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to time.Duration", temp1))
				}
			}

			m.SomeInterval = temp2

		case "some_json":
			if v == nil {
				continue
			}

			temp1, err := types.ParseJSON(v)
			if err != nil {
				return wrapError(k, v, err)
			}

			temp2, ok := temp1.(any)
			if !ok {
				if temp1 != nil {
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to any", temp1))
				}
			}

			m.SomeJSON = temp2

		case "some_jsonb":
			if v == nil {
				continue
			}

			temp1, err := types.ParseJSON(v)
			if err != nil {
				return wrapError(k, v, err)
			}

			temp2, ok := temp1.(any)
			if !ok {
				if temp1 != nil {
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to any", temp1))
				}
			}

			m.SomeJSONB = temp2

		case "some_numeric":
			if v == nil {
				continue
			}

			temp1, err := types.ParseFloat(v)
			if err != nil {
				return wrapError(k, v, err)
			}

			temp2, ok := temp1.(float64)
			if !ok {
				if temp1 != nil {
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to float64", temp1))
				}
			}

			m.SomeNumeric = temp2

		case "some_numeric_array":
			if v == nil {
				continue
			}

			temp1, err := types.ParseFloatArray(v)
			if err != nil {
				return wrapError(k, v, err)
			}

			temp2, ok := temp1.([]float64)
			if !ok {
				if temp1 != nil {
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to []float64", temp1))
				}
			}

			m.SomeNumericArray = temp2

		case "some_point":
			if v == nil {
				continue
			}

			temp1, err := types.ParsePoint(v)
			if err != nil {
				return wrapError(k, v, err)
			}

			temp2, ok := temp1.(pgtype.Vec2)
			if !ok {
				if temp1 != nil {
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to pgtype.Vec2", temp1))
				}
			}

			m.SomePoint = temp2

		case "some_polygon":
			if v == nil {
				continue
			}

			temp1, err := types.ParsePolygon(v)
			if err != nil {
				return wrapError(k, v, err)
			}

			temp2, ok := temp1.([]pgtype.Vec2)
			if !ok {
				if temp1 != nil {
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to []pgtype.Vec2", temp1))
				}
			}

			m.SomePolygon = temp2

		case "some_real":
			if v == nil {
				continue
			}

			temp1, err := types.ParseFloat(v)
			if err != nil {
				return wrapError(k, v, err)
			}

			temp2, ok := temp1.(float64)
			if !ok {
				if temp1 != nil {
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to float64", temp1))
				}
			}

			m.SomeReal = temp2

		case "some_real_array":
			if v == nil {
				continue
			}

			temp1, err := types.ParseFloatArray(v)
			if err != nil {
				return wrapError(k, v, err)
			}

			temp2, ok := temp1.([]float64)
			if !ok {
				if temp1 != nil {
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to []float64", temp1))
				}
			}

			m.SomeRealArray = temp2

		case "some_smallint":
			if v == nil {
				continue
			}

			temp1, err := types.ParseInt(v)
			if err != nil {
				return wrapError(k, v, err)
			}

			temp2, ok := temp1.(int64)
			if !ok {
				if temp1 != nil {
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to int64", temp1))
				}
			}

			m.SomeSmallint = temp2

		case "some_smallint_array":
			if v == nil {
				continue
			}

			temp1, err := types.ParseIntArray(v)
			if err != nil {
				return wrapError(k, v, err)
			}

			temp2, ok := temp1.([]int64)
			if !ok {
				if temp1 != nil {
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to []int64", temp1))
				}
			}

			m.SomeSmallintArray = temp2

		case "some_text":
			if v == nil {
				continue
			}

			temp1, err := types.ParseString(v)
			if err != nil {
				return wrapError(k, v, err)
			}

			temp2, ok := temp1.(string)
			if !ok {
				if temp1 != nil {
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to string", temp1))
				}
			}

			m.SomeText = temp2

		case "some_text_array":
			if v == nil {
				continue
			}

			temp1, err := types.ParseStringArray(v)
			if err != nil {
				return wrapError(k, v, err)
			}

			temp2, ok := temp1.([]string)
			if !ok {
				if temp1 != nil {
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to []string", temp1))
				}
			}

			m.SomeTextArray = temp2

		case "some_timestamptz":
			if v == nil {
				continue
			}

			temp1, err := types.ParseTime(v)
			if err != nil {
				return wrapError(k, v, err)
			}

			temp2, ok := temp1.(time.Time)
			if !ok {
				if temp1 != nil {
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to time.Time", temp1))
				}
			}

			m.SomeTimestamptz = temp2

		case "some_timestamp":
			if v == nil {
				continue
			}

			temp1, err := types.ParseTime(v)
			if err != nil {
				return wrapError(k, v, err)
			}

			temp2, ok := temp1.(time.Time)
			if !ok {
				if temp1 != nil {
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to time.Time", temp1))
				}
			}

			m.SomeTimestamp = temp2

		case "some_tsvector":
			if v == nil {
				continue
			}

			temp1, err := types.ParseTSVector(v)
			if err != nil {
				return wrapError(k, v, err)
			}

			temp2, ok := temp1.(map[string][]int)
			if !ok {
				if temp1 != nil {
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to map[string][]int", temp1))
				}
			}

			m.SomeTsvector = temp2

		case "some_uuid":
			if v == nil {
				continue
			}

			temp1, err := types.ParseUUID(v)
			if err != nil {
				return wrapError(k, v, err)
			}

			temp2, ok := temp1.(uuid.UUID)
			if !ok {
				if temp1 != nil {
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uuid.UUID", temp1))
				}
			}

			m.SomeUUID = temp2

		}
	}

	return nil
}

func (m *NotNullFuzz) Reload(
	ctx context.Context,
	tx *sqlx.Tx,
	includeDeleteds ...bool,
) error {
	extraWhere := ""
	if len(includeDeleteds) > 0 && includeDeleteds[0] {
		if slices.Contains(NotNullFuzzTableColumns, "deleted_at") {
			extraWhere = "\n    AND (deleted_at IS null OR deleted_at IS NOT null)"
		}
	}

	ctx, cleanup := query.WithQueryID(ctx)
	defer cleanup()

	t, err := SelectNotNullFuzz(
		ctx,
		tx,
		fmt.Sprintf("%v = $1%v", m.GetPrimaryKeyColumn(), extraWhere),
		m.GetPrimaryKeyValue(),
	)
	if err != nil {
		return err
	}

	m.ID = t.ID
	m.SomeBigint = t.SomeBigint
	m.SomeBigintArray = t.SomeBigintArray
	m.SomeBoolean = t.SomeBoolean
	m.SomeBooleanArray = t.SomeBooleanArray
	m.SomeBytea = t.SomeBytea
	m.SomeCharacterVarying = t.SomeCharacterVarying
	m.SomeCharacterVaryingArray = t.SomeCharacterVaryingArray
	m.SomeDoublePrecision = t.SomeDoublePrecision
	m.SomeDoublePrecisionArray = t.SomeDoublePrecisionArray
	m.SomeFloat = t.SomeFloat
	m.SomeFloatArray = t.SomeFloatArray
	m.SomeGeometryPointZ = t.SomeGeometryPointZ
	m.SomeHstore = t.SomeHstore
	m.SomeInet = t.SomeInet
	m.SomeInteger = t.SomeInteger
	m.SomeIntegerArray = t.SomeIntegerArray
	m.SomeInterval = t.SomeInterval
	m.SomeJSON = t.SomeJSON
	m.SomeJSONB = t.SomeJSONB
	m.SomeNumeric = t.SomeNumeric
	m.SomeNumericArray = t.SomeNumericArray
	m.SomePoint = t.SomePoint
	m.SomePolygon = t.SomePolygon
	m.SomeReal = t.SomeReal
	m.SomeRealArray = t.SomeRealArray
	m.SomeSmallint = t.SomeSmallint
	m.SomeSmallintArray = t.SomeSmallintArray
	m.SomeText = t.SomeText
	m.SomeTextArray = t.SomeTextArray
	m.SomeTimestamptz = t.SomeTimestamptz
	m.SomeTimestamp = t.SomeTimestamp
	m.SomeTsvector = t.SomeTsvector
	m.SomeUUID = t.SomeUUID

	return nil
}

func (m *NotNullFuzz) Insert(
	ctx context.Context,
	tx *sqlx.Tx,
	setPrimaryKey bool,
	setZeroValues bool,
	forceSetValuesForFields ...string,
) error {
	columns := make([]string, 0)
	values := make([]any, 0)

	if setPrimaryKey && (setZeroValues || !types.IsZeroUUID(m.ID)) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableIDColumn) || isRequired(NotNullFuzzTableColumnLookup, NotNullFuzzTableIDColumn) {
		columns = append(columns, NotNullFuzzTableIDColumn)

		v, err := types.FormatUUID(m.ID)
		if err != nil {
			return fmt.Errorf("failed to handle m.ID: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroInt(m.SomeBigint) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeBigintColumn) || isRequired(NotNullFuzzTableColumnLookup, NotNullFuzzTableSomeBigintColumn) {
		columns = append(columns, NotNullFuzzTableSomeBigintColumn)

		v, err := types.FormatInt(m.SomeBigint)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeBigint: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroIntArray(m.SomeBigintArray) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeBigintArrayColumn) || isRequired(NotNullFuzzTableColumnLookup, NotNullFuzzTableSomeBigintArrayColumn) {
		columns = append(columns, NotNullFuzzTableSomeBigintArrayColumn)

		v, err := types.FormatIntArray(m.SomeBigintArray)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeBigintArray: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroBool(m.SomeBoolean) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeBooleanColumn) || isRequired(NotNullFuzzTableColumnLookup, NotNullFuzzTableSomeBooleanColumn) {
		columns = append(columns, NotNullFuzzTableSomeBooleanColumn)

		v, err := types.FormatBool(m.SomeBoolean)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeBoolean: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroBoolArray(m.SomeBooleanArray) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeBooleanArrayColumn) || isRequired(NotNullFuzzTableColumnLookup, NotNullFuzzTableSomeBooleanArrayColumn) {
		columns = append(columns, NotNullFuzzTableSomeBooleanArrayColumn)

		v, err := types.FormatBoolArray(m.SomeBooleanArray)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeBooleanArray: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroBytes(m.SomeBytea) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeByteaColumn) || isRequired(NotNullFuzzTableColumnLookup, NotNullFuzzTableSomeByteaColumn) {
		columns = append(columns, NotNullFuzzTableSomeByteaColumn)

		v, err := types.FormatBytes(m.SomeBytea)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeBytea: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroString(m.SomeCharacterVarying) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeCharacterVaryingColumn) || isRequired(NotNullFuzzTableColumnLookup, NotNullFuzzTableSomeCharacterVaryingColumn) {
		columns = append(columns, NotNullFuzzTableSomeCharacterVaryingColumn)

		v, err := types.FormatString(m.SomeCharacterVarying)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeCharacterVarying: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroStringArray(m.SomeCharacterVaryingArray) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeCharacterVaryingArrayColumn) || isRequired(NotNullFuzzTableColumnLookup, NotNullFuzzTableSomeCharacterVaryingArrayColumn) {
		columns = append(columns, NotNullFuzzTableSomeCharacterVaryingArrayColumn)

		v, err := types.FormatStringArray(m.SomeCharacterVaryingArray)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeCharacterVaryingArray: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroFloat(m.SomeDoublePrecision) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeDoublePrecisionColumn) || isRequired(NotNullFuzzTableColumnLookup, NotNullFuzzTableSomeDoublePrecisionColumn) {
		columns = append(columns, NotNullFuzzTableSomeDoublePrecisionColumn)

		v, err := types.FormatFloat(m.SomeDoublePrecision)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeDoublePrecision: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroFloatArray(m.SomeDoublePrecisionArray) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeDoublePrecisionArrayColumn) || isRequired(NotNullFuzzTableColumnLookup, NotNullFuzzTableSomeDoublePrecisionArrayColumn) {
		columns = append(columns, NotNullFuzzTableSomeDoublePrecisionArrayColumn)

		v, err := types.FormatFloatArray(m.SomeDoublePrecisionArray)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeDoublePrecisionArray: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroFloat(m.SomeFloat) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeFloatColumn) || isRequired(NotNullFuzzTableColumnLookup, NotNullFuzzTableSomeFloatColumn) {
		columns = append(columns, NotNullFuzzTableSomeFloatColumn)

		v, err := types.FormatFloat(m.SomeFloat)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeFloat: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroFloatArray(m.SomeFloatArray) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeFloatArrayColumn) || isRequired(NotNullFuzzTableColumnLookup, NotNullFuzzTableSomeFloatArrayColumn) {
		columns = append(columns, NotNullFuzzTableSomeFloatArrayColumn)

		v, err := types.FormatFloatArray(m.SomeFloatArray)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeFloatArray: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroGeometry(m.SomeGeometryPointZ) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeGeometryPointZColumn) || isRequired(NotNullFuzzTableColumnLookup, NotNullFuzzTableSomeGeometryPointZColumn) {
		columns = append(columns, NotNullFuzzTableSomeGeometryPointZColumn)

		v, err := types.FormatGeometry(m.SomeGeometryPointZ)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeGeometryPointZ: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroHstore(m.SomeHstore) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeHstoreColumn) || isRequired(NotNullFuzzTableColumnLookup, NotNullFuzzTableSomeHstoreColumn) {
		columns = append(columns, NotNullFuzzTableSomeHstoreColumn)

		v, err := types.FormatHstore(m.SomeHstore)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeHstore: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroInet(m.SomeInet) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeInetColumn) || isRequired(NotNullFuzzTableColumnLookup, NotNullFuzzTableSomeInetColumn) {
		columns = append(columns, NotNullFuzzTableSomeInetColumn)

		v, err := types.FormatInet(m.SomeInet)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeInet: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroInt(m.SomeInteger) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeIntegerColumn) || isRequired(NotNullFuzzTableColumnLookup, NotNullFuzzTableSomeIntegerColumn) {
		columns = append(columns, NotNullFuzzTableSomeIntegerColumn)

		v, err := types.FormatInt(m.SomeInteger)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeInteger: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroIntArray(m.SomeIntegerArray) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeIntegerArrayColumn) || isRequired(NotNullFuzzTableColumnLookup, NotNullFuzzTableSomeIntegerArrayColumn) {
		columns = append(columns, NotNullFuzzTableSomeIntegerArrayColumn)

		v, err := types.FormatIntArray(m.SomeIntegerArray)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeIntegerArray: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroDuration(m.SomeInterval) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeIntervalColumn) || isRequired(NotNullFuzzTableColumnLookup, NotNullFuzzTableSomeIntervalColumn) {
		columns = append(columns, NotNullFuzzTableSomeIntervalColumn)

		v, err := types.FormatDuration(m.SomeInterval)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeInterval: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroJSON(m.SomeJSON) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeJSONColumn) || isRequired(NotNullFuzzTableColumnLookup, NotNullFuzzTableSomeJSONColumn) {
		columns = append(columns, NotNullFuzzTableSomeJSONColumn)

		v, err := types.FormatJSON(m.SomeJSON)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeJSON: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroJSON(m.SomeJSONB) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeJSONBColumn) || isRequired(NotNullFuzzTableColumnLookup, NotNullFuzzTableSomeJSONBColumn) {
		columns = append(columns, NotNullFuzzTableSomeJSONBColumn)

		v, err := types.FormatJSON(m.SomeJSONB)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeJSONB: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroFloat(m.SomeNumeric) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeNumericColumn) || isRequired(NotNullFuzzTableColumnLookup, NotNullFuzzTableSomeNumericColumn) {
		columns = append(columns, NotNullFuzzTableSomeNumericColumn)

		v, err := types.FormatFloat(m.SomeNumeric)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeNumeric: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroFloatArray(m.SomeNumericArray) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeNumericArrayColumn) || isRequired(NotNullFuzzTableColumnLookup, NotNullFuzzTableSomeNumericArrayColumn) {
		columns = append(columns, NotNullFuzzTableSomeNumericArrayColumn)

		v, err := types.FormatFloatArray(m.SomeNumericArray)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeNumericArray: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroPoint(m.SomePoint) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomePointColumn) || isRequired(NotNullFuzzTableColumnLookup, NotNullFuzzTableSomePointColumn) {
		columns = append(columns, NotNullFuzzTableSomePointColumn)

		v, err := types.FormatPoint(m.SomePoint)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomePoint: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroPolygon(m.SomePolygon) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomePolygonColumn) || isRequired(NotNullFuzzTableColumnLookup, NotNullFuzzTableSomePolygonColumn) {
		columns = append(columns, NotNullFuzzTableSomePolygonColumn)

		v, err := types.FormatPolygon(m.SomePolygon)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomePolygon: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroFloat(m.SomeReal) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeRealColumn) || isRequired(NotNullFuzzTableColumnLookup, NotNullFuzzTableSomeRealColumn) {
		columns = append(columns, NotNullFuzzTableSomeRealColumn)

		v, err := types.FormatFloat(m.SomeReal)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeReal: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroFloatArray(m.SomeRealArray) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeRealArrayColumn) || isRequired(NotNullFuzzTableColumnLookup, NotNullFuzzTableSomeRealArrayColumn) {
		columns = append(columns, NotNullFuzzTableSomeRealArrayColumn)

		v, err := types.FormatFloatArray(m.SomeRealArray)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeRealArray: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroInt(m.SomeSmallint) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeSmallintColumn) || isRequired(NotNullFuzzTableColumnLookup, NotNullFuzzTableSomeSmallintColumn) {
		columns = append(columns, NotNullFuzzTableSomeSmallintColumn)

		v, err := types.FormatInt(m.SomeSmallint)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeSmallint: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroIntArray(m.SomeSmallintArray) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeSmallintArrayColumn) || isRequired(NotNullFuzzTableColumnLookup, NotNullFuzzTableSomeSmallintArrayColumn) {
		columns = append(columns, NotNullFuzzTableSomeSmallintArrayColumn)

		v, err := types.FormatIntArray(m.SomeSmallintArray)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeSmallintArray: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroString(m.SomeText) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeTextColumn) || isRequired(NotNullFuzzTableColumnLookup, NotNullFuzzTableSomeTextColumn) {
		columns = append(columns, NotNullFuzzTableSomeTextColumn)

		v, err := types.FormatString(m.SomeText)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeText: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroStringArray(m.SomeTextArray) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeTextArrayColumn) || isRequired(NotNullFuzzTableColumnLookup, NotNullFuzzTableSomeTextArrayColumn) {
		columns = append(columns, NotNullFuzzTableSomeTextArrayColumn)

		v, err := types.FormatStringArray(m.SomeTextArray)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeTextArray: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.SomeTimestamptz) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeTimestamptzColumn) || isRequired(NotNullFuzzTableColumnLookup, NotNullFuzzTableSomeTimestamptzColumn) {
		columns = append(columns, NotNullFuzzTableSomeTimestamptzColumn)

		v, err := types.FormatTime(m.SomeTimestamptz)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeTimestamptz: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.SomeTimestamp) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeTimestampColumn) || isRequired(NotNullFuzzTableColumnLookup, NotNullFuzzTableSomeTimestampColumn) {
		columns = append(columns, NotNullFuzzTableSomeTimestampColumn)

		v, err := types.FormatTime(m.SomeTimestamp)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeTimestamp: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTSVector(m.SomeTsvector) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeTsvectorColumn) || isRequired(NotNullFuzzTableColumnLookup, NotNullFuzzTableSomeTsvectorColumn) {
		columns = append(columns, NotNullFuzzTableSomeTsvectorColumn)

		v, err := types.FormatTSVector(m.SomeTsvector)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeTsvector: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroUUID(m.SomeUUID) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeUUIDColumn) || isRequired(NotNullFuzzTableColumnLookup, NotNullFuzzTableSomeUUIDColumn) {
		columns = append(columns, NotNullFuzzTableSomeUUIDColumn)

		v, err := types.FormatUUID(m.SomeUUID)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeUUID: %v", err)
		}

		values = append(values, v)
	}

	ctx, cleanup := query.WithQueryID(ctx)
	defer cleanup()

	item, err := query.Insert(
		ctx,
		tx,
		NotNullFuzzTable,
		columns,
		nil,
		false,
		false,
		NotNullFuzzTableColumns,
		values...,
	)
	if err != nil {
		return fmt.Errorf("failed to insert %#+v: %v", m, err)
	}
	v := item[NotNullFuzzTableIDColumn]

	if v == nil {
		return fmt.Errorf("failed to find %v in %#+v", NotNullFuzzTableIDColumn, item)
	}

	wrapError := func(err error) error {
		return fmt.Errorf(
			"failed to treat %v: %#+v as uuid.UUID: %v",
			NotNullFuzzTableIDColumn,
			item[NotNullFuzzTableIDColumn],
			err,
		)
	}

	temp1, err := types.ParseUUID(v)
	if err != nil {
		return wrapError(err)
	}

	temp2, ok := temp1.(uuid.UUID)
	if !ok {
		return wrapError(fmt.Errorf("failed to cast to uuid.UUID"))
	}

	m.ID = temp2

	err = m.Reload(ctx, tx, slices.Contains(forceSetValuesForFields, "deleted_at"))
	if err != nil {
		return fmt.Errorf("failed to reload after insert")
	}

	return nil
}

func (m *NotNullFuzz) Update(
	ctx context.Context,
	tx *sqlx.Tx,
	setZeroValues bool,
	forceSetValuesForFields ...string,
) error {
	columns := make([]string, 0)
	values := make([]any, 0)

	if setZeroValues || !types.IsZeroInt(m.SomeBigint) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeBigintColumn) {
		columns = append(columns, NotNullFuzzTableSomeBigintColumn)

		v, err := types.FormatInt(m.SomeBigint)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeBigint: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroIntArray(m.SomeBigintArray) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeBigintArrayColumn) {
		columns = append(columns, NotNullFuzzTableSomeBigintArrayColumn)

		v, err := types.FormatIntArray(m.SomeBigintArray)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeBigintArray: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroBool(m.SomeBoolean) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeBooleanColumn) {
		columns = append(columns, NotNullFuzzTableSomeBooleanColumn)

		v, err := types.FormatBool(m.SomeBoolean)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeBoolean: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroBoolArray(m.SomeBooleanArray) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeBooleanArrayColumn) {
		columns = append(columns, NotNullFuzzTableSomeBooleanArrayColumn)

		v, err := types.FormatBoolArray(m.SomeBooleanArray)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeBooleanArray: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroBytes(m.SomeBytea) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeByteaColumn) {
		columns = append(columns, NotNullFuzzTableSomeByteaColumn)

		v, err := types.FormatBytes(m.SomeBytea)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeBytea: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroString(m.SomeCharacterVarying) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeCharacterVaryingColumn) {
		columns = append(columns, NotNullFuzzTableSomeCharacterVaryingColumn)

		v, err := types.FormatString(m.SomeCharacterVarying)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeCharacterVarying: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroStringArray(m.SomeCharacterVaryingArray) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeCharacterVaryingArrayColumn) {
		columns = append(columns, NotNullFuzzTableSomeCharacterVaryingArrayColumn)

		v, err := types.FormatStringArray(m.SomeCharacterVaryingArray)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeCharacterVaryingArray: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroFloat(m.SomeDoublePrecision) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeDoublePrecisionColumn) {
		columns = append(columns, NotNullFuzzTableSomeDoublePrecisionColumn)

		v, err := types.FormatFloat(m.SomeDoublePrecision)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeDoublePrecision: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroFloatArray(m.SomeDoublePrecisionArray) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeDoublePrecisionArrayColumn) {
		columns = append(columns, NotNullFuzzTableSomeDoublePrecisionArrayColumn)

		v, err := types.FormatFloatArray(m.SomeDoublePrecisionArray)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeDoublePrecisionArray: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroFloat(m.SomeFloat) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeFloatColumn) {
		columns = append(columns, NotNullFuzzTableSomeFloatColumn)

		v, err := types.FormatFloat(m.SomeFloat)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeFloat: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroFloatArray(m.SomeFloatArray) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeFloatArrayColumn) {
		columns = append(columns, NotNullFuzzTableSomeFloatArrayColumn)

		v, err := types.FormatFloatArray(m.SomeFloatArray)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeFloatArray: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroGeometry(m.SomeGeometryPointZ) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeGeometryPointZColumn) {
		columns = append(columns, NotNullFuzzTableSomeGeometryPointZColumn)

		v, err := types.FormatGeometry(m.SomeGeometryPointZ)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeGeometryPointZ: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroHstore(m.SomeHstore) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeHstoreColumn) {
		columns = append(columns, NotNullFuzzTableSomeHstoreColumn)

		v, err := types.FormatHstore(m.SomeHstore)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeHstore: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroInet(m.SomeInet) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeInetColumn) {
		columns = append(columns, NotNullFuzzTableSomeInetColumn)

		v, err := types.FormatInet(m.SomeInet)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeInet: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroInt(m.SomeInteger) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeIntegerColumn) {
		columns = append(columns, NotNullFuzzTableSomeIntegerColumn)

		v, err := types.FormatInt(m.SomeInteger)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeInteger: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroIntArray(m.SomeIntegerArray) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeIntegerArrayColumn) {
		columns = append(columns, NotNullFuzzTableSomeIntegerArrayColumn)

		v, err := types.FormatIntArray(m.SomeIntegerArray)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeIntegerArray: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroDuration(m.SomeInterval) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeIntervalColumn) {
		columns = append(columns, NotNullFuzzTableSomeIntervalColumn)

		v, err := types.FormatDuration(m.SomeInterval)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeInterval: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroJSON(m.SomeJSON) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeJSONColumn) {
		columns = append(columns, NotNullFuzzTableSomeJSONColumn)

		v, err := types.FormatJSON(m.SomeJSON)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeJSON: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroJSON(m.SomeJSONB) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeJSONBColumn) {
		columns = append(columns, NotNullFuzzTableSomeJSONBColumn)

		v, err := types.FormatJSON(m.SomeJSONB)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeJSONB: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroFloat(m.SomeNumeric) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeNumericColumn) {
		columns = append(columns, NotNullFuzzTableSomeNumericColumn)

		v, err := types.FormatFloat(m.SomeNumeric)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeNumeric: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroFloatArray(m.SomeNumericArray) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeNumericArrayColumn) {
		columns = append(columns, NotNullFuzzTableSomeNumericArrayColumn)

		v, err := types.FormatFloatArray(m.SomeNumericArray)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeNumericArray: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroPoint(m.SomePoint) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomePointColumn) {
		columns = append(columns, NotNullFuzzTableSomePointColumn)

		v, err := types.FormatPoint(m.SomePoint)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomePoint: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroPolygon(m.SomePolygon) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomePolygonColumn) {
		columns = append(columns, NotNullFuzzTableSomePolygonColumn)

		v, err := types.FormatPolygon(m.SomePolygon)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomePolygon: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroFloat(m.SomeReal) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeRealColumn) {
		columns = append(columns, NotNullFuzzTableSomeRealColumn)

		v, err := types.FormatFloat(m.SomeReal)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeReal: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroFloatArray(m.SomeRealArray) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeRealArrayColumn) {
		columns = append(columns, NotNullFuzzTableSomeRealArrayColumn)

		v, err := types.FormatFloatArray(m.SomeRealArray)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeRealArray: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroInt(m.SomeSmallint) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeSmallintColumn) {
		columns = append(columns, NotNullFuzzTableSomeSmallintColumn)

		v, err := types.FormatInt(m.SomeSmallint)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeSmallint: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroIntArray(m.SomeSmallintArray) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeSmallintArrayColumn) {
		columns = append(columns, NotNullFuzzTableSomeSmallintArrayColumn)

		v, err := types.FormatIntArray(m.SomeSmallintArray)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeSmallintArray: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroString(m.SomeText) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeTextColumn) {
		columns = append(columns, NotNullFuzzTableSomeTextColumn)

		v, err := types.FormatString(m.SomeText)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeText: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroStringArray(m.SomeTextArray) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeTextArrayColumn) {
		columns = append(columns, NotNullFuzzTableSomeTextArrayColumn)

		v, err := types.FormatStringArray(m.SomeTextArray)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeTextArray: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.SomeTimestamptz) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeTimestamptzColumn) {
		columns = append(columns, NotNullFuzzTableSomeTimestamptzColumn)

		v, err := types.FormatTime(m.SomeTimestamptz)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeTimestamptz: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.SomeTimestamp) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeTimestampColumn) {
		columns = append(columns, NotNullFuzzTableSomeTimestampColumn)

		v, err := types.FormatTime(m.SomeTimestamp)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeTimestamp: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTSVector(m.SomeTsvector) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeTsvectorColumn) {
		columns = append(columns, NotNullFuzzTableSomeTsvectorColumn)

		v, err := types.FormatTSVector(m.SomeTsvector)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeTsvector: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroUUID(m.SomeUUID) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeUUIDColumn) {
		columns = append(columns, NotNullFuzzTableSomeUUIDColumn)

		v, err := types.FormatUUID(m.SomeUUID)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeUUID: %v", err)
		}

		values = append(values, v)
	}

	v, err := types.FormatUUID(m.ID)
	if err != nil {
		return fmt.Errorf("failed to handle m.ID: %v", err)
	}

	values = append(values, v)

	ctx, cleanup := query.WithQueryID(ctx)
	defer cleanup()

	_, err = query.Update(
		ctx,
		tx,
		NotNullFuzzTable,
		columns,
		fmt.Sprintf("%v = $$??", NotNullFuzzTableIDColumn),
		NotNullFuzzTableColumns,
		values...,
	)
	if err != nil {
		return fmt.Errorf("failed to update %#+v: %v", m, err)
	}

	err = m.Reload(ctx, tx, slices.Contains(forceSetValuesForFields, "deleted_at"))
	if err != nil {
		return fmt.Errorf("failed to reload after update")
	}

	return nil
}

func (m *NotNullFuzz) Delete(
	ctx context.Context,
	tx *sqlx.Tx,
	hardDeletes ...bool,
) error {
	/* soft-delete not applicable */

	values := make([]any, 0)
	v, err := types.FormatUUID(m.ID)
	if err != nil {
		return fmt.Errorf("failed to handle m.ID: %v", err)
	}

	values = append(values, v)

	ctx, cleanup := query.WithQueryID(ctx)
	defer cleanup()

	err = query.Delete(
		ctx,
		tx,
		NotNullFuzzTable,
		fmt.Sprintf("%v = $$??", NotNullFuzzTableIDColumn),
		values...,
	)
	if err != nil {
		return fmt.Errorf("failed to delete %#+v: %v", m, err)
	}

	_ = m.Reload(ctx, tx, true)

	return nil
}

func SelectNotNullFuzzes(
	ctx context.Context,
	tx *sqlx.Tx,
	where string,
	orderBy *string,
	limit *int,
	offset *int,
	values ...any,
) ([]*NotNullFuzz, error) {
	if slices.Contains(NotNullFuzzTableColumns, "deleted_at") {
		if !strings.Contains(where, "deleted_at") {
			if where != "" {
				where += "\n    AND "
			}

			where += "deleted_at IS null"
		}
	}

	ctx, cleanup := query.WithQueryID(ctx)
	defer cleanup()

	items, err := query.Select(
		ctx,
		tx,
		NotNullFuzzTableColumnsWithTypeCasts,
		NotNullFuzzTable,
		where,
		orderBy,
		limit,
		offset,
		values...,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to call SelectNotNullFuzzs; err: %v", err)
	}

	objects := make([]*NotNullFuzz, 0)

	for _, item := range items {
		object := &NotNullFuzz{}

		err = object.FromItem(item)
		if err != nil {
			return nil, err
		}

		objects = append(objects, object)
	}

	return objects, nil
}

func SelectNotNullFuzz(
	ctx context.Context,
	tx *sqlx.Tx,
	where string,
	values ...any,
) (*NotNullFuzz, error) {
	ctx, cleanup := query.WithQueryID(ctx)
	defer cleanup()

	objects, err := SelectNotNullFuzzes(
		ctx,
		tx,
		where,
		nil,
		helpers.Ptr(2),
		helpers.Ptr(0),
		values...,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to call SelectNotNullFuzz; err: %v", err)
	}

	if len(objects) > 1 {
		return nil, fmt.Errorf("attempt to call SelectNotNullFuzz returned more than 1 row")
	}

	if len(objects) < 1 {
		return nil, sql.ErrNoRows
	}

	object := objects[0]

	return object, nil
}

func handleGetNotNullFuzzes(w http.ResponseWriter, r *http.Request, db *sqlx.DB, redisPool *redis.Pool, objectMiddlewares []server.ObjectMiddleware) {
	ctx := r.Context()

	insaneOrderParams := make([]string, 0)
	hadInsaneOrderParams := false

	unrecognizedParams := make([]string, 0)
	hadUnrecognizedParams := false

	unparseableParams := make([]string, 0)
	hadUnparseableParams := false

	var orderByDirection *string
	orderBys := make([]string, 0)

	includes := make([]string, 0)

	values := make([]any, 0)
	wheres := make([]string, 0)
	for rawKey, rawValues := range r.URL.Query() {
		if rawKey == "limit" || rawKey == "offset" {
			continue
		}

		parts := strings.Split(rawKey, "__")
		isUnrecognized := len(parts) != 2

		comparison := ""
		isSliceComparison := false
		isNullComparison := false
		IsLikeComparison := false

		if !isUnrecognized {
			column := NotNullFuzzTableColumnLookup[parts[0]]
			if column == nil {
				if parts[0] != "load" {
					isUnrecognized = true
				}
			} else {
				switch parts[1] {
				case "eq":
					comparison = "="
				case "ne":
					comparison = "!="
				case "gt":
					comparison = ">"
				case "gte":
					comparison = ">="
				case "lt":
					comparison = "<"
				case "lte":
					comparison = "<="
				case "in":
					comparison = "IN"
					isSliceComparison = true
				case "nin", "notin":
					comparison = "NOT IN"
					isSliceComparison = true
				case "isnull":
					comparison = "IS NULL"
					isNullComparison = true
				case "nisnull", "isnotnull":
					comparison = "IS NOT NULL"
					isNullComparison = true
				case "l", "like":
					comparison = "LIKE"
					IsLikeComparison = true
				case "nl", "nlike", "notlike":
					comparison = "NOT LIKE"
					IsLikeComparison = true
				case "il", "ilike":
					comparison = "ILIKE"
					IsLikeComparison = true
				case "nil", "nilike", "notilike":
					comparison = "NOT ILIKE"
					IsLikeComparison = true
				case "desc":
					if orderByDirection != nil && *orderByDirection != "DESC" {
						hadInsaneOrderParams = true
						insaneOrderParams = append(insaneOrderParams, rawKey)
						continue
					}

					orderByDirection = helpers.Ptr("DESC")
					orderBys = append(orderBys, parts[0])
					continue
				case "asc":
					if orderByDirection != nil && *orderByDirection != "ASC" {
						hadInsaneOrderParams = true
						insaneOrderParams = append(insaneOrderParams, rawKey)
						continue
					}

					orderByDirection = helpers.Ptr("ASC")
					orderBys = append(orderBys, parts[0])
					continue
				case "load":
					includes = append(includes, parts[0])
					_ = includes

					continue
				default:
					isUnrecognized = true
				}
			}
		}

		if isNullComparison {
			wheres = append(wheres, fmt.Sprintf("%s %s", parts[0], comparison))
			continue
		}

		for _, rawValue := range rawValues {
			if isUnrecognized {
				unrecognizedParams = append(unrecognizedParams, fmt.Sprintf("%s=%s", rawKey, rawValue))
				hadUnrecognizedParams = true
				continue
			}

			if hadUnrecognizedParams {
				continue
			}

			attempts := make([]string, 0)

			if !IsLikeComparison {
				attempts = append(attempts, rawValue)
			}

			if isSliceComparison {
				attempts = append(attempts, fmt.Sprintf("[%s]", rawValue))

				vs := make([]string, 0)
				for _, v := range strings.Split(rawValue, ",") {
					vs = append(vs, fmt.Sprintf("\"%s\"", v))
				}

				attempts = append(attempts, fmt.Sprintf("[%s]", strings.Join(vs, ",")))
			}

			if IsLikeComparison {
				attempts = append(attempts, fmt.Sprintf("\"%%%s%%\"", rawValue))
			} else {
				attempts = append(attempts, fmt.Sprintf("\"%s\"", rawValue))
			}

			var err error

			for _, attempt := range attempts {
				var value any
				err = json.Unmarshal([]byte(attempt), &value)
				if err == nil {
					if isSliceComparison {
						sliceValues, ok := value.([]any)
						if !ok {
							err = fmt.Errorf("failed to cast %#+v to []string", value)
							break
						}

						values = append(values, sliceValues...)

						sliceWheres := make([]string, 0)
						for range values {
							sliceWheres = append(sliceWheres, "$$??")
						}

						wheres = append(wheres, fmt.Sprintf("%s %s (%s)", parts[0], comparison, strings.Join(sliceWheres, ", ")))
					} else {
						values = append(values, value)
						wheres = append(wheres, fmt.Sprintf("%s %s $$??", parts[0], comparison))
					}

					break
				}
			}

			if err != nil {
				unparseableParams = append(unparseableParams, fmt.Sprintf("%s=%s", rawKey, rawValue))
				hadUnparseableParams = true
				continue
			}
		}
	}

	if hadUnrecognizedParams {
		helpers.HandleErrorResponse(
			w,
			http.StatusInternalServerError,
			fmt.Errorf("unrecognized params %s", strings.Join(unrecognizedParams, ", ")),
		)
		return
	}

	if hadUnparseableParams {
		helpers.HandleErrorResponse(
			w,
			http.StatusInternalServerError,
			fmt.Errorf("unparseable params %s", strings.Join(unparseableParams, ", ")),
		)
		return
	}

	if hadInsaneOrderParams {
		helpers.HandleErrorResponse(
			w,
			http.StatusInternalServerError,
			fmt.Errorf("insane order params (e.g. conflicting asc / desc) %s", strings.Join(insaneOrderParams, ", ")),
		)
		return
	}

	limit := 2000
	rawLimit := r.URL.Query().Get("limit")
	if rawLimit != "" {
		possibleLimit, err := strconv.ParseInt(rawLimit, 10, 64)
		if err != nil {
			helpers.HandleErrorResponse(
				w,
				http.StatusInternalServerError,
				fmt.Errorf("failed to parse param limit=%s as int: %v", rawLimit, err),
			)
			return
		}

		limit = int(possibleLimit)
	}

	offset := 0
	rawOffset := r.URL.Query().Get("offset")
	if rawOffset != "" {
		possibleOffset, err := strconv.ParseInt(rawOffset, 10, 64)
		if err != nil {
			helpers.HandleErrorResponse(
				w,
				http.StatusInternalServerError,
				fmt.Errorf("failed to parse param offset=%s as int: %v", rawOffset, err),
			)
			return
		}

		offset = int(possibleOffset)
	}

	hashableOrderBy := ""
	var orderBy *string
	if len(orderBys) > 0 {
		hashableOrderBy = strings.Join(orderBys, ", ")
		if len(orderBys) > 1 {
			hashableOrderBy = fmt.Sprintf("(%v)", hashableOrderBy)
		}
		hashableOrderBy = fmt.Sprintf("%v %v", hashableOrderBy, *orderByDirection)
		orderBy = &hashableOrderBy
	}

	requestHash, err := helpers.GetRequestHash(NotNullFuzzTable, wheres, hashableOrderBy, limit, offset, values, nil)
	if err != nil {
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	redisConn := redisPool.Get()
	defer func() {
		_ = redisConn.Close()
	}()

	cacheHit, err := helpers.AttemptCachedResponse(requestHash, redisConn, w)
	if err != nil {
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	if cacheHit {
		return
	}

	tx, err := db.BeginTxx(r.Context(), nil)
	if err != nil {
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	defer func() {
		_ = tx.Rollback()
	}()

	where := strings.Join(wheres, "\n    AND ")

	objects, err := SelectNotNullFuzzes(ctx, tx, where, orderBy, &limit, &offset, values...)
	if err != nil {
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	err = tx.Commit()
	if err != nil {
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	returnedObjectsAsJSON := helpers.HandleObjectsResponse(w, http.StatusOK, objects)

	err = helpers.StoreCachedResponse(requestHash, redisConn, string(returnedObjectsAsJSON))
	if err != nil {
		log.Printf("warning: %v", err)
	}
}

func handleGetNotNullFuzz(w http.ResponseWriter, r *http.Request, db *sqlx.DB, redisPool *redis.Pool, objectMiddlewares []server.ObjectMiddleware, primaryKey string) {
	ctx := r.Context()

	wheres := []string{fmt.Sprintf("%s = $$??", NotNullFuzzTablePrimaryKeyColumn)}
	values := []any{primaryKey}

	requestHash, err := helpers.GetRequestHash(NotNullFuzzTable, wheres, "", 2, 0, values, primaryKey)
	if err != nil {
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	redisConn := redisPool.Get()
	defer func() {
		_ = redisConn.Close()
	}()

	cacheHit, err := helpers.AttemptCachedResponse(requestHash, redisConn, w)
	if err != nil {
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	if cacheHit {
		return
	}

	tx, err := db.BeginTxx(r.Context(), nil)
	if err != nil {
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	defer func() {
		_ = tx.Rollback()
	}()

	where := strings.Join(wheres, "\n    AND ")

	object, err := SelectNotNullFuzz(ctx, tx, where, values...)
	if err != nil {
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	err = tx.Commit()
	if err != nil {
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	returnedObjectsAsJSON := helpers.HandleObjectsResponse(w, http.StatusOK, []*NotNullFuzz{object})

	err = helpers.StoreCachedResponse(requestHash, redisConn, string(returnedObjectsAsJSON))
	if err != nil {
		log.Printf("warning: %v", err)
	}
}

func handlePostNotNullFuzzs(w http.ResponseWriter, r *http.Request, db *sqlx.DB, redisPool *redis.Pool, objectMiddlewares []server.ObjectMiddleware, waitForChange server.WaitForChange) {
	_ = redisPool

	b, err := io.ReadAll(r.Body)
	if err != nil {
		err = fmt.Errorf("failed to read body of HTTP request: %v", err)
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	var allItems []map[string]any
	err = json.Unmarshal(b, &allItems)
	if err != nil {
		err = fmt.Errorf("failed to unmarshal %#+v as JSON list of objects: %v", string(b), err)
		helpers.HandleErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	forceSetValuesForFieldsByObjectIndex := make([][]string, 0)
	objects := make([]*NotNullFuzz, 0)
	for _, item := range allItems {
		forceSetValuesForFields := make([]string, 0)
		for _, possibleField := range maps.Keys(item) {
			if !slices.Contains(NotNullFuzzTableColumns, possibleField) {
				continue
			}

			forceSetValuesForFields = append(forceSetValuesForFields, possibleField)
		}
		forceSetValuesForFieldsByObjectIndex = append(forceSetValuesForFieldsByObjectIndex, forceSetValuesForFields)

		object := &NotNullFuzz{}
		err = object.FromItem(item)
		if err != nil {
			err = fmt.Errorf("failed to interpret %#+v as NotNullFuzz in item form: %v", item, err)
			helpers.HandleErrorResponse(w, http.StatusBadRequest, err)
			return
		}

		objects = append(objects, object)
	}

	tx, err := db.BeginTxx(r.Context(), nil)
	if err != nil {
		err = fmt.Errorf("failed to begin DB transaction: %v", err)
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	defer func() {
		_ = tx.Rollback()
	}()

	xid, err := query.GetXid(r.Context(), tx)
	if err != nil {
		err = fmt.Errorf("failed to get xid: %v", err)
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}
	_ = xid

	for i, object := range objects {
		err = object.Insert(r.Context(), tx, false, false, forceSetValuesForFieldsByObjectIndex[i]...)
		if err != nil {
			err = fmt.Errorf("failed to insert %#+v: %v", object, err)
			helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
			return
		}

		objects[i] = object
	}

	errs := make(chan error, 1)
	go func() {
		_, err = waitForChange(r.Context(), []stream.Action{stream.INSERT}, NotNullFuzzTable, xid)
		if err != nil {
			err = fmt.Errorf("failed to wait for change: %v", err)
			errs <- err
			return
		}

		errs <- nil
	}()

	err = tx.Commit()
	if err != nil {
		err = fmt.Errorf("failed to commit DB transaction: %v", err)
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	select {
	case <-r.Context().Done():
		err = fmt.Errorf("context canceled")
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	case err = <-errs:
		if err != nil {
			helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
			return
		}
	}

	helpers.HandleObjectsResponse(w, http.StatusCreated, objects)
}

func handlePutNotNullFuzz(w http.ResponseWriter, r *http.Request, db *sqlx.DB, redisPool *redis.Pool, objectMiddlewares []server.ObjectMiddleware, waitForChange server.WaitForChange, primaryKey string) {
	_ = redisPool

	b, err := io.ReadAll(r.Body)
	if err != nil {
		err = fmt.Errorf("failed to read body of HTTP request: %v", err)
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	var item map[string]any
	err = json.Unmarshal(b, &item)
	if err != nil {
		err = fmt.Errorf("failed to unmarshal %#+v as JSON object: %v", string(b), err)
		helpers.HandleErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	item[NotNullFuzzTablePrimaryKeyColumn] = primaryKey

	object := &NotNullFuzz{}
	err = object.FromItem(item)
	if err != nil {
		err = fmt.Errorf("failed to interpret %#+v as NotNullFuzz in item form: %v", item, err)
		helpers.HandleErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	tx, err := db.BeginTxx(r.Context(), nil)
	if err != nil {
		err = fmt.Errorf("failed to begin DB transaction: %v", err)
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	defer func() {
		_ = tx.Rollback()
	}()

	xid, err := query.GetXid(r.Context(), tx)
	if err != nil {
		err = fmt.Errorf("failed to get xid: %v", err)
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}
	_ = xid

	err = object.Update(r.Context(), tx, true)
	if err != nil {
		err = fmt.Errorf("failed to update %#+v: %v", object, err)
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	errs := make(chan error, 1)
	go func() {
		_, err = waitForChange(r.Context(), []stream.Action{stream.UPDATE, stream.SOFT_DELETE, stream.SOFT_RESTORE, stream.SOFT_UPDATE}, NotNullFuzzTable, xid)
		if err != nil {
			err = fmt.Errorf("failed to wait for change: %v", err)
			errs <- err
			return
		}

		errs <- nil
	}()

	err = tx.Commit()
	if err != nil {
		err = fmt.Errorf("failed to commit DB transaction: %v", err)
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	select {
	case <-r.Context().Done():
		err = fmt.Errorf("context canceled")
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	case err = <-errs:
		if err != nil {
			helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
			return
		}
	}

	helpers.HandleObjectsResponse(w, http.StatusOK, []*NotNullFuzz{object})
}

func handlePatchNotNullFuzz(w http.ResponseWriter, r *http.Request, db *sqlx.DB, redisPool *redis.Pool, objectMiddlewares []server.ObjectMiddleware, waitForChange server.WaitForChange, primaryKey string) {
	_ = redisPool

	b, err := io.ReadAll(r.Body)
	if err != nil {
		err = fmt.Errorf("failed to read body of HTTP request: %v", err)
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	var item map[string]any
	err = json.Unmarshal(b, &item)
	if err != nil {
		err = fmt.Errorf("failed to unmarshal %#+v as JSON object: %v", string(b), err)
		helpers.HandleErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	forceSetValuesForFields := make([]string, 0)
	for _, possibleField := range maps.Keys(item) {
		if !slices.Contains(NotNullFuzzTableColumns, possibleField) {
			continue
		}

		forceSetValuesForFields = append(forceSetValuesForFields, possibleField)
	}

	item[NotNullFuzzTablePrimaryKeyColumn] = primaryKey

	object := &NotNullFuzz{}
	err = object.FromItem(item)
	if err != nil {
		err = fmt.Errorf("failed to interpret %#+v as NotNullFuzz in item form: %v", item, err)
		helpers.HandleErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	tx, err := db.BeginTxx(r.Context(), nil)
	if err != nil {
		err = fmt.Errorf("failed to begin DB transaction: %v", err)
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	defer func() {
		_ = tx.Rollback()
	}()

	xid, err := query.GetXid(r.Context(), tx)
	if err != nil {
		err = fmt.Errorf("failed to get xid: %v", err)
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}
	_ = xid

	err = object.Update(r.Context(), tx, false, forceSetValuesForFields...)
	if err != nil {
		err = fmt.Errorf("failed to update %#+v: %v", object, err)
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	errs := make(chan error, 1)
	go func() {
		_, err = waitForChange(r.Context(), []stream.Action{stream.UPDATE, stream.SOFT_DELETE, stream.SOFT_RESTORE, stream.SOFT_UPDATE}, NotNullFuzzTable, xid)
		if err != nil {
			err = fmt.Errorf("failed to wait for change: %v", err)
			errs <- err
			return
		}

		errs <- nil
	}()

	err = tx.Commit()
	if err != nil {
		err = fmt.Errorf("failed to commit DB transaction: %v", err)
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	select {
	case <-r.Context().Done():
		err = fmt.Errorf("context canceled")
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	case err = <-errs:
		if err != nil {
			helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
			return
		}
	}

	helpers.HandleObjectsResponse(w, http.StatusOK, []*NotNullFuzz{object})
}

func handleDeleteNotNullFuzz(w http.ResponseWriter, r *http.Request, db *sqlx.DB, redisPool *redis.Pool, objectMiddlewares []server.ObjectMiddleware, waitForChange server.WaitForChange, primaryKey string) {
	_ = redisPool

	var item = make(map[string]any)

	item[NotNullFuzzTablePrimaryKeyColumn] = primaryKey

	object := &NotNullFuzz{}
	err := object.FromItem(item)
	if err != nil {
		err = fmt.Errorf("failed to interpret %#+v as NotNullFuzz in item form: %v", item, err)
		helpers.HandleErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	tx, err := db.BeginTxx(r.Context(), nil)
	if err != nil {
		err = fmt.Errorf("failed to begin DB transaction: %v", err)
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	defer func() {
		_ = tx.Rollback()
	}()

	xid, err := query.GetXid(r.Context(), tx)
	if err != nil {
		err = fmt.Errorf("failed to get xid: %v", err)
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}
	_ = xid

	err = object.Delete(r.Context(), tx)
	if err != nil {
		err = fmt.Errorf("failed to delete %#+v: %v", object, err)
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	errs := make(chan error, 1)
	go func() {
		_, err = waitForChange(r.Context(), []stream.Action{stream.DELETE, stream.SOFT_DELETE}, NotNullFuzzTable, xid)
		if err != nil {
			err = fmt.Errorf("failed to wait for change: %v", err)
			errs <- err
			return
		}

		errs <- nil
	}()

	err = tx.Commit()
	if err != nil {
		err = fmt.Errorf("failed to commit DB transaction: %v", err)
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	select {
	case <-r.Context().Done():
		err = fmt.Errorf("context canceled")
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	case err = <-errs:
		if err != nil {
			helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
			return
		}
	}

	helpers.HandleObjectsResponse(w, http.StatusNoContent, nil)
}

func GetNotNullFuzzRouter(db *sqlx.DB, redisPool *redis.Pool, httpMiddlewares []server.HTTPMiddleware, objectMiddlewares []server.ObjectMiddleware, waitForChange server.WaitForChange) chi.Router {
	r := chi.NewRouter()

	for _, m := range httpMiddlewares {
		r.Use(m)
	}

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		handleGetNotNullFuzzes(w, r, db, redisPool, objectMiddlewares)
	})

	r.Get("/{primaryKey}", func(w http.ResponseWriter, r *http.Request) {
		handleGetNotNullFuzz(w, r, db, redisPool, objectMiddlewares, chi.URLParam(r, "primaryKey"))
	})

	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		handlePostNotNullFuzzs(w, r, db, redisPool, objectMiddlewares, waitForChange)
	})

	r.Put("/{primaryKey}", func(w http.ResponseWriter, r *http.Request) {
		handlePutNotNullFuzz(w, r, db, redisPool, objectMiddlewares, waitForChange, chi.URLParam(r, "primaryKey"))
	})

	r.Patch("/{primaryKey}", func(w http.ResponseWriter, r *http.Request) {
		handlePatchNotNullFuzz(w, r, db, redisPool, objectMiddlewares, waitForChange, chi.URLParam(r, "primaryKey"))
	})

	r.Delete("/{primaryKey}", func(w http.ResponseWriter, r *http.Request) {
		handleDeleteNotNullFuzz(w, r, db, redisPool, objectMiddlewares, waitForChange, chi.URLParam(r, "primaryKey"))
	})

	return r
}

func NewNotNullFuzzFromItem(item map[string]any) (any, error) {
	object := &NotNullFuzz{}

	err := object.FromItem(item)
	if err != nil {
		return nil, err
	}

	return object, nil
}

func init() {
	register(
		NotNullFuzzTable,
		NotNullFuzz{},
		NewNotNullFuzzFromItem,
		"/not-null-fuzzes",
		GetNotNullFuzzRouter,
	)
}
