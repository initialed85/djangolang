package model_generated

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

type LogicalThingClaim struct {
	ID                    uuid.UUID     `json:"id"`
	ClaimedFor            string        `json:"claimed_for"`
	ClaimedUntil          *time.Time    `json:"claimed_until"`
	ClaimedBy             *uuid.UUID    `json:"claimed_by"`
	LogicalThingsID       uuid.UUID     `json:"logical_things_id"`
	LogicalThingsIDObject *LogicalThing `json:"logical_things_id_object"`
}

var LogicalThingClaimTable = "logical_thing_claims"

var LogicalThingClaimTableWithSchema = fmt.Sprintf("%s.%s", schema, LogicalThingClaimTable)

var LogicalThingClaimTableNamespaceID int32 = 1337 + 4

var (
	LogicalThingClaimTableIDColumn              = "id"
	LogicalThingClaimTableClaimedForColumn      = "claimed_for"
	LogicalThingClaimTableClaimedUntilColumn    = "claimed_until"
	LogicalThingClaimTableClaimedByColumn       = "claimed_by"
	LogicalThingClaimTableLogicalThingsIDColumn = "logical_things_id"
)

var (
	LogicalThingClaimTableIDColumnWithTypeCast              = `"id" AS id`
	LogicalThingClaimTableClaimedForColumnWithTypeCast      = `"claimed_for" AS claimed_for`
	LogicalThingClaimTableClaimedUntilColumnWithTypeCast    = `"claimed_until" AS claimed_until`
	LogicalThingClaimTableClaimedByColumnWithTypeCast       = `"claimed_by" AS claimed_by`
	LogicalThingClaimTableLogicalThingsIDColumnWithTypeCast = `"logical_things_id" AS logical_things_id`
)

var LogicalThingClaimTableColumns = []string{
	LogicalThingClaimTableIDColumn,
	LogicalThingClaimTableClaimedForColumn,
	LogicalThingClaimTableClaimedUntilColumn,
	LogicalThingClaimTableClaimedByColumn,
	LogicalThingClaimTableLogicalThingsIDColumn,
}

var LogicalThingClaimTableColumnsWithTypeCasts = []string{
	LogicalThingClaimTableIDColumnWithTypeCast,
	LogicalThingClaimTableClaimedForColumnWithTypeCast,
	LogicalThingClaimTableClaimedUntilColumnWithTypeCast,
	LogicalThingClaimTableClaimedByColumnWithTypeCast,
	LogicalThingClaimTableLogicalThingsIDColumnWithTypeCast,
}

var LogicalThingClaimIntrospectedTable *introspect.Table

var LogicalThingClaimTableColumnLookup map[string]*introspect.Column

var (
	LogicalThingClaimTablePrimaryKeyColumn = LogicalThingClaimTableIDColumn
)

func init() {
	LogicalThingClaimIntrospectedTable = tableByName[LogicalThingClaimTable]

	/* only needed during templating */
	if LogicalThingClaimIntrospectedTable == nil {
		LogicalThingClaimIntrospectedTable = &introspect.Table{}
	}

	LogicalThingClaimTableColumnLookup = LogicalThingClaimIntrospectedTable.ColumnByName
}

type LogicalThingClaimOnePathParams struct {
	PrimaryKey uuid.UUID `json:"primaryKey"`
}

type LogicalThingClaimLoadQueryParams struct {
	Depth *int `json:"depth"`
}

type LogicalThingClaimClaimRequest struct {
	For            string    `json:"for"`
	Until          time.Time `json:"until"`
	By             uuid.UUID `json:"by"`
	TimeoutSeconds float64   `json:"timeout_seconds"`
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

func (m *LogicalThingClaim) GetPrimaryKeyColumn() string {
	return LogicalThingClaimTablePrimaryKeyColumn
}

func (m *LogicalThingClaim) GetPrimaryKeyValue() any {
	return m.ID
}

func (m *LogicalThingClaim) FromItem(item map[string]any) error {
	if item == nil {
		return fmt.Errorf(
			"item unexpectedly nil during LogicalThingClaimFromItem",
		)
	}

	if len(item) == 0 {
		return fmt.Errorf(
			"item unexpectedly empty during LogicalThingClaimFromItem",
		)
	}

	wrapError := func(k string, v any, err error) error {
		return fmt.Errorf("%v: %#+v; error; %v", k, v, err)
	}

	for k, v := range item {
		_, ok := LogicalThingClaimTableColumnLookup[k]
		if !ok {
			return fmt.Errorf(
				"item contained unexpected key %#+v during LogicalThingClaimFromItem; item: %#+v",
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

		case "claimed_for":
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
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uuclaimed_for.UUID", temp1))
				}
			}

			m.ClaimedFor = temp2

		case "claimed_until":
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
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uuclaimed_until.UUID", temp1))
				}
			}

			m.ClaimedUntil = &temp2

		case "claimed_by":
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
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uuclaimed_by.UUID", temp1))
				}
			}

			m.ClaimedBy = &temp2

		case "logical_things_id":
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
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uulogical_things_id.UUID", temp1))
				}
			}

			m.LogicalThingsID = temp2

		}
	}

	return nil
}

