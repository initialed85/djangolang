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

type LogicalThing struct {
	ID                          uuid.UUID          `json:"id"`
	CreatedAt                   time.Time          `json:"created_at"`
	UpdatedAt                   time.Time          `json:"updated_at"`
	DeletedAt                   *time.Time         `json:"deleted_at"`
	ExternalID                  *string            `json:"external_id"`
	Name                        string             `json:"name"`
	Type                        string             `json:"type"`
	Tags                        []string           `json:"tags"`
	Metadata                    map[string]*string `json:"metadata"`
	RawData                     *any               `json:"raw_data"`
	ParentPhysicalThingID       *uuid.UUID         `json:"parent_physical_thing_id"`
	ParentPhysicalThingIDObject *PhysicalThing     `json:"-"`
	ParentLogicalThingID        *uuid.UUID         `json:"parent_logical_thing_id"`
	ParentLogicalThingIDObject  *LogicalThing      `json:"-"`
}

var LogicalThingTable = "logical_things"

var (
	LogicalThingTableIDColumn                    = "id"
	LogicalThingTableCreatedAtColumn             = "created_at"
	LogicalThingTableUpdatedAtColumn             = "updated_at"
	LogicalThingTableDeletedAtColumn             = "deleted_at"
	LogicalThingTableExternalIDColumn            = "external_id"
	LogicalThingTableNameColumn                  = "name"
	LogicalThingTableTypeColumn                  = "type"
	LogicalThingTableTagsColumn                  = "tags"
	LogicalThingTableMetadataColumn              = "metadata"
	LogicalThingTableRawDataColumn               = "raw_data"
	LogicalThingTableParentPhysicalThingIDColumn = "parent_physical_thing_id"
	LogicalThingTableParentLogicalThingIDColumn  = "parent_logical_thing_id"
)

var LogicalThingTableColumns = []string{
	LogicalThingTableIDColumn,
	LogicalThingTableCreatedAtColumn,
	LogicalThingTableUpdatedAtColumn,
	LogicalThingTableDeletedAtColumn,
	LogicalThingTableExternalIDColumn,
	LogicalThingTableNameColumn,
	LogicalThingTableTypeColumn,
	LogicalThingTableTagsColumn,
	LogicalThingTableMetadataColumn,
	LogicalThingTableRawDataColumn,
	LogicalThingTableParentPhysicalThingIDColumn,
	LogicalThingTableParentLogicalThingIDColumn,
}

var LogicalThingTableColumnLookup = map[string]*introspect.Column{
	LogicalThingTableIDColumn:                    new(introspect.Column),
	LogicalThingTableCreatedAtColumn:             new(introspect.Column),
	LogicalThingTableUpdatedAtColumn:             new(introspect.Column),
	LogicalThingTableDeletedAtColumn:             new(introspect.Column),
	LogicalThingTableExternalIDColumn:            new(introspect.Column),
	LogicalThingTableNameColumn:                  new(introspect.Column),
	LogicalThingTableTypeColumn:                  new(introspect.Column),
	LogicalThingTableTagsColumn:                  new(introspect.Column),
	LogicalThingTableMetadataColumn:              new(introspect.Column),
	LogicalThingTableRawDataColumn:               new(introspect.Column),
	LogicalThingTableParentPhysicalThingIDColumn: new(introspect.Column),
	LogicalThingTableParentLogicalThingIDColumn:  new(introspect.Column),
}

var (
	LogicalThingTablePrimaryKeyColumn = LogicalThingTableIDColumn
)

var (
	_ = time.Time{}
	_ = uuid.UUID{}
	_ = pq.StringArray{}
	_ = hstore.Hstore{}
	_ = geom.Point{}
)

func (m *LogicalThing) GetPrimaryKeyColumn() string {
	return LogicalThingTablePrimaryKeyColumn
}

func (m *LogicalThing) GetPrimaryKeyValue() any {
	return m.ID
}

