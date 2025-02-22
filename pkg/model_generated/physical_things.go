package model_generated

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"net/http"
	"net/netip"
	"slices"
	"strings"
	"time"

	"github.com/cridenour/go-postgis"
	"github.com/go-chi/chi/v5"
	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
	"github.com/initialed85/djangolang/pkg/config"
	"github.com/initialed85/djangolang/pkg/helpers"
	"github.com/initialed85/djangolang/pkg/introspect"
	"github.com/initialed85/djangolang/pkg/query"
	"github.com/initialed85/djangolang/pkg/server"
	"github.com/initialed85/djangolang/pkg/stream"
	"github.com/initialed85/djangolang/pkg/types"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
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

var PhysicalThingTableWithSchema = fmt.Sprintf("%s.%s", schema, PhysicalThingTable)

var PhysicalThingTableNamespaceID int32 = 1337 + 7

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
	PhysicalThingTableIDColumnWithTypeCast         = `"id" AS id`
	PhysicalThingTableCreatedAtColumnWithTypeCast  = `"created_at" AS created_at`
	PhysicalThingTableUpdatedAtColumnWithTypeCast  = `"updated_at" AS updated_at`
	PhysicalThingTableDeletedAtColumnWithTypeCast  = `"deleted_at" AS deleted_at`
	PhysicalThingTableExternalIDColumnWithTypeCast = `"external_id" AS external_id`
	PhysicalThingTableNameColumnWithTypeCast       = `"name" AS name`
	PhysicalThingTableTypeColumnWithTypeCast       = `"type" AS type`
	PhysicalThingTableTagsColumnWithTypeCast       = `"tags" AS tags`
	PhysicalThingTableMetadataColumnWithTypeCast   = `"metadata" AS metadata`
	PhysicalThingTableRawDataColumnWithTypeCast    = `"raw_data" AS raw_data`
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

var PhysicalThingIntrospectedTable *introspect.Table

var PhysicalThingTableColumnLookup map[string]*introspect.Column

var (
	PhysicalThingTablePrimaryKeyColumn = PhysicalThingTableIDColumn
)

func init() {
	PhysicalThingIntrospectedTable = tableByName[PhysicalThingTable]

	/* only needed during templating */
	if PhysicalThingIntrospectedTable == nil {
		PhysicalThingIntrospectedTable = &introspect.Table{}
	}

	PhysicalThingTableColumnLookup = PhysicalThingIntrospectedTable.ColumnByName
}

type PhysicalThingOnePathParams struct {
	PrimaryKey uuid.UUID `json:"primaryKey"`
}

type PhysicalThingLoadQueryParams struct {
	Depth *int `json:"depth"`
}

/*
TODO: find a way to not need this- there is a piece in the templating logic
that uses goimports but pending where the code is built, it may resolve
the packages to import to the wrong ones (causing odd failures)
these are just here to ensure we don't get unused imports
*/
var _ = []any{
	time.Time{},
	uuid.UUID{},
	pgtype.Hstore{},
	postgis.PointZ{},
	netip.Prefix{},
	errors.Is,
	sql.ErrNoRows,
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
		return fmt.Errorf("%v: %#+v; error; %v", k, v, err)
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
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uucreated_at.UUID", temp1))
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
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uuupdated_at.UUID", temp1))
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
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uudeleted_at.UUID", temp1))
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
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uuexternal_id.UUID", temp1))
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
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uuname.UUID", temp1))
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
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uutype.UUID", temp1))
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
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uutags.UUID", temp1))
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
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uumetadata.UUID", temp1))
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

			temp2, ok := temp1, true
			if !ok {
				if temp1 != nil {
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uuraw_data.UUID", temp1))
				}
			}

			m.RawData = &temp2

		}
	}

	return nil
}

func (m *PhysicalThing) ToItem() map[string]any {
	item := make(map[string]any)

	b, err := json.Marshal(m)
	if err != nil {
		panic(fmt.Sprintf("%T.ToItem() failed intermediate marshal to JSON: %s", m, err))
	}

	err = json.Unmarshal(b, &item)
	if err != nil {
		panic(fmt.Sprintf("%T.ToItem() failed intermediate unmarshal from JSON: %s", m, err))
	}

	return item
}

func (m *PhysicalThing) Reload(ctx context.Context, tx pgx.Tx, includeDeleteds ...bool) error {
	extraWhere := ""
	if len(includeDeleteds) > 0 && includeDeleteds[0] {
		if slices.Contains(PhysicalThingTableColumns, "deleted_at") {
			extraWhere = "\n    AND (deleted_at IS null OR deleted_at IS NOT null)"
		}
	}

	ctx, cleanup := query.WithQueryID(ctx)
	defer cleanup()

	ctx = query.WithMaxDepth(ctx, nil)

	o, _, _, _, _, err := SelectPhysicalThing(
		ctx,
		tx,
		fmt.Sprintf("%v = $1%v", m.GetPrimaryKeyColumn(), extraWhere),
		m.GetPrimaryKeyValue(),
	)
	if err != nil {
		return err
	}

	m.ID = o.ID
	m.CreatedAt = o.CreatedAt
	m.UpdatedAt = o.UpdatedAt
	m.DeletedAt = o.DeletedAt
	m.ExternalID = o.ExternalID
	m.Name = o.Name
	m.Type = o.Type
	m.Tags = o.Tags
	m.Metadata = o.Metadata
	m.RawData = o.RawData
	m.ReferencedByLocationHistoryParentPhysicalThingIDObjects = o.ReferencedByLocationHistoryParentPhysicalThingIDObjects
	m.ReferencedByLogicalThingParentPhysicalThingIDObjects = o.ReferencedByLogicalThingParentPhysicalThingIDObjects

	return nil
}

