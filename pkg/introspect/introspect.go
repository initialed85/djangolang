package introspect

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/initialed85/djangolang/pkg/helpers"
	"github.com/initialed85/djangolang/pkg/types"
	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"

	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

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

	dataType := column.DataType
	dataType = strings.Trim(strings.Split(dataType, "(")[0], `"`)

	theType, err := types.GetTypeForDataType(dataType)
	if err != nil {
		return nil, "", "", "", fmt.Errorf(
			"failed to work out Go type details for Postgres type %#+v (adjusted to %#+v) (%v.%v)",
			column.DataType, dataType, column.TableName, column.Name,
		)
	}

	return theType.ZeroType, theType.QueryTypeTemplate, theType.StreamTypeTemplate, theType.TypeTemplate, nil
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

	b, _ := json.MarshalIndent(tableByName, "", "  ")
	log.Printf("%v", string(b))

	return nil
}
