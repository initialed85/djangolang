package templates

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/chanced/caps"
	"github.com/initialed85/djangolang/pkg/helpers"
	"github.com/jmoiron/sqlx"
)

type SelectFunc = func(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]any, error)
type Handler = func(w http.ResponseWriter, r *http.Request)
type SelectHandler = Handler

var (
	logger                 = helpers.GetLogger("djangolang")
	mu                     = new(sync.RWMutex)
	actualDebug            = false
	selectFuncByTableName  = make(map[string]SelectFunc)
	columnNamesByTableName = make(map[string][]string)
)

func init() {
	mu.Lock()
	defer mu.Unlock()
	rawDesiredDebug := os.Getenv("DJANGOLANG_DEBUG")
	actualDebug = rawDesiredDebug == "1"
	logger.Printf("DJANGOLANG_DEBUG=%v, debugging enabled: %v", rawDesiredDebug, actualDebug)

	selectFuncByTableName["firsts"] = genericSelectFirsts
	columnNamesByTableName["firsts"] = FirstsColumns
}

func SetDebug(desiredDebug bool) {
	mu.Lock()
	defer mu.Unlock()
	actualDebug = desiredDebug
	logger.Printf("runtime SetDebug() called, debugging enabled: %v", actualDebug)
}

func Descending(columns ...string) *string {
	return helpers.Ptr(
		fmt.Sprintf(
			"(%v) DESC",
			strings.Join(columns, ", "),
		),
	)
}

func Ascending(columns ...string) *string {
	return helpers.Ptr(
		fmt.Sprintf(
			"(%v) ASC",
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
		return nil, fmt.Errorf("no selectFuncByTableName entry for tableName %v", tableName)
	}

	columns, ok := columnNamesByTableName[tableName]
	if !ok {
		return nil, fmt.Errorf("no columnNamesByTableName entry for tableName %v", tableName)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		status := http.StatusOK
		var body []byte
		var err error

		defer func() {
			w.Header().Add("Content-Type", "application/json")

			w.WriteHeader(status)

			if err != nil {
				_, _ = w.Write([]byte(fmt.Sprintf("{\"error\": %#+v}", err.Error())))
				return
			}

			_, _ = w.Write(body)
		}()

		if r.Method != http.MethodGet {
			status = http.StatusMethodNotAllowed
			err = fmt.Errorf("%v; wanted %v, got %v", "StatusMethodNotAllowed", http.MethodGet, r.Method)
			return
		}

		limit := helpers.Ptr(50)
		rawLimit := strings.TrimSpace(r.URL.Query().Get("limit"))
		if len(rawLimit) > 0 {
			var parsedLimit int64
			parsedLimit, err = strconv.ParseInt(rawLimit, 10, 64)
			if err != nil {
				status = http.StatusBadRequest
				err = fmt.Errorf("failed to parse limit %#+v as int; err: %v", limit, err.Error())
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
				err = fmt.Errorf("failed to parse offset %#+v as int; err: %v", offset, err.Error())
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
			err = fmt.Errorf("query failed; err: %v", err.Error())
			return
		}

		body, err = json.Marshal(items)
		if err != nil {
			status = http.StatusInternalServerError
			err = fmt.Errorf("serialization failed; err: %v", err.Error())
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

var TableNames = []string{"firsts", "firsts_thirds", "seconds", "seconds_thirds", "thirds"}
