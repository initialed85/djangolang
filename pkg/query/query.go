package query

import (
	"context"
	"fmt"
	_log "log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/initialed85/djangolang/pkg/config"
	"github.com/initialed85/djangolang/pkg/helpers"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var log = helpers.GetLogger("query")

func ThisLogger() *_log.Logger {
	return log
}

func FormatObjectName(objectName string) string {
	return fmt.Sprintf("\"%v\"", objectName)
}

func FormatTableName(tableName string) string {
	// TODO: fix this hack (all best are off for a join or *gulp* a table with the name of " JOIN " lol)
	if strings.Contains(tableName, " JOIN ") {
		return tableName
	}

	if strings.Contains(tableName, ".") {
		parts := strings.Split(tableName, ".")
		return fmt.Sprintf(`%s.%s`, FormatObjectName(parts[0]), FormatObjectName(parts[1]))
	}

	return FormatObjectName(tableName)
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
	before := time.Now()

	if config.Debug() {
		log.Printf("SELECT FROM %s entered", table)
	}

	i := 1
	for strings.Contains(where, "$$??") {
		where = strings.Replace(where, "$$??", fmt.Sprintf("$%d", i), 1)
		i++
	}

	if limit != nil {
		if *limit < 0 {
			if config.Debug() {
				log.Printf("SELECT FROM %s exited; query cache not reached yet; failed query in %s", table, time.Since(before))
			}

			return nil, 0, 0, 0, 0, fmt.Errorf(
				"invalid limit during Select; %v",
				fmt.Errorf("limit must not be negative"),
			)
		}
	}

	if offset != nil {
		if *offset < 0 {
			if config.Debug() {
				log.Printf("SELECT FROM %s exited; query cache not reached yet; failed query in %s", table, time.Since(before))
			}

			return nil, 0, 0, 0, 0, fmt.Errorf(
				"invalid offset during Select; %v",
				fmt.Errorf("offset must not be negative"),
			)
		}
	}

	sql := strings.TrimSpace(fmt.Sprintf(
		"SELECT\n    %v\nFROM\n    %v%v%v%v;",
		JoinObjectNames(FormatObjectNames(columns)),
		FormatTableName(table),
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
		if config.Debug() {
			log.Printf("SELECT FROM %s exited; query cache not reached yet; failed query in %s", table, time.Since(before))
		}

		return nil, 0, 0, 0, 0, fmt.Errorf(
			"failed to call getCacheKey during Select; %v, sql: %#+v",
			err, sql,
		)
	}

	if config.QueryDebug() {
		rawValues := ""

		for i, v := range values {
			rawValues += fmt.Sprintf("$%d = %#+v\n", i+1, v)
		}

		log.Printf("\n\n%s\n\n%s\n", sql, rawValues)
	}

	cachedItems, ok := getCachedItems(cacheKey)
	if ok {
		if config.Debug() {
			log.Printf("SELECT FROM %s exited; query cache hit; lookup in %s\n", table, time.Since(before))
		}

		return &cachedItems.Items, cachedItems.Count, cachedItems.TotalCount, 1, 1, nil
	}

	sqlForRowEstimate := strings.TrimSpace(fmt.Sprintf(
		"SELECT\n    %v\nFROM\n    %v%v%v;",
		JoinObjectNames(FormatObjectNames(columns)),
		FormatTableName(table),
		GetWhere(where),
		GetOrderBy(orderBy),
	))

	totalCount, err := GetRowEstimate(ctx, tx, sqlForRowEstimate, values...)
	if err != nil {
		if config.Debug() {
			log.Printf("SELECT FROM %s exited; query cache miss; failed query in %s", table, time.Since(before))
		}

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
		if config.Debug() {
			log.Printf("SELECT FROM %s exited; query cache miss; failed query in %s", table, time.Since(before))
		}

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
			if config.Debug() {
				log.Printf("SELECT FROM %s exited; query cache miss; failed query in %s", table, time.Since(before))
			}

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
		if config.Debug() {
			log.Printf("SELECT FROM %s exited; query cache miss; failed query in %s", table, time.Since(before))
		}

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

	if config.Debug() {
		log.Printf("SELECT FROM %s exited; query cache miss; successful query in %s", table, time.Since(before))
	}

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
		FormatTableName(table),
		JoinObjectNames(FormatObjectNames(columns)),
		strings.Join(placeholders, ",\n    "),
		JoinObjectNames(FormatObjectNames(returning, true)),
	))

	if config.QueryDebug() {
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

		if config.Debug() {
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
		FormatTableName(table),
		JoinObjectNames(FormatObjectNames(columns)),
		strings.Join(placeholders, ",\n    "),
		GetWhere(where),
		JoinObjectNames(FormatObjectNames(returning)),
	))

	if config.QueryDebug() {
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
		FormatTableName(table),
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

func LockTable(ctx context.Context, tx pgx.Tx, tableName string, timeouts ...time.Duration) error {
	var timeout *time.Duration
	if len(timeouts) > 0 {
		timeout = &timeouts[0]
	}

	noWaitInfix := ""

	if timeout != nil {
		if *timeout < time.Duration(0) {
			return fmt.Errorf("timeout may not be negative")
		}

		if *timeout == time.Duration(0) {
			noWaitInfix = " NOWAIT"
		} else {
			_, setErr := tx.Exec(ctx, fmt.Sprintf("SET LOCAL statement_timeout = '%fs';", (*timeout).Seconds()))
			if setErr != nil {
				return setErr
			}

			defer func() {
				_, unsetErr := tx.Exec(ctx, "RESET statement_timeout;")
				if unsetErr != nil {
					log.Printf("warning: failed to clear statement timeout: %v", unsetErr)
				}
			}()
		}
	}

	rawSavePointID, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	savePointID := fmt.Sprintf("savepoint_%s", strings.ReplaceAll(rawSavePointID.String(), "-", ""))

	_, err = tx.Exec(ctx, fmt.Sprintf("SAVEPOINT %s;", savePointID))
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, fmt.Sprintf("LOCK TABLE %v IN ACCESS EXCLUSIVE MODE%s;", tableName, noWaitInfix))
	if err != nil {
		_, rollbackErr := tx.Exec(ctx, fmt.Sprintf("ROLLBACK TO SAVEPOINT %s;", savePointID))
		if rollbackErr != nil {
			return fmt.Errorf("lock table attempt returned %v, subsequent rollback attempt returned %v", err, rollbackErr)
		}

		return err
	}

	return nil
}

func LockTableWithRetries(ctx context.Context, tx pgx.Tx, tableName string, overallTimeout time.Duration, individualAttemptTimeouts ...time.Duration) error {
	individualAttemptTimeout := overallTimeout / 10
	if len(individualAttemptTimeouts) > 0 {
		individualAttemptTimeout = individualAttemptTimeouts[0]
	}

	if overallTimeout <= individualAttemptTimeout {
		return fmt.Errorf("overallTimeout %s may not be less than individualAttemptTimeout %s", overallTimeout, individualAttemptTimeout)
	}

	if overallTimeout == time.Duration(0) {
		return LockTable(ctx, tx, tableName, overallTimeout)
	}

	expiry := time.Now().Add(overallTimeout)

	for time.Now().Before(expiry) {
		ok, err := func() (bool, error) {
			err := LockTable(ctx, tx, tableName, individualAttemptTimeout)
			if err != nil {
				if strings.Contains(err.Error(), "canceling statement due to statement timeout") {
					return false, nil
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

func AdvisoryLock(ctx context.Context, tx pgx.Tx, key1 int32, key2 int32, timeouts ...time.Duration) error {
	var timeout *time.Duration
	if len(timeouts) > 0 {
		timeout = &timeouts[0]
	}

	tryPrefix := ""

	if timeout != nil {
		if *timeout < time.Duration(0) {
			return fmt.Errorf("timeout may not be negative")
		}

		if *timeout == time.Duration(0) {
			tryPrefix = "try_"
		} else {
			_, setErr := tx.Exec(ctx, fmt.Sprintf("SET LOCAL statement_timeout = '%fs';", (*timeout).Seconds()))
			if setErr != nil {
				return setErr
			}

			defer func() {
				_, unsetErr := tx.Exec(ctx, "RESET statement_timeout;")
				if unsetErr != nil {
					log.Printf("warning: failed to clear statement timeout: %v", unsetErr)
				}
			}()
		}
	}

	rawSavePointID, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	savePointID := fmt.Sprintf("savepoint_%s", strings.ReplaceAll(rawSavePointID.String(), "-", ""))

	_, err = tx.Exec(ctx, fmt.Sprintf("SAVEPOINT %s;", savePointID))
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, fmt.Sprintf("SELECT pg_%sadvisory_xact_lock(%d, %d);", tryPrefix, key1, key2))
	if err != nil {
		_, rollbackErr := tx.Exec(ctx, fmt.Sprintf("ROLLBACK TO SAVEPOINT %s;", savePointID))
		if rollbackErr != nil {
			return fmt.Errorf("advisory lock attempt returned %v, subsequent rollback attempt returned %v", err, rollbackErr)
		}

		return err
	}

	return nil
}

func AdvisoryLockWithRetries(ctx context.Context, tx pgx.Tx, key1 int32, key2 int32, overallTimeout time.Duration, individualAttemptTimeouts ...time.Duration) error {
	individualAttemptTimeout := overallTimeout / 10
	if len(individualAttemptTimeouts) > 0 {
		individualAttemptTimeout = individualAttemptTimeouts[0]
	}

	if overallTimeout <= individualAttemptTimeout {
		return fmt.Errorf("overallTimeout %s may not be less than individualAttemptTimeout %s", overallTimeout, individualAttemptTimeout)
	}

	if overallTimeout == time.Duration(0) {
		return AdvisoryLock(ctx, tx, key1, key2, overallTimeout)
	}

	expiry := time.Now().Add(overallTimeout)

	for time.Now().Before(expiry) {
		ok, err := func() (bool, error) {
			err := AdvisoryLock(ctx, tx, key1, key2, individualAttemptTimeout)
			if err != nil {
				if strings.Contains(err.Error(), "canceling statement due to statement timeout") {
					return false, nil
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

	return fmt.Errorf("timed out waiting for advisory lock (%#+v, %#+v) after %s ", key1, key2, overallTimeout)
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
