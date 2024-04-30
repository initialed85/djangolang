package model_reference

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
)

type PhysicalThing struct {
	ID         uuid.UUID          `json:"id"`
	CreatedAt  time.Time          `json:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at"`
	DeletedAt  *time.Time         `json:"deleted_at"`
	ExternalID *string            `json:"external_id"`
	Name       string             `json:"name"`
	Type       string             `json:"type"`
	Tags       []string           `json:"tags"`
	Metadata   map[string]*string `json:"metadata"`
	RawData    any                `json:"raw_data"`
}

var (
	PhysicalThingTable = "physical_things"

	PhysicalThingTableIDColumn         = "id"
	PhysicalThingTableCreatedAtColumn  = "created_at"
	PhysicalThingTableUpdatedAtColumn  = "updated_at"
	PhysicalThingTableDeletedAtColumn  = "deleted_at"
	PhysicalThingTableExternalIDColumn = "external_id"
	PhysicalThingTableNameColumn       = "name"
	PhysicalThingTableTypeColumn       = "type"
	PhysicalThingTableTagsColumn       = "tags"
	PhysicalThingTableMetadataColumn   = "metadata"
	PhysicalThingTableRawDataColumn    = "raw_data"

	PhysicalThingTablePrimaryKeyColumn = PhysicalThingTableIDColumn
)

var PhysicalThingTableColumns = []string{
	PhysicalThingTableIDColumn,
	PhysicalThingTableCreatedAtColumn,
	PhysicalThingTableUpdatedAtColumn,
	PhysicalThingTableDeletedAtColumn,
	PhysicalThingTableExternalIDColumn,
	PhysicalThingTableNameColumn,
	PhysicalThingTableTypeColumn,
	PhysicalThingTableTagsColumn,
	PhysicalThingTableMetadataColumn,
	PhysicalThingTableRawDataColumn,
}

var PhysicalThingTableColumnLookup = map[string]*introspect.Column{
	PhysicalThingTableIDColumn:         new(introspect.Column),
	PhysicalThingTableCreatedAtColumn:  new(introspect.Column),
	PhysicalThingTableUpdatedAtColumn:  new(introspect.Column),
	PhysicalThingTableDeletedAtColumn:  new(introspect.Column),
	PhysicalThingTableExternalIDColumn: new(introspect.Column),
	PhysicalThingTableNameColumn:       new(introspect.Column),
	PhysicalThingTableTypeColumn:       new(introspect.Column),
	PhysicalThingTableTagsColumn:       new(introspect.Column),
	PhysicalThingTableMetadataColumn:   new(introspect.Column),
	PhysicalThingTableRawDataColumn:    new(introspect.Column),
}

func (m *PhysicalThing) GetPrimaryKeyColumn() string {
	return PhysicalThingTablePrimaryKeyColumn
}

func (m *PhysicalThing) GetPrimaryKeyValue() any {
	return m.ID
}

func (m *PhysicalThing) FromItem(item map[string]any) error {
	if item == nil {
		return fmt.Errorf(
			"item unexpectedly nil during PhysicalThingFromItem",
		)
	}

	if len(item) == 0 {
		return fmt.Errorf(
			"item unexpectedly empty during PhysicalThingFromItem",
		)
	}

	wrapError := func(k string, err error) error {
		return fmt.Errorf("%#+v: %v; item: %#+v", k, err, item)
	}

	var err error

	for k, v := range item {
		_, ok := PhysicalThingTableColumnLookup[k]
		if !ok {
			return fmt.Errorf(
				"item contained unexpected key %#+v during PhysicalThingFromItem; item: %#+v",
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
		}
	}

	return nil
}

func (m *PhysicalThing) Reload(
	ctx context.Context,
	tx *sqlx.Tx,
) error {
	t, err := SelectPhysicalThing(
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

	return nil
}

func SelectPhysicalThings(
	ctx context.Context,
	tx *sqlx.Tx,
	where string,
	limit *int,
	offset *int,
	values ...any,
) ([]*PhysicalThing, error) {
	items, err := query.Select(
		ctx,
		tx,
		PhysicalThingTableColumns,
		PhysicalThingTable,
		where,
		limit,
		offset,
		values...,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to call SelectPhysicalThings; err: %v", err)
	}

	objects := make([]*PhysicalThing, 0)

	for _, item := range items {
		object := &PhysicalThing{}

		err = object.FromItem(item)
		if err != nil {
			return nil, fmt.Errorf("failed to call PhysicalThing.FromItem; err: %v", err)
		}

		objects = append(objects, object)
	}

	return objects, nil
}

func SelectPhysicalThing(
	ctx context.Context,
	tx *sqlx.Tx,
	where string,
	values ...any,
) (*PhysicalThing, error) {
	objects, err := SelectPhysicalThings(
		ctx,
		tx,
		where,
		helpers.Ptr(2),
		helpers.Ptr(0),
		values...,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to call SelectPhysicalThing; err: %v", err)
	}

	if len(objects) > 1 {
		return nil, fmt.Errorf("attempt to call SelectPhysicalThing returned more than 1 row")
	}

	if len(objects) < 1 {
		return nil, fmt.Errorf("attempt to call SelectPhysicalThing returned no rows")
	}

	object := objects[0]

	return object, nil
}
