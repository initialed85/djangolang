package test

import (
	"context"
	"testing"

	"github.com/initialed85/djangolang/pkg/helpers"
	"github.com/initialed85/djangolang/pkg/some_db"
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
		cameras, err := some_db.SelectCameras(
			ctx,
			db,
			some_db.CameraColumns,
			nil,
			nil,
			nil,
		)
		require.NoError(t, err)
		require.Len(t, cameras, 0)

		camera := some_db.Camera{
			Name:      "SomeCamera",
			StreamURL: "http://some-url.org",
		}

		require.Equal(t, int64(0), camera.ID)

		err = camera.Insert(ctx, db)
		require.NoError(t, err)
		require.NotEqual(t, 0, camera.ID)

		err = camera.Insert(ctx, db)
		require.Error(t, err)

		camera.Name = "OtherCamera"
		err = camera.Update(ctx, db)
		require.NoError(t, err)
		require.Equal(t, "OtherCamera", camera.Name)

		cameras, err = some_db.SelectCameras(
			ctx,
			db,
			some_db.CameraColumns,
			nil,
			nil,
			nil,
		)
		require.NoError(t, err)
		require.Len(t, cameras, 1)
		require.Equal(t, "OtherCamera", cameras[0].Name)

		err = camera.Delete(ctx, db)
		require.NoError(t, err)

		cameras, err = some_db.SelectCameras(
			ctx,
			db,
			some_db.CameraColumns,
			nil,
			nil,
			nil,
		)
		require.NoError(t, err)
		require.Len(t, cameras, 0)
	})
}