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

type Change struct {
	ID                 uuid.UUID   `json:"id"`
	CreatedAt          time.Time   `json:"created_at"`
	UpdatedAt          time.Time   `json:"updated_at"`
	DeletedAt          *time.Time  `json:"deleted_at"`
	CommitHash         string      `json:"commit_hash"`
	BranchName         string      `json:"branch_name"`
	Message            string      `json:"message"`
	AuthoredBy         string      `json:"authored_by"`
	AuthoredAt         time.Time   `json:"authored_at"`
	CommittedBy        string      `json:"committed_by"`
	CommittedAt        time.Time   `json:"committed_at"`
	TriggersProducedAt *time.Time  `json:"triggers_produced_at"`
	RepositoryID       uuid.UUID   `json:"repository_id"`
	RepositoryIDObject *Repository `json:"repository_id_object"`
}

var ChangeTable = "change"

var ChangeTableNamespaceID int32 = 1337 + 1

var (
	ChangeTableIDColumn                 = "id"
	ChangeTableCreatedAtColumn          = "created_at"
	ChangeTableUpdatedAtColumn          = "updated_at"
	ChangeTableDeletedAtColumn          = "deleted_at"
	ChangeTableCommitHashColumn         = "commit_hash"
	ChangeTableBranchNameColumn         = "branch_name"
	ChangeTableMessageColumn            = "message"
	ChangeTableAuthoredByColumn         = "authored_by"
	ChangeTableAuthoredAtColumn         = "authored_at"
	ChangeTableCommittedByColumn        = "committed_by"
	ChangeTableCommittedAtColumn        = "committed_at"
	ChangeTableTriggersProducedAtColumn = "triggers_produced_at"
	ChangeTableRepositoryIDColumn       = "repository_id"
)

var (
	ChangeTableIDColumnWithTypeCast                 = `"id" AS id`
	ChangeTableCreatedAtColumnWithTypeCast          = `"created_at" AS created_at`
	ChangeTableUpdatedAtColumnWithTypeCast          = `"updated_at" AS updated_at`
	ChangeTableDeletedAtColumnWithTypeCast          = `"deleted_at" AS deleted_at`
	ChangeTableCommitHashColumnWithTypeCast         = `"commit_hash" AS commit_hash`
	ChangeTableBranchNameColumnWithTypeCast         = `"branch_name" AS branch_name`
	ChangeTableMessageColumnWithTypeCast            = `"message" AS message`
	ChangeTableAuthoredByColumnWithTypeCast         = `"authored_by" AS authored_by`
	ChangeTableAuthoredAtColumnWithTypeCast         = `"authored_at" AS authored_at`
	ChangeTableCommittedByColumnWithTypeCast        = `"committed_by" AS committed_by`
	ChangeTableCommittedAtColumnWithTypeCast        = `"committed_at" AS committed_at`
	ChangeTableTriggersProducedAtColumnWithTypeCast = `"triggers_produced_at" AS triggers_produced_at`
	ChangeTableRepositoryIDColumnWithTypeCast       = `"repository_id" AS repository_id`
)

var ChangeTableColumns = []string{
	ChangeTableIDColumn,
	ChangeTableCreatedAtColumn,
	ChangeTableUpdatedAtColumn,
	ChangeTableDeletedAtColumn,
	ChangeTableCommitHashColumn,
	ChangeTableBranchNameColumn,
	ChangeTableMessageColumn,
	ChangeTableAuthoredByColumn,
	ChangeTableAuthoredAtColumn,
	ChangeTableCommittedByColumn,
	ChangeTableCommittedAtColumn,
	ChangeTableTriggersProducedAtColumn,
	ChangeTableRepositoryIDColumn,
}

var ChangeTableColumnsWithTypeCasts = []string{
	ChangeTableIDColumnWithTypeCast,
	ChangeTableCreatedAtColumnWithTypeCast,
	ChangeTableUpdatedAtColumnWithTypeCast,
	ChangeTableDeletedAtColumnWithTypeCast,
	ChangeTableCommitHashColumnWithTypeCast,
	ChangeTableBranchNameColumnWithTypeCast,
	ChangeTableMessageColumnWithTypeCast,
	ChangeTableAuthoredByColumnWithTypeCast,
	ChangeTableAuthoredAtColumnWithTypeCast,
	ChangeTableCommittedByColumnWithTypeCast,
	ChangeTableCommittedAtColumnWithTypeCast,
	ChangeTableTriggersProducedAtColumnWithTypeCast,
	ChangeTableRepositoryIDColumnWithTypeCast,
}

var ChangeIntrospectedTable *introspect.Table

var ChangeTableColumnLookup map[string]*introspect.Column

var (
	ChangeTablePrimaryKeyColumn = ChangeTableIDColumn
)

func init() {
	ChangeIntrospectedTable = tableByName[ChangeTable]

	/* only needed during templating */
	if ChangeIntrospectedTable == nil {
		ChangeIntrospectedTable = &introspect.Table{}
	}

	ChangeTableColumnLookup = ChangeIntrospectedTable.ColumnByName
}

