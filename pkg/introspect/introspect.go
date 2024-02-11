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
	"github.com/initialed85/djangolang/pkg/helpers"
	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
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

func getType(column *Column) (any, string, error) {
	if column == nil {
		return nil, "", fmt.Errorf("column unexpectedly nil")
	}

	var zeroType any
	typeTemplate := ""

	dataType := column.DataType
	dataType = strings.Split(dataType, "(")[0]

	// TODO: add more of these as we come across them
	switch dataType {
	case "uuid":
		zeroType = uuid.UUID{}
		typeTemplate = "uuid.UUID"
	case "timestamp without time zone":
		fallthrough
	case "timestamp with time zone":
		zeroType = time.Time{}
		typeTemplate = "time.Time"
	case "json":
		fallthrough
	case "jsonb":
		zeroType = nil
		typeTemplate = "any"
	case "character varying":
		fallthrough
	case "text":
		zeroType = helpers.Deref(new(string))
		typeTemplate = "string"
	case "text[]":
		zeroType = make([]string, 0)
		typeTemplate = "[]string"
	case "integer":
		fallthrough
	case "bigint":
		zeroType = helpers.Deref(new(int64))
		typeTemplate = "int64"
	case "numeric":
		fallthrough
	case "double precision":
		zeroType = helpers.Deref(new(float64))
		typeTemplate = "float64"
	case "boolean":
		zeroType = helpers.Deref(new(bool))
		typeTemplate = "bool"
	case "tsvector": // TODO: as required
		zeroType = nil
		typeTemplate = ""
	default:
		return nil, "", fmt.Errorf(
			"failed to work out Go type details for Postgres type %#+v (%v.%v)",
			column.DataType, column.Name, column.TableName,
		)
	}

	if zeroType != nil || (zeroType == nil && typeTemplate == "any") {
		if !column.NotNull {
			zeroType = helpers.Ptr(zeroType)
			typeTemplate = fmt.Sprintf("*%v", typeTemplate)
		}
	}

	return zeroType, typeTemplate, nil
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
			column.ZeroType, column.TypeTemplate, err = getType(column)
			if err != nil {
				return nil, err
			}

			table.ColumnByName[column.Name] = column
		}

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
	rows, err := db.QueryContext(introspectCtx, introspectTablesSQL)
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
		if table.Name == "v_introspect_schemas" ||
			table.Name == "v_introspect_tables" ||
			table.Name == "v_introspect_columns" ||
			table.Name == "v_introspect_table_oids" {
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
		for _, column := range table.ColumnByName {
			if column.ForeignTable != nil && column.ForeignColumn != nil {
				logger.Printf(
					"%v.%v -> %v.%v | %v = %v",
					table.Name,
					column.Name,
					column.ForeignTable.Name,
					column.ForeignColumn.Name,
					column.DataType,
					column.TypeTemplate,
				)
			} else {
				logger.Printf(
					"%v.%v | %v = %v",
					table.Name,
					column.Name,
					column.DataType,
					column.TypeTemplate,
				)
			}
		}
	}

	return nil
}
