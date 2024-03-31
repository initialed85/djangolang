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

var AggregatedDetectionTable = "aggregated_detection"
var AggregatedDetectionTableIDColumn = "id"
var AggregatedDetectionTableStartTimestampColumn = "start_timestamp"
var AggregatedDetectionTableEndTimestampColumn = "end_timestamp"
var AggregatedDetectionTableClassIDColumn = "class_id"
var AggregatedDetectionTableClassNameColumn = "class_name"
var AggregatedDetectionTableScoreColumn = "score"
var AggregatedDetectionTableCountColumn = "count"
var AggregatedDetectionTableWeightedScoreColumn = "weighted_score"
var AggregatedDetectionTableEventIDColumn = "event_id"
var AggregatedDetectionColumns = []string{"id", "start_timestamp", "end_timestamp", "class_id", "class_name", "score", "count", "weighted_score", "event_id"}
var AggregatedDetectionTransformedColumns = []string{"id", "start_timestamp", "end_timestamp", "class_id", "class_name", "score", "count", "weighted_score", "event_id"}
var AggregatedDetectionInsertColumns = []string{"start_timestamp", "end_timestamp", "class_id", "class_name", "score", "count", "weighted_score", "event_id"}

func SelectAggregatedDetections(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]*AggregatedDetection, error) {
	mu.RLock()
	debug := actualDebug
	mu.RUnlock()

	key := "AggregatedDetection"

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
		columns = AggregatedDetectionTransformedColumns
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

func genericSelectAggregatedDetections(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]DjangolangObject, error) {
	items, err := SelectAggregatedDetections(ctx, db, columns, orderBy, limit, offset, wheres...)
	if err != nil {
		return nil, err
	}

	genericItems := make([]DjangolangObject, 0)
	for _, item := range items {
		genericItems = append(genericItems, item)
	}

	return genericItems, nil
}

func DeserializeAggregatedDetection(b []byte) (DjangolangObject, error) {
	var object AggregatedDetection

	err := json.Unmarshal(b, &object)
	if err != nil {
		return nil, err
	}

	return &object, nil
}

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

func (a *AggregatedDetection) Insert(ctx context.Context, db *sqlx.DB, columns ...string) error {
	if len(columns) > 1 {
		return fmt.Errorf("assertion failed: 'columns' variadic argument(s) must be missing or singular; got %v", len(columns))
	}

	if len(columns) == 0 {
		columns = AggregatedDetectionInsertColumns
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
			"INSERT INTO aggregated_detection (%v) VALUES (%v) RETURNING %v",
			strings.Join(columns, ", "),
			strings.Join(names, ", "),
			strings.Join(AggregatedDetectionColumns, ", "),
		),
		a,
	)
	if err != nil {
		return err
	}

	_ = result.Next()

	err = result.StructScan(a)
	if err != nil {
		return err
	}

	execStop = time.Now().UnixNano()

	return nil
}

func genericInsertAggregatedDetection(ctx context.Context, db *sqlx.DB, object DjangolangObject, columns ...string) (DjangolangObject, error) {
	if object == nil {
		return nil, fmt.Errorf("object given for insertion was unexpectedly nil")
	}

	err := object.Insert(ctx, db, columns...)
	if err != nil {
		return nil, err
	}

	return object, nil
}

func (a *AggregatedDetection) Delete(ctx context.Context, db *sqlx.DB) error {
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
			"DELETE FROM aggregated_detection WHERE id = %v",
			a.ID,
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
		return fmt.Errorf("expected exactly 1 affected row after deleting %#+v; got %v", a, rowsAffected)
	}

	return nil
}

func genericDeleteAggregatedDetection(ctx context.Context, db *sqlx.DB, object DjangolangObject) error {
	if object == nil {
		return fmt.Errorf("object given for deletion was unexpectedly nil")
	}

	err := object.Delete(ctx, db)
	if err != nil {
		return err
	}

	return nil
}
