package cameranator

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

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

type Camera struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	StreamURL string `json:"stream_url"`
}

var CameraTable = "camera"

var (
	CameraTableIDColumn        = "id"
	CameraTableNameColumn      = "name"
	CameraTableStreamURLColumn = "stream_url"
)

var (
	CameraTableIDColumnWithTypeCast        = fmt.Sprintf(`"id" AS id`)
	CameraTableNameColumnWithTypeCast      = fmt.Sprintf(`"name" AS name`)
	CameraTableStreamURLColumnWithTypeCast = fmt.Sprintf(`"stream_url" AS stream_url`)
)

var CameraTableColumns = []string{
	CameraTableIDColumn,
	CameraTableNameColumn,
	CameraTableStreamURLColumn,
}

var CameraTableColumnsWithTypeCasts = []string{
	CameraTableIDColumnWithTypeCast,
	CameraTableNameColumnWithTypeCast,
	CameraTableStreamURLColumnWithTypeCast,
}

var CameraTableColumnLookup = map[string]*introspect.Column{
	CameraTableIDColumn:        new(introspect.Column),
	CameraTableNameColumn:      new(introspect.Column),
	CameraTableStreamURLColumn: new(introspect.Column),
}

var (
	CameraTablePrimaryKeyColumn = CameraTableIDColumn
)

var (
	_ = time.Time{}
	_ = uuid.UUID{}
	_ = pq.StringArray{}
	_ = hstore.Hstore{}
	_ = geojson.Point{}
	_ = pgtype.Point{}
	_ = _pgtype.Point{}
)

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

	wrapError := func(k string, err error) error {
		return fmt.Errorf("%#+v: %v; item: %#+v", k, err, item)
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

			temp1, err := types.ParseInt(v)
			if err != nil {
				return wrapError(k, err)
			}

			temp2, ok := temp1.(int64)
			if !ok {
				if temp1 != nil {
					return wrapError(k, fmt.Errorf("failed to cast %#+v to int64", temp1))
				}
			}

			m.ID = temp2

		case "name":
			if v == nil {
				continue
			}

			temp1, err := types.ParseString(v)
			if err != nil {
				return wrapError(k, err)
			}

			temp2, ok := temp1.(string)
			if !ok {
				if temp1 != nil {
					return wrapError(k, fmt.Errorf("failed to cast %#+v to string", temp1))
				}
			}

			m.Name = temp2

		case "stream_url":
			if v == nil {
				continue
			}

			temp1, err := types.ParseString(v)
			if err != nil {
				return wrapError(k, err)
			}

			temp2, ok := temp1.(string)
			if !ok {
				if temp1 != nil {
					return wrapError(k, fmt.Errorf("failed to cast %#+v to string", temp1))
				}
			}

			m.StreamURL = temp2

		}
	}

	return nil
}

func (m *Camera) Reload(
	ctx context.Context,
	tx *sqlx.Tx,
) error {
	t, err := SelectCamera(
		ctx,
		tx,
		fmt.Sprintf("%v = $1", m.GetPrimaryKeyColumn()),
		m.GetPrimaryKeyValue(),
	)
	if err != nil {
		return err
	}

	m.ID = t.ID
	m.Name = t.Name
	m.StreamURL = t.StreamURL

	return nil
}

