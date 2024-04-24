package test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"sync"
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

	port, err := helpers.GetPort()
	require.NoError(t, err)

	go func() {
		err = some_db.RunServer(ctx)
		require.NoError(t, err)
	}()
	runtime.Gosched()
	time.Sleep(time.Second * 1)

	t.Run("TestPostGetPatchDelete", func(t *testing.T) {
		httpClient := http.Client{
			Timeout: time.Second * 5,
		}

		for i := 0; i < 50; i++ {
			cameras := make([]*some_db.Camera, 0)

			camera := &some_db.Camera{
				Name:      "SomeCamera",
				StreamURL: "http://some-url.org",
			}
			cameraJSON, err := json.Marshal(camera)
			require.NoError(t, err)

			otherCamera := &some_db.Camera{
				Name:      "OtherCamera",
				StreamURL: "http://some-url.org",
			}
			otherCameraJSON, err := json.Marshal(otherCamera)
			require.NoError(t, err)

			resp, err := httpClient.Get(fmt.Sprintf("http://localhost:%v/camera", port))
			require.NoError(t, err)
			b, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			_ = resp.Body.Close()
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
			_ = resp.Body.Close()
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
			_ = resp.Body.Close()
			require.Equal(t, http.StatusConflict, resp.StatusCode, string(b))
			err = json.Unmarshal(b, &camera)
			require.NoError(t, err)
			require.NotZero(t, camera.ID)
			require.Equal(t, "SomeCamera", camera.Name)
			require.Equal(t, "http://some-url.org", camera.StreamURL)

			req, err := http.NewRequest(
				http.MethodPatch,
				fmt.Sprintf("http://localhost:%v/camera/%v", port, camera.ID),
				bytes.NewBuffer(otherCameraJSON),
			)
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			resp, err = httpClient.Do(req)
			require.NoError(t, err)
			b, err = io.ReadAll(resp.Body)
			require.NoError(t, err)
			_ = resp.Body.Close()
			require.Equal(t, http.StatusOK, resp.StatusCode, string(b))
			err = json.Unmarshal(b, &camera)
			require.NoError(t, err)
			require.Equal(t, "OtherCamera", camera.Name)

			resp, err = httpClient.Post(
				fmt.Sprintf("http://localhost:%v/camera", port),
				"application/json",
				bytes.NewBuffer(otherCameraJSON),
			)
			require.NoError(t, err)
			b, err = io.ReadAll(resp.Body)
			require.NoError(t, err)
			_ = resp.Body.Close()
			require.Equal(t, http.StatusConflict, resp.StatusCode, string(b))
			err = json.Unmarshal(b, &camera)
			require.NoError(t, err)
			require.NotZero(t, camera.ID)
			require.Equal(t, "OtherCamera", camera.Name)
			require.Equal(t, "http://some-url.org", camera.StreamURL)

			req, err = http.NewRequest(
				http.MethodDelete,
				fmt.Sprintf("http://localhost:%v/camera/%v", port, camera.ID),
				nil,
			)
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			resp, err = httpClient.Do(req)
			require.NoError(t, err)
			b, err = io.ReadAll(resp.Body)
			require.NoError(t, err)
			_ = resp.Body.Close()
			require.Equal(t, http.StatusNoContent, resp.StatusCode, string(b))
			err = json.Unmarshal(b, &struct{}{})
			require.Error(t, err)

			resp, err = httpClient.Get(fmt.Sprintf("http://localhost:%v/camera", port))
			require.NoError(t, err)
			b, err = io.ReadAll(resp.Body)
			require.NoError(t, err)
			_ = resp.Body.Close()
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
			_ = resp.Body.Close()
			require.Equal(t, http.StatusNotFound, resp.StatusCode, string(b))
		}
	})

	t.Run("TestPerformance", func(t *testing.T) {
		wg := new(sync.WaitGroup)

		limiter := make(chan bool, 50)
		for i := 0; i < 50; i++ {
			limiter <- true
		}

		//
		// post
		//

		for i := 0; i < 1000; i++ {
			wg.Add(1)

			go func(i int) {
				defer wg.Done()

				<-limiter
				defer func() {
					limiter <- true
				}()

				httpClient := http.Client{
					Timeout: time.Second * 5,
				}

				camera := &some_db.Camera{
					Name:      fmt.Sprintf("SomeCamera-%v", i+1),
					StreamURL: "http://some-url.org",
				}

				b, err := json.Marshal(camera)
				require.NoError(t, err)

				resp, err := httpClient.Post(
					fmt.Sprintf("http://localhost:%v/camera", port),
					"application/json",
					bytes.NewBuffer(b),
				)
				require.NoError(t, err)

				b, err = io.ReadAll(resp.Body)
				require.NoError(t, err)
				defer func() {
					_ = resp.Body.Close()
				}()
				require.Equal(t, http.StatusCreated, resp.StatusCode, string(b))
			}(i)
		}

		wg.Wait()

		//
		// get
		//

		httpClient := http.Client{
			Timeout: time.Second * 5,
		}

		resp, err := httpClient.Get(fmt.Sprintf("http://localhost:%v/camera?limit=2000&order_by=id&order=asc", port))
		require.NoError(t, err)

		b, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		defer func() {
			_ = resp.Body.Close()
		}()
		require.Equal(t, http.StatusOK, resp.StatusCode, string(b))

		cameras := make([]*some_db.Camera, 0)
		err = json.Unmarshal(b, &cameras)
		require.NoError(t, err)
		require.Equal(t, 1000, len(cameras))

		//
		// patch
		//

		for i, camera := range cameras {
			wg.Add(1)

			go func(i int, camera *some_db.Camera) {
				defer wg.Done()

				<-limiter
				defer func() {
					limiter <- true
				}()

				httpClient := http.Client{
					Timeout: time.Second * 5,
				}

				camera.Name = fmt.Sprintf("OtherCamera-%v", i+1)

				b, err := json.Marshal(camera)
				require.NoError(t, err)

				req, err := http.NewRequest(
					http.MethodPatch,
					fmt.Sprintf("http://localhost:%v/camera/%v", port, camera.ID),
					bytes.NewBuffer(b),
				)
				require.NoError(t, err)
				req.Header.Set("Content-Type", "application/json")

				resp, err := httpClient.Do(req)
				require.NoError(t, err)

				b, err = io.ReadAll(resp.Body)
				require.NoError(t, err)
				defer func() {
					_ = resp.Body.Close()
				}()
				require.Equal(t, http.StatusOK, resp.StatusCode, string(b))

				err = json.Unmarshal(b, &camera)
				require.NoError(t, err)
				require.Equal(t, fmt.Sprintf("OtherCamera-%v", i+1), camera.Name)
			}(i, camera)
		}

		wg.Wait()

		//
		// delete
		//

		for _, camera := range cameras {
			wg.Add(1)

			go func(camera *some_db.Camera) {
				defer wg.Done()

				<-limiter
				defer func() {
					limiter <- true
				}()

				httpClient := http.Client{
					Timeout: time.Second * 5,
				}

				req, err := http.NewRequest(
					http.MethodDelete,
					fmt.Sprintf("http://localhost:%v/camera/%v", port, camera.ID),
					nil,
				)
				require.NoError(t, err)

				resp, err = httpClient.Do(req)
				require.NoError(t, err)

				require.Equal(t, http.StatusNoContent, resp.StatusCode, string(b))
			}(camera)
		}

		wg.Wait()
	})
}
