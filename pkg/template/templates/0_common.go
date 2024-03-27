package templates

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/chanced/caps"
	"github.com/initialed85/djangolang/pkg/helpers"
	"github.com/initialed85/djangolang/pkg/introspect"
	"github.com/jmoiron/sqlx"
)

type SelectFunc = func(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]any, error)
type Handler = func(w http.ResponseWriter, r *http.Request)
type SelectHandler = Handler

var (
	logger                            = helpers.GetLogger("djangolang")
	mu                                = new(sync.RWMutex)
	actualDebug                       = false
	selectFuncByTableName             = make(map[string]SelectFunc)
	columnNamesByTableName            = make(map[string][]string)
	transformedColumnNamesByTableName = make(map[string][]string)
	tableByName                       = make(map[string]*introspect.Table)
)

var rawTableByName = []byte(`
{
  "aggregated_detection": {
    "tablename": "aggregated_detection",
    "oid": "62863",
    "schema": "public",
    "reltuples": 7386,
    "relkind": "r",
    "relam": "2",
    "relacl": null,
    "reltype": "62865",
    "relowner": "10",
    "relhasindex": true,
    "columns": [
      {
        "column": "id",
        "datatype": "bigint",
        "table": "aggregated_detection",
        "pos": 1,
        "typeid": "20",
        "typelen": 8,
        "typemod": -1,
        "notnull": true,
        "hasdefault": true,
        "hasmissing": false,
        "ispkey": true,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "62863",
        "zero_type": 0,
        "type_template": "int64"
      },
      {
        "column": "start_timestamp",
        "datatype": "timestamp with time zone",
        "table": "aggregated_detection",
        "pos": 2,
        "typeid": "1184",
        "typelen": 8,
        "typemod": -1,
        "notnull": true,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "62863",
        "zero_type": "0001-01-01T00:00:00Z",
        "type_template": "time.Time"
      },
      {
        "column": "end_timestamp",
        "datatype": "timestamp with time zone",
        "table": "aggregated_detection",
        "pos": 3,
        "typeid": "1184",
        "typelen": 8,
        "typemod": -1,
        "notnull": true,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "62863",
        "zero_type": "0001-01-01T00:00:00Z",
        "type_template": "time.Time"
      },
      {
        "column": "class_id",
        "datatype": "bigint",
        "table": "aggregated_detection",
        "pos": 4,
        "typeid": "20",
        "typelen": 8,
        "typemod": -1,
        "notnull": true,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "62863",
        "zero_type": 0,
        "type_template": "int64"
      },
      {
        "column": "class_name",
        "datatype": "text",
        "table": "aggregated_detection",
        "pos": 5,
        "typeid": "25",
        "typelen": -1,
        "typemod": -1,
        "notnull": true,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "62863",
        "zero_type": "",
        "type_template": "string"
      },
      {
        "column": "score",
        "datatype": "double precision",
        "table": "aggregated_detection",
        "pos": 6,
        "typeid": "701",
        "typelen": 8,
        "typemod": -1,
        "notnull": true,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "62863",
        "zero_type": 0,
        "type_template": "float64"
      },
      {
        "column": "count",
        "datatype": "bigint",
        "table": "aggregated_detection",
        "pos": 7,
        "typeid": "20",
        "typelen": 8,
        "typemod": -1,
        "notnull": true,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "62863",
        "zero_type": 0,
        "type_template": "int64"
      },
      {
        "column": "weighted_score",
        "datatype": "double precision",
        "table": "aggregated_detection",
        "pos": 8,
        "typeid": "701",
        "typelen": 8,
        "typemod": -1,
        "notnull": true,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "62863",
        "zero_type": 0,
        "type_template": "float64"
      },
      {
        "column": "event_id",
        "datatype": "bigint",
        "table": "aggregated_detection",
        "pos": 9,
        "typeid": "20",
        "typelen": 8,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": "event",
        "fcolumn": "id",
        "parent_id": "62863",
        "zero_type": 0,
        "type_template": "*int64"
      }
    ]
  },
  "camera": {
    "tablename": "camera",
    "oid": "20275",
    "schema": "public",
    "reltuples": 3,
    "relkind": "r",
    "relam": "2",
    "relacl": null,
    "reltype": "20277",
    "relowner": "10",
    "relhasindex": true,
    "columns": [
      {
        "column": "id",
        "datatype": "bigint",
        "table": "camera",
        "pos": 1,
        "typeid": "20",
        "typelen": 8,
        "typemod": -1,
        "notnull": true,
        "hasdefault": true,
        "hasmissing": false,
        "ispkey": true,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20275",
        "zero_type": 0,
        "type_template": "int64"
      },
      {
        "column": "name",
        "datatype": "text",
        "table": "camera",
        "pos": 2,
        "typeid": "25",
        "typelen": -1,
        "typemod": -1,
        "notnull": true,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20275",
        "zero_type": "",
        "type_template": "string"
      },
      {
        "column": "stream_url",
        "datatype": "text",
        "table": "camera",
        "pos": 3,
        "typeid": "25",
        "typelen": -1,
        "typemod": -1,
        "notnull": true,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20275",
        "zero_type": "",
        "type_template": "string"
      }
    ]
  },
  "cameras": {
    "tablename": "cameras",
    "oid": "20402",
    "schema": "public",
    "reltuples": -1,
    "relkind": "v",
    "relam": "0",
    "relacl": null,
    "reltype": "20404",
    "relowner": "10",
    "relhasindex": false,
    "columns": [
      {
        "column": "id",
        "datatype": "bigint",
        "table": "cameras",
        "pos": 1,
        "typeid": "20",
        "typelen": 8,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20402",
        "zero_type": 0,
        "type_template": "*int64"
      },
      {
        "column": "name",
        "datatype": "text",
        "table": "cameras",
        "pos": 2,
        "typeid": "25",
        "typelen": -1,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20402",
        "zero_type": "",
        "type_template": "*string"
      },
      {
        "column": "stream_url",
        "datatype": "text",
        "table": "cameras",
        "pos": 3,
        "typeid": "25",
        "typelen": -1,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20402",
        "zero_type": "",
        "type_template": "*string"
      }
    ]
  },
  "detection": {
    "tablename": "detection",
    "oid": "20328",
    "schema": "public",
    "reltuples": 6408489,
    "relkind": "r",
    "relam": "2",
    "relacl": null,
    "reltype": "20330",
    "relowner": "10",
    "relhasindex": true,
    "columns": [
      {
        "column": "id",
        "datatype": "bigint",
        "table": "detection",
        "pos": 1,
        "typeid": "20",
        "typelen": 8,
        "typemod": -1,
        "notnull": true,
        "hasdefault": true,
        "hasmissing": false,
        "ispkey": true,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20328",
        "zero_type": 0,
        "type_template": "int64"
      },
      {
        "column": "timestamp",
        "datatype": "timestamp with time zone",
        "table": "detection",
        "pos": 2,
        "typeid": "1184",
        "typelen": 8,
        "typemod": -1,
        "notnull": true,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20328",
        "zero_type": "0001-01-01T00:00:00Z",
        "type_template": "time.Time"
      },
      {
        "column": "class_id",
        "datatype": "bigint",
        "table": "detection",
        "pos": 3,
        "typeid": "20",
        "typelen": 8,
        "typemod": -1,
        "notnull": true,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20328",
        "zero_type": 0,
        "type_template": "int64"
      },
      {
        "column": "class_name",
        "datatype": "text",
        "table": "detection",
        "pos": 4,
        "typeid": "25",
        "typelen": -1,
        "typemod": -1,
        "notnull": true,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20328",
        "zero_type": "",
        "type_template": "string"
      },
      {
        "column": "score",
        "datatype": "double precision",
        "table": "detection",
        "pos": 5,
        "typeid": "701",
        "typelen": 8,
        "typemod": -1,
        "notnull": true,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20328",
        "zero_type": 0,
        "type_template": "float64"
      },
      {
        "column": "centroid",
        "datatype": "point",
        "table": "detection",
        "pos": 6,
        "typeid": "600",
        "typelen": 16,
        "typemod": -1,
        "notnull": true,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20328",
        "zero_type": null,
        "type_template": "any"
      },
      {
        "column": "bounding_box",
        "datatype": "polygon",
        "table": "detection",
        "pos": 7,
        "typeid": "604",
        "typelen": -1,
        "typemod": -1,
        "notnull": true,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20328",
        "zero_type": null,
        "type_template": "any"
      },
      {
        "column": "camera_id",
        "datatype": "bigint",
        "table": "detection",
        "pos": 8,
        "typeid": "20",
        "typelen": 8,
        "typemod": -1,
        "notnull": true,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": "camera",
        "fcolumn": "id",
        "parent_id": "20328",
        "zero_type": 0,
        "type_template": "int64"
      },
      {
        "column": "event_id",
        "datatype": "bigint",
        "table": "detection",
        "pos": 9,
        "typeid": "20",
        "typelen": 8,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": "event",
        "fcolumn": "id",
        "parent_id": "20328",
        "zero_type": 0,
        "type_template": "*int64"
      },
      {
        "column": "object_id",
        "datatype": "bigint",
        "table": "detection",
        "pos": 10,
        "typeid": "20",
        "typelen": 8,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": "object",
        "fcolumn": "id",
        "parent_id": "20328",
        "zero_type": 0,
        "type_template": "*int64"
      },
      {
        "column": "colour",
        "datatype": "geometry(PointZ)",
        "table": "detection",
        "pos": 12,
        "typeid": "18047",
        "typelen": -1,
        "typemod": 6,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20328",
        "zero_type": null,
        "type_template": "any"
      }
    ]
  },
  "detections": {
    "tablename": "detections",
    "oid": "27901",
    "schema": "public",
    "reltuples": -1,
    "relkind": "v",
    "relam": "0",
    "relacl": null,
    "reltype": "27903",
    "relowner": "10",
    "relhasindex": false,
    "columns": [
      {
        "column": "id",
        "datatype": "bigint",
        "table": "detections",
        "pos": 1,
        "typeid": "20",
        "typelen": 8,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "27901",
        "zero_type": 0,
        "type_template": "*int64"
      },
      {
        "column": "timestamp",
        "datatype": "timestamp with time zone",
        "table": "detections",
        "pos": 2,
        "typeid": "1184",
        "typelen": 8,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "27901",
        "zero_type": "0001-01-01T00:00:00Z",
        "type_template": "*time.Time"
      },
      {
        "column": "class_id",
        "datatype": "bigint",
        "table": "detections",
        "pos": 3,
        "typeid": "20",
        "typelen": 8,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "27901",
        "zero_type": 0,
        "type_template": "*int64"
      },
      {
        "column": "class_name",
        "datatype": "text",
        "table": "detections",
        "pos": 4,
        "typeid": "25",
        "typelen": -1,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "27901",
        "zero_type": "",
        "type_template": "*string"
      },
      {
        "column": "score",
        "datatype": "double precision",
        "table": "detections",
        "pos": 5,
        "typeid": "701",
        "typelen": 8,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "27901",
        "zero_type": 0,
        "type_template": "*float64"
      },
      {
        "column": "centroid",
        "datatype": "point",
        "table": "detections",
        "pos": 6,
        "typeid": "600",
        "typelen": 16,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "27901",
        "zero_type": null,
        "type_template": "any"
      },
      {
        "column": "bounding_box",
        "datatype": "polygon",
        "table": "detections",
        "pos": 7,
        "typeid": "604",
        "typelen": -1,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "27901",
        "zero_type": null,
        "type_template": "any"
      },
      {
        "column": "camera_id",
        "datatype": "bigint",
        "table": "detections",
        "pos": 8,
        "typeid": "20",
        "typelen": 8,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "27901",
        "zero_type": 0,
        "type_template": "*int64"
      },
      {
        "column": "event_id",
        "datatype": "bigint",
        "table": "detections",
        "pos": 9,
        "typeid": "20",
        "typelen": 8,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "27901",
        "zero_type": 0,
        "type_template": "*int64"
      },
      {
        "column": "object_id",
        "datatype": "bigint",
        "table": "detections",
        "pos": 10,
        "typeid": "20",
        "typelen": 8,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "27901",
        "zero_type": 0,
        "type_template": "*int64"
      },
      {
        "column": "colour",
        "datatype": "geometry(PointZ)",
        "table": "detections",
        "pos": 11,
        "typeid": "18047",
        "typelen": -1,
        "typemod": 6,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "27901",
        "zero_type": null,
        "type_template": "any"
      }
    ]
  },
  "event": {
    "tablename": "event",
    "oid": "20286",
    "schema": "public",
    "reltuples": 35949,
    "relkind": "r",
    "relam": "2",
    "relacl": null,
    "reltype": "20288",
    "relowner": "10",
    "relhasindex": true,
    "columns": [
      {
        "column": "id",
        "datatype": "bigint",
        "table": "event",
        "pos": 1,
        "typeid": "20",
        "typelen": 8,
        "typemod": -1,
        "notnull": true,
        "hasdefault": true,
        "hasmissing": false,
        "ispkey": true,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20286",
        "zero_type": 0,
        "type_template": "int64"
      },
      {
        "column": "start_timestamp",
        "datatype": "timestamp with time zone",
        "table": "event",
        "pos": 2,
        "typeid": "1184",
        "typelen": 8,
        "typemod": -1,
        "notnull": true,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20286",
        "zero_type": "0001-01-01T00:00:00Z",
        "type_template": "time.Time"
      },
      {
        "column": "end_timestamp",
        "datatype": "timestamp with time zone",
        "table": "event",
        "pos": 3,
        "typeid": "1184",
        "typelen": 8,
        "typemod": -1,
        "notnull": true,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20286",
        "zero_type": "0001-01-01T00:00:00Z",
        "type_template": "time.Time"
      },
      {
        "column": "duration",
        "datatype": "interval",
        "table": "event",
        "pos": 4,
        "typeid": "1186",
        "typelen": 16,
        "typemod": -1,
        "notnull": true,
        "hasdefault": true,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20286",
        "zero_type": 0,
        "type_template": "time.Duration"
      },
      {
        "column": "original_video_id",
        "datatype": "bigint",
        "table": "event",
        "pos": 5,
        "typeid": "20",
        "typelen": 8,
        "typemod": -1,
        "notnull": true,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": "video",
        "fcolumn": "id",
        "parent_id": "20286",
        "zero_type": 0,
        "type_template": "int64"
      },
      {
        "column": "thumbnail_image_id",
        "datatype": "bigint",
        "table": "event",
        "pos": 6,
        "typeid": "20",
        "typelen": 8,
        "typemod": -1,
        "notnull": true,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": "image",
        "fcolumn": "id",
        "parent_id": "20286",
        "zero_type": 0,
        "type_template": "int64"
      },
      {
        "column": "processed_video_id",
        "datatype": "bigint",
        "table": "event",
        "pos": 7,
        "typeid": "20",
        "typelen": 8,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": "video",
        "fcolumn": "id",
        "parent_id": "20286",
        "zero_type": 0,
        "type_template": "*int64"
      },
      {
        "column": "source_camera_id",
        "datatype": "bigint",
        "table": "event",
        "pos": 8,
        "typeid": "20",
        "typelen": 8,
        "typemod": -1,
        "notnull": true,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": "camera",
        "fcolumn": "id",
        "parent_id": "20286",
        "zero_type": 0,
        "type_template": "int64"
      },
      {
        "column": "status",
        "datatype": "text",
        "table": "event",
        "pos": 9,
        "typeid": "25",
        "typelen": -1,
        "typemod": -1,
        "notnull": true,
        "hasdefault": true,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20286",
        "zero_type": "",
        "type_template": "string"
      }
    ]
  },
  "event_with_detection": {
    "tablename": "event_with_detection",
    "oid": "62858",
    "schema": "public",
    "reltuples": -1,
    "relkind": "v",
    "relam": "0",
    "relacl": null,
    "reltype": "62860",
    "relowner": "10",
    "relhasindex": false,
    "columns": [
      {
        "column": "id",
        "datatype": "bigint",
        "table": "event_with_detection",
        "pos": 1,
        "typeid": "20",
        "typelen": 8,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "62858",
        "zero_type": 0,
        "type_template": "*int64"
      },
      {
        "column": "start_timestamp",
        "datatype": "timestamp with time zone",
        "table": "event_with_detection",
        "pos": 2,
        "typeid": "1184",
        "typelen": 8,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "62858",
        "zero_type": "0001-01-01T00:00:00Z",
        "type_template": "*time.Time"
      },
      {
        "column": "end_timestamp",
        "datatype": "timestamp with time zone",
        "table": "event_with_detection",
        "pos": 3,
        "typeid": "1184",
        "typelen": 8,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "62858",
        "zero_type": "0001-01-01T00:00:00Z",
        "type_template": "*time.Time"
      },
      {
        "column": "duration",
        "datatype": "interval",
        "table": "event_with_detection",
        "pos": 4,
        "typeid": "1186",
        "typelen": 16,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "62858",
        "zero_type": 0,
        "type_template": "*time.Duration"
      },
      {
        "column": "original_video_id",
        "datatype": "bigint",
        "table": "event_with_detection",
        "pos": 5,
        "typeid": "20",
        "typelen": 8,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "62858",
        "zero_type": 0,
        "type_template": "*int64"
      },
      {
        "column": "thumbnail_image_id",
        "datatype": "bigint",
        "table": "event_with_detection",
        "pos": 6,
        "typeid": "20",
        "typelen": 8,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "62858",
        "zero_type": 0,
        "type_template": "*int64"
      },
      {
        "column": "processed_video_id",
        "datatype": "bigint",
        "table": "event_with_detection",
        "pos": 7,
        "typeid": "20",
        "typelen": 8,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "62858",
        "zero_type": 0,
        "type_template": "*int64"
      },
      {
        "column": "source_camera_id",
        "datatype": "bigint",
        "table": "event_with_detection",
        "pos": 8,
        "typeid": "20",
        "typelen": 8,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "62858",
        "zero_type": 0,
        "type_template": "*int64"
      },
      {
        "column": "status",
        "datatype": "text",
        "table": "event_with_detection",
        "pos": 9,
        "typeid": "25",
        "typelen": -1,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "62858",
        "zero_type": "",
        "type_template": "*string"
      },
      {
        "column": "event_id",
        "datatype": "bigint",
        "table": "event_with_detection",
        "pos": 10,
        "typeid": "20",
        "typelen": 8,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "62858",
        "zero_type": 0,
        "type_template": "*int64"
      },
      {
        "column": "class_id",
        "datatype": "bigint",
        "table": "event_with_detection",
        "pos": 11,
        "typeid": "20",
        "typelen": 8,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "62858",
        "zero_type": 0,
        "type_template": "*int64"
      },
      {
        "column": "class_name",
        "datatype": "text",
        "table": "event_with_detection",
        "pos": 12,
        "typeid": "25",
        "typelen": -1,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "62858",
        "zero_type": "",
        "type_template": "*string"
      },
      {
        "column": "score",
        "datatype": "double precision",
        "table": "event_with_detection",
        "pos": 13,
        "typeid": "701",
        "typelen": 8,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "62858",
        "zero_type": 0,
        "type_template": "*float64"
      },
      {
        "column": "count",
        "datatype": "bigint",
        "table": "event_with_detection",
        "pos": 14,
        "typeid": "20",
        "typelen": 8,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "62858",
        "zero_type": 0,
        "type_template": "*int64"
      },
      {
        "column": "weighted_score",
        "datatype": "double precision",
        "table": "event_with_detection",
        "pos": 15,
        "typeid": "701",
        "typelen": 8,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "62858",
        "zero_type": 0,
        "type_template": "*float64"
      }
    ]
  },
  "events": {
    "tablename": "events",
    "oid": "20410",
    "schema": "public",
    "reltuples": -1,
    "relkind": "v",
    "relam": "0",
    "relacl": null,
    "reltype": "20412",
    "relowner": "10",
    "relhasindex": false,
    "columns": [
      {
        "column": "id",
        "datatype": "bigint",
        "table": "events",
        "pos": 1,
        "typeid": "20",
        "typelen": 8,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20410",
        "zero_type": 0,
        "type_template": "*int64"
      },
      {
        "column": "start_timestamp",
        "datatype": "timestamp with time zone",
        "table": "events",
        "pos": 2,
        "typeid": "1184",
        "typelen": 8,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20410",
        "zero_type": "0001-01-01T00:00:00Z",
        "type_template": "*time.Time"
      },
      {
        "column": "end_timestamp",
        "datatype": "timestamp with time zone",
        "table": "events",
        "pos": 3,
        "typeid": "1184",
        "typelen": 8,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20410",
        "zero_type": "0001-01-01T00:00:00Z",
        "type_template": "*time.Time"
      },
      {
        "column": "duration",
        "datatype": "interval",
        "table": "events",
        "pos": 4,
        "typeid": "1186",
        "typelen": 16,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20410",
        "zero_type": 0,
        "type_template": "*time.Duration"
      },
      {
        "column": "original_video_id",
        "datatype": "bigint",
        "table": "events",
        "pos": 5,
        "typeid": "20",
        "typelen": 8,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20410",
        "zero_type": 0,
        "type_template": "*int64"
      },
      {
        "column": "thumbnail_image_id",
        "datatype": "bigint",
        "table": "events",
        "pos": 6,
        "typeid": "20",
        "typelen": 8,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20410",
        "zero_type": 0,
        "type_template": "*int64"
      },
      {
        "column": "processed_video_id",
        "datatype": "bigint",
        "table": "events",
        "pos": 7,
        "typeid": "20",
        "typelen": 8,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20410",
        "zero_type": 0,
        "type_template": "*int64"
      },
      {
        "column": "source_camera_id",
        "datatype": "bigint",
        "table": "events",
        "pos": 8,
        "typeid": "20",
        "typelen": 8,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20410",
        "zero_type": 0,
        "type_template": "*int64"
      },
      {
        "column": "status",
        "datatype": "text",
        "table": "events",
        "pos": 9,
        "typeid": "25",
        "typelen": -1,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20410",
        "zero_type": "",
        "type_template": "*string"
      }
    ]
  },
  "geography_columns": {
    "tablename": "geography_columns",
    "oid": "18774",
    "schema": "public",
    "reltuples": -1,
    "relkind": "v",
    "relam": "0",
    "relacl": [
      "postgres=arwdDxt/postgres",
      "=r/postgres"
    ],
    "reltype": "18776",
    "relowner": "10",
    "relhasindex": false,
    "columns": [
      {
        "column": "f_table_catalog",
        "datatype": "name",
        "table": "geography_columns",
        "pos": 1,
        "typeid": "19",
        "typelen": 64,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "18774",
        "zero_type": "",
        "type_template": "*string"
      },
      {
        "column": "f_table_schema",
        "datatype": "name",
        "table": "geography_columns",
        "pos": 2,
        "typeid": "19",
        "typelen": 64,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "18774",
        "zero_type": "",
        "type_template": "*string"
      },
      {
        "column": "f_table_name",
        "datatype": "name",
        "table": "geography_columns",
        "pos": 3,
        "typeid": "19",
        "typelen": 64,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "18774",
        "zero_type": "",
        "type_template": "*string"
      },
      {
        "column": "f_geography_column",
        "datatype": "name",
        "table": "geography_columns",
        "pos": 4,
        "typeid": "19",
        "typelen": 64,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "18774",
        "zero_type": "",
        "type_template": "*string"
      },
      {
        "column": "coord_dimension",
        "datatype": "integer",
        "table": "geography_columns",
        "pos": 5,
        "typeid": "23",
        "typelen": 4,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "18774",
        "zero_type": 0,
        "type_template": "*int64"
      },
      {
        "column": "srid",
        "datatype": "integer",
        "table": "geography_columns",
        "pos": 6,
        "typeid": "23",
        "typelen": 4,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "18774",
        "zero_type": 0,
        "type_template": "*int64"
      },
      {
        "column": "type",
        "datatype": "text",
        "table": "geography_columns",
        "pos": 7,
        "typeid": "25",
        "typelen": -1,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "18774",
        "zero_type": "",
        "type_template": "*string"
      }
    ]
  },
  "geometry_columns": {
    "tablename": "geometry_columns",
    "oid": "18923",
    "schema": "public",
    "reltuples": -1,
    "relkind": "v",
    "relam": "0",
    "relacl": [
      "postgres=arwdDxt/postgres",
      "=r/postgres"
    ],
    "reltype": "18925",
    "relowner": "10",
    "relhasindex": false,
    "columns": [
      {
        "column": "f_table_catalog",
        "datatype": "character varying(256)",
        "table": "geometry_columns",
        "pos": 1,
        "typeid": "1043",
        "typelen": -1,
        "typemod": 260,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "18923",
        "zero_type": "",
        "type_template": "*string"
      },
      {
        "column": "f_table_schema",
        "datatype": "name",
        "table": "geometry_columns",
        "pos": 2,
        "typeid": "19",
        "typelen": 64,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "18923",
        "zero_type": "",
        "type_template": "*string"
      },
      {
        "column": "f_table_name",
        "datatype": "name",
        "table": "geometry_columns",
        "pos": 3,
        "typeid": "19",
        "typelen": 64,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "18923",
        "zero_type": "",
        "type_template": "*string"
      },
      {
        "column": "f_geometry_column",
        "datatype": "name",
        "table": "geometry_columns",
        "pos": 4,
        "typeid": "19",
        "typelen": 64,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "18923",
        "zero_type": "",
        "type_template": "*string"
      },
      {
        "column": "coord_dimension",
        "datatype": "integer",
        "table": "geometry_columns",
        "pos": 5,
        "typeid": "23",
        "typelen": 4,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "18923",
        "zero_type": 0,
        "type_template": "*int64"
      },
      {
        "column": "srid",
        "datatype": "integer",
        "table": "geometry_columns",
        "pos": 6,
        "typeid": "23",
        "typelen": 4,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "18923",
        "zero_type": 0,
        "type_template": "*int64"
      },
      {
        "column": "type",
        "datatype": "character varying(30)",
        "table": "geometry_columns",
        "pos": 7,
        "typeid": "1043",
        "typelen": -1,
        "typemod": 34,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "18923",
        "zero_type": "",
        "type_template": "*string"
      }
    ]
  },
  "image": {
    "tablename": "image",
    "oid": "20309",
    "schema": "public",
    "reltuples": 36183,
    "relkind": "r",
    "relam": "2",
    "relacl": null,
    "reltype": "20311",
    "relowner": "10",
    "relhasindex": true,
    "columns": [
      {
        "column": "id",
        "datatype": "bigint",
        "table": "image",
        "pos": 1,
        "typeid": "20",
        "typelen": 8,
        "typemod": -1,
        "notnull": true,
        "hasdefault": true,
        "hasmissing": false,
        "ispkey": true,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20309",
        "zero_type": 0,
        "type_template": "int64"
      },
      {
        "column": "timestamp",
        "datatype": "timestamp with time zone",
        "table": "image",
        "pos": 2,
        "typeid": "1184",
        "typelen": 8,
        "typemod": -1,
        "notnull": true,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20309",
        "zero_type": "0001-01-01T00:00:00Z",
        "type_template": "time.Time"
      },
      {
        "column": "size",
        "datatype": "double precision",
        "table": "image",
        "pos": 3,
        "typeid": "701",
        "typelen": 8,
        "typemod": -1,
        "notnull": true,
        "hasdefault": true,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20309",
        "zero_type": 0,
        "type_template": "float64"
      },
      {
        "column": "file_path",
        "datatype": "text",
        "table": "image",
        "pos": 4,
        "typeid": "25",
        "typelen": -1,
        "typemod": -1,
        "notnull": true,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20309",
        "zero_type": "",
        "type_template": "string"
      },
      {
        "column": "camera_id",
        "datatype": "bigint",
        "table": "image",
        "pos": 5,
        "typeid": "20",
        "typelen": 8,
        "typemod": -1,
        "notnull": true,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": "camera",
        "fcolumn": "id",
        "parent_id": "20309",
        "zero_type": 0,
        "type_template": "int64"
      },
      {
        "column": "event_id",
        "datatype": "bigint",
        "table": "image",
        "pos": 6,
        "typeid": "20",
        "typelen": 8,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": "event",
        "fcolumn": "id",
        "parent_id": "20309",
        "zero_type": 0,
        "type_template": "*int64"
      }
    ]
  },
  "images": {
    "tablename": "images",
    "oid": "20414",
    "schema": "public",
    "reltuples": -1,
    "relkind": "v",
    "relam": "0",
    "relacl": null,
    "reltype": "20416",
    "relowner": "10",
    "relhasindex": false,
    "columns": [
      {
        "column": "id",
        "datatype": "bigint",
        "table": "images",
        "pos": 1,
        "typeid": "20",
        "typelen": 8,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20414",
        "zero_type": 0,
        "type_template": "*int64"
      },
      {
        "column": "timestamp",
        "datatype": "timestamp with time zone",
        "table": "images",
        "pos": 2,
        "typeid": "1184",
        "typelen": 8,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20414",
        "zero_type": "0001-01-01T00:00:00Z",
        "type_template": "*time.Time"
      },
      {
        "column": "size",
        "datatype": "double precision",
        "table": "images",
        "pos": 3,
        "typeid": "701",
        "typelen": 8,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20414",
        "zero_type": 0,
        "type_template": "*float64"
      },
      {
        "column": "file_path",
        "datatype": "text",
        "table": "images",
        "pos": 4,
        "typeid": "25",
        "typelen": -1,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20414",
        "zero_type": "",
        "type_template": "*string"
      },
      {
        "column": "camera_id",
        "datatype": "bigint",
        "table": "images",
        "pos": 5,
        "typeid": "20",
        "typelen": 8,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20414",
        "zero_type": 0,
        "type_template": "*int64"
      },
      {
        "column": "event_id",
        "datatype": "bigint",
        "table": "images",
        "pos": 6,
        "typeid": "20",
        "typelen": 8,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20414",
        "zero_type": 0,
        "type_template": "*int64"
      }
    ]
  },
  "object": {
    "tablename": "object",
    "oid": "20319",
    "schema": "public",
    "reltuples": 0,
    "relkind": "r",
    "relam": "2",
    "relacl": null,
    "reltype": "20321",
    "relowner": "10",
    "relhasindex": true,
    "columns": [
      {
        "column": "id",
        "datatype": "bigint",
        "table": "object",
        "pos": 1,
        "typeid": "20",
        "typelen": 8,
        "typemod": -1,
        "notnull": true,
        "hasdefault": true,
        "hasmissing": false,
        "ispkey": true,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20319",
        "zero_type": 0,
        "type_template": "int64"
      },
      {
        "column": "start_timestamp",
        "datatype": "timestamp with time zone",
        "table": "object",
        "pos": 2,
        "typeid": "1184",
        "typelen": 8,
        "typemod": -1,
        "notnull": true,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20319",
        "zero_type": "0001-01-01T00:00:00Z",
        "type_template": "time.Time"
      },
      {
        "column": "end_timestamp",
        "datatype": "timestamp with time zone",
        "table": "object",
        "pos": 3,
        "typeid": "1184",
        "typelen": 8,
        "typemod": -1,
        "notnull": true,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20319",
        "zero_type": "0001-01-01T00:00:00Z",
        "type_template": "time.Time"
      },
      {
        "column": "class_id",
        "datatype": "bigint",
        "table": "object",
        "pos": 4,
        "typeid": "20",
        "typelen": 8,
        "typemod": -1,
        "notnull": true,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20319",
        "zero_type": 0,
        "type_template": "int64"
      },
      {
        "column": "class_name",
        "datatype": "text",
        "table": "object",
        "pos": 5,
        "typeid": "25",
        "typelen": -1,
        "typemod": -1,
        "notnull": true,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20319",
        "zero_type": "",
        "type_template": "string"
      },
      {
        "column": "camera_id",
        "datatype": "bigint",
        "table": "object",
        "pos": 6,
        "typeid": "20",
        "typelen": 8,
        "typemod": -1,
        "notnull": true,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": "camera",
        "fcolumn": "id",
        "parent_id": "20319",
        "zero_type": 0,
        "type_template": "int64"
      },
      {
        "column": "event_id",
        "datatype": "bigint",
        "table": "object",
        "pos": 7,
        "typeid": "20",
        "typelen": 8,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": "event",
        "fcolumn": "id",
        "parent_id": "20319",
        "zero_type": 0,
        "type_template": "*int64"
      }
    ]
  },
  "objects": {
    "tablename": "objects",
    "oid": "20418",
    "schema": "public",
    "reltuples": -1,
    "relkind": "v",
    "relam": "0",
    "relacl": null,
    "reltype": "20420",
    "relowner": "10",
    "relhasindex": false,
    "columns": [
      {
        "column": "id",
        "datatype": "bigint",
        "table": "objects",
        "pos": 1,
        "typeid": "20",
        "typelen": 8,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20418",
        "zero_type": 0,
        "type_template": "*int64"
      },
      {
        "column": "start_timestamp",
        "datatype": "timestamp with time zone",
        "table": "objects",
        "pos": 2,
        "typeid": "1184",
        "typelen": 8,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20418",
        "zero_type": "0001-01-01T00:00:00Z",
        "type_template": "*time.Time"
      },
      {
        "column": "end_timestamp",
        "datatype": "timestamp with time zone",
        "table": "objects",
        "pos": 3,
        "typeid": "1184",
        "typelen": 8,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20418",
        "zero_type": "0001-01-01T00:00:00Z",
        "type_template": "*time.Time"
      },
      {
        "column": "class_id",
        "datatype": "bigint",
        "table": "objects",
        "pos": 4,
        "typeid": "20",
        "typelen": 8,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20418",
        "zero_type": 0,
        "type_template": "*int64"
      },
      {
        "column": "class_name",
        "datatype": "text",
        "table": "objects",
        "pos": 5,
        "typeid": "25",
        "typelen": -1,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20418",
        "zero_type": "",
        "type_template": "*string"
      },
      {
        "column": "camera_id",
        "datatype": "bigint",
        "table": "objects",
        "pos": 6,
        "typeid": "20",
        "typelen": 8,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20418",
        "zero_type": 0,
        "type_template": "*int64"
      },
      {
        "column": "event_id",
        "datatype": "bigint",
        "table": "objects",
        "pos": 7,
        "typeid": "20",
        "typelen": 8,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20418",
        "zero_type": 0,
        "type_template": "*int64"
      }
    ]
  },
  "raster_columns": {
    "tablename": "raster_columns",
    "oid": "20250",
    "schema": "public",
    "reltuples": -1,
    "relkind": "v",
    "relam": "0",
    "relacl": [
      "postgres=arwdDxt/postgres",
      "=r/postgres"
    ],
    "reltype": "20252",
    "relowner": "10",
    "relhasindex": false,
    "columns": [
      {
        "column": "r_table_catalog",
        "datatype": "name",
        "table": "raster_columns",
        "pos": 1,
        "typeid": "19",
        "typelen": 64,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20250",
        "zero_type": "",
        "type_template": "*string"
      },
      {
        "column": "r_table_schema",
        "datatype": "name",
        "table": "raster_columns",
        "pos": 2,
        "typeid": "19",
        "typelen": 64,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20250",
        "zero_type": "",
        "type_template": "*string"
      },
      {
        "column": "r_table_name",
        "datatype": "name",
        "table": "raster_columns",
        "pos": 3,
        "typeid": "19",
        "typelen": 64,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20250",
        "zero_type": "",
        "type_template": "*string"
      },
      {
        "column": "r_raster_column",
        "datatype": "name",
        "table": "raster_columns",
        "pos": 4,
        "typeid": "19",
        "typelen": 64,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20250",
        "zero_type": "",
        "type_template": "*string"
      },
      {
        "column": "srid",
        "datatype": "integer",
        "table": "raster_columns",
        "pos": 5,
        "typeid": "23",
        "typelen": 4,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20250",
        "zero_type": 0,
        "type_template": "*int64"
      },
      {
        "column": "scale_x",
        "datatype": "double precision",
        "table": "raster_columns",
        "pos": 6,
        "typeid": "701",
        "typelen": 8,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20250",
        "zero_type": 0,
        "type_template": "*float64"
      },
      {
        "column": "scale_y",
        "datatype": "double precision",
        "table": "raster_columns",
        "pos": 7,
        "typeid": "701",
        "typelen": 8,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20250",
        "zero_type": 0,
        "type_template": "*float64"
      },
      {
        "column": "blocksize_x",
        "datatype": "integer",
        "table": "raster_columns",
        "pos": 8,
        "typeid": "23",
        "typelen": 4,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20250",
        "zero_type": 0,
        "type_template": "*int64"
      },
      {
        "column": "blocksize_y",
        "datatype": "integer",
        "table": "raster_columns",
        "pos": 9,
        "typeid": "23",
        "typelen": 4,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20250",
        "zero_type": 0,
        "type_template": "*int64"
      },
      {
        "column": "same_alignment",
        "datatype": "boolean",
        "table": "raster_columns",
        "pos": 10,
        "typeid": "16",
        "typelen": 1,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20250",
        "zero_type": false,
        "type_template": "*bool"
      },
      {
        "column": "regular_blocking",
        "datatype": "boolean",
        "table": "raster_columns",
        "pos": 11,
        "typeid": "16",
        "typelen": 1,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20250",
        "zero_type": false,
        "type_template": "*bool"
      },
      {
        "column": "num_bands",
        "datatype": "integer",
        "table": "raster_columns",
        "pos": 12,
        "typeid": "23",
        "typelen": 4,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20250",
        "zero_type": 0,
        "type_template": "*int64"
      },
      {
        "column": "pixel_types",
        "datatype": "text[]",
        "table": "raster_columns",
        "pos": 13,
        "typeid": "1009",
        "typelen": -1,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20250",
        "zero_type": [],
        "type_template": "*[]string"
      },
      {
        "column": "nodata_values",
        "datatype": "double precision[]",
        "table": "raster_columns",
        "pos": 14,
        "typeid": "1022",
        "typelen": -1,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20250",
        "zero_type": [],
        "type_template": "*[]float64"
      },
      {
        "column": "out_db",
        "datatype": "boolean[]",
        "table": "raster_columns",
        "pos": 15,
        "typeid": "1000",
        "typelen": -1,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20250",
        "zero_type": [],
        "type_template": "*[]bool"
      },
      {
        "column": "extent",
        "datatype": "geometry",
        "table": "raster_columns",
        "pos": 16,
        "typeid": "18047",
        "typelen": -1,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20250",
        "zero_type": null,
        "type_template": "any"
      },
      {
        "column": "spatial_index",
        "datatype": "boolean",
        "table": "raster_columns",
        "pos": 17,
        "typeid": "16",
        "typelen": 1,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20250",
        "zero_type": false,
        "type_template": "*bool"
      }
    ]
  },
  "raster_overviews": {
    "tablename": "raster_overviews",
    "oid": "20259",
    "schema": "public",
    "reltuples": -1,
    "relkind": "v",
    "relam": "0",
    "relacl": [
      "postgres=arwdDxt/postgres",
      "=r/postgres"
    ],
    "reltype": "20261",
    "relowner": "10",
    "relhasindex": false,
    "columns": [
      {
        "column": "o_table_catalog",
        "datatype": "name",
        "table": "raster_overviews",
        "pos": 1,
        "typeid": "19",
        "typelen": 64,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20259",
        "zero_type": "",
        "type_template": "*string"
      },
      {
        "column": "o_table_schema",
        "datatype": "name",
        "table": "raster_overviews",
        "pos": 2,
        "typeid": "19",
        "typelen": 64,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20259",
        "zero_type": "",
        "type_template": "*string"
      },
      {
        "column": "o_table_name",
        "datatype": "name",
        "table": "raster_overviews",
        "pos": 3,
        "typeid": "19",
        "typelen": 64,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20259",
        "zero_type": "",
        "type_template": "*string"
      },
      {
        "column": "o_raster_column",
        "datatype": "name",
        "table": "raster_overviews",
        "pos": 4,
        "typeid": "19",
        "typelen": 64,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20259",
        "zero_type": "",
        "type_template": "*string"
      },
      {
        "column": "r_table_catalog",
        "datatype": "name",
        "table": "raster_overviews",
        "pos": 5,
        "typeid": "19",
        "typelen": 64,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20259",
        "zero_type": "",
        "type_template": "*string"
      },
      {
        "column": "r_table_schema",
        "datatype": "name",
        "table": "raster_overviews",
        "pos": 6,
        "typeid": "19",
        "typelen": 64,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20259",
        "zero_type": "",
        "type_template": "*string"
      },
      {
        "column": "r_table_name",
        "datatype": "name",
        "table": "raster_overviews",
        "pos": 7,
        "typeid": "19",
        "typelen": 64,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20259",
        "zero_type": "",
        "type_template": "*string"
      },
      {
        "column": "r_raster_column",
        "datatype": "name",
        "table": "raster_overviews",
        "pos": 8,
        "typeid": "19",
        "typelen": 64,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20259",
        "zero_type": "",
        "type_template": "*string"
      },
      {
        "column": "overview_factor",
        "datatype": "integer",
        "table": "raster_overviews",
        "pos": 9,
        "typeid": "23",
        "typelen": 4,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20259",
        "zero_type": 0,
        "type_template": "*int64"
      }
    ]
  },
  "spatial_ref_sys": {
    "tablename": "spatial_ref_sys",
    "oid": "18358",
    "schema": "public",
    "reltuples": 8500,
    "relkind": "r",
    "relam": "2",
    "relacl": [
      "postgres=arwdDxt/postgres",
      "=r/postgres"
    ],
    "reltype": "18360",
    "relowner": "10",
    "relhasindex": true,
    "columns": [
      {
        "column": "srid",
        "datatype": "integer",
        "table": "spatial_ref_sys",
        "pos": 1,
        "typeid": "23",
        "typelen": 4,
        "typemod": -1,
        "notnull": true,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": true,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "18358",
        "zero_type": 0,
        "type_template": "int64"
      },
      {
        "column": "auth_name",
        "datatype": "character varying(256)",
        "table": "spatial_ref_sys",
        "pos": 2,
        "typeid": "1043",
        "typelen": -1,
        "typemod": 260,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "18358",
        "zero_type": "",
        "type_template": "*string"
      },
      {
        "column": "auth_srid",
        "datatype": "integer",
        "table": "spatial_ref_sys",
        "pos": 3,
        "typeid": "23",
        "typelen": 4,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "18358",
        "zero_type": 0,
        "type_template": "*int64"
      },
      {
        "column": "srtext",
        "datatype": "character varying(2048)",
        "table": "spatial_ref_sys",
        "pos": 4,
        "typeid": "1043",
        "typelen": -1,
        "typemod": 2052,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "18358",
        "zero_type": "",
        "type_template": "*string"
      },
      {
        "column": "proj4text",
        "datatype": "character varying(2048)",
        "table": "spatial_ref_sys",
        "pos": 5,
        "typeid": "1043",
        "typelen": -1,
        "typemod": 2052,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "18358",
        "zero_type": "",
        "type_template": "*string"
      }
    ]
  },
  "video": {
    "tablename": "video",
    "oid": "20298",
    "schema": "public",
    "reltuples": 69749,
    "relkind": "r",
    "relam": "2",
    "relacl": null,
    "reltype": "20300",
    "relowner": "10",
    "relhasindex": true,
    "columns": [
      {
        "column": "id",
        "datatype": "bigint",
        "table": "video",
        "pos": 1,
        "typeid": "20",
        "typelen": 8,
        "typemod": -1,
        "notnull": true,
        "hasdefault": true,
        "hasmissing": false,
        "ispkey": true,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20298",
        "zero_type": 0,
        "type_template": "int64"
      },
      {
        "column": "start_timestamp",
        "datatype": "timestamp with time zone",
        "table": "video",
        "pos": 2,
        "typeid": "1184",
        "typelen": 8,
        "typemod": -1,
        "notnull": true,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20298",
        "zero_type": "0001-01-01T00:00:00Z",
        "type_template": "time.Time"
      },
      {
        "column": "end_timestamp",
        "datatype": "timestamp with time zone",
        "table": "video",
        "pos": 3,
        "typeid": "1184",
        "typelen": 8,
        "typemod": -1,
        "notnull": true,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20298",
        "zero_type": "0001-01-01T00:00:00Z",
        "type_template": "time.Time"
      },
      {
        "column": "duration",
        "datatype": "interval",
        "table": "video",
        "pos": 4,
        "typeid": "1186",
        "typelen": 16,
        "typemod": -1,
        "notnull": true,
        "hasdefault": true,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20298",
        "zero_type": 0,
        "type_template": "time.Duration"
      },
      {
        "column": "size",
        "datatype": "double precision",
        "table": "video",
        "pos": 5,
        "typeid": "701",
        "typelen": 8,
        "typemod": -1,
        "notnull": true,
        "hasdefault": true,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20298",
        "zero_type": 0,
        "type_template": "float64"
      },
      {
        "column": "file_path",
        "datatype": "text",
        "table": "video",
        "pos": 6,
        "typeid": "25",
        "typelen": -1,
        "typemod": -1,
        "notnull": true,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20298",
        "zero_type": "",
        "type_template": "string"
      },
      {
        "column": "camera_id",
        "datatype": "bigint",
        "table": "video",
        "pos": 7,
        "typeid": "20",
        "typelen": 8,
        "typemod": -1,
        "notnull": true,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": "camera",
        "fcolumn": "id",
        "parent_id": "20298",
        "zero_type": 0,
        "type_template": "int64"
      },
      {
        "column": "event_id",
        "datatype": "bigint",
        "table": "video",
        "pos": 8,
        "typeid": "20",
        "typelen": 8,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": "event",
        "fcolumn": "id",
        "parent_id": "20298",
        "zero_type": 0,
        "type_template": "*int64"
      }
    ]
  },
  "videos": {
    "tablename": "videos",
    "oid": "20422",
    "schema": "public",
    "reltuples": -1,
    "relkind": "v",
    "relam": "0",
    "relacl": null,
    "reltype": "20424",
    "relowner": "10",
    "relhasindex": false,
    "columns": [
      {
        "column": "id",
        "datatype": "bigint",
        "table": "videos",
        "pos": 1,
        "typeid": "20",
        "typelen": 8,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20422",
        "zero_type": 0,
        "type_template": "*int64"
      },
      {
        "column": "start_timestamp",
        "datatype": "timestamp with time zone",
        "table": "videos",
        "pos": 2,
        "typeid": "1184",
        "typelen": 8,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20422",
        "zero_type": "0001-01-01T00:00:00Z",
        "type_template": "*time.Time"
      },
      {
        "column": "end_timestamp",
        "datatype": "timestamp with time zone",
        "table": "videos",
        "pos": 3,
        "typeid": "1184",
        "typelen": 8,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20422",
        "zero_type": "0001-01-01T00:00:00Z",
        "type_template": "*time.Time"
      },
      {
        "column": "duration",
        "datatype": "interval",
        "table": "videos",
        "pos": 4,
        "typeid": "1186",
        "typelen": 16,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20422",
        "zero_type": 0,
        "type_template": "*time.Duration"
      },
      {
        "column": "size",
        "datatype": "double precision",
        "table": "videos",
        "pos": 5,
        "typeid": "701",
        "typelen": 8,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20422",
        "zero_type": 0,
        "type_template": "*float64"
      },
      {
        "column": "file_path",
        "datatype": "text",
        "table": "videos",
        "pos": 6,
        "typeid": "25",
        "typelen": -1,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20422",
        "zero_type": "",
        "type_template": "*string"
      },
      {
        "column": "camera_id",
        "datatype": "bigint",
        "table": "videos",
        "pos": 7,
        "typeid": "20",
        "typelen": 8,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20422",
        "zero_type": 0,
        "type_template": "*int64"
      },
      {
        "column": "event_id",
        "datatype": "bigint",
        "table": "videos",
        "pos": 8,
        "typeid": "20",
        "typelen": 8,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20422",
        "zero_type": 0,
        "type_template": "*int64"
      }
    ]
  }
}
`)

