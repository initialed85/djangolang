package model_generated_from_schema

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"maps"
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
)

type Trigger struct {
	ID                                    uuid.UUID    `json:"id"`
	CreatedAt                             time.Time    `json:"created_at"`
	UpdatedAt                             time.Time    `json:"updated_at"`
	DeletedAt                             *time.Time   `json:"deleted_at"`
	JobExecutorClaimedUntil               time.Time    `json:"job_executor_claimed_until"`
	RuleID                                uuid.UUID    `json:"rule_id"`
	RuleIDObject                          *Rule        `json:"rule_id_object"`
	JobID                                 uuid.UUID    `json:"job_id"`
	JobIDObject                           *Job         `json:"job_id_object"`
	ReferencedByExecutionTriggerIDObjects []*Execution `json:"referenced_by_execution_trigger_id_objects"`
}

var TriggerTable = "trigger"

var TriggerTableWithSchema = fmt.Sprintf("%s.%s", schema, TriggerTable)

var TriggerTableNamespaceID int32 = 1337 + 9

var (
	TriggerTableIDColumn                      = "id"
	TriggerTableCreatedAtColumn               = "created_at"
	TriggerTableUpdatedAtColumn               = "updated_at"
	TriggerTableDeletedAtColumn               = "deleted_at"
	TriggerTableJobExecutorClaimedUntilColumn = "job_executor_claimed_until"
	TriggerTableRuleIDColumn                  = "rule_id"
	TriggerTableJobIDColumn                   = "job_id"
)

var (
	TriggerTableIDColumnWithTypeCast                      = `"id" AS id`
	TriggerTableCreatedAtColumnWithTypeCast               = `"created_at" AS created_at`
	TriggerTableUpdatedAtColumnWithTypeCast               = `"updated_at" AS updated_at`
	TriggerTableDeletedAtColumnWithTypeCast               = `"deleted_at" AS deleted_at`
	TriggerTableJobExecutorClaimedUntilColumnWithTypeCast = `"job_executor_claimed_until" AS job_executor_claimed_until`
	TriggerTableRuleIDColumnWithTypeCast                  = `"rule_id" AS rule_id`
	TriggerTableJobIDColumnWithTypeCast                   = `"job_id" AS job_id`
)

var TriggerTableColumns = []string{
	TriggerTableIDColumn,
	TriggerTableCreatedAtColumn,
	TriggerTableUpdatedAtColumn,
	TriggerTableDeletedAtColumn,
	TriggerTableJobExecutorClaimedUntilColumn,
	TriggerTableRuleIDColumn,
	TriggerTableJobIDColumn,
}

var TriggerTableColumnsWithTypeCasts = []string{
	TriggerTableIDColumnWithTypeCast,
	TriggerTableCreatedAtColumnWithTypeCast,
	TriggerTableUpdatedAtColumnWithTypeCast,
	TriggerTableDeletedAtColumnWithTypeCast,
	TriggerTableJobExecutorClaimedUntilColumnWithTypeCast,
	TriggerTableRuleIDColumnWithTypeCast,
	TriggerTableJobIDColumnWithTypeCast,
}

var TriggerIntrospectedTable *introspect.Table

var TriggerTableColumnLookup map[string]*introspect.Column

var (
	TriggerTablePrimaryKeyColumn = TriggerTableIDColumn
)

func init() {
	TriggerIntrospectedTable = tableByName[TriggerTable]

	/* only needed during templating */
	if TriggerIntrospectedTable == nil {
		TriggerIntrospectedTable = &introspect.Table{}
	}

	TriggerTableColumnLookup = TriggerIntrospectedTable.ColumnByName
}

type TriggerOnePathParams struct {
	PrimaryKey uuid.UUID `json:"primaryKey"`
}

type TriggerLoadQueryParams struct {
	Depth *int `json:"depth"`
}

type TriggerJobExecutorClaimRequest struct {
	Until          time.Time `json:"until"`
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

func (m *Trigger) GetPrimaryKeyColumn() string {
	return TriggerTablePrimaryKeyColumn
}

func (m *Trigger) GetPrimaryKeyValue() any {
	return m.ID
}

func (m *Trigger) FromItem(item map[string]any) error {
	if item == nil {
		return fmt.Errorf(
			"item unexpectedly nil during TriggerFromItem",
		)
	}

	if len(item) == 0 {
		return fmt.Errorf(
			"item unexpectedly empty during TriggerFromItem",
		)
	}

	wrapError := func(k string, v any, err error) error {
		return fmt.Errorf("%v: %#+v; error; %v", k, v, err)
	}

	for k, v := range item {
		_, ok := TriggerTableColumnLookup[k]
		if !ok {
			return fmt.Errorf(
				"item contained unexpected key %#+v during TriggerFromItem; item: %#+v",
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

		case "job_executor_claimed_until":
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
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uujob_executor_claimed_until.UUID", temp1))
				}
			}

			m.JobExecutorClaimedUntil = temp2

		case "rule_id":
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
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uurule_id.UUID", temp1))
				}
			}

			m.RuleID = temp2

		case "job_id":
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
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uujob_id.UUID", temp1))
				}
			}

			m.JobID = temp2

		}
	}

	return nil
}

