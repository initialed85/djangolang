package query

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/initialed85/djangolang/pkg/helpers"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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

func Select(ctx context.Context, tx pgx.Tx, columns []string, table string, where string, orderBy *string, limit *int, offset *int, values ...any) (*[]map[string]any, int64, int64, int64, int64, error) {
	i := 1
	for strings.Contains(where, "$$??") {
		where = strings.Replace(where, "$$??", fmt.Sprintf("$%d", i), 1)
		i++
	}

	if limit != nil {
		if *limit < 0 {
			return nil, 0, 0, 0, 0, fmt.Errorf(
				"invalid limit during Select; %v",
				fmt.Errorf("limit must not be negative"),
			)
		}
	}

	if offset != nil {
		if *offset < 0 {
			return nil, 0, 0, 0, 0, fmt.Errorf(
				"invalid offset during Select; %v",
				fmt.Errorf("offset must not be negative"),
			)
		}
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
		return nil, 0, 0, 0, 0, fmt.Errorf(
			"failed to call getCacheKey during Select; %v, sql: %#+v",
			err, sql,
		)
	}

	cachedItems, ok := getCachedItems(cacheKey)
	if ok {
		return &cachedItems.Items, cachedItems.Count, cachedItems.TotalCount, 1, 1, nil
	}

	if helpers.IsDebug() {
		rawValues := ""

		for i, v := range values {
			rawValues += fmt.Sprintf("$%d = %#+v\n", i+1, v)
		}

		log.Printf("\n\n%s\n\n%s\n", sql, rawValues)
	}

	sqlForRowEstimate := strings.TrimSpace(fmt.Sprintf(
		"SELECT\n    %v\nFROM\n    %v%v%v;",
		JoinObjectNames(FormatObjectNames(columns)),
		FormatObjectName(table),
		GetWhere(where),
		GetOrderBy(orderBy),
	))

	totalCount, err := GetRowEstimate(ctx, tx, sqlForRowEstimate, values...)
	if err != nil {
		return nil, 0, 0, 0, 0, fmt.Errorf(
			"failed to call GetRowEstimate during Select; %v, sql: %#+v",
			err, sqlForRowEstimate,
		)
	}

	rows, err := tx.Query(
		ctx,
		sql,
		values...,
	)
	if err != nil {
		return nil, 0, 0, 0, 0, fmt.Errorf(
			"failed to call tx.Query during Select; %v, sql: %#+v",
			err, sql,
		)
	}

	defer func() {
		rows.Close()
	}()

	items := make([]map[string]any, 0)

	for rows.Next() {
		item := make(map[string]any)

		values, err := rows.Values()
		if err != nil {
			return nil, 0, 0, 0, 0, fmt.Errorf(
				"failed to call rows.Values during Select; %v, sql: %#+v, item: %#+v",
				err, sql, item,
			)
		}

		fieldDescriptions := rows.FieldDescriptions()
		if fieldDescriptions == nil {
			fieldDescriptions = make([]pgconn.FieldDescription, len(values))
		}

		for i, v := range values {
			f := fieldDescriptions[i]
			item[f.Name] = v
		}

		items = append(items, item)
	}

	err = rows.Err()
	if err != nil {
		return nil, 0, 0, 0, 0, fmt.Errorf(
			"failed to call rows.Err after Select; %v, sql: %#+v",
			err, sql,
		)
	}

	count := int64(len(items))

	actualLimit := int64(0)
	if limit != nil {
		actualLimit = int64(*limit)
	}

	actualOffset := int64(0)
	if offset != nil {
		actualOffset = int64(*offset)
	}

	setCachedItems(queryID, cacheKey, &CachedItems{
		Items:      items,
		Count:      count,
		TotalCount: totalCount,
		Limit:      actualLimit,
		Offset:     actualOffset,
	})

	return &items, count, totalCount, int64(1), int64(1), nil
}

