package model_generated_from_schema

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

type Execution struct {
	ID                                              uuid.UUID            `json:"id"`
	CreatedAt                                       time.Time            `json:"created_at"`
	UpdatedAt                                       time.Time            `json:"updated_at"`
	DeletedAt                                       *time.Time           `json:"deleted_at"`
	Status                                          string               `json:"status"`
	StartedAt                                       *time.Time           `json:"started_at"`
	EndedAt                                         *time.Time           `json:"ended_at"`
	JobExecutorClaimedUntil                         time.Time            `json:"job_executor_claimed_until"`
	TaskID                                          uuid.UUID            `json:"task_id"`
	TaskIDObject                                    *Task                `json:"task_id_object"`
	ReferencedByM2mRuleTriggerJobExecutionIDObjects []*M2mRuleTriggerJob `json:"referenced_by_m2m_rule_trigger_job_execution_id_objects"`
}

var ExecutionTable = "execution"

var ExecutionTableWithSchema = fmt.Sprintf("%s.%s", schema, ExecutionTable)

var ExecutionTableNamespaceID int32 = 1337 + 2

var (
	ExecutionTableIDColumn                      = "id"
	ExecutionTableCreatedAtColumn               = "created_at"
	ExecutionTableUpdatedAtColumn               = "updated_at"
	ExecutionTableDeletedAtColumn               = "deleted_at"
	ExecutionTableStatusColumn                  = "status"
	ExecutionTableStartedAtColumn               = "started_at"
	ExecutionTableEndedAtColumn                 = "ended_at"
	ExecutionTableJobExecutorClaimedUntilColumn = "job_executor_claimed_until"
	ExecutionTableTaskIDColumn                  = "task_id"
)

var (
	ExecutionTableIDColumnWithTypeCast                      = `"id" AS id`
	ExecutionTableCreatedAtColumnWithTypeCast               = `"created_at" AS created_at`
	ExecutionTableUpdatedAtColumnWithTypeCast               = `"updated_at" AS updated_at`
	ExecutionTableDeletedAtColumnWithTypeCast               = `"deleted_at" AS deleted_at`
	ExecutionTableStatusColumnWithTypeCast                  = `"status" AS status`
	ExecutionTableStartedAtColumnWithTypeCast               = `"started_at" AS started_at`
	ExecutionTableEndedAtColumnWithTypeCast                 = `"ended_at" AS ended_at`
	ExecutionTableJobExecutorClaimedUntilColumnWithTypeCast = `"job_executor_claimed_until" AS job_executor_claimed_until`
	ExecutionTableTaskIDColumnWithTypeCast                  = `"task_id" AS task_id`
)

var ExecutionTableColumns = []string{
	ExecutionTableIDColumn,
	ExecutionTableCreatedAtColumn,
	ExecutionTableUpdatedAtColumn,
	ExecutionTableDeletedAtColumn,
	ExecutionTableStatusColumn,
	ExecutionTableStartedAtColumn,
	ExecutionTableEndedAtColumn,
	ExecutionTableJobExecutorClaimedUntilColumn,
	ExecutionTableTaskIDColumn,
}

var ExecutionTableColumnsWithTypeCasts = []string{
	ExecutionTableIDColumnWithTypeCast,
	ExecutionTableCreatedAtColumnWithTypeCast,
	ExecutionTableUpdatedAtColumnWithTypeCast,
	ExecutionTableDeletedAtColumnWithTypeCast,
	ExecutionTableStatusColumnWithTypeCast,
	ExecutionTableStartedAtColumnWithTypeCast,
	ExecutionTableEndedAtColumnWithTypeCast,
	ExecutionTableJobExecutorClaimedUntilColumnWithTypeCast,
	ExecutionTableTaskIDColumnWithTypeCast,
}

var ExecutionIntrospectedTable *introspect.Table

var ExecutionTableColumnLookup map[string]*introspect.Column

var (
	ExecutionTablePrimaryKeyColumn = ExecutionTableIDColumn
)

func init() {
	ExecutionIntrospectedTable = tableByName[ExecutionTable]

	/* only needed during templating */
	if ExecutionIntrospectedTable == nil {
		ExecutionIntrospectedTable = &introspect.Table{}
	}

	ExecutionTableColumnLookup = ExecutionIntrospectedTable.ColumnByName
}

type ExecutionOnePathParams struct {
	PrimaryKey uuid.UUID `json:"primaryKey"`
}

type ExecutionLoadQueryParams struct {
	Depth *int `json:"depth"`
}

