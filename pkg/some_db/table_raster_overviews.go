package some_db

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/initialed85/djangolang/pkg/types"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
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
var RasterOverviewViewInsertColumns = []string{"o_table_catalog", "o_table_schema", "o_table_name", "o_raster_column", "r_table_catalog", "r_table_schema", "r_table_name", "r_raster_column", "overview_factor"}

func SelectRasterOverviewsView(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...types.Fragment) ([]*RasterOverviewView, error) {
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

	if columns == nil {
		columns = RasterOverviewViewTransformedColumns
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

	whereSQLs := make([]string, 0)
	whereValues := make([]any, 0)
	for _, fragment := range wheres {
		whereSQLs = append(whereSQLs, fragment.SQL)
		whereValues = append(whereValues, fragment.Values...)
	}

	where := strings.TrimSpace(strings.Join(whereSQLs, " AND "))
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
		whereValues...,
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

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return items, nil
}

func genericSelectRasterOverviewsView(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...types.Fragment) ([]types.DjangolangObject, error) {
	items, err := SelectRasterOverviewsView(ctx, db, columns, orderBy, limit, offset, wheres...)
	if err != nil {
		return nil, err
	}

	genericItems := make([]types.DjangolangObject, 0)
	for _, item := range items {
		genericItems = append(genericItems, item)
	}

	return genericItems, nil
}

func DeserializeRasterOverviewView(b []byte) (types.DjangolangObject, error) {
	var object RasterOverviewView

	err := json.Unmarshal(b, &object)
	if err != nil {
		return nil, err
	}

	return &object, nil
}

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

func (r *RasterOverviewView) Insert(ctx context.Context, db *sqlx.DB, columns ...string) error {
	return fmt.Errorf("not implemented (table has no primary key)")
}

func genericInsertRasterOverviewView(ctx context.Context, db *sqlx.DB, object types.DjangolangObject, columns ...string) (types.DjangolangObject, error) {
	if object == nil {
		return nil, fmt.Errorf("object given for insertion was unexpectedly nil")
	}

	err := object.Insert(ctx, db, columns...)
	if err != nil {
		return nil, err
	}

	return object, nil
}

func (r *RasterOverviewView) GetPrimaryKey() (any, error) {
	return nil, fmt.Errorf("not implemented (table has no primary key)")
}

func (r *RasterOverviewView) SetPrimaryKey(value any) error {
	return fmt.Errorf("not implemented (table has no primary key)")
}

func (r *RasterOverviewView) Update(ctx context.Context, db *sqlx.DB, columns ...string) error {
	return fmt.Errorf("not implemented (table has no primary key)")
}

func genericUpdateRasterOverviewView(ctx context.Context, db *sqlx.DB, object types.DjangolangObject, columns ...string) (types.DjangolangObject, error) {
	if object == nil {
		return nil, fmt.Errorf("object given for update was unexpectedly nil")
	}

	err := object.Update(ctx, db, columns...)
	if err != nil {
		return nil, err
	}

	return object, nil
}

func (r *RasterOverviewView) Delete(ctx context.Context, db *sqlx.DB) error {
	return fmt.Errorf("not implemented (table has no primary key)")
}

func genericDeleteRasterOverviewView(ctx context.Context, db *sqlx.DB, object types.DjangolangObject) error {
	if object == nil {
		return fmt.Errorf("object given for deletion was unexpectedly nil")
	}

	err := object.Delete(ctx, db)
	if err != nil {
		return err
	}

	return nil
}
