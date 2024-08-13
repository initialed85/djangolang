package query

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/initialed85/djangolang/pkg/helpers"
	"github.com/jmoiron/sqlx"
)

func FormatObjectName(objectName string) string {
	return fmt.Sprintf("\"%v\"", objectName)
}

func FormatObjectNames(objectNames []string, forceAliases ...bool) []string {
	forceAlias := len(forceAliases) > 0 && forceAliases[0]

	formattedObjectNames := make([]string, 0)

	for _, objectName := range objectNames {
		if strings.Contains(objectName, `"`) && strings.Contains(objectName, " AS ") {
			formattedObjectNames = append(
				formattedObjectNames,
				objectName,
			)
			continue
		}

		formattedObjectName := FormatObjectName(objectName)

		if forceAlias {
			formattedObjectName += fmt.Sprintf(" AS %s", strings.Trim(objectName, `"`))
		}

		formattedObjectNames = append(
			formattedObjectNames,
			formattedObjectName,
		)
	}

	return formattedObjectNames
}

func JoinObjectNames(objectNames []string) string {
	return strings.Join(objectNames, ",\n    ")
}

func GetWhere(where string) string {
	where = strings.TrimSpace(where)
	if len(where) == 0 {
		return ""
	}

	return fmt.Sprintf("\nWHERE\n    %v", where)
}

func GetOrderBy(orderBy *string) string {
	if orderBy == nil {
		return ""
	}

	return fmt.Sprintf("\nORDER BY\n    %v", strings.TrimSpace(*orderBy))
}

func GetLimitAndOffset(limit *int, offset *int) string {
	if limit == nil {
		return ""
	}

	if offset == nil {
		return fmt.Sprintf("\nLIMIT %v", *limit)
	}

	return fmt.Sprintf("\nLIMIT %v\nOFFSET %v", *limit, *offset)
}

func Select(
	ctx context.Context,
	tx *sqlx.Tx,
	columns []string,
	table string,
	where string,
	orderBy *string,
	limit *int,
	offset *int,
	values ...any,
) ([]map[string]any, error) {
	i := 1
	for strings.Contains(where, "$$??") {
		where = strings.Replace(where, "$$??", fmt.Sprintf("$%d", i), 1)
		i++
	}

	sql := strings.TrimSpace(fmt.Sprintf(
		"SELECT\n    %v\nFROM\n    %v%v%v%v;",
		JoinObjectNames(FormatObjectNames(columns)),
		FormatObjectName(table),
		GetWhere(where),
		GetOrderBy(orderBy),
		GetLimitAndOffset(limit, offset),
	))

	queryID := GetQueryID(ctx)

	cacheKey, err := getCacheKey(
		queryID,
		columns,
		table,
		where,
		orderBy,
		limit,
		offset,
		values,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to call getCacheKey during Select; err: %v, sql: %#+v",
			err, sql,
		)
	}

	cachedItems, ok := getCachedItems(cacheKey)
	if ok {
		return cachedItems, nil
	}

	if helpers.IsDebug() {
		rawValues := ""

		for i, v := range values {
			rawValues += fmt.Sprintf("$%d = %#+v\n", i+1, v)
		}

		log.Printf("\n\n%s\n\n%s\n", sql, rawValues)
	}

	rows, err := tx.QueryxContext(
		ctx,
		sql,
		values...,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to call tx.QueryxContext during Select; err: %v, sql: %#+v",
			err, sql,
		)
	}

	defer func() {
		_ = rows.Close()
	}()

	items := make([]map[string]any, 0)

	for rows.Next() {
		item := make(map[string]any)

		err = rows.MapScan(item)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to call rows.MapScan during Select; err: %v, sql: %#+v, item: %#+v",
				err, sql, item,
			)
		}

		items = append(items, item)
	}

	setCachedItems(queryID, cacheKey, items)

	return items, nil
}

func Insert(
	ctx context.Context,
	tx *sqlx.Tx,
	table string,
	columns []string,
	conflictColumnNames []string, // TODO
	onConflictDoNothing bool, // TODO
	onConflictUpdate bool, // TODO
	returning []string,
	values ...any,
) (map[string]any, error) {
	placeholders := []string{}
	for i := range values {
		placeholders = append(placeholders, fmt.Sprintf("$%v", i+1))
	}

	sql := strings.TrimSpace(fmt.Sprintf(
		"INSERT INTO %v (\n    %v\n) VALUES (\n    %v\n) RETURNING \n    %v;",
		FormatObjectName(table),
		JoinObjectNames(FormatObjectNames(columns)),
		strings.Join(placeholders, ",\n    "),
		JoinObjectNames(FormatObjectNames(returning, true)),
	))

	if helpers.IsDebug() {
		rawValues := ""

		for i, v := range values {
			rawValues += fmt.Sprintf("$%d = %#+v\n", i+1, v)
		}

		log.Printf("\n\n%s\n\n%s\n", sql, rawValues)
	}

	rows, err := tx.QueryxContext(
		ctx,
		sql,
		values...,
	)
	if err != nil {
		err = fmt.Errorf(
			"failed to call tx.QueryxContext during Insert; err: %v, sql: %#+v",
			err, sql,
		)

		if helpers.IsDebug() {
			log.Printf("<<<\n\n%s", err.Error())
		}

		return nil, err
	}

	defer func() {
		_ = rows.Close()
	}()

	items := make([]map[string]any, 0)

	for rows.Next() {
		item := make(map[string]any)

		err = rows.MapScan(item)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to call rows.MapScan during Insert; err: %v, sql: %#+v, item: %#+v",
				err, sql, item,
			)
		}

		items = append(items, item)
	}

	if len(items) < 1 {
		return nil, fmt.Errorf(
			"unexpectedly got no returned rows after Insert; err: %v, sql: %#+v",
			err, sql,
		)
	}

	if len(items) > 1 {
		return nil, fmt.Errorf(
			"unexpectedly got more than 1 returned row after Insert; err: %v, sql: %#+v",
			err, sql,
		)
	}

	return items[0], nil
}

