package model_generated

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/netip"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/cridenour/go-postgis"
	"github.com/go-chi/chi/v5"
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
	RawData                     any                `json:"raw_data"`
	ParentPhysicalThingID       *uuid.UUID         `json:"parent_physical_thing_id"`
	ParentPhysicalThingIDObject *PhysicalThing     `json:"parent_physical_thing_id_object"`
	ParentLogicalThingID        *uuid.UUID         `json:"parent_logical_thing_id"`
	ParentLogicalThingIDObject  *LogicalThing      `json:"parent_logical_thing_id_object"`
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

var (
	LogicalThingTableIDColumnWithTypeCast                    = fmt.Sprintf(`"id" AS id`)
	LogicalThingTableCreatedAtColumnWithTypeCast             = fmt.Sprintf(`"created_at" AS created_at`)
	LogicalThingTableUpdatedAtColumnWithTypeCast             = fmt.Sprintf(`"updated_at" AS updated_at`)
	LogicalThingTableDeletedAtColumnWithTypeCast             = fmt.Sprintf(`"deleted_at" AS deleted_at`)
	LogicalThingTableExternalIDColumnWithTypeCast            = fmt.Sprintf(`"external_id" AS external_id`)
	LogicalThingTableNameColumnWithTypeCast                  = fmt.Sprintf(`"name" AS name`)
	LogicalThingTableTypeColumnWithTypeCast                  = fmt.Sprintf(`"type" AS type`)
	LogicalThingTableTagsColumnWithTypeCast                  = fmt.Sprintf(`"tags" AS tags`)
	LogicalThingTableMetadataColumnWithTypeCast              = fmt.Sprintf(`"metadata" AS metadata`)
	LogicalThingTableRawDataColumnWithTypeCast               = fmt.Sprintf(`"raw_data" AS raw_data`)
	LogicalThingTableParentPhysicalThingIDColumnWithTypeCast = fmt.Sprintf(`"parent_physical_thing_id" AS parent_physical_thing_id`)
	LogicalThingTableParentLogicalThingIDColumnWithTypeCast  = fmt.Sprintf(`"parent_logical_thing_id" AS parent_logical_thing_id`)
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

var LogicalThingTableColumnsWithTypeCasts = []string{
	LogicalThingTableIDColumnWithTypeCast,
	LogicalThingTableCreatedAtColumnWithTypeCast,
	LogicalThingTableUpdatedAtColumnWithTypeCast,
	LogicalThingTableDeletedAtColumnWithTypeCast,
	LogicalThingTableExternalIDColumnWithTypeCast,
	LogicalThingTableNameColumnWithTypeCast,
	LogicalThingTableTypeColumnWithTypeCast,
	LogicalThingTableTagsColumnWithTypeCast,
	LogicalThingTableMetadataColumnWithTypeCast,
	LogicalThingTableRawDataColumnWithTypeCast,
	LogicalThingTableParentPhysicalThingIDColumnWithTypeCast,
	LogicalThingTableParentLogicalThingIDColumnWithTypeCast,
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
	_ = geojson.Point{}
	_ = pgtype.Point{}
	_ = _pgtype.Point{}
	_ = postgis.PointZ{}
	_ = netip.Prefix{}
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

	wrapError := func(k string, v any, err error) error {
		return fmt.Errorf("%v: %#+v; error: %v", k, v, err)
	}

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
			if v == nil {
				continue
			}

			temp1, err := types.ParseUUID(v)
			if err != nil {
				return wrapError(k, v, err)
			}

			temp2, ok := temp1.(uuid.UUID)
			if !ok {
				if temp1 != nil {
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uuid.UUID", temp1))
				}
			}

			m.ID = temp2

		case "created_at":
			if v == nil {
				continue
			}

			temp1, err := types.ParseTime(v)
			if err != nil {
				return wrapError(k, v, err)
			}

			temp2, ok := temp1.(time.Time)
			if !ok {
				if temp1 != nil {
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to time.Time", temp1))
				}
			}

			m.CreatedAt = temp2

		case "updated_at":
			if v == nil {
				continue
			}

			temp1, err := types.ParseTime(v)
			if err != nil {
				return wrapError(k, v, err)
			}

			temp2, ok := temp1.(time.Time)
			if !ok {
				if temp1 != nil {
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to time.Time", temp1))
				}
			}

			m.UpdatedAt = temp2

		case "deleted_at":
			if v == nil {
				continue
			}

			temp1, err := types.ParseTime(v)
			if err != nil {
				return wrapError(k, v, err)
			}

			temp2, ok := temp1.(time.Time)
			if !ok {
				if temp1 != nil {
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to time.Time", temp1))
				}
			}

			m.DeletedAt = &temp2

		case "external_id":
			if v == nil {
				continue
			}

			temp1, err := types.ParseString(v)
			if err != nil {
				return wrapError(k, v, err)
			}

			temp2, ok := temp1.(string)
			if !ok {
				if temp1 != nil {
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to string", temp1))
				}
			}

			m.ExternalID = &temp2

		case "name":
			if v == nil {
				continue
			}

			temp1, err := types.ParseString(v)
			if err != nil {
				return wrapError(k, v, err)
			}

			temp2, ok := temp1.(string)
			if !ok {
				if temp1 != nil {
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to string", temp1))
				}
			}

			m.Name = temp2

		case "type":
			if v == nil {
				continue
			}

			temp1, err := types.ParseString(v)
			if err != nil {
				return wrapError(k, v, err)
			}

			temp2, ok := temp1.(string)
			if !ok {
				if temp1 != nil {
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to string", temp1))
				}
			}

			m.Type = temp2

		case "tags":
			if v == nil {
				continue
			}

			temp1, err := types.ParseStringArray(v)
			if err != nil {
				return wrapError(k, v, err)
			}

			temp2, ok := temp1.([]string)
			if !ok {
				if temp1 != nil {
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to []string", temp1))
				}
			}

			m.Tags = temp2

		case "metadata":
			if v == nil {
				continue
			}

			temp1, err := types.ParseHstore(v)
			if err != nil {
				return wrapError(k, v, err)
			}

			temp2, ok := temp1.(map[string]*string)
			if !ok {
				if temp1 != nil {
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to map[string]*string", temp1))
				}
			}

			m.Metadata = temp2

		case "raw_data":
			if v == nil {
				continue
			}

			temp1, err := types.ParseJSON(v)
			if err != nil {
				return wrapError(k, v, err)
			}

			temp2, ok := temp1.(any)
			if !ok {
				if temp1 != nil {
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to any", temp1))
				}
			}

			m.RawData = &temp2

		case "parent_physical_thing_id":
			if v == nil {
				continue
			}

			temp1, err := types.ParseUUID(v)
			if err != nil {
				return wrapError(k, v, err)
			}

			temp2, ok := temp1.(uuid.UUID)
			if !ok {
				if temp1 != nil {
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uuid.UUID", temp1))
				}
			}

			m.ParentPhysicalThingID = &temp2

		case "parent_logical_thing_id":
			if v == nil {
				continue
			}

			temp1, err := types.ParseUUID(v)
			if err != nil {
				return wrapError(k, v, err)
			}

			temp2, ok := temp1.(uuid.UUID)
			if !ok {
				if temp1 != nil {
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uuid.UUID", temp1))
				}
			}

			m.ParentLogicalThingID = &temp2

		}
	}

	return nil
}