func (m *Trigger) ToItem() map[string]any {
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

func (m *Trigger) Reload(ctx context.Context, tx pgx.Tx, includeDeleteds ...bool) error {
	extraWhere := ""
	if len(includeDeleteds) > 0 && includeDeleteds[0] {
		if slices.Contains(TriggerTableColumns, "deleted_at") {
			extraWhere = "\n    AND (deleted_at IS null OR deleted_at IS NOT null)"
		}
	}

	ctx, cleanup := query.WithQueryID(ctx)
	defer cleanup()

	ctx = query.WithMaxDepth(ctx, nil)

	o, _, _, _, _, err := SelectTrigger(
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
	m.JobExecutorClaimedUntil = o.JobExecutorClaimedUntil
	m.RuleID = o.RuleID
	m.RuleIDObject = o.RuleIDObject
	m.JobID = o.JobID
	m.JobIDObject = o.JobIDObject
	m.ReferencedByExecutionTriggerIDObjects = o.ReferencedByExecutionTriggerIDObjects

	return nil
}

func (m *Trigger) GetColumnsAndValues(setPrimaryKey bool, setZeroValues bool, forceSetValuesForFields ...string) ([]string, []any, error) {
	columns := make([]string, 0)
	values := make([]any, 0)

	if setPrimaryKey && (setZeroValues || !types.IsZeroUUID(m.ID) || slices.Contains(forceSetValuesForFields, TriggerTableIDColumn) || isRequired(TriggerTableColumnLookup, TriggerTableIDColumn)) {
		columns = append(columns, TriggerTableIDColumn)

		v, err := types.FormatUUID(m.ID)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to handle m.ID; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.CreatedAt) || slices.Contains(forceSetValuesForFields, TriggerTableCreatedAtColumn) || isRequired(TriggerTableColumnLookup, TriggerTableCreatedAtColumn) {
		columns = append(columns, TriggerTableCreatedAtColumn)

		v, err := types.FormatTime(m.CreatedAt)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to handle m.CreatedAt; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.UpdatedAt) || slices.Contains(forceSetValuesForFields, TriggerTableUpdatedAtColumn) || isRequired(TriggerTableColumnLookup, TriggerTableUpdatedAtColumn) {
		columns = append(columns, TriggerTableUpdatedAtColumn)

		v, err := types.FormatTime(m.UpdatedAt)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to handle m.UpdatedAt; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.DeletedAt) || slices.Contains(forceSetValuesForFields, TriggerTableDeletedAtColumn) || isRequired(TriggerTableColumnLookup, TriggerTableDeletedAtColumn) {
		columns = append(columns, TriggerTableDeletedAtColumn)

		v, err := types.FormatTime(m.DeletedAt)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to handle m.DeletedAt; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.JobExecutorClaimedUntil) || slices.Contains(forceSetValuesForFields, TriggerTableJobExecutorClaimedUntilColumn) || isRequired(TriggerTableColumnLookup, TriggerTableJobExecutorClaimedUntilColumn) {
		columns = append(columns, TriggerTableJobExecutorClaimedUntilColumn)

		v, err := types.FormatTime(m.JobExecutorClaimedUntil)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to handle m.JobExecutorClaimedUntil; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroUUID(m.RuleID) || slices.Contains(forceSetValuesForFields, TriggerTableRuleIDColumn) || isRequired(TriggerTableColumnLookup, TriggerTableRuleIDColumn) {
		columns = append(columns, TriggerTableRuleIDColumn)

		v, err := types.FormatUUID(m.RuleID)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to handle m.RuleID; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroUUID(m.JobID) || slices.Contains(forceSetValuesForFields, TriggerTableJobIDColumn) || isRequired(TriggerTableColumnLookup, TriggerTableJobIDColumn) {
		columns = append(columns, TriggerTableJobIDColumn)

		v, err := types.FormatUUID(m.JobID)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to handle m.JobID; %v", err)
		}

		values = append(values, v)
	}

	return columns, values, nil
}

func (m *Trigger) Insert(ctx context.Context, tx pgx.Tx, setPrimaryKey bool, setZeroValues bool, forceSetValuesForFields ...string) error {
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
		TriggerTableWithSchema,
		columns,
		nil,
		nil,
		nil,
		false,
		false,
		TriggerTableColumns,
		values...,
	)
	if err != nil {
		return fmt.Errorf("failed to insert %#+v; %v", m, err)
	}
	v := (*item)[TriggerTableIDColumn]

	if v == nil {
		return fmt.Errorf("failed to find %v in %#+v", TriggerTableIDColumn, item)
	}

	wrapError := func(err error) error {
		return fmt.Errorf(
			"failed to treat %v: %#+v as uuid.UUID: %v",
			TriggerTableIDColumn,
			(*item)[TriggerTableIDColumn],
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

func (m *Trigger) Update(ctx context.Context, tx pgx.Tx, setZeroValues bool, forceSetValuesForFields ...string) error {
	columns := make([]string, 0)
	values := make([]any, 0)

	if setZeroValues || !types.IsZeroTime(m.CreatedAt) || slices.Contains(forceSetValuesForFields, TriggerTableCreatedAtColumn) {
		columns = append(columns, TriggerTableCreatedAtColumn)

		v, err := types.FormatTime(m.CreatedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.CreatedAt; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.UpdatedAt) || slices.Contains(forceSetValuesForFields, TriggerTableUpdatedAtColumn) {
		columns = append(columns, TriggerTableUpdatedAtColumn)

		v, err := types.FormatTime(m.UpdatedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.UpdatedAt; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.DeletedAt) || slices.Contains(forceSetValuesForFields, TriggerTableDeletedAtColumn) {
		columns = append(columns, TriggerTableDeletedAtColumn)

		v, err := types.FormatTime(m.DeletedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.DeletedAt; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.JobExecutorClaimedUntil) || slices.Contains(forceSetValuesForFields, TriggerTableJobExecutorClaimedUntilColumn) {
		columns = append(columns, TriggerTableJobExecutorClaimedUntilColumn)

		v, err := types.FormatTime(m.JobExecutorClaimedUntil)
		if err != nil {
			return fmt.Errorf("failed to handle m.JobExecutorClaimedUntil; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroUUID(m.RuleID) || slices.Contains(forceSetValuesForFields, TriggerTableRuleIDColumn) {
		columns = append(columns, TriggerTableRuleIDColumn)

		v, err := types.FormatUUID(m.RuleID)
		if err != nil {
			return fmt.Errorf("failed to handle m.RuleID; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroUUID(m.JobID) || slices.Contains(forceSetValuesForFields, TriggerTableJobIDColumn) {
		columns = append(columns, TriggerTableJobIDColumn)

		v, err := types.FormatUUID(m.JobID)
		if err != nil {
			return fmt.Errorf("failed to handle m.JobID; %v", err)
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
		TriggerTableWithSchema,
		columns,
		fmt.Sprintf("%v = $$??", TriggerTableIDColumn),
		TriggerTableColumns,
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

func (m *Trigger) Delete(ctx context.Context, tx pgx.Tx, hardDeletes ...bool) error {
	hardDelete := false
	if len(hardDeletes) > 0 {
		hardDelete = hardDeletes[0]
	}

	if !hardDelete && slices.Contains(TriggerTableColumns, "deleted_at") {
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
		TriggerTableWithSchema,
		fmt.Sprintf("%v = $$??", TriggerTableIDColumn),
		values...,
	)
	if err != nil {
		return fmt.Errorf("failed to delete %#+v; %v", m, err)
	}

	_ = m.Reload(ctx, tx, true)

	return nil
}

func (m *Trigger) LockTable(ctx context.Context, tx pgx.Tx, timeouts ...time.Duration) error {
	return query.LockTable(ctx, tx, TriggerTableWithSchema, timeouts...)
}

func (m *Trigger) LockTableWithRetries(ctx context.Context, tx pgx.Tx, overallTimeout time.Duration, individualAttempttimeout time.Duration) error {
	return query.LockTableWithRetries(ctx, tx, TriggerTableWithSchema, overallTimeout, individualAttempttimeout)
}

func (m *Trigger) AdvisoryLock(ctx context.Context, tx pgx.Tx, key int32, timeouts ...time.Duration) error {
	return query.AdvisoryLock(ctx, tx, TriggerTableNamespaceID, key, timeouts...)
}

func (m *Trigger) AdvisoryLockWithRetries(ctx context.Context, tx pgx.Tx, key int32, overallTimeout time.Duration, individualAttempttimeout time.Duration) error {
	return query.AdvisoryLockWithRetries(ctx, tx, TriggerTableNamespaceID, key, overallTimeout, individualAttempttimeout)
}

func (m *Trigger) JobExecutorClaim(ctx context.Context, tx pgx.Tx, until time.Time, timeout time.Duration) error {
	err := m.AdvisoryLockWithRetries(ctx, tx, math.MinInt32, timeout, time.Second*1)
	if err != nil {
		return fmt.Errorf("failed to claim (advisory lock): %s", err.Error())
	}

	_, _, _, _, _, err = SelectTrigger(
		ctx,
		tx,
		fmt.Sprintf(
			"%s = $$?? AND (job_executor_claimed_until IS null OR job_executor_claimed_until < now())",
			TriggerTablePrimaryKeyColumn,
		),
		m.GetPrimaryKeyValue(),
	)
	if err != nil {
		return fmt.Errorf("failed to claim (select): %s", err.Error())
	}

	m.JobExecutorClaimedUntil = until

	err = m.Update(ctx, tx, false)
	if err != nil {
		return fmt.Errorf("failed to claim (update): %s", err.Error())
	}

	return nil
}

func SelectTriggers(ctx context.Context, tx pgx.Tx, where string, orderBy *string, limit *int, offset *int, values ...any) ([]*Trigger, int64, int64, int64, int64, error) {
	before := time.Now()

	if config.Debug() {
		log.Printf("entered SelectTriggers")

		defer func() {
			log.Printf("exited SelectTriggers in %s", time.Since(before))
		}()
	}
	if slices.Contains(TriggerTableColumns, "deleted_at") {
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

	shouldLoad := query.ShouldLoad(ctx, TriggerTable) || query.ShouldLoad(ctx, fmt.Sprintf("referenced_by_%s", TriggerTable))

	var ok bool
	ctx, ok = query.HandleQueryPathGraphCycles(ctx, fmt.Sprintf("%s{%v}", TriggerTable, nil), !isLoadQuery)
	if !ok && !shouldLoad {
		if config.Debug() {
			log.Printf("skipping SelectTrigger early (query.ShouldLoad(): %v, query.HandleQueryPathGraphCycles(): %v)", shouldLoad, ok)
		}
		return []*Trigger{}, 0, 0, 0, 0, nil
	}

	var items *[]map[string]any
	var count int64
	var totalCount int64
	var page int64
	var totalPages int64
	var err error

	useInstead, shouldSkip := query.ShouldSkip[Trigger](ctx)
	if !shouldSkip {
		items, count, totalCount, page, totalPages, err = query.Select(
			ctx,
			tx,
			TriggerTableColumnsWithTypeCasts,
			TriggerTableWithSchema,
			where,
			orderBy,
			limit,
			offset,
			values...,
		)
		if err != nil {
			return nil, 0, 0, 0, 0, fmt.Errorf("failed to call SelectTriggers; %v", err)
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

	objects := make([]*Trigger, 0)

	for _, item := range *items {
		var object *Trigger

		if !shouldSkip {
			object = &Trigger{}
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

		if !types.IsZeroUUID(object.RuleID) {
			ctx, ok := query.HandleQueryPathGraphCycles(ctx, fmt.Sprintf("%s{%v}", RuleTable, object.RuleID), true)
			shouldLoad := query.ShouldLoad(ctx, RuleTable)
			if ok || shouldLoad {
				thisBefore := time.Now()

				if config.Debug() {
					log.Printf("loading SelectTriggers->SelectRule for object.RuleIDObject{%s: %v}", RuleTablePrimaryKeyColumn, object.RuleID)
				}

				object.RuleIDObject, _, _, _, _, err = SelectRule(
					ctx,
					tx,
					fmt.Sprintf("%v = $1", RuleTablePrimaryKeyColumn),
					object.RuleID,
				)
				if err != nil {
					if !errors.Is(err, sql.ErrNoRows) {
						return nil, 0, 0, 0, 0, err
					}
				}

				if config.Debug() {
					log.Printf("loaded SelectTriggers->SelectRule for object.RuleIDObject in %s", time.Since(thisBefore))
				}
			}
		}

		if !types.IsZeroUUID(object.JobID) {
			ctx, ok := query.HandleQueryPathGraphCycles(ctx, fmt.Sprintf("%s{%v}", JobTable, object.JobID), true)
			shouldLoad := query.ShouldLoad(ctx, JobTable)
			if ok || shouldLoad {
				thisBefore := time.Now()

				if config.Debug() {
					log.Printf("loading SelectTriggers->SelectJob for object.JobIDObject{%s: %v}", JobTablePrimaryKeyColumn, object.JobID)
				}

				object.JobIDObject, _, _, _, _, err = SelectJob(
					ctx,
					tx,
					fmt.Sprintf("%v = $1", JobTablePrimaryKeyColumn),
					object.JobID,
				)
				if err != nil {
					if !errors.Is(err, sql.ErrNoRows) {
						return nil, 0, 0, 0, 0, err
					}
				}

				if config.Debug() {
					log.Printf("loaded SelectTriggers->SelectJob for object.JobIDObject in %s", time.Since(thisBefore))
				}
			}
		}

		err = func() error {
			shouldLoad := query.ShouldLoad(ctx, fmt.Sprintf("referenced_by_%s", ExecutionTable))
			ctx, ok := query.HandleQueryPathGraphCycles(ctx, fmt.Sprintf("__ReferencedBy__%s{%v}", ExecutionTable, object.GetPrimaryKeyValue()), true)
			if ok || shouldLoad {
				thisBefore := time.Now()

				if config.Debug() {
					log.Printf("loading SelectTriggers->SelectExecutions for object.ReferencedByExecutionTriggerIDObjects")
				}

				object.ReferencedByExecutionTriggerIDObjects, _, _, _, _, err = SelectExecutions(
					ctx,
					tx,
					fmt.Sprintf("%v = $1", ExecutionTableTriggerIDColumn),
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
					log.Printf("loaded SelectTriggers->SelectExecutions for object.ReferencedByExecutionTriggerIDObjects in %s", time.Since(thisBefore))
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

func SelectTrigger(ctx context.Context, tx pgx.Tx, where string, values ...any) (*Trigger, int64, int64, int64, int64, error) {
	ctx, cleanup := query.WithQueryID(ctx)
	defer cleanup()

	ctx = query.WithMaxDepth(ctx, nil)

	objects, _, _, _, _, err := SelectTriggers(
		ctx,
		tx,
		where,
		nil,
		helpers.Ptr(2),
		helpers.Ptr(0),
		values...,
	)
	if err != nil {
		return nil, 0, 0, 0, 0, fmt.Errorf("failed to call SelectTrigger; %v", err)
	}

	if len(objects) > 1 {
		return nil, 0, 0, 0, 0, fmt.Errorf("attempt to call SelectTrigger returned more than 1 row")
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

func InsertTriggers(ctx context.Context, tx pgx.Tx, objects []*Trigger, setPrimaryKey bool, setZeroValues bool, forceSetValuesForFields ...string) ([]*Trigger, error) {
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
		TriggerTableWithSchema,
		columns,
		nil,
		nil,
		nil,
		false,
		false,
		TriggerTableColumns,
		values...,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to bulk insert %d objects; %v", len(objects), err)
	}

	returnedObjects := make([]*Trigger, 0)

	for _, item := range items {
		v := &Trigger{}
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

func JobExecutorClaimTrigger(ctx context.Context, tx pgx.Tx, until time.Time, timeout time.Duration, where string, values ...any) (*Trigger, error) {
	m := &Trigger{}

	err := m.AdvisoryLockWithRetries(ctx, tx, math.MinInt32, timeout, time.Second*1)
	if err != nil {
		return nil, fmt.Errorf("failed to claim: %s", err.Error())
	}

	if strings.TrimSpace(where) != "" {
		where += " AND\n"
	}

	where += "    (job_executor_claimed_until IS null OR job_executor_claimed_until < now())"

	ms, _, _, _, _, err := SelectTriggers(
		ctx,
		tx,
		where,
		helpers.Ptr(
			"job_executor_claimed_until ASC",
		),
		helpers.Ptr(1),
		nil,
		values...,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to claim: %s", err.Error())
	}

	if len(ms) == 0 {
		return nil, nil
	}

	m = ms[0]

	m.JobExecutorClaimedUntil = until

	err = m.Update(ctx, tx, false)
	if err != nil {
		return nil, fmt.Errorf("failed to claim: %s", err.Error())
	}

	return m, nil
}

func handleGetTriggers(arguments *server.SelectManyArguments, db *pgxpool.Pool) ([]*Trigger, int64, int64, int64, int64, error) {
	tx, err := db.Begin(arguments.Ctx)
	if err != nil {
		return nil, 0, 0, 0, 0, err
	}

	defer func() {
		_ = tx.Rollback(arguments.Ctx)
	}()

	objects, count, totalCount, page, totalPages, err := SelectTriggers(arguments.Ctx, tx, arguments.Where, arguments.OrderBy, arguments.Limit, arguments.Offset, arguments.Values...)
	if err != nil {
		return nil, 0, 0, 0, 0, err
	}

	err = tx.Commit(arguments.Ctx)
	if err != nil {
		return nil, 0, 0, 0, 0, err
	}

	return objects, count, totalCount, page, totalPages, nil
}

func handleGetTrigger(arguments *server.SelectOneArguments, db *pgxpool.Pool, primaryKey uuid.UUID) ([]*Trigger, int64, int64, int64, int64, error) {
	tx, err := db.Begin(arguments.Ctx)
	if err != nil {
		return nil, 0, 0, 0, 0, err
	}

	defer func() {
		_ = tx.Rollback(arguments.Ctx)
	}()

	object, count, totalCount, page, totalPages, err := SelectTrigger(arguments.Ctx, tx, arguments.Where, arguments.Values...)
	if err != nil {
		return nil, 0, 0, 0, 0, err
	}

	err = tx.Commit(arguments.Ctx)
	if err != nil {
		return nil, 0, 0, 0, 0, err
	}

	return []*Trigger{object}, count, totalCount, page, totalPages, nil
}

func handlePostTrigger(arguments *server.LoadArguments, db *pgxpool.Pool, waitForChange server.WaitForChange, objects []*Trigger, forceSetValuesForFieldsByObjectIndex [][]string) ([]*Trigger, int64, int64, int64, int64, error) {
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

	returnedObjects, err := InsertTriggers(arguments.Ctx, tx, objects, false, false, slices.Collect(maps.Keys(forceSetValuesForFieldsByObjectIndexMaximal))...)
	if err != nil {
		err = fmt.Errorf("failed to insert %d objects; %v", len(objects), err)
		return nil, 0, 0, 0, 0, err
	}

	copy(objects, returnedObjects)

	errs := make(chan error, 1)
	go func() {
		_, err := waitForChange(arguments.Ctx, []stream.Action{stream.INSERT}, TriggerTable, xid)
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

func handlePutTrigger(arguments *server.LoadArguments, db *pgxpool.Pool, waitForChange server.WaitForChange, object *Trigger) ([]*Trigger, int64, int64, int64, int64, error) {
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
		_, err := waitForChange(arguments.Ctx, []stream.Action{stream.UPDATE, stream.SOFT_DELETE, stream.SOFT_RESTORE, stream.SOFT_UPDATE}, TriggerTable, xid)
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

	return []*Trigger{object}, count, totalCount, page, totalPages, nil
}

func handlePatchTrigger(arguments *server.LoadArguments, db *pgxpool.Pool, waitForChange server.WaitForChange, object *Trigger, forceSetValuesForFields []string) ([]*Trigger, int64, int64, int64, int64, error) {
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
		_, err := waitForChange(arguments.Ctx, []stream.Action{stream.UPDATE, stream.SOFT_DELETE, stream.SOFT_RESTORE, stream.SOFT_UPDATE}, TriggerTable, xid)
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

	return []*Trigger{object}, count, totalCount, page, totalPages, nil
}

func handleDeleteTrigger(arguments *server.LoadArguments, db *pgxpool.Pool, waitForChange server.WaitForChange, object *Trigger) error {
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
		_, err := waitForChange(arguments.Ctx, []stream.Action{stream.DELETE, stream.SOFT_DELETE}, TriggerTable, xid)
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

func MutateRouterForTrigger(r chi.Router, db *pgxpool.Pool, redisPool *redis.Pool, objectMiddlewares []server.ObjectMiddleware, waitForChange server.WaitForChange) {

	func() {
		postHandlerForJobExecutorClaim, err := getHTTPHandler(
			http.MethodPost,
			"/job-executor-claim-trigger",
			http.StatusOK,
			func(
				ctx context.Context,
				pathParams server.EmptyPathParams,
				queryParams server.EmptyQueryParams,
				req TriggerJobExecutorClaimRequest,
				rawReq any,
			) (server.Response[Trigger], error) {
				tx, err := db.Begin(ctx)
				if err != nil {
					return server.Response[Trigger]{}, err
				}

				defer func() {
					_ = tx.Rollback(ctx)
				}()

				object, err := JobExecutorClaimTrigger(ctx, tx, req.Until, time.Millisecond*time.Duration(req.TimeoutSeconds*1000), "")
				if err != nil {
					return server.Response[Trigger]{}, err
				}

				count := int64(0)

				totalCount := int64(0)

				limit := int64(0)

				offset := int64(0)

				if object == nil {
					return server.Response[Trigger]{
						Status:     http.StatusOK,
						Success:    true,
						Error:      nil,
						Objects:    []*Trigger{},
						Count:      count,
						TotalCount: totalCount,
						Limit:      limit,
						Offset:     offset,
					}, nil
				}

				err = tx.Commit(ctx)
				if err != nil {
					return server.Response[Trigger]{}, err
				}

				return server.Response[Trigger]{
					Status:     http.StatusOK,
					Success:    true,
					Error:      nil,
					Objects:    []*Trigger{object},
					Count:      count,
					TotalCount: totalCount,
					Limit:      limit,
					Offset:     offset,
				}, nil
			},
			Trigger{},
			TriggerIntrospectedTable,
		)
		if err != nil {
			panic(err)
		}
		r.Post(postHandlerForJobExecutorClaim.FullPath, postHandlerForJobExecutorClaim.ServeHTTP)

		postHandlerForJobExecutorClaimOne, err := getHTTPHandler(
			http.MethodPost,
			"/triggers/{primaryKey}/job-executor-claim",
			http.StatusOK,
			func(
				ctx context.Context,
				pathParams TriggerOnePathParams,
				queryParams TriggerLoadQueryParams,
				req TriggerJobExecutorClaimRequest,
				rawReq any,
			) (server.Response[Trigger], error) {
				before := time.Now()

				redisConn := redisPool.Get()
				defer func() {
					_ = redisConn.Close()
				}()

				arguments, err := server.GetSelectOneArguments(ctx, queryParams.Depth, TriggerIntrospectedTable, pathParams.PrimaryKey, nil, nil)
				if err != nil {
					if config.Debug() {
						log.Printf("request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return server.Response[Trigger]{}, err
				}

				/* note: deliberately no attempt at a cache hit */

				var object *Trigger
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

					object, count, totalCount, _, _, err = SelectTrigger(arguments.Ctx, tx, arguments.Where, arguments.Values...)
					if err != nil {
						return fmt.Errorf("failed to select object to claim: %s", err.Error())
					}

					err = object.JobExecutorClaim(arguments.Ctx, tx, req.Until, time.Millisecond*time.Duration(req.TimeoutSeconds*1000))
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

					return server.Response[Trigger]{}, err
				}

				limit := int64(0)

				offset := int64(0)

				response := server.Response[Trigger]{
					Status:     http.StatusOK,
					Success:    true,
					Error:      nil,
					Objects:    []*Trigger{object},
					Count:      count,
					TotalCount: totalCount,
					Limit:      limit,
					Offset:     offset,
				}

				return response, nil
			},
			Trigger{},
			TriggerIntrospectedTable,
		)
		if err != nil {
			panic(err)
		}
		r.Post(postHandlerForJobExecutorClaimOne.FullPath, postHandlerForJobExecutorClaimOne.ServeHTTP)
	}()

	func() {
		getManyHandler, err := getHTTPHandler(
			http.MethodGet,
			"/triggers",
			http.StatusOK,
			func(
				ctx context.Context,
				pathParams server.EmptyPathParams,
				queryParams map[string]any,
				req server.EmptyRequest,
				rawReq any,
			) (server.Response[Trigger], error) {
				before := time.Now()

				redisConn := redisPool.Get()
				defer func() {
					_ = redisConn.Close()
				}()

				arguments, err := server.GetSelectManyArguments(ctx, queryParams, TriggerIntrospectedTable, nil, nil)
				if err != nil {
					if config.Debug() {
						log.Printf("request cache not yet reached; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return server.Response[Trigger]{}, err
				}

				cachedResponseAsJSON, cacheHit, err := server.GetCachedResponseAsJSON(arguments.RequestHash, redisConn)
				if err != nil {
					if config.Debug() {
						log.Printf("request cache failed; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return server.Response[Trigger]{}, err
				}

				if cacheHit {
					var cachedResponse server.Response[Trigger]

					/* TODO: it'd be nice to be able to avoid this (i.e. just pass straight through) */
					err = json.Unmarshal(cachedResponseAsJSON, &cachedResponse)
					if err != nil {
						if config.Debug() {
							log.Printf("request cache hit but failed unmarshal; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
						}

						return server.Response[Trigger]{}, err
					}

					if config.Debug() {
						log.Printf("request cache hit; request succeeded in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return cachedResponse, nil
				}

				objects, count, totalCount, _, _, err := handleGetTriggers(arguments, db)
				if err != nil {
					if config.Debug() {
						log.Printf("request cache missed; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return server.Response[Trigger]{}, err
				}

				limit := int64(0)
				if arguments.Limit != nil {
					limit = int64(*arguments.Limit)
				}

				offset := int64(0)
				if arguments.Offset != nil {
					offset = int64(*arguments.Offset)
				}

				response := server.Response[Trigger]{
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

					return server.Response[Trigger]{}, err
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
			Trigger{},
			TriggerIntrospectedTable,
		)
		if err != nil {
			panic(err)
		}
		r.Get(getManyHandler.FullPath, getManyHandler.ServeHTTP)
	}()

	func() {
		getOneHandler, err := getHTTPHandler(
			http.MethodGet,
			"/triggers/{primaryKey}",
			http.StatusOK,
			func(
				ctx context.Context,
				pathParams TriggerOnePathParams,
				queryParams TriggerLoadQueryParams,
				req server.EmptyRequest,
				rawReq any,
			) (server.Response[Trigger], error) {
				before := time.Now()

				redisConn := redisPool.Get()
				defer func() {
					_ = redisConn.Close()
				}()

				arguments, err := server.GetSelectOneArguments(ctx, queryParams.Depth, TriggerIntrospectedTable, pathParams.PrimaryKey, nil, nil)
				if err != nil {
					if config.Debug() {
						log.Printf("request cache not yet reached; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return server.Response[Trigger]{}, err
				}

				cachedResponseAsJSON, cacheHit, err := server.GetCachedResponseAsJSON(arguments.RequestHash, redisConn)
				if err != nil {
					if config.Debug() {
						log.Printf("request cache failed; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return server.Response[Trigger]{}, err
				}

				if cacheHit {
					var cachedResponse server.Response[Trigger]

					/* TODO: it'd be nice to be able to avoid this (i.e. just pass straight through) */
					err = json.Unmarshal(cachedResponseAsJSON, &cachedResponse)
					if err != nil {
						if config.Debug() {
							log.Printf("request cache hit but failed unmarshal; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
						}

						return server.Response[Trigger]{}, err
					}

					if config.Debug() {
						log.Printf("request cache hit; request succeeded in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return cachedResponse, nil
				}

				objects, count, totalCount, _, _, err := handleGetTrigger(arguments, db, pathParams.PrimaryKey)
				if err != nil {
					if config.Debug() {
						log.Printf("request cache missed; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return server.Response[Trigger]{}, err
				}

				limit := int64(0)

				offset := int64(0)

				response := server.Response[Trigger]{
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

					return server.Response[Trigger]{}, err
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
			Trigger{},
			TriggerIntrospectedTable,
		)
		if err != nil {
			panic(err)
		}
		r.Get(getOneHandler.FullPath, getOneHandler.ServeHTTP)
	}()

	func() {
		postHandler, err := getHTTPHandler(
			http.MethodPost,
			"/triggers",
			http.StatusCreated,
			func(
				ctx context.Context,
				pathParams server.EmptyPathParams,
				queryParams TriggerLoadQueryParams,
				req []*Trigger,
				rawReq any,
			) (server.Response[Trigger], error) {
				allRawItems, ok := rawReq.([]any)
				if !ok {
					return server.Response[Trigger]{}, fmt.Errorf("failed to cast %#+v to []map[string]any", rawReq)
				}

				allItems := make([]map[string]any, 0)
				for _, rawItem := range allRawItems {
					item, ok := rawItem.(map[string]any)
					if !ok {
						return server.Response[Trigger]{}, fmt.Errorf("failed to cast %#+v to map[string]any", rawItem)
					}

					allItems = append(allItems, item)
				}

				forceSetValuesForFieldsByObjectIndex := make([][]string, 0)
				for _, item := range allItems {
					forceSetValuesForFields := make([]string, 0)
					for _, possibleField := range slices.Collect(maps.Keys(item)) {
						if !slices.Contains(TriggerTableColumns, possibleField) {
							continue
						}

						forceSetValuesForFields = append(forceSetValuesForFields, possibleField)
					}
					forceSetValuesForFieldsByObjectIndex = append(forceSetValuesForFieldsByObjectIndex, forceSetValuesForFields)
				}

				arguments, err := server.GetLoadArguments(ctx, queryParams.Depth)
				if err != nil {
					return server.Response[Trigger]{}, err
				}

				objects, count, totalCount, _, _, err := handlePostTrigger(arguments, db, waitForChange, req, forceSetValuesForFieldsByObjectIndex)
				if err != nil {
					return server.Response[Trigger]{}, err
				}

				limit := int64(0)

				offset := int64(0)

				return server.Response[Trigger]{
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
			Trigger{},
			TriggerIntrospectedTable,
		)
		if err != nil {
			panic(err)
		}
		r.Post(postHandler.FullPath, postHandler.ServeHTTP)
	}()

	func() {
		putHandler, err := getHTTPHandler(
			http.MethodPatch,
			"/triggers/{primaryKey}",
			http.StatusOK,
			func(
				ctx context.Context,
				pathParams TriggerOnePathParams,
				queryParams TriggerLoadQueryParams,
				req Trigger,
				rawReq any,
			) (server.Response[Trigger], error) {
				item, ok := rawReq.(map[string]any)
				if !ok {
					return server.Response[Trigger]{}, fmt.Errorf("failed to cast %#+v to map[string]any", item)
				}

				arguments, err := server.GetLoadArguments(ctx, queryParams.Depth)
				if err != nil {
					return server.Response[Trigger]{}, err
				}

				object := &req
				object.ID = pathParams.PrimaryKey

				objects, count, totalCount, _, _, err := handlePutTrigger(arguments, db, waitForChange, object)
				if err != nil {
					return server.Response[Trigger]{}, err
				}

				limit := int64(0)

				offset := int64(0)

				return server.Response[Trigger]{
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
			Trigger{},
			TriggerIntrospectedTable,
		)
		if err != nil {
			panic(err)
		}
		r.Put(putHandler.FullPath, putHandler.ServeHTTP)
	}()

	func() {
		patchHandler, err := getHTTPHandler(
			http.MethodPatch,
			"/triggers/{primaryKey}",
			http.StatusOK,
			func(
				ctx context.Context,
				pathParams TriggerOnePathParams,
				queryParams TriggerLoadQueryParams,
				req Trigger,
				rawReq any,
			) (server.Response[Trigger], error) {
				item, ok := rawReq.(map[string]any)
				if !ok {
					return server.Response[Trigger]{}, fmt.Errorf("failed to cast %#+v to map[string]any", item)
				}

				forceSetValuesForFields := make([]string, 0)
				for _, possibleField := range slices.Collect(maps.Keys(item)) {
					if !slices.Contains(TriggerTableColumns, possibleField) {
						continue
					}

					forceSetValuesForFields = append(forceSetValuesForFields, possibleField)
				}

				arguments, err := server.GetLoadArguments(ctx, queryParams.Depth)
				if err != nil {
					return server.Response[Trigger]{}, err
				}

				object := &req
				object.ID = pathParams.PrimaryKey

				objects, count, totalCount, _, _, err := handlePatchTrigger(arguments, db, waitForChange, object, forceSetValuesForFields)
				if err != nil {
					return server.Response[Trigger]{}, err
				}

				limit := int64(0)

				offset := int64(0)

				return server.Response[Trigger]{
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
			Trigger{},
			TriggerIntrospectedTable,
		)
		if err != nil {
			panic(err)
		}
		r.Patch(patchHandler.FullPath, patchHandler.ServeHTTP)
	}()

	func() {
		deleteHandler, err := getHTTPHandler(
			http.MethodDelete,
			"/triggers/{primaryKey}",
			http.StatusNoContent,
			func(
				ctx context.Context,
				pathParams TriggerOnePathParams,
				queryParams TriggerLoadQueryParams,
				req server.EmptyRequest,
				rawReq any,
			) (server.EmptyResponse, error) {
				arguments, err := server.GetLoadArguments(ctx, queryParams.Depth)
				if err != nil {
					return server.EmptyResponse{}, err
				}

				object := &Trigger{}
				object.ID = pathParams.PrimaryKey

				err = handleDeleteTrigger(arguments, db, waitForChange, object)
				if err != nil {
					return server.EmptyResponse{}, err
				}

				return server.EmptyResponse{}, nil
			},
			Trigger{},
			TriggerIntrospectedTable,
		)
		if err != nil {
			panic(err)
		}
		r.Delete(deleteHandler.FullPath, deleteHandler.ServeHTTP)
	}()
}

func NewTriggerFromItem(item map[string]any) (any, error) {
	object := &Trigger{}

	err := object.FromItem(item)
	if err != nil {
		return nil, err
	}

	return object, nil
}

func init() {
	register(
		TriggerTable,
		Trigger{},
		NewTriggerFromItem,
		"/triggers",
		MutateRouterForTrigger,
	)
}
