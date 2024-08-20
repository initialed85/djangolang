package model_generated

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/netip"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/cridenour/go-postgis"
	"github.com/go-chi/chi/v5"
	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
	"github.com/initialed85/djangolang/pkg/helpers"
	"github.com/initialed85/djangolang/pkg/introspect"
	"github.com/initialed85/djangolang/pkg/query"
	"github.com/initialed85/djangolang/pkg/server"
	"github.com/initialed85/djangolang/pkg/stream"
	"github.com/initialed85/djangolang/pkg/types"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/lib/pq/hstore"
	"golang.org/x/exp/maps"
)

type PhysicalThing struct {
	ID                                                      uuid.UUID          `json:"id"`
	CreatedAt                                               time.Time          `json:"created_at"`
	UpdatedAt                                               time.Time          `json:"updated_at"`
	DeletedAt                                               *time.Time         `json:"deleted_at"`
	ExternalID                                              *string            `json:"external_id"`
	Name                                                    string             `json:"name"`
	Type                                                    string             `json:"type"`
	Tags                                                    []string           `json:"tags"`
	Metadata                                                map[string]*string `json:"metadata"`
	RawData                                                 any                `json:"raw_data"`
	ReferencedByLocationHistoryParentPhysicalThingIDObjects []*LocationHistory `json:"referenced_by_location_history_parent_physical_thing_id_objects"`
	ReferencedByLogicalThingParentPhysicalThingIDObjects    []*LogicalThing    `json:"referenced_by_logical_thing_parent_physical_thing_id_objects"`
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
	PhysicalThingTableIDColumn:         {Name: PhysicalThingTableIDColumn, NotNull: true, HasDefault: true},
	PhysicalThingTableCreatedAtColumn:  {Name: PhysicalThingTableCreatedAtColumn, NotNull: true, HasDefault: true},
	PhysicalThingTableUpdatedAtColumn:  {Name: PhysicalThingTableUpdatedAtColumn, NotNull: true, HasDefault: true},
	PhysicalThingTableDeletedAtColumn:  {Name: PhysicalThingTableDeletedAtColumn, NotNull: false, HasDefault: false},
	PhysicalThingTableExternalIDColumn: {Name: PhysicalThingTableExternalIDColumn, NotNull: false, HasDefault: false},
	PhysicalThingTableNameColumn:       {Name: PhysicalThingTableNameColumn, NotNull: true, HasDefault: false},
	PhysicalThingTableTypeColumn:       {Name: PhysicalThingTableTypeColumn, NotNull: true, HasDefault: false},
	PhysicalThingTableTagsColumn:       {Name: PhysicalThingTableTagsColumn, NotNull: true, HasDefault: true},
	PhysicalThingTableMetadataColumn:   {Name: PhysicalThingTableMetadataColumn, NotNull: true, HasDefault: true},
	PhysicalThingTableRawDataColumn:    {Name: PhysicalThingTableRawDataColumn, NotNull: false, HasDefault: false},
}

var (
	PhysicalThingTablePrimaryKeyColumn = PhysicalThingTableIDColumn
)
var _ = []any{
	time.Time{},
	time.Duration(0),
	nil,
	pq.StringArray{},
	string(""),
	pq.Int64Array{},
	int64(0),
	pq.Float64Array{},
	float64(0),
	pq.BoolArray{},
	bool(false),
	map[string][]int{},
	uuid.UUID{},
	hstore.Hstore{},
	pgtype.Point{},
	pgtype.Polygon{},
	postgis.PointZ{},
	netip.Prefix{},
	[]byte{},
	errors.Is,
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

	wrapError := func(k string, v any, err error) error {
		return fmt.Errorf("%v: %#+v; error: %v", k, v, err)
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

		}
	}

	return nil
}

