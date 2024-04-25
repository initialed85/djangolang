package types

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
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
