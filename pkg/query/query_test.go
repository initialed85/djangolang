package query

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"testing"
	"time"

	_helpers "github.com/initialed85/djangolang/internal/helpers"
	"github.com/initialed85/djangolang/pkg/helpers"
	"github.com/jackc/pgtype"
	"github.com/stretchr/testify/require"
)

func TestQuery(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := helpers.GetDBFromEnvironment(ctx)
	if err != nil {
		require.NoError(t, err)
	}
	defer func() {
		db.Close()
	}()

	t.Run("Select", func(t *testing.T) {
		tx, err := db.Begin(ctx)
		require.NoError(t, err)
		defer func() {
			_ = tx.Rollback(ctx)
		}()

		physicalExternalID := "QuerySelectSomePhysicalThingExternalID"
		physicalThingName := "QuerySelectSomePhysicalThingName"
		physicalThingType := "QuerySelectSomePhysicalThingType"
		physicalThingTags := `'{tag1,tag2,tag3,"isn''t this, \"complicated\""}'`
		physicalThingMetadata := `'key1=>1, key2=>"a", key3=>true, key4=>NULL, key5=>"isn''t this, \"complicated\""'`
		physicalThingRawData := `'{"key1": 1, "key2": "a", "key3": true, "key4": null, "key5": "isn''t this, \"complicated\""}'`

		cleanup := func() {
			_, err = db.Exec(
				ctx,
				`DELETE FROM physical_things WHERE name = $1;`,
				physicalThingName,
			)
			require.NoError(t, err)
		}
		defer cleanup()

		_, err = db.Exec(
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

		items, err := Select(
			ctx,
			tx,
			[]string{
				"id",
				"created_at",
				"updated_at",
				"deleted_at",
				"external_id",
				"name",
				"type",
				"tags",
				"metadata",
				"raw_data",
			},
			"physical_things",
			"external_id = $$??\n    AND name = $$??\n    AND type = $$??\n    AND deleted_at IS null",
			helpers.Ptr("created_at ASC"),
			nil,
			nil,
			physicalExternalID,
			physicalThingName,
			physicalThingType,
		)
		require.NoError(t, err)
		require.Len(t, *items, 1)

		item := (*items)[0]

		log.Printf("item: %v", _helpers.UnsafeJSONPrettyFormat(item))

		require.IsType(t, [16]uint8{}, item["id"], "id")
		require.IsType(t, time.Time{}, item["created_at"], "created_at")
		require.IsType(t, time.Time{}, item["updated_at"], "updated_at")
		require.IsType(t, nil, item["deleted_at"], "deleted_at")
		require.IsType(t, "", item["external_id"], "external_id")
		require.IsType(t, "", item["name"], "name")
		require.IsType(t, "", item["type"], "type")
		require.IsType(t, []any{}, item["tags"], "tags")
		require.IsType(t, "", item["metadata"], "metadata")
		require.IsType(t, map[string]any{}, item["raw_data"], "raw_data")

		err = tx.Commit(ctx)
		require.NoError(t, err)
	})

	t.Run("Insert", func(t *testing.T) {
		tx, err := db.Begin(ctx)
		require.NoError(t, err)
		defer func() {
			_ = tx.Rollback(ctx)
		}()

		physicalExternalID := "QueryInsertSomePhysicalThingExternalID"
		physicalThingName := "QueryInsertSomePhysicalThingName"
		physicalThingType := "QueryInsertSomePhysicalThingType"
		physicalThingTags := []string{
			"tag1",
			"tag2",
			"tag3",
			"isn't this, \"complicated\"",
		}

		physicalThingMetadata := pgtype.Hstore{
			Map: map[string]pgtype.Text{
				"key1": {String: "1", Status: pgtype.Present},
				"key2": {String: "a", Status: pgtype.Present},
				"key3": {String: "true", Status: pgtype.Present},
				"key4": {String: "", Status: pgtype.Null},
				"key5": {String: "isn't this, \"complicated\"", Status: pgtype.Present},
			},
			Status: pgtype.Present,
		}

		rawPhysicalThingRawData := map[string]any{
			"key1": 1,
			"key2": "a",
			"key3": true,
			"key4": nil,
			"key5": "isn't this, \"complicated\"",
		}
		physicalThingRawData, err := json.Marshal(rawPhysicalThingRawData)
		require.NoError(t, err)

		cleanup := func() {
			_, err = db.Exec(
				ctx,
				`DELETE FROM physical_things WHERE name = $1;`,
				physicalThingName,
			)
			require.NoError(t, err)
		}
		defer cleanup()

		item, err := Insert(
			ctx,
			tx,
			"physical_things",
			[]string{
				"external_id",
				"name",
				"type",
				"tags",
				"metadata",
				"raw_data",
			},
			nil,
			false,
			false,
			[]string{
				"id",
				"created_at",
				"updated_at",
				"deleted_at",
				"external_id",
				"name",
				"type",
				"tags",
				"metadata",
				"raw_data",
			},
			physicalExternalID,
			physicalThingName,
			physicalThingType,
			physicalThingTags,
			physicalThingMetadata,
			physicalThingRawData,
		)
		require.NoError(t, err)
		require.NotNil(t, item)

		err = tx.Commit(ctx)
		require.NoError(t, err)
	})

	t.Run("Update", func(t *testing.T) {
		tx, err := db.Begin(ctx)
		require.NoError(t, err)
		defer func() {
			_ = tx.Rollback(ctx)
		}()

		insertPhysicalExternalID := "QueryUpdateSomePhysicalThingExternalID1"
		insertPhysicalThingName := "QueryUpdateSomePhysicalThingName1"
		insertPhysicalThingType := "QueryUpdateSomePhysicalThingType1"
		insertPhysicalThingTags := `'{tag1,tag2,tag3,"isn''t this, \"complicated\""}'`
		insertPhysicalThingMetadata := `'key1=>1, key2=>"a", key3=>true, key4=>NULL, key5=>"isn''t this, \"complicated\""'`
		insertPhysicalThingRawData := `'{"key1": 1, "key2": "a", "key3": true, "key4": null, "key5": "isn''t this, \"complicated\""}'`

		physicalExternalID := "QueryUpdateSomePhysicalThingExternalID2"
		physicalThingName := "QueryUpdateSomePhysicalThingName2"
		physicalThingType := "QueryUpdateSomePhysicalThingType2"
		physicalThingTags := []string{
			"tag1",
			"tag2",
			"tag3",
			"isn't this, \"complicated\"",
		}

		cleanup := func() {
			_, err = db.Exec(
				ctx,
				`DELETE FROM physical_things WHERE name = $1 OR name = $2;`,
				insertPhysicalThingName,
				physicalThingName,
			)
			require.NoError(t, err)
		}
		defer cleanup()

		_, err = db.Exec(
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
				insertPhysicalExternalID,
				insertPhysicalThingName,
				insertPhysicalThingType,
				insertPhysicalThingTags,
				insertPhysicalThingMetadata,
				insertPhysicalThingRawData,
			),
		)
		require.NoError(t, err)

		physicalThingMetadata := pgtype.Hstore{
			Map: map[string]pgtype.Text{
				"key1": {String: "1", Status: pgtype.Present},
				"key2": {String: "a", Status: pgtype.Present},
				"key3": {String: "true", Status: pgtype.Present},
				"key4": {String: "", Status: pgtype.Null},
				"key5": {String: "isn't this, \"complicated\"", Status: pgtype.Present},
			},
			Status: pgtype.Present,
		}

		rawPhysicalThingRawData := map[string]any{
			"key1": 1,
			"key2": "a",
			"key3": true,
			"key4": nil,
			"key5": "isn't this, \"complicated\"",
		}
		physicalThingRawData, err := json.Marshal(rawPhysicalThingRawData)
		require.NoError(t, err)

		item, err := Update(
			ctx,
			tx,
			"physical_things",
			[]string{
				"external_id",
				"name",
				"type",
				"tags",
				"metadata",
				"raw_data",
			},
			"external_id = $$??\n    AND name = $$??\n    AND type = $$??",
			[]string{
				"id",
				"created_at",
				"updated_at",
				"deleted_at",
				"external_id",
				"name",
				"type",
				"tags",
				"metadata",
				"raw_data",
			},
			physicalExternalID,
			physicalThingName,
			physicalThingType,
			physicalThingTags,
			physicalThingMetadata,
			physicalThingRawData,
			insertPhysicalExternalID,
			insertPhysicalThingName,
			insertPhysicalThingType,
		)
		require.NoError(t, err)
		require.NotNil(t, item)

		err = tx.Commit(ctx)
		require.NoError(t, err)
	})

	t.Run("Delete", func(t *testing.T) {
		tx, err := db.Begin(ctx)
		require.NoError(t, err)
		defer func() {
			_ = tx.Rollback(ctx)
		}()

		physicalExternalID := "QueryDeleteSomePhysicalThingExternalID1"
		physicalThingName := "QueryDeleteSomePhysicalThingName1"
		physicalThingType := "QueryDeleteSomePhysicalThingType1"
		physicalThingTags := `'{tag1,tag2,tag3,"isn''t this, \"complicated\""}'`
		physicalThingMetadata := `'key1=>1, key2=>"a", key3=>true, key4=>NULL, key5=>"isn''t this, \"complicated\""'`
		physicalThingRawData := `'{"key1": 1, "key2": "a", "key3": true, "key4": null, "key5": "isn''t this, \"complicated\""}'`

		cleanup := func() {
			_, err = db.Exec(
				ctx,
				`DELETE FROM physical_things WHERE name = $1;`,
				physicalThingName,
			)
			require.NoError(t, err)
		}
		defer cleanup()

		_, err = db.Exec(
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

		err = Delete(
			ctx,
			tx,
			"physical_things",
			"external_id = $$??\n    AND name = $$??\n    AND type = $$??",
			physicalExternalID,
			physicalThingName,
			physicalThingType,
		)
		require.NoError(t, err)

		err = tx.Commit(ctx)
		require.NoError(t, err)
	})

	t.Run("GetXid", func(t *testing.T) {
		tx, err := db.Begin(ctx)
		require.NoError(t, err)
		defer func() {
			_ = tx.Rollback(ctx)
		}()

		xid, err := GetXid(ctx, tx)
		require.NoError(t, err)
		require.NotNil(t, xid)
		require.Greater(t, xid, uint32(0))

		err = tx.Commit(ctx)
		require.NoError(t, err)
	})

	t.Run("LockTable", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(ctx, time.Second*5)
		defer cancel()

		checkSelect := func(shouldFail bool) {
			otherCtx, otherCancel := context.WithTimeout(ctx, time.Second*1)
			defer otherCancel()

			otherTx, err := db.Begin(otherCtx)
			require.NoError(t, err)

			defer func() {
				_ = otherTx.Rollback(otherCtx)
			}()

			_, err = Select(
				otherCtx,
				otherTx,
				[]string{
					"id",
				},
				"logical_things",
				"",
				nil,
				nil,
				nil,
			)
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

			err = LockTable(
				otherCtx,
				otherTx,
				"logical_things",
				true,
			)
			if shouldFail {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		}

		tx, err := db.Begin(ctx)
		require.NoError(t, err)
		defer func() {
			_ = tx.Rollback(ctx)
		}()

		err = LockTable(ctx, tx, "logical_things", true)
		require.NoError(t, err)

		checkSelect(true)
		checkLock(true)

		err = tx.Commit(ctx)
		require.NoError(t, err)

		checkSelect(false)
		checkLock(false)
	})
}
