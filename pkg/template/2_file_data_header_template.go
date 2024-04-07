package template

import "strings"

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

type DjangolangObject interface {
	GetPrimaryKey() (any, error)
	SetPrimaryKey(value any) error
	Insert(ctx context.Context, db *sqlx.DB, columns ...string) error
	Update(ctx context.Context, db *sqlx.DB, columns ...string) error
	Delete(ctx context.Context, db *sqlx.DB) error
}

type SelectFunc = func(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]DjangolangObject, error)
type MutateFunc = func(ctx context.Context, db *sqlx.DB, object DjangolangObject, columns ...string) (DjangolangObject, error)
type InsertFunc = MutateFunc
type UpdateFunc = MutateFunc
type DeleteFunc = func(ctx context.Context, db *sqlx.DB, object DjangolangObject) error

type DeserializeFunc = func(b []byte) (DjangolangObject, error)

type Handler = func(w http.ResponseWriter, r *http.Request)
type SelectHandler = Handler
type InsertHandler = Handler
type UpdateHandler = Handler
type DeleteHandler = Handler

var (
	dbName                            = "%v"
	logger                            = helpers.GetLogger(fmt.Sprintf("djangolang/%%v", dbName))
	mu                                = new(sync.RWMutex)
	actualDebug                       = false
	selectFuncByTableName             = make(map[string]SelectFunc)
	insertFuncByTableName             = make(map[string]InsertFunc)
	updateFuncByTableName             = make(map[string]UpdateFunc)
	deleteFuncByTableName             = make(map[string]DeleteFunc)
	deserializeFuncByTableName        = make(map[string]DeserializeFunc)
	columnNamesByTableName            = make(map[string][]string)
	transformedColumnNamesByTableName = make(map[string][]string)
	tableByName                       = make(map[string]*introspect.Table)
)

var rawTableByName = []byte(%v)

%v

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

func GetInsertFuncByTableName() map[string]InsertFunc {
	thisInsertFuncByTableName := make(map[string]InsertFunc)

	for tableName, insertFunc := range insertFuncByTableName {
		thisInsertFuncByTableName[tableName] = insertFunc
	}

	return thisInsertFuncByTableName
}

func GetUpdateFuncByTableName() map[string]UpdateFunc {
	thisUpdateFuncByTableName := make(map[string]UpdateFunc)

	for tableName, updateFunc := range updateFuncByTableName {
		thisUpdateFuncByTableName[tableName] = updateFunc
	}

	return thisUpdateFuncByTableName
}

func GetDeleteFuncByTableName() map[string]DeleteFunc {
	thisDeleteFuncByTableName := make(map[string]DeleteFunc)

	for tableName, deleteFunc := range deleteFuncByTableName {
		thisDeleteFuncByTableName[tableName] = deleteFunc
	}

	return thisDeleteFuncByTableName
}
`) + "\n\n"
