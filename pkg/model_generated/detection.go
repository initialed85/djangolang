package model_generated

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/netip"
	"slices"
	"strings"
	"time"

	"github.com/cridenour/go-postgis"
	"github.com/go-chi/chi/v5"
	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
	"github.com/initialed85/djangolang/pkg/config"
	"github.com/initialed85/djangolang/pkg/helpers"
	"github.com/initialed85/djangolang/pkg/introspect"
	"github.com/initialed85/djangolang/pkg/query"
	"github.com/initialed85/djangolang/pkg/server"
	"github.com/initialed85/djangolang/pkg/stream"
	"github.com/initialed85/djangolang/pkg/types"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/exp/maps"
)

type Detection struct {
	ID             uuid.UUID     `json:"id"`
	CreatedAt      time.Time     `json:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at"`
	DeletedAt      *time.Time    `json:"deleted_at"`
	SeenAt         time.Time     `json:"seen_at"`
	ClassID        int64         `json:"class_id"`
	ClassName      string        `json:"class_name"`
	Score          float64       `json:"score"`
	Centroid       pgtype.Vec2   `json:"centroid"`
	BoundingBox    []pgtype.Vec2 `json:"bounding_box"`
	VideoID        uuid.UUID     `json:"video_id"`
	VideoIDObject  *Video        `json:"video_id_object"`
	CameraID       uuid.UUID     `json:"camera_id"`
	CameraIDObject *Camera       `json:"camera_id_object"`
}

var DetectionTable = "detection"

var DetectionTableNamespaceID int32 = 1337 + 2

var (
	DetectionTableIDColumn          = "id"
	DetectionTableCreatedAtColumn   = "created_at"
	DetectionTableUpdatedAtColumn   = "updated_at"
	DetectionTableDeletedAtColumn   = "deleted_at"
	DetectionTableSeenAtColumn      = "seen_at"
	DetectionTableClassIDColumn     = "class_id"
	DetectionTableClassNameColumn   = "class_name"
	DetectionTableScoreColumn       = "score"
	DetectionTableCentroidColumn    = "centroid"
	DetectionTableBoundingBoxColumn = "bounding_box"
	DetectionTableVideoIDColumn     = "video_id"
	DetectionTableCameraIDColumn    = "camera_id"
)

var (
	DetectionTableIDColumnWithTypeCast          = `"id" AS id`
	DetectionTableCreatedAtColumnWithTypeCast   = `"created_at" AS created_at`
	DetectionTableUpdatedAtColumnWithTypeCast   = `"updated_at" AS updated_at`
	DetectionTableDeletedAtColumnWithTypeCast   = `"deleted_at" AS deleted_at`
	DetectionTableSeenAtColumnWithTypeCast      = `"seen_at" AS seen_at`
	DetectionTableClassIDColumnWithTypeCast     = `"class_id" AS class_id`
	DetectionTableClassNameColumnWithTypeCast   = `"class_name" AS class_name`
	DetectionTableScoreColumnWithTypeCast       = `"score" AS score`
	DetectionTableCentroidColumnWithTypeCast    = `"centroid" AS centroid`
	DetectionTableBoundingBoxColumnWithTypeCast = `"bounding_box" AS bounding_box`
	DetectionTableVideoIDColumnWithTypeCast     = `"video_id" AS video_id`
	DetectionTableCameraIDColumnWithTypeCast    = `"camera_id" AS camera_id`
)

var DetectionTableColumns = []string{
	DetectionTableIDColumn,
	DetectionTableCreatedAtColumn,
	DetectionTableUpdatedAtColumn,
	DetectionTableDeletedAtColumn,
	DetectionTableSeenAtColumn,
	DetectionTableClassIDColumn,
	DetectionTableClassNameColumn,
	DetectionTableScoreColumn,
	DetectionTableCentroidColumn,
	DetectionTableBoundingBoxColumn,
	DetectionTableVideoIDColumn,
	DetectionTableCameraIDColumn,
}

var DetectionTableColumnsWithTypeCasts = []string{
	DetectionTableIDColumnWithTypeCast,
	DetectionTableCreatedAtColumnWithTypeCast,
	DetectionTableUpdatedAtColumnWithTypeCast,
	DetectionTableDeletedAtColumnWithTypeCast,
	DetectionTableSeenAtColumnWithTypeCast,
	DetectionTableClassIDColumnWithTypeCast,
	DetectionTableClassNameColumnWithTypeCast,
	DetectionTableScoreColumnWithTypeCast,
	DetectionTableCentroidColumnWithTypeCast,
	DetectionTableBoundingBoxColumnWithTypeCast,
	DetectionTableVideoIDColumnWithTypeCast,
	DetectionTableCameraIDColumnWithTypeCast,
}

var DetectionIntrospectedTable *introspect.Table

var DetectionTableColumnLookup map[string]*introspect.Column

var (
	DetectionTablePrimaryKeyColumn = DetectionTableIDColumn
)

func init() {
	DetectionIntrospectedTable = tableByName[DetectionTable]

	/* only needed during templating */
	if DetectionIntrospectedTable == nil {
		DetectionIntrospectedTable = &introspect.Table{}
	}

	DetectionTableColumnLookup = DetectionIntrospectedTable.ColumnByName
}

type DetectionOnePathParams struct {
	PrimaryKey uuid.UUID `json:"primaryKey"`
}

type DetectionLoadQueryParams struct {
	Depth *int `json:"depth"`
}

