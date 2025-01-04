package model_reference

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"math"
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
	"golang.org/x/exp/maps"
)

type LogicalThing struct {
	ID                                                  uuid.UUID          `json:"id"`
	CreatedAt                                           time.Time          `json:"created_at"`
	UpdatedAt                                           time.Time          `json:"updated_at"`
	DeletedAt                                           *time.Time         `json:"deleted_at"`
	ExternalID                                          *string            `json:"external_id"`
	Name                                                string             `json:"name"`
	Type                                                string             `json:"type"`
	Tags                                                []string           `json:"tags"`
	Metadata                                            map[string]*string `json:"metadata"`
	RawData                                             *any               `json:"raw_data"`
	ClaimedUntil                                        time.Time          `json:"claimed_until"`
	ParentPhysicalThingID                               *uuid.UUID         `json:"parent_physical_thing_id"`
	ParentPhysicalThingIDObject                         *PhysicalThing     `json:"parent_physical_thing_object"`
	ParentLogicalThingID                                *uuid.UUID         `json:"parent_logical_thing_id"`
	ParentLogicalThingIDObject                          *LogicalThing      `json:"parent_logical_thing_object"`
	ReferencedByLogicalThingParentLogicalThingIDObjects []*LogicalThing    `json:"referenced_by_logical_thing_parent_logical_thing_id_objects"`
}

var LogicalThingTable = "logical_things"

var LogicalThingTableWithSchema = fmt.Sprintf("%s.%s", schema, LogicalThingTable)

var LogicalThingTableNamespaceID int32 = 1337 // LogicalThingTableNamespaceID

var ( // ColumnVariables
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
	LogicalThingTableAgeColumn                   = "age"
	LogicalThingTableOptionalAgeColumn           = "optional_age"
	LogicalThingTableCountColumn                 = "count"
	LogicalThingTableOptionalCountColumn         = "optional_count"
	LogicalThingTableParentPhysicalThingIDColumn = "parent_physical_thing_id"
	LogicalThingTableParentLogicalThingIDColumn  = "parent_logical_thing_id"
)

var ( // ColumnVariablesWithTypeCasts
	LogicalThingTableIDColumnWithTypeCast                    = `"id" AS id`
	LogicalThingTableCreatedAtColumnWithTypeCast             = `"created_at" AS created_at`
	LogicalThingTableUpdatedAtColumnWithTypeCast             = `"updated_at" AS updated_at`
	LogicalThingTableDeletedAtColumnWithTypeCast             = `"deleted_at" AS deleted_at`
	LogicalThingTableExternalIDColumnWithTypeCast            = `"external_id" AS external_id`
	LogicalThingTableNameColumnWithTypeCast                  = `"name" AS name`
	LogicalThingTableTypeColumnWithTypeCast                  = `"type" AS type`
	LogicalThingTableTagsColumnWithTypeCast                  = `"tags" AS tags`
	LogicalThingTableMetadataColumnWithTypeCast              = `"metadata" AS metadata`
	LogicalThingTableRawDataColumnWithTypeCast               = `"raw_data" AS raw_data`
	LogicalThingTableAgeColumnWithTypeCast                   = `"age" AS age`
	LogicalThingTableOptionalAgeColumnWithTypeCast           = `"optional_age" AS optional_age`
	LogicalThingTableCountColumnWithTypeCast                 = `"count" AS count`
	LogicalThingTableOptionalCountColumnWithTypeCast         = `"optional_count" AS optional_count`
	LogicalThingTableParentPhysicalThingIDColumnWithTypeCast = `"parent_physical_thing_id" AS parent_physical_thing_id`
	LogicalThingTableParentLogicalThingIDColumnWithTypeCast  = `"parent_logical_thing_id" AS parent_logical_thing_id`
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
	LogicalThingTableAgeColumn,
	LogicalThingTableOptionalAgeColumn,
	LogicalThingTableCountColumn,
	LogicalThingTableOptionalCountColumn,
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
	LogicalThingTableAgeColumnWithTypeCast,
	LogicalThingTableOptionalAgeColumnWithTypeCast,
	LogicalThingTableCountColumnWithTypeCast,
	LogicalThingTableOptionalCountColumnWithTypeCast,
	LogicalThingTableParentPhysicalThingIDColumnWithTypeCast,
	LogicalThingTableParentLogicalThingIDColumnWithTypeCast,
}

var LogicalThingIntrospectedTable *introspect.Table

var LogicalThingTableColumnLookup map[string]*introspect.Column

var ( // PrimaryKeyColumn
	LogicalThingTablePrimaryKeyColumn = LogicalThingTableIDColumn
)

func init() {
	LogicalThingIntrospectedTable = tableByName[LogicalThingTable]

	/* only needed during templating */
	if LogicalThingIntrospectedTable == nil {
		LogicalThingIntrospectedTable = &introspect.Table{}
	}

	LogicalThingTableColumnLookup = LogicalThingIntrospectedTable.ColumnByName
}

type LogicalThingOnePathParams struct {
	PrimaryKey uuid.UUID `json:"primaryKey"`
}

type LogicalThingLoadQueryParams struct {
	Depth *int `json:"depth"`
}

// <claim-request>

type LogicalThingClaimRequest struct {
	Until          time.Time `json:"until"`
	TimeoutSeconds float64   `json:"timeout_seconds"`
} // </claim-request>

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
		return fmt.Errorf("%v: %#+v; error; %v", k, v, err)
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
				return wrapError(k, v, fmt.Errorf("failed to cast to time.Time"))
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
				return wrapError(k, v, fmt.Errorf("failed to cast to time.Time"))
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
				return wrapError(k, v, fmt.Errorf("failed to cast to string"))
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
				return wrapError(k, v, fmt.Errorf("failed to cast to string"))
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
				return wrapError(k, v, fmt.Errorf("failed to cast to string"))
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
				return wrapError(k, v, fmt.Errorf("failed to cast to []string"))
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
				return wrapError(k, v, fmt.Errorf("failed to cast to map[string]*string"))
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
				return wrapError(k, v, fmt.Errorf("failed to cast to any"))
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
				return wrapError(k, v, fmt.Errorf("failed to cast to uuid.UUID"))
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
				return wrapError(k, v, fmt.Errorf("failed to cast to uuid.UUID"))
			}

			m.ParentLogicalThingID = &temp2

		}
	}

	return nil
}

