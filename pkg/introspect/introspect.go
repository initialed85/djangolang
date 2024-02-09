package introspect

import (
	"context"
	_ "embed"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/initialed85/djangolang/internal/helpers"
	"github.com/jackc/pgx/v5"
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

func Introspect(ctx context.Context, conn *pgx.Conn, schema string) (map[string]*Table, error) {
	needToPushViews := false

	mu.Lock()
	needToPushViews = !havePushedViews
	mu.Unlock()

	if needToPushViews {
		logger.Printf("need to push introspection views, pushing...")

		pushViewsCtx, cancel := context.WithTimeout(ctx, time.Second*60)
		defer cancel()
		results, err := conn.PgConn().Exec(pushViewsCtx, createIntrospectViewsSQL).ReadAll()
		if err != nil {
			return nil, err
		}

		for _, result := range results {
			if result.Err != nil {
				return nil, result.Err
			}
		}
		logger.Printf("done.")

		mu.Lock()
		needToPushViews = false
		mu.Unlock()
	}

	logger.Printf("running introspection query...")
	introspectTablesSQL := fmt.Sprintf(introspectTablesTemplateSQL, schema)
	introspectCtx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()
	rows, err := conn.Query(introspectCtx, introspectTablesSQL)
	if err != nil {
		return nil, err
	}
	logger.Printf("done.")

	logger.Printf("scanning rows to structs...")
	tables, err := pgx.CollectRows[*Table](
		rows,
		func(row pgx.CollectableRow) (*Table, error) {
			var table Table
			err = row.Scan(&table)
			if err != nil {
				return nil, err
			}

			return &table, nil
		},
	)
	if err != nil {
		return nil, err
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

		table.ColumnByName = make(map[string]*Column)
		for _, column := range table.Columns {
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
	conn, err := helpers.GetDBConnFromEnvironment(ctx)
	if err != nil {
		return err
	}
	defer func() {
		closeCtx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		_ = conn.Close(closeCtx)
	}()

	schema := helpers.GetSchema()

	tableByName, err := Introspect(ctx, conn, schema)
	if err != nil {
		return err
	}

	for _, table := range tableByName {
		for _, column := range table.ColumnByName {
			if column.ForeignTable != nil && column.ForeignColumn != nil {
				logger.Printf(
					"%v.%v (%v) -> %v.%v (%v)",
					table.Name,
					column.Name,
					column.DataType,
					column.ForeignTable.Name,
					column.ForeignColumn.Name,
					column.ForeignColumn.DataType,
				)
			} else {
				logger.Printf(
					"%v.%v (%v)",
					table.Name,
					column.Name,
					column.DataType,
				)
			}
		}
	}

	return nil
}
