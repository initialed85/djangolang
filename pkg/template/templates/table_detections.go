package templates

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

var DetectionViewTable = "detections"
var DetectionViewIDColumn = "id"
var DetectionViewTimestampColumn = "timestamp"
var DetectionViewClassIDColumn = "class_id"
var DetectionViewClassNameColumn = "class_name"
var DetectionViewScoreColumn = "score"
var DetectionViewCentroidColumn = "centroid"
var DetectionViewBoundingBoxColumn = "bounding_box"
var DetectionViewCameraIDColumn = "camera_id"
var DetectionViewEventIDColumn = "event_id"
var DetectionViewObjectIDColumn = "object_id"
var DetectionViewColourColumn = "colour"
var DetectionViewColumns = []string{"id", "timestamp", "class_id", "class_name", "score", "centroid", "bounding_box", "camera_id", "event_id", "object_id", "colour"}
var DetectionViewTransformedColumns = []string{"id", "timestamp", "class_id", "class_name", "score", "ST_AsGeoJSON(centroid::geometry)::jsonb AS centroid", "ST_AsGeoJSON(bounding_box::geometry)::jsonb AS bounding_box", "camera_id", "event_id", "object_id", "colour"}

type DetectionView struct {
	ID          *int64     `json:"id" db:"id"`
	Timestamp   *time.Time `json:"timestamp" db:"timestamp"`
	ClassID     *int64     `json:"class_id" db:"class_id"`
	ClassName   *string    `json:"class_name" db:"class_name"`
	Score       *float64   `json:"score" db:"score"`
	Centroid    any        `json:"centroid" db:"centroid"`
	BoundingBox any        `json:"bounding_box" db:"bounding_box"`
	CameraID    *int64     `json:"camera_id" db:"camera_id"`
	EventID     *int64     `json:"event_id" db:"event_id"`
	ObjectID    *int64     `json:"object_id" db:"object_id"`
	Colour      any        `json:"colour" db:"colour"`
}

func SelectDetectionsView(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]*DetectionView, error) {
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
		"SELECT %v FROM detections%v",
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

	items := make([]*DetectionView, 0)
	for rows.Next() {
		rowCount++

		var item DetectionView
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

		// this is where foreign object ID map populations would appear (if required)

		items = append(items, &item)
	}

	scanStop = time.Now().UnixNano()

	foreignObjectStart = time.Now().UnixNano()

	// this is where any foreign object loading would appear (if required)

	foreignObjectStop = time.Now().UnixNano()

	return items, nil
}

func genericSelectDetectionsView(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]any, error) {
	items, err := SelectDetectionsView(ctx, db, columns, orderBy, limit, offset, wheres...)
	if err != nil {
		return nil, err
	}

	genericItems := make([]any, 0)
	for _, item := range items {
		genericItems = append(genericItems, item)
	}

	return genericItems, nil
}