func (m *LogicalThing) FromItem(item map[string]any) error {
	if item == nil {
		return fmt.Errorf(
			"item unexpectedly nil during LogicalThingFromItem",
		)
	}

	if len(item) == 0 {
		return fmt.Errorf(
			"item unexpectedly empty during LogicalThingFromItem",
		)
	}

	wrapError := func(k string, err error) error {
		return fmt.Errorf("%#+v: %v; item: %#+v", k, err, item)
	}

	var err error

	for k, v := range item {
		_, ok := LogicalThingTableColumnLookup[k]
		if !ok {
			return fmt.Errorf(
				"item contained unexpected key %#+v during LogicalThingFromItem; item: %#+v",
				k, item,
			)
		}

		switch k {
		case "id":
			m.ID, err = types.ParseUUID(v)
			if err != nil {
				return wrapError(k, err)
			}
		case "created_at":
			m.CreatedAt, err = types.ParseTime(v)
			if err != nil {
				return wrapError(k, err)
			}
		case "updated_at":
			m.UpdatedAt, err = types.ParseTime(v)
			if err != nil {
				return wrapError(k, err)
			}
		case "deleted_at":
			m.DeletedAt, err = types.ParsePtr(types.ParseTime, v)
			if err != nil {
				return wrapError(k, err)
			}
		case "external_id":
			m.ExternalID, err = types.ParsePtr(types.ParseString, v)
			if err != nil {
				return wrapError(k, err)
			}
		case "name":
			m.Name, err = types.ParseString(v)
			if err != nil {
				return wrapError(k, err)
			}
		case "type":
			m.Type, err = types.ParseString(v)
			if err != nil {
				return wrapError(k, err)
			}
		case "tags":
			m.Tags, err = types.ParseStringArray(v)
			if err != nil {
				return wrapError(k, err)
			}
		case "metadata":
			m.Metadata, err = types.ParseHstore(v)
			if err != nil {
				return wrapError(k, err)
			}
		case "raw_data":
			m.RawData, err = types.ParsePtr(types.ParseJSON, v)
			if err != nil {
				return wrapError(k, err)
			}
		case "parent_physical_thing_id":
			m.ParentPhysicalThingID, err = types.ParsePtr(types.ParseUUID, v)
			if err != nil {
				return wrapError(k, err)
			}
		case "parent_logical_thing_id":
			m.ParentLogicalThingID, err = types.ParsePtr(types.ParseUUID, v)
			if err != nil {
				return wrapError(k, err)
			}
		}
	}

	return nil
}

func (m *LogicalThing) Reload(
	ctx context.Context,
	tx *sqlx.Tx,
) error {
	t, err := SelectLogicalThing(
		ctx,
		tx,
		fmt.Sprintf("%v = $1", m.GetPrimaryKeyColumn()),
		m.GetPrimaryKeyValue(),
	)
	if err != nil {
		return err
	}
	m.ID = t.ID
	m.CreatedAt = t.CreatedAt
	m.UpdatedAt = t.UpdatedAt
	m.DeletedAt = t.DeletedAt
	m.ExternalID = t.ExternalID
	m.Name = t.Name
	m.Type = t.Type
	m.Tags = t.Tags
	m.Metadata = t.Metadata
	m.RawData = t.RawData
	m.ParentPhysicalThingID = t.ParentPhysicalThingID
	m.ParentPhysicalThingIDObject = t.ParentPhysicalThingIDObject
	m.ParentLogicalThingID = t.ParentLogicalThingID
	m.ParentLogicalThingIDObject = t.ParentLogicalThingIDObject

	return nil
}

func SelectLogicalThings(
	ctx context.Context,
	tx *sqlx.Tx,
	where string,
	limit *int,
	offset *int,
	values ...any,
) ([]*LogicalThing, error) {
	items, err := query.Select(
		ctx,
		tx,
		LogicalThingTableColumns,
		LogicalThingTable,
		where,
		limit,
		offset,
		values...,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to call SelectLogicalThings; err: %v", err)
	}

	objects := make([]*LogicalThing, 0)

	for _, item := range items {
		object := &LogicalThing{}

		err = object.FromItem(item)
		if err != nil {
			return nil, fmt.Errorf("failed to call LogicalThing.FromItem; err: %v", err)
		}
		if object.ParentPhysicalThingID != nil {
			object.ParentPhysicalThingIDObject, err = SelectPhysicalThing(
				ctx,
				tx,
				fmt.Sprintf("%v = $1", PhysicalThingTablePrimaryKeyColumn),
				object.ParentPhysicalThingID,
			)
			if err != nil {
				return nil, fmt.Errorf("failed to load <no value>.ParentPhysicalThingIDObject; err: %v", err)
			}
		}
		if object.ParentLogicalThingID != nil {
			object.ParentLogicalThingIDObject, err = SelectLogicalThing(
				ctx,
				tx,
				fmt.Sprintf("%v = $1", PhysicalThingTablePrimaryKeyColumn),
				object.ParentLogicalThingID,
			)
			if err != nil {
				return nil, fmt.Errorf("failed to load <no value>.ParentLogicalThingIDObject; err: %v", err)
			}
		}

		objects = append(objects, object)
	}

	return objects, nil
}

func SelectLogicalThing(
	ctx context.Context,
	tx *sqlx.Tx,
	where string,
	values ...any,
) (*LogicalThing, error) {
	objects, err := SelectLogicalThings(
		ctx,
		tx,
		where,
		helpers.Ptr(2),
		helpers.Ptr(0),
		values...,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to call SelectLogicalThing; err: %v", err)
	}

	if len(objects) > 1 {
		return nil, fmt.Errorf("attempt to call SelectLogicalThing returned more than 1 row")
	}

	if len(objects) < 1 {
		return nil, fmt.Errorf("attempt to call SelectLogicalThing returned no rows")
	}

	object := objects[0]

	return object, nil
}
