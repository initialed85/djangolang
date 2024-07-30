package server

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/initialed85/djangolang/pkg/stream"
	"github.com/jmoiron/sqlx"
)

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
	ID        uuid.UUID      `json:"id"`
	Action    stream.Action  `json:"action"`
	TableName string         `json:"table_name"`
	Item      map[string]any `json:"-"`
	Object    any            `json:"object"`
}

func (c *Change) String() string {
	b, _ := json.Marshal(c.Object)

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

	return fmt.Sprintf(
		"(%s) %s %s %s%s",
		c.ID, c.Action, c.TableName, primaryKeySummary, string(b),
	)
}
