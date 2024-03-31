package some_db

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/initialed85/djangolang/pkg/helpers"
	"github.com/initialed85/djangolang/pkg/introspect"
	"github.com/jmoiron/sqlx"
)

type DjangolangObject interface {
	Insert(ctx context.Context, db *sqlx.DB, columns ...string) error
	Delete(ctx context.Context, db *sqlx.DB) error
}

type SelectFunc = func(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]DjangolangObject, error)
type InsertFunc = func(ctx context.Context, db *sqlx.DB, object DjangolangObject, columns ...string) (DjangolangObject, error)
type DeleteFunc = func(ctx context.Context, db *sqlx.DB, object DjangolangObject) error
type DeserializeFunc = func(b []byte) (DjangolangObject, error)
type Handler = func(w http.ResponseWriter, r *http.Request)
type SelectHandler = Handler
type InsertHandler = Handler
type DeleteHandler = Handler

var (
	logger                            = helpers.GetLogger("djangolang/some_db")
	mu                                = new(sync.RWMutex)
	actualDebug                       = false
	selectFuncByTableName             = make(map[string]SelectFunc)
	insertFuncByTableName             = make(map[string]InsertFunc)
	deleteFuncByTableName             = make(map[string]DeleteFunc)
	deserializeFuncByTableName        = make(map[string]DeserializeFunc)
	columnNamesByTableName            = make(map[string][]string)
	transformedColumnNamesByTableName = make(map[string][]string)
	tableByName                       = make(map[string]*introspect.Table)
)

var rawTableByName = []byte(`
{
  "aggregated_detection": {
    "tablename": "aggregated_detection",
    "oid": "19708",
    "schema": "public",
    "reltuples": -1,
    "relkind": "r",
    "relam": "2",
    "relacl": null,
    "reltype": "19710",
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
        "parent_id": "19708",
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
        "parent_id": "19708",
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
        "parent_id": "19708",
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
        "parent_id": "19708",
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
        "parent_id": "19708",
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
        "parent_id": "19708",
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
        "parent_id": "19708",
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
        "parent_id": "19708",
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
        "parent_id": "19708",
        "zero_type": 0,
        "type_template": "*int64"
      }
    ]
  },
  "camera": {
    "tablename": "camera",
    "oid": "19717",
    "schema": "public",
    "reltuples": -1,
    "relkind": "r",
    "relam": "2",
    "relacl": null,
    "reltype": "19719",
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
        "parent_id": "19717",
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
        "parent_id": "19717",
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
        "parent_id": "19717",
        "zero_type": "",
        "type_template": "string"
      }
    ]
  },
  "detection": {
    "tablename": "detection",
    "oid": "19727",
    "schema": "public",
    "reltuples": -1,
    "relkind": "r",
    "relam": "2",
    "relacl": null,
    "reltype": "19729",
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
        "parent_id": "19727",
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
        "parent_id": "19727",
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
        "parent_id": "19727",
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
        "parent_id": "19727",
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
        "parent_id": "19727",
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
        "parent_id": "19727",
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
        "parent_id": "19727",
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
        "parent_id": "19727",
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
        "parent_id": "19727",
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
        "parent_id": "19727",
        "zero_type": 0,
        "type_template": "*int64"
      },
      {
        "column": "colour",
        "datatype": "geometry(PointZ)",
        "table": "detection",
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
        "parent_id": "19727",
        "zero_type": null,
        "type_template": "any"
      }
    ]
  },
  "event": {
    "tablename": "event",
    "oid": "19739",
    "schema": "public",
    "reltuples": -1,
    "relkind": "r",
    "relam": "2",
    "relacl": null,
    "reltype": "19741",
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
        "parent_id": "19739",
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
        "parent_id": "19739",
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
        "parent_id": "19739",
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
        "parent_id": "19739",
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
        "parent_id": "19739",
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
        "parent_id": "19739",
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
        "parent_id": "19739",
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
        "parent_id": "19739",
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
        "parent_id": "19739",
        "zero_type": "",
        "type_template": "string"
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
    "oid": "19755",
    "schema": "public",
    "reltuples": -1,
    "relkind": "r",
    "relam": "2",
    "relacl": null,
    "reltype": "19757",
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
        "parent_id": "19755",
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
        "parent_id": "19755",
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
        "parent_id": "19755",
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
        "parent_id": "19755",
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
        "parent_id": "19755",
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
        "parent_id": "19755",
        "zero_type": 0,
        "type_template": "*int64"
      }
    ]
  },
  "object": {
    "tablename": "object",
    "oid": "19767",
    "schema": "public",
    "reltuples": -1,
    "relkind": "r",
    "relam": "2",
    "relacl": null,
    "reltype": "19769",
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
        "parent_id": "19767",
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
        "parent_id": "19767",
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
        "parent_id": "19767",
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
        "parent_id": "19767",
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
        "parent_id": "19767",
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
        "parent_id": "19767",
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
        "parent_id": "19767",
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
    "oid": "19781",
    "schema": "public",
    "reltuples": -1,
    "relkind": "r",
    "relam": "2",
    "relacl": null,
    "reltype": "19783",
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
        "parent_id": "19781",
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
        "parent_id": "19781",
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
        "parent_id": "19781",
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
        "parent_id": "19781",
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
        "parent_id": "19781",
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
        "parent_id": "19781",
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
        "parent_id": "19781",
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
        "parent_id": "19781",
        "zero_type": 0,
        "type_template": "*int64"
      }
    ]
  }
}
`)