func (m *LogicalThingClaim) Reload(ctx context.Context, tx pgx.Tx, includeDeleteds ...bool) error {
	extraWhere := ""
	if len(includeDeleteds) > 0 && includeDeleteds[0] {
		if slices.Contains(LogicalThingClaimTableColumns, "deleted_at") {
			extraWhere = "\n    AND (deleted_at IS null OR deleted_at IS NOT null)"
		}
	}

	ctx, cleanup := query.WithQueryID(ctx)
	defer cleanup()

	ctx = query.WithMaxDepth(ctx, nil)

	o, _, _, _, _, err := SelectLogicalThingClaim(
		ctx,
		tx,
		fmt.Sprintf("%v = $1%v", m.GetPrimaryKeyColumn(), extraWhere),
		m.GetPrimaryKeyValue(),
	)
	if err != nil {
		return err
	}

	m.ID = o.ID
	m.ClaimedFor = o.ClaimedFor
	m.ClaimedUntil = o.ClaimedUntil
	m.ClaimedBy = o.ClaimedBy
	m.LogicalThingsID = o.LogicalThingsID
	m.LogicalThingsIDObject = o.LogicalThingsIDObject

	return nil
}

func (m *LogicalThingClaim) Insert(ctx context.Context, tx pgx.Tx, setPrimaryKey bool, setZeroValues bool, forceSetValuesForFields ...string) error {
	columns := make([]string, 0)
	values := make([]any, 0)

	if setPrimaryKey && (setZeroValues || !types.IsZeroUUID(m.ID) || slices.Contains(forceSetValuesForFields, LogicalThingClaimTableIDColumn) || isRequired(LogicalThingClaimTableColumnLookup, LogicalThingClaimTableIDColumn)) {
		columns = append(columns, LogicalThingClaimTableIDColumn)

		v, err := types.FormatUUID(m.ID)
		if err != nil {
			return fmt.Errorf("failed to handle m.ID; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroString(m.ClaimedFor) || slices.Contains(forceSetValuesForFields, LogicalThingClaimTableClaimedForColumn) || isRequired(LogicalThingClaimTableColumnLookup, LogicalThingClaimTableClaimedForColumn) {
		columns = append(columns, LogicalThingClaimTableClaimedForColumn)

		v, err := types.FormatString(m.ClaimedFor)
		if err != nil {
			return fmt.Errorf("failed to handle m.ClaimedFor; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.ClaimedUntil) || slices.Contains(forceSetValuesForFields, LogicalThingClaimTableClaimedUntilColumn) || isRequired(LogicalThingClaimTableColumnLookup, LogicalThingClaimTableClaimedUntilColumn) {
		columns = append(columns, LogicalThingClaimTableClaimedUntilColumn)

		v, err := types.FormatTime(m.ClaimedUntil)
		if err != nil {
			return fmt.Errorf("failed to handle m.ClaimedUntil; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroUUID(m.ClaimedBy) || slices.Contains(forceSetValuesForFields, LogicalThingClaimTableClaimedByColumn) || isRequired(LogicalThingClaimTableColumnLookup, LogicalThingClaimTableClaimedByColumn) {
		columns = append(columns, LogicalThingClaimTableClaimedByColumn)

		v, err := types.FormatUUID(m.ClaimedBy)
		if err != nil {
			return fmt.Errorf("failed to handle m.ClaimedBy; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroUUID(m.LogicalThingsID) || slices.Contains(forceSetValuesForFields, LogicalThingClaimTableLogicalThingsIDColumn) || isRequired(LogicalThingClaimTableColumnLookup, LogicalThingClaimTableLogicalThingsIDColumn) {
		columns = append(columns, LogicalThingClaimTableLogicalThingsIDColumn)

		v, err := types.FormatUUID(m.LogicalThingsID)
		if err != nil {
			return fmt.Errorf("failed to handle m.LogicalThingsID; %v", err)
		}

		values = append(values, v)
	}

	ctx, cleanup := query.WithQueryID(ctx)
	defer cleanup()

	ctx = query.WithMaxDepth(ctx, nil)

	item, err := query.Insert(
		ctx,
		tx,
		LogicalThingClaimTableWithSchema,
		columns,
		nil,
		false,
		false,
		LogicalThingClaimTableColumns,
		values...,
	)
	if err != nil {
		return fmt.Errorf("failed to insert %#+v; %v", m, err)
	}
	v := (*item)[LogicalThingClaimTableIDColumn]

	if v == nil {
		return fmt.Errorf("failed to find %v in %#+v", LogicalThingClaimTableIDColumn, item)
	}

	wrapError := func(err error) error {
		return fmt.Errorf(
			"failed to treat %v: %#+v as uuid.UUID: %v",
			LogicalThingClaimTableIDColumn,
			(*item)[LogicalThingClaimTableIDColumn],
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

func (m *LogicalThingClaim) Update(ctx context.Context, tx pgx.Tx, setZeroValues bool, forceSetValuesForFields ...string) error {
	columns := make([]string, 0)
	values := make([]any, 0)

	if setZeroValues || !types.IsZeroString(m.ClaimedFor) || slices.Contains(forceSetValuesForFields, LogicalThingClaimTableClaimedForColumn) {
		columns = append(columns, LogicalThingClaimTableClaimedForColumn)

		v, err := types.FormatString(m.ClaimedFor)
		if err != nil {
			return fmt.Errorf("failed to handle m.ClaimedFor; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.ClaimedUntil) || slices.Contains(forceSetValuesForFields, LogicalThingClaimTableClaimedUntilColumn) {
		columns = append(columns, LogicalThingClaimTableClaimedUntilColumn)

		v, err := types.FormatTime(m.ClaimedUntil)
		if err != nil {
			return fmt.Errorf("failed to handle m.ClaimedUntil; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroUUID(m.ClaimedBy) || slices.Contains(forceSetValuesForFields, LogicalThingClaimTableClaimedByColumn) {
		columns = append(columns, LogicalThingClaimTableClaimedByColumn)

		v, err := types.FormatUUID(m.ClaimedBy)
		if err != nil {
			return fmt.Errorf("failed to handle m.ClaimedBy; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroUUID(m.LogicalThingsID) || slices.Contains(forceSetValuesForFields, LogicalThingClaimTableLogicalThingsIDColumn) {
		columns = append(columns, LogicalThingClaimTableLogicalThingsIDColumn)

		v, err := types.FormatUUID(m.LogicalThingsID)
		if err != nil {
			return fmt.Errorf("failed to handle m.LogicalThingsID; %v", err)
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
		LogicalThingClaimTableWithSchema,
		columns,
		fmt.Sprintf("%v = $$??", LogicalThingClaimTableIDColumn),
		LogicalThingClaimTableColumns,
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

func (m *LogicalThingClaim) Delete(ctx context.Context, tx pgx.Tx, hardDeletes ...bool) error {
	/* soft-delete not applicable */

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
		LogicalThingClaimTableWithSchema,
		fmt.Sprintf("%v = $$??", LogicalThingClaimTableIDColumn),
		values...,
	)
	if err != nil {
		return fmt.Errorf("failed to delete %#+v; %v", m, err)
	}

	_ = m.Reload(ctx, tx, true)

	return nil
}

func (m *LogicalThingClaim) LockTable(ctx context.Context, tx pgx.Tx, timeouts ...time.Duration) error {
	return query.LockTable(ctx, tx, LogicalThingClaimTableWithSchema, timeouts...)
}

func (m *LogicalThingClaim) LockTableWithRetries(ctx context.Context, tx pgx.Tx, overallTimeout time.Duration, individualAttempttimeout time.Duration) error {
	return query.LockTableWithRetries(ctx, tx, LogicalThingClaimTableWithSchema, overallTimeout, individualAttempttimeout)
}

func (m *LogicalThingClaim) AdvisoryLock(ctx context.Context, tx pgx.Tx, key int32, timeouts ...time.Duration) error {
	return query.AdvisoryLock(ctx, tx, LogicalThingClaimTableNamespaceID, key, timeouts...)
}

func (m *LogicalThingClaim) AdvisoryLockWithRetries(ctx context.Context, tx pgx.Tx, key int32, overallTimeout time.Duration, individualAttempttimeout time.Duration) error {
	return query.AdvisoryLockWithRetries(ctx, tx, LogicalThingClaimTableNamespaceID, key, overallTimeout, individualAttempttimeout)
}

func (m *LogicalThingClaim) Claim(ctx context.Context, tx pgx.Tx, until time.Time, by uuid.UUID, timeout time.Duration) error {
	claimTableName := fmt.Sprintf("%s_claim", LogicalThingClaimTable)
	if !slices.Contains(maps.Keys(tableByName), claimTableName) {
		return fmt.Errorf("cannot invoke claim for LogicalThingClaim without \"%s\" table", claimTableName)
	}

	err := m.AdvisoryLockWithRetries(ctx, tx, math.MinInt32, timeout, time.Second*1)
	if err != nil {
		return fmt.Errorf("failed to claim (advisory lock): %s", err.Error())
	}

	x, _, _, _, _, err := SelectLogicalThingClaim(
		ctx,
		tx,
		fmt.Sprintf(
			"%s = $$?? AND (claimed_by = $$?? OR (claimed_until IS null OR claimed_until < now()))",
			LogicalThingClaimTablePrimaryKeyColumn,
		),
		m.GetPrimaryKeyValue(),
		by,
	)
	if err != nil {
		return fmt.Errorf("failed to claim (select): %s", err.Error())
	}

	_ = x

	err = m.Update(ctx, tx, false)
	if err != nil {
		return fmt.Errorf("failed to claim (update): %s", err.Error())
	}

	return nil
}

func SelectLogicalThingClaims(ctx context.Context, tx pgx.Tx, where string, orderBy *string, limit *int, offset *int, values ...any) ([]*LogicalThingClaim, int64, int64, int64, int64, error) {
	before := time.Now()

	if config.Debug() {
		log.Printf("entered SelectLogicalThingClaims")

		defer func() {
			log.Printf("exited SelectLogicalThingClaims in %s", time.Since(before))
		}()
	}
	if slices.Contains(LogicalThingClaimTableColumns, "deleted_at") {
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

	shouldLoad := query.ShouldLoad(ctx, LogicalThingClaimTable) || query.ShouldLoad(ctx, fmt.Sprintf("referenced_by_%s", LogicalThingClaimTable))

	var ok bool
	ctx, ok = query.HandleQueryPathGraphCycles(ctx, fmt.Sprintf("%s{%v}", LogicalThingClaimTable, nil), !isLoadQuery)
	if !ok && !shouldLoad {
		if config.Debug() {
			log.Printf("skipping SelectLogicalThingClaim early (query.ShouldLoad(): %v, query.HandleQueryPathGraphCycles(): %v)", shouldLoad, ok)
		}
		return []*LogicalThingClaim{}, 0, 0, 0, 0, nil
	}

	items, count, totalCount, page, totalPages, err := query.Select(
		ctx,
		tx,
		LogicalThingClaimTableColumnsWithTypeCasts,
		LogicalThingClaimTableWithSchema,
		where,
		orderBy,
		limit,
		offset,
		values...,
	)
	if err != nil {
		return nil, 0, 0, 0, 0, fmt.Errorf("failed to call SelectLogicalThingClaims; %v", err)
	}

	objects := make([]*LogicalThingClaim, 0)

	for _, item := range *items {
		object := &LogicalThingClaim{}

		err = object.FromItem(item)
		if err != nil {
			return nil, 0, 0, 0, 0, err
		}

		if !types.IsZeroUUID(object.LogicalThingsID) {
			ctx, ok := query.HandleQueryPathGraphCycles(ctx, fmt.Sprintf("%s{%v}", LogicalThingClaimTable, object.LogicalThingsID), true)
			shouldLoad := query.ShouldLoad(ctx, LogicalThingClaimTable)
			if ok || shouldLoad {
				thisBefore := time.Now()

				if config.Debug() {
					log.Printf("loading SelectLogicalThingClaims->SelectLogicalThing for object.LogicalThingsIDObject{%s: %v}", LogicalThingClaimTablePrimaryKeyColumn, object.LogicalThingsID)
				}

				object.LogicalThingsIDObject, _, _, _, _, err = SelectLogicalThing(
					ctx,
					tx,
					fmt.Sprintf("%v = $1", LogicalThingClaimTablePrimaryKeyColumn),
					object.LogicalThingsID,
				)
				if err != nil {
					if !errors.Is(err, sql.ErrNoRows) {
						return nil, 0, 0, 0, 0, err
					}
				}

				if config.Debug() {
					log.Printf("loaded SelectLogicalThingClaims->SelectLogicalThing for object.LogicalThingsIDObject in %s", time.Since(thisBefore))
				}
			}
		}

		objects = append(objects, object)
	}

	return objects, count, totalCount, page, totalPages, nil
}

func SelectLogicalThingClaim(ctx context.Context, tx pgx.Tx, where string, values ...any) (*LogicalThingClaim, int64, int64, int64, int64, error) {
	ctx, cleanup := query.WithQueryID(ctx)
	defer cleanup()

	ctx = query.WithMaxDepth(ctx, nil)

	objects, _, _, _, _, err := SelectLogicalThingClaims(
		ctx,
		tx,
		where,
		nil,
		helpers.Ptr(2),
		helpers.Ptr(0),
		values...,
	)
	if err != nil {
		return nil, 0, 0, 0, 0, fmt.Errorf("failed to call SelectLogicalThingClaim; %v", err)
	}

	if len(objects) > 1 {
		return nil, 0, 0, 0, 0, fmt.Errorf("attempt to call SelectLogicalThingClaim returned more than 1 row")
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

func ClaimLogicalThingClaim(ctx context.Context, tx pgx.Tx, claimedFor string, claimedUntil time.Time, claimedBy uuid.UUID, timeout time.Duration, wheres ...string) (*LogicalThingClaim, error) {
	m := &LogicalThingClaim{}

	err := m.AdvisoryLockWithRetries(ctx, tx, math.MinInt32, timeout, time.Second*1)
	if err != nil {
		return nil, fmt.Errorf("failed to claim: %s", err.Error())
	}

	extraWhere := ""
	if len(wheres) > 0 {
		extraWhere = fmt.Sprintf("AND %s", extraWhere)
	}

	itemsPtr, _, _, _, _, err := query.Select(
		ctx,
		tx,
		LogicalThingClaimTableColumns,
		fmt.Sprintf(
			"%s LEFT JOIN %s ON %s = %s AND %s = $$?? AND (%s = $$?? OR %s < now())",
			LogicalThingClaimTable,
			LogicalThingClaimTable,
			LogicalThingClaimTableLogicalThingsIDColumn,
			LogicalThingTableIDColumn,
			LogicalThingClaimTableClaimedForColumn,
			LogicalThingClaimTableClaimedByColumn,
			LogicalThingClaimTableClaimedUntilColumn,
		),
		fmt.Sprintf(
			"%s IS null OR %s < now()",
			LogicalThingClaimTableClaimedUntilColumn,
			LogicalThingClaimTableClaimedUntilColumn,
		),
		helpers.Ptr(fmt.Sprintf(
			"coalesce(%s, '0001-01-01'::timestamptz) ASC",
			LogicalThingClaimTableClaimedUntilColumn,
		)),
		helpers.Ptr(1),
		helpers.Ptr(0),
		claimedFor,
		claimedBy,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to claim: %s", err.Error())
	}

	if itemsPtr == nil {
		return nil, fmt.Errorf("failed to claim: %s", errors.New("itemsPtr unexpectedly nil"))
	}

	items := *itemsPtr

	if items != nil && len(items) == 0 {
		return nil, nil
	}

	claims, _, _, _, _, err := SelectLogicalThingClaims(
		ctx,
		tx,
		fmt.Sprintf(
			"%s = $$?? AND (%s IS null OR %s < now())",
			LogicalThingClaimTableClaimedForColumn,
			LogicalThingClaimTableClaimedByColumn,
			LogicalThingClaimTableClaimedUntilColumn,
		),
		helpers.Ptr(
			fmt.Sprintf(
				"%s ASC",
				LogicalThingClaimTableClaimedUntilColumn,
			),
		),
		helpers.Ptr(1),
		nil,
		claimedFor,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to claim: %s", err.Error())
	}

	if len(claims) > 0 {
		possibleM, _, _, _, _, err := SelectLogicalThing(
			ctx,
			tx,
			fmt.Sprintf(
				"(%s = $$??)%s",
				extraWhere,
			),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to claim: %s", err.Error())
		}

		m = possibleM
	} else {

	}

	ms, _, _, _, _, err := SelectLogicalThingClaims(
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

	return m, nil
}

func handleGetLogicalThingClaims(arguments *server.SelectManyArguments, db *pgxpool.Pool) ([]*LogicalThingClaim, int64, int64, int64, int64, error) {
	tx, err := db.Begin(arguments.Ctx)
	if err != nil {
		return nil, 0, 0, 0, 0, err
	}

	defer func() {
		_ = tx.Rollback(arguments.Ctx)
	}()

	objects, count, totalCount, page, totalPages, err := SelectLogicalThingClaims(arguments.Ctx, tx, arguments.Where, arguments.OrderBy, arguments.Limit, arguments.Offset, arguments.Values...)
	if err != nil {
		return nil, 0, 0, 0, 0, err
	}

	err = tx.Commit(arguments.Ctx)
	if err != nil {
		return nil, 0, 0, 0, 0, err
	}

	return objects, count, totalCount, page, totalPages, nil
}

func handleGetLogicalThingClaim(arguments *server.SelectOneArguments, db *pgxpool.Pool, primaryKey uuid.UUID) ([]*LogicalThingClaim, int64, int64, int64, int64, error) {
	tx, err := db.Begin(arguments.Ctx)
	if err != nil {
		return nil, 0, 0, 0, 0, err
	}

	defer func() {
		_ = tx.Rollback(arguments.Ctx)
	}()

	object, count, totalCount, page, totalPages, err := SelectLogicalThingClaim(arguments.Ctx, tx, arguments.Where, arguments.Values...)
	if err != nil {
		return nil, 0, 0, 0, 0, err
	}

	err = tx.Commit(arguments.Ctx)
	if err != nil {
		return nil, 0, 0, 0, 0, err
	}

	return []*LogicalThingClaim{object}, count, totalCount, page, totalPages, nil
}

func handlePostLogicalThingClaim(arguments *server.LoadArguments, db *pgxpool.Pool, waitForChange server.WaitForChange, objects []*LogicalThingClaim, forceSetValuesForFieldsByObjectIndex [][]string) ([]*LogicalThingClaim, int64, int64, int64, int64, error) {
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
		_, err := waitForChange(arguments.Ctx, []stream.Action{stream.INSERT}, LogicalThingClaimTable, xid)
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

func handlePutLogicalThingClaim(arguments *server.LoadArguments, db *pgxpool.Pool, waitForChange server.WaitForChange, object *LogicalThingClaim) ([]*LogicalThingClaim, int64, int64, int64, int64, error) {
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
		_, err := waitForChange(arguments.Ctx, []stream.Action{stream.UPDATE, stream.SOFT_DELETE, stream.SOFT_RESTORE, stream.SOFT_UPDATE}, LogicalThingClaimTable, xid)
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

	return []*LogicalThingClaim{object}, count, totalCount, page, totalPages, nil
}

func handlePatchLogicalThingClaim(arguments *server.LoadArguments, db *pgxpool.Pool, waitForChange server.WaitForChange, object *LogicalThingClaim, forceSetValuesForFields []string) ([]*LogicalThingClaim, int64, int64, int64, int64, error) {
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
		_, err := waitForChange(arguments.Ctx, []stream.Action{stream.UPDATE, stream.SOFT_DELETE, stream.SOFT_RESTORE, stream.SOFT_UPDATE}, LogicalThingClaimTable, xid)
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

	return []*LogicalThingClaim{object}, count, totalCount, page, totalPages, nil
}

func handleDeleteLogicalThingClaim(arguments *server.LoadArguments, db *pgxpool.Pool, waitForChange server.WaitForChange, object *LogicalThingClaim) error {
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
		_, err := waitForChange(arguments.Ctx, []stream.Action{stream.DELETE, stream.SOFT_DELETE}, LogicalThingClaimTable, xid)
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

func MutateRouterForLogicalThingClaim(r chi.Router, db *pgxpool.Pool, redisPool *redis.Pool, objectMiddlewares []server.ObjectMiddleware, waitForChange server.WaitForChange) {
	claimTableName := fmt.Sprintf("%s_claim", LogicalThingClaimTable)
	if slices.Contains(maps.Keys(tableByName), claimTableName) {
		func() {
			postHandlerForClaim, err := getHTTPHandler(
				http.MethodPost,
				"/claim-logical-thing-claim",
				http.StatusOK,
				func(
					ctx context.Context,
					pathParams server.EmptyPathParams,
					queryParams server.EmptyQueryParams,
					req LogicalThingClaimClaimRequest,
					rawReq any,
				) (server.Response[LogicalThingClaim], error) {
					tx, err := db.Begin(ctx)
					if err != nil {
						return server.Response[LogicalThingClaim]{}, err
					}

					defer func() {
						_ = tx.Rollback(ctx)
					}()

					object, err := ClaimLogicalThingClaim(ctx, tx, req.For, req.Until, req.By, time.Millisecond*time.Duration(req.TimeoutSeconds*1000))
					if err != nil {
						return server.Response[LogicalThingClaim]{}, err
					}

					count := int64(0)

					totalCount := int64(0)

					limit := int64(0)

					offset := int64(0)

					if object == nil {
						return server.Response[LogicalThingClaim]{
							Status:     http.StatusOK,
							Success:    true,
							Error:      nil,
							Objects:    []*LogicalThingClaim{},
							Count:      count,
							TotalCount: totalCount,
							Limit:      limit,
							Offset:     offset,
						}, nil
					}

					err = tx.Commit(ctx)
					if err != nil {
						return server.Response[LogicalThingClaim]{}, err
					}

					return server.Response[LogicalThingClaim]{
						Status:     http.StatusOK,
						Success:    true,
						Error:      nil,
						Objects:    []*LogicalThingClaim{object},
						Count:      count,
						TotalCount: totalCount,
						Limit:      limit,
						Offset:     offset,
					}, nil
				},
				LogicalThingClaim{},
				LogicalThingClaimIntrospectedTable,
			)
			if err != nil {
				panic(err)
			}
			r.Post(postHandlerForClaim.FullPath, postHandlerForClaim.ServeHTTP)

			postHandlerForClaimOne, err := getHTTPHandler(
				http.MethodPost,
				"/logical-thing-claims/{primaryKey}/claim",
				http.StatusOK,
				func(
					ctx context.Context,
					pathParams LogicalThingClaimOnePathParams,
					queryParams LogicalThingClaimLoadQueryParams,
					req LogicalThingClaimClaimRequest,
					rawReq any,
				) (server.Response[LogicalThingClaim], error) {
					before := time.Now()

					redisConn := redisPool.Get()
					defer func() {
						_ = redisConn.Close()
					}()

					arguments, err := server.GetSelectOneArguments(ctx, queryParams.Depth, LogicalThingClaimIntrospectedTable, pathParams.PrimaryKey, nil, nil)
					if err != nil {
						if config.Debug() {
							log.Printf("request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
						}

						return server.Response[LogicalThingClaim]{}, err
					}

					/* note: deliberately no attempt at a cache hit */

					var object *LogicalThingClaim
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

						object, count, totalCount, _, _, err = SelectLogicalThingClaim(arguments.Ctx, tx, arguments.Where, arguments.Values...)
						if err != nil {
							return fmt.Errorf("failed to select object to claim: %s", err.Error())
						}

						err = object.Claim(arguments.Ctx, tx, req.Until, req.By, time.Millisecond*time.Duration(req.TimeoutSeconds*1000))
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

						return server.Response[LogicalThingClaim]{}, err
					}

					limit := int64(0)

					offset := int64(0)

					response := server.Response[LogicalThingClaim]{
						Status:     http.StatusOK,
						Success:    true,
						Error:      nil,
						Objects:    []*LogicalThingClaim{object},
						Count:      count,
						TotalCount: totalCount,
						Limit:      limit,
						Offset:     offset,
					}

					return response, nil
				},
				LogicalThingClaim{},
				LogicalThingClaimIntrospectedTable,
			)
			if err != nil {
				panic(err)
			}
			r.Post(postHandlerForClaimOne.FullPath, postHandlerForClaimOne.ServeHTTP)
		}()
	}

	func() {
		getManyHandler, err := getHTTPHandler(
			http.MethodGet,
			"/logical-thing-claims",
			http.StatusOK,
			func(
				ctx context.Context,
				pathParams server.EmptyPathParams,
				queryParams map[string]any,
				req server.EmptyRequest,
				rawReq any,
			) (server.Response[LogicalThingClaim], error) {
				before := time.Now()

				redisConn := redisPool.Get()
				defer func() {
					_ = redisConn.Close()
				}()

				arguments, err := server.GetSelectManyArguments(ctx, queryParams, LogicalThingClaimIntrospectedTable, nil, nil)
				if err != nil {
					if config.Debug() {
						log.Printf("request cache not yet reached; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return server.Response[LogicalThingClaim]{}, err
				}

				cachedResponseAsJSON, cacheHit, err := server.GetCachedResponseAsJSON(arguments.RequestHash, redisConn)
				if err != nil {
					if config.Debug() {
						log.Printf("request cache failed; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return server.Response[LogicalThingClaim]{}, err
				}

				if cacheHit {
					var cachedResponse server.Response[LogicalThingClaim]

					/* TODO: it'd be nice to be able to avoid this (i.e. just pass straight through) */
					err = json.Unmarshal(cachedResponseAsJSON, &cachedResponse)
					if err != nil {
						if config.Debug() {
							log.Printf("request cache hit but failed unmarshal; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
						}

						return server.Response[LogicalThingClaim]{}, err
					}

					if config.Debug() {
						log.Printf("request cache hit; request succeeded in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return cachedResponse, nil
				}

				objects, count, totalCount, _, _, err := handleGetLogicalThingClaims(arguments, db)
				if err != nil {
					if config.Debug() {
						log.Printf("request cache missed; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return server.Response[LogicalThingClaim]{}, err
				}

				limit := int64(0)
				if arguments.Limit != nil {
					limit = int64(*arguments.Limit)
				}

				offset := int64(0)
				if arguments.Offset != nil {
					offset = int64(*arguments.Offset)
				}

				response := server.Response[LogicalThingClaim]{
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

					return server.Response[LogicalThingClaim]{}, err
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
			LogicalThingClaim{},
			LogicalThingClaimIntrospectedTable,
		)
		if err != nil {
			panic(err)
		}
		r.Get(getManyHandler.FullPath, getManyHandler.ServeHTTP)
	}()

	func() {
		getOneHandler, err := getHTTPHandler(
			http.MethodGet,
			"/logical-thing-claims/{primaryKey}",
			http.StatusOK,
			func(
				ctx context.Context,
				pathParams LogicalThingClaimOnePathParams,
				queryParams LogicalThingClaimLoadQueryParams,
				req server.EmptyRequest,
				rawReq any,
			) (server.Response[LogicalThingClaim], error) {
				before := time.Now()

				redisConn := redisPool.Get()
				defer func() {
					_ = redisConn.Close()
				}()

				arguments, err := server.GetSelectOneArguments(ctx, queryParams.Depth, LogicalThingClaimIntrospectedTable, pathParams.PrimaryKey, nil, nil)
				if err != nil {
					if config.Debug() {
						log.Printf("request cache not yet reached; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return server.Response[LogicalThingClaim]{}, err
				}

				cachedResponseAsJSON, cacheHit, err := server.GetCachedResponseAsJSON(arguments.RequestHash, redisConn)
				if err != nil {
					if config.Debug() {
						log.Printf("request cache failed; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return server.Response[LogicalThingClaim]{}, err
				}

				if cacheHit {
					var cachedResponse server.Response[LogicalThingClaim]

					/* TODO: it'd be nice to be able to avoid this (i.e. just pass straight through) */
					err = json.Unmarshal(cachedResponseAsJSON, &cachedResponse)
					if err != nil {
						if config.Debug() {
							log.Printf("request cache hit but failed unmarshal; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
						}

						return server.Response[LogicalThingClaim]{}, err
					}

					if config.Debug() {
						log.Printf("request cache hit; request succeeded in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return cachedResponse, nil
				}

				objects, count, totalCount, _, _, err := handleGetLogicalThingClaim(arguments, db, pathParams.PrimaryKey)
				if err != nil {
					if config.Debug() {
						log.Printf("request cache missed; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return server.Response[LogicalThingClaim]{}, err
				}

				limit := int64(0)

				offset := int64(0)

				response := server.Response[LogicalThingClaim]{
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

					return server.Response[LogicalThingClaim]{}, err
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
			LogicalThingClaim{},
			LogicalThingClaimIntrospectedTable,
		)
		if err != nil {
			panic(err)
		}
		r.Get(getOneHandler.FullPath, getOneHandler.ServeHTTP)
	}()

	func() {
		postHandler, err := getHTTPHandler(
			http.MethodPost,
			"/logical-thing-claims",
			http.StatusCreated,
			func(
				ctx context.Context,
				pathParams server.EmptyPathParams,
				queryParams LogicalThingClaimLoadQueryParams,
				req []*LogicalThingClaim,
				rawReq any,
			) (server.Response[LogicalThingClaim], error) {
				allRawItems, ok := rawReq.([]any)
				if !ok {
					return server.Response[LogicalThingClaim]{}, fmt.Errorf("failed to cast %#+v to []map[string]any", rawReq)
				}

				allItems := make([]map[string]any, 0)
				for _, rawItem := range allRawItems {
					item, ok := rawItem.(map[string]any)
					if !ok {
						return server.Response[LogicalThingClaim]{}, fmt.Errorf("failed to cast %#+v to map[string]any", rawItem)
					}

					allItems = append(allItems, item)
				}

				forceSetValuesForFieldsByObjectIndex := make([][]string, 0)
				for _, item := range allItems {
					forceSetValuesForFields := make([]string, 0)
					for _, possibleField := range maps.Keys(item) {
						if !slices.Contains(LogicalThingClaimTableColumns, possibleField) {
							continue
						}

						forceSetValuesForFields = append(forceSetValuesForFields, possibleField)
					}
					forceSetValuesForFieldsByObjectIndex = append(forceSetValuesForFieldsByObjectIndex, forceSetValuesForFields)
				}

				arguments, err := server.GetLoadArguments(ctx, queryParams.Depth)
				if err != nil {
					return server.Response[LogicalThingClaim]{}, err
				}

				objects, count, totalCount, _, _, err := handlePostLogicalThingClaim(arguments, db, waitForChange, req, forceSetValuesForFieldsByObjectIndex)
				if err != nil {
					return server.Response[LogicalThingClaim]{}, err
				}

				limit := int64(0)

				offset := int64(0)

				return server.Response[LogicalThingClaim]{
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
			LogicalThingClaim{},
			LogicalThingClaimIntrospectedTable,
		)
		if err != nil {
			panic(err)
		}
		r.Post(postHandler.FullPath, postHandler.ServeHTTP)
	}()

	func() {
		putHandler, err := getHTTPHandler(
			http.MethodPatch,
			"/logical-thing-claims/{primaryKey}",
			http.StatusOK,
			func(
				ctx context.Context,
				pathParams LogicalThingClaimOnePathParams,
				queryParams LogicalThingClaimLoadQueryParams,
				req LogicalThingClaim,
				rawReq any,
			) (server.Response[LogicalThingClaim], error) {
				item, ok := rawReq.(map[string]any)
				if !ok {
					return server.Response[LogicalThingClaim]{}, fmt.Errorf("failed to cast %#+v to map[string]any", item)
				}

				arguments, err := server.GetLoadArguments(ctx, queryParams.Depth)
				if err != nil {
					return server.Response[LogicalThingClaim]{}, err
				}

				object := &req
				object.ID = pathParams.PrimaryKey

				objects, count, totalCount, _, _, err := handlePutLogicalThingClaim(arguments, db, waitForChange, object)
				if err != nil {
					return server.Response[LogicalThingClaim]{}, err
				}

				limit := int64(0)

				offset := int64(0)

				return server.Response[LogicalThingClaim]{
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
			LogicalThingClaim{},
			LogicalThingClaimIntrospectedTable,
		)
		if err != nil {
			panic(err)
		}
		r.Put(putHandler.FullPath, putHandler.ServeHTTP)
	}()

	func() {
		patchHandler, err := getHTTPHandler(
			http.MethodPatch,
			"/logical-thing-claims/{primaryKey}",
			http.StatusOK,
			func(
				ctx context.Context,
				pathParams LogicalThingClaimOnePathParams,
				queryParams LogicalThingClaimLoadQueryParams,
				req LogicalThingClaim,
				rawReq any,
			) (server.Response[LogicalThingClaim], error) {
				item, ok := rawReq.(map[string]any)
				if !ok {
					return server.Response[LogicalThingClaim]{}, fmt.Errorf("failed to cast %#+v to map[string]any", item)
				}

				forceSetValuesForFields := make([]string, 0)
				for _, possibleField := range maps.Keys(item) {
					if !slices.Contains(LogicalThingClaimTableColumns, possibleField) {
						continue
					}

					forceSetValuesForFields = append(forceSetValuesForFields, possibleField)
				}

				arguments, err := server.GetLoadArguments(ctx, queryParams.Depth)
				if err != nil {
					return server.Response[LogicalThingClaim]{}, err
				}

				object := &req
				object.ID = pathParams.PrimaryKey

				objects, count, totalCount, _, _, err := handlePatchLogicalThingClaim(arguments, db, waitForChange, object, forceSetValuesForFields)
				if err != nil {
					return server.Response[LogicalThingClaim]{}, err
				}

				limit := int64(0)

				offset := int64(0)

				return server.Response[LogicalThingClaim]{
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
			LogicalThingClaim{},
			LogicalThingClaimIntrospectedTable,
		)
		if err != nil {
			panic(err)
		}
		r.Patch(patchHandler.FullPath, patchHandler.ServeHTTP)
	}()

	func() {
		deleteHandler, err := getHTTPHandler(
			http.MethodDelete,
			"/logical-thing-claims/{primaryKey}",
			http.StatusNoContent,
			func(
				ctx context.Context,
				pathParams LogicalThingClaimOnePathParams,
				queryParams LogicalThingClaimLoadQueryParams,
				req server.EmptyRequest,
				rawReq any,
			) (server.EmptyResponse, error) {
				arguments, err := server.GetLoadArguments(ctx, queryParams.Depth)
				if err != nil {
					return server.EmptyResponse{}, err
				}

				object := &LogicalThingClaim{}
				object.ID = pathParams.PrimaryKey

				err = handleDeleteLogicalThingClaim(arguments, db, waitForChange, object)
				if err != nil {
					return server.EmptyResponse{}, err
				}

				return server.EmptyResponse{}, nil
			},
			LogicalThingClaim{},
			LogicalThingClaimIntrospectedTable,
		)
		if err != nil {
			panic(err)
		}
		r.Delete(deleteHandler.FullPath, deleteHandler.ServeHTTP)
	}()
}

func NewLogicalThingClaimFromItem(item map[string]any) (any, error) {
	object := &LogicalThingClaim{}

	err := object.FromItem(item)
	if err != nil {
		return nil, err
	}

	return object, nil
}

func init() {
	register(
		LogicalThingClaimTable,
		LogicalThingClaim{},
		NewLogicalThingClaimFromItem,
		"/logical-thing-claims",
		MutateRouterForLogicalThingClaim,
	)
}