func init() {
	mu.Lock()
	defer mu.Unlock()
	rawDesiredDebug := os.Getenv("DJANGOLANG_DEBUG")
	actualDebug = rawDesiredDebug == "1"
	logger.Printf("DJANGOLANG_DEBUG=%v, debugging enabled: %v", rawDesiredDebug, actualDebug)

	var err error
	tableByName, err = GetTableByName()
	if err != nil {
		log.Panic(err)
	}

	selectFuncByTableName["aggregated_detection"] = genericSelectAggregatedDetections
	columnNamesByTableName["aggregated_detection"] = AggregatedDetectionColumns
	transformedColumnNamesByTableName["aggregated_detection"] = AggregatedDetectionTransformedColumns
	selectFuncByTableName["camera"] = genericSelectCameras
	columnNamesByTableName["camera"] = CameraColumns
	transformedColumnNamesByTableName["camera"] = CameraTransformedColumns
	selectFuncByTableName["cameras"] = genericSelectCamerasView
	columnNamesByTableName["cameras"] = CameraViewColumns
	transformedColumnNamesByTableName["cameras"] = CameraViewTransformedColumns
	selectFuncByTableName["detection"] = genericSelectDetections
	columnNamesByTableName["detection"] = DetectionColumns
	transformedColumnNamesByTableName["detection"] = DetectionTransformedColumns
	selectFuncByTableName["detections"] = genericSelectDetectionsView
	columnNamesByTableName["detections"] = DetectionViewColumns
	transformedColumnNamesByTableName["detections"] = DetectionViewTransformedColumns
	selectFuncByTableName["event"] = genericSelectEvents
	columnNamesByTableName["event"] = EventColumns
	transformedColumnNamesByTableName["event"] = EventTransformedColumns
	selectFuncByTableName["event_with_detection"] = genericSelectEventWithDetectionsView
	columnNamesByTableName["event_with_detection"] = EventWithDetectionViewColumns
	transformedColumnNamesByTableName["event_with_detection"] = EventWithDetectionViewTransformedColumns
	selectFuncByTableName["events"] = genericSelectEventsView
	columnNamesByTableName["events"] = EventViewColumns
	transformedColumnNamesByTableName["events"] = EventViewTransformedColumns
	selectFuncByTableName["geography_columns"] = genericSelectGeographyColumnsView
	columnNamesByTableName["geography_columns"] = GeographyColumnViewColumns
	transformedColumnNamesByTableName["geography_columns"] = GeographyColumnViewTransformedColumns
	selectFuncByTableName["geometry_columns"] = genericSelectGeometryColumnsView
	columnNamesByTableName["geometry_columns"] = GeometryColumnViewColumns
	transformedColumnNamesByTableName["geometry_columns"] = GeometryColumnViewTransformedColumns
	selectFuncByTableName["image"] = genericSelectImages
	columnNamesByTableName["image"] = ImageColumns
	transformedColumnNamesByTableName["image"] = ImageTransformedColumns
	selectFuncByTableName["images"] = genericSelectImagesView
	columnNamesByTableName["images"] = ImageViewColumns
	transformedColumnNamesByTableName["images"] = ImageViewTransformedColumns
	selectFuncByTableName["object"] = genericSelectObjects
	columnNamesByTableName["object"] = ObjectColumns
	transformedColumnNamesByTableName["object"] = ObjectTransformedColumns
	selectFuncByTableName["objects"] = genericSelectObjectsView
	columnNamesByTableName["objects"] = ObjectViewColumns
	transformedColumnNamesByTableName["objects"] = ObjectViewTransformedColumns
	selectFuncByTableName["raster_columns"] = genericSelectRasterColumnsView
	columnNamesByTableName["raster_columns"] = RasterColumnViewColumns
	transformedColumnNamesByTableName["raster_columns"] = RasterColumnViewTransformedColumns
	selectFuncByTableName["raster_overviews"] = genericSelectRasterOverviewsView
	columnNamesByTableName["raster_overviews"] = RasterOverviewViewColumns
	transformedColumnNamesByTableName["raster_overviews"] = RasterOverviewViewTransformedColumns
	selectFuncByTableName["spatial_ref_sys"] = genericSelectSpatialRefSys
	columnNamesByTableName["spatial_ref_sys"] = SpatialRefSyColumns
	transformedColumnNamesByTableName["spatial_ref_sys"] = SpatialRefSyTransformedColumns
	selectFuncByTableName["video"] = genericSelectVideos
	columnNamesByTableName["video"] = VideoColumns
	transformedColumnNamesByTableName["video"] = VideoTransformedColumns
	selectFuncByTableName["videos"] = genericSelectVideosView
	columnNamesByTableName["videos"] = VideoViewColumns
	transformedColumnNamesByTableName["videos"] = VideoViewTransformedColumns
}

