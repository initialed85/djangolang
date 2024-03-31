package template

import (
	"strings"
)

var genericInsertFuncTemplate = strings.TrimSpace(`
func genericInsert%v(ctx context.Context, db *sqlx.DB, object DjangolangObject, columns ...string) (DjangolangObject, error) {
	if object == nil {
		return nil, fmt.Errorf("object given for insertion was unexpectedly nil")
	}

	err := object.Insert(ctx, db, columns...)
	if err != nil {
		return nil, err
	}

	return object, nil
}
`) + "\n\n"

var insertFuncTemplate = strings.TrimSpace(`
func (%v *%v) Insert(ctx context.Context, db *sqlx.DB, columns ...string) error {
	if len(columns) > 1 {
		return fmt.Errorf("assertion failed: 'columns' variadic argument(s) must be missing or singular; got %%v", len(columns))
	}

	if len(columns) == 0 {
		columns = %vInsertColumns
	}

	mu.RLock()
	debug := actualDebug
	mu.RUnlock()

	var buildStart int64
	var buildStop int64
	var execStart int64
	var execStop int64

	var sql string
	var rowCount int64

	defer func() {
		if !debug {
			return
		}

		buildDuration := 0.0
		execDuration := 0.0

		if buildStop > 0 {
			buildDuration = float64(buildStop-buildStart) * 1e-9
		}

		if execStop > 0 {
			execDuration = float64(execStop-execStart) * 1e-9
		}

		if debug {
			logger.Printf(
				"inserted %%v rows; %%.3f seconds to build, %%.3f seconds to execute; sql:\n%%v",
				rowCount, buildDuration, execDuration, sql,
			)
		}
	}()

	buildStart = time.Now().UnixNano()

	buildStop = time.Now().UnixNano()

	execStart = time.Now().UnixNano()

	insertCtx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()

	names := make([]string, 0)
	for _, column := range columns {
		names = append(names, fmt.Sprintf(":%%v", column))
	}

	result, err := db.NamedQueryContext(
		insertCtx,
		fmt.Sprintf(
			"INSERT INTO %v (%%v) VALUES (%%v) RETURNING %%v",
			strings.Join(columns, ", "),
			strings.Join(names, ", "),
			strings.Join(%vColumns, ", "),
		),
		%v,
	)
	if err != nil {
		return err
	}

	_ = result.Next()

	err = result.StructScan(%v)
	if err != nil {
		return err
	}

	execStop = time.Now().UnixNano()

	return nil
}
`) + "\n\n"

var insertFuncTemplateNotImplemented = strings.TrimSpace(`
func (%v *%v) Insert(ctx context.Context, db *sqlx.DB, columns ...string) error {
	return fmt.Errorf("not implemented (table has no primary key)")
}
`) + "\n\n"
