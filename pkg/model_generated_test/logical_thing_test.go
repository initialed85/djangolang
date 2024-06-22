package model_generated_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	_helpers "github.com/initialed85/djangolang/internal/helpers"
	"github.com/initialed85/djangolang/pkg/helpers"
	"github.com/initialed85/djangolang/pkg/model_generated"
	"github.com/initialed85/djangolang/pkg/server"
	"github.com/initialed85/djangolang/pkg/stream"
	"github.com/stretchr/testify/require"
)

type HTTPClient struct {
	httpClient *http.Client
}

func (h *HTTPClient) Get(url string) (resp *http.Response, err error) {
	return h.httpClient.Get(url)
}

func (h *HTTPClient) Post(url, contentType string, body io.Reader) (*http.Response, error) {
	return h.httpClient.Post(url, contentType, body)
}

func (h *HTTPClient) Put(url, contentType string, body io.Reader) (*http.Response, error) {
	r, err := http.NewRequest(http.MethodPut, url, body)
	if err != nil {
		return nil, err
	}

	r.Header.Set("Content-Type", "application/json")

	return h.httpClient.Do(r)
}

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

	httpClient := &HTTPClient{
		httpClient: &http.Client{
			Timeout: time.Second * 5,
		},
	}

	changes := make(chan server.Change, 1024)
	mu := new(sync.Mutex)
	lastChangeByTableName := make(map[string]server.Change)

	go func() {
		os.Setenv("DJANGOLANG_NODE_NAME", "model_generated_logical_thing")
		_ = model_generated.RunServer(ctx, changes, "127.0.0.1:5050", db)
	}()
	runtime.Gosched()

	time.Sleep(time.Second * 1)

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

	t.Run("Select", func(t *testing.T) {
		physicalExternalID := "SelectSomePhysicalThingExternalID"
		physicalThingName := "SelectSomePhysicalThingName"
		physicalThingType := "SelectSomePhysicalThingType"
		logicalExternalID := "SelectSomeLogicalThingExternalID"
		logicalThingName := "SelectSomeLogicalThingName"
		logicalThingType := "SelectSomeLogicalThingType"
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
			_, err = db.ExecContext(
				ctx,
				`DELETE FROM physical_things
			WHERE
				name = $1;`,
				physicalThingName,
			)
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

		var logicalThing *model_generated.LogicalThing
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

			logicalThing, err = model_generated.SelectLogicalThing(
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
		}, time.Second*1, time.Millisecond*10)

		logicalThingFromLastChange := &model_generated.LogicalThing{}
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

	t.Run("Insert", func(t *testing.T) {
		physicalExternalID := "InsertSomePhysicalThingExternalID"
		physicalThingName := "InsertSomePhysicalThingName"
		physicalThingType := "InsertSomePhysicalThingType"
		logicalExternalID := "InsertSomeLogicalThingExternalID"
		logicalThingName := "InsertSomeLogicalThingName"
		logicalThingType := "InsertSomeLogicalThingType"
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
			_, err = db.ExecContext(
				ctx,
				`DELETE FROM logical_things
			WHERE
				name = $1;`,
				logicalThingName,
			)
			_, err = db.ExecContext(
				ctx,
				`DELETE FROM physical_things
			WHERE
				name = $1;`,
				physicalThingName,
			)
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
			tx, _ := db.BeginTxx(ctx, nil)
			defer tx.Rollback()
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

			_ = tx.Commit()
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
		}, time.Second*1, time.Millisecond*10)

		logicalThingFromLastChange := &model_generated.LogicalThing{}
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

		cleanup := func() {
			_, err = db.ExecContext(
				ctx,
				`DELETE FROM logical_things
			WHERE
				name = $1 OR name = $2;`,
				insertLogicalThingName,
				updateLogicalThingName,
			)
			_, err = db.ExecContext(
				ctx,
				`DELETE FROM physical_things
			WHERE
				name = $1;`,
				insertPhysicalThingName,
			)
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
			ExternalID: &insertLogicalExternalID,
			Name:       insertLogicalThingName,
			Type:       insertLogicalThingType,
			Tags:       insertPhysicalAndLogicalThingTags,
			Metadata:   insertPhysicalAndLogicalThingMetadata,
			RawData:    insertPhysicalAndLogicalThingRawData,
		}

		func() {
			tx, _ := db.BeginTxx(ctx, nil)
			defer tx.Rollback()
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
			err = logicalThing.Update(ctx, tx, false)
			require.NoError(t, err)
			require.NotNil(t, logicalThing)

			_ = tx.Commit()
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

			if lastChange.Action != stream.UPDATE {
				return false
			}

			return true
		}, time.Second*1, time.Millisecond*10)

		logicalThingFromLastChange := &model_generated.LogicalThing{}
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

	t.Run("Delete", func(t *testing.T) {
		physicalExternalID := "DeleteSomePhysicalThingExternalID"
		physicalThingName := "DeleteSomePhysicalThingName"
		physicalThingType := "DeleteSomePhysicalThingType"
		logicalExternalID := "DeleteSomeLogicalThingExternalID"
		logicalThingName := "DeleteSomeLogicalThingName"
		logicalThingType := "DeleteSomeLogicalThingType"
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
			_, err = db.ExecContext(
				ctx,
				`DELETE FROM logical_things
			WHERE
				name = $1;`,
				logicalThingName,
			)
			_, err = db.ExecContext(
				ctx,
				`DELETE FROM physical_things
			WHERE
				name = $1;`,
				physicalThingName,
			)
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
			tx, _ := db.BeginTxx(ctx, nil)
			defer tx.Rollback()
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

			_ = tx.Commit()
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

			if lastChange.Action != stream.DELETE && lastChange.Action != stream.UPDATE {
				return false
			}

			return true
		}, time.Second*1, time.Millisecond*10)

		logicalThingFromLastChange := &model_generated.LogicalThing{}
		err = logicalThingFromLastChange.FromItem(lastChange.Item)
		require.NoError(t, err)

		log.Printf("logicalThingFromLastChange: %v", _helpers.UnsafeJSONPrettyFormat(logicalThingFromLastChange))

		require.Equal(t, logicalThing.ID, logicalThingFromLastChange.ID)
		// require.Equal(t, logicalThing.DeletedAt, logicalThingFromLastChange.DeletedAt)
		require.Nil(t, logicalThing.DeletedAt)
		require.NotNil(t, logicalThingFromLastChange.DeletedAt)

		func() {
			tx, _ := db.BeginTxx(ctx, nil)
			defer tx.Rollback()
			err = logicalThingFromLastChange.Reload(ctx, tx)
			require.Error(t, err)
			_ = tx.Commit()
		}()

		log.Printf("logicalThingFromLastChangeAfterReloading: %v", _helpers.UnsafeJSONPrettyFormat(logicalThingFromLastChange))
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
			_, err = db.ExecContext(
				ctx,
				`DELETE FROM logical_things
			WHERE
				name = $1;`,
				logicalThingName,
			)
			_, err = db.ExecContext(
				ctx,
				`DELETE FROM logical_things
			WHERE
				name = $1;`,
				logicalThingName+"-2",
			)
			_, err = db.ExecContext(
				ctx,
				`DELETE FROM physical_things
			WHERE
				name = $1;`,
				physicalThingName,
			)
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
			tx, _ := db.BeginTxx(ctx, nil)
			defer tx.Rollback()
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

			_ = tx.Commit()
		}()

		r, err := httpClient.Get("http://127.0.0.1:5050/logical-things")
		require.NoError(t, err)
		b, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, r.StatusCode, string(b))

		var response helpers.Response
		err = json.Unmarshal(b, &response)
		require.NoError(t, err)

		require.True(t, response.Success)
		require.Empty(t, response.Error)
		require.NotEmpty(t, response.Objects)

		objects, ok := response.Objects.([]any)
		require.True(t, ok)

		tempObjects := make([]any, 0)
		for _, object := range objects {
			object, ok := object.(map[string]any)
			require.True(t, ok)

			if object["name"] != logicalThing1.Name && object["name"] != logicalThing2.Name {
				continue
			}

			tempObjects = append(tempObjects, object)
		}
		objects = tempObjects

		require.Equal(t, len(objects), 2)

		object1, ok := objects[0].(map[string]any)
		require.True(t, ok)

		require.Equal(
			t,
			map[string]any{
				"created_at":  logicalThing1.CreatedAt.Format(time.RFC3339Nano),
				"deleted_at":  nil,
				"external_id": "RouterGetManySomeLogicalThingExternalID",
				"id":          logicalThing1.ID.String(),
				"metadata": map[string]any{
					"key1": "1",
					"key2": "a",
					"key3": "true",
					"key4": nil,
					"key5": "isn''t this, \"complicated\"",
				},
				"name":                           "RouterGetManySomeLogicalThingName",
				"parent_logical_thing_id":        nil,
				"parent_logical_thing_id_object": nil,
				"parent_physical_thing_id":       logicalThing1.ParentPhysicalThingIDObject.ID.String(),
				"parent_physical_thing_id_object": map[string]any{
					"created_at":  logicalThing1.ParentPhysicalThingIDObject.CreatedAt.Format(time.RFC3339Nano),
					"deleted_at":  nil,
					"external_id": "RouterGetManySomePhysicalThingExternalID",
					"id":          logicalThing1.ParentPhysicalThingIDObject.ID.String(),
					"metadata": map[string]any{
						"key1": "1",
						"key2": "a",
						"key3": "true",
						"key4": nil,
						"key5": "isn''t this, \"complicated\"",
					},
					"name": "RouterGetManySomePhysicalThingName",
					"raw_data": map[string]any{
						"key1": "1",
						"key2": "a",
						"key3": true,
						"key4": nil,
						"key5": "isn''t this, \"complicated\"",
					},
					"tags": []any{
						"tag1",
						"tag2",
						"tag3",
						"isn''t this, \"complicated\"",
					},
					"type":       "RouterGetManySomePhysicalThingType",
					"updated_at": logicalThing1.ParentPhysicalThingIDObject.UpdatedAt.Format(time.RFC3339Nano),
				},
				"raw_data": map[string]any{
					"key1": "1",
					"key2": "a",
					"key3": true,
					"key4": nil,
					"key5": "isn''t this, \"complicated\"",
				},
				"tags": []any{
					"tag1",
					"tag2",
					"tag3",
					"isn''t this, \"complicated\"",
				},
				"type":       "RouterGetManySomeLogicalThingType",
				"updated_at": logicalThing1.CreatedAt.Format(time.RFC3339Nano),
			},
			object1,
		)

		object2, ok := objects[1].(map[string]any)
		require.True(t, ok)

		require.Equal(
			t,
			map[string]any{
				"created_at":  logicalThing2.CreatedAt.Format(time.RFC3339Nano),
				"deleted_at":  nil,
				"external_id": "RouterGetManySomeLogicalThingExternalID-2",
				"id":          logicalThing2.ID.String(),
				"metadata": map[string]any{
					"key1": "1",
					"key2": "a",
					"key3": "true",
					"key4": nil,
					"key5": "isn''t this, \"complicated\"",
				},
				"name":                           "RouterGetManySomeLogicalThingName-2",
				"parent_logical_thing_id":        nil,
				"parent_logical_thing_id_object": nil,
				"parent_physical_thing_id":       logicalThing2.ParentPhysicalThingIDObject.ID.String(),
				"parent_physical_thing_id_object": map[string]any{
					"created_at":  logicalThing2.ParentPhysicalThingIDObject.CreatedAt.Format(time.RFC3339Nano),
					"deleted_at":  nil,
					"external_id": "RouterGetManySomePhysicalThingExternalID",
					"id":          logicalThing2.ParentPhysicalThingIDObject.ID.String(),
					"metadata": map[string]any{
						"key1": "1",
						"key2": "a",
						"key3": "true",
						"key4": nil,
						"key5": "isn''t this, \"complicated\"",
					},
					"name": "RouterGetManySomePhysicalThingName",
					"raw_data": map[string]any{
						"key1": "1",
						"key2": "a",
						"key3": true,
						"key4": nil,
						"key5": "isn''t this, \"complicated\"",
					},
					"tags": []any{
						"tag1",
						"tag2",
						"tag3",
						"isn''t this, \"complicated\"",
					},
					"type":       "RouterGetManySomePhysicalThingType",
					"updated_at": logicalThing2.ParentPhysicalThingIDObject.UpdatedAt.Format(time.RFC3339Nano),
				},
				"raw_data": map[string]any{
					"key1": "1",
					"key2": "a",
					"key3": true,
					"key4": nil,
					"key5": "isn''t this, \"complicated\"",
				},
				"tags": []any{
					"tag1",
					"tag2",
					"tag3",
					"isn''t this, \"complicated\"",
				},
				"type":       "RouterGetManySomeLogicalThingType-2",
				"updated_at": logicalThing2.CreatedAt.Format(time.RFC3339Nano),
			},
			object2,
		)
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
			_, err = db.ExecContext(
				ctx,
				`DELETE FROM logical_things
			WHERE
				name = $1;`,
				logicalThingName,
			)
			_, err = db.ExecContext(
				ctx,
				`DELETE FROM logical_things
			WHERE
				name = $1;`,
				logicalThingName+"-2",
			)
			_, err = db.ExecContext(
				ctx,
				`DELETE FROM physical_things
			WHERE
				name = $1;`,
				physicalThingName,
			)
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
			func() {
				tx, _ := db.BeginTxx(ctx, nil)
				defer tx.Rollback()
				err = physicalThing.Insert(
					ctx,
					tx,
					false,
					false,
				)
				require.NoError(t, err)
				require.NotNil(t, physicalThing)
				_ = tx.Commit()
			}()

			func() {
				tx, _ := db.BeginTxx(ctx, nil)
				defer tx.Rollback()
				logicalThing1.ParentPhysicalThingID = helpers.Ptr(physicalThing.ID)
				err = logicalThing1.Insert(
					ctx,
					tx,
					false,
					false,
				)
				require.NoError(t, err)
				require.NotNil(t, logicalThing1)
				log.Printf("%v", logicalThing1.CreatedAt.Format(time.RFC3339Nano))
				_ = tx.Commit()
			}()

			time.Sleep(time.Millisecond * 10)

			func() {
				tx, _ := db.BeginTxx(ctx, nil)
				defer tx.Rollback()
				logicalThing2.ParentPhysicalThingID = helpers.Ptr(physicalThing.ID)
				err = logicalThing2.Insert(
					ctx,
					tx,
					false,
					false,
				)
				require.NoError(t, err)
				require.NotNil(t, logicalThing2)
				log.Printf("%v", logicalThing2.CreatedAt.Format(time.RFC3339Nano))
				_ = tx.Commit()
			}()
		}()

		doReq := func(params string) []any {
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

			var response helpers.Response
			err = json.Unmarshal(b, &response)
			require.NoError(t, err)

			require.Equal(t, http.StatusOK, response.Status)
			require.True(t, response.Success)
			require.Empty(t, response.Error)
			require.NotEmpty(t, response.Objects)

			objects, ok := response.Objects.([]any)
			require.True(t, ok)

			return objects
		}

		allParams := []string{
			fmt.Sprintf("id__eq=%s", logicalThing1.ID.String()),
			fmt.Sprintf("id__ne=%s", logicalThing2.ID.String()),
			fmt.Sprintf("id__in=%s,%s,%s", logicalThing1.ID.String(), uuid.Must(uuid.NewRandom()).String(), uuid.Must(uuid.NewRandom()).String()),
			fmt.Sprintf("id__nin=%s,%s,%s", logicalThing2.ID.String(), uuid.Must(uuid.NewRandom()).String(), uuid.Must(uuid.NewRandom()).String()),
			fmt.Sprintf("id__eq=%s&external_id__eq=%s", logicalThing1.ID.String(), *logicalThing1.ExternalID),
			fmt.Sprintf("created_at__lt=%s", logicalThing2.CreatedAt.Format(time.RFC3339Nano)),
			fmt.Sprintf("created_at__lte=%s", logicalThing2.CreatedAt.Add(-time.Millisecond*10).Format(time.RFC3339Nano)),
			fmt.Sprintf("id__eq=%s&created_at__gt=%s", logicalThing1.ID.String(), "2020-03-27T08:30:00%%2b08:00"),
			fmt.Sprintf("id__eq=%s&created_at__gte=%s", logicalThing1.ID.String(), "2020-03-27T08:30:00%%2b08:00"),
			fmt.Sprintf("id__eq=%s&parent_logical_thing_id__isnull=", logicalThing1.ID.String()),
			fmt.Sprintf("id__eq=%s&parent_logical_thing_id__isnull&parent_physical_thing_id__isnotnull=", logicalThing1.ID.String()),
			fmt.Sprintf("name__il=%s&name__nil=%s", logicalThing1.Name[4:], logicalThing2.Name[4:]),
		}

		for _, params := range allParams {
			objects := doReq(params)

			require.Equal(t, 1, len(objects), params)

			object1, ok := objects[0].(map[string]any)
			require.True(t, ok)

			require.Equal(
				t,
				map[string]any{
					"created_at":  logicalThing1.CreatedAt.Format(time.RFC3339Nano),
					"deleted_at":  nil,
					"external_id": "RouterGetManySomeLogicalThingExternalID",
					"id":          logicalThing1.ID.String(),
					"metadata": map[string]any{
						"key1": "1",
						"key2": "a",
						"key3": "true",
						"key4": nil,
						"key5": "isn''t this, \"complicated\"",
					},
					"name":                           "RouterGetManySomeLogicalThingName",
					"parent_logical_thing_id":        nil,
					"parent_logical_thing_id_object": nil,
					"parent_physical_thing_id":       logicalThing1.ParentPhysicalThingIDObject.ID.String(),
					"parent_physical_thing_id_object": map[string]any{
						"created_at":  logicalThing1.ParentPhysicalThingIDObject.CreatedAt.Format(time.RFC3339Nano),
						"deleted_at":  nil,
						"external_id": "RouterGetManySomePhysicalThingExternalID",
						"id":          logicalThing1.ParentPhysicalThingIDObject.ID.String(),
						"metadata": map[string]any{
							"key1": "1",
							"key2": "a",
							"key3": "true",
							"key4": nil,
							"key5": "isn''t this, \"complicated\"",
						},
						"name": "RouterGetManySomePhysicalThingName",
						"raw_data": map[string]any{
							"key1": "1",
							"key2": "a",
							"key3": true,
							"key4": nil,
							"key5": "isn''t this, \"complicated\"",
						},
						"tags": []any{
							"tag1",
							"tag2",
							"tag3",
							"isn''t this, \"complicated\"",
						},
						"type":       "RouterGetManySomePhysicalThingType",
						"updated_at": logicalThing1.ParentPhysicalThingIDObject.UpdatedAt.Format(time.RFC3339Nano),
					},
					"raw_data": map[string]any{
						"key1": "1",
						"key2": "a",
						"key3": true,
						"key4": nil,
						"key5": "isn''t this, \"complicated\"",
					},
					"tags": []any{
						"tag1",
						"tag2",
						"tag3",
						"isn''t this, \"complicated\"",
					},
					"type":       "RouterGetManySomeLogicalThingType",
					"updated_at": logicalThing1.CreatedAt.Format(time.RFC3339Nano),
				},
				object1,
			)
		}
	})

	t.Run("RouterGetOne", func(t *testing.T) {
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
			_, err = db.ExecContext(
				ctx,
				`DELETE FROM logical_things
			WHERE
				name = $1;`,
				logicalThingName,
			)
			_, err = db.ExecContext(
				ctx,
				`DELETE FROM logical_things
			WHERE
				name = $1;`,
				logicalThingName+"-2",
			)
			_, err = db.ExecContext(
				ctx,
				`DELETE FROM physical_things
			WHERE
				name = $1;`,
				physicalThingName,
			)
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
			tx, _ := db.BeginTxx(ctx, nil)
			defer tx.Rollback()
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

			_ = tx.Commit()
		}()

		r, err := httpClient.Get(fmt.Sprintf("http://127.0.0.1:5050/logical-things/%v", logicalThing1.ID.String()))
		require.NoError(t, err)
		b, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, r.StatusCode, string(b))

		var response helpers.Response
		err = json.Unmarshal(b, &response)
		require.NoError(t, err)

		require.True(t, response.Success)
		require.Empty(t, response.Error)
		require.NotEmpty(t, response.Objects)

		objects, ok := response.Objects.([]any)
		require.True(t, ok)
		require.Equal(t, 1, len(objects))

		object1, ok := objects[0].(map[string]any)
		require.True(t, ok)

		require.Equal(
			t,
			map[string]any{
				"created_at":  logicalThing1.CreatedAt.Format(time.RFC3339Nano),
				"deleted_at":  nil,
				"external_id": "RouterGetManySomeLogicalThingExternalID",
				"id":          logicalThing1.ID.String(),
				"metadata": map[string]any{
					"key1": "1",
					"key2": "a",
					"key3": "true",
					"key4": nil,
					"key5": "isn''t this, \"complicated\"",
				},
				"name":                           "RouterGetManySomeLogicalThingName",
				"parent_logical_thing_id":        nil,
				"parent_logical_thing_id_object": nil,
				"parent_physical_thing_id":       logicalThing1.ParentPhysicalThingIDObject.ID.String(),
				"parent_physical_thing_id_object": map[string]any{
					"created_at":  logicalThing1.ParentPhysicalThingIDObject.CreatedAt.Format(time.RFC3339Nano),
					"deleted_at":  nil,
					"external_id": "RouterGetManySomePhysicalThingExternalID",
					"id":          logicalThing1.ParentPhysicalThingIDObject.ID.String(),
					"metadata": map[string]any{
						"key1": "1",
						"key2": "a",
						"key3": "true",
						"key4": nil,
						"key5": "isn''t this, \"complicated\"",
					},
					"name": "RouterGetManySomePhysicalThingName",
					"raw_data": map[string]any{
						"key1": "1",
						"key2": "a",
						"key3": true,
						"key4": nil,
						"key5": "isn''t this, \"complicated\"",
					},
					"tags": []any{
						"tag1",
						"tag2",
						"tag3",
						"isn''t this, \"complicated\"",
					},
					"type":       "RouterGetManySomePhysicalThingType",
					"updated_at": logicalThing1.ParentPhysicalThingIDObject.UpdatedAt.Format(time.RFC3339Nano),
				},
				"raw_data": map[string]any{
					"key1": "1",
					"key2": "a",
					"key3": true,
					"key4": nil,
					"key5": "isn''t this, \"complicated\"",
				},
				"tags": []any{
					"tag1",
					"tag2",
					"tag3",
					"isn''t this, \"complicated\"",
				},
				"type":       "RouterGetManySomeLogicalThingType",
				"updated_at": logicalThing1.CreatedAt.Format(time.RFC3339Nano),
			},
			object1,
		)
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
			_, err = db.ExecContext(
				ctx,
				`DELETE FROM logical_things
			WHERE
				name = $1;`,
				logicalThingName,
			)
			_, err = db.ExecContext(
				ctx,
				`DELETE FROM logical_things
			WHERE
				name = $1;`,
				logicalThingName+"-2",
			)
			_, err = db.ExecContext(
				ctx,
				`DELETE FROM physical_things
			WHERE
				name = $1;`,
				physicalThingName,
			)
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
			tx, _ := db.BeginTxx(ctx, nil)
			defer tx.Rollback()
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

			_ = tx.Commit()
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

		r, err := httpClient.Post("http://127.0.0.1:5050/logical-things", "application/json", bytes.NewReader(payload))
		require.NoError(t, err)
		b, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, r.StatusCode, string(b))

		var response helpers.Response
		err = json.Unmarshal(b, &response)
		require.NoError(t, err)

		require.True(t, response.Success)
		require.Empty(t, response.Error)
		require.NotEmpty(t, response.Objects)

		objects, ok := response.Objects.([]any)
		require.True(t, ok)
		require.Equal(t, 1, len(objects))

		object1, ok := objects[0].(map[string]any)
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
			_, err = db.ExecContext(
				ctx,
				`DELETE FROM logical_things
			WHERE
				name = $1;`,
				logicalThingName,
			)
			_, err = db.ExecContext(
				ctx,
				`DELETE FROM logical_things
			WHERE
				name = $1;`,
				logicalThingName+"-2",
			)
			_, err = db.ExecContext(
				ctx,
				`DELETE FROM physical_things
			WHERE
				name = $1;`,
				physicalThingName,
			)
		}
		// defer cleanup()
		cleanup()

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
			tx, _ := db.BeginTxx(ctx, nil)
			defer tx.Rollback()
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

			_ = tx.Commit()
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
		}

		payload, err := json.Marshal(rawItem)
		require.NoError(t, err)

		r, err := httpClient.Put(fmt.Sprintf("http://127.0.0.1:5050/logical-things/%s", logicalThing2.ID), "application/json", bytes.NewReader(payload))
		require.NoError(t, err)
		b, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, r.StatusCode, string(b))

		var response helpers.Response
		err = json.Unmarshal(b, &response)
		require.NoError(t, err)

		require.True(t, response.Success)
		require.Empty(t, response.Error)
		require.NotEmpty(t, response.Objects)

		objects, ok := response.Objects.([]any)
		require.True(t, ok)
		require.Equal(t, 1, len(objects))

		object1, ok := objects[0].(map[string]any)
		require.True(t, ok)

		require.Equal(t, "RouterPutOneSomeLogicalThingName-2", object1["name"])
		require.Equal(t, "RouterPutOneSomeLogicalThingType-2", object1["type"])
		require.Equal(t, "RouterPutOneSomeLogicalThingExternalID-2", object1["external_id"])
		require.Equal(t, []interface{}([]interface{}{"tag1", "tag2", "tag3", "isn''t this, \"complicated\""}), object1["tags"])
		require.Equal(t, map[string]interface{}(map[string]interface{}{"key1": "1", "key2": "a", "key3": "true", "key4": interface{}(nil), "key5": "isn''t this, \"complicated\""}), object1["metadata"])
		require.Equal(t, map[string]interface{}(map[string]interface{}{"key1": "1", "key2": "a", "key3": true, "key4": interface{}(nil), "key5": "isn''t this, \"complicated\""}), object1["raw_data"])
		require.NotNil(t, object1["parent_physical_thing_id"])
		require.NotNil(t, object1["parent_physical_thing_id_object"])
		require.Nil(t, object1["parent_logical_thing_id"])
		require.Nil(t, object1["parent_logical_thing_id_object"])
	})

}
