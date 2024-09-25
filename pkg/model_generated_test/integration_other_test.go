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

	"github.com/initialed85/djangolang/pkg/config"
	"github.com/initialed85/djangolang/pkg/helpers"
	"github.com/initialed85/djangolang/pkg/model_generated"
	"github.com/initialed85/djangolang/pkg/query"
	"github.com/initialed85/djangolang/pkg/server"
	"github.com/stretchr/testify/require"
)

func TestIntegrationOther(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctx = query.WithMaxDepth(ctx, helpers.Ptr(0))

	db, err := config.GetDBFromEnvironment(ctx)
	if err != nil {
		require.NoError(t, err)
	}
	defer func() {
		db.Close()
	}()

	redisPool, err := config.GetRedisFromEnvironment()
	if err != nil {
		require.NoError(t, err)
	}
	defer func() {
		redisPool.Close()
	}()

	redisConn := redisPool.Get()
	defer func() {
		_ = redisConn.Close()
	}()

	httpClient := &HTTPClient{
		httpClient: &http.Client{
			Timeout: time.Second * 10,
		},
	}

	changes := make(chan server.Change, 1024)
	mu := new(sync.Mutex)
	lastChangeByTableName := make(map[string]server.Change)

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

	getLastChangeForTableName := func(tableName string) *server.Change {
		mu.Lock()
		defer mu.Unlock()

		change, ok := lastChangeByTableName[tableName]
		if !ok {
			return nil
		}

		return &change
	}

	go func() {
		os.Setenv("DJANGOLANG_NODE_NAME", "model_generated_integration_other_test")
		err := model_generated.RunServer(ctx, changes, "127.0.0.1:2020", db, redisPool, nil, nil, nil)
		if err != nil {
			log.Printf("model_generated.RunServer failed; %v", err)
		}
	}()
	runtime.Gosched()
	time.Sleep(time.Second * 5)

	require.Eventually(
		t,
		func() bool {
			resp, err := httpClient.Get("http://localhost:2020/cameras")
			if err != nil {
				return false
			}

			if resp.StatusCode != http.StatusOK {
				return false
			}

			return true
		},
		time.Second*10,
		time.Millisecond*100,
	)
	time.Sleep(time.Second * 1)

	cleanup := func() {
		_, _ = db.Exec(
			ctx,
			`DELETE FROM camera CASCADE;`,
		)
		if redisConn != nil {
			_, _ = redisConn.Do("FLUSHALL")
		}
	}
	defer cleanup()

	camera1Name := "IntegrationCamera1Name"
	camera1StreamURLA := "IntegrationCamera1StreamURL"
	camera1StreamURLB := "IntegrationCamera1StreamURL"
	camera1ItemA := []map[string]any{
		{
			"name":       camera1Name,
			"stream_url": camera1StreamURLA,
		},
	}

	videoTimestamp, err := time.Parse(time.RFC3339, "2024-07-19T11:45:00+08:00")
	require.NoError(t, err)

	video1Name := "IntegrationVideo1Name"
	video1ItemA := []map[string]any{
		{
			"file_name":  video1Name,
			"started_at": videoTimestamp.Format(time.RFC3339Nano),
			"camera_id":  nil,
		},
	}

	detectionTimestamp, err := time.Parse(time.RFC3339, "2024-07-19T11:45:00+08:00")
	require.NoError(t, err)

	detection1Centroid := map[string]any{
		"X": 1.337,
		"Y": 69.420,
	}

	detection1BoundingBox := []map[string]any{
		{
			"X": 1.337,
			"Y": 69.420,
		},
	}

	detection1Item := []map[string]any{
		{
			"seen_at":      detectionTimestamp,
			"centroid":     detection1Centroid,
			"bounding_box": detection1BoundingBox,
			"camera_id":    nil,
			"video_id":     nil,
		},
	}

	// detection2Centroid := map[string]any{
	// 	"X": 1.337,
	// 	"Y": 69.420,
	// }

	// detection2BoundingBox := []map[string]any{
	// 	{
	// 		"X": 0.0,
	// 		"Y": 0.0,
	// 	},
	// 	{
	// 		"X": 1.0,
	// 		"Y": 0.0,
	// 	},
	// 	{
	// 		"X": 1.0,
	// 		"Y": 1.0,
	// 	},
	// 	{
	// 		"X": 0.0,
	// 		"Y": 1.0,
	// 	},
	// 	{
	// 		"X": 0.0,
	// 		"Y": 0.0,
	// 	},
	// }

	// detection2Item := []map[string]any{
	// 	{
	// 		"seen_at":      detectionTimestamp,
	// 		"centroid":     detection2Centroid,
	// 		"bounding_box": detection2BoundingBox,
	// 		"camera_id":    nil,
	// 		"video_id":     nil,
	// 	},
	// }

	setup := func() (*model_generated.Camera, *model_generated.Video, *model_generated.Detection, *model_generated.Detection) {
		//
		// Camera
		//

		cameraItemJSON, err := json.Marshal(camera1ItemA)
		require.NoError(t, err)

		resp, err := httpClient.Post(
			"http://localhost:2020/cameras",
			"application/json",
			bytes.NewReader(cameraItemJSON),
		)
		require.NoError(t, err)
		respBody, _ := io.ReadAll(resp.Body)
		require.NoError(t, err, string(respBody))
		require.Equal(t, http.StatusCreated, resp.StatusCode, string(respBody))

		require.Eventually(
			t,
			func() bool {
				change := getLastChangeForTableName(model_generated.CameraTable)
				if change == nil {
					return false
				}

				if change.Item["name"] != camera1Name {
					return false
				}

				return true
			},
			time.Second*10,
			time.Millisecond*10,
			"failed to confirm Camera",
		)

		tx, err := db.Begin(ctx)
		require.NoError(t, err)
		defer func() {
			_ = tx.Rollback(ctx)
		}()

		camera1, count, totalCount, page, totalPages, err := model_generated.SelectCamera(ctx, tx, "name = $$??", camera1Name)
		require.NoError(t, err)
		require.Equal(t, camera1Name, camera1Name)
		require.Equal(t, int64(1), count)
		require.Equal(t, int64(1), totalCount)
		require.Equal(t, int64(1), page)
		require.Equal(t, int64(1), totalPages)

		err = tx.Commit(ctx)
		require.NoError(t, err)

		//
		// Video
		//

		video1ItemA[0]["camera_id"] = &camera1.ID
		videoItemJSON, err := json.Marshal(video1ItemA)
		require.NoError(t, err)

		resp, err = httpClient.Post(
			"http://localhost:2020/videos",
			"application/json",
			bytes.NewReader(videoItemJSON),
		)
		respBody, _ = io.ReadAll(resp.Body)
		require.NoError(t, err, string(respBody))
		require.Equal(t, http.StatusCreated, resp.StatusCode, string(respBody))

		require.Eventually(
			t,
			func() bool {
				change := getLastChangeForTableName(model_generated.VideoTable)
				if change == nil {
					return false
				}

				if change.Item["file_name"] != video1Name {
					return false
				}

				return true
			},
			time.Second*10,
			time.Millisecond*10,
			"failed to confirm Video",
		)

		tx, err = db.Begin(ctx)
		require.NoError(t, err)
		defer func() {
			_ = tx.Rollback(ctx)
		}()

		video1, count, totalCount, page, totalPages, err := model_generated.SelectVideo(ctx, tx, "file_name = $$??", video1Name)
		require.NoError(t, err)
		require.Equal(t, video1Name, video1Name)
		require.Equal(t, int64(1), count)
		require.Equal(t, int64(1), totalCount)
		require.Equal(t, int64(1), page)
		require.Equal(t, int64(1), totalPages)

		err = tx.Commit(ctx)
		require.NoError(t, err)

		//
		// Detection1
		//

		detection1Item[0]["camera_id"] = &camera1.ID
		detection1Item[0]["video_id"] = &video1.ID
		detection1ItemJSON, err := json.Marshal(detection1Item)
		require.NoError(t, err)

		resp, err = httpClient.Post(
			"http://localhost:2020/detections",
			"application/json",
			bytes.NewReader(detection1ItemJSON),
		)
		respBody, _ = io.ReadAll(resp.Body)
		var x any
		_ = json.Unmarshal(respBody, &x)
		respBody, _ = json.MarshalIndent(x, "", "  ")
		require.NoError(t, err, string(respBody))
		require.Equal(t, http.StatusCreated, resp.StatusCode, string(respBody))

		require.Eventually(
			t,
			func() bool {
				change := getLastChangeForTableName(model_generated.DetectionTable)
				if change == nil {
					return false
				}

				object, ok := change.Object.(*model_generated.Detection)
				if !ok {
					return false
				}

				if object.CameraID != camera1.ID {
					return false
				}

				return true
			},
			time.Second*10,
			time.Millisecond*10,
			"failed to confirm Detection",
		)

		tx, err = db.Begin(ctx)
		require.NoError(t, err)
		defer func() {
			_ = tx.Rollback(ctx)
		}()

		detection1, count, totalCount, page, totalPages, err := model_generated.SelectDetection(ctx, tx, "camera_id = $$??", camera1.ID)
		require.NoError(t, err)
		require.Equal(t, int64(1), count)
		require.Equal(t, int64(1), totalCount)
		require.Equal(t, int64(1), page)
		require.Equal(t, int64(1), totalPages)

		require.NotNil(t, detection1.CameraIDObject)
		require.Equal(t, camera1.ID, detection1.CameraIDObject.ID)

		err = tx.Commit(ctx)
		require.NoError(t, err)

		// //
		// // Detection2
		// //

		// detection2Item[0]["camera_id"] = &camera1.ID
		// detection2Item[0]["video_id"] = &video1.ID
		// detection2ItemJSON, err := json.Marshal(detection2Item)
		// require.NoError(t, err)

		// resp, err = httpClient.Post(
		// 	"http://localhost:2020/detections",
		// 	"application/json",
		// 	bytes.NewReader(detection2ItemJSON),
		// )
		// respBody, _ = io.ReadAll(resp.Body)
		// _ = json.Unmarshal(respBody, &x)
		// respBody, _ = json.MarshalIndent(x, "", "  ")
		// require.NoError(t, err, string(respBody))
		// require.Equal(t, http.StatusCreated, resp.StatusCode, string(respBody))

		// require.Eventually(
		// 	t,
		// 	func() bool {
		// 		change := getLastChangeForTableName(model_generated.DetectionTable)
		// 		if change == nil {
		// 			return false
		// 		}

		// 		object, ok := change.Object.(*model_generated.Detection)
		// 		if !ok {
		// 			return false
		// 		}

		// 		if object.CameraID != camera1.ID {
		// 			return false
		// 		}

		// 		return true
		// 	},
		// 	time.Second*1,
		// 	time.Millisecond*10,
		// 	"failed to confirm Detection",
		// )

		// tx, err = db.Begin(ctx)
		// require.NoError(t, err)
		// defer func() {
		// 	_ = tx.Rollback(ctx)
		// }()

		// 2, count, totalCount, page, totalPages, err := model_generated.SelectDetection(ctx, tx, "camera_id = $$?? AND id != $$??", camera1.ID, detection1.ID)
		// require.NoError(t, err)
		// require.NotNil(t, detection2.CameraIDObject)
		// require.Equal(t, camera1.ID, detection1.CameraIDObject.ID)

		// err = tx.Commit(ctx)
		// require.NoError(t, err)

		//
		// reloads
		//

		tx, err = db.Begin(ctx)
		require.NoError(t, err)
		defer func() {
			_ = tx.Rollback(ctx)
		}()

		_ = camera1.Reload(ctx, tx)
		_ = video1.Reload(ctx, tx)
		_ = detection1.Reload(ctx, tx)

		_ = tx.Commit(ctx)

		var detection2 *model_generated.Detection

		return camera1, video1, detection1, detection2
	}
	camera1, _, _, _ := setup()

	t.Run("Camera", func(t *testing.T) {
		resp, err := httpClient.Get(
			fmt.Sprintf("http://localhost:2020/cameras/%s", camera1.ID.String()),
		)
		respBody, _ := io.ReadAll(resp.Body)
		require.NoError(t, err, string(respBody))
		require.Equal(t, http.StatusOK, resp.StatusCode, string(respBody))

		camera1ItemB := map[string]any{
			"stream_url": camera1StreamURLB,
		}

		cameraItemBJSON, err := json.Marshal(camera1ItemB)
		require.NoError(t, err)

		_ = cameraItemBJSON

		resp, err = httpClient.Patch(
			fmt.Sprintf("http://localhost:2020/cameras/%s", camera1.ID.String()),
			"application/json",
			bytes.NewReader(cameraItemBJSON),
		)
		respBody, _ = io.ReadAll(resp.Body)
		require.NoError(t, err, string(respBody))
		require.Equal(t, http.StatusOK, resp.StatusCode, string(respBody))

		tx, err := db.Begin(ctx)
		require.NoError(t, err)
		defer func() {
			_ = tx.Rollback(ctx)
		}()
		err = camera1.Reload(ctx, tx)
		require.NoError(t, err)
		err = tx.Commit(ctx)
		require.NoError(t, err)

		require.Equal(t, camera1StreamURLB, camera1.StreamURL)

		resp, err = httpClient.Get("http://localhost:2020/cameras")
		respBody, _ = io.ReadAll(resp.Body)
		require.NoError(t, err, string(respBody))
		require.Equal(t, http.StatusOK, resp.StatusCode, string(respBody))

		var items any
		err = json.Unmarshal(respBody, &items)
		require.NoError(t, err, string(respBody))

		b, _ := json.MarshalIndent(items, "", "  ")
		log.Printf("b: %v", string(b))
	})
}
