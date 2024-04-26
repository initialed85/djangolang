package types

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/initialed85/djangolang/pkg/helpers"
	"github.com/jmoiron/sqlx"
)

type Fragment struct {
	SQL    string
	Values []any
}

func Clause(clause string, value any) Fragment {
	return Fragment{
		SQL:    clause,
		Values: []any{value},
	}
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

type DjangolangObject interface {
	GetPrimaryKey() (any, error)
	SetPrimaryKey(value any) error
	Insert(ctx context.Context, db *sqlx.DB, columns ...string) error
	Update(ctx context.Context, db *sqlx.DB, columns ...string) error
	Delete(ctx context.Context, db *sqlx.DB) error
}

type SelectFunc = func(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...Fragment) ([]DjangolangObject, error)
type MutateFunc = func(ctx context.Context, db *sqlx.DB, object DjangolangObject, columns ...string) (DjangolangObject, error)
type InsertFunc = MutateFunc
type UpdateFunc = MutateFunc
type DeleteFunc = func(ctx context.Context, db *sqlx.DB, object DjangolangObject) error

type DeserializeFunc = func(b []byte) (DjangolangObject, error)
type ConvertFunc = func(map[string]any) (DjangolangObject, error)

type Handler = func(w http.ResponseWriter, r *http.Request)
type SelectHandler = Handler
type InsertHandler = Handler
type UpdateHandler = Handler
type DeleteHandler = Handler

type Action string

const (
	INSERT   Action = "INSERT"
	UPDATE   Action = "UPDATE"
	DELETE   Action = "DELETE"
	TRUNCATE Action = "TRUNCATE"
)

type Change struct {
	ID               uuid.UUID `json:"id"`
	Action           Action    `json:"action"`
	TableName        string    `json:"table_name"`
	PrimaryKeyColumn string    `json:"primary_key_column"`
	PrimaryKeyValue  any       `json:"primary_key_value"`
	Object           any       `json:"object"`
}

func (c *Change) String() string {
	b, _ := json.Marshal(c.PrimaryKeyValue)
	primaryKeyValue := string(b)

	object, _ := json.Marshal(c.Object)

	return fmt.Sprintf(
		"%v %v (%v: %v): %v",
		c.Action, c.TableName, c.PrimaryKeyColumn, primaryKeyValue, string(object),
	)
}
