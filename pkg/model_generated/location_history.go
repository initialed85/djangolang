package model_generated

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/netip"
	"slices"
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
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/exp/maps"
)

type LocationHistory struct {
	ID                          uuid.UUID      `json:"id"`
	CreatedAt                   time.Time      `json:"created_at"`
	UpdatedAt                   time.Time      `json:"updated_at"`
	DeletedAt                   *time.Time     `json:"deleted_at"`
	Timestamp                   time.Time      `json:"timestamp"`
	Point                       *pgtype.Vec2   `json:"point"`
	Polygon                     *[]pgtype.Vec2 `json:"polygon"`
	ParentPhysicalThingID       *uuid.UUID     `json:"parent_physical_thing_id"`
	ParentPhysicalThingIDObject *PhysicalThing `json:"parent_physical_thing_id_object"`
}

var LocationHistoryTable = "location_history"

var (
	LocationHistoryTableIDColumn                    = "id"
	LocationHistoryTableCreatedAtColumn             = "created_at"
	LocationHistoryTableUpdatedAtColumn             = "updated_at"
	LocationHistoryTableDeletedAtColumn             = "deleted_at"
	LocationHistoryTableTimestampColumn             = "timestamp"
	LocationHistoryTablePointColumn                 = "point"
	LocationHistoryTablePolygonColumn               = "polygon"
	LocationHistoryTableParentPhysicalThingIDColumn = "parent_physical_thing_id"
)

var (
	LocationHistoryTableIDColumnWithTypeCast                    = `"id" AS id`
	LocationHistoryTableCreatedAtColumnWithTypeCast             = `"created_at" AS created_at`
	LocationHistoryTableUpdatedAtColumnWithTypeCast             = `"updated_at" AS updated_at`
	LocationHistoryTableDeletedAtColumnWithTypeCast             = `"deleted_at" AS deleted_at`
	LocationHistoryTableTimestampColumnWithTypeCast             = `"timestamp" AS timestamp`
	LocationHistoryTablePointColumnWithTypeCast                 = `"point" AS point`
	LocationHistoryTablePolygonColumnWithTypeCast               = `"polygon" AS polygon`
	LocationHistoryTableParentPhysicalThingIDColumnWithTypeCast = `"parent_physical_thing_id" AS parent_physical_thing_id`
)

var LocationHistoryTableColumns = []string{
	LocationHistoryTableIDColumn,
	LocationHistoryTableCreatedAtColumn,
	LocationHistoryTableUpdatedAtColumn,
	LocationHistoryTableDeletedAtColumn,
	LocationHistoryTableTimestampColumn,
	LocationHistoryTablePointColumn,
	LocationHistoryTablePolygonColumn,
	LocationHistoryTableParentPhysicalThingIDColumn,
}

var LocationHistoryTableColumnsWithTypeCasts = []string{
	LocationHistoryTableIDColumnWithTypeCast,
	LocationHistoryTableCreatedAtColumnWithTypeCast,
	LocationHistoryTableUpdatedAtColumnWithTypeCast,
	LocationHistoryTableDeletedAtColumnWithTypeCast,
	LocationHistoryTableTimestampColumnWithTypeCast,
	LocationHistoryTablePointColumnWithTypeCast,
	LocationHistoryTablePolygonColumnWithTypeCast,
	LocationHistoryTableParentPhysicalThingIDColumnWithTypeCast,
}

var LocationHistoryIntrospectedTable *introspect.Table

var LocationHistoryTableColumnLookup map[string]*introspect.Column

var (
	LocationHistoryTablePrimaryKeyColumn = LocationHistoryTableIDColumn
)

func init() {
	LocationHistoryIntrospectedTable = tableByName[LocationHistoryTable]

	/* only needed during templating */
	if LocationHistoryIntrospectedTable == nil {
		LocationHistoryIntrospectedTable = &introspect.Table{}
	}

	LocationHistoryTableColumnLookup = LocationHistoryIntrospectedTable.ColumnByName
}

type LocationHistoryOnePathParams struct {
	PrimaryKey uuid.UUID `json:"primaryKey"`
}

