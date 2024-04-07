package some_db

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

var SpatialRefSyTable = "spatial_ref_sys"
var SpatialRefSyTableSridColumn = "srid"
var SpatialRefSyTableAuthNameColumn = "auth_name"
var SpatialRefSyTableAuthSridColumn = "auth_srid"
var SpatialRefSyTableSrtextColumn = "srtext"
var SpatialRefSyTableProj4textColumn = "proj4text"
var SpatialRefSyColumns = []string{"srid", "auth_name", "auth_srid", "srtext", "proj4text"}
var SpatialRefSyTransformedColumns = []string{"srid", "auth_name", "auth_srid", "srtext", "proj4text"}
var SpatialRefSyInsertColumns = []string{"auth_name", "auth_srid", "srtext", "proj4text"}

func SelectSpatialRefSys(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]*SpatialRefSy, error) {
	mu.RLock()
	debug := actualDebug
	mu.RUnlock()

	key := "SpatialRefSy"

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
		columns = SpatialRefSyTransformedColumns
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
		"SELECT %v FROM spatial_ref_sys%v",
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

	// this is where any foreign object ID map declarations would appear (if required)

	items := make([]*SpatialRefSy, 0)
	for rows.Next() {
		rowCount++

		var item SpatialRefSy
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

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return items, nil
}

func genericSelectSpatialRefSys(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]DjangolangObject, error) {
	items, err := SelectSpatialRefSys(ctx, db, columns, orderBy, limit, offset, wheres...)
	if err != nil {
		return nil, err
	}

	genericItems := make([]DjangolangObject, 0)
	for _, item := range items {
		genericItems = append(genericItems, item)
	}

	return genericItems, nil
}

func DeserializeSpatialRefSy(b []byte) (DjangolangObject, error) {
	var object SpatialRefSy

	err := json.Unmarshal(b, &object)
	if err != nil {
		return nil, err
	}

	return &object, nil
}

type SpatialRefSy struct {
	Srid      int64   `json:"srid" db:"srid"`
	AuthName  *string `json:"auth_name" db:"auth_name"`
	AuthSrid  *int64  `json:"auth_srid" db:"auth_srid"`
	Srtext    *string `json:"srtext" db:"srtext"`
	Proj4text *string `json:"proj4text" db:"proj4text"`
}

func (s *SpatialRefSy) Insert(ctx context.Context, db *sqlx.DB, columns ...string) error {
	if len(columns) > 1 {
		return fmt.Errorf("assertion failed: 'columns' variadic argument(s) must be missing or singular; got %v", len(columns))
	}

	if len(columns) == 0 {
		columns = SpatialRefSyInsertColumns
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
		"INSERT INTO spatial_ref_sys (%v) VALUES (%v) RETURNING %v",
		strings.Join(columns, ", "),
		strings.Join(names, ", "),
		strings.Join(SpatialRefSyColumns, ", "),
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
		s,
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

	err = result.StructScan(s)
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

func genericInsertSpatialRefSy(ctx context.Context, db *sqlx.DB, object DjangolangObject, columns ...string) (DjangolangObject, error) {
	if object == nil {
		return nil, fmt.Errorf("object given for insertion was unexpectedly nil")
	}

	err := object.Insert(ctx, db, columns...)
	if err != nil {
		return nil, err
	}

	return object, nil
}

func (s *SpatialRefSy) GetPrimaryKey() (any, error) {
	return s.Srid, nil
}

func (s *SpatialRefSy) SetPrimaryKey(value any) error {
	s.Srid = value.(int64)

	return nil
}

func (s *SpatialRefSy) Update(ctx context.Context, db *sqlx.DB, columns ...string) error {
	if len(columns) > 1 {
		return fmt.Errorf("assertion failed: 'columns' variadic argument(s) must be missing or singular; got %v", len(columns))
	}

	if len(columns) == 0 {
		columns = SpatialRefSyInsertColumns
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
		s.Srid,
		strings.Join(SpatialRefSyColumns, ", "),
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
		s,
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

	err = result.StructScan(s)
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

func genericUpdateSpatialRefSy(ctx context.Context, db *sqlx.DB, object DjangolangObject, columns ...string) (DjangolangObject, error) {
	if object == nil {
		return nil, fmt.Errorf("object given for update was unexpectedly nil")
	}

	err := object.Update(ctx, db, columns...)
	if err != nil {
		return nil, err
	}

	return object, nil
}

func (s *SpatialRefSy) Delete(ctx context.Context, db *sqlx.DB) error {
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
		"DELETE FROM spatial_ref_sys WHERE srid = %v",
		s.Srid,
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
		return fmt.Errorf("expected exactly 1 affected row after deleting %#+v; got %v", s, rowCount)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

func genericDeleteSpatialRefSy(ctx context.Context, db *sqlx.DB, object DjangolangObject) error {
	if object == nil {
		return fmt.Errorf("object given for deletion was unexpectedly nil")
	}

	err := object.Delete(ctx, db)
	if err != nil {
		return err
	}

	return nil
}
