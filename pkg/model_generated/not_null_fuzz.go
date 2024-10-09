package model_generated

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"net/netip"
	"slices"
	"strings"
	"time"

	"github.com/cridenour/go-postgis"
	"github.com/go-chi/chi/v5"
	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
	"github.com/initialed85/djangolang/pkg/config"
	"github.com/initialed85/djangolang/pkg/helpers"
	"github.com/initialed85/djangolang/pkg/introspect"
	"github.com/initialed85/djangolang/pkg/query"
	"github.com/initialed85/djangolang/pkg/server"
	"github.com/initialed85/djangolang/pkg/stream"
	"github.com/initialed85/djangolang/pkg/types"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/exp/maps"
)

type NotNullFuzz struct {
	MrPrimary                                      int64              `json:"mr_primary"`
	SomeBigint                                     int64              `json:"some_bigint"`
	SomeBigintArray                                []int64            `json:"some_bigint_array"`
	SomeBoolean                                    bool               `json:"some_boolean"`
	SomeBooleanArray                               []bool             `json:"some_boolean_array"`
	SomeBytea                                      []byte             `json:"some_bytea"`
	SomeCharacterVarying                           string             `json:"some_character_varying"`
	SomeCharacterVaryingArray                      []string           `json:"some_character_varying_array"`
	SomeDoublePrecision                            float64            `json:"some_double_precision"`
	SomeDoublePrecisionArray                       []float64          `json:"some_double_precision_array"`
	SomeFloat                                      float64            `json:"some_float"`
	SomeFloatArray                                 []float64          `json:"some_float_array"`
	SomeGeometryPointZ                             postgis.PointZ     `json:"some_geometry_point_z"`
	SomeHstore                                     map[string]*string `json:"some_hstore"`
	SomeInet                                       netip.Prefix       `json:"some_inet"`
	SomeInteger                                    int64              `json:"some_integer"`
	SomeIntegerArray                               []int64            `json:"some_integer_array"`
	SomeInterval                                   time.Duration      `json:"some_interval"`
	SomeJSON                                       any                `json:"some_json"`
	SomeJSONB                                      any                `json:"some_jsonb"`
	SomeNumeric                                    float64            `json:"some_numeric"`
	SomeNumericArray                               []float64          `json:"some_numeric_array"`
	SomePoint                                      pgtype.Vec2        `json:"some_point"`
	SomePolygon                                    []pgtype.Vec2      `json:"some_polygon"`
	SomeReal                                       float64            `json:"some_real"`
	SomeRealArray                                  []float64          `json:"some_real_array"`
	SomeSmallint                                   int64              `json:"some_smallint"`
	SomeSmallintArray                              []int64            `json:"some_smallint_array"`
	SomeText                                       string             `json:"some_text"`
	SomeTextArray                                  []string           `json:"some_text_array"`
	SomeTimestamptz                                time.Time          `json:"some_timestamptz"`
	SomeTimestamp                                  time.Time          `json:"some_timestamp"`
	SomeTsvector                                   map[string][]int   `json:"some_tsvector"`
	SomeUUID                                       uuid.UUID          `json:"some_uuid"`
	OtherNotNullFuzz                               *int64             `json:"other_not_null_fuzz"`
	OtherNotNullFuzzObject                         *NotNullFuzz       `json:"other_not_null_fuzz_object"`
	ReferencedByNotNullFuzzOtherNotNullFuzzObjects []*NotNullFuzz     `json:"referenced_by_not_null_fuzz_other_not_null_fuzz_objects"`
}

var NotNullFuzzTable = "not_null_fuzz"

var NotNullFuzzTableNamespaceID int32 = 1337 + 5

var (
	NotNullFuzzTableMrPrimaryColumn                 = "mr_primary"
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
	NotNullFuzzTableOtherNotNullFuzzColumn          = "other_not_null_fuzz"
)

var (
	NotNullFuzzTableMrPrimaryColumnWithTypeCast                 = `"mr_primary" AS mr_primary`
	NotNullFuzzTableSomeBigintColumnWithTypeCast                = `"some_bigint" AS some_bigint`
	NotNullFuzzTableSomeBigintArrayColumnWithTypeCast           = `"some_bigint_array" AS some_bigint_array`
	NotNullFuzzTableSomeBooleanColumnWithTypeCast               = `"some_boolean" AS some_boolean`
	NotNullFuzzTableSomeBooleanArrayColumnWithTypeCast          = `"some_boolean_array" AS some_boolean_array`
	NotNullFuzzTableSomeByteaColumnWithTypeCast                 = `"some_bytea" AS some_bytea`
	NotNullFuzzTableSomeCharacterVaryingColumnWithTypeCast      = `"some_character_varying" AS some_character_varying`
	NotNullFuzzTableSomeCharacterVaryingArrayColumnWithTypeCast = `"some_character_varying_array" AS some_character_varying_array`
	NotNullFuzzTableSomeDoublePrecisionColumnWithTypeCast       = `"some_double_precision" AS some_double_precision`
	NotNullFuzzTableSomeDoublePrecisionArrayColumnWithTypeCast  = `"some_double_precision_array" AS some_double_precision_array`
	NotNullFuzzTableSomeFloatColumnWithTypeCast                 = `"some_float" AS some_float`
	NotNullFuzzTableSomeFloatArrayColumnWithTypeCast            = `"some_float_array" AS some_float_array`
	NotNullFuzzTableSomeGeometryPointZColumnWithTypeCast        = `"some_geometry_point_z" AS some_geometry_point_z`
	NotNullFuzzTableSomeHstoreColumnWithTypeCast                = `"some_hstore" AS some_hstore`
	NotNullFuzzTableSomeInetColumnWithTypeCast                  = `"some_inet" AS some_inet`
	NotNullFuzzTableSomeIntegerColumnWithTypeCast               = `"some_integer" AS some_integer`
	NotNullFuzzTableSomeIntegerArrayColumnWithTypeCast          = `"some_integer_array" AS some_integer_array`
	NotNullFuzzTableSomeIntervalColumnWithTypeCast              = `"some_interval" AS some_interval`
	NotNullFuzzTableSomeJSONColumnWithTypeCast                  = `"some_json" AS some_json`
	NotNullFuzzTableSomeJSONBColumnWithTypeCast                 = `"some_jsonb" AS some_jsonb`
	NotNullFuzzTableSomeNumericColumnWithTypeCast               = `"some_numeric" AS some_numeric`
	NotNullFuzzTableSomeNumericArrayColumnWithTypeCast          = `"some_numeric_array" AS some_numeric_array`
	NotNullFuzzTableSomePointColumnWithTypeCast                 = `"some_point" AS some_point`
	NotNullFuzzTableSomePolygonColumnWithTypeCast               = `"some_polygon" AS some_polygon`
	NotNullFuzzTableSomeRealColumnWithTypeCast                  = `"some_real" AS some_real`
	NotNullFuzzTableSomeRealArrayColumnWithTypeCast             = `"some_real_array" AS some_real_array`
	NotNullFuzzTableSomeSmallintColumnWithTypeCast              = `"some_smallint" AS some_smallint`
	NotNullFuzzTableSomeSmallintArrayColumnWithTypeCast         = `"some_smallint_array" AS some_smallint_array`
	NotNullFuzzTableSomeTextColumnWithTypeCast                  = `"some_text" AS some_text`
	NotNullFuzzTableSomeTextArrayColumnWithTypeCast             = `"some_text_array" AS some_text_array`
	NotNullFuzzTableSomeTimestamptzColumnWithTypeCast           = `"some_timestamptz" AS some_timestamptz`
	NotNullFuzzTableSomeTimestampColumnWithTypeCast             = `"some_timestamp" AS some_timestamp`
	NotNullFuzzTableSomeTsvectorColumnWithTypeCast              = `"some_tsvector" AS some_tsvector`
	NotNullFuzzTableSomeUUIDColumnWithTypeCast                  = `"some_uuid" AS some_uuid`
	NotNullFuzzTableOtherNotNullFuzzColumnWithTypeCast          = `"other_not_null_fuzz" AS other_not_null_fuzz`
)

var NotNullFuzzTableColumns = []string{
	NotNullFuzzTableMrPrimaryColumn,
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
	NotNullFuzzTableOtherNotNullFuzzColumn,
}

var NotNullFuzzTableColumnsWithTypeCasts = []string{
	NotNullFuzzTableMrPrimaryColumnWithTypeCast,
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
	NotNullFuzzTableOtherNotNullFuzzColumnWithTypeCast,
}

var NotNullFuzzIntrospectedTable *introspect.Table

var NotNullFuzzTableColumnLookup map[string]*introspect.Column

var (
	NotNullFuzzTablePrimaryKeyColumn = NotNullFuzzTableMrPrimaryColumn
)

func init() {
	NotNullFuzzIntrospectedTable = tableByName[NotNullFuzzTable]

	/* only needed during templating */
	if NotNullFuzzIntrospectedTable == nil {
		NotNullFuzzIntrospectedTable = &introspect.Table{}
	}

	NotNullFuzzTableColumnLookup = NotNullFuzzIntrospectedTable.ColumnByName
}

type NotNullFuzzOnePathParams struct {
	PrimaryKey int64 `json:"primaryKey"`
}

type NotNullFuzzLoadQueryParams struct {
	Depth *int `json:"depth"`
}

type NotNullFuzzClaimRequest struct {
	Until          time.Time `json:"until"`
	By             uuid.UUID `json:"by"`
	TimeoutSeconds float64   `json:"timeout_seconds"`
}

