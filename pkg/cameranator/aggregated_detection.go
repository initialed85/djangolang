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

type AggregatedDetection struct {
	ID             int64     `json:"id"`
	StartTimestamp time.Time `json:"start_timestamp"`
	EndTimestamp   time.Time `json:"end_timestamp"`
	ClassID        int64     `json:"class_id"`
	ClassName      string    `json:"class_name"`
	Score          float64   `json:"score"`
	Count          int64     `json:"count"`
	WeightedScore  float64   `json:"weighted_score"`
	EventID        *int64    `json:"event_id"`
	EventIDObject  *Event    `json:"event_id_object"`
}

var AggregatedDetectionTable = "aggregated_detection"

var (
	AggregatedDetectionTableIDColumn             = "id"
	AggregatedDetectionTableStartTimestampColumn = "start_timestamp"
	AggregatedDetectionTableEndTimestampColumn   = "end_timestamp"
	AggregatedDetectionTableClassIDColumn        = "class_id"
	AggregatedDetectionTableClassNameColumn      = "class_name"
	AggregatedDetectionTableScoreColumn          = "score"
	AggregatedDetectionTableCountColumn          = "count"
	AggregatedDetectionTableWeightedScoreColumn  = "weighted_score"
	AggregatedDetectionTableEventIDColumn        = "event_id"
)

var (
	AggregatedDetectionTableIDColumnWithTypeCast             = fmt.Sprintf(`"id" AS id`)
	AggregatedDetectionTableStartTimestampColumnWithTypeCast = fmt.Sprintf(`"start_timestamp" AS start_timestamp`)
	AggregatedDetectionTableEndTimestampColumnWithTypeCast   = fmt.Sprintf(`"end_timestamp" AS end_timestamp`)
	AggregatedDetectionTableClassIDColumnWithTypeCast        = fmt.Sprintf(`"class_id" AS class_id`)
	AggregatedDetectionTableClassNameColumnWithTypeCast      = fmt.Sprintf(`"class_name" AS class_name`)
	AggregatedDetectionTableScoreColumnWithTypeCast          = fmt.Sprintf(`"score" AS score`)
	AggregatedDetectionTableCountColumnWithTypeCast          = fmt.Sprintf(`"count" AS count`)
	AggregatedDetectionTableWeightedScoreColumnWithTypeCast  = fmt.Sprintf(`"weighted_score" AS weighted_score`)
	AggregatedDetectionTableEventIDColumnWithTypeCast        = fmt.Sprintf(`"event_id" AS event_id`)
)

var AggregatedDetectionTableColumns = []string{
	AggregatedDetectionTableIDColumn,
	AggregatedDetectionTableStartTimestampColumn,
	AggregatedDetectionTableEndTimestampColumn,
	AggregatedDetectionTableClassIDColumn,
	AggregatedDetectionTableClassNameColumn,
	AggregatedDetectionTableScoreColumn,
	AggregatedDetectionTableCountColumn,
	AggregatedDetectionTableWeightedScoreColumn,
	AggregatedDetectionTableEventIDColumn,
}

var AggregatedDetectionTableColumnsWithTypeCasts = []string{
	AggregatedDetectionTableIDColumnWithTypeCast,
	AggregatedDetectionTableStartTimestampColumnWithTypeCast,
	AggregatedDetectionTableEndTimestampColumnWithTypeCast,
	AggregatedDetectionTableClassIDColumnWithTypeCast,
	AggregatedDetectionTableClassNameColumnWithTypeCast,
	AggregatedDetectionTableScoreColumnWithTypeCast,
	AggregatedDetectionTableCountColumnWithTypeCast,
	AggregatedDetectionTableWeightedScoreColumnWithTypeCast,
	AggregatedDetectionTableEventIDColumnWithTypeCast,
}