func SetDebug(desiredDebug bool) {
	mu.Lock()
	defer mu.Unlock()
	actualDebug = desiredDebug
	logger.Printf("runtime SetDebug() called, debugging enabled: %v", actualDebug)
}

func Descending(columns ...string) *string {
	return helpers.Ptr(
		fmt.Sprintf(
			"(%v) DESC",
			strings.Join(columns, ", "),
		),
	)
}

func Ascending(columns ...string) *string {
	return helpers.Ptr(
		fmt.Sprintf(
			"(%v) ASC",
			strings.Join(columns, ", "),
		),
	)
}

func Columns(includeColumns []string, excludeColumns ...string) []string {
	excludeColumnLookup := make(map[string]bool)
	for _, column := range excludeColumns {
		excludeColumnLookup[column] = true
	}

	columns := make([]string, 0)
	for _, column := range includeColumns {
		_, ok := excludeColumnLookup[column]
		if ok {
			continue
		}

		columns = append(columns, column)
	}

	return columns
}

func GetRawTableByName() []byte {
	return rawTableByName
}

func GetTableByName() (map[string]*introspect.Table, error) {
	thisTableByName := make(map[string]*introspect.Table)

	err := json.Unmarshal(rawTableByName, &thisTableByName)
	if err != nil {
		return nil, err
	}

	thisTableByName, err = introspect.MapTableByName(thisTableByName)
	if err != nil {
		return nil, err
	}

	return thisTableByName, nil
}

