package model_generated_from_schema

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

type Output struct {
	ID                             uuid.UUID  `json:"id"`
	CreatedAt                      time.Time  `json:"created_at"`
	UpdatedAt                      time.Time  `json:"updated_at"`
	DeletedAt                      *time.Time `json:"deleted_at"`
	Status                         string     `json:"status"`
	StartedAt                      *time.Time `json:"started_at"`
	EndedAt                        *time.Time `json:"ended_at"`
	ExitStatus                     int64      `json:"exit_status"`
	Error                          *string    `json:"error"`
	TaskID                         uuid.UUID  `json:"task_id"`
	TaskIDObject                   *Task      `json:"task_id_object"`
	Logid                          uuid.UUID  `json:"logid"`
	LogidObject                    *Log       `json:"logid_object"`
	ReferencedByLogOutputIDObjects []*Log     `json:"referenced_by_log_output_id_objects"`
}

var OutputTable = "output"

var OutputTableWithSchema = fmt.Sprintf("%s.%s", schema, OutputTable)

var OutputTableNamespaceID int32 = 1337 + 6

var (
	OutputTableIDColumn         = "id"
	OutputTableCreatedAtColumn  = "created_at"
	OutputTableUpdatedAtColumn  = "updated_at"
	OutputTableDeletedAtColumn  = "deleted_at"
	OutputTableStatusColumn     = "status"
	OutputTableStartedAtColumn  = "started_at"
	OutputTableEndedAtColumn    = "ended_at"
	OutputTableExitStatusColumn = "exit_status"
	OutputTableErrorColumn      = "error"
	OutputTableTaskIDColumn     = "task_id"
	OutputTableLogidColumn      = "logid"
)

var (
	OutputTableIDColumnWithTypeCast         = `"id" AS id`
	OutputTableCreatedAtColumnWithTypeCast  = `"created_at" AS created_at`
	OutputTableUpdatedAtColumnWithTypeCast  = `"updated_at" AS updated_at`
	OutputTableDeletedAtColumnWithTypeCast  = `"deleted_at" AS deleted_at`
	OutputTableStatusColumnWithTypeCast     = `"status" AS status`
	OutputTableStartedAtColumnWithTypeCast  = `"started_at" AS started_at`
	OutputTableEndedAtColumnWithTypeCast    = `"ended_at" AS ended_at`
	OutputTableExitStatusColumnWithTypeCast = `"exit_status" AS exit_status`
	OutputTableErrorColumnWithTypeCast      = `"error" AS error`
	OutputTableTaskIDColumnWithTypeCast     = `"task_id" AS task_id`
	OutputTableLogidColumnWithTypeCast      = `"logid" AS logid`
)

var OutputTableColumns = []string{
	OutputTableIDColumn,
	OutputTableCreatedAtColumn,
	OutputTableUpdatedAtColumn,
	OutputTableDeletedAtColumn,
	OutputTableStatusColumn,
	OutputTableStartedAtColumn,
	OutputTableEndedAtColumn,
	OutputTableExitStatusColumn,
	OutputTableErrorColumn,
	OutputTableTaskIDColumn,
	OutputTableLogidColumn,
}

var OutputTableColumnsWithTypeCasts = []string{
	OutputTableIDColumnWithTypeCast,
	OutputTableCreatedAtColumnWithTypeCast,
	OutputTableUpdatedAtColumnWithTypeCast,
	OutputTableDeletedAtColumnWithTypeCast,
	OutputTableStatusColumnWithTypeCast,
	OutputTableStartedAtColumnWithTypeCast,
	OutputTableEndedAtColumnWithTypeCast,
	OutputTableExitStatusColumnWithTypeCast,
	OutputTableErrorColumnWithTypeCast,
	OutputTableTaskIDColumnWithTypeCast,
	OutputTableLogidColumnWithTypeCast,
}

var OutputIntrospectedTable *introspect.Table

var OutputTableColumnLookup map[string]*introspect.Column

var (
	OutputTablePrimaryKeyColumn = OutputTableIDColumn
)

func init() {
	OutputIntrospectedTable = tableByName[OutputTable]

	/* only needed during templating */
	if OutputIntrospectedTable == nil {
		OutputIntrospectedTable = &introspect.Table{}
	}

	OutputTableColumnLookup = OutputIntrospectedTable.ColumnByName
}

