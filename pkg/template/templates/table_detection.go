package templates

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
var DetectionIDColumn = "id"
var DetectionTimestampColumn = "timestamp"
var DetectionClassIDColumn = "class_id"
var DetectionClassNameColumn = "class_name"
var DetectionScoreColumn = "score"
var DetectionCentroidColumn = "centroid"
var DetectionBoundingBoxColumn = "bounding_box"
var DetectionCameraIDColumn = "camera_id"
var DetectionEventIDColumn = "event_id"
var DetectionObjectIDColumn = "object_id"
var DetectionColourColumn = "colour"
var DetectionColumns = []string{"id", "timestamp", "class_id", "class_name", "score", "centroid", "bounding_box", "camera_id", "event_id", "object_id", "colour"}
var DetectionTransformedColumns = []string{"id", "timestamp", "class_id", "class_name", "score", "ST_AsGeoJSON(centroid::geometry)::jsonb AS centroid", "ST_AsGeoJSON(bounding_box::geometry)::jsonb AS bounding_box", "camera_id", "event_id", "object_id", "colour"}

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

func SelectDetections(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]*Detection, error) {
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
		idsForCameraID = append(idsForCameraID, fmt.Sprintf("%v", id))
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
		idsForEventID = append(idsForEventID, fmt.Sprintf("%v", id))
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
		idsForObjectID = append(idsForObjectID, fmt.Sprintf("%v", id))
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

	return items, nil
}

func genericSelectDetections(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]any, error) {
	items, err := SelectDetections(ctx, db, columns, orderBy, limit, offset, wheres...)
	if err != nil {
		return nil, err
	}

	genericItems := make([]any, 0)
	for _, item := range items {
		genericItems = append(genericItems, item)
	}

	return genericItems, nil
}
