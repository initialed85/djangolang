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

type M2mRuleTriggerJob struct {
	ID                                              uuid.UUID    `json:"id"`
	CreatedAt                                       time.Time    `json:"created_at"`
	UpdatedAt                                       time.Time    `json:"updated_at"`
	DeletedAt                                       *time.Time   `json:"deleted_at"`
	ExecutionsProducedAt                            *time.Time   `json:"executions_produced_at"`
	JobExecutorClaimedUntil                         time.Time    `json:"job_executor_claimed_until"`
	JobID                                           uuid.UUID    `json:"job_id"`
	JobIDObject                                     *Job         `json:"job_id_object"`
	RuleID                                          uuid.UUID    `json:"rule_id"`
	RuleIDObject                                    *Rule        `json:"rule_id_object"`
	ReferencedByExecutionM2mRuleTriggerJobIDObjects []*Execution `json:"referenced_by_execution_m2m_rule_trigger_job_id_objects"`
}

var M2mRuleTriggerJobTable = "m2m_rule_trigger_job"

var M2mRuleTriggerJobTableWithSchema = fmt.Sprintf("%s.%s", schema, M2mRuleTriggerJobTable)

var M2mRuleTriggerJobTableNamespaceID int32 = 1337 + 5

var (
	M2mRuleTriggerJobTableIDColumn                      = "id"
	M2mRuleTriggerJobTableCreatedAtColumn               = "created_at"
	M2mRuleTriggerJobTableUpdatedAtColumn               = "updated_at"
	M2mRuleTriggerJobTableDeletedAtColumn               = "deleted_at"
	M2mRuleTriggerJobTableExecutionsProducedAtColumn    = "executions_produced_at"
	M2mRuleTriggerJobTableJobExecutorClaimedUntilColumn = "job_executor_claimed_until"
	M2mRuleTriggerJobTableJobIDColumn                   = "job_id"
	M2mRuleTriggerJobTableRuleIDColumn                  = "rule_id"
)

var (
	M2mRuleTriggerJobTableIDColumnWithTypeCast                      = `"id" AS id`
	M2mRuleTriggerJobTableCreatedAtColumnWithTypeCast               = `"created_at" AS created_at`
	M2mRuleTriggerJobTableUpdatedAtColumnWithTypeCast               = `"updated_at" AS updated_at`
	M2mRuleTriggerJobTableDeletedAtColumnWithTypeCast               = `"deleted_at" AS deleted_at`
	M2mRuleTriggerJobTableExecutionsProducedAtColumnWithTypeCast    = `"executions_produced_at" AS executions_produced_at`
	M2mRuleTriggerJobTableJobExecutorClaimedUntilColumnWithTypeCast = `"job_executor_claimed_until" AS job_executor_claimed_until`
	M2mRuleTriggerJobTableJobIDColumnWithTypeCast                   = `"job_id" AS job_id`
	M2mRuleTriggerJobTableRuleIDColumnWithTypeCast                  = `"rule_id" AS rule_id`
)

var M2mRuleTriggerJobTableColumns = []string{
	M2mRuleTriggerJobTableIDColumn,
	M2mRuleTriggerJobTableCreatedAtColumn,
	M2mRuleTriggerJobTableUpdatedAtColumn,
	M2mRuleTriggerJobTableDeletedAtColumn,
	M2mRuleTriggerJobTableExecutionsProducedAtColumn,
	M2mRuleTriggerJobTableJobExecutorClaimedUntilColumn,
	M2mRuleTriggerJobTableJobIDColumn,
	M2mRuleTriggerJobTableRuleIDColumn,
}

var M2mRuleTriggerJobTableColumnsWithTypeCasts = []string{
	M2mRuleTriggerJobTableIDColumnWithTypeCast,
	M2mRuleTriggerJobTableCreatedAtColumnWithTypeCast,
	M2mRuleTriggerJobTableUpdatedAtColumnWithTypeCast,
	M2mRuleTriggerJobTableDeletedAtColumnWithTypeCast,
	M2mRuleTriggerJobTableExecutionsProducedAtColumnWithTypeCast,
	M2mRuleTriggerJobTableJobExecutorClaimedUntilColumnWithTypeCast,
	M2mRuleTriggerJobTableJobIDColumnWithTypeCast,
	M2mRuleTriggerJobTableRuleIDColumnWithTypeCast,
}

var M2mRuleTriggerJobIntrospectedTable *introspect.Table

var M2mRuleTriggerJobTableColumnLookup map[string]*introspect.Column

var (
	M2mRuleTriggerJobTablePrimaryKeyColumn = M2mRuleTriggerJobTableIDColumn
)

func init() {
	M2mRuleTriggerJobIntrospectedTable = tableByName[M2mRuleTriggerJobTable]

	/* only needed during templating */
	if M2mRuleTriggerJobIntrospectedTable == nil {
		M2mRuleTriggerJobIntrospectedTable = &introspect.Table{}
	}

	M2mRuleTriggerJobTableColumnLookup = M2mRuleTriggerJobIntrospectedTable.ColumnByName
}

type M2mRuleTriggerJobOnePathParams struct {
	PrimaryKey uuid.UUID `json:"primaryKey"`
}

type M2mRuleTriggerJobLoadQueryParams struct {
	Depth *int `json:"depth"`
}

