package some_db

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/initialed85/djangolang/pkg/helpers"
	"github.com/initialed85/djangolang/pkg/types"
	"github.com/jmoiron/sqlx"
	"golang.org/x/exp/maps"
)

var EventTable = "event"
var EventTableIDColumn = "id"
var EventTableStartTimestampColumn = "start_timestamp"
var EventTableEndTimestampColumn = "end_timestamp"
var EventTableDurationColumn = "duration"
var EventTableOriginalVideoIDColumn = "original_video_id"
var EventTableThumbnailImageIDColumn = "thumbnail_image_id"
var EventTableProcessedVideoIDColumn = "processed_video_id"
var EventTableSourceCameraIDColumn = "source_camera_id"
var EventTableStatusColumn = "status"
var EventColumns = []string{"id", "start_timestamp", "end_timestamp", "duration", "original_video_id", "thumbnail_image_id", "processed_video_id", "source_camera_id", "status"}
var EventTransformedColumns = []string{"id", "start_timestamp", "end_timestamp", "extract(microseconds FROM duration)::numeric * 1000 AS duration", "original_video_id", "thumbnail_image_id", "processed_video_id", "source_camera_id", "status"}
var EventInsertColumns = []string{"start_timestamp", "end_timestamp", "duration", "original_video_id", "thumbnail_image_id", "processed_video_id", "source_camera_id", "status"}

func SelectEvents(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]*Event, error) {
	mu.RLock()
	debug := actualDebug
	mu.RUnlock()

	key := "Event"

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

	if columns == nil {
		columns = EventTransformedColumns
	}

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
				"selected %v column(s), %v row(s); %.3f seconds to build, %.3f seconds to execute, %.3f seconds to scan, %.3f seconds to load foreign objects; sql:\n%v\n\n",
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

	queryCtx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()

	tx, err := db.BeginTxx(queryCtx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	rows, err := tx.QueryxContext(
		queryCtx,
		sql,
	)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

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
		b, err := json.Marshal(id)
		if err != nil {
			return nil, err
		}

		s := strings.ReplaceAll(string(b), "\"", "'")

		idsForOriginalVideoID = append(idsForOriginalVideoID, s)
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
		b, err := json.Marshal(id)
		if err != nil {
			return nil, err
		}

		s := strings.ReplaceAll(string(b), "\"", "'")

		idsForThumbnailImageID = append(idsForThumbnailImageID, s)
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
		b, err := json.Marshal(id)
		if err != nil {
			return nil, err
		}

		s := strings.ReplaceAll(string(b), "\"", "'")

		idsForProcessedVideoID = append(idsForProcessedVideoID, s)
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
		b, err := json.Marshal(id)
		if err != nil {
			return nil, err
		}

		s := strings.ReplaceAll(string(b), "\"", "'")

		idsForSourceCameraID = append(idsForSourceCameraID, s)
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

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return items, nil
}

func genericSelectEvents(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]types.DjangolangObject, error) {
	items, err := SelectEvents(ctx, db, columns, orderBy, limit, offset, wheres...)
	if err != nil {
		return nil, err
	}

	genericItems := make([]types.DjangolangObject, 0)
	for _, item := range items {
		genericItems = append(genericItems, item)
	}

	return genericItems, nil
}

func DeserializeEvent(b []byte) (types.DjangolangObject, error) {
	var object Event

	err := json.Unmarshal(b, &object)
	if err != nil {
		return nil, err
	}

	return &object, nil
}

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