func GetSelectFuncByTableName() map[string]SelectFunc {
	thisSelectFuncByTableName := make(map[string]SelectFunc)

	for tableName, selectFunc := range selectFuncByTableName {
		thisSelectFuncByTableName[tableName] = selectFunc
	}

	return thisSelectFuncByTableName
}

func GetSelectHandlerForTableName(tableName string, db *sqlx.DB) (SelectHandler, error) {
	selectFunc, ok := selectFuncByTableName[tableName]
	if !ok {
		return nil, fmt.Errorf("no selectFuncByTableName entry for tableName %v", tableName)
	}

	columns, ok := transformedColumnNamesByTableName[tableName]
	if !ok {
		return nil, fmt.Errorf("no columnNamesByTableName entry for tableName %v", tableName)
	}

	table, ok := tableByName[tableName]
	if !ok {
		return nil, fmt.Errorf("no tableByName entry for tableName %v", tableName)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		status := http.StatusOK
		var body []byte
		var err error

		defer func() {
			w.Header().Add("Content-Type", "application/json")

			w.WriteHeader(status)

			if err != nil {
				_, _ = w.Write([]byte(fmt.Sprintf("{\"error\": %#+v}", err.Error())))
				return
			}

			_, _ = w.Write(body)
		}()

		if r.Method != http.MethodGet {
			status = http.StatusMethodNotAllowed
			err = fmt.Errorf("%v; wanted %v, got %v", "StatusMethodNotAllowed", http.MethodGet, r.Method)
			return
		}

		limit := helpers.Ptr(50)
		rawLimit := strings.TrimSpace(r.URL.Query().Get("limit"))
		if len(rawLimit) > 0 {
			var parsedLimit int64
			parsedLimit, err = strconv.ParseInt(rawLimit, 10, 64)
			if err != nil {
				status = http.StatusBadRequest
				err = fmt.Errorf("failed to parse limit %#+v as int; err: %v", limit, err.Error())
				return
			}
			limit = helpers.Ptr(int(parsedLimit))
		}

		offset := helpers.Ptr(50)
		rawOffset := strings.TrimSpace(r.URL.Query().Get("offset"))
		if len(rawOffset) > 0 {
			var parsedOffset int64
			parsedOffset, err = strconv.ParseInt(rawOffset, 10, 64)
			if err != nil {
				status = http.StatusBadRequest
				err = fmt.Errorf("failed to parse offset %#+v as int; err: %v", offset, err.Error())
				return
			}
			offset = helpers.Ptr(int(parsedOffset))
		}

		descending := strings.HasPrefix(strings.ToLower(strings.TrimSpace(r.URL.Query().Get("order"))), "desc")

		orderBy := make([]string, 0)
		for k, vs := range r.URL.Query() {
			if k != "order_by" {
				continue
			}

			for _, v := range vs {
				_, ok := table.ColumnByName[v]
				if !ok {
					status = http.StatusBadRequest
					err = fmt.Errorf("bad order by; no table.ColumnByName entry for columnName %v", v)
					return
				}

				found := false

				for _, existingOrderBy := range orderBy {
					if existingOrderBy == v {
						found = true
						break
					}
				}

				if found {
					continue
				}

				orderBy = append(orderBy, v)
			}
		}

		var order *string
		if len(orderBy) > 0 {
			if descending {
				order = Descending(orderBy...)
			} else {
				order = Ascending(orderBy...)
			}
		}

		wheres := make([]string, 0)
		for k, vs := range r.URL.Query() {
			if k == "order" || k == "order_by" || k == "limit" || k == "offset" {
				continue
			}

			field := k
			matcher := "="

			parts := strings.Split(k, "__")
			if len(parts) > 1 {
				field = strings.Join(parts[0:len(parts)-1], "__")
				lastPart := strings.ToLower(parts[len(parts)-1])
				if lastPart == "eq" {
					matcher = "="
				} else if lastPart == "ne" {
					matcher = "!="
				} else if lastPart == "gt" {
					matcher = ">"
				} else if lastPart == "lt" {
					matcher = "<"
				} else if lastPart == "gte" {
					matcher = ">="
				} else if lastPart == "lte" {
					matcher = "<="
				} else if lastPart == "ilike" {
					matcher = "ILIKE"
				} else if lastPart == "in" {
					matcher = "IN"
				} else if lastPart == "nin" || lastPart == "not_in" {
					matcher = "NOT IN"
				}
			}

			innerWheres := make([]string, 0)

			for _, v := range vs {
				column, ok := table.ColumnByName[field]
				if !ok {
					status = http.StatusBadRequest
					err = fmt.Errorf("bad filter; no table.ColumnByName entry for columnName %v", field)
					return
				}

				if matcher != "IN" && matcher != "NOT IN" {
					if column.TypeTemplate == "string" || column.TypeTemplate == "uuid" {
						v = fmt.Sprintf("%#+v", v)
						innerWheres = append(
							innerWheres,
							fmt.Sprintf("%v %v '%v'", field, matcher, v[1:len(v)-1]),
						)
					} else {
						innerWheres = append(
							innerWheres,
							fmt.Sprintf("%v %v %v", field, matcher, v),
						)
					}
				} else {
					itemsBefore := strings.Split(v, ",")
					itemsAfter := make([]string, 0)

					for _, x := range itemsBefore {
						x = strings.TrimSpace(x)
						if x == "" {
							continue
						}

						if column.TypeTemplate == "string" || column.TypeTemplate == "uuid" {
							x = fmt.Sprintf("%#+v", x)
							itemsAfter = append(itemsAfter, fmt.Sprintf("'%v'", x[1:len(x)-1]))
						} else {
							itemsAfter = append(itemsAfter, x)
						}
					}

					innerWheres = append(
						innerWheres,
						fmt.Sprintf("%v %v (%v)", field, matcher, strings.Join(itemsAfter, ", ")),
					)
				}
			}

			wheres = append(wheres, fmt.Sprintf("(%v)", strings.Join(innerWheres, " OR ")))
		}

		var items []any
		items, err = selectFunc(r.Context(), db, columns, order, limit, offset, wheres...)
		if err != nil {
			status = http.StatusInternalServerError
			err = fmt.Errorf("query failed; err: %v", err.Error())
			return
		}

		body, err = json.Marshal(items)
		if err != nil {
			status = http.StatusInternalServerError
			err = fmt.Errorf("serialization failed; err: %v", err.Error())
		}
	}, nil
}

func GetSelectHandlerByEndpointName(db *sqlx.DB) (map[string]SelectHandler, error) {
	thisSelectHandlerByEndpointName := make(map[string]SelectHandler)

	for _, tableName := range TableNames {
		endpointName := caps.ToKebab[string](tableName)

		selectHandler, err := GetSelectHandlerForTableName(tableName, db)
		if err != nil {
			return nil, err
		}

		thisSelectHandlerByEndpointName[endpointName] = selectHandler
	}

	return thisSelectHandlerByEndpointName, nil
}

var TableNames = []string{"aggregated_detection", "camera", "cameras", "detection", "detections", "event", "event_with_detection", "events", "geography_columns", "geometry_columns", "image", "images", "object", "objects", "raster_columns", "raster_overviews", "spatial_ref_sys", "video", "videos"}
