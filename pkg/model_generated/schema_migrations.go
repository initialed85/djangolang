package model_generated

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/initialed85/djangolang/pkg/helpers"
	"github.com/initialed85/djangolang/pkg/introspect"
	"github.com/initialed85/djangolang/pkg/query"
	"github.com/initialed85/djangolang/pkg/types"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/lib/pq/hstore"
	"github.com/twpayne/go-geom"
)

type SchemaMigration struct {
	Version int64 `json:"version"`
	Dirty   bool  `json:"dirty"`
}

var SchemaMigrationTable = "schema_migrations"

var (
	SchemaMigrationTableVersionColumn = "version"
	SchemaMigrationTableDirtyColumn   = "dirty"
)

var SchemaMigrationTableColumns = []string{
	SchemaMigrationTableVersionColumn,
	SchemaMigrationTableDirtyColumn,
}

var SchemaMigrationTableColumnLookup = map[string]*introspect.Column{
	SchemaMigrationTableVersionColumn: new(introspect.Column),
	SchemaMigrationTableDirtyColumn:   new(introspect.Column),
}

var (
	SchemaMigrationTablePrimaryKeyColumn = SchemaMigrationTableVersionColumn
)

var (
	_ = time.Time{}
	_ = uuid.UUID{}
	_ = pq.StringArray{}
	_ = hstore.Hstore{}
	_ = geom.Point{}
)

func (m *SchemaMigration) GetPrimaryKeyColumn() string {
	return SchemaMigrationTablePrimaryKeyColumn
}

func (m *SchemaMigration) GetPrimaryKeyValue() any {
	return m.Version
}

func (m *SchemaMigration) FromItem(item map[string]any) error {
	if item == nil {
		return fmt.Errorf(
			"item unexpectedly nil during SchemaMigrationFromItem",
		)
	}

	if len(item) == 0 {
		return fmt.Errorf(
			"item unexpectedly empty during SchemaMigrationFromItem",
		)
	}

	wrapError := func(k string, err error) error {
		return fmt.Errorf("%#+v: %v; item: %#+v", k, err, item)
	}

	var err error

	for k, v := range item {
		_, ok := SchemaMigrationTableColumnLookup[k]
		if !ok {
			return fmt.Errorf(
				"item contained unexpected key %#+v during SchemaMigrationFromItem; item: %#+v",
				k, item,
			)
		}

		switch k {
		case "version":
			m.Version, err = types.ParseInt(v)
			if err != nil {
				return wrapError(k, err)
			}
		case "dirty":
			m.Dirty, err = types.ParseBool(v)
			if err != nil {
				return wrapError(k, err)
			}
		}
	}

	return nil
}

func (m *SchemaMigration) Reload(
	ctx context.Context,
	tx *sqlx.Tx,
) error {
	t, err := SelectSchemaMigration(
		ctx,
		tx,
		fmt.Sprintf("%v = $1", m.GetPrimaryKeyColumn()),
		m.GetPrimaryKeyValue(),
	)
	if err != nil {
		return err
	}
	m.Version = t.Version
	m.Dirty = t.Dirty

	return nil
}

func SelectSchemaMigrations(
	ctx context.Context,
	tx *sqlx.Tx,
	where string,
	limit *int,
	offset *int,
	values ...any,
) ([]*SchemaMigration, error) {
	items, err := query.Select(
		ctx,
		tx,
		SchemaMigrationTableColumns,
		SchemaMigrationTable,
		where,
		limit,
		offset,
		values...,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to call SelectSchemaMigrations; err: %v", err)
	}

	objects := make([]*SchemaMigration, 0)

	for _, item := range items {
		object := &SchemaMigration{}

		err = object.FromItem(item)
		if err != nil {
			return nil, fmt.Errorf("failed to call SchemaMigration.FromItem; err: %v", err)
		}

		objects = append(objects, object)
	}

	return objects, nil
}

func SelectSchemaMigration(
	ctx context.Context,
	tx *sqlx.Tx,
	where string,
	values ...any,
) (*SchemaMigration, error) {
	objects, err := SelectSchemaMigrations(
		ctx,
		tx,
		where,
		helpers.Ptr(2),
		helpers.Ptr(0),
		values...,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to call SelectSchemaMigration; err: %v", err)
	}

	if len(objects) > 1 {
		return nil, fmt.Errorf("attempt to call SelectSchemaMigration returned more than 1 row")
	}

	if len(objects) < 1 {
		return nil, fmt.Errorf("attempt to call SelectSchemaMigration returned no rows")
	}

	object := objects[0]

	return object, nil
}