type ChangeOnePathParams struct {
	PrimaryKey uuid.UUID `json:"primaryKey"`
}

type ChangeLoadQueryParams struct {
	Depth *int `json:"depth"`
}

type ChangeClaimRequest struct {
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

func (m *Change) GetPrimaryKeyColumn() string {
	return ChangeTablePrimaryKeyColumn
}

func (m *Change) GetPrimaryKeyValue() any {
	return m.ID
}

func (m *Change) FromItem(item map[string]any) error {
	if item == nil {
		return fmt.Errorf(
			"item unexpectedly nil during ChangeFromItem",
		)
	}

	if len(item) == 0 {
		return fmt.Errorf(
			"item unexpectedly empty during ChangeFromItem",
		)
	}

	wrapError := func(k string, v any, err error) error {
		return fmt.Errorf("%v: %#+v; error; %v", k, v, err)
	}

	for k, v := range item {
		_, ok := ChangeTableColumnLookup[k]
		if !ok {
			return fmt.Errorf(
				"item contained unexpected key %#+v during ChangeFromItem; item: %#+v",
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

		case "commit_hash":
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
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uucommit_hash.UUID", temp1))
				}
			}

			m.CommitHash = temp2

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

			m.BranchName = temp2

		case "message":
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
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uumessage.UUID", temp1))
				}
			}

			m.Message = temp2

		case "authored_by":
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
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uuauthored_by.UUID", temp1))
				}
			}

			m.AuthoredBy = temp2

		case "authored_at":
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
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uuauthored_at.UUID", temp1))
				}
			}

			m.AuthoredAt = temp2

		case "committed_by":
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
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uucommitted_by.UUID", temp1))
				}
			}

			m.CommittedBy = temp2

		case "committed_at":
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
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uucommitted_at.UUID", temp1))
				}
			}

			m.CommittedAt = temp2

		case "triggers_produced_at":
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
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uutriggers_produced_at.UUID", temp1))
				}
			}

			m.TriggersProducedAt = &temp2

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

