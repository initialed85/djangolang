package model_reference

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
	"github.com/initialed85/djangolang/pkg/stream"
	"github.com/stretchr/testify/require"
)

func TestLogicalThings(t *testing.T) {
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
		os.Setenv("DJANGOLANG_NODE_NAME", "model_reference_logical_thing")
		err = stream.Run(ctx, changes, tableByName)
		require.NoError(t, err)
	}()
	runtime.Gosched()

	time.Sleep(time.Second * 1)

	t.Run("Select", func(t *testing.T) {
		physicalExternalID := "ReferenceSomePhysicalThingExternalID"
		physicalThingName := "ReferenceSomePhysicalThingName"
		physicalThingType := "ReferenceSomePhysicalThingType"
		logicalExternalID := "ReferenceSomeLogicalThingExternalID"
		logicalThingName := "ReferenceSomeLogicalThingName"
		logicalThingType := "ReferenceSomeLogicalThingType"
		physicalAndLogicalThingTags := `'{tag1,tag2,tag3,"isn''t this, \"complicated\""}'`
		physicalAndLogicalThingMetadata := `'key1=>1, key2=>"a", key3=>true, key4=>NULL, key5=>"isn''t this, \"complicated\""'`
		physicalAndLogicalThingRawData := `'{"key1": 1, "key2": "a", "key3": true, "key4": null, "key5": "isn''t this, \"complicated\""}'`

		cleanup := func() {
			_, err = db.ExecContext(
				ctx,
				`DELETE FROM logical_things
			WHERE
				name = $1;`,
				logicalThingName,
			)
			require.NoError(t, err)
			_, err = db.ExecContext(
				ctx,
				`DELETE FROM physical_things
			WHERE
				name = $1;`,
				physicalThingName,
			)
			require.NoError(t, err)
		}
		defer cleanup()

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
				physicalAndLogicalThingTags,
				physicalAndLogicalThingMetadata,
				physicalAndLogicalThingRawData,
			),
		)
		require.NoError(t, err)

		var logicalThing *LogicalThing
		var physicalThing *PhysicalThing
		var err error

		func() {
			tx, _ := db.BeginTxx(ctx, nil)
			defer tx.Rollback()
			physicalThing, err = SelectPhysicalThing(
				ctx,
				tx,
				fmt.Sprintf(
					"%v = $1 AND %v = $2 AND %v = $3",
					PhysicalThingTableExternalIDColumn,
					PhysicalThingTableNameColumn,
					PhysicalThingTableTypeColumn,
				),
				physicalExternalID,
				physicalThingName,
				physicalThingType,
			)
			require.NoError(t, err)
			require.NotNil(t, physicalThing)

			_, err = db.ExecContext(
				ctx,
				fmt.Sprintf(`INSERT INTO logical_things (
				external_id,
				name,
				type,
				tags,
				metadata,
				raw_data,
				parent_physical_thing_id
			)
			VALUES (
				'%v',
				'%v',
				'%v',
				%v,
				%v,
				%v,
				'%v'
			);`,
					logicalExternalID,
					logicalThingName,
					logicalThingType,
					physicalAndLogicalThingTags,
					physicalAndLogicalThingMetadata,
					physicalAndLogicalThingRawData,
					physicalThing.ID.String(),
				),
			)
			require.NoError(t, err)

			logicalThing, err = SelectLogicalThing(
				ctx,
				tx,
				"external_id = $1 AND name = $2 AND type = $3",
				logicalExternalID,
				logicalThingName,
				logicalThingType,
			)
			require.NoError(t, err)
			require.NotNil(t, logicalThing)
			_ = tx.Commit()
			require.NotNil(t, physicalThing)
		}()

		log.Printf("logicalThing: %v", _helpers.UnsafeJSONPrettyFormat(logicalThing))

		require.IsType(t, uuid.UUID{}, logicalThing.ID, "ID")
		require.IsType(t, time.Time{}, logicalThing.CreatedAt, "CreatedAt")
		require.IsType(t, time.Time{}, logicalThing.UpdatedAt, "UpdatedAt")
		require.IsType(t, helpers.Nil(time.Time{}), logicalThing.DeletedAt, "DeletedAt")
		require.IsType(t, helpers.Ptr(""), logicalThing.ExternalID, "ExternalID")
		require.IsType(t, "", logicalThing.Name, "Name")
		require.IsType(t, "", logicalThing.Type, "Type")
		require.IsType(t, []string{}, logicalThing.Tags, "Tags")
		require.IsType(t, map[string]*string{}, logicalThing.Metadata, "Metadata")
		require.IsType(t, new(any), logicalThing.RawData, "RawData")
		require.IsType(t, helpers.Ptr(uuid.UUID{}), logicalThing.ParentPhysicalThingID, "ID")
		require.IsType(t, helpers.Nil(uuid.UUID{}), logicalThing.ParentLogicalThingID, "ID")

		require.IsType(t, uuid.UUID{}, logicalThing.ParentPhysicalThingIDObject.ID, "ID")
		require.IsType(t, time.Time{}, logicalThing.ParentPhysicalThingIDObject.CreatedAt, "CreatedAt")
		require.IsType(t, time.Time{}, logicalThing.ParentPhysicalThingIDObject.UpdatedAt, "UpdatedAt")
		require.IsType(t, helpers.Nil(time.Time{}), logicalThing.ParentPhysicalThingIDObject.DeletedAt, "DeletedAt")
		require.IsType(t, helpers.Ptr(""), logicalThing.ParentPhysicalThingIDObject.ExternalID, "ExternalID")
		require.IsType(t, "", logicalThing.ParentPhysicalThingIDObject.Name, "Name")
		require.IsType(t, "", logicalThing.ParentPhysicalThingIDObject.Type, "Type")
		require.IsType(t, []string{}, logicalThing.ParentPhysicalThingIDObject.Tags, "Tags")
		require.IsType(t, map[string]*string{}, logicalThing.ParentPhysicalThingIDObject.Metadata, "Metadata")
		require.IsType(t, new(any), logicalThing.ParentPhysicalThingIDObject.RawData, "ParentPhysicalThingIDObject")

		require.IsType(t, helpers.Nil(LogicalThing{}), logicalThing.ParentLogicalThingIDObject, "ParentLogicalThingIDObject")

		var lastChange stream.Change
		require.Eventually(t, func() bool {
			mu.Lock()
			defer mu.Unlock()

			var ok bool
			lastChange, ok = lastChangeByTableName[LogicalThingTable]
			if !ok {
				return false
			}

			if lastChange.Action != stream.INSERT {
				return false
			}

			return true
		}, time.Second*1, time.Millisecond*10)

		logicalThingFromLastChange := &LogicalThing{}
		err = logicalThingFromLastChange.FromItem(lastChange.Item)
		require.NoError(t, err)

		log.Printf("logicalThingFromLastChange: %v", _helpers.UnsafeJSONPrettyFormat(logicalThingFromLastChange))

		require.Equal(t, logicalThing.ID, logicalThingFromLastChange.ID)
		require.Equal(t, logicalThing.CreatedAt.Unix(), logicalThingFromLastChange.CreatedAt.Unix())
		require.Equal(t, logicalThing.UpdatedAt.Unix(), logicalThingFromLastChange.UpdatedAt.Unix())
		require.Equal(t, logicalThing.DeletedAt, logicalThingFromLastChange.DeletedAt)
		require.Equal(t, logicalThing.ExternalID, logicalThingFromLastChange.ExternalID)
		require.Equal(t, logicalThing.Name, logicalThingFromLastChange.Name)
		require.Equal(t, logicalThing.Type, logicalThingFromLastChange.Type)
		require.Equal(t, logicalThing.Tags, logicalThingFromLastChange.Tags)
		require.Equal(t, logicalThing.Metadata, logicalThingFromLastChange.Metadata)
		require.Equal(t, logicalThing.RawData, logicalThingFromLastChange.RawData)
		require.Equal(t, logicalThing.ParentPhysicalThingID, logicalThingFromLastChange.ParentPhysicalThingID)
		require.Equal(t, logicalThing.ParentLogicalThingID, logicalThingFromLastChange.ParentLogicalThingID)
		require.NotEqual(t, logicalThing.ParentPhysicalThingIDObject, logicalThingFromLastChange.ParentPhysicalThingIDObject)
		require.Equal(t, logicalThing.ParentLogicalThingIDObject, logicalThingFromLastChange.ParentLogicalThingIDObject)

		func() {
			tx, _ := db.BeginTxx(ctx, nil)
			defer tx.Rollback()
			err = logicalThingFromLastChange.Reload(ctx, tx)
			require.NoError(t, err)
			require.Equal(t, logicalThing.ParentPhysicalThingIDObject, logicalThingFromLastChange.ParentPhysicalThingIDObject)
			_ = tx.Commit()
		}()

		log.Printf("logicalThingFromLastChangeAfterReloading: %v", _helpers.UnsafeJSONPrettyFormat(logicalThingFromLastChange))
	})
}
