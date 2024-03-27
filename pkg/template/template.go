package template

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
	_pluralize "github.com/gertd/go-pluralize"
	_helpers "github.com/initialed85/djangolang/internal/helpers"
	"github.com/initialed85/djangolang/pkg/helpers"
	"github.com/initialed85/djangolang/pkg/introspect"
	"golang.org/x/exp/maps"
)

var (
	logger    = helpers.GetLogger("generate")
	pluralize = _pluralize.NewClient()
)

var headerTemplate = strings.TrimSpace(`
package templates

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/initialed85/djangolang/pkg/helpers"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"golang.org/x/exp/maps"
)
`) + "\n\n"

var fileDataHeaderTemplate = strings.TrimSpace(`
package templates

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/chanced/caps"
	"github.com/initialed85/djangolang/pkg/helpers"
	"github.com/initialed85/djangolang/pkg/introspect"
	"github.com/jmoiron/sqlx"
)

type SelectFunc = func(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]any, error)
type Handler = func(w http.ResponseWriter, r *http.Request)
type SelectHandler = Handler

var (
	logger                            = helpers.GetLogger("djangolang")
	mu                                = new(sync.RWMutex)
	actualDebug                       = false
	selectFuncByTableName             = make(map[string]SelectFunc)
	columnNamesByTableName            = make(map[string][]string)
	transformedColumnNamesByTableName = make(map[string][]string)
	tableByName                       = make(map[string]*introspect.Table)
)

var rawTableByName = []byte(%v)

func init() {
	mu.Lock()
	defer mu.Unlock()
	rawDesiredDebug := os.Getenv("DJANGOLANG_DEBUG")
	actualDebug = rawDesiredDebug == "1"
	logger.Printf("DJANGOLANG_DEBUG=%%v, debugging enabled: %%v", rawDesiredDebug, actualDebug)

	var err error
	tableByName, err = GetTableByName()
	if err != nil {
		log.Panic(err)
	}

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

func GetRawTableByName() []byte {
	return rawTableByName
}

func GetTableByName() (map[string]*introspect.Table, error) {
	thisTableByName := make(map[string]*introspect.Table)

	err := json.Unmarshal(rawTableByName, &thisTableByName)
	if err != nil {
		return nil, err
	}

	thisTableByName, err = introspect.MapTableByName(thisTableByName)
	if err != nil {
		return nil, err
	}

	return thisTableByName, nil
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

	columns, ok := transformedColumnNamesByTableName[tableName]
	if !ok {
		return nil, fmt.Errorf("no columnNamesByTableName entry for tableName %%v", tableName)
	}

	table, ok := tableByName[tableName]
	if !ok {
		return nil, fmt.Errorf("no tableByName entry for tableName %%v", tableName)
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

		descending := strings.HasPrefix(strings.ToLower(strings.TrimSpace(r.URL.Query().Get("order"))), "desc")

		orderBy := make([]string, 0)
		for k, vs := range r.URL.Query() {
			if k != "order_by" {
				continue
			}

			for _, v := range vs {
				_, ok := table.ColumnByName[v]
				if !ok {
					status = http.StatusBadRequest
					err = fmt.Errorf("bad order by; no table.ColumnByName entry for columnName %%v", v)
					return
				}

				found := false

				for _, existingOrderBy := range orderBy {
					if existingOrderBy == v {
						found = true
						break
					}
				}

				if found {
					continue
				}

				orderBy = append(orderBy, v)
			}
		}

		var order *string
		if len(orderBy) > 0 {
			if descending {
				order = Descending(orderBy...)
			} else {
				order = Ascending(orderBy...)
			}
		}

		wheres := make([]string, 0)
		for k, vs := range r.URL.Query() {
			if k == "order" || k == "order_by" || k == "limit" || k == "offset" {
				continue
			}

			field := k
			matcher := "="

			parts := strings.Split(k, "__")
			if len(parts) > 1 {
				field = strings.Join(parts[0:len(parts)-1], "__")
				lastPart := strings.ToLower(parts[len(parts)-1])
				if lastPart == "eq" {
					matcher = "="
				} else if lastPart == "ne" {
					matcher = "!="
				} else if lastPart == "gt" {
					matcher = ">"
				} else if lastPart == "lt" {
					matcher = "<"
				} else if lastPart == "gte" {
					matcher = ">="
				} else if lastPart == "lte" {
					matcher = "<="
				} else if lastPart == "ilike" {
					matcher = "ILIKE"
				} else if lastPart == "in" {
					matcher = "IN"
				} else if lastPart == "nin" || lastPart == "not_in" {
					matcher = "NOT IN"
				}
			}

			innerWheres := make([]string, 0)

			for _, v := range vs {
				column, ok := table.ColumnByName[field]
				if !ok {
					status = http.StatusBadRequest
					err = fmt.Errorf("bad filter; no table.ColumnByName entry for columnName %%v", field)
					return
				}

				if matcher != "IN" && matcher != "NOT IN" {
					if column.TypeTemplate == "string" || column.TypeTemplate == "uuid" {
						v = fmt.Sprintf("%%#+v", v)
						innerWheres = append(
							innerWheres,
							fmt.Sprintf("%%v %%v '%%v'", field, matcher, v[1:len(v)-1]),
						)
					} else {
						innerWheres = append(
							innerWheres,
							fmt.Sprintf("%%v %%v %%v", field, matcher, v),
						)
					}
				} else {
					itemsBefore := strings.Split(v, ",")
					itemsAfter := make([]string, 0)

					for _, x := range itemsBefore {
						x = strings.TrimSpace(x)
						if x == "" {
							continue
						}

						if column.TypeTemplate == "string" || column.TypeTemplate == "uuid" {
							x = fmt.Sprintf("%%#+v", x)
							itemsAfter = append(itemsAfter, fmt.Sprintf("'%%v'", x[1:len(x)-1]))
						} else {
							itemsAfter = append(itemsAfter, x)
						}
					}

					innerWheres = append(
						innerWheres,
						fmt.Sprintf("%%v %%v (%%v)", field, matcher, strings.Join(itemsAfter, ", ")),
					)
				}
			}

			wheres = append(wheres, fmt.Sprintf("(%%v)", strings.Join(innerWheres, " OR ")))
		}

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

%v
`) + "\n\n"