func (m *Change) Reload(ctx context.Context, tx pgx.Tx, includeDeleteds ...bool) error {
	extraWhere := ""
	if len(includeDeleteds) > 0 && includeDeleteds[0] {
		if slices.Contains(ChangeTableColumns, "deleted_at") {
			extraWhere = "\n    AND (deleted_at IS null OR deleted_at IS NOT null)"
		}
	}

	ctx, cleanup := query.WithQueryID(ctx)
	defer cleanup()

	ctx = query.WithMaxDepth(ctx, nil)

	o, _, _, _, _, err := SelectChange(
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
	m.CommitHash = o.CommitHash
	m.BranchName = o.BranchName
	m.Message = o.Message
	m.AuthoredBy = o.AuthoredBy
	m.AuthoredAt = o.AuthoredAt
	m.CommittedBy = o.CommittedBy
	m.CommittedAt = o.CommittedAt
	m.TriggersProducedAt = o.TriggersProducedAt
	m.RepositoryID = o.RepositoryID
	m.RepositoryIDObject = o.RepositoryIDObject

	return nil
}

func (m *Change) Insert(ctx context.Context, tx pgx.Tx, setPrimaryKey bool, setZeroValues bool, forceSetValuesForFields ...string) error {
	columns := make([]string, 0)
	values := make([]any, 0)

	if setPrimaryKey && (setZeroValues || !types.IsZeroUUID(m.ID) || slices.Contains(forceSetValuesForFields, ChangeTableIDColumn) || isRequired(ChangeTableColumnLookup, ChangeTableIDColumn)) {
		columns = append(columns, ChangeTableIDColumn)

		v, err := types.FormatUUID(m.ID)
		if err != nil {
			return fmt.Errorf("failed to handle m.ID; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.CreatedAt) || slices.Contains(forceSetValuesForFields, ChangeTableCreatedAtColumn) || isRequired(ChangeTableColumnLookup, ChangeTableCreatedAtColumn) {
		columns = append(columns, ChangeTableCreatedAtColumn)

		v, err := types.FormatTime(m.CreatedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.CreatedAt; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.UpdatedAt) || slices.Contains(forceSetValuesForFields, ChangeTableUpdatedAtColumn) || isRequired(ChangeTableColumnLookup, ChangeTableUpdatedAtColumn) {
		columns = append(columns, ChangeTableUpdatedAtColumn)

		v, err := types.FormatTime(m.UpdatedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.UpdatedAt; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.DeletedAt) || slices.Contains(forceSetValuesForFields, ChangeTableDeletedAtColumn) || isRequired(ChangeTableColumnLookup, ChangeTableDeletedAtColumn) {
		columns = append(columns, ChangeTableDeletedAtColumn)

		v, err := types.FormatTime(m.DeletedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.DeletedAt; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroString(m.CommitHash) || slices.Contains(forceSetValuesForFields, ChangeTableCommitHashColumn) || isRequired(ChangeTableColumnLookup, ChangeTableCommitHashColumn) {
		columns = append(columns, ChangeTableCommitHashColumn)

		v, err := types.FormatString(m.CommitHash)
		if err != nil {
			return fmt.Errorf("failed to handle m.CommitHash; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroString(m.BranchName) || slices.Contains(forceSetValuesForFields, ChangeTableBranchNameColumn) || isRequired(ChangeTableColumnLookup, ChangeTableBranchNameColumn) {
		columns = append(columns, ChangeTableBranchNameColumn)

		v, err := types.FormatString(m.BranchName)
		if err != nil {
			return fmt.Errorf("failed to handle m.BranchName; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroString(m.Message) || slices.Contains(forceSetValuesForFields, ChangeTableMessageColumn) || isRequired(ChangeTableColumnLookup, ChangeTableMessageColumn) {
		columns = append(columns, ChangeTableMessageColumn)

		v, err := types.FormatString(m.Message)
		if err != nil {
			return fmt.Errorf("failed to handle m.Message; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroString(m.AuthoredBy) || slices.Contains(forceSetValuesForFields, ChangeTableAuthoredByColumn) || isRequired(ChangeTableColumnLookup, ChangeTableAuthoredByColumn) {
		columns = append(columns, ChangeTableAuthoredByColumn)

		v, err := types.FormatString(m.AuthoredBy)
		if err != nil {
			return fmt.Errorf("failed to handle m.AuthoredBy; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.AuthoredAt) || slices.Contains(forceSetValuesForFields, ChangeTableAuthoredAtColumn) || isRequired(ChangeTableColumnLookup, ChangeTableAuthoredAtColumn) {
		columns = append(columns, ChangeTableAuthoredAtColumn)

		v, err := types.FormatTime(m.AuthoredAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.AuthoredAt; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroString(m.CommittedBy) || slices.Contains(forceSetValuesForFields, ChangeTableCommittedByColumn) || isRequired(ChangeTableColumnLookup, ChangeTableCommittedByColumn) {
		columns = append(columns, ChangeTableCommittedByColumn)

		v, err := types.FormatString(m.CommittedBy)
		if err != nil {
			return fmt.Errorf("failed to handle m.CommittedBy; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.CommittedAt) || slices.Contains(forceSetValuesForFields, ChangeTableCommittedAtColumn) || isRequired(ChangeTableColumnLookup, ChangeTableCommittedAtColumn) {
		columns = append(columns, ChangeTableCommittedAtColumn)

		v, err := types.FormatTime(m.CommittedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.CommittedAt; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.TriggersProducedAt) || slices.Contains(forceSetValuesForFields, ChangeTableTriggersProducedAtColumn) || isRequired(ChangeTableColumnLookup, ChangeTableTriggersProducedAtColumn) {
		columns = append(columns, ChangeTableTriggersProducedAtColumn)

		v, err := types.FormatTime(m.TriggersProducedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.TriggersProducedAt; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroUUID(m.RepositoryID) || slices.Contains(forceSetValuesForFields, ChangeTableRepositoryIDColumn) || isRequired(ChangeTableColumnLookup, ChangeTableRepositoryIDColumn) {
		columns = append(columns, ChangeTableRepositoryIDColumn)

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
		ChangeTable,
		columns,
		nil,
		false,
		false,
		ChangeTableColumns,
		values...,
	)
	if err != nil {
		return fmt.Errorf("failed to insert %#+v; %v", m, err)
	}
	v := (*item)[ChangeTableIDColumn]

	if v == nil {
		return fmt.Errorf("failed to find %v in %#+v", ChangeTableIDColumn, item)
	}

	wrapError := func(err error) error {
		return fmt.Errorf(
			"failed to treat %v: %#+v as uuid.UUID: %v",
			ChangeTableIDColumn,
			(*item)[ChangeTableIDColumn],
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

func (m *Change) Update(ctx context.Context, tx pgx.Tx, setZeroValues bool, forceSetValuesForFields ...string) error {
	columns := make([]string, 0)
	values := make([]any, 0)

	if setZeroValues || !types.IsZeroTime(m.CreatedAt) || slices.Contains(forceSetValuesForFields, ChangeTableCreatedAtColumn) {
		columns = append(columns, ChangeTableCreatedAtColumn)

		v, err := types.FormatTime(m.CreatedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.CreatedAt; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.UpdatedAt) || slices.Contains(forceSetValuesForFields, ChangeTableUpdatedAtColumn) {
		columns = append(columns, ChangeTableUpdatedAtColumn)

		v, err := types.FormatTime(m.UpdatedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.UpdatedAt; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.DeletedAt) || slices.Contains(forceSetValuesForFields, ChangeTableDeletedAtColumn) {
		columns = append(columns, ChangeTableDeletedAtColumn)

		v, err := types.FormatTime(m.DeletedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.DeletedAt; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroString(m.CommitHash) || slices.Contains(forceSetValuesForFields, ChangeTableCommitHashColumn) {
		columns = append(columns, ChangeTableCommitHashColumn)

		v, err := types.FormatString(m.CommitHash)
		if err != nil {
			return fmt.Errorf("failed to handle m.CommitHash; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroString(m.BranchName) || slices.Contains(forceSetValuesForFields, ChangeTableBranchNameColumn) {
		columns = append(columns, ChangeTableBranchNameColumn)

		v, err := types.FormatString(m.BranchName)
		if err != nil {
			return fmt.Errorf("failed to handle m.BranchName; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroString(m.Message) || slices.Contains(forceSetValuesForFields, ChangeTableMessageColumn) {
		columns = append(columns, ChangeTableMessageColumn)

		v, err := types.FormatString(m.Message)
		if err != nil {
			return fmt.Errorf("failed to handle m.Message; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroString(m.AuthoredBy) || slices.Contains(forceSetValuesForFields, ChangeTableAuthoredByColumn) {
		columns = append(columns, ChangeTableAuthoredByColumn)

		v, err := types.FormatString(m.AuthoredBy)
		if err != nil {
			return fmt.Errorf("failed to handle m.AuthoredBy; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.AuthoredAt) || slices.Contains(forceSetValuesForFields, ChangeTableAuthoredAtColumn) {
		columns = append(columns, ChangeTableAuthoredAtColumn)

		v, err := types.FormatTime(m.AuthoredAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.AuthoredAt; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroString(m.CommittedBy) || slices.Contains(forceSetValuesForFields, ChangeTableCommittedByColumn) {
		columns = append(columns, ChangeTableCommittedByColumn)

		v, err := types.FormatString(m.CommittedBy)
		if err != nil {
			return fmt.Errorf("failed to handle m.CommittedBy; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.CommittedAt) || slices.Contains(forceSetValuesForFields, ChangeTableCommittedAtColumn) {
		columns = append(columns, ChangeTableCommittedAtColumn)

		v, err := types.FormatTime(m.CommittedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.CommittedAt; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.TriggersProducedAt) || slices.Contains(forceSetValuesForFields, ChangeTableTriggersProducedAtColumn) {
		columns = append(columns, ChangeTableTriggersProducedAtColumn)

		v, err := types.FormatTime(m.TriggersProducedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.TriggersProducedAt; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroUUID(m.RepositoryID) || slices.Contains(forceSetValuesForFields, ChangeTableRepositoryIDColumn) {
		columns = append(columns, ChangeTableRepositoryIDColumn)

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
		ChangeTable,
		columns,
		fmt.Sprintf("%v = $$??", ChangeTableIDColumn),
		ChangeTableColumns,
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

func (m *Change) Delete(ctx context.Context, tx pgx.Tx, hardDeletes ...bool) error {
	hardDelete := false
	if len(hardDeletes) > 0 {
		hardDelete = hardDeletes[0]
	}

	if !hardDelete && slices.Contains(ChangeTableColumns, "deleted_at") {
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
		ChangeTable,
		fmt.Sprintf("%v = $$??", ChangeTableIDColumn),
		values...,
	)
	if err != nil {
		return fmt.Errorf("failed to delete %#+v; %v", m, err)
	}

	_ = m.Reload(ctx, tx, true)

	return nil
}

func (m *Change) LockTable(ctx context.Context, tx pgx.Tx, timeouts ...time.Duration) error {
	return query.LockTable(ctx, tx, ChangeTable, timeouts...)
}

func (m *Change) LockTableWithRetries(ctx context.Context, tx pgx.Tx, overallTimeout time.Duration, individualAttempttimeout time.Duration) error {
	return query.LockTableWithRetries(ctx, tx, ChangeTable, overallTimeout, individualAttempttimeout)
}

func (m *Change) AdvisoryLock(ctx context.Context, tx pgx.Tx, key int32, timeouts ...time.Duration) error {
	return query.AdvisoryLock(ctx, tx, ChangeTableNamespaceID, key, timeouts...)
}

func (m *Change) AdvisoryLockWithRetries(ctx context.Context, tx pgx.Tx, key int32, overallTimeout time.Duration, individualAttempttimeout time.Duration) error {
	return query.AdvisoryLockWithRetries(ctx, tx, ChangeTableNamespaceID, key, overallTimeout, individualAttempttimeout)
}

func (m *Change) Claim(ctx context.Context, tx pgx.Tx, until time.Time, by uuid.UUID, timeout time.Duration) error {
	if !(slices.Contains(ChangeTableColumns, "claimed_until") && slices.Contains(ChangeTableColumns, "claimed_by")) {
		return fmt.Errorf("can only invoke Claim for tables with 'claimed_until' and 'claimed_by' columns")
	}

	err := m.AdvisoryLockWithRetries(ctx, tx, math.MinInt32, timeout, time.Second*1)
	if err != nil {
		return fmt.Errorf("failed to claim (advisory lock): %s", err.Error())
	}

	x, _, _, _, _, err := SelectChange(
		ctx,
		tx,
		fmt.Sprintf(
			"%s = $$?? AND (claimed_by = $$?? OR (claimed_until IS null OR claimed_until < now()))",
			ChangeTablePrimaryKeyColumn,
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

func SelectChanges(ctx context.Context, tx pgx.Tx, where string, orderBy *string, limit *int, offset *int, values ...any) ([]*Change, int64, int64, int64, int64, error) {
	before := time.Now()

	if config.Debug() {
		log.Printf("entered SelectChanges")

		defer func() {
			log.Printf("exited SelectChanges in %s", time.Since(before))
		}()
	}
	if slices.Contains(ChangeTableColumns, "deleted_at") {
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

	shouldLoad := query.ShouldLoad(ctx, ChangeTable) || query.ShouldLoad(ctx, fmt.Sprintf("referenced_by_%s", ChangeTable))

	var ok bool
	ctx, ok = query.HandleQueryPathGraphCycles(ctx, fmt.Sprintf("%s{%v}", ChangeTable, nil), !isLoadQuery)
	if !ok && !shouldLoad {
		if config.Debug() {
			log.Printf("skipping SelectChange early (query.ShouldLoad(): %v, query.HandleQueryPathGraphCycles(): %v)", shouldLoad, ok)
		}
		return []*Change{}, 0, 0, 0, 0, nil
	}

	items, count, totalCount, page, totalPages, err := query.Select(
		ctx,
		tx,
		ChangeTableColumnsWithTypeCasts,
		ChangeTable,
		where,
		orderBy,
		limit,
		offset,
		values...,
	)
	if err != nil {
		return nil, 0, 0, 0, 0, fmt.Errorf("failed to call SelectChanges; %v", err)
	}

	objects := make([]*Change, 0)

	for _, item := range *items {
		object := &Change{}

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
					log.Printf("loading SelectChanges->SelectRepository for object.RepositoryIDObject{%s: %v}", RepositoryTablePrimaryKeyColumn, object.RepositoryID)
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
					log.Printf("loaded SelectChanges->SelectRepository for object.RepositoryIDObject in %s", time.Since(thisBefore))
				}
			}
		}

		objects = append(objects, object)
	}

	return objects, count, totalCount, page, totalPages, nil
}

func SelectChange(ctx context.Context, tx pgx.Tx, where string, values ...any) (*Change, int64, int64, int64, int64, error) {
	ctx, cleanup := query.WithQueryID(ctx)
	defer cleanup()

	ctx = query.WithMaxDepth(ctx, nil)

	objects, _, _, _, _, err := SelectChanges(
		ctx,
		tx,
		where,
		nil,
		helpers.Ptr(2),
		helpers.Ptr(0),
		values...,
	)
	if err != nil {
		return nil, 0, 0, 0, 0, fmt.Errorf("failed to call SelectChange; %v", err)
	}

	if len(objects) > 1 {
		return nil, 0, 0, 0, 0, fmt.Errorf("attempt to call SelectChange returned more than 1 row")
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

func ClaimChange(ctx context.Context, tx pgx.Tx, until time.Time, by uuid.UUID, timeout time.Duration, wheres ...string) (*Change, error) {
	if !(slices.Contains(ChangeTableColumns, "claimed_until") && slices.Contains(ChangeTableColumns, "claimed_by")) {
		return nil, fmt.Errorf("can only invoke Claim for tables with 'claimed_until' and 'claimed_by' columns")
	}

	m := &Change{}

	err := m.AdvisoryLockWithRetries(ctx, tx, math.MinInt32, timeout, time.Second*1)
	if err != nil {
		return nil, fmt.Errorf("failed to claim: %s", err.Error())
	}

	extraWhere := ""
	if len(wheres) > 0 {
		extraWhere = fmt.Sprintf("AND %s", extraWhere)
	}

	ms, _, _, _, _, err := SelectChanges(
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

func handleGetChanges(arguments *server.SelectManyArguments, db *pgxpool.Pool) ([]*Change, int64, int64, int64, int64, error) {
	tx, err := db.Begin(arguments.Ctx)
	if err != nil {
		return nil, 0, 0, 0, 0, err
	}

	defer func() {
		_ = tx.Rollback(arguments.Ctx)
	}()

	objects, count, totalCount, page, totalPages, err := SelectChanges(arguments.Ctx, tx, arguments.Where, arguments.OrderBy, arguments.Limit, arguments.Offset, arguments.Values...)
	if err != nil {
		return nil, 0, 0, 0, 0, err
	}

	err = tx.Commit(arguments.Ctx)
	if err != nil {
		return nil, 0, 0, 0, 0, err
	}

	return objects, count, totalCount, page, totalPages, nil
}

func handleGetChange(arguments *server.SelectOneArguments, db *pgxpool.Pool, primaryKey uuid.UUID) ([]*Change, int64, int64, int64, int64, error) {
	tx, err := db.Begin(arguments.Ctx)
	if err != nil {
		return nil, 0, 0, 0, 0, err
	}

	defer func() {
		_ = tx.Rollback(arguments.Ctx)
	}()

	object, count, totalCount, page, totalPages, err := SelectChange(arguments.Ctx, tx, arguments.Where, arguments.Values...)
	if err != nil {
		return nil, 0, 0, 0, 0, err
	}

	err = tx.Commit(arguments.Ctx)
	if err != nil {
		return nil, 0, 0, 0, 0, err
	}

	return []*Change{object}, count, totalCount, page, totalPages, nil
}

func handlePostChanges(arguments *server.LoadArguments, db *pgxpool.Pool, waitForChange server.WaitForChange, objects []*Change, forceSetValuesForFieldsByObjectIndex [][]string) ([]*Change, int64, int64, int64, int64, error) {
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
		_, err := waitForChange(arguments.Ctx, []stream.Action{stream.INSERT}, ChangeTable, xid)
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

func handlePutChange(arguments *server.LoadArguments, db *pgxpool.Pool, waitForChange server.WaitForChange, object *Change) ([]*Change, int64, int64, int64, int64, error) {
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
		_, err := waitForChange(arguments.Ctx, []stream.Action{stream.UPDATE, stream.SOFT_DELETE, stream.SOFT_RESTORE, stream.SOFT_UPDATE}, ChangeTable, xid)
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

	return []*Change{object}, count, totalCount, page, totalPages, nil
}

func handlePatchChange(arguments *server.LoadArguments, db *pgxpool.Pool, waitForChange server.WaitForChange, object *Change, forceSetValuesForFields []string) ([]*Change, int64, int64, int64, int64, error) {
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
		_, err := waitForChange(arguments.Ctx, []stream.Action{stream.UPDATE, stream.SOFT_DELETE, stream.SOFT_RESTORE, stream.SOFT_UPDATE}, ChangeTable, xid)
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

	return []*Change{object}, count, totalCount, page, totalPages, nil
}

func handleDeleteChange(arguments *server.LoadArguments, db *pgxpool.Pool, waitForChange server.WaitForChange, object *Change) error {
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
		_, err := waitForChange(arguments.Ctx, []stream.Action{stream.DELETE, stream.SOFT_DELETE}, ChangeTable, xid)
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

func MutateRouterForChange(r chi.Router, db *pgxpool.Pool, redisPool *redis.Pool, objectMiddlewares []server.ObjectMiddleware, waitForChange server.WaitForChange) {
	if slices.Contains(ChangeTableColumns, "claimed_until") && slices.Contains(ChangeTableColumns, "claimed_by") {
		func() {
			postHandlerForClaim, err := getHTTPHandler(
				http.MethodPost,
				"/claim-change",
				http.StatusOK,
				func(
					ctx context.Context,
					pathParams server.EmptyPathParams,
					queryParams server.EmptyQueryParams,
					req ChangeClaimRequest,
					rawReq any,
				) (server.Response[Change], error) {
					tx, err := db.Begin(ctx)
					if err != nil {
						return server.Response[Change]{}, err
					}

					defer func() {
						_ = tx.Rollback(ctx)
					}()

					object, err := ClaimChange(ctx, tx, req.Until, req.By, time.Millisecond*time.Duration(req.TimeoutSeconds*1000))
					if err != nil {
						return server.Response[Change]{}, err
					}

					count := int64(0)

					totalCount := int64(0)

					limit := int64(0)

					offset := int64(0)

					if object == nil {
						return server.Response[Change]{
							Status:     http.StatusOK,
							Success:    true,
							Error:      nil,
							Objects:    []*Change{},
							Count:      count,
							TotalCount: totalCount,
							Limit:      limit,
							Offset:     offset,
						}, nil
					}

					err = tx.Commit(ctx)
					if err != nil {
						return server.Response[Change]{}, err
					}

					return server.Response[Change]{
						Status:     http.StatusOK,
						Success:    true,
						Error:      nil,
						Objects:    []*Change{object},
						Count:      count,
						TotalCount: totalCount,
						Limit:      limit,
						Offset:     offset,
					}, nil
				},
				Change{},
				ChangeIntrospectedTable,
			)
			if err != nil {
				panic(err)
			}
			r.Post(postHandlerForClaim.FullPath, postHandlerForClaim.ServeHTTP)

			postHandlerForClaimOne, err := getHTTPHandler(
				http.MethodPost,
				"/changes/{primaryKey}/claim",
				http.StatusOK,
				func(
					ctx context.Context,
					pathParams ChangeOnePathParams,
					queryParams ChangeLoadQueryParams,
					req ChangeClaimRequest,
					rawReq any,
				) (server.Response[Change], error) {
					before := time.Now()

					redisConn := redisPool.Get()
					defer func() {
						_ = redisConn.Close()
					}()

					arguments, err := server.GetSelectOneArguments(ctx, queryParams.Depth, ChangeIntrospectedTable, pathParams.PrimaryKey, nil, nil)
					if err != nil {
						if config.Debug() {
							log.Printf("request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
						}

						return server.Response[Change]{}, err
					}

					/* note: deliberately no attempt at a cache hit */

					var object *Change
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

						object, count, totalCount, _, _, err = SelectChange(arguments.Ctx, tx, arguments.Where, arguments.Values...)
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

						return server.Response[Change]{}, err
					}

					limit := int64(0)

					offset := int64(0)

					response := server.Response[Change]{
						Status:     http.StatusOK,
						Success:    true,
						Error:      nil,
						Objects:    []*Change{object},
						Count:      count,
						TotalCount: totalCount,
						Limit:      limit,
						Offset:     offset,
					}

					return response, nil
				},
				Change{},
				ChangeIntrospectedTable,
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
			"/changes",
			http.StatusOK,
			func(
				ctx context.Context,
				pathParams server.EmptyPathParams,
				queryParams map[string]any,
				req server.EmptyRequest,
				rawReq any,
			) (server.Response[Change], error) {
				before := time.Now()

				redisConn := redisPool.Get()
				defer func() {
					_ = redisConn.Close()
				}()

				arguments, err := server.GetSelectManyArguments(ctx, queryParams, ChangeIntrospectedTable, nil, nil)
				if err != nil {
					if config.Debug() {
						log.Printf("request cache not yet reached; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return server.Response[Change]{}, err
				}

				cachedResponseAsJSON, cacheHit, err := server.GetCachedResponseAsJSON(arguments.RequestHash, redisConn)
				if err != nil {
					if config.Debug() {
						log.Printf("request cache failed; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return server.Response[Change]{}, err
				}

				if cacheHit {
					var cachedResponse server.Response[Change]

					/* TODO: it'd be nice to be able to avoid this (i.e. just pass straight through) */
					err = json.Unmarshal(cachedResponseAsJSON, &cachedResponse)
					if err != nil {
						if config.Debug() {
							log.Printf("request cache hit but failed unmarshal; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
						}

						return server.Response[Change]{}, err
					}

					if config.Debug() {
						log.Printf("request cache hit; request succeeded in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return cachedResponse, nil
				}

				objects, count, totalCount, _, _, err := handleGetChanges(arguments, db)
				if err != nil {
					if config.Debug() {
						log.Printf("request cache missed; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return server.Response[Change]{}, err
				}

				limit := int64(0)
				if arguments.Limit != nil {
					limit = int64(*arguments.Limit)
				}

				offset := int64(0)
				if arguments.Offset != nil {
					offset = int64(*arguments.Offset)
				}

				response := server.Response[Change]{
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

					return server.Response[Change]{}, err
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
			Change{},
			ChangeIntrospectedTable,
		)
		if err != nil {
			panic(err)
		}
		r.Get(getManyHandler.FullPath, getManyHandler.ServeHTTP)
	}()

	func() {
		getOneHandler, err := getHTTPHandler(
			http.MethodGet,
			"/changes/{primaryKey}",
			http.StatusOK,
			func(
				ctx context.Context,
				pathParams ChangeOnePathParams,
				queryParams ChangeLoadQueryParams,
				req server.EmptyRequest,
				rawReq any,
			) (server.Response[Change], error) {
				before := time.Now()

				redisConn := redisPool.Get()
				defer func() {
					_ = redisConn.Close()
				}()

				arguments, err := server.GetSelectOneArguments(ctx, queryParams.Depth, ChangeIntrospectedTable, pathParams.PrimaryKey, nil, nil)
				if err != nil {
					if config.Debug() {
						log.Printf("request cache not yet reached; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return server.Response[Change]{}, err
				}

				cachedResponseAsJSON, cacheHit, err := server.GetCachedResponseAsJSON(arguments.RequestHash, redisConn)
				if err != nil {
					if config.Debug() {
						log.Printf("request cache failed; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return server.Response[Change]{}, err
				}

				if cacheHit {
					var cachedResponse server.Response[Change]

					/* TODO: it'd be nice to be able to avoid this (i.e. just pass straight through) */
					err = json.Unmarshal(cachedResponseAsJSON, &cachedResponse)
					if err != nil {
						if config.Debug() {
							log.Printf("request cache hit but failed unmarshal; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
						}

						return server.Response[Change]{}, err
					}

					if config.Debug() {
						log.Printf("request cache hit; request succeeded in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return cachedResponse, nil
				}

				objects, count, totalCount, _, _, err := handleGetChange(arguments, db, pathParams.PrimaryKey)
				if err != nil {
					if config.Debug() {
						log.Printf("request cache missed; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return server.Response[Change]{}, err
				}

				limit := int64(0)

				offset := int64(0)

				response := server.Response[Change]{
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

					return server.Response[Change]{}, err
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
			Change{},
			ChangeIntrospectedTable,
		)
		if err != nil {
			panic(err)
		}
		r.Get(getOneHandler.FullPath, getOneHandler.ServeHTTP)
	}()

	func() {
		postHandler, err := getHTTPHandler(
			http.MethodPost,
			"/changes",
			http.StatusCreated,
			func(
				ctx context.Context,
				pathParams server.EmptyPathParams,
				queryParams ChangeLoadQueryParams,
				req []*Change,
				rawReq any,
			) (server.Response[Change], error) {
				allRawItems, ok := rawReq.([]any)
				if !ok {
					return server.Response[Change]{}, fmt.Errorf("failed to cast %#+v to []map[string]any", rawReq)
				}

				allItems := make([]map[string]any, 0)
				for _, rawItem := range allRawItems {
					item, ok := rawItem.(map[string]any)
					if !ok {
						return server.Response[Change]{}, fmt.Errorf("failed to cast %#+v to map[string]any", rawItem)
					}

					allItems = append(allItems, item)
				}

				forceSetValuesForFieldsByObjectIndex := make([][]string, 0)
				for _, item := range allItems {
					forceSetValuesForFields := make([]string, 0)
					for _, possibleField := range maps.Keys(item) {
						if !slices.Contains(ChangeTableColumns, possibleField) {
							continue
						}

						forceSetValuesForFields = append(forceSetValuesForFields, possibleField)
					}
					forceSetValuesForFieldsByObjectIndex = append(forceSetValuesForFieldsByObjectIndex, forceSetValuesForFields)
				}

				arguments, err := server.GetLoadArguments(ctx, queryParams.Depth)
				if err != nil {
					return server.Response[Change]{}, err
				}

				objects, count, totalCount, _, _, err := handlePostChanges(arguments, db, waitForChange, req, forceSetValuesForFieldsByObjectIndex)
				if err != nil {
					return server.Response[Change]{}, err
				}

				limit := int64(0)

				offset := int64(0)

				return server.Response[Change]{
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
			Change{},
			ChangeIntrospectedTable,
		)
		if err != nil {
			panic(err)
		}
		r.Post(postHandler.FullPath, postHandler.ServeHTTP)
	}()

	func() {
		putHandler, err := getHTTPHandler(
			http.MethodPatch,
			"/changes/{primaryKey}",
			http.StatusOK,
			func(
				ctx context.Context,
				pathParams ChangeOnePathParams,
				queryParams ChangeLoadQueryParams,
				req Change,
				rawReq any,
			) (server.Response[Change], error) {
				item, ok := rawReq.(map[string]any)
				if !ok {
					return server.Response[Change]{}, fmt.Errorf("failed to cast %#+v to map[string]any", item)
				}

				arguments, err := server.GetLoadArguments(ctx, queryParams.Depth)
				if err != nil {
					return server.Response[Change]{}, err
				}

				object := &req
				object.ID = pathParams.PrimaryKey

				objects, count, totalCount, _, _, err := handlePutChange(arguments, db, waitForChange, object)
				if err != nil {
					return server.Response[Change]{}, err
				}

				limit := int64(0)

				offset := int64(0)

				return server.Response[Change]{
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
			Change{},
			ChangeIntrospectedTable,
		)
		if err != nil {
			panic(err)
		}
		r.Put(putHandler.FullPath, putHandler.ServeHTTP)
	}()

	func() {
		patchHandler, err := getHTTPHandler(
			http.MethodPatch,
			"/changes/{primaryKey}",
			http.StatusOK,
			func(
				ctx context.Context,
				pathParams ChangeOnePathParams,
				queryParams ChangeLoadQueryParams,
				req Change,
				rawReq any,
			) (server.Response[Change], error) {
				item, ok := rawReq.(map[string]any)
				if !ok {
					return server.Response[Change]{}, fmt.Errorf("failed to cast %#+v to map[string]any", item)
				}

				forceSetValuesForFields := make([]string, 0)
				for _, possibleField := range maps.Keys(item) {
					if !slices.Contains(ChangeTableColumns, possibleField) {
						continue
					}

					forceSetValuesForFields = append(forceSetValuesForFields, possibleField)
				}

				arguments, err := server.GetLoadArguments(ctx, queryParams.Depth)
				if err != nil {
					return server.Response[Change]{}, err
				}

				object := &req
				object.ID = pathParams.PrimaryKey

				objects, count, totalCount, _, _, err := handlePatchChange(arguments, db, waitForChange, object, forceSetValuesForFields)
				if err != nil {
					return server.Response[Change]{}, err
				}

				limit := int64(0)

				offset := int64(0)

				return server.Response[Change]{
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
			Change{},
			ChangeIntrospectedTable,
		)
		if err != nil {
			panic(err)
		}
		r.Patch(patchHandler.FullPath, patchHandler.ServeHTTP)
	}()

	func() {
		deleteHandler, err := getHTTPHandler(
			http.MethodDelete,
			"/changes/{primaryKey}",
			http.StatusNoContent,
			func(
				ctx context.Context,
				pathParams ChangeOnePathParams,
				queryParams ChangeLoadQueryParams,
				req server.EmptyRequest,
				rawReq any,
			) (server.EmptyResponse, error) {
				arguments, err := server.GetLoadArguments(ctx, queryParams.Depth)
				if err != nil {
					return server.EmptyResponse{}, err
				}

				object := &Change{}
				object.ID = pathParams.PrimaryKey

				err = handleDeleteChange(arguments, db, waitForChange, object)
				if err != nil {
					return server.EmptyResponse{}, err
				}

				return server.EmptyResponse{}, nil
			},
			Change{},
			ChangeIntrospectedTable,
		)
		if err != nil {
			panic(err)
		}
		r.Delete(deleteHandler.FullPath, deleteHandler.ServeHTTP)
	}()
}

func NewChangeFromItem(item map[string]any) (any, error) {
	object := &Change{}

	err := object.FromItem(item)
	if err != nil {
		return nil, err
	}

	return object, nil
}

func init() {
	register(
		ChangeTable,
		Change{},
		NewChangeFromItem,
		"/changes",
		MutateRouterForChange,
	)
}
