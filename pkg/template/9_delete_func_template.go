package template

import "strings"

var genericDeleteFuncTemplate = strings.TrimSpace(`
func genericDelete%v(ctx context.Context, db *sqlx.DB, object DjangolangObject) error {
	if object == nil {
		return fmt.Errorf("object given for deletion was unexpectedly nil")
	}

	err := object.Delete(ctx, db)
	if err != nil {
		return err
	}

	return nil
}
`) + "\n\n"

var deleteFuncTemplate = strings.TrimSpace(`
func (%v *%v) Delete(ctx context.Context, db *sqlx.DB) error {
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
				"deleted %%v row(s); %%.3f seconds to build, %%.3f seconds to execute; sql:\n%%v\n\n",
				rowCount, buildDuration, execDuration, sql,
			)
		}
	}()

	buildStart = time.Now().UnixNano()

	buildStop = time.Now().UnixNano()

	execStart = time.Now().UnixNano()

	deleteCtx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()

	sql = fmt.Sprintf(
		"DELETE FROM %v WHERE %v = %%v",
		%v.%v,
	)

	result, err := db.ExecContext(
		deleteCtx,
		sql,
	)
	if err != nil {
		return err
	}

	rowCount, err = result.RowsAffected()
	if err != nil {
		return err
	}

	execStop = time.Now().UnixNano()

	if rowCount != 1 {
		return fmt.Errorf("expected exactly 1 affected row after deleting %%#+v; got %%v", %v, rowCount)
	}

	return nil
}
`) + "\n\n"

var deleteFuncTemplateNotImplemented = strings.TrimSpace(`
func (%v *%v) Delete(ctx context.Context, db *sqlx.DB) error {
	return fmt.Errorf("not implemented (table has no primary key)")
}`) + "\n\n"