/*
TODO: find a way to not need this- there is a piece in the templating logic
that uses goimports but pending where the code is built, it may resolve
the packages to import to the wrong ones (causing odd failures)
these are just here to ensure we don't get unused imports
*/
var _ = []any{
	time.Time{},
	uuid.UUID{},
	pgtype.Hstore{},
	postgis.PointZ{},
	netip.Prefix{},
	errors.Is,
	sql.ErrNoRows,
}

func (m *NotNullFuzz) GetPrimaryKeyColumn() string {
	return NotNullFuzzTablePrimaryKeyColumn
}

func (m *NotNullFuzz) GetPrimaryKeyValue() any {
	return m.MrPrimary
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
		return fmt.Errorf("%v: %#+v; error; %v", k, v, err)
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
		case "mr_primary":
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
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uumr_primary.UUID", temp1))
				}
			}

			m.MrPrimary = temp2

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
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uusome_bigint.UUID", temp1))
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
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uusome_bigint_array.UUID", temp1))
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
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uusome_boolean.UUID", temp1))
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
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uusome_boolean_array.UUID", temp1))
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
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uusome_bytea.UUID", temp1))
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
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uusome_character_varying.UUID", temp1))
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
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uusome_character_varying_array.UUID", temp1))
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
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uusome_double_precision.UUID", temp1))
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
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uusome_double_precision_array.UUID", temp1))
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
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uusome_float.UUID", temp1))
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
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uusome_float_array.UUID", temp1))
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
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uusome_geometry_point_z.UUID", temp1))
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
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uusome_hstore.UUID", temp1))
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
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uusome_inet.UUID", temp1))
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
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uusome_integer.UUID", temp1))
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
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uusome_integer_array.UUID", temp1))
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
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uusome_interval.UUID", temp1))
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

			temp2, ok := temp1, true
			if !ok {
				if temp1 != nil {
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uusome_json.UUID", temp1))
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

			temp2, ok := temp1, true
			if !ok {
				if temp1 != nil {
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uusome_jsonb.UUID", temp1))
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
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uusome_numeric.UUID", temp1))
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
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uusome_numeric_array.UUID", temp1))
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
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uusome_point.UUID", temp1))
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
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uusome_polygon.UUID", temp1))
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
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uusome_real.UUID", temp1))
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
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uusome_real_array.UUID", temp1))
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
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uusome_smallint.UUID", temp1))
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
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uusome_smallint_array.UUID", temp1))
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
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uusome_text.UUID", temp1))
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
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uusome_text_array.UUID", temp1))
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
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uusome_timestamptz.UUID", temp1))
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
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uusome_timestamp.UUID", temp1))
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
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uusome_tsvector.UUID", temp1))
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
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uusome_uuid.UUID", temp1))
				}
			}

			m.SomeUUID = temp2

		case "other_not_null_fuzz":
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
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uuother_not_null_fuzz.UUID", temp1))
				}
			}

			m.OtherNotNullFuzz = &temp2

		}
	}

	return nil
}

func (m *NotNullFuzz) Reload(ctx context.Context, tx pgx.Tx, includeDeleteds ...bool) error {
	extraWhere := ""
	if len(includeDeleteds) > 0 && includeDeleteds[0] {
		if slices.Contains(NotNullFuzzTableColumns, "deleted_at") {
			extraWhere = "\n    AND (deleted_at IS null OR deleted_at IS NOT null)"
		}
	}

	ctx, cleanup := query.WithQueryID(ctx)
	defer cleanup()

	ctx = query.WithMaxDepth(ctx, nil)

	o, _, _, _, _, err := SelectNotNullFuzz(
		ctx,
		tx,
		fmt.Sprintf("%v = $1%v", m.GetPrimaryKeyColumn(), extraWhere),
		m.GetPrimaryKeyValue(),
	)
	if err != nil {
		return err
	}

	m.MrPrimary = o.MrPrimary
	m.SomeBigint = o.SomeBigint
	m.SomeBigintArray = o.SomeBigintArray
	m.SomeBoolean = o.SomeBoolean
	m.SomeBooleanArray = o.SomeBooleanArray
	m.SomeBytea = o.SomeBytea
	m.SomeCharacterVarying = o.SomeCharacterVarying
	m.SomeCharacterVaryingArray = o.SomeCharacterVaryingArray
	m.SomeDoublePrecision = o.SomeDoublePrecision
	m.SomeDoublePrecisionArray = o.SomeDoublePrecisionArray
	m.SomeFloat = o.SomeFloat
	m.SomeFloatArray = o.SomeFloatArray
	m.SomeGeometryPointZ = o.SomeGeometryPointZ
	m.SomeHstore = o.SomeHstore
	m.SomeInet = o.SomeInet
	m.SomeInteger = o.SomeInteger
	m.SomeIntegerArray = o.SomeIntegerArray
	m.SomeInterval = o.SomeInterval
	m.SomeJSON = o.SomeJSON
	m.SomeJSONB = o.SomeJSONB
	m.SomeNumeric = o.SomeNumeric
	m.SomeNumericArray = o.SomeNumericArray
	m.SomePoint = o.SomePoint
	m.SomePolygon = o.SomePolygon
	m.SomeReal = o.SomeReal
	m.SomeRealArray = o.SomeRealArray
	m.SomeSmallint = o.SomeSmallint
	m.SomeSmallintArray = o.SomeSmallintArray
	m.SomeText = o.SomeText
	m.SomeTextArray = o.SomeTextArray
	m.SomeTimestamptz = o.SomeTimestamptz
	m.SomeTimestamp = o.SomeTimestamp
	m.SomeTsvector = o.SomeTsvector
	m.SomeUUID = o.SomeUUID
	m.OtherNotNullFuzz = o.OtherNotNullFuzz
	m.OtherNotNullFuzzObject = o.OtherNotNullFuzzObject
	m.ReferencedByNotNullFuzzOtherNotNullFuzzObjects = o.ReferencedByNotNullFuzzOtherNotNullFuzzObjects

	return nil
}

