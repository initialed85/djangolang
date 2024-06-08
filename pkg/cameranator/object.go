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

type Object struct {
	ID             int64     `json:"id"`
	StartTimestamp time.Time `json:"start_timestamp"`
	EndTimestamp   time.Time `json:"end_timestamp"`
	ClassID        int64     `json:"class_id"`
	ClassName      string    `json:"class_name"`
	CameraID       int64     `json:"camera_id"`
	CameraIDObject *Camera   `json:"camera_id_object"`
	EventID        *int64    `json:"event_id"`
	EventIDObject  *Event    `json:"event_id_object"`
}

var ObjectTable = "object"

var (
	ObjectTableIDColumn             = "id"
	ObjectTableStartTimestampColumn = "start_timestamp"
	ObjectTableEndTimestampColumn   = "end_timestamp"
	ObjectTableClassIDColumn        = "class_id"
	ObjectTableClassNameColumn      = "class_name"
	ObjectTableCameraIDColumn       = "camera_id"
	ObjectTableEventIDColumn        = "event_id"
)

var (
	ObjectTableIDColumnWithTypeCast             = fmt.Sprintf(`"id" AS id`)
	ObjectTableStartTimestampColumnWithTypeCast = fmt.Sprintf(`"start_timestamp" AS start_timestamp`)
	ObjectTableEndTimestampColumnWithTypeCast   = fmt.Sprintf(`"end_timestamp" AS end_timestamp`)
	ObjectTableClassIDColumnWithTypeCast        = fmt.Sprintf(`"class_id" AS class_id`)
	ObjectTableClassNameColumnWithTypeCast      = fmt.Sprintf(`"class_name" AS class_name`)
	ObjectTableCameraIDColumnWithTypeCast       = fmt.Sprintf(`"camera_id" AS camera_id`)
	ObjectTableEventIDColumnWithTypeCast        = fmt.Sprintf(`"event_id" AS event_id`)
)

var ObjectTableColumns = []string{
	ObjectTableIDColumn,
	ObjectTableStartTimestampColumn,
	ObjectTableEndTimestampColumn,
	ObjectTableClassIDColumn,
	ObjectTableClassNameColumn,
	ObjectTableCameraIDColumn,
	ObjectTableEventIDColumn,
}

var ObjectTableColumnsWithTypeCasts = []string{
	ObjectTableIDColumnWithTypeCast,
	ObjectTableStartTimestampColumnWithTypeCast,
	ObjectTableEndTimestampColumnWithTypeCast,
	ObjectTableClassIDColumnWithTypeCast,
	ObjectTableClassNameColumnWithTypeCast,
	ObjectTableCameraIDColumnWithTypeCast,
	ObjectTableEventIDColumnWithTypeCast,
}

var ObjectTableColumnLookup = map[string]*introspect.Column{
	ObjectTableIDColumn:             new(introspect.Column),
	ObjectTableStartTimestampColumn: new(introspect.Column),
	ObjectTableEndTimestampColumn:   new(introspect.Column),
	ObjectTableClassIDColumn:        new(introspect.Column),
	ObjectTableClassNameColumn:      new(introspect.Column),
	ObjectTableCameraIDColumn:       new(introspect.Column),
	ObjectTableEventIDColumn:        new(introspect.Column),
}

