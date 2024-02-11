package generate

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/chanced/caps"
	"github.com/initialed85/djangolang/pkg/helpers"
	"github.com/initialed85/djangolang/pkg/introspect"
	"golang.org/x/exp/maps"
)

var (
	logger = helpers.GetLogger("generate")
)

var fileDataHeaderTemplate = strings.TrimSpace(`
package djangolang

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/initialed85/djangolang/pkg/helpers"
	"github.com/jmoiron/sqlx"
)

type SelectFunc = func(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]any, error)
type Handler = func(w http.ResponseWriter, r *http.Request)
type SelectHandler = Handler

var (
	logger                = helpers.GetLogger("djangolang")
	mu                    = new(sync.RWMutex)
	actualDebug           = false
	selectFuncByTableName = make(map[string]SelectFunc)
	columnsByTableName = make(map[string][]string)
)

func init() {
	mu.Lock()
	defer mu.Unlock()
	rawDesiredDebug := os.Getenv("DJANGOLANG_DEBUG")
	actualDebug = rawDesiredDebug == "1"
	logger.Printf("DJANGOLANG_DEBUG=%%v, debugging enabled: %%v", rawDesiredDebug, actualDebug)

	%v
}

func SetDebug(desiredDebug bool) {
	mu.Lock()
	defer mu.Unlock()
	actualDebug = desiredDebug
	logger.Printf("runtime SetDebug() called, debugging enabled: %%v", actualDebug)
}

func Descending(columns ...string) *string {
	return helpers.Ptr(
		fmt.Sprintf(
			"(%%v) DESC",
			strings.Join(columns, ", "),
		),
	)
}

func Ascending(columns ...string) *string {
	return helpers.Ptr(
		fmt.Sprintf(
			"(%%v) ASC",
			strings.Join(columns, ", "),
		),
	)
}

func Columns(includeColumns []string, excludeColumns ...string) []string {
	excludeColumnLookup := make(map[string]bool)
	for _, column := range excludeColumns {
		excludeColumnLookup[column] = true
	}

	columns := make([]string, 0)
	for _, column := range includeColumns {
		_, ok := excludeColumnLookup[column]
		if ok {
			continue
		}

		columns = append(columns, column)
	}

	return columns
}

func GetSelectFuncByTableName() map[string]SelectFunc {
	thisSelectFuncByTableName := make(map[string]SelectFunc)

	for tableName, selectFunc := range selectFuncByTableName {
		thisSelectFuncByTableName[tableName] = selectFunc
	}

	return thisSelectFuncByTableName
}

func GetSelectHandlerForTableName(tableName string, db *sqlx.DB) (SelectHandler, error) {
	selectFunc, ok := selectFuncByTableName[tableName]
	if !ok {
		return nil, fmt.Errorf("no selectFuncByTableName entry for tableName %%v", tableName)
	}

	columns, ok := columnsByTableName[tableName]
	if !ok {
		return nil, fmt.Errorf("no columnsByTableName entry for tableName %%v", tableName)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		status := http.StatusOK
		var body []byte
		var err error

		defer func() {
			w.Header().Add("Content-Type", "application/json")

			w.WriteHeader(status)

			if err != nil {
				_, _ = w.Write([]byte(fmt.Sprintf("{\"error\": %%#+v}", err.Error())))
				return
			}

			_, _ = w.Write(body)
		}()

		if r.Method != http.MethodGet {
			status = http.StatusMethodNotAllowed
			err = fmt.Errorf("%%v; wanted %%v, got %%v", "StatusMethodNotAllowed", http.MethodGet, r.Method)
			return
		}

		limit := helpers.Ptr(50)
		rawLimit := strings.TrimSpace(r.URL.Query().Get("limit"))
		if len(rawLimit) > 0 {
			var parsedLimit int64
			parsedLimit, err = strconv.ParseInt(rawLimit, 10, 64)
			if err != nil {
				status = http.StatusBadRequest
				err = fmt.Errorf("failed to parse limit %%#+v as int; err: %%v", limit, err.Error())
				return
			}
			limit = helpers.Ptr(int(parsedLimit))
		}

		offset := helpers.Ptr(50)
		rawOffset := strings.TrimSpace(r.URL.Query().Get("offset"))
		if len(rawOffset) > 0 {
			var parsedOffset int64
			parsedOffset, err = strconv.ParseInt(rawOffset, 10, 64)
			if err != nil {
				status = http.StatusBadRequest
				err = fmt.Errorf("failed to parse offset %%#+v as int; err: %%v", offset, err.Error())
				return
			}
			offset = helpers.Ptr(int(parsedOffset))
		}

		var order *string

		wheres := make([]string, 0)

		var items []any
		items, err = selectFunc(r.Context(), db, columns, order, limit, offset, wheres...)
		if err != nil {
			status = http.StatusInternalServerError
			err = fmt.Errorf("query failed; err: %%v", err.Error())
			return
		}

		body, err = json.Marshal(items)
		if err != nil {
			status = http.StatusInternalServerError
			err = fmt.Errorf("serialization failed; err: %%v", err.Error())
		}
	}, nil
}

func GetSelectHandlerByEndpointName(db *sqlx.DB) (map[string]SelectHandler, error) {
	thisSelectHandlerByEndpointName := make(map[string]SelectHandler)

	for _, tableName := range TableNames {
		endpointName := caps.ToKebab[string](tableName)

		selectHandler, err := GetSelectHandlerForTableName(tableName, db)
		if err != nil {
			return nil, err
		}

		thisSelectHandlerByEndpointName[endpointName] = selectHandler
	}

	return thisSelectHandlerByEndpointName, nil
}
`) + "\n\n"

