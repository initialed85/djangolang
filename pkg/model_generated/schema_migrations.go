package model_generated

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/netip"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/cridenour/go-postgis"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/initialed85/djangolang/pkg/helpers"
	"github.com/initialed85/djangolang/pkg/introspect"
	"github.com/initialed85/djangolang/pkg/query"
	"github.com/initialed85/djangolang/pkg/types"
	_pgtype "github.com/jackc/pgtype"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/lib/pq/hstore"
	"github.com/paulmach/orb/geojson"
)

type SchemaMigration struct {
	Version int64 `json:"version"`
	Dirty   bool  `json:"dirty"`
}

var SchemaMigrationTable = "schema_migrations"

var (
	SchemaMigrationTableVersionColumn = "version"
	SchemaMigrationTableDirtyColumn   = "dirty"
)

var (
	SchemaMigrationTableVersionColumnWithTypeCast = fmt.Sprintf(`"version" AS version`)
	SchemaMigrationTableDirtyColumnWithTypeCast   = fmt.Sprintf(`"dirty" AS dirty`)
)

var SchemaMigrationTableColumns = []string{
	SchemaMigrationTableVersionColumn,
	SchemaMigrationTableDirtyColumn,
}

var SchemaMigrationTableColumnsWithTypeCasts = []string{
	SchemaMigrationTableVersionColumnWithTypeCast,
	SchemaMigrationTableDirtyColumnWithTypeCast,
}

var SchemaMigrationTableColumnLookup = map[string]*introspect.Column{
	SchemaMigrationTableVersionColumn: new(introspect.Column),
	SchemaMigrationTableDirtyColumn:   new(introspect.Column),
}

var (
	SchemaMigrationTablePrimaryKeyColumn = SchemaMigrationTableVersionColumn
)

var (
	_ = time.Time{}
	_ = uuid.UUID{}
	_ = pq.StringArray{}
	_ = hstore.Hstore{}
	_ = geojson.Point{}
	_ = pgtype.Point{}
	_ = _pgtype.Point{}
	_ = postgis.PointZ{}
	_ = netip.Prefix{}
)

func (m *SchemaMigration) GetPrimaryKeyColumn() string {
	return SchemaMigrationTablePrimaryKeyColumn
}

func (m *SchemaMigration) GetPrimaryKeyValue() any {
	return m.Version
}

func (m *SchemaMigration) FromItem(item map[string]any) error {
	if item == nil {
		return fmt.Errorf(
			"item unexpectedly nil during SchemaMigrationFromItem",
		)
	}

	if len(item) == 0 {
		return fmt.Errorf(
			"item unexpectedly empty during SchemaMigrationFromItem",
		)
	}

	wrapError := func(k string, v any, err error) error {
		return fmt.Errorf("%v: %#+v; error: %v", k, v, err)
	}

	for k, v := range item {
		_, ok := SchemaMigrationTableColumnLookup[k]
		if !ok {
			return fmt.Errorf(
				"item contained unexpected key %#+v during SchemaMigrationFromItem; item: %#+v",
				k, item,
			)
		}

		switch k {
		case "version":
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
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to int64", temp1))
				}
			}

			m.Version = temp2

		case "dirty":
			if v == nil {
				continue
			}

			temp1, err := types.ParseBool(v)
			if err != nil {
				return wrapError(k, v, err)
			}

			temp2, ok := temp1.(bool)
			if !ok {
				if temp1 != nil {
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to bool", temp1))
				}
			}

			m.Dirty = temp2

		}
	}

	return nil
}

func (m *SchemaMigration) Reload(
	ctx context.Context,
	tx *sqlx.Tx,
	includeDeleteds ...bool,
) error {
	extraWhere := ""
	if len(includeDeleteds) > 0 {
		if slices.Contains(SchemaMigrationTableColumns, "deleted_at") {
			extraWhere = "\n    AND (deleted_at IS null OR deleted_at IS NOT null)"
		}
	}

	t, err := SelectSchemaMigration(
		ctx,
		tx,
		fmt.Sprintf("%v = $1%v", m.GetPrimaryKeyColumn(), extraWhere),
		m.GetPrimaryKeyValue(),
	)
	if err != nil {
		return err
	}

	m.Version = t.Version
	m.Dirty = t.Dirty

	return nil
}