type ExecutionJobExecutorClaimRequest struct {
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

func (m *Execution) GetPrimaryKeyColumn() string {
	return ExecutionTablePrimaryKeyColumn
}

func (m *Execution) GetPrimaryKeyValue() any {
	return m.ID
}

func (m *Execution) FromItem(item map[string]any) error {
	if item == nil {
		return fmt.Errorf(
			"item unexpectedly nil during ExecutionFromItem",
		)
	}

	if len(item) == 0 {
		return fmt.Errorf(
			"item unexpectedly empty during ExecutionFromItem",
		)
	}

	wrapError := func(k string, v any, err error) error {
		return fmt.Errorf("%v: %#+v; error; %v", k, v, err)
	}

	for k, v := range item {
		_, ok := ExecutionTableColumnLookup[k]
		if !ok {
			return fmt.Errorf(
				"item contained unexpected key %#+v during ExecutionFromItem; item: %#+v",
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

		case "status":
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
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uustatus.UUID", temp1))
				}
			}

			m.Status = temp2

		case "started_at":
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
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uustarted_at.UUID", temp1))
				}
			}

			m.StartedAt = &temp2

		case "ended_at":
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
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uuended_at.UUID", temp1))
				}
			}

			m.EndedAt = &temp2

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

		case "task_id":
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
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uutask_id.UUID", temp1))
				}
			}

			m.TaskID = temp2

		}
	}

	return nil
}

