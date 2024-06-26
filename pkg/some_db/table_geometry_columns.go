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

var GeometryColumnViewTable = "geometry_columns"
var GeometryColumnViewTableFTableCatalogColumn = "f_table_catalog"
var GeometryColumnViewTableFTableSchemaColumn = "f_table_schema"
var GeometryColumnViewTableFTableNameColumn = "f_table_name"
var GeometryColumnViewTableFGeometryColumnColumn = "f_geometry_column"
var GeometryColumnViewTableCoordDimensionColumn = "coord_dimension"
var GeometryColumnViewTableSridColumn = "srid"
var GeometryColumnViewTableTypeColumn = "type"
var GeometryColumnViewColumns = []string{"f_table_catalog", "f_table_schema", "f_table_name", "f_geometry_column", "coord_dimension", "srid", "type"}
var GeometryColumnViewTransformedColumns = []string{"f_table_catalog", "f_table_schema", "f_table_name", "f_geometry_column", "coord_dimension", "srid", "type"}
var GeometryColumnViewInsertColumns = []string{"f_table_catalog", "f_table_schema", "f_table_name", "f_geometry_column", "coord_dimension", "srid", "type"}

func SelectGeometryColumnsView(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...types.Fragment) ([]*GeometryColumnView, error) {
	mu.RLock()
	debug := actualDebug
	mu.RUnlock()

	key := "GeometryColumnView"

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
		columns = GeometryColumnViewTransformedColumns
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
		"SELECT %v FROM geometry_columns%v",
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

	items := make([]*GeometryColumnView, 0)
	for rows.Next() {
		rowCount++

		var item GeometryColumnView
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

func genericSelectGeometryColumnsView(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...types.Fragment) ([]types.DjangolangObject, error) {
	items, err := SelectGeometryColumnsView(ctx, db, columns, orderBy, limit, offset, wheres...)
	if err != nil {
		return nil, err
	}

	genericItems := make([]types.DjangolangObject, 0)
	for _, item := range items {
		genericItems = append(genericItems, item)
	}

	return genericItems, nil
}

func DeserializeGeometryColumnView(b []byte) (types.DjangolangObject, error) {
	var object GeometryColumnView

	err := json.Unmarshal(b, &object)
	if err != nil {
		return nil, err
	}

	return &object, nil
}

type GeometryColumnView struct {
	FTableCatalog   *string `json:"f_table_catalog" db:"f_table_catalog"`
	FTableSchema    *string `json:"f_table_schema" db:"f_table_schema"`
	FTableName      *string `json:"f_table_name" db:"f_table_name"`
	FGeometryColumn *string `json:"f_geometry_column" db:"f_geometry_column"`
	CoordDimension  *int64  `json:"coord_dimension" db:"coord_dimension"`
	Srid            *int64  `json:"srid" db:"srid"`
	Type            *string `json:"type" db:"type"`
}

func (g *GeometryColumnView) Insert(ctx context.Context, db *sqlx.DB, columns ...string) error {
	return fmt.Errorf("not implemented (table has no primary key)")
}

func genericInsertGeometryColumnView(ctx context.Context, db *sqlx.DB, object types.DjangolangObject, columns ...string) (types.DjangolangObject, error) {
	if object == nil {
		return nil, fmt.Errorf("object given for insertion was unexpectedly nil")
	}

	err := object.Insert(ctx, db, columns...)
	if err != nil {
		return nil, err
	}

	return object, nil
}

func (g *GeometryColumnView) GetPrimaryKey() (any, error) {
	return nil, fmt.Errorf("not implemented (table has no primary key)")
}

func (g *GeometryColumnView) SetPrimaryKey(value any) error {
	return fmt.Errorf("not implemented (table has no primary key)")
}

func (g *GeometryColumnView) Update(ctx context.Context, db *sqlx.DB, columns ...string) error {
	return fmt.Errorf("not implemented (table has no primary key)")
}

func genericUpdateGeometryColumnView(ctx context.Context, db *sqlx.DB, object types.DjangolangObject, columns ...string) (types.DjangolangObject, error) {
	if object == nil {
		return nil, fmt.Errorf("object given for update was unexpectedly nil")
	}

	err := object.Update(ctx, db, columns...)
	if err != nil {
		return nil, err
	}

	return object, nil
}

func (g *GeometryColumnView) Delete(ctx context.Context, db *sqlx.DB) error {
	return fmt.Errorf("not implemented (table has no primary key)")
}

func genericDeleteGeometryColumnView(ctx context.Context, db *sqlx.DB, object types.DjangolangObject) error {
	if object == nil {
		return fmt.Errorf("object given for deletion was unexpectedly nil")
	}

	err := object.Delete(ctx, db)
	if err != nil {
		return err
	}

	return nil
}
