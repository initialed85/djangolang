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

type Event struct {
	ID                     int64         `json:"id"`
	StartTimestamp         time.Time     `json:"start_timestamp"`
	EndTimestamp           time.Time     `json:"end_timestamp"`
	Duration               time.Duration `json:"duration"`
	OriginalVideoID        int64         `json:"original_video_id"`
	OriginalVideoIDObject  *Video        `json:"original_video_id_object"`
	ThumbnailImageID       int64         `json:"thumbnail_image_id"`
	ThumbnailImageIDObject *Image        `json:"thumbnail_image_id_object"`
	ProcessedVideoID       *int64        `json:"processed_video_id"`
	ProcessedVideoIDObject *Video        `json:"processed_video_id_object"`
	SourceCameraID         int64         `json:"source_camera_id"`
	SourceCameraIDObject   *Camera       `json:"source_camera_id_object"`
	Status                 string        `json:"status"`
}

var EventTable = "event"

var (
	EventTableIDColumn               = "id"
	EventTableStartTimestampColumn   = "start_timestamp"
	EventTableEndTimestampColumn     = "end_timestamp"
	EventTableDurationColumn         = "duration"
	EventTableOriginalVideoIDColumn  = "original_video_id"
	EventTableThumbnailImageIDColumn = "thumbnail_image_id"
	EventTableProcessedVideoIDColumn = "processed_video_id"
	EventTableSourceCameraIDColumn   = "source_camera_id"
	EventTableStatusColumn           = "status"
)

var (
	EventTableIDColumnWithTypeCast               = fmt.Sprintf(`"id" AS id`)
	EventTableStartTimestampColumnWithTypeCast   = fmt.Sprintf(`"start_timestamp" AS start_timestamp`)
	EventTableEndTimestampColumnWithTypeCast     = fmt.Sprintf(`"end_timestamp" AS end_timestamp`)
	EventTableDurationColumnWithTypeCast         = fmt.Sprintf(`"duration" AS duration`)
	EventTableOriginalVideoIDColumnWithTypeCast  = fmt.Sprintf(`"original_video_id" AS original_video_id`)
	EventTableThumbnailImageIDColumnWithTypeCast = fmt.Sprintf(`"thumbnail_image_id" AS thumbnail_image_id`)
	EventTableProcessedVideoIDColumnWithTypeCast = fmt.Sprintf(`"processed_video_id" AS processed_video_id`)
	EventTableSourceCameraIDColumnWithTypeCast   = fmt.Sprintf(`"source_camera_id" AS source_camera_id`)
	EventTableStatusColumnWithTypeCast           = fmt.Sprintf(`"status" AS status`)
)

var EventTableColumns = []string{
	EventTableIDColumn,
	EventTableStartTimestampColumn,
	EventTableEndTimestampColumn,
	EventTableDurationColumn,
	EventTableOriginalVideoIDColumn,
	EventTableThumbnailImageIDColumn,
	EventTableProcessedVideoIDColumn,
	EventTableSourceCameraIDColumn,
	EventTableStatusColumn,
}

var EventTableColumnsWithTypeCasts = []string{
	EventTableIDColumnWithTypeCast,
	EventTableStartTimestampColumnWithTypeCast,
	EventTableEndTimestampColumnWithTypeCast,
	EventTableDurationColumnWithTypeCast,
	EventTableOriginalVideoIDColumnWithTypeCast,
	EventTableThumbnailImageIDColumnWithTypeCast,
	EventTableProcessedVideoIDColumnWithTypeCast,
	EventTableSourceCameraIDColumnWithTypeCast,
	EventTableStatusColumnWithTypeCast,
}

var EventTableColumnLookup = map[string]*introspect.Column{
	EventTableIDColumn:               new(introspect.Column),
	EventTableStartTimestampColumn:   new(introspect.Column),
	EventTableEndTimestampColumn:     new(introspect.Column),
	EventTableDurationColumn:         new(introspect.Column),
	EventTableOriginalVideoIDColumn:  new(introspect.Column),
	EventTableThumbnailImageIDColumn: new(introspect.Column),
	EventTableProcessedVideoIDColumn: new(introspect.Column),
	EventTableSourceCameraIDColumn:   new(introspect.Column),
	EventTableStatusColumn:           new(introspect.Column),
}