func (m *LogicalThing) Reload(ctx context.Context, tx pgx.Tx, includeDeleteds ...bool) error {
	extraWhere := ""

	// <reload-soft-delete>
	if len(includeDeleteds) > 0 && includeDeleteds[0] {
		if slices.Contains(LogicalThingTableColumns, "deleted_at") {
			extraWhere = "\n    AND (deleted_at IS null OR deleted_at IS NOT null)"
		}
	}
	// </reload-soft-delete>

	ctx, cleanup := query.WithQueryID(ctx)
	defer cleanup()

	ctx = query.WithMaxDepth(ctx, nil)

	o, _, _, _, _, err := SelectLogicalThing(
		ctx,
		tx,
		fmt.Sprintf("%v = $1%v", m.GetPrimaryKeyColumn(), extraWhere),
		m.GetPrimaryKeyValue(),
	)
	if err != nil {
		return err
	}

	// <reload-set-fields>
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
	m.ParentPhysicalThingID = o.ParentPhysicalThingID
	m.ParentPhysicalThingIDObject = o.ParentPhysicalThingIDObject
	m.ParentLogicalThingID = o.ParentLogicalThingID
	m.ParentLogicalThingIDObject = o.ParentLogicalThingIDObject
	// </reload-set-fields>

	return nil
}