type OutputOnePathParams struct {
	PrimaryKey uuid.UUID `json:"primaryKey"`
}

type OutputLoadQueryParams struct {
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

func (m *Output) GetPrimaryKeyColumn() string {
	return OutputTablePrimaryKeyColumn
}

func (m *Output) GetPrimaryKeyValue() any {
	return m.ID
}

func (m *Output) FromItem(item map[string]any) error {
	if item == nil {
		return fmt.Errorf(
			"item unexpectedly nil during OutputFromItem",
		)
	}

	if len(item) == 0 {
		return fmt.Errorf(
			"item unexpectedly empty during OutputFromItem",
		)
	}

	wrapError := func(k string, v any, err error) error {
		return fmt.Errorf("%v: %#+v; error; %v", k, v, err)
	}

	for k, v := range item {
		_, ok := OutputTableColumnLookup[k]
		if !ok {
			return fmt.Errorf(
				"item contained unexpected key %#+v during OutputFromItem; item: %#+v",
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

		case "exit_status":
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
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uuexit_status.UUID", temp1))
				}
			}

			m.ExitStatus = temp2

		case "error":
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
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uuerror.UUID", temp1))
				}
			}

			m.Error = &temp2

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

		case "logid":
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
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uulogid.UUID", temp1))
				}
			}

			m.Logid = temp2

		}
	}

	return nil
}

