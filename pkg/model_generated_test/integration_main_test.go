package model_generated_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/initialed85/djangolang/pkg/helpers"
	"github.com/initialed85/djangolang/pkg/model_generated"
	"github.com/initialed85/djangolang/pkg/query"
	"github.com/initialed85/djangolang/pkg/server"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

func testIntegration(
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
			`DELETE FROM logical_things CASCADE
		WHERE
			name ILIKE $1;`,
			"%Integration%",
		)
		_, _ = db.Exec(
			ctx,
			`DELETE FROM physical_things CASCADE
		WHERE
			name ILIKE $1 OR name ILIKE $2;`,
			"%Integration%",
			"Performance-%",
		)
		if redisConn != nil {
			_, _ = redisConn.Do("FLUSHALL")
		}
	}
	defer cleanup()

	physicalThing1Name := "IntegrationPhysicalThing1Name"
	physicalThing1TypeA := "IntegrationPhysicalThing1Type"
	physicalThing1TypeB := "IntegrationPhysicalThing1Type"
	physicalThing1ItemA := []map[string]any{
		{
			"name": physicalThing1Name,
			"type": physicalThing1TypeA,
		},
	}

	logicalThing1Name := "IntegrationLogicalThing1Name"
	logicalThing1TypeA := "IntegrationLogicalThing1Type"
	logicalThing1ItemA := []map[string]any{
		{
			"name":                     logicalThing1Name,
			"type":                     logicalThing1TypeA,
			"count":                    0,
			"parent_physical_thing_id": nil,
		},
	}

	locationHistoryTimestamp, err := time.Parse(time.RFC3339, "2024-07-19T11:45:00+08:00")
	require.NoError(t, err)

	locationHistory1Point := map[string]any{
		"P": map[string]any{
			"X": 1.337,
			"Y": 69.420,
		},
	}
	locationHistory1Item := []map[string]any{
		{
			"timestamp":                locationHistoryTimestamp,
			"point":                    locationHistory1Point,
			"parent_physical_thing_id": nil,
		},
	}

	locationHistory2Polygon := []map[string]any{
		{
			"P": map[string]any{
				"X": 0.0,
				"Y": 0.0,
			},
		},
		{
			"P": map[string]any{
				"X": 1.0,
				"Y": 0.0,
			},
		},
		{
			"P": map[string]any{
				"X": 1.0,
				"Y": 1.0,
			},
		},
		{
			"P": map[string]any{
				"X": 0.0,
				"Y": 1.0,
			},
		},
		{
			"P": map[string]any{
				"X": 0.0,
				"Y": 0.0,
			},
		},
	}
	locationHistory2Item := []map[string]any{
		{
			"timestamp":                locationHistoryTimestamp,
			"polygon":                  locationHistory2Polygon,
			"parent_physical_thing_id": nil,
		},
	}

	setup := func() (*model_generated.PhysicalThing, *model_generated.LogicalThing, *model_generated.LocationHistory, *model_generated.LocationHistory) {
		//
		// PhysicalThing
		//

		physicalThingItemJSON, err := json.Marshal(physicalThing1ItemA)
		require.NoError(t, err)

		resp, err := httpClient.Post(
			"http://localhost:5050/physical-things",
			"application/json",
			bytes.NewReader(physicalThingItemJSON),
		)
		require.NoError(t, err)
		respBody, _ := io.ReadAll(resp.Body)
		require.NoError(t, err, string(respBody))
		require.Equal(t, http.StatusCreated, resp.StatusCode, string(respBody))

		require.Eventually(
			t,
			func() bool {
				change := getLastChangeForTableName(model_generated.PhysicalThingTable)
				if change == nil {
					return false
				}

				if change.Item["name"] != physicalThing1Name {
					return false
				}

				return true
			},
			time.Second*10,
			time.Millisecond*10,
			"failed to confirm PhysicalThing",
		)

		tx, err := db.Begin(ctx)
		require.NoError(t, err)
		defer func() {
			_ = tx.Rollback(ctx)
		}()

		physicalThing1, count, totalCount, page, totalPages, err := model_generated.SelectPhysicalThing(ctx, tx, "name = $$??", physicalThing1Name)
		require.NoError(t, err)
		require.Equal(t, physicalThing1Name, physicalThing1Name)
		require.Equal(t, int64(1), count)
		require.Equal(t, int64(1), totalCount)
		require.Equal(t, int64(1), page)
		require.Equal(t, int64(1), totalPages)

		err = tx.Commit(ctx)
		require.NoError(t, err)

		//
		// LogicalThing
		//

		logicalThing1ItemA[0]["parent_physical_thing_id"] = &physicalThing1.ID
		logicalThingItemJSON, err := json.Marshal(logicalThing1ItemA)
		require.NoError(t, err)

		resp, err = httpClient.Post(
			"http://localhost:5050/logical-things",
			"application/json",
			bytes.NewReader(logicalThingItemJSON),
		)
		respBody, _ = io.ReadAll(resp.Body)
		require.NoError(t, err, string(respBody))
		require.Equal(t, http.StatusCreated, resp.StatusCode, string(respBody))

		require.Eventually(
			t,
			func() bool {
				change := getLastChangeForTableName(model_generated.LogicalThingTable)
				if change == nil {
					return false
				}

				if change.Item["name"] != logicalThing1Name {
					return false
				}

				return true
			},
			time.Second*10,
			time.Millisecond*10,
			"failed to confirm LogicalThing",
		)

		tx, err = db.Begin(ctx)
		require.NoError(t, err)
		defer func() {
			_ = tx.Rollback(ctx)
		}()

		logicalThing1, count, totalCount, page, totalPages, err := model_generated.SelectLogicalThing(ctx, tx, "name = $$??", logicalThing1Name)
		require.NoError(t, err)
		require.Equal(t, logicalThing1Name, logicalThing1Name)
		require.Equal(t, int64(1), count)
		require.Equal(t, int64(1), totalCount)
		require.Equal(t, int64(1), page)
		require.Equal(t, int64(1), totalPages)

		err = tx.Commit(ctx)
		require.NoError(t, err)

		//
		// LocationHistory1
		//

		locationHistory1Item[0]["parent_physical_thing_id"] = &physicalThing1.ID
		locationHistory1ItemJSON, err := json.Marshal(locationHistory1Item)
		require.NoError(t, err)

		resp, err = httpClient.Post(
			"http://localhost:5050/location-histories",
			"application/json",
			bytes.NewReader(locationHistory1ItemJSON),
		)
		respBody, _ = io.ReadAll(resp.Body)
		var x any
		_ = json.Unmarshal(respBody, &x)
		respBody, _ = json.MarshalIndent(x, "", "  ")
		require.NoError(t, err, string(respBody))
		require.Equal(t, http.StatusCreated, resp.StatusCode, string(respBody))

		require.Eventually(
			t,
			func() bool {
				change := getLastChangeForTableName(model_generated.LocationHistoryTable)
				if change == nil {
					return false
				}

				object, ok := change.Object.(*model_generated.LocationHistory)
				if !ok {
					return false
				}

				if object.ParentPhysicalThingID == nil {
					return false
				}

				if *object.ParentPhysicalThingID != physicalThing1.ID {
					return false
				}

				return true
			},
			time.Second*10,
			time.Millisecond*10,
			"failed to confirm LocationHistory",
		)

		tx, err = db.Begin(ctx)
		require.NoError(t, err)
		defer func() {
			_ = tx.Rollback(ctx)
		}()

		ctx = query.WithMaxDepth(ctx, helpers.Ptr(0))

		locationHistory1, count, totalCount, page, totalPages, err := model_generated.SelectLocationHistory(ctx, tx, "parent_physical_thing_id = $$??", physicalThing1.ID)
		require.NoError(t, err)
		require.Equal(t, int64(1), count)
		require.Equal(t, int64(1), totalCount)
		require.Equal(t, int64(1), page)
		require.Equal(t, int64(1), totalPages)

		b, _ := json.MarshalIndent(locationHistory1, "", "  ")
		log.Printf("b: %v", string(b))

		require.NotNil(t, locationHistory1.ParentPhysicalThingIDObject)
		require.Equal(t, physicalThing1.ID, locationHistory1.ParentPhysicalThingIDObject.ID)

		err = tx.Commit(ctx)
		require.NoError(t, err)

		//
		// LocationHistory2
		//

		locationHistory2Item[0]["parent_physical_thing_id"] = &physicalThing1.ID
		locationHistory2ItemJSON, err := json.Marshal(locationHistory2Item)
		require.NoError(t, err)

		resp, err = httpClient.Post(
			"http://localhost:5050/location-histories",
			"application/json",
			bytes.NewReader(locationHistory2ItemJSON),
		)
		respBody, _ = io.ReadAll(resp.Body)
		_ = json.Unmarshal(respBody, &x)
		respBody, _ = json.MarshalIndent(x, "", "  ")
		require.NoError(t, err, string(respBody))
		require.Equal(t, http.StatusCreated, resp.StatusCode, string(respBody))

		require.Eventually(
			t,
			func() bool {
				change := getLastChangeForTableName(model_generated.LocationHistoryTable)
				if change == nil {
					return false
				}

				object, ok := change.Object.(*model_generated.LocationHistory)
				if !ok {
					return false
				}

				if object.ParentPhysicalThingID == nil {
					return false
				}

				if *object.ParentPhysicalThingID != physicalThing1.ID {
					return false
				}

				return true
			},
			time.Second*10,
			time.Millisecond*10,
			"failed to confirm LocationHistory",
		)

		tx, err = db.Begin(ctx)
		require.NoError(t, err)
		defer func() {
			_ = tx.Rollback(ctx)
		}()

		locationHistory2, count, totalCount, page, totalPages, err := model_generated.SelectLocationHistory(ctx, tx, "parent_physical_thing_id = $$?? AND id != $$??", physicalThing1.ID, locationHistory1.ID)
		require.NoError(t, err)
		require.NotNil(t, locationHistory2.ParentPhysicalThingIDObject)
		require.Equal(t, physicalThing1.ID, locationHistory1.ParentPhysicalThingIDObject.ID)
		require.Equal(t, int64(1), count)
		require.Equal(t, int64(1), totalCount)
		require.Equal(t, int64(1), page)
		require.Equal(t, int64(1), totalPages)

		err = tx.Commit(ctx)
		require.NoError(t, err)

		//
		// reloads
		//

		tx, err = db.Begin(ctx)
		require.NoError(t, err)
		defer func() {
			_ = tx.Rollback(ctx)
		}()

		_ = physicalThing1.Reload(ctx, tx)
		_ = logicalThing1.Reload(ctx, tx)
		_ = locationHistory1.Reload(ctx, tx)

		_ = tx.Commit(ctx)

		return physicalThing1, logicalThing1, locationHistory1, locationHistory2
	}
	physicalThing1, logicalThing1, locationHistory1, locationHistory2 := setup()

	_ = physicalThing1
	_ = logicalThing1
	_ = locationHistory1
	_ = locationHistory2

	t.Run("PhysicalThing", func(t *testing.T) {
		resp, err := httpClient.Get(
			fmt.Sprintf("http://localhost:5050/physical-things/%s", physicalThing1.ID.String()),
		)
		respBody, _ := io.ReadAll(resp.Body)
		require.NoError(t, err, string(respBody))
		require.Equal(t, http.StatusOK, resp.StatusCode, string(respBody))

		physicalThing1ItemB := map[string]any{
			"type": physicalThing1TypeB,
		}

		physicalThingItemBJSON, err := json.Marshal(physicalThing1ItemB)
		require.NoError(t, err)

		_ = physicalThingItemBJSON

		resp, err = httpClient.Patch(
			fmt.Sprintf("http://localhost:5050/physical-things/%s", physicalThing1.ID.String()),
			"application/json",
			bytes.NewReader(physicalThingItemBJSON),
		)
		respBody, _ = io.ReadAll(resp.Body)
		require.NoError(t, err, string(respBody))
		require.Equal(t, http.StatusOK, resp.StatusCode, string(respBody))

		tx, err := db.Begin(ctx)
		require.NoError(t, err)
		defer func() {
			_ = tx.Rollback(ctx)
		}()
		err = physicalThing1.Reload(ctx, tx)
		require.NoError(t, err)
		err = tx.Commit(ctx)
		require.NoError(t, err)

		require.Equal(t, physicalThing1TypeB, physicalThing1.Type)
	})

	t.Run("Performance", func(t *testing.T) {
		t.Skipf("TOOD")

		physicalThings := make([]model_generated.PhysicalThing, 0)
		for i := 0; i < 1024; i++ {
			physicalThings = append(physicalThings, model_generated.PhysicalThing{
				Name:     fmt.Sprintf("Performance-%d", i),
				Type:     "Performance",
				Tags:     []string{},
				Metadata: map[string]*string{},
			})
		}

		physicalThingItemsJSON, err := json.Marshal(physicalThings)
		require.NoError(t, err)

		resp, err := httpClient.Post(
			"http://localhost:5050/physical-things",
			"application/json",
			bytes.NewReader(physicalThingItemsJSON),
		)
		require.NoError(t, err)
		respBody, _ := io.ReadAll(resp.Body)
		require.NoError(t, err, string(respBody))
		require.Equal(t, http.StatusCreated, resp.StatusCode, string(respBody))
	})
}
