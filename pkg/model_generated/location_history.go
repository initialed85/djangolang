package model_generated

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/netip"
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

type LocationHistory struct {
	ID                          uuid.UUID      `json:"id"`
	CreatedAt                   time.Time      `json:"created_at"`
	UpdatedAt                   time.Time      `json:"updated_at"`
	DeletedAt                   *time.Time     `json:"deleted_at"`
	Timestamp                   time.Time      `json:"timestamp"`
	Point                       *pgtype.Vec2   `json:"point"`
	Polygon                     *[]pgtype.Vec2 `json:"polygon"`
	ParentPhysicalThingID       *uuid.UUID     `json:"parent_physical_thing_id"`
	ParentPhysicalThingIDObject *PhysicalThing `json:"parent_physical_thing_id_object"`
}

var LocationHistoryTable = "location_history"

var (
	LocationHistoryTableIDColumn                    = "id"
	LocationHistoryTableCreatedAtColumn             = "created_at"
	LocationHistoryTableUpdatedAtColumn             = "updated_at"
	LocationHistoryTableDeletedAtColumn             = "deleted_at"
	LocationHistoryTableTimestampColumn             = "timestamp"
	LocationHistoryTablePointColumn                 = "point"
	LocationHistoryTablePolygonColumn               = "polygon"
	LocationHistoryTableParentPhysicalThingIDColumn = "parent_physical_thing_id"
)

var (
	LocationHistoryTableIDColumnWithTypeCast                    = fmt.Sprintf(`"id" AS id`)
	LocationHistoryTableCreatedAtColumnWithTypeCast             = fmt.Sprintf(`"created_at" AS created_at`)
	LocationHistoryTableUpdatedAtColumnWithTypeCast             = fmt.Sprintf(`"updated_at" AS updated_at`)
	LocationHistoryTableDeletedAtColumnWithTypeCast             = fmt.Sprintf(`"deleted_at" AS deleted_at`)
	LocationHistoryTableTimestampColumnWithTypeCast             = fmt.Sprintf(`"timestamp" AS timestamp`)
	LocationHistoryTablePointColumnWithTypeCast                 = fmt.Sprintf(`"point" AS point`)
	LocationHistoryTablePolygonColumnWithTypeCast               = fmt.Sprintf(`"polygon" AS polygon`)
	LocationHistoryTableParentPhysicalThingIDColumnWithTypeCast = fmt.Sprintf(`"parent_physical_thing_id" AS parent_physical_thing_id`)
)

var LocationHistoryTableColumns = []string{
	LocationHistoryTableIDColumn,
	LocationHistoryTableCreatedAtColumn,
	LocationHistoryTableUpdatedAtColumn,
	LocationHistoryTableDeletedAtColumn,
	LocationHistoryTableTimestampColumn,
	LocationHistoryTablePointColumn,
	LocationHistoryTablePolygonColumn,
	LocationHistoryTableParentPhysicalThingIDColumn,
}

var LocationHistoryTableColumnsWithTypeCasts = []string{
	LocationHistoryTableIDColumnWithTypeCast,
	LocationHistoryTableCreatedAtColumnWithTypeCast,
	LocationHistoryTableUpdatedAtColumnWithTypeCast,
	LocationHistoryTableDeletedAtColumnWithTypeCast,
	LocationHistoryTableTimestampColumnWithTypeCast,
	LocationHistoryTablePointColumnWithTypeCast,
	LocationHistoryTablePolygonColumnWithTypeCast,
	LocationHistoryTableParentPhysicalThingIDColumnWithTypeCast,
}

var LocationHistoryTableColumnLookup = map[string]*introspect.Column{
	LocationHistoryTableIDColumn:                    new(introspect.Column),
	LocationHistoryTableCreatedAtColumn:             new(introspect.Column),
	LocationHistoryTableUpdatedAtColumn:             new(introspect.Column),
	LocationHistoryTableDeletedAtColumn:             new(introspect.Column),
	LocationHistoryTableTimestampColumn:             new(introspect.Column),
	LocationHistoryTablePointColumn:                 new(introspect.Column),
	LocationHistoryTablePolygonColumn:               new(introspect.Column),
	LocationHistoryTableParentPhysicalThingIDColumn: new(introspect.Column),
}

