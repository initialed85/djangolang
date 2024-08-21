package model_generated

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/netip"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/cridenour/go-postgis"
	"github.com/go-chi/chi/v5"
	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
	"github.com/initialed85/djangolang/pkg/helpers"
	"github.com/initialed85/djangolang/pkg/introspect"
	"github.com/initialed85/djangolang/pkg/query"
	"github.com/initialed85/djangolang/pkg/server"
	"github.com/initialed85/djangolang/pkg/stream"
	"github.com/initialed85/djangolang/pkg/types"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/lib/pq/hstore"
	"golang.org/x/exp/maps"
)

type Video struct {
	ID                                  uuid.UUID      `json:"id"`
	CreatedAt                           time.Time      `json:"created_at"`
	UpdatedAt                           time.Time      `json:"updated_at"`
	DeletedAt                           *time.Time     `json:"deleted_at"`
	FileName                            string         `json:"file_name"`
	StartedAt                           time.Time      `json:"started_at"`
	EndedAt                             *time.Time     `json:"ended_at"`
	Duration                            *time.Duration `json:"duration"`
	FileSize                            *float64       `json:"file_size"`
	ThumbnailName                       *string        `json:"thumbnail_name"`
	Status                              *string        `json:"status"`
	CameraID                            uuid.UUID      `json:"camera_id"`
	CameraIDObject                      *Camera        `json:"camera_id_object"`
	ReferencedByDetectionVideoIDObjects []*Detection   `json:"referenced_by_detection_video_id_objects"`
}

var VideoTable = "video"

var (
	VideoTableIDColumn            = "id"
	VideoTableCreatedAtColumn     = "created_at"
	VideoTableUpdatedAtColumn     = "updated_at"
	VideoTableDeletedAtColumn     = "deleted_at"
	VideoTableFileNameColumn      = "file_name"
	VideoTableStartedAtColumn     = "started_at"
	VideoTableEndedAtColumn       = "ended_at"
	VideoTableDurationColumn      = "duration"
	VideoTableFileSizeColumn      = "file_size"
	VideoTableThumbnailNameColumn = "thumbnail_name"
	VideoTableStatusColumn        = "status"
	VideoTableCameraIDColumn      = "camera_id"
)

var (
	VideoTableIDColumnWithTypeCast            = fmt.Sprintf(`"id" AS id`)
	VideoTableCreatedAtColumnWithTypeCast     = fmt.Sprintf(`"created_at" AS created_at`)
	VideoTableUpdatedAtColumnWithTypeCast     = fmt.Sprintf(`"updated_at" AS updated_at`)
	VideoTableDeletedAtColumnWithTypeCast     = fmt.Sprintf(`"deleted_at" AS deleted_at`)
	VideoTableFileNameColumnWithTypeCast      = fmt.Sprintf(`"file_name" AS file_name`)
	VideoTableStartedAtColumnWithTypeCast     = fmt.Sprintf(`"started_at" AS started_at`)
	VideoTableEndedAtColumnWithTypeCast       = fmt.Sprintf(`"ended_at" AS ended_at`)
	VideoTableDurationColumnWithTypeCast      = fmt.Sprintf(`"duration" AS duration`)
	VideoTableFileSizeColumnWithTypeCast      = fmt.Sprintf(`"file_size" AS file_size`)
	VideoTableThumbnailNameColumnWithTypeCast = fmt.Sprintf(`"thumbnail_name" AS thumbnail_name`)
	VideoTableStatusColumnWithTypeCast        = fmt.Sprintf(`"status" AS status`)
	VideoTableCameraIDColumnWithTypeCast      = fmt.Sprintf(`"camera_id" AS camera_id`)
)

var VideoTableColumns = []string{
	VideoTableIDColumn,
	VideoTableCreatedAtColumn,
	VideoTableUpdatedAtColumn,
	VideoTableDeletedAtColumn,
	VideoTableFileNameColumn,
	VideoTableStartedAtColumn,
	VideoTableEndedAtColumn,
	VideoTableDurationColumn,
	VideoTableFileSizeColumn,
	VideoTableThumbnailNameColumn,
	VideoTableStatusColumn,
	VideoTableCameraIDColumn,
}

var VideoTableColumnsWithTypeCasts = []string{
	VideoTableIDColumnWithTypeCast,
	VideoTableCreatedAtColumnWithTypeCast,
	VideoTableUpdatedAtColumnWithTypeCast,
	VideoTableDeletedAtColumnWithTypeCast,
	VideoTableFileNameColumnWithTypeCast,
	VideoTableStartedAtColumnWithTypeCast,
	VideoTableEndedAtColumnWithTypeCast,
	VideoTableDurationColumnWithTypeCast,
	VideoTableFileSizeColumnWithTypeCast,
	VideoTableThumbnailNameColumnWithTypeCast,
	VideoTableStatusColumnWithTypeCast,
	VideoTableCameraIDColumnWithTypeCast,
}

