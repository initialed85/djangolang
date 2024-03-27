package templates

import (
	"context"
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

type SpatialRefSy struct {
	Srid      int64   `json:"srid" db:"srid"`
	AuthName  *string `json:"auth_name" db:"auth_name"`
	AuthSrid  *int64  `json:"auth_srid" db:"auth_srid"`
	Srtext    *string `json:"srtext" db:"srtext"`
	Proj4text *string `json:"proj4text" db:"proj4text"`
}

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

	selectCtx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()

	rows, err := db.QueryxContext(selectCtx, sql)
	if err != nil {
		return nil, err
	}

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

	return items, nil
}

func genericSelectSpatialRefSys(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]any, error) {
	items, err := SelectSpatialRefSys(ctx, db, columns, orderBy, limit, offset, wheres...)
	if err != nil {
		return nil, err
	}

	genericItems := make([]any, 0)
	for _, item := range items {
		genericItems = append(genericItems, item)
	}

	return genericItems, nil
}
