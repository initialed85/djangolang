package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/initialed85/djangolang/pkg/helpers"
	"github.com/initialed85/djangolang/pkg/pg_types"
	"github.com/initialed85/djangolang/pkg/some_db"
	"github.com/initialed85/djangolang/pkg/types"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func TestSQLHelpers(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := helpers.GetDBFromEnvironment(ctx)
	if err != nil {
		require.NoError(t, err)
	}
	defer func() {
		_ = db.Close()
	}()

	t.Run("TestInsertSelectUpdateDelete", func(t *testing.T) {
		physicalThings, err := some_db.SelectPhysicalThings(
			ctx,
			db,
			some_db.PhysicalThingColumns,
			nil,
			nil,
			nil,
		)
		require.NoError(t, err)
		require.Len(t, physicalThings, 0)

		physicalThing := some_db.PhysicalThing{
			Name:     "SomePhysicalThing",
			Type:     "http://some-url.org",
			Tags:     make(pq.StringArray, 0),
			Metadata: make(pg_types.Hstore),
		}

		require.Equal(t, uuid.UUID{}, physicalThing.ID)

		err = physicalThing.Insert(ctx, db)
		require.NoError(t, err)
		require.NotEqual(t, 0, physicalThing.ID)

		err = physicalThing.Insert(ctx, db)
		require.Error(t, err)

		physicalThing.Name = "OtherPhysicalThing"
		err = physicalThing.Update(ctx, db)
		require.NoError(t, err)
		require.Equal(t, "OtherPhysicalThing", physicalThing.Name)

		physicalThings, err = some_db.SelectPhysicalThings(
			ctx,
			db,
			some_db.PhysicalThingColumns,
			nil,
			nil,
			nil,
		)
		require.NoError(t, err)
		require.Len(t, physicalThings, 1)
		require.Equal(t, "OtherPhysicalThing", physicalThings[0].Name)

		physicalThings, err = some_db.SelectPhysicalThings(
			ctx,
			db,
			some_db.PhysicalThingColumns,
			nil,
			nil,
			nil,
			types.Clause(fmt.Sprintf("%v = $1", some_db.PhysicalThingTableIDColumn), physicalThing.ID),
		)
		require.NoError(t, err)
		require.Len(t, physicalThings, 1)
		require.Equal(t, "OtherPhysicalThing", physicalThings[0].Name)

		physicalThings, err = some_db.SelectPhysicalThings(
			ctx,
			db,
			some_db.PhysicalThingColumns,
			nil,
			nil,
			nil,
			types.Clause(fmt.Sprintf("%v IN ($1)", some_db.PhysicalThingTableNameColumn), physicalThing.Name),
		)
		require.NoError(t, err)
		require.Len(t, physicalThings, 1)
		require.Equal(t, "OtherPhysicalThing", physicalThings[0].Name)

		// TODO
		// locationHistory := some_db.LocationHistory{
		// 	Timestamp: time.Now(),
		// 	Point:     geom.NewPointEmpty(geom.XYZ).MustSetCoords(geom.Coord{1.0, 2.0, 3.0}),
		// }
		//
		// err = locationHistory.Insert(ctx, db)
		// require.NoError(t, err)

		err = physicalThing.Delete(ctx, db)
		require.NoError(t, err)

		physicalThings, err = some_db.SelectPhysicalThings(
			ctx,
			db,
			some_db.PhysicalThingColumns,
			nil,
			nil,
			nil,
		)
		require.NoError(t, err)
		require.Len(t, physicalThings, 0)
	})
}
