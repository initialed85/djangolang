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
	"github.com/google/uuid"
	"github.com/initialed85/djangolang/internal/hack"
	"github.com/initialed85/djangolang/pkg/helpers"
	"github.com/initialed85/djangolang/pkg/model_generated"
	"github.com/initialed85/djangolang/pkg/query"
	"github.com/initialed85/djangolang/pkg/server"
	"github.com/initialed85/djangolang/pkg/stream"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

func testLogicalThings(
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

	t.Run("Select", func(t *testing.T) {
		physicalExternalID := "SelectSomePhysicalThingExternalID"
		physicalThingName := "SelectSomePhysicalThingName"
		physicalThingType := "SelectSomePhysicalThingType"
		logicalExternalID := "SelectSomeLogicalThingExternalID"
		logicalThingName := "SelectSomeLogicalThingName"
		logicalThingType := "SelectSomeLogicalThingType"
		logicalThingAge := time.Millisecond * 1111
		logicalThingCount := 0
		physicalAndLogicalThingTags := `'{tag1,tag2,tag3,"isn''t this, \"complicated\""}'`
		physicalAndLogicalThingMetadata := `'key1=>1, key2=>"a", key3=>true, key4=>NULL, key5=>"isn''t this, \"complicated\""'`
		physicalAndLogicalThingRawData := `'{"key1": 1, "key2": "a", "key3": true, "key4": null, "key5": "isn''t this, \"complicated\""}'`

		cleanup := func() {
			_, _ = db.Exec(
				ctx,
				`DELETE FROM logical_things
			WHERE
				name = $1;`,
				logicalThingName,
			)
			_, _ = db.Exec(
				ctx,
				`DELETE FROM physical_things
			WHERE
				name = $1;`,
				physicalThingName,
			)
			if redisConn != nil {
				_, _ = redisConn.Do("FLUSHALL")
			}
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
				physicalAndLogicalThingTags,
				physicalAndLogicalThingMetadata,
				physicalAndLogicalThingRawData,
			),
		)
		require.NoError(t, err)

		var logicalThing *model_generated.LogicalThing
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
				fmt.Sprintf(`INSERT INTO logical_things (
				external_id,
				name,
				type,
				tags,
				metadata,
				raw_data,
				age,
				optional_age,
				count,
				parent_physical_thing_id
			)
			VALUES (
				'%v',
				'%v',
				'%v',
				%v,
				%v,
				%v,
				interval '%v microseconds',
				interval '%v microseconds',
				%v,
				'%v'
			);`,
					logicalExternalID,
					logicalThingName,
					logicalThingType,
					physicalAndLogicalThingTags,
					physicalAndLogicalThingMetadata,
					physicalAndLogicalThingRawData,
					logicalThingAge.Microseconds(),
					logicalThingAge.Microseconds(),
					logicalThingCount,
					physicalThing.ID.String(),
				),
			)
			require.NoError(t, err)

			logicalThing, count, totalCount, page, totalPages, err = model_generated.SelectLogicalThing(
				ctx,
				tx,
				"external_id = $1 AND name = $2 AND type = $3",
				logicalExternalID,
				logicalThingName,
				logicalThingType,
			)
			require.NoError(t, err)
			require.NotNil(t, logicalThing)
			_ = tx.Commit(ctx)
			require.NotNil(t, physicalThing)
		}()

		log.Printf("logicalThing: %v", hack.UnsafeJSONPrettyFormat(logicalThing))

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

		require.IsType(t, helpers.Nil(model_generated.LogicalThing{}), logicalThing.ParentLogicalThingIDObject, "ParentLogicalThingIDObject")

		var lastChange server.Change
		require.Eventually(t, func() bool {
			mu.Lock()
			defer mu.Unlock()

			var ok bool
			lastChange, ok = lastChangeByTableName[model_generated.LogicalThingTable]
			if !ok {
				return false
			}

			if lastChange.Action != stream.INSERT {
				return false
			}

			return true
		}, time.Second*10, time.Millisecond*10)

		logicalThingFromLastChange := &model_generated.LogicalThing{}
		err = logicalThingFromLastChange.FromItem(lastChange.Item)
		require.NoError(t, err)

		log.Printf("logicalThingFromLastChange: %v", hack.UnsafeJSONPrettyFormat(logicalThingFromLastChange))

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
			tx, _ := db.Begin(ctx)
			defer tx.Rollback(ctx)
			err = logicalThingFromLastChange.Reload(ctx, tx)
			require.NoError(t, err)
			require.Equal(t, logicalThing.ParentPhysicalThingIDObject, logicalThingFromLastChange.ParentPhysicalThingIDObject)
			_ = tx.Commit(ctx)
		}()

		log.Printf("logicalThingFromLastChangeAfterReloading: %v", hack.UnsafeJSONPrettyFormat(logicalThingFromLastChange))
	})

	t.Run("Insert", func(t *testing.T) {
		physicalExternalID := "InsertSomePhysicalThingExternalID"
		physicalThingName := "InsertSomePhysicalThingName"
		physicalThingType := "InsertSomePhysicalThingType"
		logicalExternalID := "InsertSomeLogicalThingExternalID"
		logicalThingName := "InsertSomeLogicalThingName"
		logicalThingType := "InsertSomeLogicalThingType"
		logicalThingAge := time.Millisecond * 2222
		physicalAndLogicalThingTags := []string{"tag1", "tag2", "tag3", "isn''t this, \"complicated\""}
		physicalAndLogicalThingMetadata := map[string]*string{
			"key1": helpers.Ptr("1"),
			"key2": helpers.Ptr("a"),
			"key3": helpers.Ptr("true"),
			"key4": nil,
			"key5": helpers.Ptr("isn''t this, \"complicated\""),
		}
		physicalAndLogicalThingRawData := map[string]any{
			"key1": "1",
			"key2": "a",
			"key3": true,
			"key4": nil,
			"key5": "isn''t this, \"complicated\"",
		}

		cleanup := func() {
			_, _ = db.Exec(
				ctx,
				`DELETE FROM logical_things
			WHERE
				name = $1;`,
				logicalThingName,
			)
			_, _ = db.Exec(
				ctx,
				`DELETE FROM physical_things
			WHERE
				name = $1;`,
				physicalThingName,
			)
			if redisConn != nil {
				_, _ = redisConn.Do("FLUSHALL")
			}

		}
		defer cleanup()

		physicalThing := &model_generated.PhysicalThing{
			ExternalID: &physicalExternalID,
			Name:       physicalThingName,
			Type:       physicalThingType,
			Tags:       physicalAndLogicalThingTags,
			Metadata:   physicalAndLogicalThingMetadata,
			RawData:    physicalAndLogicalThingRawData,
		}

		logicalThing := &model_generated.LogicalThing{
			ExternalID:  &logicalExternalID,
			Name:        logicalThingName,
			Type:        logicalThingType,
			Tags:        physicalAndLogicalThingTags,
			Metadata:    physicalAndLogicalThingMetadata,
			RawData:     physicalAndLogicalThingRawData,
			Age:         logicalThingAge,
			OptionalAge: &logicalThingAge,
		}

		func() {
			tx, _ := db.Begin(ctx)
			defer tx.Rollback(ctx)
			err := physicalThing.Insert(
				ctx,
				tx,
				false,
				false,
			)
			require.NoError(t, err)
			require.NotNil(t, physicalThing)

			logicalThing.ParentPhysicalThingID = helpers.Ptr(physicalThing.ID)
			err = logicalThing.Insert(
				ctx,
				tx,
				false,
				false,
			)
			require.NoError(t, err)
			require.NotNil(t, logicalThing)

			_ = tx.Commit(ctx)
		}()

		log.Printf("logicalThing: %v", hack.UnsafeJSONPrettyFormat(logicalThing))

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
		require.IsType(t, time.Duration(0), logicalThing.Age, "Age")
		require.IsType(t, helpers.Ptr(time.Duration(0)), logicalThing.OptionalAge, "Age")
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

		require.IsType(t, helpers.Nil(model_generated.LogicalThing{}), logicalThing.ParentLogicalThingIDObject, "ParentLogicalThingIDObject")

		require.Equal(t, &logicalExternalID, logicalThing.ExternalID)
		require.Equal(t, logicalThingName, logicalThing.Name)
		require.Equal(t, logicalThingType, logicalThing.Type)
		require.Equal(t, physicalAndLogicalThingTags, logicalThing.Tags)
		require.Equal(t, physicalAndLogicalThingMetadata, logicalThing.Metadata)

		// TODO
		// require.Equal(t, physicalAndLogicalThingRawData, logicalThing.RawData)

		require.Equal(t, logicalThingAge, logicalThing.Age)
		require.Equal(t, &logicalThingAge, logicalThing.OptionalAge)

		var lastChange server.Change
		require.Eventually(t, func() bool {
			mu.Lock()
			defer mu.Unlock()

			var ok bool
			lastChange, ok = lastChangeByTableName[model_generated.LogicalThingTable]
			if !ok {
				return false
			}

			if lastChange.Action != stream.INSERT {
				return false
			}

			return true
		}, time.Second*10, time.Millisecond*10)

		logicalThingFromLastChange := &model_generated.LogicalThing{}
		err := logicalThingFromLastChange.FromItem(lastChange.Item)
		require.NoError(t, err)

		log.Printf("logicalThingFromLastChange: %v", hack.UnsafeJSONPrettyFormat(logicalThingFromLastChange))

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
		require.Equal(t, logicalThing.Age, logicalThingFromLastChange.Age)
		require.Equal(t, logicalThing.ParentPhysicalThingID, logicalThingFromLastChange.ParentPhysicalThingID)
		require.Equal(t, logicalThing.ParentLogicalThingID, logicalThingFromLastChange.ParentLogicalThingID)
		require.NotEqual(t, logicalThing.ParentPhysicalThingIDObject, logicalThingFromLastChange.ParentPhysicalThingIDObject)
		require.Equal(t, logicalThing.ParentLogicalThingIDObject, logicalThingFromLastChange.ParentLogicalThingIDObject)

		func() {
			tx, _ := db.Begin(ctx)
			defer tx.Rollback(ctx)
			err = logicalThingFromLastChange.Reload(ctx, tx)
			require.NoError(t, err)
			require.Equal(t, logicalThing.ParentPhysicalThingIDObject, logicalThingFromLastChange.ParentPhysicalThingIDObject)
			_ = tx.Commit(ctx)
		}()

		log.Printf("logicalThingFromLastChangeAfterReloading: %v", hack.UnsafeJSONPrettyFormat(logicalThingFromLastChange))
	})

	t.Run("Update", func(t *testing.T) {
		insertPhysicalExternalID := "UpdateSomePhysicalThingExternalID"
		insertPhysicalThingName := "UpdateSomePhysicalThingName"
		insertPhysicalThingType := "UpdateSomePhysicalThingType"
		insertLogicalExternalID := "UpdateSomeLogicalThingExternalID"
		insertLogicalThingName := "UpdateSomeLogicalThingName"
		insertLogicalThingType := "UpdateSomeLogicalThingType"
		insertPhysicalAndLogicalThingTags := []string{"tag1", "tag2", "tag3", "isn''t this, \"complicated\""}
		insertPhysicalAndLogicalThingMetadata := map[string]*string{
			"key1": helpers.Ptr("1"),
			"key2": helpers.Ptr("a"),
			"key3": helpers.Ptr("true"),
			"key4": nil,
			"key5": helpers.Ptr("isn''t this, \"complicated\""),
		}
		insertPhysicalAndLogicalThingRawData := map[string]any{
			"key1": "1",
			"key2": "a",
			"key3": true,
			"key4": nil,
			"key5": "isn''t this, \"complicated\"",
		}
		insertLogicalThingAge := time.Millisecond * 3333

		updateLogicalExternalID := "UpdateSomeLogicalThingExternalID2"
		updateLogicalThingName := "UpdateSomeLogicalThingName2"
		updateLogicalThingType := "UpdateSomeLogicalThingType2"
		updatePhysicalAndLogicalThingTags := []string{"tag1", "tag2", "tag3", "isn''t this, \"complicated\"", "ham"}
		updatePhysicalAndLogicalThingMetadata := map[string]*string{
			"key1": helpers.Ptr("1"),
			"key2": helpers.Ptr("a"),
			"key3": helpers.Ptr("true"),
			"key4": nil,
			"key5": helpers.Ptr("isn''t this, \"complicated\""),
			"key6": helpers.Ptr("ham"),
		}
		updatePhysicalAndLogicalThingRawData := map[string]any{
			"key1": "1",
			"key2": "a",
			"key3": true,
			"key4": nil,
			"key5": "isn''t this, \"complicated\"",
			"key6": "ham",
		}
		updateLogicalThingAge := time.Millisecond * 4444

		cleanup := func() {
			_, _ = db.Exec(
				ctx,
				`DELETE FROM logical_things
			WHERE
				name = $1 OR name = $2;`,
				insertLogicalThingName,
				updateLogicalThingName,
			)
			_, _ = db.Exec(
				ctx,
				`DELETE FROM physical_things
			WHERE
				name = $1;`,
				insertPhysicalThingName,
			)
			if redisConn != nil {
				_, _ = redisConn.Do("FLUSHALL")
			}

		}
		defer cleanup()

		var err error

		physicalThing := &model_generated.PhysicalThing{
			ExternalID: &insertPhysicalExternalID,
			Name:       insertPhysicalThingName,
			Type:       insertPhysicalThingType,
			Tags:       insertPhysicalAndLogicalThingTags,
			Metadata:   insertPhysicalAndLogicalThingMetadata,
			RawData:    insertPhysicalAndLogicalThingRawData,
		}

		logicalThing := &model_generated.LogicalThing{
			ExternalID:  &insertLogicalExternalID,
			Name:        insertLogicalThingName,
			Type:        insertLogicalThingType,
			Tags:        insertPhysicalAndLogicalThingTags,
			Metadata:    insertPhysicalAndLogicalThingMetadata,
			RawData:     insertPhysicalAndLogicalThingRawData,
			Age:         insertLogicalThingAge,
			OptionalAge: &insertLogicalThingAge,
		}

		func() {
			tx, _ := db.Begin(ctx)
			defer tx.Rollback(ctx)
			err = physicalThing.Insert(
				ctx,
				tx,
				false,
				false,
			)
			require.NoError(t, err)
			require.NotNil(t, physicalThing)

			logicalThing.ParentPhysicalThingID = helpers.Ptr(physicalThing.ID)
			err = logicalThing.Insert(
				ctx,
				tx,
				false,
				false,
			)
			require.NoError(t, err)
			require.NotNil(t, logicalThing)

			logicalThing.ExternalID = &updateLogicalExternalID
			logicalThing.Name = updateLogicalThingName
			logicalThing.Type = updateLogicalThingType
			logicalThing.Tags = updatePhysicalAndLogicalThingTags
			logicalThing.Metadata = updatePhysicalAndLogicalThingMetadata
			logicalThing.RawData = updatePhysicalAndLogicalThingRawData
			logicalThing.Age = updateLogicalThingAge
			logicalThing.OptionalAge = &updateLogicalThingAge
			err = logicalThing.Update(ctx, tx, false)
			require.NoError(t, err)
			require.NotNil(t, logicalThing)

			_ = tx.Commit(ctx)
		}()

		log.Printf("logicalThing: %v", hack.UnsafeJSONPrettyFormat(logicalThing))

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
		require.IsType(t, time.Duration(0), logicalThing.Age, "Age")
		require.IsType(t, helpers.Ptr(time.Duration(0)), logicalThing.OptionalAge, "Age")
		require.IsType(t, helpers.Ptr(uuid.UUID{}), logicalThing.ParentPhysicalThingID, "ParentPhysicalThingID")
		require.IsType(t, helpers.Nil(uuid.UUID{}), logicalThing.ParentLogicalThingID, "ParentLogicalThingID")

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

		require.IsType(t, helpers.Nil(model_generated.LogicalThing{}), logicalThing.ParentLogicalThingIDObject, "ParentLogicalThingIDObject")

		require.Equal(t, &updateLogicalExternalID, logicalThing.ExternalID)
		require.Equal(t, updateLogicalThingName, logicalThing.Name)
		require.Equal(t, updateLogicalThingType, logicalThing.Type)
		require.Equal(t, updatePhysicalAndLogicalThingTags, logicalThing.Tags)
		require.Equal(t, updatePhysicalAndLogicalThingMetadata, logicalThing.Metadata)

		// TODO
		// require.Equal(t, physicalAndLogicalThingRawData, logicalThing.RawData)

		require.Equal(t, updateLogicalThingAge, logicalThing.Age)
		require.Equal(t, &updateLogicalThingAge, logicalThing.OptionalAge)

		var lastChange server.Change
		require.Eventually(t, func() bool {
			mu.Lock()
			defer mu.Unlock()

			var ok bool
			lastChange, ok = lastChangeByTableName[model_generated.LogicalThingTable]
			if !ok {
				return false
			}

			if lastChange.Action != stream.UPDATE {
				return false
			}

			return true
		}, time.Second*10, time.Millisecond*10)

		logicalThingFromLastChange := &model_generated.LogicalThing{}
		err = logicalThingFromLastChange.FromItem(lastChange.Item)
		require.NoError(t, err)

		log.Printf("logicalThingFromLastChange: %v", hack.UnsafeJSONPrettyFormat(logicalThingFromLastChange))

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
		require.Equal(t, logicalThing.Age, logicalThingFromLastChange.Age)
		require.Equal(t, logicalThing.ParentPhysicalThingID, logicalThingFromLastChange.ParentPhysicalThingID)
		require.Equal(t, logicalThing.ParentLogicalThingID, logicalThingFromLastChange.ParentLogicalThingID)
		require.NotEqual(t, logicalThing.ParentPhysicalThingIDObject, logicalThingFromLastChange.ParentPhysicalThingIDObject)
		require.Equal(t, logicalThing.ParentLogicalThingIDObject, logicalThingFromLastChange.ParentLogicalThingIDObject)

		func() {
			tx, _ := db.Begin(ctx)
			defer tx.Rollback(ctx)
			err = logicalThingFromLastChange.Reload(ctx, tx)
			require.NoError(t, err)
			require.Equal(t, logicalThing.ParentPhysicalThingIDObject, logicalThingFromLastChange.ParentPhysicalThingIDObject)
			_ = tx.Commit(ctx)
		}()

		log.Printf("logicalThingFromLastChangeAfterReloading: %v", hack.UnsafeJSONPrettyFormat(logicalThingFromLastChange))
	})

	t.Run("SoftDelete", func(t *testing.T) {
		physicalExternalID := "SoftDeleteSomePhysicalThingExternalID"
		physicalThingName := "SoftDeleteSomePhysicalThingName"
		physicalThingType := "SoftDeleteSomePhysicalThingType"
		logicalExternalID := "SoftDeleteSomeLogicalThingExternalID"
		logicalThingName := "SoftDeleteSomeLogicalThingName"
		logicalThingType := "SoftDeleteSomeLogicalThingType"
		physicalAndLogicalThingTags := []string{"tag1", "tag2", "tag3", "isn''t this, \"complicated\""}
		physicalAndLogicalThingMetadata := map[string]*string{
			"key1": helpers.Ptr("1"),
			"key2": helpers.Ptr("a"),
			"key3": helpers.Ptr("true"),
			"key4": nil,
			"key5": helpers.Ptr("isn''t this, \"complicated\""),
		}
		physicalAndLogicalThingRawData := map[string]any{
			"key1": "1",
			"key2": "a",
			"key3": true,
			"key4": nil,
			"key5": "isn''t this, \"complicated\"",
		}

		cleanup := func() {
			_, _ = db.Exec(
				ctx,
				`DELETE FROM logical_things
			WHERE
				name = $1;`,
				logicalThingName,
			)
			_, _ = db.Exec(
				ctx,
				`DELETE FROM physical_things
			WHERE
				name = $1;`,
				physicalThingName,
			)
			if redisConn != nil {
				_, _ = redisConn.Do("FLUSHALL")
			}

		}
		defer cleanup()

		var err error

		physicalThing := &model_generated.PhysicalThing{
			ExternalID: &physicalExternalID,
			Name:       physicalThingName,
			Type:       physicalThingType,
			Tags:       physicalAndLogicalThingTags,
			Metadata:   physicalAndLogicalThingMetadata,
			RawData:    physicalAndLogicalThingRawData,
		}

		logicalThing := &model_generated.LogicalThing{
			ExternalID: &logicalExternalID,
			Name:       logicalThingName,
			Type:       logicalThingType,
			Tags:       physicalAndLogicalThingTags,
			Metadata:   physicalAndLogicalThingMetadata,
			RawData:    physicalAndLogicalThingRawData,
		}

		func() {
			tx, _ := db.Begin(ctx)
			defer tx.Rollback(ctx)
			err = physicalThing.Insert(
				ctx,
				tx,
				false,
				false,
			)
			require.NoError(t, err)
			require.NotNil(t, physicalThing)

			logicalThing.ParentPhysicalThingID = helpers.Ptr(physicalThing.ID)
			err = logicalThing.Insert(
				ctx,
				tx,
				false,
				false,
			)
			require.NoError(t, err)
			require.NotNil(t, logicalThing)

			err = logicalThing.Delete(ctx, tx)
			require.NoError(t, err)

			_ = tx.Commit(ctx)
		}()

		log.Printf("logicalThing: %v", hack.UnsafeJSONPrettyFormat(logicalThing))

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

		require.IsType(t, helpers.Nil(model_generated.LogicalThing{}), logicalThing.ParentLogicalThingIDObject, "ParentLogicalThingIDObject")

		var lastChange server.Change
		require.Eventually(t, func() bool {
			mu.Lock()
			defer mu.Unlock()

			var ok bool
			lastChange, ok = lastChangeByTableName[model_generated.LogicalThingTable]
			if !ok {
				return false
			}

			if lastChange.Action != stream.SOFT_DELETE {
				return false
			}

			logicalThingFromLastChange := &model_generated.LogicalThing{}
			err = logicalThingFromLastChange.FromItem(lastChange.Item)
			if err != nil {
				return false
			}

			if logicalThingFromLastChange.ID != logicalThing.ID {
				return false
			}

			return true
		}, time.Second*10, time.Millisecond*10)

		logicalThingFromLastChange := &model_generated.LogicalThing{}
		err = logicalThingFromLastChange.FromItem(lastChange.Item)
		require.NoError(t, err)

		log.Printf("logicalThingFromLastChange: %v", hack.UnsafeJSONPrettyFormat(logicalThingFromLastChange))

		require.Equal(t, logicalThing.ID, logicalThingFromLastChange.ID)
		require.NotNil(t, logicalThing.DeletedAt)
		require.NotNil(t, logicalThingFromLastChange.DeletedAt)

		func() {
			tx, _ := db.Begin(ctx)
			defer tx.Rollback(ctx)
			err = logicalThingFromLastChange.Reload(ctx, tx)
			require.Error(t, err)
			_ = tx.Commit(ctx)
		}()

		log.Printf("logicalThingFromLastChangeAfterReloading: %v", hack.UnsafeJSONPrettyFormat(logicalThingFromLastChange))
	})

	t.Run("HardDelete", func(t *testing.T) {
		physicalExternalID := "HardDeleteSomePhysicalThingExternalID"
		physicalThingName := "HardDeleteSomePhysicalThingName"
		physicalThingType := "HardDeleteSomePhysicalThingType"
		logicalExternalID := "HardDeleteSomeLogicalThingExternalID"
		logicalThingName := "HardDeleteSomeLogicalThingName"
		logicalThingType := "HardDeleteSomeLogicalThingType"
		physicalAndLogicalThingTags := []string{"tag1", "tag2", "tag3", "isn''t this, \"complicated\""}
		physicalAndLogicalThingMetadata := map[string]*string{
			"key1": helpers.Ptr("1"),
			"key2": helpers.Ptr("a"),
			"key3": helpers.Ptr("true"),
			"key4": nil,
			"key5": helpers.Ptr("isn''t this, \"complicated\""),
		}
		physicalAndLogicalThingRawData := map[string]any{
			"key1": "1",
			"key2": "a",
			"key3": true,
			"key4": nil,
			"key5": "isn''t this, \"complicated\"",
		}

		cleanup := func() {
			_, _ = db.Exec(
				ctx,
				`DELETE FROM logical_things
			WHERE
				name = $1;`,
				logicalThingName,
			)
			_, _ = db.Exec(
				ctx,
				`DELETE FROM physical_things
			WHERE
				name = $1;`,
				physicalThingName,
			)
			if redisConn != nil {
				_, _ = redisConn.Do("FLUSHALL")
			}

		}
		defer cleanup()

		var err error

		physicalThing := &model_generated.PhysicalThing{
			ExternalID: &physicalExternalID,
			Name:       physicalThingName,
			Type:       physicalThingType,
			Tags:       physicalAndLogicalThingTags,
			Metadata:   physicalAndLogicalThingMetadata,
			RawData:    physicalAndLogicalThingRawData,
		}

		logicalThing := &model_generated.LogicalThing{
			ExternalID: &logicalExternalID,
			Name:       logicalThingName,
			Type:       logicalThingType,
			Tags:       physicalAndLogicalThingTags,
			Metadata:   physicalAndLogicalThingMetadata,
			RawData:    physicalAndLogicalThingRawData,
		}

		func() {
			tx, _ := db.Begin(ctx)
			defer tx.Rollback(ctx)
			err = physicalThing.Insert(
				ctx,
				tx,
				false,
				false,
			)
			require.NoError(t, err)
			require.NotNil(t, physicalThing)

			logicalThing.ParentPhysicalThingID = helpers.Ptr(physicalThing.ID)
			err = logicalThing.Insert(
				ctx,
				tx,
				false,
				false,
			)
			require.NoError(t, err)
			require.NotNil(t, logicalThing)

			err = logicalThing.Delete(ctx, tx, true)
			require.NoError(t, err)

			_ = tx.Commit(ctx)
		}()

		log.Printf("logicalThing: %v", hack.UnsafeJSONPrettyFormat(logicalThing))

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

		require.IsType(t, helpers.Nil(model_generated.LogicalThing{}), logicalThing.ParentLogicalThingIDObject, "ParentLogicalThingIDObject")

		var lastChange server.Change
		require.Eventually(t, func() bool {
			mu.Lock()
			defer mu.Unlock()

			var ok bool
			lastChange, ok = lastChangeByTableName[model_generated.LogicalThingTable]
			if !ok {
				return false
			}

			if lastChange.Action != stream.SOFT_DELETE {
				return false
			}

			logicalThingFromLastChange := &model_generated.LogicalThing{}
			err = logicalThingFromLastChange.FromItem(lastChange.Item)
			if err != nil {
				return false
			}

			if logicalThingFromLastChange.ID != logicalThing.ID {
				return false
			}

			return true
		}, time.Second*10, time.Millisecond*10)

		logicalThingFromLastChange := &model_generated.LogicalThing{}
		err = logicalThingFromLastChange.FromItem(lastChange.Item)
		require.NoError(t, err)

		log.Printf("logicalThingFromLastChange: %v", hack.UnsafeJSONPrettyFormat(logicalThingFromLastChange))

		require.Equal(t, logicalThing.ID, logicalThingFromLastChange.ID)
		require.NotNil(t, logicalThing.DeletedAt)
		require.NotNil(t, logicalThingFromLastChange.DeletedAt)

		func() {
			tx, _ := db.Begin(ctx)
			defer tx.Rollback(ctx)
			err = logicalThingFromLastChange.Reload(ctx, tx)
			require.Error(t, err)
			_ = tx.Commit(ctx)
		}()

		log.Printf("logicalThingFromLastChangeAfterReloading: %v", hack.UnsafeJSONPrettyFormat(logicalThingFromLastChange))
	})

	t.Run("LockTable", func(t *testing.T) {
		physicalExternalID := "LockTableSomePhysicalThingExternalID"
		physicalThingName := "LockTableSomePhysicalThingName"
		physicalThingType := "LockTableSomePhysicalThingType"
		logicalExternalID := "LockTableSomeLogicalThingExternalID"
		logicalThingName := "LockTableSomeLogicalThingName"
		logicalThingType := "LockTableSomeLogicalThingType"
		logicalThingAge := time.Millisecond * 1111
		logicalThingCount := 0
		physicalAndLogicalThingTags := `'{tag1,tag2,tag3,"isn''t this, \"complicated\""}'`
		physicalAndLogicalThingMetadata := `'key1=>1, key2=>"a", key3=>true, key4=>NULL, key5=>"isn''t this, \"complicated\""'`
		physicalAndLogicalThingRawData := `'{"key1": 1, "key2": "a", "key3": true, "key4": null, "key5": "isn''t this, \"complicated\""}'`

		cleanup := func() {
			_, _ = db.Exec(
				ctx,
				`DELETE FROM logical_things
			WHERE
				name = $1;`,
				logicalThingName,
			)
			_, _ = db.Exec(
				ctx,
				`DELETE FROM physical_things
			WHERE
				name = $1;`,
				physicalThingName,
			)
			if redisConn != nil {
				_, _ = redisConn.Do("FLUSHALL")
			}
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
				physicalAndLogicalThingTags,
				physicalAndLogicalThingMetadata,
				physicalAndLogicalThingRawData,
			),
		)
		require.NoError(t, err)

		var logicalThing *model_generated.LogicalThing
		var physicalThing *model_generated.PhysicalThing

		func() {
			ctx, cancel := context.WithTimeout(ctx, time.Second*5)
			defer cancel()

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
				fmt.Sprintf(`INSERT INTO logical_things (
				external_id,
				name,
				type,
				tags,
				metadata,
				raw_data,
				age,
				optional_age,
				count,
				parent_physical_thing_id
			)
			VALUES (
				'%v',
				'%v',
				'%v',
				%v,
				%v,
				%v,
				interval '%v microseconds',
				interval '%v microseconds',
				%v,
				'%v'
			);`,
					logicalExternalID,
					logicalThingName,
					logicalThingType,
					physicalAndLogicalThingTags,
					physicalAndLogicalThingMetadata,
					physicalAndLogicalThingRawData,
					logicalThingAge.Microseconds(),
					logicalThingAge.Microseconds(),
					logicalThingCount,
					physicalThing.ID.String(),
				),
			)
			require.NoError(t, err)

			logicalThing, count, totalCount, page, totalPages, err = model_generated.SelectLogicalThing(
				ctx,
				tx,
				"external_id = $1 AND name = $2 AND type = $3",
				logicalExternalID,
				logicalThingName,
				logicalThingType,
			)
			require.NoError(t, err)
			require.NotNil(t, logicalThing)

			checkSelect := func(shouldFail bool) {
				otherCtx, otherCancel := context.WithTimeout(ctx, time.Second*1)
				defer otherCancel()

				otherTx, err := db.Begin(otherCtx)
				require.NoError(t, err)

				defer func() {
					_ = otherTx.Rollback(otherCtx)
				}()

				err = logicalThing.Reload(otherCtx, otherTx)
				if shouldFail {
					require.Error(t, err)
				} else {
					require.NoError(t, err)
				}
			}

			checkLock := func(shouldFail bool) {
				otherCtx, otherCancel := context.WithTimeout(ctx, time.Second*1)
				defer otherCancel()

				otherTx, err := db.Begin(otherCtx)
				require.NoError(t, err)

				defer func() {
					_ = otherTx.Rollback(otherCtx)
				}()

				err = logicalThing.LockTable(
					otherCtx,
					otherTx,
					time.Second*0,
				)
				if shouldFail {
					require.Error(t, err)
				} else {
					require.NoError(t, err)
				}
			}

			err = tx.Commit(ctx)
			require.NoError(t, err)

			tx, err := db.Begin(ctx)
			require.NoError(t, err)
			defer func() {
				_ = tx.Rollback(ctx)
			}()

			err = logicalThing.LockTable(ctx, tx, time.Second*0)
			require.NoError(t, err)

			checkSelect(true)
			checkLock(true)

			err = tx.Commit(ctx)
			require.NoError(t, err)

			checkSelect(false)
			checkLock(false)
		}()
	})

	t.Run("RouterGetMany", func(t *testing.T) {
		physicalExternalID := "RouterGetManySomePhysicalThingExternalID"
		physicalThingName := "RouterGetManySomePhysicalThingName"
		physicalThingType := "RouterGetManySomePhysicalThingType"
		logicalExternalID := "RouterGetManySomeLogicalThingExternalID"
		logicalThingName := "RouterGetManySomeLogicalThingName"
		logicalThingType := "RouterGetManySomeLogicalThingType"
		physicalAndLogicalThingTags := []string{"tag1", "tag2", "tag3", "isn''t this, \"complicated\""}
		physicalAndLogicalThingMetadata := map[string]*string{
			"key1": helpers.Ptr("1"),
			"key2": helpers.Ptr("a"),
			"key3": helpers.Ptr("true"),
			"key4": nil,
			"key5": helpers.Ptr("isn''t this, \"complicated\""),
		}
		physicalAndLogicalThingRawData := map[string]any{
			"key1": "1",
			"key2": "a",
			"key3": true,
			"key4": nil,
			"key5": "isn''t this, \"complicated\"",
		}

		cleanup := func() {
			_, _ = db.Exec(
				ctx,
				`DELETE FROM logical_things
			WHERE
				name = $1;`,
				logicalThingName,
			)
			_, _ = db.Exec(
				ctx,
				`DELETE FROM logical_things
			WHERE
				name = $1;`,
				logicalThingName+"-2",
			)
			_, _ = db.Exec(
				ctx,
				`DELETE FROM physical_things
			WHERE
				name = $1;`,
				physicalThingName,
			)
			if redisConn != nil {
				_, _ = redisConn.Do("FLUSHALL")
			}

		}
		defer cleanup()

		var err error

		physicalThing := &model_generated.PhysicalThing{
			ExternalID: &physicalExternalID,
			Name:       physicalThingName,
			Type:       physicalThingType,
			Tags:       physicalAndLogicalThingTags,
			Metadata:   physicalAndLogicalThingMetadata,
			RawData:    physicalAndLogicalThingRawData,
		}

		logicalThing1 := &model_generated.LogicalThing{
			ExternalID: &logicalExternalID,
			Name:       logicalThingName,
			Type:       logicalThingType,
			Tags:       physicalAndLogicalThingTags,
			Metadata:   physicalAndLogicalThingMetadata,
			RawData:    physicalAndLogicalThingRawData,
		}

		logicalThing2 := &model_generated.LogicalThing{
			ExternalID:  helpers.Ptr(logicalExternalID + "-2"),
			Name:        logicalThingName + "-2",
			Type:        logicalThingType + "-2",
			Tags:        physicalAndLogicalThingTags,
			Metadata:    physicalAndLogicalThingMetadata,
			RawData:     physicalAndLogicalThingRawData,
			Age:         time.Millisecond * 5555,
			OptionalAge: helpers.Ptr(time.Millisecond * 6666),
		}

		func() {
			tx, _ := db.Begin(ctx)
			defer tx.Rollback(ctx)
			err = physicalThing.Insert(
				ctx,
				tx,
				false,
				false,
			)
			require.NoError(t, err)
			require.NotNil(t, physicalThing)

			logicalThing1.ParentPhysicalThingID = helpers.Ptr(physicalThing.ID)
			err = logicalThing1.Insert(
				ctx,
				tx,
				false,
				false,
			)
			require.NoError(t, err)
			require.NotNil(t, logicalThing1)

			logicalThing2.ParentPhysicalThingID = helpers.Ptr(physicalThing.ID)
			err = logicalThing2.Insert(
				ctx,
				tx,
				false,
				false,
			)
			require.NoError(t, err)
			require.NotNil(t, logicalThing2)

			_ = tx.Commit(ctx)
		}()

		r, err := httpClient.Get("http://127.0.0.1:5050/logical-things")
		require.NoError(t, err)
		b, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, r.StatusCode, string(b))

		var response server.Response[any]
		err = json.Unmarshal(b, &response)
		require.NoError(t, err)

		require.True(t, response.Success)
		require.Empty(t, response.Error)
		require.NotEmpty(t, response.Objects)

		objects := response.Objects

		tempObjects := make([]*any, 0)
		for _, object := range objects {
			object, ok := (*object).(map[string]any)
			require.True(t, ok)

			if object["name"] != logicalThing1.Name && object["name"] != logicalThing2.Name {
				continue
			}

			var x any = object
			tempObjects = append(tempObjects, &x)
		}
		objects = tempObjects

		require.Equal(t, len(objects), 2)
	})

	t.Run("RouterGetManyWithFilters", func(t *testing.T) {
		physicalExternalID := "RouterGetManySomePhysicalThingExternalID"
		physicalThingName := "RouterGetManySomePhysicalThingName"
		physicalThingType := "RouterGetManySomePhysicalThingType"
		logicalExternalID := "RouterGetManySomeLogicalThingExternalID"
		logicalThingName := "RouterGetManySomeLogicalThingName"
		logicalThingType := "RouterGetManySomeLogicalThingType"
		physicalAndLogicalThingTags := []string{"tag1", "tag2", "tag3", "isn''t this, \"complicated\""}
		physicalAndLogicalThingMetadata := map[string]*string{
			"key1": helpers.Ptr("1"),
			"key2": helpers.Ptr("a"),
			"key3": helpers.Ptr("true"),
			"key4": nil,
			"key5": helpers.Ptr("isn''t this, \"complicated\""),
		}
		physicalAndLogicalThingRawData := map[string]any{
			"key1": "1",
			"key2": "a",
			"key3": true,
			"key4": nil,
			"key5": "isn''t this, \"complicated\"",
		}

		cleanup := func() {
			_, _ = db.Exec(
				ctx,
				`DELETE FROM logical_things
			WHERE
				name = $1;`,
				logicalThingName,
			)
			_, _ = db.Exec(
				ctx,
				`DELETE FROM logical_things
			WHERE
				name = $1;`,
				logicalThingName+"-2",
			)
			_, _ = db.Exec(
				ctx,
				`DELETE FROM physical_things
			WHERE
				name = $1;`,
				physicalThingName,
			)
			if redisConn != nil {
				_, _ = redisConn.Do("FLUSHALL")
			}

		}
		defer cleanup()

		var _ error

		physicalThing := &model_generated.PhysicalThing{
			ExternalID: &physicalExternalID,
			Name:       physicalThingName,
			Type:       physicalThingType,
			Tags:       physicalAndLogicalThingTags,
			Metadata:   physicalAndLogicalThingMetadata,
			RawData:    physicalAndLogicalThingRawData,
		}

		logicalThing1 := &model_generated.LogicalThing{
			ExternalID: &logicalExternalID,
			Name:       logicalThingName,
			Type:       logicalThingType,
			Tags:       physicalAndLogicalThingTags,
			Metadata:   physicalAndLogicalThingMetadata,
			RawData:    physicalAndLogicalThingRawData,
		}

		logicalThing2 := &model_generated.LogicalThing{
			ExternalID: helpers.Ptr(logicalExternalID + "-2"),
			Name:       logicalThingName + "-2",
			Type:       logicalThingType + "-2",
			Tags:       physicalAndLogicalThingTags,
			Metadata:   physicalAndLogicalThingMetadata,
			RawData:    physicalAndLogicalThingRawData,
		}

		func() {
			func() {
				tx, _ := db.Begin(ctx)
				defer tx.Rollback(ctx)
				err := physicalThing.Insert(
					ctx,
					tx,
					false,
					false,
				)
				require.NoError(t, err)
				require.NotNil(t, physicalThing)
				_ = tx.Commit(ctx)
			}()

			func() {
				tx, _ := db.Begin(ctx)
				defer tx.Rollback(ctx)
				logicalThing1.ParentPhysicalThingID = helpers.Ptr(physicalThing.ID)
				err := logicalThing1.Insert(
					ctx,
					tx,
					false,
					false,
				)
				require.NoError(t, err)
				require.NotNil(t, logicalThing1)
				_ = tx.Commit(ctx)
			}()

			func() {
				tx, _ := db.Begin(ctx)
				defer tx.Rollback(ctx)
				logicalThing2.ParentPhysicalThingID = helpers.Ptr(physicalThing.ID)
				err := logicalThing2.Insert(
					ctx,
					tx,
					false,
					false,
				)
				require.NoError(t, err)
				require.NotNil(t, logicalThing2)
				_ = tx.Commit(ctx)
			}()
		}()

		doReq := func(params string) []*any {
			r, err := httpClient.Get(
				fmt.Sprintf(
					"http://127.0.0.1:5050/logical-things?%s",
					params,
				),
			)
			require.NoError(t, err)
			b, err := io.ReadAll(r.Body)
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, r.StatusCode, string(b))

			var response server.Response[any]
			err = json.Unmarshal(b, &response)
			require.NoError(t, err)

			require.Equal(t, http.StatusOK, response.Status)
			require.True(t, response.Success)
			require.Empty(t, response.Error)
			require.NotEmpty(t, response.Objects, fmt.Sprintf("%v | %v", params, hack.UnsafeJSONPrettyFormat(response.Objects)))

			objects := response.Objects

			return objects
		}

		allParams := []string{
			fmt.Sprintf("id__eq=%s", logicalThing1.ID.String()),
			fmt.Sprintf("id__ne=%s", logicalThing2.ID.String()),
			fmt.Sprintf("id__in=%s,%s,%s", logicalThing1.ID.String(), uuid.Must(uuid.NewRandom()).String(), uuid.Must(uuid.NewRandom()).String()),
			fmt.Sprintf("id__notin=%s,%s,%s", logicalThing2.ID.String(), uuid.Must(uuid.NewRandom()).String(), uuid.Must(uuid.NewRandom()).String()),
			fmt.Sprintf("id__eq=%s&external_id__eq=%s", logicalThing1.ID.String(), *logicalThing1.ExternalID),
			fmt.Sprintf("created_at__lt=%s", logicalThing2.CreatedAt.Format(time.RFC3339Nano)),
			fmt.Sprintf("created_at__lte=%s", logicalThing2.CreatedAt.Format(time.RFC3339Nano)),
			fmt.Sprintf("id__eq=%s&created_at__gt=%s", logicalThing1.ID.String(), "2020-03-27T08:30:00%%2b08:00"),
			fmt.Sprintf("id__eq=%s&created_at__gte=%s", logicalThing1.ID.String(), "2020-03-27T08:30:00%%2b08:00"),
			fmt.Sprintf("id__eq=%s&parent_logical_thing_id__isnull=", logicalThing1.ID.String()),
			fmt.Sprintf("id__eq=%s&parent_logical_thing_id__isnull&parent_physical_thing_id__isnotnull=", logicalThing1.ID.String()),
			fmt.Sprintf("name__ilike=%s&name__notilike=%s", logicalThing1.Name[4:], logicalThing2.Name[4:]),
		}

		allCounts := []int{
			1,
			1,
			1,
			2,
			1,
			1,
			2,
			1,
			1,
			1,
			1,
			1,
		}

		for i, params := range allParams {
			objects := doReq(params)

			require.Equal(t, allCounts[i], len(objects), fmt.Sprintf("%v | %v", params, hack.UnsafeJSONPrettyFormat(objects)))

			// object1, ok := objects[0].(map[string]any)
			// require.True(t, ok)

			// require.Equal(
			// 	t,
			// 	map[string]any{
			// 		"age":          float64(time.Duration(0)),
			// 		"optional_age": nil,
			// 		"created_at":   logicalThing1.CreatedAt.Format(time.RFC3339Nano),
			// 		"deleted_at":   nil,
			// 		"external_id":  "RouterGetManySomeLogicalThingExternalID",
			// 		"id":           logicalThing1.ID.String(),
			// 		"metadata": map[string]any{
			// 			"key1": "1",
			// 			"key2": "a",
			// 			"key3": "true",
			// 			"key4": nil,
			// 			"key5": "isn''t this, \"complicated\"",
			// 		},
			// 		"name":                           "RouterGetManySomeLogicalThingName",
			// 		"parent_logical_thing_id":        nil,
			// 		"parent_logical_thing_id_object": nil,
			// 		"parent_physical_thing_id":       logicalThing1.ParentPhysicalThingIDObject.ID.String(),
			// 		"parent_physical_thing_id_object": map[string]any{
			// 			"created_at":  logicalThing1.ParentPhysicalThingIDObject.CreatedAt.Format(time.RFC3339Nano),
			// 			"deleted_at":  nil,
			// 			"external_id": "RouterGetManySomePhysicalThingExternalID",
			// 			"id":          logicalThing1.ParentPhysicalThingIDObject.ID.String(),
			// 			"metadata": map[string]any{
			// 				"key1": "1",
			// 				"key2": "a",
			// 				"key3": "true",
			// 				"key4": nil,
			// 				"key5": "isn''t this, \"complicated\"",
			// 			},
			// 			"name": "RouterGetManySomePhysicalThingName",
			// 			"raw_data": map[string]any{
			// 				"key1": "1",
			// 				"key2": "a",
			// 				"key3": true,
			// 				"key4": nil,
			// 				"key5": "isn''t this, \"complicated\"",
			// 			},
			// 			"tags": []any{
			// 				"tag1",
			// 				"tag2",
			// 				"tag3",
			// 				"isn''t this, \"complicated\"",
			// 			},
			// 			"type":       "RouterGetManySomePhysicalThingType",
			// 			"updated_at": logicalThing1.ParentPhysicalThingIDObject.UpdatedAt.Format(time.RFC3339Nano),
			// 			"referenced_by_location_history_parent_physical_thing_id_objects": []any{},
			// 			"referenced_by_logical_thing_parent_physical_thing_id_objects":    nil,
			// 		},
			// 		"raw_data": map[string]any{
			// 			"key1": "1",
			// 			"key2": "a",
			// 			"key3": true,
			// 			"key4": nil,
			// 			"key5": "isn''t this, \"complicated\"",
			// 		},
			// 		"tags": []any{
			// 			"tag1",
			// 			"tag2",
			// 			"tag3",
			// 			"isn''t this, \"complicated\"",
			// 		},
			// 		"type":       "RouterGetManySomeLogicalThingType",
			// 		"updated_at": logicalThing1.CreatedAt.Format(time.RFC3339Nano),
			// 		"referenced_by_logical_thing_parent_logical_thing_id_objects": nil,
			// 	},
			// 	object1,
			// )
		}
	})

	t.Run("RouterGetOne", func(t *testing.T) {
		physicalExternalID := "RouterGetOneSomePhysicalThingExternalID"
		physicalThingName := "RouterGetOneSomePhysicalThingName"
		physicalThingType := "RouterGetOneSomePhysicalThingType"
		logicalExternalID := "RouterGetOneSomeLogicalThingExternalID"
		logicalThingName := "RouterGetOneSomeLogicalThingName"
		logicalThingType := "RouterGetOneSomeLogicalThingType"
		physicalAndLogicalThingTags := []string{"tag1", "tag2", "tag3", "isn''t this, \"complicated\""}
		physicalAndLogicalThingMetadata := map[string]*string{
			"key1": helpers.Ptr("1"),
			"key2": helpers.Ptr("a"),
			"key3": helpers.Ptr("true"),
			"key4": nil,
			"key5": helpers.Ptr("isn''t this, \"complicated\""),
		}
		physicalAndLogicalThingRawData := map[string]any{
			"key1": "1",
			"key2": "a",
			"key3": true,
			"key4": nil,
			"key5": "isn''t this, \"complicated\"",
		}

		cleanup := func() {
			_, _ = db.Exec(
				ctx,
				`DELETE FROM logical_things
			WHERE
				name = $1;`,
				logicalThingName,
			)
			_, _ = db.Exec(
				ctx,
				`DELETE FROM logical_things
			WHERE
				name = $1;`,
				logicalThingName+"-2",
			)
			_, _ = db.Exec(
				ctx,
				`DELETE FROM physical_things
			WHERE
				name = $1;`,
				physicalThingName,
			)
			if redisConn != nil {
				_, _ = redisConn.Do("FLUSHALL")
			}

		}
		defer cleanup()

		var err error

		physicalThing := &model_generated.PhysicalThing{
			ExternalID: &physicalExternalID,
			Name:       physicalThingName,
			Type:       physicalThingType,
			Tags:       physicalAndLogicalThingTags,
			Metadata:   physicalAndLogicalThingMetadata,
			RawData:    physicalAndLogicalThingRawData,
		}

		logicalThing1 := &model_generated.LogicalThing{
			ExternalID: &logicalExternalID,
			Name:       logicalThingName,
			Type:       logicalThingType,
			Tags:       physicalAndLogicalThingTags,
			Metadata:   physicalAndLogicalThingMetadata,
			RawData:    physicalAndLogicalThingRawData,
		}

		logicalThing2 := &model_generated.LogicalThing{
			ExternalID: helpers.Ptr(logicalExternalID + "-2"),
			Name:       logicalThingName + "-2",
			Type:       logicalThingType + "-2",
			Tags:       physicalAndLogicalThingTags,
			Metadata:   physicalAndLogicalThingMetadata,
			RawData:    physicalAndLogicalThingRawData,
		}

		func() {
			tx, _ := db.Begin(ctx)
			defer tx.Rollback(ctx)
			err = physicalThing.Insert(
				ctx,
				tx,
				false,
				false,
			)
			require.NoError(t, err)
			require.NotNil(t, physicalThing)

			logicalThing1.ParentPhysicalThingID = helpers.Ptr(physicalThing.ID)
			err = logicalThing1.Insert(
				ctx,
				tx,
				false,
				false,
			)
			require.NoError(t, err)
			require.NotNil(t, logicalThing1)

			logicalThing2.ParentPhysicalThingID = helpers.Ptr(physicalThing.ID)
			err = logicalThing2.Insert(
				ctx,
				tx,
				false,
				false,
			)
			require.NoError(t, err)
			require.NotNil(t, logicalThing2)

			_ = tx.Commit(ctx)
		}()

		r, err := httpClient.Get(fmt.Sprintf("http://127.0.0.1:5050/logical-things/%v", logicalThing1.ID.String()))
		require.NoError(t, err)
		b, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, r.StatusCode, []string{string(b), logicalThing1.ID.String()})

		var response server.Response[any]
		err = json.Unmarshal(b, &response)
		require.NoError(t, err)

		require.True(t, response.Success)
		require.Empty(t, response.Error)
		require.NotEmpty(t, response.Objects)

		objects := response.Objects
		require.Equal(t, 1, len(objects))

		// object1, ok := objects[0].(map[string]any)
		// require.True(t, ok)

		// require.Equal(
		// 	t,
		// 	map[string]any{
		// 		"age":          float64(time.Duration(0)),
		// 		"optional_age": nil,
		// 		"created_at":   logicalThing1.CreatedAt.Format(time.RFC3339Nano),
		// 		"deleted_at":   nil,
		// 		"external_id":  "RouterGetOneSomeLogicalThingExternalID",
		// 		"id":           logicalThing1.ID.String(),
		// 		"metadata": map[string]any{
		// 			"key1": "1",
		// 			"key2": "a",
		// 			"key3": "true",
		// 			"key4": nil,
		// 			"key5": "isn''t this, \"complicated\"",
		// 		},
		// 		"name":                           "RouterGetOneSomeLogicalThingName",
		// 		"parent_logical_thing_id":        nil,
		// 		"parent_logical_thing_id_object": nil,
		// 		"parent_physical_thing_id":       logicalThing1.ParentPhysicalThingIDObject.ID.String(),
		// 		"parent_physical_thing_id_object": map[string]any{
		// 			"created_at":  logicalThing1.ParentPhysicalThingIDObject.CreatedAt.Format(time.RFC3339Nano),
		// 			"deleted_at":  nil,
		// 			"external_id": "RouterGetOneSomePhysicalThingExternalID",
		// 			"id":          logicalThing1.ParentPhysicalThingIDObject.ID.String(),
		// 			"metadata": map[string]any{
		// 				"key1": "1",
		// 				"key2": "a",
		// 				"key3": "true",
		// 				"key4": nil,
		// 				"key5": "isn''t this, \"complicated\"",
		// 			},
		// 			"name": "RouterGetOneSomePhysicalThingName",
		// 			"raw_data": map[string]any{
		// 				"key1": "1",
		// 				"key2": "a",
		// 				"key3": true,
		// 				"key4": nil,
		// 				"key5": "isn''t this, \"complicated\"",
		// 			},
		// 			"tags": []any{
		// 				"tag1",
		// 				"tag2",
		// 				"tag3",
		// 				"isn''t this, \"complicated\"",
		// 			},
		// 			"type":       "RouterGetOneSomePhysicalThingType",
		// 			"updated_at": logicalThing1.ParentPhysicalThingIDObject.UpdatedAt.Format(time.RFC3339Nano),
		// 			"referenced_by_location_history_parent_physical_thing_id_objects": []any{},
		// 			"referenced_by_logical_thing_parent_physical_thing_id_objects":    nil,
		// 		},
		// 		"raw_data": map[string]any{
		// 			"key1": "1",
		// 			"key2": "a",
		// 			"key3": true,
		// 			"key4": nil,
		// 			"key5": "isn''t this, \"complicated\"",
		// 		},
		// 		"tags": []any{
		// 			"tag1",
		// 			"tag2",
		// 			"tag3",
		// 			"isn''t this, \"complicated\"",
		// 		},
		// 		"type":       "RouterGetOneSomeLogicalThingType",
		// 		"updated_at": logicalThing1.CreatedAt.Format(time.RFC3339Nano),
		// 		"referenced_by_logical_thing_parent_logical_thing_id_objects": nil,
		// 	},
		// 	object1,
		// )
	})

	t.Run("RouterPostOne", func(t *testing.T) {
		physicalExternalID := "RouterPostOneSomePhysicalThingExternalID"
		physicalThingName := "RouterPostOneSomePhysicalThingName"
		physicalThingType := "RouterPostOneSomePhysicalThingType"
		logicalExternalID := "RouterPostOneSomeLogicalThingExternalID"
		logicalThingName := "RouterPostOneSomeLogicalThingName"
		logicalThingType := "RouterPostOneSomeLogicalThingType"
		physicalAndLogicalThingTags := []string{"tag1", "tag2", "tag3", "isn''t this, \"complicated\""}
		physicalAndLogicalThingMetadata := map[string]*string{
			"key1": helpers.Ptr("1"),
			"key2": helpers.Ptr("a"),
			"key3": helpers.Ptr("true"),
			"key4": nil,
			"key5": helpers.Ptr("isn''t this, \"complicated\""),
		}
		physicalAndLogicalThingRawData := map[string]any{
			"key1": "1",
			"key2": "a",
			"key3": true,
			"key4": nil,
			"key5": "isn''t this, \"complicated\"",
		}

		cleanup := func() {
			_, _ = db.Exec(
				ctx,
				`DELETE FROM logical_things
			WHERE
				name = $1;`,
				logicalThingName,
			)
			_, _ = db.Exec(
				ctx,
				`DELETE FROM logical_things
			WHERE
				name = $1;`,
				logicalThingName+"-2",
			)
			_, _ = db.Exec(
				ctx,
				`DELETE FROM physical_things
			WHERE
				name = $1;`,
				physicalThingName,
			)
			if redisConn != nil {
				_, _ = redisConn.Do("FLUSHALL")
			}

		}
		defer cleanup()

		var err error

		physicalThing := &model_generated.PhysicalThing{
			ExternalID: &physicalExternalID,
			Name:       physicalThingName,
			Type:       physicalThingType,
			Tags:       physicalAndLogicalThingTags,
			Metadata:   physicalAndLogicalThingMetadata,
			RawData:    physicalAndLogicalThingRawData,
		}

		logicalThing1 := &model_generated.LogicalThing{
			ExternalID: &logicalExternalID,
			Name:       logicalThingName,
			Type:       logicalThingType,
			Tags:       physicalAndLogicalThingTags,
			Metadata:   physicalAndLogicalThingMetadata,
			RawData:    physicalAndLogicalThingRawData,
		}

		logicalThing2 := &model_generated.LogicalThing{
			ExternalID: helpers.Ptr(logicalExternalID + "-2"),
			Name:       logicalThingName + "-2",
			Type:       logicalThingType + "-2",
			Tags:       physicalAndLogicalThingTags,
			Metadata:   physicalAndLogicalThingMetadata,
			RawData:    physicalAndLogicalThingRawData,
		}

		func() {
			tx, _ := db.Begin(ctx)
			defer tx.Rollback(ctx)

			ctx = query.WithMaxDepth(ctx, helpers.Ptr(0))

			err = physicalThing.Insert(
				ctx,
				tx,
				false,
				false,
			)
			require.NoError(t, err)
			require.NotNil(t, physicalThing)

			logicalThing1.ParentPhysicalThingID = helpers.Ptr(physicalThing.ID)
			err = logicalThing1.Insert(
				ctx,
				tx,
				false,
				false,
			)
			require.NoError(t, err)
			require.NotNil(t, logicalThing1)

			logicalThing2.ParentPhysicalThingID = helpers.Ptr(physicalThing.ID)

			_ = tx.Commit(ctx)
		}()

		rawItem := map[string]any{
			"id":                       logicalThing2.ID,
			"external_id":              logicalThing2.ExternalID,
			"name":                     logicalThing2.Name,
			"type":                     logicalThing2.Type,
			"tags":                     logicalThing2.Tags,
			"metadata":                 logicalThing2.Metadata,
			"raw_data":                 logicalThing2.RawData,
			"parent_physical_thing_id": logicalThing2.ParentPhysicalThingID,
			"parent_logical_thing_id":  logicalThing2.ParentLogicalThingID,
		}

		payload, err := json.Marshal([]any{rawItem})
		require.NoError(t, err)

		r, err := httpClient.Post("http://127.0.0.1:5050/logical-things?depth=0", "application/json", bytes.NewReader(payload))
		require.NoError(t, err)
		b, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, r.StatusCode, string(b))

		var response server.Response[any]
		err = json.Unmarshal(b, &response)
		require.NoError(t, err)

		require.True(t, response.Success)
		require.Empty(t, response.Error)
		require.NotEmpty(t, response.Objects)

		objects := response.Objects
		require.Equal(t, 1, len(objects))

		object1, ok := (*objects[0]).(map[string]any)
		require.True(t, ok)

		require.Equal(t, "RouterPostOneSomeLogicalThingName-2", object1["name"])
		require.Equal(t, "RouterPostOneSomeLogicalThingType-2", object1["type"])
		require.Equal(t, "RouterPostOneSomeLogicalThingExternalID-2", object1["external_id"])
		require.Equal(t, []interface{}([]interface{}{"tag1", "tag2", "tag3", "isn''t this, \"complicated\""}), object1["tags"])
		require.Equal(t, map[string]interface{}(map[string]interface{}{"key1": "1", "key2": "a", "key3": "true", "key4": interface{}(nil), "key5": "isn''t this, \"complicated\""}), object1["metadata"])
		require.Equal(t, map[string]interface{}(map[string]interface{}{"key1": "1", "key2": "a", "key3": true, "key4": interface{}(nil), "key5": "isn''t this, \"complicated\""}), object1["raw_data"])
		require.NotNil(t, object1["parent_physical_thing_id"])
		require.NotNil(t, object1["parent_physical_thing_id_object"])
		require.Nil(t, object1["parent_logical_thing_id"])
		require.Nil(t, object1["parent_logical_thing_id_object"])
	})

	t.Run("RouterPutOne", func(t *testing.T) {
		physicalExternalID := "RouterPutOneSomePhysicalThingExternalID"
		physicalThingName := "RouterPutOneSomePhysicalThingName"
		physicalThingType := "RouterPutOneSomePhysicalThingType"
		logicalExternalID := "RouterPutOneSomeLogicalThingExternalID"
		logicalThingName := "RouterPutOneSomeLogicalThingName"
		logicalThingType := "RouterPutOneSomeLogicalThingType"
		physicalAndLogicalThingTags := []string{"tag1", "tag2", "tag3", "isn''t this, \"complicated\""}
		physicalAndLogicalThingMetadata := map[string]*string{
			"key1": helpers.Ptr("1"),
			"key2": helpers.Ptr("a"),
			"key3": helpers.Ptr("true"),
			"key4": nil,
			"key5": helpers.Ptr("isn''t this, \"complicated\""),
		}
		physicalAndLogicalThingRawData := map[string]any{
			"key1": "1",
			"key2": "a",
			"key3": true,
			"key4": nil,
			"key5": "isn''t this, \"complicated\"",
		}

		cleanup := func() {
			_, _ = db.Exec(
				ctx,
				`DELETE FROM logical_things
			WHERE
				name = $1;`,
				logicalThingName,
			)
			_, _ = db.Exec(
				ctx,
				`DELETE FROM logical_things
			WHERE
				name = $1;`,
				logicalThingName+"-2",
			)
			_, _ = db.Exec(
				ctx,
				`DELETE FROM physical_things
			WHERE
				name = $1;`,
				physicalThingName,
			)
			if redisConn != nil {
				_, _ = redisConn.Do("FLUSHALL")
			}

		}
		defer cleanup()

		var err error

		physicalThing := &model_generated.PhysicalThing{
			ExternalID: &physicalExternalID,
			Name:       physicalThingName,
			Type:       physicalThingType,
			Tags:       physicalAndLogicalThingTags,
			Metadata:   physicalAndLogicalThingMetadata,
			RawData:    physicalAndLogicalThingRawData,
		}

		logicalThing1 := &model_generated.LogicalThing{
			ExternalID: &logicalExternalID,
			Name:       logicalThingName,
			Type:       logicalThingType,
			Tags:       physicalAndLogicalThingTags,
			Metadata:   physicalAndLogicalThingMetadata,
			RawData:    physicalAndLogicalThingRawData,
		}

		logicalThing2 := &model_generated.LogicalThing{
			Name: logicalThingName + "-2",
			Type: logicalThingType + "-2",
		}

		func() {
			tx, _ := db.Begin(ctx)
			defer tx.Rollback(ctx)
			err = physicalThing.Insert(
				ctx,
				tx,
				false,
				false,
			)
			require.NoError(t, err)
			require.NotNil(t, physicalThing)

			logicalThing1.ParentPhysicalThingID = helpers.Ptr(physicalThing.ID)
			err = logicalThing1.Insert(
				ctx,
				tx,
				false,
				false,
			)
			require.NoError(t, err)
			require.NotNil(t, logicalThing1)

			err = logicalThing2.Insert(
				ctx,
				tx,
				false,
				false,
			)

			_ = tx.Commit(ctx)
		}()

		logicalThing2.ExternalID = helpers.Ptr(logicalExternalID + "-2")
		logicalThing2.Tags = physicalAndLogicalThingTags
		logicalThing2.Metadata = physicalAndLogicalThingMetadata
		logicalThing2.RawData = physicalAndLogicalThingRawData
		logicalThing2.ParentPhysicalThingID = helpers.Ptr(physicalThing.ID)

		rawItem := map[string]any{
			"id":                       logicalThing2.ID,
			"external_id":              logicalThing2.ExternalID,
			"name":                     logicalThing2.Name,
			"type":                     logicalThing2.Type,
			"tags":                     logicalThing2.Tags,
			"metadata":                 logicalThing2.Metadata,
			"raw_data":                 logicalThing2.RawData,
			"parent_physical_thing_id": logicalThing2.ParentPhysicalThingID,
			"parent_logical_thing_id":  logicalThing2.ParentLogicalThingID,
			"age":                      (time.Hour * 24) + time.Second*1337,
			"optional_age":             (time.Hour * 24) + time.Second*1337,
		}

		payload, err := json.Marshal(rawItem)
		require.NoError(t, err)

		r, err := httpClient.Put(fmt.Sprintf("http://127.0.0.1:5050/logical-things/%s?depth=0", logicalThing2.ID), "application/json", bytes.NewReader(payload))
		require.NoError(t, err)
		b, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, r.StatusCode, string(b))

		var response server.Response[any]
		err = json.Unmarshal(b, &response)
		require.NoError(t, err)

		require.True(t, response.Success)
		require.Empty(t, response.Error)
		require.NotEmpty(t, response.Objects)

		objects := response.Objects
		require.Equal(t, 1, len(objects))

		object1, ok := (*objects[0]).(map[string]any)
		require.True(t, ok)

		require.Equal(t, "RouterPutOneSomeLogicalThingName-2", object1["name"])
		require.Equal(t, "RouterPutOneSomeLogicalThingType-2", object1["type"])
		require.Equal(t, "RouterPutOneSomeLogicalThingExternalID-2", object1["external_id"])
		require.Equal(t, []interface{}([]interface{}{"tag1", "tag2", "tag3", "isn''t this, \"complicated\""}), object1["tags"])
		require.Equal(t, map[string]interface{}(map[string]interface{}{"key1": "1", "key2": "a", "key3": "true", "key4": interface{}(nil), "key5": "isn''t this, \"complicated\""}), object1["metadata"])
		require.Equal(t, map[string]interface{}(map[string]interface{}{"key1": "1", "key2": "a", "key3": true, "key4": interface{}(nil), "key5": "isn''t this, \"complicated\""}), object1["raw_data"])
		require.Equal(t, float64((time.Hour*24)+time.Second*1337), object1["age"].(float64))
		require.NotNil(t, object1["parent_physical_thing_id"])
		require.NotNil(t, object1["parent_physical_thing_id_object"])
		require.Nil(t, object1["parent_logical_thing_id"])
		require.Nil(t, object1["parent_logical_thing_id_object"])
	})

	t.Run("RouterPatchOne", func(t *testing.T) {
		physicalExternalID := "RouterPatchOneSomePhysicalThingExternalID"
		physicalThingName := "RouterPatchOneSomePhysicalThingName"
		physicalThingType := "RouterPatchOneSomePhysicalThingType"
		logicalExternalID := "RouterPatchOneSomeLogicalThingExternalID"
		logicalThingName := "RouterPatchOneSomeLogicalThingName"
		logicalThingType := "RouterPatchOneSomeLogicalThingType"
		physicalAndLogicalThingTags := []string{"tag1", "tag2", "tag3", "isn''t this, \"complicated\""}
		physicalAndLogicalThingMetadata := map[string]*string{
			"key1": helpers.Ptr("1"),
			"key2": helpers.Ptr("a"),
			"key3": helpers.Ptr("true"),
			"key4": nil,
			"key5": helpers.Ptr("isn''t this, \"complicated\""),
		}
		physicalAndLogicalThingRawData := map[string]any{
			"key1": "1",
			"key2": "a",
			"key3": true,
			"key4": nil,
			"key5": "isn''t this, \"complicated\"",
		}

		cleanup := func() {
			_, _ = db.Exec(
				ctx,
				`DELETE FROM logical_things
			WHERE
				name = $1;`,
				logicalThingName,
			)
			_, _ = db.Exec(
				ctx,
				`DELETE FROM logical_things
			WHERE
				name = $1;`,
				logicalThingName+"-2",
			)
			_, _ = db.Exec(
				ctx,
				`DELETE FROM physical_things
			WHERE
				name = $1;`,
				physicalThingName,
			)
			if redisConn != nil {
				_, _ = redisConn.Do("FLUSHALL")
			}

		}
		defer cleanup()

		var err error

		physicalThing := &model_generated.PhysicalThing{
			ExternalID: &physicalExternalID,
			Name:       physicalThingName,
			Type:       physicalThingType,
			Tags:       physicalAndLogicalThingTags,
			Metadata:   physicalAndLogicalThingMetadata,
			RawData:    physicalAndLogicalThingRawData,
		}

		logicalThing1 := &model_generated.LogicalThing{
			ExternalID: &logicalExternalID,
			Name:       logicalThingName,
			Type:       logicalThingType,
			Tags:       physicalAndLogicalThingTags,
			Metadata:   physicalAndLogicalThingMetadata,
			RawData:    physicalAndLogicalThingRawData,
		}

		logicalThing2 := &model_generated.LogicalThing{
			Name: logicalThingName + "-2",
			Type: logicalThingType + "-2",
		}

		func() {
			tx, _ := db.Begin(ctx)
			defer tx.Rollback(ctx)
			err = physicalThing.Insert(
				ctx,
				tx,
				false,
				false,
			)
			require.NoError(t, err)
			require.NotNil(t, physicalThing)

			logicalThing1.ParentPhysicalThingID = helpers.Ptr(physicalThing.ID)
			err = logicalThing1.Insert(
				ctx,
				tx,
				false,
				false,
			)
			require.NoError(t, err)
			require.NotNil(t, logicalThing1)

			err = logicalThing2.Insert(
				ctx,
				tx,
				false,
				false,
			)

			_ = tx.Commit(ctx)
		}()

		logicalThing2.ExternalID = helpers.Ptr(logicalExternalID + "-2")
		logicalThing2.Tags = physicalAndLogicalThingTags
		logicalThing2.Metadata = physicalAndLogicalThingMetadata
		logicalThing2.RawData = physicalAndLogicalThingRawData
		logicalThing2.ParentPhysicalThingID = helpers.Ptr(physicalThing.ID)

		rawItem := map[string]any{
			"external_id":              logicalThing2.ExternalID,
			"tags":                     logicalThing2.Tags,
			"metadata":                 logicalThing2.Metadata,
			"raw_data":                 logicalThing2.RawData,
			"parent_physical_thing_id": logicalThing2.ParentPhysicalThingID,
		}

		payload, err := json.Marshal(rawItem)
		require.NoError(t, err)

		r, err := httpClient.Patch(fmt.Sprintf("http://127.0.0.1:5050/logical-things/%s?depth=0", logicalThing2.ID), "application/json", bytes.NewReader(payload))
		require.NoError(t, err)
		b, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, r.StatusCode, string(b))

		var response server.Response[any]
		err = json.Unmarshal(b, &response)
		require.NoError(t, err)

		require.True(t, response.Success)
		require.Empty(t, response.Error)
		require.NotEmpty(t, response.Objects)

		objects := response.Objects
		require.Equal(t, 1, len(objects))

		object1, ok := (*objects[0]).(map[string]any)
		require.True(t, ok)

		require.Equal(t, "RouterPatchOneSomeLogicalThingName-2", object1["name"])
		require.Equal(t, "RouterPatchOneSomeLogicalThingType-2", object1["type"])
		require.Equal(t, "RouterPatchOneSomeLogicalThingExternalID-2", object1["external_id"])
		require.Equal(t, []interface{}([]interface{}{"tag1", "tag2", "tag3", "isn''t this, \"complicated\""}), object1["tags"])
		require.Equal(t, map[string]interface{}(map[string]interface{}{"key1": "1", "key2": "a", "key3": "true", "key4": interface{}(nil), "key5": "isn''t this, \"complicated\""}), object1["metadata"])
		require.Equal(t, map[string]interface{}(map[string]interface{}{"key1": "1", "key2": "a", "key3": true, "key4": interface{}(nil), "key5": "isn''t this, \"complicated\""}), object1["raw_data"])
		require.NotNil(t, object1["parent_physical_thing_id"])
		require.NotNil(t, object1["parent_physical_thing_id_object"])
		require.Nil(t, object1["parent_logical_thing_id"])
		require.Nil(t, object1["parent_logical_thing_id_object"])
	})

	t.Run("RouterPatchOneForUndelete", func(t *testing.T) {
		physicalExternalID := "RouterPatchOneForUndeleteSomePhysicalThingExternalID"
		physicalThingName := "RouterPatchOneForUndeleteSomePhysicalThingName"
		physicalThingType := "RouterPatchOneForUndeleteSomePhysicalThingType"
		logicalExternalID := "RouterPatchOneForUndeleteSomeLogicalThingExternalID"
		logicalThingName := "RouterPatchOneForUndeleteSomeLogicalThingName"
		logicalThingType := "RouterPatchOneForUndeleteSomeLogicalThingType"
		physicalAndLogicalThingTags := []string{"tag1", "tag2", "tag3", "isn''t this, \"complicated\""}
		physicalAndLogicalThingMetadata := map[string]*string{
			"key1": helpers.Ptr("1"),
			"key2": helpers.Ptr("a"),
			"key3": helpers.Ptr("true"),
			"key4": nil,
			"key5": helpers.Ptr("isn''t this, \"complicated\""),
		}
		physicalAndLogicalThingRawData := map[string]any{
			"key1": "1",
			"key2": "a",
			"key3": true,
			"key4": nil,
			"key5": "isn''t this, \"complicated\"",
		}

		cleanup := func() {
			_, _ = db.Exec(
				ctx,
				`DELETE FROM logical_things
			WHERE
				name = $1;`,
				logicalThingName,
			)
			_, _ = db.Exec(
				ctx,
				`DELETE FROM logical_things
			WHERE
				name = $1;`,
				logicalThingName+"-2",
			)
			_, _ = db.Exec(
				ctx,
				`DELETE FROM physical_things
			WHERE
				name = $1;`,
				physicalThingName,
			)
			if redisConn != nil {
				_, _ = redisConn.Do("FLUSHALL")
			}

		}
		defer cleanup()

		var err error

		physicalThing := &model_generated.PhysicalThing{
			ExternalID: &physicalExternalID,
			Name:       physicalThingName,
			Type:       physicalThingType,
			Tags:       physicalAndLogicalThingTags,
			Metadata:   physicalAndLogicalThingMetadata,
			RawData:    physicalAndLogicalThingRawData,
		}

		logicalThing1 := &model_generated.LogicalThing{
			ExternalID: &logicalExternalID,
			Name:       logicalThingName,
			Type:       logicalThingType,
			Tags:       physicalAndLogicalThingTags,
			Metadata:   physicalAndLogicalThingMetadata,
			RawData:    physicalAndLogicalThingRawData,
		}

		logicalThing2 := &model_generated.LogicalThing{
			Name: logicalThingName + "-2",
			Type: logicalThingType + "-2",
		}

		func() {
			tx, _ := db.Begin(ctx)
			defer tx.Rollback(ctx)
			err = physicalThing.Insert(
				ctx,
				tx,
				false,
				false,
			)
			require.NoError(t, err)
			require.NotNil(t, physicalThing)

			logicalThing1.ParentPhysicalThingID = helpers.Ptr(physicalThing.ID)
			err = logicalThing1.Insert(
				ctx,
				tx,
				false,
				false,
			)
			require.NoError(t, err)
			require.NotNil(t, logicalThing1)

			err = logicalThing2.Insert(
				ctx,
				tx,
				false,
				false,
			)

			_ = tx.Commit(ctx)
		}()

		r, err := httpClient.Delete(fmt.Sprintf("http://127.0.0.1:5050/logical-things/%v?depth=0", logicalThing2.ID.String()))
		require.NoError(t, err)
		b, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		require.Equal(t, http.StatusNoContent, r.StatusCode, string(b))

		var lastChange server.Change
		require.Eventually(t, func() bool {
			mu.Lock()
			defer mu.Unlock()

			var ok bool
			lastChange, ok = lastChangeByTableName[model_generated.LogicalThingTable]
			if !ok {
				return false
			}

			if lastChange.Action != stream.SOFT_DELETE {
				return false
			}

			return true
		}, time.Second*10, time.Millisecond*10)

		rawItem := map[string]any{
			"deleted_at": nil,
		}

		payload, err := json.Marshal(rawItem)
		require.NoError(t, err)

		r, err = httpClient.Patch(fmt.Sprintf("http://127.0.0.1:5050/logical-things/%s?depth=0", logicalThing2.ID), "application/json", bytes.NewReader(payload))
		require.NoError(t, err)
		b, err = io.ReadAll(r.Body)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, r.StatusCode, string(b))

		require.Eventually(t, func() bool {
			mu.Lock()
			defer mu.Unlock()

			var ok bool
			lastChange, ok = lastChangeByTableName[model_generated.LogicalThingTable]
			if !ok {
				return false
			}

			if lastChange.Action != stream.SOFT_RESTORE {
				return false
			}

			return true
		}, time.Second*10, time.Millisecond*10)

		var response server.Response[any]
		err = json.Unmarshal(b, &response)
		require.NoError(t, err)

		require.True(t, response.Success)
		require.Empty(t, response.Error)
		require.NotEmpty(t, response.Objects)

		objects := response.Objects
		require.Equal(t, 1, len(objects))

		object1, ok := (*objects[0]).(map[string]any)
		require.True(t, ok)

		require.Equal(t, "RouterPatchOneForUndeleteSomeLogicalThingName-2", object1["name"])
		require.Equal(t, "RouterPatchOneForUndeleteSomeLogicalThingType-2", object1["type"])
		require.Equal(t, nil, object1["deleted_at"])

	})

	t.Run("RouterDeleteOne", func(t *testing.T) {
		physicalExternalID := "RouterDeleteOneSomePhysicalThingExternalID"
		physicalThingName := "RouterDeleteOneSomePhysicalThingName"
		physicalThingType := "RouterDeleteOneSomePhysicalThingType"
		logicalExternalID := "RouterDeleteOneSomeLogicalThingExternalID"
		logicalThingName := "RouterDeleteOneSomeLogicalThingName"
		logicalThingType := "RouterDeleteOneSomeLogicalThingType"
		physicalAndLogicalThingTags := []string{"tag1", "tag2", "tag3", "isn''t this, \"complicated\""}
		physicalAndLogicalThingMetadata := map[string]*string{
			"key1": helpers.Ptr("1"),
			"key2": helpers.Ptr("a"),
			"key3": helpers.Ptr("true"),
			"key4": nil,
			"key5": helpers.Ptr("isn''t this, \"complicated\""),
		}
		physicalAndLogicalThingRawData := map[string]any{
			"key1": "1",
			"key2": "a",
			"key3": true,
			"key4": nil,
			"key5": "isn''t this, \"complicated\"",
		}

		cleanup := func() {
			_, _ = db.Exec(
				ctx,
				`DELETE FROM logical_things
			WHERE
				name = $1;`,
				logicalThingName,
			)
			_, _ = db.Exec(
				ctx,
				`DELETE FROM logical_things
			WHERE
				name = $1;`,
				logicalThingName+"-2",
			)
			_, _ = db.Exec(
				ctx,
				`DELETE FROM physical_things
			WHERE
				name = $1;`,
				physicalThingName,
			)
			if redisConn != nil {
				_, _ = redisConn.Do("FLUSHALL")
			}

		}
		defer cleanup()

		var err error

		physicalThing := &model_generated.PhysicalThing{
			ExternalID: &physicalExternalID,
			Name:       physicalThingName,
			Type:       physicalThingType,
			Tags:       physicalAndLogicalThingTags,
			Metadata:   physicalAndLogicalThingMetadata,
			RawData:    physicalAndLogicalThingRawData,
		}

		logicalThing1 := &model_generated.LogicalThing{
			ExternalID: &logicalExternalID,
			Name:       logicalThingName,
			Type:       logicalThingType,
			Tags:       physicalAndLogicalThingTags,
			Metadata:   physicalAndLogicalThingMetadata,
			RawData:    physicalAndLogicalThingRawData,
		}

		logicalThing2 := &model_generated.LogicalThing{
			ExternalID: helpers.Ptr(logicalExternalID + "-2"),
			Name:       logicalThingName + "-2",
			Type:       logicalThingType + "-2",
			Tags:       physicalAndLogicalThingTags,
			Metadata:   physicalAndLogicalThingMetadata,
			RawData:    physicalAndLogicalThingRawData,
		}

		func() {
			tx, _ := db.Begin(ctx)
			defer tx.Rollback(ctx)
			err = physicalThing.Insert(
				ctx,
				tx,
				false,
				false,
			)
			require.NoError(t, err)
			require.NotNil(t, physicalThing)

			logicalThing1.ParentPhysicalThingID = helpers.Ptr(physicalThing.ID)
			err = logicalThing1.Insert(
				ctx,
				tx,
				false,
				false,
			)
			require.NoError(t, err)
			require.NotNil(t, logicalThing1)

			logicalThing2.ParentPhysicalThingID = helpers.Ptr(physicalThing.ID)
			err = logicalThing2.Insert(
				ctx,
				tx,
				false,
				false,
			)
			require.NoError(t, err)
			require.NotNil(t, logicalThing2)

			_ = tx.Commit(ctx)
		}()

		r, err := httpClient.Get(fmt.Sprintf("http://127.0.0.1:5050/logical-things/%v", logicalThing1.ID.String()))
		require.NoError(t, err)
		b, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, r.StatusCode, string(b))

		r, err = httpClient.Delete(fmt.Sprintf("http://127.0.0.1:5050/logical-things/%v", logicalThing1.ID.String()))
		require.NoError(t, err)
		b, err = io.ReadAll(r.Body)
		require.NoError(t, err)
		require.Equal(t, http.StatusNoContent, r.StatusCode, string(b))

		var lastChange server.Change
		require.Eventually(t, func() bool {
			mu.Lock()
			defer mu.Unlock()

			var ok bool
			lastChange, ok = lastChangeByTableName[model_generated.LogicalThingTable]
			if !ok {
				return false
			}

			if lastChange.Action != stream.SOFT_DELETE {
				return false
			}

			return true
		}, time.Second*10, time.Millisecond*10)

		logicalThingFromLastChange := &model_generated.LogicalThing{}
		err = logicalThingFromLastChange.FromItem(lastChange.Item)
		require.NoError(t, err)

		log.Printf("logicalThingFromLastChange: %v", hack.UnsafeJSONPrettyFormat(logicalThingFromLastChange))

		require.Equal(t, logicalThing1.ID, logicalThingFromLastChange.ID)
		require.Nil(t, logicalThing1.DeletedAt)
		require.NotNil(t, logicalThingFromLastChange.DeletedAt)
	})
}
