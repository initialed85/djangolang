package test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"testing"
	"time"

	"github.com/initialed85/djangolang/pkg/helpers"
	"github.com/initialed85/djangolang/pkg/some_db"
	"github.com/stretchr/testify/require"
)

func TestEndpoints(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := helpers.GetDBFromEnvironment(ctx)
	if err != nil {
		require.NoError(t, err)
	}
	defer func() {
		_ = db.Close()
	}()

	_, err = db.Query("TRUNCATE TABLE camera CASCADE;")
	require.NoError(t, err)

	port, err := helpers.GetPort()
	require.NoError(t, err)

	go func() {
		err = some_db.RunServer(ctx)
		require.NoError(t, err)
	}()
	runtime.Gosched()
	time.Sleep(time.Second * 1)

	httpClient := http.Client{
		Timeout: time.Second * 5,
	}

	t.Run("TestPostGetDelete", func(t *testing.T) {
		cameras := make([]*some_db.Camera, 0)

		camera := &some_db.Camera{
			Name:      "SomeCamera",
			StreamURL: "http://some-url.org",
		}

		cameraJSON, err := json.Marshal(camera)
		require.NoError(t, err)

		resp, err := httpClient.Get(fmt.Sprintf("http://localhost:%v/camera", port))
		require.NoError(t, err)
		b, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		defer func() {
			_ = resp.Body.Close()
		}()
		require.Equal(t, http.StatusOK, resp.StatusCode, string(b))
		err = json.Unmarshal(b, &cameras)
		require.NoError(t, err)
		require.Len(t, cameras, 0)

		resp, err = httpClient.Post(
			fmt.Sprintf("http://localhost:%v/camera", port),
			"application/json",
			bytes.NewBuffer(cameraJSON),
		)
		require.NoError(t, err)
		b, err = io.ReadAll(resp.Body)
		require.NoError(t, err)
		defer func() {
			_ = resp.Body.Close()
		}()
		require.Equal(t, http.StatusCreated, resp.StatusCode, string(b))
		err = json.Unmarshal(b, &camera)
		require.NoError(t, err)
		require.NotZero(t, camera.ID)
		require.Equal(t, "SomeCamera", camera.Name)
		require.Equal(t, "http://some-url.org", camera.StreamURL)

		resp, err = httpClient.Post(
			fmt.Sprintf("http://localhost:%v/camera", port),
			"application/json",
			bytes.NewBuffer(cameraJSON),
		)
		require.NoError(t, err)
		b, err = io.ReadAll(resp.Body)
		require.NoError(t, err)
		defer func() {
			_ = resp.Body.Close()
		}()
		require.Equal(t, http.StatusConflict, resp.StatusCode, string(b))
		err = json.Unmarshal(b, &camera)
		require.NoError(t, err)
		require.NotZero(t, camera.ID)
		require.Equal(t, "SomeCamera", camera.Name)
		require.Equal(t, "http://some-url.org", camera.StreamURL)

		req, err := http.NewRequest(
			http.MethodDelete,
			fmt.Sprintf("http://localhost:%v/camera/%v", port, camera.ID),
			bytes.NewBuffer(cameraJSON),
		)
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		resp, err = httpClient.Do(req)
		require.NoError(t, err)
		b, err = io.ReadAll(resp.Body)
		require.NoError(t, err)
		defer func() {
			_ = resp.Body.Close()
		}()
		require.Equal(t, http.StatusNoContent, resp.StatusCode, string(b))
		err = json.Unmarshal(b, &struct{}{})
		require.Error(t, err)

		resp, err = httpClient.Get(fmt.Sprintf("http://localhost:%v/camera", port))
		require.NoError(t, err)
		b, err = io.ReadAll(resp.Body)
		require.NoError(t, err)
		defer func() {
			_ = resp.Body.Close()
		}()
		require.Equal(t, resp.StatusCode, http.StatusOK, string(b))
		err = json.Unmarshal(b, &cameras)
		require.NoError(t, err)
		require.Len(t, cameras, 0)

		req, err = http.NewRequest(
			http.MethodDelete,
			fmt.Sprintf("http://localhost:%v/camera/%v", port, camera.ID),
			bytes.NewBuffer(cameraJSON),
		)
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		resp, err = httpClient.Do(req)
		require.NoError(t, err)
		b, err = io.ReadAll(resp.Body)
		require.NoError(t, err)
		defer func() {
			_ = resp.Body.Close()
		}()
		require.Equal(t, http.StatusNotFound, resp.StatusCode, string(b))
	})
}
