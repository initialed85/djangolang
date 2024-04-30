package query

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	_helpers "github.com/initialed85/djangolang/internal/helpers"
	"github.com/initialed85/djangolang/pkg/helpers"
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

	tx, err := db.BeginTxx(ctx, nil)
	require.NoError(t, err)
	defer func() {
		_ = tx.Rollback()
	}()

	t.Run("Select", func(t *testing.T) {
		physicalExternalID := "SomePhysicalThingExternalID1"
		physicalThingName := "SomePhysicalThingName1"
		physicalThingType := "SomePhysicalThingType1"
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
}
