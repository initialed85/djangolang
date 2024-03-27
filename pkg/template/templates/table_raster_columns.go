package templates

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

var RasterColumnViewTable = "raster_columns"
var RasterColumnViewRTableCatalogColumn = "r_table_catalog"
var RasterColumnViewRTableSchemaColumn = "r_table_schema"
var RasterColumnViewRTableNameColumn = "r_table_name"
var RasterColumnViewRRasterColumnColumn = "r_raster_column"
var RasterColumnViewSridColumn = "srid"
var RasterColumnViewScaleXColumn = "scale_x"
var RasterColumnViewScaleYColumn = "scale_y"
var RasterColumnViewBlocksizeXColumn = "blocksize_x"
var RasterColumnViewBlocksizeYColumn = "blocksize_y"
var RasterColumnViewSameAlignmentColumn = "same_alignment"
var RasterColumnViewRegularBlockingColumn = "regular_blocking"
var RasterColumnViewNumBandsColumn = "num_bands"
var RasterColumnViewPixelTypesColumn = "pixel_types"
var RasterColumnViewNodataValuesColumn = "nodata_values"
var RasterColumnViewOutDbColumn = "out_db"
var RasterColumnViewExtentColumn = "extent"
var RasterColumnViewSpatialIndexColumn = "spatial_index"
var RasterColumnViewColumns = []string{"r_table_catalog", "r_table_schema", "r_table_name", "r_raster_column", "srid", "scale_x", "scale_y", "blocksize_x", "blocksize_y", "same_alignment", "regular_blocking", "num_bands", "pixel_types", "nodata_values", "out_db", "extent", "spatial_index"}
var RasterColumnViewTransformedColumns = []string{"r_table_catalog", "r_table_schema", "r_table_name", "r_raster_column", "srid", "scale_x", "scale_y", "blocksize_x", "blocksize_y", "same_alignment", "regular_blocking", "num_bands", "pixel_types", "nodata_values", "out_db", "ST_AsGeoJSON(extent::geometry)::jsonb AS extent", "spatial_index"}

type RasterColumnView struct {
	RTableCatalog   *string    `json:"r_table_catalog" db:"r_table_catalog"`
	RTableSchema    *string    `json:"r_table_schema" db:"r_table_schema"`
	RTableName      *string    `json:"r_table_name" db:"r_table_name"`
	RRasterColumn   *string    `json:"r_raster_column" db:"r_raster_column"`
	Srid            *int64     `json:"srid" db:"srid"`
	ScaleX          *float64   `json:"scale_x" db:"scale_x"`
	ScaleY          *float64   `json:"scale_y" db:"scale_y"`
	BlocksizeX      *int64     `json:"blocksize_x" db:"blocksize_x"`
	BlocksizeY      *int64     `json:"blocksize_y" db:"blocksize_y"`
	SameAlignment   *bool      `json:"same_alignment" db:"same_alignment"`
	RegularBlocking *bool      `json:"regular_blocking" db:"regular_blocking"`
	NumBands        *int64     `json:"num_bands" db:"num_bands"`
	PixelTypes      *[]string  `json:"pixel_types" db:"pixel_types"`
	NodataValues    *[]float64 `json:"nodata_values" db:"nodata_values"`
	OutDb           *[]bool    `json:"out_db" db:"out_db"`
	Extent          any        `json:"extent" db:"extent"`
	SpatialIndex    *bool      `json:"spatial_index" db:"spatial_index"`
}

func SelectRasterColumnsView(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]*RasterColumnView, error) {
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

	selectCtx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()

	rows, err := db.QueryxContext(selectCtx, sql)
	if err != nil {
		return nil, err
	}

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

		var temp1Extent any
		var temp2Extent []any

		if item.Extent != nil {
			err = json.Unmarshal(item.Extent.([]byte), &temp1Extent)
			if err != nil {
				err = json.Unmarshal(item.Extent.([]byte), &temp2Extent)
				if err != nil {
					item.Extent = nil
				} else {
					item.Extent = &temp2Extent
				}
			} else {
				item.Extent = &temp1Extent
			}
		}

		item.Extent = &temp1Extent

		// this is where foreign object ID map populations would appear (if required)

		items = append(items, &item)
	}

	scanStop = time.Now().UnixNano()

	foreignObjectStart = time.Now().UnixNano()

	// this is where any foreign object loading would appear (if required)

	foreignObjectStop = time.Now().UnixNano()

	return items, nil
}

func genericSelectRasterColumnsView(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]any, error) {
	items, err := SelectRasterColumnsView(ctx, db, columns, orderBy, limit, offset, wheres...)
	if err != nil {
		return nil, err
	}

	genericItems := make([]any, 0)
	for _, item := range items {
		genericItems = append(genericItems, item)
	}

	return genericItems, nil
}