/*
TODO: find a way to not need this- there is a piece in the templating logic
that uses goimports but pending where the code is built, it may resolve
the packages to import to the wrong ones (causing odd failures)
these are just here to ensure we don't get unused imports
*/
var _ = []any{
	time.Time{},
	uuid.UUID{},
	pgtype.Hstore{},
	postgis.PointZ{},
	netip.Prefix{},
	errors.Is,
	sql.ErrNoRows,
}

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

	wrapError := func(k string, v any, err error) error {
		return fmt.Errorf("%v: %#+v; error; %v", k, v, err)
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

			temp1, err := types.ParseUUID(v)
			if err != nil {
				return wrapError(k, v, err)
			}

			temp2, ok := temp1.(uuid.UUID)
			if !ok {
				if temp1 != nil {
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uuid.UUID", temp1))
				}
			}

			m.ID = temp2

		case "created_at":
			if v == nil {
				continue
			}

			temp1, err := types.ParseTime(v)
			if err != nil {
				return wrapError(k, v, err)
			}

			temp2, ok := temp1.(time.Time)
			if !ok {
				if temp1 != nil {
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uucreated_at.UUID", temp1))
				}
			}

			m.CreatedAt = temp2

		case "updated_at":
			if v == nil {
				continue
			}

			temp1, err := types.ParseTime(v)
			if err != nil {
				return wrapError(k, v, err)
			}

			temp2, ok := temp1.(time.Time)
			if !ok {
				if temp1 != nil {
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uuupdated_at.UUID", temp1))
				}
			}

			m.UpdatedAt = temp2

		case "deleted_at":
			if v == nil {
				continue
			}

			temp1, err := types.ParseTime(v)
			if err != nil {
				return wrapError(k, v, err)
			}

			temp2, ok := temp1.(time.Time)
			if !ok {
				if temp1 != nil {
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uudeleted_at.UUID", temp1))
				}
			}

			m.DeletedAt = &temp2

		case "seen_at":
			if v == nil {
				continue
			}

			temp1, err := types.ParseTime(v)
			if err != nil {
				return wrapError(k, v, err)
			}

			temp2, ok := temp1.(time.Time)
			if !ok {
				if temp1 != nil {
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uuseen_at.UUID", temp1))
				}
			}

			m.SeenAt = temp2

		case "class_id":
			if v == nil {
				continue
			}

			temp1, err := types.ParseInt(v)
			if err != nil {
				return wrapError(k, v, err)
			}

			temp2, ok := temp1.(int64)
			if !ok {
				if temp1 != nil {
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uuclass_id.UUID", temp1))
				}
			}

			m.ClassID = temp2

		case "class_name":
			if v == nil {
				continue
			}

			temp1, err := types.ParseString(v)
			if err != nil {
				return wrapError(k, v, err)
			}

			temp2, ok := temp1.(string)
			if !ok {
				if temp1 != nil {
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uuclass_name.UUID", temp1))
				}
			}

			m.ClassName = temp2

		case "score":
			if v == nil {
				continue
			}

			temp1, err := types.ParseFloat(v)
			if err != nil {
				return wrapError(k, v, err)
			}

			temp2, ok := temp1.(float64)
			if !ok {
				if temp1 != nil {
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uuscore.UUID", temp1))
				}
			}

			m.Score = temp2

		case "centroid":
			if v == nil {
				continue
			}

			temp1, err := types.ParsePoint(v)
			if err != nil {
				return wrapError(k, v, err)
			}

			temp2, ok := temp1.(pgtype.Vec2)
			if !ok {
				if temp1 != nil {
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uucentroid.UUID", temp1))
				}
			}

			m.Centroid = temp2

		case "bounding_box":
			if v == nil {
				continue
			}

			temp1, err := types.ParsePolygon(v)
			if err != nil {
				return wrapError(k, v, err)
			}

			temp2, ok := temp1.([]pgtype.Vec2)
			if !ok {
				if temp1 != nil {
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uubounding_box.UUID", temp1))
				}
			}

			m.BoundingBox = temp2

		case "video_id":
			if v == nil {
				continue
			}

			temp1, err := types.ParseUUID(v)
			if err != nil {
				return wrapError(k, v, err)
			}

			temp2, ok := temp1.(uuid.UUID)
			if !ok {
				if temp1 != nil {
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uuvideo_id.UUID", temp1))
				}
			}

			m.VideoID = temp2

		case "camera_id":
			if v == nil {
				continue
			}

			temp1, err := types.ParseUUID(v)
			if err != nil {
				return wrapError(k, v, err)
			}

			temp2, ok := temp1.(uuid.UUID)
			if !ok {
				if temp1 != nil {
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uucamera_id.UUID", temp1))
				}
			}

			m.CameraID = temp2

		}
	}

	return nil
}

func (m *Detection) Reload(ctx context.Context, tx pgx.Tx, includeDeleteds ...bool) error {
	extraWhere := ""
	if len(includeDeleteds) > 0 && includeDeleteds[0] {
		if slices.Contains(DetectionTableColumns, "deleted_at") {
			extraWhere = "\n    AND (deleted_at IS null OR deleted_at IS NOT null)"
		}
	}

	ctx, cleanup := query.WithQueryID(ctx)
	defer cleanup()

	ctx = query.WithMaxDepth(ctx, nil)

	o, _, _, _, _, err := SelectDetection(
		ctx,
		tx,
		fmt.Sprintf("%v = $1%v", m.GetPrimaryKeyColumn(), extraWhere),
		m.GetPrimaryKeyValue(),
	)
	if err != nil {
		return err
	}

	m.ID = o.ID
	m.CreatedAt = o.CreatedAt
	m.UpdatedAt = o.UpdatedAt
	m.DeletedAt = o.DeletedAt
	m.SeenAt = o.SeenAt
	m.ClassID = o.ClassID
	m.ClassName = o.ClassName
	m.Score = o.Score
	m.Centroid = o.Centroid
	m.BoundingBox = o.BoundingBox
	m.VideoID = o.VideoID
	m.VideoIDObject = o.VideoIDObject
	m.CameraID = o.CameraID
	m.CameraIDObject = o.CameraIDObject

	return nil
}