func Insert(ctx context.Context, tx pgx.Tx, table string, columns []string, conflictColumnNames []string, onConflictDoNothing bool, onConflictUpdate bool, returning []string, values ...any) (*map[string]any, error) {
	placeholders := make([]string, 0)
	adjustedValues := make([]any, 0)
	i := 0
	for _, v := range values {
		s, ok := v.(string)
		if ok {
			// TODO: this is a hack, it's called for in the types/functions.go
			if strings.HasPrefix(s, "ST_PointZ (") {
				placeholders = append(placeholders, s)
				continue
			}
		}

		placeholders = append(placeholders, fmt.Sprintf("$%v", i+1))
		adjustedValues = append(adjustedValues, v)

		i++
	}

	values = adjustedValues

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

	rows, err := tx.Query(
		ctx,
		sql,
		values...,
	)
	if err != nil {
		err = fmt.Errorf(
			"failed to call tx.Query during Insert; %v, sql: %#+v",
			err, sql,
		)

		if helpers.IsDebug() {
			log.Printf("<<<\n\n%s", err.Error())
		}

		return nil, err
	}

	defer func() {
		rows.Close()
	}()

	items := make([]map[string]any, 0)

	for rows.Next() {
		item := make(map[string]any)

		values, err := rows.Values()
		if err != nil {
			return nil, fmt.Errorf(
				"failed to call rows.Values during Insert; %v, sql: %#+v, item: %#+v",
				err, sql, item,
			)
		}

		fieldDescriptions := rows.FieldDescriptions()
		if fieldDescriptions == nil {
			fieldDescriptions = make([]pgconn.FieldDescription, len(values))
		}

		for i, v := range values {
			f := fieldDescriptions[i]
			item[f.Name] = v
		}

		items = append(items, item)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf(
			"failed to call rows.Err after Insert; %v, sql: %#+v, values: %#+v, items: %#+v",
			err, sql, values, items,
		)
	}

	if len(items) < 1 {
		return nil, fmt.Errorf(
			"unexpectedly got no returned rows after Insert; %v, sql: %#+v",
			err, sql,
		)
	}

	if len(items) > 1 {
		return nil, fmt.Errorf(
			"unexpectedly got more than 1 returned row after Insert; %v, sql: %#+v",
			err, sql,
		)
	}

	return &items[0], nil
}

func Update(ctx context.Context, tx pgx.Tx, table string, columns []string, where string, returning []string, values ...any) (*map[string]any, error) {
	placeholders := []string{}
	adjustedValues := make([]any, 0)
	j := 0
	for i := 0; i < len(values)-strings.Count(where, "$$??"); i++ {
		v := values[i]

		s, ok := v.(string)
		if ok {
			// TODO: this is a hack, it's called for in the types/functions.go
			if strings.HasPrefix(s, "ST_PointZ (") {
				placeholders = append(placeholders, s)
				continue
			}
		}

		placeholders = append(placeholders, fmt.Sprintf("$%v", j+1))
		adjustedValues = append(adjustedValues, v)
		j++
	}
	values = append(adjustedValues, values[len(values)-strings.Count(where, "$$??"):]...)

	i := j + 1
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

	rows, err := tx.Query(
		ctx,
		sql,
		values...,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to call tx.Query during Update; %v, sql: %#+v",
			err, sql,
		)
	}

	defer func() {
		rows.Close()
	}()

	items := make([]map[string]any, 0)

	for rows.Next() {
		item := make(map[string]any)

		values, err := rows.Values()
		if err != nil {
			return nil, fmt.Errorf(
				"failed to call rows.Values during Select; %v, sql: %#+v, item: %#+v",
				err, sql, item,
			)
		}

		fieldDescriptions := rows.FieldDescriptions()
		if fieldDescriptions == nil {
			fieldDescriptions = make([]pgconn.FieldDescription, len(values))
		}

		for i, v := range values {
			f := fieldDescriptions[i]
			item[f.Name] = v
		}

		items = append(items, item)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf(
			"failed to call rows.Err after Update; %v, sql: %#+v",
			err, sql,
		)
	}

	if len(items) < 1 {
		return nil, fmt.Errorf(
			"unexpectedly got no returned rows after Update; %v, sql: %#+v",
			err, sql,
		)
	}

	if len(items) > 1 {
		return nil, fmt.Errorf(
			"unexpectedly got more than 1 returned row after Update; %v, sql: %#+v",
			err, sql,
		)
	}

	return &items[0], nil
}

func Delete(ctx context.Context, tx pgx.Tx, table string, where string, values ...any) error {
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

	result, err := tx.Exec(
		ctx,
		sql,
		values...,
	)
	if err != nil {
		return fmt.Errorf(
			"failed to call tx.Query during Delete; %v, sql: %#+v",
			err, sql,
		)
	}

	// TODO: cleaner handling for soft delete vs hard delete
	// rowsAffected, err := result.RowsAffected()
	// if err != nil {
	// 	return fmt.Errorf(
	// 		"failed to call result.RowsAffected during Delete; %v, sql: %#+v",
	// 		err, sql,
	// 	)
	// }

	// TODO: cleaner handling for soft delete vs hard delete
	// if rowsAffected <= 0 {
	// 	return fmt.Errorf(
	// 		"result.RowsAffected did not return a positive number during Delete; %v, sql: %#+v",
	// 		err, sql,
	// 	)
	// }

	// TODO: cleaner handling for soft delete vs hard delete
	_ = result

	return nil
}