var (
	ObjectTablePrimaryKeyColumn = ObjectTableIDColumn
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

func (m *Object) GetPrimaryKeyColumn() string {
	return ObjectTablePrimaryKeyColumn
}

func (m *Object) GetPrimaryKeyValue() any {
	return m.ID
}

func (m *Object) FromItem(item map[string]any) error {
	if item == nil {
		return fmt.Errorf(
			"item unexpectedly nil during ObjectFromItem",
		)
	}

	if len(item) == 0 {
		return fmt.Errorf(
			"item unexpectedly empty during ObjectFromItem",
		)
	}

	wrapError := func(k string, err error) error {
		return fmt.Errorf("%#+v: %v; item: %#+v", k, err, item)
	}

	for k, v := range item {
		_, ok := ObjectTableColumnLookup[k]
		if !ok {
			return fmt.Errorf(
				"item contained unexpected key %#+v during ObjectFromItem; item: %#+v",
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

		case "start_timestamp":
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

			m.StartTimestamp = temp2

		case "end_timestamp":
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

			m.EndTimestamp = temp2

		case "class_id":
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

			m.ClassID = temp2

		case "class_name":
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

			m.ClassName = temp2

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

func (m *Object) Reload(
	ctx context.Context,
	tx *sqlx.Tx,
) error {
	t, err := SelectObject(
		ctx,
		tx,
		fmt.Sprintf("%v = $1", m.GetPrimaryKeyColumn()),
		m.GetPrimaryKeyValue(),
	)
	if err != nil {
		return err
	}

	m.ID = t.ID
	m.StartTimestamp = t.StartTimestamp
	m.EndTimestamp = t.EndTimestamp
	m.ClassID = t.ClassID
	m.ClassName = t.ClassName
	m.CameraID = t.CameraID
	m.CameraIDObject = t.CameraIDObject
	m.EventID = t.EventID
	m.EventIDObject = t.EventIDObject

	return nil
}

func (m *Object) Insert(
	ctx context.Context,
	tx *sqlx.Tx,
	setPrimaryKey bool,
	setZeroValues bool,
) error {
	columns := make([]string, 0)
	values := make([]any, 0)

	if setPrimaryKey && (setZeroValues || !types.IsZeroInt(m.ID)) {
		columns = append(columns, ObjectTableIDColumn)

		v, err := types.FormatInt(m.ID)
		if err != nil {
			return fmt.Errorf("failed to handle m.ID: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.StartTimestamp) {
		columns = append(columns, ObjectTableStartTimestampColumn)

		v, err := types.FormatTime(m.StartTimestamp)
		if err != nil {
			return fmt.Errorf("failed to handle m.StartTimestamp: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.EndTimestamp) {
		columns = append(columns, ObjectTableEndTimestampColumn)

		v, err := types.FormatTime(m.EndTimestamp)
		if err != nil {
			return fmt.Errorf("failed to handle m.EndTimestamp: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroInt(m.ClassID) {
		columns = append(columns, ObjectTableClassIDColumn)

		v, err := types.FormatInt(m.ClassID)
		if err != nil {
			return fmt.Errorf("failed to handle m.ClassID: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroString(m.ClassName) {
		columns = append(columns, ObjectTableClassNameColumn)

		v, err := types.FormatString(m.ClassName)
		if err != nil {
			return fmt.Errorf("failed to handle m.ClassName: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroInt(m.CameraID) {
		columns = append(columns, ObjectTableCameraIDColumn)

		v, err := types.FormatInt(m.CameraID)
		if err != nil {
			return fmt.Errorf("failed to handle m.CameraID: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroInt(m.EventID) {
		columns = append(columns, ObjectTableEventIDColumn)

		v, err := types.FormatInt(m.EventID)
		if err != nil {
			return fmt.Errorf("failed to handle m.EventID: %v", err)
		}

		values = append(values, v)
	}

	item, err := query.Insert(
		ctx,
		tx,
		ObjectTable,
		columns,
		nil,
		false,
		false,
		ObjectTableColumns,
		values...,
	)
	if err != nil {
		return fmt.Errorf("failed to insert %#+v: %v", m, err)
	}
	v := item[ObjectTableIDColumn]

	if v == nil {
		return fmt.Errorf("failed to find %v in %#+v", ObjectTableIDColumn, item)
	}

	wrapError := func(err error) error {
		return fmt.Errorf(
			"failed to treat %v: %#+v as int64: %v",
			ObjectTableIDColumn,
			item[ObjectTableIDColumn],
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

func (m *Object) Update(
	ctx context.Context,
	tx *sqlx.Tx,
	setZeroValues bool,
) error {
	columns := make([]string, 0)
	values := make([]any, 0)

	if setZeroValues || !types.IsZeroTime(m.StartTimestamp) {
		columns = append(columns, ObjectTableStartTimestampColumn)

		v, err := types.FormatTime(m.StartTimestamp)
		if err != nil {
			return fmt.Errorf("failed to handle m.StartTimestamp: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.EndTimestamp) {
		columns = append(columns, ObjectTableEndTimestampColumn)

		v, err := types.FormatTime(m.EndTimestamp)
		if err != nil {
			return fmt.Errorf("failed to handle m.EndTimestamp: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroInt(m.ClassID) {
		columns = append(columns, ObjectTableClassIDColumn)

		v, err := types.FormatInt(m.ClassID)
		if err != nil {
			return fmt.Errorf("failed to handle m.ClassID: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroString(m.ClassName) {
		columns = append(columns, ObjectTableClassNameColumn)

		v, err := types.FormatString(m.ClassName)
		if err != nil {
			return fmt.Errorf("failed to handle m.ClassName: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroInt(m.CameraID) {
		columns = append(columns, ObjectTableCameraIDColumn)

		v, err := types.FormatInt(m.CameraID)
		if err != nil {
			return fmt.Errorf("failed to handle m.CameraID: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroInt(m.EventID) {
		columns = append(columns, ObjectTableEventIDColumn)

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
		ObjectTable,
		columns,
		fmt.Sprintf("%v = $$??", ObjectTableIDColumn),
		ObjectTableColumns,
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

func (m *Object) Delete(
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
		ObjectTable,
		fmt.Sprintf("%v = $$??", ObjectTableIDColumn),
		values...,
	)
	if err != nil {
		return fmt.Errorf("failed to delete %#+v: %v", m, err)
	}

	return nil
}

func SelectObjects(
	ctx context.Context,
	tx *sqlx.Tx,
	where string,
	limit *int,
	offset *int,
	values ...any,
) ([]*Object, error) {
	items, err := query.Select(
		ctx,
		tx,
		ObjectTableColumnsWithTypeCasts,
		ObjectTable,
		where,
		limit,
		offset,
		values...,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to call SelectObjects; err: %v", err)
	}

	objects := make([]*Object, 0)

	for _, item := range items {
		object := &Object{}

		err = object.FromItem(item)
		if err != nil {
			return nil, fmt.Errorf("failed to call Object.FromItem; err: %v", err)
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

func SelectObject(
	ctx context.Context,
	tx *sqlx.Tx,
	where string,
	values ...any,
) (*Object, error) {
	objects, err := SelectObjects(
		ctx,
		tx,
		where,
		helpers.Ptr(2),
		helpers.Ptr(0),
		values...,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to call SelectObject; err: %v", err)
	}

	if len(objects) > 1 {
		return nil, fmt.Errorf("attempt to call SelectObject returned more than 1 row")
	}

	if len(objects) < 1 {
		return nil, fmt.Errorf("attempt to call SelectObject returned no rows")
	}

	object := objects[0]

	return object, nil
}

func handleGetObjects(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
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
			column := ObjectTableColumnLookup[parts[0]]
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

	objects, err := SelectObjects(ctx, tx, where, &limit, &offset, values...)
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

func handleGetObject(w http.ResponseWriter, r *http.Request, db *sqlx.DB, primaryKey string) {
	ctx := r.Context()

	where := fmt.Sprintf("%s = $$??", ObjectTablePrimaryKeyColumn)

	values := []any{primaryKey}

	tx, err := db.BeginTxx(r.Context(), nil)
	if err != nil {
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	object, err := SelectObject(ctx, tx, where, values...)
	if err != nil {
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	err = tx.Commit()
	if err != nil {
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	helpers.HandleObjectsResponse(w, http.StatusOK, []*Object{object})
}

func handlePostObjects(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
}

func handlePutObject(w http.ResponseWriter, r *http.Request, db *sqlx.DB, primaryKey string) {
}

func handlePatchObject(w http.ResponseWriter, r *http.Request, db *sqlx.DB, primaryKey string) {
}

func GetObjectRouter(db *sqlx.DB, middlewares ...func(http.Handler) http.Handler) chi.Router {
	r := chi.NewRouter()

	for _, m := range middlewares {
		r.Use(m)
	}

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		handleGetObjects(w, r, db)
	})

	r.Get("/{primaryKey}", func(w http.ResponseWriter, r *http.Request) {
		handleGetObject(w, r, db, chi.URLParam(r, "primaryKey"))
	})

	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		handlePostObjects(w, r, db)
	})

	r.Put("/{primaryKey}", func(w http.ResponseWriter, r *http.Request) {
		handlePutObject(w, r, db, chi.URLParam(r, "primaryKey"))
	})

	r.Patch("/{primaryKey}", func(w http.ResponseWriter, r *http.Request) {
		handlePatchObject(w, r, db, chi.URLParam(r, "primaryKey"))
	})

	return r
}

func GetObjectHandlerFunc(db *sqlx.DB, middlewares ...func(http.Handler) http.Handler) http.HandlerFunc {
	r := chi.NewRouter()

	r.Mount("/objects", GetObjectRouter(db, middlewares...))

	return r.ServeHTTP
}

func NewObjectFromItem(item map[string]any) (any, error) {
	object := &Object{}

	err := object.FromItem(item)
	if err != nil {
		return nil, err
	}

	return object, nil
}

func init() {
	register(
		ObjectTable,
		NewObjectFromItem,
		"/objects",
		GetObjectRouter,
	)
}
