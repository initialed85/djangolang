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

var DetectionTable = "detection"
var DetectionTableIDColumn = "id"
var DetectionTableTimestampColumn = "timestamp"
var DetectionTableClassIDColumn = "class_id"
var DetectionTableClassNameColumn = "class_name"
var DetectionTableScoreColumn = "score"
var DetectionTableCentroidColumn = "centroid"
var DetectionTableBoundingBoxColumn = "bounding_box"
var DetectionTableCameraIDColumn = "camera_id"
var DetectionTableEventIDColumn = "event_id"
var DetectionTableObjectIDColumn = "object_id"
var DetectionTableColourColumn = "colour"
var DetectionColumns = []string{"id", "timestamp", "class_id", "class_name", "score", "centroid", "bounding_box", "camera_id", "event_id", "object_id", "colour"}
var DetectionTransformedColumns = []string{"id", "timestamp", "class_id", "class_name", "score", "ST_AsGeoJSON(centroid::geometry)::jsonb AS centroid", "ST_AsGeoJSON(bounding_box::geometry)::jsonb AS bounding_box", "camera_id", "event_id", "object_id", "colour"}
var DetectionInsertColumns = []string{"timestamp", "class_id", "class_name", "score", "centroid", "bounding_box", "camera_id", "event_id", "object_id", "colour"}

func SelectDetections(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]*Detection, error) {
	mu.RLock()
	debug := actualDebug
	mu.RUnlock()

	key := "Detection"

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
		columns = DetectionTransformedColumns
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
		"SELECT %v FROM detection%v",
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

	cameraByCameraID := make(map[int64]*Camera, 0)
	eventByEventID := make(map[int64]*Event, 0)
	objectByObjectID := make(map[int64]*Object, 0)

	items := make([]*Detection, 0)
	for rows.Next() {
		rowCount++

		var item Detection
		err = rows.StructScan(&item)
		if err != nil {
			return nil, err
		}

		var temp1Centroid any
		var temp2Centroid []any

		if item.Centroid != nil {
			err = json.Unmarshal(item.Centroid.([]byte), &temp1Centroid)
			if err != nil {
				err = json.Unmarshal(item.Centroid.([]byte), &temp2Centroid)
				if err != nil {
					item.Centroid = nil
				} else {
					item.Centroid = &temp2Centroid
				}
			} else {
				item.Centroid = &temp1Centroid
			}
		}

		item.Centroid = &temp1Centroid

		var temp1BoundingBox any
		var temp2BoundingBox []any

		if item.BoundingBox != nil {
			err = json.Unmarshal(item.BoundingBox.([]byte), &temp1BoundingBox)
			if err != nil {
				err = json.Unmarshal(item.BoundingBox.([]byte), &temp2BoundingBox)
				if err != nil {
					item.BoundingBox = nil
				} else {
					item.BoundingBox = &temp2BoundingBox
				}
			} else {
				item.BoundingBox = &temp1BoundingBox
			}
		}

		item.BoundingBox = &temp1BoundingBox

		var temp1Colour any
		var temp2Colour []any

		if item.Colour != nil {
			err = json.Unmarshal(item.Colour.([]byte), &temp1Colour)
			if err != nil {
				err = json.Unmarshal(item.Colour.([]byte), &temp2Colour)
				if err != nil {
					item.Colour = nil
				} else {
					item.Colour = &temp2Colour
				}
			} else {
				item.Colour = &temp1Colour
			}
		}

		item.Colour = &temp1Colour

		cameraByCameraID[item.CameraID] = nil
		eventByEventID[helpers.Deref(item.EventID)] = nil
		objectByObjectID[helpers.Deref(item.ObjectID)] = nil

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

	idsForObjectID := make([]string, 0)
	for _, id := range maps.Keys(objectByObjectID) {
		b, err := json.Marshal(id)
		if err != nil {
			return nil, err
		}

		s := strings.ReplaceAll(string(b), "\"", "'")

		idsForObjectID = append(idsForObjectID, s)
	}

	if len(idsForObjectID) > 0 {
		rowsForObjectID, err := SelectObjects(
			ctx,
			db,
			ObjectTransformedColumns,
			nil,
			nil,
			nil,
			fmt.Sprintf("id IN (%v)", strings.Join(idsForObjectID, ", ")),
		)
		if err != nil {
			return nil, err
		}

		for _, row := range rowsForObjectID {
			objectByObjectID[row.ID] = row
		}

		for _, item := range items {
			item.ObjectIDObject = objectByObjectID[helpers.Deref(item.ObjectID)]
		}
	}

	foreignObjectStop = time.Now().UnixNano()

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return items, nil
}

func genericSelectDetections(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]DjangolangObject, error) {
	items, err := SelectDetections(ctx, db, columns, orderBy, limit, offset, wheres...)
	if err != nil {
		return nil, err
	}

	genericItems := make([]DjangolangObject, 0)
	for _, item := range items {
		genericItems = append(genericItems, item)
	}

	return genericItems, nil
}

