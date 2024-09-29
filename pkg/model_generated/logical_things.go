package model_generated

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
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
	RawData                                             any                `json:"raw_data"`
	Age                                                 time.Duration      `json:"age"`
	OptionalAge                                         *time.Duration     `json:"optional_age"`
	Count                                               int64              `json:"count"`
	OptionalCount                                       *int64             `json:"optional_count"`
	ParentPhysicalThingID                               *uuid.UUID         `json:"parent_physical_thing_id"`
	ParentPhysicalThingIDObject                         *PhysicalThing     `json:"parent_physical_thing_id_object"`
	ParentLogicalThingID                                *uuid.UUID         `json:"parent_logical_thing_id"`
	ParentLogicalThingIDObject                          *LogicalThing      `json:"parent_logical_thing_id_object"`
	ReferencedByLogicalThingParentLogicalThingIDObjects []*LogicalThing    `json:"referenced_by_logical_thing_parent_logical_thing_id_objects"`
}

var LogicalThingTable = "logical_things"

var LogicalThingTableNamespaceID int32 = 1337 + 4

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
	LogicalThingTableAgeColumn                   = "age"
	LogicalThingTableOptionalAgeColumn           = "optional_age"
	LogicalThingTableCountColumn                 = "count"
	LogicalThingTableOptionalCountColumn         = "optional_count"
	LogicalThingTableParentPhysicalThingIDColumn = "parent_physical_thing_id"
	LogicalThingTableParentLogicalThingIDColumn  = "parent_logical_thing_id"
)

var (
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

var (
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

		case "age":
			if v == nil {
				continue
			}

			temp1, err := types.ParseDuration(v)
			if err != nil {
				return wrapError(k, v, err)
			}

			temp2, ok := temp1.(time.Duration)
			if !ok {
				if temp1 != nil {
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uuage.UUID", temp1))
				}
			}

			m.Age = temp2

		case "optional_age":
			if v == nil {
				continue
			}

			temp1, err := types.ParseDuration(v)
			if err != nil {
				return wrapError(k, v, err)
			}

			temp2, ok := temp1.(time.Duration)
			if !ok {
				if temp1 != nil {
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uuoptional_age.UUID", temp1))
				}
			}

			m.OptionalAge = &temp2

		case "count":
			if v == nil {
				continue
			}

			temp1, err := types.ParseInt(v)
			if err != nil {
				return wrapError(k, v, err)
			}

			temp2, ok := temp1.(int64)
			if !ok {
				if temp1 != nil {
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uucount.UUID", temp1))
				}
			}

			m.Count = temp2

		case "optional_count":
			if v == nil {
				continue
			}

			temp1, err := types.ParseInt(v)
			if err != nil {
				return wrapError(k, v, err)
			}

			temp2, ok := temp1.(int64)
			if !ok {
				if temp1 != nil {
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uuoptional_count.UUID", temp1))
				}
			}

			m.OptionalCount = &temp2

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
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uuparent_physical_thing_id.UUID", temp1))
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
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uuparent_logical_thing_id.UUID", temp1))
				}
			}

			m.ParentLogicalThingID = &temp2

		}
	}

	return nil
}

func (m *LogicalThing) Reload(ctx context.Context, tx pgx.Tx, includeDeleteds ...bool) error {
	extraWhere := ""
	if len(includeDeleteds) > 0 && includeDeleteds[0] {
		if slices.Contains(LogicalThingTableColumns, "deleted_at") {
			extraWhere = "\n    AND (deleted_at IS null OR deleted_at IS NOT null)"
		}
	}

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
	m.Age = o.Age
	m.OptionalAge = o.OptionalAge
	m.Count = o.Count
	m.OptionalCount = o.OptionalCount
	m.ParentPhysicalThingID = o.ParentPhysicalThingID
	m.ParentPhysicalThingIDObject = o.ParentPhysicalThingIDObject
	m.ParentLogicalThingID = o.ParentLogicalThingID
	m.ParentLogicalThingIDObject = o.ParentLogicalThingIDObject
	m.ReferencedByLogicalThingParentLogicalThingIDObjects = o.ReferencedByLogicalThingParentLogicalThingIDObjects

	return nil
}

