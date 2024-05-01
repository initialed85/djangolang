package query

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"testing"
	"time"

	_helpers "github.com/initialed85/djangolang/internal/helpers"
	"github.com/initialed85/djangolang/pkg/helpers"
	"github.com/lib/pq"
	"github.com/lib/pq/hstore"
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
		_ = db.Close()
	}()

	t.Run("Select", func(t *testing.T) {
		tx, err := db.BeginTxx(ctx, nil)
		require.NoError(t, err)
		defer func() {
			_ = tx.Rollback()
		}()

		physicalExternalID := "QuerySelectSomePhysicalThingExternalID"
		physicalThingName := "QuerySelectSomePhysicalThingName"
		physicalThingType := "QuerySelectSomePhysicalThingType"
		physicalThingTags := `'{tag1,tag2,tag3,"isn''t this, \"complicated\""}'`
		physicalThingMetadata := `'key1=>1, key2=>"a", key3=>true, key4=>NULL, key5=>"isn''t this, \"complicated\""'`
		physicalThingRawData := `'{"key1": 1, "key2": "a", "key3": true, "key4": null, "key5": "isn''t this, \"complicated\""}'`

		cleanup := func() {
			_, err = db.ExecContext(
				ctx,
				`DELETE FROM physical_things WHERE name = $1;`,
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
			"external_id = $1 AND name = $2 AND type = $3",
			nil,
			nil,
			physicalExternalID,
			physicalThingName,
			physicalThingType,
		)
		require.NoError(t, err)
		require.Len(t, items, 1)

		item := items[0]

		log.Printf("item: %v", _helpers.UnsafeJSONPrettyFormat(item))

		require.IsType(t, []uint8{}, item["id"], "id")
		require.IsType(t, time.Time{}, item["created_at"], "created_at")
		require.IsType(t, time.Time{}, item["updated_at"], "updated_at")
		require.IsType(t, nil, item["deleted_at"], "deleted_at")
		require.IsType(t, "", item["external_id"], "external_id")
		require.IsType(t, "", item["name"], "name")
		require.IsType(t, "", item["type"], "type")
		require.IsType(t, []uint8{}, item["tags"], "tags")
		require.IsType(t, []uint8{}, item["metadata"], "metadata")
		require.IsType(t, []uint8{}, item["raw_data"], "raw_data")

		err = tx.Commit()
		require.NoError(t, err)
	})

	t.Run("Insert", func(t *testing.T) {
		tx, err := db.BeginTxx(ctx, nil)
		require.NoError(t, err)
		defer func() {
			_ = tx.Rollback()
		}()

		physicalExternalID := "QueryInsertSomePhysicalThingExternalID"
		physicalThingName := "QueryInsertSomePhysicalThingName"
		physicalThingType := "QueryInsertSomePhysicalThingType"
		physicalThingTags := pq.Array([]string{
			"tag1",
			"tag2",
			"tag3",
			"isn't this, \"complicated\"",
		})

		physicalThingMetadata := hstore.Hstore{
			Map: map[string]sql.NullString{
				"key1": {String: "1", Valid: true},
				"key2": {String: "a", Valid: true},
				"key3": {String: "true", Valid: true},
				"key4": {},
				"key5": {String: "isn't this, \"complicated\"", Valid: true},
			},
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
			_, err = db.ExecContext(
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

		err = tx.Commit()
		require.NoError(t, err)
	})

	t.Run("Update", func(t *testing.T) {
		tx, err := db.BeginTxx(ctx, nil)
		require.NoError(t, err)
		defer func() {
			_ = tx.Rollback()
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
		physicalThingTags := pq.Array([]string{
			"tag1",
			"tag2",
			"tag3",
			"isn't this, \"complicated\"",
		})

		cleanup := func() {
			_, err = db.ExecContext(
				ctx,
				`DELETE FROM physical_things WHERE name = $1 OR name = $2;`,
				insertPhysicalThingName,
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
				insertPhysicalExternalID,
				insertPhysicalThingName,
				insertPhysicalThingType,
				insertPhysicalThingTags,
				insertPhysicalThingMetadata,
				insertPhysicalThingRawData,
			),
		)
		require.NoError(t, err)

		physicalThingMetadata := hstore.Hstore{
			Map: map[string]sql.NullString{
				"key1": {String: "1", Valid: true},
				"key2": {String: "a", Valid: true},
				"key3": {String: "true", Valid: true},
				"key4": {},
				"key5": {String: "isn't this, \"complicated\"", Valid: true},
			},
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
			fmt.Sprintf(
				"external_id = '%v' AND name = '%v' AND type = '%v'",
				insertPhysicalExternalID,
				insertPhysicalThingName,
				insertPhysicalThingType,
			),
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

		err = tx.Commit()
		require.NoError(t, err)
	})
}