var loadTemplate = strings.TrimSpace(`
	idsFor%v := make([]string, 0)
	for _, id := range maps.Keys(%v) {
		idsFor%v = append(idsFor%v, fmt.Sprintf("%%v", id))
	}

	if len(idsFor%v) > 0 {
		rowsFor%v, err := Select%v(
			ctx,
			db,
			%vTransformedColumns,
			nil,
			nil,
			nil,
			fmt.Sprintf("%v IN (%%v)", strings.Join(idsFor%v, ", ")),
		)
		if err != nil {
			return nil, err
		}

		for _, row := range rowsFor%v {
			%v[row.ID] = row
		}

		for _, item := range items {
			item.%vObject = %v[%v]
		}
	}
`)

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

		if item.%v != nil {
			err = json.Unmarshal(item.%v.([]byte), &temp1%v)
			if err != nil {
				err = json.Unmarshal(item.%v.([]byte), &temp2%v)
				if err != nil {
					item.%v = nil
				} else {
					item.%v = &temp2%v
				}
			} else {
				item.%v = &temp1%v
			}
		}

		item.%v = &temp1%v
`, "\n") + "\n\n"

func init() {
	converter, _ := caps.DefaultConverter.(caps.StdConverter)
	converter.Set("Dob", "DOB")
	converter.Set("Cpo", "CPO")
	converter.Set("Mwh", "MWH")
	converter.Set("Kwh", "KWH")
	converter.Set("Wh", "WH")
	converter.Set("Json", "JSON")
	converter.Set("Jsonb", "JSONB")
	converter.Set("Mac", "MAC")
	converter.Set("Ip", "IP")
}

func Run(ctx context.Context) error {
	relativeOutputPath := "./pkg/template/templates"
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

	for _, tableName := range tableNames {
		tableFileData := ""

		table := tableByName[tableName]

		structNameSingular := pluralize.Singular(caps.ToCamel(table.Name))
		structNamePlural := pluralize.Plural(caps.ToCamel(table.Name))

		if table.RelKind == "v" {
			structNameSingular += "View"
			structNamePlural += "View"
		}

		initData += fmt.Sprintf(
			"    selectFuncByTableName[%#+v] = genericSelect%v\n",
			table.Name,
			structNamePlural,
		)

		initData += fmt.Sprintf(
			"    columnNamesByTableName[%#+v] = %vColumns\n",
			table.Name,
			structNameSingular,
		)

		initData += fmt.Sprintf(
			"    transformedColumnNamesByTableName[%#+v] = %vTransformedColumns\n",
			table.Name,
			structNameSingular,
		)

		nameData := fmt.Sprintf(
			"var %vTable = \"%v\"\n",
			structNameSingular,
			table.Name,
		)

		structData := fmt.Sprintf(
			"type %v struct {\n",
			structNameSingular,
		)

		postprocessData := ""

		columnNames := make([]string, 0)
		transformedColumnNames := make([]string, 0)

		foreignObjectIDMapDeclarations := ""
		foreignObjectIDMapPopulations := ""
		foreignObjectLoading := ""

		for _, column := range table.Columns {
			if column.ZeroType == nil && column.TypeTemplate != "any" {
				continue
			}

			fieldName := caps.ToCamel(column.Name)

			columnNames = append(columnNames, column.Name)

			transformedColumnName := column.Name

			if column.DataType == "point" || column.DataType == "polygon" || column.DataType == "geometry" {
				transformedColumnName = fmt.Sprintf("ST_AsGeoJSON(%v::geometry)::jsonb AS %v", column.Name, column.Name)
			} else if column.DataType == "tsvector" {
				transformedColumnName = fmt.Sprintf("array_to_json(tsvector_to_array(%v))::jsonb AS %v", column.Name, column.Name)
			} else if column.DataType == "interval" {
				transformedColumnName = fmt.Sprintf("extract(microseconds FROM %v)::numeric * 1000 AS %v", column.Name, column.Name)
			}

			transformedColumnNames = append(transformedColumnNames, transformedColumnName)

			nameData += fmt.Sprintf(
				"var %v%vColumn = \"%v\"\n",
				structNameSingular,
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

			if column.ForeignColumn != nil {
				foreignObjectNameSingular := pluralize.Singular(caps.ToCamel(column.ForeignTable.Name))
				foreignObjectNamePlural := pluralize.Plural(caps.ToCamel(column.ForeignTable.Name))

				structData += fmt.Sprintf(
					"    %vObject *%v `json:\"%v_object,omitempty\"`\n",
					fieldName,
					foreignObjectNameSingular,
					column.Name,
				)

				foreignObjectMapName := fmt.Sprintf(
					"%vBy%v",
					caps.ToLowerCamel(foreignObjectNameSingular),
					fieldName,
				)

				foreignObjectIDMapDeclarations += fmt.Sprintf(
					"%v := make(map[%v]*%v, 0)\n    ",
					foreignObjectMapName,
					strings.TrimLeft(column.TypeTemplate, "*"),
					foreignObjectNameSingular,
				)

				if strings.HasPrefix(column.TypeTemplate, "*") {
					foreignObjectIDMapPopulations += fmt.Sprintf(
						"%v[helpers.Deref(item.%v)] = nil\n        ",
						foreignObjectMapName,
						fieldName,
					)
				} else {
					foreignObjectIDMapPopulations += fmt.Sprintf(
						"%v[item.%v] = nil\n        ",
						foreignObjectMapName,
						fieldName,
					)
				}

				possiblePointerFieldName := ""
				if strings.HasPrefix(column.TypeTemplate, "*") {
					possiblePointerFieldName = fmt.Sprintf("helpers.Deref(item.%v)", fieldName)
				} else {
					possiblePointerFieldName = fmt.Sprintf("item.%v", fieldName)
				}

				foreignObjectLoading += fmt.Sprintf(
					loadTemplate,
					fieldName,
					foreignObjectMapName,
					fieldName,
					fieldName,
					fieldName,
					fieldName,
					foreignObjectNamePlural,
					foreignObjectNameSingular,
					column.ForeignColumn.Name,
					fieldName,
					fieldName,
					foreignObjectMapName,
					fieldName,
					foreignObjectMapName,
					possiblePointerFieldName,
				) + "\n    \n    "
			}

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
					fieldName,
				)
			}
		}

		if postprocessData == "" {
			postprocessData = "        // this is where any post-scan processing would appear (if required)"
		}

		nameData += fmt.Sprintf("var %vColumns = %#+v\n", structNameSingular, columnNames)
		nameData += fmt.Sprintf("var %vTransformedColumns = %#+v\n", structNameSingular, transformedColumnNames)
		nameData += "\n"

		if foreignObjectIDMapDeclarations == "" {
			foreignObjectIDMapDeclarations = "        // this is where any foreign object ID map declarations would appear (if required)"
		}

		if foreignObjectIDMapPopulations == "" {
			foreignObjectIDMapPopulations = "        // this is where foreign object ID map populations would appear (if required)"
		}

		if foreignObjectLoading == "" {
			foreignObjectLoading = "        // this is where any foreign object loading would appear (if required)"
		}

		selectData := fmt.Sprintf(
			selectFuncTemplate,
			structNamePlural,
			structNameSingular,
			table.Name,
			foreignObjectIDMapDeclarations,
			structNameSingular,
			structNameSingular,
			postprocessData,
			foreignObjectIDMapPopulations,
			foreignObjectLoading,
			structNamePlural,
			structNamePlural,
		)

		structData += "}\n\n"

		tableFileData += nameData
		tableFileData += structData
		tableFileData += selectData

		tableFileData = headerTemplate + tableFileData

		fullOutputPath := filepath.Join(outputPath, fmt.Sprintf("table_%v.go", table.Name))

		err = os.WriteFile(
			fullOutputPath,
			[]byte(tableFileData),
			0o777,
		)
		if err != nil {
			return err
		}
	}

	header := fmt.Sprintf(
		fileDataHeaderTemplate,
		fmt.Sprintf("`\n%v\n`", _helpers.UnsafeJSONPrettyFormat(tableByName)),
		strings.TrimSpace(initData),
		fmt.Sprintf("var TableNames = %#+v\n\n", tableNames),
	)

	fileData = header + fileData

	fullOutputPath := filepath.Join(outputPath, "0_common.go")

	err = os.WriteFile(
		fullOutputPath,
		[]byte(fileData),
		0o777,
	)
	if err != nil {
		return err
	}

	command1Ctx, command1Cancel := context.WithTimeout(ctx, time.Second*10)
	defer command1Cancel()
	out, err := exec.CommandContext(command1Ctx, "goimports", "-w", outputPath).CombinedOutput()
	if err != nil {
		logger.Printf("error: goimports failed; err: %v; output follows:\n\n%v", err, string(out))
		return err
	}

	command2Ctx, command2Cancel := context.WithTimeout(ctx, time.Second*10)
	defer command2Cancel()
	out, err = exec.CommandContext(command2Ctx, "gofmt", outputPath).CombinedOutput()
	if err != nil {
		logger.Printf("error: gofmt failed; err: %v; output follows:\n\n%v", err, string(out))
		return err
	}

	command3Ctx, command3Cancel := context.WithTimeout(ctx, time.Second*10)
	defer command3Cancel()
	out, err = exec.CommandContext(command3Ctx, "go", "vet", outputPath).CombinedOutput()
	if err != nil {
		logger.Printf("error: go vet failed; err: %v; output follows:\n\n%v", err, string(out))
		return err
	}

	return nil
}
