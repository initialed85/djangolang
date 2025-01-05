package model_generated_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"net/netip"
	"sync"
	"testing"
	"time"

	"github.com/cridenour/go-postgis"
	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
	"github.com/initialed85/djangolang/internal/hack"
	"github.com/initialed85/djangolang/pkg/helpers"
	"github.com/initialed85/djangolang/pkg/model_generated"
	"github.com/initialed85/djangolang/pkg/server"
	"github.com/initialed85/djangolang/pkg/stream"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

type CollectMrPrimariesResponse struct {
	MrPrimaries []int `json:"mr_primaries"`
}

func testNotNullFuzz(
	t *testing.T,
	ctx context.Context,
	db *pgxpool.Pool,
	redisConn redis.Conn,
	mu *sync.Mutex,
	lastChangeByTableName map[string]*server.Change,
	httpClient *HTTPClient,
	getLastChangeForTableName func(tableName string) *server.Change,
) {
	cleanup := func() {
		_, _ = db.Exec(
			ctx,
			`TRUNCATE TABLE not_null_fuzz`,
		)
		if redisConn != nil {
			_, _ = redisConn.Do("FLUSHALL")
		}
	}
	defer cleanup()

	prefix, err := netip.ParsePrefix("192.168.137.222/24")
	require.NoError(t, err)

	_, err = db.Exec(
		ctx,
		`INSERT INTO not_null_fuzz DEFAULT VALUES;`,
	)
	require.NoError(t, err)

	t.Run("ObjectInteractions", func(t *testing.T) {
		tx, err := db.Begin(ctx)
		require.NoError(t, err)
		defer func() {
			_ = tx.Rollback(ctx)
		}()

		notNullFuzzes, count, totalCount, page, totalPages, err := model_generated.SelectNotNullFuzzes(ctx, tx, "", nil, nil, nil)
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(notNullFuzzes), 1)
		require.GreaterOrEqual(t, count, int64(1))
		require.GreaterOrEqual(t, totalCount, int64(10))
		require.Equal(t, int64(1), page)
		require.Equal(t, int64(1), totalPages)

		notNullFuzz := notNullFuzzes[0]
		err = notNullFuzz.Update(ctx, tx, false)
		require.NoError(t, err)

		err = tx.Commit(ctx)
		require.NoError(t, err)

		require.Equal(t, int64(1), notNullFuzz.SomeBigint, "SomeBigint")
		require.Equal(t, []int64{1}, notNullFuzz.SomeBigintArray, "SomeBigintArray")
		require.Equal(t, true, notNullFuzz.SomeBoolean, "SomeBoolean")
		require.Equal(t, []bool{true}, notNullFuzz.SomeBooleanArray, "SomeBooleanArray")
		require.Equal(t, []byte{0x65}, notNullFuzz.SomeBytea, "SomeBytea")
		require.Equal(t, "A", notNullFuzz.SomeCharacterVarying, "SomeCharacterVarying")
		require.Equal(t, []string{"A"}, notNullFuzz.SomeCharacterVaryingArray, "SomeCharacterVaryingArray")
		require.Equal(t, float64(1.0), notNullFuzz.SomeDoublePrecision, "SomeDoublePrecision")
		require.Equal(t, []float64{1.0}, notNullFuzz.SomeDoublePrecisionArray, "SomeDoublePrecisionArray")
		require.Equal(t, float64(1.0), notNullFuzz.SomeFloat, "SomeFloat")
		require.Equal(t, []float64{1.0}, notNullFuzz.SomeFloatArray, "SomeFloatArray")
		require.Equal(t, postgis.PointZ{X: 1.337, Y: 69.42, Z: 800.8135}, notNullFuzz.SomeGeometryPointZ, "SomeGeometryPointZ")
		require.Equal(t, map[string]*string{"A": helpers.Ptr("1")}, notNullFuzz.SomeHstore, "SomeHstore")
		require.Equal(t, prefix, notNullFuzz.SomeInet, "SomeInet")
		require.Equal(t, int64(1), notNullFuzz.SomeInteger, "SomeInteger")
		require.Equal(t, []int64{1}, notNullFuzz.SomeIntegerArray, "SomeIntegerArray")
		require.Equal(t, time.Millisecond*1337, notNullFuzz.SomeInterval, "SomeInterval")
		require.Equal(t, map[string]interface{}{"some": "data"}, notNullFuzz.SomeJSON, "SomeJSON")
		require.Equal(t, map[string]interface{}{"some": "data"}, notNullFuzz.SomeJSONB, "SomeJSONB")
		require.Equal(t, float64(1.0), notNullFuzz.SomeNumeric, "SomeNumeric")
		require.Equal(t, []float64{1.0}, notNullFuzz.SomeNumericArray, "SomeNumericArray")
		require.Equal(t, pgtype.Vec2{X: 1.337, Y: 69.42}, notNullFuzz.SomePoint, "SomePoint")
		require.Equal(t, []pgtype.Vec2{{X: 75, Y: 29}, {X: 77, Y: 29}, {X: 77, Y: 29}, {X: 75, Y: 29}}, notNullFuzz.SomePolygon, "SomePolygon")
		require.Equal(t, float64(1.0), notNullFuzz.SomeReal, "SomeReal")
		require.Equal(t, []float64{1.0}, notNullFuzz.SomeRealArray, "SomeRealArray")
		require.Equal(t, int64(1), notNullFuzz.SomeSmallint, "SomeSmallint")
		require.Equal(t, []int64{1}, notNullFuzz.SomeSmallintArray, "SomeSmallintArray")
		require.Equal(t, "A", notNullFuzz.SomeText, "SomeText")
		require.Equal(t, []string{"A"}, notNullFuzz.SomeTextArray, "SomeTextArray")
		require.Equal(t, "2024-07-19 11:45:00 +0800 AWST", notNullFuzz.SomeTimestamptz.String(), "SomeTimestamptz")
		require.Equal(t, "2020-03-27 08:30:00 +0000 UTC", notNullFuzz.SomeTimestamp.String(), "SomeTimestamp")
		require.Equal(t, map[string][]int(map[string][]int{"a": []int(nil)}), notNullFuzz.SomeTsvector, "SomeTsvector")
		require.Equal(t, uuid.UUID{0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11}, notNullFuzz.SomeUUID, "SomeUUID")

		require.Eventually(
			t,
			func() bool {
				return getLastChangeForTableName(model_generated.NotNullFuzzTable) != nil
			},
			time.Second*30,
			time.Millisecond*10,
		)

		lastChange := getLastChangeForTableName(model_generated.NotNullFuzzTable)
		require.NotNil(t, lastChange)

		object, ok := lastChange.Object.(*model_generated.NotNullFuzz)
		require.True(t, ok)

		require.Equal(t, int64(1), object.SomeBigint, "SomeBigint")
		require.Equal(t, []int64{1}, object.SomeBigintArray, "SomeBigintArray")
		require.Equal(t, true, object.SomeBoolean, "SomeBoolean")
		require.Equal(t, []bool{true}, object.SomeBooleanArray, "SomeBooleanArray")
		require.Equal(t, []byte{0x65}, object.SomeBytea, "SomeBytea")
		require.Equal(t, "A", object.SomeCharacterVarying, "SomeCharacterVarying")
		require.Equal(t, []string{"A"}, object.SomeCharacterVaryingArray, "SomeCharacterVaryingArray")
		require.Equal(t, float64(1.0), object.SomeDoublePrecision, "SomeDoublePrecision")
		require.Equal(t, []float64{1.0}, object.SomeDoublePrecisionArray, "SomeDoublePrecisionArray")
		require.Equal(t, float64(1.0), object.SomeFloat, "SomeFloat")
		require.Equal(t, []float64{1.0}, object.SomeFloatArray, "SomeFloatArray")
		require.Equal(t, postgis.PointZ{X: 1.337, Y: 69.42, Z: 800.8135}, object.SomeGeometryPointZ, "SomeGeometryPointZ")
		require.Equal(t, map[string]*string{"A": helpers.Ptr("1")}, object.SomeHstore, "SomeHstore")
		require.Equal(t, prefix, object.SomeInet, "SomeInet")
		require.Equal(t, int64(1), object.SomeInteger, "SomeInteger")
		require.Equal(t, []int64{1}, object.SomeIntegerArray, "SomeIntegerArray")
		require.Equal(t, time.Millisecond*1337, object.SomeInterval, "SomeInterval")
		require.Equal(t, map[string]interface{}{"some": "data"}, object.SomeJSON, "SomeJSON")
		require.Equal(t, map[string]interface{}{"some": "data"}, object.SomeJSONB, "SomeJSONB")
		require.Equal(t, float64(1.0), object.SomeNumeric, "SomeNumeric")
		require.Equal(t, []float64{1.0}, object.SomeNumericArray, "SomeNumericArray")
		require.Equal(t, pgtype.Vec2{X: 1.337, Y: 69.42}, object.SomePoint, "SomePoint")
		require.Equal(t, []pgtype.Vec2{{X: 75, Y: 29}, {X: 77, Y: 29}, {X: 77, Y: 29}, {X: 75, Y: 29}}, object.SomePolygon, "SomePolygon")
		require.Equal(t, float64(1.0), object.SomeReal, "SomeReal")
		require.Equal(t, []float64{1.0}, object.SomeRealArray, "SomeRealArray")
		require.Equal(t, int64(1), object.SomeSmallint, "SomeSmallint")
		require.Equal(t, []int64{1}, object.SomeSmallintArray, "SomeSmallintArray")
		require.Equal(t, "A", object.SomeText, "SomeText")
		require.Equal(t, []string{"A"}, object.SomeTextArray, "SomeTextArray")
		// if you .Reload() you'll get this, otherwise you'll get the UTC one below it
		// require.Equal(t, "2024-07-19 11:45:00 +0800 AWST", object.SomeTimestamptz.String(), "SomeTimestamptz")
		require.Equal(t, "2024-07-19 03:45:00 +0000 +0000", object.SomeTimestamptz.String(), "SomeTimestamptz")
		require.Equal(t, "2020-03-27 08:30:00 +0000 UTC", object.SomeTimestamp.String(), "SomeTimestamp")
		require.Equal(t, map[string][]int(map[string][]int{"a": []int(nil)}), object.SomeTsvector, "SomeTsvector")
		require.Equal(t, uuid.UUID{0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11}, object.SomeUUID, "SomeUUID")

		require.Equal(t, int64(1), lastChange.Item["some_bigint"], "SomeBigint")
		require.Equal(t, []any{int64(1)}, lastChange.Item["some_bigint_array"], "SomeBigintArray")
		require.Equal(t, true, lastChange.Item["some_boolean"], "SomeBoolean")
		require.Equal(t, []any{true}, lastChange.Item["some_boolean_array"], "SomeBooleanArray")
		require.Equal(t, []byte{0x5a, 0x51, 0x3d, 0x3d}, lastChange.Item["some_bytea"], "SomeBytea")
		require.Equal(t, "A", lastChange.Item["some_character_varying"], "SomeCharacterVarying")
		require.Equal(t, []any{"A"}, lastChange.Item["some_character_varying_array"], "SomeCharacterVaryingArray")
		require.Equal(t, float64(1.0), lastChange.Item["some_double_precision"], "SomeDoublePrecision")
		require.Equal(t, []any{1.0}, lastChange.Item["some_double_precision_array"], "SomeDoublePrecisionArray")
		require.Equal(t, float64(1.0), lastChange.Item["some_float"], "SomeFloat")
		require.Equal(t, []any{1.0}, lastChange.Item["some_float_array"], "SomeFloatArray")
		require.Equal(t, "01010000803108AC1C5A64F53F7B14AE47E15A51405EBA490C82068940", lastChange.Item["some_geometry_point_z"], "SomeGeometryPointZ")
		require.Equal(t, "\"A\"=>\"1\"", lastChange.Item["some_hstore"], "SomeHstore")
		require.Equal(t, prefix, lastChange.Item["some_inet"], "SomeInet")
		require.Equal(t, int32(1), lastChange.Item["some_integer"], "SomeInteger")
		require.Equal(t, []any{int32(1)}, lastChange.Item["some_integer_array"], "SomeIntegerArray")
		require.Equal(t, pgtype.Interval{Microseconds: 1337000, Days: 0, Months: 0, Valid: true}, lastChange.Item["some_interval"], "SomeInterval")
		require.Equal(t, map[string]interface{}{"some": "data"}, lastChange.Item["some_json"], "SomeJSON")
		require.Equal(t, map[string]interface{}{"some": "data"}, lastChange.Item["some_jsonb"], "SomeJSONB")
		require.Equal(t, pgtype.Numeric{Int: big.NewInt(1), Exp: 0, NaN: false, InfinityModifier: 0, Valid: true}, lastChange.Item["some_numeric"], "SomeNumeric")
		require.Equal(t, []interface{}{pgtype.Numeric{Int: big.NewInt(1), Exp: 0, NaN: false, InfinityModifier: 0, Valid: true}}, lastChange.Item["some_numeric_array"], "SomeNumericArray")
		require.Equal(t, pgtype.Point{P: pgtype.Vec2{X: 1.337, Y: 69.42}, Valid: true}, lastChange.Item["some_point"], "SomePoint")
		require.Equal(t, pgtype.Polygon{P: []pgtype.Vec2{{X: 75, Y: 29}, {X: 77, Y: 29}, {X: 77, Y: 29}, {X: 75, Y: 29}}, Valid: true}, lastChange.Item["some_polygon"], "SomePolygon")
		require.Equal(t, float32(1.0), lastChange.Item["some_real"], "SomeReal")
		require.Equal(t, []any{float32(1.0)}, lastChange.Item["some_real_array"], "SomeRealArray")
		require.Equal(t, int16(1), lastChange.Item["some_smallint"], "SomeSmallint")
		require.Equal(t, []any{int16(1)}, lastChange.Item["some_smallint_array"], "SomeSmallintArray")
		require.Equal(t, "A", lastChange.Item["some_text"], "SomeText")
		require.Equal(t, []any{"A"}, lastChange.Item["some_text_array"], "SomeTextArray")
		require.Equal(t, "2024-07-19 03:45:00 +0000 +0000", lastChange.Item["some_timestamptz"].(time.Time).String(), "SomeTimestamptz")
		require.Equal(t, "2020-03-27 08:30:00 +0000 UTC", lastChange.Item["some_timestamp"].(time.Time).String(), "SomeTimestamp")
		require.Equal(t, "'a'", lastChange.Item["some_tsvector"], "SomeTsvector")
		require.Equal(t, [16]byte{0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11}, lastChange.Item["some_uuid"], "SomeUUID")

		tx, err = db.Begin(ctx)
		require.NoError(t, err)
		defer func() {
			_ = tx.Rollback(ctx)
		}()

		notNullFuzzes, count, totalCount, page, totalPages, err = model_generated.SelectNotNullFuzzes(ctx, tx, "", nil, nil, nil)
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(notNullFuzzes), 1)
		require.GreaterOrEqual(t, int64(count), int64(1))
		require.GreaterOrEqual(t, int64(totalCount), int64(1))
		require.Equal(t, int64(1), int64(page))
		require.Equal(t, int64(1), int64(totalPages))

		err = tx.Commit(ctx)
		require.NoError(t, err)

		notNullFuzz = notNullFuzzes[0]

		tx, err = db.Begin(ctx)
		require.NoError(t, err)
		defer func() {
			_ = tx.Rollback(ctx)
		}()

		notNullFuzz.SomeBigint = 2
		notNullFuzz.SomeBigintArray = append(notNullFuzz.SomeBigintArray, 2)

		err = notNullFuzz.Insert(ctx, tx, false, false)
		require.NoError(t, err)

		err = tx.Commit(ctx)
		require.NoError(t, err)

		require.Equal(t, int64(2), notNullFuzz.SomeBigint, "SomeBigint")
		require.Equal(t, []int64{1, 2}, notNullFuzz.SomeBigintArray, "SomeBigintArray")
		require.Equal(t, true, notNullFuzz.SomeBoolean, "SomeBoolean")
		require.Equal(t, []bool{true}, notNullFuzz.SomeBooleanArray, "SomeBooleanArray")
		require.Equal(t, []byte{0x65}, notNullFuzz.SomeBytea, "SomeBytea")
		require.Equal(t, "A", notNullFuzz.SomeCharacterVarying, "SomeCharacterVarying")
		require.Equal(t, []string{"A"}, notNullFuzz.SomeCharacterVaryingArray, "SomeCharacterVaryingArray")
		require.Equal(t, float64(1.0), notNullFuzz.SomeDoublePrecision, "SomeDoublePrecision")
		require.Equal(t, []float64{1.0}, notNullFuzz.SomeDoublePrecisionArray, "SomeDoublePrecisionArray")
		require.Equal(t, float64(1.0), notNullFuzz.SomeFloat, "SomeFloat")
		require.Equal(t, []float64{1.0}, notNullFuzz.SomeFloatArray, "SomeFloatArray")
		require.Equal(t, postgis.PointZ{X: 1.337, Y: 69.42, Z: 800.8135}, notNullFuzz.SomeGeometryPointZ, "SomeGeometryPointZ")
		require.Equal(t, map[string]*string{"A": helpers.Ptr("1")}, notNullFuzz.SomeHstore, "SomeHstore")
		require.Equal(t, prefix, notNullFuzz.SomeInet, "SomeInet")
		require.Equal(t, int64(1), notNullFuzz.SomeInteger, "SomeInteger")
		require.Equal(t, []int64{1}, notNullFuzz.SomeIntegerArray, "SomeIntegerArray")
		require.Equal(t, time.Millisecond*1337, notNullFuzz.SomeInterval, "SomeInterval")
		require.Equal(t, map[string]interface{}{"some": "data"}, notNullFuzz.SomeJSON, "SomeJSON")
		require.Equal(t, map[string]interface{}{"some": "data"}, notNullFuzz.SomeJSONB, "SomeJSONB")
		require.Equal(t, float64(1.0), notNullFuzz.SomeNumeric, "SomeNumeric")
		require.Equal(t, []float64{1.0}, notNullFuzz.SomeNumericArray, "SomeNumericArray")
		require.Equal(t, pgtype.Vec2{X: 1.337, Y: 69.42}, notNullFuzz.SomePoint, "SomePoint")
		require.Equal(t, []pgtype.Vec2{{X: 75, Y: 29}, {X: 77, Y: 29}, {X: 77, Y: 29}, {X: 75, Y: 29}}, notNullFuzz.SomePolygon, "SomePolygon")
		require.Equal(t, float64(1.0), notNullFuzz.SomeReal, "SomeReal")
		require.Equal(t, []float64{1.0}, notNullFuzz.SomeRealArray, "SomeRealArray")
		require.Equal(t, int64(1), notNullFuzz.SomeSmallint, "SomeSmallint")
		require.Equal(t, []int64{1}, notNullFuzz.SomeSmallintArray, "SomeSmallintArray")
		require.Equal(t, "A", notNullFuzz.SomeText, "SomeText")
		require.Equal(t, []string{"A"}, notNullFuzz.SomeTextArray, "SomeTextArray")
		require.Equal(t, "2024-07-19 11:45:00 +0800 AWST", notNullFuzz.SomeTimestamptz.String(), "SomeTimestamptz")
		require.Equal(t, "2020-03-27 08:30:00 +0000 UTC", notNullFuzz.SomeTimestamp.String(), "SomeTimestamp")
		require.Equal(t, map[string][]int(map[string][]int{"a": []int(nil)}), notNullFuzz.SomeTsvector, "SomeTsvector")
		require.Equal(t, uuid.UUID{0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11}, notNullFuzz.SomeUUID, "SomeUUID")

		tx, err = db.Begin(ctx)
		require.NoError(t, err)
		defer func() {
			_ = tx.Rollback(ctx)
		}()

		notNullFuzzes, count, totalCount, page, totalPages, err = model_generated.SelectNotNullFuzzes(ctx, tx, fmt.Sprintf("%s = $$??", model_generated.NotNullFuzzTableSomeBigintColumn), nil, nil, nil, 2)
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(notNullFuzzes), 1)
		require.GreaterOrEqual(t, int64(count), int64(1))
		require.GreaterOrEqual(t, int64(totalCount), int64(1))
		require.Equal(t, int64(1), int64(page))
		require.Equal(t, int64(1), int64(totalPages))

		err = tx.Commit(ctx)
		require.NoError(t, err)

		notNullFuzz = notNullFuzzes[len(notNullFuzzes)-1]

		require.Equal(t, int64(2), notNullFuzz.SomeBigint, "SomeBigint")
		require.Equal(t, []int64{1, 2}, notNullFuzz.SomeBigintArray, "SomeBigintArray")
		require.Equal(t, true, notNullFuzz.SomeBoolean, "SomeBoolean")
		require.Equal(t, []bool{true}, notNullFuzz.SomeBooleanArray, "SomeBooleanArray")
		require.Equal(t, []byte{0x65}, notNullFuzz.SomeBytea, "SomeBytea")
		require.Equal(t, "A", notNullFuzz.SomeCharacterVarying, "SomeCharacterVarying")
		require.Equal(t, []string{"A"}, notNullFuzz.SomeCharacterVaryingArray, "SomeCharacterVaryingArray")
		require.Equal(t, float64(1.0), notNullFuzz.SomeDoublePrecision, "SomeDoublePrecision")
		require.Equal(t, []float64{1.0}, notNullFuzz.SomeDoublePrecisionArray, "SomeDoublePrecisionArray")
		require.Equal(t, float64(1.0), notNullFuzz.SomeFloat, "SomeFloat")
		require.Equal(t, []float64{1.0}, notNullFuzz.SomeFloatArray, "SomeFloatArray")
		require.Equal(t, postgis.PointZ{X: 1.337, Y: 69.42, Z: 800.8135}, notNullFuzz.SomeGeometryPointZ, "SomeGeometryPointZ")
		require.Equal(t, map[string]*string{"A": helpers.Ptr("1")}, notNullFuzz.SomeHstore, "SomeHstore")
		require.Equal(t, prefix, notNullFuzz.SomeInet, "SomeInet")
		require.Equal(t, int64(1), notNullFuzz.SomeInteger, "SomeInteger")
		require.Equal(t, []int64{1}, notNullFuzz.SomeIntegerArray, "SomeIntegerArray")
		require.Equal(t, time.Millisecond*1337, notNullFuzz.SomeInterval, "SomeInterval")
		require.Equal(t, map[string]interface{}{"some": "data"}, notNullFuzz.SomeJSON, "SomeJSON")
		require.Equal(t, map[string]interface{}{"some": "data"}, notNullFuzz.SomeJSONB, "SomeJSONB")
		require.Equal(t, float64(1.0), notNullFuzz.SomeNumeric, "SomeNumeric")
		require.Equal(t, []float64{1.0}, notNullFuzz.SomeNumericArray, "SomeNumericArray")
		require.Equal(t, pgtype.Vec2{X: 1.337, Y: 69.42}, notNullFuzz.SomePoint, "SomePoint")
		require.Equal(t, []pgtype.Vec2{{X: 75, Y: 29}, {X: 77, Y: 29}, {X: 77, Y: 29}, {X: 75, Y: 29}}, notNullFuzz.SomePolygon, "SomePolygon")
		require.Equal(t, float64(1.0), notNullFuzz.SomeReal, "SomeReal")
		require.Equal(t, []float64{1.0}, notNullFuzz.SomeRealArray, "SomeRealArray")
		require.Equal(t, int64(1), notNullFuzz.SomeSmallint, "SomeSmallint")
		require.Equal(t, []int64{1}, notNullFuzz.SomeSmallintArray, "SomeSmallintArray")
		require.Equal(t, "A", notNullFuzz.SomeText, "SomeText")
		require.Equal(t, []string{"A"}, notNullFuzz.SomeTextArray, "SomeTextArray")
		require.Equal(t, "2024-07-19 11:45:00 +0800 AWST", notNullFuzz.SomeTimestamptz.String(), "SomeTimestamptz")
		require.Equal(t, "2020-03-27 08:30:00 +0000 UTC", notNullFuzz.SomeTimestamp.String(), "SomeTimestamp")
		require.Equal(t, map[string][]int(map[string][]int{"a": []int(nil)}), notNullFuzz.SomeTsvector, "SomeTsvector")
		require.Equal(t, uuid.UUID{0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11}, notNullFuzz.SomeUUID, "SomeUUID")

		tx, err = db.Begin(ctx)
		require.NoError(t, err)
		defer func() {
			_ = tx.Rollback(ctx)
		}()

		notNullFuzz.SomeCharacterVarying = "B"
		err = notNullFuzz.Update(ctx, tx, false)
		require.NoError(t, err)

		err = tx.Commit(ctx)
		require.NoError(t, err)

		require.Equal(t, int64(2), notNullFuzz.SomeBigint, "SomeBigint")
		require.Equal(t, []int64{1, 2}, notNullFuzz.SomeBigintArray, "SomeBigintArray")
		require.Equal(t, true, notNullFuzz.SomeBoolean, "SomeBoolean")
		require.Equal(t, []bool{true}, notNullFuzz.SomeBooleanArray, "SomeBooleanArray")
		require.Equal(t, []byte{0x65}, notNullFuzz.SomeBytea, "SomeBytea")
		require.Equal(t, "B", notNullFuzz.SomeCharacterVarying, "SomeCharacterVarying")
		require.Equal(t, []string{"A"}, notNullFuzz.SomeCharacterVaryingArray, "SomeCharacterVaryingArray")
		require.Equal(t, float64(1.0), notNullFuzz.SomeDoublePrecision, "SomeDoublePrecision")
		require.Equal(t, []float64{1.0}, notNullFuzz.SomeDoublePrecisionArray, "SomeDoublePrecisionArray")
		require.Equal(t, float64(1.0), notNullFuzz.SomeFloat, "SomeFloat")
		require.Equal(t, []float64{1.0}, notNullFuzz.SomeFloatArray, "SomeFloatArray")
		require.Equal(t, postgis.PointZ{X: 1.337, Y: 69.42, Z: 800.8135}, notNullFuzz.SomeGeometryPointZ, "SomeGeometryPointZ")
		require.Equal(t, map[string]*string{"A": helpers.Ptr("1")}, notNullFuzz.SomeHstore, "SomeHstore")
		require.Equal(t, prefix, notNullFuzz.SomeInet, "SomeInet")
		require.Equal(t, int64(1), notNullFuzz.SomeInteger, "SomeInteger")
		require.Equal(t, []int64{1}, notNullFuzz.SomeIntegerArray, "SomeIntegerArray")
		require.Equal(t, time.Millisecond*1337, notNullFuzz.SomeInterval, "SomeInterval")
		require.Equal(t, map[string]interface{}{"some": "data"}, notNullFuzz.SomeJSON, "SomeJSON")
		require.Equal(t, map[string]interface{}{"some": "data"}, notNullFuzz.SomeJSONB, "SomeJSONB")
		require.Equal(t, float64(1.0), notNullFuzz.SomeNumeric, "SomeNumeric")
		require.Equal(t, []float64{1.0}, notNullFuzz.SomeNumericArray, "SomeNumericArray")
		require.Equal(t, pgtype.Vec2{X: 1.337, Y: 69.42}, notNullFuzz.SomePoint, "SomePoint")
		require.Equal(t, []pgtype.Vec2{{X: 75, Y: 29}, {X: 77, Y: 29}, {X: 77, Y: 29}, {X: 75, Y: 29}}, notNullFuzz.SomePolygon, "SomePolygon")
		require.Equal(t, float64(1.0), notNullFuzz.SomeReal, "SomeReal")
		require.Equal(t, []float64{1.0}, notNullFuzz.SomeRealArray, "SomeRealArray")
		require.Equal(t, int64(1), notNullFuzz.SomeSmallint, "SomeSmallint")
		require.Equal(t, []int64{1}, notNullFuzz.SomeSmallintArray, "SomeSmallintArray")
		require.Equal(t, "A", notNullFuzz.SomeText, "SomeText")
		require.Equal(t, []string{"A"}, notNullFuzz.SomeTextArray, "SomeTextArray")
		require.Equal(t, "2024-07-19 11:45:00 +0800 AWST", notNullFuzz.SomeTimestamptz.String(), "SomeTimestamptz")
		require.Equal(t, "2020-03-27 08:30:00 +0000 UTC", notNullFuzz.SomeTimestamp.String(), "SomeTimestamp")
		require.Equal(t, map[string][]int(map[string][]int{"a": []int(nil)}), notNullFuzz.SomeTsvector, "SomeTsvector")
		require.Equal(t, uuid.UUID{0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11}, notNullFuzz.SomeUUID, "SomeUUID")

		tx, err = db.Begin(ctx)
		require.NoError(t, err)
		defer func() {
			_ = tx.Rollback(ctx)
		}()

		err = notNullFuzz.Delete(ctx, tx)
		require.NoError(t, err)

		err = tx.Commit(ctx)
		require.NoError(t, err)

		require.Eventually(
			t,
			func() bool {
				return getLastChangeForTableName(model_generated.NotNullFuzzTable) != nil
			},
			time.Second*30,
			time.Millisecond*10,
		)

		lastChange = getLastChangeForTableName(model_generated.NotNullFuzzTable)
		require.NotNil(t, lastChange)
		require.Equal(t, stream.DELETE, lastChange.Action)
	})

	t.Run("EndpointInteractions", func(t *testing.T) {
		resp, err := httpClient.Get("http://localhost:5050/not-null-fuzzes?some_bigint__eq=1")
		require.NoError(t, err)
		respBody, _ := io.ReadAll(resp.Body)
		require.NoError(t, err, string(respBody))
		require.Equal(t, http.StatusOK, resp.StatusCode, string(respBody))

		var rawItems any
		err = json.Unmarshal(respBody, &rawItems)
		require.NoError(t, err)

		items := rawItems.(map[string]interface{})
		objects := items["objects"].([]interface{})
		object := objects[0].(map[string]interface{})
		delete(object, "id")

		expectedObjectIsh := map[string]interface{}{
			"some_bigint":                  1.0,
			"some_bigint_array":            []interface{}{1.0},
			"some_boolean":                 true,
			"some_boolean_array":           []interface{}{true},
			"some_bytea":                   "ZQ==",
			"some_character_varying":       "A",
			"some_character_varying_array": []interface{}{"A"},
			"some_double_precision":        1.0,
			"some_double_precision_array":  []interface{}{1.0},
			"some_float":                   1.0,
			"some_float_array":             []interface{}{1.0},
			"some_geometry_point_z": map[string]interface{}{
				"X": 1.337,
				"Y": 69.42,
				"Z": 800.8135,
			},
			"some_hstore": map[string]interface{}{
				"A": "1",
			},
			"some_inet":          "192.168.137.222/24",
			"some_integer":       1.0,
			"some_integer_array": []interface{}{1.0},
			"some_interval":      1.337e+09,
			"some_json": map[string]interface{}{
				"some": "data",
			},
			"some_jsonb": map[string]interface{}{
				"some": "data",
			},
			"some_numeric":       1.0,
			"some_numeric_array": []interface{}{1.0},
			"some_point":         map[string]interface{}{"X": 1.337, "Y": 69.42},
			"some_polygon": []interface{}{
				map[string]interface{}{
					"X": 75.0, "Y": 29.0,
				},
				map[string]interface{}{
					"X": 77.0, "Y": 29.0,
				},
				map[string]interface{}{
					"X": 77.0, "Y": 29.0,
				},
				map[string]interface{}{
					"X": 75.0, "Y": 29.0},
			},
			"some_real":           1.0,
			"some_real_array":     []interface{}{1.0},
			"some_smallint":       1.0,
			"some_smallint_array": []interface{}{1.0},
			"some_text":           "A",
			"some_text_array":     []interface{}{"A"},
			"some_timestamp":      "2020-03-27T08:30:00Z",
			"some_timestamptz":    "2024-07-19T11:45:00+08:00",
			"some_tsvector":       map[string]interface{}{"a": interface{}(nil)},
			"some_uuid":           "11111111-1111-1111-1111-111111111111",
		}

		for k, v := range expectedObjectIsh {
			require.Equal(t, v, object[k], k)
		}

		expectedObjectIsh["some_bigint"] = float64(2)

		reqBody := []any{expectedObjectIsh}
		b, err := json.MarshalIndent(reqBody, "", "  ")
		require.NoError(t, err)

		resp, err = httpClient.Post("http://localhost:5050/not-null-fuzzes", "application/json", bytes.NewReader(b))
		require.NoError(t, err)
		respBody, _ = io.ReadAll(resp.Body)
		require.NoError(t, err, string(respBody))
		require.Equal(t, http.StatusCreated, resp.StatusCode, string(respBody))

		err = json.Unmarshal(respBody, &rawItems)
		require.NoError(t, err)

		items = rawItems.(map[string]interface{})
		objects = items["objects"].([]interface{})
		object = objects[0].(map[string]interface{})

		for k, v := range expectedObjectIsh {
			require.Equal(t, v, object[k], k)
		}

		b, err = json.MarshalIndent(expectedObjectIsh, "", "  ")
		require.NoError(t, err)

		resp, err = httpClient.Put(fmt.Sprintf("http://localhost:5050/not-null-fuzzes/%v", object["mr_primary"]), "application/json", bytes.NewReader(b))
		require.NoError(t, err)
		respBody, _ = io.ReadAll(resp.Body)
		require.NoError(t, err, string(respBody))
		require.Equal(t, http.StatusOK, resp.StatusCode, string(respBody))

		err = json.Unmarshal(respBody, &rawItems)
		require.NoError(t, err)

		items = rawItems.(map[string]interface{})
		objects, ok := items["objects"].([]interface{})
		require.True(t, ok)
		object = objects[0].(map[string]interface{})
		for k, v := range expectedObjectIsh {
			require.Equal(t, v, object[k], k)
		}

		resp, err = httpClient.Get("http://localhost:5050/custom/collect-mr-primaries")
		require.NoError(t, err)
		respBody, _ = io.ReadAll(resp.Body)
		require.NoError(t, err, string(respBody))
		require.Equal(t, http.StatusOK, resp.StatusCode, string(respBody))
		err = json.Unmarshal(respBody, &rawItems)
		require.NoError(t, err)
		require.NotNil(t, rawItems)
		log.Printf("rawItems: %s", hack.UnsafeJSONPrettyFormat(rawItems))
		fmt.Printf("\n")

		resp, err = httpClient.Get("http://localhost:5050/not-null-fuzzes?some_bigint__in=1")
		require.NoError(t, err)
		respBody, _ = io.ReadAll(resp.Body)
		require.NoError(t, err, string(respBody))
		require.Equal(t, http.StatusOK, resp.StatusCode, string(respBody))
		err = json.Unmarshal(respBody, &rawItems)
		require.NoError(t, err)
		require.NotNil(t, rawItems)
		log.Printf("rawItems: %s", hack.UnsafeJSONPrettyFormat(rawItems))
		fmt.Printf("\n")
		items = rawItems.(map[string]interface{})
		objects, ok = items["objects"].([]interface{})
		require.True(t, ok)
		require.GreaterOrEqual(t, len(objects), 1)
		fmt.Printf("\n")

		resp, err = httpClient.Get("http://localhost:5050/not-null-fuzzes?some_bigint__in=69")
		require.NoError(t, err)
		respBody, _ = io.ReadAll(resp.Body)
		require.NoError(t, err, string(respBody))
		require.Equal(t, http.StatusOK, resp.StatusCode, string(respBody))
		err = json.Unmarshal(respBody, &rawItems)
		require.NoError(t, err)
		require.NotNil(t, rawItems)
		log.Printf("rawItems: %s", hack.UnsafeJSONPrettyFormat(rawItems))
		fmt.Printf("\n")
		items = rawItems.(map[string]interface{})
		require.Nil(t, items["objects"])
		fmt.Printf("\n")

		resp, err = httpClient.Get("http://localhost:5050/not-null-fuzzes?some_character_varying__in=A")
		require.NoError(t, err)
		respBody, _ = io.ReadAll(resp.Body)
		require.NoError(t, err, string(respBody))
		require.Equal(t, http.StatusOK, resp.StatusCode, string(respBody))
		err = json.Unmarshal(respBody, &rawItems)
		require.NoError(t, err)
		require.NotNil(t, rawItems)
		log.Printf("rawItems: %s", hack.UnsafeJSONPrettyFormat(rawItems))
		fmt.Printf("\n")
		items = rawItems.(map[string]interface{})
		objects = items["objects"].([]interface{})
		require.GreaterOrEqual(t, len(objects), 1)
		fmt.Printf("\n")

		resp, err = httpClient.Get("http://localhost:5050/not-null-fuzzes?some_character_varying__in=Z")
		require.NoError(t, err)
		respBody, _ = io.ReadAll(resp.Body)
		require.NoError(t, err, string(respBody))
		require.Equal(t, http.StatusOK, resp.StatusCode, string(respBody))
		err = json.Unmarshal(respBody, &rawItems)
		require.NoError(t, err)
		require.NotNil(t, rawItems)
		log.Printf("rawItems: %s", hack.UnsafeJSONPrettyFormat(rawItems))
		fmt.Printf("\n")
		items = rawItems.(map[string]interface{})
		require.Nil(t, items["objects"])
		fmt.Printf("\n")

		resp, err = httpClient.Get("http://localhost:5050/not-null-fuzzes?some_boolean__in=true")
		require.NoError(t, err)
		respBody, _ = io.ReadAll(resp.Body)
		require.NoError(t, err, string(respBody))
		require.Equal(t, http.StatusOK, resp.StatusCode, string(respBody))
		err = json.Unmarshal(respBody, &rawItems)
		require.NoError(t, err)
		require.NotNil(t, rawItems)
		log.Printf("rawItems: %s", hack.UnsafeJSONPrettyFormat(rawItems))
		fmt.Printf("\n")
		items = rawItems.(map[string]interface{})
		objects = items["objects"].([]interface{})
		require.GreaterOrEqual(t, len(objects), 1)
		fmt.Printf("\n")

		resp, err = httpClient.Get("http://localhost:5050/not-null-fuzzes?some_boolean__in=false")
		require.NoError(t, err)
		respBody, _ = io.ReadAll(resp.Body)
		require.NoError(t, err, string(respBody))
		require.Equal(t, http.StatusOK, resp.StatusCode, string(respBody))
		err = json.Unmarshal(respBody, &rawItems)
		require.NoError(t, err)
		require.NotNil(t, rawItems)
		log.Printf("rawItems: %s", hack.UnsafeJSONPrettyFormat(rawItems))
		fmt.Printf("\n")
		items = rawItems.(map[string]interface{})
		require.Nil(t, items["objects"])
		fmt.Printf("\n")

		resp, err = httpClient.Get("http://localhost:5050/not-null-fuzzes?some_timestamp__in=2020-03-27T08:30:00Z")
		require.NoError(t, err)
		respBody, _ = io.ReadAll(resp.Body)
		require.NoError(t, err, string(respBody))
		require.Equal(t, http.StatusOK, resp.StatusCode, string(respBody))
		err = json.Unmarshal(respBody, &rawItems)
		require.NoError(t, err)
		require.NotNil(t, rawItems)
		log.Printf("rawItems: %s", hack.UnsafeJSONPrettyFormat(rawItems))
		fmt.Printf("\n")
		items = rawItems.(map[string]interface{})
		objects = items["objects"].([]interface{})
		require.GreaterOrEqual(t, len(objects), 1)
		fmt.Printf("\n")

		resp, err = httpClient.Get("http://localhost:5050/not-null-fuzzes?some_timestamp__in=2020-04-27T08:30:00Z")
		require.NoError(t, err)
		respBody, _ = io.ReadAll(resp.Body)
		require.NoError(t, err, string(respBody))
		require.Equal(t, http.StatusOK, resp.StatusCode, string(respBody))
		err = json.Unmarshal(respBody, &rawItems)
		require.NoError(t, err)
		require.NotNil(t, rawItems)
		log.Printf("rawItems: %s", hack.UnsafeJSONPrettyFormat(rawItems))
		fmt.Printf("\n")
		items = rawItems.(map[string]interface{})
		require.Nil(t, items["objects"])
		fmt.Printf("\n")

		resp, err = httpClient.Get("http://localhost:5050/not-null-fuzzes?some_timestamptz__in=2024-07-19T11:45:00+08:00")
		require.NoError(t, err)
		respBody, _ = io.ReadAll(resp.Body)
		require.NoError(t, err, string(respBody))
		require.Equal(t, http.StatusOK, resp.StatusCode, string(respBody))
		err = json.Unmarshal(respBody, &rawItems)
		require.NoError(t, err)
		require.NotNil(t, rawItems)
		log.Printf("rawItems: %s", hack.UnsafeJSONPrettyFormat(rawItems))
		fmt.Printf("\n")
		items = rawItems.(map[string]interface{})
		objects = items["objects"].([]interface{})
		require.GreaterOrEqual(t, len(objects), 1)
		fmt.Printf("\n")

		resp, err = httpClient.Get("http://localhost:5050/not-null-fuzzes?some_timestamptz__in=2023-07-19T11:45:00+08:00")
		require.NoError(t, err)
		respBody, _ = io.ReadAll(resp.Body)
		require.NoError(t, err, string(respBody))
		require.Equal(t, http.StatusOK, resp.StatusCode, string(respBody))
		err = json.Unmarshal(respBody, &rawItems)
		require.NoError(t, err)
		require.NotNil(t, rawItems)
		log.Printf("rawItems: %s", hack.UnsafeJSONPrettyFormat(rawItems))
		fmt.Printf("\n")
		items = rawItems.(map[string]interface{})
		require.Nil(t, items["objects"])
		fmt.Printf("\n")
	})
}