var TableNames = []string{"aggregated_detection", "camera", "detection", "event", "geography_columns", "geometry_columns", "image", "object", "spatial_ref_sys", "video"}

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
	insertFuncByTableName["aggregated_detection"] = genericInsertAggregatedDetection
	deleteFuncByTableName["aggregated_detection"] = genericDeleteAggregatedDetection
	deserializeFuncByTableName["aggregated_detection"] = DeserializeAggregatedDetection
	columnNamesByTableName["aggregated_detection"] = AggregatedDetectionColumns
	transformedColumnNamesByTableName["aggregated_detection"] = AggregatedDetectionTransformedColumns
	selectFuncByTableName["camera"] = genericSelectCameras
	insertFuncByTableName["camera"] = genericInsertCamera
	deleteFuncByTableName["camera"] = genericDeleteCamera
	deserializeFuncByTableName["camera"] = DeserializeCamera
	columnNamesByTableName["camera"] = CameraColumns
	transformedColumnNamesByTableName["camera"] = CameraTransformedColumns
	selectFuncByTableName["detection"] = genericSelectDetections
	insertFuncByTableName["detection"] = genericInsertDetection
	deleteFuncByTableName["detection"] = genericDeleteDetection
	deserializeFuncByTableName["detection"] = DeserializeDetection
	columnNamesByTableName["detection"] = DetectionColumns
	transformedColumnNamesByTableName["detection"] = DetectionTransformedColumns
	selectFuncByTableName["event"] = genericSelectEvents
	insertFuncByTableName["event"] = genericInsertEvent
	deleteFuncByTableName["event"] = genericDeleteEvent
	deserializeFuncByTableName["event"] = DeserializeEvent
	columnNamesByTableName["event"] = EventColumns
	transformedColumnNamesByTableName["event"] = EventTransformedColumns
	selectFuncByTableName["geography_columns"] = genericSelectGeographyColumnsView
	insertFuncByTableName["geography_columns"] = genericInsertGeographyColumnView
	deleteFuncByTableName["geography_columns"] = genericDeleteGeographyColumnView
	deserializeFuncByTableName["geography_columns"] = DeserializeGeographyColumnView
	columnNamesByTableName["geography_columns"] = GeographyColumnViewColumns
	transformedColumnNamesByTableName["geography_columns"] = GeographyColumnViewTransformedColumns
	selectFuncByTableName["geometry_columns"] = genericSelectGeometryColumnsView
	insertFuncByTableName["geometry_columns"] = genericInsertGeometryColumnView
	deleteFuncByTableName["geometry_columns"] = genericDeleteGeometryColumnView
	deserializeFuncByTableName["geometry_columns"] = DeserializeGeometryColumnView
	columnNamesByTableName["geometry_columns"] = GeometryColumnViewColumns
	transformedColumnNamesByTableName["geometry_columns"] = GeometryColumnViewTransformedColumns
	selectFuncByTableName["image"] = genericSelectImages
	insertFuncByTableName["image"] = genericInsertImage
	deleteFuncByTableName["image"] = genericDeleteImage
	deserializeFuncByTableName["image"] = DeserializeImage
	columnNamesByTableName["image"] = ImageColumns
	transformedColumnNamesByTableName["image"] = ImageTransformedColumns
	selectFuncByTableName["object"] = genericSelectObjects
	insertFuncByTableName["object"] = genericInsertObject
	deleteFuncByTableName["object"] = genericDeleteObject
	deserializeFuncByTableName["object"] = DeserializeObject
	columnNamesByTableName["object"] = ObjectColumns
	transformedColumnNamesByTableName["object"] = ObjectTransformedColumns
	selectFuncByTableName["spatial_ref_sys"] = genericSelectSpatialRefSys
	insertFuncByTableName["spatial_ref_sys"] = genericInsertSpatialRefSy
	deleteFuncByTableName["spatial_ref_sys"] = genericDeleteSpatialRefSy
	deserializeFuncByTableName["spatial_ref_sys"] = DeserializeSpatialRefSy
	columnNamesByTableName["spatial_ref_sys"] = SpatialRefSyColumns
	transformedColumnNamesByTableName["spatial_ref_sys"] = SpatialRefSyTransformedColumns
	selectFuncByTableName["video"] = genericSelectVideos
	insertFuncByTableName["video"] = genericInsertVideo
	deleteFuncByTableName["video"] = genericDeleteVideo
	deserializeFuncByTableName["video"] = DeserializeVideo
	columnNamesByTableName["video"] = VideoColumns
	transformedColumnNamesByTableName["video"] = VideoTransformedColumns
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

func GetInsertFuncByTableName() map[string]InsertFunc {
	thisInsertFuncByTableName := make(map[string]InsertFunc)

	for tableName, insertFunc := range insertFuncByTableName {
		thisInsertFuncByTableName[tableName] = insertFunc
	}

	return thisInsertFuncByTableName
}

func GetDeleteFuncByTableName() map[string]DeleteFunc {
	thisDeleteFuncByTableName := make(map[string]DeleteFunc)

	for tableName, deleteFunc := range deleteFuncByTableName {
		thisDeleteFuncByTableName[tableName] = deleteFunc
	}

	return thisDeleteFuncByTableName
}
