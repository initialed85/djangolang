package cameranator

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
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

type Video struct {
	ID             int64         `json:"id"`
	StartTimestamp time.Time     `json:"start_timestamp"`
	EndTimestamp   time.Time     `json:"end_timestamp"`
	Duration       time.Duration `json:"duration"`
	Size           float64       `json:"size"`
	FilePath       string        `json:"file_path"`
	CameraID       int64         `json:"camera_id"`
	CameraIDObject *Camera       `json:"camera_id_object"`
	EventID        *int64        `json:"event_id"`
	EventIDObject  *Event        `json:"event_id_object"`
}

var VideoTable = "video"

var (
	VideoTableIDColumn             = "id"
	VideoTableStartTimestampColumn = "start_timestamp"
	VideoTableEndTimestampColumn   = "end_timestamp"
	VideoTableDurationColumn       = "duration"
	VideoTableSizeColumn           = "size"
	VideoTableFilePathColumn       = "file_path"
	VideoTableCameraIDColumn       = "camera_id"
	VideoTableEventIDColumn        = "event_id"
)

var (
	VideoTableIDColumnWithTypeCast             = fmt.Sprintf(`"id" AS id`)
	VideoTableStartTimestampColumnWithTypeCast = fmt.Sprintf(`"start_timestamp" AS start_timestamp`)
	VideoTableEndTimestampColumnWithTypeCast   = fmt.Sprintf(`"end_timestamp" AS end_timestamp`)
	VideoTableDurationColumnWithTypeCast       = fmt.Sprintf(`"duration" AS duration`)
	VideoTableSizeColumnWithTypeCast           = fmt.Sprintf(`"size" AS size`)
	VideoTableFilePathColumnWithTypeCast       = fmt.Sprintf(`"file_path" AS file_path`)
	VideoTableCameraIDColumnWithTypeCast       = fmt.Sprintf(`"camera_id" AS camera_id`)
	VideoTableEventIDColumnWithTypeCast        = fmt.Sprintf(`"event_id" AS event_id`)
)

var VideoTableColumns = []string{
	VideoTableIDColumn,
	VideoTableStartTimestampColumn,
	VideoTableEndTimestampColumn,
	VideoTableDurationColumn,
	VideoTableSizeColumn,
	VideoTableFilePathColumn,
	VideoTableCameraIDColumn,
	VideoTableEventIDColumn,
}

var VideoTableColumnsWithTypeCasts = []string{
	VideoTableIDColumnWithTypeCast,
	VideoTableStartTimestampColumnWithTypeCast,
	VideoTableEndTimestampColumnWithTypeCast,
	VideoTableDurationColumnWithTypeCast,
	VideoTableSizeColumnWithTypeCast,
	VideoTableFilePathColumnWithTypeCast,
	VideoTableCameraIDColumnWithTypeCast,
	VideoTableEventIDColumnWithTypeCast,
}

var VideoTableColumnLookup = map[string]*introspect.Column{
	VideoTableIDColumn:             new(introspect.Column),
	VideoTableStartTimestampColumn: new(introspect.Column),
	VideoTableEndTimestampColumn:   new(introspect.Column),
	VideoTableDurationColumn:       new(introspect.Column),
	VideoTableSizeColumn:           new(introspect.Column),
	VideoTableFilePathColumn:       new(introspect.Column),
	VideoTableCameraIDColumn:       new(introspect.Column),
	VideoTableEventIDColumn:        new(introspect.Column),
}

var (
	VideoTablePrimaryKeyColumn = VideoTableIDColumn
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
)

func (m *Video) GetPrimaryKeyColumn() string {
	return VideoTablePrimaryKeyColumn
}

func (m *Video) GetPrimaryKeyValue() any {
	return m.ID
}