func (m *LogicalThing) Insert(ctx context.Context, tx pgx.Tx, setPrimaryKey bool, setZeroValues bool, forceSetValuesForFields ...string) error {
	columns := make([]string, 0)
	values := make([]any, 0)

	// <insert-set-fields-primary-key>
	// <insert-set-field-primary-key>
	if setPrimaryKey && (setZeroValues || !types.IsZeroUUID(m.ID) || slices.Contains(forceSetValuesForFields, LogicalThingTableIDColumn) || isRequired(LogicalThingTableColumnLookup, LogicalThingTableIDColumn)) {
		columns = append(columns, LogicalThingTableIDColumn)

		v, err := types.FormatUUID(m.ID)
		if err != nil {
			return fmt.Errorf("failed to handle m.ID; %v", err)
		}

		values = append(values, v)
	}
	// </insert-set-field-primary-key>
	// </insert-set-fields-primary-key>

	// <insert-set-fields>
	// <insert-set-field>
	if setZeroValues || !types.IsZeroTime(m.CreatedAt) || slices.Contains(forceSetValuesForFields, LogicalThingTableCreatedAtColumn) || isRequired(LogicalThingTableColumnLookup, LogicalThingTableCreatedAtColumn) {
		columns = append(columns, LogicalThingTableCreatedAtColumn)

		v, err := types.FormatTime(m.CreatedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.CreatedAt; %v", err)
		}

		values = append(values, v)
	}
	// </insert-set-field>

	if setZeroValues || !types.IsZeroTime(m.UpdatedAt) || slices.Contains(forceSetValuesForFields, LogicalThingTableUpdatedAtColumn) {
		columns = append(columns, LogicalThingTableUpdatedAtColumn)
		values = append(values, m.UpdatedAt)
	}

	if setZeroValues || !types.IsZeroTime(m.DeletedAt) || slices.Contains(forceSetValuesForFields, LogicalThingTableDeletedAtColumn) {
		columns = append(columns, LogicalThingTableDeletedAtColumn)
		values = append(values, m.DeletedAt)
	}

	if setZeroValues || !types.IsZeroString(m.ExternalID) || slices.Contains(forceSetValuesForFields, LogicalThingTableExternalIDColumn) {
		columns = append(columns, LogicalThingTableExternalIDColumn)
		values = append(values, m.ExternalID)
	}

	if setZeroValues || !types.IsZeroString(m.Name) || slices.Contains(forceSetValuesForFields, LogicalThingTableNameColumn) {
		columns = append(columns, LogicalThingTableNameColumn)
		values = append(values, m.Name)
	}

	if setZeroValues || !types.IsZeroString(m.Type) || slices.Contains(forceSetValuesForFields, LogicalThingTableTypeColumn) {
		columns = append(columns, LogicalThingTableTypeColumn)
		values = append(values, m.Type)
	}

	if setZeroValues || !types.IsZeroStringArray(m.Tags) || slices.Contains(forceSetValuesForFields, LogicalThingTableTagsColumn) {
		columns = append(columns, LogicalThingTableTagsColumn)
		values = append(values, m.Tags)
	}

	if setZeroValues || !types.IsZeroHstore(m.Metadata) || slices.Contains(forceSetValuesForFields, LogicalThingTableMetadataColumn) {
		columns = append(columns, LogicalThingTableMetadataColumn)
		values = append(values, m.Metadata)
	}

	if setZeroValues || !types.IsZeroJSON(m.RawData) || slices.Contains(forceSetValuesForFields, LogicalThingTableRawDataColumn) {
		columns = append(columns, LogicalThingTableRawDataColumn)
		values = append(values, m.RawData)
	}

	if setZeroValues || !types.IsZeroUUID(m.ParentPhysicalThingID) || slices.Contains(forceSetValuesForFields, LogicalThingTableParentPhysicalThingIDColumn) {
		columns = append(columns, LogicalThingTableParentPhysicalThingIDColumn)
		values = append(values, m.ParentPhysicalThingID)
	}

	if setZeroValues || !types.IsZeroUUID(m.ParentLogicalThingID) || slices.Contains(forceSetValuesForFields, LogicalThingTableParentLogicalThingIDColumn) {
		columns = append(columns, LogicalThingTableParentLogicalThingIDColumn)
		values = append(values, m.ParentLogicalThingID)
	}
	// </insert-set-fields>

	ctx, cleanup := query.WithQueryID(ctx)
	defer cleanup()

	ctx = query.WithMaxDepth(ctx, nil)

	item, err := query.Insert(
		ctx,
		tx,
		LogicalThingTableWithSchema,
		columns,
		nil,
		false,
		false,
		LogicalThingTableColumns,
		values...,
	)
	if err != nil {
		return fmt.Errorf("failed to insert %#+v; %v", m, err)
	}

	// <insert-set-primary-key>
	v := (*item)[LogicalThingTablePrimaryKeyColumn]

	if v == nil {
		return fmt.Errorf("failed to find %v in %#+v", LogicalThingTablePrimaryKeyColumn, item)
	}

	wrapError := func(err error) error {
		return fmt.Errorf(
			"failed to treat %v: %#+v as uuid.UUID: %v",
			LogicalThingTablePrimaryKeyColumn,
			(*item)[LogicalThingTablePrimaryKeyColumn],
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
	// </insert-set-primary-key>

	err = m.Reload(ctx, tx, slices.Contains(forceSetValuesForFields, "deleted_at"))
	if err != nil {
		return fmt.Errorf("failed to reload after insert; %v", err)
	}

	return nil
}

func (m *LogicalThing) Update(ctx context.Context, tx pgx.Tx, setZeroValues bool, forceSetValuesForFields ...string) error {
	columns := make([]string, 0)
	values := make([]any, 0)

	// <update-set-fields>
	// <update-set-field>
	if setZeroValues || !types.IsZeroTime(m.CreatedAt) || slices.Contains(forceSetValuesForFields, LogicalThingTableCreatedAtColumn) {
		columns = append(columns, LogicalThingTableCreatedAtColumn)

		v, err := types.FormatTime(m.CreatedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.CreatedAt; %v", err)
		}

		values = append(values, v)
	}
	// </update-set-field>

	if setZeroValues || !types.IsZeroTime(m.UpdatedAt) || slices.Contains(forceSetValuesForFields, LogicalThingTableUpdatedAtColumn) {
		columns = append(columns, LogicalThingTableUpdatedAtColumn)
		values = append(values, m.UpdatedAt)
	}

	if setZeroValues || !types.IsZeroTime(m.DeletedAt) || slices.Contains(forceSetValuesForFields, LogicalThingTableDeletedAtColumn) {
		columns = append(columns, LogicalThingTableDeletedAtColumn)
		values = append(values, m.DeletedAt)
	}

	if setZeroValues || !types.IsZeroString(m.ExternalID) || slices.Contains(forceSetValuesForFields, LogicalThingTableExternalIDColumn) {
		columns = append(columns, LogicalThingTableExternalIDColumn)
		values = append(values, m.ExternalID)
	}

	if setZeroValues || !types.IsZeroString(m.Name) || slices.Contains(forceSetValuesForFields, LogicalThingTableNameColumn) {
		columns = append(columns, LogicalThingTableNameColumn)
		values = append(values, m.Name)
	}

	if setZeroValues || !types.IsZeroString(m.Type) || slices.Contains(forceSetValuesForFields, LogicalThingTableTypeColumn) {
		columns = append(columns, LogicalThingTableTypeColumn)
		values = append(values, m.Type)
	}

	if setZeroValues || !types.IsZeroStringArray(m.Tags) || slices.Contains(forceSetValuesForFields, LogicalThingTableTagsColumn) {
		columns = append(columns, LogicalThingTableTagsColumn)
		values = append(values, m.Tags)
	}

	if setZeroValues || !types.IsZeroHstore(m.Metadata) || slices.Contains(forceSetValuesForFields, LogicalThingTableMetadataColumn) {
		columns = append(columns, LogicalThingTableMetadataColumn)
		values = append(values, m.Metadata)
	}

	if setZeroValues || !types.IsZeroJSON(m.RawData) || slices.Contains(forceSetValuesForFields, LogicalThingTableRawDataColumn) {
		columns = append(columns, LogicalThingTableRawDataColumn)
		values = append(values, m.RawData)
	}

	if setZeroValues || !types.IsZeroUUID(m.ParentPhysicalThingID) || slices.Contains(forceSetValuesForFields, LogicalThingTableParentPhysicalThingIDColumn) {
		columns = append(columns, LogicalThingTableParentPhysicalThingIDColumn)
		values = append(values, m.ParentPhysicalThingID)
	}

	if setZeroValues || !types.IsZeroUUID(m.ParentLogicalThingID) || slices.Contains(forceSetValuesForFields, LogicalThingTableParentLogicalThingIDColumn) {
		columns = append(columns, LogicalThingTableParentLogicalThingIDColumn)
		values = append(values, m.ParentLogicalThingID)
	}
	// </update-set-fields>

	// <update-set-primary-key>

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
		LogicalThingTableWithSchema,
		columns,
		fmt.Sprintf("%v = $$??", LogicalThingTablePrimaryKeyColumn),
		LogicalThingTableColumns,
		values...,
	)
	if err != nil {
		return fmt.Errorf("failed to update %#+v; %v", m, err)
	}
	// </update-set-primary-key>

	err = m.Reload(ctx, tx, slices.Contains(forceSetValuesForFields, "deleted_at"))
	if err != nil {
		return fmt.Errorf("failed to reload after update")
	}

	return nil
}

func (m *LogicalThing) Delete(ctx context.Context, tx pgx.Tx, hardDeletes ...bool) error {
	// <delete-soft-delete>
	hardDelete := false
	if len(hardDeletes) > 0 {
		hardDelete = hardDeletes[0]
	}

	if !hardDelete && slices.Contains(LogicalThingTableColumns, "deleted_at") {
		m.DeletedAt = helpers.Ptr(time.Now().UTC())
		err := m.Update(ctx, tx, false, "deleted_at")
		if err != nil {
			return fmt.Errorf("failed to soft-delete (update) %#+v; %v", m, err)
		}
	}
	// </delete-soft-delete>

	values := make([]any, 0)

	// <delete-set-primary-key>
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
		LogicalThingTableWithSchema,
		fmt.Sprintf("%v = $$??", LogicalThingTablePrimaryKeyColumn),
		values...,
	)
	if err != nil {
		return fmt.Errorf("failed to delete %#+v; %v", m, err)
	}
	// </delete-set-primary-key>

	_ = m.Reload(ctx, tx, true)

	return nil
}