type LocationHistoryLoadQueryParams struct {
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

func (m *LocationHistory) GetPrimaryKeyColumn() string {
	return LocationHistoryTablePrimaryKeyColumn
}

func (m *LocationHistory) GetPrimaryKeyValue() any {
	return m.ID
}

func (m *LocationHistory) FromItem(item map[string]any) error {
	if item == nil {
		return fmt.Errorf(
			"item unexpectedly nil during LocationHistoryFromItem",
		)
	}

	if len(item) == 0 {
		return fmt.Errorf(
			"item unexpectedly empty during LocationHistoryFromItem",
		)
	}

	wrapError := func(k string, v any, err error) error {
		return fmt.Errorf("%v: %#+v; error; %v", k, v, err)
	}

	for k, v := range item {
		_, ok := LocationHistoryTableColumnLookup[k]
		if !ok {
			return fmt.Errorf(
				"item contained unexpected key %#+v during LocationHistoryFromItem; item: %#+v",
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

		case "timestamp":
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
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uutimestamp.UUID", temp1))
				}
			}

			m.Timestamp = temp2

		case "point":
			if v == nil {
				continue
			}

			temp1, err := types.ParsePoint(v)
			if err != nil {
				return wrapError(k, v, err)
			}

			temp2, ok := temp1.(pgtype.Vec2)
			if !ok {
				if temp1 != nil {
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uupoint.UUID", temp1))
				}
			}

			m.Point = &temp2

		case "polygon":
			if v == nil {
				continue
			}

			temp1, err := types.ParsePolygon(v)
			if err != nil {
				return wrapError(k, v, err)
			}

			temp2, ok := temp1.([]pgtype.Vec2)
			if !ok {
				if temp1 != nil {
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uupolygon.UUID", temp1))
				}
			}

			m.Polygon = &temp2

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

		}
	}

	return nil
}

