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

var AggregatedDetectionTable = "aggregated_detection"
var AggregatedDetectionIDColumn = "id"
var AggregatedDetectionStartTimestampColumn = "start_timestamp"
var AggregatedDetectionEndTimestampColumn = "end_timestamp"
var AggregatedDetectionClassIDColumn = "class_id"
var AggregatedDetectionClassNameColumn = "class_name"
var AggregatedDetectionScoreColumn = "score"
var AggregatedDetectionCountColumn = "count"
var AggregatedDetectionWeightedScoreColumn = "weighted_score"
var AggregatedDetectionEventIDColumn = "event_id"
var AggregatedDetectionColumns = []string{"id", "start_timestamp", "end_timestamp", "class_id", "class_name", "score", "count", "weighted_score", "event_id"}
var AggregatedDetectionTransformedColumns = []string{"id", "start_timestamp", "end_timestamp", "class_id", "class_name", "score", "count", "weighted_score", "event_id"}

type AggregatedDetection struct {
	ID             int64     `json:"id" db:"id"`
	StartTimestamp time.Time `json:"start_timestamp" db:"start_timestamp"`
	EndTimestamp   time.Time `json:"end_timestamp" db:"end_timestamp"`
	ClassID        int64     `json:"class_id" db:"class_id"`
	ClassName      string    `json:"class_name" db:"class_name"`
	Score          float64   `json:"score" db:"score"`
	Count          int64     `json:"count" db:"count"`
	WeightedScore  float64   `json:"weighted_score" db:"weighted_score"`
	EventID        *int64    `json:"event_id" db:"event_id"`
	EventIDObject  *Event    `json:"event_id_object,omitempty"`
}

func SelectAggregatedDetections(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]*AggregatedDetection, error) {
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
		"SELECT %v FROM aggregated_detection%v",
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

	eventByEventID := make(map[int64]*Event, 0)

	items := make([]*AggregatedDetection, 0)
	for rows.Next() {
		rowCount++

		var item AggregatedDetection
		err = rows.StructScan(&item)
		if err != nil {
			return nil, err
		}

		// this is where any post-scan processing would appear (if required)

		eventByEventID[helpers.Deref(item.EventID)] = nil

		items = append(items, &item)
	}

	scanStop = time.Now().UnixNano()

	foreignObjectStart = time.Now().UnixNano()

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

	foreignObjectStop = time.Now().UnixNano()

	return items, nil
}

func genericSelectAggregatedDetections(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]any, error) {
	items, err := SelectAggregatedDetections(ctx, db, columns, orderBy, limit, offset, wheres...)
	if err != nil {
		return nil, err
	}

	genericItems := make([]any, 0)
	for _, item := range items {
		genericItems = append(genericItems, item)
	}

	return genericItems, nil
}
