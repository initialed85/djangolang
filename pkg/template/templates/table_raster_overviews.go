package templates

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

var RasterOverviewViewTable = "raster_overviews"
var RasterOverviewViewTableOTableCatalogColumn = "o_table_catalog"
var RasterOverviewViewTableOTableSchemaColumn = "o_table_schema"
var RasterOverviewViewTableOTableNameColumn = "o_table_name"
var RasterOverviewViewTableORasterColumnColumn = "o_raster_column"
var RasterOverviewViewTableRTableCatalogColumn = "r_table_catalog"
var RasterOverviewViewTableRTableSchemaColumn = "r_table_schema"
var RasterOverviewViewTableRTableNameColumn = "r_table_name"
var RasterOverviewViewTableRRasterColumnColumn = "r_raster_column"
var RasterOverviewViewTableOverviewFactorColumn = "overview_factor"
var RasterOverviewViewColumns = []string{"o_table_catalog", "o_table_schema", "o_table_name", "o_raster_column", "r_table_catalog", "r_table_schema", "r_table_name", "r_raster_column", "overview_factor"}
var RasterOverviewViewTransformedColumns = []string{"o_table_catalog", "o_table_schema", "o_table_name", "o_raster_column", "r_table_catalog", "r_table_schema", "r_table_name", "r_raster_column", "overview_factor"}

type RasterOverviewView struct {
	OTableCatalog  *string `json:"o_table_catalog" db:"o_table_catalog"`
	OTableSchema   *string `json:"o_table_schema" db:"o_table_schema"`
	OTableName     *string `json:"o_table_name" db:"o_table_name"`
	ORasterColumn  *string `json:"o_raster_column" db:"o_raster_column"`
	RTableCatalog  *string `json:"r_table_catalog" db:"r_table_catalog"`
	RTableSchema   *string `json:"r_table_schema" db:"r_table_schema"`
	RTableName     *string `json:"r_table_name" db:"r_table_name"`
	RRasterColumn  *string `json:"r_raster_column" db:"r_raster_column"`
	OverviewFactor *int64  `json:"overview_factor" db:"overview_factor"`
}

func SelectRasterOverviewsView(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]*RasterOverviewView, error) {
	mu.RLock()
	debug := actualDebug
	mu.RUnlock()

	key := "RasterOverviewView"

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
		"SELECT %v FROM raster_overviews%v",
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

	items := make([]*RasterOverviewView, 0)
	for rows.Next() {
		rowCount++

		var item RasterOverviewView
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

func genericSelectRasterOverviewsView(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]any, error) {
	items, err := SelectRasterOverviewsView(ctx, db, columns, orderBy, limit, offset, wheres...)
	if err != nil {
		return nil, err
	}

	genericItems := make([]any, 0)
	for _, item := range items {
		genericItems = append(genericItems, item)
	}

	return genericItems, nil
}
