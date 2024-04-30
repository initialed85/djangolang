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
	_pgtype "github.com/jackc/pgtype"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/lib/pq/hstore"
	"github.com/paulmach/orb/geojson"
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
	RawData    *any               `json:"raw_data"`
}

var PhysicalThingTable = "physical_things"

var (
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
)

var (
	PhysicalThingTableIDColumnWithTypeCast         = fmt.Sprintf(`"id" AS id`)
	PhysicalThingTableCreatedAtColumnWithTypeCast  = fmt.Sprintf(`"created_at" AS created_at`)
	PhysicalThingTableUpdatedAtColumnWithTypeCast  = fmt.Sprintf(`"updated_at" AS updated_at`)
	PhysicalThingTableDeletedAtColumnWithTypeCast  = fmt.Sprintf(`"deleted_at" AS deleted_at`)
	PhysicalThingTableExternalIDColumnWithTypeCast = fmt.Sprintf(`"external_id" AS external_id`)
	PhysicalThingTableNameColumnWithTypeCast       = fmt.Sprintf(`"name" AS name`)
	PhysicalThingTableTypeColumnWithTypeCast       = fmt.Sprintf(`"type" AS type`)
	PhysicalThingTableTagsColumnWithTypeCast       = fmt.Sprintf(`"tags" AS tags`)
	PhysicalThingTableMetadataColumnWithTypeCast   = fmt.Sprintf(`"metadata" AS metadata`)
	PhysicalThingTableRawDataColumnWithTypeCast    = fmt.Sprintf(`"raw_data" AS raw_data`)
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

var PhysicalThingTableColumnsWithTypeCasts = []string{
	PhysicalThingTableIDColumnWithTypeCast,
	PhysicalThingTableCreatedAtColumnWithTypeCast,
	PhysicalThingTableUpdatedAtColumnWithTypeCast,
	PhysicalThingTableDeletedAtColumnWithTypeCast,
	PhysicalThingTableExternalIDColumnWithTypeCast,
	PhysicalThingTableNameColumnWithTypeCast,
	PhysicalThingTableTypeColumnWithTypeCast,
	PhysicalThingTableTagsColumnWithTypeCast,
	PhysicalThingTableMetadataColumnWithTypeCast,
	PhysicalThingTableRawDataColumnWithTypeCast,
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

var (
	PhysicalThingTablePrimaryKeyColumn = PhysicalThingTableIDColumn
)

var (
	_ = time.Time{}
	_ = uuid.UUID{}
	_ = pq.StringArray{}
	_ = hstore.Hstore{}
	_ = geojson.Point{}
	_ = pgtype.Point{}
	_ = _pgtype.Point{}
)

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
			if v == nil {
				continue
			}

			temp1, err := types.ParseUUID(v)
			if err != nil {
				return wrapError(k, err)
			}

			temp2, ok := temp1.(uuid.UUID)
			if !ok {
				return wrapError(k, fmt.Errorf("failed to cast to uuid.UUID"))
			}

			m.ID = temp2
		case "created_at":
			if v == nil {
				continue
			}

			temp1, err := types.ParseTime(v)
			if err != nil {
				return wrapError(k, err)
			}

			temp2, ok := temp1.(time.Time)
			if !ok {
				return wrapError(k, fmt.Errorf("failed to cast to time.Time"))
			}

			m.CreatedAt = temp2
		case "updated_at":
			if v == nil {
				continue
			}

			temp1, err := types.ParseTime(v)
			if err != nil {
				return wrapError(k, err)
			}

			temp2, ok := temp1.(time.Time)
			if !ok {
				return wrapError(k, fmt.Errorf("failed to cast to time.Time"))
			}

			m.UpdatedAt = temp2
		case "deleted_at":
			if v == nil {
				continue
			}

			temp1, err := types.ParseTime(v)
			if err != nil {
				return wrapError(k, err)
			}

			temp2, ok := temp1.(time.Time)
			if !ok {
				return wrapError(k, fmt.Errorf("failed to cast to time.Time"))
			}

			m.DeletedAt = &temp2
		case "external_id":
			if v == nil {
				continue
			}

			temp1, err := types.ParseString(v)
			if err != nil {
				return wrapError(k, err)
			}

			temp2, ok := temp1.(string)
			if !ok {
				return wrapError(k, fmt.Errorf("failed to cast to string"))
			}

			m.ExternalID = &temp2
		case "name":
			if v == nil {
				continue
			}

			temp1, err := types.ParseString(v)
			if err != nil {
				return wrapError(k, err)
			}

			temp2, ok := temp1.(string)
			if !ok {
				return wrapError(k, fmt.Errorf("failed to cast to string"))
			}

			m.Name = temp2
		case "type":
			if v == nil {
				continue
			}

			temp1, err := types.ParseString(v)
			if err != nil {
				return wrapError(k, err)
			}

			temp2, ok := temp1.(string)
			if !ok {
				return wrapError(k, fmt.Errorf("failed to cast to string"))
			}

			m.Type = temp2
		case "tags":
			if v == nil {
				continue
			}

			temp1, err := types.ParseStringArray(v)
			if err != nil {
				return wrapError(k, err)
			}

			temp2, ok := temp1.([]string)
			if !ok {
				return wrapError(k, fmt.Errorf("failed to cast to []string"))
			}

			m.Tags = temp2
		case "metadata":
			if v == nil {
				continue
			}

			temp1, err := types.ParseHstore(v)
			if err != nil {
				return wrapError(k, err)
			}

			temp2, ok := temp1.(map[string]*string)
			if !ok {
				return wrapError(k, fmt.Errorf("failed to cast to map[string]*string"))
			}

			m.Metadata = temp2
		case "raw_data":
			if v == nil {
				continue
			}

			temp1, err := types.ParseJSON(v)
			if err != nil {
				return wrapError(k, err)
			}

			temp2, ok := temp1.(any)
			if !ok {
				return wrapError(k, fmt.Errorf("failed to cast to any"))
			}

			m.RawData = &temp2
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
		PhysicalThingTableColumnsWithTypeCasts,
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
