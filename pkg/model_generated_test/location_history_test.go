package model_generated_test

import (
	"context"
	"fmt"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
	"github.com/initialed85/djangolang/internal/hack"
	"github.com/initialed85/djangolang/pkg/helpers"
	"github.com/initialed85/djangolang/pkg/model_generated"
	"github.com/initialed85/djangolang/pkg/query"
	"github.com/initialed85/djangolang/pkg/server"
	"github.com/initialed85/djangolang/pkg/stream"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

func testLocationHistory(
	t *testing.T,
	ctx context.Context,
	db *pgxpool.Pool,
	redisConn redis.Conn,
	mu *sync.Mutex,
	lastChangeByTableName map[string]server.Change,
	httpClient *HTTPClient,
	getLastChangeForTableName func(tableName string) *server.Change,
) {
	count := int64(-1)
	totalCount := int64(-1)
	page := int64(-1)
	totalPages := int64(-1)

	_ = count
	_ = totalCount
	_ = page
	_ = totalPages

	t.Run("SelectWithPoint", func(t *testing.T) {
		physicalExternalID := "SelectWithPointSomePhysicalThingExternalID"
		physicalThingName := "SelectWithPointSomePhysicalThingName"
		physicalThingType := "SelectWithPointSomePhysicalThingType"
		physicalThingTags := `'{tag1,tag2,tag3,"isn''t this, \"complicated\""}'`
		physicalThingMetadata := `'key1=>1, key2=>"a", key3=>true, key4=>NULL, key5=>"isn''t this, \"complicated\""'`
		physicalThingRawData := `'{"key1": 1, "key2": "a", "key3": true, "key4": null, "key5": "isn''t this, \"complicated\""}'`

		cleanup := func() {
			_, err := db.Exec(
				ctx,
				`DELETE FROM location_history;`,
			)
			require.NoError(t, err)
			_, err = db.Exec(
				ctx,
				`DELETE FROM physical_things CASCADE
			WHERE
				name = $1;`,
				physicalThingName,
			)
			require.NoError(t, err)
		}
		defer cleanup()

		_, err := db.Exec(
			ctx,
			fmt.Sprintf(`INSERT INTO physical_things (
				external_id,
				name,
				type,
				tags,
				metadata,
				raw_data
			)
			VALUES (
				'%v',
				'%v',
				'%v',
				%v,
				%v,
				%v
			);`,
				physicalExternalID,
				physicalThingName,
				physicalThingType,
				physicalThingTags,
				physicalThingMetadata,
				physicalThingRawData,
			),
		)
		require.NoError(t, err)

		var locationHistory *model_generated.LocationHistory
		var physicalThing *model_generated.PhysicalThing

		func() {
			tx, _ := db.Begin(ctx)
			defer tx.Rollback(ctx)

			ctx = query.WithMaxDepth(ctx, helpers.Ptr(0))

			physicalThing, count, totalCount, page, totalPages, err = model_generated.SelectPhysicalThing(
				ctx,
				tx,
				fmt.Sprintf(
					"%v = $1 AND %v = $2 AND %v = $3",
					model_generated.PhysicalThingTableExternalIDColumn,
					model_generated.PhysicalThingTableNameColumn,
					model_generated.PhysicalThingTableTypeColumn,
				),
				physicalExternalID,
				physicalThingName,
				physicalThingType,
			)
			require.NoError(t, err)
			require.NotNil(t, physicalThing)

			_, err = db.Exec(
				ctx,
				fmt.Sprintf(`INSERT INTO location_history (
				timestamp,
				point,
				parent_physical_thing_id
			)
			VALUES (
				now(),
				ST_MakePoint(1.01, 2.02, 3.03)::point,
				'%v'
			);`,
					physicalThing.ID.String(),
				),
			)
			require.NoError(t, err)

			locationHistory, count, totalCount, page, totalPages, err = model_generated.SelectLocationHistory(
				ctx,
				tx,
				"parent_physical_thing_id = $1",
				physicalThing.ID.String(),
			)
			require.NoError(t, err)
			require.NotNil(t, locationHistory)
			_ = tx.Commit(ctx)
			require.NotNil(t, physicalThing)
		}()

		log.Printf("locationHistory: %v", hack.UnsafeJSONPrettyFormat(locationHistory))

		require.IsType(t, uuid.UUID{}, locationHistory.ID, "ID")
		require.IsType(t, time.Time{}, locationHistory.CreatedAt, "CreatedAt")
		require.IsType(t, time.Time{}, locationHistory.UpdatedAt, "UpdatedAt")
		require.IsType(t, helpers.Nil(time.Time{}), locationHistory.DeletedAt, "DeletedAt")
		require.IsType(t, time.Time{}, locationHistory.Timestamp, "Timestamp")
		require.IsType(t, &pgtype.Vec2{}, locationHistory.Point, "Point")
		require.IsType(t, &[]pgtype.Vec2{}, locationHistory.Polygon, "Polygon")
		require.IsType(t, helpers.Ptr(uuid.UUID{}), locationHistory.ParentPhysicalThingID, "ID")

		require.IsType(t, uuid.UUID{}, locationHistory.ParentPhysicalThingIDObject.ID, "ID")
		require.IsType(t, time.Time{}, locationHistory.ParentPhysicalThingIDObject.CreatedAt, "CreatedAt")
		require.IsType(t, time.Time{}, locationHistory.ParentPhysicalThingIDObject.UpdatedAt, "UpdatedAt")
		require.IsType(t, helpers.Nil(time.Time{}), locationHistory.ParentPhysicalThingIDObject.DeletedAt, "DeletedAt")
		require.IsType(t, helpers.Ptr(""), locationHistory.ParentPhysicalThingIDObject.ExternalID, "ExternalID")
		require.IsType(t, "", locationHistory.ParentPhysicalThingIDObject.Name, "Name")
		require.IsType(t, "", locationHistory.ParentPhysicalThingIDObject.Type, "Type")
		require.IsType(t, []string{}, locationHistory.ParentPhysicalThingIDObject.Tags, "Tags")
		require.IsType(t, map[string]*string{}, locationHistory.ParentPhysicalThingIDObject.Metadata, "Metadata")
		require.IsType(t, new(any), locationHistory.ParentPhysicalThingIDObject.RawData, "ParentPhysicalThingIDObject")

		var lastChange server.Change
		require.Eventually(t, func() bool {
			mu.Lock()
			defer mu.Unlock()

			var ok bool
			lastChange, ok = lastChangeByTableName[model_generated.LocationHistoryTable]
			if !ok {
				return false
			}

			if lastChange.Action != stream.INSERT {
				return false
			}

			return true
		}, time.Second*10, time.Millisecond*10)

		locationHistoryFromLastChange := &model_generated.LocationHistory{}
		err = locationHistoryFromLastChange.FromItem(lastChange.Item)
		require.NoError(t, err)

		log.Printf("locationHistoryFromLastChange: %v", hack.UnsafeJSONPrettyFormat(locationHistoryFromLastChange))

		require.Equal(t, locationHistory.ID, locationHistoryFromLastChange.ID)
		require.Equal(t, locationHistory.CreatedAt.Unix(), locationHistoryFromLastChange.CreatedAt.Unix())
		require.Equal(t, locationHistory.UpdatedAt.Unix(), locationHistoryFromLastChange.UpdatedAt.Unix())
		require.Equal(t, locationHistory.DeletedAt, locationHistoryFromLastChange.DeletedAt)
		require.Equal(t, locationHistory.Timestamp.Unix(), locationHistoryFromLastChange.Timestamp.Unix())
		require.Equal(t, locationHistory.Point, locationHistoryFromLastChange.Point)
		require.Equal(t, locationHistory.Polygon, locationHistoryFromLastChange.Polygon)
		require.Equal(t, locationHistory.ParentPhysicalThingID, locationHistoryFromLastChange.ParentPhysicalThingID)
		require.NotEqual(t, locationHistory.ParentPhysicalThingIDObject, locationHistoryFromLastChange.ParentPhysicalThingIDObject)

		func() {
			tx, _ := db.Begin(ctx)
			defer tx.Rollback(ctx)
			err = locationHistoryFromLastChange.Reload(ctx, tx)
			require.NoError(t, err)
			require.Equal(t, locationHistory.ParentPhysicalThingIDObject, locationHistoryFromLastChange.ParentPhysicalThingIDObject)
			_ = tx.Commit(ctx)
		}()

		log.Printf("locationHistoryFromLastChangeAfterReloading: %v", hack.UnsafeJSONPrettyFormat(locationHistoryFromLastChange))
		time.Sleep(time.Millisecond * 100)

		func() {
			tx, _ := db.Begin(ctx)
			defer tx.Rollback(ctx)

			parentPhysicalThing := locationHistory.ParentPhysicalThingIDObject
			err = parentPhysicalThing.Reload(ctx, tx)
			require.NoError(t, err)
			require.Len(t, parentPhysicalThing.ReferencedByLocationHistoryParentPhysicalThingIDObjects, 1)
			require.Equal(t, locationHistory.ID, parentPhysicalThing.ReferencedByLocationHistoryParentPhysicalThingIDObjects[0].ID)

			_ = tx.Commit(ctx)
		}()
	})

	t.Run("SelectWithPolygon", func(t *testing.T) {
		physicalExternalID := "SelectWithPolygonSomePhysicalThingExternalID"
		physicalThingName := "SelectWithPolygonSomePhysicalThingName"
		physicalThingType := "SelectWithPolygonSomePhysicalThingType"
		physicalThingTags := `'{tag1,tag2,tag3,"isn''t this, \"complicated\""}'`
		physicalThingMetadata := `'key1=>1, key2=>"a", key3=>true, key4=>NULL, key5=>"isn''t this, \"complicated\""'`
		physicalThingRawData := `'{"key1": 1, "key2": "a", "key3": true, "key4": null, "key5": "isn''t this, \"complicated\""}'`

		cleanup := func() {
			_, err := db.Exec(
				ctx,
				`DELETE FROM location_history;`,
			)
			require.NoError(t, err)
			_, err = db.Exec(
				ctx,
				`DELETE FROM physical_things CASCADE
			WHERE
				name = $1;`,
				physicalThingName,
			)
			require.NoError(t, err)
		}
		defer cleanup()

		_, err := db.Exec(
			ctx,
			fmt.Sprintf(`INSERT INTO physical_things (
				external_id,
				name,
				type,
				tags,
				metadata,
				raw_data
			)
			VALUES (
				'%v',
				'%v',
				'%v',
				%v,
				%v,
				%v
			);`,
				physicalExternalID,
				physicalThingName,
				physicalThingType,
				physicalThingTags,
				physicalThingMetadata,
				physicalThingRawData,
			),
		)
		require.NoError(t, err)

		var locationHistory *model_generated.LocationHistory
		var physicalThing *model_generated.PhysicalThing

		func() {
			tx, _ := db.Begin(ctx)
			defer tx.Rollback(ctx)
			physicalThing, count, totalCount, page, totalPages, err = model_generated.SelectPhysicalThing(
				ctx,
				tx,
				fmt.Sprintf(
					"%v = $1 AND %v = $2 AND %v = $3",
					model_generated.PhysicalThingTableExternalIDColumn,
					model_generated.PhysicalThingTableNameColumn,
					model_generated.PhysicalThingTableTypeColumn,
				),
				physicalExternalID,
				physicalThingName,
				physicalThingType,
			)
			require.NoError(t, err)
			require.NotNil(t, physicalThing)

			_, err = db.Exec(
				ctx,
				fmt.Sprintf(`INSERT INTO location_history (
				timestamp,
				polygon,
				parent_physical_thing_id
			)
			VALUES (
				now(),
				ST_MakePolygon( ST_GeomFromText('LINESTRING(75 29,77 29,77 29, 75 29)'))::polygon,
				'%v'
			);`,
					physicalThing.ID.String(),
				),
			)
			require.NoError(t, err)

			locationHistory, count, totalCount, page, totalPages, err = model_generated.SelectLocationHistory(
				ctx,
				tx,
				"parent_physical_thing_id = $1",
				physicalThing.ID.String(),
			)
			require.NoError(t, err)
			require.NotNil(t, locationHistory)
			_ = tx.Commit(ctx)
			require.NotNil(t, physicalThing)
		}()

		log.Printf("locationHistory: %v", hack.UnsafeJSONPrettyFormat(locationHistory))

		require.IsType(t, uuid.UUID{}, locationHistory.ID, "ID")
		require.IsType(t, time.Time{}, locationHistory.CreatedAt, "CreatedAt")
		require.IsType(t, time.Time{}, locationHistory.UpdatedAt, "UpdatedAt")
		require.IsType(t, helpers.Nil(time.Time{}), locationHistory.DeletedAt, "DeletedAt")
		require.IsType(t, time.Time{}, locationHistory.Timestamp, "Timestamp")
		require.IsType(t, &pgtype.Vec2{}, locationHistory.Point, "Point")
		require.IsType(t, &[]pgtype.Vec2{}, locationHistory.Polygon, "Polygon")
		require.IsType(t, helpers.Ptr(uuid.UUID{}), locationHistory.ParentPhysicalThingID, "ID")

		require.IsType(t, uuid.UUID{}, locationHistory.ParentPhysicalThingIDObject.ID, "ID")
		require.IsType(t, time.Time{}, locationHistory.ParentPhysicalThingIDObject.CreatedAt, "CreatedAt")
		require.IsType(t, time.Time{}, locationHistory.ParentPhysicalThingIDObject.UpdatedAt, "UpdatedAt")
		require.IsType(t, helpers.Nil(time.Time{}), locationHistory.ParentPhysicalThingIDObject.DeletedAt, "DeletedAt")
		require.IsType(t, helpers.Ptr(""), locationHistory.ParentPhysicalThingIDObject.ExternalID, "ExternalID")
		require.IsType(t, "", locationHistory.ParentPhysicalThingIDObject.Name, "Name")
		require.IsType(t, "", locationHistory.ParentPhysicalThingIDObject.Type, "Type")
		require.IsType(t, []string{}, locationHistory.ParentPhysicalThingIDObject.Tags, "Tags")
		require.IsType(t, map[string]*string{}, locationHistory.ParentPhysicalThingIDObject.Metadata, "Metadata")
		require.IsType(t, new(any), locationHistory.ParentPhysicalThingIDObject.RawData, "ParentPhysicalThingIDObject")

		var lastChange server.Change
		require.Eventually(t, func() bool {
			mu.Lock()
			defer mu.Unlock()

			var ok bool
			lastChange, ok = lastChangeByTableName[model_generated.LocationHistoryTable]
			if !ok {
				return false
			}

			if lastChange.Action != stream.INSERT {
				return false
			}

			return true
		}, time.Second*10, time.Millisecond*10)

		locationHistoryFromLastChange := &model_generated.LocationHistory{}
		err = locationHistoryFromLastChange.FromItem(lastChange.Item)
		require.NoError(t, err)

		log.Printf("locationHistoryFromLastChange: %v", hack.UnsafeJSONPrettyFormat(locationHistoryFromLastChange))

		require.Equal(t, locationHistory.ID, locationHistoryFromLastChange.ID)
		require.Equal(t, locationHistory.CreatedAt.Unix(), locationHistoryFromLastChange.CreatedAt.Unix())
		require.Equal(t, locationHistory.UpdatedAt.Unix(), locationHistoryFromLastChange.UpdatedAt.Unix())
		require.Equal(t, locationHistory.DeletedAt, locationHistoryFromLastChange.DeletedAt)
		require.Equal(t, locationHistory.Timestamp.Unix(), locationHistoryFromLastChange.Timestamp.Unix())
		require.Equal(t, locationHistory.Point, locationHistoryFromLastChange.Point)
		require.Equal(t, locationHistory.Polygon, locationHistoryFromLastChange.Polygon)
		require.Equal(t, locationHistory.ParentPhysicalThingID, locationHistoryFromLastChange.ParentPhysicalThingID)
		require.NotEqual(t, locationHistory.ParentPhysicalThingIDObject, locationHistoryFromLastChange.ParentPhysicalThingIDObject)

		func() {
			tx, _ := db.Begin(ctx)
			defer tx.Rollback(ctx)
			err = locationHistoryFromLastChange.Reload(ctx, tx)
			require.NoError(t, err)
			require.Equal(t, locationHistory.ParentPhysicalThingIDObject, locationHistoryFromLastChange.ParentPhysicalThingIDObject)
			_ = tx.Commit(ctx)
		}()

		log.Printf("locationHistoryFromLastChangeAfterReloading: %v", hack.UnsafeJSONPrettyFormat(locationHistoryFromLastChange))
	})
}
