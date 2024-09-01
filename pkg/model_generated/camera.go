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

type Camera struct {
	ID                                   uuid.UUID    `json:"id"`
	CreatedAt                            time.Time    `json:"created_at"`
	UpdatedAt                            time.Time    `json:"updated_at"`
	DeletedAt                            *time.Time   `json:"deleted_at"`
	Name                                 string       `json:"name"`
	StreamURL                            string       `json:"stream_url"`
	LastSeen                             *time.Time   `json:"last_seen"`
	ReferencedByVideoCameraIDObjects     []*Video     `json:"referenced_by_video_camera_id_objects"`
	ReferencedByDetectionCameraIDObjects []*Detection `json:"referenced_by_detection_camera_id_objects"`
}

var CameraTable = "camera"

var (
	CameraTableIDColumn        = "id"
	CameraTableCreatedAtColumn = "created_at"
	CameraTableUpdatedAtColumn = "updated_at"
	CameraTableDeletedAtColumn = "deleted_at"
	CameraTableNameColumn      = "name"
	CameraTableStreamURLColumn = "stream_url"
	CameraTableLastSeenColumn  = "last_seen"
)

var (
	CameraTableIDColumnWithTypeCast        = `"id" AS id`
	CameraTableCreatedAtColumnWithTypeCast = `"created_at" AS created_at`
	CameraTableUpdatedAtColumnWithTypeCast = `"updated_at" AS updated_at`
	CameraTableDeletedAtColumnWithTypeCast = `"deleted_at" AS deleted_at`
	CameraTableNameColumnWithTypeCast      = `"name" AS name`
	CameraTableStreamURLColumnWithTypeCast = `"stream_url" AS stream_url`
	CameraTableLastSeenColumnWithTypeCast  = `"last_seen" AS last_seen`
)

var CameraTableColumns = []string{
	CameraTableIDColumn,
	CameraTableCreatedAtColumn,
	CameraTableUpdatedAtColumn,
	CameraTableDeletedAtColumn,
	CameraTableNameColumn,
	CameraTableStreamURLColumn,
	CameraTableLastSeenColumn,
}

var CameraTableColumnsWithTypeCasts = []string{
	CameraTableIDColumnWithTypeCast,
	CameraTableCreatedAtColumnWithTypeCast,
	CameraTableUpdatedAtColumnWithTypeCast,
	CameraTableDeletedAtColumnWithTypeCast,
	CameraTableNameColumnWithTypeCast,
	CameraTableStreamURLColumnWithTypeCast,
	CameraTableLastSeenColumnWithTypeCast,
}

var CameraIntrospectedTable *introspect.Table

var CameraTableColumnLookup map[string]*introspect.Column

var (
	CameraTablePrimaryKeyColumn = CameraTableIDColumn
)

func init() {
	CameraIntrospectedTable = tableByName[CameraTable]

	/* only needed during templating */
	if CameraIntrospectedTable == nil {
		CameraIntrospectedTable = &introspect.Table{}
	}

	CameraTableColumnLookup = CameraIntrospectedTable.ColumnByName
}

type CameraOnePathParams struct {
	PrimaryKey uuid.UUID `json:"primaryKey"`
}

type CameraLoadQueryParams struct {
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

func (m *Camera) GetPrimaryKeyColumn() string {
	return CameraTablePrimaryKeyColumn
}

func (m *Camera) GetPrimaryKeyValue() any {
	return m.ID
}

func (m *Camera) FromItem(item map[string]any) error {
	if item == nil {
		return fmt.Errorf(
			"item unexpectedly nil during CameraFromItem",
		)
	}

	if len(item) == 0 {
		return fmt.Errorf(
			"item unexpectedly empty during CameraFromItem",
		)
	}

	wrapError := func(k string, v any, err error) error {
		return fmt.Errorf("%v: %#+v; error; %v", k, v, err)
	}

	for k, v := range item {
		_, ok := CameraTableColumnLookup[k]
		if !ok {
			return fmt.Errorf(
				"item contained unexpected key %#+v during CameraFromItem; item: %#+v",
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

		case "stream_url":
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
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uustream_url.UUID", temp1))
				}
			}

			m.StreamURL = temp2

		case "last_seen":
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
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uulast_seen.UUID", temp1))
				}
			}

			m.LastSeen = &temp2

		}
	}

	return nil
}