var (
	LocationHistoryTablePrimaryKeyColumn = LocationHistoryTableIDColumn
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

func (m *LocationHistory) GetPrimaryKeyColumn() string {
	return LocationHistoryTablePrimaryKeyColumn
}

func (m *LocationHistory) GetPrimaryKeyValue() any {
	return m.ID
}

func (m *LocationHistory) FromItem(item map[string]any) error {
	if item == nil {
		return fmt.Errorf(
			"item unexpectedly nil during LocationHistoryFromItem",
		)
	}

	if len(item) == 0 {
		return fmt.Errorf(
			"item unexpectedly empty during LocationHistoryFromItem",
		)
	}

	wrapError := func(k string, err error) error {
		return fmt.Errorf("%#+v: %v; item: %#+v", k, err, item)
	}

	for k, v := range item {
		_, ok := LocationHistoryTableColumnLookup[k]
		if !ok {
			return fmt.Errorf(
				"item contained unexpected key %#+v during LocationHistoryFromItem; item: %#+v",
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
				return wrapError(k, err)
			}

			temp2, ok := temp1.(uuid.UUID)
			if !ok {
				if temp1 != nil {
					return wrapError(k, fmt.Errorf("failed to cast %#+v to uuid.UUID", temp1))
				}
			}

			m.ID = temp2

		case "created_at":
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

			m.CreatedAt = temp2

		case "updated_at":
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

			m.UpdatedAt = temp2

		case "deleted_at":
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

			m.DeletedAt = &temp2

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

		case "point":
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

			m.Point = &temp2

		case "polygon":
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

			m.Polygon = &temp2

		case "parent_physical_thing_id":
			if v == nil {
				continue
			}

			temp1, err := types.ParseUUID(v)
			if err != nil {
				return wrapError(k, err)
			}

			temp2, ok := temp1.(uuid.UUID)
			if !ok {
				if temp1 != nil {
					return wrapError(k, fmt.Errorf("failed to cast %#+v to uuid.UUID", temp1))
				}
			}

			m.ParentPhysicalThingID = &temp2

		}
	}

	return nil
}

func (m *LocationHistory) Reload(
	ctx context.Context,
	tx *sqlx.Tx,
) error {
	t, err := SelectLocationHistory(
		ctx,
		tx,
		fmt.Sprintf("%v = $1", m.GetPrimaryKeyColumn()),
		m.GetPrimaryKeyValue(),
	)
	if err != nil {
		return err
	}

	m.ID = t.ID
	m.CreatedAt = t.CreatedAt
	m.UpdatedAt = t.UpdatedAt
	m.DeletedAt = t.DeletedAt
	m.Timestamp = t.Timestamp
	m.Point = t.Point
	m.Polygon = t.Polygon
	m.ParentPhysicalThingID = t.ParentPhysicalThingID
	m.ParentPhysicalThingIDObject = t.ParentPhysicalThingIDObject

	return nil
}

func (m *LocationHistory) Insert(
	ctx context.Context,
	tx *sqlx.Tx,
	setPrimaryKey bool,
	setZeroValues bool,
) error {
	columns := make([]string, 0)
	values := make([]any, 0)

	if setPrimaryKey && (setZeroValues || !types.IsZeroUUID(m.ID)) {
		columns = append(columns, LocationHistoryTableIDColumn)

		v, err := types.FormatUUID(m.ID)
		if err != nil {
			return fmt.Errorf("failed to handle m.ID: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.CreatedAt) {
		columns = append(columns, LocationHistoryTableCreatedAtColumn)

		v, err := types.FormatTime(m.CreatedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.CreatedAt: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.UpdatedAt) {
		columns = append(columns, LocationHistoryTableUpdatedAtColumn)

		v, err := types.FormatTime(m.UpdatedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.UpdatedAt: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.DeletedAt) {
		columns = append(columns, LocationHistoryTableDeletedAtColumn)

		v, err := types.FormatTime(m.DeletedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.DeletedAt: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.Timestamp) {
		columns = append(columns, LocationHistoryTableTimestampColumn)

		v, err := types.FormatTime(m.Timestamp)
		if err != nil {
			return fmt.Errorf("failed to handle m.Timestamp: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroPoint(m.Point) {
		columns = append(columns, LocationHistoryTablePointColumn)

		v, err := types.FormatPoint(m.Point)
		if err != nil {
			return fmt.Errorf("failed to handle m.Point: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroPolygon(m.Polygon) {
		columns = append(columns, LocationHistoryTablePolygonColumn)

		v, err := types.FormatPolygon(m.Polygon)
		if err != nil {
			return fmt.Errorf("failed to handle m.Polygon: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroUUID(m.ParentPhysicalThingID) {
		columns = append(columns, LocationHistoryTableParentPhysicalThingIDColumn)

		v, err := types.FormatUUID(m.ParentPhysicalThingID)
		if err != nil {
			return fmt.Errorf("failed to handle m.ParentPhysicalThingID: %v", err)
		}

		values = append(values, v)
	}

	item, err := query.Insert(
		ctx,
		tx,
		LocationHistoryTable,
		columns,
		nil,
		false,
		false,
		LocationHistoryTableColumns,
		values...,
	)
	if err != nil {
		return fmt.Errorf("failed to insert %#+v: %v", m, err)
	}
	v := item[LocationHistoryTableIDColumn]

	if v == nil {
		return fmt.Errorf("failed to find %v in %#+v", LocationHistoryTableIDColumn, item)
	}

	wrapError := func(err error) error {
		return fmt.Errorf(
			"failed to treat %v: %#+v as uuid.UUID: %v",
			LocationHistoryTableIDColumn,
			item[LocationHistoryTableIDColumn],
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

	err = m.Reload(ctx, tx)
	if err != nil {
		return fmt.Errorf("failed to reload after insert")
	}

	return nil
}

func (m *LocationHistory) Update(
	ctx context.Context,
	tx *sqlx.Tx,
	setZeroValues bool,
) error {
	columns := make([]string, 0)
	values := make([]any, 0)

	if setZeroValues || !types.IsZeroTime(m.CreatedAt) {
		columns = append(columns, LocationHistoryTableCreatedAtColumn)

		v, err := types.FormatTime(m.CreatedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.CreatedAt: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.UpdatedAt) {
		columns = append(columns, LocationHistoryTableUpdatedAtColumn)

		v, err := types.FormatTime(m.UpdatedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.UpdatedAt: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.DeletedAt) {
		columns = append(columns, LocationHistoryTableDeletedAtColumn)

		v, err := types.FormatTime(m.DeletedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.DeletedAt: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.Timestamp) {
		columns = append(columns, LocationHistoryTableTimestampColumn)

		v, err := types.FormatTime(m.Timestamp)
		if err != nil {
			return fmt.Errorf("failed to handle m.Timestamp: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroPoint(m.Point) {
		columns = append(columns, LocationHistoryTablePointColumn)

		v, err := types.FormatPoint(m.Point)
		if err != nil {
			return fmt.Errorf("failed to handle m.Point: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroPolygon(m.Polygon) {
		columns = append(columns, LocationHistoryTablePolygonColumn)

		v, err := types.FormatPolygon(m.Polygon)
		if err != nil {
			return fmt.Errorf("failed to handle m.Polygon: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroUUID(m.ParentPhysicalThingID) {
		columns = append(columns, LocationHistoryTableParentPhysicalThingIDColumn)

		v, err := types.FormatUUID(m.ParentPhysicalThingID)
		if err != nil {
			return fmt.Errorf("failed to handle m.ParentPhysicalThingID: %v", err)
		}

		values = append(values, v)
	}

	v, err := types.FormatUUID(m.ID)
	if err != nil {
		return fmt.Errorf("failed to handle m.ID: %v", err)
	}

	values = append(values, v)

	_, err = query.Update(
		ctx,
		tx,
		LocationHistoryTable,
		columns,
		fmt.Sprintf("%v = $$??", LocationHistoryTableIDColumn),
		LocationHistoryTableColumns,
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

func (m *LocationHistory) Delete(
	ctx context.Context,
	tx *sqlx.Tx,
) error {
	values := make([]any, 0)
	v, err := types.FormatUUID(m.ID)
	if err != nil {
		return fmt.Errorf("failed to handle m.ID: %v", err)
	}

	values = append(values, v)

	err = query.Delete(
		ctx,
		tx,
		LocationHistoryTable,
		fmt.Sprintf("%v = $$??", LocationHistoryTableIDColumn),
		values...,
	)
	if err != nil {
		return fmt.Errorf("failed to delete %#+v: %v", m, err)
	}

	return nil
}

func SelectLocationHistorys(
	ctx context.Context,
	tx *sqlx.Tx,
	where string,
	limit *int,
	offset *int,
	values ...any,
) ([]*LocationHistory, error) {
	items, err := query.Select(
		ctx,
		tx,
		LocationHistoryTableColumnsWithTypeCasts,
		LocationHistoryTable,
		where,
		limit,
		offset,
		values...,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to call SelectLocationHistorys; err: %v", err)
	}

	objects := make([]*LocationHistory, 0)

	for _, item := range items {
		object := &LocationHistory{}

		err = object.FromItem(item)
		if err != nil {
			return nil, fmt.Errorf("failed to call LocationHistory.FromItem; err: %v", err)
		}

		if !types.IsZeroUUID(object.ParentPhysicalThingID) {
			object.ParentPhysicalThingIDObject, _ = SelectPhysicalThing(
				ctx,
				tx,
				fmt.Sprintf("%v = $1", PhysicalThingTablePrimaryKeyColumn),
				object.ParentPhysicalThingID,
			)
		}

		objects = append(objects, object)
	}

	return objects, nil
}

func SelectLocationHistory(
	ctx context.Context,
	tx *sqlx.Tx,
	where string,
	values ...any,
) (*LocationHistory, error) {
	objects, err := SelectLocationHistorys(
		ctx,
		tx,
		where,
		helpers.Ptr(2),
		helpers.Ptr(0),
		values...,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to call SelectLocationHistory; err: %v", err)
	}

	if len(objects) > 1 {
		return nil, fmt.Errorf("attempt to call SelectLocationHistory returned more than 1 row")
	}

	if len(objects) < 1 {
		return nil, fmt.Errorf("attempt to call SelectLocationHistory returned no rows")
	}

	object := objects[0]

	return object, nil
}

func handleGetLocationHistorys(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
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
			column := LocationHistoryTableColumnLookup[parts[0]]
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

	objects, err := SelectLocationHistorys(ctx, tx, where, &limit, &offset, values...)
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

func handleGetLocationHistory(w http.ResponseWriter, r *http.Request, db *sqlx.DB, primaryKey string) {
	ctx := r.Context()

	where := fmt.Sprintf("%s = $$??", LocationHistoryTablePrimaryKeyColumn)

	values := []any{primaryKey}

	tx, err := db.BeginTxx(r.Context(), nil)
	if err != nil {
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	object, err := SelectLocationHistory(ctx, tx, where, values...)
	if err != nil {
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	err = tx.Commit()
	if err != nil {
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	helpers.HandleObjectsResponse(w, http.StatusOK, []*LocationHistory{object})
}

func handlePostLocationHistorys(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
}

func handlePutLocationHistory(w http.ResponseWriter, r *http.Request, db *sqlx.DB, primaryKey string) {
}

func handlePatchLocationHistory(w http.ResponseWriter, r *http.Request, db *sqlx.DB, primaryKey string) {
}

func GetLocationHistoryRouter(db *sqlx.DB, middlewares ...func(http.Handler) http.Handler) chi.Router {
	r := chi.NewRouter()

	for _, m := range middlewares {
		r.Use(m)
	}

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		handleGetLocationHistorys(w, r, db)
	})

	r.Get("/{primaryKey}", func(w http.ResponseWriter, r *http.Request) {
		handleGetLocationHistory(w, r, db, chi.URLParam(r, "primaryKey"))
	})

	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		handlePostLocationHistorys(w, r, db)
	})

	r.Put("/{primaryKey}", func(w http.ResponseWriter, r *http.Request) {
		handlePutLocationHistory(w, r, db, chi.URLParam(r, "primaryKey"))
	})

	r.Patch("/{primaryKey}", func(w http.ResponseWriter, r *http.Request) {
		handlePatchLocationHistory(w, r, db, chi.URLParam(r, "primaryKey"))
	})

	return r
}

func GetLocationHistoryHandlerFunc(db *sqlx.DB, middlewares ...func(http.Handler) http.Handler) http.HandlerFunc {
	r := chi.NewRouter()

	r.Mount("/location-histories", GetLocationHistoryRouter(db, middlewares...))

	return r.ServeHTTP
}

func NewLocationHistoryFromItem(item map[string]any) (any, error) {
	object := &LocationHistory{}

	err := object.FromItem(item)
	if err != nil {
		return nil, err
	}

	return object, nil
}

func init() {
	register(
		LocationHistoryTable,
		NewLocationHistoryFromItem,
		"/location-histories",
		GetLocationHistoryRouter,
	)
}