func (m *Output) Reload(ctx context.Context, tx pgx.Tx, includeDeleteds ...bool) error {
	extraWhere := ""
	if len(includeDeleteds) > 0 && includeDeleteds[0] {
		if slices.Contains(OutputTableColumns, "deleted_at") {
			extraWhere = "\n    AND (deleted_at IS null OR deleted_at IS NOT null)"
		}
	}

	ctx, cleanup := query.WithQueryID(ctx)
	defer cleanup()

	ctx = query.WithMaxDepth(ctx, nil)

	o, _, _, _, _, err := SelectOutput(
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
	m.ExitStatus = o.ExitStatus
	m.Error = o.Error
	m.TaskID = o.TaskID
	m.TaskIDObject = o.TaskIDObject
	m.Logid = o.Logid
	m.LogidObject = o.LogidObject
	m.ReferencedByLogOutputIDObjects = o.ReferencedByLogOutputIDObjects

	return nil
}

func (m *Output) Insert(ctx context.Context, tx pgx.Tx, setPrimaryKey bool, setZeroValues bool, forceSetValuesForFields ...string) error {
	columns := make([]string, 0)
	values := make([]any, 0)

	if setPrimaryKey && (setZeroValues || !types.IsZeroUUID(m.ID) || slices.Contains(forceSetValuesForFields, OutputTableIDColumn) || isRequired(OutputTableColumnLookup, OutputTableIDColumn)) {
		columns = append(columns, OutputTableIDColumn)

		v, err := types.FormatUUID(m.ID)
		if err != nil {
			return fmt.Errorf("failed to handle m.ID; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.CreatedAt) || slices.Contains(forceSetValuesForFields, OutputTableCreatedAtColumn) || isRequired(OutputTableColumnLookup, OutputTableCreatedAtColumn) {
		columns = append(columns, OutputTableCreatedAtColumn)

		v, err := types.FormatTime(m.CreatedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.CreatedAt; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.UpdatedAt) || slices.Contains(forceSetValuesForFields, OutputTableUpdatedAtColumn) || isRequired(OutputTableColumnLookup, OutputTableUpdatedAtColumn) {
		columns = append(columns, OutputTableUpdatedAtColumn)

		v, err := types.FormatTime(m.UpdatedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.UpdatedAt; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.DeletedAt) || slices.Contains(forceSetValuesForFields, OutputTableDeletedAtColumn) || isRequired(OutputTableColumnLookup, OutputTableDeletedAtColumn) {
		columns = append(columns, OutputTableDeletedAtColumn)

		v, err := types.FormatTime(m.DeletedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.DeletedAt; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroString(m.Status) || slices.Contains(forceSetValuesForFields, OutputTableStatusColumn) || isRequired(OutputTableColumnLookup, OutputTableStatusColumn) {
		columns = append(columns, OutputTableStatusColumn)

		v, err := types.FormatString(m.Status)
		if err != nil {
			return fmt.Errorf("failed to handle m.Status; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.StartedAt) || slices.Contains(forceSetValuesForFields, OutputTableStartedAtColumn) || isRequired(OutputTableColumnLookup, OutputTableStartedAtColumn) {
		columns = append(columns, OutputTableStartedAtColumn)

		v, err := types.FormatTime(m.StartedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.StartedAt; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.EndedAt) || slices.Contains(forceSetValuesForFields, OutputTableEndedAtColumn) || isRequired(OutputTableColumnLookup, OutputTableEndedAtColumn) {
		columns = append(columns, OutputTableEndedAtColumn)

		v, err := types.FormatTime(m.EndedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.EndedAt; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroInt(m.ExitStatus) || slices.Contains(forceSetValuesForFields, OutputTableExitStatusColumn) || isRequired(OutputTableColumnLookup, OutputTableExitStatusColumn) {
		columns = append(columns, OutputTableExitStatusColumn)

		v, err := types.FormatInt(m.ExitStatus)
		if err != nil {
			return fmt.Errorf("failed to handle m.ExitStatus; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroString(m.Error) || slices.Contains(forceSetValuesForFields, OutputTableErrorColumn) || isRequired(OutputTableColumnLookup, OutputTableErrorColumn) {
		columns = append(columns, OutputTableErrorColumn)

		v, err := types.FormatString(m.Error)
		if err != nil {
			return fmt.Errorf("failed to handle m.Error; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroUUID(m.TaskID) || slices.Contains(forceSetValuesForFields, OutputTableTaskIDColumn) || isRequired(OutputTableColumnLookup, OutputTableTaskIDColumn) {
		columns = append(columns, OutputTableTaskIDColumn)

		v, err := types.FormatUUID(m.TaskID)
		if err != nil {
			return fmt.Errorf("failed to handle m.TaskID; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroUUID(m.Logid) || slices.Contains(forceSetValuesForFields, OutputTableLogidColumn) || isRequired(OutputTableColumnLookup, OutputTableLogidColumn) {
		columns = append(columns, OutputTableLogidColumn)

		v, err := types.FormatUUID(m.Logid)
		if err != nil {
			return fmt.Errorf("failed to handle m.Logid; %v", err)
		}

		values = append(values, v)
	}

	ctx, cleanup := query.WithQueryID(ctx)
	defer cleanup()

	ctx = query.WithMaxDepth(ctx, nil)

	item, err := query.Insert(
		ctx,
		tx,
		OutputTableWithSchema,
		columns,
		nil,
		false,
		false,
		OutputTableColumns,
		values...,
	)
	if err != nil {
		return fmt.Errorf("failed to insert %#+v; %v", m, err)
	}
	v := (*item)[OutputTableIDColumn]

	if v == nil {
		return fmt.Errorf("failed to find %v in %#+v", OutputTableIDColumn, item)
	}

	wrapError := func(err error) error {
		return fmt.Errorf(
			"failed to treat %v: %#+v as uuid.UUID: %v",
			OutputTableIDColumn,
			(*item)[OutputTableIDColumn],
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

func (m *Output) Update(ctx context.Context, tx pgx.Tx, setZeroValues bool, forceSetValuesForFields ...string) error {
	columns := make([]string, 0)
	values := make([]any, 0)

	if setZeroValues || !types.IsZeroTime(m.CreatedAt) || slices.Contains(forceSetValuesForFields, OutputTableCreatedAtColumn) {
		columns = append(columns, OutputTableCreatedAtColumn)

		v, err := types.FormatTime(m.CreatedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.CreatedAt; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.UpdatedAt) || slices.Contains(forceSetValuesForFields, OutputTableUpdatedAtColumn) {
		columns = append(columns, OutputTableUpdatedAtColumn)

		v, err := types.FormatTime(m.UpdatedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.UpdatedAt; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.DeletedAt) || slices.Contains(forceSetValuesForFields, OutputTableDeletedAtColumn) {
		columns = append(columns, OutputTableDeletedAtColumn)

		v, err := types.FormatTime(m.DeletedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.DeletedAt; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroString(m.Status) || slices.Contains(forceSetValuesForFields, OutputTableStatusColumn) {
		columns = append(columns, OutputTableStatusColumn)

		v, err := types.FormatString(m.Status)
		if err != nil {
			return fmt.Errorf("failed to handle m.Status; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.StartedAt) || slices.Contains(forceSetValuesForFields, OutputTableStartedAtColumn) {
		columns = append(columns, OutputTableStartedAtColumn)

		v, err := types.FormatTime(m.StartedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.StartedAt; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.EndedAt) || slices.Contains(forceSetValuesForFields, OutputTableEndedAtColumn) {
		columns = append(columns, OutputTableEndedAtColumn)

		v, err := types.FormatTime(m.EndedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.EndedAt; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroInt(m.ExitStatus) || slices.Contains(forceSetValuesForFields, OutputTableExitStatusColumn) {
		columns = append(columns, OutputTableExitStatusColumn)

		v, err := types.FormatInt(m.ExitStatus)
		if err != nil {
			return fmt.Errorf("failed to handle m.ExitStatus; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroString(m.Error) || slices.Contains(forceSetValuesForFields, OutputTableErrorColumn) {
		columns = append(columns, OutputTableErrorColumn)

		v, err := types.FormatString(m.Error)
		if err != nil {
			return fmt.Errorf("failed to handle m.Error; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroUUID(m.TaskID) || slices.Contains(forceSetValuesForFields, OutputTableTaskIDColumn) {
		columns = append(columns, OutputTableTaskIDColumn)

		v, err := types.FormatUUID(m.TaskID)
		if err != nil {
			return fmt.Errorf("failed to handle m.TaskID; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroUUID(m.Logid) || slices.Contains(forceSetValuesForFields, OutputTableLogidColumn) {
		columns = append(columns, OutputTableLogidColumn)

		v, err := types.FormatUUID(m.Logid)
		if err != nil {
			return fmt.Errorf("failed to handle m.Logid; %v", err)
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
		OutputTableWithSchema,
		columns,
		fmt.Sprintf("%v = $$??", OutputTableIDColumn),
		OutputTableColumns,
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

func (m *Output) Delete(ctx context.Context, tx pgx.Tx, hardDeletes ...bool) error {
	hardDelete := false
	if len(hardDeletes) > 0 {
		hardDelete = hardDeletes[0]
	}

	if !hardDelete && slices.Contains(OutputTableColumns, "deleted_at") {
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
		OutputTableWithSchema,
		fmt.Sprintf("%v = $$??", OutputTableIDColumn),
		values...,
	)
	if err != nil {
		return fmt.Errorf("failed to delete %#+v; %v", m, err)
	}

	_ = m.Reload(ctx, tx, true)

	return nil
}

func (m *Output) LockTable(ctx context.Context, tx pgx.Tx, timeouts ...time.Duration) error {
	return query.LockTable(ctx, tx, OutputTableWithSchema, timeouts...)
}

func (m *Output) LockTableWithRetries(ctx context.Context, tx pgx.Tx, overallTimeout time.Duration, individualAttempttimeout time.Duration) error {
	return query.LockTableWithRetries(ctx, tx, OutputTableWithSchema, overallTimeout, individualAttempttimeout)
}

func (m *Output) AdvisoryLock(ctx context.Context, tx pgx.Tx, key int32, timeouts ...time.Duration) error {
	return query.AdvisoryLock(ctx, tx, OutputTableNamespaceID, key, timeouts...)
}

func (m *Output) AdvisoryLockWithRetries(ctx context.Context, tx pgx.Tx, key int32, overallTimeout time.Duration, individualAttempttimeout time.Duration) error {
	return query.AdvisoryLockWithRetries(ctx, tx, OutputTableNamespaceID, key, overallTimeout, individualAttempttimeout)
}

func SelectOutputs(ctx context.Context, tx pgx.Tx, where string, orderBy *string, limit *int, offset *int, values ...any) ([]*Output, int64, int64, int64, int64, error) {
	before := time.Now()

	if config.Debug() {
		log.Printf("entered SelectOutputs")

		defer func() {
			log.Printf("exited SelectOutputs in %s", time.Since(before))
		}()
	}
	if slices.Contains(OutputTableColumns, "deleted_at") {
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

	shouldLoad := query.ShouldLoad(ctx, OutputTable) || query.ShouldLoad(ctx, fmt.Sprintf("referenced_by_%s", OutputTable))

	var ok bool
	ctx, ok = query.HandleQueryPathGraphCycles(ctx, fmt.Sprintf("%s{%v}", OutputTable, nil), !isLoadQuery)
	if !ok && !shouldLoad {
		if config.Debug() {
			log.Printf("skipping SelectOutput early (query.ShouldLoad(): %v, query.HandleQueryPathGraphCycles(): %v)", shouldLoad, ok)
		}
		return []*Output{}, 0, 0, 0, 0, nil
	}

	items, count, totalCount, page, totalPages, err := query.Select(
		ctx,
		tx,
		OutputTableColumnsWithTypeCasts,
		OutputTableWithSchema,
		where,
		orderBy,
		limit,
		offset,
		values...,
	)
	if err != nil {
		return nil, 0, 0, 0, 0, fmt.Errorf("failed to call SelectOutputs; %v", err)
	}

	objects := make([]*Output, 0)

	for _, item := range *items {
		object := &Output{}

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
					log.Printf("loading SelectOutputs->SelectTask for object.TaskIDObject{%s: %v}", TaskTablePrimaryKeyColumn, object.TaskID)
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
					log.Printf("loaded SelectOutputs->SelectTask for object.TaskIDObject in %s", time.Since(thisBefore))
				}
			}
		}

		if !types.IsZeroUUID(object.Logid) {
			ctx, ok := query.HandleQueryPathGraphCycles(ctx, fmt.Sprintf("%s{%v}", LogTable, object.Logid), true)
			shouldLoad := query.ShouldLoad(ctx, LogTable)
			if ok || shouldLoad {
				thisBefore := time.Now()

				if config.Debug() {
					log.Printf("loading SelectOutputs->SelectLog for object.LogidObject{%s: %v}", LogTablePrimaryKeyColumn, object.Logid)
				}

				object.LogidObject, _, _, _, _, err = SelectLog(
					ctx,
					tx,
					fmt.Sprintf("%v = $1", LogTablePrimaryKeyColumn),
					object.Logid,
				)
				if err != nil {
					if !errors.Is(err, sql.ErrNoRows) {
						return nil, 0, 0, 0, 0, err
					}
				}

				if config.Debug() {
					log.Printf("loaded SelectOutputs->SelectLog for object.LogidObject in %s", time.Since(thisBefore))
				}
			}
		}

		err = func() error {
			shouldLoad := query.ShouldLoad(ctx, fmt.Sprintf("referenced_by_%s", LogTable))
			ctx, ok := query.HandleQueryPathGraphCycles(ctx, fmt.Sprintf("__ReferencedBy__%s{%v}", LogTable, object.GetPrimaryKeyValue()), true)
			if ok || shouldLoad {
				thisBefore := time.Now()

				if config.Debug() {
					log.Printf("loading SelectOutputs->SelectLogs for object.ReferencedByLogOutputIDObjects")
				}

				object.ReferencedByLogOutputIDObjects, _, _, _, _, err = SelectLogs(
					ctx,
					tx,
					fmt.Sprintf("%v = $1", LogTableOutputIDColumn),
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
					log.Printf("loaded SelectOutputs->SelectLogs for object.ReferencedByLogOutputIDObjects in %s", time.Since(thisBefore))
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

func SelectOutput(ctx context.Context, tx pgx.Tx, where string, values ...any) (*Output, int64, int64, int64, int64, error) {
	ctx, cleanup := query.WithQueryID(ctx)
	defer cleanup()

	ctx = query.WithMaxDepth(ctx, nil)

	objects, _, _, _, _, err := SelectOutputs(
		ctx,
		tx,
		where,
		nil,
		helpers.Ptr(2),
		helpers.Ptr(0),
		values...,
	)
	if err != nil {
		return nil, 0, 0, 0, 0, fmt.Errorf("failed to call SelectOutput; %v", err)
	}

	if len(objects) > 1 {
		return nil, 0, 0, 0, 0, fmt.Errorf("attempt to call SelectOutput returned more than 1 row")
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

func handleGetOutputs(arguments *server.SelectManyArguments, db *pgxpool.Pool) ([]*Output, int64, int64, int64, int64, error) {
	tx, err := db.Begin(arguments.Ctx)
	if err != nil {
		return nil, 0, 0, 0, 0, err
	}

	defer func() {
		_ = tx.Rollback(arguments.Ctx)
	}()

	objects, count, totalCount, page, totalPages, err := SelectOutputs(arguments.Ctx, tx, arguments.Where, arguments.OrderBy, arguments.Limit, arguments.Offset, arguments.Values...)
	if err != nil {
		return nil, 0, 0, 0, 0, err
	}

	err = tx.Commit(arguments.Ctx)
	if err != nil {
		return nil, 0, 0, 0, 0, err
	}

	return objects, count, totalCount, page, totalPages, nil
}

func handleGetOutput(arguments *server.SelectOneArguments, db *pgxpool.Pool, primaryKey uuid.UUID) ([]*Output, int64, int64, int64, int64, error) {
	tx, err := db.Begin(arguments.Ctx)
	if err != nil {
		return nil, 0, 0, 0, 0, err
	}

	defer func() {
		_ = tx.Rollback(arguments.Ctx)
	}()

	object, count, totalCount, page, totalPages, err := SelectOutput(arguments.Ctx, tx, arguments.Where, arguments.Values...)
	if err != nil {
		return nil, 0, 0, 0, 0, err
	}

	err = tx.Commit(arguments.Ctx)
	if err != nil {
		return nil, 0, 0, 0, 0, err
	}

	return []*Output{object}, count, totalCount, page, totalPages, nil
}

func handlePostOutput(arguments *server.LoadArguments, db *pgxpool.Pool, waitForChange server.WaitForChange, objects []*Output, forceSetValuesForFieldsByObjectIndex [][]string) ([]*Output, int64, int64, int64, int64, error) {
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
		_, err := waitForChange(arguments.Ctx, []stream.Action{stream.INSERT}, OutputTable, xid)
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

func handlePutOutput(arguments *server.LoadArguments, db *pgxpool.Pool, waitForChange server.WaitForChange, object *Output) ([]*Output, int64, int64, int64, int64, error) {
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
		_, err := waitForChange(arguments.Ctx, []stream.Action{stream.UPDATE, stream.SOFT_DELETE, stream.SOFT_RESTORE, stream.SOFT_UPDATE}, OutputTable, xid)
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

	return []*Output{object}, count, totalCount, page, totalPages, nil
}

func handlePatchOutput(arguments *server.LoadArguments, db *pgxpool.Pool, waitForChange server.WaitForChange, object *Output, forceSetValuesForFields []string) ([]*Output, int64, int64, int64, int64, error) {
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
		_, err := waitForChange(arguments.Ctx, []stream.Action{stream.UPDATE, stream.SOFT_DELETE, stream.SOFT_RESTORE, stream.SOFT_UPDATE}, OutputTable, xid)
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

	return []*Output{object}, count, totalCount, page, totalPages, nil
}

func handleDeleteOutput(arguments *server.LoadArguments, db *pgxpool.Pool, waitForChange server.WaitForChange, object *Output) error {
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
		_, err := waitForChange(arguments.Ctx, []stream.Action{stream.DELETE, stream.SOFT_DELETE}, OutputTable, xid)
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

func MutateRouterForOutput(r chi.Router, db *pgxpool.Pool, redisPool *redis.Pool, objectMiddlewares []server.ObjectMiddleware, waitForChange server.WaitForChange) {

	func() {
		getManyHandler, err := getHTTPHandler(
			http.MethodGet,
			"/outputs",
			http.StatusOK,
			func(
				ctx context.Context,
				pathParams server.EmptyPathParams,
				queryParams map[string]any,
				req server.EmptyRequest,
				rawReq any,
			) (server.Response[Output], error) {
				before := time.Now()

				redisConn := redisPool.Get()
				defer func() {
					_ = redisConn.Close()
				}()

				arguments, err := server.GetSelectManyArguments(ctx, queryParams, OutputIntrospectedTable, nil, nil)
				if err != nil {
					if config.Debug() {
						log.Printf("request cache not yet reached; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return server.Response[Output]{}, err
				}

				cachedResponseAsJSON, cacheHit, err := server.GetCachedResponseAsJSON(arguments.RequestHash, redisConn)
				if err != nil {
					if config.Debug() {
						log.Printf("request cache failed; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return server.Response[Output]{}, err
				}

				if cacheHit {
					var cachedResponse server.Response[Output]

					/* TODO: it'd be nice to be able to avoid this (i.e. just pass straight through) */
					err = json.Unmarshal(cachedResponseAsJSON, &cachedResponse)
					if err != nil {
						if config.Debug() {
							log.Printf("request cache hit but failed unmarshal; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
						}

						return server.Response[Output]{}, err
					}

					if config.Debug() {
						log.Printf("request cache hit; request succeeded in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return cachedResponse, nil
				}

				objects, count, totalCount, _, _, err := handleGetOutputs(arguments, db)
				if err != nil {
					if config.Debug() {
						log.Printf("request cache missed; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return server.Response[Output]{}, err
				}

				limit := int64(0)
				if arguments.Limit != nil {
					limit = int64(*arguments.Limit)
				}

				offset := int64(0)
				if arguments.Offset != nil {
					offset = int64(*arguments.Offset)
				}

				response := server.Response[Output]{
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

					return server.Response[Output]{}, err
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
			Output{},
			OutputIntrospectedTable,
		)
		if err != nil {
			panic(err)
		}
		r.Get(getManyHandler.FullPath, getManyHandler.ServeHTTP)
	}()

	func() {
		getOneHandler, err := getHTTPHandler(
			http.MethodGet,
			"/outputs/{primaryKey}",
			http.StatusOK,
			func(
				ctx context.Context,
				pathParams OutputOnePathParams,
				queryParams OutputLoadQueryParams,
				req server.EmptyRequest,
				rawReq any,
			) (server.Response[Output], error) {
				before := time.Now()

				redisConn := redisPool.Get()
				defer func() {
					_ = redisConn.Close()
				}()

				arguments, err := server.GetSelectOneArguments(ctx, queryParams.Depth, OutputIntrospectedTable, pathParams.PrimaryKey, nil, nil)
				if err != nil {
					if config.Debug() {
						log.Printf("request cache not yet reached; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return server.Response[Output]{}, err
				}

				cachedResponseAsJSON, cacheHit, err := server.GetCachedResponseAsJSON(arguments.RequestHash, redisConn)
				if err != nil {
					if config.Debug() {
						log.Printf("request cache failed; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return server.Response[Output]{}, err
				}

				if cacheHit {
					var cachedResponse server.Response[Output]

					/* TODO: it'd be nice to be able to avoid this (i.e. just pass straight through) */
					err = json.Unmarshal(cachedResponseAsJSON, &cachedResponse)
					if err != nil {
						if config.Debug() {
							log.Printf("request cache hit but failed unmarshal; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
						}

						return server.Response[Output]{}, err
					}

					if config.Debug() {
						log.Printf("request cache hit; request succeeded in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return cachedResponse, nil
				}

				objects, count, totalCount, _, _, err := handleGetOutput(arguments, db, pathParams.PrimaryKey)
				if err != nil {
					if config.Debug() {
						log.Printf("request cache missed; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return server.Response[Output]{}, err
				}

				limit := int64(0)

				offset := int64(0)

				response := server.Response[Output]{
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

					return server.Response[Output]{}, err
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
			Output{},
			OutputIntrospectedTable,
		)
		if err != nil {
			panic(err)
		}
		r.Get(getOneHandler.FullPath, getOneHandler.ServeHTTP)
	}()

	func() {
		postHandler, err := getHTTPHandler(
			http.MethodPost,
			"/outputs",
			http.StatusCreated,
			func(
				ctx context.Context,
				pathParams server.EmptyPathParams,
				queryParams OutputLoadQueryParams,
				req []*Output,
				rawReq any,
			) (server.Response[Output], error) {
				allRawItems, ok := rawReq.([]any)
				if !ok {
					return server.Response[Output]{}, fmt.Errorf("failed to cast %#+v to []map[string]any", rawReq)
				}

				allItems := make([]map[string]any, 0)
				for _, rawItem := range allRawItems {
					item, ok := rawItem.(map[string]any)
					if !ok {
						return server.Response[Output]{}, fmt.Errorf("failed to cast %#+v to map[string]any", rawItem)
					}

					allItems = append(allItems, item)
				}

				forceSetValuesForFieldsByObjectIndex := make([][]string, 0)
				for _, item := range allItems {
					forceSetValuesForFields := make([]string, 0)
					for _, possibleField := range maps.Keys(item) {
						if !slices.Contains(OutputTableColumns, possibleField) {
							continue
						}

						forceSetValuesForFields = append(forceSetValuesForFields, possibleField)
					}
					forceSetValuesForFieldsByObjectIndex = append(forceSetValuesForFieldsByObjectIndex, forceSetValuesForFields)
				}

				arguments, err := server.GetLoadArguments(ctx, queryParams.Depth)
				if err != nil {
					return server.Response[Output]{}, err
				}

				objects, count, totalCount, _, _, err := handlePostOutput(arguments, db, waitForChange, req, forceSetValuesForFieldsByObjectIndex)
				if err != nil {
					return server.Response[Output]{}, err
				}

				limit := int64(0)

				offset := int64(0)

				return server.Response[Output]{
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
			Output{},
			OutputIntrospectedTable,
		)
		if err != nil {
			panic(err)
		}
		r.Post(postHandler.FullPath, postHandler.ServeHTTP)
	}()

	func() {
		putHandler, err := getHTTPHandler(
			http.MethodPatch,
			"/outputs/{primaryKey}",
			http.StatusOK,
			func(
				ctx context.Context,
				pathParams OutputOnePathParams,
				queryParams OutputLoadQueryParams,
				req Output,
				rawReq any,
			) (server.Response[Output], error) {
				item, ok := rawReq.(map[string]any)
				if !ok {
					return server.Response[Output]{}, fmt.Errorf("failed to cast %#+v to map[string]any", item)
				}

				arguments, err := server.GetLoadArguments(ctx, queryParams.Depth)
				if err != nil {
					return server.Response[Output]{}, err
				}

				object := &req
				object.ID = pathParams.PrimaryKey

				objects, count, totalCount, _, _, err := handlePutOutput(arguments, db, waitForChange, object)
				if err != nil {
					return server.Response[Output]{}, err
				}

				limit := int64(0)

				offset := int64(0)

				return server.Response[Output]{
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
			Output{},
			OutputIntrospectedTable,
		)
		if err != nil {
			panic(err)
		}
		r.Put(putHandler.FullPath, putHandler.ServeHTTP)
	}()

	func() {
		patchHandler, err := getHTTPHandler(
			http.MethodPatch,
			"/outputs/{primaryKey}",
			http.StatusOK,
			func(
				ctx context.Context,
				pathParams OutputOnePathParams,
				queryParams OutputLoadQueryParams,
				req Output,
				rawReq any,
			) (server.Response[Output], error) {
				item, ok := rawReq.(map[string]any)
				if !ok {
					return server.Response[Output]{}, fmt.Errorf("failed to cast %#+v to map[string]any", item)
				}

				forceSetValuesForFields := make([]string, 0)
				for _, possibleField := range maps.Keys(item) {
					if !slices.Contains(OutputTableColumns, possibleField) {
						continue
					}

					forceSetValuesForFields = append(forceSetValuesForFields, possibleField)
				}

				arguments, err := server.GetLoadArguments(ctx, queryParams.Depth)
				if err != nil {
					return server.Response[Output]{}, err
				}

				object := &req
				object.ID = pathParams.PrimaryKey

				objects, count, totalCount, _, _, err := handlePatchOutput(arguments, db, waitForChange, object, forceSetValuesForFields)
				if err != nil {
					return server.Response[Output]{}, err
				}

				limit := int64(0)

				offset := int64(0)

				return server.Response[Output]{
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
			Output{},
			OutputIntrospectedTable,
		)
		if err != nil {
			panic(err)
		}
		r.Patch(patchHandler.FullPath, patchHandler.ServeHTTP)
	}()

	func() {
		deleteHandler, err := getHTTPHandler(
			http.MethodDelete,
			"/outputs/{primaryKey}",
			http.StatusNoContent,
			func(
				ctx context.Context,
				pathParams OutputOnePathParams,
				queryParams OutputLoadQueryParams,
				req server.EmptyRequest,
				rawReq any,
			) (server.EmptyResponse, error) {
				arguments, err := server.GetLoadArguments(ctx, queryParams.Depth)
				if err != nil {
					return server.EmptyResponse{}, err
				}

				object := &Output{}
				object.ID = pathParams.PrimaryKey

				err = handleDeleteOutput(arguments, db, waitForChange, object)
				if err != nil {
					return server.EmptyResponse{}, err
				}

				return server.EmptyResponse{}, nil
			},
			Output{},
			OutputIntrospectedTable,
		)
		if err != nil {
			panic(err)
		}
		r.Delete(deleteHandler.FullPath, deleteHandler.ServeHTTP)
	}()
}

func NewOutputFromItem(item map[string]any) (any, error) {
	object := &Output{}

	err := object.FromItem(item)
	if err != nil {
		return nil, err
	}

	return object, nil
}

func init() {
	register(
		OutputTable,
		Output{},
		NewOutputFromItem,
		"/outputs",
		MutateRouterForOutput,
	)
}