func (m *LocationHistory) Reload(ctx context.Context, tx pgx.Tx, includeDeleteds ...bool) error {
	extraWhere := ""
	if len(includeDeleteds) > 0 && includeDeleteds[0] {
		if slices.Contains(LocationHistoryTableColumns, "deleted_at") {
			extraWhere = "\n    AND (deleted_at IS null OR deleted_at IS NOT null)"
		}
	}

	ctx, cleanup := query.WithQueryID(ctx)
	defer cleanup()

	t, err := SelectLocationHistory(
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
	m.Timestamp = t.Timestamp
	m.Point = t.Point
	m.Polygon = t.Polygon
	m.ParentPhysicalThingID = t.ParentPhysicalThingID
	m.ParentPhysicalThingIDObject = t.ParentPhysicalThingIDObject

	return nil
}

func (m *LocationHistory) Insert(ctx context.Context, tx pgx.Tx, setPrimaryKey bool, setZeroValues bool, forceSetValuesForFields ...string) error {
	columns := make([]string, 0)
	values := make([]any, 0)

	if setPrimaryKey && (setZeroValues || !types.IsZeroUUID(m.ID) || slices.Contains(forceSetValuesForFields, LocationHistoryTableIDColumn) || isRequired(LocationHistoryTableColumnLookup, LocationHistoryTableIDColumn)) {
		columns = append(columns, LocationHistoryTableIDColumn)

		v, err := types.FormatUUID(m.ID)
		if err != nil {
			return fmt.Errorf("failed to handle m.ID: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.CreatedAt) || slices.Contains(forceSetValuesForFields, LocationHistoryTableCreatedAtColumn) || isRequired(LocationHistoryTableColumnLookup, LocationHistoryTableCreatedAtColumn) {
		columns = append(columns, LocationHistoryTableCreatedAtColumn)

		v, err := types.FormatTime(m.CreatedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.CreatedAt: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.UpdatedAt) || slices.Contains(forceSetValuesForFields, LocationHistoryTableUpdatedAtColumn) || isRequired(LocationHistoryTableColumnLookup, LocationHistoryTableUpdatedAtColumn) {
		columns = append(columns, LocationHistoryTableUpdatedAtColumn)

		v, err := types.FormatTime(m.UpdatedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.UpdatedAt: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.DeletedAt) || slices.Contains(forceSetValuesForFields, LocationHistoryTableDeletedAtColumn) || isRequired(LocationHistoryTableColumnLookup, LocationHistoryTableDeletedAtColumn) {
		columns = append(columns, LocationHistoryTableDeletedAtColumn)

		v, err := types.FormatTime(m.DeletedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.DeletedAt: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.Timestamp) || slices.Contains(forceSetValuesForFields, LocationHistoryTableTimestampColumn) || isRequired(LocationHistoryTableColumnLookup, LocationHistoryTableTimestampColumn) {
		columns = append(columns, LocationHistoryTableTimestampColumn)

		v, err := types.FormatTime(m.Timestamp)
		if err != nil {
			return fmt.Errorf("failed to handle m.Timestamp: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroPoint(m.Point) || slices.Contains(forceSetValuesForFields, LocationHistoryTablePointColumn) || isRequired(LocationHistoryTableColumnLookup, LocationHistoryTablePointColumn) {
		columns = append(columns, LocationHistoryTablePointColumn)

		v, err := types.FormatPoint(m.Point)
		if err != nil {
			return fmt.Errorf("failed to handle m.Point: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroPolygon(m.Polygon) || slices.Contains(forceSetValuesForFields, LocationHistoryTablePolygonColumn) || isRequired(LocationHistoryTableColumnLookup, LocationHistoryTablePolygonColumn) {
		columns = append(columns, LocationHistoryTablePolygonColumn)

		v, err := types.FormatPolygon(m.Polygon)
		if err != nil {
			return fmt.Errorf("failed to handle m.Polygon: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroUUID(m.ParentPhysicalThingID) || slices.Contains(forceSetValuesForFields, LocationHistoryTableParentPhysicalThingIDColumn) || isRequired(LocationHistoryTableColumnLookup, LocationHistoryTableParentPhysicalThingIDColumn) {
		columns = append(columns, LocationHistoryTableParentPhysicalThingIDColumn)

		v, err := types.FormatUUID(m.ParentPhysicalThingID)
		if err != nil {
			return fmt.Errorf("failed to handle m.ParentPhysicalThingID: %v", err)
		}

		values = append(values, v)
	}

	ctx, cleanup := query.WithQueryID(ctx)
	defer cleanup()

	item, err := query.Insert(
		ctx,
		tx,
		LocationHistoryTable,
		columns,
		nil,
		false,
		false,
		LocationHistoryTableColumns,
		values...,
	)
	if err != nil {
		return fmt.Errorf("failed to insert %#+v; %v", m, err)
	}
	v := (*item)[LocationHistoryTableIDColumn]

	if v == nil {
		return fmt.Errorf("failed to find %v in %#+v", LocationHistoryTableIDColumn, item)
	}

	wrapError := func(err error) error {
		return fmt.Errorf(
			"failed to treat %v: %#+v as uuid.UUID: %v",
			LocationHistoryTableIDColumn,
			(*item)[LocationHistoryTableIDColumn],
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
		return fmt.Errorf("failed to reload after insert: %v", err)
	}

	return nil
}

func (m *LocationHistory) Update(ctx context.Context, tx pgx.Tx, setZeroValues bool, forceSetValuesForFields ...string) error {
	columns := make([]string, 0)
	values := make([]any, 0)

	if setZeroValues || !types.IsZeroTime(m.CreatedAt) || slices.Contains(forceSetValuesForFields, LocationHistoryTableCreatedAtColumn) {
		columns = append(columns, LocationHistoryTableCreatedAtColumn)

		v, err := types.FormatTime(m.CreatedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.CreatedAt: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.UpdatedAt) || slices.Contains(forceSetValuesForFields, LocationHistoryTableUpdatedAtColumn) {
		columns = append(columns, LocationHistoryTableUpdatedAtColumn)

		v, err := types.FormatTime(m.UpdatedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.UpdatedAt: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.DeletedAt) || slices.Contains(forceSetValuesForFields, LocationHistoryTableDeletedAtColumn) {
		columns = append(columns, LocationHistoryTableDeletedAtColumn)

		v, err := types.FormatTime(m.DeletedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.DeletedAt: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.Timestamp) || slices.Contains(forceSetValuesForFields, LocationHistoryTableTimestampColumn) {
		columns = append(columns, LocationHistoryTableTimestampColumn)

		v, err := types.FormatTime(m.Timestamp)
		if err != nil {
			return fmt.Errorf("failed to handle m.Timestamp: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroPoint(m.Point) || slices.Contains(forceSetValuesForFields, LocationHistoryTablePointColumn) {
		columns = append(columns, LocationHistoryTablePointColumn)

		v, err := types.FormatPoint(m.Point)
		if err != nil {
			return fmt.Errorf("failed to handle m.Point: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroPolygon(m.Polygon) || slices.Contains(forceSetValuesForFields, LocationHistoryTablePolygonColumn) {
		columns = append(columns, LocationHistoryTablePolygonColumn)

		v, err := types.FormatPolygon(m.Polygon)
		if err != nil {
			return fmt.Errorf("failed to handle m.Polygon: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroUUID(m.ParentPhysicalThingID) || slices.Contains(forceSetValuesForFields, LocationHistoryTableParentPhysicalThingIDColumn) {
		columns = append(columns, LocationHistoryTableParentPhysicalThingIDColumn)

		v, err := types.FormatUUID(m.ParentPhysicalThingID)
		if err != nil {
			return fmt.Errorf("failed to handle m.ParentPhysicalThingID: %v", err)
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
		LocationHistoryTable,
		columns,
		fmt.Sprintf("%v = $$??", LocationHistoryTableIDColumn),
		LocationHistoryTableColumns,
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

func (m *LocationHistory) Delete(ctx context.Context, tx pgx.Tx, hardDeletes ...bool) error {
	hardDelete := false
	if len(hardDeletes) > 0 {
		hardDelete = hardDeletes[0]
	}

	if !hardDelete && slices.Contains(LocationHistoryTableColumns, "deleted_at") {
		m.DeletedAt = helpers.Ptr(time.Now().UTC())
		err := m.Update(ctx, tx, false, "deleted_at")
		if err != nil {
			return fmt.Errorf("failed to soft-delete (update) %#+v; %v", m, err)
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
		LocationHistoryTable,
		fmt.Sprintf("%v = $$??", LocationHistoryTableIDColumn),
		values...,
	)
	if err != nil {
		return fmt.Errorf("failed to delete %#+v; %v", m, err)
	}

	_ = m.Reload(ctx, tx, true)

	return nil
}

func (m *LocationHistory) LockTable(ctx context.Context, tx pgx.Tx, noWait bool) error {
	return query.LockTable(ctx, tx, LocationHistoryTable, noWait)
}

func SelectLocationHistories(ctx context.Context, tx pgx.Tx, where string, orderBy *string, limit *int, offset *int, values ...any) ([]*LocationHistory, error) {
	if slices.Contains(LocationHistoryTableColumns, "deleted_at") {
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
		LocationHistoryTableColumnsWithTypeCasts,
		LocationHistoryTable,
		where,
		orderBy,
		limit,
		offset,
		values...,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to call SelectLocationHistorys; err: %v", err)
	}

	objects := make([]*LocationHistory, 0)

	for _, item := range *items {
		object := &LocationHistory{}

		err = object.FromItem(item)
		if err != nil {
			return nil, err
		}

		thatCtx := ctx

		thatCtx, ok1 := query.HandleQueryPathGraphCycles(ctx, fmt.Sprintf("%s{%v}", LocationHistoryTable, object.GetPrimaryKeyValue()))
		thatCtx, ok2 := query.HandleQueryPathGraphCycles(thatCtx, fmt.Sprintf("__ReferencedBy__%s{%v}", LocationHistoryTable, object.GetPrimaryKeyValue()))
		if !(ok1 && ok2) {
			continue
		}

		_ = thatCtx

		if !types.IsZeroUUID(object.ParentPhysicalThingID) {
			thisCtx := thatCtx
			thisCtx, ok1 := query.HandleQueryPathGraphCycles(thisCtx, fmt.Sprintf("%s{%v}", PhysicalThingTable, object.ParentPhysicalThingID))
			thisCtx, ok2 := query.HandleQueryPathGraphCycles(thisCtx, fmt.Sprintf("__ReferencedBy__%s{%v}", PhysicalThingTable, object.ParentPhysicalThingID))
			if ok1 && ok2 {
				object.ParentPhysicalThingIDObject, err = SelectPhysicalThing(
					thisCtx,
					tx,
					fmt.Sprintf("%v = $1", PhysicalThingTablePrimaryKeyColumn),
					object.ParentPhysicalThingID,
				)
				if err != nil {
					if !errors.Is(err, sql.ErrNoRows) {
						return nil, err
					}
				}
			}
		}

		objects = append(objects, object)
	}

	return objects, nil
}

func SelectLocationHistory(ctx context.Context, tx pgx.Tx, where string, values ...any) (*LocationHistory, error) {
	ctx, cleanup := query.WithQueryID(ctx)
	defer cleanup()

	objects, err := SelectLocationHistories(
		ctx,
		tx,
		where,
		nil,
		helpers.Ptr(2),
		helpers.Ptr(0),
		values...,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to call SelectLocationHistory; err: %v", err)
	}

	if len(objects) > 1 {
		return nil, fmt.Errorf("attempt to call SelectLocationHistory returned more than 1 row")
	}

	if len(objects) < 1 {
		return nil, sql.ErrNoRows
	}

	object := objects[0]

	return object, nil
}

func handleGetLocationHistories(arguments *server.SelectManyArguments, db *pgxpool.Pool) ([]*LocationHistory, error) {
	tx, err := db.Begin(arguments.Ctx)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = tx.Rollback(arguments.Ctx)
	}()

	objects, err := SelectLocationHistories(arguments.Ctx, tx, arguments.Where, arguments.OrderBy, arguments.Limit, arguments.Offset, arguments.Values...)
	if err != nil {
		return nil, err
	}

	err = tx.Commit(arguments.Ctx)
	if err != nil {
		return nil, err
	}

	return objects, nil
}

func handleGetLocationHistory(arguments *server.SelectOneArguments, db *pgxpool.Pool, primaryKey uuid.UUID) ([]*LocationHistory, error) {
	tx, err := db.Begin(arguments.Ctx)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = tx.Rollback(arguments.Ctx)
	}()

	object, err := SelectLocationHistory(arguments.Ctx, tx, arguments.Where, arguments.Values...)
	if err != nil {
		return nil, err
	}

	err = tx.Commit(arguments.Ctx)
	if err != nil {
		return nil, err
	}

	return []*LocationHistory{object}, nil
}

func handlePostLocationHistorys(arguments *server.LoadArguments, db *pgxpool.Pool, waitForChange server.WaitForChange, objects []*LocationHistory, forceSetValuesForFieldsByObjectIndex [][]string) ([]*LocationHistory, error) {
	tx, err := db.Begin(arguments.Ctx)
	if err != nil {
		err = fmt.Errorf("failed to begin DB transaction: %v", err)
		return nil, err
	}

	defer func() {
		_ = tx.Rollback(arguments.Ctx)
	}()

	xid, err := query.GetXid(arguments.Ctx, tx)
	if err != nil {
		err = fmt.Errorf("failed to get xid: %v", err)
		return nil, err
	}
	_ = xid

	for i, object := range objects {
		err = object.Insert(arguments.Ctx, tx, false, false, forceSetValuesForFieldsByObjectIndex[i]...)
		if err != nil {
			err = fmt.Errorf("failed to insert %#+v; %v", object, err)
			return nil, err
		}

		objects[i] = object
	}

	errs := make(chan error, 1)
	go func() {
		_, err = waitForChange(arguments.Ctx, []stream.Action{stream.INSERT}, LocationHistoryTable, xid)
		if err != nil {
			err = fmt.Errorf("failed to wait for change: %v", err)
			errs <- err
			return
		}

		errs <- nil
	}()

	err = tx.Commit(arguments.Ctx)
	if err != nil {
		err = fmt.Errorf("failed to commit DB transaction: %v", err)
		return nil, err
	}

	select {
	case <-arguments.Ctx.Done():
		err = fmt.Errorf("context canceled")
		return nil, err
	case err = <-errs:
		if err != nil {
			return nil, err
		}
	}

	return objects, nil
}

func handlePutLocationHistory(arguments *server.LoadArguments, db *pgxpool.Pool, waitForChange server.WaitForChange, object *LocationHistory) ([]*LocationHistory, error) {
	tx, err := db.Begin(arguments.Ctx)
	if err != nil {
		err = fmt.Errorf("failed to begin DB transaction: %v", err)
		return nil, err
	}

	defer func() {
		_ = tx.Rollback(arguments.Ctx)
	}()

	xid, err := query.GetXid(arguments.Ctx, tx)
	if err != nil {
		err = fmt.Errorf("failed to get xid: %v", err)
		return nil, err
	}
	_ = xid

	err = object.Update(arguments.Ctx, tx, true)
	if err != nil {
		err = fmt.Errorf("failed to update %#+v; %v", object, err)
		return nil, err
	}

	errs := make(chan error, 1)
	go func() {
		_, err = waitForChange(arguments.Ctx, []stream.Action{stream.UPDATE, stream.SOFT_DELETE, stream.SOFT_RESTORE, stream.SOFT_UPDATE}, LocationHistoryTable, xid)
		if err != nil {
			err = fmt.Errorf("failed to wait for change: %v", err)
			errs <- err
			return
		}

		errs <- nil
	}()

	err = tx.Commit(arguments.Ctx)
	if err != nil {
		err = fmt.Errorf("failed to commit DB transaction: %v", err)
		return nil, err
	}

	select {
	case <-arguments.Ctx.Done():
		err = fmt.Errorf("context canceled")
		return nil, err
	case err = <-errs:
		if err != nil {
			return nil, err
		}
	}

	return []*LocationHistory{object}, nil
}

func handlePatchLocationHistory(arguments *server.LoadArguments, db *pgxpool.Pool, waitForChange server.WaitForChange, object *LocationHistory, forceSetValuesForFields []string) ([]*LocationHistory, error) {
	tx, err := db.Begin(arguments.Ctx)
	if err != nil {
		err = fmt.Errorf("failed to begin DB transaction: %v", err)
		return nil, err
	}

	defer func() {
		_ = tx.Rollback(arguments.Ctx)
	}()

	xid, err := query.GetXid(arguments.Ctx, tx)
	if err != nil {
		err = fmt.Errorf("failed to get xid: %v", err)
		return nil, err
	}
	_ = xid

	err = object.Update(arguments.Ctx, tx, false, forceSetValuesForFields...)
	if err != nil {
		err = fmt.Errorf("failed to update %#+v; %v", object, err)
		return nil, err
	}

	errs := make(chan error, 1)
	go func() {
		_, err = waitForChange(arguments.Ctx, []stream.Action{stream.UPDATE, stream.SOFT_DELETE, stream.SOFT_RESTORE, stream.SOFT_UPDATE}, LocationHistoryTable, xid)
		if err != nil {
			err = fmt.Errorf("failed to wait for change: %v", err)
			errs <- err
			return
		}

		errs <- nil
	}()

	err = tx.Commit(arguments.Ctx)
	if err != nil {
		err = fmt.Errorf("failed to commit DB transaction: %v", err)
		return nil, err
	}

	select {
	case <-arguments.Ctx.Done():
		err = fmt.Errorf("context canceled")
		return nil, err
	case err = <-errs:
		if err != nil {
			return nil, err
		}
	}

	return []*LocationHistory{object}, nil
}

func handleDeleteLocationHistory(arguments *server.LoadArguments, db *pgxpool.Pool, waitForChange server.WaitForChange, object *LocationHistory) error {
	tx, err := db.Begin(arguments.Ctx)
	if err != nil {
		err = fmt.Errorf("failed to begin DB transaction: %v", err)
		return err
	}

	defer func() {
		_ = tx.Rollback(arguments.Ctx)
	}()

	xid, err := query.GetXid(arguments.Ctx, tx)
	if err != nil {
		err = fmt.Errorf("failed to get xid: %v", err)
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
		_, err = waitForChange(arguments.Ctx, []stream.Action{stream.DELETE, stream.SOFT_DELETE}, LocationHistoryTable, xid)
		if err != nil {
			err = fmt.Errorf("failed to wait for change: %v", err)
			errs <- err
			return
		}

		errs <- nil
	}()

	err = tx.Commit(arguments.Ctx)
	if err != nil {
		err = fmt.Errorf("failed to commit DB transaction: %v", err)
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

func GetLocationHistoryRouter(db *pgxpool.Pool, redisPool *redis.Pool, httpMiddlewares []server.HTTPMiddleware, objectMiddlewares []server.ObjectMiddleware, waitForChange server.WaitForChange) chi.Router {
	r := chi.NewRouter()

	for _, m := range httpMiddlewares {
		r.Use(m)
	}

	getManyHandler, err := server.GetCustomHTTPHandler(
		http.MethodGet,
		"/",
		http.StatusOK,
		func(
			ctx context.Context,
			pathParams server.EmptyPathParams,
			queryParams map[string]any,
			req server.EmptyRequest,
			rawReq any,
		) (*helpers.TypedResponse[LocationHistory], error) {
			redisConn := redisPool.Get()
			defer func() {
				_ = redisConn.Close()
			}()

			arguments, err := server.GetSelectManyArguments(ctx, queryParams, LocationHistoryIntrospectedTable, nil, nil)
			if err != nil {
				return nil, err
			}

			cachedObjectsAsJSON, cacheHit, err := helpers.GetCachedObjectsAsJSON(arguments.RequestHash, redisConn)
			if err != nil {
				return nil, err
			}

			if cacheHit {
				var cachedObjects []*LocationHistory
				err = json.Unmarshal(cachedObjectsAsJSON, &cachedObjects)
				if err != nil {
					return nil, err
				}

				return &helpers.TypedResponse[LocationHistory]{
					Status:  http.StatusOK,
					Success: true,
					Error:   nil,
					Objects: cachedObjects,
				}, nil
			}

			objects, err := handleGetLocationHistories(arguments, db)
			if err != nil {
				return nil, err
			}

			objectsAsJSON, err := json.Marshal(objects)
			if err != nil {
				return nil, err
			}

			err = helpers.StoreCachedResponse(arguments.RequestHash, redisConn, string(objectsAsJSON))
			if err != nil {
				log.Printf("warning: %v", err)
			}

			return &helpers.TypedResponse[LocationHistory]{
				Status:  http.StatusOK,
				Success: true,
				Error:   nil,
				Objects: objects,
			}, nil
		},
	)
	if err != nil {
		panic(err)
	}
	r.Get("/", getManyHandler.ServeHTTP)

	getOneHandler, err := server.GetCustomHTTPHandler(
		http.MethodGet,
		"/{primaryKey}",
		http.StatusOK,
		func(
			ctx context.Context,
			pathParams LocationHistoryOnePathParams,
			queryParams LocationHistoryLoadQueryParams,
			req server.EmptyRequest,
			rawReq any,
		) (*helpers.TypedResponse[LocationHistory], error) {
			redisConn := redisPool.Get()
			defer func() {
				_ = redisConn.Close()
			}()

			arguments, err := server.GetSelectOneArguments(ctx, queryParams.Depth, LocationHistoryIntrospectedTable, pathParams.PrimaryKey, nil, nil)
			if err != nil {
				return nil, err
			}

			cachedObjectsAsJSON, cacheHit, err := helpers.GetCachedObjectsAsJSON(arguments.RequestHash, redisConn)
			if err != nil {
				return nil, err
			}

			if cacheHit {
				var cachedObjects []*LocationHistory
				err = json.Unmarshal(cachedObjectsAsJSON, &cachedObjects)
				if err != nil {
					return nil, err
				}

				return &helpers.TypedResponse[LocationHistory]{
					Status:  http.StatusOK,
					Success: true,
					Error:   nil,
					Objects: cachedObjects,
				}, nil
			}

			objects, err := handleGetLocationHistory(arguments, db, pathParams.PrimaryKey)
			if err != nil {
				return nil, err
			}

			objectsAsJSON, err := json.Marshal(objects)
			if err != nil {
				return nil, err
			}

			err = helpers.StoreCachedResponse(arguments.RequestHash, redisConn, string(objectsAsJSON))
			if err != nil {
				log.Printf("warning: %v", err)
			}

			return &helpers.TypedResponse[LocationHistory]{
				Status:  http.StatusOK,
				Success: true,
				Error:   nil,
				Objects: objects,
			}, nil
		},
	)
	if err != nil {
		panic(err)
	}
	r.Get("/{primaryKey}", getOneHandler.ServeHTTP)

	postHandler, err := server.GetCustomHTTPHandler(
		http.MethodPost,
		"/",
		http.StatusCreated,
		func(
			ctx context.Context,
			pathParams server.EmptyPathParams,
			queryParams LocationHistoryLoadQueryParams,
			req []*LocationHistory,
			rawReq any,
		) (*helpers.TypedResponse[LocationHistory], error) {
			allRawItems, ok := rawReq.([]any)
			if !ok {
				return nil, fmt.Errorf("failed to cast %#+v to []map[string]any", rawReq)
			}

			allItems := make([]map[string]any, 0)
			for _, rawItem := range allRawItems {
				item, ok := rawItem.(map[string]any)
				if !ok {
					return nil, fmt.Errorf("failed to cast %#+v to map[string]any", rawItem)
				}

				allItems = append(allItems, item)
			}

			forceSetValuesForFieldsByObjectIndex := make([][]string, 0)
			for _, item := range allItems {
				forceSetValuesForFields := make([]string, 0)
				for _, possibleField := range maps.Keys(item) {
					if !slices.Contains(LocationHistoryTableColumns, possibleField) {
						continue
					}

					forceSetValuesForFields = append(forceSetValuesForFields, possibleField)
				}
				forceSetValuesForFieldsByObjectIndex = append(forceSetValuesForFieldsByObjectIndex, forceSetValuesForFields)
			}

			arguments, err := server.GetLoadArguments(ctx, queryParams.Depth)
			if err != nil {
				return nil, err
			}

			objects, err := handlePostLocationHistorys(arguments, db, waitForChange, req, forceSetValuesForFieldsByObjectIndex)
			if err != nil {
				return nil, err
			}

			return &helpers.TypedResponse[LocationHistory]{
				Status:  http.StatusCreated,
				Success: true,
				Error:   nil,
				Objects: objects,
			}, nil
		},
	)
	if err != nil {
		panic(err)
	}
	r.Post("/", postHandler.ServeHTTP)

	putHandler, err := server.GetCustomHTTPHandler(
		http.MethodPatch,
		"/{primaryKey}",
		http.StatusOK,
		func(
			ctx context.Context,
			pathParams LocationHistoryOnePathParams,
			queryParams LocationHistoryLoadQueryParams,
			req LocationHistory,
			rawReq any,
		) (*helpers.TypedResponse[LocationHistory], error) {
			item, ok := rawReq.(map[string]any)
			if !ok {
				return nil, fmt.Errorf("failed to cast %#+v to map[string]any", item)
			}

			arguments, err := server.GetLoadArguments(ctx, queryParams.Depth)
			if err != nil {
				return nil, err
			}

			object := &req
			object.ID = pathParams.PrimaryKey

			objects, err := handlePutLocationHistory(arguments, db, waitForChange, object)
			if err != nil {
				return nil, err
			}

			return &helpers.TypedResponse[LocationHistory]{
				Status:  http.StatusOK,
				Success: true,
				Error:   nil,
				Objects: objects,
			}, nil
		},
	)
	if err != nil {
		panic(err)
	}
	r.Put("/{primaryKey}", putHandler.ServeHTTP)

	patchHandler, err := server.GetCustomHTTPHandler(
		http.MethodPatch,
		"/{primaryKey}",
		http.StatusOK,
		func(
			ctx context.Context,
			pathParams LocationHistoryOnePathParams,
			queryParams LocationHistoryLoadQueryParams,
			req LocationHistory,
			rawReq any,
		) (*helpers.TypedResponse[LocationHistory], error) {
			item, ok := rawReq.(map[string]any)
			if !ok {
				return nil, fmt.Errorf("failed to cast %#+v to map[string]any", item)
			}

			forceSetValuesForFields := make([]string, 0)
			for _, possibleField := range maps.Keys(item) {
				if !slices.Contains(LocationHistoryTableColumns, possibleField) {
					continue
				}

				forceSetValuesForFields = append(forceSetValuesForFields, possibleField)
			}

			arguments, err := server.GetLoadArguments(ctx, queryParams.Depth)
			if err != nil {
				return nil, err
			}

			object := &req
			object.ID = pathParams.PrimaryKey

			objects, err := handlePatchLocationHistory(arguments, db, waitForChange, object, forceSetValuesForFields)
			if err != nil {
				return nil, err
			}

			return &helpers.TypedResponse[LocationHistory]{
				Status:  http.StatusOK,
				Success: true,
				Error:   nil,
				Objects: objects,
			}, nil
		},
	)
	if err != nil {
		panic(err)
	}
	r.Patch("/{primaryKey}", patchHandler.ServeHTTP)

	deleteHandler, err := server.GetCustomHTTPHandler(
		http.MethodDelete,
		"/{primaryKey}",
		http.StatusNoContent,
		func(
			ctx context.Context,
			pathParams LocationHistoryOnePathParams,
			queryParams LocationHistoryLoadQueryParams,
			req server.EmptyRequest,
			rawReq any,
		) (*server.EmptyResponse, error) {
			arguments := &server.LoadArguments{
				Ctx: ctx,
			}

			object := &LocationHistory{}
			object.ID = pathParams.PrimaryKey

			err := handleDeleteLocationHistory(arguments, db, waitForChange, object)
			if err != nil {
				return nil, err
			}

			return nil, nil
		},
	)
	if err != nil {
		panic(err)
	}
	r.Delete("/{primaryKey}", deleteHandler.ServeHTTP)

	return r
}

func NewLocationHistoryFromItem(item map[string]any) (any, error) {
	object := &LocationHistory{}

	err := object.FromItem(item)
	if err != nil {
		return nil, err
	}

	return object, nil
}

func init() {
	register(
		LocationHistoryTable,
		LocationHistory{},
		NewLocationHistoryFromItem,
		"/location-histories",
		GetLocationHistoryRouter,
	)
}
