package templates

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/initialed85/djangolang/pkg/helpers"
	"github.com/jmoiron/sqlx"
	"golang.org/x/exp/maps"
)

var EventTable = "event"
var EventIDColumn = "id"
var EventStartTimestampColumn = "start_timestamp"
var EventEndTimestampColumn = "end_timestamp"
var EventDurationColumn = "duration"
var EventOriginalVideoIDColumn = "original_video_id"
var EventThumbnailImageIDColumn = "thumbnail_image_id"
var EventProcessedVideoIDColumn = "processed_video_id"
var EventSourceCameraIDColumn = "source_camera_id"
var EventStatusColumn = "status"
var EventColumns = []string{"id", "start_timestamp", "end_timestamp", "duration", "original_video_id", "thumbnail_image_id", "processed_video_id", "source_camera_id", "status"}
var EventTransformedColumns = []string{"id", "start_timestamp", "end_timestamp", "extract(microseconds FROM duration)::numeric * 1000 AS duration", "original_video_id", "thumbnail_image_id", "processed_video_id", "source_camera_id", "status"}

type Event struct {
	ID                     int64         `json:"id" db:"id"`
	StartTimestamp         time.Time     `json:"start_timestamp" db:"start_timestamp"`
	EndTimestamp           time.Time     `json:"end_timestamp" db:"end_timestamp"`
	Duration               time.Duration `json:"duration" db:"duration"`
	OriginalVideoID        int64         `json:"original_video_id" db:"original_video_id"`
	OriginalVideoIDObject  *Video        `json:"original_video_id_object,omitempty"`
	ThumbnailImageID       int64         `json:"thumbnail_image_id" db:"thumbnail_image_id"`
	ThumbnailImageIDObject *Image        `json:"thumbnail_image_id_object,omitempty"`
	ProcessedVideoID       *int64        `json:"processed_video_id" db:"processed_video_id"`
	ProcessedVideoIDObject *Video        `json:"processed_video_id_object,omitempty"`
	SourceCameraID         int64         `json:"source_camera_id" db:"source_camera_id"`
	SourceCameraIDObject   *Camera       `json:"source_camera_id_object,omitempty"`
	Status                 string        `json:"status" db:"status"`
}