func Update(
	ctx context.Context,
	tx *sqlx.Tx,
	table string,
	columns []string,
	where string,
	returning []string,
	values ...any,
) (map[string]any, error) {
	placeholders := []string{}
	for i := 0; i < len(values)-strings.Count(where, "$$??"); i++ {
		placeholders = append(placeholders, fmt.Sprintf("$%v", i+1))
	}

	i := len(placeholders) + 1
	for strings.Contains(where, "$$??") {
		where = strings.Replace(where, "$$??", fmt.Sprintf("$%d", i), 1)
		i++
	}

	var sql string

	if len(columns) > 1 {
		sql = "UPDATE %v\nSET (\n    %v\n) = (\n    %v\n)%v\nRETURNING\n    %v;"
	} else {
		sql = "UPDATE %v\nSET \n    %v\n = \n    %v\n%v\nRETURNING\n    %v;"
	}

	sql = strings.TrimSpace(fmt.Sprintf(
		sql,
		FormatObjectName(table),
		JoinObjectNames(FormatObjectNames(columns)),
		strings.Join(placeholders, ",\n    "),
		GetWhere(where),
		JoinObjectNames(FormatObjectNames(returning)),
	))

	if helpers.IsDebug() {
		rawValues := ""

		for i, v := range values {
			rawValues += fmt.Sprintf("$%d = %#+v\n", i+1, v)
		}

		log.Printf("\n\n%s\n\n%s\n", sql, rawValues)
	}

	rows, err := tx.QueryxContext(
		ctx,
		sql,
		values...,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to call tx.QueryxContext during Update; err: %v, sql: %#+v",
			err, sql,
		)
	}

	defer func() {
		_ = rows.Close()
	}()

	items := make([]map[string]any, 0)

	for rows.Next() {
		item := make(map[string]any)

		err = rows.MapScan(item)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to call rows.MapScan during Update; err: %v, sql: %#+v, item: %#+v",
				err, sql, item,
			)
		}

		items = append(items, item)
	}

	if len(items) < 1 {
		return nil, fmt.Errorf(
			"unexpectedly got no returned rows after Update; err: %v, sql: %#+v",
			err, sql,
		)
	}

	if len(items) > 1 {
		return nil, fmt.Errorf(
			"unexpectedly got more than 1 returned row after Update; err: %v, sql: %#+v",
			err, sql,
		)
	}

	return items[0], nil
}

func Delete(
	ctx context.Context,
	tx *sqlx.Tx,
	table string,
	where string,
	values ...any,
) error {
	i := 1
	for strings.Contains(where, "$$??") {
		where = strings.Replace(where, "$$??", fmt.Sprintf("$%d", i), 1)
		i++
	}

	sql := fmt.Sprintf(
		"DELETE FROM %v%v;",
		FormatObjectName(table),
		GetWhere(where),
	)

	result, err := tx.ExecContext(
		ctx,
		sql,
		values...,
	)
	if err != nil {
		return fmt.Errorf(
			"failed to call tx.QueryxContext during Delete; err: %v, sql: %#+v",
			err, sql,
		)
	}

	// TODO: cleaner handling for soft delete vs hard delete
	// rowsAffected, err := result.RowsAffected()
	// if err != nil {
	// 	return fmt.Errorf(
	// 		"failed to call result.RowsAffected during Delete; err: %v, sql: %#+v",
	// 		err, sql,
	// 	)
	// }

	// TODO: cleaner handling for soft delete vs hard delete
	// if rowsAffected <= 0 {
	// 	return fmt.Errorf(
	// 		"result.RowsAffected did not return a positive number during Delete; err: %v, sql: %#+v",
	// 		err, sql,
	// 	)
	// }

	// TODO: cleaner handling for soft delete vs hard delete
	_ = result

	return nil
}

func GetXid(
	ctx context.Context,
	tx *sqlx.Tx,
) (uint32, error) {
	var xid uint32
	row := tx.QueryRowContext(ctx, "SELECT txid_current();")

	err := row.Err()
	if err != nil {
		return 0, fmt.Errorf("failed to call txid_current(): %v", err)
	}

	err = row.Scan(&xid)
	if err != nil {
		return 0, fmt.Errorf("failed to call txid_current(): %v", err)
	}

	return xid, nil
}
