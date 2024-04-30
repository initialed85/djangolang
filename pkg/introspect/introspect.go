package introspect

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	_helpers "github.com/initialed85/djangolang/internal/helpers"
	"github.com/initialed85/djangolang/pkg/helpers"
	"github.com/jmoiron/sqlx"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"github.com/lib/pq/hstore"

	geom "github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/geojson"
)

var (
	//go:embed introspect.sql
	createIntrospectViewsSQL string
	mu                       *sync.Mutex = new(sync.Mutex)
	havePushedViews          bool        = false
	logger                               = helpers.GetLogger("introspect")
)

var introspectTablesTemplateSQL = strings.TrimSpace(`
SELECT
    row_to_json(v_introspect_tables)
FROM
    v_introspect_tables
WHERE
    schema = '%v'
    AND relkind IN ('r', 'v');
`)

func getTypeForPartialColumn(column *Column) (any, string, string, string, error) {
	if column == nil {
		return nil, "", "", "", fmt.Errorf("column unexpectedly nil")
	}

	var zeroType any
	queryTypeTemplate := ""
	streamTypeTemplate := ""
	typeTemplate := ""

	dataType := column.DataType
	dataType = strings.Trim(strings.Split(dataType, "(")[0], `"`)

	// TODO: add more of these as we come across them
	switch dataType {

	//
	// slices
	//

	case "timestamp without time zone[]":
		fallthrough
	case "timestamp with time zone[]":
		zeroType = make([]time.Time, 0)
		queryTypeTemplate = "[]time.Time"
	case "timestamp without time zone":
		fallthrough
	case "timestamp with time zone":
		zeroType = time.Time{}
		queryTypeTemplate = "time.Time"

	case "interval[]":
		zeroType = make([]time.Duration, 0)
		queryTypeTemplate = "[]time.Duration"
	case "interval":
		zeroType = helpers.Deref(new(time.Duration))
		queryTypeTemplate = "time.Duration"

	case "json[]":
		fallthrough
	case "jsonb[]":
		zeroType = make([]any, 0)
		queryTypeTemplate = "[]any"
	case "json":
		fallthrough
	case "jsonb":
		zeroType = nil
		queryTypeTemplate = "any"

	case "character varying[]":
		fallthrough
	case "text[]":
		zeroType = make(pq.StringArray, 0)
		queryTypeTemplate = "pq.StringArray"
		typeTemplate = "[]string"
	case "character varying":
		fallthrough
	case "text":
		zeroType = helpers.Deref(new(string))
		queryTypeTemplate = "string"

	case "smallint[]":
		fallthrough
	case "integer[]":
		fallthrough
	case "bigint[]":
		zeroType = make(pq.Int64Array, 0)
		queryTypeTemplate = "pq.Int64Array"
	case "smallint":
		fallthrough
	case "integer":
		fallthrough
	case "bigint":
		zeroType = helpers.Deref(new(int64))
		queryTypeTemplate = "int64"

	case "real[]":
		fallthrough
	case "float[]":
		fallthrough
	case "numeric[]":
		fallthrough
	case "double precision[]":
		zeroType = make(pq.Float64Array, 0)
		queryTypeTemplate = "pq.Float64Array"
	case "float":
		fallthrough
	case "real":
		fallthrough
	case "numeric":
		fallthrough
	case "double precision":
		zeroType = helpers.Deref(new(float64))
		queryTypeTemplate = "float64"

	case "boolean[]":
		zeroType = make(pq.BoolArray, 0)
		queryTypeTemplate = "pq.BoolArray"
	case "boolean":
		zeroType = helpers.Deref(new(bool))
		queryTypeTemplate = "bool"

	case "tsvector[]":
		zeroType = make([]any, 0)
		queryTypeTemplate = "[]any"
	case "tsvector":
		zeroType = nil
		queryTypeTemplate = "any"

	case "uuid[]":
		zeroType = make([]uuid.UUID, 0)
		queryTypeTemplate = "[]uuid.UUID"
		streamTypeTemplate = "[][16]uint8"
	case "uuid":
		zeroType = uuid.UUID{}
		queryTypeTemplate = "uuid.UUID"
		streamTypeTemplate = "[16]uint8"

	case "point[]":
		zeroType = make([]geom.Point, 0)
		queryTypeTemplate = "[]geom.Point"
	case "polygon[]":
		zeroType = make([]geom.Polygon, 0)
		queryTypeTemplate = "[]geom.Polygon"
	case "collection[]":
		zeroType = make([]geom.GeometryCollection, 0)
		queryTypeTemplate = "[]geom.GeometryCollection"
	case "geometry[]":
		zeroType = make([]geojson.Geometry, 0)
		queryTypeTemplate = "[]geojson.Geometry"

	case "hstore[]":
		zeroType = make([]hstore.Hstore, 0)
		queryTypeTemplate = "[]hstore.Hstore"
		typeTemplate = "[]map[string]*string"
	case "hstore":
		zeroType = hstore.Hstore{}
		queryTypeTemplate = "hstore.Hstore"
		typeTemplate = "map[string]*string"

	case "point":
		zeroType = geom.Point{}
		queryTypeTemplate = "geom.Point"
	case "polygon":
		zeroType = geom.Polygon{}
		queryTypeTemplate = "geom.Polygon"
	case "collection":
		zeroType = geom.GeometryCollection{}
		queryTypeTemplate = "geom.GeometryCollection"
	case "geometry":
		zeroType = geojson.Geometry{}
		queryTypeTemplate = "geojson.Geometry"

	default:
		logger.Printf("column = %v", _helpers.UnsafeJSONPrettyFormat(column))
		return nil, "", "", "", fmt.Errorf(
			"failed to work out Go type details for Postgres type %#+v (%v.%v)",
			column.DataType, column.Name, column.TableName,
		)
	}

	if streamTypeTemplate == "" {
		streamTypeTemplate = queryTypeTemplate
	}

	if typeTemplate == "" {
		typeTemplate = queryTypeTemplate
	}

	return zeroType, queryTypeTemplate, streamTypeTemplate, typeTemplate, nil
}