func SelectEvents(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]*Event, error) {
	mu.RLock()
	debug := actualDebug
	mu.RUnlock()

	var buildStart int64
	var buildStop int64
	var execStart int64
	var execStop int64
	var scanStart int64
	var scanStop int64
	var foreignObjectStop int64
	var foreignObjectStart int64

	var sql string
	var columnCount int = len(columns)
	var rowCount int64
	defer func() {
		if !debug {
			return
		}

		buildDuration := 0.0
		execDuration := 0.0
		scanDuration := 0.0
		foreignObjectDuration := 0.0

		if buildStop > 0 {
			buildDuration = float64(buildStop-buildStart) * 1e-9
		}

		if execStop > 0 {
			execDuration = float64(execStop-execStart) * 1e-9
		}

		if scanStop > 0 {
			scanDuration = float64(scanStop-scanStart) * 1e-9
		}

		if foreignObjectStop > 0 {
			foreignObjectDuration = float64(foreignObjectStop-foreignObjectStart) * 1e-9
		}

		if debug {
			logger.Printf(
				"selected %v columns, %v rows; %.3f seconds to build, %.3f seconds to execute, %.3f seconds to scan, %.3f seconds to load foreign objects; sql:\n%v",
				columnCount, rowCount, buildDuration, execDuration, scanDuration, foreignObjectDuration, sql,
			)
		}
	}()

	buildStart = time.Now().UnixNano()

	where := strings.TrimSpace(strings.Join(wheres, " AND "))
	if len(where) > 0 {
		where = " WHERE " + where
	}

	sql = fmt.Sprintf(
		"SELECT %v FROM event%v",
		strings.Join(columns, ", "),
		where,
	)

	if orderBy != nil {
		actualOrderBy := strings.TrimSpace(*orderBy)
		if len(actualOrderBy) > 0 {
			sql += fmt.Sprintf(" ORDER BY %v", actualOrderBy)
		}
	}

	if limit != nil && *limit >= 0 {
		sql += fmt.Sprintf(" LIMIT %v", *limit)
	}

	buildStop = time.Now().UnixNano()

	execStart = time.Now().UnixNano()

	selectCtx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()

	rows, err := db.QueryxContext(selectCtx, sql)
	if err != nil {
		return nil, err
	}

	execStop = time.Now().UnixNano()

	scanStart = time.Now().UnixNano()

	videoByOriginalVideoID := make(map[int64]*Video, 0)
	imageByThumbnailImageID := make(map[int64]*Image, 0)
	videoByProcessedVideoID := make(map[int64]*Video, 0)
	cameraBySourceCameraID := make(map[int64]*Camera, 0)

	items := make([]*Event, 0)
	for rows.Next() {
		rowCount++

		var item Event
		err = rows.StructScan(&item)
		if err != nil {
			return nil, err
		}

		// this is where any post-scan processing would appear (if required)

		videoByOriginalVideoID[item.OriginalVideoID] = nil
		imageByThumbnailImageID[item.ThumbnailImageID] = nil
		videoByProcessedVideoID[helpers.Deref(item.ProcessedVideoID)] = nil
		cameraBySourceCameraID[item.SourceCameraID] = nil

		items = append(items, &item)
	}

	scanStop = time.Now().UnixNano()

	foreignObjectStart = time.Now().UnixNano()

	idsForOriginalVideoID := make([]string, 0)
	for _, id := range maps.Keys(videoByOriginalVideoID) {
		idsForOriginalVideoID = append(idsForOriginalVideoID, fmt.Sprintf("%v", id))
	}

	if len(idsForOriginalVideoID) > 0 {
		rowsForOriginalVideoID, err := SelectVideos(
			ctx,
			db,
			VideoTransformedColumns,
			nil,
			nil,
			nil,
			fmt.Sprintf("id IN (%v)", strings.Join(idsForOriginalVideoID, ", ")),
		)
		if err != nil {
			return nil, err
		}

		for _, row := range rowsForOriginalVideoID {
			videoByOriginalVideoID[row.ID] = row
		}

		for _, item := range items {
			item.OriginalVideoIDObject = videoByOriginalVideoID[item.OriginalVideoID]
		}
	}

	idsForThumbnailImageID := make([]string, 0)
	for _, id := range maps.Keys(imageByThumbnailImageID) {
		idsForThumbnailImageID = append(idsForThumbnailImageID, fmt.Sprintf("%v", id))
	}

	if len(idsForThumbnailImageID) > 0 {
		rowsForThumbnailImageID, err := SelectImages(
			ctx,
			db,
			ImageTransformedColumns,
			nil,
			nil,
			nil,
			fmt.Sprintf("id IN (%v)", strings.Join(idsForThumbnailImageID, ", ")),
		)
		if err != nil {
			return nil, err
		}

		for _, row := range rowsForThumbnailImageID {
			imageByThumbnailImageID[row.ID] = row
		}

		for _, item := range items {
			item.ThumbnailImageIDObject = imageByThumbnailImageID[item.ThumbnailImageID]
		}
	}

	idsForProcessedVideoID := make([]string, 0)
	for _, id := range maps.Keys(videoByProcessedVideoID) {
		idsForProcessedVideoID = append(idsForProcessedVideoID, fmt.Sprintf("%v", id))
	}

	if len(idsForProcessedVideoID) > 0 {
		rowsForProcessedVideoID, err := SelectVideos(
			ctx,
			db,
			VideoTransformedColumns,
			nil,
			nil,
			nil,
			fmt.Sprintf("id IN (%v)", strings.Join(idsForProcessedVideoID, ", ")),
		)
		if err != nil {
			return nil, err
		}

		for _, row := range rowsForProcessedVideoID {
			videoByProcessedVideoID[row.ID] = row
		}

		for _, item := range items {
			item.ProcessedVideoIDObject = videoByProcessedVideoID[helpers.Deref(item.ProcessedVideoID)]
		}
	}

	idsForSourceCameraID := make([]string, 0)
	for _, id := range maps.Keys(cameraBySourceCameraID) {
		idsForSourceCameraID = append(idsForSourceCameraID, fmt.Sprintf("%v", id))
	}

	if len(idsForSourceCameraID) > 0 {
		rowsForSourceCameraID, err := SelectCameras(
			ctx,
			db,
			CameraTransformedColumns,
			nil,
			nil,
			nil,
			fmt.Sprintf("id IN (%v)", strings.Join(idsForSourceCameraID, ", ")),
		)
		if err != nil {
			return nil, err
		}

		for _, row := range rowsForSourceCameraID {
			cameraBySourceCameraID[row.ID] = row
		}

		for _, item := range items {
			item.SourceCameraIDObject = cameraBySourceCameraID[item.SourceCameraID]
		}
	}

	foreignObjectStop = time.Now().UnixNano()

	return items, nil
}

func genericSelectEvents(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]any, error) {
	items, err := SelectEvents(ctx, db, columns, orderBy, limit, offset, wheres...)
	if err != nil {
		return nil, err
	}

	genericItems := make([]any, 0)
	for _, item := range items {
		genericItems = append(genericItems, item)
	}

	return genericItems, nil
}
