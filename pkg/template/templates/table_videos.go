package templates

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

var VideoViewTable = "videos"
var VideoViewTableIDColumn = "id"
var VideoViewTableStartTimestampColumn = "start_timestamp"
var VideoViewTableEndTimestampColumn = "end_timestamp"
var VideoViewTableDurationColumn = "duration"
var VideoViewTableSizeColumn = "size"
var VideoViewTableFilePathColumn = "file_path"
var VideoViewTableCameraIDColumn = "camera_id"
var VideoViewTableEventIDColumn = "event_id"
var VideoViewColumns = []string{"id", "start_timestamp", "end_timestamp", "duration", "size", "file_path", "camera_id", "event_id"}
var VideoViewTransformedColumns = []string{"id", "start_timestamp", "end_timestamp", "extract(microseconds FROM duration)::numeric * 1000 AS duration", "size", "file_path", "camera_id", "event_id"}

type VideoView struct {
	ID             *int64         `json:"id" db:"id"`
	StartTimestamp *time.Time     `json:"start_timestamp" db:"start_timestamp"`
	EndTimestamp   *time.Time     `json:"end_timestamp" db:"end_timestamp"`
	Duration       *time.Duration `json:"duration" db:"duration"`
	Size           *float64       `json:"size" db:"size"`
	FilePath       *string        `json:"file_path" db:"file_path"`
	CameraID       *int64         `json:"camera_id" db:"camera_id"`
	EventID        *int64         `json:"event_id" db:"event_id"`
}

func SelectVideosView(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]*VideoView, error) {
	mu.RLock()
	debug := actualDebug
	mu.RUnlock()

	key := "VideoView"

	path, _ := ctx.Value("path").(map[string]struct{})
	if path == nil {
		path = make(map[string]struct{}, 0)
	}

	// to avoid a stack overflow in the case of a recursive schema
	_, ok := path[key]
	if ok {
		return nil, nil
	}

	path[key] = struct{}{}

	ctx = context.WithValue(ctx, "path", path)

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
		"SELECT %v FROM videos%v",
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

	// this is where any foreign object ID map declarations would appear (if required)

	items := make([]*VideoView, 0)
	for rows.Next() {
		rowCount++

		var item VideoView
		err = rows.StructScan(&item)
		if err != nil {
			return nil, err
		}

		// this is where any post-scan processing would appear (if required)

		// this is where foreign object ID map populations would appear (if required)

		items = append(items, &item)
	}

	scanStop = time.Now().UnixNano()

	foreignObjectStart = time.Now().UnixNano()

	// this is where any foreign object loading would appear (if required)

	foreignObjectStop = time.Now().UnixNano()

	return items, nil
}

func genericSelectVideosView(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]any, error) {
	items, err := SelectVideosView(ctx, db, columns, orderBy, limit, offset, wheres...)
	if err != nil {
		return nil, err
	}

	genericItems := make([]any, 0)
	for _, item := range items {
		genericItems = append(genericItems, item)
	}

	return genericItems, nil
}