type M2mRuleTriggerJobJobExecutorClaimRequest struct {
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

func (m *M2mRuleTriggerJob) GetPrimaryKeyColumn() string {
	return M2mRuleTriggerJobTablePrimaryKeyColumn
}

func (m *M2mRuleTriggerJob) GetPrimaryKeyValue() any {
	return m.ID
}

func (m *M2mRuleTriggerJob) FromItem(item map[string]any) error {
	if item == nil {
		return fmt.Errorf(
			"item unexpectedly nil during M2mRuleTriggerJobFromItem",
		)
	}

	if len(item) == 0 {
		return fmt.Errorf(
			"item unexpectedly empty during M2mRuleTriggerJobFromItem",
		)
	}

	wrapError := func(k string, v any, err error) error {
		return fmt.Errorf("%v: %#+v; error; %v", k, v, err)
	}

	for k, v := range item {
		_, ok := M2mRuleTriggerJobTableColumnLookup[k]
		if !ok {
			return fmt.Errorf(
				"item contained unexpected key %#+v during M2mRuleTriggerJobFromItem; item: %#+v",
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

		case "executions_produced_at":
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
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uuexecutions_produced_at.UUID", temp1))
				}
			}

			m.ExecutionsProducedAt = &temp2

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

		}
	}

	return nil
}