func (m *NotNullFuzz) Insert(ctx context.Context, tx pgx.Tx, setPrimaryKey bool, setZeroValues bool, forceSetValuesForFields ...string) error {
	columns := make([]string, 0)
	values := make([]any, 0)

	if setPrimaryKey && (setZeroValues || !types.IsZeroInt(m.MrPrimary) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableMrPrimaryColumn) || isRequired(NotNullFuzzTableColumnLookup, NotNullFuzzTableMrPrimaryColumn)) {
		columns = append(columns, NotNullFuzzTableMrPrimaryColumn)

		v, err := types.FormatInt(m.MrPrimary)
		if err != nil {
			return fmt.Errorf("failed to handle m.MrPrimary; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroInt(m.SomeBigint) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeBigintColumn) || isRequired(NotNullFuzzTableColumnLookup, NotNullFuzzTableSomeBigintColumn) {
		columns = append(columns, NotNullFuzzTableSomeBigintColumn)

		v, err := types.FormatInt(m.SomeBigint)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeBigint; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroIntArray(m.SomeBigintArray) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeBigintArrayColumn) || isRequired(NotNullFuzzTableColumnLookup, NotNullFuzzTableSomeBigintArrayColumn) {
		columns = append(columns, NotNullFuzzTableSomeBigintArrayColumn)

		v, err := types.FormatIntArray(m.SomeBigintArray)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeBigintArray; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroBool(m.SomeBoolean) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeBooleanColumn) || isRequired(NotNullFuzzTableColumnLookup, NotNullFuzzTableSomeBooleanColumn) {
		columns = append(columns, NotNullFuzzTableSomeBooleanColumn)

		v, err := types.FormatBool(m.SomeBoolean)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeBoolean; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroBoolArray(m.SomeBooleanArray) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeBooleanArrayColumn) || isRequired(NotNullFuzzTableColumnLookup, NotNullFuzzTableSomeBooleanArrayColumn) {
		columns = append(columns, NotNullFuzzTableSomeBooleanArrayColumn)

		v, err := types.FormatBoolArray(m.SomeBooleanArray)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeBooleanArray; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroBytes(m.SomeBytea) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeByteaColumn) || isRequired(NotNullFuzzTableColumnLookup, NotNullFuzzTableSomeByteaColumn) {
		columns = append(columns, NotNullFuzzTableSomeByteaColumn)

		v, err := types.FormatBytes(m.SomeBytea)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeBytea; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroString(m.SomeCharacterVarying) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeCharacterVaryingColumn) || isRequired(NotNullFuzzTableColumnLookup, NotNullFuzzTableSomeCharacterVaryingColumn) {
		columns = append(columns, NotNullFuzzTableSomeCharacterVaryingColumn)

		v, err := types.FormatString(m.SomeCharacterVarying)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeCharacterVarying; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroStringArray(m.SomeCharacterVaryingArray) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeCharacterVaryingArrayColumn) || isRequired(NotNullFuzzTableColumnLookup, NotNullFuzzTableSomeCharacterVaryingArrayColumn) {
		columns = append(columns, NotNullFuzzTableSomeCharacterVaryingArrayColumn)

		v, err := types.FormatStringArray(m.SomeCharacterVaryingArray)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeCharacterVaryingArray; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroFloat(m.SomeDoublePrecision) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeDoublePrecisionColumn) || isRequired(NotNullFuzzTableColumnLookup, NotNullFuzzTableSomeDoublePrecisionColumn) {
		columns = append(columns, NotNullFuzzTableSomeDoublePrecisionColumn)

		v, err := types.FormatFloat(m.SomeDoublePrecision)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeDoublePrecision; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroFloatArray(m.SomeDoublePrecisionArray) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeDoublePrecisionArrayColumn) || isRequired(NotNullFuzzTableColumnLookup, NotNullFuzzTableSomeDoublePrecisionArrayColumn) {
		columns = append(columns, NotNullFuzzTableSomeDoublePrecisionArrayColumn)

		v, err := types.FormatFloatArray(m.SomeDoublePrecisionArray)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeDoublePrecisionArray; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroFloat(m.SomeFloat) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeFloatColumn) || isRequired(NotNullFuzzTableColumnLookup, NotNullFuzzTableSomeFloatColumn) {
		columns = append(columns, NotNullFuzzTableSomeFloatColumn)

		v, err := types.FormatFloat(m.SomeFloat)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeFloat; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroFloatArray(m.SomeFloatArray) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeFloatArrayColumn) || isRequired(NotNullFuzzTableColumnLookup, NotNullFuzzTableSomeFloatArrayColumn) {
		columns = append(columns, NotNullFuzzTableSomeFloatArrayColumn)

		v, err := types.FormatFloatArray(m.SomeFloatArray)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeFloatArray; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroGeometry(m.SomeGeometryPointZ) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeGeometryPointZColumn) || isRequired(NotNullFuzzTableColumnLookup, NotNullFuzzTableSomeGeometryPointZColumn) {
		columns = append(columns, NotNullFuzzTableSomeGeometryPointZColumn)

		v, err := types.FormatGeometry(m.SomeGeometryPointZ)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeGeometryPointZ; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroHstore(m.SomeHstore) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeHstoreColumn) || isRequired(NotNullFuzzTableColumnLookup, NotNullFuzzTableSomeHstoreColumn) {
		columns = append(columns, NotNullFuzzTableSomeHstoreColumn)

		v, err := types.FormatHstore(m.SomeHstore)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeHstore; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroInet(m.SomeInet) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeInetColumn) || isRequired(NotNullFuzzTableColumnLookup, NotNullFuzzTableSomeInetColumn) {
		columns = append(columns, NotNullFuzzTableSomeInetColumn)

		v, err := types.FormatInet(m.SomeInet)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeInet; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroInt(m.SomeInteger) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeIntegerColumn) || isRequired(NotNullFuzzTableColumnLookup, NotNullFuzzTableSomeIntegerColumn) {
		columns = append(columns, NotNullFuzzTableSomeIntegerColumn)

		v, err := types.FormatInt(m.SomeInteger)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeInteger; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroIntArray(m.SomeIntegerArray) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeIntegerArrayColumn) || isRequired(NotNullFuzzTableColumnLookup, NotNullFuzzTableSomeIntegerArrayColumn) {
		columns = append(columns, NotNullFuzzTableSomeIntegerArrayColumn)

		v, err := types.FormatIntArray(m.SomeIntegerArray)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeIntegerArray; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroDuration(m.SomeInterval) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeIntervalColumn) || isRequired(NotNullFuzzTableColumnLookup, NotNullFuzzTableSomeIntervalColumn) {
		columns = append(columns, NotNullFuzzTableSomeIntervalColumn)

		v, err := types.FormatDuration(m.SomeInterval)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeInterval; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroJSON(m.SomeJSON) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeJSONColumn) || isRequired(NotNullFuzzTableColumnLookup, NotNullFuzzTableSomeJSONColumn) {
		columns = append(columns, NotNullFuzzTableSomeJSONColumn)

		v, err := types.FormatJSON(m.SomeJSON)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeJSON; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroJSON(m.SomeJSONB) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeJSONBColumn) || isRequired(NotNullFuzzTableColumnLookup, NotNullFuzzTableSomeJSONBColumn) {
		columns = append(columns, NotNullFuzzTableSomeJSONBColumn)

		v, err := types.FormatJSON(m.SomeJSONB)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeJSONB; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroFloat(m.SomeNumeric) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeNumericColumn) || isRequired(NotNullFuzzTableColumnLookup, NotNullFuzzTableSomeNumericColumn) {
		columns = append(columns, NotNullFuzzTableSomeNumericColumn)

		v, err := types.FormatFloat(m.SomeNumeric)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeNumeric; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroFloatArray(m.SomeNumericArray) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeNumericArrayColumn) || isRequired(NotNullFuzzTableColumnLookup, NotNullFuzzTableSomeNumericArrayColumn) {
		columns = append(columns, NotNullFuzzTableSomeNumericArrayColumn)

		v, err := types.FormatFloatArray(m.SomeNumericArray)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeNumericArray; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroPoint(m.SomePoint) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomePointColumn) || isRequired(NotNullFuzzTableColumnLookup, NotNullFuzzTableSomePointColumn) {
		columns = append(columns, NotNullFuzzTableSomePointColumn)

		v, err := types.FormatPoint(m.SomePoint)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomePoint; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroPolygon(m.SomePolygon) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomePolygonColumn) || isRequired(NotNullFuzzTableColumnLookup, NotNullFuzzTableSomePolygonColumn) {
		columns = append(columns, NotNullFuzzTableSomePolygonColumn)

		v, err := types.FormatPolygon(m.SomePolygon)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomePolygon; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroFloat(m.SomeReal) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeRealColumn) || isRequired(NotNullFuzzTableColumnLookup, NotNullFuzzTableSomeRealColumn) {
		columns = append(columns, NotNullFuzzTableSomeRealColumn)

		v, err := types.FormatFloat(m.SomeReal)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeReal; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroFloatArray(m.SomeRealArray) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeRealArrayColumn) || isRequired(NotNullFuzzTableColumnLookup, NotNullFuzzTableSomeRealArrayColumn) {
		columns = append(columns, NotNullFuzzTableSomeRealArrayColumn)

		v, err := types.FormatFloatArray(m.SomeRealArray)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeRealArray; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroInt(m.SomeSmallint) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeSmallintColumn) || isRequired(NotNullFuzzTableColumnLookup, NotNullFuzzTableSomeSmallintColumn) {
		columns = append(columns, NotNullFuzzTableSomeSmallintColumn)

		v, err := types.FormatInt(m.SomeSmallint)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeSmallint; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroIntArray(m.SomeSmallintArray) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeSmallintArrayColumn) || isRequired(NotNullFuzzTableColumnLookup, NotNullFuzzTableSomeSmallintArrayColumn) {
		columns = append(columns, NotNullFuzzTableSomeSmallintArrayColumn)

		v, err := types.FormatIntArray(m.SomeSmallintArray)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeSmallintArray; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroString(m.SomeText) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeTextColumn) || isRequired(NotNullFuzzTableColumnLookup, NotNullFuzzTableSomeTextColumn) {
		columns = append(columns, NotNullFuzzTableSomeTextColumn)

		v, err := types.FormatString(m.SomeText)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeText; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroStringArray(m.SomeTextArray) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeTextArrayColumn) || isRequired(NotNullFuzzTableColumnLookup, NotNullFuzzTableSomeTextArrayColumn) {
		columns = append(columns, NotNullFuzzTableSomeTextArrayColumn)

		v, err := types.FormatStringArray(m.SomeTextArray)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeTextArray; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.SomeTimestamptz) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeTimestamptzColumn) || isRequired(NotNullFuzzTableColumnLookup, NotNullFuzzTableSomeTimestamptzColumn) {
		columns = append(columns, NotNullFuzzTableSomeTimestamptzColumn)

		v, err := types.FormatTime(m.SomeTimestamptz)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeTimestamptz; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.SomeTimestamp) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeTimestampColumn) || isRequired(NotNullFuzzTableColumnLookup, NotNullFuzzTableSomeTimestampColumn) {
		columns = append(columns, NotNullFuzzTableSomeTimestampColumn)

		v, err := types.FormatTime(m.SomeTimestamp)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeTimestamp; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTSVector(m.SomeTsvector) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeTsvectorColumn) || isRequired(NotNullFuzzTableColumnLookup, NotNullFuzzTableSomeTsvectorColumn) {
		columns = append(columns, NotNullFuzzTableSomeTsvectorColumn)

		v, err := types.FormatTSVector(m.SomeTsvector)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeTsvector; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroUUID(m.SomeUUID) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeUUIDColumn) || isRequired(NotNullFuzzTableColumnLookup, NotNullFuzzTableSomeUUIDColumn) {
		columns = append(columns, NotNullFuzzTableSomeUUIDColumn)

		v, err := types.FormatUUID(m.SomeUUID)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeUUID; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroInt(m.OtherNotNullFuzz) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableOtherNotNullFuzzColumn) || isRequired(NotNullFuzzTableColumnLookup, NotNullFuzzTableOtherNotNullFuzzColumn) {
		columns = append(columns, NotNullFuzzTableOtherNotNullFuzzColumn)

		v, err := types.FormatInt(m.OtherNotNullFuzz)
		if err != nil {
			return fmt.Errorf("failed to handle m.OtherNotNullFuzz; %v", err)
		}

		values = append(values, v)
	}

	ctx, cleanup := query.WithQueryID(ctx)
	defer cleanup()

	ctx = query.WithMaxDepth(ctx, nil)

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
		return fmt.Errorf("failed to insert %#+v; %v", m, err)
	}
	v := (*item)[NotNullFuzzTableMrPrimaryColumn]

	if v == nil {
		return fmt.Errorf("failed to find %v in %#+v", NotNullFuzzTableMrPrimaryColumn, item)
	}

	wrapError := func(err error) error {
		return fmt.Errorf(
			"failed to treat %v: %#+v as int64: %v",
			NotNullFuzzTableMrPrimaryColumn,
			(*item)[NotNullFuzzTableMrPrimaryColumn],
			err,
		)
	}

	temp1, err := types.ParseInt(v)
	if err != nil {
		return wrapError(err)
	}

	temp2, ok := temp1.(int64)
	if !ok {
		return wrapError(fmt.Errorf("failed to cast to int64"))
	}

	m.MrPrimary = temp2

	err = m.Reload(ctx, tx, slices.Contains(forceSetValuesForFields, "deleted_at"))
	if err != nil {
		return fmt.Errorf("failed to reload after insert; %v", err)
	}

	return nil
}

