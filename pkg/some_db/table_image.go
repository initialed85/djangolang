package some_db

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/initialed85/djangolang/pkg/helpers"
	"github.com/jmoiron/sqlx"
	"golang.org/x/exp/maps"
)

var ImageTable = "image"
var ImageTableIDColumn = "id"
var ImageTableTimestampColumn = "timestamp"
var ImageTableSizeColumn = "size"
var ImageTableFilePathColumn = "file_path"
var ImageTableCameraIDColumn = "camera_id"
var ImageTableEventIDColumn = "event_id"
var ImageColumns = []string{"id", "timestamp", "size", "file_path", "camera_id", "event_id"}
var ImageTransformedColumns = []string{"id", "timestamp", "size", "file_path", "camera_id", "event_id"}
var ImageInsertColumns = []string{"timestamp", "size", "file_path", "camera_id", "event_id"}

func SelectImages(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]*Image, error) {
	mu.RLock()
	debug := actualDebug
	mu.RUnlock()

	key := "Image"

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
		columns = ImageTransformedColumns
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
		"SELECT %v FROM image%v",
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

	cameraByCameraID := make(map[int64]*Camera, 0)
	eventByEventID := make(map[int64]*Event, 0)

	items := make([]*Image, 0)
	for rows.Next() {
		rowCount++

		var item Image
		err = rows.StructScan(&item)
		if err != nil {
			return nil, err
		}

		// this is where any post-scan processing would appear (if required)

		cameraByCameraID[item.CameraID] = nil
		eventByEventID[helpers.Deref(item.EventID)] = nil

		items = append(items, &item)
	}

	scanStop = time.Now().UnixNano()

	foreignObjectStart = time.Now().UnixNano()

	idsForCameraID := make([]string, 0)
	for _, id := range maps.Keys(cameraByCameraID) {
		b, err := json.Marshal(id)
		if err != nil {
			return nil, err
		}

		s := strings.ReplaceAll(string(b), "\"", "'")

		idsForCameraID = append(idsForCameraID, s)
	}

	if len(idsForCameraID) > 0 {
		rowsForCameraID, err := SelectCameras(
			ctx,
			db,
			CameraTransformedColumns,
			nil,
			nil,
			nil,
			fmt.Sprintf("id IN (%v)", strings.Join(idsForCameraID, ", ")),
		)
		if err != nil {
			return nil, err
		}

		for _, row := range rowsForCameraID {
			cameraByCameraID[row.ID] = row
		}

		for _, item := range items {
			item.CameraIDObject = cameraByCameraID[item.CameraID]
		}
	}

	idsForEventID := make([]string, 0)
	for _, id := range maps.Keys(eventByEventID) {
		b, err := json.Marshal(id)
		if err != nil {
			return nil, err
		}

		s := strings.ReplaceAll(string(b), "\"", "'")

		idsForEventID = append(idsForEventID, s)
	}

	if len(idsForEventID) > 0 {
		rowsForEventID, err := SelectEvents(
			ctx,
			db,
			EventTransformedColumns,
			nil,
			nil,
			nil,
			fmt.Sprintf("id IN (%v)", strings.Join(idsForEventID, ", ")),
		)
		if err != nil {
			return nil, err
		}

		for _, row := range rowsForEventID {
			eventByEventID[row.ID] = row
		}

		for _, item := range items {
			item.EventIDObject = eventByEventID[helpers.Deref(item.EventID)]
		}
	}

	foreignObjectStop = time.Now().UnixNano()

	return items, nil
}

func genericSelectImages(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]DjangolangObject, error) {
	items, err := SelectImages(ctx, db, columns, orderBy, limit, offset, wheres...)
	if err != nil {
		return nil, err
	}

	genericItems := make([]DjangolangObject, 0)
	for _, item := range items {
		genericItems = append(genericItems, item)
	}

	return genericItems, nil
}

func DeserializeImage(b []byte) (DjangolangObject, error) {
	var object Image

	err := json.Unmarshal(b, &object)
	if err != nil {
		return nil, err
	}

	return &object, nil
}

type Image struct {
	ID             int64     `json:"id" db:"id"`
	Timestamp      time.Time `json:"timestamp" db:"timestamp"`
	Size           float64   `json:"size" db:"size"`
	FilePath       string    `json:"file_path" db:"file_path"`
	CameraID       int64     `json:"camera_id" db:"camera_id"`
	CameraIDObject *Camera   `json:"camera_id_object,omitempty"`
	EventID        *int64    `json:"event_id" db:"event_id"`
	EventIDObject  *Event    `json:"event_id_object,omitempty"`
}

func (i *Image) Insert(ctx context.Context, db *sqlx.DB, columns ...string) error {
	if len(columns) > 1 {
		return fmt.Errorf("assertion failed: 'columns' variadic argument(s) must be missing or singular; got %v", len(columns))
	}

	if len(columns) == 0 {
		columns = ImageInsertColumns
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
				"inserted %v rows; %.3f seconds to build, %.3f seconds to execute; sql:\n%v",
				rowCount, buildDuration, execDuration, sql,
			)
		}
	}()

	buildStart = time.Now().UnixNano()

	buildStop = time.Now().UnixNano()

	execStart = time.Now().UnixNano()

	insertCtx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()

	names := make([]string, 0)
	for _, column := range columns {
		names = append(names, fmt.Sprintf(":%v", column))
	}

	result, err := db.NamedQueryContext(
		insertCtx,
		fmt.Sprintf(
			"INSERT INTO image (%v) VALUES (%v) RETURNING %v",
			strings.Join(columns, ", "),
			strings.Join(names, ", "),
			strings.Join(ImageColumns, ", "),
		),
		i,
	)
	if err != nil {
		return err
	}

	_ = result.Next()

	err = result.StructScan(i)
	if err != nil {
		return err
	}

	execStop = time.Now().UnixNano()

	return nil
}

func genericInsertImage(ctx context.Context, db *sqlx.DB, object DjangolangObject, columns ...string) (DjangolangObject, error) {
	if object == nil {
		return nil, fmt.Errorf("object given for insertion was unexpectedly nil")
	}

	err := object.Insert(ctx, db, columns...)
	if err != nil {
		return nil, err
	}

	return object, nil
}

func (i *Image) Delete(ctx context.Context, db *sqlx.DB) error {
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
				"deleted %v rows; %.3f seconds to build, %.3f seconds to execute; sql:\n%v",
				rowCount, buildDuration, execDuration, sql,
			)
		}
	}()

	buildStart = time.Now().UnixNano()

	buildStop = time.Now().UnixNano()

	execStart = time.Now().UnixNano()

	deleteCtx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()

	result, err := db.ExecContext(
		deleteCtx,
		fmt.Sprintf(
			"DELETE FROM image WHERE id = %v",
			i.ID,
		),
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	execStop = time.Now().UnixNano()

	if rowsAffected != 1 {
		return fmt.Errorf("expected exactly 1 affected row after deleting %#+v; got %v", i, rowsAffected)
	}

	return nil
}

func genericDeleteImage(ctx context.Context, db *sqlx.DB, object DjangolangObject) error {
	if object == nil {
		return fmt.Errorf("object given for deletion was unexpectedly nil")
	}

	err := object.Delete(ctx, db)
	if err != nil {
		return err
	}

	return nil
}
