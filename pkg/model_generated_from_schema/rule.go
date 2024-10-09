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

type Rule struct {
	ID                                         uuid.UUID            `json:"id"`
	CreatedAt                                  time.Time            `json:"created_at"`
	UpdatedAt                                  time.Time            `json:"updated_at"`
	DeletedAt                                  *time.Time           `json:"deleted_at"`
	BranchName                                 *string              `json:"branch_name"`
	RepositoryID                               uuid.UUID            `json:"repository_id"`
	RepositoryIDObject                         *Repository          `json:"repository_id_object"`
	ReferencedByM2mRuleTriggerJobRuleIDObjects []*M2mRuleTriggerJob `json:"referenced_by_m2m_rule_trigger_job_rule_id_objects"`
}

var RuleTable = "rule"

var RuleTableNamespaceID int32 = 1337 + 8

var (
	RuleTableIDColumn           = "id"
	RuleTableCreatedAtColumn    = "created_at"
	RuleTableUpdatedAtColumn    = "updated_at"
	RuleTableDeletedAtColumn    = "deleted_at"
	RuleTableBranchNameColumn   = "branch_name"
	RuleTableRepositoryIDColumn = "repository_id"
)

var (
	RuleTableIDColumnWithTypeCast           = `"id" AS id`
	RuleTableCreatedAtColumnWithTypeCast    = `"created_at" AS created_at`
	RuleTableUpdatedAtColumnWithTypeCast    = `"updated_at" AS updated_at`
	RuleTableDeletedAtColumnWithTypeCast    = `"deleted_at" AS deleted_at`
	RuleTableBranchNameColumnWithTypeCast   = `"branch_name" AS branch_name`
	RuleTableRepositoryIDColumnWithTypeCast = `"repository_id" AS repository_id`
)

var RuleTableColumns = []string{
	RuleTableIDColumn,
	RuleTableCreatedAtColumn,
	RuleTableUpdatedAtColumn,
	RuleTableDeletedAtColumn,
	RuleTableBranchNameColumn,
	RuleTableRepositoryIDColumn,
}

var RuleTableColumnsWithTypeCasts = []string{
	RuleTableIDColumnWithTypeCast,
	RuleTableCreatedAtColumnWithTypeCast,
	RuleTableUpdatedAtColumnWithTypeCast,
	RuleTableDeletedAtColumnWithTypeCast,
	RuleTableBranchNameColumnWithTypeCast,
	RuleTableRepositoryIDColumnWithTypeCast,
}

var RuleIntrospectedTable *introspect.Table

var RuleTableColumnLookup map[string]*introspect.Column

var (
	RuleTablePrimaryKeyColumn = RuleTableIDColumn
)

func init() {
	RuleIntrospectedTable = tableByName[RuleTable]

	/* only needed during templating */
	if RuleIntrospectedTable == nil {
		RuleIntrospectedTable = &introspect.Table{}
	}

	RuleTableColumnLookup = RuleIntrospectedTable.ColumnByName
}

type RuleOnePathParams struct {
	PrimaryKey uuid.UUID `json:"primaryKey"`
}

type RuleLoadQueryParams struct {
	Depth *int `json:"depth"`
}

type RuleClaimRequest struct {
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

func (m *Rule) GetPrimaryKeyColumn() string {
	return RuleTablePrimaryKeyColumn
}

func (m *Rule) GetPrimaryKeyValue() any {
	return m.ID
}

func (m *Rule) FromItem(item map[string]any) error {
	if item == nil {
		return fmt.Errorf(
			"item unexpectedly nil during RuleFromItem",
		)
	}

	if len(item) == 0 {
		return fmt.Errorf(
			"item unexpectedly empty during RuleFromItem",
		)
	}

	wrapError := func(k string, v any, err error) error {
		return fmt.Errorf("%v: %#+v; error; %v", k, v, err)
	}

	for k, v := range item {
		_, ok := RuleTableColumnLookup[k]
		if !ok {
			return fmt.Errorf(
				"item contained unexpected key %#+v during RuleFromItem; item: %#+v",
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

		case "branch_name":
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
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uubranch_name.UUID", temp1))
				}
			}

			m.BranchName = &temp2

		case "repository_id":
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
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uurepository_id.UUID", temp1))
				}
			}

			m.RepositoryID = temp2

		}
	}

	return nil
}

