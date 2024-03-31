package template

import "strings"

var genericSelectFuncTemplate = strings.TrimSpace(`
func genericSelect%v(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres...string) ([]DjangolangObject, error) {
	items, err := Select%v(ctx, db, columns, orderBy, limit, offset, wheres...)
	if err != nil {
		return nil, err
	}

	genericItems := make([]DjangolangObject, 0)
	for _, item := range items {
		genericItems = append(genericItems, item)
	}

	return genericItems, nil
}
`) + "\n\n"

var selectFuncTemplate = strings.TrimSpace(`
func Select%v(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres...string) ([]*%v, error) {
	mu.RLock()
	debug := actualDebug
	mu.RUnlock()

	key := "%v"

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
		columns = %vTransformedColumns
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
				"selected %%v columns, %%v rows; %%.3f seconds to build, %%.3f seconds to execute, %%.3f seconds to scan, %%.3f seconds to load foreign objects; sql:\n%%v",
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
		"SELECT %%v FROM %v%%v",
		strings.Join(columns, ", "),
		where,
	)

	if orderBy != nil {
		actualOrderBy := strings.TrimSpace(*orderBy)
		if len(actualOrderBy) > 0 {
			sql += fmt.Sprintf(" ORDER BY %%v", actualOrderBy)
		}
	}

	if limit != nil && *limit >= 0 {
		sql += fmt.Sprintf(" LIMIT %%v", *limit)
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

	%v

	items := make([]*%v, 0)
	for rows.Next() {
		rowCount++

		var item %v
		err = rows.StructScan(&item)
		if err != nil {
			return nil, err
		}

		%v

		%v

		items = append(items, &item)
	}

	scanStop = time.Now().UnixNano()

	foreignObjectStart = time.Now().UnixNano()

	%v

	foreignObjectStop = time.Now().UnixNano()

	return items, nil
}
`) + "\n\n"

// TODO: as required
// var selectFuncTemplateNotImplemented = strings.TrimSpace(`
// func Select%v(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres...string) ([]*%v, error) {
// 	return nil, fmt.Errorf("not implemented (table has no primary key)")
// }
// `) + "\n\n"