func (m *Video) FromItem(item map[string]any) error {
	if item == nil {
		return fmt.Errorf(
			"item unexpectedly nil during VideoFromItem",
		)
	}

	if len(item) == 0 {
		return fmt.Errorf(
			"item unexpectedly empty during VideoFromItem",
		)
	}

	wrapError := func(k string, err error) error {
		return fmt.Errorf("%#+v: %v; item: %#+v", k, err, item)
	}

	for k, v := range item {
		_, ok := VideoTableColumnLookup[k]
		if !ok {
			return fmt.Errorf(
				"item contained unexpected key %#+v during VideoFromItem; item: %#+v",
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

		case "duration":
			if v == nil {
				continue
			}

			temp1, err := types.ParseDuration(v)
			if err != nil {
				return wrapError(k, err)
			}

			temp2, ok := temp1.(time.Duration)
			if !ok {
				if temp1 != nil {
					return wrapError(k, fmt.Errorf("failed to cast %#+v to time.Duration", temp1))
				}
			}

			m.Duration = temp2

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

func (m *Video) Reload(
	ctx context.Context,
	tx *sqlx.Tx,
) error {
	t, err := SelectVideo(
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
	m.Duration = t.Duration
	m.Size = t.Size
	m.FilePath = t.FilePath
	m.CameraID = t.CameraID
	m.CameraIDObject = t.CameraIDObject
	m.EventID = t.EventID
	m.EventIDObject = t.EventIDObject

	return nil
}

func (m *Video) Insert(
	ctx context.Context,
	tx *sqlx.Tx,
	setPrimaryKey bool,
	setZeroValues bool,
) error {
	columns := make([]string, 0)
	values := make([]any, 0)

	if setPrimaryKey && (setZeroValues || !types.IsZeroInt(m.ID)) {
		columns = append(columns, VideoTableIDColumn)

		v, err := types.FormatInt(m.ID)
		if err != nil {
			return fmt.Errorf("failed to handle m.ID: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.StartTimestamp) {
		columns = append(columns, VideoTableStartTimestampColumn)

		v, err := types.FormatTime(m.StartTimestamp)
		if err != nil {
			return fmt.Errorf("failed to handle m.StartTimestamp: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.EndTimestamp) {
		columns = append(columns, VideoTableEndTimestampColumn)

		v, err := types.FormatTime(m.EndTimestamp)
		if err != nil {
			return fmt.Errorf("failed to handle m.EndTimestamp: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroDuration(m.Duration) {
		columns = append(columns, VideoTableDurationColumn)

		v, err := types.FormatDuration(m.Duration)
		if err != nil {
			return fmt.Errorf("failed to handle m.Duration: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroFloat(m.Size) {
		columns = append(columns, VideoTableSizeColumn)

		v, err := types.FormatFloat(m.Size)
		if err != nil {
			return fmt.Errorf("failed to handle m.Size: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroString(m.FilePath) {
		columns = append(columns, VideoTableFilePathColumn)

		v, err := types.FormatString(m.FilePath)
		if err != nil {
			return fmt.Errorf("failed to handle m.FilePath: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroInt(m.CameraID) {
		columns = append(columns, VideoTableCameraIDColumn)

		v, err := types.FormatInt(m.CameraID)
		if err != nil {
			return fmt.Errorf("failed to handle m.CameraID: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroInt(m.EventID) {
		columns = append(columns, VideoTableEventIDColumn)

		v, err := types.FormatInt(m.EventID)
		if err != nil {
			return fmt.Errorf("failed to handle m.EventID: %v", err)
		}

		values = append(values, v)
	}

	item, err := query.Insert(
		ctx,
		tx,
		VideoTable,
		columns,
		nil,
		false,
		false,
		VideoTableColumns,
		values...,
	)
	if err != nil {
		return fmt.Errorf("failed to insert %#+v: %v", m, err)
	}
	v := item[VideoTableIDColumn]

	if v == nil {
		return fmt.Errorf("failed to find %v in %#+v", VideoTableIDColumn, item)
	}

	wrapError := func(err error) error {
		return fmt.Errorf(
			"failed to treat %v: %#+v as int64: %v",
			VideoTableIDColumn,
			item[VideoTableIDColumn],
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

func (m *Video) Update(
	ctx context.Context,
	tx *sqlx.Tx,
	setZeroValues bool,
) error {
	columns := make([]string, 0)
	values := make([]any, 0)

	if setZeroValues || !types.IsZeroTime(m.StartTimestamp) {
		columns = append(columns, VideoTableStartTimestampColumn)

		v, err := types.FormatTime(m.StartTimestamp)
		if err != nil {
			return fmt.Errorf("failed to handle m.StartTimestamp: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.EndTimestamp) {
		columns = append(columns, VideoTableEndTimestampColumn)

		v, err := types.FormatTime(m.EndTimestamp)
		if err != nil {
			return fmt.Errorf("failed to handle m.EndTimestamp: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroDuration(m.Duration) {
		columns = append(columns, VideoTableDurationColumn)

		v, err := types.FormatDuration(m.Duration)
		if err != nil {
			return fmt.Errorf("failed to handle m.Duration: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroFloat(m.Size) {
		columns = append(columns, VideoTableSizeColumn)

		v, err := types.FormatFloat(m.Size)
		if err != nil {
			return fmt.Errorf("failed to handle m.Size: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroString(m.FilePath) {
		columns = append(columns, VideoTableFilePathColumn)

		v, err := types.FormatString(m.FilePath)
		if err != nil {
			return fmt.Errorf("failed to handle m.FilePath: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroInt(m.CameraID) {
		columns = append(columns, VideoTableCameraIDColumn)

		v, err := types.FormatInt(m.CameraID)
		if err != nil {
			return fmt.Errorf("failed to handle m.CameraID: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroInt(m.EventID) {
		columns = append(columns, VideoTableEventIDColumn)

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
		VideoTable,
		columns,
		fmt.Sprintf("%v = $$??", VideoTableIDColumn),
		VideoTableColumns,
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

func (m *Video) Delete(
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
		VideoTable,
		fmt.Sprintf("%v = $$??", VideoTableIDColumn),
		values...,
	)
	if err != nil {
		return fmt.Errorf("failed to delete %#+v: %v", m, err)
	}

	return nil
}

func SelectVideos(
	ctx context.Context,
	tx *sqlx.Tx,
	where string,
	limit *int,
	offset *int,
	values ...any,
) ([]*Video, error) {
	items, err := query.Select(
		ctx,
		tx,
		VideoTableColumnsWithTypeCasts,
		VideoTable,
		where,
		limit,
		offset,
		values...,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to call SelectVideos; err: %v", err)
	}

	objects := make([]*Video, 0)

	for _, item := range items {
		object := &Video{}

		err = object.FromItem(item)
		if err != nil {
			return nil, fmt.Errorf("failed to call Video.FromItem; err: %v", err)
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

func SelectVideo(
	ctx context.Context,
	tx *sqlx.Tx,
	where string,
	values ...any,
) (*Video, error) {
	objects, err := SelectVideos(
		ctx,
		tx,
		where,
		helpers.Ptr(2),
		helpers.Ptr(0),
		values...,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to call SelectVideo; err: %v", err)
	}

	if len(objects) > 1 {
		return nil, fmt.Errorf("attempt to call SelectVideo returned more than 1 row")
	}

	if len(objects) < 1 {
		return nil, fmt.Errorf("attempt to call SelectVideo returned no rows")
	}

	object := objects[0]

	return object, nil
}

func handleGetVideos(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
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
			column := VideoTableColumnLookup[parts[0]]
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

	objects, err := SelectVideos(ctx, tx, where, &limit, &offset, values...)
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

func handleGetVideo(w http.ResponseWriter, r *http.Request, db *sqlx.DB, primaryKey string) {
	ctx := r.Context()

	where := fmt.Sprintf("%s = $$??", VideoTablePrimaryKeyColumn)

	values := []any{primaryKey}

	tx, err := db.BeginTxx(r.Context(), nil)
	if err != nil {
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	object, err := SelectVideo(ctx, tx, where, values...)
	if err != nil {
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	err = tx.Commit()
	if err != nil {
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	helpers.HandleObjectsResponse(w, http.StatusOK, []*Video{object})
}

func handlePostVideos(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
}

func handlePutVideo(w http.ResponseWriter, r *http.Request, db *sqlx.DB, primaryKey string) {
}

func handlePatchVideo(w http.ResponseWriter, r *http.Request, db *sqlx.DB, primaryKey string) {
}

func GetVideoRouter(db *sqlx.DB, middlewares ...func(http.Handler) http.Handler) chi.Router {
	r := chi.NewRouter()

	for _, m := range middlewares {
		r.Use(m)
	}

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		handleGetVideos(w, r, db)
	})

	r.Get("/{primaryKey}", func(w http.ResponseWriter, r *http.Request) {
		handleGetVideo(w, r, db, chi.URLParam(r, "primaryKey"))
	})

	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		handlePostVideos(w, r, db)
	})

	r.Put("/{primaryKey}", func(w http.ResponseWriter, r *http.Request) {
		handlePutVideo(w, r, db, chi.URLParam(r, "primaryKey"))
	})

	r.Patch("/{primaryKey}", func(w http.ResponseWriter, r *http.Request) {
		handlePatchVideo(w, r, db, chi.URLParam(r, "primaryKey"))
	})

	return r
}

func GetVideoHandlerFunc(db *sqlx.DB, middlewares ...func(http.Handler) http.Handler) http.HandlerFunc {
	r := chi.NewRouter()

	r.Mount("/videos", GetVideoRouter(db, middlewares...))

	return r.ServeHTTP
}

func NewVideoFromItem(item map[string]any) (any, error) {
	object := &Video{}

	err := object.FromItem(item)
	if err != nil {
		return nil, err
	}

	return object, nil
}

func init() {
	register(
		VideoTable,
		NewVideoFromItem,
		"/videos",
		GetVideoRouter,
	)
}