func (m *Camera) Reload(ctx context.Context, tx pgx.Tx, includeDeleteds ...bool) error {
	extraWhere := ""
	if len(includeDeleteds) > 0 && includeDeleteds[0] {
		if slices.Contains(CameraTableColumns, "deleted_at") {
			extraWhere = "\n    AND (deleted_at IS null OR deleted_at IS NOT null)"
		}
	}

	ctx, cleanup := query.WithQueryID(ctx)
	defer cleanup()

	t, err := SelectCamera(
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
	m.Name = t.Name
	m.StreamURL = t.StreamURL
	m.LastSeen = t.LastSeen
	m.ReferencedByVideoCameraIDObjects = t.ReferencedByVideoCameraIDObjects
	m.ReferencedByDetectionCameraIDObjects = t.ReferencedByDetectionCameraIDObjects

	return nil
}

func (m *Camera) Insert(ctx context.Context, tx pgx.Tx, setPrimaryKey bool, setZeroValues bool, forceSetValuesForFields ...string) error {
	columns := make([]string, 0)
	values := make([]any, 0)

	if setPrimaryKey && (setZeroValues || !types.IsZeroUUID(m.ID) || slices.Contains(forceSetValuesForFields, CameraTableIDColumn) || isRequired(CameraTableColumnLookup, CameraTableIDColumn)) {
		columns = append(columns, CameraTableIDColumn)

		v, err := types.FormatUUID(m.ID)
		if err != nil {
			return fmt.Errorf("failed to handle m.ID: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.CreatedAt) || slices.Contains(forceSetValuesForFields, CameraTableCreatedAtColumn) || isRequired(CameraTableColumnLookup, CameraTableCreatedAtColumn) {
		columns = append(columns, CameraTableCreatedAtColumn)

		v, err := types.FormatTime(m.CreatedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.CreatedAt: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.UpdatedAt) || slices.Contains(forceSetValuesForFields, CameraTableUpdatedAtColumn) || isRequired(CameraTableColumnLookup, CameraTableUpdatedAtColumn) {
		columns = append(columns, CameraTableUpdatedAtColumn)

		v, err := types.FormatTime(m.UpdatedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.UpdatedAt: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.DeletedAt) || slices.Contains(forceSetValuesForFields, CameraTableDeletedAtColumn) || isRequired(CameraTableColumnLookup, CameraTableDeletedAtColumn) {
		columns = append(columns, CameraTableDeletedAtColumn)

		v, err := types.FormatTime(m.DeletedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.DeletedAt: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroString(m.Name) || slices.Contains(forceSetValuesForFields, CameraTableNameColumn) || isRequired(CameraTableColumnLookup, CameraTableNameColumn) {
		columns = append(columns, CameraTableNameColumn)

		v, err := types.FormatString(m.Name)
		if err != nil {
			return fmt.Errorf("failed to handle m.Name: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroString(m.StreamURL) || slices.Contains(forceSetValuesForFields, CameraTableStreamURLColumn) || isRequired(CameraTableColumnLookup, CameraTableStreamURLColumn) {
		columns = append(columns, CameraTableStreamURLColumn)

		v, err := types.FormatString(m.StreamURL)
		if err != nil {
			return fmt.Errorf("failed to handle m.StreamURL: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.LastSeen) || slices.Contains(forceSetValuesForFields, CameraTableLastSeenColumn) || isRequired(CameraTableColumnLookup, CameraTableLastSeenColumn) {
		columns = append(columns, CameraTableLastSeenColumn)

		v, err := types.FormatTime(m.LastSeen)
		if err != nil {
			return fmt.Errorf("failed to handle m.LastSeen: %v", err)
		}

		values = append(values, v)
	}

	ctx, cleanup := query.WithQueryID(ctx)
	defer cleanup()

	item, err := query.Insert(
		ctx,
		tx,
		CameraTable,
		columns,
		nil,
		false,
		false,
		CameraTableColumns,
		values...,
	)
	if err != nil {
		return fmt.Errorf("failed to insert %#+v; %v", m, err)
	}
	v := (*item)[CameraTableIDColumn]

	if v == nil {
		return fmt.Errorf("failed to find %v in %#+v", CameraTableIDColumn, item)
	}

	wrapError := func(err error) error {
		return fmt.Errorf(
			"failed to treat %v: %#+v as uuid.UUID: %v",
			CameraTableIDColumn,
			(*item)[CameraTableIDColumn],
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

func (m *Camera) Update(ctx context.Context, tx pgx.Tx, setZeroValues bool, forceSetValuesForFields ...string) error {
	columns := make([]string, 0)
	values := make([]any, 0)

	if setZeroValues || !types.IsZeroTime(m.CreatedAt) || slices.Contains(forceSetValuesForFields, CameraTableCreatedAtColumn) {
		columns = append(columns, CameraTableCreatedAtColumn)

		v, err := types.FormatTime(m.CreatedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.CreatedAt: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.UpdatedAt) || slices.Contains(forceSetValuesForFields, CameraTableUpdatedAtColumn) {
		columns = append(columns, CameraTableUpdatedAtColumn)

		v, err := types.FormatTime(m.UpdatedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.UpdatedAt: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.DeletedAt) || slices.Contains(forceSetValuesForFields, CameraTableDeletedAtColumn) {
		columns = append(columns, CameraTableDeletedAtColumn)

		v, err := types.FormatTime(m.DeletedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.DeletedAt: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroString(m.Name) || slices.Contains(forceSetValuesForFields, CameraTableNameColumn) {
		columns = append(columns, CameraTableNameColumn)

		v, err := types.FormatString(m.Name)
		if err != nil {
			return fmt.Errorf("failed to handle m.Name: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroString(m.StreamURL) || slices.Contains(forceSetValuesForFields, CameraTableStreamURLColumn) {
		columns = append(columns, CameraTableStreamURLColumn)

		v, err := types.FormatString(m.StreamURL)
		if err != nil {
			return fmt.Errorf("failed to handle m.StreamURL: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.LastSeen) || slices.Contains(forceSetValuesForFields, CameraTableLastSeenColumn) {
		columns = append(columns, CameraTableLastSeenColumn)

		v, err := types.FormatTime(m.LastSeen)
		if err != nil {
			return fmt.Errorf("failed to handle m.LastSeen: %v", err)
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
		CameraTable,
		columns,
		fmt.Sprintf("%v = $$??", CameraTableIDColumn),
		CameraTableColumns,
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

func (m *Camera) Delete(ctx context.Context, tx pgx.Tx, hardDeletes ...bool) error {
	hardDelete := false
	if len(hardDeletes) > 0 {
		hardDelete = hardDeletes[0]
	}

	if !hardDelete && slices.Contains(CameraTableColumns, "deleted_at") {
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
		CameraTable,
		fmt.Sprintf("%v = $$??", CameraTableIDColumn),
		values...,
	)
	if err != nil {
		return fmt.Errorf("failed to delete %#+v; %v", m, err)
	}

	_ = m.Reload(ctx, tx, true)

	return nil
}

func (m *Camera) LockTable(ctx context.Context, tx pgx.Tx, noWait bool) error {
	return query.LockTable(ctx, tx, CameraTable, noWait)
}

func SelectCameras(ctx context.Context, tx pgx.Tx, where string, orderBy *string, limit *int, offset *int, values ...any) ([]*Camera, error) {
	if slices.Contains(CameraTableColumns, "deleted_at") {
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
		CameraTableColumnsWithTypeCasts,
		CameraTable,
		where,
		orderBy,
		limit,
		offset,
		values...,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to call SelectCameras; err: %v", err)
	}

	objects := make([]*Camera, 0)

	for _, item := range *items {
		object := &Camera{}

		err = object.FromItem(item)
		if err != nil {
			return nil, err
		}

		thatCtx := ctx

		thatCtx, ok1 := query.HandleQueryPathGraphCycles(ctx, fmt.Sprintf("%s{%v}", CameraTable, object.GetPrimaryKeyValue()))
		thatCtx, ok2 := query.HandleQueryPathGraphCycles(thatCtx, fmt.Sprintf("__ReferencedBy__%s{%v}", CameraTable, object.GetPrimaryKeyValue()))
		if !(ok1 && ok2) {
			continue
		}

		_ = thatCtx

		err = func() error {
			thisCtx := thatCtx
			thisCtx, ok1 := query.HandleQueryPathGraphCycles(thisCtx, fmt.Sprintf("%s{%v}", CameraTable, object.GetPrimaryKeyValue()))
			thisCtx, ok2 := query.HandleQueryPathGraphCycles(thisCtx, fmt.Sprintf("__ReferencedBy__%s{%v}", CameraTable, object.GetPrimaryKeyValue()))

			if ok1 && ok2 {
				object.ReferencedByVideoCameraIDObjects, err = SelectVideos(
					thisCtx,
					tx,
					fmt.Sprintf("%v = $1", VideoTableCameraIDColumn),
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
			}

			return nil
		}()
		if err != nil {
			return nil, err
		}

		err = func() error {
			thisCtx := thatCtx
			thisCtx, ok1 := query.HandleQueryPathGraphCycles(thisCtx, fmt.Sprintf("%s{%v}", CameraTable, object.GetPrimaryKeyValue()))
			thisCtx, ok2 := query.HandleQueryPathGraphCycles(thisCtx, fmt.Sprintf("__ReferencedBy__%s{%v}", CameraTable, object.GetPrimaryKeyValue()))

			if ok1 && ok2 {
				object.ReferencedByDetectionCameraIDObjects, err = SelectDetections(
					thisCtx,
					tx,
					fmt.Sprintf("%v = $1", DetectionTableCameraIDColumn),
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

func SelectCamera(ctx context.Context, tx pgx.Tx, where string, values ...any) (*Camera, error) {
	ctx, cleanup := query.WithQueryID(ctx)
	defer cleanup()

	objects, err := SelectCameras(
		ctx,
		tx,
		where,
		nil,
		helpers.Ptr(2),
		helpers.Ptr(0),
		values...,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to call SelectCamera; err: %v", err)
	}

	if len(objects) > 1 {
		return nil, fmt.Errorf("attempt to call SelectCamera returned more than 1 row")
	}

	if len(objects) < 1 {
		return nil, sql.ErrNoRows
	}

	object := objects[0]

	return object, nil
}

func handleGetCameras(arguments *server.SelectManyArguments, db *pgxpool.Pool) ([]*Camera, error) {
	tx, err := db.Begin(arguments.Ctx)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = tx.Rollback(arguments.Ctx)
	}()

	objects, err := SelectCameras(arguments.Ctx, tx, arguments.Where, arguments.OrderBy, arguments.Limit, arguments.Offset, arguments.Values...)
	if err != nil {
		return nil, err
	}

	err = tx.Commit(arguments.Ctx)
	if err != nil {
		return nil, err
	}

	return objects, nil
}

func handleGetCamera(arguments *server.SelectOneArguments, db *pgxpool.Pool, primaryKey uuid.UUID) ([]*Camera, error) {
	tx, err := db.Begin(arguments.Ctx)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = tx.Rollback(arguments.Ctx)
	}()

	object, err := SelectCamera(arguments.Ctx, tx, arguments.Where, arguments.Values...)
	if err != nil {
		return nil, err
	}

	err = tx.Commit(arguments.Ctx)
	if err != nil {
		return nil, err
	}

	return []*Camera{object}, nil
}

func handlePostCameras(arguments *server.LoadArguments, db *pgxpool.Pool, waitForChange server.WaitForChange, objects []*Camera, forceSetValuesForFieldsByObjectIndex [][]string) ([]*Camera, error) {
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
		_, err = waitForChange(arguments.Ctx, []stream.Action{stream.INSERT}, CameraTable, xid)
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

func handlePutCamera(arguments *server.LoadArguments, db *pgxpool.Pool, waitForChange server.WaitForChange, object *Camera) ([]*Camera, error) {
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
		_, err = waitForChange(arguments.Ctx, []stream.Action{stream.UPDATE, stream.SOFT_DELETE, stream.SOFT_RESTORE, stream.SOFT_UPDATE}, CameraTable, xid)
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

	return []*Camera{object}, nil
}

func handlePatchCamera(arguments *server.LoadArguments, db *pgxpool.Pool, waitForChange server.WaitForChange, object *Camera, forceSetValuesForFields []string) ([]*Camera, error) {
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
		_, err = waitForChange(arguments.Ctx, []stream.Action{stream.UPDATE, stream.SOFT_DELETE, stream.SOFT_RESTORE, stream.SOFT_UPDATE}, CameraTable, xid)
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

	return []*Camera{object}, nil
}

func handleDeleteCamera(arguments *server.LoadArguments, db *pgxpool.Pool, waitForChange server.WaitForChange, object *Camera) error {
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
		_, err = waitForChange(arguments.Ctx, []stream.Action{stream.DELETE, stream.SOFT_DELETE}, CameraTable, xid)
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

func GetCameraRouter(db *pgxpool.Pool, redisPool *redis.Pool, httpMiddlewares []server.HTTPMiddleware, objectMiddlewares []server.ObjectMiddleware, waitForChange server.WaitForChange) chi.Router {
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
		) (*helpers.TypedResponse[Camera], error) {
			redisConn := redisPool.Get()
			defer func() {
				_ = redisConn.Close()
			}()

			arguments, err := server.GetSelectManyArguments(ctx, queryParams, CameraIntrospectedTable, nil, nil)
			if err != nil {
				return nil, err
			}

			cachedObjectsAsJSON, cacheHit, err := helpers.GetCachedObjectsAsJSON(arguments.RequestHash, redisConn)
			if err != nil {
				return nil, err
			}

			if cacheHit {
				var cachedObjects []*Camera
				err = json.Unmarshal(cachedObjectsAsJSON, &cachedObjects)
				if err != nil {
					return nil, err
				}

				return &helpers.TypedResponse[Camera]{
					Status:  http.StatusOK,
					Success: true,
					Error:   nil,
					Objects: cachedObjects,
				}, nil
			}

			objects, err := handleGetCameras(arguments, db)
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

			return &helpers.TypedResponse[Camera]{
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
			pathParams CameraOnePathParams,
			queryParams CameraLoadQueryParams,
			req server.EmptyRequest,
			rawReq any,
		) (*helpers.TypedResponse[Camera], error) {
			redisConn := redisPool.Get()
			defer func() {
				_ = redisConn.Close()
			}()

			arguments, err := server.GetSelectOneArguments(ctx, queryParams.Depth, CameraIntrospectedTable, pathParams.PrimaryKey, nil, nil)
			if err != nil {
				return nil, err
			}

			cachedObjectsAsJSON, cacheHit, err := helpers.GetCachedObjectsAsJSON(arguments.RequestHash, redisConn)
			if err != nil {
				return nil, err
			}

			if cacheHit {
				var cachedObjects []*Camera
				err = json.Unmarshal(cachedObjectsAsJSON, &cachedObjects)
				if err != nil {
					return nil, err
				}

				return &helpers.TypedResponse[Camera]{
					Status:  http.StatusOK,
					Success: true,
					Error:   nil,
					Objects: cachedObjects,
				}, nil
			}

			objects, err := handleGetCamera(arguments, db, pathParams.PrimaryKey)
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

			return &helpers.TypedResponse[Camera]{
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
			queryParams CameraLoadQueryParams,
			req []*Camera,
			rawReq any,
		) (*helpers.TypedResponse[Camera], error) {
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
					if !slices.Contains(CameraTableColumns, possibleField) {
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

			objects, err := handlePostCameras(arguments, db, waitForChange, req, forceSetValuesForFieldsByObjectIndex)
			if err != nil {
				return nil, err
			}

			return &helpers.TypedResponse[Camera]{
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
			pathParams CameraOnePathParams,
			queryParams CameraLoadQueryParams,
			req Camera,
			rawReq any,
		) (*helpers.TypedResponse[Camera], error) {
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

			objects, err := handlePutCamera(arguments, db, waitForChange, object)
			if err != nil {
				return nil, err
			}

			return &helpers.TypedResponse[Camera]{
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
			pathParams CameraOnePathParams,
			queryParams CameraLoadQueryParams,
			req Camera,
			rawReq any,
		) (*helpers.TypedResponse[Camera], error) {
			item, ok := rawReq.(map[string]any)
			if !ok {
				return nil, fmt.Errorf("failed to cast %#+v to map[string]any", item)
			}

			forceSetValuesForFields := make([]string, 0)
			for _, possibleField := range maps.Keys(item) {
				if !slices.Contains(CameraTableColumns, possibleField) {
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

			objects, err := handlePatchCamera(arguments, db, waitForChange, object, forceSetValuesForFields)
			if err != nil {
				return nil, err
			}

			return &helpers.TypedResponse[Camera]{
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
			pathParams CameraOnePathParams,
			queryParams CameraLoadQueryParams,
			req server.EmptyRequest,
			rawReq any,
		) (*server.EmptyResponse, error) {
			arguments := &server.LoadArguments{
				Ctx: ctx,
			}

			object := &Camera{}
			object.ID = pathParams.PrimaryKey

			err := handleDeleteCamera(arguments, db, waitForChange, object)
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

func NewCameraFromItem(item map[string]any) (any, error) {
	object := &Camera{}

	err := object.FromItem(item)
	if err != nil {
		return nil, err
	}

	return object, nil
}

func init() {
	register(
		CameraTable,
		Camera{},
		NewCameraFromItem,
		"/cameras",
		GetCameraRouter,
	)
}