func MapTableByName(originalTableByName map[string]*Table) (map[string]*Table, error) {
	var err error

	tableByName := make(map[string]*Table)
	for _, table := range originalTableByName {
		tableByName[table.Name] = table
	}

	for _, table := range tableByName {
		table.ColumnByName = make(map[string]*Column)
		for _, column := range table.Columns {
			column.ZeroType, column.QueryTypeTemplate, column.StreamTypeTemplate, column.TypeTemplate, err = getTypeForPartialColumn(column)
			if err != nil {
				return nil, err
			}

			table.ColumnByName[column.Name] = column

			if column.IsPrimaryKey {
				table.PrimaryKeyColumn = column
			}
		}

		table.ForeignTables = make([]*Table, 0)
		table.ForeignTableByName = make(map[string]*Table)

		tableByName[table.Name] = table
	}

	for _, table := range tableByName {
		for _, column := range table.ColumnByName {
			if column.ForeignTableName == nil || column.ForeignColumnName == nil {
				continue
			}

			foreignTable, ok := tableByName[*column.ForeignTableName]
			if !ok {
				return nil, fmt.Errorf(
					"%v.%v has foreign key %v.%v but we failed to find that table",
					table.Name, column.Name,
					*column.ForeignTableName,
					*column.ForeignColumnName,
				)
			}
			column.ForeignTable = foreignTable
			table.ForeignTables = append(table.ForeignTables, foreignTable)
			table.ForeignTableByName[foreignTable.Name] = foreignTable

			foreignColumn, ok := foreignTable.ColumnByName[*column.ForeignColumnName]
			if !ok {
				return nil, fmt.Errorf(
					"%v.%v has foreign key %v.%v but we failed to find that table's column",
					table.Name, column.Name,
					*column.ForeignTableName,
					*column.ForeignColumnName,
				)
			}

			column.ForeignColumn = foreignColumn

			table.ColumnByName[column.Name] = column
		}
	}

	return tableByName, nil
}

