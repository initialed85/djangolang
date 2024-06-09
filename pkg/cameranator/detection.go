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

type Detection struct {
	ID             int64           `json:"id"`
	Timestamp      time.Time       `json:"timestamp"`
	ClassID        int64           `json:"class_id"`
	ClassName      string          `json:"class_name"`
	Score          float64         `json:"score"`
	Centroid       pgtype.Vec2     `json:"centroid"`
	BoundingBox    []pgtype.Vec2   `json:"bounding_box"`
	CameraID       int64           `json:"camera_id"`
	CameraIDObject *Camera         `json:"camera_id_object"`
	EventID        *int64          `json:"event_id"`
	EventIDObject  *Event          `json:"event_id_object"`
	ObjectID       *int64          `json:"object_id"`
	ObjectIDObject *Object         `json:"object_id_object"`
	Colour         *postgis.PointZ `json:"colour"`
}

var DetectionTable = "detection"

var (
	DetectionTableIDColumn          = "id"
	DetectionTableTimestampColumn   = "timestamp"
	DetectionTableClassIDColumn     = "class_id"
	DetectionTableClassNameColumn   = "class_name"
	DetectionTableScoreColumn       = "score"
	DetectionTableCentroidColumn    = "centroid"
	DetectionTableBoundingBoxColumn = "bounding_box"
	DetectionTableCameraIDColumn    = "camera_id"
	DetectionTableEventIDColumn     = "event_id"
	DetectionTableObjectIDColumn    = "object_id"
	DetectionTableColourColumn      = "colour"
)

var (
	DetectionTableIDColumnWithTypeCast          = fmt.Sprintf(`"id" AS id`)
	DetectionTableTimestampColumnWithTypeCast   = fmt.Sprintf(`"timestamp" AS timestamp`)
	DetectionTableClassIDColumnWithTypeCast     = fmt.Sprintf(`"class_id" AS class_id`)
	DetectionTableClassNameColumnWithTypeCast   = fmt.Sprintf(`"class_name" AS class_name`)
	DetectionTableScoreColumnWithTypeCast       = fmt.Sprintf(`"score" AS score`)
	DetectionTableCentroidColumnWithTypeCast    = fmt.Sprintf(`"centroid" AS centroid`)
	DetectionTableBoundingBoxColumnWithTypeCast = fmt.Sprintf(`"bounding_box" AS bounding_box`)
	DetectionTableCameraIDColumnWithTypeCast    = fmt.Sprintf(`"camera_id" AS camera_id`)
	DetectionTableEventIDColumnWithTypeCast     = fmt.Sprintf(`"event_id" AS event_id`)
	DetectionTableObjectIDColumnWithTypeCast    = fmt.Sprintf(`"object_id" AS object_id`)
	DetectionTableColourColumnWithTypeCast      = fmt.Sprintf(`"colour" AS colour`)
)

var DetectionTableColumns = []string{
	DetectionTableIDColumn,
	DetectionTableTimestampColumn,
	DetectionTableClassIDColumn,
	DetectionTableClassNameColumn,
	DetectionTableScoreColumn,
	DetectionTableCentroidColumn,
	DetectionTableBoundingBoxColumn,
	DetectionTableCameraIDColumn,
	DetectionTableEventIDColumn,
	DetectionTableObjectIDColumn,
	DetectionTableColourColumn,
}

var DetectionTableColumnsWithTypeCasts = []string{
	DetectionTableIDColumnWithTypeCast,
	DetectionTableTimestampColumnWithTypeCast,
	DetectionTableClassIDColumnWithTypeCast,
	DetectionTableClassNameColumnWithTypeCast,
	DetectionTableScoreColumnWithTypeCast,
	DetectionTableCentroidColumnWithTypeCast,
	DetectionTableBoundingBoxColumnWithTypeCast,
	DetectionTableCameraIDColumnWithTypeCast,
	DetectionTableEventIDColumnWithTypeCast,
	DetectionTableObjectIDColumnWithTypeCast,
	DetectionTableColourColumnWithTypeCast,
}

var DetectionTableColumnLookup = map[string]*introspect.Column{
	DetectionTableIDColumn:          new(introspect.Column),
	DetectionTableTimestampColumn:   new(introspect.Column),
	DetectionTableClassIDColumn:     new(introspect.Column),
	DetectionTableClassNameColumn:   new(introspect.Column),
	DetectionTableScoreColumn:       new(introspect.Column),
	DetectionTableCentroidColumn:    new(introspect.Column),
	DetectionTableBoundingBoxColumn: new(introspect.Column),
	DetectionTableCameraIDColumn:    new(introspect.Column),
	DetectionTableEventIDColumn:     new(introspect.Column),
	DetectionTableObjectIDColumn:    new(introspect.Column),
	DetectionTableColourColumn:      new(introspect.Column),
}

