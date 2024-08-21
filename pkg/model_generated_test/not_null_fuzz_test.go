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
	"os"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/cridenour/go-postgis"
	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
	"github.com/initialed85/djangolang/pkg/helpers"
	"github.com/initialed85/djangolang/pkg/model_generated"
	"github.com/initialed85/djangolang/pkg/query"
	"github.com/initialed85/djangolang/pkg/server"
	"github.com/initialed85/djangolang/pkg/stream"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func TestNotNullFuzz(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctx = query.WithMaxDepth(ctx, helpers.Ptr(0))

	db, err := helpers.GetDBFromEnvironment(ctx)
	if err != nil {
		require.NoError(t, err)
	}
	defer func() {
		_ = db.Close()
	}()

	redisURL := helpers.GetRedisURL()

	var redisPool *redis.Pool
	var redisConn redis.Conn
	if redisURL != "" {
		redisPool = &redis.Pool{
			DialContext: func(ctx context.Context) (redis.Conn, error) {
				return redis.DialURLContext(ctx, redisURL)
			},
			MaxIdle:         2,
			MaxActive:       100,
			IdleTimeout:     time.Second * 300,
			Wait:            false,
			MaxConnLifetime: time.Hour * 24,
		}

		defer func() {
			_ = redisPool.Close()
		}()

		redisConn = redisPool.Get()
		defer func() {
			_ = redisConn.Close()
		}()
	}

	httpClient := &HTTPClient{
		httpClient: &http.Client{
			Timeout: time.Second * 10,
		},
	}

	changes := make(chan server.Change, 1024)
	mu := new(sync.Mutex)
	lastChangeByTableName := make(map[string]server.Change)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case change := <-changes:
				mu.Lock()
				lastChangeByTableName[change.TableName] = change
				mu.Unlock()
			}
		}
	}()
	runtime.Gosched()

	getLastChangeForTableName := func(tableName string) *server.Change {
		mu.Lock()
		defer mu.Unlock()

		change, ok := lastChangeByTableName[tableName]
		if !ok {
			return nil
		}

		return &change
	}

	_ = getLastChangeForTableName

	go func() {
		os.Setenv("DJANGOLANG_NODE_NAME", "model_generated_not_null_fuzz_test")
		err = model_generated.RunServer(ctx, changes, "127.0.0.1:3030", db, redisPool, nil, nil)
		if err != nil {
			log.Printf("model_generated.RunServer failed: %v", err)
		}
	}()
	runtime.Gosched()

	require.Eventually(
		t,
		func() bool {
			resp, err := httpClient.Get("http://localhost:3030/logical-things")
			if err != nil {
				return false
			}

			if resp.StatusCode != http.StatusOK {
				return false
			}

			return true
		},
		time.Second*1,
		time.Millisecond*100,
	)

	cleanup := func() {
		_, _ = db.ExecContext(
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

	_, err = db.ExecContext(
		ctx,
		`INSERT INTO not_null_fuzz DEFAULT VALUES;`,
	)
	require.NoError(t, err)

	t.Run("ObjectInteractions", func(t *testing.T) {
		tx, err := db.BeginTxx(ctx, nil)
		require.NoError(t, err)
		defer func() {
			_ = tx.Rollback()
		}()

		notNullFuzzes, err := model_generated.SelectNotNullFuzzes(ctx, tx, "", nil, nil, nil)
		require.NoError(t, err)
		require.Len(t, notNullFuzzes, 1)

		err = tx.Commit()
		require.NoError(t, err)

		notNullFuzz := notNullFuzzes[0]

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
		require.Equal(t, "2024-07-19 03:45:00 +0000 UTC", notNullFuzz.SomeTimestamptz.String(), "SomeTimestamptz")
		require.Equal(t, "2020-03-27 08:30:00 +0000 +0000", notNullFuzz.SomeTimestamp.String(), "SomeTimestamp")
		require.Equal(t, map[string][]int(map[string][]int{"a": []int(nil)}), notNullFuzz.SomeTsvector, "SomeTsvector")
		require.Equal(t, uuid.UUID{0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11}, notNullFuzz.SomeUUID, "SomeUUID")

		require.Eventually(
			t,
			func() bool {
				return getLastChangeForTableName(model_generated.NotNullFuzzTable) != nil
			},
			time.Second*10,
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
		require.Equal(t, "2024-07-19 03:45:00 +0000 UTC", object.SomeTimestamptz.String(), "SomeTimestamptz")
		require.Equal(t, "2020-03-27 08:30:00 +0000 +0000", object.SomeTimestamp.String(), "SomeTimestamp")
		require.Equal(t, map[string][]int(map[string][]int{"a": []int(nil)}), object.SomeTsvector, "SomeTsvector")
		require.Equal(t, uuid.UUID{0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11}, object.SomeUUID, "SomeUUID")

		require.Equal(t, int64(1), lastChange.Item["some_bigint"], "SomeBigint")
		require.Equal(t, []any{int64(1)}, lastChange.Item["some_bigint_array"], "SomeBigintArray")
		require.Equal(t, true, lastChange.Item["some_boolean"], "SomeBoolean")
		require.Equal(t, []any{true}, lastChange.Item["some_boolean_array"], "SomeBooleanArray")
		require.Equal(t, []byte{0x65}, lastChange.Item["some_bytea"], "SomeBytea")
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
		require.Equal(t, pgtype.Numeric{Int: big.NewInt(10), Exp: -1, NaN: false, InfinityModifier: 0, Valid: true}, lastChange.Item["some_numeric"], "SomeNumeric")
		require.Equal(t, []any{pgtype.Numeric{Int: big.NewInt(10), Exp: -1, NaN: false, InfinityModifier: 0, Valid: true}}, lastChange.Item["some_numeric_array"], "SomeNumericArray")
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

		tx, err = db.BeginTxx(ctx, nil)
		require.NoError(t, err)
		defer func() {
			_ = tx.Rollback()
		}()

		notNullFuzzes, err = model_generated.SelectNotNullFuzzes(ctx, tx, "", nil, nil, nil)
		require.NoError(t, err)
		require.Len(t, notNullFuzzes, 1)

		err = tx.Commit()
		require.NoError(t, err)

		notNullFuzz = notNullFuzzes[0]

		tx, err = db.BeginTxx(ctx, nil)
		require.NoError(t, err)
		defer func() {
			_ = tx.Rollback()
		}()

		notNullFuzz.SomeBigint = 2
		notNullFuzz.SomeBigintArray = append(notNullFuzz.SomeBigintArray, 2)

		err = notNullFuzz.Insert(ctx, tx, false, false)
		require.NoError(t, err)

		err = tx.Commit()
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
		require.Equal(t, "2024-07-19 03:45:00 +0000 UTC", notNullFuzz.SomeTimestamptz.String(), "SomeTimestamptz")
		require.Equal(t, "2020-03-27 08:30:00 +0000 +0000", notNullFuzz.SomeTimestamp.String(), "SomeTimestamp")
		require.Equal(t, map[string][]int(map[string][]int{"a": []int(nil)}), notNullFuzz.SomeTsvector, "SomeTsvector")
		require.Equal(t, uuid.UUID{0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11}, notNullFuzz.SomeUUID, "SomeUUID")

		tx, err = db.BeginTxx(ctx, nil)
		require.NoError(t, err)
		defer func() {
			_ = tx.Rollback()
		}()

		notNullFuzzes, err = model_generated.SelectNotNullFuzzes(ctx, tx, fmt.Sprintf("%s = $$??", model_generated.NotNullFuzzTableSomeBigintColumn), nil, nil, nil, 2)
		require.NoError(t, err)
		require.Len(t, notNullFuzzes, 1)

		err = tx.Commit()
		require.NoError(t, err)

		notNullFuzz = notNullFuzzes[0]

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
		require.Equal(t, "2024-07-19 03:45:00 +0000 UTC", notNullFuzz.SomeTimestamptz.String(), "SomeTimestamptz")
		require.Equal(t, "2020-03-27 08:30:00 +0000 +0000", notNullFuzz.SomeTimestamp.String(), "SomeTimestamp")
		require.Equal(t, map[string][]int(map[string][]int{"a": []int(nil)}), notNullFuzz.SomeTsvector, "SomeTsvector")
		require.Equal(t, uuid.UUID{0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11}, notNullFuzz.SomeUUID, "SomeUUID")

		tx, err = db.BeginTxx(ctx, nil)
		require.NoError(t, err)
		defer func() {
			_ = tx.Rollback()
		}()

		notNullFuzz.SomeCharacterVarying = "B"
		err = notNullFuzz.Update(ctx, tx, false)
		require.NoError(t, err)

		err = tx.Commit()
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
		require.Equal(t, "2024-07-19 03:45:00 +0000 UTC", notNullFuzz.SomeTimestamptz.String(), "SomeTimestamptz")
		require.Equal(t, "2020-03-27 08:30:00 +0000 +0000", notNullFuzz.SomeTimestamp.String(), "SomeTimestamp")
		require.Equal(t, map[string][]int(map[string][]int{"a": []int(nil)}), notNullFuzz.SomeTsvector, "SomeTsvector")
		require.Equal(t, uuid.UUID{0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11}, notNullFuzz.SomeUUID, "SomeUUID")

		tx, err = db.BeginTxx(ctx, nil)
		require.NoError(t, err)
		defer func() {
			_ = tx.Rollback()
		}()

		err = notNullFuzz.Delete(ctx, tx)
		require.NoError(t, err)

		err = tx.Commit()
		require.NoError(t, err)

		require.Eventually(
			t,
			func() bool {
				return getLastChangeForTableName(model_generated.NotNullFuzzTable) != nil
			},
			time.Second*10,
			time.Millisecond*10,
		)

		lastChange = getLastChangeForTableName(model_generated.NotNullFuzzTable)
		require.NotNil(t, lastChange)
		require.Equal(t, stream.DELETE, lastChange.Action)
	})

	t.Run("EndpointInteractions", func(t *testing.T) {
		resp, err := httpClient.Get("http://localhost:3030/not-null-fuzzes?some_bigint__eq=1")
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
			"some_timestamptz":    "2024-07-19T03:45:00Z",
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

		resp, err = httpClient.Post("http://localhost:3030/not-null-fuzzes", "application/json", bytes.NewReader(b))
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

		resp, err = httpClient.Put(fmt.Sprintf("http://localhost:3030/not-null-fuzzes/%v", object["id"]), "application/json", bytes.NewReader(b))
		require.NoError(t, err)
		respBody, _ = io.ReadAll(resp.Body)
		require.NoError(t, err, string(respBody))
		require.Equal(t, http.StatusOK, resp.StatusCode, string(respBody))

		err = json.Unmarshal(respBody, &rawItems)
		require.NoError(t, err)

		items = rawItems.(map[string]interface{})
		objects = items["objects"].([]interface{})
		object = objects[0].(map[string]interface{})

		for k, v := range expectedObjectIsh {
			require.Equal(t, v, object[k], k)
		}
	})
}