func DeserializeDetection(b []byte) (DjangolangObject, error) {
	var object Detection

	err := json.Unmarshal(b, &object)
	if err != nil {
		return nil, err
	}

	return &object, nil
}

type Detection struct {
	ID             int64     `json:"id" db:"id"`
	Timestamp      time.Time `json:"timestamp" db:"timestamp"`
	ClassID        int64     `json:"class_id" db:"class_id"`
	ClassName      string    `json:"class_name" db:"class_name"`
	Score          float64   `json:"score" db:"score"`
	Centroid       any       `json:"centroid" db:"centroid"`
	BoundingBox    any       `json:"bounding_box" db:"bounding_box"`
	CameraID       int64     `json:"camera_id" db:"camera_id"`
	CameraIDObject *Camera   `json:"camera_id_object,omitempty"`
	EventID        *int64    `json:"event_id" db:"event_id"`
	EventIDObject  *Event    `json:"event_id_object,omitempty"`
	ObjectID       *int64    `json:"object_id" db:"object_id"`
	ObjectIDObject *Object   `json:"object_id_object,omitempty"`
	Colour         any       `json:"colour" db:"colour"`
}

func (d *Detection) Insert(ctx context.Context, db *sqlx.DB, columns ...string) error {
	if len(columns) > 1 {
		return fmt.Errorf("assertion failed: 'columns' variadic argument(s) must be missing or singular; got %v", len(columns))
	}

	if len(columns) == 0 {
		columns = DetectionInsertColumns
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
		"INSERT INTO detection (%v) VALUES (%v) RETURNING %v",
		strings.Join(columns, ", "),
		strings.Join(names, ", "),
		strings.Join(DetectionColumns, ", "),
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
		d,
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

	err = result.StructScan(d)
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

func genericInsertDetection(ctx context.Context, db *sqlx.DB, object DjangolangObject, columns ...string) (DjangolangObject, error) {
	if object == nil {
		return nil, fmt.Errorf("object given for insertion was unexpectedly nil")
	}

	err := object.Insert(ctx, db, columns...)
	if err != nil {
		return nil, err
	}

	return object, nil
}

func (d *Detection) GetPrimaryKey() (any, error) {
	return d.ID, nil
}

func (d *Detection) SetPrimaryKey(value any) error {
	d.ID = value.(int64)

	return nil
}

func (d *Detection) Update(ctx context.Context, db *sqlx.DB, columns ...string) error {
	if len(columns) > 1 {
		return fmt.Errorf("assertion failed: 'columns' variadic argument(s) must be missing or singular; got %v", len(columns))
	}

	if len(columns) == 0 {
		columns = DetectionInsertColumns
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
		d.ID,
		strings.Join(DetectionColumns, ", "),
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
		d,
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

	err = result.StructScan(d)
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

func genericUpdateDetection(ctx context.Context, db *sqlx.DB, object DjangolangObject, columns ...string) (DjangolangObject, error) {
	if object == nil {
		return nil, fmt.Errorf("object given for update was unexpectedly nil")
	}

	err := object.Update(ctx, db, columns...)
	if err != nil {
		return nil, err
	}

	return object, nil
}

func (d *Detection) Delete(ctx context.Context, db *sqlx.DB) error {
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
		"DELETE FROM detection WHERE id = %v",
		d.ID,
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
		return fmt.Errorf("expected exactly 1 affected row after deleting %#+v; got %v", d, rowCount)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

func genericDeleteDetection(ctx context.Context, db *sqlx.DB, object DjangolangObject) error {
	if object == nil {
		return fmt.Errorf("object given for deletion was unexpectedly nil")
	}

	err := object.Delete(ctx, db)
	if err != nil {
		return err
	}

	return nil
}