func (m *M2mRuleTriggerJob) Reload(ctx context.Context, tx pgx.Tx, includeDeleteds ...bool) error {
	extraWhere := ""
	if len(includeDeleteds) > 0 && includeDeleteds[0] {
		if slices.Contains(M2mRuleTriggerJobTableColumns, "deleted_at") {
			extraWhere = "\n    AND (deleted_at IS null OR deleted_at IS NOT null)"
		}
	}

	ctx, cleanup := query.WithQueryID(ctx)
	defer cleanup()

	ctx = query.WithMaxDepth(ctx, nil)

	o, _, _, _, _, err := SelectM2mRuleTriggerJob(
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
	m.ExecutionsProducedAt = o.ExecutionsProducedAt
	m.JobExecutorClaimedUntil = o.JobExecutorClaimedUntil
	m.JobID = o.JobID
	m.JobIDObject = o.JobIDObject
	m.RuleID = o.RuleID
	m.RuleIDObject = o.RuleIDObject
	m.ReferencedByExecutionM2mRuleTriggerJobIDObjects = o.ReferencedByExecutionM2mRuleTriggerJobIDObjects

	return nil
}

func (m *M2mRuleTriggerJob) Insert(ctx context.Context, tx pgx.Tx, setPrimaryKey bool, setZeroValues bool, forceSetValuesForFields ...string) error {
	columns := make([]string, 0)
	values := make([]any, 0)

	if setPrimaryKey && (setZeroValues || !types.IsZeroUUID(m.ID) || slices.Contains(forceSetValuesForFields, M2mRuleTriggerJobTableIDColumn) || isRequired(M2mRuleTriggerJobTableColumnLookup, M2mRuleTriggerJobTableIDColumn)) {
		columns = append(columns, M2mRuleTriggerJobTableIDColumn)

		v, err := types.FormatUUID(m.ID)
		if err != nil {
			return fmt.Errorf("failed to handle m.ID; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.CreatedAt) || slices.Contains(forceSetValuesForFields, M2mRuleTriggerJobTableCreatedAtColumn) || isRequired(M2mRuleTriggerJobTableColumnLookup, M2mRuleTriggerJobTableCreatedAtColumn) {
		columns = append(columns, M2mRuleTriggerJobTableCreatedAtColumn)

		v, err := types.FormatTime(m.CreatedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.CreatedAt; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.UpdatedAt) || slices.Contains(forceSetValuesForFields, M2mRuleTriggerJobTableUpdatedAtColumn) || isRequired(M2mRuleTriggerJobTableColumnLookup, M2mRuleTriggerJobTableUpdatedAtColumn) {
		columns = append(columns, M2mRuleTriggerJobTableUpdatedAtColumn)

		v, err := types.FormatTime(m.UpdatedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.UpdatedAt; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.DeletedAt) || slices.Contains(forceSetValuesForFields, M2mRuleTriggerJobTableDeletedAtColumn) || isRequired(M2mRuleTriggerJobTableColumnLookup, M2mRuleTriggerJobTableDeletedAtColumn) {
		columns = append(columns, M2mRuleTriggerJobTableDeletedAtColumn)

		v, err := types.FormatTime(m.DeletedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.DeletedAt; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.ExecutionsProducedAt) || slices.Contains(forceSetValuesForFields, M2mRuleTriggerJobTableExecutionsProducedAtColumn) || isRequired(M2mRuleTriggerJobTableColumnLookup, M2mRuleTriggerJobTableExecutionsProducedAtColumn) {
		columns = append(columns, M2mRuleTriggerJobTableExecutionsProducedAtColumn)

		v, err := types.FormatTime(m.ExecutionsProducedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.ExecutionsProducedAt; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.JobExecutorClaimedUntil) || slices.Contains(forceSetValuesForFields, M2mRuleTriggerJobTableJobExecutorClaimedUntilColumn) || isRequired(M2mRuleTriggerJobTableColumnLookup, M2mRuleTriggerJobTableJobExecutorClaimedUntilColumn) {
		columns = append(columns, M2mRuleTriggerJobTableJobExecutorClaimedUntilColumn)

		v, err := types.FormatTime(m.JobExecutorClaimedUntil)
		if err != nil {
			return fmt.Errorf("failed to handle m.JobExecutorClaimedUntil; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroUUID(m.JobID) || slices.Contains(forceSetValuesForFields, M2mRuleTriggerJobTableJobIDColumn) || isRequired(M2mRuleTriggerJobTableColumnLookup, M2mRuleTriggerJobTableJobIDColumn) {
		columns = append(columns, M2mRuleTriggerJobTableJobIDColumn)

		v, err := types.FormatUUID(m.JobID)
		if err != nil {
			return fmt.Errorf("failed to handle m.JobID; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroUUID(m.RuleID) || slices.Contains(forceSetValuesForFields, M2mRuleTriggerJobTableRuleIDColumn) || isRequired(M2mRuleTriggerJobTableColumnLookup, M2mRuleTriggerJobTableRuleIDColumn) {
		columns = append(columns, M2mRuleTriggerJobTableRuleIDColumn)

		v, err := types.FormatUUID(m.RuleID)
		if err != nil {
			return fmt.Errorf("failed to handle m.RuleID; %v", err)
		}

		values = append(values, v)
	}

	ctx, cleanup := query.WithQueryID(ctx)
	defer cleanup()

	ctx = query.WithMaxDepth(ctx, nil)

	item, err := query.Insert(
		ctx,
		tx,
		M2mRuleTriggerJobTableWithSchema,
		columns,
		nil,
		false,
		false,
		M2mRuleTriggerJobTableColumns,
		values...,
	)
	if err != nil {
		return fmt.Errorf("failed to insert %#+v; %v", m, err)
	}
	v := (*item)[M2mRuleTriggerJobTableIDColumn]

	if v == nil {
		return fmt.Errorf("failed to find %v in %#+v", M2mRuleTriggerJobTableIDColumn, item)
	}

	wrapError := func(err error) error {
		return fmt.Errorf(
			"failed to treat %v: %#+v as uuid.UUID: %v",
			M2mRuleTriggerJobTableIDColumn,
			(*item)[M2mRuleTriggerJobTableIDColumn],
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

func (m *M2mRuleTriggerJob) Update(ctx context.Context, tx pgx.Tx, setZeroValues bool, forceSetValuesForFields ...string) error {
	columns := make([]string, 0)
	values := make([]any, 0)

	if setZeroValues || !types.IsZeroTime(m.CreatedAt) || slices.Contains(forceSetValuesForFields, M2mRuleTriggerJobTableCreatedAtColumn) {
		columns = append(columns, M2mRuleTriggerJobTableCreatedAtColumn)

		v, err := types.FormatTime(m.CreatedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.CreatedAt; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.UpdatedAt) || slices.Contains(forceSetValuesForFields, M2mRuleTriggerJobTableUpdatedAtColumn) {
		columns = append(columns, M2mRuleTriggerJobTableUpdatedAtColumn)

		v, err := types.FormatTime(m.UpdatedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.UpdatedAt; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.DeletedAt) || slices.Contains(forceSetValuesForFields, M2mRuleTriggerJobTableDeletedAtColumn) {
		columns = append(columns, M2mRuleTriggerJobTableDeletedAtColumn)

		v, err := types.FormatTime(m.DeletedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.DeletedAt; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.ExecutionsProducedAt) || slices.Contains(forceSetValuesForFields, M2mRuleTriggerJobTableExecutionsProducedAtColumn) {
		columns = append(columns, M2mRuleTriggerJobTableExecutionsProducedAtColumn)

		v, err := types.FormatTime(m.ExecutionsProducedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.ExecutionsProducedAt; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.JobExecutorClaimedUntil) || slices.Contains(forceSetValuesForFields, M2mRuleTriggerJobTableJobExecutorClaimedUntilColumn) {
		columns = append(columns, M2mRuleTriggerJobTableJobExecutorClaimedUntilColumn)

		v, err := types.FormatTime(m.JobExecutorClaimedUntil)
		if err != nil {
			return fmt.Errorf("failed to handle m.JobExecutorClaimedUntil; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroUUID(m.JobID) || slices.Contains(forceSetValuesForFields, M2mRuleTriggerJobTableJobIDColumn) {
		columns = append(columns, M2mRuleTriggerJobTableJobIDColumn)

		v, err := types.FormatUUID(m.JobID)
		if err != nil {
			return fmt.Errorf("failed to handle m.JobID; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroUUID(m.RuleID) || slices.Contains(forceSetValuesForFields, M2mRuleTriggerJobTableRuleIDColumn) {
		columns = append(columns, M2mRuleTriggerJobTableRuleIDColumn)

		v, err := types.FormatUUID(m.RuleID)
		if err != nil {
			return fmt.Errorf("failed to handle m.RuleID; %v", err)
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
		M2mRuleTriggerJobTableWithSchema,
		columns,
		fmt.Sprintf("%v = $$??", M2mRuleTriggerJobTableIDColumn),
		M2mRuleTriggerJobTableColumns,
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

func (m *M2mRuleTriggerJob) Delete(ctx context.Context, tx pgx.Tx, hardDeletes ...bool) error {
	hardDelete := false
	if len(hardDeletes) > 0 {
		hardDelete = hardDeletes[0]
	}

	if !hardDelete && slices.Contains(M2mRuleTriggerJobTableColumns, "deleted_at") {
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
		M2mRuleTriggerJobTableWithSchema,
		fmt.Sprintf("%v = $$??", M2mRuleTriggerJobTableIDColumn),
		values...,
	)
	if err != nil {
		return fmt.Errorf("failed to delete %#+v; %v", m, err)
	}

	_ = m.Reload(ctx, tx, true)

	return nil
}

func (m *M2mRuleTriggerJob) LockTable(ctx context.Context, tx pgx.Tx, timeouts ...time.Duration) error {
	return query.LockTable(ctx, tx, M2mRuleTriggerJobTableWithSchema, timeouts...)
}

func (m *M2mRuleTriggerJob) LockTableWithRetries(ctx context.Context, tx pgx.Tx, overallTimeout time.Duration, individualAttempttimeout time.Duration) error {
	return query.LockTableWithRetries(ctx, tx, M2mRuleTriggerJobTableWithSchema, overallTimeout, individualAttempttimeout)
}

func (m *M2mRuleTriggerJob) AdvisoryLock(ctx context.Context, tx pgx.Tx, key int32, timeouts ...time.Duration) error {
	return query.AdvisoryLock(ctx, tx, M2mRuleTriggerJobTableNamespaceID, key, timeouts...)
}

func (m *M2mRuleTriggerJob) AdvisoryLockWithRetries(ctx context.Context, tx pgx.Tx, key int32, overallTimeout time.Duration, individualAttempttimeout time.Duration) error {
	return query.AdvisoryLockWithRetries(ctx, tx, M2mRuleTriggerJobTableNamespaceID, key, overallTimeout, individualAttempttimeout)
}

func (m *M2mRuleTriggerJob) JobExecutorClaim(ctx context.Context, tx pgx.Tx, until time.Time, timeout time.Duration) error {
	err := m.AdvisoryLockWithRetries(ctx, tx, math.MinInt32, timeout, time.Second*1)
	if err != nil {
		return fmt.Errorf("failed to claim (advisory lock): %s", err.Error())
	}

	_, _, _, _, _, err = SelectM2mRuleTriggerJob(
		ctx,
		tx,
		fmt.Sprintf(
			"%s = $$?? AND (job_executor_claimed_until IS null OR job_executor_claimed_until < now())",
			M2mRuleTriggerJobTablePrimaryKeyColumn,
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

func SelectM2mRuleTriggerJobs(ctx context.Context, tx pgx.Tx, where string, orderBy *string, limit *int, offset *int, values ...any) ([]*M2mRuleTriggerJob, int64, int64, int64, int64, error) {
	before := time.Now()

	if config.Debug() {
		log.Printf("entered SelectM2mRuleTriggerJobs")

		defer func() {
			log.Printf("exited SelectM2mRuleTriggerJobs in %s", time.Since(before))
		}()
	}
	if slices.Contains(M2mRuleTriggerJobTableColumns, "deleted_at") {
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

	shouldLoad := query.ShouldLoad(ctx, M2mRuleTriggerJobTable) || query.ShouldLoad(ctx, fmt.Sprintf("referenced_by_%s", M2mRuleTriggerJobTable))

	var ok bool
	ctx, ok = query.HandleQueryPathGraphCycles(ctx, fmt.Sprintf("%s{%v}", M2mRuleTriggerJobTable, nil), !isLoadQuery)
	if !ok && !shouldLoad {
		if config.Debug() {
			log.Printf("skipping SelectM2mRuleTriggerJob early (query.ShouldLoad(): %v, query.HandleQueryPathGraphCycles(): %v)", shouldLoad, ok)
		}
		return []*M2mRuleTriggerJob{}, 0, 0, 0, 0, nil
	}

	items, count, totalCount, page, totalPages, err := query.Select(
		ctx,
		tx,
		M2mRuleTriggerJobTableColumnsWithTypeCasts,
		M2mRuleTriggerJobTableWithSchema,
		where,
		orderBy,
		limit,
		offset,
		values...,
	)
	if err != nil {
		return nil, 0, 0, 0, 0, fmt.Errorf("failed to call SelectM2mRuleTriggerJobs; %v", err)
	}

	objects := make([]*M2mRuleTriggerJob, 0)

	for _, item := range *items {
		object := &M2mRuleTriggerJob{}

		err = object.FromItem(item)
		if err != nil {
			return nil, 0, 0, 0, 0, err
		}

		if !types.IsZeroUUID(object.JobID) {
			ctx, ok := query.HandleQueryPathGraphCycles(ctx, fmt.Sprintf("%s{%v}", JobTable, object.JobID), true)
			shouldLoad := query.ShouldLoad(ctx, JobTable)
			if ok || shouldLoad {
				thisBefore := time.Now()

				if config.Debug() {
					log.Printf("loading SelectM2mRuleTriggerJobs->SelectJob for object.JobIDObject{%s: %v}", JobTablePrimaryKeyColumn, object.JobID)
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
					log.Printf("loaded SelectM2mRuleTriggerJobs->SelectJob for object.JobIDObject in %s", time.Since(thisBefore))
				}
			}
		}

		if !types.IsZeroUUID(object.RuleID) {
			ctx, ok := query.HandleQueryPathGraphCycles(ctx, fmt.Sprintf("%s{%v}", RuleTable, object.RuleID), true)
			shouldLoad := query.ShouldLoad(ctx, RuleTable)
			if ok || shouldLoad {
				thisBefore := time.Now()

				if config.Debug() {
					log.Printf("loading SelectM2mRuleTriggerJobs->SelectRule for object.RuleIDObject{%s: %v}", RuleTablePrimaryKeyColumn, object.RuleID)
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
					log.Printf("loaded SelectM2mRuleTriggerJobs->SelectRule for object.RuleIDObject in %s", time.Since(thisBefore))
				}
			}
		}

		err = func() error {
			shouldLoad := query.ShouldLoad(ctx, fmt.Sprintf("referenced_by_%s", ExecutionTable))
			ctx, ok := query.HandleQueryPathGraphCycles(ctx, fmt.Sprintf("__ReferencedBy__%s{%v}", ExecutionTable, object.GetPrimaryKeyValue()), true)
			if ok || shouldLoad {
				thisBefore := time.Now()

				if config.Debug() {
					log.Printf("loading SelectM2mRuleTriggerJobs->SelectExecutions for object.ReferencedByExecutionM2mRuleTriggerJobIDObjects")
				}

				object.ReferencedByExecutionM2mRuleTriggerJobIDObjects, _, _, _, _, err = SelectExecutions(
					ctx,
					tx,
					fmt.Sprintf("%v = $1", ExecutionTableM2mRuleTriggerJobIDColumn),
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
					log.Printf("loaded SelectM2mRuleTriggerJobs->SelectExecutions for object.ReferencedByExecutionM2mRuleTriggerJobIDObjects in %s", time.Since(thisBefore))
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

func SelectM2mRuleTriggerJob(ctx context.Context, tx pgx.Tx, where string, values ...any) (*M2mRuleTriggerJob, int64, int64, int64, int64, error) {
	ctx, cleanup := query.WithQueryID(ctx)
	defer cleanup()

	ctx = query.WithMaxDepth(ctx, nil)

	objects, _, _, _, _, err := SelectM2mRuleTriggerJobs(
		ctx,
		tx,
		where,
		nil,
		helpers.Ptr(2),
		helpers.Ptr(0),
		values...,
	)
	if err != nil {
		return nil, 0, 0, 0, 0, fmt.Errorf("failed to call SelectM2mRuleTriggerJob; %v", err)
	}

	if len(objects) > 1 {
		return nil, 0, 0, 0, 0, fmt.Errorf("attempt to call SelectM2mRuleTriggerJob returned more than 1 row")
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

func JobExecutorClaimM2mRuleTriggerJob(ctx context.Context, tx pgx.Tx, until time.Time, timeout time.Duration, wheres ...string) (*M2mRuleTriggerJob, error) {
	m := &M2mRuleTriggerJob{}

	err := m.AdvisoryLockWithRetries(ctx, tx, math.MinInt32, timeout, time.Second*1)
	if err != nil {
		return nil, fmt.Errorf("failed to claim: %s", err.Error())
	}

	extraWhere := ""
	if len(wheres) > 0 {
		extraWhere = fmt.Sprintf(" AND\n    %s", strings.Join(wheres, " AND\n    "))
	}

	ms, _, _, _, _, err := SelectM2mRuleTriggerJobs(
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

func handleGetM2mRuleTriggerJobs(arguments *server.SelectManyArguments, db *pgxpool.Pool) ([]*M2mRuleTriggerJob, int64, int64, int64, int64, error) {
	tx, err := db.Begin(arguments.Ctx)
	if err != nil {
		return nil, 0, 0, 0, 0, err
	}

	defer func() {
		_ = tx.Rollback(arguments.Ctx)
	}()

	objects, count, totalCount, page, totalPages, err := SelectM2mRuleTriggerJobs(arguments.Ctx, tx, arguments.Where, arguments.OrderBy, arguments.Limit, arguments.Offset, arguments.Values...)
	if err != nil {
		return nil, 0, 0, 0, 0, err
	}

	err = tx.Commit(arguments.Ctx)
	if err != nil {
		return nil, 0, 0, 0, 0, err
	}

	return objects, count, totalCount, page, totalPages, nil
}

func handleGetM2mRuleTriggerJob(arguments *server.SelectOneArguments, db *pgxpool.Pool, primaryKey uuid.UUID) ([]*M2mRuleTriggerJob, int64, int64, int64, int64, error) {
	tx, err := db.Begin(arguments.Ctx)
	if err != nil {
		return nil, 0, 0, 0, 0, err
	}

	defer func() {
		_ = tx.Rollback(arguments.Ctx)
	}()

	object, count, totalCount, page, totalPages, err := SelectM2mRuleTriggerJob(arguments.Ctx, tx, arguments.Where, arguments.Values...)
	if err != nil {
		return nil, 0, 0, 0, 0, err
	}

	err = tx.Commit(arguments.Ctx)
	if err != nil {
		return nil, 0, 0, 0, 0, err
	}

	return []*M2mRuleTriggerJob{object}, count, totalCount, page, totalPages, nil
}

func handlePostM2mRuleTriggerJob(arguments *server.LoadArguments, db *pgxpool.Pool, waitForChange server.WaitForChange, objects []*M2mRuleTriggerJob, forceSetValuesForFieldsByObjectIndex [][]string) ([]*M2mRuleTriggerJob, int64, int64, int64, int64, error) {
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
		_, err := waitForChange(arguments.Ctx, []stream.Action{stream.INSERT}, M2mRuleTriggerJobTable, xid)
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

func handlePutM2mRuleTriggerJob(arguments *server.LoadArguments, db *pgxpool.Pool, waitForChange server.WaitForChange, object *M2mRuleTriggerJob) ([]*M2mRuleTriggerJob, int64, int64, int64, int64, error) {
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
		_, err := waitForChange(arguments.Ctx, []stream.Action{stream.UPDATE, stream.SOFT_DELETE, stream.SOFT_RESTORE, stream.SOFT_UPDATE}, M2mRuleTriggerJobTable, xid)
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

	return []*M2mRuleTriggerJob{object}, count, totalCount, page, totalPages, nil
}

func handlePatchM2mRuleTriggerJob(arguments *server.LoadArguments, db *pgxpool.Pool, waitForChange server.WaitForChange, object *M2mRuleTriggerJob, forceSetValuesForFields []string) ([]*M2mRuleTriggerJob, int64, int64, int64, int64, error) {
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
		_, err := waitForChange(arguments.Ctx, []stream.Action{stream.UPDATE, stream.SOFT_DELETE, stream.SOFT_RESTORE, stream.SOFT_UPDATE}, M2mRuleTriggerJobTable, xid)
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

	return []*M2mRuleTriggerJob{object}, count, totalCount, page, totalPages, nil
}

func handleDeleteM2mRuleTriggerJob(arguments *server.LoadArguments, db *pgxpool.Pool, waitForChange server.WaitForChange, object *M2mRuleTriggerJob) error {
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
		_, err := waitForChange(arguments.Ctx, []stream.Action{stream.DELETE, stream.SOFT_DELETE}, M2mRuleTriggerJobTable, xid)
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

func MutateRouterForM2mRuleTriggerJob(r chi.Router, db *pgxpool.Pool, redisPool *redis.Pool, objectMiddlewares []server.ObjectMiddleware, waitForChange server.WaitForChange) {

	func() {
		postHandlerForJobExecutorClaim, err := getHTTPHandler(
			http.MethodPost,
			"/job-executor-claim-m2m-rule-trigger-job",
			http.StatusOK,
			func(
				ctx context.Context,
				pathParams server.EmptyPathParams,
				queryParams server.EmptyQueryParams,
				req M2mRuleTriggerJobJobExecutorClaimRequest,
				rawReq any,
			) (server.Response[M2mRuleTriggerJob], error) {
				tx, err := db.Begin(ctx)
				if err != nil {
					return server.Response[M2mRuleTriggerJob]{}, err
				}

				defer func() {
					_ = tx.Rollback(ctx)
				}()

				object, err := JobExecutorClaimM2mRuleTriggerJob(ctx, tx, req.Until, time.Millisecond*time.Duration(req.TimeoutSeconds*1000))
				if err != nil {
					return server.Response[M2mRuleTriggerJob]{}, err
				}

				count := int64(0)

				totalCount := int64(0)

				limit := int64(0)

				offset := int64(0)

				if object == nil {
					return server.Response[M2mRuleTriggerJob]{
						Status:     http.StatusOK,
						Success:    true,
						Error:      nil,
						Objects:    []*M2mRuleTriggerJob{},
						Count:      count,
						TotalCount: totalCount,
						Limit:      limit,
						Offset:     offset,
					}, nil
				}

				err = tx.Commit(ctx)
				if err != nil {
					return server.Response[M2mRuleTriggerJob]{}, err
				}

				return server.Response[M2mRuleTriggerJob]{
					Status:     http.StatusOK,
					Success:    true,
					Error:      nil,
					Objects:    []*M2mRuleTriggerJob{object},
					Count:      count,
					TotalCount: totalCount,
					Limit:      limit,
					Offset:     offset,
				}, nil
			},
			M2mRuleTriggerJob{},
			M2mRuleTriggerJobIntrospectedTable,
		)
		if err != nil {
			panic(err)
		}
		r.Post(postHandlerForJobExecutorClaim.FullPath, postHandlerForJobExecutorClaim.ServeHTTP)

		postHandlerForJobExecutorClaimOne, err := getHTTPHandler(
			http.MethodPost,
			"/m2m-rule-trigger-jobs/{primaryKey}/job-executor-claim",
			http.StatusOK,
			func(
				ctx context.Context,
				pathParams M2mRuleTriggerJobOnePathParams,
				queryParams M2mRuleTriggerJobLoadQueryParams,
				req M2mRuleTriggerJobJobExecutorClaimRequest,
				rawReq any,
			) (server.Response[M2mRuleTriggerJob], error) {
				before := time.Now()

				redisConn := redisPool.Get()
				defer func() {
					_ = redisConn.Close()
				}()

				arguments, err := server.GetSelectOneArguments(ctx, queryParams.Depth, M2mRuleTriggerJobIntrospectedTable, pathParams.PrimaryKey, nil, nil)
				if err != nil {
					if config.Debug() {
						log.Printf("request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return server.Response[M2mRuleTriggerJob]{}, err
				}

				/* note: deliberately no attempt at a cache hit */

				var object *M2mRuleTriggerJob
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

					object, count, totalCount, _, _, err = SelectM2mRuleTriggerJob(arguments.Ctx, tx, arguments.Where, arguments.Values...)
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

					return server.Response[M2mRuleTriggerJob]{}, err
				}

				limit := int64(0)

				offset := int64(0)

				response := server.Response[M2mRuleTriggerJob]{
					Status:     http.StatusOK,
					Success:    true,
					Error:      nil,
					Objects:    []*M2mRuleTriggerJob{object},
					Count:      count,
					TotalCount: totalCount,
					Limit:      limit,
					Offset:     offset,
				}

				return response, nil
			},
			M2mRuleTriggerJob{},
			M2mRuleTriggerJobIntrospectedTable,
		)
		if err != nil {
			panic(err)
		}
		r.Post(postHandlerForJobExecutorClaimOne.FullPath, postHandlerForJobExecutorClaimOne.ServeHTTP)
	}()

	func() {
		getManyHandler, err := getHTTPHandler(
			http.MethodGet,
			"/m2m-rule-trigger-jobs",
			http.StatusOK,
			func(
				ctx context.Context,
				pathParams server.EmptyPathParams,
				queryParams map[string]any,
				req server.EmptyRequest,
				rawReq any,
			) (server.Response[M2mRuleTriggerJob], error) {
				before := time.Now()

				redisConn := redisPool.Get()
				defer func() {
					_ = redisConn.Close()
				}()

				arguments, err := server.GetSelectManyArguments(ctx, queryParams, M2mRuleTriggerJobIntrospectedTable, nil, nil)
				if err != nil {
					if config.Debug() {
						log.Printf("request cache not yet reached; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return server.Response[M2mRuleTriggerJob]{}, err
				}

				cachedResponseAsJSON, cacheHit, err := server.GetCachedResponseAsJSON(arguments.RequestHash, redisConn)
				if err != nil {
					if config.Debug() {
						log.Printf("request cache failed; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return server.Response[M2mRuleTriggerJob]{}, err
				}

				if cacheHit {
					var cachedResponse server.Response[M2mRuleTriggerJob]

					/* TODO: it'd be nice to be able to avoid this (i.e. just pass straight through) */
					err = json.Unmarshal(cachedResponseAsJSON, &cachedResponse)
					if err != nil {
						if config.Debug() {
							log.Printf("request cache hit but failed unmarshal; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
						}

						return server.Response[M2mRuleTriggerJob]{}, err
					}

					if config.Debug() {
						log.Printf("request cache hit; request succeeded in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return cachedResponse, nil
				}

				objects, count, totalCount, _, _, err := handleGetM2mRuleTriggerJobs(arguments, db)
				if err != nil {
					if config.Debug() {
						log.Printf("request cache missed; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return server.Response[M2mRuleTriggerJob]{}, err
				}

				limit := int64(0)
				if arguments.Limit != nil {
					limit = int64(*arguments.Limit)
				}

				offset := int64(0)
				if arguments.Offset != nil {
					offset = int64(*arguments.Offset)
				}

				response := server.Response[M2mRuleTriggerJob]{
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

					return server.Response[M2mRuleTriggerJob]{}, err
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
			M2mRuleTriggerJob{},
			M2mRuleTriggerJobIntrospectedTable,
		)
		if err != nil {
			panic(err)
		}
		r.Get(getManyHandler.FullPath, getManyHandler.ServeHTTP)
	}()

	func() {
		getOneHandler, err := getHTTPHandler(
			http.MethodGet,
			"/m2m-rule-trigger-jobs/{primaryKey}",
			http.StatusOK,
			func(
				ctx context.Context,
				pathParams M2mRuleTriggerJobOnePathParams,
				queryParams M2mRuleTriggerJobLoadQueryParams,
				req server.EmptyRequest,
				rawReq any,
			) (server.Response[M2mRuleTriggerJob], error) {
				before := time.Now()

				redisConn := redisPool.Get()
				defer func() {
					_ = redisConn.Close()
				}()

				arguments, err := server.GetSelectOneArguments(ctx, queryParams.Depth, M2mRuleTriggerJobIntrospectedTable, pathParams.PrimaryKey, nil, nil)
				if err != nil {
					if config.Debug() {
						log.Printf("request cache not yet reached; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return server.Response[M2mRuleTriggerJob]{}, err
				}

				cachedResponseAsJSON, cacheHit, err := server.GetCachedResponseAsJSON(arguments.RequestHash, redisConn)
				if err != nil {
					if config.Debug() {
						log.Printf("request cache failed; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return server.Response[M2mRuleTriggerJob]{}, err
				}

				if cacheHit {
					var cachedResponse server.Response[M2mRuleTriggerJob]

					/* TODO: it'd be nice to be able to avoid this (i.e. just pass straight through) */
					err = json.Unmarshal(cachedResponseAsJSON, &cachedResponse)
					if err != nil {
						if config.Debug() {
							log.Printf("request cache hit but failed unmarshal; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
						}

						return server.Response[M2mRuleTriggerJob]{}, err
					}

					if config.Debug() {
						log.Printf("request cache hit; request succeeded in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return cachedResponse, nil
				}

				objects, count, totalCount, _, _, err := handleGetM2mRuleTriggerJob(arguments, db, pathParams.PrimaryKey)
				if err != nil {
					if config.Debug() {
						log.Printf("request cache missed; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return server.Response[M2mRuleTriggerJob]{}, err
				}

				limit := int64(0)

				offset := int64(0)

				response := server.Response[M2mRuleTriggerJob]{
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

					return server.Response[M2mRuleTriggerJob]{}, err
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
			M2mRuleTriggerJob{},
			M2mRuleTriggerJobIntrospectedTable,
		)
		if err != nil {
			panic(err)
		}
		r.Get(getOneHandler.FullPath, getOneHandler.ServeHTTP)
	}()

	func() {
		postHandler, err := getHTTPHandler(
			http.MethodPost,
			"/m2m-rule-trigger-jobs",
			http.StatusCreated,
			func(
				ctx context.Context,
				pathParams server.EmptyPathParams,
				queryParams M2mRuleTriggerJobLoadQueryParams,
				req []*M2mRuleTriggerJob,
				rawReq any,
			) (server.Response[M2mRuleTriggerJob], error) {
				allRawItems, ok := rawReq.([]any)
				if !ok {
					return server.Response[M2mRuleTriggerJob]{}, fmt.Errorf("failed to cast %#+v to []map[string]any", rawReq)
				}

				allItems := make([]map[string]any, 0)
				for _, rawItem := range allRawItems {
					item, ok := rawItem.(map[string]any)
					if !ok {
						return server.Response[M2mRuleTriggerJob]{}, fmt.Errorf("failed to cast %#+v to map[string]any", rawItem)
					}

					allItems = append(allItems, item)
				}

				forceSetValuesForFieldsByObjectIndex := make([][]string, 0)
				for _, item := range allItems {
					forceSetValuesForFields := make([]string, 0)
					for _, possibleField := range maps.Keys(item) {
						if !slices.Contains(M2mRuleTriggerJobTableColumns, possibleField) {
							continue
						}

						forceSetValuesForFields = append(forceSetValuesForFields, possibleField)
					}
					forceSetValuesForFieldsByObjectIndex = append(forceSetValuesForFieldsByObjectIndex, forceSetValuesForFields)
				}

				arguments, err := server.GetLoadArguments(ctx, queryParams.Depth)
				if err != nil {
					return server.Response[M2mRuleTriggerJob]{}, err
				}

				objects, count, totalCount, _, _, err := handlePostM2mRuleTriggerJob(arguments, db, waitForChange, req, forceSetValuesForFieldsByObjectIndex)
				if err != nil {
					return server.Response[M2mRuleTriggerJob]{}, err
				}

				limit := int64(0)

				offset := int64(0)

				return server.Response[M2mRuleTriggerJob]{
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
			M2mRuleTriggerJob{},
			M2mRuleTriggerJobIntrospectedTable,
		)
		if err != nil {
			panic(err)
		}
		r.Post(postHandler.FullPath, postHandler.ServeHTTP)
	}()

	func() {
		putHandler, err := getHTTPHandler(
			http.MethodPatch,
			"/m2m-rule-trigger-jobs/{primaryKey}",
			http.StatusOK,
			func(
				ctx context.Context,
				pathParams M2mRuleTriggerJobOnePathParams,
				queryParams M2mRuleTriggerJobLoadQueryParams,
				req M2mRuleTriggerJob,
				rawReq any,
			) (server.Response[M2mRuleTriggerJob], error) {
				item, ok := rawReq.(map[string]any)
				if !ok {
					return server.Response[M2mRuleTriggerJob]{}, fmt.Errorf("failed to cast %#+v to map[string]any", item)
				}

				arguments, err := server.GetLoadArguments(ctx, queryParams.Depth)
				if err != nil {
					return server.Response[M2mRuleTriggerJob]{}, err
				}

				object := &req
				object.ID = pathParams.PrimaryKey

				objects, count, totalCount, _, _, err := handlePutM2mRuleTriggerJob(arguments, db, waitForChange, object)
				if err != nil {
					return server.Response[M2mRuleTriggerJob]{}, err
				}

				limit := int64(0)

				offset := int64(0)

				return server.Response[M2mRuleTriggerJob]{
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
			M2mRuleTriggerJob{},
			M2mRuleTriggerJobIntrospectedTable,
		)
		if err != nil {
			panic(err)
		}
		r.Put(putHandler.FullPath, putHandler.ServeHTTP)
	}()

	func() {
		patchHandler, err := getHTTPHandler(
			http.MethodPatch,
			"/m2m-rule-trigger-jobs/{primaryKey}",
			http.StatusOK,
			func(
				ctx context.Context,
				pathParams M2mRuleTriggerJobOnePathParams,
				queryParams M2mRuleTriggerJobLoadQueryParams,
				req M2mRuleTriggerJob,
				rawReq any,
			) (server.Response[M2mRuleTriggerJob], error) {
				item, ok := rawReq.(map[string]any)
				if !ok {
					return server.Response[M2mRuleTriggerJob]{}, fmt.Errorf("failed to cast %#+v to map[string]any", item)
				}

				forceSetValuesForFields := make([]string, 0)
				for _, possibleField := range maps.Keys(item) {
					if !slices.Contains(M2mRuleTriggerJobTableColumns, possibleField) {
						continue
					}

					forceSetValuesForFields = append(forceSetValuesForFields, possibleField)
				}

				arguments, err := server.GetLoadArguments(ctx, queryParams.Depth)
				if err != nil {
					return server.Response[M2mRuleTriggerJob]{}, err
				}

				object := &req
				object.ID = pathParams.PrimaryKey

				objects, count, totalCount, _, _, err := handlePatchM2mRuleTriggerJob(arguments, db, waitForChange, object, forceSetValuesForFields)
				if err != nil {
					return server.Response[M2mRuleTriggerJob]{}, err
				}

				limit := int64(0)

				offset := int64(0)

				return server.Response[M2mRuleTriggerJob]{
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
			M2mRuleTriggerJob{},
			M2mRuleTriggerJobIntrospectedTable,
		)
		if err != nil {
			panic(err)
		}
		r.Patch(patchHandler.FullPath, patchHandler.ServeHTTP)
	}()

	func() {
		deleteHandler, err := getHTTPHandler(
			http.MethodDelete,
			"/m2m-rule-trigger-jobs/{primaryKey}",
			http.StatusNoContent,
			func(
				ctx context.Context,
				pathParams M2mRuleTriggerJobOnePathParams,
				queryParams M2mRuleTriggerJobLoadQueryParams,
				req server.EmptyRequest,
				rawReq any,
			) (server.EmptyResponse, error) {
				arguments, err := server.GetLoadArguments(ctx, queryParams.Depth)
				if err != nil {
					return server.EmptyResponse{}, err
				}

				object := &M2mRuleTriggerJob{}
				object.ID = pathParams.PrimaryKey

				err = handleDeleteM2mRuleTriggerJob(arguments, db, waitForChange, object)
				if err != nil {
					return server.EmptyResponse{}, err
				}

				return server.EmptyResponse{}, nil
			},
			M2mRuleTriggerJob{},
			M2mRuleTriggerJobIntrospectedTable,
		)
		if err != nil {
			panic(err)
		}
		r.Delete(deleteHandler.FullPath, deleteHandler.ServeHTTP)
	}()
}

func NewM2mRuleTriggerJobFromItem(item map[string]any) (any, error) {
	object := &M2mRuleTriggerJob{}

	err := object.FromItem(item)
	if err != nil {
		return nil, err
	}

	return object, nil
}

func init() {
	register(
		M2mRuleTriggerJobTable,
		M2mRuleTriggerJob{},
		NewM2mRuleTriggerJobFromItem,
		"/m2m-rule-trigger-jobs",
		MutateRouterForM2mRuleTriggerJob,
	)
}