var (
	EventTablePrimaryKeyColumn = EventTableIDColumn
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

func (m *Event) GetPrimaryKeyColumn() string {
	return EventTablePrimaryKeyColumn
}

func (m *Event) GetPrimaryKeyValue() any {
	return m.ID
}

func (m *Event) FromItem(item map[string]any) error {
	if item == nil {
		return fmt.Errorf(
			"item unexpectedly nil during EventFromItem",
		)
	}

	if len(item) == 0 {
		return fmt.Errorf(
			"item unexpectedly empty during EventFromItem",
		)
	}

	wrapError := func(k string, err error) error {
		return fmt.Errorf("%#+v: %v; item: %#+v", k, err, item)
	}

	for k, v := range item {
		_, ok := EventTableColumnLookup[k]
		if !ok {
			return fmt.Errorf(
				"item contained unexpected key %#+v during EventFromItem; item: %#+v",
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

			temp1, err := types.ParseNotImplemented(v)
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

		case "original_video_id":
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

			m.OriginalVideoID = temp2

		case "thumbnail_image_id":
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

			m.ThumbnailImageID = temp2

		case "processed_video_id":
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

			m.ProcessedVideoID = &temp2

		case "source_camera_id":
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

			m.SourceCameraID = temp2

		case "status":
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

			m.Status = temp2

		}
	}

	return nil
}

func (m *Event) Reload(
	ctx context.Context,
	tx *sqlx.Tx,
) error {
	t, err := SelectEvent(
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
	m.OriginalVideoID = t.OriginalVideoID
	m.OriginalVideoIDObject = t.OriginalVideoIDObject
	m.ThumbnailImageID = t.ThumbnailImageID
	m.ThumbnailImageIDObject = t.ThumbnailImageIDObject
	m.ProcessedVideoID = t.ProcessedVideoID
	m.ProcessedVideoIDObject = t.ProcessedVideoIDObject
	m.SourceCameraID = t.SourceCameraID
	m.SourceCameraIDObject = t.SourceCameraIDObject
	m.Status = t.Status

	return nil
}

func (m *Event) Insert(
	ctx context.Context,
	tx *sqlx.Tx,
	setPrimaryKey bool,
	setZeroValues bool,
) error {
	columns := make([]string, 0)
	values := make([]any, 0)

	if setPrimaryKey && (setZeroValues || !types.IsZeroInt(m.ID)) {
		columns = append(columns, EventTableIDColumn)

		v, err := types.FormatInt(m.ID)
		if err != nil {
			return fmt.Errorf("failed to handle m.ID: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.StartTimestamp) {
		columns = append(columns, EventTableStartTimestampColumn)

		v, err := types.FormatTime(m.StartTimestamp)
		if err != nil {
			return fmt.Errorf("failed to handle m.StartTimestamp: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.EndTimestamp) {
		columns = append(columns, EventTableEndTimestampColumn)

		v, err := types.FormatTime(m.EndTimestamp)
		if err != nil {
			return fmt.Errorf("failed to handle m.EndTimestamp: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroNotImplemented(m.Duration) {
		columns = append(columns, EventTableDurationColumn)

		v, err := types.FormatNotImplemented(m.Duration)
		if err != nil {
			return fmt.Errorf("failed to handle m.Duration: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroInt(m.OriginalVideoID) {
		columns = append(columns, EventTableOriginalVideoIDColumn)

		v, err := types.FormatInt(m.OriginalVideoID)
		if err != nil {
			return fmt.Errorf("failed to handle m.OriginalVideoID: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroInt(m.ThumbnailImageID) {
		columns = append(columns, EventTableThumbnailImageIDColumn)

		v, err := types.FormatInt(m.ThumbnailImageID)
		if err != nil {
			return fmt.Errorf("failed to handle m.ThumbnailImageID: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroInt(m.ProcessedVideoID) {
		columns = append(columns, EventTableProcessedVideoIDColumn)

		v, err := types.FormatInt(m.ProcessedVideoID)
		if err != nil {
			return fmt.Errorf("failed to handle m.ProcessedVideoID: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroInt(m.SourceCameraID) {
		columns = append(columns, EventTableSourceCameraIDColumn)

		v, err := types.FormatInt(m.SourceCameraID)
		if err != nil {
			return fmt.Errorf("failed to handle m.SourceCameraID: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroString(m.Status) {
		columns = append(columns, EventTableStatusColumn)

		v, err := types.FormatString(m.Status)
		if err != nil {
			return fmt.Errorf("failed to handle m.Status: %v", err)
		}

		values = append(values, v)
	}

	item, err := query.Insert(
		ctx,
		tx,
		EventTable,
		columns,
		nil,
		false,
		false,
		EventTableColumns,
		values...,
	)
	if err != nil {
		return fmt.Errorf("failed to insert %#+v: %v", m, err)
	}
	v := item[EventTableIDColumn]

	if v == nil {
		return fmt.Errorf("failed to find %v in %#+v", EventTableIDColumn, item)
	}

	wrapError := func(err error) error {
		return fmt.Errorf(
			"failed to treat %v: %#+v as int64: %v",
			EventTableIDColumn,
			item[EventTableIDColumn],
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

func (m *Event) Update(
	ctx context.Context,
	tx *sqlx.Tx,
	setZeroValues bool,
) error {
	columns := make([]string, 0)
	values := make([]any, 0)

	if setZeroValues || !types.IsZeroTime(m.StartTimestamp) {
		columns = append(columns, EventTableStartTimestampColumn)

		v, err := types.FormatTime(m.StartTimestamp)
		if err != nil {
			return fmt.Errorf("failed to handle m.StartTimestamp: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.EndTimestamp) {
		columns = append(columns, EventTableEndTimestampColumn)

		v, err := types.FormatTime(m.EndTimestamp)
		if err != nil {
			return fmt.Errorf("failed to handle m.EndTimestamp: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroNotImplemented(m.Duration) {
		columns = append(columns, EventTableDurationColumn)

		v, err := types.FormatNotImplemented(m.Duration)
		if err != nil {
			return fmt.Errorf("failed to handle m.Duration: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroInt(m.OriginalVideoID) {
		columns = append(columns, EventTableOriginalVideoIDColumn)

		v, err := types.FormatInt(m.OriginalVideoID)
		if err != nil {
			return fmt.Errorf("failed to handle m.OriginalVideoID: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroInt(m.ThumbnailImageID) {
		columns = append(columns, EventTableThumbnailImageIDColumn)

		v, err := types.FormatInt(m.ThumbnailImageID)
		if err != nil {
			return fmt.Errorf("failed to handle m.ThumbnailImageID: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroInt(m.ProcessedVideoID) {
		columns = append(columns, EventTableProcessedVideoIDColumn)

		v, err := types.FormatInt(m.ProcessedVideoID)
		if err != nil {
			return fmt.Errorf("failed to handle m.ProcessedVideoID: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroInt(m.SourceCameraID) {
		columns = append(columns, EventTableSourceCameraIDColumn)

		v, err := types.FormatInt(m.SourceCameraID)
		if err != nil {
			return fmt.Errorf("failed to handle m.SourceCameraID: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroString(m.Status) {
		columns = append(columns, EventTableStatusColumn)

		v, err := types.FormatString(m.Status)
		if err != nil {
			return fmt.Errorf("failed to handle m.Status: %v", err)
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
		EventTable,
		columns,
		fmt.Sprintf("%v = $$??", EventTableIDColumn),
		EventTableColumns,
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

func (m *Event) Delete(
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
		EventTable,
		fmt.Sprintf("%v = $$??", EventTableIDColumn),
		values...,
	)
	if err != nil {
		return fmt.Errorf("failed to delete %#+v: %v", m, err)
	}

	return nil
}

func SelectEvents(
	ctx context.Context,
	tx *sqlx.Tx,
	where string,
	limit *int,
	offset *int,
	values ...any,
) ([]*Event, error) {
	items, err := query.Select(
		ctx,
		tx,
		EventTableColumnsWithTypeCasts,
		EventTable,
		where,
		limit,
		offset,
		values...,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to call SelectEvents; err: %v", err)
	}

	objects := make([]*Event, 0)

	for _, item := range items {
		object := &Event{}

		err = object.FromItem(item)
		if err != nil {
			return nil, fmt.Errorf("failed to call Event.FromItem; err: %v", err)
		}

		if !types.IsZeroInt(object.OriginalVideoID) {
			object.OriginalVideoIDObject, _ = SelectVideo(
				ctx,
				tx,
				fmt.Sprintf("%v = $1", VideoTablePrimaryKeyColumn),
				object.OriginalVideoID,
			)
		}

		if !types.IsZeroInt(object.ThumbnailImageID) {
			object.ThumbnailImageIDObject, _ = SelectImage(
				ctx,
				tx,
				fmt.Sprintf("%v = $1", ImageTablePrimaryKeyColumn),
				object.ThumbnailImageID,
			)
		}

		if !types.IsZeroInt(object.ProcessedVideoID) {
			object.ProcessedVideoIDObject, _ = SelectVideo(
				ctx,
				tx,
				fmt.Sprintf("%v = $1", VideoTablePrimaryKeyColumn),
				object.ProcessedVideoID,
			)
		}

		if !types.IsZeroInt(object.SourceCameraID) {
			object.SourceCameraIDObject, _ = SelectCamera(
				ctx,
				tx,
				fmt.Sprintf("%v = $1", CameraTablePrimaryKeyColumn),
				object.SourceCameraID,
			)
		}

		objects = append(objects, object)
	}

	return objects, nil
}

func SelectEvent(
	ctx context.Context,
	tx *sqlx.Tx,
	where string,
	values ...any,
) (*Event, error) {
	objects, err := SelectEvents(
		ctx,
		tx,
		where,
		helpers.Ptr(2),
		helpers.Ptr(0),
		values...,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to call SelectEvent; err: %v", err)
	}

	if len(objects) > 1 {
		return nil, fmt.Errorf("attempt to call SelectEvent returned more than 1 row")
	}

	if len(objects) < 1 {
		return nil, fmt.Errorf("attempt to call SelectEvent returned no rows")
	}

	object := objects[0]

	return object, nil
}

func handleGetEvents(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
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
			column := EventTableColumnLookup[parts[0]]
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

	objects, err := SelectEvents(ctx, tx, where, &limit, &offset, values...)
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

func handleGetEvent(w http.ResponseWriter, r *http.Request, db *sqlx.DB, primaryKey string) {
	ctx := r.Context()

	where := fmt.Sprintf("%s = $$??", EventTablePrimaryKeyColumn)

	values := []any{primaryKey}

	tx, err := db.BeginTxx(r.Context(), nil)
	if err != nil {
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	object, err := SelectEvent(ctx, tx, where, values...)
	if err != nil {
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	err = tx.Commit()
	if err != nil {
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	helpers.HandleObjectsResponse(w, http.StatusOK, []*Event{object})
}

func handlePostEvents(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
}

func handlePutEvent(w http.ResponseWriter, r *http.Request, db *sqlx.DB, primaryKey string) {
}

func handlePatchEvent(w http.ResponseWriter, r *http.Request, db *sqlx.DB, primaryKey string) {
}

func GetEventRouter(db *sqlx.DB, middlewares ...func(http.Handler) http.Handler) chi.Router {
	r := chi.NewRouter()

	for _, m := range middlewares {
		r.Use(m)
	}

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		handleGetEvents(w, r, db)
	})

	r.Get("/{primaryKey}", func(w http.ResponseWriter, r *http.Request) {
		handleGetEvent(w, r, db, chi.URLParam(r, "primaryKey"))
	})

	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		handlePostEvents(w, r, db)
	})

	r.Put("/{primaryKey}", func(w http.ResponseWriter, r *http.Request) {
		handlePutEvent(w, r, db, chi.URLParam(r, "primaryKey"))
	})

	r.Patch("/{primaryKey}", func(w http.ResponseWriter, r *http.Request) {
		handlePatchEvent(w, r, db, chi.URLParam(r, "primaryKey"))
	})

	return r
}

func GetEventHandlerFunc(db *sqlx.DB, middlewares ...func(http.Handler) http.Handler) http.HandlerFunc {
	r := chi.NewRouter()

	r.Mount("/events", GetEventRouter(db, middlewares...))

	return r.ServeHTTP
}

func NewEventFromItem(item map[string]any) (any, error) {
	object := &Event{}

	err := object.FromItem(item)
	if err != nil {
		return nil, err
	}

	return object, nil
}

func init() {
	register(
		EventTable,
		NewEventFromItem,
		"/events",
		GetEventRouter,
	)
}