func (m *LogicalThing) Insert(ctx context.Context, tx pgx.Tx, setPrimaryKey bool, setZeroValues bool, forceSetValuesForFields ...string) error {
	columns := make([]string, 0)
	values := make([]any, 0)

	if setPrimaryKey && (setZeroValues || !types.IsZeroUUID(m.ID) || slices.Contains(forceSetValuesForFields, LogicalThingTableIDColumn) || isRequired(LogicalThingTableColumnLookup, LogicalThingTableIDColumn)) {
		columns = append(columns, LogicalThingTableIDColumn)

		v, err := types.FormatUUID(m.ID)
		if err != nil {
			return fmt.Errorf("failed to handle m.ID; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.CreatedAt) || slices.Contains(forceSetValuesForFields, LogicalThingTableCreatedAtColumn) || isRequired(LogicalThingTableColumnLookup, LogicalThingTableCreatedAtColumn) {
		columns = append(columns, LogicalThingTableCreatedAtColumn)

		v, err := types.FormatTime(m.CreatedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.CreatedAt; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.UpdatedAt) || slices.Contains(forceSetValuesForFields, LogicalThingTableUpdatedAtColumn) || isRequired(LogicalThingTableColumnLookup, LogicalThingTableUpdatedAtColumn) {
		columns = append(columns, LogicalThingTableUpdatedAtColumn)

		v, err := types.FormatTime(m.UpdatedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.UpdatedAt; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.DeletedAt) || slices.Contains(forceSetValuesForFields, LogicalThingTableDeletedAtColumn) || isRequired(LogicalThingTableColumnLookup, LogicalThingTableDeletedAtColumn) {
		columns = append(columns, LogicalThingTableDeletedAtColumn)

		v, err := types.FormatTime(m.DeletedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.DeletedAt; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroString(m.ExternalID) || slices.Contains(forceSetValuesForFields, LogicalThingTableExternalIDColumn) || isRequired(LogicalThingTableColumnLookup, LogicalThingTableExternalIDColumn) {
		columns = append(columns, LogicalThingTableExternalIDColumn)

		v, err := types.FormatString(m.ExternalID)
		if err != nil {
			return fmt.Errorf("failed to handle m.ExternalID; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroString(m.Name) || slices.Contains(forceSetValuesForFields, LogicalThingTableNameColumn) || isRequired(LogicalThingTableColumnLookup, LogicalThingTableNameColumn) {
		columns = append(columns, LogicalThingTableNameColumn)

		v, err := types.FormatString(m.Name)
		if err != nil {
			return fmt.Errorf("failed to handle m.Name; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroString(m.Type) || slices.Contains(forceSetValuesForFields, LogicalThingTableTypeColumn) || isRequired(LogicalThingTableColumnLookup, LogicalThingTableTypeColumn) {
		columns = append(columns, LogicalThingTableTypeColumn)

		v, err := types.FormatString(m.Type)
		if err != nil {
			return fmt.Errorf("failed to handle m.Type; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroStringArray(m.Tags) || slices.Contains(forceSetValuesForFields, LogicalThingTableTagsColumn) || isRequired(LogicalThingTableColumnLookup, LogicalThingTableTagsColumn) {
		columns = append(columns, LogicalThingTableTagsColumn)

		v, err := types.FormatStringArray(m.Tags)
		if err != nil {
			return fmt.Errorf("failed to handle m.Tags; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroHstore(m.Metadata) || slices.Contains(forceSetValuesForFields, LogicalThingTableMetadataColumn) || isRequired(LogicalThingTableColumnLookup, LogicalThingTableMetadataColumn) {
		columns = append(columns, LogicalThingTableMetadataColumn)

		v, err := types.FormatHstore(m.Metadata)
		if err != nil {
			return fmt.Errorf("failed to handle m.Metadata; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroJSON(m.RawData) || slices.Contains(forceSetValuesForFields, LogicalThingTableRawDataColumn) || isRequired(LogicalThingTableColumnLookup, LogicalThingTableRawDataColumn) {
		columns = append(columns, LogicalThingTableRawDataColumn)

		v, err := types.FormatJSON(m.RawData)
		if err != nil {
			return fmt.Errorf("failed to handle m.RawData; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroDuration(m.Age) || slices.Contains(forceSetValuesForFields, LogicalThingTableAgeColumn) || isRequired(LogicalThingTableColumnLookup, LogicalThingTableAgeColumn) {
		columns = append(columns, LogicalThingTableAgeColumn)

		v, err := types.FormatDuration(m.Age)
		if err != nil {
			return fmt.Errorf("failed to handle m.Age; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroDuration(m.OptionalAge) || slices.Contains(forceSetValuesForFields, LogicalThingTableOptionalAgeColumn) || isRequired(LogicalThingTableColumnLookup, LogicalThingTableOptionalAgeColumn) {
		columns = append(columns, LogicalThingTableOptionalAgeColumn)

		v, err := types.FormatDuration(m.OptionalAge)
		if err != nil {
			return fmt.Errorf("failed to handle m.OptionalAge; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroInt(m.Count) || slices.Contains(forceSetValuesForFields, LogicalThingTableCountColumn) || isRequired(LogicalThingTableColumnLookup, LogicalThingTableCountColumn) {
		columns = append(columns, LogicalThingTableCountColumn)

		v, err := types.FormatInt(m.Count)
		if err != nil {
			return fmt.Errorf("failed to handle m.Count; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroInt(m.OptionalCount) || slices.Contains(forceSetValuesForFields, LogicalThingTableOptionalCountColumn) || isRequired(LogicalThingTableColumnLookup, LogicalThingTableOptionalCountColumn) {
		columns = append(columns, LogicalThingTableOptionalCountColumn)

		v, err := types.FormatInt(m.OptionalCount)
		if err != nil {
			return fmt.Errorf("failed to handle m.OptionalCount; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroUUID(m.ParentPhysicalThingID) || slices.Contains(forceSetValuesForFields, LogicalThingTableParentPhysicalThingIDColumn) || isRequired(LogicalThingTableColumnLookup, LogicalThingTableParentPhysicalThingIDColumn) {
		columns = append(columns, LogicalThingTableParentPhysicalThingIDColumn)

		v, err := types.FormatUUID(m.ParentPhysicalThingID)
		if err != nil {
			return fmt.Errorf("failed to handle m.ParentPhysicalThingID; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroUUID(m.ParentLogicalThingID) || slices.Contains(forceSetValuesForFields, LogicalThingTableParentLogicalThingIDColumn) || isRequired(LogicalThingTableColumnLookup, LogicalThingTableParentLogicalThingIDColumn) {
		columns = append(columns, LogicalThingTableParentLogicalThingIDColumn)

		v, err := types.FormatUUID(m.ParentLogicalThingID)
		if err != nil {
			return fmt.Errorf("failed to handle m.ParentLogicalThingID; %v", err)
		}

		values = append(values, v)
	}

	ctx, cleanup := query.WithQueryID(ctx)
	defer cleanup()

	ctx = query.WithMaxDepth(ctx, nil)

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
		return fmt.Errorf("failed to insert %#+v; %v", m, err)
	}
	v := (*item)[LogicalThingTableIDColumn]

	if v == nil {
		return fmt.Errorf("failed to find %v in %#+v", LogicalThingTableIDColumn, item)
	}

	wrapError := func(err error) error {
		return fmt.Errorf(
			"failed to treat %v: %#+v as uuid.UUID: %v",
			LogicalThingTableIDColumn,
			(*item)[LogicalThingTableIDColumn],
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

func (m *LogicalThing) Update(ctx context.Context, tx pgx.Tx, setZeroValues bool, forceSetValuesForFields ...string) error {
	columns := make([]string, 0)
	values := make([]any, 0)

	if setZeroValues || !types.IsZeroTime(m.CreatedAt) || slices.Contains(forceSetValuesForFields, LogicalThingTableCreatedAtColumn) {
		columns = append(columns, LogicalThingTableCreatedAtColumn)

		v, err := types.FormatTime(m.CreatedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.CreatedAt; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.UpdatedAt) || slices.Contains(forceSetValuesForFields, LogicalThingTableUpdatedAtColumn) {
		columns = append(columns, LogicalThingTableUpdatedAtColumn)

		v, err := types.FormatTime(m.UpdatedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.UpdatedAt; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.DeletedAt) || slices.Contains(forceSetValuesForFields, LogicalThingTableDeletedAtColumn) {
		columns = append(columns, LogicalThingTableDeletedAtColumn)

		v, err := types.FormatTime(m.DeletedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.DeletedAt; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroString(m.ExternalID) || slices.Contains(forceSetValuesForFields, LogicalThingTableExternalIDColumn) {
		columns = append(columns, LogicalThingTableExternalIDColumn)

		v, err := types.FormatString(m.ExternalID)
		if err != nil {
			return fmt.Errorf("failed to handle m.ExternalID; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroString(m.Name) || slices.Contains(forceSetValuesForFields, LogicalThingTableNameColumn) {
		columns = append(columns, LogicalThingTableNameColumn)

		v, err := types.FormatString(m.Name)
		if err != nil {
			return fmt.Errorf("failed to handle m.Name; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroString(m.Type) || slices.Contains(forceSetValuesForFields, LogicalThingTableTypeColumn) {
		columns = append(columns, LogicalThingTableTypeColumn)

		v, err := types.FormatString(m.Type)
		if err != nil {
			return fmt.Errorf("failed to handle m.Type; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroStringArray(m.Tags) || slices.Contains(forceSetValuesForFields, LogicalThingTableTagsColumn) {
		columns = append(columns, LogicalThingTableTagsColumn)

		v, err := types.FormatStringArray(m.Tags)
		if err != nil {
			return fmt.Errorf("failed to handle m.Tags; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroHstore(m.Metadata) || slices.Contains(forceSetValuesForFields, LogicalThingTableMetadataColumn) {
		columns = append(columns, LogicalThingTableMetadataColumn)

		v, err := types.FormatHstore(m.Metadata)
		if err != nil {
			return fmt.Errorf("failed to handle m.Metadata; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroJSON(m.RawData) || slices.Contains(forceSetValuesForFields, LogicalThingTableRawDataColumn) {
		columns = append(columns, LogicalThingTableRawDataColumn)

		v, err := types.FormatJSON(m.RawData)
		if err != nil {
			return fmt.Errorf("failed to handle m.RawData; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroDuration(m.Age) || slices.Contains(forceSetValuesForFields, LogicalThingTableAgeColumn) {
		columns = append(columns, LogicalThingTableAgeColumn)

		v, err := types.FormatDuration(m.Age)
		if err != nil {
			return fmt.Errorf("failed to handle m.Age; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroDuration(m.OptionalAge) || slices.Contains(forceSetValuesForFields, LogicalThingTableOptionalAgeColumn) {
		columns = append(columns, LogicalThingTableOptionalAgeColumn)

		v, err := types.FormatDuration(m.OptionalAge)
		if err != nil {
			return fmt.Errorf("failed to handle m.OptionalAge; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroInt(m.Count) || slices.Contains(forceSetValuesForFields, LogicalThingTableCountColumn) {
		columns = append(columns, LogicalThingTableCountColumn)

		v, err := types.FormatInt(m.Count)
		if err != nil {
			return fmt.Errorf("failed to handle m.Count; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroInt(m.OptionalCount) || slices.Contains(forceSetValuesForFields, LogicalThingTableOptionalCountColumn) {
		columns = append(columns, LogicalThingTableOptionalCountColumn)

		v, err := types.FormatInt(m.OptionalCount)
		if err != nil {
			return fmt.Errorf("failed to handle m.OptionalCount; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroUUID(m.ParentPhysicalThingID) || slices.Contains(forceSetValuesForFields, LogicalThingTableParentPhysicalThingIDColumn) {
		columns = append(columns, LogicalThingTableParentPhysicalThingIDColumn)

		v, err := types.FormatUUID(m.ParentPhysicalThingID)
		if err != nil {
			return fmt.Errorf("failed to handle m.ParentPhysicalThingID; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroUUID(m.ParentLogicalThingID) || slices.Contains(forceSetValuesForFields, LogicalThingTableParentLogicalThingIDColumn) {
		columns = append(columns, LogicalThingTableParentLogicalThingIDColumn)

		v, err := types.FormatUUID(m.ParentLogicalThingID)
		if err != nil {
			return fmt.Errorf("failed to handle m.ParentLogicalThingID; %v", err)
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
		LogicalThingTable,
		columns,
		fmt.Sprintf("%v = $$??", LogicalThingTableIDColumn),
		LogicalThingTableColumns,
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

func (m *LogicalThing) Delete(ctx context.Context, tx pgx.Tx, hardDeletes ...bool) error {
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
		LogicalThingTable,
		fmt.Sprintf("%v = $$??", LogicalThingTableIDColumn),
		values...,
	)
	if err != nil {
		return fmt.Errorf("failed to delete %#+v; %v", m, err)
	}

	_ = m.Reload(ctx, tx, true)

	return nil
}

func (m *LogicalThing) LockTable(ctx context.Context, tx pgx.Tx, timeouts ...time.Duration) error {
	return query.LockTable(ctx, tx, LogicalThingTable, timeouts...)
}

func (m *LogicalThing) LockTableWithRetries(ctx context.Context, tx pgx.Tx, overallTimeout time.Duration, individualAttempttimeout time.Duration) error {
	return query.LockTableWithRetries(ctx, tx, LogicalThingTable, overallTimeout, individualAttempttimeout)
}

func (m *LogicalThing) AdvisoryLock(ctx context.Context, tx pgx.Tx, key int32, timeouts ...time.Duration) error {
	return query.AdvisoryLock(ctx, tx, LogicalThingTableNamespaceID, key, timeouts...)
}

func (m *LogicalThing) AdvisoryLockWithRetries(ctx context.Context, tx pgx.Tx, key int32, overallTimeout time.Duration, individualAttempttimeout time.Duration) error {
	return query.AdvisoryLockWithRetries(ctx, tx, LogicalThingTableNamespaceID, key, overallTimeout, individualAttempttimeout)
}

func SelectLogicalThings(ctx context.Context, tx pgx.Tx, where string, orderBy *string, limit *int, offset *int, values ...any) ([]*LogicalThing, int64, int64, int64, int64, error) {
	before := time.Now()

	if config.Debug() {
		log.Printf("entered SelectLogicalThings")

		defer func() {
			log.Printf("exited SelectLogicalThings in %s", time.Since(before))
		}()
	}
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

	forceLoad := query.ShouldLoad(ctx, LogicalThingTableColumnLookup[LogicalThingTablePrimaryKeyColumn], nil) ||
		query.ShouldLoad(ctx, nil, LogicalThingTableColumnLookup[LogicalThingTablePrimaryKeyColumn])

	var ok bool
	ctx, ok = query.HandleQueryPathGraphCycles(ctx, fmt.Sprintf("%s{%v}", LogicalThingTable, nil), !isLoadQuery)
	if !ok && !forceLoad {
		if config.Debug() {
			log.Printf("skipping SelectLogicalThing early (query.ShouldLoad(): %v, query.HandleQueryPathGraphCycles(): %v)", forceLoad, ok)
		}
		return []*LogicalThing{}, 0, 0, 0, 0, nil
	}

	items, count, totalCount, page, totalPages, err := query.Select(
		ctx,
		tx,
		LogicalThingTableColumnsWithTypeCasts,
		LogicalThingTable,
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

		if !types.IsZeroUUID(object.ParentPhysicalThingID) {
			ctx, ok := query.HandleQueryPathGraphCycles(ctx, fmt.Sprintf("%s{%v}", PhysicalThingTable, object.ParentPhysicalThingID), true)
			shouldLoad := query.ShouldLoad(ctx, LogicalThingTableColumnLookup[LogicalThingTableParentPhysicalThingIDColumn], nil)
			if ok || shouldLoad {
				thisBefore := time.Now()

				if config.Debug() {
					log.Printf("loading SelectLogicalThings->SelectPhysicalThing for object.ParentPhysicalThingIDObject")
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
					log.Printf("loaded SelectLogicalThings->SelectPhysicalThing for object.ParentPhysicalThingIDObject in %s", time.Since(thisBefore))
				}
			}
		}

		if !types.IsZeroUUID(object.ParentLogicalThingID) {
			ctx, ok := query.HandleQueryPathGraphCycles(ctx, fmt.Sprintf("%s{%v}", LogicalThingTable, object.ParentLogicalThingID), true)
			shouldLoad := query.ShouldLoad(ctx, LogicalThingTableColumnLookup[LogicalThingTableParentLogicalThingIDColumn], nil)
			if ok || shouldLoad {
				thisBefore := time.Now()

				if config.Debug() {
					log.Printf("loading SelectLogicalThings->SelectLogicalThing for object.ParentLogicalThingIDObject")
				}

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

				if config.Debug() {
					log.Printf("loaded SelectLogicalThings->SelectLogicalThing for object.ParentLogicalThingIDObject in %s", time.Since(thisBefore))
				}
			}
		}

		/*
			err = func() error {
				shouldLoad := query.ShouldLoad(ctx, nil, LogicalThingTableColumnLookup[LogicalThingTableParentLogicalThingIDColumn])
				ctx, ok := query.HandleQueryPathGraphCycles(ctx, fmt.Sprintf("__ReferencedBy__%s{%v}", LogicalThingTable, object.GetPrimaryKeyValue()), true)
				if ok || shouldLoad {
					thisBefore := time.Now()

					if config.Debug() {
						log.Printf("loading SelectLogicalThings->SelectLogicalThings for object.ReferencedByLogicalThingParentLogicalThingIDObjects")
					}

					object.ReferencedByLogicalThingParentLogicalThingIDObjects, _, _, _, _, err = SelectLogicalThings(
						ctx,
						tx,
						fmt.Sprintf("%v = $1", LogicalThingTableParentLogicalThingIDColumn),
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
						log.Printf("loaded SelectLogicalThings->SelectLogicalThings for object.ReferencedByLogicalThingParentLogicalThingIDObjects in %s", time.Since(thisBefore))
					}

				}

				return nil
			}()
			if err != nil {
				return nil, 0, 0, 0, 0, err
			}
		*/

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

func handleGetLogicalThings(arguments *server.SelectManyArguments, db *pgxpool.Pool) ([]*LogicalThing, int64, int64, int64, int64, error) {
	tx, err := db.Begin(arguments.Ctx)
	if err != nil {
		if config.Debug() {
			log.Printf("")
		}

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

func GetLogicalThingRouter(db *pgxpool.Pool, redisPool *redis.Pool, httpMiddlewares []server.HTTPMiddleware, objectMiddlewares []server.ObjectMiddleware, waitForChange server.WaitForChange) chi.Router {
	r := chi.NewRouter()

	for _, m := range httpMiddlewares {
		r.Use(m)
	}

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
			panic(err)
		}
		r.Get(getManyHandler.PathWithinRouter, getManyHandler.ServeHTTP)
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
			panic(err)
		}
		r.Get(getOneHandler.PathWithinRouter, getOneHandler.ServeHTTP)
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
			panic(err)
		}
		r.Post(postHandler.PathWithinRouter, postHandler.ServeHTTP)
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
			panic(err)
		}
		r.Put(putHandler.PathWithinRouter, putHandler.ServeHTTP)
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
			panic(err)
		}
		r.Patch(patchHandler.PathWithinRouter, patchHandler.ServeHTTP)
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
			panic(err)
		}
		r.Delete(deleteHandler.PathWithinRouter, deleteHandler.ServeHTTP)
	}()

	return r
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
		GetLogicalThingRouter,
	)
}