func (m *Detection) Insert(ctx context.Context, tx pgx.Tx, setPrimaryKey bool, setZeroValues bool, forceSetValuesForFields ...string) error {
	columns := make([]string, 0)
	values := make([]any, 0)

	if setPrimaryKey && (setZeroValues || !types.IsZeroUUID(m.ID) || slices.Contains(forceSetValuesForFields, DetectionTableIDColumn) || isRequired(DetectionTableColumnLookup, DetectionTableIDColumn)) {
		columns = append(columns, DetectionTableIDColumn)

		v, err := types.FormatUUID(m.ID)
		if err != nil {
			return fmt.Errorf("failed to handle m.ID; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.CreatedAt) || slices.Contains(forceSetValuesForFields, DetectionTableCreatedAtColumn) || isRequired(DetectionTableColumnLookup, DetectionTableCreatedAtColumn) {
		columns = append(columns, DetectionTableCreatedAtColumn)

		v, err := types.FormatTime(m.CreatedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.CreatedAt; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.UpdatedAt) || slices.Contains(forceSetValuesForFields, DetectionTableUpdatedAtColumn) || isRequired(DetectionTableColumnLookup, DetectionTableUpdatedAtColumn) {
		columns = append(columns, DetectionTableUpdatedAtColumn)

		v, err := types.FormatTime(m.UpdatedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.UpdatedAt; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.DeletedAt) || slices.Contains(forceSetValuesForFields, DetectionTableDeletedAtColumn) || isRequired(DetectionTableColumnLookup, DetectionTableDeletedAtColumn) {
		columns = append(columns, DetectionTableDeletedAtColumn)

		v, err := types.FormatTime(m.DeletedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.DeletedAt; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.SeenAt) || slices.Contains(forceSetValuesForFields, DetectionTableSeenAtColumn) || isRequired(DetectionTableColumnLookup, DetectionTableSeenAtColumn) {
		columns = append(columns, DetectionTableSeenAtColumn)

		v, err := types.FormatTime(m.SeenAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.SeenAt; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroInt(m.ClassID) || slices.Contains(forceSetValuesForFields, DetectionTableClassIDColumn) || isRequired(DetectionTableColumnLookup, DetectionTableClassIDColumn) {
		columns = append(columns, DetectionTableClassIDColumn)

		v, err := types.FormatInt(m.ClassID)
		if err != nil {
			return fmt.Errorf("failed to handle m.ClassID; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroString(m.ClassName) || slices.Contains(forceSetValuesForFields, DetectionTableClassNameColumn) || isRequired(DetectionTableColumnLookup, DetectionTableClassNameColumn) {
		columns = append(columns, DetectionTableClassNameColumn)

		v, err := types.FormatString(m.ClassName)
		if err != nil {
			return fmt.Errorf("failed to handle m.ClassName; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroFloat(m.Score) || slices.Contains(forceSetValuesForFields, DetectionTableScoreColumn) || isRequired(DetectionTableColumnLookup, DetectionTableScoreColumn) {
		columns = append(columns, DetectionTableScoreColumn)

		v, err := types.FormatFloat(m.Score)
		if err != nil {
			return fmt.Errorf("failed to handle m.Score; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroPoint(m.Centroid) || slices.Contains(forceSetValuesForFields, DetectionTableCentroidColumn) || isRequired(DetectionTableColumnLookup, DetectionTableCentroidColumn) {
		columns = append(columns, DetectionTableCentroidColumn)

		v, err := types.FormatPoint(m.Centroid)
		if err != nil {
			return fmt.Errorf("failed to handle m.Centroid; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroPolygon(m.BoundingBox) || slices.Contains(forceSetValuesForFields, DetectionTableBoundingBoxColumn) || isRequired(DetectionTableColumnLookup, DetectionTableBoundingBoxColumn) {
		columns = append(columns, DetectionTableBoundingBoxColumn)

		v, err := types.FormatPolygon(m.BoundingBox)
		if err != nil {
			return fmt.Errorf("failed to handle m.BoundingBox; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroUUID(m.VideoID) || slices.Contains(forceSetValuesForFields, DetectionTableVideoIDColumn) || isRequired(DetectionTableColumnLookup, DetectionTableVideoIDColumn) {
		columns = append(columns, DetectionTableVideoIDColumn)

		v, err := types.FormatUUID(m.VideoID)
		if err != nil {
			return fmt.Errorf("failed to handle m.VideoID; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroUUID(m.CameraID) || slices.Contains(forceSetValuesForFields, DetectionTableCameraIDColumn) || isRequired(DetectionTableColumnLookup, DetectionTableCameraIDColumn) {
		columns = append(columns, DetectionTableCameraIDColumn)

		v, err := types.FormatUUID(m.CameraID)
		if err != nil {
			return fmt.Errorf("failed to handle m.CameraID; %v", err)
		}

		values = append(values, v)
	}

	ctx, cleanup := query.WithQueryID(ctx)
	defer cleanup()

	ctx = query.WithMaxDepth(ctx, nil)

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
		return fmt.Errorf("failed to insert %#+v; %v", m, err)
	}
	v := (*item)[DetectionTableIDColumn]

	if v == nil {
		return fmt.Errorf("failed to find %v in %#+v", DetectionTableIDColumn, item)
	}

	wrapError := func(err error) error {
		return fmt.Errorf(
			"failed to treat %v: %#+v as uuid.UUID: %v",
			DetectionTableIDColumn,
			(*item)[DetectionTableIDColumn],
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

	err = m.Reload(ctx, tx, slices.Contains(forceSetValuesForFields, "deleted_at"))
	if err != nil {
		return fmt.Errorf("failed to reload after insert; %v", err)
	}

	return nil
}

func (m *Detection) Update(ctx context.Context, tx pgx.Tx, setZeroValues bool, forceSetValuesForFields ...string) error {
	columns := make([]string, 0)
	values := make([]any, 0)

	if setZeroValues || !types.IsZeroTime(m.CreatedAt) || slices.Contains(forceSetValuesForFields, DetectionTableCreatedAtColumn) {
		columns = append(columns, DetectionTableCreatedAtColumn)

		v, err := types.FormatTime(m.CreatedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.CreatedAt; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.UpdatedAt) || slices.Contains(forceSetValuesForFields, DetectionTableUpdatedAtColumn) {
		columns = append(columns, DetectionTableUpdatedAtColumn)

		v, err := types.FormatTime(m.UpdatedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.UpdatedAt; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.DeletedAt) || slices.Contains(forceSetValuesForFields, DetectionTableDeletedAtColumn) {
		columns = append(columns, DetectionTableDeletedAtColumn)

		v, err := types.FormatTime(m.DeletedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.DeletedAt; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.SeenAt) || slices.Contains(forceSetValuesForFields, DetectionTableSeenAtColumn) {
		columns = append(columns, DetectionTableSeenAtColumn)

		v, err := types.FormatTime(m.SeenAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.SeenAt; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroInt(m.ClassID) || slices.Contains(forceSetValuesForFields, DetectionTableClassIDColumn) {
		columns = append(columns, DetectionTableClassIDColumn)

		v, err := types.FormatInt(m.ClassID)
		if err != nil {
			return fmt.Errorf("failed to handle m.ClassID; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroString(m.ClassName) || slices.Contains(forceSetValuesForFields, DetectionTableClassNameColumn) {
		columns = append(columns, DetectionTableClassNameColumn)

		v, err := types.FormatString(m.ClassName)
		if err != nil {
			return fmt.Errorf("failed to handle m.ClassName; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroFloat(m.Score) || slices.Contains(forceSetValuesForFields, DetectionTableScoreColumn) {
		columns = append(columns, DetectionTableScoreColumn)

		v, err := types.FormatFloat(m.Score)
		if err != nil {
			return fmt.Errorf("failed to handle m.Score; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroPoint(m.Centroid) || slices.Contains(forceSetValuesForFields, DetectionTableCentroidColumn) {
		columns = append(columns, DetectionTableCentroidColumn)

		v, err := types.FormatPoint(m.Centroid)
		if err != nil {
			return fmt.Errorf("failed to handle m.Centroid; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroPolygon(m.BoundingBox) || slices.Contains(forceSetValuesForFields, DetectionTableBoundingBoxColumn) {
		columns = append(columns, DetectionTableBoundingBoxColumn)

		v, err := types.FormatPolygon(m.BoundingBox)
		if err != nil {
			return fmt.Errorf("failed to handle m.BoundingBox; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroUUID(m.VideoID) || slices.Contains(forceSetValuesForFields, DetectionTableVideoIDColumn) {
		columns = append(columns, DetectionTableVideoIDColumn)

		v, err := types.FormatUUID(m.VideoID)
		if err != nil {
			return fmt.Errorf("failed to handle m.VideoID; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroUUID(m.CameraID) || slices.Contains(forceSetValuesForFields, DetectionTableCameraIDColumn) {
		columns = append(columns, DetectionTableCameraIDColumn)

		v, err := types.FormatUUID(m.CameraID)
		if err != nil {
			return fmt.Errorf("failed to handle m.CameraID; %v", err)
		}

		values = append(values, v)
	}

	v, err := types.FormatUUID(m.ID)
	if err != nil {
		return fmt.Errorf("failed to handle m.ID; %v", err)
	}

	values = append(values, v)

	ctx, cleanup := query.WithQueryID(ctx)
	defer cleanup()

	ctx = query.WithMaxDepth(ctx, nil)

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
		return fmt.Errorf("failed to update %#+v; %v", m, err)
	}

	err = m.Reload(ctx, tx, slices.Contains(forceSetValuesForFields, "deleted_at"))
	if err != nil {
		return fmt.Errorf("failed to reload after update")
	}

	return nil
}

func (m *Detection) Delete(ctx context.Context, tx pgx.Tx, hardDeletes ...bool) error {
	hardDelete := false
	if len(hardDeletes) > 0 {
		hardDelete = hardDeletes[0]
	}

	if !hardDelete && slices.Contains(DetectionTableColumns, "deleted_at") {
		m.DeletedAt = helpers.Ptr(time.Now().UTC())
		err := m.Update(ctx, tx, false, "deleted_at")
		if err != nil {
			return fmt.Errorf("failed to soft-delete (update) %#+v; %v", m, err)
		}
	}

	values := make([]any, 0)
	v, err := types.FormatUUID(m.ID)
	if err != nil {
		return fmt.Errorf("failed to handle m.ID; %v", err)
	}

	values = append(values, v)

	ctx, cleanup := query.WithQueryID(ctx)
	defer cleanup()

	ctx = query.WithMaxDepth(ctx, nil)

	err = query.Delete(
		ctx,
		tx,
		DetectionTable,
		fmt.Sprintf("%v = $$??", DetectionTableIDColumn),
		values...,
	)
	if err != nil {
		return fmt.Errorf("failed to delete %#+v; %v", m, err)
	}

	_ = m.Reload(ctx, tx, true)

	return nil
}

func (m *Detection) LockTable(ctx context.Context, tx pgx.Tx, timeouts ...time.Duration) error {
	return query.LockTable(ctx, tx, DetectionTable, timeouts...)
}

func (m *Detection) LockTableWithRetries(ctx context.Context, tx pgx.Tx, overallTimeout time.Duration, individualAttempttimeout time.Duration) error {
	return query.LockTableWithRetries(ctx, tx, DetectionTable, overallTimeout, individualAttempttimeout)
}

func (m *Detection) AdvisoryLock(ctx context.Context, tx pgx.Tx, key int32, timeouts ...time.Duration) error {
	return query.AdvisoryLock(ctx, tx, DetectionTableNamespaceID, key, timeouts...)
}

func (m *Detection) AdvisoryLockWithRetries(ctx context.Context, tx pgx.Tx, key int32, overallTimeout time.Duration, individualAttempttimeout time.Duration) error {
	return query.AdvisoryLockWithRetries(ctx, tx, DetectionTableNamespaceID, key, overallTimeout, individualAttempttimeout)
}

func SelectDetections(ctx context.Context, tx pgx.Tx, where string, orderBy *string, limit *int, offset *int, values ...any) ([]*Detection, int64, int64, int64, int64, error) {
	before := time.Now()

	if config.Debug() {
		log.Printf("entered SelectDetections")

		defer func() {
			log.Printf("exited SelectDetections in %s", time.Since(before))
		}()
	}
	if slices.Contains(DetectionTableColumns, "deleted_at") {
		if !strings.Contains(where, "deleted_at") {
			if where != "" {
				where += "\n    AND "
			}

			where += "deleted_at IS null"
		}
	}

	ctx, cleanup := query.WithQueryID(ctx)
	defer cleanup()

	possiblePathValue := query.GetCurrentPathValue(ctx)
	isLoadQuery := possiblePathValue != nil && len(possiblePathValue.VisitedTableNames) > 0

	shouldLoad := query.ShouldLoad(ctx, DetectionTable) || query.ShouldLoad(ctx, fmt.Sprintf("referenced_by_%s", DetectionTable))

	var ok bool
	ctx, ok = query.HandleQueryPathGraphCycles(ctx, fmt.Sprintf("%s{%v}", DetectionTable, nil), !isLoadQuery)
	if !ok && !shouldLoad {
		if config.Debug() {
			log.Printf("skipping SelectDetection early (query.ShouldLoad(): %v, query.HandleQueryPathGraphCycles(): %v)", shouldLoad, ok)
		}
		return []*Detection{}, 0, 0, 0, 0, nil
	}

	items, count, totalCount, page, totalPages, err := query.Select(
		ctx,
		tx,
		DetectionTableColumnsWithTypeCasts,
		DetectionTable,
		where,
		orderBy,
		limit,
		offset,
		values...,
	)
	if err != nil {
		return nil, 0, 0, 0, 0, fmt.Errorf("failed to call SelectDetections; %v", err)
	}

	objects := make([]*Detection, 0)

	for _, item := range *items {
		object := &Detection{}

		err = object.FromItem(item)
		if err != nil {
			return nil, 0, 0, 0, 0, err
		}

		if !types.IsZeroUUID(object.VideoID) {
			ctx, ok := query.HandleQueryPathGraphCycles(ctx, fmt.Sprintf("%s{%v}", VideoTable, object.VideoID), true)
			shouldLoad := query.ShouldLoad(ctx, VideoTable)
			if ok || shouldLoad {
				thisBefore := time.Now()

				if config.Debug() {
					log.Printf("loading SelectDetections->SelectVideo for object.VideoIDObject{%s: %v}", VideoTablePrimaryKeyColumn, object.VideoID)
				}

				object.VideoIDObject, _, _, _, _, err = SelectVideo(
					ctx,
					tx,
					fmt.Sprintf("%v = $1", VideoTablePrimaryKeyColumn),
					object.VideoID,
				)
				if err != nil {
					if !errors.Is(err, sql.ErrNoRows) {
						return nil, 0, 0, 0, 0, err
					}
				}

				if config.Debug() {
					log.Printf("loaded SelectDetections->SelectVideo for object.VideoIDObject in %s", time.Since(thisBefore))
				}
			}
		}

		if !types.IsZeroUUID(object.CameraID) {
			ctx, ok := query.HandleQueryPathGraphCycles(ctx, fmt.Sprintf("%s{%v}", CameraTable, object.CameraID), true)
			shouldLoad := query.ShouldLoad(ctx, CameraTable)
			if ok || shouldLoad {
				thisBefore := time.Now()

				if config.Debug() {
					log.Printf("loading SelectDetections->SelectCamera for object.CameraIDObject{%s: %v}", CameraTablePrimaryKeyColumn, object.CameraID)
				}

				object.CameraIDObject, _, _, _, _, err = SelectCamera(
					ctx,
					tx,
					fmt.Sprintf("%v = $1", CameraTablePrimaryKeyColumn),
					object.CameraID,
				)
				if err != nil {
					if !errors.Is(err, sql.ErrNoRows) {
						return nil, 0, 0, 0, 0, err
					}
				}

				if config.Debug() {
					log.Printf("loaded SelectDetections->SelectCamera for object.CameraIDObject in %s", time.Since(thisBefore))
				}
			}
		}

		objects = append(objects, object)
	}

	return objects, count, totalCount, page, totalPages, nil
}

func SelectDetection(ctx context.Context, tx pgx.Tx, where string, values ...any) (*Detection, int64, int64, int64, int64, error) {
	ctx, cleanup := query.WithQueryID(ctx)
	defer cleanup()

	ctx = query.WithMaxDepth(ctx, nil)

	objects, _, _, _, _, err := SelectDetections(
		ctx,
		tx,
		where,
		nil,
		helpers.Ptr(2),
		helpers.Ptr(0),
		values...,
	)
	if err != nil {
		return nil, 0, 0, 0, 0, fmt.Errorf("failed to call SelectDetection; %v", err)
	}

	if len(objects) > 1 {
		return nil, 0, 0, 0, 0, fmt.Errorf("attempt to call SelectDetection returned more than 1 row")
	}

	if len(objects) < 1 {
		return nil, 0, 0, 0, 0, sql.ErrNoRows
	}

	object := objects[0]

	count := int64(1)
	totalCount := count
	page := int64(1)
	totalPages := page

	return object, count, totalCount, page, totalPages, nil
}

func handleGetDetections(arguments *server.SelectManyArguments, db *pgxpool.Pool) ([]*Detection, int64, int64, int64, int64, error) {
	tx, err := db.Begin(arguments.Ctx)
	if err != nil {
		return nil, 0, 0, 0, 0, err
	}

	defer func() {
		_ = tx.Rollback(arguments.Ctx)
	}()

	objects, count, totalCount, page, totalPages, err := SelectDetections(arguments.Ctx, tx, arguments.Where, arguments.OrderBy, arguments.Limit, arguments.Offset, arguments.Values...)
	if err != nil {
		return nil, 0, 0, 0, 0, err
	}

	err = tx.Commit(arguments.Ctx)
	if err != nil {
		return nil, 0, 0, 0, 0, err
	}

	return objects, count, totalCount, page, totalPages, nil
}

func handleGetDetection(arguments *server.SelectOneArguments, db *pgxpool.Pool, primaryKey uuid.UUID) ([]*Detection, int64, int64, int64, int64, error) {
	tx, err := db.Begin(arguments.Ctx)
	if err != nil {
		return nil, 0, 0, 0, 0, err
	}

	defer func() {
		_ = tx.Rollback(arguments.Ctx)
	}()

	object, count, totalCount, page, totalPages, err := SelectDetection(arguments.Ctx, tx, arguments.Where, arguments.Values...)
	if err != nil {
		return nil, 0, 0, 0, 0, err
	}

	err = tx.Commit(arguments.Ctx)
	if err != nil {
		return nil, 0, 0, 0, 0, err
	}

	return []*Detection{object}, count, totalCount, page, totalPages, nil
}

func handlePostDetections(arguments *server.LoadArguments, db *pgxpool.Pool, waitForChange server.WaitForChange, objects []*Detection, forceSetValuesForFieldsByObjectIndex [][]string) ([]*Detection, int64, int64, int64, int64, error) {
	tx, err := db.Begin(arguments.Ctx)
	if err != nil {
		err = fmt.Errorf("failed to begin DB transaction; %v", err)
		return nil, 0, 0, 0, 0, err
	}

	defer func() {
		_ = tx.Rollback(arguments.Ctx)
	}()

	xid, err := query.GetXid(arguments.Ctx, tx)
	if err != nil {
		err = fmt.Errorf("failed to get xid; %v", err)
		return nil, 0, 0, 0, 0, err
	}
	_ = xid

	for i, object := range objects {
		err = object.Insert(arguments.Ctx, tx, false, false, forceSetValuesForFieldsByObjectIndex[i]...)
		if err != nil {
			err = fmt.Errorf("failed to insert %#+v; %v", object, err)
			return nil, 0, 0, 0, 0, err
		}

		objects[i] = object
	}

	errs := make(chan error, 1)
	go func() {
		_, err := waitForChange(arguments.Ctx, []stream.Action{stream.INSERT}, DetectionTable, xid)
		if err != nil {
			err = fmt.Errorf("failed to wait for change; %v", err)
			errs <- err
			return
		}

		errs <- nil
	}()

	err = tx.Commit(arguments.Ctx)
	if err != nil {
		err = fmt.Errorf("failed to commit DB transaction; %v", err)
		return nil, 0, 0, 0, 0, err
	}

	select {
	case <-arguments.Ctx.Done():
		err = fmt.Errorf("context canceled")
		return nil, 0, 0, 0, 0, err
	case err = <-errs:
		if err != nil {
			return nil, 0, 0, 0, 0, err
		}
	}

	count := int64(len(objects))
	totalCount := count
	page := int64(1)
	totalPages := page

	return objects, count, totalCount, page, totalPages, nil
}

func handlePutDetection(arguments *server.LoadArguments, db *pgxpool.Pool, waitForChange server.WaitForChange, object *Detection) ([]*Detection, int64, int64, int64, int64, error) {
	tx, err := db.Begin(arguments.Ctx)
	if err != nil {
		err = fmt.Errorf("failed to begin DB transaction; %v", err)
		return nil, 0, 0, 0, 0, err
	}

	defer func() {
		_ = tx.Rollback(arguments.Ctx)
	}()

	xid, err := query.GetXid(arguments.Ctx, tx)
	if err != nil {
		err = fmt.Errorf("failed to get xid; %v", err)
		return nil, 0, 0, 0, 0, err
	}
	_ = xid

	err = object.Update(arguments.Ctx, tx, true)
	if err != nil {
		err = fmt.Errorf("failed to update %#+v; %v", object, err)
		return nil, 0, 0, 0, 0, err
	}

	errs := make(chan error, 1)
	go func() {
		_, err := waitForChange(arguments.Ctx, []stream.Action{stream.UPDATE, stream.SOFT_DELETE, stream.SOFT_RESTORE, stream.SOFT_UPDATE}, DetectionTable, xid)
		if err != nil {
			err = fmt.Errorf("failed to wait for change; %v", err)
			errs <- err
			return
		}

		errs <- nil
	}()

	err = tx.Commit(arguments.Ctx)
	if err != nil {
		err = fmt.Errorf("failed to commit DB transaction; %v", err)
		return nil, 0, 0, 0, 0, err
	}

	select {
	case <-arguments.Ctx.Done():
		err = fmt.Errorf("context canceled")
		return nil, 0, 0, 0, 0, err
	case err = <-errs:
		if err != nil {
			return nil, 0, 0, 0, 0, err
		}
	}

	count := int64(1)
	totalCount := count
	page := int64(1)
	totalPages := page

	return []*Detection{object}, count, totalCount, page, totalPages, nil
}

func handlePatchDetection(arguments *server.LoadArguments, db *pgxpool.Pool, waitForChange server.WaitForChange, object *Detection, forceSetValuesForFields []string) ([]*Detection, int64, int64, int64, int64, error) {
	tx, err := db.Begin(arguments.Ctx)
	if err != nil {
		err = fmt.Errorf("failed to begin DB transaction; %v", err)
		return nil, 0, 0, 0, 0, err
	}

	defer func() {
		_ = tx.Rollback(arguments.Ctx)
	}()

	xid, err := query.GetXid(arguments.Ctx, tx)
	if err != nil {
		err = fmt.Errorf("failed to get xid; %v", err)
		return nil, 0, 0, 0, 0, err
	}
	_ = xid

	err = object.Update(arguments.Ctx, tx, false, forceSetValuesForFields...)
	if err != nil {
		err = fmt.Errorf("failed to update %#+v; %v", object, err)
		return nil, 0, 0, 0, 0, err
	}

	errs := make(chan error, 1)
	go func() {
		_, err := waitForChange(arguments.Ctx, []stream.Action{stream.UPDATE, stream.SOFT_DELETE, stream.SOFT_RESTORE, stream.SOFT_UPDATE}, DetectionTable, xid)
		if err != nil {
			err = fmt.Errorf("failed to wait for change; %v", err)
			errs <- err
			return
		}

		errs <- nil
	}()

	err = tx.Commit(arguments.Ctx)
	if err != nil {
		err = fmt.Errorf("failed to commit DB transaction; %v", err)
		return nil, 0, 0, 0, 0, err
	}

	select {
	case <-arguments.Ctx.Done():
		err = fmt.Errorf("context canceled")
		return nil, 0, 0, 0, 0, err
	case err = <-errs:
		if err != nil {
			return nil, 0, 0, 0, 0, err
		}
	}

	count := int64(1)
	totalCount := count
	page := int64(1)
	totalPages := page

	return []*Detection{object}, count, totalCount, page, totalPages, nil
}

func handleDeleteDetection(arguments *server.LoadArguments, db *pgxpool.Pool, waitForChange server.WaitForChange, object *Detection) error {
	tx, err := db.Begin(arguments.Ctx)
	if err != nil {
		err = fmt.Errorf("failed to begin DB transaction; %v", err)
		return err
	}

	defer func() {
		_ = tx.Rollback(arguments.Ctx)
	}()

	xid, err := query.GetXid(arguments.Ctx, tx)
	if err != nil {
		err = fmt.Errorf("failed to get xid; %v", err)
		return err
	}
	_ = xid

	err = object.Delete(arguments.Ctx, tx)
	if err != nil {
		err = fmt.Errorf("failed to delete %#+v; %v", object, err)
		return err
	}

	errs := make(chan error, 1)
	go func() {
		_, err := waitForChange(arguments.Ctx, []stream.Action{stream.DELETE, stream.SOFT_DELETE}, DetectionTable, xid)
		if err != nil {
			err = fmt.Errorf("failed to wait for change; %v", err)
			errs <- err
			return
		}

		errs <- nil
	}()

	err = tx.Commit(arguments.Ctx)
	if err != nil {
		err = fmt.Errorf("failed to commit DB transaction; %v", err)
		return err
	}

	select {
	case <-arguments.Ctx.Done():
		err = fmt.Errorf("context canceled")
		return err
	case err = <-errs:
		if err != nil {
			return err
		}
	}

	return nil
}

func GetDetectionRouter(db *pgxpool.Pool, redisPool *redis.Pool, httpMiddlewares []server.HTTPMiddleware, objectMiddlewares []server.ObjectMiddleware, waitForChange server.WaitForChange) chi.Router {
	r := chi.NewRouter()

	for _, m := range httpMiddlewares {
		r.Use(m)
	}

	func() {
		getManyHandler, err := getHTTPHandler(
			http.MethodGet,
			"/detections",
			http.StatusOK,
			func(
				ctx context.Context,
				pathParams server.EmptyPathParams,
				queryParams map[string]any,
				req server.EmptyRequest,
				rawReq any,
			) (server.Response[Detection], error) {
				before := time.Now()

				redisConn := redisPool.Get()
				defer func() {
					_ = redisConn.Close()
				}()

				arguments, err := server.GetSelectManyArguments(ctx, queryParams, DetectionIntrospectedTable, nil, nil)
				if err != nil {
					if config.Debug() {
						log.Printf("request cache not yet reached; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return server.Response[Detection]{}, err
				}

				cachedResponseAsJSON, cacheHit, err := server.GetCachedResponseAsJSON(arguments.RequestHash, redisConn)
				if err != nil {
					if config.Debug() {
						log.Printf("request cache failed; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return server.Response[Detection]{}, err
				}

				if cacheHit {
					var cachedResponse server.Response[Detection]

					/* TODO: it'd be nice to be able to avoid this (i.e. just pass straight through) */
					err = json.Unmarshal(cachedResponseAsJSON, &cachedResponse)
					if err != nil {
						if config.Debug() {
							log.Printf("request cache hit but failed unmarshal; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
						}

						return server.Response[Detection]{}, err
					}

					if config.Debug() {
						log.Printf("request cache hit; request succeeded in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return cachedResponse, nil
				}

				objects, count, totalCount, _, _, err := handleGetDetections(arguments, db)
				if err != nil {
					if config.Debug() {
						log.Printf("request cache missed; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return server.Response[Detection]{}, err
				}

				limit := int64(0)
				if arguments.Limit != nil {
					limit = int64(*arguments.Limit)
				}

				offset := int64(0)
				if arguments.Offset != nil {
					offset = int64(*arguments.Offset)
				}

				response := server.Response[Detection]{
					Status:     http.StatusOK,
					Success:    true,
					Error:      nil,
					Objects:    objects,
					Count:      count,
					TotalCount: totalCount,
					Limit:      limit,
					Offset:     offset,
				}

				/* TODO: it'd be nice to be able to avoid this (i.e. just marshal once, further out) */
				responseAsJSON, err := json.Marshal(response)
				if err != nil {
					if config.Debug() {
						log.Printf("request cache missed; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return server.Response[Detection]{}, err
				}

				err = server.StoreCachedResponse(arguments.RequestHash, redisConn, responseAsJSON)
				if err != nil {
					log.Printf("warning; %v", err)
				}

				if config.Debug() {
					log.Printf("request cache missed; request succeeded in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
				}

				return response, nil
			},
			Detection{},
			DetectionIntrospectedTable,
		)
		if err != nil {
			panic(err)
		}
		r.Get(getManyHandler.PathWithinRouter, getManyHandler.ServeHTTP)
	}()

	func() {
		getOneHandler, err := getHTTPHandler(
			http.MethodGet,
			"/detections/{primaryKey}",
			http.StatusOK,
			func(
				ctx context.Context,
				pathParams DetectionOnePathParams,
				queryParams DetectionLoadQueryParams,
				req server.EmptyRequest,
				rawReq any,
			) (server.Response[Detection], error) {
				before := time.Now()

				redisConn := redisPool.Get()
				defer func() {
					_ = redisConn.Close()
				}()

				arguments, err := server.GetSelectOneArguments(ctx, queryParams.Depth, DetectionIntrospectedTable, pathParams.PrimaryKey, nil, nil)
				if err != nil {
					if config.Debug() {
						log.Printf("request cache not yet reached; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return server.Response[Detection]{}, err
				}

				cachedResponseAsJSON, cacheHit, err := server.GetCachedResponseAsJSON(arguments.RequestHash, redisConn)
				if err != nil {
					if config.Debug() {
						log.Printf("request cache failed; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return server.Response[Detection]{}, err
				}

				if cacheHit {
					var cachedResponse server.Response[Detection]

					/* TODO: it'd be nice to be able to avoid this (i.e. just pass straight through) */
					err = json.Unmarshal(cachedResponseAsJSON, &cachedResponse)
					if err != nil {
						if config.Debug() {
							log.Printf("request cache hit but failed unmarshal; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
						}

						return server.Response[Detection]{}, err
					}

					if config.Debug() {
						log.Printf("request cache hit; request succeeded in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return cachedResponse, nil
				}

				objects, count, totalCount, _, _, err := handleGetDetection(arguments, db, pathParams.PrimaryKey)
				if err != nil {
					if config.Debug() {
						log.Printf("request cache missed; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return server.Response[Detection]{}, err
				}

				limit := int64(0)

				offset := int64(0)

				response := server.Response[Detection]{
					Status:     http.StatusOK,
					Success:    true,
					Error:      nil,
					Objects:    objects,
					Count:      count,
					TotalCount: totalCount,
					Limit:      limit,
					Offset:     offset,
				}

				/* TODO: it'd be nice to be able to avoid this (i.e. just marshal once, further out) */
				responseAsJSON, err := json.Marshal(response)
				if err != nil {
					if config.Debug() {
						log.Printf("request cache missed; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return server.Response[Detection]{}, err
				}

				err = server.StoreCachedResponse(arguments.RequestHash, redisConn, responseAsJSON)
				if err != nil {
					log.Printf("warning; %v", err)
				}

				if config.Debug() {
					log.Printf("request cache hit; request succeeded in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
				}

				return response, nil
			},
			Detection{},
			DetectionIntrospectedTable,
		)
		if err != nil {
			panic(err)
		}
		r.Get(getOneHandler.PathWithinRouter, getOneHandler.ServeHTTP)
	}()

	func() {
		postHandler, err := getHTTPHandler(
			http.MethodPost,
			"/detections",
			http.StatusCreated,
			func(
				ctx context.Context,
				pathParams server.EmptyPathParams,
				queryParams DetectionLoadQueryParams,
				req []*Detection,
				rawReq any,
			) (server.Response[Detection], error) {
				allRawItems, ok := rawReq.([]any)
				if !ok {
					return server.Response[Detection]{}, fmt.Errorf("failed to cast %#+v to []map[string]any", rawReq)
				}

				allItems := make([]map[string]any, 0)
				for _, rawItem := range allRawItems {
					item, ok := rawItem.(map[string]any)
					if !ok {
						return server.Response[Detection]{}, fmt.Errorf("failed to cast %#+v to map[string]any", rawItem)
					}

					allItems = append(allItems, item)
				}

				forceSetValuesForFieldsByObjectIndex := make([][]string, 0)
				for _, item := range allItems {
					forceSetValuesForFields := make([]string, 0)
					for _, possibleField := range maps.Keys(item) {
						if !slices.Contains(DetectionTableColumns, possibleField) {
							continue
						}

						forceSetValuesForFields = append(forceSetValuesForFields, possibleField)
					}
					forceSetValuesForFieldsByObjectIndex = append(forceSetValuesForFieldsByObjectIndex, forceSetValuesForFields)
				}

				arguments, err := server.GetLoadArguments(ctx, queryParams.Depth)
				if err != nil {
					return server.Response[Detection]{}, err
				}

				objects, count, totalCount, _, _, err := handlePostDetections(arguments, db, waitForChange, req, forceSetValuesForFieldsByObjectIndex)
				if err != nil {
					return server.Response[Detection]{}, err
				}

				limit := int64(0)

				offset := int64(0)

				return server.Response[Detection]{
					Status:     http.StatusOK,
					Success:    true,
					Error:      nil,
					Objects:    objects,
					Count:      count,
					TotalCount: totalCount,
					Limit:      limit,
					Offset:     offset,
				}, nil
			},
			Detection{},
			DetectionIntrospectedTable,
		)
		if err != nil {
			panic(err)
		}
		r.Post(postHandler.PathWithinRouter, postHandler.ServeHTTP)
	}()

	func() {
		putHandler, err := getHTTPHandler(
			http.MethodPatch,
			"/detections/{primaryKey}",
			http.StatusOK,
			func(
				ctx context.Context,
				pathParams DetectionOnePathParams,
				queryParams DetectionLoadQueryParams,
				req Detection,
				rawReq any,
			) (server.Response[Detection], error) {
				item, ok := rawReq.(map[string]any)
				if !ok {
					return server.Response[Detection]{}, fmt.Errorf("failed to cast %#+v to map[string]any", item)
				}

				arguments, err := server.GetLoadArguments(ctx, queryParams.Depth)
				if err != nil {
					return server.Response[Detection]{}, err
				}

				object := &req
				object.ID = pathParams.PrimaryKey

				objects, count, totalCount, _, _, err := handlePutDetection(arguments, db, waitForChange, object)
				if err != nil {
					return server.Response[Detection]{}, err
				}

				limit := int64(0)

				offset := int64(0)

				return server.Response[Detection]{
					Status:     http.StatusOK,
					Success:    true,
					Error:      nil,
					Objects:    objects,
					Count:      count,
					TotalCount: totalCount,
					Limit:      limit,
					Offset:     offset,
				}, nil
			},
			Detection{},
			DetectionIntrospectedTable,
		)
		if err != nil {
			panic(err)
		}
		r.Put(putHandler.PathWithinRouter, putHandler.ServeHTTP)
	}()

	func() {
		patchHandler, err := getHTTPHandler(
			http.MethodPatch,
			"/detections/{primaryKey}",
			http.StatusOK,
			func(
				ctx context.Context,
				pathParams DetectionOnePathParams,
				queryParams DetectionLoadQueryParams,
				req Detection,
				rawReq any,
			) (server.Response[Detection], error) {
				item, ok := rawReq.(map[string]any)
				if !ok {
					return server.Response[Detection]{}, fmt.Errorf("failed to cast %#+v to map[string]any", item)
				}

				forceSetValuesForFields := make([]string, 0)
				for _, possibleField := range maps.Keys(item) {
					if !slices.Contains(DetectionTableColumns, possibleField) {
						continue
					}

					forceSetValuesForFields = append(forceSetValuesForFields, possibleField)
				}

				arguments, err := server.GetLoadArguments(ctx, queryParams.Depth)
				if err != nil {
					return server.Response[Detection]{}, err
				}

				object := &req
				object.ID = pathParams.PrimaryKey

				objects, count, totalCount, _, _, err := handlePatchDetection(arguments, db, waitForChange, object, forceSetValuesForFields)
				if err != nil {
					return server.Response[Detection]{}, err
				}

				limit := int64(0)

				offset := int64(0)

				return server.Response[Detection]{
					Status:     http.StatusOK,
					Success:    true,
					Error:      nil,
					Objects:    objects,
					Count:      count,
					TotalCount: totalCount,
					Limit:      limit,
					Offset:     offset,
				}, nil
			},
			Detection{},
			DetectionIntrospectedTable,
		)
		if err != nil {
			panic(err)
		}
		r.Patch(patchHandler.PathWithinRouter, patchHandler.ServeHTTP)
	}()

	func() {
		deleteHandler, err := getHTTPHandler(
			http.MethodDelete,
			"/detections/{primaryKey}",
			http.StatusNoContent,
			func(
				ctx context.Context,
				pathParams DetectionOnePathParams,
				queryParams DetectionLoadQueryParams,
				req server.EmptyRequest,
				rawReq any,
			) (server.EmptyResponse, error) {
				arguments, err := server.GetLoadArguments(ctx, queryParams.Depth)
				if err != nil {
					return server.EmptyResponse{}, err
				}

				object := &Detection{}
				object.ID = pathParams.PrimaryKey

				err = handleDeleteDetection(arguments, db, waitForChange, object)
				if err != nil {
					return server.EmptyResponse{}, err
				}

				return server.EmptyResponse{}, nil
			},
			Detection{},
			DetectionIntrospectedTable,
		)
		if err != nil {
			panic(err)
		}
		r.Delete(deleteHandler.PathWithinRouter, deleteHandler.ServeHTTP)
	}()

	return r
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
		Detection{},
		NewDetectionFromItem,
		"/detections",
		GetDetectionRouter,
	)
}