var (
	DetectionTablePrimaryKeyColumn = DetectionTableIDColumn
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

func (m *Detection) GetPrimaryKeyColumn() string {
	return DetectionTablePrimaryKeyColumn
}

func (m *Detection) GetPrimaryKeyValue() any {
	return m.ID
}

func (m *Detection) FromItem(item map[string]any) error {
	if item == nil {
		return fmt.Errorf(
			"item unexpectedly nil during DetectionFromItem",
		)
	}

	if len(item) == 0 {
		return fmt.Errorf(
			"item unexpectedly empty during DetectionFromItem",
		)
	}

	wrapError := func(k string, err error) error {
		return fmt.Errorf("%#+v: %v; item: %#+v", k, err, item)
	}

	for k, v := range item {
		_, ok := DetectionTableColumnLookup[k]
		if !ok {
			return fmt.Errorf(
				"item contained unexpected key %#+v during DetectionFromItem; item: %#+v",
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

		case "score":
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

			m.Score = temp2

		case "centroid":
			if v == nil {
				continue
			}

			temp1, err := types.ParsePoint(v)
			if err != nil {
				return wrapError(k, err)
			}

			temp2, ok := temp1.(pgtype.Vec2)
			if !ok {
				if temp1 != nil {
					return wrapError(k, fmt.Errorf("failed to cast %#+v to pgtype.Vec2", temp1))
				}
			}

			m.Centroid = temp2

		case "bounding_box":
			if v == nil {
				continue
			}

			temp1, err := types.ParsePolygon(v)
			if err != nil {
				return wrapError(k, err)
			}

			temp2, ok := temp1.([]pgtype.Vec2)
			if !ok {
				if temp1 != nil {
					return wrapError(k, fmt.Errorf("failed to cast %#+v to []pgtype.Vec2", temp1))
				}
			}

			m.BoundingBox = temp2

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

		case "object_id":
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

			m.ObjectID = &temp2

		case "colour":
			if v == nil {
				continue
			}

			temp1, err := types.ParseGeometry(v)
			if err != nil {
				return wrapError(k, err)
			}

			temp2, ok := temp1.(postgis.PointZ)
			if !ok {
				if temp1 != nil {
					return wrapError(k, fmt.Errorf("failed to cast %#+v to postgis.PointZ", temp1))
				}
			}

			m.Colour = &temp2

		}
	}

	return nil
}

func (m *Detection) Reload(
	ctx context.Context,
	tx *sqlx.Tx,
) error {
	t, err := SelectDetection(
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
	m.ClassID = t.ClassID
	m.ClassName = t.ClassName
	m.Score = t.Score
	m.Centroid = t.Centroid
	m.BoundingBox = t.BoundingBox
	m.CameraID = t.CameraID
	m.CameraIDObject = t.CameraIDObject
	m.EventID = t.EventID
	m.EventIDObject = t.EventIDObject
	m.ObjectID = t.ObjectID
	m.ObjectIDObject = t.ObjectIDObject
	m.Colour = t.Colour

	return nil
}

func (m *Detection) Insert(
	ctx context.Context,
	tx *sqlx.Tx,
	setPrimaryKey bool,
	setZeroValues bool,
) error {
	columns := make([]string, 0)
	values := make([]any, 0)

	if setPrimaryKey && (setZeroValues || !types.IsZeroInt(m.ID)) {
		columns = append(columns, DetectionTableIDColumn)

		v, err := types.FormatInt(m.ID)
		if err != nil {
			return fmt.Errorf("failed to handle m.ID: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.Timestamp) {
		columns = append(columns, DetectionTableTimestampColumn)

		v, err := types.FormatTime(m.Timestamp)
		if err != nil {
			return fmt.Errorf("failed to handle m.Timestamp: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroInt(m.ClassID) {
		columns = append(columns, DetectionTableClassIDColumn)

		v, err := types.FormatInt(m.ClassID)
		if err != nil {
			return fmt.Errorf("failed to handle m.ClassID: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroString(m.ClassName) {
		columns = append(columns, DetectionTableClassNameColumn)

		v, err := types.FormatString(m.ClassName)
		if err != nil {
			return fmt.Errorf("failed to handle m.ClassName: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroFloat(m.Score) {
		columns = append(columns, DetectionTableScoreColumn)

		v, err := types.FormatFloat(m.Score)
		if err != nil {
			return fmt.Errorf("failed to handle m.Score: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroPoint(m.Centroid) {
		columns = append(columns, DetectionTableCentroidColumn)

		v, err := types.FormatPoint(m.Centroid)
		if err != nil {
			return fmt.Errorf("failed to handle m.Centroid: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroPolygon(m.BoundingBox) {
		columns = append(columns, DetectionTableBoundingBoxColumn)

		v, err := types.FormatPolygon(m.BoundingBox)
		if err != nil {
			return fmt.Errorf("failed to handle m.BoundingBox: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroInt(m.CameraID) {
		columns = append(columns, DetectionTableCameraIDColumn)

		v, err := types.FormatInt(m.CameraID)
		if err != nil {
			return fmt.Errorf("failed to handle m.CameraID: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroInt(m.EventID) {
		columns = append(columns, DetectionTableEventIDColumn)

		v, err := types.FormatInt(m.EventID)
		if err != nil {
			return fmt.Errorf("failed to handle m.EventID: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroInt(m.ObjectID) {
		columns = append(columns, DetectionTableObjectIDColumn)

		v, err := types.FormatInt(m.ObjectID)
		if err != nil {
			return fmt.Errorf("failed to handle m.ObjectID: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroGeometry(m.Colour) {
		columns = append(columns, DetectionTableColourColumn)

		v, err := types.FormatGeometry(m.Colour)
		if err != nil {
			return fmt.Errorf("failed to handle m.Colour: %v", err)
		}

		values = append(values, v)
	}

	item, err := query.Insert(
		ctx,
		tx,
		DetectionTable,
		columns,
		nil,
		false,
		false,
		DetectionTableColumns,
		values...,
	)
	if err != nil {
		return fmt.Errorf("failed to insert %#+v: %v", m, err)
	}
	v := item[DetectionTableIDColumn]

	if v == nil {
		return fmt.Errorf("failed to find %v in %#+v", DetectionTableIDColumn, item)
	}

	wrapError := func(err error) error {
		return fmt.Errorf(
			"failed to treat %v: %#+v as int64: %v",
			DetectionTableIDColumn,
			item[DetectionTableIDColumn],
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

func (m *Detection) Update(
	ctx context.Context,
	tx *sqlx.Tx,
	setZeroValues bool,
) error {
	columns := make([]string, 0)
	values := make([]any, 0)

	if setZeroValues || !types.IsZeroTime(m.Timestamp) {
		columns = append(columns, DetectionTableTimestampColumn)

		v, err := types.FormatTime(m.Timestamp)
		if err != nil {
			return fmt.Errorf("failed to handle m.Timestamp: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroInt(m.ClassID) {
		columns = append(columns, DetectionTableClassIDColumn)

		v, err := types.FormatInt(m.ClassID)
		if err != nil {
			return fmt.Errorf("failed to handle m.ClassID: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroString(m.ClassName) {
		columns = append(columns, DetectionTableClassNameColumn)

		v, err := types.FormatString(m.ClassName)
		if err != nil {
			return fmt.Errorf("failed to handle m.ClassName: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroFloat(m.Score) {
		columns = append(columns, DetectionTableScoreColumn)

		v, err := types.FormatFloat(m.Score)
		if err != nil {
			return fmt.Errorf("failed to handle m.Score: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroPoint(m.Centroid) {
		columns = append(columns, DetectionTableCentroidColumn)

		v, err := types.FormatPoint(m.Centroid)
		if err != nil {
			return fmt.Errorf("failed to handle m.Centroid: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroPolygon(m.BoundingBox) {
		columns = append(columns, DetectionTableBoundingBoxColumn)

		v, err := types.FormatPolygon(m.BoundingBox)
		if err != nil {
			return fmt.Errorf("failed to handle m.BoundingBox: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroInt(m.CameraID) {
		columns = append(columns, DetectionTableCameraIDColumn)

		v, err := types.FormatInt(m.CameraID)
		if err != nil {
			return fmt.Errorf("failed to handle m.CameraID: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroInt(m.EventID) {
		columns = append(columns, DetectionTableEventIDColumn)

		v, err := types.FormatInt(m.EventID)
		if err != nil {
			return fmt.Errorf("failed to handle m.EventID: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroInt(m.ObjectID) {
		columns = append(columns, DetectionTableObjectIDColumn)

		v, err := types.FormatInt(m.ObjectID)
		if err != nil {
			return fmt.Errorf("failed to handle m.ObjectID: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroGeometry(m.Colour) {
		columns = append(columns, DetectionTableColourColumn)

		v, err := types.FormatGeometry(m.Colour)
		if err != nil {
			return fmt.Errorf("failed to handle m.Colour: %v", err)
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
		DetectionTable,
		columns,
		fmt.Sprintf("%v = $$??", DetectionTableIDColumn),
		DetectionTableColumns,
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

func (m *Detection) Delete(
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
		DetectionTable,
		fmt.Sprintf("%v = $$??", DetectionTableIDColumn),
		values...,
	)
	if err != nil {
		return fmt.Errorf("failed to delete %#+v: %v", m, err)
	}

	return nil
}

func SelectDetections(
	ctx context.Context,
	tx *sqlx.Tx,
	where string,
	limit *int,
	offset *int,
	values ...any,
) ([]*Detection, error) {
	items, err := query.Select(
		ctx,
		tx,
		DetectionTableColumnsWithTypeCasts,
		DetectionTable,
		where,
		limit,
		offset,
		values...,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to call SelectDetections; err: %v", err)
	}

	objects := make([]*Detection, 0)

	for _, item := range items {
		object := &Detection{}

		err = object.FromItem(item)
		if err != nil {
			return nil, fmt.Errorf("failed to call Detection.FromItem; err: %v", err)
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

		if !types.IsZeroInt(object.ObjectID) {
			object.ObjectIDObject, _ = SelectObject(
				ctx,
				tx,
				fmt.Sprintf("%v = $1", ObjectTablePrimaryKeyColumn),
				object.ObjectID,
			)
		}

		objects = append(objects, object)
	}

	return objects, nil
}

func SelectDetection(
	ctx context.Context,
	tx *sqlx.Tx,
	where string,
	values ...any,
) (*Detection, error) {
	objects, err := SelectDetections(
		ctx,
		tx,
		where,
		helpers.Ptr(2),
		helpers.Ptr(0),
		values...,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to call SelectDetection; err: %v", err)
	}

	if len(objects) > 1 {
		return nil, fmt.Errorf("attempt to call SelectDetection returned more than 1 row")
	}

	if len(objects) < 1 {
		return nil, fmt.Errorf("attempt to call SelectDetection returned no rows")
	}

	object := objects[0]

	return object, nil
}

func handleGetDetections(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
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
			column := DetectionTableColumnLookup[parts[0]]
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

	objects, err := SelectDetections(ctx, tx, where, &limit, &offset, values...)
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

func handleGetDetection(w http.ResponseWriter, r *http.Request, db *sqlx.DB, primaryKey string) {
	ctx := r.Context()

	where := fmt.Sprintf("%s = $$??", DetectionTablePrimaryKeyColumn)

	values := []any{primaryKey}

	tx, err := db.BeginTxx(r.Context(), nil)
	if err != nil {
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	object, err := SelectDetection(ctx, tx, where, values...)
	if err != nil {
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	err = tx.Commit()
	if err != nil {
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	helpers.HandleObjectsResponse(w, http.StatusOK, []*Detection{object})
}

func handlePostDetections(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
}

func handlePutDetection(w http.ResponseWriter, r *http.Request, db *sqlx.DB, primaryKey string) {
}

func handlePatchDetection(w http.ResponseWriter, r *http.Request, db *sqlx.DB, primaryKey string) {
}

func GetDetectionRouter(db *sqlx.DB, middlewares ...func(http.Handler) http.Handler) chi.Router {
	r := chi.NewRouter()

	for _, m := range middlewares {
		r.Use(m)
	}

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		handleGetDetections(w, r, db)
	})

	r.Get("/{primaryKey}", func(w http.ResponseWriter, r *http.Request) {
		handleGetDetection(w, r, db, chi.URLParam(r, "primaryKey"))
	})

	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		handlePostDetections(w, r, db)
	})

	r.Put("/{primaryKey}", func(w http.ResponseWriter, r *http.Request) {
		handlePutDetection(w, r, db, chi.URLParam(r, "primaryKey"))
	})

	r.Patch("/{primaryKey}", func(w http.ResponseWriter, r *http.Request) {
		handlePatchDetection(w, r, db, chi.URLParam(r, "primaryKey"))
	})

	return r
}

func GetDetectionHandlerFunc(db *sqlx.DB, middlewares ...func(http.Handler) http.Handler) http.HandlerFunc {
	r := chi.NewRouter()

	r.Mount("/detections", GetDetectionRouter(db, middlewares...))

	return r.ServeHTTP
}

func NewDetectionFromItem(item map[string]any) (any, error) {
	object := &Detection{}

	err := object.FromItem(item)
	if err != nil {
		return nil, err
	}

	return object, nil
}

func init() {
	register(
		DetectionTable,
		NewDetectionFromItem,
		"/detections",
		GetDetectionRouter,
	)
}