func (m *Execution) Reload(ctx context.Context, tx pgx.Tx, includeDeleteds ...bool) error {
	extraWhere := ""
	if len(includeDeleteds) > 0 && includeDeleteds[0] {
		if slices.Contains(ExecutionTableColumns, "deleted_at") {
			extraWhere = "\n    AND (deleted_at IS null OR deleted_at IS NOT null)"
		}
	}

	ctx, cleanup := query.WithQueryID(ctx)
	defer cleanup()

	ctx = query.WithMaxDepth(ctx, nil)

	o, _, _, _, _, err := SelectExecution(
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
	m.Status = o.Status
	m.StartedAt = o.StartedAt
	m.EndedAt = o.EndedAt
	m.JobExecutorClaimedUntil = o.JobExecutorClaimedUntil
	m.TaskID = o.TaskID
	m.TaskIDObject = o.TaskIDObject
	m.ReferencedByM2mRuleTriggerJobExecutionIDObjects = o.ReferencedByM2mRuleTriggerJobExecutionIDObjects

	return nil
}

func (m *Execution) Insert(ctx context.Context, tx pgx.Tx, setPrimaryKey bool, setZeroValues bool, forceSetValuesForFields ...string) error {
	columns := make([]string, 0)
	values := make([]any, 0)

	if setPrimaryKey && (setZeroValues || !types.IsZeroUUID(m.ID) || slices.Contains(forceSetValuesForFields, ExecutionTableIDColumn) || isRequired(ExecutionTableColumnLookup, ExecutionTableIDColumn)) {
		columns = append(columns, ExecutionTableIDColumn)

		v, err := types.FormatUUID(m.ID)
		if err != nil {
			return fmt.Errorf("failed to handle m.ID; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.CreatedAt) || slices.Contains(forceSetValuesForFields, ExecutionTableCreatedAtColumn) || isRequired(ExecutionTableColumnLookup, ExecutionTableCreatedAtColumn) {
		columns = append(columns, ExecutionTableCreatedAtColumn)

		v, err := types.FormatTime(m.CreatedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.CreatedAt; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.UpdatedAt) || slices.Contains(forceSetValuesForFields, ExecutionTableUpdatedAtColumn) || isRequired(ExecutionTableColumnLookup, ExecutionTableUpdatedAtColumn) {
		columns = append(columns, ExecutionTableUpdatedAtColumn)

		v, err := types.FormatTime(m.UpdatedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.UpdatedAt; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.DeletedAt) || slices.Contains(forceSetValuesForFields, ExecutionTableDeletedAtColumn) || isRequired(ExecutionTableColumnLookup, ExecutionTableDeletedAtColumn) {
		columns = append(columns, ExecutionTableDeletedAtColumn)

		v, err := types.FormatTime(m.DeletedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.DeletedAt; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroString(m.Status) || slices.Contains(forceSetValuesForFields, ExecutionTableStatusColumn) || isRequired(ExecutionTableColumnLookup, ExecutionTableStatusColumn) {
		columns = append(columns, ExecutionTableStatusColumn)

		v, err := types.FormatString(m.Status)
		if err != nil {
			return fmt.Errorf("failed to handle m.Status; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.StartedAt) || slices.Contains(forceSetValuesForFields, ExecutionTableStartedAtColumn) || isRequired(ExecutionTableColumnLookup, ExecutionTableStartedAtColumn) {
		columns = append(columns, ExecutionTableStartedAtColumn)

		v, err := types.FormatTime(m.StartedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.StartedAt; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.EndedAt) || slices.Contains(forceSetValuesForFields, ExecutionTableEndedAtColumn) || isRequired(ExecutionTableColumnLookup, ExecutionTableEndedAtColumn) {
		columns = append(columns, ExecutionTableEndedAtColumn)

		v, err := types.FormatTime(m.EndedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.EndedAt; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.JobExecutorClaimedUntil) || slices.Contains(forceSetValuesForFields, ExecutionTableJobExecutorClaimedUntilColumn) || isRequired(ExecutionTableColumnLookup, ExecutionTableJobExecutorClaimedUntilColumn) {
		columns = append(columns, ExecutionTableJobExecutorClaimedUntilColumn)

		v, err := types.FormatTime(m.JobExecutorClaimedUntil)
		if err != nil {
			return fmt.Errorf("failed to handle m.JobExecutorClaimedUntil; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroUUID(m.TaskID) || slices.Contains(forceSetValuesForFields, ExecutionTableTaskIDColumn) || isRequired(ExecutionTableColumnLookup, ExecutionTableTaskIDColumn) {
		columns = append(columns, ExecutionTableTaskIDColumn)

		v, err := types.FormatUUID(m.TaskID)
		if err != nil {
			return fmt.Errorf("failed to handle m.TaskID; %v", err)
		}

		values = append(values, v)
	}

	ctx, cleanup := query.WithQueryID(ctx)
	defer cleanup()

	ctx = query.WithMaxDepth(ctx, nil)

	item, err := query.Insert(
		ctx,
		tx,
		ExecutionTableWithSchema,
		columns,
		nil,
		false,
		false,
		ExecutionTableColumns,
		values...,
	)
	if err != nil {
		return fmt.Errorf("failed to insert %#+v; %v", m, err)
	}
	v := (*item)[ExecutionTableIDColumn]

	if v == nil {
		return fmt.Errorf("failed to find %v in %#+v", ExecutionTableIDColumn, item)
	}

	wrapError := func(err error) error {
		return fmt.Errorf(
			"failed to treat %v: %#+v as uuid.UUID: %v",
			ExecutionTableIDColumn,
			(*item)[ExecutionTableIDColumn],
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

func (m *Execution) Update(ctx context.Context, tx pgx.Tx, setZeroValues bool, forceSetValuesForFields ...string) error {
	columns := make([]string, 0)
	values := make([]any, 0)

	if setZeroValues || !types.IsZeroTime(m.CreatedAt) || slices.Contains(forceSetValuesForFields, ExecutionTableCreatedAtColumn) {
		columns = append(columns, ExecutionTableCreatedAtColumn)

		v, err := types.FormatTime(m.CreatedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.CreatedAt; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.UpdatedAt) || slices.Contains(forceSetValuesForFields, ExecutionTableUpdatedAtColumn) {
		columns = append(columns, ExecutionTableUpdatedAtColumn)

		v, err := types.FormatTime(m.UpdatedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.UpdatedAt; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.DeletedAt) || slices.Contains(forceSetValuesForFields, ExecutionTableDeletedAtColumn) {
		columns = append(columns, ExecutionTableDeletedAtColumn)

		v, err := types.FormatTime(m.DeletedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.DeletedAt; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroString(m.Status) || slices.Contains(forceSetValuesForFields, ExecutionTableStatusColumn) {
		columns = append(columns, ExecutionTableStatusColumn)

		v, err := types.FormatString(m.Status)
		if err != nil {
			return fmt.Errorf("failed to handle m.Status; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.StartedAt) || slices.Contains(forceSetValuesForFields, ExecutionTableStartedAtColumn) {
		columns = append(columns, ExecutionTableStartedAtColumn)

		v, err := types.FormatTime(m.StartedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.StartedAt; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.EndedAt) || slices.Contains(forceSetValuesForFields, ExecutionTableEndedAtColumn) {
		columns = append(columns, ExecutionTableEndedAtColumn)

		v, err := types.FormatTime(m.EndedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.EndedAt; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.JobExecutorClaimedUntil) || slices.Contains(forceSetValuesForFields, ExecutionTableJobExecutorClaimedUntilColumn) {
		columns = append(columns, ExecutionTableJobExecutorClaimedUntilColumn)

		v, err := types.FormatTime(m.JobExecutorClaimedUntil)
		if err != nil {
			return fmt.Errorf("failed to handle m.JobExecutorClaimedUntil; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroUUID(m.TaskID) || slices.Contains(forceSetValuesForFields, ExecutionTableTaskIDColumn) {
		columns = append(columns, ExecutionTableTaskIDColumn)

		v, err := types.FormatUUID(m.TaskID)
		if err != nil {
			return fmt.Errorf("failed to handle m.TaskID; %v", err)
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
		ExecutionTableWithSchema,
		columns,
		fmt.Sprintf("%v = $$??", ExecutionTableIDColumn),
		ExecutionTableColumns,
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

func (m *Execution) Delete(ctx context.Context, tx pgx.Tx, hardDeletes ...bool) error {
	hardDelete := false
	if len(hardDeletes) > 0 {
		hardDelete = hardDeletes[0]
	}

	if !hardDelete && slices.Contains(ExecutionTableColumns, "deleted_at") {
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
		ExecutionTableWithSchema,
		fmt.Sprintf("%v = $$??", ExecutionTableIDColumn),
		values...,
	)
	if err != nil {
		return fmt.Errorf("failed to delete %#+v; %v", m, err)
	}

	_ = m.Reload(ctx, tx, true)

	return nil
}

func (m *Execution) LockTable(ctx context.Context, tx pgx.Tx, timeouts ...time.Duration) error {
	return query.LockTable(ctx, tx, ExecutionTableWithSchema, timeouts...)
}

func (m *Execution) LockTableWithRetries(ctx context.Context, tx pgx.Tx, overallTimeout time.Duration, individualAttempttimeout time.Duration) error {
	return query.LockTableWithRetries(ctx, tx, ExecutionTableWithSchema, overallTimeout, individualAttempttimeout)
}

func (m *Execution) AdvisoryLock(ctx context.Context, tx pgx.Tx, key int32, timeouts ...time.Duration) error {
	return query.AdvisoryLock(ctx, tx, ExecutionTableNamespaceID, key, timeouts...)
}

func (m *Execution) AdvisoryLockWithRetries(ctx context.Context, tx pgx.Tx, key int32, overallTimeout time.Duration, individualAttempttimeout time.Duration) error {
	return query.AdvisoryLockWithRetries(ctx, tx, ExecutionTableNamespaceID, key, overallTimeout, individualAttempttimeout)
}

func (m *Execution) JobExecutorClaim(ctx context.Context, tx pgx.Tx, until time.Time, timeout time.Duration) error {
	err := m.AdvisoryLockWithRetries(ctx, tx, math.MinInt32, timeout, time.Second*1)
	if err != nil {
		return fmt.Errorf("failed to claim (advisory lock): %s", err.Error())
	}

	_, _, _, _, _, err = SelectExecution(
		ctx,
		tx,
		fmt.Sprintf(
			"%s = $$?? AND (job_executor_claimed_until IS null OR job_executor_claimed_until < now())",
			ExecutionTablePrimaryKeyColumn,
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

func SelectExecutions(ctx context.Context, tx pgx.Tx, where string, orderBy *string, limit *int, offset *int, values ...any) ([]*Execution, int64, int64, int64, int64, error) {
	before := time.Now()

	if config.Debug() {
		log.Printf("entered SelectExecutions")

		defer func() {
			log.Printf("exited SelectExecutions in %s", time.Since(before))
		}()
	}
	if slices.Contains(ExecutionTableColumns, "deleted_at") {
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

	shouldLoad := query.ShouldLoad(ctx, ExecutionTable) || query.ShouldLoad(ctx, fmt.Sprintf("referenced_by_%s", ExecutionTable))

	var ok bool
	ctx, ok = query.HandleQueryPathGraphCycles(ctx, fmt.Sprintf("%s{%v}", ExecutionTable, nil), !isLoadQuery)
	if !ok && !shouldLoad {
		if config.Debug() {
			log.Printf("skipping SelectExecution early (query.ShouldLoad(): %v, query.HandleQueryPathGraphCycles(): %v)", shouldLoad, ok)
		}
		return []*Execution{}, 0, 0, 0, 0, nil
	}

	items, count, totalCount, page, totalPages, err := query.Select(
		ctx,
		tx,
		ExecutionTableColumnsWithTypeCasts,
		ExecutionTableWithSchema,
		where,
		orderBy,
		limit,
		offset,
		values...,
	)
	if err != nil {
		return nil, 0, 0, 0, 0, fmt.Errorf("failed to call SelectExecutions; %v", err)
	}

	objects := make([]*Execution, 0)

	for _, item := range *items {
		object := &Execution{}

		err = object.FromItem(item)
		if err != nil {
			return nil, 0, 0, 0, 0, err
		}

		if !types.IsZeroUUID(object.TaskID) {
			ctx, ok := query.HandleQueryPathGraphCycles(ctx, fmt.Sprintf("%s{%v}", TaskTable, object.TaskID), true)
			shouldLoad := query.ShouldLoad(ctx, TaskTable)
			if ok || shouldLoad {
				thisBefore := time.Now()

				if config.Debug() {
					log.Printf("loading SelectExecutions->SelectTask for object.TaskIDObject{%s: %v}", TaskTablePrimaryKeyColumn, object.TaskID)
				}

				object.TaskIDObject, _, _, _, _, err = SelectTask(
					ctx,
					tx,
					fmt.Sprintf("%v = $1", TaskTablePrimaryKeyColumn),
					object.TaskID,
				)
				if err != nil {
					if !errors.Is(err, sql.ErrNoRows) {
						return nil, 0, 0, 0, 0, err
					}
				}

				if config.Debug() {
					log.Printf("loaded SelectExecutions->SelectTask for object.TaskIDObject in %s", time.Since(thisBefore))
				}
			}
		}

		err = func() error {
			shouldLoad := query.ShouldLoad(ctx, fmt.Sprintf("referenced_by_%s", M2mRuleTriggerJobTable))
			ctx, ok := query.HandleQueryPathGraphCycles(ctx, fmt.Sprintf("__ReferencedBy__%s{%v}", M2mRuleTriggerJobTable, object.GetPrimaryKeyValue()), true)
			if ok || shouldLoad {
				thisBefore := time.Now()

				if config.Debug() {
					log.Printf("loading SelectExecutions->SelectM2mRuleTriggerJobs for object.ReferencedByM2mRuleTriggerJobExecutionIDObjects")
				}

				object.ReferencedByM2mRuleTriggerJobExecutionIDObjects, _, _, _, _, err = SelectM2mRuleTriggerJobs(
					ctx,
					tx,
					fmt.Sprintf("%v = $1", M2mRuleTriggerJobTableExecutionIDColumn),
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
					log.Printf("loaded SelectExecutions->SelectM2mRuleTriggerJobs for object.ReferencedByM2mRuleTriggerJobExecutionIDObjects in %s", time.Since(thisBefore))
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

func SelectExecution(ctx context.Context, tx pgx.Tx, where string, values ...any) (*Execution, int64, int64, int64, int64, error) {
	ctx, cleanup := query.WithQueryID(ctx)
	defer cleanup()

	ctx = query.WithMaxDepth(ctx, nil)

	objects, _, _, _, _, err := SelectExecutions(
		ctx,
		tx,
		where,
		nil,
		helpers.Ptr(2),
		helpers.Ptr(0),
		values...,
	)
	if err != nil {
		return nil, 0, 0, 0, 0, fmt.Errorf("failed to call SelectExecution; %v", err)
	}

	if len(objects) > 1 {
		return nil, 0, 0, 0, 0, fmt.Errorf("attempt to call SelectExecution returned more than 1 row")
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

func JobExecutorClaimExecution(ctx context.Context, tx pgx.Tx, until time.Time, timeout time.Duration, wheres ...string) (*Execution, error) {
	m := &Execution{}

	err := m.AdvisoryLockWithRetries(ctx, tx, math.MinInt32, timeout, time.Second*1)
	if err != nil {
		return nil, fmt.Errorf("failed to claim: %s", err.Error())
	}

	extraWhere := ""
	if len(wheres) > 0 {
		extraWhere = fmt.Sprintf("AND %s", extraWhere)
	}

	ms, _, _, _, _, err := SelectExecutions(
		ctx,
		tx,
		fmt.Sprintf(
			"(job_executor_claimed_until IS null OR job_executor_claimed_until < now())%s",
			extraWhere,
		),
		helpers.Ptr(
			"job_executor_claimed_until ASC",
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

	m.JobExecutorClaimedUntil = until

	err = m.Update(ctx, tx, false)
	if err != nil {
		return nil, fmt.Errorf("failed to claim: %s", err.Error())
	}

	return m, nil
}

func handleGetExecutions(arguments *server.SelectManyArguments, db *pgxpool.Pool) ([]*Execution, int64, int64, int64, int64, error) {
	tx, err := db.Begin(arguments.Ctx)
	if err != nil {
		return nil, 0, 0, 0, 0, err
	}

	defer func() {
		_ = tx.Rollback(arguments.Ctx)
	}()

	objects, count, totalCount, page, totalPages, err := SelectExecutions(arguments.Ctx, tx, arguments.Where, arguments.OrderBy, arguments.Limit, arguments.Offset, arguments.Values...)
	if err != nil {
		return nil, 0, 0, 0, 0, err
	}

	err = tx.Commit(arguments.Ctx)
	if err != nil {
		return nil, 0, 0, 0, 0, err
	}

	return objects, count, totalCount, page, totalPages, nil
}

func handleGetExecution(arguments *server.SelectOneArguments, db *pgxpool.Pool, primaryKey uuid.UUID) ([]*Execution, int64, int64, int64, int64, error) {
	tx, err := db.Begin(arguments.Ctx)
	if err != nil {
		return nil, 0, 0, 0, 0, err
	}

	defer func() {
		_ = tx.Rollback(arguments.Ctx)
	}()

	object, count, totalCount, page, totalPages, err := SelectExecution(arguments.Ctx, tx, arguments.Where, arguments.Values...)
	if err != nil {
		return nil, 0, 0, 0, 0, err
	}

	err = tx.Commit(arguments.Ctx)
	if err != nil {
		return nil, 0, 0, 0, 0, err
	}

	return []*Execution{object}, count, totalCount, page, totalPages, nil
}

func handlePostExecution(arguments *server.LoadArguments, db *pgxpool.Pool, waitForChange server.WaitForChange, objects []*Execution, forceSetValuesForFieldsByObjectIndex [][]string) ([]*Execution, int64, int64, int64, int64, error) {
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
		_, err := waitForChange(arguments.Ctx, []stream.Action{stream.INSERT}, ExecutionTable, xid)
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

func handlePutExecution(arguments *server.LoadArguments, db *pgxpool.Pool, waitForChange server.WaitForChange, object *Execution) ([]*Execution, int64, int64, int64, int64, error) {
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
		_, err := waitForChange(arguments.Ctx, []stream.Action{stream.UPDATE, stream.SOFT_DELETE, stream.SOFT_RESTORE, stream.SOFT_UPDATE}, ExecutionTable, xid)
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

	return []*Execution{object}, count, totalCount, page, totalPages, nil
}

func handlePatchExecution(arguments *server.LoadArguments, db *pgxpool.Pool, waitForChange server.WaitForChange, object *Execution, forceSetValuesForFields []string) ([]*Execution, int64, int64, int64, int64, error) {
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
		_, err := waitForChange(arguments.Ctx, []stream.Action{stream.UPDATE, stream.SOFT_DELETE, stream.SOFT_RESTORE, stream.SOFT_UPDATE}, ExecutionTable, xid)
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

	return []*Execution{object}, count, totalCount, page, totalPages, nil
}

func handleDeleteExecution(arguments *server.LoadArguments, db *pgxpool.Pool, waitForChange server.WaitForChange, object *Execution) error {
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
		_, err := waitForChange(arguments.Ctx, []stream.Action{stream.DELETE, stream.SOFT_DELETE}, ExecutionTable, xid)
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

func MutateRouterForExecution(r chi.Router, db *pgxpool.Pool, redisPool *redis.Pool, objectMiddlewares []server.ObjectMiddleware, waitForChange server.WaitForChange) {

	func() {
		postHandlerForJobExecutorClaim, err := getHTTPHandler(
			http.MethodPost,
			"/job-executor-claim-execution",
			http.StatusOK,
			func(
				ctx context.Context,
				pathParams server.EmptyPathParams,
				queryParams server.EmptyQueryParams,
				req ExecutionJobExecutorClaimRequest,
				rawReq any,
			) (server.Response[Execution], error) {
				tx, err := db.Begin(ctx)
				if err != nil {
					return server.Response[Execution]{}, err
				}

				defer func() {
					_ = tx.Rollback(ctx)
				}()

				object, err := JobExecutorClaimExecution(ctx, tx, req.Until, time.Millisecond*time.Duration(req.TimeoutSeconds*1000))
				if err != nil {
					return server.Response[Execution]{}, err
				}

				count := int64(0)

				totalCount := int64(0)

				limit := int64(0)

				offset := int64(0)

				if object == nil {
					return server.Response[Execution]{
						Status:     http.StatusOK,
						Success:    true,
						Error:      nil,
						Objects:    []*Execution{},
						Count:      count,
						TotalCount: totalCount,
						Limit:      limit,
						Offset:     offset,
					}, nil
				}

				err = tx.Commit(ctx)
				if err != nil {
					return server.Response[Execution]{}, err
				}

				return server.Response[Execution]{
					Status:     http.StatusOK,
					Success:    true,
					Error:      nil,
					Objects:    []*Execution{object},
					Count:      count,
					TotalCount: totalCount,
					Limit:      limit,
					Offset:     offset,
				}, nil
			},
			Execution{},
			ExecutionIntrospectedTable,
		)
		if err != nil {
			panic(err)
		}
		r.Post(postHandlerForJobExecutorClaim.FullPath, postHandlerForJobExecutorClaim.ServeHTTP)

		postHandlerForJobExecutorClaimOne, err := getHTTPHandler(
			http.MethodPost,
			"/executions/{primaryKey}/job-executor-claim",
			http.StatusOK,
			func(
				ctx context.Context,
				pathParams ExecutionOnePathParams,
				queryParams ExecutionLoadQueryParams,
				req ExecutionJobExecutorClaimRequest,
				rawReq any,
			) (server.Response[Execution], error) {
				before := time.Now()

				redisConn := redisPool.Get()
				defer func() {
					_ = redisConn.Close()
				}()

				arguments, err := server.GetSelectOneArguments(ctx, queryParams.Depth, ExecutionIntrospectedTable, pathParams.PrimaryKey, nil, nil)
				if err != nil {
					if config.Debug() {
						log.Printf("request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return server.Response[Execution]{}, err
				}

				/* note: deliberately no attempt at a cache hit */

				var object *Execution
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

					object, count, totalCount, _, _, err = SelectExecution(arguments.Ctx, tx, arguments.Where, arguments.Values...)
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

					return server.Response[Execution]{}, err
				}

				limit := int64(0)

				offset := int64(0)

				response := server.Response[Execution]{
					Status:     http.StatusOK,
					Success:    true,
					Error:      nil,
					Objects:    []*Execution{object},
					Count:      count,
					TotalCount: totalCount,
					Limit:      limit,
					Offset:     offset,
				}

				return response, nil
			},
			Execution{},
			ExecutionIntrospectedTable,
		)
		if err != nil {
			panic(err)
		}
		r.Post(postHandlerForJobExecutorClaimOne.FullPath, postHandlerForJobExecutorClaimOne.ServeHTTP)
	}()

	func() {
		getManyHandler, err := getHTTPHandler(
			http.MethodGet,
			"/executions",
			http.StatusOK,
			func(
				ctx context.Context,
				pathParams server.EmptyPathParams,
				queryParams map[string]any,
				req server.EmptyRequest,
				rawReq any,
			) (server.Response[Execution], error) {
				before := time.Now()

				redisConn := redisPool.Get()
				defer func() {
					_ = redisConn.Close()
				}()

				arguments, err := server.GetSelectManyArguments(ctx, queryParams, ExecutionIntrospectedTable, nil, nil)
				if err != nil {
					if config.Debug() {
						log.Printf("request cache not yet reached; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return server.Response[Execution]{}, err
				}

				cachedResponseAsJSON, cacheHit, err := server.GetCachedResponseAsJSON(arguments.RequestHash, redisConn)
				if err != nil {
					if config.Debug() {
						log.Printf("request cache failed; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return server.Response[Execution]{}, err
				}

				if cacheHit {
					var cachedResponse server.Response[Execution]

					/* TODO: it'd be nice to be able to avoid this (i.e. just pass straight through) */
					err = json.Unmarshal(cachedResponseAsJSON, &cachedResponse)
					if err != nil {
						if config.Debug() {
							log.Printf("request cache hit but failed unmarshal; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
						}

						return server.Response[Execution]{}, err
					}

					if config.Debug() {
						log.Printf("request cache hit; request succeeded in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return cachedResponse, nil
				}

				objects, count, totalCount, _, _, err := handleGetExecutions(arguments, db)
				if err != nil {
					if config.Debug() {
						log.Printf("request cache missed; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return server.Response[Execution]{}, err
				}

				limit := int64(0)
				if arguments.Limit != nil {
					limit = int64(*arguments.Limit)
				}

				offset := int64(0)
				if arguments.Offset != nil {
					offset = int64(*arguments.Offset)
				}

				response := server.Response[Execution]{
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

					return server.Response[Execution]{}, err
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
			Execution{},
			ExecutionIntrospectedTable,
		)
		if err != nil {
			panic(err)
		}
		r.Get(getManyHandler.FullPath, getManyHandler.ServeHTTP)
	}()

	func() {
		getOneHandler, err := getHTTPHandler(
			http.MethodGet,
			"/executions/{primaryKey}",
			http.StatusOK,
			func(
				ctx context.Context,
				pathParams ExecutionOnePathParams,
				queryParams ExecutionLoadQueryParams,
				req server.EmptyRequest,
				rawReq any,
			) (server.Response[Execution], error) {
				before := time.Now()

				redisConn := redisPool.Get()
				defer func() {
					_ = redisConn.Close()
				}()

				arguments, err := server.GetSelectOneArguments(ctx, queryParams.Depth, ExecutionIntrospectedTable, pathParams.PrimaryKey, nil, nil)
				if err != nil {
					if config.Debug() {
						log.Printf("request cache not yet reached; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return server.Response[Execution]{}, err
				}

				cachedResponseAsJSON, cacheHit, err := server.GetCachedResponseAsJSON(arguments.RequestHash, redisConn)
				if err != nil {
					if config.Debug() {
						log.Printf("request cache failed; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return server.Response[Execution]{}, err
				}

				if cacheHit {
					var cachedResponse server.Response[Execution]

					/* TODO: it'd be nice to be able to avoid this (i.e. just pass straight through) */
					err = json.Unmarshal(cachedResponseAsJSON, &cachedResponse)
					if err != nil {
						if config.Debug() {
							log.Printf("request cache hit but failed unmarshal; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
						}

						return server.Response[Execution]{}, err
					}

					if config.Debug() {
						log.Printf("request cache hit; request succeeded in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return cachedResponse, nil
				}

				objects, count, totalCount, _, _, err := handleGetExecution(arguments, db, pathParams.PrimaryKey)
				if err != nil {
					if config.Debug() {
						log.Printf("request cache missed; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return server.Response[Execution]{}, err
				}

				limit := int64(0)

				offset := int64(0)

				response := server.Response[Execution]{
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

					return server.Response[Execution]{}, err
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
			Execution{},
			ExecutionIntrospectedTable,
		)
		if err != nil {
			panic(err)
		}
		r.Get(getOneHandler.FullPath, getOneHandler.ServeHTTP)
	}()

	func() {
		postHandler, err := getHTTPHandler(
			http.MethodPost,
			"/executions",
			http.StatusCreated,
			func(
				ctx context.Context,
				pathParams server.EmptyPathParams,
				queryParams ExecutionLoadQueryParams,
				req []*Execution,
				rawReq any,
			) (server.Response[Execution], error) {
				allRawItems, ok := rawReq.([]any)
				if !ok {
					return server.Response[Execution]{}, fmt.Errorf("failed to cast %#+v to []map[string]any", rawReq)
				}

				allItems := make([]map[string]any, 0)
				for _, rawItem := range allRawItems {
					item, ok := rawItem.(map[string]any)
					if !ok {
						return server.Response[Execution]{}, fmt.Errorf("failed to cast %#+v to map[string]any", rawItem)
					}

					allItems = append(allItems, item)
				}

				forceSetValuesForFieldsByObjectIndex := make([][]string, 0)
				for _, item := range allItems {
					forceSetValuesForFields := make([]string, 0)
					for _, possibleField := range maps.Keys(item) {
						if !slices.Contains(ExecutionTableColumns, possibleField) {
							continue
						}

						forceSetValuesForFields = append(forceSetValuesForFields, possibleField)
					}
					forceSetValuesForFieldsByObjectIndex = append(forceSetValuesForFieldsByObjectIndex, forceSetValuesForFields)
				}

				arguments, err := server.GetLoadArguments(ctx, queryParams.Depth)
				if err != nil {
					return server.Response[Execution]{}, err
				}

				objects, count, totalCount, _, _, err := handlePostExecution(arguments, db, waitForChange, req, forceSetValuesForFieldsByObjectIndex)
				if err != nil {
					return server.Response[Execution]{}, err
				}

				limit := int64(0)

				offset := int64(0)

				return server.Response[Execution]{
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
			Execution{},
			ExecutionIntrospectedTable,
		)
		if err != nil {
			panic(err)
		}
		r.Post(postHandler.FullPath, postHandler.ServeHTTP)
	}()

	func() {
		putHandler, err := getHTTPHandler(
			http.MethodPatch,
			"/executions/{primaryKey}",
			http.StatusOK,
			func(
				ctx context.Context,
				pathParams ExecutionOnePathParams,
				queryParams ExecutionLoadQueryParams,
				req Execution,
				rawReq any,
			) (server.Response[Execution], error) {
				item, ok := rawReq.(map[string]any)
				if !ok {
					return server.Response[Execution]{}, fmt.Errorf("failed to cast %#+v to map[string]any", item)
				}

				arguments, err := server.GetLoadArguments(ctx, queryParams.Depth)
				if err != nil {
					return server.Response[Execution]{}, err
				}

				object := &req
				object.ID = pathParams.PrimaryKey

				objects, count, totalCount, _, _, err := handlePutExecution(arguments, db, waitForChange, object)
				if err != nil {
					return server.Response[Execution]{}, err
				}

				limit := int64(0)

				offset := int64(0)

				return server.Response[Execution]{
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
			Execution{},
			ExecutionIntrospectedTable,
		)
		if err != nil {
			panic(err)
		}
		r.Put(putHandler.FullPath, putHandler.ServeHTTP)
	}()

	func() {
		patchHandler, err := getHTTPHandler(
			http.MethodPatch,
			"/executions/{primaryKey}",
			http.StatusOK,
			func(
				ctx context.Context,
				pathParams ExecutionOnePathParams,
				queryParams ExecutionLoadQueryParams,
				req Execution,
				rawReq any,
			) (server.Response[Execution], error) {
				item, ok := rawReq.(map[string]any)
				if !ok {
					return server.Response[Execution]{}, fmt.Errorf("failed to cast %#+v to map[string]any", item)
				}

				forceSetValuesForFields := make([]string, 0)
				for _, possibleField := range maps.Keys(item) {
					if !slices.Contains(ExecutionTableColumns, possibleField) {
						continue
					}

					forceSetValuesForFields = append(forceSetValuesForFields, possibleField)
				}

				arguments, err := server.GetLoadArguments(ctx, queryParams.Depth)
				if err != nil {
					return server.Response[Execution]{}, err
				}

				object := &req
				object.ID = pathParams.PrimaryKey

				objects, count, totalCount, _, _, err := handlePatchExecution(arguments, db, waitForChange, object, forceSetValuesForFields)
				if err != nil {
					return server.Response[Execution]{}, err
				}

				limit := int64(0)

				offset := int64(0)

				return server.Response[Execution]{
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
			Execution{},
			ExecutionIntrospectedTable,
		)
		if err != nil {
			panic(err)
		}
		r.Patch(patchHandler.FullPath, patchHandler.ServeHTTP)
	}()

	func() {
		deleteHandler, err := getHTTPHandler(
			http.MethodDelete,
			"/executions/{primaryKey}",
			http.StatusNoContent,
			func(
				ctx context.Context,
				pathParams ExecutionOnePathParams,
				queryParams ExecutionLoadQueryParams,
				req server.EmptyRequest,
				rawReq any,
			) (server.EmptyResponse, error) {
				arguments, err := server.GetLoadArguments(ctx, queryParams.Depth)
				if err != nil {
					return server.EmptyResponse{}, err
				}

				object := &Execution{}
				object.ID = pathParams.PrimaryKey

				err = handleDeleteExecution(arguments, db, waitForChange, object)
				if err != nil {
					return server.EmptyResponse{}, err
				}

				return server.EmptyResponse{}, nil
			},
			Execution{},
			ExecutionIntrospectedTable,
		)
		if err != nil {
			panic(err)
		}
		r.Delete(deleteHandler.FullPath, deleteHandler.ServeHTTP)
	}()
}

func NewExecutionFromItem(item map[string]any) (any, error) {
	object := &Execution{}

	err := object.FromItem(item)
	if err != nil {
		return nil, err
	}

	return object, nil
}

func init() {
	register(
		ExecutionTable,
		Execution{},
		NewExecutionFromItem,
		"/executions",
		MutateRouterForExecution,
	)
}