func (m *LogicalThing) Reload(
	ctx context.Context,
	tx *sqlx.Tx,
	includeDeleteds ...bool,
) error {
	extraWhere := ""
	if len(includeDeleteds) > 0 {
		if slices.Contains(LogicalThingTableColumns, "deleted_at") {
			extraWhere = "\n    AND (deleted_at IS null OR deleted_at IS NOT null)"
		}
	}

	t, err := SelectLogicalThing(
		ctx,
		tx,
		fmt.Sprintf("%v = $1%v", m.GetPrimaryKeyColumn(), extraWhere),
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

func (m *LogicalThing) Insert(
	ctx context.Context,
	tx *sqlx.Tx,
	setPrimaryKey bool,
	setZeroValues bool,
) error {
	columns := make([]string, 0)
	values := make([]any, 0)

	if setPrimaryKey && (setZeroValues || !types.IsZeroUUID(m.ID)) {
		columns = append(columns, LogicalThingTableIDColumn)

		v, err := types.FormatUUID(m.ID)
		if err != nil {
			return fmt.Errorf("failed to handle m.ID: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.CreatedAt) {
		columns = append(columns, LogicalThingTableCreatedAtColumn)

		v, err := types.FormatTime(m.CreatedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.CreatedAt: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.UpdatedAt) {
		columns = append(columns, LogicalThingTableUpdatedAtColumn)

		v, err := types.FormatTime(m.UpdatedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.UpdatedAt: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.DeletedAt) {
		columns = append(columns, LogicalThingTableDeletedAtColumn)

		v, err := types.FormatTime(m.DeletedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.DeletedAt: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroString(m.ExternalID) {
		columns = append(columns, LogicalThingTableExternalIDColumn)

		v, err := types.FormatString(m.ExternalID)
		if err != nil {
			return fmt.Errorf("failed to handle m.ExternalID: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroString(m.Name) {
		columns = append(columns, LogicalThingTableNameColumn)

		v, err := types.FormatString(m.Name)
		if err != nil {
			return fmt.Errorf("failed to handle m.Name: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroString(m.Type) {
		columns = append(columns, LogicalThingTableTypeColumn)

		v, err := types.FormatString(m.Type)
		if err != nil {
			return fmt.Errorf("failed to handle m.Type: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroStringArray(m.Tags) {
		columns = append(columns, LogicalThingTableTagsColumn)

		v, err := types.FormatStringArray(m.Tags)
		if err != nil {
			return fmt.Errorf("failed to handle m.Tags: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroHstore(m.Metadata) {
		columns = append(columns, LogicalThingTableMetadataColumn)

		v, err := types.FormatHstore(m.Metadata)
		if err != nil {
			return fmt.Errorf("failed to handle m.Metadata: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroJSON(m.RawData) {
		columns = append(columns, LogicalThingTableRawDataColumn)

		v, err := types.FormatJSON(m.RawData)
		if err != nil {
			return fmt.Errorf("failed to handle m.RawData: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroUUID(m.ParentPhysicalThingID) {
		columns = append(columns, LogicalThingTableParentPhysicalThingIDColumn)

		v, err := types.FormatUUID(m.ParentPhysicalThingID)
		if err != nil {
			return fmt.Errorf("failed to handle m.ParentPhysicalThingID: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroUUID(m.ParentLogicalThingID) {
		columns = append(columns, LogicalThingTableParentLogicalThingIDColumn)

		v, err := types.FormatUUID(m.ParentLogicalThingID)
		if err != nil {
			return fmt.Errorf("failed to handle m.ParentLogicalThingID: %v", err)
		}

		values = append(values, v)
	}

	item, err := query.Insert(
		ctx,
		tx,
		LogicalThingTable,
		columns,
		nil,
		false,
		false,
		LogicalThingTableColumns,
		values...,
	)
	if err != nil {
		return fmt.Errorf("failed to insert %#+v: %v", m, err)
	}
	v := item[LogicalThingTableIDColumn]

	if v == nil {
		return fmt.Errorf("failed to find %v in %#+v", LogicalThingTableIDColumn, item)
	}

	wrapError := func(err error) error {
		return fmt.Errorf(
			"failed to treat %v: %#+v as uuid.UUID: %v",
			LogicalThingTableIDColumn,
			item[LogicalThingTableIDColumn],
			err,
		)
	}

	temp1, err := types.ParseUUID(v)
	if err != nil {
		return wrapError(err)
	}

	temp2, ok := temp1.(uuid.UUID)
	if !ok {
		return wrapError(fmt.Errorf("failed to cast to uuid.UUID"))
	}

	m.ID = temp2

	err = m.Reload(ctx, tx)
	if err != nil {
		return fmt.Errorf("failed to reload after insert")
	}

	return nil
}

func (m *LogicalThing) Update(
	ctx context.Context,
	tx *sqlx.Tx,
	setZeroValues bool,
) error {
	columns := make([]string, 0)
	values := make([]any, 0)

	if setZeroValues || !types.IsZeroTime(m.CreatedAt) {
		columns = append(columns, LogicalThingTableCreatedAtColumn)

		v, err := types.FormatTime(m.CreatedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.CreatedAt: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.UpdatedAt) {
		columns = append(columns, LogicalThingTableUpdatedAtColumn)

		v, err := types.FormatTime(m.UpdatedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.UpdatedAt: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.DeletedAt) {
		columns = append(columns, LogicalThingTableDeletedAtColumn)

		v, err := types.FormatTime(m.DeletedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.DeletedAt: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroString(m.ExternalID) {
		columns = append(columns, LogicalThingTableExternalIDColumn)

		v, err := types.FormatString(m.ExternalID)
		if err != nil {
			return fmt.Errorf("failed to handle m.ExternalID: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroString(m.Name) {
		columns = append(columns, LogicalThingTableNameColumn)

		v, err := types.FormatString(m.Name)
		if err != nil {
			return fmt.Errorf("failed to handle m.Name: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroString(m.Type) {
		columns = append(columns, LogicalThingTableTypeColumn)

		v, err := types.FormatString(m.Type)
		if err != nil {
			return fmt.Errorf("failed to handle m.Type: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroStringArray(m.Tags) {
		columns = append(columns, LogicalThingTableTagsColumn)

		v, err := types.FormatStringArray(m.Tags)
		if err != nil {
			return fmt.Errorf("failed to handle m.Tags: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroHstore(m.Metadata) {
		columns = append(columns, LogicalThingTableMetadataColumn)

		v, err := types.FormatHstore(m.Metadata)
		if err != nil {
			return fmt.Errorf("failed to handle m.Metadata: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroJSON(m.RawData) {
		columns = append(columns, LogicalThingTableRawDataColumn)

		v, err := types.FormatJSON(m.RawData)
		if err != nil {
			return fmt.Errorf("failed to handle m.RawData: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroUUID(m.ParentPhysicalThingID) {
		columns = append(columns, LogicalThingTableParentPhysicalThingIDColumn)

		v, err := types.FormatUUID(m.ParentPhysicalThingID)
		if err != nil {
			return fmt.Errorf("failed to handle m.ParentPhysicalThingID: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroUUID(m.ParentLogicalThingID) {
		columns = append(columns, LogicalThingTableParentLogicalThingIDColumn)

		v, err := types.FormatUUID(m.ParentLogicalThingID)
		if err != nil {
			return fmt.Errorf("failed to handle m.ParentLogicalThingID: %v", err)
		}

		values = append(values, v)
	}

	v, err := types.FormatUUID(m.ID)
	if err != nil {
		return fmt.Errorf("failed to handle m.ID: %v", err)
	}

	values = append(values, v)

	_, err = query.Update(
		ctx,
		tx,
		LogicalThingTable,
		columns,
		fmt.Sprintf("%v = $$??", LogicalThingTableIDColumn),
		LogicalThingTableColumns,
		values...,
	)
	if err != nil {
		return fmt.Errorf("failed to update %#+v: %v", m, err)
	}

	err = m.Reload(ctx, tx)
	if err != nil {
		return fmt.Errorf("failed to reload after update")
	}

	return nil
}

func (m *LogicalThing) Delete(
	ctx context.Context,
	tx *sqlx.Tx,
) error {
	values := make([]any, 0)
	v, err := types.FormatUUID(m.ID)
	if err != nil {
		return fmt.Errorf("failed to handle m.ID: %v", err)
	}

	values = append(values, v)

	err = query.Delete(
		ctx,
		tx,
		LogicalThingTable,
		fmt.Sprintf("%v = $$??", LogicalThingTableIDColumn),
		values...,
	)
	if err != nil {
		return fmt.Errorf("failed to delete %#+v: %v", m, err)
	}

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
	if slices.Contains(LogicalThingTableColumns, "deleted_at") {
		if !strings.Contains(where, "deleted_at") {
			if where != "" {
				where += "\n    AND "
			}

			where += "deleted_at IS null"
		}
	}

	items, err := query.Select(
		ctx,
		tx,
		LogicalThingTableColumnsWithTypeCasts,
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

		if !types.IsZeroUUID(object.ParentPhysicalThingID) {
			object.ParentPhysicalThingIDObject, _ = SelectPhysicalThing(
				ctx,
				tx,
				fmt.Sprintf("%v = $1", PhysicalThingTablePrimaryKeyColumn),
				object.ParentPhysicalThingID,
			)
		}

		if !types.IsZeroUUID(object.ParentLogicalThingID) {
			object.ParentLogicalThingIDObject, _ = SelectLogicalThing(
				ctx,
				tx,
				fmt.Sprintf("%v = $1", LogicalThingTablePrimaryKeyColumn),
				object.ParentLogicalThingID,
			)
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

func handleGetLogicalThings(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
	ctx := r.Context()

	unrecognizedParams := make([]string, 0)
	hadUnrecognizedParams := false

	unparseableParams := make([]string, 0)
	hadUnparseableParams := false

	values := make([]any, 0)

	wheres := make([]string, 0)
	for rawKey, rawValues := range r.URL.Query() {
		if rawKey == "limit" || rawKey == "offset" {
			continue
		}

		parts := strings.Split(rawKey, "__")
		isUnrecognized := len(parts) != 2

		comparison := ""
		isSliceComparison := false
		isNullComparison := false
		IsLikeComparison := false

		if !isUnrecognized {
			column := LogicalThingTableColumnLookup[parts[0]]
			if column == nil {
				isUnrecognized = true
			} else {
				switch parts[1] {
				case "eq":
					comparison = "="
				case "ne":
					comparison = "!="
				case "gt":
					comparison = ">"
				case "gte":
					comparison = ">="
				case "lt":
					comparison = "<"
				case "lte":
					comparison = "<="
				case "in":
					comparison = "IN"
					isSliceComparison = true
				case "nin", "notin":
					comparison = "NOT IN"
					isSliceComparison = true
				case "isnull":
					comparison = "IS NULL"
					isNullComparison = true
				case "nisnull", "isnotnull":
					comparison = "IS NOT NULL"
					isNullComparison = true
				case "l", "like":
					comparison = "LIKE"
					IsLikeComparison = true
				case "nl", "nlike", "notlike":
					comparison = "NOT LIKE"
					IsLikeComparison = true
				case "il", "ilike":
					comparison = "ILIKE"
					IsLikeComparison = true
				case "nil", "nilike", "notilike":
					comparison = "NOT ILIKE"
					IsLikeComparison = true
				default:
					isUnrecognized = true
				}
			}
		}

		if isNullComparison {
			wheres = append(wheres, fmt.Sprintf("%s %s", parts[0], comparison))
			continue
		}

		for _, rawValue := range rawValues {
			if isUnrecognized {
				unrecognizedParams = append(unrecognizedParams, fmt.Sprintf("%s=%s", rawKey, rawValue))
				hadUnrecognizedParams = true
				continue
			}

			if hadUnrecognizedParams {
				continue
			}

			attempts := make([]string, 0)

			if !IsLikeComparison {
				attempts = append(attempts, rawValue)
			}

			if isSliceComparison {
				attempts = append(attempts, fmt.Sprintf("[%s]", rawValue))

				vs := make([]string, 0)
				for _, v := range strings.Split(rawValue, ",") {
					vs = append(vs, fmt.Sprintf("\"%s\"", v))
				}

				attempts = append(attempts, fmt.Sprintf("[%s]", strings.Join(vs, ",")))
			}

			if IsLikeComparison {
				attempts = append(attempts, fmt.Sprintf("\"%%%s%%\"", rawValue))
			} else {
				attempts = append(attempts, fmt.Sprintf("\"%s\"", rawValue))
			}

			var err error

			for _, attempt := range attempts {
				var value any
				err = json.Unmarshal([]byte(attempt), &value)
				if err == nil {
					if isSliceComparison {
						sliceValues, ok := value.([]any)
						if !ok {
							err = fmt.Errorf("failed to cast %#+v to []string", value)
							break
						}

						values = append(values, sliceValues...)

						sliceWheres := make([]string, 0)
						for range values {
							sliceWheres = append(sliceWheres, "$$??")
						}

						wheres = append(wheres, fmt.Sprintf("%s %s (%s)", parts[0], comparison, strings.Join(sliceWheres, ", ")))
					} else {
						values = append(values, value)
						wheres = append(wheres, fmt.Sprintf("%s %s $$??", parts[0], comparison))
					}

					break
				}
			}

			if err != nil {
				unparseableParams = append(unparseableParams, fmt.Sprintf("%s=%s", rawKey, rawValue))
				hadUnparseableParams = true
				continue
			}
		}
	}

	where := strings.Join(wheres, "\n    AND ")

	if hadUnrecognizedParams {
		helpers.HandleErrorResponse(
			w,
			http.StatusInternalServerError,
			fmt.Errorf("unrecognized params %s", strings.Join(unrecognizedParams, ", ")),
		)
		return
	}

	if hadUnparseableParams {
		helpers.HandleErrorResponse(
			w,
			http.StatusInternalServerError,
			fmt.Errorf("unparseable params %s", strings.Join(unparseableParams, ", ")),
		)
		return
	}

	limit := 2000
	rawLimit := r.URL.Query().Get("limit")
	if rawLimit != "" {
		possibleLimit, err := strconv.ParseInt(rawLimit, 10, 64)
		if err == nil {
			helpers.HandleErrorResponse(
				w,
				http.StatusInternalServerError,
				fmt.Errorf("failed to parse param limit=%s as int: %v", rawLimit, err),
			)
			return
		}

		limit = int(possibleLimit)
	}

	offset := 0
	rawOffset := r.URL.Query().Get("offset")
	if rawOffset != "" {
		possibleOffset, err := strconv.ParseInt(rawOffset, 10, 64)
		if err == nil {
			helpers.HandleErrorResponse(
				w,
				http.StatusInternalServerError,
				fmt.Errorf("failed to parse param offset=%s as int: %v", rawOffset, err),
			)
			return
		}

		offset = int(possibleOffset)
	}

	tx, err := db.BeginTxx(r.Context(), nil)
	if err != nil {
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	objects, err := SelectLogicalThings(ctx, tx, where, &limit, &offset, values...)
	if err != nil {
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	err = tx.Commit()
	if err != nil {
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	helpers.HandleObjectsResponse(w, http.StatusOK, objects)
}

func handleGetLogicalThing(w http.ResponseWriter, r *http.Request, db *sqlx.DB, primaryKey string) {
	ctx := r.Context()

	where := fmt.Sprintf("%s = $$??", LogicalThingTablePrimaryKeyColumn)

	values := []any{primaryKey}

	tx, err := db.BeginTxx(r.Context(), nil)
	if err != nil {
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	object, err := SelectLogicalThing(ctx, tx, where, values...)
	if err != nil {
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	err = tx.Commit()
	if err != nil {
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	helpers.HandleObjectsResponse(w, http.StatusOK, []*LogicalThing{object})
}

func handlePostLogicalThings(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		err = fmt.Errorf("failed to read body of HTTP request: %v", err)
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	var allItems []map[string]any
	err = json.Unmarshal(b, &allItems)
	if err != nil {
		err = fmt.Errorf("failed to unmarshal %#+v as JSON list of objects: %v", string(b), err)
		helpers.HandleErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	objects := make([]*LogicalThing, 0)
	for _, item := range allItems {
		object := &LogicalThing{}
		err = object.FromItem(item)
		if err != nil {
			err = fmt.Errorf("failed to interpret %#+v as LogicalThing in item form: %v", item, err)
			helpers.HandleErrorResponse(w, http.StatusBadRequest, err)
			return
		}

		objects = append(objects, object)
	}

	tx, err := db.BeginTxx(r.Context(), nil)
	if err != nil {
		err = fmt.Errorf("failed to begin DB transaction: %v", err)
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	defer func() {
		_ = tx.Rollback()
	}()

	for i, object := range objects {
		err = object.Insert(r.Context(), tx, false, false)
		if err != nil {
			err = fmt.Errorf("failed to insert %#+v: %v", object, err)
			helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
			return
		}

		objects[i] = object
	}

	err = tx.Commit()
	if err != nil {
		err = fmt.Errorf("failed to commit DB transaction: %v", err)
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	helpers.HandleObjectsResponse(w, http.StatusCreated, objects)
}

func handlePutLogicalThing(w http.ResponseWriter, r *http.Request, db *sqlx.DB, primaryKey string) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		err = fmt.Errorf("failed to read body of HTTP request: %v", err)
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	var item map[string]any
	err = json.Unmarshal(b, &item)
	if err != nil {
		err = fmt.Errorf("failed to unmarshal %#+v as JSON object: %v", string(b), err)
		helpers.HandleErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	item[LogicalThingTablePrimaryKeyColumn] = primaryKey

	object := &LogicalThing{}
	err = object.FromItem(item)
	if err != nil {
		err = fmt.Errorf("failed to interpret %#+v as LogicalThing in item form: %v", item, err)
		helpers.HandleErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	tx, err := db.BeginTxx(r.Context(), nil)
	if err != nil {
		err = fmt.Errorf("failed to begin DB transaction: %v", err)
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	defer func() {
		_ = tx.Rollback()
	}()

	err = object.Update(r.Context(), tx, true)
	if err != nil {
		err = fmt.Errorf("failed to update %#+v: %v", object, err)
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	err = tx.Commit()
	if err != nil {
		err = fmt.Errorf("failed to commit DB transaction: %v", err)
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	helpers.HandleObjectsResponse(w, http.StatusOK, []*LogicalThing{object})
}

func handlePatchLogicalThing(w http.ResponseWriter, r *http.Request, db *sqlx.DB, primaryKey string) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		err = fmt.Errorf("failed to read body of HTTP request: %v", err)
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	var item map[string]any
	err = json.Unmarshal(b, &item)
	if err != nil {
		err = fmt.Errorf("failed to unmarshal %#+v as JSON object: %v", string(b), err)
		helpers.HandleErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	item[LogicalThingTablePrimaryKeyColumn] = primaryKey

	object := &LogicalThing{}
	err = object.FromItem(item)
	if err != nil {
		err = fmt.Errorf("failed to interpret %#+v as LogicalThing in item form: %v", item, err)
		helpers.HandleErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	tx, err := db.BeginTxx(r.Context(), nil)
	if err != nil {
		err = fmt.Errorf("failed to begin DB transaction: %v", err)
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	defer func() {
		_ = tx.Rollback()
	}()

	err = object.Update(r.Context(), tx, false)
	if err != nil {
		err = fmt.Errorf("failed to update %#+v: %v", object, err)
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	err = tx.Commit()
	if err != nil {
		err = fmt.Errorf("failed to commit DB transaction: %v", err)
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	helpers.HandleObjectsResponse(w, http.StatusOK, []*LogicalThing{object})
}

func GetLogicalThingRouter(db *sqlx.DB, middlewares ...func(http.Handler) http.Handler) chi.Router {
	r := chi.NewRouter()

	for _, m := range middlewares {
		r.Use(m)
	}

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		handleGetLogicalThings(w, r, db)
	})

	r.Get("/{primaryKey}", func(w http.ResponseWriter, r *http.Request) {
		handleGetLogicalThing(w, r, db, chi.URLParam(r, "primaryKey"))
	})

	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		handlePostLogicalThings(w, r, db)
	})

	r.Put("/{primaryKey}", func(w http.ResponseWriter, r *http.Request) {
		handlePutLogicalThing(w, r, db, chi.URLParam(r, "primaryKey"))
	})

	r.Patch("/{primaryKey}", func(w http.ResponseWriter, r *http.Request) {
		handlePatchLogicalThing(w, r, db, chi.URLParam(r, "primaryKey"))
	})

	return r
}

func GetLogicalThingHandlerFunc(db *sqlx.DB, middlewares ...func(http.Handler) http.Handler) http.HandlerFunc {
	r := chi.NewRouter()

	r.Mount("/logical-things", GetLogicalThingRouter(db, middlewares...))

	return r.ServeHTTP
}

func NewLogicalThingFromItem(item map[string]any) (any, error) {
	object := &LogicalThing{}

	err := object.FromItem(item)
	if err != nil {
		return nil, err
	}

	return object, nil
}

func init() {
	register(
		LogicalThingTable,
		NewLogicalThingFromItem,
		"/logical-things",
		GetLogicalThingRouter,
	)
}