var VideoTableColumnLookup = map[string]*introspect.Column{
	VideoTableIDColumn:            {Name: VideoTableIDColumn, NotNull: true, HasDefault: true},
	VideoTableCreatedAtColumn:     {Name: VideoTableCreatedAtColumn, NotNull: true, HasDefault: true},
	VideoTableUpdatedAtColumn:     {Name: VideoTableUpdatedAtColumn, NotNull: true, HasDefault: true},
	VideoTableDeletedAtColumn:     {Name: VideoTableDeletedAtColumn, NotNull: false, HasDefault: false},
	VideoTableFileNameColumn:      {Name: VideoTableFileNameColumn, NotNull: true, HasDefault: false},
	VideoTableStartedAtColumn:     {Name: VideoTableStartedAtColumn, NotNull: true, HasDefault: false},
	VideoTableEndedAtColumn:       {Name: VideoTableEndedAtColumn, NotNull: false, HasDefault: false},
	VideoTableDurationColumn:      {Name: VideoTableDurationColumn, NotNull: false, HasDefault: false},
	VideoTableFileSizeColumn:      {Name: VideoTableFileSizeColumn, NotNull: false, HasDefault: false},
	VideoTableThumbnailNameColumn: {Name: VideoTableThumbnailNameColumn, NotNull: false, HasDefault: false},
	VideoTableStatusColumn:        {Name: VideoTableStatusColumn, NotNull: false, HasDefault: false},
	VideoTableCameraIDColumn:      {Name: VideoTableCameraIDColumn, NotNull: true, HasDefault: false},
}

var (
	VideoTablePrimaryKeyColumn = VideoTableIDColumn
)
var _ = []any{
	time.Time{},
	time.Duration(0),
	nil,
	pq.StringArray{},
	string(""),
	pq.Int64Array{},
	int64(0),
	pq.Float64Array{},
	float64(0),
	pq.BoolArray{},
	bool(false),
	map[string][]int{},
	uuid.UUID{},
	hstore.Hstore{},
	pgtype.Point{},
	pgtype.Polygon{},
	postgis.PointZ{},
	netip.Prefix{},
	[]byte{},
	errors.Is,
}

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

	wrapError := func(k string, v any, err error) error {
		return fmt.Errorf("%v: %#+v; error: %v", k, v, err)
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
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to time.Time", temp1))
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
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to time.Time", temp1))
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
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to time.Time", temp1))
				}
			}

			m.DeletedAt = &temp2

		case "file_name":
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
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to string", temp1))
				}
			}

			m.FileName = temp2

		case "started_at":
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
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to time.Time", temp1))
				}
			}

			m.StartedAt = temp2

		case "ended_at":
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
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to time.Time", temp1))
				}
			}

			m.EndedAt = &temp2

		case "duration":
			if v == nil {
				continue
			}

			temp1, err := types.ParseDuration(v)
			if err != nil {
				return wrapError(k, v, err)
			}

			temp2, ok := temp1.(time.Duration)
			if !ok {
				if temp1 != nil {
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to time.Duration", temp1))
				}
			}

			m.Duration = &temp2

		case "file_size":
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
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to float64", temp1))
				}
			}

			m.FileSize = &temp2

		case "thumbnail_name":
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
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to string", temp1))
				}
			}

			m.ThumbnailName = &temp2

		case "status":
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
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to string", temp1))
				}
			}

			m.Status = &temp2

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
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uuid.UUID", temp1))
				}
			}

			m.CameraID = temp2

		}
	}

	return nil
}

func (m *Video) Reload(ctx context.Context, tx *sqlx.Tx, includeDeleteds ...bool) error {
	extraWhere := ""
	if len(includeDeleteds) > 0 && includeDeleteds[0] {
		if slices.Contains(VideoTableColumns, "deleted_at") {
			extraWhere = "\n    AND (deleted_at IS null OR deleted_at IS NOT null)"
		}
	}

	ctx, cleanup := query.WithQueryID(ctx)
	defer cleanup()

	t, err := SelectVideo(
		ctx,
		tx,
		fmt.Sprintf("%v = $1%v", m.GetPrimaryKeyColumn(), extraWhere),
		m.GetPrimaryKeyValue(),
	)
	if err != nil {
		return err
	}

	m.ID = t.ID
	m.CreatedAt = t.CreatedAt
	m.UpdatedAt = t.UpdatedAt
	m.DeletedAt = t.DeletedAt
	m.FileName = t.FileName
	m.StartedAt = t.StartedAt
	m.EndedAt = t.EndedAt
	m.Duration = t.Duration
	m.FileSize = t.FileSize
	m.ThumbnailName = t.ThumbnailName
	m.Status = t.Status
	m.CameraID = t.CameraID
	m.CameraIDObject = t.CameraIDObject
	m.ReferencedByDetectionVideoIDObjects = t.ReferencedByDetectionVideoIDObjects

	return nil
}