func Introspect(ctx context.Context, db *sqlx.DB, schema string) (map[string]*Table, error) {
	needToPushViews := false

	mu.Lock()
	needToPushViews = !havePushedViews
	mu.Unlock()

	if needToPushViews {
		logger.Printf("need to push introspection views, pushing...")

		pushViewsCtx, cancel := context.WithTimeout(ctx, time.Second*10)
		defer cancel()
		_, err := db.ExecContext(pushViewsCtx, createIntrospectViewsSQL)
		if err != nil {
			return nil, err
		}
		logger.Printf("done.")

		mu.Lock()
		needToPushViews = false
		mu.Unlock()
	}

	logger.Printf("running introspection query...")
	introspectTablesSQL := fmt.Sprintf(introspectTablesTemplateSQL, schema)
	introspectCtx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	rows, err := db.QueryxContext(introspectCtx, introspectTablesSQL)
	if err != nil {
		return nil, err
	}
	logger.Printf("done.")

	logger.Printf("scanning rows to structs...")
	tables := make([]*Table, 0)
	for rows.Next() {
		var b []byte
		err = rows.Scan(&b)
		if err != nil {
			return nil, err
		}

		var table Table
		err = json.Unmarshal(b, &table)
		if err != nil {
			return nil, err
		}
		tables = append(tables, &table)
	}
	logger.Printf("done.")

	logger.Printf("mapping structs to parents and neighbours...")

	tableByName := make(map[string]*Table)
	for _, table := range tables {
		// ignore the views we made for introspection purposes
		if table.Name == "v_introspect_schemas" ||
			table.Name == "v_introspect_tables" ||
			table.Name == "v_introspect_columns" ||
			table.Name == "v_introspect_table_oids" {
			continue
		}

		// TODO: should this be configurable?
		// ignore tables relating to PostGIS
		if table.Name == "geography_columns" ||
			table.Name == "raster_columns" ||
			table.Name == "raster_overviews" ||
			table.Name == "spatial_ref_sys" ||
			table.Name == "geometry_columns" {
			continue
		}

		tableByName[table.Name] = table
	}

	tableByName, err = MapTableByName(tableByName)
	if err != nil {
		return nil, err
	}

	logger.Printf("done.")

	numTables := len(tableByName)
	numColumns := 0
	numForeignKeys := 0

	for _, table := range tableByName {
		for _, column := range table.ColumnByName {
			numColumns++
			if column.ForeignTable != nil && column.ForeignColumn != nil {
				numForeignKeys++
			}
		}
	}

	logger.Printf("introspected %v tables, %v columns and %v foreign keys", numTables, numColumns, numForeignKeys)

	return tableByName, nil
}

func Run(ctx context.Context) error {
	db, err := helpers.GetDBFromEnvironment(ctx)
	if err != nil {
		return err
	}
	defer func() {
		_ = db.Close()
	}()

	schema := helpers.GetSchema()

	tableByName, err := Introspect(ctx, db, schema)
	if err != nil {
		return err
	}

	for _, table := range tableByName {
		if table.PrimaryKeyColumn != nil {
			logger.Printf(
				"%v.%v | %v = %v (primary key)",
				table.Name,
				table.PrimaryKeyColumn.Name,
				table.PrimaryKeyColumn.DataType,
				table.PrimaryKeyColumn.QueryTypeTemplate,
			)
		} else {
			logger.Printf(
				"%v.%v | %v = %v (primary key)",
				table.Name,
				nil,
				nil,
				nil,
			)
		}

		for _, column := range table.ColumnByName {
			if column.IsPrimaryKey {
				continue
			}

			if column.ForeignTable != nil && column.ForeignColumn != nil {
				logger.Printf(
					"%v.%v -> %v.%v | %v = %v",
					table.Name,
					column.Name,
					column.ForeignTable.Name,
					column.ForeignColumn.Name,
					column.DataType,
					column.QueryTypeTemplate,
				)
			} else {
				logger.Printf(
					"%v.%v | %v = %v",
					table.Name,
					column.Name,
					column.DataType,
					column.QueryTypeTemplate,
				)
			}
		}
	}

	return nil
}
