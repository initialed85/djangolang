package model_generated_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	_helpers "github.com/initialed85/djangolang/internal/helpers"
	"github.com/initialed85/djangolang/pkg/helpers"
	"github.com/initialed85/djangolang/pkg/introspect"
	"github.com/initialed85/djangolang/pkg/model_generated"
	"github.com/initialed85/djangolang/pkg/stream"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func TestLocationHistory(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := helpers.GetDBFromEnvironment(ctx)
	if err != nil {
		require.NoError(t, err)
	}
	defer func() {
		_ = db.Close()
	}()

	schema := helpers.GetSchema()

	tableByName, err := introspect.Introspect(ctx, db, schema)
	require.NoError(t, err)

	changes := make(chan stream.Change, 1024)
	mu := new(sync.Mutex)
	lastChangeByTableName := make(map[string]stream.Change)

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

	go func() {
		os.Setenv("DJANGOLANG_NODE_NAME", "model_generated_location_history")
		err = stream.Run(ctx, changes, tableByName)
		require.NoError(t, err)
	}()
	runtime.Gosched()

	time.Sleep(time.Second * 1)

	t.Run("SelectWithPoint", func(t *testing.T) {
		physicalExternalID := "SomePhysicalThingExternalID4"
		physicalThingName := "SomePhysicalThingName4"
		physicalThingType := "SomePhysicalThingType4"
		physicalThingTags := `'{tag1,tag2,tag3,"isn''t this, \"complicated\""}'`
		physicalThingMetadata := `'key1=>1, key2=>"a", key3=>true, key4=>NULL, key5=>"isn''t this, \"complicated\""'`
		physicalThingRawData := `'{"key1": 1, "key2": "a", "key3": true, "key4": null, "key5": "isn''t this, \"complicated\""}'`

		cleanup := func() {
			_, err = db.ExecContext(
				ctx,
				`DELETE FROM location_history;`,
			)
			_, err = db.ExecContext(
				ctx,
				`DELETE FROM physical_things CASCADE
			WHERE
				name = $1;`,
				physicalThingName,
			)
			require.NoError(t, err)
		}
		cleanup()

		_, err = db.ExecContext(
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
		var err error

		func() {
			tx, _ := db.BeginTxx(ctx, nil)
			defer tx.Rollback()
			physicalThing, err = model_generated.SelectPhysicalThing(
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

			_, err = db.ExecContext(
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

			locationHistory, err = model_generated.SelectLocationHistory(
				ctx,
				tx,
				"parent_physical_thing_id = $1",
				physicalThing.ID.String(),
			)
			require.NoError(t, err)
			require.NotNil(t, locationHistory)
			_ = tx.Commit()
			require.NotNil(t, physicalThing)
		}()

		log.Printf("locationHistory: %v", _helpers.UnsafeJSONPrettyFormat(locationHistory))

		require.IsType(t, uuid.UUID{}, locationHistory.ID, "ID")
		require.IsType(t, time.Time{}, locationHistory.CreatedAt, "CreatedAt")
		require.IsType(t, time.Time{}, locationHistory.UpdatedAt, "UpdatedAt")
		require.IsType(t, helpers.Nil(time.Time{}), locationHistory.DeletedAt, "DeletedAt")
		require.IsType(t, time.Time{}, locationHistory.Timestamp, "Timestamp")
		require.IsType(t, &pgtype.Point{}, locationHistory.Point, "Point")
		require.IsType(t, &pgtype.Polygon{}, locationHistory.Polygon, "Polygon")
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

		var lastChange stream.Change
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
		}, time.Second*1, time.Millisecond*10)

		locationHistoryFromLastChange := &model_generated.LocationHistory{}
		err = locationHistoryFromLastChange.FromItem(lastChange.Item)
		require.NoError(t, err)

		log.Printf("locationHistoryFromLastChange: %v", _helpers.UnsafeJSONPrettyFormat(locationHistoryFromLastChange))

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
			tx, _ := db.BeginTxx(ctx, nil)
			defer tx.Rollback()
			err = locationHistoryFromLastChange.Reload(ctx, tx)
			require.NoError(t, err)
			require.Equal(t, locationHistory.ParentPhysicalThingIDObject, locationHistoryFromLastChange.ParentPhysicalThingIDObject)
			_ = tx.Commit()
		}()

		log.Printf("locationHistoryFromLastChangeAfterReloading: %v", _helpers.UnsafeJSONPrettyFormat(locationHistoryFromLastChange))
	})

	t.Run("SelectWithPolygon", func(t *testing.T) {
		physicalExternalID := "SomePhysicalThingExternalID4"
		physicalThingName := "SomePhysicalThingName4"
		physicalThingType := "SomePhysicalThingType4"
		physicalThingTags := `'{tag1,tag2,tag3,"isn''t this, \"complicated\""}'`
		physicalThingMetadata := `'key1=>1, key2=>"a", key3=>true, key4=>NULL, key5=>"isn''t this, \"complicated\""'`
		physicalThingRawData := `'{"key1": 1, "key2": "a", "key3": true, "key4": null, "key5": "isn''t this, \"complicated\""}'`

		cleanup := func() {
			_, err = db.ExecContext(
				ctx,
				`DELETE FROM location_history;`,
			)
			_, err = db.ExecContext(
				ctx,
				`DELETE FROM physical_things CASCADE
			WHERE
				name = $1;`,
				physicalThingName,
			)
			require.NoError(t, err)
		}
		cleanup()

		_, err = db.ExecContext(
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
		var err error

		func() {
			tx, _ := db.BeginTxx(ctx, nil)
			defer tx.Rollback()
			physicalThing, err = model_generated.SelectPhysicalThing(
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

			_, err = db.ExecContext(
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

			locationHistory, err = model_generated.SelectLocationHistory(
				ctx,
				tx,
				"parent_physical_thing_id = $1",
				physicalThing.ID.String(),
			)
			require.NoError(t, err)
			require.NotNil(t, locationHistory)
			_ = tx.Commit()
			require.NotNil(t, physicalThing)
		}()

		log.Printf("locationHistory: %v", _helpers.UnsafeJSONPrettyFormat(locationHistory))

		require.IsType(t, uuid.UUID{}, locationHistory.ID, "ID")
		require.IsType(t, time.Time{}, locationHistory.CreatedAt, "CreatedAt")
		require.IsType(t, time.Time{}, locationHistory.UpdatedAt, "UpdatedAt")
		require.IsType(t, helpers.Nil(time.Time{}), locationHistory.DeletedAt, "DeletedAt")
		require.IsType(t, time.Time{}, locationHistory.Timestamp, "Timestamp")
		require.IsType(t, &pgtype.Point{}, locationHistory.Point, "Point")
		require.IsType(t, &pgtype.Polygon{}, locationHistory.Polygon, "Polygon")
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

		var lastChange stream.Change
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
		}, time.Second*1, time.Millisecond*10)

		locationHistoryFromLastChange := &model_generated.LocationHistory{}
		err = locationHistoryFromLastChange.FromItem(lastChange.Item)
		require.NoError(t, err)

		log.Printf("locationHistoryFromLastChange: %v", _helpers.UnsafeJSONPrettyFormat(locationHistoryFromLastChange))

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
			tx, _ := db.BeginTxx(ctx, nil)
			defer tx.Rollback()
			err = locationHistoryFromLastChange.Reload(ctx, tx)
			require.NoError(t, err)
			require.Equal(t, locationHistory.ParentPhysicalThingIDObject, locationHistoryFromLastChange.ParentPhysicalThingIDObject)
			_ = tx.Commit()
		}()

		log.Printf("locationHistoryFromLastChangeAfterReloading: %v", _helpers.UnsafeJSONPrettyFormat(locationHistoryFromLastChange))
	})
}