func (m *SchemaMigration) Insert(
	ctx context.Context,
	tx *sqlx.Tx,
	setPrimaryKey bool,
	setZeroValues bool,
) error {
	columns := make([]string, 0)
	values := make([]any, 0)

	if setPrimaryKey && (setZeroValues || !types.IsZeroInt(m.Version)) {
		columns = append(columns, SchemaMigrationTableVersionColumn)

		v, err := types.FormatInt(m.Version)
		if err != nil {
			return fmt.Errorf("failed to handle m.Version: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroBool(m.Dirty) {
		columns = append(columns, SchemaMigrationTableDirtyColumn)

		v, err := types.FormatBool(m.Dirty)
		if err != nil {
			return fmt.Errorf("failed to handle m.Dirty: %v", err)
		}

		values = append(values, v)
	}

	item, err := query.Insert(
		ctx,
		tx,
		SchemaMigrationTable,
		columns,
		nil,
		false,
		false,
		SchemaMigrationTableColumns,
		values...,
	)
	if err != nil {
		return fmt.Errorf("failed to insert %#+v: %v", m, err)
	}
	v := item[SchemaMigrationTableVersionColumn]

	if v == nil {
		return fmt.Errorf("failed to find %v in %#+v", SchemaMigrationTableVersionColumn, item)
	}

	wrapError := func(err error) error {
		return fmt.Errorf(
			"failed to treat %v: %#+v as int64: %v",
			SchemaMigrationTableVersionColumn,
			item[SchemaMigrationTableVersionColumn],
			err,
		)
	}

	temp1, err := types.ParseUUID(v)
	if err != nil {
		return wrapError(err)
	}

	temp2, ok := temp1.(int64)
	if !ok {
		return wrapError(fmt.Errorf("failed to cast to int64"))
	}

	m.Version = temp2

	err = m.Reload(ctx, tx)
	if err != nil {
		return fmt.Errorf("failed to reload after insert")
	}

	return nil
}

func (m *SchemaMigration) Update(
	ctx context.Context,
	tx *sqlx.Tx,
	setZeroValues bool,
) error {
	columns := make([]string, 0)
	values := make([]any, 0)

	if setZeroValues || !types.IsZeroBool(m.Dirty) {
		columns = append(columns, SchemaMigrationTableDirtyColumn)

		v, err := types.FormatBool(m.Dirty)
		if err != nil {
			return fmt.Errorf("failed to handle m.Dirty: %v", err)
		}

		values = append(values, v)
	}

	v, err := types.FormatInt(m.Version)
	if err != nil {
		return fmt.Errorf("failed to handle m.Version: %v", err)
	}

	values = append(values, v)

	_, err = query.Update(
		ctx,
		tx,
		SchemaMigrationTable,
		columns,
		fmt.Sprintf("%v = $$??", SchemaMigrationTableVersionColumn),
		SchemaMigrationTableColumns,
		values...,
	)
	if err != nil {
		return fmt.Errorf("failed to update %#+v: %v", m, err)
	}

	err = m.Reload(ctx, tx)
	if err != nil {
		return fmt.Errorf("failed to reload after update")
	}

	return nil
}

func (m *SchemaMigration) Delete(
	ctx context.Context,
	tx *sqlx.Tx,
) error {
	values := make([]any, 0)
	v, err := types.FormatInt(m.Version)
	if err != nil {
		return fmt.Errorf("failed to handle m.Version: %v", err)
	}

	values = append(values, v)

	err = query.Delete(
		ctx,
		tx,
		SchemaMigrationTable,
		fmt.Sprintf("%v = $$??", SchemaMigrationTableVersionColumn),
		values...,
	)
	if err != nil {
		return fmt.Errorf("failed to delete %#+v: %v", m, err)
	}

	return nil
}

func SelectSchemaMigrations(
	ctx context.Context,
	tx *sqlx.Tx,
	where string,
	limit *int,
	offset *int,
	values ...any,
) ([]*SchemaMigration, error) {
	if slices.Contains(SchemaMigrationTableColumns, "deleted_at") {
		if !strings.Contains(where, "deleted_at") {
			if where != "" {
				where += "\n    AND "
			}

			where += "deleted_at IS null"
		}
	}

	items, err := query.Select(
		ctx,
		tx,
		SchemaMigrationTableColumnsWithTypeCasts,
		SchemaMigrationTable,
		where,
		limit,
		offset,
		values...,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to call SelectSchemaMigrations; err: %v", err)
	}

	objects := make([]*SchemaMigration, 0)

	for _, item := range items {
		object := &SchemaMigration{}

		err = object.FromItem(item)
		if err != nil {
			return nil, fmt.Errorf("failed to call SchemaMigration.FromItem; err: %v", err)
		}

		objects = append(objects, object)
	}

	return objects, nil
}

func SelectSchemaMigration(
	ctx context.Context,
	tx *sqlx.Tx,
	where string,
	values ...any,
) (*SchemaMigration, error) {
	objects, err := SelectSchemaMigrations(
		ctx,
		tx,
		where,
		helpers.Ptr(2),
		helpers.Ptr(0),
		values...,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to call SelectSchemaMigration; err: %v", err)
	}

	if len(objects) > 1 {
		return nil, fmt.Errorf("attempt to call SelectSchemaMigration returned more than 1 row")
	}

	if len(objects) < 1 {
		return nil, fmt.Errorf("attempt to call SelectSchemaMigration returned no rows")
	}

	object := objects[0]

	return object, nil
}

func handleGetSchemaMigrations(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
	ctx := r.Context()

	unrecognizedParams := make([]string, 0)
	hadUnrecognizedParams := false

	unparseableParams := make([]string, 0)
	hadUnparseableParams := false

	values := make([]any, 0)

	wheres := make([]string, 0)
	for rawKey, rawValues := range r.URL.Query() {
		if rawKey == "limit" || rawKey == "offset" {
			continue
		}

		parts := strings.Split(rawKey, "__")
		isUnrecognized := len(parts) != 2

		comparison := ""
		isSliceComparison := false
		isNullComparison := false
		IsLikeComparison := false

		if !isUnrecognized {
			column := SchemaMigrationTableColumnLookup[parts[0]]
			if column == nil {
				isUnrecognized = true
			} else {
				switch parts[1] {
				case "eq":
					comparison = "="
				case "ne":
					comparison = "!="
				case "gt":
					comparison = ">"
				case "gte":
					comparison = ">="
				case "lt":
					comparison = "<"
				case "lte":
					comparison = "<="
				case "in":
					comparison = "IN"
					isSliceComparison = true
				case "nin", "notin":
					comparison = "NOT IN"
					isSliceComparison = true
				case "isnull":
					comparison = "IS NULL"
					isNullComparison = true
				case "nisnull", "isnotnull":
					comparison = "IS NOT NULL"
					isNullComparison = true
				case "l", "like":
					comparison = "LIKE"
					IsLikeComparison = true
				case "nl", "nlike", "notlike":
					comparison = "NOT LIKE"
					IsLikeComparison = true
				case "il", "ilike":
					comparison = "ILIKE"
					IsLikeComparison = true
				case "nil", "nilike", "notilike":
					comparison = "NOT ILIKE"
					IsLikeComparison = true
				default:
					isUnrecognized = true
				}
			}
		}

		if isNullComparison {
			wheres = append(wheres, fmt.Sprintf("%s %s", parts[0], comparison))
			continue
		}

		for _, rawValue := range rawValues {
			if isUnrecognized {
				unrecognizedParams = append(unrecognizedParams, fmt.Sprintf("%s=%s", rawKey, rawValue))
				hadUnrecognizedParams = true
				continue
			}

			if hadUnrecognizedParams {
				continue
			}

			attempts := make([]string, 0)

			if !IsLikeComparison {
				attempts = append(attempts, rawValue)
			}

			if isSliceComparison {
				attempts = append(attempts, fmt.Sprintf("[%s]", rawValue))

				vs := make([]string, 0)
				for _, v := range strings.Split(rawValue, ",") {
					vs = append(vs, fmt.Sprintf("\"%s\"", v))
				}

				attempts = append(attempts, fmt.Sprintf("[%s]", strings.Join(vs, ",")))
			}

			if IsLikeComparison {
				attempts = append(attempts, fmt.Sprintf("\"%%%s%%\"", rawValue))
			} else {
				attempts = append(attempts, fmt.Sprintf("\"%s\"", rawValue))
			}

			var err error

			for _, attempt := range attempts {
				var value any
				err = json.Unmarshal([]byte(attempt), &value)
				if err == nil {
					if isSliceComparison {
						sliceValues, ok := value.([]any)
						if !ok {
							err = fmt.Errorf("failed to cast %#+v to []string", value)
							break
						}

						values = append(values, sliceValues...)

						sliceWheres := make([]string, 0)
						for range values {
							sliceWheres = append(sliceWheres, "$$??")
						}

						wheres = append(wheres, fmt.Sprintf("%s %s (%s)", parts[0], comparison, strings.Join(sliceWheres, ", ")))
					} else {
						values = append(values, value)
						wheres = append(wheres, fmt.Sprintf("%s %s $$??", parts[0], comparison))
					}

					break
				}
			}

			if err != nil {
				unparseableParams = append(unparseableParams, fmt.Sprintf("%s=%s", rawKey, rawValue))
				hadUnparseableParams = true
				continue
			}
		}
	}

	where := strings.Join(wheres, "\n    AND ")

	if hadUnrecognizedParams {
		helpers.HandleErrorResponse(
			w,
			http.StatusInternalServerError,
			fmt.Errorf("unrecognized params %s", strings.Join(unrecognizedParams, ", ")),
		)
		return
	}

	if hadUnparseableParams {
		helpers.HandleErrorResponse(
			w,
			http.StatusInternalServerError,
			fmt.Errorf("unparseable params %s", strings.Join(unparseableParams, ", ")),
		)
		return
	}

	limit := 2000
	rawLimit := r.URL.Query().Get("limit")
	if rawLimit != "" {
		possibleLimit, err := strconv.ParseInt(rawLimit, 10, 64)
		if err == nil {
			helpers.HandleErrorResponse(
				w,
				http.StatusInternalServerError,
				fmt.Errorf("failed to parse param limit=%s as int: %v", rawLimit, err),
			)
			return
		}

		limit = int(possibleLimit)
	}

	offset := 0
	rawOffset := r.URL.Query().Get("offset")
	if rawOffset != "" {
		possibleOffset, err := strconv.ParseInt(rawOffset, 10, 64)
		if err == nil {
			helpers.HandleErrorResponse(
				w,
				http.StatusInternalServerError,
				fmt.Errorf("failed to parse param offset=%s as int: %v", rawOffset, err),
			)
			return
		}

		offset = int(possibleOffset)
	}

	tx, err := db.BeginTxx(r.Context(), nil)
	if err != nil {
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	objects, err := SelectSchemaMigrations(ctx, tx, where, &limit, &offset, values...)
	if err != nil {
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	err = tx.Commit()
	if err != nil {
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	helpers.HandleObjectsResponse(w, http.StatusOK, objects)
}

func handleGetSchemaMigration(w http.ResponseWriter, r *http.Request, db *sqlx.DB, primaryKey string) {
	ctx := r.Context()

	where := fmt.Sprintf("%s = $$??", SchemaMigrationTablePrimaryKeyColumn)

	values := []any{primaryKey}

	tx, err := db.BeginTxx(r.Context(), nil)
	if err != nil {
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	object, err := SelectSchemaMigration(ctx, tx, where, values...)
	if err != nil {
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	err = tx.Commit()
	if err != nil {
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	helpers.HandleObjectsResponse(w, http.StatusOK, []*SchemaMigration{object})
}

func handlePostSchemaMigrations(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		err = fmt.Errorf("failed to read body of HTTP request: %v", err)
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	var allItems []map[string]any
	err = json.Unmarshal(b, &allItems)
	if err != nil {
		err = fmt.Errorf("failed to unmarshal %#+v as JSON list of objects: %v", string(b), err)
		helpers.HandleErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	objects := make([]*SchemaMigration, 0)
	for _, item := range allItems {
		object := &SchemaMigration{}
		err = object.FromItem(item)
		if err != nil {
			err = fmt.Errorf("failed to interpret %#+v as SchemaMigration in item form: %v", item, err)
			helpers.HandleErrorResponse(w, http.StatusBadRequest, err)
			return
		}

		objects = append(objects, object)
	}

	tx, err := db.BeginTxx(r.Context(), nil)
	if err != nil {
		err = fmt.Errorf("failed to begin DB transaction: %v", err)
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	defer func() {
		_ = tx.Rollback()
	}()

	for i, object := range objects {
		err = object.Insert(r.Context(), tx, false, false)
		if err != nil {
			err = fmt.Errorf("failed to insert %#+v: %v", object, err)
			helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
			return
		}

		objects[i] = object
	}

	err = tx.Commit()
	if err != nil {
		err = fmt.Errorf("failed to commit DB transaction: %v", err)
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	helpers.HandleObjectsResponse(w, http.StatusCreated, objects)
}

func handlePutSchemaMigration(w http.ResponseWriter, r *http.Request, db *sqlx.DB, primaryKey string) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		err = fmt.Errorf("failed to read body of HTTP request: %v", err)
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	var item map[string]any
	err = json.Unmarshal(b, &item)
	if err != nil {
		err = fmt.Errorf("failed to unmarshal %#+v as JSON object: %v", string(b), err)
		helpers.HandleErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	item[SchemaMigrationTablePrimaryKeyColumn] = primaryKey

	object := &SchemaMigration{}
	err = object.FromItem(item)
	if err != nil {
		err = fmt.Errorf("failed to interpret %#+v as SchemaMigration in item form: %v", item, err)
		helpers.HandleErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	tx, err := db.BeginTxx(r.Context(), nil)
	if err != nil {
		err = fmt.Errorf("failed to begin DB transaction: %v", err)
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	defer func() {
		_ = tx.Rollback()
	}()

	err = object.Update(r.Context(), tx, true)
	if err != nil {
		err = fmt.Errorf("failed to update %#+v: %v", object, err)
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	err = tx.Commit()
	if err != nil {
		err = fmt.Errorf("failed to commit DB transaction: %v", err)
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	helpers.HandleObjectsResponse(w, http.StatusOK, []*SchemaMigration{object})
}

func handlePatchSchemaMigration(w http.ResponseWriter, r *http.Request, db *sqlx.DB, primaryKey string) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		err = fmt.Errorf("failed to read body of HTTP request: %v", err)
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	var item map[string]any
	err = json.Unmarshal(b, &item)
	if err != nil {
		err = fmt.Errorf("failed to unmarshal %#+v as JSON object: %v", string(b), err)
		helpers.HandleErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	item[SchemaMigrationTablePrimaryKeyColumn] = primaryKey

	object := &SchemaMigration{}
	err = object.FromItem(item)
	if err != nil {
		err = fmt.Errorf("failed to interpret %#+v as SchemaMigration in item form: %v", item, err)
		helpers.HandleErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	tx, err := db.BeginTxx(r.Context(), nil)
	if err != nil {
		err = fmt.Errorf("failed to begin DB transaction: %v", err)
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	defer func() {
		_ = tx.Rollback()
	}()

	err = object.Update(r.Context(), tx, false)
	if err != nil {
		err = fmt.Errorf("failed to update %#+v: %v", object, err)
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	err = tx.Commit()
	if err != nil {
		err = fmt.Errorf("failed to commit DB transaction: %v", err)
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	helpers.HandleObjectsResponse(w, http.StatusOK, []*SchemaMigration{object})
}

func handleDeleteSchemaMigration(w http.ResponseWriter, r *http.Request, db *sqlx.DB, primaryKey string) {
	var item = make(map[string]any)

	item[SchemaMigrationTablePrimaryKeyColumn] = primaryKey

	object := &SchemaMigration{}
	err := object.FromItem(item)
	if err != nil {
		err = fmt.Errorf("failed to interpret %#+v as SchemaMigration in item form: %v", item, err)
		helpers.HandleErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	tx, err := db.BeginTxx(r.Context(), nil)
	if err != nil {
		err = fmt.Errorf("failed to begin DB transaction: %v", err)
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	defer func() {
		_ = tx.Rollback()
	}()

	err = object.Delete(r.Context(), tx)
	if err != nil {
		err = fmt.Errorf("failed to delete %#+v: %v", object, err)
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	err = tx.Commit()
	if err != nil {
		err = fmt.Errorf("failed to commit DB transaction: %v", err)
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	helpers.HandleObjectsResponse(w, http.StatusNoContent, nil)
}

func GetSchemaMigrationRouter(db *sqlx.DB, middlewares ...func(http.Handler) http.Handler) chi.Router {
	r := chi.NewRouter()

	for _, m := range middlewares {
		r.Use(m)
	}

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		handleGetSchemaMigrations(w, r, db)
	})

	r.Get("/{primaryKey}", func(w http.ResponseWriter, r *http.Request) {
		handleGetSchemaMigration(w, r, db, chi.URLParam(r, "primaryKey"))
	})

	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		handlePostSchemaMigrations(w, r, db)
	})

	r.Put("/{primaryKey}", func(w http.ResponseWriter, r *http.Request) {
		handlePutSchemaMigration(w, r, db, chi.URLParam(r, "primaryKey"))
	})

	r.Patch("/{primaryKey}", func(w http.ResponseWriter, r *http.Request) {
		handlePatchSchemaMigration(w, r, db, chi.URLParam(r, "primaryKey"))
	})

	r.Delete("/{primaryKey}", func(w http.ResponseWriter, r *http.Request) {
		handleDeleteSchemaMigration(w, r, db, chi.URLParam(r, "primaryKey"))
	})

	return r
}

func GetSchemaMigrationHandlerFunc(db *sqlx.DB, middlewares ...func(http.Handler) http.Handler) http.HandlerFunc {
	r := chi.NewRouter()

	r.Mount("/schema-migrations", GetSchemaMigrationRouter(db, middlewares...))

	return r.ServeHTTP
}

func NewSchemaMigrationFromItem(item map[string]any) (any, error) {
	object := &SchemaMigration{}

	err := object.FromItem(item)
	if err != nil {
		return nil, err
	}

	return object, nil
}

func init() {
	register(
		SchemaMigrationTable,
		NewSchemaMigrationFromItem,
		"/schema-migrations",
		GetSchemaMigrationRouter,
	)
}
