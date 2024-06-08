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

type Image struct {
	ID             int64     `json:"id"`
	Timestamp      time.Time `json:"timestamp"`
	Size           float64   `json:"size"`
	FilePath       string    `json:"file_path"`
	CameraID       int64     `json:"camera_id"`
	CameraIDObject *Camera   `json:"camera_id_object"`
	EventID        *int64    `json:"event_id"`
	EventIDObject  *Event    `json:"event_id_object"`
}

var ImageTable = "image"

var (
	ImageTableIDColumn        = "id"
	ImageTableTimestampColumn = "timestamp"
	ImageTableSizeColumn      = "size"
	ImageTableFilePathColumn  = "file_path"
	ImageTableCameraIDColumn  = "camera_id"
	ImageTableEventIDColumn   = "event_id"
)

var (
	ImageTableIDColumnWithTypeCast        = fmt.Sprintf(`"id" AS id`)
	ImageTableTimestampColumnWithTypeCast = fmt.Sprintf(`"timestamp" AS timestamp`)
	ImageTableSizeColumnWithTypeCast      = fmt.Sprintf(`"size" AS size`)
	ImageTableFilePathColumnWithTypeCast  = fmt.Sprintf(`"file_path" AS file_path`)
	ImageTableCameraIDColumnWithTypeCast  = fmt.Sprintf(`"camera_id" AS camera_id`)
	ImageTableEventIDColumnWithTypeCast   = fmt.Sprintf(`"event_id" AS event_id`)
)

var ImageTableColumns = []string{
	ImageTableIDColumn,
	ImageTableTimestampColumn,
	ImageTableSizeColumn,
	ImageTableFilePathColumn,
	ImageTableCameraIDColumn,
	ImageTableEventIDColumn,
}

var ImageTableColumnsWithTypeCasts = []string{
	ImageTableIDColumnWithTypeCast,
	ImageTableTimestampColumnWithTypeCast,
	ImageTableSizeColumnWithTypeCast,
	ImageTableFilePathColumnWithTypeCast,
	ImageTableCameraIDColumnWithTypeCast,
	ImageTableEventIDColumnWithTypeCast,
}

var ImageTableColumnLookup = map[string]*introspect.Column{
	ImageTableIDColumn:        new(introspect.Column),
	ImageTableTimestampColumn: new(introspect.Column),
	ImageTableSizeColumn:      new(introspect.Column),
	ImageTableFilePathColumn:  new(introspect.Column),
	ImageTableCameraIDColumn:  new(introspect.Column),
	ImageTableEventIDColumn:   new(introspect.Column),
}