var selectFuncTemplate = strings.TrimSpace(`
func Select%v(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres...string) ([]*%v, error) {
	mu.RLock()
	debug := actualDebug
	mu.RUnlock()

	var buildStart int64
	var buildStop int64
	var execStart int64
	var execStop int64
	var scanStart int64
	var scanStop int64

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

		if buildStop > 0 {
			buildDuration = float64(buildStop-buildStart) * 1e-9
		}

		if execStop > 0 {
			execDuration = float64(execStop-execStart) * 1e-9
		}

		if scanStop > 0 {
			scanDuration = float64(scanStop-scanStart) * 1e-9
		}

		if debug {
			logger.Printf(
				"selected %%v columns, %%v rows; %%.3f seconds to build, %%.3f seconds to execute, %%.3f seconds to scan; sql:\n%%v",
				columnCount, rowCount, buildDuration, execDuration, scanDuration, sql,
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

	items := make([]*%v, 0)
	for rows.Next() {
		rowCount++

		var item %v
		err = rows.StructScan(&item)
		if err != nil {
			return nil, err
		}

		%v

		items = append(items, &item)
	}

	scanStop = time.Now().UnixNano()

	return items, nil
}

func genericSelect%v(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres...string) ([]any, error) {
	items, err := Select%v(ctx, db, columns, orderBy, limit, offset, wheres...)
	if err != nil {
		return nil, err
	}

	genericItems := make([]any, 0)
	for _, item := range items {
		genericItems = append(genericItems, item)
	}

	return genericItems, nil
}
`) + "\n\n"

var jsonBlockTemplate = strings.Trim(`
        var temp1%v any
		var temp2%v []any

		err = json.Unmarshal(item.%v.([]byte), &temp1%v)
        if err != nil {
        	err = json.Unmarshal(item.%v.([]byte), &temp2%v)
			if err != nil {
				return nil, fmt.Errorf("failed to unmarshal %%#+v as json: %%v", string(item.%v.([]byte)), err)
			} else {
				item.%v = temp2%v
			}
        } else {
			item.%v = temp1%v
		}

		item.%v = temp1%v
`, "\n") + "\n\n"

func init() {
	converter, _ := caps.DefaultConverter.(caps.StdConverter)
	converter.Set("Dob", "DOB")
	converter.Set("Cpo", "CPO")
	converter.Set("Kwh", "KWH")
}