var AggregatedDetectionTableColumnLookup = map[string]*introspect.Column{
	AggregatedDetectionTableIDColumn:             new(introspect.Column),
	AggregatedDetectionTableStartTimestampColumn: new(introspect.Column),
	AggregatedDetectionTableEndTimestampColumn:   new(introspect.Column),
	AggregatedDetectionTableClassIDColumn:        new(introspect.Column),
	AggregatedDetectionTableClassNameColumn:      new(introspect.Column),
	AggregatedDetectionTableScoreColumn:          new(introspect.Column),
	AggregatedDetectionTableCountColumn:          new(introspect.Column),
	AggregatedDetectionTableWeightedScoreColumn:  new(introspect.Column),
	AggregatedDetectionTableEventIDColumn:        new(introspect.Column),
}

var (
	AggregatedDetectionTablePrimaryKeyColumn = AggregatedDetectionTableIDColumn
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

func (m *AggregatedDetection) GetPrimaryKeyColumn() string {
	return AggregatedDetectionTablePrimaryKeyColumn
}

func (m *AggregatedDetection) GetPrimaryKeyValue() any {
	return m.ID
}

func (m *AggregatedDetection) FromItem(item map[string]any) error {
	if item == nil {
		return fmt.Errorf(
			"item unexpectedly nil during AggregatedDetectionFromItem",
		)
	}

	if len(item) == 0 {
		return fmt.Errorf(
			"item unexpectedly empty during AggregatedDetectionFromItem",
		)
	}

	wrapError := func(k string, err error) error {
		return fmt.Errorf("%#+v: %v; item: %#+v", k, err, item)
	}

	for k, v := range item {
		_, ok := AggregatedDetectionTableColumnLookup[k]
		if !ok {
			return fmt.Errorf(
				"item contained unexpected key %#+v during AggregatedDetectionFromItem; item: %#+v",
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

		case "count":
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

			m.Count = temp2

		case "weighted_score":
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

			m.WeightedScore = temp2

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

func (m *AggregatedDetection) Reload(
	ctx context.Context,
	tx *sqlx.Tx,
) error {
	t, err := SelectAggregatedDetection(
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
	m.Score = t.Score
	m.Count = t.Count
	m.WeightedScore = t.WeightedScore
	m.EventID = t.EventID
	m.EventIDObject = t.EventIDObject

	return nil
}

func (m *AggregatedDetection) Insert(
	ctx context.Context,
	tx *sqlx.Tx,
	setPrimaryKey bool,
	setZeroValues bool,
) error {
	columns := make([]string, 0)
	values := make([]any, 0)

	if setPrimaryKey && (setZeroValues || !types.IsZeroInt(m.ID)) {
		columns = append(columns, AggregatedDetectionTableIDColumn)

		v, err := types.FormatInt(m.ID)
		if err != nil {
			return fmt.Errorf("failed to handle m.ID: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.StartTimestamp) {
		columns = append(columns, AggregatedDetectionTableStartTimestampColumn)

		v, err := types.FormatTime(m.StartTimestamp)
		if err != nil {
			return fmt.Errorf("failed to handle m.StartTimestamp: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.EndTimestamp) {
		columns = append(columns, AggregatedDetectionTableEndTimestampColumn)

		v, err := types.FormatTime(m.EndTimestamp)
		if err != nil {
			return fmt.Errorf("failed to handle m.EndTimestamp: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroInt(m.ClassID) {
		columns = append(columns, AggregatedDetectionTableClassIDColumn)

		v, err := types.FormatInt(m.ClassID)
		if err != nil {
			return fmt.Errorf("failed to handle m.ClassID: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroString(m.ClassName) {
		columns = append(columns, AggregatedDetectionTableClassNameColumn)

		v, err := types.FormatString(m.ClassName)
		if err != nil {
			return fmt.Errorf("failed to handle m.ClassName: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroFloat(m.Score) {
		columns = append(columns, AggregatedDetectionTableScoreColumn)

		v, err := types.FormatFloat(m.Score)
		if err != nil {
			return fmt.Errorf("failed to handle m.Score: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroInt(m.Count) {
		columns = append(columns, AggregatedDetectionTableCountColumn)

		v, err := types.FormatInt(m.Count)
		if err != nil {
			return fmt.Errorf("failed to handle m.Count: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroFloat(m.WeightedScore) {
		columns = append(columns, AggregatedDetectionTableWeightedScoreColumn)

		v, err := types.FormatFloat(m.WeightedScore)
		if err != nil {
			return fmt.Errorf("failed to handle m.WeightedScore: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroInt(m.EventID) {
		columns = append(columns, AggregatedDetectionTableEventIDColumn)

		v, err := types.FormatInt(m.EventID)
		if err != nil {
			return fmt.Errorf("failed to handle m.EventID: %v", err)
		}

		values = append(values, v)
	}

	item, err := query.Insert(
		ctx,
		tx,
		AggregatedDetectionTable,
		columns,
		nil,
		false,
		false,
		AggregatedDetectionTableColumns,
		values...,
	)
	if err != nil {
		return fmt.Errorf("failed to insert %#+v: %v", m, err)
	}
	v := item[AggregatedDetectionTableIDColumn]

	if v == nil {
		return fmt.Errorf("failed to find %v in %#+v", AggregatedDetectionTableIDColumn, item)
	}

	wrapError := func(err error) error {
		return fmt.Errorf(
			"failed to treat %v: %#+v as int64: %v",
			AggregatedDetectionTableIDColumn,
			item[AggregatedDetectionTableIDColumn],
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

func (m *AggregatedDetection) Update(
	ctx context.Context,
	tx *sqlx.Tx,
	setZeroValues bool,
) error {
	columns := make([]string, 0)
	values := make([]any, 0)

	if setZeroValues || !types.IsZeroTime(m.StartTimestamp) {
		columns = append(columns, AggregatedDetectionTableStartTimestampColumn)

		v, err := types.FormatTime(m.StartTimestamp)
		if err != nil {
			return fmt.Errorf("failed to handle m.StartTimestamp: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.EndTimestamp) {
		columns = append(columns, AggregatedDetectionTableEndTimestampColumn)

		v, err := types.FormatTime(m.EndTimestamp)
		if err != nil {
			return fmt.Errorf("failed to handle m.EndTimestamp: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroInt(m.ClassID) {
		columns = append(columns, AggregatedDetectionTableClassIDColumn)

		v, err := types.FormatInt(m.ClassID)
		if err != nil {
			return fmt.Errorf("failed to handle m.ClassID: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroString(m.ClassName) {
		columns = append(columns, AggregatedDetectionTableClassNameColumn)

		v, err := types.FormatString(m.ClassName)
		if err != nil {
			return fmt.Errorf("failed to handle m.ClassName: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroFloat(m.Score) {
		columns = append(columns, AggregatedDetectionTableScoreColumn)

		v, err := types.FormatFloat(m.Score)
		if err != nil {
			return fmt.Errorf("failed to handle m.Score: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroInt(m.Count) {
		columns = append(columns, AggregatedDetectionTableCountColumn)

		v, err := types.FormatInt(m.Count)
		if err != nil {
			return fmt.Errorf("failed to handle m.Count: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroFloat(m.WeightedScore) {
		columns = append(columns, AggregatedDetectionTableWeightedScoreColumn)

		v, err := types.FormatFloat(m.WeightedScore)
		if err != nil {
			return fmt.Errorf("failed to handle m.WeightedScore: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroInt(m.EventID) {
		columns = append(columns, AggregatedDetectionTableEventIDColumn)

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
		AggregatedDetectionTable,
		columns,
		fmt.Sprintf("%v = $$??", AggregatedDetectionTableIDColumn),
		AggregatedDetectionTableColumns,
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

func (m *AggregatedDetection) Delete(
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
		AggregatedDetectionTable,
		fmt.Sprintf("%v = $$??", AggregatedDetectionTableIDColumn),
		values...,
	)
	if err != nil {
		return fmt.Errorf("failed to delete %#+v: %v", m, err)
	}

	return nil
}

func SelectAggregatedDetections(
	ctx context.Context,
	tx *sqlx.Tx,
	where string,
	limit *int,
	offset *int,
	values ...any,
) ([]*AggregatedDetection, error) {
	items, err := query.Select(
		ctx,
		tx,
		AggregatedDetectionTableColumnsWithTypeCasts,
		AggregatedDetectionTable,
		where,
		limit,
		offset,
		values...,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to call SelectAggregatedDetections; err: %v", err)
	}

	objects := make([]*AggregatedDetection, 0)

	for _, item := range items {
		object := &AggregatedDetection{}

		err = object.FromItem(item)
		if err != nil {
			return nil, fmt.Errorf("failed to call AggregatedDetection.FromItem; err: %v", err)
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

func SelectAggregatedDetection(
	ctx context.Context,
	tx *sqlx.Tx,
	where string,
	values ...any,
) (*AggregatedDetection, error) {
	objects, err := SelectAggregatedDetections(
		ctx,
		tx,
		where,
		helpers.Ptr(2),
		helpers.Ptr(0),
		values...,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to call SelectAggregatedDetection; err: %v", err)
	}

	if len(objects) > 1 {
		return nil, fmt.Errorf("attempt to call SelectAggregatedDetection returned more than 1 row")
	}

	if len(objects) < 1 {
		return nil, fmt.Errorf("attempt to call SelectAggregatedDetection returned no rows")
	}

	object := objects[0]

	return object, nil
}

func handleGetAggregatedDetections(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
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
			column := AggregatedDetectionTableColumnLookup[parts[0]]
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

	objects, err := SelectAggregatedDetections(ctx, tx, where, &limit, &offset, values...)
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

func handleGetAggregatedDetection(w http.ResponseWriter, r *http.Request, db *sqlx.DB, primaryKey string) {
	ctx := r.Context()

	where := fmt.Sprintf("%s = $$??", AggregatedDetectionTablePrimaryKeyColumn)

	values := []any{primaryKey}

	tx, err := db.BeginTxx(r.Context(), nil)
	if err != nil {
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	object, err := SelectAggregatedDetection(ctx, tx, where, values...)
	if err != nil {
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	err = tx.Commit()
	if err != nil {
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	helpers.HandleObjectsResponse(w, http.StatusOK, []*AggregatedDetection{object})
}

func handlePostAggregatedDetections(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
}

func handlePutAggregatedDetection(w http.ResponseWriter, r *http.Request, db *sqlx.DB, primaryKey string) {
}

func handlePatchAggregatedDetection(w http.ResponseWriter, r *http.Request, db *sqlx.DB, primaryKey string) {
}

func GetAggregatedDetectionRouter(db *sqlx.DB, middlewares ...func(http.Handler) http.Handler) chi.Router {
	r := chi.NewRouter()

	for _, m := range middlewares {
		r.Use(m)
	}

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		handleGetAggregatedDetections(w, r, db)
	})

	r.Get("/{primaryKey}", func(w http.ResponseWriter, r *http.Request) {
		handleGetAggregatedDetection(w, r, db, chi.URLParam(r, "primaryKey"))
	})

	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		handlePostAggregatedDetections(w, r, db)
	})

	r.Put("/{primaryKey}", func(w http.ResponseWriter, r *http.Request) {
		handlePutAggregatedDetection(w, r, db, chi.URLParam(r, "primaryKey"))
	})

	r.Patch("/{primaryKey}", func(w http.ResponseWriter, r *http.Request) {
		handlePatchAggregatedDetection(w, r, db, chi.URLParam(r, "primaryKey"))
	})

	return r
}

func GetAggregatedDetectionHandlerFunc(db *sqlx.DB, middlewares ...func(http.Handler) http.Handler) http.HandlerFunc {
	r := chi.NewRouter()

	r.Mount("/aggregated-detections", GetAggregatedDetectionRouter(db, middlewares...))

	return r.ServeHTTP
}

func NewAggregatedDetectionFromItem(item map[string]any) (any, error) {
	object := &AggregatedDetection{}

	err := object.FromItem(item)
	if err != nil {
		return nil, err
	}

	return object, nil
}

func init() {
	register(
		AggregatedDetectionTable,
		NewAggregatedDetectionFromItem,
		"/aggregated-detections",
		GetAggregatedDetectionRouter,
	)
}