func (m *NotNullFuzz) Update(ctx context.Context, tx pgx.Tx, setZeroValues bool, forceSetValuesForFields ...string) error {
	columns := make([]string, 0)
	values := make([]any, 0)

	if setZeroValues || !types.IsZeroInt(m.SomeBigint) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeBigintColumn) {
		columns = append(columns, NotNullFuzzTableSomeBigintColumn)

		v, err := types.FormatInt(m.SomeBigint)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeBigint; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroIntArray(m.SomeBigintArray) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeBigintArrayColumn) {
		columns = append(columns, NotNullFuzzTableSomeBigintArrayColumn)

		v, err := types.FormatIntArray(m.SomeBigintArray)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeBigintArray; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroBool(m.SomeBoolean) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeBooleanColumn) {
		columns = append(columns, NotNullFuzzTableSomeBooleanColumn)

		v, err := types.FormatBool(m.SomeBoolean)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeBoolean; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroBoolArray(m.SomeBooleanArray) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeBooleanArrayColumn) {
		columns = append(columns, NotNullFuzzTableSomeBooleanArrayColumn)

		v, err := types.FormatBoolArray(m.SomeBooleanArray)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeBooleanArray; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroBytes(m.SomeBytea) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeByteaColumn) {
		columns = append(columns, NotNullFuzzTableSomeByteaColumn)

		v, err := types.FormatBytes(m.SomeBytea)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeBytea; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroString(m.SomeCharacterVarying) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeCharacterVaryingColumn) {
		columns = append(columns, NotNullFuzzTableSomeCharacterVaryingColumn)

		v, err := types.FormatString(m.SomeCharacterVarying)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeCharacterVarying; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroStringArray(m.SomeCharacterVaryingArray) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeCharacterVaryingArrayColumn) {
		columns = append(columns, NotNullFuzzTableSomeCharacterVaryingArrayColumn)

		v, err := types.FormatStringArray(m.SomeCharacterVaryingArray)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeCharacterVaryingArray; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroFloat(m.SomeDoublePrecision) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeDoublePrecisionColumn) {
		columns = append(columns, NotNullFuzzTableSomeDoublePrecisionColumn)

		v, err := types.FormatFloat(m.SomeDoublePrecision)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeDoublePrecision; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroFloatArray(m.SomeDoublePrecisionArray) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeDoublePrecisionArrayColumn) {
		columns = append(columns, NotNullFuzzTableSomeDoublePrecisionArrayColumn)

		v, err := types.FormatFloatArray(m.SomeDoublePrecisionArray)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeDoublePrecisionArray; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroFloat(m.SomeFloat) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeFloatColumn) {
		columns = append(columns, NotNullFuzzTableSomeFloatColumn)

		v, err := types.FormatFloat(m.SomeFloat)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeFloat; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroFloatArray(m.SomeFloatArray) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeFloatArrayColumn) {
		columns = append(columns, NotNullFuzzTableSomeFloatArrayColumn)

		v, err := types.FormatFloatArray(m.SomeFloatArray)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeFloatArray; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroGeometry(m.SomeGeometryPointZ) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeGeometryPointZColumn) {
		columns = append(columns, NotNullFuzzTableSomeGeometryPointZColumn)

		v, err := types.FormatGeometry(m.SomeGeometryPointZ)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeGeometryPointZ; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroHstore(m.SomeHstore) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeHstoreColumn) {
		columns = append(columns, NotNullFuzzTableSomeHstoreColumn)

		v, err := types.FormatHstore(m.SomeHstore)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeHstore; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroInet(m.SomeInet) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeInetColumn) {
		columns = append(columns, NotNullFuzzTableSomeInetColumn)

		v, err := types.FormatInet(m.SomeInet)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeInet; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroInt(m.SomeInteger) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeIntegerColumn) {
		columns = append(columns, NotNullFuzzTableSomeIntegerColumn)

		v, err := types.FormatInt(m.SomeInteger)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeInteger; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroIntArray(m.SomeIntegerArray) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeIntegerArrayColumn) {
		columns = append(columns, NotNullFuzzTableSomeIntegerArrayColumn)

		v, err := types.FormatIntArray(m.SomeIntegerArray)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeIntegerArray; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroDuration(m.SomeInterval) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeIntervalColumn) {
		columns = append(columns, NotNullFuzzTableSomeIntervalColumn)

		v, err := types.FormatDuration(m.SomeInterval)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeInterval; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroJSON(m.SomeJSON) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeJSONColumn) {
		columns = append(columns, NotNullFuzzTableSomeJSONColumn)

		v, err := types.FormatJSON(m.SomeJSON)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeJSON; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroJSON(m.SomeJSONB) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeJSONBColumn) {
		columns = append(columns, NotNullFuzzTableSomeJSONBColumn)

		v, err := types.FormatJSON(m.SomeJSONB)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeJSONB; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroFloat(m.SomeNumeric) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeNumericColumn) {
		columns = append(columns, NotNullFuzzTableSomeNumericColumn)

		v, err := types.FormatFloat(m.SomeNumeric)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeNumeric; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroFloatArray(m.SomeNumericArray) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeNumericArrayColumn) {
		columns = append(columns, NotNullFuzzTableSomeNumericArrayColumn)

		v, err := types.FormatFloatArray(m.SomeNumericArray)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeNumericArray; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroPoint(m.SomePoint) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomePointColumn) {
		columns = append(columns, NotNullFuzzTableSomePointColumn)

		v, err := types.FormatPoint(m.SomePoint)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomePoint; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroPolygon(m.SomePolygon) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomePolygonColumn) {
		columns = append(columns, NotNullFuzzTableSomePolygonColumn)

		v, err := types.FormatPolygon(m.SomePolygon)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomePolygon; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroFloat(m.SomeReal) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeRealColumn) {
		columns = append(columns, NotNullFuzzTableSomeRealColumn)

		v, err := types.FormatFloat(m.SomeReal)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeReal; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroFloatArray(m.SomeRealArray) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeRealArrayColumn) {
		columns = append(columns, NotNullFuzzTableSomeRealArrayColumn)

		v, err := types.FormatFloatArray(m.SomeRealArray)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeRealArray; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroInt(m.SomeSmallint) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeSmallintColumn) {
		columns = append(columns, NotNullFuzzTableSomeSmallintColumn)

		v, err := types.FormatInt(m.SomeSmallint)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeSmallint; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroIntArray(m.SomeSmallintArray) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeSmallintArrayColumn) {
		columns = append(columns, NotNullFuzzTableSomeSmallintArrayColumn)

		v, err := types.FormatIntArray(m.SomeSmallintArray)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeSmallintArray; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroString(m.SomeText) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeTextColumn) {
		columns = append(columns, NotNullFuzzTableSomeTextColumn)

		v, err := types.FormatString(m.SomeText)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeText; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroStringArray(m.SomeTextArray) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeTextArrayColumn) {
		columns = append(columns, NotNullFuzzTableSomeTextArrayColumn)

		v, err := types.FormatStringArray(m.SomeTextArray)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeTextArray; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.SomeTimestamptz) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeTimestamptzColumn) {
		columns = append(columns, NotNullFuzzTableSomeTimestamptzColumn)

		v, err := types.FormatTime(m.SomeTimestamptz)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeTimestamptz; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.SomeTimestamp) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeTimestampColumn) {
		columns = append(columns, NotNullFuzzTableSomeTimestampColumn)

		v, err := types.FormatTime(m.SomeTimestamp)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeTimestamp; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTSVector(m.SomeTsvector) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeTsvectorColumn) {
		columns = append(columns, NotNullFuzzTableSomeTsvectorColumn)

		v, err := types.FormatTSVector(m.SomeTsvector)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeTsvector; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroUUID(m.SomeUUID) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableSomeUUIDColumn) {
		columns = append(columns, NotNullFuzzTableSomeUUIDColumn)

		v, err := types.FormatUUID(m.SomeUUID)
		if err != nil {
			return fmt.Errorf("failed to handle m.SomeUUID; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroInt(m.OtherNotNullFuzz) || slices.Contains(forceSetValuesForFields, NotNullFuzzTableOtherNotNullFuzzColumn) {
		columns = append(columns, NotNullFuzzTableOtherNotNullFuzzColumn)

		v, err := types.FormatInt(m.OtherNotNullFuzz)
		if err != nil {
			return fmt.Errorf("failed to handle m.OtherNotNullFuzz; %v", err)
		}

		values = append(values, v)
	}

	v, err := types.FormatInt(m.MrPrimary)
	if err != nil {
		return fmt.Errorf("failed to handle m.MrPrimary; %v", err)
	}

	values = append(values, v)

	ctx, cleanup := query.WithQueryID(ctx)
	defer cleanup()

	ctx = query.WithMaxDepth(ctx, nil)

	_, err = query.Update(
		ctx,
		tx,
		NotNullFuzzTable,
		columns,
		fmt.Sprintf("%v = $$??", NotNullFuzzTableMrPrimaryColumn),
		NotNullFuzzTableColumns,
		values...,
	)
	if err != nil {
		return fmt.Errorf("failed to update %#+v; %v", m, err)
	}

	err = m.Reload(ctx, tx, slices.Contains(forceSetValuesForFields, "deleted_at"))
	if err != nil {
		return fmt.Errorf("failed to reload after update")
	}

	return nil
}

func (m *NotNullFuzz) Delete(ctx context.Context, tx pgx.Tx, hardDeletes ...bool) error {
	/* soft-delete not applicable */

	values := make([]any, 0)
	v, err := types.FormatInt(m.MrPrimary)
	if err != nil {
		return fmt.Errorf("failed to handle m.MrPrimary; %v", err)
	}

	values = append(values, v)

	ctx, cleanup := query.WithQueryID(ctx)
	defer cleanup()

	ctx = query.WithMaxDepth(ctx, nil)

	err = query.Delete(
		ctx,
		tx,
		NotNullFuzzTable,
		fmt.Sprintf("%v = $$??", NotNullFuzzTableMrPrimaryColumn),
		values...,
	)
	if err != nil {
		return fmt.Errorf("failed to delete %#+v; %v", m, err)
	}

	_ = m.Reload(ctx, tx, true)

	return nil
}

func (m *NotNullFuzz) LockTable(ctx context.Context, tx pgx.Tx, timeouts ...time.Duration) error {
	return query.LockTable(ctx, tx, NotNullFuzzTable, timeouts...)
}

func (m *NotNullFuzz) LockTableWithRetries(ctx context.Context, tx pgx.Tx, overallTimeout time.Duration, individualAttempttimeout time.Duration) error {
	return query.LockTableWithRetries(ctx, tx, NotNullFuzzTable, overallTimeout, individualAttempttimeout)
}

func (m *NotNullFuzz) AdvisoryLock(ctx context.Context, tx pgx.Tx, key int32, timeouts ...time.Duration) error {
	return query.AdvisoryLock(ctx, tx, NotNullFuzzTableNamespaceID, key, timeouts...)
}

func (m *NotNullFuzz) AdvisoryLockWithRetries(ctx context.Context, tx pgx.Tx, key int32, overallTimeout time.Duration, individualAttempttimeout time.Duration) error {
	return query.AdvisoryLockWithRetries(ctx, tx, NotNullFuzzTableNamespaceID, key, overallTimeout, individualAttempttimeout)
}

func (m *NotNullFuzz) Claim(ctx context.Context, tx pgx.Tx, until time.Time, by uuid.UUID, timeout time.Duration) error {
	if !(slices.Contains(NotNullFuzzTableColumns, "claimed_until") && slices.Contains(NotNullFuzzTableColumns, "claimed_by")) {
		return fmt.Errorf("can only invoke Claim for tables with 'claimed_until' and 'claimed_by' columns")
	}

	err := m.AdvisoryLockWithRetries(ctx, tx, math.MinInt32, timeout, time.Second*1)
	if err != nil {
		return fmt.Errorf("failed to claim (advisory lock): %s", err.Error())
	}

	x, _, _, _, _, err := SelectNotNullFuzz(
		ctx,
		tx,
		fmt.Sprintf(
			"%s = $$?? AND (claimed_by = $$?? OR (claimed_until IS null OR claimed_until < now()))",
			NotNullFuzzTablePrimaryKeyColumn,
		),
		m.GetPrimaryKeyValue(),
		by,
	)
	if err != nil {
		return fmt.Errorf("failed to claim (select): %s", err.Error())
	}

	_ = x

	/* m.ClaimedUntil = &until */
	/* m.ClaimedBy = &by */

	err = m.Update(ctx, tx, false)
	if err != nil {
		return fmt.Errorf("failed to claim (update): %s", err.Error())
	}

	return nil
}

func SelectNotNullFuzzes(ctx context.Context, tx pgx.Tx, where string, orderBy *string, limit *int, offset *int, values ...any) ([]*NotNullFuzz, int64, int64, int64, int64, error) {
	before := time.Now()

	if config.Debug() {
		log.Printf("entered SelectNotNullFuzzes")

		defer func() {
			log.Printf("exited SelectNotNullFuzzes in %s", time.Since(before))
		}()
	}
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

	possiblePathValue := query.GetCurrentPathValue(ctx)
	isLoadQuery := possiblePathValue != nil && len(possiblePathValue.VisitedTableNames) > 0

	shouldLoad := query.ShouldLoad(ctx, NotNullFuzzTable) || query.ShouldLoad(ctx, fmt.Sprintf("referenced_by_%s", NotNullFuzzTable))

	var ok bool
	ctx, ok = query.HandleQueryPathGraphCycles(ctx, fmt.Sprintf("%s{%v}", NotNullFuzzTable, nil), !isLoadQuery)
	if !ok && !shouldLoad {
		if config.Debug() {
			log.Printf("skipping SelectNotNullFuzz early (query.ShouldLoad(): %v, query.HandleQueryPathGraphCycles(): %v)", shouldLoad, ok)
		}
		return []*NotNullFuzz{}, 0, 0, 0, 0, nil
	}

	items, count, totalCount, page, totalPages, err := query.Select(
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
		return nil, 0, 0, 0, 0, fmt.Errorf("failed to call SelectNotNullFuzzs; %v", err)
	}

	objects := make([]*NotNullFuzz, 0)

	for _, item := range *items {
		object := &NotNullFuzz{}

		err = object.FromItem(item)
		if err != nil {
			return nil, 0, 0, 0, 0, err
		}

		if !types.IsZeroInt(object.OtherNotNullFuzz) {
			ctx, ok := query.HandleQueryPathGraphCycles(ctx, fmt.Sprintf("%s{%v}", NotNullFuzzTable, object.OtherNotNullFuzz), true)
			shouldLoad := query.ShouldLoad(ctx, NotNullFuzzTable)
			if ok || shouldLoad {
				thisBefore := time.Now()

				if config.Debug() {
					log.Printf("loading SelectNotNullFuzzes->SelectNotNullFuzz for object.OtherNotNullFuzzObject{%s: %v}", NotNullFuzzTablePrimaryKeyColumn, object.OtherNotNullFuzz)
				}

				object.OtherNotNullFuzzObject, _, _, _, _, err = SelectNotNullFuzz(
					ctx,
					tx,
					fmt.Sprintf("%v = $1", NotNullFuzzTablePrimaryKeyColumn),
					object.OtherNotNullFuzz,
				)
				if err != nil {
					if !errors.Is(err, sql.ErrNoRows) {
						return nil, 0, 0, 0, 0, err
					}
				}

				if config.Debug() {
					log.Printf("loaded SelectNotNullFuzzes->SelectNotNullFuzz for object.OtherNotNullFuzzObject in %s", time.Since(thisBefore))
				}
			}
		}

		/*
			err = func() error {
				shouldLoad := query.ShouldLoad(ctx, fmt.Sprintf("referenced_by_%s", NotNullFuzzTable))
				ctx, ok := query.HandleQueryPathGraphCycles(ctx, fmt.Sprintf("__ReferencedBy__%s{%v}", NotNullFuzzTable, object.GetPrimaryKeyValue()), true)
				if ok || shouldLoad {
					thisBefore := time.Now()

					if config.Debug() {
						log.Printf("loading SelectNotNullFuzzes->SelectNotNullFuzzes for object.ReferencedByNotNullFuzzOtherNotNullFuzzObjects")
					}

					object.ReferencedByNotNullFuzzOtherNotNullFuzzObjects, _, _, _, _, err = SelectNotNullFuzzes(
						ctx,
						tx,
						fmt.Sprintf("%v = $1", NotNullFuzzTableOtherNotNullFuzzColumn),
						nil,
						nil,
						nil,
						object.GetPrimaryKeyValue(),
					)
					if err != nil {
						if !errors.Is(err, sql.ErrNoRows) {
							return err
						}
					}

					if config.Debug() {
						log.Printf("loaded SelectNotNullFuzzes->SelectNotNullFuzzes for object.ReferencedByNotNullFuzzOtherNotNullFuzzObjects in %s", time.Since(thisBefore))
					}

				}

				return nil
			}()
			if err != nil {
				return nil, 0, 0, 0, 0, err
			}
		*/

		objects = append(objects, object)
	}

	return objects, count, totalCount, page, totalPages, nil
}

func SelectNotNullFuzz(ctx context.Context, tx pgx.Tx, where string, values ...any) (*NotNullFuzz, int64, int64, int64, int64, error) {
	ctx, cleanup := query.WithQueryID(ctx)
	defer cleanup()

	ctx = query.WithMaxDepth(ctx, nil)

	objects, _, _, _, _, err := SelectNotNullFuzzes(
		ctx,
		tx,
		where,
		nil,
		helpers.Ptr(2),
		helpers.Ptr(0),
		values...,
	)
	if err != nil {
		return nil, 0, 0, 0, 0, fmt.Errorf("failed to call SelectNotNullFuzz; %v", err)
	}

	if len(objects) > 1 {
		return nil, 0, 0, 0, 0, fmt.Errorf("attempt to call SelectNotNullFuzz returned more than 1 row")
	}

	if len(objects) < 1 {
		return nil, 0, 0, 0, 0, sql.ErrNoRows
	}

	object := objects[0]

	count := int64(1)
	totalCount := count
	page := int64(1)
	totalPages := page

	return object, count, totalCount, page, totalPages, nil
}

func ClaimNotNullFuzz(ctx context.Context, tx pgx.Tx, until time.Time, by uuid.UUID, timeout time.Duration, wheres ...string) (*NotNullFuzz, error) {
	if !(slices.Contains(NotNullFuzzTableColumns, "claimed_until") && slices.Contains(NotNullFuzzTableColumns, "claimed_by")) {
		return nil, fmt.Errorf("can only invoke Claim for tables with 'claimed_until' and 'claimed_by' columns")
	}

	m := &NotNullFuzz{}

	err := m.AdvisoryLockWithRetries(ctx, tx, math.MinInt32, timeout, time.Second*1)
	if err != nil {
		return nil, fmt.Errorf("failed to claim: %s", err.Error())
	}

	extraWhere := ""
	if len(wheres) > 0 {
		extraWhere = fmt.Sprintf("AND %s", extraWhere)
	}

	ms, _, _, _, _, err := SelectNotNullFuzzes(
		ctx,
		tx,
		fmt.Sprintf(
			"(claimed_until IS null OR claimed_until < now())%s",
			extraWhere,
		),
		helpers.Ptr(
			"claimed_until ASC",
		),
		helpers.Ptr(1),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to claim: %s", err.Error())
	}

	if len(ms) == 0 {
		return nil, nil
	}

	m = ms[0]

	/* m.ClaimedUntil = &until */
	/* m.ClaimedBy = &by */

	err = m.Update(ctx, tx, false)
	if err != nil {
		return nil, fmt.Errorf("failed to claim: %s", err.Error())
	}

	return m, nil
}

func handleGetNotNullFuzzes(arguments *server.SelectManyArguments, db *pgxpool.Pool) ([]*NotNullFuzz, int64, int64, int64, int64, error) {
	tx, err := db.Begin(arguments.Ctx)
	if err != nil {
		return nil, 0, 0, 0, 0, err
	}

	defer func() {
		_ = tx.Rollback(arguments.Ctx)
	}()

	objects, count, totalCount, page, totalPages, err := SelectNotNullFuzzes(arguments.Ctx, tx, arguments.Where, arguments.OrderBy, arguments.Limit, arguments.Offset, arguments.Values...)
	if err != nil {
		return nil, 0, 0, 0, 0, err
	}

	err = tx.Commit(arguments.Ctx)
	if err != nil {
		return nil, 0, 0, 0, 0, err
	}

	return objects, count, totalCount, page, totalPages, nil
}

func handleGetNotNullFuzz(arguments *server.SelectOneArguments, db *pgxpool.Pool, primaryKey int64) ([]*NotNullFuzz, int64, int64, int64, int64, error) {
	tx, err := db.Begin(arguments.Ctx)
	if err != nil {
		return nil, 0, 0, 0, 0, err
	}

	defer func() {
		_ = tx.Rollback(arguments.Ctx)
	}()

	object, count, totalCount, page, totalPages, err := SelectNotNullFuzz(arguments.Ctx, tx, arguments.Where, arguments.Values...)
	if err != nil {
		return nil, 0, 0, 0, 0, err
	}

	err = tx.Commit(arguments.Ctx)
	if err != nil {
		return nil, 0, 0, 0, 0, err
	}

	return []*NotNullFuzz{object}, count, totalCount, page, totalPages, nil
}

func handlePostNotNullFuzzs(arguments *server.LoadArguments, db *pgxpool.Pool, waitForChange server.WaitForChange, objects []*NotNullFuzz, forceSetValuesForFieldsByObjectIndex [][]string) ([]*NotNullFuzz, int64, int64, int64, int64, error) {
	tx, err := db.Begin(arguments.Ctx)
	if err != nil {
		err = fmt.Errorf("failed to begin DB transaction; %v", err)
		return nil, 0, 0, 0, 0, err
	}

	defer func() {
		_ = tx.Rollback(arguments.Ctx)
	}()

	xid, err := query.GetXid(arguments.Ctx, tx)
	if err != nil {
		err = fmt.Errorf("failed to get xid; %v", err)
		return nil, 0, 0, 0, 0, err
	}
	_ = xid

	for i, object := range objects {
		err = object.Insert(arguments.Ctx, tx, false, false, forceSetValuesForFieldsByObjectIndex[i]...)
		if err != nil {
			err = fmt.Errorf("failed to insert %#+v; %v", object, err)
			return nil, 0, 0, 0, 0, err
		}

		objects[i] = object
	}

	errs := make(chan error, 1)
	go func() {
		_, err := waitForChange(arguments.Ctx, []stream.Action{stream.INSERT}, NotNullFuzzTable, xid)
		if err != nil {
			err = fmt.Errorf("failed to wait for change; %v", err)
			errs <- err
			return
		}

		errs <- nil
	}()

	err = tx.Commit(arguments.Ctx)
	if err != nil {
		err = fmt.Errorf("failed to commit DB transaction; %v", err)
		return nil, 0, 0, 0, 0, err
	}

	select {
	case <-arguments.Ctx.Done():
		err = fmt.Errorf("context canceled")
		return nil, 0, 0, 0, 0, err
	case err = <-errs:
		if err != nil {
			return nil, 0, 0, 0, 0, err
		}
	}

	count := int64(len(objects))
	totalCount := count
	page := int64(1)
	totalPages := page

	return objects, count, totalCount, page, totalPages, nil
}

func handlePutNotNullFuzz(arguments *server.LoadArguments, db *pgxpool.Pool, waitForChange server.WaitForChange, object *NotNullFuzz) ([]*NotNullFuzz, int64, int64, int64, int64, error) {
	tx, err := db.Begin(arguments.Ctx)
	if err != nil {
		err = fmt.Errorf("failed to begin DB transaction; %v", err)
		return nil, 0, 0, 0, 0, err
	}

	defer func() {
		_ = tx.Rollback(arguments.Ctx)
	}()

	xid, err := query.GetXid(arguments.Ctx, tx)
	if err != nil {
		err = fmt.Errorf("failed to get xid; %v", err)
		return nil, 0, 0, 0, 0, err
	}
	_ = xid

	err = object.Update(arguments.Ctx, tx, true)
	if err != nil {
		err = fmt.Errorf("failed to update %#+v; %v", object, err)
		return nil, 0, 0, 0, 0, err
	}

	errs := make(chan error, 1)
	go func() {
		_, err := waitForChange(arguments.Ctx, []stream.Action{stream.UPDATE, stream.SOFT_DELETE, stream.SOFT_RESTORE, stream.SOFT_UPDATE}, NotNullFuzzTable, xid)
		if err != nil {
			err = fmt.Errorf("failed to wait for change; %v", err)
			errs <- err
			return
		}

		errs <- nil
	}()

	err = tx.Commit(arguments.Ctx)
	if err != nil {
		err = fmt.Errorf("failed to commit DB transaction; %v", err)
		return nil, 0, 0, 0, 0, err
	}

	select {
	case <-arguments.Ctx.Done():
		err = fmt.Errorf("context canceled")
		return nil, 0, 0, 0, 0, err
	case err = <-errs:
		if err != nil {
			return nil, 0, 0, 0, 0, err
		}
	}

	count := int64(1)
	totalCount := count
	page := int64(1)
	totalPages := page

	return []*NotNullFuzz{object}, count, totalCount, page, totalPages, nil
}

func handlePatchNotNullFuzz(arguments *server.LoadArguments, db *pgxpool.Pool, waitForChange server.WaitForChange, object *NotNullFuzz, forceSetValuesForFields []string) ([]*NotNullFuzz, int64, int64, int64, int64, error) {
	tx, err := db.Begin(arguments.Ctx)
	if err != nil {
		err = fmt.Errorf("failed to begin DB transaction; %v", err)
		return nil, 0, 0, 0, 0, err
	}

	defer func() {
		_ = tx.Rollback(arguments.Ctx)
	}()

	xid, err := query.GetXid(arguments.Ctx, tx)
	if err != nil {
		err = fmt.Errorf("failed to get xid; %v", err)
		return nil, 0, 0, 0, 0, err
	}
	_ = xid

	err = object.Update(arguments.Ctx, tx, false, forceSetValuesForFields...)
	if err != nil {
		err = fmt.Errorf("failed to update %#+v; %v", object, err)
		return nil, 0, 0, 0, 0, err
	}

	errs := make(chan error, 1)
	go func() {
		_, err := waitForChange(arguments.Ctx, []stream.Action{stream.UPDATE, stream.SOFT_DELETE, stream.SOFT_RESTORE, stream.SOFT_UPDATE}, NotNullFuzzTable, xid)
		if err != nil {
			err = fmt.Errorf("failed to wait for change; %v", err)
			errs <- err
			return
		}

		errs <- nil
	}()

	err = tx.Commit(arguments.Ctx)
	if err != nil {
		err = fmt.Errorf("failed to commit DB transaction; %v", err)
		return nil, 0, 0, 0, 0, err
	}

	select {
	case <-arguments.Ctx.Done():
		err = fmt.Errorf("context canceled")
		return nil, 0, 0, 0, 0, err
	case err = <-errs:
		if err != nil {
			return nil, 0, 0, 0, 0, err
		}
	}

	count := int64(1)
	totalCount := count
	page := int64(1)
	totalPages := page

	return []*NotNullFuzz{object}, count, totalCount, page, totalPages, nil
}

func handleDeleteNotNullFuzz(arguments *server.LoadArguments, db *pgxpool.Pool, waitForChange server.WaitForChange, object *NotNullFuzz) error {
	tx, err := db.Begin(arguments.Ctx)
	if err != nil {
		err = fmt.Errorf("failed to begin DB transaction; %v", err)
		return err
	}

	defer func() {
		_ = tx.Rollback(arguments.Ctx)
	}()

	xid, err := query.GetXid(arguments.Ctx, tx)
	if err != nil {
		err = fmt.Errorf("failed to get xid; %v", err)
		return err
	}
	_ = xid

	err = object.Delete(arguments.Ctx, tx)
	if err != nil {
		err = fmt.Errorf("failed to delete %#+v; %v", object, err)
		return err
	}

	errs := make(chan error, 1)
	go func() {
		_, err := waitForChange(arguments.Ctx, []stream.Action{stream.DELETE, stream.SOFT_DELETE}, NotNullFuzzTable, xid)
		if err != nil {
			err = fmt.Errorf("failed to wait for change; %v", err)
			errs <- err
			return
		}

		errs <- nil
	}()

	err = tx.Commit(arguments.Ctx)
	if err != nil {
		err = fmt.Errorf("failed to commit DB transaction; %v", err)
		return err
	}

	select {
	case <-arguments.Ctx.Done():
		err = fmt.Errorf("context canceled")
		return err
	case err = <-errs:
		if err != nil {
			return err
		}
	}

	return nil
}

func MutateRouterForNotNullFuzz(r chi.Router, db *pgxpool.Pool, redisPool *redis.Pool, objectMiddlewares []server.ObjectMiddleware, waitForChange server.WaitForChange) {
	if slices.Contains(NotNullFuzzTableColumns, "claimed_until") && slices.Contains(NotNullFuzzTableColumns, "claimed_by") {
		func() {
			postHandlerForClaim, err := getHTTPHandler(
				http.MethodPost,
				"/claim-not-null-fuzz",
				http.StatusOK,
				func(
					ctx context.Context,
					pathParams server.EmptyPathParams,
					queryParams server.EmptyQueryParams,
					req NotNullFuzzClaimRequest,
					rawReq any,
				) (server.Response[NotNullFuzz], error) {
					tx, err := db.Begin(ctx)
					if err != nil {
						return server.Response[NotNullFuzz]{}, err
					}

					defer func() {
						_ = tx.Rollback(ctx)
					}()

					object, err := ClaimNotNullFuzz(ctx, tx, req.Until, req.By, time.Millisecond*time.Duration(req.TimeoutSeconds*1000))
					if err != nil {
						return server.Response[NotNullFuzz]{}, err
					}

					count := int64(0)

					totalCount := int64(0)

					limit := int64(0)

					offset := int64(0)

					if object == nil {
						return server.Response[NotNullFuzz]{
							Status:     http.StatusOK,
							Success:    true,
							Error:      nil,
							Objects:    []*NotNullFuzz{},
							Count:      count,
							TotalCount: totalCount,
							Limit:      limit,
							Offset:     offset,
						}, nil
					}

					err = tx.Commit(ctx)
					if err != nil {
						return server.Response[NotNullFuzz]{}, err
					}

					return server.Response[NotNullFuzz]{
						Status:     http.StatusOK,
						Success:    true,
						Error:      nil,
						Objects:    []*NotNullFuzz{object},
						Count:      count,
						TotalCount: totalCount,
						Limit:      limit,
						Offset:     offset,
					}, nil
				},
				NotNullFuzz{},
				NotNullFuzzIntrospectedTable,
			)
			if err != nil {
				panic(err)
			}
			r.Post(postHandlerForClaim.FullPath, postHandlerForClaim.ServeHTTP)

			postHandlerForClaimOne, err := getHTTPHandler(
				http.MethodPost,
				"/not-null-fuzzes/{primaryKey}/claim",
				http.StatusOK,
				func(
					ctx context.Context,
					pathParams NotNullFuzzOnePathParams,
					queryParams NotNullFuzzLoadQueryParams,
					req NotNullFuzzClaimRequest,
					rawReq any,
				) (server.Response[NotNullFuzz], error) {
					before := time.Now()

					redisConn := redisPool.Get()
					defer func() {
						_ = redisConn.Close()
					}()

					arguments, err := server.GetSelectOneArguments(ctx, queryParams.Depth, NotNullFuzzIntrospectedTable, pathParams.PrimaryKey, nil, nil)
					if err != nil {
						if config.Debug() {
							log.Printf("request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
						}

						return server.Response[NotNullFuzz]{}, err
					}

					/* note: deliberately no attempt at a cache hit */

					var object *NotNullFuzz
					var count int64
					var totalCount int64

					err = func() error {
						tx, err := db.Begin(arguments.Ctx)
						if err != nil {
							return err
						}

						defer func() {
							_ = tx.Rollback(arguments.Ctx)
						}()

						object, count, totalCount, _, _, err = SelectNotNullFuzz(arguments.Ctx, tx, arguments.Where, arguments.Values...)
						if err != nil {
							return fmt.Errorf("failed to select object to claim: %s", err.Error())
						}

						err = object.Claim(arguments.Ctx, tx, req.Until, req.By, time.Millisecond*time.Duration(req.TimeoutSeconds*1000))
						if err != nil {
							return err
						}

						err = tx.Commit(arguments.Ctx)
						if err != nil {
							return err
						}

						return nil
					}()
					if err != nil {
						if config.Debug() {
							log.Printf("request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
						}

						return server.Response[NotNullFuzz]{}, err
					}

					limit := int64(0)

					offset := int64(0)

					response := server.Response[NotNullFuzz]{
						Status:     http.StatusOK,
						Success:    true,
						Error:      nil,
						Objects:    []*NotNullFuzz{object},
						Count:      count,
						TotalCount: totalCount,
						Limit:      limit,
						Offset:     offset,
					}

					return response, nil
				},
				NotNullFuzz{},
				NotNullFuzzIntrospectedTable,
			)
			if err != nil {
				panic(err)
			}
			r.Post(postHandlerForClaimOne.FullPath, postHandlerForClaimOne.ServeHTTP)
		}()
	}

	func() {
		getManyHandler, err := getHTTPHandler(
			http.MethodGet,
			"/not-null-fuzzes",
			http.StatusOK,
			func(
				ctx context.Context,
				pathParams server.EmptyPathParams,
				queryParams map[string]any,
				req server.EmptyRequest,
				rawReq any,
			) (server.Response[NotNullFuzz], error) {
				before := time.Now()

				redisConn := redisPool.Get()
				defer func() {
					_ = redisConn.Close()
				}()

				arguments, err := server.GetSelectManyArguments(ctx, queryParams, NotNullFuzzIntrospectedTable, nil, nil)
				if err != nil {
					if config.Debug() {
						log.Printf("request cache not yet reached; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return server.Response[NotNullFuzz]{}, err
				}

				cachedResponseAsJSON, cacheHit, err := server.GetCachedResponseAsJSON(arguments.RequestHash, redisConn)
				if err != nil {
					if config.Debug() {
						log.Printf("request cache failed; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return server.Response[NotNullFuzz]{}, err
				}

				if cacheHit {
					var cachedResponse server.Response[NotNullFuzz]

					/* TODO: it'd be nice to be able to avoid this (i.e. just pass straight through) */
					err = json.Unmarshal(cachedResponseAsJSON, &cachedResponse)
					if err != nil {
						if config.Debug() {
							log.Printf("request cache hit but failed unmarshal; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
						}

						return server.Response[NotNullFuzz]{}, err
					}

					if config.Debug() {
						log.Printf("request cache hit; request succeeded in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return cachedResponse, nil
				}

				objects, count, totalCount, _, _, err := handleGetNotNullFuzzes(arguments, db)
				if err != nil {
					if config.Debug() {
						log.Printf("request cache missed; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return server.Response[NotNullFuzz]{}, err
				}

				limit := int64(0)
				if arguments.Limit != nil {
					limit = int64(*arguments.Limit)
				}

				offset := int64(0)
				if arguments.Offset != nil {
					offset = int64(*arguments.Offset)
				}

				response := server.Response[NotNullFuzz]{
					Status:     http.StatusOK,
					Success:    true,
					Error:      nil,
					Objects:    objects,
					Count:      count,
					TotalCount: totalCount,
					Limit:      limit,
					Offset:     offset,
				}

				/* TODO: it'd be nice to be able to avoid this (i.e. just marshal once, further out) */
				responseAsJSON, err := json.Marshal(response)
				if err != nil {
					if config.Debug() {
						log.Printf("request cache missed; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return server.Response[NotNullFuzz]{}, err
				}

				err = server.StoreCachedResponse(arguments.RequestHash, redisConn, responseAsJSON)
				if err != nil {
					log.Printf("warning; %v", err)
				}

				if config.Debug() {
					log.Printf("request cache missed; request succeeded in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
				}

				return response, nil
			},
			NotNullFuzz{},
			NotNullFuzzIntrospectedTable,
		)
		if err != nil {
			panic(err)
		}
		r.Get(getManyHandler.FullPath, getManyHandler.ServeHTTP)
	}()

	func() {
		getOneHandler, err := getHTTPHandler(
			http.MethodGet,
			"/not-null-fuzzes/{primaryKey}",
			http.StatusOK,
			func(
				ctx context.Context,
				pathParams NotNullFuzzOnePathParams,
				queryParams NotNullFuzzLoadQueryParams,
				req server.EmptyRequest,
				rawReq any,
			) (server.Response[NotNullFuzz], error) {
				before := time.Now()

				redisConn := redisPool.Get()
				defer func() {
					_ = redisConn.Close()
				}()

				arguments, err := server.GetSelectOneArguments(ctx, queryParams.Depth, NotNullFuzzIntrospectedTable, pathParams.PrimaryKey, nil, nil)
				if err != nil {
					if config.Debug() {
						log.Printf("request cache not yet reached; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return server.Response[NotNullFuzz]{}, err
				}

				cachedResponseAsJSON, cacheHit, err := server.GetCachedResponseAsJSON(arguments.RequestHash, redisConn)
				if err != nil {
					if config.Debug() {
						log.Printf("request cache failed; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return server.Response[NotNullFuzz]{}, err
				}

				if cacheHit {
					var cachedResponse server.Response[NotNullFuzz]

					/* TODO: it'd be nice to be able to avoid this (i.e. just pass straight through) */
					err = json.Unmarshal(cachedResponseAsJSON, &cachedResponse)
					if err != nil {
						if config.Debug() {
							log.Printf("request cache hit but failed unmarshal; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
						}

						return server.Response[NotNullFuzz]{}, err
					}

					if config.Debug() {
						log.Printf("request cache hit; request succeeded in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return cachedResponse, nil
				}

				objects, count, totalCount, _, _, err := handleGetNotNullFuzz(arguments, db, pathParams.PrimaryKey)
				if err != nil {
					if config.Debug() {
						log.Printf("request cache missed; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return server.Response[NotNullFuzz]{}, err
				}

				limit := int64(0)

				offset := int64(0)

				response := server.Response[NotNullFuzz]{
					Status:     http.StatusOK,
					Success:    true,
					Error:      nil,
					Objects:    objects,
					Count:      count,
					TotalCount: totalCount,
					Limit:      limit,
					Offset:     offset,
				}

				/* TODO: it'd be nice to be able to avoid this (i.e. just marshal once, further out) */
				responseAsJSON, err := json.Marshal(response)
				if err != nil {
					if config.Debug() {
						log.Printf("request cache missed; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return server.Response[NotNullFuzz]{}, err
				}

				err = server.StoreCachedResponse(arguments.RequestHash, redisConn, responseAsJSON)
				if err != nil {
					log.Printf("warning; %v", err)
				}

				if config.Debug() {
					log.Printf("request cache hit; request succeeded in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
				}

				return response, nil
			},
			NotNullFuzz{},
			NotNullFuzzIntrospectedTable,
		)
		if err != nil {
			panic(err)
		}
		r.Get(getOneHandler.FullPath, getOneHandler.ServeHTTP)
	}()

	func() {
		postHandler, err := getHTTPHandler(
			http.MethodPost,
			"/not-null-fuzzes",
			http.StatusCreated,
			func(
				ctx context.Context,
				pathParams server.EmptyPathParams,
				queryParams NotNullFuzzLoadQueryParams,
				req []*NotNullFuzz,
				rawReq any,
			) (server.Response[NotNullFuzz], error) {
				allRawItems, ok := rawReq.([]any)
				if !ok {
					return server.Response[NotNullFuzz]{}, fmt.Errorf("failed to cast %#+v to []map[string]any", rawReq)
				}

				allItems := make([]map[string]any, 0)
				for _, rawItem := range allRawItems {
					item, ok := rawItem.(map[string]any)
					if !ok {
						return server.Response[NotNullFuzz]{}, fmt.Errorf("failed to cast %#+v to map[string]any", rawItem)
					}

					allItems = append(allItems, item)
				}

				forceSetValuesForFieldsByObjectIndex := make([][]string, 0)
				for _, item := range allItems {
					forceSetValuesForFields := make([]string, 0)
					for _, possibleField := range maps.Keys(item) {
						if !slices.Contains(NotNullFuzzTableColumns, possibleField) {
							continue
						}

						forceSetValuesForFields = append(forceSetValuesForFields, possibleField)
					}
					forceSetValuesForFieldsByObjectIndex = append(forceSetValuesForFieldsByObjectIndex, forceSetValuesForFields)
				}

				arguments, err := server.GetLoadArguments(ctx, queryParams.Depth)
				if err != nil {
					return server.Response[NotNullFuzz]{}, err
				}

				objects, count, totalCount, _, _, err := handlePostNotNullFuzzs(arguments, db, waitForChange, req, forceSetValuesForFieldsByObjectIndex)
				if err != nil {
					return server.Response[NotNullFuzz]{}, err
				}

				limit := int64(0)

				offset := int64(0)

				return server.Response[NotNullFuzz]{
					Status:     http.StatusOK,
					Success:    true,
					Error:      nil,
					Objects:    objects,
					Count:      count,
					TotalCount: totalCount,
					Limit:      limit,
					Offset:     offset,
				}, nil
			},
			NotNullFuzz{},
			NotNullFuzzIntrospectedTable,
		)
		if err != nil {
			panic(err)
		}
		r.Post(postHandler.FullPath, postHandler.ServeHTTP)
	}()

	func() {
		putHandler, err := getHTTPHandler(
			http.MethodPatch,
			"/not-null-fuzzes/{primaryKey}",
			http.StatusOK,
			func(
				ctx context.Context,
				pathParams NotNullFuzzOnePathParams,
				queryParams NotNullFuzzLoadQueryParams,
				req NotNullFuzz,
				rawReq any,
			) (server.Response[NotNullFuzz], error) {
				item, ok := rawReq.(map[string]any)
				if !ok {
					return server.Response[NotNullFuzz]{}, fmt.Errorf("failed to cast %#+v to map[string]any", item)
				}

				arguments, err := server.GetLoadArguments(ctx, queryParams.Depth)
				if err != nil {
					return server.Response[NotNullFuzz]{}, err
				}

				object := &req
				object.MrPrimary = pathParams.PrimaryKey

				objects, count, totalCount, _, _, err := handlePutNotNullFuzz(arguments, db, waitForChange, object)
				if err != nil {
					return server.Response[NotNullFuzz]{}, err
				}

				limit := int64(0)

				offset := int64(0)

				return server.Response[NotNullFuzz]{
					Status:     http.StatusOK,
					Success:    true,
					Error:      nil,
					Objects:    objects,
					Count:      count,
					TotalCount: totalCount,
					Limit:      limit,
					Offset:     offset,
				}, nil
			},
			NotNullFuzz{},
			NotNullFuzzIntrospectedTable,
		)
		if err != nil {
			panic(err)
		}
		r.Put(putHandler.FullPath, putHandler.ServeHTTP)
	}()

	func() {
		patchHandler, err := getHTTPHandler(
			http.MethodPatch,
			"/not-null-fuzzes/{primaryKey}",
			http.StatusOK,
			func(
				ctx context.Context,
				pathParams NotNullFuzzOnePathParams,
				queryParams NotNullFuzzLoadQueryParams,
				req NotNullFuzz,
				rawReq any,
			) (server.Response[NotNullFuzz], error) {
				item, ok := rawReq.(map[string]any)
				if !ok {
					return server.Response[NotNullFuzz]{}, fmt.Errorf("failed to cast %#+v to map[string]any", item)
				}

				forceSetValuesForFields := make([]string, 0)
				for _, possibleField := range maps.Keys(item) {
					if !slices.Contains(NotNullFuzzTableColumns, possibleField) {
						continue
					}

					forceSetValuesForFields = append(forceSetValuesForFields, possibleField)
				}

				arguments, err := server.GetLoadArguments(ctx, queryParams.Depth)
				if err != nil {
					return server.Response[NotNullFuzz]{}, err
				}

				object := &req
				object.MrPrimary = pathParams.PrimaryKey

				objects, count, totalCount, _, _, err := handlePatchNotNullFuzz(arguments, db, waitForChange, object, forceSetValuesForFields)
				if err != nil {
					return server.Response[NotNullFuzz]{}, err
				}

				limit := int64(0)

				offset := int64(0)

				return server.Response[NotNullFuzz]{
					Status:     http.StatusOK,
					Success:    true,
					Error:      nil,
					Objects:    objects,
					Count:      count,
					TotalCount: totalCount,
					Limit:      limit,
					Offset:     offset,
				}, nil
			},
			NotNullFuzz{},
			NotNullFuzzIntrospectedTable,
		)
		if err != nil {
			panic(err)
		}
		r.Patch(patchHandler.FullPath, patchHandler.ServeHTTP)
	}()

	func() {
		deleteHandler, err := getHTTPHandler(
			http.MethodDelete,
			"/not-null-fuzzes/{primaryKey}",
			http.StatusNoContent,
			func(
				ctx context.Context,
				pathParams NotNullFuzzOnePathParams,
				queryParams NotNullFuzzLoadQueryParams,
				req server.EmptyRequest,
				rawReq any,
			) (server.EmptyResponse, error) {
				arguments, err := server.GetLoadArguments(ctx, queryParams.Depth)
				if err != nil {
					return server.EmptyResponse{}, err
				}

				object := &NotNullFuzz{}
				object.MrPrimary = pathParams.PrimaryKey

				err = handleDeleteNotNullFuzz(arguments, db, waitForChange, object)
				if err != nil {
					return server.EmptyResponse{}, err
				}

				return server.EmptyResponse{}, nil
			},
			NotNullFuzz{},
			NotNullFuzzIntrospectedTable,
		)
		if err != nil {
			panic(err)
		}
		r.Delete(deleteHandler.FullPath, deleteHandler.ServeHTTP)
	}()
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
		MutateRouterForNotNullFuzz,
	)
}