func (m *Video) Insert(ctx context.Context, tx *sqlx.Tx, setPrimaryKey bool, setZeroValues bool, forceSetValuesForFields ...string) error {
	columns := make([]string, 0)
	values := make([]any, 0)

	if setPrimaryKey && (setZeroValues || !types.IsZeroUUID(m.ID)) || slices.Contains(forceSetValuesForFields, VideoTableIDColumn) || isRequired(VideoTableColumnLookup, VideoTableIDColumn) {
		columns = append(columns, VideoTableIDColumn)

		v, err := types.FormatUUID(m.ID)
		if err != nil {
			return fmt.Errorf("failed to handle m.ID: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.CreatedAt) || slices.Contains(forceSetValuesForFields, VideoTableCreatedAtColumn) || isRequired(VideoTableColumnLookup, VideoTableCreatedAtColumn) {
		columns = append(columns, VideoTableCreatedAtColumn)

		v, err := types.FormatTime(m.CreatedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.CreatedAt: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.UpdatedAt) || slices.Contains(forceSetValuesForFields, VideoTableUpdatedAtColumn) || isRequired(VideoTableColumnLookup, VideoTableUpdatedAtColumn) {
		columns = append(columns, VideoTableUpdatedAtColumn)

		v, err := types.FormatTime(m.UpdatedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.UpdatedAt: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.DeletedAt) || slices.Contains(forceSetValuesForFields, VideoTableDeletedAtColumn) || isRequired(VideoTableColumnLookup, VideoTableDeletedAtColumn) {
		columns = append(columns, VideoTableDeletedAtColumn)

		v, err := types.FormatTime(m.DeletedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.DeletedAt: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroString(m.FileName) || slices.Contains(forceSetValuesForFields, VideoTableFileNameColumn) || isRequired(VideoTableColumnLookup, VideoTableFileNameColumn) {
		columns = append(columns, VideoTableFileNameColumn)

		v, err := types.FormatString(m.FileName)
		if err != nil {
			return fmt.Errorf("failed to handle m.FileName: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.StartedAt) || slices.Contains(forceSetValuesForFields, VideoTableStartedAtColumn) || isRequired(VideoTableColumnLookup, VideoTableStartedAtColumn) {
		columns = append(columns, VideoTableStartedAtColumn)

		v, err := types.FormatTime(m.StartedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.StartedAt: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.EndedAt) || slices.Contains(forceSetValuesForFields, VideoTableEndedAtColumn) || isRequired(VideoTableColumnLookup, VideoTableEndedAtColumn) {
		columns = append(columns, VideoTableEndedAtColumn)

		v, err := types.FormatTime(m.EndedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.EndedAt: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroDuration(m.Duration) || slices.Contains(forceSetValuesForFields, VideoTableDurationColumn) || isRequired(VideoTableColumnLookup, VideoTableDurationColumn) {
		columns = append(columns, VideoTableDurationColumn)

		v, err := types.FormatDuration(m.Duration)
		if err != nil {
			return fmt.Errorf("failed to handle m.Duration: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroFloat(m.FileSize) || slices.Contains(forceSetValuesForFields, VideoTableFileSizeColumn) || isRequired(VideoTableColumnLookup, VideoTableFileSizeColumn) {
		columns = append(columns, VideoTableFileSizeColumn)

		v, err := types.FormatFloat(m.FileSize)
		if err != nil {
			return fmt.Errorf("failed to handle m.FileSize: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroString(m.ThumbnailName) || slices.Contains(forceSetValuesForFields, VideoTableThumbnailNameColumn) || isRequired(VideoTableColumnLookup, VideoTableThumbnailNameColumn) {
		columns = append(columns, VideoTableThumbnailNameColumn)

		v, err := types.FormatString(m.ThumbnailName)
		if err != nil {
			return fmt.Errorf("failed to handle m.ThumbnailName: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroString(m.Status) || slices.Contains(forceSetValuesForFields, VideoTableStatusColumn) || isRequired(VideoTableColumnLookup, VideoTableStatusColumn) {
		columns = append(columns, VideoTableStatusColumn)

		v, err := types.FormatString(m.Status)
		if err != nil {
			return fmt.Errorf("failed to handle m.Status: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroUUID(m.CameraID) || slices.Contains(forceSetValuesForFields, VideoTableCameraIDColumn) || isRequired(VideoTableColumnLookup, VideoTableCameraIDColumn) {
		columns = append(columns, VideoTableCameraIDColumn)

		v, err := types.FormatUUID(m.CameraID)
		if err != nil {
			return fmt.Errorf("failed to handle m.CameraID: %v", err)
		}

		values = append(values, v)
	}

	ctx, cleanup := query.WithQueryID(ctx)
	defer cleanup()

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
			"failed to treat %v: %#+v as uuid.UUID: %v",
			VideoTableIDColumn,
			item[VideoTableIDColumn],
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
		return fmt.Errorf("failed to reload after insert: %v", err)
	}

	return nil
}

func (m *Video) Update(ctx context.Context, tx *sqlx.Tx, setZeroValues bool, forceSetValuesForFields ...string) error {
	columns := make([]string, 0)
	values := make([]any, 0)

	if setZeroValues || !types.IsZeroTime(m.CreatedAt) || slices.Contains(forceSetValuesForFields, VideoTableCreatedAtColumn) {
		columns = append(columns, VideoTableCreatedAtColumn)

		v, err := types.FormatTime(m.CreatedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.CreatedAt: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.UpdatedAt) || slices.Contains(forceSetValuesForFields, VideoTableUpdatedAtColumn) {
		columns = append(columns, VideoTableUpdatedAtColumn)

		v, err := types.FormatTime(m.UpdatedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.UpdatedAt: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.DeletedAt) || slices.Contains(forceSetValuesForFields, VideoTableDeletedAtColumn) {
		columns = append(columns, VideoTableDeletedAtColumn)

		v, err := types.FormatTime(m.DeletedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.DeletedAt: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroString(m.FileName) || slices.Contains(forceSetValuesForFields, VideoTableFileNameColumn) {
		columns = append(columns, VideoTableFileNameColumn)

		v, err := types.FormatString(m.FileName)
		if err != nil {
			return fmt.Errorf("failed to handle m.FileName: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.StartedAt) || slices.Contains(forceSetValuesForFields, VideoTableStartedAtColumn) {
		columns = append(columns, VideoTableStartedAtColumn)

		v, err := types.FormatTime(m.StartedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.StartedAt: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.EndedAt) || slices.Contains(forceSetValuesForFields, VideoTableEndedAtColumn) {
		columns = append(columns, VideoTableEndedAtColumn)

		v, err := types.FormatTime(m.EndedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.EndedAt: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroDuration(m.Duration) || slices.Contains(forceSetValuesForFields, VideoTableDurationColumn) {
		columns = append(columns, VideoTableDurationColumn)

		v, err := types.FormatDuration(m.Duration)
		if err != nil {
			return fmt.Errorf("failed to handle m.Duration: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroFloat(m.FileSize) || slices.Contains(forceSetValuesForFields, VideoTableFileSizeColumn) {
		columns = append(columns, VideoTableFileSizeColumn)

		v, err := types.FormatFloat(m.FileSize)
		if err != nil {
			return fmt.Errorf("failed to handle m.FileSize: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroString(m.ThumbnailName) || slices.Contains(forceSetValuesForFields, VideoTableThumbnailNameColumn) {
		columns = append(columns, VideoTableThumbnailNameColumn)

		v, err := types.FormatString(m.ThumbnailName)
		if err != nil {
			return fmt.Errorf("failed to handle m.ThumbnailName: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroString(m.Status) || slices.Contains(forceSetValuesForFields, VideoTableStatusColumn) {
		columns = append(columns, VideoTableStatusColumn)

		v, err := types.FormatString(m.Status)
		if err != nil {
			return fmt.Errorf("failed to handle m.Status: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroUUID(m.CameraID) || slices.Contains(forceSetValuesForFields, VideoTableCameraIDColumn) {
		columns = append(columns, VideoTableCameraIDColumn)

		v, err := types.FormatUUID(m.CameraID)
		if err != nil {
			return fmt.Errorf("failed to handle m.CameraID: %v", err)
		}

		values = append(values, v)
	}

	v, err := types.FormatUUID(m.ID)
	if err != nil {
		return fmt.Errorf("failed to handle m.ID: %v", err)
	}

	values = append(values, v)

	ctx, cleanup := query.WithQueryID(ctx)
	defer cleanup()

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

	err = m.Reload(ctx, tx, slices.Contains(forceSetValuesForFields, "deleted_at"))
	if err != nil {
		return fmt.Errorf("failed to reload after update")
	}

	return nil
}

func (m *Video) Delete(ctx context.Context, tx *sqlx.Tx, hardDeletes ...bool) error {
	hardDelete := false
	if len(hardDeletes) > 0 {
		hardDelete = hardDeletes[0]
	}

	if !hardDelete && slices.Contains(VideoTableColumns, "deleted_at") {
		m.DeletedAt = helpers.Ptr(time.Now().UTC())
		err := m.Update(ctx, tx, false, "deleted_at")
		if err != nil {
			return fmt.Errorf("failed to soft-delete (update) %#+v: %v", m, err)
		}
	}

	values := make([]any, 0)
	v, err := types.FormatUUID(m.ID)
	if err != nil {
		return fmt.Errorf("failed to handle m.ID: %v", err)
	}

	values = append(values, v)

	ctx, cleanup := query.WithQueryID(ctx)
	defer cleanup()

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

	_ = m.Reload(ctx, tx, true)

	return nil
}

func SelectVideos(ctx context.Context, tx *sqlx.Tx, where string, orderBy *string, limit *int, offset *int, values ...any) ([]*Video, error) {
	if slices.Contains(VideoTableColumns, "deleted_at") {
		if !strings.Contains(where, "deleted_at") {
			if where != "" {
				where += "\n    AND "
			}

			where += "deleted_at IS null"
		}
	}

	ctx, cleanup := query.WithQueryID(ctx)
	defer cleanup()

	items, err := query.Select(
		ctx,
		tx,
		VideoTableColumnsWithTypeCasts,
		VideoTable,
		where,
		orderBy,
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
			return nil, err
		}

		thatCtx := ctx

		thatCtx, ok1 := query.HandleQueryPathGraphCycles(ctx, fmt.Sprintf("%s{%v}", VideoTable, object.ID))
		thatCtx, ok2 := query.HandleQueryPathGraphCycles(thatCtx, fmt.Sprintf("__ReferencedBy__%s{%v}", VideoTable, object.ID))
		if !(ok1 && ok2) {
			continue
		}

		_ = thatCtx

		if !types.IsZeroUUID(object.CameraID) {
			thisCtx := thatCtx
			thisCtx, ok1 := query.HandleQueryPathGraphCycles(thisCtx, fmt.Sprintf("%s{%v}", CameraTable, object.CameraID))
			thisCtx, ok2 := query.HandleQueryPathGraphCycles(thisCtx, fmt.Sprintf("__ReferencedBy__%s{%v}", CameraTable, object.CameraID))
			if ok1 && ok2 {
				object.CameraIDObject, err = SelectCamera(
					thisCtx,
					tx,
					fmt.Sprintf("%v = $1", CameraTablePrimaryKeyColumn),
					object.CameraID,
				)
				if err != nil {
					if !errors.Is(err, sql.ErrNoRows) {
						return nil, err
					}
				}
			}
		}

		err = func() error {
			thisCtx := thatCtx
			thisCtx, ok1 := query.HandleQueryPathGraphCycles(thisCtx, fmt.Sprintf("%s{%v}", VideoTable, object.ID))
			thisCtx, ok2 := query.HandleQueryPathGraphCycles(thisCtx, fmt.Sprintf("__ReferencedBy__%s{%v}", VideoTable, object.ID))

			if ok1 && ok2 {
				object.ReferencedByDetectionVideoIDObjects, err = SelectDetections(
					thisCtx,
					tx,
					fmt.Sprintf("%v = $1", DetectionTableVideoIDColumn),
					nil,
					nil,
					nil,
					object.ID,
				)
				if err != nil {
					if !errors.Is(err, sql.ErrNoRows) {
						return err
					}
				}
			}

			return nil
		}()
		if err != nil {
			return nil, err
		}

		objects = append(objects, object)
	}

	return objects, nil
}

func SelectVideo(ctx context.Context, tx *sqlx.Tx, where string, values ...any) (*Video, error) {
	ctx, cleanup := query.WithQueryID(ctx)
	defer cleanup()

	objects, err := SelectVideos(
		ctx,
		tx,
		where,
		nil,
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
		return nil, sql.ErrNoRows
	}

	object := objects[0]

	return object, nil
}

func handleGetVideos(w http.ResponseWriter, r *http.Request, db *sqlx.DB, redisPool *redis.Pool, objectMiddlewares []server.ObjectMiddleware) {
	ctx := r.Context()

	insaneOrderParams := make([]string, 0)
	hadInsaneOrderParams := false

	unrecognizedParams := make([]string, 0)
	hadUnrecognizedParams := false

	unparseableParams := make([]string, 0)
	hadUnparseableParams := false

	var orderByDirection *string
	orderBys := make([]string, 0)

	values := make([]any, 0)
	wheres := make([]string, 0)
	for rawKey, rawValues := range r.URL.Query() {
		if rawKey == "limit" || rawKey == "offset" || rawKey == "depth" {
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
				case "desc":
					if orderByDirection != nil && *orderByDirection != "DESC" {
						hadInsaneOrderParams = true
						insaneOrderParams = append(insaneOrderParams, rawKey)
						continue
					}

					orderByDirection = helpers.Ptr("DESC")
					orderBys = append(orderBys, parts[0])
					continue
				case "asc":
					if orderByDirection != nil && *orderByDirection != "ASC" {
						hadInsaneOrderParams = true
						insaneOrderParams = append(insaneOrderParams, rawKey)
						continue
					}

					orderByDirection = helpers.Ptr("ASC")
					orderBys = append(orderBys, parts[0])
					continue
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

	if hadInsaneOrderParams {
		helpers.HandleErrorResponse(
			w,
			http.StatusInternalServerError,
			fmt.Errorf("insane order params (e.g. conflicting asc / desc) %s", strings.Join(insaneOrderParams, ", ")),
		)
		return
	}

	limit := 2000
	rawLimit := r.URL.Query().Get("limit")
	if rawLimit != "" {
		possibleLimit, err := strconv.ParseInt(rawLimit, 10, 64)
		if err != nil {
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
		if err != nil {
			helpers.HandleErrorResponse(
				w,
				http.StatusInternalServerError,
				fmt.Errorf("failed to parse param offset=%s as int: %v", rawOffset, err),
			)
			return
		}

		offset = int(possibleOffset)
	}

	depth := -1
	rawDepth := r.URL.Query().Get("depth")
	if rawDepth != "" {
		possibleDepth, err := strconv.ParseInt(rawDepth, 10, 64)
		if err != nil {
			helpers.HandleErrorResponse(
				w,
				http.StatusInternalServerError,
				fmt.Errorf("failed to parse param depth=%s as int: %v", rawDepth, err),
			)
			return
		}

		depth = int(possibleDepth)

		ctx = query.WithMaxDepth(ctx, &depth)
	}

	hashableOrderBy := ""
	var orderBy *string
	if len(orderBys) > 0 {
		hashableOrderBy = strings.Join(orderBys, ", ")
		if len(orderBys) > 1 {
			hashableOrderBy = fmt.Sprintf("(%v)", hashableOrderBy)
		}
		hashableOrderBy = fmt.Sprintf("%v %v", hashableOrderBy, *orderByDirection)
		orderBy = &hashableOrderBy
	}

	requestHash, err := helpers.GetRequestHash(VideoTable, wheres, hashableOrderBy, limit, offset, depth, values, nil)
	if err != nil {
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	redisConn := redisPool.Get()
	defer func() {
		_ = redisConn.Close()
	}()

	cacheHit, err := helpers.AttemptCachedResponse(requestHash, redisConn, w)
	if err != nil {
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	if cacheHit {
		return
	}

	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	defer func() {
		_ = tx.Rollback()
	}()

	where := strings.Join(wheres, "\n    AND ")

	objects, err := SelectVideos(ctx, tx, where, orderBy, &limit, &offset, values...)
	if err != nil {
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	err = tx.Commit()
	if err != nil {
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	returnedObjectsAsJSON := helpers.HandleObjectsResponse(w, http.StatusOK, objects)

	err = helpers.StoreCachedResponse(requestHash, redisConn, string(returnedObjectsAsJSON))
	if err != nil {
		log.Printf("warning: %v", err)
	}
}

func handleGetVideo(w http.ResponseWriter, r *http.Request, db *sqlx.DB, redisPool *redis.Pool, objectMiddlewares []server.ObjectMiddleware, primaryKey string) {
	ctx := r.Context()

	wheres := []string{fmt.Sprintf("%s = $$??", VideoTablePrimaryKeyColumn)}
	values := []any{primaryKey}

	unrecognizedParams := make([]string, 0)
	hadUnrecognizedParams := false

	for rawKey, rawValues := range r.URL.Query() {
		if rawKey == "depth" {
			continue
		}

		isUnrecognized := true

		for _, rawValue := range rawValues {
			if isUnrecognized {
				unrecognizedParams = append(unrecognizedParams, fmt.Sprintf("%s=%s", rawKey, rawValue))
				hadUnrecognizedParams = true
				continue
			}

			if hadUnrecognizedParams {
				continue
			}
		}
	}

	if hadUnrecognizedParams {
		helpers.HandleErrorResponse(
			w,
			http.StatusInternalServerError,
			fmt.Errorf("unrecognized params %s", strings.Join(unrecognizedParams, ", ")),
		)
		return
	}

	depth := -1
	rawDepth := r.URL.Query().Get("depth")
	if rawDepth != "" {
		possibleDepth, err := strconv.ParseInt(rawDepth, 10, 64)
		if err != nil {
			helpers.HandleErrorResponse(
				w,
				http.StatusInternalServerError,
				fmt.Errorf("failed to parse param depth=%s as int: %v", rawDepth, err),
			)
			return
		}

		depth = int(possibleDepth)

		ctx = query.WithMaxDepth(ctx, &depth)
	}

	requestHash, err := helpers.GetRequestHash(VideoTable, wheres, "", 2, 0, depth, values, primaryKey)
	if err != nil {
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	redisConn := redisPool.Get()
	defer func() {
		_ = redisConn.Close()
	}()

	cacheHit, err := helpers.AttemptCachedResponse(requestHash, redisConn, w)
	if err != nil {
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	if cacheHit {
		return
	}

	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	defer func() {
		_ = tx.Rollback()
	}()

	where := strings.Join(wheres, "\n    AND ")

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

	returnedObjectsAsJSON := helpers.HandleObjectsResponse(w, http.StatusOK, []*Video{object})

	err = helpers.StoreCachedResponse(requestHash, redisConn, string(returnedObjectsAsJSON))
	if err != nil {
		log.Printf("warning: %v", err)
	}
}

func handlePostVideos(w http.ResponseWriter, r *http.Request, db *sqlx.DB, redisPool *redis.Pool, objectMiddlewares []server.ObjectMiddleware, waitForChange server.WaitForChange) {
	_ = redisPool

	ctx := r.Context()

	unrecognizedParams := make([]string, 0)
	hadUnrecognizedParams := false

	for rawKey, rawValues := range r.URL.Query() {
		if rawKey == "depth" {
			continue
		}

		isUnrecognized := true

		for _, rawValue := range rawValues {
			if isUnrecognized {
				unrecognizedParams = append(unrecognizedParams, fmt.Sprintf("%s=%s", rawKey, rawValue))
				hadUnrecognizedParams = true
				continue
			}

			if hadUnrecognizedParams {
				continue
			}
		}
	}

	if hadUnrecognizedParams {
		helpers.HandleErrorResponse(
			w,
			http.StatusInternalServerError,
			fmt.Errorf("unrecognized params %s", strings.Join(unrecognizedParams, ", ")),
		)
		return
	}

	depth := -1
	rawDepth := r.URL.Query().Get("depth")
	if rawDepth != "" {
		possibleDepth, err := strconv.ParseInt(rawDepth, 10, 64)
		if err != nil {
			helpers.HandleErrorResponse(
				w,
				http.StatusInternalServerError,
				fmt.Errorf("failed to parse param depth=%s as int: %v", rawDepth, err),
			)
			return
		}

		depth = int(possibleDepth)

		ctx = query.WithMaxDepth(ctx, &depth)
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		err = fmt.Errorf("failed to read body of HTTP request: %v", err)
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	var allItems []map[string]any
	err = json.Unmarshal(b, &allItems)
	if err != nil {
		err = fmt.Errorf("failed to unmarshal %#+v as JSON list of objects: %v", string(b), err)
		helpers.HandleErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	forceSetValuesForFieldsByObjectIndex := make([][]string, 0)
	objects := make([]*Video, 0)
	for _, item := range allItems {
		forceSetValuesForFields := make([]string, 0)
		for _, possibleField := range maps.Keys(item) {
			if !slices.Contains(VideoTableColumns, possibleField) {
				continue
			}

			forceSetValuesForFields = append(forceSetValuesForFields, possibleField)
		}
		forceSetValuesForFieldsByObjectIndex = append(forceSetValuesForFieldsByObjectIndex, forceSetValuesForFields)

		object := &Video{}
		err = object.FromItem(item)
		if err != nil {
			err = fmt.Errorf("failed to interpret %#+v as Video in item form: %v", item, err)
			helpers.HandleErrorResponse(w, http.StatusBadRequest, err)
			return
		}

		objects = append(objects, object)
	}

	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		err = fmt.Errorf("failed to begin DB transaction: %v", err)
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	defer func() {
		_ = tx.Rollback()
	}()

	xid, err := query.GetXid(ctx, tx)
	if err != nil {
		err = fmt.Errorf("failed to get xid: %v", err)
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}
	_ = xid

	for i, object := range objects {
		err = object.Insert(ctx, tx, false, false, forceSetValuesForFieldsByObjectIndex[i]...)
		if err != nil {
			err = fmt.Errorf("failed to insert %#+v: %v", object, err)
			helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
			return
		}

		objects[i] = object
	}

	errs := make(chan error, 1)
	go func() {
		_, err = waitForChange(ctx, []stream.Action{stream.INSERT}, VideoTable, xid)
		if err != nil {
			err = fmt.Errorf("failed to wait for change: %v", err)
			errs <- err
			return
		}

		errs <- nil
	}()

	err = tx.Commit()
	if err != nil {
		err = fmt.Errorf("failed to commit DB transaction: %v", err)
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	select {
	case <-r.Context().Done():
		err = fmt.Errorf("context canceled")
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	case err = <-errs:
		if err != nil {
			helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
			return
		}
	}

	helpers.HandleObjectsResponse(w, http.StatusCreated, objects)
}

func handlePutVideo(w http.ResponseWriter, r *http.Request, db *sqlx.DB, redisPool *redis.Pool, objectMiddlewares []server.ObjectMiddleware, waitForChange server.WaitForChange, primaryKey string) {
	_ = redisPool

	ctx := r.Context()

	unrecognizedParams := make([]string, 0)
	hadUnrecognizedParams := false

	for rawKey, rawValues := range r.URL.Query() {
		if rawKey == "depth" {
			continue
		}

		isUnrecognized := true

		for _, rawValue := range rawValues {
			if isUnrecognized {
				unrecognizedParams = append(unrecognizedParams, fmt.Sprintf("%s=%s", rawKey, rawValue))
				hadUnrecognizedParams = true
				continue
			}

			if hadUnrecognizedParams {
				continue
			}
		}
	}

	if hadUnrecognizedParams {
		helpers.HandleErrorResponse(
			w,
			http.StatusInternalServerError,
			fmt.Errorf("unrecognized params %s", strings.Join(unrecognizedParams, ", ")),
		)
		return
	}

	depth := -1
	rawDepth := r.URL.Query().Get("depth")
	if rawDepth != "" {
		possibleDepth, err := strconv.ParseInt(rawDepth, 10, 64)
		if err != nil {
			helpers.HandleErrorResponse(
				w,
				http.StatusInternalServerError,
				fmt.Errorf("failed to parse param depth=%s as int: %v", rawDepth, err),
			)
			return
		}

		depth = int(possibleDepth)

		ctx = query.WithMaxDepth(ctx, &depth)
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		err = fmt.Errorf("failed to read body of HTTP request: %v", err)
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	var item map[string]any
	err = json.Unmarshal(b, &item)
	if err != nil {
		err = fmt.Errorf("failed to unmarshal %#+v as JSON object: %v", string(b), err)
		helpers.HandleErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	item[VideoTablePrimaryKeyColumn] = primaryKey

	object := &Video{}
	err = object.FromItem(item)
	if err != nil {
		err = fmt.Errorf("failed to interpret %#+v as Video in item form: %v", item, err)
		helpers.HandleErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		err = fmt.Errorf("failed to begin DB transaction: %v", err)
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	defer func() {
		_ = tx.Rollback()
	}()

	xid, err := query.GetXid(ctx, tx)
	if err != nil {
		err = fmt.Errorf("failed to get xid: %v", err)
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}
	_ = xid

	err = object.Update(ctx, tx, true)
	if err != nil {
		err = fmt.Errorf("failed to update %#+v: %v", object, err)
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	errs := make(chan error, 1)
	go func() {
		_, err = waitForChange(ctx, []stream.Action{stream.UPDATE, stream.SOFT_DELETE, stream.SOFT_RESTORE, stream.SOFT_UPDATE}, VideoTable, xid)
		if err != nil {
			err = fmt.Errorf("failed to wait for change: %v", err)
			errs <- err
			return
		}

		errs <- nil
	}()

	err = tx.Commit()
	if err != nil {
		err = fmt.Errorf("failed to commit DB transaction: %v", err)
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	select {
	case <-r.Context().Done():
		err = fmt.Errorf("context canceled")
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	case err = <-errs:
		if err != nil {
			helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
			return
		}
	}

	helpers.HandleObjectsResponse(w, http.StatusOK, []*Video{object})
}

func handlePatchVideo(w http.ResponseWriter, r *http.Request, db *sqlx.DB, redisPool *redis.Pool, objectMiddlewares []server.ObjectMiddleware, waitForChange server.WaitForChange, primaryKey string) {
	_ = redisPool

	ctx := r.Context()

	unrecognizedParams := make([]string, 0)
	hadUnrecognizedParams := false

	for rawKey, rawValues := range r.URL.Query() {
		if rawKey == "depth" {
			continue
		}

		isUnrecognized := true

		for _, rawValue := range rawValues {
			if isUnrecognized {
				unrecognizedParams = append(unrecognizedParams, fmt.Sprintf("%s=%s", rawKey, rawValue))
				hadUnrecognizedParams = true
				continue
			}

			if hadUnrecognizedParams {
				continue
			}
		}
	}

	if hadUnrecognizedParams {
		helpers.HandleErrorResponse(
			w,
			http.StatusInternalServerError,
			fmt.Errorf("unrecognized params %s", strings.Join(unrecognizedParams, ", ")),
		)
		return
	}

	depth := -1
	rawDepth := r.URL.Query().Get("depth")
	if rawDepth != "" {
		possibleDepth, err := strconv.ParseInt(rawDepth, 10, 64)
		if err != nil {
			helpers.HandleErrorResponse(
				w,
				http.StatusInternalServerError,
				fmt.Errorf("failed to parse param depth=%s as int: %v", rawDepth, err),
			)
			return
		}

		depth = int(possibleDepth)

		ctx = query.WithMaxDepth(ctx, &depth)
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		err = fmt.Errorf("failed to read body of HTTP request: %v", err)
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	var item map[string]any
	err = json.Unmarshal(b, &item)
	if err != nil {
		err = fmt.Errorf("failed to unmarshal %#+v as JSON object: %v", string(b), err)
		helpers.HandleErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	forceSetValuesForFields := make([]string, 0)
	for _, possibleField := range maps.Keys(item) {
		if !slices.Contains(VideoTableColumns, possibleField) {
			continue
		}

		forceSetValuesForFields = append(forceSetValuesForFields, possibleField)
	}

	item[VideoTablePrimaryKeyColumn] = primaryKey

	object := &Video{}
	err = object.FromItem(item)
	if err != nil {
		err = fmt.Errorf("failed to interpret %#+v as Video in item form: %v", item, err)
		helpers.HandleErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		err = fmt.Errorf("failed to begin DB transaction: %v", err)
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	defer func() {
		_ = tx.Rollback()
	}()

	xid, err := query.GetXid(ctx, tx)
	if err != nil {
		err = fmt.Errorf("failed to get xid: %v", err)
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}
	_ = xid

	err = object.Update(ctx, tx, false, forceSetValuesForFields...)
	if err != nil {
		err = fmt.Errorf("failed to update %#+v: %v", object, err)
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	errs := make(chan error, 1)
	go func() {
		_, err = waitForChange(ctx, []stream.Action{stream.UPDATE, stream.SOFT_DELETE, stream.SOFT_RESTORE, stream.SOFT_UPDATE}, VideoTable, xid)
		if err != nil {
			err = fmt.Errorf("failed to wait for change: %v", err)
			errs <- err
			return
		}

		errs <- nil
	}()

	err = tx.Commit()
	if err != nil {
		err = fmt.Errorf("failed to commit DB transaction: %v", err)
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	select {
	case <-r.Context().Done():
		err = fmt.Errorf("context canceled")
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	case err = <-errs:
		if err != nil {
			helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
			return
		}
	}

	helpers.HandleObjectsResponse(w, http.StatusOK, []*Video{object})
}

func handleDeleteVideo(w http.ResponseWriter, r *http.Request, db *sqlx.DB, redisPool *redis.Pool, objectMiddlewares []server.ObjectMiddleware, waitForChange server.WaitForChange, primaryKey string) {
	_ = redisPool

	ctx := r.Context()

	unrecognizedParams := make([]string, 0)
	hadUnrecognizedParams := false

	for rawKey, rawValues := range r.URL.Query() {
		if rawKey == "depth" {
			continue
		}

		isUnrecognized := true

		for _, rawValue := range rawValues {
			if isUnrecognized {
				unrecognizedParams = append(unrecognizedParams, fmt.Sprintf("%s=%s", rawKey, rawValue))
				hadUnrecognizedParams = true
				continue
			}

			if hadUnrecognizedParams {
				continue
			}
		}
	}

	if hadUnrecognizedParams {
		helpers.HandleErrorResponse(
			w,
			http.StatusInternalServerError,
			fmt.Errorf("unrecognized params %s", strings.Join(unrecognizedParams, ", ")),
		)
		return
	}

	depth := -1
	rawDepth := r.URL.Query().Get("depth")
	if rawDepth != "" {
		possibleDepth, err := strconv.ParseInt(rawDepth, 10, 64)
		if err != nil {
			helpers.HandleErrorResponse(
				w,
				http.StatusInternalServerError,
				fmt.Errorf("failed to parse param depth=%s as int: %v", rawDepth, err),
			)
			return
		}

		depth = int(possibleDepth)

		ctx = query.WithMaxDepth(ctx, &depth)
	}

	var item = make(map[string]any)

	item[VideoTablePrimaryKeyColumn] = primaryKey

	object := &Video{}
	err := object.FromItem(item)
	if err != nil {
		err = fmt.Errorf("failed to interpret %#+v as Video in item form: %v", item, err)
		helpers.HandleErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		err = fmt.Errorf("failed to begin DB transaction: %v", err)
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	defer func() {
		_ = tx.Rollback()
	}()

	xid, err := query.GetXid(ctx, tx)
	if err != nil {
		err = fmt.Errorf("failed to get xid: %v", err)
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}
	_ = xid

	err = object.Delete(ctx, tx)
	if err != nil {
		err = fmt.Errorf("failed to delete %#+v: %v", object, err)
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	errs := make(chan error, 1)
	go func() {
		_, err = waitForChange(ctx, []stream.Action{stream.DELETE, stream.SOFT_DELETE}, VideoTable, xid)
		if err != nil {
			err = fmt.Errorf("failed to wait for change: %v", err)
			errs <- err
			return
		}

		errs <- nil
	}()

	err = tx.Commit()
	if err != nil {
		err = fmt.Errorf("failed to commit DB transaction: %v", err)
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	select {
	case <-r.Context().Done():
		err = fmt.Errorf("context canceled")
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	case err = <-errs:
		if err != nil {
			helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
			return
		}
	}

	helpers.HandleObjectsResponse(w, http.StatusNoContent, nil)
}

func GetVideoRouter(db *sqlx.DB, redisPool *redis.Pool, httpMiddlewares []server.HTTPMiddleware, objectMiddlewares []server.ObjectMiddleware, waitForChange server.WaitForChange) chi.Router {
	r := chi.NewRouter()

	for _, m := range httpMiddlewares {
		r.Use(m)
	}

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		handleGetVideos(w, r, db, redisPool, objectMiddlewares)
	})

	r.Get("/{primaryKey}", func(w http.ResponseWriter, r *http.Request) {
		handleGetVideo(w, r, db, redisPool, objectMiddlewares, chi.URLParam(r, "primaryKey"))
	})

	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		handlePostVideos(w, r, db, redisPool, objectMiddlewares, waitForChange)
	})

	r.Put("/{primaryKey}", func(w http.ResponseWriter, r *http.Request) {
		handlePutVideo(w, r, db, redisPool, objectMiddlewares, waitForChange, chi.URLParam(r, "primaryKey"))
	})

	r.Patch("/{primaryKey}", func(w http.ResponseWriter, r *http.Request) {
		handlePatchVideo(w, r, db, redisPool, objectMiddlewares, waitForChange, chi.URLParam(r, "primaryKey"))
	})

	r.Delete("/{primaryKey}", func(w http.ResponseWriter, r *http.Request) {
		handleDeleteVideo(w, r, db, redisPool, objectMiddlewares, waitForChange, chi.URLParam(r, "primaryKey"))
	})

	return r
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
		Video{},
		NewVideoFromItem,
		"/videos",
		GetVideoRouter,
	)
}