var (
	ImageTablePrimaryKeyColumn = ImageTableIDColumn
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

func (m *Image) GetPrimaryKeyColumn() string {
	return ImageTablePrimaryKeyColumn
}

func (m *Image) GetPrimaryKeyValue() any {
	return m.ID
}

func (m *Image) FromItem(item map[string]any) error {
	if item == nil {
		return fmt.Errorf(
			"item unexpectedly nil during ImageFromItem",
		)
	}

	if len(item) == 0 {
		return fmt.Errorf(
			"item unexpectedly empty during ImageFromItem",
		)
	}

	wrapError := func(k string, err error) error {
		return fmt.Errorf("%#+v: %v; item: %#+v", k, err, item)
	}

	for k, v := range item {
		_, ok := ImageTableColumnLookup[k]
		if !ok {
			return fmt.Errorf(
				"item contained unexpected key %#+v during ImageFromItem; item: %#+v",
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

		case "timestamp":
			if v == nil {
				continue
			}

			temp1, err := types.ParseTime(v)
			if err != nil {
				return wrapError(k, err)
			}

			temp2, ok := temp1.(time.Time)
			if !ok {
				if temp1 != nil {
					return wrapError(k, fmt.Errorf("failed to cast %#+v to time.Time", temp1))
				}
			}

			m.Timestamp = temp2

		case "size":
			if v == nil {
				continue
			}

			temp1, err := types.ParseFloat(v)
			if err != nil {
				return wrapError(k, err)
			}

			temp2, ok := temp1.(float64)
			if !ok {
				if temp1 != nil {
					return wrapError(k, fmt.Errorf("failed to cast %#+v to float64", temp1))
				}
			}

			m.Size = temp2

		case "file_path":
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

			m.FilePath = temp2

		case "camera_id":
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

			m.CameraID = temp2

		case "event_id":
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

			m.EventID = &temp2

		}
	}

	return nil
}

func (m *Image) Reload(
	ctx context.Context,
	tx *sqlx.Tx,
) error {
	t, err := SelectImage(
		ctx,
		tx,
		fmt.Sprintf("%v = $1", m.GetPrimaryKeyColumn()),
		m.GetPrimaryKeyValue(),
	)
	if err != nil {
		return err
	}

	m.ID = t.ID
	m.Timestamp = t.Timestamp
	m.Size = t.Size
	m.FilePath = t.FilePath
	m.CameraID = t.CameraID
	m.CameraIDObject = t.CameraIDObject
	m.EventID = t.EventID
	m.EventIDObject = t.EventIDObject

	return nil
}

func (m *Image) Insert(
	ctx context.Context,
	tx *sqlx.Tx,
	setPrimaryKey bool,
	setZeroValues bool,
) error {
	columns := make([]string, 0)
	values := make([]any, 0)

	if setPrimaryKey && (setZeroValues || !types.IsZeroInt(m.ID)) {
		columns = append(columns, ImageTableIDColumn)

		v, err := types.FormatInt(m.ID)
		if err != nil {
			return fmt.Errorf("failed to handle m.ID: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.Timestamp) {
		columns = append(columns, ImageTableTimestampColumn)

		v, err := types.FormatTime(m.Timestamp)
		if err != nil {
			return fmt.Errorf("failed to handle m.Timestamp: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroFloat(m.Size) {
		columns = append(columns, ImageTableSizeColumn)

		v, err := types.FormatFloat(m.Size)
		if err != nil {
			return fmt.Errorf("failed to handle m.Size: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroString(m.FilePath) {
		columns = append(columns, ImageTableFilePathColumn)

		v, err := types.FormatString(m.FilePath)
		if err != nil {
			return fmt.Errorf("failed to handle m.FilePath: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroInt(m.CameraID) {
		columns = append(columns, ImageTableCameraIDColumn)

		v, err := types.FormatInt(m.CameraID)
		if err != nil {
			return fmt.Errorf("failed to handle m.CameraID: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroInt(m.EventID) {
		columns = append(columns, ImageTableEventIDColumn)

		v, err := types.FormatInt(m.EventID)
		if err != nil {
			return fmt.Errorf("failed to handle m.EventID: %v", err)
		}

		values = append(values, v)
	}

	item, err := query.Insert(
		ctx,
		tx,
		ImageTable,
		columns,
		nil,
		false,
		false,
		ImageTableColumns,
		values...,
	)
	if err != nil {
		return fmt.Errorf("failed to insert %#+v: %v", m, err)
	}
	v := item[ImageTableIDColumn]

	if v == nil {
		return fmt.Errorf("failed to find %v in %#+v", ImageTableIDColumn, item)
	}

	wrapError := func(err error) error {
		return fmt.Errorf(
			"failed to treat %v: %#+v as int64: %v",
			ImageTableIDColumn,
			item[ImageTableIDColumn],
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

func (m *Image) Update(
	ctx context.Context,
	tx *sqlx.Tx,
	setZeroValues bool,
) error {
	columns := make([]string, 0)
	values := make([]any, 0)

	if setZeroValues || !types.IsZeroTime(m.Timestamp) {
		columns = append(columns, ImageTableTimestampColumn)

		v, err := types.FormatTime(m.Timestamp)
		if err != nil {
			return fmt.Errorf("failed to handle m.Timestamp: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroFloat(m.Size) {
		columns = append(columns, ImageTableSizeColumn)

		v, err := types.FormatFloat(m.Size)
		if err != nil {
			return fmt.Errorf("failed to handle m.Size: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroString(m.FilePath) {
		columns = append(columns, ImageTableFilePathColumn)

		v, err := types.FormatString(m.FilePath)
		if err != nil {
			return fmt.Errorf("failed to handle m.FilePath: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroInt(m.CameraID) {
		columns = append(columns, ImageTableCameraIDColumn)

		v, err := types.FormatInt(m.CameraID)
		if err != nil {
			return fmt.Errorf("failed to handle m.CameraID: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroInt(m.EventID) {
		columns = append(columns, ImageTableEventIDColumn)

		v, err := types.FormatInt(m.EventID)
		if err != nil {
			return fmt.Errorf("failed to handle m.EventID: %v", err)
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
		ImageTable,
		columns,
		fmt.Sprintf("%v = $$??", ImageTableIDColumn),
		ImageTableColumns,
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

func (m *Image) Delete(
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
		ImageTable,
		fmt.Sprintf("%v = $$??", ImageTableIDColumn),
		values...,
	)
	if err != nil {
		return fmt.Errorf("failed to delete %#+v: %v", m, err)
	}

	return nil
}

func SelectImages(
	ctx context.Context,
	tx *sqlx.Tx,
	where string,
	limit *int,
	offset *int,
	values ...any,
) ([]*Image, error) {
	items, err := query.Select(
		ctx,
		tx,
		ImageTableColumnsWithTypeCasts,
		ImageTable,
		where,
		limit,
		offset,
		values...,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to call SelectImages; err: %v", err)
	}

	objects := make([]*Image, 0)

	for _, item := range items {
		object := &Image{}

		err = object.FromItem(item)
		if err != nil {
			return nil, fmt.Errorf("failed to call Image.FromItem; err: %v", err)
		}

		if !types.IsZeroInt(object.CameraID) {
			object.CameraIDObject, _ = SelectCamera(
				ctx,
				tx,
				fmt.Sprintf("%v = $1", CameraTablePrimaryKeyColumn),
				object.CameraID,
			)
		}

		if !types.IsZeroInt(object.EventID) {
			object.EventIDObject, _ = SelectEvent(
				ctx,
				tx,
				fmt.Sprintf("%v = $1", EventTablePrimaryKeyColumn),
				object.EventID,
			)
		}

		objects = append(objects, object)
	}

	return objects, nil
}

func SelectImage(
	ctx context.Context,
	tx *sqlx.Tx,
	where string,
	values ...any,
) (*Image, error) {
	objects, err := SelectImages(
		ctx,
		tx,
		where,
		helpers.Ptr(2),
		helpers.Ptr(0),
		values...,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to call SelectImage; err: %v", err)
	}

	if len(objects) > 1 {
		return nil, fmt.Errorf("attempt to call SelectImage returned more than 1 row")
	}

	if len(objects) < 1 {
		return nil, fmt.Errorf("attempt to call SelectImage returned no rows")
	}

	object := objects[0]

	return object, nil
}

func handleGetImages(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
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
			column := ImageTableColumnLookup[parts[0]]
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

	objects, err := SelectImages(ctx, tx, where, &limit, &offset, values...)
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

func handleGetImage(w http.ResponseWriter, r *http.Request, db *sqlx.DB, primaryKey string) {
	ctx := r.Context()

	where := fmt.Sprintf("%s = $$??", ImageTablePrimaryKeyColumn)

	values := []any{primaryKey}

	tx, err := db.BeginTxx(r.Context(), nil)
	if err != nil {
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	object, err := SelectImage(ctx, tx, where, values...)
	if err != nil {
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	err = tx.Commit()
	if err != nil {
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	helpers.HandleObjectsResponse(w, http.StatusOK, []*Image{object})
}

func handlePostImages(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
}

func handlePutImage(w http.ResponseWriter, r *http.Request, db *sqlx.DB, primaryKey string) {
}

func handlePatchImage(w http.ResponseWriter, r *http.Request, db *sqlx.DB, primaryKey string) {
}

func GetImageRouter(db *sqlx.DB, middlewares ...func(http.Handler) http.Handler) chi.Router {
	r := chi.NewRouter()

	for _, m := range middlewares {
		r.Use(m)
	}

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		handleGetImages(w, r, db)
	})

	r.Get("/{primaryKey}", func(w http.ResponseWriter, r *http.Request) {
		handleGetImage(w, r, db, chi.URLParam(r, "primaryKey"))
	})

	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		handlePostImages(w, r, db)
	})

	r.Put("/{primaryKey}", func(w http.ResponseWriter, r *http.Request) {
		handlePutImage(w, r, db, chi.URLParam(r, "primaryKey"))
	})

	r.Patch("/{primaryKey}", func(w http.ResponseWriter, r *http.Request) {
		handlePatchImage(w, r, db, chi.URLParam(r, "primaryKey"))
	})

	return r
}

func GetImageHandlerFunc(db *sqlx.DB, middlewares ...func(http.Handler) http.Handler) http.HandlerFunc {
	r := chi.NewRouter()

	r.Mount("/images", GetImageRouter(db, middlewares...))

	return r.ServeHTTP
}

func NewImageFromItem(item map[string]any) (any, error) {
	object := &Image{}

	err := object.FromItem(item)
	if err != nil {
		return nil, err
	}

	return object, nil
}

func init() {
	register(
		ImageTable,
		NewImageFromItem,
		"/images",
		GetImageRouter,
	)
}
