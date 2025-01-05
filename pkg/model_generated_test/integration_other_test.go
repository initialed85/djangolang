package model_generated_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/initialed85/djangolang/pkg/helpers"
	"github.com/initialed85/djangolang/pkg/model_generated"
	"github.com/initialed85/djangolang/pkg/query"
	"github.com/initialed85/djangolang/pkg/server"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

func testIntegrationOther(
	t *testing.T,
	ctx context.Context,
	db *pgxpool.Pool,
	redisConn redis.Conn,
	mu *sync.Mutex,
	lastChangeByTableName map[string]*server.Change,
	httpClient *HTTPClient,
	getLastChangeForTableName func(tableName string) *server.Change,
) {
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

	setup := func() (*model_generated.Camera, *model_generated.Video, *model_generated.Detection, *model_generated.Detection) {
		ctx := query.WithMaxDepth(ctx, helpers.Ptr(0))

		//
		// Camera
		//

		cameraItemJSON, err := json.Marshal(camera1ItemA)
		require.NoError(t, err)

		resp, err := httpClient.Post(
			"http://localhost:5050/cameras",
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
			"http://localhost:5050/videos",
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
			"http://localhost:5050/detections",
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
			fmt.Sprintf("http://localhost:5050/cameras/%s", camera1.ID.String()),
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
			fmt.Sprintf("http://localhost:5050/cameras/%s", camera1.ID.String()),
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

		getVideos := func(rawQueryParams string) *server.Response[model_generated.Video] {
			resp, err := httpClient.Get(fmt.Sprintf("http://localhost:5050/videos%s", rawQueryParams))
			respBody, _ = io.ReadAll(resp.Body)
			require.NoError(t, err, string(respBody))
			require.Equal(t, http.StatusOK, resp.StatusCode, string(respBody))

			var typedResp *server.Response[model_generated.Video]
			err = json.Unmarshal(respBody, &typedResp)
			require.NoError(t, err, string(respBody))
			require.GreaterOrEqual(t, len(typedResp.Objects), 1)

			b, _ := json.MarshalIndent(typedResp, "", "  ")
			// log.Printf("b: %v", string(b))

			_ = b

			return typedResp
		}

		var videosResp *server.Response[model_generated.Video]

		videosResp = getVideos("")
		require.Nil(t, videosResp.Objects[0].CameraIDObject)

		videosResp = getVideos("?depth=0")
		require.NotNil(t, videosResp.Objects[0].CameraIDObject)
		require.NotNil(t, videosResp.Objects[0].CameraIDObject.ReferencedByVideoCameraIDObjects)

		videosResp = getVideos("?depth=1")
		require.Nil(t, videosResp.Objects[0].CameraIDObject)

		videosResp = getVideos("?depth=2")
		require.NotNil(t, videosResp.Objects[0].CameraIDObject)
		require.Nil(t, videosResp.Objects[0].CameraIDObject.ReferencedByVideoCameraIDObjects)

		videosResp = getVideos("?depth=3")
		require.NotNil(t, videosResp.Objects[0].CameraIDObject)
		require.NotNil(t, videosResp.Objects[0].CameraIDObject.ReferencedByVideoCameraIDObjects)

		videosResp = getVideos("?camera__load=")
		require.NotNil(t, videosResp.Objects[0].CameraIDObject)
		require.Nil(t, videosResp.Objects[0].ReferencedByDetectionVideoIDObjects)

		videosResp = getVideos("?referenced_by_detection__load=")
		require.Nil(t, videosResp.Objects[0].CameraIDObject)
		require.NotNil(t, videosResp.Objects[0].ReferencedByDetectionVideoIDObjects)
	})

	t.Run("RouterClaim", func(t *testing.T) {
		now := time.Now().UTC()
		later := now.Add(time.Second * 1)
		idA := "00000000-0000-0000-0000-000000000001"
		idB := "00000000-0000-0000-0000-000000000002"

		//
		// can't claim something that isn't claimable
		//

		r, err := httpClient.Post(
			"http://127.0.0.1:5050/claim-logical-thing",
			"application/json",
			bytes.NewReader([]byte(fmt.Sprintf(`{"until": "%s", "by": "%s", "timeout_seconds": 2}`, later.Format(time.RFC3339Nano), idA))),
		)
		require.NoError(t, err)
		require.Equal(t, http.StatusNotFound, r.StatusCode)

		//
		// can claim something that is
		//

		r, err = httpClient.Post(
			"http://127.0.0.1:5050/claim-camera",
			"application/json",
			bytes.NewReader([]byte(fmt.Sprintf(`{"until": "%s", "by": "%s", "timeout_seconds": 2}`, later.Format(time.RFC3339Nano), idA))),
		)
		require.NoError(t, err)
		b, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, r.StatusCode, []string{string(b), camera1.ID.String()})

		var response server.Response[any]
		err = json.Unmarshal(b, &response)
		require.NoError(t, err)

		require.True(t, response.Success)
		require.Empty(t, response.Error)
		require.NotEmpty(t, response.Objects)

		objects := response.Objects
		require.Equal(t, 1, len(objects))

		//
		// can't claim if there's nothing left
		//

		r, err = httpClient.Post(
			"http://127.0.0.1:5050/claim-camera",
			"application/json",
			bytes.NewReader([]byte(fmt.Sprintf(`{"until": "%s", "by": "%s", "timeout_seconds": 2}`, later.Format(time.RFC3339Nano), idB))),
		)
		require.NoError(t, err)
		b, err = io.ReadAll(r.Body)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, r.StatusCode, []string{string(b), camera1.ID.String()})

		response = server.Response[any]{}
		err = json.Unmarshal(b, &response)
		require.NoError(t, err)

		require.True(t, response.Success)
		require.Empty(t, response.Error)
		require.Empty(t, response.Objects)

		//
		// can't reclaim via this mechanism
		//

		r, err = httpClient.Post(
			"http://127.0.0.1:5050/claim-camera",
			"application/json",
			bytes.NewReader([]byte(fmt.Sprintf(`{"until": "%s", "by": "%s", "timeout_seconds": 2}`, later.Format(time.RFC3339Nano), idA))),
		)
		require.NoError(t, err)
		b, err = io.ReadAll(r.Body)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, r.StatusCode, []string{string(b), camera1.ID.String()})

		response = server.Response[any]{}
		err = json.Unmarshal(b, &response)
		require.NoError(t, err)

		require.True(t, response.Success)
		require.Empty(t, response.Error)
		require.Empty(t, response.Objects)

		//
		// can't claim the same thing if it's claimed by somebody else
		//

		r, err = httpClient.Post(
			fmt.Sprintf("http://127.0.0.1:5050/cameras/%s/claim", camera1.ID),
			"application/json",
			bytes.NewReader([]byte(fmt.Sprintf(`{"until": "%s", "by": "%s", "timeout_seconds": 2}`, later.Format(time.RFC3339Nano), idB))),
		)
		require.NoError(t, err)
		b, err = io.ReadAll(r.Body)
		require.NoError(t, err)
		require.Equal(t, http.StatusInternalServerError, r.StatusCode, []string{string(b), camera1.ID.String()})

		response = server.Response[any]{}
		err = json.Unmarshal(b, &response)
		require.NoError(t, err)

		require.False(t, response.Success)
		require.NotEmpty(t, response.Error)
		require.Nil(t, response.Objects)

		//
		// can reclaim via this mechanism even if it isn't ours, once the existing claim has expired
		//

		time.Sleep(time.Second * 3)

		r, err = httpClient.Post(
			fmt.Sprintf("http://127.0.0.1:5050/cameras/%s/claim", camera1.ID),
			"application/json",
			bytes.NewReader([]byte(fmt.Sprintf(`{"until": "%s", "by": "%s", "timeout_seconds": 2}`, later.Format(time.RFC3339Nano), idB))),
		)
		require.NoError(t, err)
		b, err = io.ReadAll(r.Body)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, r.StatusCode, []string{string(b), camera1.ID.String()})

		response = server.Response[any]{}
		err = json.Unmarshal(b, &response)
		require.NoError(t, err)

		require.True(t, response.Success)
		require.Empty(t, response.Error)
		require.NotEmpty(t, response.Objects)
	})
}