func (m *LogicalThing) LockTable(ctx context.Context, tx pgx.Tx, timeouts ...time.Duration) error {
	return query.LockTable(ctx, tx, LogicalThingTableWithSchema, timeouts...)
}

func (m *LogicalThing) LockTableWithRetries(ctx context.Context, tx pgx.Tx, overallTimeout time.Duration, individualAttempttimeout time.Duration) error {
	return query.LockTableWithRetries(ctx, tx, LogicalThingTableWithSchema, overallTimeout, individualAttempttimeout)
}

func (m *LogicalThing) AdvisoryLock(ctx context.Context, tx pgx.Tx, key int32, timeouts ...time.Duration) error {
	return query.AdvisoryLock(ctx, tx, LogicalThingTableNamespaceID, key, timeouts...)
}

func (m *LogicalThing) AdvisoryLockWithRetries(ctx context.Context, tx pgx.Tx, key int32, overallTimeout time.Duration, individualAttempttimeout time.Duration) error {
	return query.AdvisoryLockWithRetries(ctx, tx, LogicalThingTableNamespaceID, key, overallTimeout, individualAttempttimeout)
}

// <claim-method>

func (m *LogicalThing) Claim(ctx context.Context, tx pgx.Tx, until time.Time, timeout time.Duration) error {
	err := m.AdvisoryLockWithRetries(ctx, tx, math.MinInt32, timeout, time.Second*1)
	if err != nil {
		return fmt.Errorf("failed to claim (advisory lock): %s", err.Error())
	}

	_, _, _, _, _, err = SelectLogicalThing(
		ctx,
		tx,
		fmt.Sprintf(
			"%s = $$?? AND (claimed_until IS null OR claimed_until < now())",
			LogicalThingTablePrimaryKeyColumn,
		),
		m.GetPrimaryKeyValue(),
	)
	if err != nil {
		return fmt.Errorf("failed to claim (select): %s", err.Error())
	}

	m.ClaimedUntil = until

	err = m.Update(ctx, tx, false)
	if err != nil {
		return fmt.Errorf("failed to claim (update): %s", err.Error())
	}

	return nil
} // </claim-method>