func (m *PhysicalThing) GetColumnsAndValues(setPrimaryKey bool, setZeroValues bool, forceSetValuesForFields ...string) ([]string, []any, error) {
	columns := make([]string, 0)
	values := make([]any, 0)

	if setPrimaryKey && (setZeroValues || !types.IsZeroUUID(m.ID) || slices.Contains(forceSetValuesForFields, PhysicalThingTableIDColumn) || isRequired(PhysicalThingTableColumnLookup, PhysicalThingTableIDColumn)) {
		columns = append(columns, PhysicalThingTableIDColumn)

		v, err := types.FormatUUID(m.ID)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to handle m.ID; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.CreatedAt) || slices.Contains(forceSetValuesForFields, PhysicalThingTableCreatedAtColumn) || isRequired(PhysicalThingTableColumnLookup, PhysicalThingTableCreatedAtColumn) {
		columns = append(columns, PhysicalThingTableCreatedAtColumn)

		v, err := types.FormatTime(m.CreatedAt)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to handle m.CreatedAt; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.UpdatedAt) || slices.Contains(forceSetValuesForFields, PhysicalThingTableUpdatedAtColumn) || isRequired(PhysicalThingTableColumnLookup, PhysicalThingTableUpdatedAtColumn) {
		columns = append(columns, PhysicalThingTableUpdatedAtColumn)

		v, err := types.FormatTime(m.UpdatedAt)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to handle m.UpdatedAt; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.DeletedAt) || slices.Contains(forceSetValuesForFields, PhysicalThingTableDeletedAtColumn) || isRequired(PhysicalThingTableColumnLookup, PhysicalThingTableDeletedAtColumn) {
		columns = append(columns, PhysicalThingTableDeletedAtColumn)

		v, err := types.FormatTime(m.DeletedAt)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to handle m.DeletedAt; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroString(m.ExternalID) || slices.Contains(forceSetValuesForFields, PhysicalThingTableExternalIDColumn) || isRequired(PhysicalThingTableColumnLookup, PhysicalThingTableExternalIDColumn) {
		columns = append(columns, PhysicalThingTableExternalIDColumn)

		v, err := types.FormatString(m.ExternalID)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to handle m.ExternalID; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroString(m.Name) || slices.Contains(forceSetValuesForFields, PhysicalThingTableNameColumn) || isRequired(PhysicalThingTableColumnLookup, PhysicalThingTableNameColumn) {
		columns = append(columns, PhysicalThingTableNameColumn)

		v, err := types.FormatString(m.Name)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to handle m.Name; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroString(m.Type) || slices.Contains(forceSetValuesForFields, PhysicalThingTableTypeColumn) || isRequired(PhysicalThingTableColumnLookup, PhysicalThingTableTypeColumn) {
		columns = append(columns, PhysicalThingTableTypeColumn)

		v, err := types.FormatString(m.Type)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to handle m.Type; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroStringArray(m.Tags) || slices.Contains(forceSetValuesForFields, PhysicalThingTableTagsColumn) || isRequired(PhysicalThingTableColumnLookup, PhysicalThingTableTagsColumn) {
		columns = append(columns, PhysicalThingTableTagsColumn)

		v, err := types.FormatStringArray(m.Tags)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to handle m.Tags; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroHstore(m.Metadata) || slices.Contains(forceSetValuesForFields, PhysicalThingTableMetadataColumn) || isRequired(PhysicalThingTableColumnLookup, PhysicalThingTableMetadataColumn) {
		columns = append(columns, PhysicalThingTableMetadataColumn)

		v, err := types.FormatHstore(m.Metadata)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to handle m.Metadata; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroJSON(m.RawData) || slices.Contains(forceSetValuesForFields, PhysicalThingTableRawDataColumn) || isRequired(PhysicalThingTableColumnLookup, PhysicalThingTableRawDataColumn) {
		columns = append(columns, PhysicalThingTableRawDataColumn)

		v, err := types.FormatJSON(m.RawData)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to handle m.RawData; %v", err)
		}

		values = append(values, v)
	}

	return columns, values, nil
}

func (m *PhysicalThing) Insert(ctx context.Context, tx pgx.Tx, setPrimaryKey bool, setZeroValues bool, forceSetValuesForFields ...string) error {
	columns, values, err := m.GetColumnsAndValues(setPrimaryKey, setZeroValues, forceSetValuesForFields...)
	if err != nil {
		return fmt.Errorf("failed to get columns and values to insert %#+v; %v", m, err)
	}

	ctx, cleanup := query.WithQueryID(ctx)
	defer cleanup()

	ctx = query.WithMaxDepth(ctx, nil)

	item, err := query.Insert(
		ctx,
		tx,
		PhysicalThingTableWithSchema,
		columns,
		nil,
		nil,
		nil,
		false,
		false,
		PhysicalThingTableColumns,
		values...,
	)
	if err != nil {
		return fmt.Errorf("failed to insert %#+v; %v", m, err)
	}
	v := (*item)[PhysicalThingTableIDColumn]

	if v == nil {
		return fmt.Errorf("failed to find %v in %#+v", PhysicalThingTableIDColumn, item)
	}

	wrapError := func(err error) error {
		return fmt.Errorf(
			"failed to treat %v: %#+v as uuid.UUID: %v",
			PhysicalThingTableIDColumn,
			(*item)[PhysicalThingTableIDColumn],
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
		return fmt.Errorf("failed to reload after insert; %v", err)
	}

	return nil
}

func (m *PhysicalThing) Update(ctx context.Context, tx pgx.Tx, setZeroValues bool, forceSetValuesForFields ...string) error {
	columns := make([]string, 0)
	values := make([]any, 0)

	if setZeroValues || !types.IsZeroTime(m.CreatedAt) || slices.Contains(forceSetValuesForFields, PhysicalThingTableCreatedAtColumn) {
		columns = append(columns, PhysicalThingTableCreatedAtColumn)

		v, err := types.FormatTime(m.CreatedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.CreatedAt; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.UpdatedAt) || slices.Contains(forceSetValuesForFields, PhysicalThingTableUpdatedAtColumn) {
		columns = append(columns, PhysicalThingTableUpdatedAtColumn)

		v, err := types.FormatTime(m.UpdatedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.UpdatedAt; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.DeletedAt) || slices.Contains(forceSetValuesForFields, PhysicalThingTableDeletedAtColumn) {
		columns = append(columns, PhysicalThingTableDeletedAtColumn)

		v, err := types.FormatTime(m.DeletedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.DeletedAt; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroString(m.ExternalID) || slices.Contains(forceSetValuesForFields, PhysicalThingTableExternalIDColumn) {
		columns = append(columns, PhysicalThingTableExternalIDColumn)

		v, err := types.FormatString(m.ExternalID)
		if err != nil {
			return fmt.Errorf("failed to handle m.ExternalID; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroString(m.Name) || slices.Contains(forceSetValuesForFields, PhysicalThingTableNameColumn) {
		columns = append(columns, PhysicalThingTableNameColumn)

		v, err := types.FormatString(m.Name)
		if err != nil {
			return fmt.Errorf("failed to handle m.Name; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroString(m.Type) || slices.Contains(forceSetValuesForFields, PhysicalThingTableTypeColumn) {
		columns = append(columns, PhysicalThingTableTypeColumn)

		v, err := types.FormatString(m.Type)
		if err != nil {
			return fmt.Errorf("failed to handle m.Type; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroStringArray(m.Tags) || slices.Contains(forceSetValuesForFields, PhysicalThingTableTagsColumn) {
		columns = append(columns, PhysicalThingTableTagsColumn)

		v, err := types.FormatStringArray(m.Tags)
		if err != nil {
			return fmt.Errorf("failed to handle m.Tags; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroHstore(m.Metadata) || slices.Contains(forceSetValuesForFields, PhysicalThingTableMetadataColumn) {
		columns = append(columns, PhysicalThingTableMetadataColumn)

		v, err := types.FormatHstore(m.Metadata)
		if err != nil {
			return fmt.Errorf("failed to handle m.Metadata; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroJSON(m.RawData) || slices.Contains(forceSetValuesForFields, PhysicalThingTableRawDataColumn) {
		columns = append(columns, PhysicalThingTableRawDataColumn)

		v, err := types.FormatJSON(m.RawData)
		if err != nil {
			return fmt.Errorf("failed to handle m.RawData; %v", err)
		}

		values = append(values, v)
	}

	v, err := types.FormatUUID(m.ID)
	if err != nil {
		return fmt.Errorf("failed to handle m.ID; %v", err)
	}

	values = append(values, v)

	ctx, cleanup := query.WithQueryID(ctx)
	defer cleanup()

	ctx = query.WithMaxDepth(ctx, nil)

	_, err = query.Update(
		ctx,
		tx,
		PhysicalThingTableWithSchema,
		columns,
		fmt.Sprintf("%v = $$??", PhysicalThingTableIDColumn),
		PhysicalThingTableColumns,
		values...,
	)
	if err != nil {
		return fmt.Errorf("failed to update %#+v; %v", m, err)
	}

	err = m.Reload(ctx, tx, slices.Contains(forceSetValuesForFields, "deleted_at"))
	if err != nil {
		return fmt.Errorf("failed to reload after update")
	}

	return nil
}

func (m *PhysicalThing) Delete(ctx context.Context, tx pgx.Tx, hardDeletes ...bool) error {
	hardDelete := false
	if len(hardDeletes) > 0 {
		hardDelete = hardDeletes[0]
	}

	if !hardDelete && slices.Contains(PhysicalThingTableColumns, "deleted_at") {
		m.DeletedAt = helpers.Ptr(time.Now().UTC())
		err := m.Update(ctx, tx, false, "deleted_at")
		if err != nil {
			return fmt.Errorf("failed to soft-delete (update) %#+v; %v", m, err)
		}
	}

	values := make([]any, 0)
	v, err := types.FormatUUID(m.ID)
	if err != nil {
		return fmt.Errorf("failed to handle m.ID; %v", err)
	}

	values = append(values, v)

	ctx, cleanup := query.WithQueryID(ctx)
	defer cleanup()

	ctx = query.WithMaxDepth(ctx, nil)

	err = query.Delete(
		ctx,
		tx,
		PhysicalThingTableWithSchema,
		fmt.Sprintf("%v = $$??", PhysicalThingTableIDColumn),
		values...,
	)
	if err != nil {
		return fmt.Errorf("failed to delete %#+v; %v", m, err)
	}

	_ = m.Reload(ctx, tx, true)

	return nil
}

func (m *PhysicalThing) LockTable(ctx context.Context, tx pgx.Tx, timeouts ...time.Duration) error {
	return query.LockTable(ctx, tx, PhysicalThingTableWithSchema, timeouts...)
}

func (m *PhysicalThing) LockTableWithRetries(ctx context.Context, tx pgx.Tx, overallTimeout time.Duration, individualAttempttimeout time.Duration) error {
	return query.LockTableWithRetries(ctx, tx, PhysicalThingTableWithSchema, overallTimeout, individualAttempttimeout)
}

func (m *PhysicalThing) AdvisoryLock(ctx context.Context, tx pgx.Tx, key int32, timeouts ...time.Duration) error {
	return query.AdvisoryLock(ctx, tx, PhysicalThingTableNamespaceID, key, timeouts...)
}

func (m *PhysicalThing) AdvisoryLockWithRetries(ctx context.Context, tx pgx.Tx, key int32, overallTimeout time.Duration, individualAttempttimeout time.Duration) error {
	return query.AdvisoryLockWithRetries(ctx, tx, PhysicalThingTableNamespaceID, key, overallTimeout, individualAttempttimeout)
}

func SelectPhysicalThings(ctx context.Context, tx pgx.Tx, where string, orderBy *string, limit *int, offset *int, values ...any) ([]*PhysicalThing, int64, int64, int64, int64, error) {
	before := time.Now()

	if config.Debug() {
		log.Printf("entered SelectPhysicalThings")

		defer func() {
			log.Printf("exited SelectPhysicalThings in %s", time.Since(before))
		}()
	}
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

	possiblePathValue := query.GetCurrentPathValue(ctx)
	isLoadQuery := possiblePathValue != nil && len(possiblePathValue.VisitedTableNames) > 0

	shouldLoad := query.ShouldLoad(ctx, PhysicalThingTable) || query.ShouldLoad(ctx, fmt.Sprintf("referenced_by_%s", PhysicalThingTable))

	var ok bool
	ctx, ok = query.HandleQueryPathGraphCycles(ctx, fmt.Sprintf("%s{%v}", PhysicalThingTable, nil), !isLoadQuery)
	if !ok && !shouldLoad {
		if config.Debug() {
			log.Printf("skipping SelectPhysicalThing early (query.ShouldLoad(): %v, query.HandleQueryPathGraphCycles(): %v)", shouldLoad, ok)
		}
		return []*PhysicalThing{}, 0, 0, 0, 0, nil
	}

	var items *[]map[string]any
	var count int64
	var totalCount int64
	var page int64
	var totalPages int64
	var err error

	useInstead, shouldSkip := query.ShouldSkip[PhysicalThing](ctx)
	if !shouldSkip {
		items, count, totalCount, page, totalPages, err = query.Select(
			ctx,
			tx,
			PhysicalThingTableColumnsWithTypeCasts,
			PhysicalThingTableWithSchema,
			where,
			orderBy,
			limit,
			offset,
			values...,
		)
		if err != nil {
			return nil, 0, 0, 0, 0, fmt.Errorf("failed to call SelectPhysicalThings; %v", err)
		}
	} else {
		ctx = query.WithoutSkip(ctx)
		count = 1
		totalCount = 1
		page = 1
		totalPages = 1
		items = &[]map[string]any{
			nil,
		}
	}

	objects := make([]*PhysicalThing, 0)

	for _, item := range *items {
		var object *PhysicalThing

		if !shouldSkip {
			object = &PhysicalThing{}
			err = object.FromItem(item)
			if err != nil {
				return nil, 0, 0, 0, 0, err
			}
		} else {
			object = useInstead
		}

		if object == nil {
			return nil, 0, 0, 0, 0, fmt.Errorf("assertion failed: object unexpectedly nil")
		}

		err = func() error {
			shouldLoad := query.ShouldLoad(ctx, fmt.Sprintf("referenced_by_%s", LocationHistoryTable))
			ctx, ok := query.HandleQueryPathGraphCycles(ctx, fmt.Sprintf("__ReferencedBy__%s{%v}", LocationHistoryTable, object.GetPrimaryKeyValue()), true)
			if ok || shouldLoad {
				thisBefore := time.Now()

				if config.Debug() {
					log.Printf("loading SelectPhysicalThings->SelectLocationHistories for object.ReferencedByLocationHistoryParentPhysicalThingIDObjects")
				}

				object.ReferencedByLocationHistoryParentPhysicalThingIDObjects, _, _, _, _, err = SelectLocationHistories(
					ctx,
					tx,
					fmt.Sprintf("%v = $1", LocationHistoryTableParentPhysicalThingIDColumn),
					nil,
					nil,
					nil,
					object.GetPrimaryKeyValue(),
				)
				if err != nil {
					if !errors.Is(err, sql.ErrNoRows) {
						return err
					}
				}

				if config.Debug() {
					log.Printf("loaded SelectPhysicalThings->SelectLocationHistories for object.ReferencedByLocationHistoryParentPhysicalThingIDObjects in %s", time.Since(thisBefore))
				}

			}

			return nil
		}()
		if err != nil {
			return nil, 0, 0, 0, 0, err
		}

		err = func() error {
			shouldLoad := query.ShouldLoad(ctx, fmt.Sprintf("referenced_by_%s", PhysicalThingTable))
			ctx, ok := query.HandleQueryPathGraphCycles(ctx, fmt.Sprintf("__ReferencedBy__%s{%v}", PhysicalThingTable, object.GetPrimaryKeyValue()), true)
			if ok || shouldLoad {
				thisBefore := time.Now()

				if config.Debug() {
					log.Printf("loading SelectPhysicalThings->SelectLogicalThings for object.ReferencedByLogicalThingParentPhysicalThingIDObjects")
				}

				object.ReferencedByLogicalThingParentPhysicalThingIDObjects, _, _, _, _, err = SelectLogicalThings(
					ctx,
					tx,
					fmt.Sprintf("%v = $1", LogicalThingTableParentPhysicalThingIDColumn),
					nil,
					nil,
					nil,
					object.GetPrimaryKeyValue(),
				)
				if err != nil {
					if !errors.Is(err, sql.ErrNoRows) {
						return err
					}
				}

				if config.Debug() {
					log.Printf("loaded SelectPhysicalThings->SelectLogicalThings for object.ReferencedByLogicalThingParentPhysicalThingIDObjects in %s", time.Since(thisBefore))
				}

			}

			return nil
		}()
		if err != nil {
			return nil, 0, 0, 0, 0, err
		}

		objects = append(objects, object)
	}

	return objects, count, totalCount, page, totalPages, nil
}

func SelectPhysicalThing(ctx context.Context, tx pgx.Tx, where string, values ...any) (*PhysicalThing, int64, int64, int64, int64, error) {
	ctx, cleanup := query.WithQueryID(ctx)
	defer cleanup()

	ctx = query.WithMaxDepth(ctx, nil)

	objects, _, _, _, _, err := SelectPhysicalThings(
		ctx,
		tx,
		where,
		nil,
		helpers.Ptr(2),
		helpers.Ptr(0),
		values...,
	)
	if err != nil {
		return nil, 0, 0, 0, 0, fmt.Errorf("failed to call SelectPhysicalThing; %v", err)
	}

	if len(objects) > 1 {
		return nil, 0, 0, 0, 0, fmt.Errorf("attempt to call SelectPhysicalThing returned more than 1 row")
	}

	if len(objects) < 1 {
		return nil, 0, 0, 0, 0, sql.ErrNoRows
	}

	object := objects[0]

	count := int64(1)
	totalCount := count
	page := int64(1)
	totalPages := page

	return object, count, totalCount, page, totalPages, nil
}

func InsertPhysicalThings(ctx context.Context, tx pgx.Tx, objects []*PhysicalThing, setPrimaryKey bool, setZeroValues bool, forceSetValuesForFields ...string) ([]*PhysicalThing, error) {
	var columns []string
	values := make([]any, 0)

	for i, object := range objects {
		thisColumns, thisValues, err := object.GetColumnsAndValues(setPrimaryKey, setZeroValues, forceSetValuesForFields...)
		if err != nil {
			return nil, err
		}

		if columns == nil {
			columns = thisColumns
		} else {
			if len(columns) != len(thisColumns) {
				return nil, fmt.Errorf(
					"assertion failed: call 1 of object.GetColumnsAndValues() gave %d columns but call %d gave %d columns",
					len(columns),
					i+1,
					len(thisColumns),
				)
			}
		}

		values = append(values, thisValues...)
	}

	ctx, cleanup := query.WithQueryID(ctx)
	defer cleanup()

	ctx = query.WithMaxDepth(ctx, nil)

	items, err := query.BulkInsert(
		ctx,
		tx,
		PhysicalThingTableWithSchema,
		columns,
		nil,
		nil,
		nil,
		false,
		false,
		PhysicalThingTableColumns,
		values...,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to bulk insert %d objects; %v", len(objects), err)
	}

	returnedObjects := make([]*PhysicalThing, 0)

	for _, item := range items {
		v := &PhysicalThing{}
		err = v.FromItem(*item)
		if err != nil {
			return nil, fmt.Errorf("failed %T.FromItem for %#+v; %v", *item, *item, err)
		}

		err = v.Reload(query.WithSkip(ctx, v), tx)
		if err != nil {
			return nil, fmt.Errorf("failed %T.Reload for %#+v; %v", *item, *item, err)
		}

		returnedObjects = append(returnedObjects, v)
	}

	return returnedObjects, nil
}

func handleGetPhysicalThings(arguments *server.SelectManyArguments, db *pgxpool.Pool) ([]*PhysicalThing, int64, int64, int64, int64, error) {
	tx, err := db.Begin(arguments.Ctx)
	if err != nil {
		return nil, 0, 0, 0, 0, err
	}

	defer func() {
		_ = tx.Rollback(arguments.Ctx)
	}()

	objects, count, totalCount, page, totalPages, err := SelectPhysicalThings(arguments.Ctx, tx, arguments.Where, arguments.OrderBy, arguments.Limit, arguments.Offset, arguments.Values...)
	if err != nil {
		return nil, 0, 0, 0, 0, err
	}

	err = tx.Commit(arguments.Ctx)
	if err != nil {
		return nil, 0, 0, 0, 0, err
	}

	return objects, count, totalCount, page, totalPages, nil
}

func handleGetPhysicalThing(arguments *server.SelectOneArguments, db *pgxpool.Pool, primaryKey uuid.UUID) ([]*PhysicalThing, int64, int64, int64, int64, error) {
	tx, err := db.Begin(arguments.Ctx)
	if err != nil {
		return nil, 0, 0, 0, 0, err
	}

	defer func() {
		_ = tx.Rollback(arguments.Ctx)
	}()

	object, count, totalCount, page, totalPages, err := SelectPhysicalThing(arguments.Ctx, tx, arguments.Where, arguments.Values...)
	if err != nil {
		return nil, 0, 0, 0, 0, err
	}

	err = tx.Commit(arguments.Ctx)
	if err != nil {
		return nil, 0, 0, 0, 0, err
	}

	return []*PhysicalThing{object}, count, totalCount, page, totalPages, nil
}

func handlePostPhysicalThing(arguments *server.LoadArguments, db *pgxpool.Pool, waitForChange server.WaitForChange, objects []*PhysicalThing, forceSetValuesForFieldsByObjectIndex [][]string) ([]*PhysicalThing, int64, int64, int64, int64, error) {
	tx, err := db.Begin(arguments.Ctx)
	if err != nil {
		err = fmt.Errorf("failed to begin DB transaction; %v", err)
		return nil, 0, 0, 0, 0, err
	}

	defer func() {
		_ = tx.Rollback(arguments.Ctx)
	}()

	xid, err := query.GetXid(arguments.Ctx, tx)
	if err != nil {
		err = fmt.Errorf("failed to get xid; %v", err)
		return nil, 0, 0, 0, 0, err
	}

	/* TODO: problematic- basically the bulks insert insists all rows have the same schema, which they usually should */
	forceSetValuesForFieldsByObjectIndexMaximal := make(map[string]struct{})
	for _, forceSetforceSetValuesForFields := range forceSetValuesForFieldsByObjectIndex {
		for _, field := range forceSetforceSetValuesForFields {
			forceSetValuesForFieldsByObjectIndexMaximal[field] = struct{}{}
		}
	}

	returnedObjects, err := InsertPhysicalThings(arguments.Ctx, tx, objects, false, false, slices.Collect(maps.Keys(forceSetValuesForFieldsByObjectIndexMaximal))...)
	if err != nil {
		err = fmt.Errorf("failed to insert %d objects; %v", len(objects), err)
		return nil, 0, 0, 0, 0, err
	}

	copy(objects, returnedObjects)

	errs := make(chan error, 1)
	go func() {
		_, err := waitForChange(arguments.Ctx, []stream.Action{stream.INSERT}, PhysicalThingTable, xid)
		if err != nil {
			err = fmt.Errorf("failed to wait for change; %v", err)
			errs <- err
			return
		}

		errs <- nil
	}()

	err = tx.Commit(arguments.Ctx)
	if err != nil {
		err = fmt.Errorf("failed to commit DB transaction; %v", err)
		return nil, 0, 0, 0, 0, err
	}

	select {
	case <-arguments.Ctx.Done():
		err = fmt.Errorf("context canceled")
		return nil, 0, 0, 0, 0, err
	case err = <-errs:
		if err != nil {
			return nil, 0, 0, 0, 0, err
		}
	}

	count := int64(len(objects))
	totalCount := count
	page := int64(1)
	totalPages := page

	return objects, count, totalCount, page, totalPages, nil
}

func handlePutPhysicalThing(arguments *server.LoadArguments, db *pgxpool.Pool, waitForChange server.WaitForChange, object *PhysicalThing) ([]*PhysicalThing, int64, int64, int64, int64, error) {
	tx, err := db.Begin(arguments.Ctx)
	if err != nil {
		err = fmt.Errorf("failed to begin DB transaction; %v", err)
		return nil, 0, 0, 0, 0, err
	}

	defer func() {
		_ = tx.Rollback(arguments.Ctx)
	}()

	xid, err := query.GetXid(arguments.Ctx, tx)
	if err != nil {
		err = fmt.Errorf("failed to get xid; %v", err)
		return nil, 0, 0, 0, 0, err
	}
	_ = xid

	err = object.Update(arguments.Ctx, tx, true)
	if err != nil {
		err = fmt.Errorf("failed to update %#+v; %v", object, err)
		return nil, 0, 0, 0, 0, err
	}

	errs := make(chan error, 1)
	go func() {
		_, err := waitForChange(arguments.Ctx, []stream.Action{stream.UPDATE, stream.SOFT_DELETE, stream.SOFT_RESTORE, stream.SOFT_UPDATE}, PhysicalThingTable, xid)
		if err != nil {
			err = fmt.Errorf("failed to wait for change; %v", err)
			errs <- err
			return
		}

		errs <- nil
	}()

	err = tx.Commit(arguments.Ctx)
	if err != nil {
		err = fmt.Errorf("failed to commit DB transaction; %v", err)
		return nil, 0, 0, 0, 0, err
	}

	select {
	case <-arguments.Ctx.Done():
		err = fmt.Errorf("context canceled")
		return nil, 0, 0, 0, 0, err
	case err = <-errs:
		if err != nil {
			return nil, 0, 0, 0, 0, err
		}
	}

	count := int64(1)
	totalCount := count
	page := int64(1)
	totalPages := page

	return []*PhysicalThing{object}, count, totalCount, page, totalPages, nil
}

func handlePatchPhysicalThing(arguments *server.LoadArguments, db *pgxpool.Pool, waitForChange server.WaitForChange, object *PhysicalThing, forceSetValuesForFields []string) ([]*PhysicalThing, int64, int64, int64, int64, error) {
	tx, err := db.Begin(arguments.Ctx)
	if err != nil {
		err = fmt.Errorf("failed to begin DB transaction; %v", err)
		return nil, 0, 0, 0, 0, err
	}

	defer func() {
		_ = tx.Rollback(arguments.Ctx)
	}()

	xid, err := query.GetXid(arguments.Ctx, tx)
	if err != nil {
		err = fmt.Errorf("failed to get xid; %v", err)
		return nil, 0, 0, 0, 0, err
	}
	_ = xid

	err = object.Update(arguments.Ctx, tx, false, forceSetValuesForFields...)
	if err != nil {
		err = fmt.Errorf("failed to update %#+v; %v", object, err)
		return nil, 0, 0, 0, 0, err
	}

	errs := make(chan error, 1)
	go func() {
		_, err := waitForChange(arguments.Ctx, []stream.Action{stream.UPDATE, stream.SOFT_DELETE, stream.SOFT_RESTORE, stream.SOFT_UPDATE}, PhysicalThingTable, xid)
		if err != nil {
			err = fmt.Errorf("failed to wait for change; %v", err)
			errs <- err
			return
		}

		errs <- nil
	}()

	err = tx.Commit(arguments.Ctx)
	if err != nil {
		err = fmt.Errorf("failed to commit DB transaction; %v", err)
		return nil, 0, 0, 0, 0, err
	}

	select {
	case <-arguments.Ctx.Done():
		err = fmt.Errorf("context canceled")
		return nil, 0, 0, 0, 0, err
	case err = <-errs:
		if err != nil {
			return nil, 0, 0, 0, 0, err
		}
	}

	count := int64(1)
	totalCount := count
	page := int64(1)
	totalPages := page

	return []*PhysicalThing{object}, count, totalCount, page, totalPages, nil
}

func handleDeletePhysicalThing(arguments *server.LoadArguments, db *pgxpool.Pool, waitForChange server.WaitForChange, object *PhysicalThing) error {
	tx, err := db.Begin(arguments.Ctx)
	if err != nil {
		err = fmt.Errorf("failed to begin DB transaction; %v", err)
		return err
	}

	defer func() {
		_ = tx.Rollback(arguments.Ctx)
	}()

	xid, err := query.GetXid(arguments.Ctx, tx)
	if err != nil {
		err = fmt.Errorf("failed to get xid; %v", err)
		return err
	}
	_ = xid

	err = object.Delete(arguments.Ctx, tx)
	if err != nil {
		err = fmt.Errorf("failed to delete %#+v; %v", object, err)
		return err
	}

	errs := make(chan error, 1)
	go func() {
		_, err := waitForChange(arguments.Ctx, []stream.Action{stream.DELETE, stream.SOFT_DELETE}, PhysicalThingTable, xid)
		if err != nil {
			err = fmt.Errorf("failed to wait for change; %v", err)
			errs <- err
			return
		}

		errs <- nil
	}()

	err = tx.Commit(arguments.Ctx)
	if err != nil {
		err = fmt.Errorf("failed to commit DB transaction; %v", err)
		return err
	}

	select {
	case <-arguments.Ctx.Done():
		err = fmt.Errorf("context canceled")
		return err
	case err = <-errs:
		if err != nil {
			return err
		}
	}

	return nil
}

func MutateRouterForPhysicalThing(r chi.Router, db *pgxpool.Pool, redisPool *redis.Pool, objectMiddlewares []server.ObjectMiddleware, waitForChange server.WaitForChange) {

	func() {
		getManyHandler, err := getHTTPHandler(
			http.MethodGet,
			"/physical-things",
			http.StatusOK,
			func(
				ctx context.Context,
				pathParams server.EmptyPathParams,
				queryParams map[string]any,
				req server.EmptyRequest,
				rawReq any,
			) (server.Response[PhysicalThing], error) {
				before := time.Now()

				redisConn := redisPool.Get()
				defer func() {
					_ = redisConn.Close()
				}()

				arguments, err := server.GetSelectManyArguments(ctx, queryParams, PhysicalThingIntrospectedTable, nil, nil)
				if err != nil {
					if config.Debug() {
						log.Printf("request cache not yet reached; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return server.Response[PhysicalThing]{}, err
				}

				cachedResponseAsJSON, cacheHit, err := server.GetCachedResponseAsJSON(arguments.RequestHash, redisConn)
				if err != nil {
					if config.Debug() {
						log.Printf("request cache failed; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return server.Response[PhysicalThing]{}, err
				}

				if cacheHit {
					var cachedResponse server.Response[PhysicalThing]

					/* TODO: it'd be nice to be able to avoid this (i.e. just pass straight through) */
					err = json.Unmarshal(cachedResponseAsJSON, &cachedResponse)
					if err != nil {
						if config.Debug() {
							log.Printf("request cache hit but failed unmarshal; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
						}

						return server.Response[PhysicalThing]{}, err
					}

					if config.Debug() {
						log.Printf("request cache hit; request succeeded in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return cachedResponse, nil
				}

				objects, count, totalCount, _, _, err := handleGetPhysicalThings(arguments, db)
				if err != nil {
					if config.Debug() {
						log.Printf("request cache missed; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return server.Response[PhysicalThing]{}, err
				}

				limit := int64(0)
				if arguments.Limit != nil {
					limit = int64(*arguments.Limit)
				}

				offset := int64(0)
				if arguments.Offset != nil {
					offset = int64(*arguments.Offset)
				}

				response := server.Response[PhysicalThing]{
					Status:     http.StatusOK,
					Success:    true,
					Error:      nil,
					Objects:    objects,
					Count:      count,
					TotalCount: totalCount,
					Limit:      limit,
					Offset:     offset,
				}

				/* TODO: it'd be nice to be able to avoid this (i.e. just marshal once, further out) */
				responseAsJSON, err := json.Marshal(response)
				if err != nil {
					if config.Debug() {
						log.Printf("request cache missed; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return server.Response[PhysicalThing]{}, err
				}

				err = server.StoreCachedResponse(arguments.RequestHash, redisConn, responseAsJSON)
				if err != nil {
					log.Printf("warning; %v", err)
				}

				if config.Debug() {
					log.Printf("request cache missed; request succeeded in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
				}

				return response, nil
			},
			PhysicalThing{},
			PhysicalThingIntrospectedTable,
		)
		if err != nil {
			panic(err)
		}
		r.Get(getManyHandler.FullPath, getManyHandler.ServeHTTP)
	}()

	func() {
		getOneHandler, err := getHTTPHandler(
			http.MethodGet,
			"/physical-things/{primaryKey}",
			http.StatusOK,
			func(
				ctx context.Context,
				pathParams PhysicalThingOnePathParams,
				queryParams PhysicalThingLoadQueryParams,
				req server.EmptyRequest,
				rawReq any,
			) (server.Response[PhysicalThing], error) {
				before := time.Now()

				redisConn := redisPool.Get()
				defer func() {
					_ = redisConn.Close()
				}()

				arguments, err := server.GetSelectOneArguments(ctx, queryParams.Depth, PhysicalThingIntrospectedTable, pathParams.PrimaryKey, nil, nil)
				if err != nil {
					if config.Debug() {
						log.Printf("request cache not yet reached; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return server.Response[PhysicalThing]{}, err
				}

				cachedResponseAsJSON, cacheHit, err := server.GetCachedResponseAsJSON(arguments.RequestHash, redisConn)
				if err != nil {
					if config.Debug() {
						log.Printf("request cache failed; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return server.Response[PhysicalThing]{}, err
				}

				if cacheHit {
					var cachedResponse server.Response[PhysicalThing]

					/* TODO: it'd be nice to be able to avoid this (i.e. just pass straight through) */
					err = json.Unmarshal(cachedResponseAsJSON, &cachedResponse)
					if err != nil {
						if config.Debug() {
							log.Printf("request cache hit but failed unmarshal; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
						}

						return server.Response[PhysicalThing]{}, err
					}

					if config.Debug() {
						log.Printf("request cache hit; request succeeded in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return cachedResponse, nil
				}

				objects, count, totalCount, _, _, err := handleGetPhysicalThing(arguments, db, pathParams.PrimaryKey)
				if err != nil {
					if config.Debug() {
						log.Printf("request cache missed; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return server.Response[PhysicalThing]{}, err
				}

				limit := int64(0)

				offset := int64(0)

				response := server.Response[PhysicalThing]{
					Status:     http.StatusOK,
					Success:    true,
					Error:      nil,
					Objects:    objects,
					Count:      count,
					TotalCount: totalCount,
					Limit:      limit,
					Offset:     offset,
				}

				/* TODO: it'd be nice to be able to avoid this (i.e. just marshal once, further out) */
				responseAsJSON, err := json.Marshal(response)
				if err != nil {
					if config.Debug() {
						log.Printf("request cache missed; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return server.Response[PhysicalThing]{}, err
				}

				err = server.StoreCachedResponse(arguments.RequestHash, redisConn, responseAsJSON)
				if err != nil {
					log.Printf("warning; %v", err)
				}

				if config.Debug() {
					log.Printf("request cache hit; request succeeded in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
				}

				return response, nil
			},
			PhysicalThing{},
			PhysicalThingIntrospectedTable,
		)
		if err != nil {
			panic(err)
		}
		r.Get(getOneHandler.FullPath, getOneHandler.ServeHTTP)
	}()

	func() {
		postHandler, err := getHTTPHandler(
			http.MethodPost,
			"/physical-things",
			http.StatusCreated,
			func(
				ctx context.Context,
				pathParams server.EmptyPathParams,
				queryParams PhysicalThingLoadQueryParams,
				req []*PhysicalThing,
				rawReq any,
			) (server.Response[PhysicalThing], error) {
				allRawItems, ok := rawReq.([]any)
				if !ok {
					return server.Response[PhysicalThing]{}, fmt.Errorf("failed to cast %#+v to []map[string]any", rawReq)
				}

				allItems := make([]map[string]any, 0)
				for _, rawItem := range allRawItems {
					item, ok := rawItem.(map[string]any)
					if !ok {
						return server.Response[PhysicalThing]{}, fmt.Errorf("failed to cast %#+v to map[string]any", rawItem)
					}

					allItems = append(allItems, item)
				}

				forceSetValuesForFieldsByObjectIndex := make([][]string, 0)
				for _, item := range allItems {
					forceSetValuesForFields := make([]string, 0)
					for _, possibleField := range slices.Collect(maps.Keys(item)) {
						if !slices.Contains(PhysicalThingTableColumns, possibleField) {
							continue
						}

						forceSetValuesForFields = append(forceSetValuesForFields, possibleField)
					}
					forceSetValuesForFieldsByObjectIndex = append(forceSetValuesForFieldsByObjectIndex, forceSetValuesForFields)
				}

				arguments, err := server.GetLoadArguments(ctx, queryParams.Depth)
				if err != nil {
					return server.Response[PhysicalThing]{}, err
				}

				objects, count, totalCount, _, _, err := handlePostPhysicalThing(arguments, db, waitForChange, req, forceSetValuesForFieldsByObjectIndex)
				if err != nil {
					return server.Response[PhysicalThing]{}, err
				}

				limit := int64(0)

				offset := int64(0)

				return server.Response[PhysicalThing]{
					Status:     http.StatusOK,
					Success:    true,
					Error:      nil,
					Objects:    objects,
					Count:      count,
					TotalCount: totalCount,
					Limit:      limit,
					Offset:     offset,
				}, nil
			},
			PhysicalThing{},
			PhysicalThingIntrospectedTable,
		)
		if err != nil {
			panic(err)
		}
		r.Post(postHandler.FullPath, postHandler.ServeHTTP)
	}()

	func() {
		putHandler, err := getHTTPHandler(
			http.MethodPatch,
			"/physical-things/{primaryKey}",
			http.StatusOK,
			func(
				ctx context.Context,
				pathParams PhysicalThingOnePathParams,
				queryParams PhysicalThingLoadQueryParams,
				req PhysicalThing,
				rawReq any,
			) (server.Response[PhysicalThing], error) {
				item, ok := rawReq.(map[string]any)
				if !ok {
					return server.Response[PhysicalThing]{}, fmt.Errorf("failed to cast %#+v to map[string]any", item)
				}

				arguments, err := server.GetLoadArguments(ctx, queryParams.Depth)
				if err != nil {
					return server.Response[PhysicalThing]{}, err
				}

				object := &req
				object.ID = pathParams.PrimaryKey

				objects, count, totalCount, _, _, err := handlePutPhysicalThing(arguments, db, waitForChange, object)
				if err != nil {
					return server.Response[PhysicalThing]{}, err
				}

				limit := int64(0)

				offset := int64(0)

				return server.Response[PhysicalThing]{
					Status:     http.StatusOK,
					Success:    true,
					Error:      nil,
					Objects:    objects,
					Count:      count,
					TotalCount: totalCount,
					Limit:      limit,
					Offset:     offset,
				}, nil
			},
			PhysicalThing{},
			PhysicalThingIntrospectedTable,
		)
		if err != nil {
			panic(err)
		}
		r.Put(putHandler.FullPath, putHandler.ServeHTTP)
	}()

	func() {
		patchHandler, err := getHTTPHandler(
			http.MethodPatch,
			"/physical-things/{primaryKey}",
			http.StatusOK,
			func(
				ctx context.Context,
				pathParams PhysicalThingOnePathParams,
				queryParams PhysicalThingLoadQueryParams,
				req PhysicalThing,
				rawReq any,
			) (server.Response[PhysicalThing], error) {
				item, ok := rawReq.(map[string]any)
				if !ok {
					return server.Response[PhysicalThing]{}, fmt.Errorf("failed to cast %#+v to map[string]any", item)
				}

				forceSetValuesForFields := make([]string, 0)
				for _, possibleField := range slices.Collect(maps.Keys(item)) {
					if !slices.Contains(PhysicalThingTableColumns, possibleField) {
						continue
					}

					forceSetValuesForFields = append(forceSetValuesForFields, possibleField)
				}

				arguments, err := server.GetLoadArguments(ctx, queryParams.Depth)
				if err != nil {
					return server.Response[PhysicalThing]{}, err
				}

				object := &req
				object.ID = pathParams.PrimaryKey

				objects, count, totalCount, _, _, err := handlePatchPhysicalThing(arguments, db, waitForChange, object, forceSetValuesForFields)
				if err != nil {
					return server.Response[PhysicalThing]{}, err
				}

				limit := int64(0)

				offset := int64(0)

				return server.Response[PhysicalThing]{
					Status:     http.StatusOK,
					Success:    true,
					Error:      nil,
					Objects:    objects,
					Count:      count,
					TotalCount: totalCount,
					Limit:      limit,
					Offset:     offset,
				}, nil
			},
			PhysicalThing{},
			PhysicalThingIntrospectedTable,
		)
		if err != nil {
			panic(err)
		}
		r.Patch(patchHandler.FullPath, patchHandler.ServeHTTP)
	}()

	func() {
		deleteHandler, err := getHTTPHandler(
			http.MethodDelete,
			"/physical-things/{primaryKey}",
			http.StatusNoContent,
			func(
				ctx context.Context,
				pathParams PhysicalThingOnePathParams,
				queryParams PhysicalThingLoadQueryParams,
				req server.EmptyRequest,
				rawReq any,
			) (server.EmptyResponse, error) {
				arguments, err := server.GetLoadArguments(ctx, queryParams.Depth)
				if err != nil {
					return server.EmptyResponse{}, err
				}

				object := &PhysicalThing{}
				object.ID = pathParams.PrimaryKey

				err = handleDeletePhysicalThing(arguments, db, waitForChange, object)
				if err != nil {
					return server.EmptyResponse{}, err
				}

				return server.EmptyResponse{}, nil
			},
			PhysicalThing{},
			PhysicalThingIntrospectedTable,
		)
		if err != nil {
			panic(err)
		}
		r.Delete(deleteHandler.FullPath, deleteHandler.ServeHTTP)
	}()
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
		MutateRouterForPhysicalThing,
	)
}
