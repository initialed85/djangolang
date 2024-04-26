package some_db

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/initialed85/djangolang/pkg/types"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"github.com/twpayne/go-geom/encoding/geojson"
)

var RasterColumnViewTable = "raster_columns"
var RasterColumnViewTableRTableCatalogColumn = "r_table_catalog"
var RasterColumnViewTableRTableSchemaColumn = "r_table_schema"
var RasterColumnViewTableRTableNameColumn = "r_table_name"
var RasterColumnViewTableRRasterColumnColumn = "r_raster_column"
var RasterColumnViewTableSridColumn = "srid"
var RasterColumnViewTableScaleXColumn = "scale_x"
var RasterColumnViewTableScaleYColumn = "scale_y"
var RasterColumnViewTableBlocksizeXColumn = "blocksize_x"
var RasterColumnViewTableBlocksizeYColumn = "blocksize_y"
var RasterColumnViewTableSameAlignmentColumn = "same_alignment"
var RasterColumnViewTableRegularBlockingColumn = "regular_blocking"
var RasterColumnViewTableNumBandsColumn = "num_bands"
var RasterColumnViewTablePixelTypesColumn = "pixel_types"
var RasterColumnViewTableNodataValuesColumn = "nodata_values"
var RasterColumnViewTableOutDbColumn = "out_db"
var RasterColumnViewTableExtentColumn = "extent"
var RasterColumnViewTableSpatialIndexColumn = "spatial_index"
var RasterColumnViewColumns = []string{"r_table_catalog", "r_table_schema", "r_table_name", "r_raster_column", "srid", "scale_x", "scale_y", "blocksize_x", "blocksize_y", "same_alignment", "regular_blocking", "num_bands", "pixel_types", "nodata_values", "out_db", "extent", "spatial_index"}
var RasterColumnViewTransformedColumns = []string{"r_table_catalog", "r_table_schema", "r_table_name", "r_raster_column", "srid", "scale_x", "scale_y", "blocksize_x", "blocksize_y", "same_alignment", "regular_blocking", "num_bands", "pixel_types", "nodata_values", "out_db", "ST_AsGeoJSON(extent::geometry)::jsonb AS extent", "spatial_index"}
var RasterColumnViewInsertColumns = []string{"r_table_catalog", "r_table_schema", "r_table_name", "r_raster_column", "srid", "scale_x", "scale_y", "blocksize_x", "blocksize_y", "same_alignment", "regular_blocking", "num_bands", "pixel_types", "nodata_values", "out_db", "extent", "spatial_index"}

func SelectRasterColumnsView(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...types.Fragment) ([]*RasterColumnView, error) {
	mu.RLock()
	debug := actualDebug
	mu.RUnlock()

	key := "RasterColumnView"

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
		columns = RasterColumnViewTransformedColumns
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
		"SELECT %v FROM raster_columns%v",
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

	items := make([]*RasterColumnView, 0)
	for rows.Next() {
		rowCount++

		var item RasterColumnView
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

func genericSelectRasterColumnsView(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...types.Fragment) ([]types.DjangolangObject, error) {
	items, err := SelectRasterColumnsView(ctx, db, columns, orderBy, limit, offset, wheres...)
	if err != nil {
		return nil, err
	}

	genericItems := make([]types.DjangolangObject, 0)
	for _, item := range items {
		genericItems = append(genericItems, item)
	}

	return genericItems, nil
}

func DeserializeRasterColumnView(b []byte) (types.DjangolangObject, error) {
	var object RasterColumnView

	err := json.Unmarshal(b, &object)
	if err != nil {
		return nil, err
	}

	return &object, nil
}

type RasterColumnView struct {
	RTableCatalog   *string           `json:"r_table_catalog" db:"r_table_catalog"`
	RTableSchema    *string           `json:"r_table_schema" db:"r_table_schema"`
	RTableName      *string           `json:"r_table_name" db:"r_table_name"`
	RRasterColumn   *string           `json:"r_raster_column" db:"r_raster_column"`
	Srid            *int64            `json:"srid" db:"srid"`
	ScaleX          *float64          `json:"scale_x" db:"scale_x"`
	ScaleY          *float64          `json:"scale_y" db:"scale_y"`
	BlocksizeX      *int64            `json:"blocksize_x" db:"blocksize_x"`
	BlocksizeY      *int64            `json:"blocksize_y" db:"blocksize_y"`
	SameAlignment   *bool             `json:"same_alignment" db:"same_alignment"`
	RegularBlocking *bool             `json:"regular_blocking" db:"regular_blocking"`
	NumBands        *int64            `json:"num_bands" db:"num_bands"`
	PixelTypes      *pq.StringArray   `json:"pixel_types" db:"pixel_types"`
	NodataValues    *pq.Float64Array  `json:"nodata_values" db:"nodata_values"`
	OutDb           *pq.BoolArray     `json:"out_db" db:"out_db"`
	Extent          *geojson.Geometry `json:"extent" db:"extent"`
	SpatialIndex    *bool             `json:"spatial_index" db:"spatial_index"`
}

func (r *RasterColumnView) Insert(ctx context.Context, db *sqlx.DB, columns ...string) error {
	return fmt.Errorf("not implemented (table has no primary key)")
}

func genericInsertRasterColumnView(ctx context.Context, db *sqlx.DB, object types.DjangolangObject, columns ...string) (types.DjangolangObject, error) {
	if object == nil {
		return nil, fmt.Errorf("object given for insertion was unexpectedly nil")
	}

	err := object.Insert(ctx, db, columns...)
	if err != nil {
		return nil, err
	}

	return object, nil
}

func (r *RasterColumnView) GetPrimaryKey() (any, error) {
	return nil, fmt.Errorf("not implemented (table has no primary key)")
}

func (r *RasterColumnView) SetPrimaryKey(value any) error {
	return fmt.Errorf("not implemented (table has no primary key)")
}

func (r *RasterColumnView) Update(ctx context.Context, db *sqlx.DB, columns ...string) error {
	return fmt.Errorf("not implemented (table has no primary key)")
}

func genericUpdateRasterColumnView(ctx context.Context, db *sqlx.DB, object types.DjangolangObject, columns ...string) (types.DjangolangObject, error) {
	if object == nil {
		return nil, fmt.Errorf("object given for update was unexpectedly nil")
	}

	err := object.Update(ctx, db, columns...)
	if err != nil {
		return nil, err
	}

	return object, nil
}

func (r *RasterColumnView) Delete(ctx context.Context, db *sqlx.DB) error {
	return fmt.Errorf("not implemented (table has no primary key)")
}

func genericDeleteRasterColumnView(ctx context.Context, db *sqlx.DB, object types.DjangolangObject) error {
	if object == nil {
		return fmt.Errorf("object given for deletion was unexpectedly nil")
	}

	err := object.Delete(ctx, db)
	if err != nil {
		return err
	}

	return nil
}