func GetXid(ctx context.Context, tx pgx.Tx) (uint32, error) {
	var xid uint32
	row := tx.QueryRow(ctx, "SELECT txid_current();")

	err := row.Scan(&xid)
	if err != nil {
		return 0, fmt.Errorf("failed to call txid_current(): %v", err)
	}

	return xid, nil
}

func LockTable(ctx context.Context, tx pgx.Tx, tableName string, noWait bool) error {
	noWaitInfix := ""
	if noWait {
		noWaitInfix = " NOWAIT"
	}

	_, err := tx.Exec(ctx, fmt.Sprintf("LOCK TABLE %v IN ACCESS EXCLUSIVE MODE%s;", tableName, noWaitInfix))
	if err != nil {
		return err
	}

	return nil
}

func LockTableWithRetries(ctx context.Context, tx pgx.Tx, tableName string, overallTimeout time.Duration, individualAttemptTimeouts ...time.Duration) error {
	individualAttemptTimeout := time.Second * 1
	if len(individualAttemptTimeouts) > 0 {
		individualAttemptTimeout = individualAttemptTimeouts[0]
	}

	expiry := time.Now().Add(overallTimeout)

	rawSavePointID, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	savePointID := fmt.Sprintf("savepoint_%s", strings.ReplaceAll(rawSavePointID.String(), "-", ""))

	_, err = tx.Exec(ctx, fmt.Sprintf("SAVEPOINT %s;", savePointID))
	if err != nil {
		return err
	}

	for time.Now().Before(expiry) {
		ok, err := func() (bool, error) {
			_, err = tx.Exec(ctx, fmt.Sprintf("SET LOCAL lock_timeout = '%fs';", individualAttemptTimeout.Seconds()))
			if err != nil {
				return false, err
			}

			err = LockTable(ctx, tx, tableName, false)
			if err != nil {
				if strings.Contains(err.Error(), "canceling statement due to lock timeout") {
					return false, nil
				}

				_, rollbackErr := tx.Exec(ctx, fmt.Sprintf("ROLLBACK TO SAVEPOINT %s;", savePointID))
				if rollbackErr != nil {
					return false, rollbackErr
				}

				return false, err
			}

			return true, nil
		}()
		if err != nil {
			return err
		}

		if !ok {
			continue
		}

		return nil
	}

	return fmt.Errorf("timed out waiting to lock table %#+v after %s ", tableName, overallTimeout)
}

func Explain(ctx context.Context, tx pgx.Tx, sql string, values ...any) ([]map[string]any, error) {
	rows, err := tx.Query(
		ctx,
		fmt.Sprintf("EXPLAIN (FORMAT JSON) %s;", strings.TrimRight(sql, ";")),
		values...,
	)
	if err != nil {
		return nil, err
	}

	explanations := make([]map[string]any, 0)

	for rows.Next() {
		err = rows.Scan(&explanations)
		if err != nil {
			return nil, err
		}
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return explanations, nil
}

func GetRowEstimate(ctx context.Context, tx pgx.Tx, sql string, values ...any) (int64, error) {
	explanations, err := Explain(ctx, tx, sql, values...)
	if err != nil {
		return 0, err
	}

	if len(explanations) != 1 {
		return 0, fmt.Errorf("wanted exactly 1 item, got %d for explanations: %#+v", len(explanations), explanations)
	}

	rawPlan, ok := explanations[0]["Plan"]
	if !ok {
		return 0, fmt.Errorf("failed to find explanations[0]['Plan'] key in explanations: %#+v", explanations)
	}

	plan, ok := rawPlan.(map[string]any)
	if !ok {
		return 0, fmt.Errorf("could not cast explanations[0]['Plan'] to map[string]any for explanations: %#+v", explanations)
	}

	rawRowEstimate, ok := plan["Plan Rows"]
	if !ok {
		return 0, fmt.Errorf("failed to find explanations[0]['Plan']['Plan Rows'] key in explanations: %#+v", explanations)
	}

	rowEstimate, ok := rawRowEstimate.(float64)
	if !ok {
		return 0, fmt.Errorf("could not cast explanations[0]['Plan']['Plan Rows'] to float64 for explanations: %#+v", explanations)
	}

	return int64(rowEstimate), nil
}
