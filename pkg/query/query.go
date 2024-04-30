package query

import (
	"context"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

func FormatObjectName(objectName string) string {
	return fmt.Sprintf("\"%v\"", objectName)
}

func FormatObjectNames(objectNames []string) []string {
	formattedObjectNames := make([]string, 0)

	for _, objectName := range objectNames {
		if strings.Contains(objectName, `"`) && strings.Contains(objectName, " AS ") {
			formattedObjectNames = append(
				formattedObjectNames,
				objectName,
			)
			continue
		}

		formattedObjectNames = append(
			formattedObjectNames,
			FormatObjectName(objectName),
		)
	}

	return formattedObjectNames
}

func JoinObjectNames(objectNames []string) string {
	return strings.Join(objectNames, ", ")
}

func GetWhere(where string) string {
	where = strings.TrimSpace(where)
	if len(where) == 0 {
		return ""
	}

	return fmt.Sprintf(" WHERE %v", where)
}

func GetLimitAndOffset(limit *int, offset *int) string {
	if limit == nil {
		return ""
	}

	if offset == nil {
		return fmt.Sprintf(" LIMIT %v", *limit)
	}

	return fmt.Sprintf(" LIMIT %v OFFSET %v", *limit, *offset)
}

func Select(
	ctx context.Context,
	tx *sqlx.Tx,
	columns []string,
	table string,
	where string,
	limit *int,
	offset *int,
	values ...any,
) ([]map[string]any, error) {
	sql := fmt.Sprintf(
		"SELECT %v FROM %v%v%v;",
		JoinObjectNames(FormatObjectNames(columns)),
		FormatObjectName(table),
		GetWhere(where),
		GetLimitAndOffset(limit, offset),
	)

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

	sql := fmt.Sprintf(
		"INSERT INTO %v (%v) VALUES (%v) RETURNING (%v);",
		FormatObjectName(table),
		JoinObjectNames(FormatObjectNames(columns)),
		strings.Join(placeholders, ", "),
		JoinObjectNames(FormatObjectNames(returning)),
	)

	rows, err := tx.QueryxContext(
		ctx,
		sql,
		values...,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to call tx.QueryxContext during Insert; err: %v, sql: %#+v",
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
	for i := range values {
		placeholders = append(placeholders, fmt.Sprintf("$%v", i+1))
	}

	sql := fmt.Sprintf(
		"UPDATE %v SET (%v) = (%v)%v RETURNING (%v);",
		FormatObjectName(table),
		JoinObjectNames(FormatObjectNames(columns)),
		strings.Join(placeholders, ", "),
		GetWhere(where),
		JoinObjectNames(FormatObjectNames(returning)),
	)

	rows, err := tx.QueryxContext(
		ctx,
		sql,
		values...,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to call tx.QueryxContext during Insert; err: %v, sql: %#+v",
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

func Delete(
	ctx context.Context,
	db *sqlx.Tx,
	table string,
	where string,
	values []any,
) error {
	return nil
}