func (m *PhysicalThing) Reload(
	ctx context.Context,
	tx *sqlx.Tx,
	includeDeleteds ...bool,
) error {
	extraWhere := ""
	if len(includeDeleteds) > 0 && includeDeleteds[0] {
		if slices.Contains(PhysicalThingTableColumns, "deleted_at") {
			extraWhere = "\n    AND (deleted_at IS null OR deleted_at IS NOT null)"
		}
	}

	ctx, cleanup := query.WithQueryID(ctx)
	defer cleanup()

	t, err := SelectPhysicalThing(
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
	m.ReferencedByLocationHistoryParentPhysicalThingIDObjects = t.ReferencedByLocationHistoryParentPhysicalThingIDObjects
	m.ReferencedByLogicalThingParentPhysicalThingIDObjects = t.ReferencedByLogicalThingParentPhysicalThingIDObjects

	return nil
}

func (m *PhysicalThing) Insert(
	ctx context.Context,
	tx *sqlx.Tx,
	setPrimaryKey bool,
	setZeroValues bool,
	forceSetValuesForFields ...string,
) error {
	columns := make([]string, 0)
	values := make([]any, 0)

	if setPrimaryKey && (setZeroValues || !types.IsZeroUUID(m.ID)) || slices.Contains(forceSetValuesForFields, PhysicalThingTableIDColumn) || isRequired(PhysicalThingTableColumnLookup, PhysicalThingTableIDColumn) {
		columns = append(columns, PhysicalThingTableIDColumn)

		v, err := types.FormatUUID(m.ID)
		if err != nil {
			return fmt.Errorf("failed to handle m.ID: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.CreatedAt) || slices.Contains(forceSetValuesForFields, PhysicalThingTableCreatedAtColumn) || isRequired(PhysicalThingTableColumnLookup, PhysicalThingTableCreatedAtColumn) {
		columns = append(columns, PhysicalThingTableCreatedAtColumn)

		v, err := types.FormatTime(m.CreatedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.CreatedAt: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.UpdatedAt) || slices.Contains(forceSetValuesForFields, PhysicalThingTableUpdatedAtColumn) || isRequired(PhysicalThingTableColumnLookup, PhysicalThingTableUpdatedAtColumn) {
		columns = append(columns, PhysicalThingTableUpdatedAtColumn)

		v, err := types.FormatTime(m.UpdatedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.UpdatedAt: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.DeletedAt) || slices.Contains(forceSetValuesForFields, PhysicalThingTableDeletedAtColumn) || isRequired(PhysicalThingTableColumnLookup, PhysicalThingTableDeletedAtColumn) {
		columns = append(columns, PhysicalThingTableDeletedAtColumn)

		v, err := types.FormatTime(m.DeletedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.DeletedAt: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroString(m.ExternalID) || slices.Contains(forceSetValuesForFields, PhysicalThingTableExternalIDColumn) || isRequired(PhysicalThingTableColumnLookup, PhysicalThingTableExternalIDColumn) {
		columns = append(columns, PhysicalThingTableExternalIDColumn)

		v, err := types.FormatString(m.ExternalID)
		if err != nil {
			return fmt.Errorf("failed to handle m.ExternalID: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroString(m.Name) || slices.Contains(forceSetValuesForFields, PhysicalThingTableNameColumn) || isRequired(PhysicalThingTableColumnLookup, PhysicalThingTableNameColumn) {
		columns = append(columns, PhysicalThingTableNameColumn)

		v, err := types.FormatString(m.Name)
		if err != nil {
			return fmt.Errorf("failed to handle m.Name: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroString(m.Type) || slices.Contains(forceSetValuesForFields, PhysicalThingTableTypeColumn) || isRequired(PhysicalThingTableColumnLookup, PhysicalThingTableTypeColumn) {
		columns = append(columns, PhysicalThingTableTypeColumn)

		v, err := types.FormatString(m.Type)
		if err != nil {
			return fmt.Errorf("failed to handle m.Type: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroStringArray(m.Tags) || slices.Contains(forceSetValuesForFields, PhysicalThingTableTagsColumn) || isRequired(PhysicalThingTableColumnLookup, PhysicalThingTableTagsColumn) {
		columns = append(columns, PhysicalThingTableTagsColumn)

		v, err := types.FormatStringArray(m.Tags)
		if err != nil {
			return fmt.Errorf("failed to handle m.Tags: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroHstore(m.Metadata) || slices.Contains(forceSetValuesForFields, PhysicalThingTableMetadataColumn) || isRequired(PhysicalThingTableColumnLookup, PhysicalThingTableMetadataColumn) {
		columns = append(columns, PhysicalThingTableMetadataColumn)

		v, err := types.FormatHstore(m.Metadata)
		if err != nil {
			return fmt.Errorf("failed to handle m.Metadata: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroJSON(m.RawData) || slices.Contains(forceSetValuesForFields, PhysicalThingTableRawDataColumn) || isRequired(PhysicalThingTableColumnLookup, PhysicalThingTableRawDataColumn) {
		columns = append(columns, PhysicalThingTableRawDataColumn)

		v, err := types.FormatJSON(m.RawData)
		if err != nil {
			return fmt.Errorf("failed to handle m.RawData: %v", err)
		}

		values = append(values, v)
	}

	ctx, cleanup := query.WithQueryID(ctx)
	defer cleanup()

	item, err := query.Insert(
		ctx,
		tx,
		PhysicalThingTable,
		columns,
		nil,
		false,
		false,
		PhysicalThingTableColumns,
		values...,
	)
	if err != nil {
		return fmt.Errorf("failed to insert %#+v: %v", m, err)
	}
	v := item[PhysicalThingTableIDColumn]

	if v == nil {
		return fmt.Errorf("failed to find %v in %#+v", PhysicalThingTableIDColumn, item)
	}

	wrapError := func(err error) error {
		return fmt.Errorf(
			"failed to treat %v: %#+v as uuid.UUID: %v",
			PhysicalThingTableIDColumn,
			item[PhysicalThingTableIDColumn],
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

	err = m.Reload(ctx, tx, slices.Contains(forceSetValuesForFields, "deleted_at"))
	if err != nil {
		return fmt.Errorf("failed to reload after insert")
	}

	return nil
}

func (m *PhysicalThing) Update(
	ctx context.Context,
	tx *sqlx.Tx,
	setZeroValues bool,
	forceSetValuesForFields ...string,
) error {
	columns := make([]string, 0)
	values := make([]any, 0)

	if setZeroValues || !types.IsZeroTime(m.CreatedAt) || slices.Contains(forceSetValuesForFields, PhysicalThingTableCreatedAtColumn) {
		columns = append(columns, PhysicalThingTableCreatedAtColumn)

		v, err := types.FormatTime(m.CreatedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.CreatedAt: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.UpdatedAt) || slices.Contains(forceSetValuesForFields, PhysicalThingTableUpdatedAtColumn) {
		columns = append(columns, PhysicalThingTableUpdatedAtColumn)

		v, err := types.FormatTime(m.UpdatedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.UpdatedAt: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.DeletedAt) || slices.Contains(forceSetValuesForFields, PhysicalThingTableDeletedAtColumn) {
		columns = append(columns, PhysicalThingTableDeletedAtColumn)

		v, err := types.FormatTime(m.DeletedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.DeletedAt: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroString(m.ExternalID) || slices.Contains(forceSetValuesForFields, PhysicalThingTableExternalIDColumn) {
		columns = append(columns, PhysicalThingTableExternalIDColumn)

		v, err := types.FormatString(m.ExternalID)
		if err != nil {
			return fmt.Errorf("failed to handle m.ExternalID: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroString(m.Name) || slices.Contains(forceSetValuesForFields, PhysicalThingTableNameColumn) {
		columns = append(columns, PhysicalThingTableNameColumn)

		v, err := types.FormatString(m.Name)
		if err != nil {
			return fmt.Errorf("failed to handle m.Name: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroString(m.Type) || slices.Contains(forceSetValuesForFields, PhysicalThingTableTypeColumn) {
		columns = append(columns, PhysicalThingTableTypeColumn)

		v, err := types.FormatString(m.Type)
		if err != nil {
			return fmt.Errorf("failed to handle m.Type: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroStringArray(m.Tags) || slices.Contains(forceSetValuesForFields, PhysicalThingTableTagsColumn) {
		columns = append(columns, PhysicalThingTableTagsColumn)

		v, err := types.FormatStringArray(m.Tags)
		if err != nil {
			return fmt.Errorf("failed to handle m.Tags: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroHstore(m.Metadata) || slices.Contains(forceSetValuesForFields, PhysicalThingTableMetadataColumn) {
		columns = append(columns, PhysicalThingTableMetadataColumn)

		v, err := types.FormatHstore(m.Metadata)
		if err != nil {
			return fmt.Errorf("failed to handle m.Metadata: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroJSON(m.RawData) || slices.Contains(forceSetValuesForFields, PhysicalThingTableRawDataColumn) {
		columns = append(columns, PhysicalThingTableRawDataColumn)

		v, err := types.FormatJSON(m.RawData)
		if err != nil {
			return fmt.Errorf("failed to handle m.RawData: %v", err)
		}

		values = append(values, v)
	}

	v, err := types.FormatUUID(m.ID)
	if err != nil {
		return fmt.Errorf("failed to handle m.ID: %v", err)
	}

	values = append(values, v)

	ctx, cleanup := query.WithQueryID(ctx)
	defer cleanup()

	_, err = query.Update(
		ctx,
		tx,
		PhysicalThingTable,
		columns,
		fmt.Sprintf("%v = $$??", PhysicalThingTableIDColumn),
		PhysicalThingTableColumns,
		values...,
	)
	if err != nil {
		return fmt.Errorf("failed to update %#+v: %v", m, err)
	}

	err = m.Reload(ctx, tx, slices.Contains(forceSetValuesForFields, "deleted_at"))
	if err != nil {
		return fmt.Errorf("failed to reload after update")
	}

	return nil
}

func (m *PhysicalThing) Delete(
	ctx context.Context,
	tx *sqlx.Tx,
	hardDeletes ...bool,
) error {
	hardDelete := false
	if len(hardDeletes) > 0 {
		hardDelete = hardDeletes[0]
	}

	if !hardDelete && slices.Contains(PhysicalThingTableColumns, "deleted_at") {
		m.DeletedAt = helpers.Ptr(time.Now().UTC())
		err := m.Update(ctx, tx, false, "deleted_at")
		if err != nil {
			return fmt.Errorf("failed to soft-delete (update) %#+v: %v", m, err)
		}
	}

	values := make([]any, 0)
	v, err := types.FormatUUID(m.ID)
	if err != nil {
		return fmt.Errorf("failed to handle m.ID: %v", err)
	}

	values = append(values, v)

	ctx, cleanup := query.WithQueryID(ctx)
	defer cleanup()

	err = query.Delete(
		ctx,
		tx,
		PhysicalThingTable,
		fmt.Sprintf("%v = $$??", PhysicalThingTableIDColumn),
		values...,
	)
	if err != nil {
		return fmt.Errorf("failed to delete %#+v: %v", m, err)
	}

	_ = m.Reload(ctx, tx, true)

	return nil
}

func SelectPhysicalThings(
	ctx context.Context,
	tx *sqlx.Tx,
	where string,
	orderBy *string,
	limit *int,
	offset *int,
	values ...any,
) ([]*PhysicalThing, error) {
	if slices.Contains(PhysicalThingTableColumns, "deleted_at") {
		if !strings.Contains(where, "deleted_at") {
			if where != "" {
				where += "\n    AND "
			}

			where += "deleted_at IS null"
		}
	}

	ctx, cleanup := query.WithQueryID(ctx)
	defer cleanup()

	items, err := query.Select(
		ctx,
		tx,
		PhysicalThingTableColumnsWithTypeCasts,
		PhysicalThingTable,
		where,
		orderBy,
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
			return nil, err
		}

		thatCtx, ok := query.HandleQueryPathGraphCycles(ctx, fmt.Sprintf("%s{%v}", PhysicalThingTable, object.ID))
		if !ok {
			continue
		}

		_ = thatCtx

		err = func() error {
			var ok bool
			thisCtx, ok := query.HandleQueryPathGraphCycles(thatCtx, fmt.Sprintf("%s{%v}", PhysicalThingTable, object.ID))

			if ok {
				object.ReferencedByLocationHistoryParentPhysicalThingIDObjects, err = SelectLocationHistories(
					thisCtx,
					tx,
					fmt.Sprintf("%v = $1", LocationHistoryTableParentPhysicalThingIDColumn),
					nil,
					nil,
					nil,
					object.ID,
				)
				if err != nil {
					if !errors.Is(err, sql.ErrNoRows) {
						return err
					}
				}
			}

			return nil
		}()
		if err != nil {
			return nil, err
		}

		err = func() error {
			var ok bool
			thisCtx, ok := query.HandleQueryPathGraphCycles(thatCtx, fmt.Sprintf("%s{%v}", PhysicalThingTable, object.ID))

			if ok {
				object.ReferencedByLogicalThingParentPhysicalThingIDObjects, err = SelectLogicalThings(
					thisCtx,
					tx,
					fmt.Sprintf("%v = $1", LogicalThingTableParentPhysicalThingIDColumn),
					nil,
					nil,
					nil,
					object.ID,
				)
				if err != nil {
					if !errors.Is(err, sql.ErrNoRows) {
						return err
					}
				}
			}

			return nil
		}()
		if err != nil {
			return nil, err
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
	ctx, cleanup := query.WithQueryID(ctx)
	defer cleanup()

	objects, err := SelectPhysicalThings(
		ctx,
		tx,
		where,
		nil,
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
		return nil, sql.ErrNoRows
	}

	object := objects[0]

	return object, nil
}

func handleGetPhysicalThings(w http.ResponseWriter, r *http.Request, db *sqlx.DB, redisPool *redis.Pool, objectMiddlewares []server.ObjectMiddleware) {
	ctx := r.Context()

	insaneOrderParams := make([]string, 0)
	hadInsaneOrderParams := false

	unrecognizedParams := make([]string, 0)
	hadUnrecognizedParams := false

	unparseableParams := make([]string, 0)
	hadUnparseableParams := false

	var orderByDirection *string
	orderBys := make([]string, 0)

	includes := make([]string, 0)

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
			column := PhysicalThingTableColumnLookup[parts[0]]
			if column == nil {
				if parts[0] != "load" {
					isUnrecognized = true
				}
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
				case "desc":
					if orderByDirection != nil && *orderByDirection != "DESC" {
						hadInsaneOrderParams = true
						insaneOrderParams = append(insaneOrderParams, rawKey)
						continue
					}

					orderByDirection = helpers.Ptr("DESC")
					orderBys = append(orderBys, parts[0])
					continue
				case "asc":
					if orderByDirection != nil && *orderByDirection != "ASC" {
						hadInsaneOrderParams = true
						insaneOrderParams = append(insaneOrderParams, rawKey)
						continue
					}

					orderByDirection = helpers.Ptr("ASC")
					orderBys = append(orderBys, parts[0])
					continue
				case "load":
					includes = append(includes, parts[0])
					_ = includes

					continue
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

	if hadInsaneOrderParams {
		helpers.HandleErrorResponse(
			w,
			http.StatusInternalServerError,
			fmt.Errorf("insane order params (e.g. conflicting asc / desc) %s", strings.Join(insaneOrderParams, ", ")),
		)
		return
	}

	limit := 2000
	rawLimit := r.URL.Query().Get("limit")
	if rawLimit != "" {
		possibleLimit, err := strconv.ParseInt(rawLimit, 10, 64)
		if err != nil {
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
		if err != nil {
			helpers.HandleErrorResponse(
				w,
				http.StatusInternalServerError,
				fmt.Errorf("failed to parse param offset=%s as int: %v", rawOffset, err),
			)
			return
		}

		offset = int(possibleOffset)
	}

	hashableOrderBy := ""
	var orderBy *string
	if len(orderBys) > 0 {
		hashableOrderBy = strings.Join(orderBys, ", ")
		if len(orderBys) > 1 {
			hashableOrderBy = fmt.Sprintf("(%v)", hashableOrderBy)
		}
		hashableOrderBy = fmt.Sprintf("%v %v", hashableOrderBy, *orderByDirection)
		orderBy = &hashableOrderBy
	}

	requestHash, err := helpers.GetRequestHash(PhysicalThingTable, wheres, hashableOrderBy, limit, offset, values, nil)
	if err != nil {
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	redisConn := redisPool.Get()
	defer func() {
		_ = redisConn.Close()
	}()

	cacheHit, err := helpers.AttemptCachedResponse(requestHash, redisConn, w)
	if err != nil {
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	if cacheHit {
		return
	}

	tx, err := db.BeginTxx(r.Context(), nil)
	if err != nil {
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	defer func() {
		_ = tx.Rollback()
	}()

	where := strings.Join(wheres, "\n    AND ")

	objects, err := SelectPhysicalThings(ctx, tx, where, orderBy, &limit, &offset, values...)
	if err != nil {
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	err = tx.Commit()
	if err != nil {
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	returnedObjectsAsJSON := helpers.HandleObjectsResponse(w, http.StatusOK, objects)

	err = helpers.StoreCachedResponse(requestHash, redisConn, string(returnedObjectsAsJSON))
	if err != nil {
		log.Printf("warning: %v", err)
	}
}

func handleGetPhysicalThing(w http.ResponseWriter, r *http.Request, db *sqlx.DB, redisPool *redis.Pool, objectMiddlewares []server.ObjectMiddleware, primaryKey string) {
	ctx := r.Context()

	wheres := []string{fmt.Sprintf("%s = $$??", PhysicalThingTablePrimaryKeyColumn)}
	values := []any{primaryKey}

	requestHash, err := helpers.GetRequestHash(PhysicalThingTable, wheres, "", 2, 0, values, primaryKey)
	if err != nil {
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	redisConn := redisPool.Get()
	defer func() {
		_ = redisConn.Close()
	}()

	cacheHit, err := helpers.AttemptCachedResponse(requestHash, redisConn, w)
	if err != nil {
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	if cacheHit {
		return
	}

	tx, err := db.BeginTxx(r.Context(), nil)
	if err != nil {
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	defer func() {
		_ = tx.Rollback()
	}()

	where := strings.Join(wheres, "\n    AND ")

	object, err := SelectPhysicalThing(ctx, tx, where, values...)
	if err != nil {
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	err = tx.Commit()
	if err != nil {
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	returnedObjectsAsJSON := helpers.HandleObjectsResponse(w, http.StatusOK, []*PhysicalThing{object})

	err = helpers.StoreCachedResponse(requestHash, redisConn, string(returnedObjectsAsJSON))
	if err != nil {
		log.Printf("warning: %v", err)
	}
}

func handlePostPhysicalThings(w http.ResponseWriter, r *http.Request, db *sqlx.DB, redisPool *redis.Pool, objectMiddlewares []server.ObjectMiddleware, waitForChange server.WaitForChange) {
	_ = redisPool

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

	forceSetValuesForFieldsByObjectIndex := make([][]string, 0)
	objects := make([]*PhysicalThing, 0)
	for _, item := range allItems {
		forceSetValuesForFields := make([]string, 0)
		for _, possibleField := range maps.Keys(item) {
			if !slices.Contains(PhysicalThingTableColumns, possibleField) {
				continue
			}

			forceSetValuesForFields = append(forceSetValuesForFields, possibleField)
		}
		forceSetValuesForFieldsByObjectIndex = append(forceSetValuesForFieldsByObjectIndex, forceSetValuesForFields)

		object := &PhysicalThing{}
		err = object.FromItem(item)
		if err != nil {
			err = fmt.Errorf("failed to interpret %#+v as PhysicalThing in item form: %v", item, err)
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

	xid, err := query.GetXid(r.Context(), tx)
	if err != nil {
		err = fmt.Errorf("failed to get xid: %v", err)
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}
	_ = xid

	for i, object := range objects {
		err = object.Insert(r.Context(), tx, false, false, forceSetValuesForFieldsByObjectIndex[i]...)
		if err != nil {
			err = fmt.Errorf("failed to insert %#+v: %v", object, err)
			helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
			return
		}

		objects[i] = object
	}

	errs := make(chan error, 1)
	go func() {
		_, err = waitForChange(r.Context(), []stream.Action{stream.INSERT}, PhysicalThingTable, xid)
		if err != nil {
			err = fmt.Errorf("failed to wait for change: %v", err)
			errs <- err
			return
		}

		errs <- nil
	}()

	err = tx.Commit()
	if err != nil {
		err = fmt.Errorf("failed to commit DB transaction: %v", err)
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	select {
	case <-r.Context().Done():
		err = fmt.Errorf("context canceled")
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	case err = <-errs:
		if err != nil {
			helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
			return
		}
	}

	helpers.HandleObjectsResponse(w, http.StatusCreated, objects)
}

func handlePutPhysicalThing(w http.ResponseWriter, r *http.Request, db *sqlx.DB, redisPool *redis.Pool, objectMiddlewares []server.ObjectMiddleware, waitForChange server.WaitForChange, primaryKey string) {
	_ = redisPool

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

	item[PhysicalThingTablePrimaryKeyColumn] = primaryKey

	object := &PhysicalThing{}
	err = object.FromItem(item)
	if err != nil {
		err = fmt.Errorf("failed to interpret %#+v as PhysicalThing in item form: %v", item, err)
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

	xid, err := query.GetXid(r.Context(), tx)
	if err != nil {
		err = fmt.Errorf("failed to get xid: %v", err)
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}
	_ = xid

	err = object.Update(r.Context(), tx, true)
	if err != nil {
		err = fmt.Errorf("failed to update %#+v: %v", object, err)
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	errs := make(chan error, 1)
	go func() {
		_, err = waitForChange(r.Context(), []stream.Action{stream.UPDATE, stream.SOFT_DELETE, stream.SOFT_RESTORE, stream.SOFT_UPDATE}, PhysicalThingTable, xid)
		if err != nil {
			err = fmt.Errorf("failed to wait for change: %v", err)
			errs <- err
			return
		}

		errs <- nil
	}()

	err = tx.Commit()
	if err != nil {
		err = fmt.Errorf("failed to commit DB transaction: %v", err)
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	select {
	case <-r.Context().Done():
		err = fmt.Errorf("context canceled")
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	case err = <-errs:
		if err != nil {
			helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
			return
		}
	}

	helpers.HandleObjectsResponse(w, http.StatusOK, []*PhysicalThing{object})
}

func handlePatchPhysicalThing(w http.ResponseWriter, r *http.Request, db *sqlx.DB, redisPool *redis.Pool, objectMiddlewares []server.ObjectMiddleware, waitForChange server.WaitForChange, primaryKey string) {
	_ = redisPool

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

	forceSetValuesForFields := make([]string, 0)
	for _, possibleField := range maps.Keys(item) {
		if !slices.Contains(PhysicalThingTableColumns, possibleField) {
			continue
		}

		forceSetValuesForFields = append(forceSetValuesForFields, possibleField)
	}

	item[PhysicalThingTablePrimaryKeyColumn] = primaryKey

	object := &PhysicalThing{}
	err = object.FromItem(item)
	if err != nil {
		err = fmt.Errorf("failed to interpret %#+v as PhysicalThing in item form: %v", item, err)
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

	xid, err := query.GetXid(r.Context(), tx)
	if err != nil {
		err = fmt.Errorf("failed to get xid: %v", err)
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}
	_ = xid

	err = object.Update(r.Context(), tx, false, forceSetValuesForFields...)
	if err != nil {
		err = fmt.Errorf("failed to update %#+v: %v", object, err)
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	errs := make(chan error, 1)
	go func() {
		_, err = waitForChange(r.Context(), []stream.Action{stream.UPDATE, stream.SOFT_DELETE, stream.SOFT_RESTORE, stream.SOFT_UPDATE}, PhysicalThingTable, xid)
		if err != nil {
			err = fmt.Errorf("failed to wait for change: %v", err)
			errs <- err
			return
		}

		errs <- nil
	}()

	err = tx.Commit()
	if err != nil {
		err = fmt.Errorf("failed to commit DB transaction: %v", err)
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	select {
	case <-r.Context().Done():
		err = fmt.Errorf("context canceled")
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	case err = <-errs:
		if err != nil {
			helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
			return
		}
	}

	helpers.HandleObjectsResponse(w, http.StatusOK, []*PhysicalThing{object})
}

func handleDeletePhysicalThing(w http.ResponseWriter, r *http.Request, db *sqlx.DB, redisPool *redis.Pool, objectMiddlewares []server.ObjectMiddleware, waitForChange server.WaitForChange, primaryKey string) {
	_ = redisPool

	var item = make(map[string]any)

	item[PhysicalThingTablePrimaryKeyColumn] = primaryKey

	object := &PhysicalThing{}
	err := object.FromItem(item)
	if err != nil {
		err = fmt.Errorf("failed to interpret %#+v as PhysicalThing in item form: %v", item, err)
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

	xid, err := query.GetXid(r.Context(), tx)
	if err != nil {
		err = fmt.Errorf("failed to get xid: %v", err)
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}
	_ = xid

	err = object.Delete(r.Context(), tx)
	if err != nil {
		err = fmt.Errorf("failed to delete %#+v: %v", object, err)
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	errs := make(chan error, 1)
	go func() {
		_, err = waitForChange(r.Context(), []stream.Action{stream.DELETE, stream.SOFT_DELETE}, PhysicalThingTable, xid)
		if err != nil {
			err = fmt.Errorf("failed to wait for change: %v", err)
			errs <- err
			return
		}

		errs <- nil
	}()

	err = tx.Commit()
	if err != nil {
		err = fmt.Errorf("failed to commit DB transaction: %v", err)
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	select {
	case <-r.Context().Done():
		err = fmt.Errorf("context canceled")
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	case err = <-errs:
		if err != nil {
			helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
			return
		}
	}

	helpers.HandleObjectsResponse(w, http.StatusNoContent, nil)
}

func GetPhysicalThingRouter(db *sqlx.DB, redisPool *redis.Pool, httpMiddlewares []server.HTTPMiddleware, objectMiddlewares []server.ObjectMiddleware, waitForChange server.WaitForChange) chi.Router {
	r := chi.NewRouter()

	for _, m := range httpMiddlewares {
		r.Use(m)
	}

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		handleGetPhysicalThings(w, r, db, redisPool, objectMiddlewares)
	})

	r.Get("/{primaryKey}", func(w http.ResponseWriter, r *http.Request) {
		handleGetPhysicalThing(w, r, db, redisPool, objectMiddlewares, chi.URLParam(r, "primaryKey"))
	})

	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		handlePostPhysicalThings(w, r, db, redisPool, objectMiddlewares, waitForChange)
	})

	r.Put("/{primaryKey}", func(w http.ResponseWriter, r *http.Request) {
		handlePutPhysicalThing(w, r, db, redisPool, objectMiddlewares, waitForChange, chi.URLParam(r, "primaryKey"))
	})

	r.Patch("/{primaryKey}", func(w http.ResponseWriter, r *http.Request) {
		handlePatchPhysicalThing(w, r, db, redisPool, objectMiddlewares, waitForChange, chi.URLParam(r, "primaryKey"))
	})

	r.Delete("/{primaryKey}", func(w http.ResponseWriter, r *http.Request) {
		handleDeletePhysicalThing(w, r, db, redisPool, objectMiddlewares, waitForChange, chi.URLParam(r, "primaryKey"))
	})

	return r
}

func NewPhysicalThingFromItem(item map[string]any) (any, error) {
	object := &PhysicalThing{}

	err := object.FromItem(item)
	if err != nil {
		return nil, err
	}

	return object, nil
}

func init() {
	register(
		PhysicalThingTable,
		PhysicalThing{},
		NewPhysicalThingFromItem,
		"/physical-things",
		GetPhysicalThingRouter,
	)
}