func SelectLogicalThings(ctx context.Context, tx pgx.Tx, where string, orderBy *string, limit *int, offset *int, values ...any) ([]*LogicalThing, int64, int64, int64, int64, error) {
	before := time.Now()

	if config.Debug() {
		log.Printf("entered SelectLogicalThings")

		defer func() {
			log.Printf("exited SelectLogicalThings in %s", time.Since(before))
		}()
	}

	// TODO: change this lazy implicit logic to be informed by introspection at generation time
	if slices.Contains(LogicalThingTableColumns, "deleted_at") {
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

	shouldLoad := query.ShouldLoad(ctx, LogicalThingTable) || query.ShouldLoad(ctx, fmt.Sprintf("referenced_by_%s", LogicalThingTable))

	var ok bool
	ctx, ok = query.HandleQueryPathGraphCycles(ctx, fmt.Sprintf("%s{%v}", LogicalThingTable, nil), !isLoadQuery)
	if !ok && !shouldLoad {
		if config.Debug() {
			log.Printf("skipping SelectLogicalThings early (query.ShouldLoad(): %v, query.HandleQueryPathGraphCycles(): %v)", shouldLoad, ok)
		}
		return []*LogicalThing{}, 0, 0, 0, 0, nil
	}

	items, count, totalCount, page, totalPages, err := query.Select(
		ctx,
		tx,
		LogicalThingTableColumnsWithTypeCasts,
		LogicalThingTableWithSchema,
		where,
		orderBy,
		limit,
		offset,
		values...,
	)
	if err != nil {
		return nil, 0, 0, 0, 0, fmt.Errorf("failed to call SelectLogicalThings; %v", err)
	}

	objects := make([]*LogicalThing, 0)

	for _, item := range *items {
		object := &LogicalThing{}

		err = object.FromItem(item)
		if err != nil {
			return nil, 0, 0, 0, 0, err
		}

		// <select-load-foreign-objects>
		// <select-load-foreign-object>
		if !types.IsZeroUUID(object.ParentPhysicalThingID) {
			ctx, ok := query.HandleQueryPathGraphCycles(ctx, fmt.Sprintf("%s{%v}", PhysicalThingTable, object.ParentPhysicalThingID), true)
			shouldLoad := query.ShouldLoad(ctx, PhysicalThingTable)
			if ok || shouldLoad {
				thisBefore := time.Now()

				if config.Debug() {
					log.Printf("loading __this_function__->SelectPhysicalThing for object.ParentPhysicalThingIDObject{%s: %v}", PhysicalThingTablePrimaryKeyColumn, object.ParentPhysicalThingID)
				}

				object.ParentPhysicalThingIDObject, _, _, _, _, err = SelectPhysicalThing(
					ctx,
					tx,
					fmt.Sprintf("%v = $1", PhysicalThingTablePrimaryKeyColumn),
					object.ParentPhysicalThingID,
				)
				if err != nil {
					if !errors.Is(err, sql.ErrNoRows) {
						return nil, 0, 0, 0, 0, err
					}
				}

				if config.Debug() {
					log.Printf("loaded __this_function__->SelectPhysicalThing for object.ParentPhysicalThingIDObject in %s", time.Since(thisBefore))
				}
			}
		}
		// </select-load-foreign-object>

		if !types.IsZeroUUID(object.ParentLogicalThingID) {
			ctx, ok := query.HandleQueryPathGraphCycles(ctx, fmt.Sprintf("%s{%v}", LogicalThingTable, object.ParentLogicalThingID), true)
			if ok {
				object.ParentLogicalThingIDObject, _, _, _, _, err = SelectLogicalThing(
					ctx,
					tx,
					fmt.Sprintf("%v = $1", LogicalThingTablePrimaryKeyColumn),
					object.ParentLogicalThingID,
				)
				if err != nil {
					if !errors.Is(err, sql.ErrNoRows) {
						return nil, 0, 0, 0, 0, err
					}
				}
			}
		}
		// </select-load-foreign-objects>

		// <select-load-referenced-by-objects>
		// <select-load-referenced-by-object>
		err = func() error {
			shouldLoad := query.ShouldLoad(ctx, fmt.Sprintf("referenced_by_%s", LogicalThingTable))
			ctx, ok := query.HandleQueryPathGraphCycles(ctx, fmt.Sprintf("__ReferencedBy__%s{%v}", LogicalThingTable, object.GetPrimaryKeyValue()), true)
			if ok || shouldLoad {
				thisBefore := time.Now()

				if config.Debug() {
					log.Printf("loading __this_function__->SelectLogicalThings for object.ReferencedByLogicalThingParentLogicalThingIDObjects")
				}

				object.ReferencedByLogicalThingParentLogicalThingIDObjects, _, _, _, _, err = SelectLogicalThings(
					ctx,
					tx,
					fmt.Sprintf("%v = $1", LogicalThingTablePrimaryKeyColumn), // referenced-by
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
					log.Printf("loaded __this_function__->SelectLogicalThings for object.ReferencedByLogicalThingParentLogicalThingIDObjects in %s", time.Since(thisBefore))
				}

			}

			return nil
		}()
		if err != nil {
			return nil, 0, 0, 0, 0, err
		}

		// </select-load-referenced-by-object>
		// </select-load-referenced-by-objects>

		objects = append(objects, object)
	}

	return objects, count, totalCount, page, totalPages, nil
}

func SelectLogicalThing(ctx context.Context, tx pgx.Tx, where string, values ...any) (*LogicalThing, int64, int64, int64, int64, error) {
	ctx, cleanup := query.WithQueryID(ctx)
	defer cleanup()

	ctx = query.WithMaxDepth(ctx, nil)

	objects, _, _, _, _, err := SelectLogicalThings(
		ctx,
		tx,
		where,
		nil,
		helpers.Ptr(2),
		helpers.Ptr(0),
		values...,
	)
	if err != nil {
		return nil, 0, 0, 0, 0, fmt.Errorf("failed to call SelectLogicalThing; %v", err)
	}

	if len(objects) > 1 {
		return nil, 0, 0, 0, 0, fmt.Errorf("attempt to call SelectLogicalThing returned more than 1 row")
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

// <claim-func>

func ClaimLogicalThing(ctx context.Context, tx pgx.Tx, until time.Time, timeout time.Duration, wheres ...string) (*LogicalThing, error) {
	m := &LogicalThing{}

	err := m.AdvisoryLockWithRetries(ctx, tx, math.MinInt32, timeout, time.Second*1)
	if err != nil {
		return nil, fmt.Errorf("failed to claim: %s", err.Error())
	}

	extraWhere := ""
	if len(wheres) > 0 {
		extraWhere = fmt.Sprintf("AND %s", extraWhere)
	}

	ms, _, _, _, _, err := SelectLogicalThings(
		ctx,
		tx,
		fmt.Sprintf(
			"(claimed_until IS null OR claimed_until < now())%s",
			extraWhere,
		),
		helpers.Ptr(
			"claimed_until ASC",
		),
		helpers.Ptr(1),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to claim: %s", err.Error())
	}

	if len(ms) == 0 {
		return nil, nil
	}

	m = ms[0]

	m.ClaimedUntil = until

	err = m.Update(ctx, tx, false)
	if err != nil {
		return nil, fmt.Errorf("failed to claim: %s", err.Error())
	}

	return m, nil
} // </claim-func>

func handleGetLogicalThings(arguments *server.SelectManyArguments, db *pgxpool.Pool) ([]*LogicalThing, int64, int64, int64, int64, error) {
	tx, err := db.Begin(arguments.Ctx)
	if err != nil {
		return nil, 0, 0, 0, 0, err
	}

	defer func() {
		_ = tx.Rollback(arguments.Ctx)
	}()

	objects, count, totalCount, page, totalPages, err := SelectLogicalThings(arguments.Ctx, tx, arguments.Where, arguments.OrderBy, arguments.Limit, arguments.Offset, arguments.Values...)
	if err != nil {
		return nil, 0, 0, 0, 0, err
	}

	err = tx.Commit(arguments.Ctx)
	if err != nil {
		return nil, 0, 0, 0, 0, err
	}

	return objects, count, totalCount, page, totalPages, nil
}

func handleGetLogicalThing(arguments *server.SelectOneArguments, db *pgxpool.Pool, primaryKey uuid.UUID) ([]*LogicalThing, int64, int64, int64, int64, error) {
	tx, err := db.Begin(arguments.Ctx)
	if err != nil {
		return nil, 0, 0, 0, 0, err
	}

	defer func() {
		_ = tx.Rollback(arguments.Ctx)
	}()

	object, count, totalCount, page, totalPages, err := SelectLogicalThing(arguments.Ctx, tx, arguments.Where, arguments.Values...)
	if err != nil {
		return nil, 0, 0, 0, 0, err
	}

	err = tx.Commit(arguments.Ctx)
	if err != nil {
		return nil, 0, 0, 0, 0, err
	}

	return []*LogicalThing{object}, count, totalCount, page, totalPages, nil
}

func handlePostLogicalThings(arguments *server.LoadArguments, db *pgxpool.Pool, waitForChange server.WaitForChange, objects []*LogicalThing, forceSetValuesForFieldsByObjectIndex [][]string) ([]*LogicalThing, int64, int64, int64, int64, error) {
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

	for i, object := range objects {
		err = object.Insert(arguments.Ctx, tx, false, false, forceSetValuesForFieldsByObjectIndex[i]...)
		if err != nil {
			err = fmt.Errorf("failed to insert %#+v; %v", object, err)
			return nil, 0, 0, 0, 0, err
		}

		objects[i] = object
	}

	errs := make(chan error, 1)
	go func() {
		_, err := waitForChange(arguments.Ctx, []stream.Action{stream.INSERT}, LogicalThingTable, xid)
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

func handlePutLogicalThing(arguments *server.LoadArguments, db *pgxpool.Pool, waitForChange server.WaitForChange, object *LogicalThing) ([]*LogicalThing, int64, int64, int64, int64, error) {
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
		_, err := waitForChange(arguments.Ctx, []stream.Action{stream.UPDATE, stream.SOFT_DELETE, stream.SOFT_RESTORE, stream.SOFT_UPDATE}, LogicalThingTable, xid)
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

	return []*LogicalThing{object}, count, totalCount, page, totalPages, nil
}

func handlePatchLogicalThing(arguments *server.LoadArguments, db *pgxpool.Pool, waitForChange server.WaitForChange, object *LogicalThing, forceSetValuesForFields []string) ([]*LogicalThing, int64, int64, int64, int64, error) {
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
		_, err := waitForChange(arguments.Ctx, []stream.Action{stream.UPDATE, stream.SOFT_DELETE, stream.SOFT_RESTORE, stream.SOFT_UPDATE}, LogicalThingTable, xid)
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

	return []*LogicalThing{object}, count, totalCount, page, totalPages, nil
}

func handleDeleteLogicalThing(arguments *server.LoadArguments, db *pgxpool.Pool, waitForChange server.WaitForChange, object *LogicalThing) error {
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
		_, err := waitForChange(arguments.Ctx, []stream.Action{stream.DELETE, stream.SOFT_DELETE}, LogicalThingTable, xid)
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

func MutateRouterForLogicalThing(r chi.Router, db *pgxpool.Pool, redisPool *redis.Pool, objectMiddlewares []server.ObjectMiddleware, waitForChange server.WaitForChange) {
	// <claim-handlers>

	func() {
		postHandlerForClaim, err := getHTTPHandler(
			http.MethodPost,
			"/claim-logical-thing",
			http.StatusOK,
			func(
				ctx context.Context,
				pathParams server.EmptyPathParams,
				queryParams server.EmptyQueryParams,
				req LogicalThingClaimRequest,
				rawReq any,
			) (server.Response[LogicalThing], error) {
				tx, err := db.Begin(ctx)
				if err != nil {
					return server.Response[LogicalThing]{}, err
				}

				defer func() {
					_ = tx.Rollback(ctx)
				}()

				object, err := ClaimLogicalThing(ctx, tx, req.Until, time.Millisecond*time.Duration(req.TimeoutSeconds*1000))
				if err != nil {
					return server.Response[LogicalThing]{}, err
				}

				count := int64(0)

				totalCount := int64(0)

				limit := int64(0)

				offset := int64(0)

				if object == nil {
					return server.Response[LogicalThing]{
						Status:     http.StatusOK,
						Success:    true,
						Error:      nil,
						Objects:    []*LogicalThing{},
						Count:      count,
						TotalCount: totalCount,
						Limit:      limit,
						Offset:     offset,
					}, nil
				}

				err = tx.Commit(ctx)
				if err != nil {
					return server.Response[LogicalThing]{}, err
				}

				return server.Response[LogicalThing]{
					Status:     http.StatusOK,
					Success:    true,
					Error:      nil,
					Objects:    []*LogicalThing{object},
					Count:      count,
					TotalCount: totalCount,
					Limit:      limit,
					Offset:     offset,
				}, nil
			},
			LogicalThing{},
			LogicalThingIntrospectedTable,
		)
		if err != nil {
			panic(err) // TODO
		}
		r.Post(postHandlerForClaim.FullPath, postHandlerForClaim.ServeHTTP)

		postHandlerForClaimOne, err := getHTTPHandler(
			http.MethodPost,
			"/logical-things/{primaryKey}/claim",
			http.StatusOK,
			func(
				ctx context.Context,
				pathParams LogicalThingOnePathParams,
				queryParams LogicalThingLoadQueryParams,
				req LogicalThingClaimRequest,
				rawReq any,
			) (server.Response[LogicalThing], error) {
				before := time.Now()

				redisConn := redisPool.Get()
				defer func() {
					_ = redisConn.Close()
				}()

				arguments, err := server.GetSelectOneArguments(ctx, queryParams.Depth, LogicalThingIntrospectedTable, pathParams.PrimaryKey, nil, nil)
				if err != nil {
					if config.Debug() {
						log.Printf("request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return server.Response[LogicalThing]{}, err
				}

				/* note: deliberately no attempt at a cache hit */

				var object *LogicalThing
				var count int64
				var totalCount int64

				err = func() error {
					tx, err := db.Begin(arguments.Ctx)
					if err != nil {
						return err
					}

					defer func() {
						_ = tx.Rollback(arguments.Ctx)
					}()

					object, count, totalCount, _, _, err = SelectLogicalThing(arguments.Ctx, tx, arguments.Where, arguments.Values...)
					if err != nil {
						return fmt.Errorf("failed to select object to claim: %s", err.Error())
					}

					err = object.Claim(arguments.Ctx, tx, req.Until, time.Millisecond*time.Duration(req.TimeoutSeconds*1000))
					if err != nil {
						return err
					}

					err = tx.Commit(arguments.Ctx)
					if err != nil {
						return err
					}

					return nil
				}()
				if err != nil {
					if config.Debug() {
						log.Printf("request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return server.Response[LogicalThing]{}, err
				}

				limit := int64(0)

				offset := int64(0)

				response := server.Response[LogicalThing]{
					Status:     http.StatusOK,
					Success:    true,
					Error:      nil,
					Objects:    []*LogicalThing{object},
					Count:      count,
					TotalCount: totalCount,
					Limit:      limit,
					Offset:     offset,
				}

				return response, nil
			},
			LogicalThing{},
			LogicalThingIntrospectedTable,
		)
		if err != nil {
			panic(err) // TODO
		}
		r.Post(postHandlerForClaimOne.FullPath, postHandlerForClaimOne.ServeHTTP)
	}()
	// </claim-handlers>

	func() {
		getManyHandler, err := getHTTPHandler(
			http.MethodGet,
			"/logical-things",
			http.StatusOK,
			func(
				ctx context.Context,
				pathParams server.EmptyPathParams,
				queryParams map[string]any,
				req server.EmptyRequest,
				rawReq any,
			) (server.Response[LogicalThing], error) {
				before := time.Now()

				redisConn := redisPool.Get()
				defer func() {
					_ = redisConn.Close()
				}()

				arguments, err := server.GetSelectManyArguments(ctx, queryParams, LogicalThingIntrospectedTable, nil, nil)
				if err != nil {
					if config.Debug() {
						log.Printf("request cache not yet reached; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return server.Response[LogicalThing]{}, err
				}

				cachedResponseAsJSON, cacheHit, err := server.GetCachedResponseAsJSON(arguments.RequestHash, redisConn)
				if err != nil {
					if config.Debug() {
						log.Printf("request cache failed; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return server.Response[LogicalThing]{}, err
				}

				if cacheHit {
					var cachedResponse server.Response[LogicalThing]

					/* TODO: it'd be nice to be able to avoid this (i.e. just pass straight through) */
					err = json.Unmarshal(cachedResponseAsJSON, &cachedResponse)
					if err != nil {
						if config.Debug() {
							log.Printf("request cache hit but failed unmarshal; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
						}

						return server.Response[LogicalThing]{}, err
					}

					if config.Debug() {
						log.Printf("request cache hit; request succeeded in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return cachedResponse, nil
				}

				objects, count, totalCount, _, _, err := handleGetLogicalThings(arguments, db)
				if err != nil {
					if config.Debug() {
						log.Printf("request cache missed; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return server.Response[LogicalThing]{}, err
				}

				limit := int64(0)
				if arguments.Limit != nil {
					limit = int64(*arguments.Limit)
				}

				offset := int64(0)
				if arguments.Offset != nil {
					offset = int64(*arguments.Offset)
				}

				response := server.Response[LogicalThing]{
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

					return server.Response[LogicalThing]{}, err
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
			LogicalThing{},
			LogicalThingIntrospectedTable,
		)
		if err != nil {
			panic(err) // TODO
		}
		r.Get(getManyHandler.FullPath, getManyHandler.ServeHTTP)
	}()

	func() {
		getOneHandler, err := getHTTPHandler(
			http.MethodGet,
			"/logical-things/{primaryKey}",
			http.StatusOK,
			func(
				ctx context.Context,
				pathParams LogicalThingOnePathParams,
				queryParams LogicalThingLoadQueryParams,
				req server.EmptyRequest,
				rawReq any,
			) (server.Response[LogicalThing], error) {
				before := time.Now()

				redisConn := redisPool.Get()
				defer func() {
					_ = redisConn.Close()
				}()

				arguments, err := server.GetSelectOneArguments(ctx, queryParams.Depth, LogicalThingIntrospectedTable, pathParams.PrimaryKey, nil, nil)
				if err != nil {
					if config.Debug() {
						log.Printf("request cache not yet reached; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return server.Response[LogicalThing]{}, err
				}

				cachedResponseAsJSON, cacheHit, err := server.GetCachedResponseAsJSON(arguments.RequestHash, redisConn)
				if err != nil {
					if config.Debug() {
						log.Printf("request cache failed; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return server.Response[LogicalThing]{}, err
				}

				if cacheHit {
					var cachedResponse server.Response[LogicalThing]

					/* TODO: it'd be nice to be able to avoid this (i.e. just pass straight through) */
					err = json.Unmarshal(cachedResponseAsJSON, &cachedResponse)
					if err != nil {
						if config.Debug() {
							log.Printf("request cache hit but failed unmarshal; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
						}

						return server.Response[LogicalThing]{}, err
					}

					if config.Debug() {
						log.Printf("request cache hit; request succeeded in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return cachedResponse, nil
				}

				objects, count, totalCount, _, _, err := handleGetLogicalThing(arguments, db, pathParams.PrimaryKey)
				if err != nil {
					if config.Debug() {
						log.Printf("request cache missed; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return server.Response[LogicalThing]{}, err
				}

				limit := int64(0)

				offset := int64(0)

				response := server.Response[LogicalThing]{
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

					return server.Response[LogicalThing]{}, err
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
			LogicalThing{},
			LogicalThingIntrospectedTable,
		)
		if err != nil {
			panic(err) // TODO
		}
		r.Get(getOneHandler.FullPath, getOneHandler.ServeHTTP)
	}()

	func() {
		postHandler, err := getHTTPHandler(
			http.MethodPost,
			"/logical-things",
			http.StatusCreated,
			func(
				ctx context.Context,
				pathParams server.EmptyPathParams,
				queryParams LogicalThingLoadQueryParams,
				req []*LogicalThing,
				rawReq any,
			) (server.Response[LogicalThing], error) {
				allRawItems, ok := rawReq.([]any)
				if !ok {
					return server.Response[LogicalThing]{}, fmt.Errorf("failed to cast %#+v to []map[string]any", rawReq)
				}

				allItems := make([]map[string]any, 0)
				for _, rawItem := range allRawItems {
					item, ok := rawItem.(map[string]any)
					if !ok {
						return server.Response[LogicalThing]{}, fmt.Errorf("failed to cast %#+v to map[string]any", rawItem)
					}

					allItems = append(allItems, item)
				}

				forceSetValuesForFieldsByObjectIndex := make([][]string, 0)
				for _, item := range allItems {
					forceSetValuesForFields := make([]string, 0)
					for _, possibleField := range maps.Keys(item) {
						if !slices.Contains(LogicalThingTableColumns, possibleField) {
							continue
						}

						forceSetValuesForFields = append(forceSetValuesForFields, possibleField)
					}
					forceSetValuesForFieldsByObjectIndex = append(forceSetValuesForFieldsByObjectIndex, forceSetValuesForFields)
				}

				arguments, err := server.GetLoadArguments(ctx, queryParams.Depth)
				if err != nil {
					return server.Response[LogicalThing]{}, err
				}

				objects, count, totalCount, _, _, err := handlePostLogicalThings(arguments, db, waitForChange, req, forceSetValuesForFieldsByObjectIndex)
				if err != nil {
					return server.Response[LogicalThing]{}, err
				}

				limit := int64(0)

				offset := int64(0)

				return server.Response[LogicalThing]{
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
			LogicalThing{},
			LogicalThingIntrospectedTable,
		)
		if err != nil {
			panic(err) // TODO
		}
		r.Post(postHandler.FullPath, postHandler.ServeHTTP)
	}()

	func() {
		putHandler, err := getHTTPHandler(
			http.MethodPatch,
			"/logical-things/{primaryKey}",
			http.StatusOK,
			func(
				ctx context.Context,
				pathParams LogicalThingOnePathParams,
				queryParams LogicalThingLoadQueryParams,
				req LogicalThing,
				rawReq any,
			) (server.Response[LogicalThing], error) {
				item, ok := rawReq.(map[string]any)
				if !ok {
					return server.Response[LogicalThing]{}, fmt.Errorf("failed to cast %#+v to map[string]any", item)
				}

				arguments, err := server.GetLoadArguments(ctx, queryParams.Depth)
				if err != nil {
					return server.Response[LogicalThing]{}, err
				}

				object := &req
				object.ID = pathParams.PrimaryKey

				objects, count, totalCount, _, _, err := handlePutLogicalThing(arguments, db, waitForChange, object)
				if err != nil {
					return server.Response[LogicalThing]{}, err
				}

				limit := int64(0)

				offset := int64(0)

				return server.Response[LogicalThing]{
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
			LogicalThing{},
			LogicalThingIntrospectedTable,
		)
		if err != nil {
			panic(err) // TODO
		}
		r.Put(putHandler.FullPath, putHandler.ServeHTTP)
	}()

	func() {
		patchHandler, err := getHTTPHandler(
			http.MethodPatch,
			"/logical-things/{primaryKey}",
			http.StatusOK,
			func(
				ctx context.Context,
				pathParams LogicalThingOnePathParams,
				queryParams LogicalThingLoadQueryParams,
				req LogicalThing,
				rawReq any,
			) (server.Response[LogicalThing], error) {
				item, ok := rawReq.(map[string]any)
				if !ok {
					return server.Response[LogicalThing]{}, fmt.Errorf("failed to cast %#+v to map[string]any", item)
				}

				forceSetValuesForFields := make([]string, 0)
				for _, possibleField := range maps.Keys(item) {
					if !slices.Contains(LogicalThingTableColumns, possibleField) {
						continue
					}

					forceSetValuesForFields = append(forceSetValuesForFields, possibleField)
				}

				arguments, err := server.GetLoadArguments(ctx, queryParams.Depth)
				if err != nil {
					return server.Response[LogicalThing]{}, err
				}

				object := &req
				object.ID = pathParams.PrimaryKey

				objects, count, totalCount, _, _, err := handlePatchLogicalThing(arguments, db, waitForChange, object, forceSetValuesForFields)
				if err != nil {
					return server.Response[LogicalThing]{}, err
				}

				limit := int64(0)

				offset := int64(0)

				return server.Response[LogicalThing]{
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
			LogicalThing{},
			LogicalThingIntrospectedTable,
		)
		if err != nil {
			panic(err) // TODO
		}
		r.Patch(patchHandler.FullPath, patchHandler.ServeHTTP)
	}()

	func() {
		deleteHandler, err := getHTTPHandler(
			http.MethodDelete,
			"/logical-things/{primaryKey}",
			http.StatusNoContent,
			func(
				ctx context.Context,
				pathParams LogicalThingOnePathParams,
				queryParams LogicalThingLoadQueryParams,
				req server.EmptyRequest,
				rawReq any,
			) (server.EmptyResponse, error) {
				arguments, err := server.GetLoadArguments(ctx, queryParams.Depth)
				if err != nil {
					return server.EmptyResponse{}, err
				}

				object := &LogicalThing{}
				object.ID = pathParams.PrimaryKey

				err = handleDeleteLogicalThing(arguments, db, waitForChange, object)
				if err != nil {
					return server.EmptyResponse{}, err
				}

				return server.EmptyResponse{}, nil
			},
			LogicalThing{},
			LogicalThingIntrospectedTable,
		)
		if err != nil {
			panic(err) // TODO
		}
		r.Delete(deleteHandler.FullPath, deleteHandler.ServeHTTP)
	}()
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
		LogicalThing{},
		NewLogicalThingFromItem,
		"/logical-things",
		MutateRouterForLogicalThing,
	)
}