func (m *Camera) Insert(
	ctx context.Context,
	tx *sqlx.Tx,
	setPrimaryKey bool,
	setZeroValues bool,
) error {
	columns := make([]string, 0)
	values := make([]any, 0)

	if setPrimaryKey && (setZeroValues || !types.IsZeroInt(m.ID)) {
		columns = append(columns, CameraTableIDColumn)

		v, err := types.FormatInt(m.ID)
		if err != nil {
			return fmt.Errorf("failed to handle m.ID: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroString(m.Name) {
		columns = append(columns, CameraTableNameColumn)

		v, err := types.FormatString(m.Name)
		if err != nil {
			return fmt.Errorf("failed to handle m.Name: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroString(m.StreamURL) {
		columns = append(columns, CameraTableStreamURLColumn)

		v, err := types.FormatString(m.StreamURL)
		if err != nil {
			return fmt.Errorf("failed to handle m.StreamURL: %v", err)
		}

		values = append(values, v)
	}

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
		return fmt.Errorf("failed to insert %#+v: %v", m, err)
	}
	v := item[CameraTableIDColumn]

	if v == nil {
		return fmt.Errorf("failed to find %v in %#+v", CameraTableIDColumn, item)
	}

	wrapError := func(err error) error {
		return fmt.Errorf(
			"failed to treat %v: %#+v as int64: %v",
			CameraTableIDColumn,
			item[CameraTableIDColumn],
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

	m.ID = temp2

	err = m.Reload(ctx, tx)
	if err != nil {
		return fmt.Errorf("failed to reload after insert")
	}

	return nil
}

func (m *Camera) Update(
	ctx context.Context,
	tx *sqlx.Tx,
	setZeroValues bool,
) error {
	columns := make([]string, 0)
	values := make([]any, 0)

	if setZeroValues || !types.IsZeroString(m.Name) {
		columns = append(columns, CameraTableNameColumn)

		v, err := types.FormatString(m.Name)
		if err != nil {
			return fmt.Errorf("failed to handle m.Name: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroString(m.StreamURL) {
		columns = append(columns, CameraTableStreamURLColumn)

		v, err := types.FormatString(m.StreamURL)
		if err != nil {
			return fmt.Errorf("failed to handle m.StreamURL: %v", err)
		}

		values = append(values, v)
	}

	v, err := types.FormatInt(m.ID)
	if err != nil {
		return fmt.Errorf("failed to handle m.ID: %v", err)
	}

	values = append(values, v)

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
		return fmt.Errorf("failed to update %#+v: %v", m, err)
	}

	err = m.Reload(ctx, tx)
	if err != nil {
		return fmt.Errorf("failed to reload after update")
	}

	return nil
}

func (m *Camera) Delete(
	ctx context.Context,
	tx *sqlx.Tx,
) error {
	values := make([]any, 0)
	v, err := types.FormatInt(m.ID)
	if err != nil {
		return fmt.Errorf("failed to handle m.ID: %v", err)
	}

	values = append(values, v)

	err = query.Delete(
		ctx,
		tx,
		CameraTable,
		fmt.Sprintf("%v = $$??", CameraTableIDColumn),
		values...,
	)
	if err != nil {
		return fmt.Errorf("failed to delete %#+v: %v", m, err)
	}

	return nil
}

func SelectCameras(
	ctx context.Context,
	tx *sqlx.Tx,
	where string,
	limit *int,
	offset *int,
	values ...any,
) ([]*Camera, error) {
	items, err := query.Select(
		ctx,
		tx,
		CameraTableColumnsWithTypeCasts,
		CameraTable,
		where,
		limit,
		offset,
		values...,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to call SelectCameras; err: %v", err)
	}

	objects := make([]*Camera, 0)

	for _, item := range items {
		object := &Camera{}

		err = object.FromItem(item)
		if err != nil {
			return nil, fmt.Errorf("failed to call Camera.FromItem; err: %v", err)
		}

		objects = append(objects, object)
	}

	return objects, nil
}

func SelectCamera(
	ctx context.Context,
	tx *sqlx.Tx,
	where string,
	values ...any,
) (*Camera, error) {
	objects, err := SelectCameras(
		ctx,
		tx,
		where,
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
		return nil, fmt.Errorf("attempt to call SelectCamera returned no rows")
	}

	object := objects[0]

	return object, nil
}

func handleGetCameras(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
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
			column := CameraTableColumnLookup[parts[0]]
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

	objects, err := SelectCameras(ctx, tx, where, &limit, &offset, values...)
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

func handleGetCamera(w http.ResponseWriter, r *http.Request, db *sqlx.DB, primaryKey string) {
	ctx := r.Context()

	where := fmt.Sprintf("%s = $$??", CameraTablePrimaryKeyColumn)

	values := []any{primaryKey}

	tx, err := db.BeginTxx(r.Context(), nil)
	if err != nil {
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	object, err := SelectCamera(ctx, tx, where, values...)
	if err != nil {
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	err = tx.Commit()
	if err != nil {
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	helpers.HandleObjectsResponse(w, http.StatusOK, []*Camera{object})
}

func handlePostCameras(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
}

func handlePutCamera(w http.ResponseWriter, r *http.Request, db *sqlx.DB, primaryKey string) {
}

func handlePatchCamera(w http.ResponseWriter, r *http.Request, db *sqlx.DB, primaryKey string) {
}

func GetCameraRouter(db *sqlx.DB, middlewares ...func(http.Handler) http.Handler) chi.Router {
	r := chi.NewRouter()

	for _, m := range middlewares {
		r.Use(m)
	}

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		handleGetCameras(w, r, db)
	})

	r.Get("/{primaryKey}", func(w http.ResponseWriter, r *http.Request) {
		handleGetCamera(w, r, db, chi.URLParam(r, "primaryKey"))
	})

	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		handlePostCameras(w, r, db)
	})

	r.Put("/{primaryKey}", func(w http.ResponseWriter, r *http.Request) {
		handlePutCamera(w, r, db, chi.URLParam(r, "primaryKey"))
	})

	r.Patch("/{primaryKey}", func(w http.ResponseWriter, r *http.Request) {
		handlePatchCamera(w, r, db, chi.URLParam(r, "primaryKey"))
	})

	return r
}

func GetCameraHandlerFunc(db *sqlx.DB, middlewares ...func(http.Handler) http.Handler) http.HandlerFunc {
	r := chi.NewRouter()

	r.Mount("/cameras", GetCameraRouter(db, middlewares...))

	return r.ServeHTTP
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
		NewCameraFromItem,
		"/cameras",
		GetCameraRouter,
	)
}
