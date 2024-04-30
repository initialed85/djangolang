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
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/lib/pq/hstore"
	"github.com/twpayne/go-geom"
)

type LocationHistory struct {
	ID                          uuid.UUID      `json:"id"`
	CreatedAt                   time.Time      `json:"created_at"`
	UpdatedAt                   time.Time      `json:"updated_at"`
	DeletedAt                   *time.Time     `json:"deleted_at"`
	Timestamp                   time.Time      `json:"timestamp"`
	Point                       *geom.Point    `json:"point"`
	Polygon                     *geom.Polygon  `json:"polygon"`
	ParentPhysicalThingID       *uuid.UUID     `json:"parent_physical_thing_id"`
	ParentPhysicalThingIDObject *PhysicalThing `json:"-"`
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
	_ = geom.Point{}
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

	var err error

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
			m.ID, err = types.ParseUUID(v)
			if err != nil {
				return wrapError(k, err)
			}
		case "created_at":
			m.CreatedAt, err = types.ParseTime(v)
			if err != nil {
				return wrapError(k, err)
			}
		case "updated_at":
			m.UpdatedAt, err = types.ParseTime(v)
			if err != nil {
				return wrapError(k, err)
			}
		case "deleted_at":
			m.DeletedAt, err = types.ParsePtr(types.ParseTime, v)
			if err != nil {
				return wrapError(k, err)
			}
		case "timestamp":
			m.Timestamp, err = types.ParseTime(v)
			if err != nil {
				return wrapError(k, err)
			}
		case "point":
			m.Point, err = types.ParsePtr(types.ParsePoint, v)
			if err != nil {
				return wrapError(k, err)
			}
		case "polygon":
			m.Polygon, err = types.ParsePtr(types.ParsePolygon, v)
			if err != nil {
				return wrapError(k, err)
			}
		case "parent_physical_thing_id":
			m.ParentPhysicalThingID, err = types.ParsePtr(types.ParseUUID, v)
			if err != nil {
				return wrapError(k, err)
			}
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
		LocationHistoryTableColumns,
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
