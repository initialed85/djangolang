package model_generated

import (
	"context"
	"fmt"
	"time"

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
	ID                          uuid.UUID       `json:"id"`
	CreatedAt                   time.Time       `json:"created_at"`
	UpdatedAt                   time.Time       `json:"updated_at"`
	DeletedAt                   *time.Time      `json:"deleted_at"`
	Timestamp                   time.Time       `json:"timestamp"`
	Point                       *pgtype.Point   `json:"point"`
	Polygon                     *pgtype.Polygon `json:"polygon"`
	ParentPhysicalThingID       *uuid.UUID      `json:"parent_physical_thing_id"`
	ParentPhysicalThingIDObject *PhysicalThing  `json:"-"`
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
				return wrapError(k, fmt.Errorf("failed to cast to uuid.UUID"))
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
				return wrapError(k, fmt.Errorf("failed to cast to time.Time"))
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
				return wrapError(k, fmt.Errorf("failed to cast to time.Time"))
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
				return wrapError(k, fmt.Errorf("failed to cast to time.Time"))
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
				return wrapError(k, fmt.Errorf("failed to cast to time.Time"))
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

			temp2, ok := temp1.(pgtype.Point)
			if !ok {
				return wrapError(k, fmt.Errorf("failed to cast to pgtype.Point"))
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

			temp2, ok := temp1.(pgtype.Polygon)
			if !ok {
				return wrapError(k, fmt.Errorf("failed to cast to pgtype.Polygon"))
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
				return wrapError(k, fmt.Errorf("failed to cast to uuid.UUID"))
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

		if object.ParentPhysicalThingID != nil {
			object.ParentPhysicalThingIDObject, err = SelectPhysicalThing(
				ctx,
				tx,
				fmt.Sprintf("%v = $1", PhysicalThingTablePrimaryKeyColumn),
				object.ParentPhysicalThingID,
			)
			if err != nil {
				return nil, fmt.Errorf("failed to load <no value>.ParentPhysicalThingIDObject; err: %v", err)
			}
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

func (l *LocationHistory) Insert(ctx context.Context, db *sqlx.DB, columns ...string) error {
	return nil
}