func (m *Rule) Reload(ctx context.Context, tx pgx.Tx, includeDeleteds ...bool) error {
	extraWhere := ""
	if len(includeDeleteds) > 0 && includeDeleteds[0] {
		if slices.Contains(RuleTableColumns, "deleted_at") {
			extraWhere = "\n    AND (deleted_at IS null OR deleted_at IS NOT null)"
		}
	}

	ctx, cleanup := query.WithQueryID(ctx)
	defer cleanup()

	ctx = query.WithMaxDepth(ctx, nil)

	o, _, _, _, _, err := SelectRule(
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
	m.BranchName = o.BranchName
	m.RepositoryID = o.RepositoryID
	m.RepositoryIDObject = o.RepositoryIDObject
	m.ReferencedByM2mRuleTriggerJobRuleIDObjects = o.ReferencedByM2mRuleTriggerJobRuleIDObjects

	return nil
}

func (m *Rule) Insert(ctx context.Context, tx pgx.Tx, setPrimaryKey bool, setZeroValues bool, forceSetValuesForFields ...string) error {
	columns := make([]string, 0)
	values := make([]any, 0)

	if setPrimaryKey && (setZeroValues || !types.IsZeroUUID(m.ID) || slices.Contains(forceSetValuesForFields, RuleTableIDColumn) || isRequired(RuleTableColumnLookup, RuleTableIDColumn)) {
		columns = append(columns, RuleTableIDColumn)

		v, err := types.FormatUUID(m.ID)
		if err != nil {
			return fmt.Errorf("failed to handle m.ID; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.CreatedAt) || slices.Contains(forceSetValuesForFields, RuleTableCreatedAtColumn) || isRequired(RuleTableColumnLookup, RuleTableCreatedAtColumn) {
		columns = append(columns, RuleTableCreatedAtColumn)

		v, err := types.FormatTime(m.CreatedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.CreatedAt; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.UpdatedAt) || slices.Contains(forceSetValuesForFields, RuleTableUpdatedAtColumn) || isRequired(RuleTableColumnLookup, RuleTableUpdatedAtColumn) {
		columns = append(columns, RuleTableUpdatedAtColumn)

		v, err := types.FormatTime(m.UpdatedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.UpdatedAt; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.DeletedAt) || slices.Contains(forceSetValuesForFields, RuleTableDeletedAtColumn) || isRequired(RuleTableColumnLookup, RuleTableDeletedAtColumn) {
		columns = append(columns, RuleTableDeletedAtColumn)

		v, err := types.FormatTime(m.DeletedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.DeletedAt; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroString(m.BranchName) || slices.Contains(forceSetValuesForFields, RuleTableBranchNameColumn) || isRequired(RuleTableColumnLookup, RuleTableBranchNameColumn) {
		columns = append(columns, RuleTableBranchNameColumn)

		v, err := types.FormatString(m.BranchName)
		if err != nil {
			return fmt.Errorf("failed to handle m.BranchName; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroUUID(m.RepositoryID) || slices.Contains(forceSetValuesForFields, RuleTableRepositoryIDColumn) || isRequired(RuleTableColumnLookup, RuleTableRepositoryIDColumn) {
		columns = append(columns, RuleTableRepositoryIDColumn)

		v, err := types.FormatUUID(m.RepositoryID)
		if err != nil {
			return fmt.Errorf("failed to handle m.RepositoryID; %v", err)
		}

		values = append(values, v)
	}

	ctx, cleanup := query.WithQueryID(ctx)
	defer cleanup()

	ctx = query.WithMaxDepth(ctx, nil)

	item, err := query.Insert(
		ctx,
		tx,
		RuleTable,
		columns,
		nil,
		false,
		false,
		RuleTableColumns,
		values...,
	)
	if err != nil {
		return fmt.Errorf("failed to insert %#+v; %v", m, err)
	}
	v := (*item)[RuleTableIDColumn]

	if v == nil {
		return fmt.Errorf("failed to find %v in %#+v", RuleTableIDColumn, item)
	}

	wrapError := func(err error) error {
		return fmt.Errorf(
			"failed to treat %v: %#+v as uuid.UUID: %v",
			RuleTableIDColumn,
			(*item)[RuleTableIDColumn],
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

func (m *Rule) Update(ctx context.Context, tx pgx.Tx, setZeroValues bool, forceSetValuesForFields ...string) error {
	columns := make([]string, 0)
	values := make([]any, 0)

	if setZeroValues || !types.IsZeroTime(m.CreatedAt) || slices.Contains(forceSetValuesForFields, RuleTableCreatedAtColumn) {
		columns = append(columns, RuleTableCreatedAtColumn)

		v, err := types.FormatTime(m.CreatedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.CreatedAt; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.UpdatedAt) || slices.Contains(forceSetValuesForFields, RuleTableUpdatedAtColumn) {
		columns = append(columns, RuleTableUpdatedAtColumn)

		v, err := types.FormatTime(m.UpdatedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.UpdatedAt; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.DeletedAt) || slices.Contains(forceSetValuesForFields, RuleTableDeletedAtColumn) {
		columns = append(columns, RuleTableDeletedAtColumn)

		v, err := types.FormatTime(m.DeletedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.DeletedAt; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroString(m.BranchName) || slices.Contains(forceSetValuesForFields, RuleTableBranchNameColumn) {
		columns = append(columns, RuleTableBranchNameColumn)

		v, err := types.FormatString(m.BranchName)
		if err != nil {
			return fmt.Errorf("failed to handle m.BranchName; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroUUID(m.RepositoryID) || slices.Contains(forceSetValuesForFields, RuleTableRepositoryIDColumn) {
		columns = append(columns, RuleTableRepositoryIDColumn)

		v, err := types.FormatUUID(m.RepositoryID)
		if err != nil {
			return fmt.Errorf("failed to handle m.RepositoryID; %v", err)
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
		RuleTable,
		columns,
		fmt.Sprintf("%v = $$??", RuleTableIDColumn),
		RuleTableColumns,
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

func (m *Rule) Delete(ctx context.Context, tx pgx.Tx, hardDeletes ...bool) error {
	hardDelete := false
	if len(hardDeletes) > 0 {
		hardDelete = hardDeletes[0]
	}

	if !hardDelete && slices.Contains(RuleTableColumns, "deleted_at") {
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
		RuleTable,
		fmt.Sprintf("%v = $$??", RuleTableIDColumn),
		values...,
	)
	if err != nil {
		return fmt.Errorf("failed to delete %#+v; %v", m, err)
	}

	_ = m.Reload(ctx, tx, true)

	return nil
}

func (m *Rule) LockTable(ctx context.Context, tx pgx.Tx, timeouts ...time.Duration) error {
	return query.LockTable(ctx, tx, RuleTable, timeouts...)
}

func (m *Rule) LockTableWithRetries(ctx context.Context, tx pgx.Tx, overallTimeout time.Duration, individualAttempttimeout time.Duration) error {
	return query.LockTableWithRetries(ctx, tx, RuleTable, overallTimeout, individualAttempttimeout)
}

func (m *Rule) AdvisoryLock(ctx context.Context, tx pgx.Tx, key int32, timeouts ...time.Duration) error {
	return query.AdvisoryLock(ctx, tx, RuleTableNamespaceID, key, timeouts...)
}

func (m *Rule) AdvisoryLockWithRetries(ctx context.Context, tx pgx.Tx, key int32, overallTimeout time.Duration, individualAttempttimeout time.Duration) error {
	return query.AdvisoryLockWithRetries(ctx, tx, RuleTableNamespaceID, key, overallTimeout, individualAttempttimeout)
}

func (m *Rule) Claim(ctx context.Context, tx pgx.Tx, until time.Time, by uuid.UUID, timeout time.Duration) error {
	if !(slices.Contains(RuleTableColumns, "claimed_until") && slices.Contains(RuleTableColumns, "claimed_by")) {
		return fmt.Errorf("can only invoke Claim for tables with 'claimed_until' and 'claimed_by' columns")
	}

	err := m.AdvisoryLockWithRetries(ctx, tx, math.MinInt32, timeout, time.Second*1)
	if err != nil {
		return fmt.Errorf("failed to claim (advisory lock): %s", err.Error())
	}

	x, _, _, _, _, err := SelectRule(
		ctx,
		tx,
		fmt.Sprintf(
			"%s = $$?? AND (claimed_by = $$?? OR (claimed_until IS null OR claimed_until < now()))",
			RuleTablePrimaryKeyColumn,
		),
		m.GetPrimaryKeyValue(),
		by,
	)
	if err != nil {
		return fmt.Errorf("failed to claim (select): %s", err.Error())
	}

	_ = x

	/* m.ClaimedUntil = &until */
	/* m.ClaimedBy = &by */

	err = m.Update(ctx, tx, false)
	if err != nil {
		return fmt.Errorf("failed to claim (update): %s", err.Error())
	}

	return nil
}

func SelectRules(ctx context.Context, tx pgx.Tx, where string, orderBy *string, limit *int, offset *int, values ...any) ([]*Rule, int64, int64, int64, int64, error) {
	before := time.Now()

	if config.Debug() {
		log.Printf("entered SelectRules")

		defer func() {
			log.Printf("exited SelectRules in %s", time.Since(before))
		}()
	}
	if slices.Contains(RuleTableColumns, "deleted_at") {
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

	shouldLoad := query.ShouldLoad(ctx, RuleTable) || query.ShouldLoad(ctx, fmt.Sprintf("referenced_by_%s", RuleTable))

	var ok bool
	ctx, ok = query.HandleQueryPathGraphCycles(ctx, fmt.Sprintf("%s{%v}", RuleTable, nil), !isLoadQuery)
	if !ok && !shouldLoad {
		if config.Debug() {
			log.Printf("skipping SelectRule early (query.ShouldLoad(): %v, query.HandleQueryPathGraphCycles(): %v)", shouldLoad, ok)
		}
		return []*Rule{}, 0, 0, 0, 0, nil
	}

	items, count, totalCount, page, totalPages, err := query.Select(
		ctx,
		tx,
		RuleTableColumnsWithTypeCasts,
		RuleTable,
		where,
		orderBy,
		limit,
		offset,
		values...,
	)
	if err != nil {
		return nil, 0, 0, 0, 0, fmt.Errorf("failed to call SelectRules; %v", err)
	}

	objects := make([]*Rule, 0)

	for _, item := range *items {
		object := &Rule{}

		err = object.FromItem(item)
		if err != nil {
			return nil, 0, 0, 0, 0, err
		}

		if !types.IsZeroUUID(object.RepositoryID) {
			ctx, ok := query.HandleQueryPathGraphCycles(ctx, fmt.Sprintf("%s{%v}", RepositoryTable, object.RepositoryID), true)
			shouldLoad := query.ShouldLoad(ctx, RepositoryTable)
			if ok || shouldLoad {
				thisBefore := time.Now()

				if config.Debug() {
					log.Printf("loading SelectRules->SelectRepository for object.RepositoryIDObject{%s: %v}", RepositoryTablePrimaryKeyColumn, object.RepositoryID)
				}

				object.RepositoryIDObject, _, _, _, _, err = SelectRepository(
					ctx,
					tx,
					fmt.Sprintf("%v = $1", RepositoryTablePrimaryKeyColumn),
					object.RepositoryID,
				)
				if err != nil {
					if !errors.Is(err, sql.ErrNoRows) {
						return nil, 0, 0, 0, 0, err
					}
				}

				if config.Debug() {
					log.Printf("loaded SelectRules->SelectRepository for object.RepositoryIDObject in %s", time.Since(thisBefore))
				}
			}
		}

		err = func() error {
			shouldLoad := query.ShouldLoad(ctx, fmt.Sprintf("referenced_by_%s", M2mRuleTriggerJobTable))
			ctx, ok := query.HandleQueryPathGraphCycles(ctx, fmt.Sprintf("__ReferencedBy__%s{%v}", M2mRuleTriggerJobTable, object.GetPrimaryKeyValue()), true)
			if ok || shouldLoad {
				thisBefore := time.Now()

				if config.Debug() {
					log.Printf("loading SelectRules->SelectM2mRuleTriggerJobs for object.ReferencedByM2mRuleTriggerJobRuleIDObjects")
				}

				object.ReferencedByM2mRuleTriggerJobRuleIDObjects, _, _, _, _, err = SelectM2mRuleTriggerJobs(
					ctx,
					tx,
					fmt.Sprintf("%v = $1", M2mRuleTriggerJobTableRuleIDColumn),
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
					log.Printf("loaded SelectRules->SelectM2mRuleTriggerJobs for object.ReferencedByM2mRuleTriggerJobRuleIDObjects in %s", time.Since(thisBefore))
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

func SelectRule(ctx context.Context, tx pgx.Tx, where string, values ...any) (*Rule, int64, int64, int64, int64, error) {
	ctx, cleanup := query.WithQueryID(ctx)
	defer cleanup()

	ctx = query.WithMaxDepth(ctx, nil)

	objects, _, _, _, _, err := SelectRules(
		ctx,
		tx,
		where,
		nil,
		helpers.Ptr(2),
		helpers.Ptr(0),
		values...,
	)
	if err != nil {
		return nil, 0, 0, 0, 0, fmt.Errorf("failed to call SelectRule; %v", err)
	}

	if len(objects) > 1 {
		return nil, 0, 0, 0, 0, fmt.Errorf("attempt to call SelectRule returned more than 1 row")
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

func ClaimRule(ctx context.Context, tx pgx.Tx, until time.Time, by uuid.UUID, timeout time.Duration, wheres ...string) (*Rule, error) {
	if !(slices.Contains(RuleTableColumns, "claimed_until") && slices.Contains(RuleTableColumns, "claimed_by")) {
		return nil, fmt.Errorf("can only invoke Claim for tables with 'claimed_until' and 'claimed_by' columns")
	}

	m := &Rule{}

	err := m.AdvisoryLockWithRetries(ctx, tx, math.MinInt32, timeout, time.Second*1)
	if err != nil {
		return nil, fmt.Errorf("failed to claim: %s", err.Error())
	}

	extraWhere := ""
	if len(wheres) > 0 {
		extraWhere = fmt.Sprintf("AND %s", extraWhere)
	}

	ms, _, _, _, _, err := SelectRules(
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

	/* m.ClaimedUntil = &until */
	/* m.ClaimedBy = &by */

	err = m.Update(ctx, tx, false)
	if err != nil {
		return nil, fmt.Errorf("failed to claim: %s", err.Error())
	}

	return m, nil
}

func handleGetRules(arguments *server.SelectManyArguments, db *pgxpool.Pool) ([]*Rule, int64, int64, int64, int64, error) {
	tx, err := db.Begin(arguments.Ctx)
	if err != nil {
		return nil, 0, 0, 0, 0, err
	}

	defer func() {
		_ = tx.Rollback(arguments.Ctx)
	}()

	objects, count, totalCount, page, totalPages, err := SelectRules(arguments.Ctx, tx, arguments.Where, arguments.OrderBy, arguments.Limit, arguments.Offset, arguments.Values...)
	if err != nil {
		return nil, 0, 0, 0, 0, err
	}

	err = tx.Commit(arguments.Ctx)
	if err != nil {
		return nil, 0, 0, 0, 0, err
	}

	return objects, count, totalCount, page, totalPages, nil
}

func handleGetRule(arguments *server.SelectOneArguments, db *pgxpool.Pool, primaryKey uuid.UUID) ([]*Rule, int64, int64, int64, int64, error) {
	tx, err := db.Begin(arguments.Ctx)
	if err != nil {
		return nil, 0, 0, 0, 0, err
	}

	defer func() {
		_ = tx.Rollback(arguments.Ctx)
	}()

	object, count, totalCount, page, totalPages, err := SelectRule(arguments.Ctx, tx, arguments.Where, arguments.Values...)
	if err != nil {
		return nil, 0, 0, 0, 0, err
	}

	err = tx.Commit(arguments.Ctx)
	if err != nil {
		return nil, 0, 0, 0, 0, err
	}

	return []*Rule{object}, count, totalCount, page, totalPages, nil
}

func handlePostRules(arguments *server.LoadArguments, db *pgxpool.Pool, waitForChange server.WaitForChange, objects []*Rule, forceSetValuesForFieldsByObjectIndex [][]string) ([]*Rule, int64, int64, int64, int64, error) {
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
		_, err := waitForChange(arguments.Ctx, []stream.Action{stream.INSERT}, RuleTable, xid)
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

func handlePutRule(arguments *server.LoadArguments, db *pgxpool.Pool, waitForChange server.WaitForChange, object *Rule) ([]*Rule, int64, int64, int64, int64, error) {
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
		_, err := waitForChange(arguments.Ctx, []stream.Action{stream.UPDATE, stream.SOFT_DELETE, stream.SOFT_RESTORE, stream.SOFT_UPDATE}, RuleTable, xid)
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

	return []*Rule{object}, count, totalCount, page, totalPages, nil
}

func handlePatchRule(arguments *server.LoadArguments, db *pgxpool.Pool, waitForChange server.WaitForChange, object *Rule, forceSetValuesForFields []string) ([]*Rule, int64, int64, int64, int64, error) {
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
		_, err := waitForChange(arguments.Ctx, []stream.Action{stream.UPDATE, stream.SOFT_DELETE, stream.SOFT_RESTORE, stream.SOFT_UPDATE}, RuleTable, xid)
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

	return []*Rule{object}, count, totalCount, page, totalPages, nil
}

func handleDeleteRule(arguments *server.LoadArguments, db *pgxpool.Pool, waitForChange server.WaitForChange, object *Rule) error {
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
		_, err := waitForChange(arguments.Ctx, []stream.Action{stream.DELETE, stream.SOFT_DELETE}, RuleTable, xid)
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

func MutateRouterForRule(r chi.Router, db *pgxpool.Pool, redisPool *redis.Pool, objectMiddlewares []server.ObjectMiddleware, waitForChange server.WaitForChange) {
	if slices.Contains(RuleTableColumns, "claimed_until") && slices.Contains(RuleTableColumns, "claimed_by") {
		func() {
			postHandlerForClaim, err := getHTTPHandler(
				http.MethodPost,
				"/claim-rule",
				http.StatusOK,
				func(
					ctx context.Context,
					pathParams server.EmptyPathParams,
					queryParams server.EmptyQueryParams,
					req RuleClaimRequest,
					rawReq any,
				) (server.Response[Rule], error) {
					tx, err := db.Begin(ctx)
					if err != nil {
						return server.Response[Rule]{}, err
					}

					defer func() {
						_ = tx.Rollback(ctx)
					}()

					object, err := ClaimRule(ctx, tx, req.Until, req.By, time.Millisecond*time.Duration(req.TimeoutSeconds*1000))
					if err != nil {
						return server.Response[Rule]{}, err
					}

					count := int64(0)

					totalCount := int64(0)

					limit := int64(0)

					offset := int64(0)

					if object == nil {
						return server.Response[Rule]{
							Status:     http.StatusOK,
							Success:    true,
							Error:      nil,
							Objects:    []*Rule{},
							Count:      count,
							TotalCount: totalCount,
							Limit:      limit,
							Offset:     offset,
						}, nil
					}

					err = tx.Commit(ctx)
					if err != nil {
						return server.Response[Rule]{}, err
					}

					return server.Response[Rule]{
						Status:     http.StatusOK,
						Success:    true,
						Error:      nil,
						Objects:    []*Rule{object},
						Count:      count,
						TotalCount: totalCount,
						Limit:      limit,
						Offset:     offset,
					}, nil
				},
				Rule{},
				RuleIntrospectedTable,
			)
			if err != nil {
				panic(err)
			}
			r.Post(postHandlerForClaim.FullPath, postHandlerForClaim.ServeHTTP)

			postHandlerForClaimOne, err := getHTTPHandler(
				http.MethodPost,
				"/rules/{primaryKey}/claim",
				http.StatusOK,
				func(
					ctx context.Context,
					pathParams RuleOnePathParams,
					queryParams RuleLoadQueryParams,
					req RuleClaimRequest,
					rawReq any,
				) (server.Response[Rule], error) {
					before := time.Now()

					redisConn := redisPool.Get()
					defer func() {
						_ = redisConn.Close()
					}()

					arguments, err := server.GetSelectOneArguments(ctx, queryParams.Depth, RuleIntrospectedTable, pathParams.PrimaryKey, nil, nil)
					if err != nil {
						if config.Debug() {
							log.Printf("request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
						}

						return server.Response[Rule]{}, err
					}

					/* note: deliberately no attempt at a cache hit */

					var object *Rule
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

						object, count, totalCount, _, _, err = SelectRule(arguments.Ctx, tx, arguments.Where, arguments.Values...)
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

						return server.Response[Rule]{}, err
					}

					limit := int64(0)

					offset := int64(0)

					response := server.Response[Rule]{
						Status:     http.StatusOK,
						Success:    true,
						Error:      nil,
						Objects:    []*Rule{object},
						Count:      count,
						TotalCount: totalCount,
						Limit:      limit,
						Offset:     offset,
					}

					return response, nil
				},
				Rule{},
				RuleIntrospectedTable,
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
			"/rules",
			http.StatusOK,
			func(
				ctx context.Context,
				pathParams server.EmptyPathParams,
				queryParams map[string]any,
				req server.EmptyRequest,
				rawReq any,
			) (server.Response[Rule], error) {
				before := time.Now()

				redisConn := redisPool.Get()
				defer func() {
					_ = redisConn.Close()
				}()

				arguments, err := server.GetSelectManyArguments(ctx, queryParams, RuleIntrospectedTable, nil, nil)
				if err != nil {
					if config.Debug() {
						log.Printf("request cache not yet reached; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return server.Response[Rule]{}, err
				}

				cachedResponseAsJSON, cacheHit, err := server.GetCachedResponseAsJSON(arguments.RequestHash, redisConn)
				if err != nil {
					if config.Debug() {
						log.Printf("request cache failed; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return server.Response[Rule]{}, err
				}

				if cacheHit {
					var cachedResponse server.Response[Rule]

					/* TODO: it'd be nice to be able to avoid this (i.e. just pass straight through) */
					err = json.Unmarshal(cachedResponseAsJSON, &cachedResponse)
					if err != nil {
						if config.Debug() {
							log.Printf("request cache hit but failed unmarshal; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
						}

						return server.Response[Rule]{}, err
					}

					if config.Debug() {
						log.Printf("request cache hit; request succeeded in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return cachedResponse, nil
				}

				objects, count, totalCount, _, _, err := handleGetRules(arguments, db)
				if err != nil {
					if config.Debug() {
						log.Printf("request cache missed; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return server.Response[Rule]{}, err
				}

				limit := int64(0)
				if arguments.Limit != nil {
					limit = int64(*arguments.Limit)
				}

				offset := int64(0)
				if arguments.Offset != nil {
					offset = int64(*arguments.Offset)
				}

				response := server.Response[Rule]{
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

					return server.Response[Rule]{}, err
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
			Rule{},
			RuleIntrospectedTable,
		)
		if err != nil {
			panic(err)
		}
		r.Get(getManyHandler.FullPath, getManyHandler.ServeHTTP)
	}()

	func() {
		getOneHandler, err := getHTTPHandler(
			http.MethodGet,
			"/rules/{primaryKey}",
			http.StatusOK,
			func(
				ctx context.Context,
				pathParams RuleOnePathParams,
				queryParams RuleLoadQueryParams,
				req server.EmptyRequest,
				rawReq any,
			) (server.Response[Rule], error) {
				before := time.Now()

				redisConn := redisPool.Get()
				defer func() {
					_ = redisConn.Close()
				}()

				arguments, err := server.GetSelectOneArguments(ctx, queryParams.Depth, RuleIntrospectedTable, pathParams.PrimaryKey, nil, nil)
				if err != nil {
					if config.Debug() {
						log.Printf("request cache not yet reached; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return server.Response[Rule]{}, err
				}

				cachedResponseAsJSON, cacheHit, err := server.GetCachedResponseAsJSON(arguments.RequestHash, redisConn)
				if err != nil {
					if config.Debug() {
						log.Printf("request cache failed; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return server.Response[Rule]{}, err
				}

				if cacheHit {
					var cachedResponse server.Response[Rule]

					/* TODO: it'd be nice to be able to avoid this (i.e. just pass straight through) */
					err = json.Unmarshal(cachedResponseAsJSON, &cachedResponse)
					if err != nil {
						if config.Debug() {
							log.Printf("request cache hit but failed unmarshal; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
						}

						return server.Response[Rule]{}, err
					}

					if config.Debug() {
						log.Printf("request cache hit; request succeeded in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return cachedResponse, nil
				}

				objects, count, totalCount, _, _, err := handleGetRule(arguments, db, pathParams.PrimaryKey)
				if err != nil {
					if config.Debug() {
						log.Printf("request cache missed; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return server.Response[Rule]{}, err
				}

				limit := int64(0)

				offset := int64(0)

				response := server.Response[Rule]{
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

					return server.Response[Rule]{}, err
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
			Rule{},
			RuleIntrospectedTable,
		)
		if err != nil {
			panic(err)
		}
		r.Get(getOneHandler.FullPath, getOneHandler.ServeHTTP)
	}()

	func() {
		postHandler, err := getHTTPHandler(
			http.MethodPost,
			"/rules",
			http.StatusCreated,
			func(
				ctx context.Context,
				pathParams server.EmptyPathParams,
				queryParams RuleLoadQueryParams,
				req []*Rule,
				rawReq any,
			) (server.Response[Rule], error) {
				allRawItems, ok := rawReq.([]any)
				if !ok {
					return server.Response[Rule]{}, fmt.Errorf("failed to cast %#+v to []map[string]any", rawReq)
				}

				allItems := make([]map[string]any, 0)
				for _, rawItem := range allRawItems {
					item, ok := rawItem.(map[string]any)
					if !ok {
						return server.Response[Rule]{}, fmt.Errorf("failed to cast %#+v to map[string]any", rawItem)
					}

					allItems = append(allItems, item)
				}

				forceSetValuesForFieldsByObjectIndex := make([][]string, 0)
				for _, item := range allItems {
					forceSetValuesForFields := make([]string, 0)
					for _, possibleField := range maps.Keys(item) {
						if !slices.Contains(RuleTableColumns, possibleField) {
							continue
						}

						forceSetValuesForFields = append(forceSetValuesForFields, possibleField)
					}
					forceSetValuesForFieldsByObjectIndex = append(forceSetValuesForFieldsByObjectIndex, forceSetValuesForFields)
				}

				arguments, err := server.GetLoadArguments(ctx, queryParams.Depth)
				if err != nil {
					return server.Response[Rule]{}, err
				}

				objects, count, totalCount, _, _, err := handlePostRules(arguments, db, waitForChange, req, forceSetValuesForFieldsByObjectIndex)
				if err != nil {
					return server.Response[Rule]{}, err
				}

				limit := int64(0)

				offset := int64(0)

				return server.Response[Rule]{
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
			Rule{},
			RuleIntrospectedTable,
		)
		if err != nil {
			panic(err)
		}
		r.Post(postHandler.FullPath, postHandler.ServeHTTP)
	}()

	func() {
		putHandler, err := getHTTPHandler(
			http.MethodPatch,
			"/rules/{primaryKey}",
			http.StatusOK,
			func(
				ctx context.Context,
				pathParams RuleOnePathParams,
				queryParams RuleLoadQueryParams,
				req Rule,
				rawReq any,
			) (server.Response[Rule], error) {
				item, ok := rawReq.(map[string]any)
				if !ok {
					return server.Response[Rule]{}, fmt.Errorf("failed to cast %#+v to map[string]any", item)
				}

				arguments, err := server.GetLoadArguments(ctx, queryParams.Depth)
				if err != nil {
					return server.Response[Rule]{}, err
				}

				object := &req
				object.ID = pathParams.PrimaryKey

				objects, count, totalCount, _, _, err := handlePutRule(arguments, db, waitForChange, object)
				if err != nil {
					return server.Response[Rule]{}, err
				}

				limit := int64(0)

				offset := int64(0)

				return server.Response[Rule]{
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
			Rule{},
			RuleIntrospectedTable,
		)
		if err != nil {
			panic(err)
		}
		r.Put(putHandler.FullPath, putHandler.ServeHTTP)
	}()

	func() {
		patchHandler, err := getHTTPHandler(
			http.MethodPatch,
			"/rules/{primaryKey}",
			http.StatusOK,
			func(
				ctx context.Context,
				pathParams RuleOnePathParams,
				queryParams RuleLoadQueryParams,
				req Rule,
				rawReq any,
			) (server.Response[Rule], error) {
				item, ok := rawReq.(map[string]any)
				if !ok {
					return server.Response[Rule]{}, fmt.Errorf("failed to cast %#+v to map[string]any", item)
				}

				forceSetValuesForFields := make([]string, 0)
				for _, possibleField := range maps.Keys(item) {
					if !slices.Contains(RuleTableColumns, possibleField) {
						continue
					}

					forceSetValuesForFields = append(forceSetValuesForFields, possibleField)
				}

				arguments, err := server.GetLoadArguments(ctx, queryParams.Depth)
				if err != nil {
					return server.Response[Rule]{}, err
				}

				object := &req
				object.ID = pathParams.PrimaryKey

				objects, count, totalCount, _, _, err := handlePatchRule(arguments, db, waitForChange, object, forceSetValuesForFields)
				if err != nil {
					return server.Response[Rule]{}, err
				}

				limit := int64(0)

				offset := int64(0)

				return server.Response[Rule]{
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
			Rule{},
			RuleIntrospectedTable,
		)
		if err != nil {
			panic(err)
		}
		r.Patch(patchHandler.FullPath, patchHandler.ServeHTTP)
	}()

	func() {
		deleteHandler, err := getHTTPHandler(
			http.MethodDelete,
			"/rules/{primaryKey}",
			http.StatusNoContent,
			func(
				ctx context.Context,
				pathParams RuleOnePathParams,
				queryParams RuleLoadQueryParams,
				req server.EmptyRequest,
				rawReq any,
			) (server.EmptyResponse, error) {
				arguments, err := server.GetLoadArguments(ctx, queryParams.Depth)
				if err != nil {
					return server.EmptyResponse{}, err
				}

				object := &Rule{}
				object.ID = pathParams.PrimaryKey

				err = handleDeleteRule(arguments, db, waitForChange, object)
				if err != nil {
					return server.EmptyResponse{}, err
				}

				return server.EmptyResponse{}, nil
			},
			Rule{},
			RuleIntrospectedTable,
		)
		if err != nil {
			panic(err)
		}
		r.Delete(deleteHandler.FullPath, deleteHandler.ServeHTTP)
	}()
}

func NewRuleFromItem(item map[string]any) (any, error) {
	object := &Rule{}

	err := object.FromItem(item)
	if err != nil {
		return nil, err
	}

	return object, nil
}

func init() {
	register(
		RuleTable,
		Rule{},
		NewRuleFromItem,
		"/rules",
		MutateRouterForRule,
	)
}