func Run(ctx context.Context) error {
	relativeOutputPath := "./pkg/djangolang"
	if len(os.Args) > 2 {
		relativeOutputPath = os.Args[2]
	}

	logger.Printf("resolving %v...", relativeOutputPath)
	outputPath, err := filepath.Abs(relativeOutputPath)
	if err != nil {
		return err
	}
	logger.Printf("done.")

	logger.Printf("cleaning out %v...", outputPath)
	_ = os.RemoveAll(outputPath)
	err = os.MkdirAll(outputPath, 0o777)
	if err != nil {
		return err
	}
	logger.Printf("done.")

	conn, err := helpers.GetDBFromEnvironment(ctx)
	if err != nil {
		return err
	}
	defer func() {
		_ = conn.Close()
	}()

	schema := helpers.GetSchema()

	tableByName, err := introspect.Introspect(ctx, conn, schema)
	if err != nil {
		return err
	}

	tableNames := maps.Keys(tableByName)
	slices.Sort(tableNames)

	fileData := ""

	initData := ""

	tableData := fmt.Sprintf("var TableNames = %#+v\n\n", tableNames)

	fileData += tableData

	for _, tableName := range tableNames {
		table := tableByName[tableName]

		structName := caps.ToCamel(table.Name)

		initData += fmt.Sprintf(
			"    selectFuncByTableName[%#+v] = genericSelect%v\n",
			table.Name,
			structName,
		)

		initData += fmt.Sprintf(
			"    columnsByTableName[%#+v] = %vColumns\n",
			table.Name,
			structName,
		)

		nameData := fmt.Sprintf(
			"var %vTable = \"%v\"\n",
			structName,
			table.Name,
		)

		structData := fmt.Sprintf(
			"type %v struct {\n",
			structName,
		)

		postprocessData := ""

		columnNames := make([]string, 0)

		for _, column := range table.Columns {
			if column.ZeroType == nil && !(column.TypeTemplate == "any" || column.TypeTemplate == "*any") {
				continue
			}

			fieldName := caps.ToCamel(column.Name)

			columnNames = append(columnNames, column.Name)

			nameData += fmt.Sprintf(
				"var %v%vColumn = \"%v\"\n",
				structName,
				fieldName,
				column.Name,
			)

			structData += fmt.Sprintf(
				"    %v %v `json:\"%v\" db:\"%v\"`\n",
				fieldName,
				column.TypeTemplate,
				column.Name,
				column.Name,
			)

			if column.ZeroType == nil && (column.TypeTemplate == "any" || column.TypeTemplate == "*any") {
				postprocessData += fmt.Sprintf(
					jsonBlockTemplate,
					fieldName,
					fieldName,
					fieldName,
					fieldName,
					fieldName,
					fieldName,
					fieldName,
					fieldName,
					fieldName,
					fieldName,
					fieldName,
					fieldName,
					fieldName,
				)
			}
		}

		if postprocessData == "" {
			postprocessData = "        // this is where any post-scan processing would appear (if required)"
		}

		nameData += fmt.Sprintf("var %vColumns = %#+v\n", structName, columnNames) + "\n"

		selectData := fmt.Sprintf(
			selectFuncTemplate,
			structName,
			structName,
			table.Name,
			structName,
			structName,
			postprocessData,
			structName,
			structName,
		)

		structData += "}\n\n"

		fileData += nameData
		fileData += structData
		fileData += selectData
	}

	header := fmt.Sprintf(fileDataHeaderTemplate, strings.TrimSpace(initData))

	finalFileData := header + fileData

	fullOutputPath := filepath.Join(outputPath, "models.go")

	err = os.WriteFile(
		fullOutputPath,
		[]byte(finalFileData),
		0o777,
	)
	if err != nil {
		return err
	}

	command1Ctx, command1Cancel := context.WithTimeout(ctx, time.Second*10)
	defer command1Cancel()
	out, err := exec.CommandContext(command1Ctx, "goimports", "-w", fullOutputPath).CombinedOutput()
	if err != nil {
		logger.Printf("error: goimports failed; err: %v; output follows:\n\n%v", err, string(out))
		return err
	}

	command2Ctx, command2Cancel := context.WithTimeout(ctx, time.Second*10)
	defer command2Cancel()
	out, err = exec.CommandContext(command2Ctx, "gofmt", fullOutputPath).CombinedOutput()
	if err != nil {
		logger.Printf("error: gofmt failed; err: %v; output follows:\n\n%v", err, string(out))
		return err
	}

	command3Ctx, command3Cancel := context.WithTimeout(ctx, time.Second*10)
	defer command3Cancel()
	out, err = exec.CommandContext(command3Ctx, "go", "vet", fullOutputPath).CombinedOutput()
	if err != nil {
		logger.Printf("error: go vet failed; err: %v; output follows:\n\n%v", err, string(out))
		return err
	}

	return nil
}
