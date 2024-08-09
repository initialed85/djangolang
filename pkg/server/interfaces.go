package server

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
	"github.com/initialed85/djangolang/pkg/stream"
	"github.com/jmoiron/sqlx"
)

type HTTPMiddleware func(http.Handler) http.Handler
type ObjectMiddleware func()
type GetRouterFn func(*sqlx.DB, *redis.Pool, []HTTPMiddleware, []ObjectMiddleware, WaitForChange) chi.Router
type WaitForChange func(context.Context, []stream.Action, string, uint32) (*Change, error)

type WithReload interface {
	Reload(context.Context, *sqlx.Tx, ...bool) error
}

type WithPrimaryKey interface {
	GetPrimaryKeyColumn() string
	GetPrimaryKeyValue() any
}

type WithInsert interface {
	Insert(context.Context, sqlx.Tx) error
}

type Change struct {
	Timestamp time.Time      `json:"timestamp"`
	ID        uuid.UUID      `json:"id"`
	Action    stream.Action  `json:"action"`
	TableName string         `json:"table_name"`
	Item      map[string]any `json:"-"`
	Object    any            `json:"object"`
	Xid       uint32         `json:"xid"`
}

func (c *Change) String() string {
	primaryKeySummary := ""
	if c.Object != nil {
		object, ok := c.Object.(WithPrimaryKey)
		if ok {
			primaryKeySummary = fmt.Sprintf(
				"(%s = %s) ",
				object.GetPrimaryKeyColumn(),
				object.GetPrimaryKeyValue(),
			)
		}
	}

	return strings.TrimSpace(fmt.Sprintf(
		"(%s / %d) @ %s; %s %s %s",
		c.ID, c.Xid, c.Timestamp, c.Action, c.TableName, primaryKeySummary,
	))
}