func (e *Event) Insert(ctx context.Context, db *sqlx.DB, columns ...string) error {
	if len(columns) > 1 {
		return fmt.Errorf("assertion failed: 'columns' variadic argument(s) must be missing or singular; got %v", len(columns))
	}

	if len(columns) == 0 {
		columns = EventInsertColumns
	}

	mu.RLock()
	debug := actualDebug
	mu.RUnlock()

	var buildStart int64
	var buildStop int64
	var execStart int64
	var execStop int64

	var sql string
	var rowCount int64

	defer func() {
		if !debug {
			return
		}

		buildDuration := 0.0
		execDuration := 0.0

		if buildStop > 0 {
			buildDuration = float64(buildStop-buildStart) * 1e-9
		}

		if execStop > 0 {
			execDuration = float64(execStop-execStart) * 1e-9
		}

		if debug {
			logger.Printf(
				"inserted %v row(s); %.3f seconds to build, %.3f seconds to execute; sql:\n%v\n\n",
				rowCount, buildDuration, execDuration, sql,
			)
		}
	}()

	buildStart = time.Now().UnixNano()

	buildStop = time.Now().UnixNano()

	execStart = time.Now().UnixNano()

	names := make([]string, 0)
	for _, column := range columns {
		names = append(names, fmt.Sprintf(":%v", column))
	}

	sql = fmt.Sprintf(
		"INSERT INTO event (%v) VALUES (%v) RETURNING %v",
		strings.Join(columns, ", "),
		strings.Join(names, ", "),
		strings.Join(EventColumns, ", "),
	)

	queryCtx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()

	tx, err := db.BeginTxx(queryCtx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	result, err := tx.NamedQuery(
		sql,
		e,
	)
	if err != nil {
		return err
	}
	defer func() {
		_ = result.Close()
	}()

	ok := result.Next()
	if !ok {
		return fmt.Errorf("insert w/ returning unexpectedly returned nothing")
	}

	err = result.StructScan(e)
	if err != nil {
		return err
	}

	rowCount = 1

	execStop = time.Now().UnixNano()

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

func genericInsertEvent(ctx context.Context, db *sqlx.DB, object types.DjangolangObject, columns ...string) (types.DjangolangObject, error) {
	if object == nil {
		return nil, fmt.Errorf("object given for insertion was unexpectedly nil")
	}

	err := object.Insert(ctx, db, columns...)
	if err != nil {
		return nil, err
	}

	return object, nil
}

func (e *Event) GetPrimaryKey() (any, error) {
	return e.ID, nil
}

func (e *Event) SetPrimaryKey(value any) error {
	e.ID = value.(int64)

	return nil
}

func (e *Event) Update(ctx context.Context, db *sqlx.DB, columns ...string) error {
	if len(columns) > 1 {
		return fmt.Errorf("assertion failed: 'columns' variadic argument(s) must be missing or singular; got %v", len(columns))
	}

	if len(columns) == 0 {
		columns = EventInsertColumns
	}

	mu.RLock()
	debug := actualDebug
	mu.RUnlock()

	var buildStart int64
	var buildStop int64
	var execStart int64
	var execStop int64

	var sql string
	var rowCount int64

	defer func() {
		if !debug {
			return
		}

		buildDuration := 0.0
		execDuration := 0.0

		if buildStop > 0 {
			buildDuration = float64(buildStop-buildStart) * 1e-9
		}

		if execStop > 0 {
			execDuration = float64(execStop-execStart) * 1e-9
		}

		if debug {
			logger.Printf(
				"updated %v row(s); %.3f seconds to build, %.3f seconds to execute; sql:\n%v\n\n",
				rowCount, buildDuration, execDuration, sql,
			)
		}
	}()

	buildStart = time.Now().UnixNano()

	buildStop = time.Now().UnixNano()

	execStart = time.Now().UnixNano()

	names := make([]string, 0)
	for _, column := range columns {
		names = append(names, fmt.Sprintf(":%v", column))
	}

	sql = fmt.Sprintf(
		"UPDATE camera SET (%v) = (%v) WHERE id = %v RETURNING %v",
		strings.Join(columns, ", "),
		strings.Join(names, ", "),
		e.ID,
		strings.Join(EventColumns, ", "),
	)

	queryCtx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()

	tx, err := db.BeginTxx(queryCtx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	result, err := tx.NamedQuery(
		sql,
		e,
	)
	if err != nil {
		return err
	}
	defer func() {
		_ = result.Close()
	}()

	ok := result.Next()
	if !ok {
		return fmt.Errorf("update w/ returning unexpectedly returned nothing")
	}

	err = result.StructScan(e)
	if err != nil {
		return err
	}

	rowCount = 1

	execStop = time.Now().UnixNano()

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

func genericUpdateEvent(ctx context.Context, db *sqlx.DB, object types.DjangolangObject, columns ...string) (types.DjangolangObject, error) {
	if object == nil {
		return nil, fmt.Errorf("object given for update was unexpectedly nil")
	}

	err := object.Update(ctx, db, columns...)
	if err != nil {
		return nil, err
	}

	return object, nil
}

func (e *Event) Delete(ctx context.Context, db *sqlx.DB) error {
	mu.RLock()
	debug := actualDebug
	mu.RUnlock()

	var buildStart int64
	var buildStop int64
	var execStart int64
	var execStop int64

	var sql string
	var rowCount int64

	defer func() {
		if !debug {
			return
		}

		buildDuration := 0.0
		execDuration := 0.0

		if buildStop > 0 {
			buildDuration = float64(buildStop-buildStart) * 1e-9
		}

		if execStop > 0 {
			execDuration = float64(execStop-execStart) * 1e-9
		}

		if debug {
			logger.Printf(
				"deleted %v row(s); %.3f seconds to build, %.3f seconds to execute; sql:\n%v\n\n",
				rowCount, buildDuration, execDuration, sql,
			)
		}
	}()

	buildStart = time.Now().UnixNano()

	buildStop = time.Now().UnixNano()

	execStart = time.Now().UnixNano()

	sql = fmt.Sprintf(
		"DELETE FROM event WHERE id = %v",
		e.ID,
	)

	queryCtx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()

	tx, err := db.BeginTxx(queryCtx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	result, err := tx.ExecContext(
		queryCtx,
		sql,
	)
	if err != nil {
		return err
	}

	rowCount, err = result.RowsAffected()
	if err != nil {
		return err
	}

	execStop = time.Now().UnixNano()

	if rowCount != 1 {
		return fmt.Errorf("expected exactly 1 affected row after deleting %#+v; got %v", e, rowCount)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

func genericDeleteEvent(ctx context.Context, db *sqlx.DB, object types.DjangolangObject) error {
	if object == nil {
		return fmt.Errorf("object given for deletion was unexpectedly nil")
	}

	err := object.Delete(ctx, db)
	if err != nil {
		return err
	}

	return nil
}
