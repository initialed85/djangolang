package some_db

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/initialed85/djangolang/pkg/helpers"
	"github.com/initialed85/djangolang/pkg/introspect"
	"github.com/initialed85/djangolang/pkg/types"
)

var (
	dbName                            = "some_db"
	logger                            = helpers.GetLogger(fmt.Sprintf("djangolang/%v", dbName))
	mu                                = new(sync.RWMutex)
	actualDebug                       = false
	selectFuncByTableName             = make(map[string]types.SelectFunc)
	insertFuncByTableName             = make(map[string]types.InsertFunc)
	updateFuncByTableName             = make(map[string]types.UpdateFunc)
	deleteFuncByTableName             = make(map[string]types.DeleteFunc)
	deserializeFuncByTableName        = make(map[string]types.DeserializeFunc)
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
    "oid": "19729",
    "schema": "public",
    "reltuples": -1,
    "relkind": "r",
    "relam": "2",
    "relacl": null,
    "reltype": "19731",
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
        "parent_id": "19729",
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
        "parent_id": "19729",
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
        "parent_id": "19729",
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
        "parent_id": "19729",
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
        "parent_id": "19729",
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
        "parent_id": "19729",
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
        "parent_id": "19729",
        "zero_type": null,
        "type_template": "any"
      },
      {
        "column": "colour",
        "datatype": "geometry(PointZ)",
        "table": "detection",
        "pos": 8,
        "typeid": "18047",
        "typelen": -1,
        "typemod": 6,
        "notnull": true,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "19729",
        "zero_type": null,
        "type_template": "any"
      },
      {
        "column": "camera_id",
        "datatype": "bigint",
        "table": "detection",
        "pos": 9,
        "typeid": "20",
        "typelen": 8,
        "typemod": -1,
        "notnull": true,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": "camera",
        "fcolumn": "id",
        "parent_id": "19729",
        "zero_type": 0,
        "type_template": "int64"
      },
      {
        "column": "event_id",
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
        "ftable": "event",
        "fcolumn": "id",
        "parent_id": "19729",
        "zero_type": 0,
        "type_template": "*int64"
      },
      {
        "column": "object_id",
        "datatype": "bigint",
        "table": "detection",
        "pos": 11,
        "typeid": "20",
        "typelen": 8,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": "object",
        "fcolumn": "id",
        "parent_id": "19729",
        "zero_type": 0,
        "type_template": "*int64"
      }
    ]
  },
  "event": {
    "tablename": "event",
    "oid": "19741",
    "schema": "public",
    "reltuples": -1,
    "relkind": "r",
    "relam": "2",
    "relacl": null,
    "reltype": "19743",
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
        "parent_id": "19741",
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
        "parent_id": "19741",
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
        "parent_id": "19741",
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
        "parent_id": "19741",
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
        "parent_id": "19741",
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
        "parent_id": "19741",
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
        "parent_id": "19741",
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
        "parent_id": "19741",
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
        "parent_id": "19741",
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
    "oid": "19757",
    "schema": "public",
    "reltuples": -1,
    "relkind": "r",
    "relam": "2",
    "relacl": null,
    "reltype": "19759",
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
        "parent_id": "19757",
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
        "parent_id": "19757",
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
        "parent_id": "19757",
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
        "parent_id": "19757",
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
        "parent_id": "19757",
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
        "parent_id": "19757",
        "zero_type": 0,
        "type_template": "*int64"
      }
    ]
  },
  "object": {
    "tablename": "object",
    "oid": "19769",
    "schema": "public",
    "reltuples": -1,
    "relkind": "r",
    "relam": "2",
    "relacl": null,
    "reltype": "19771",
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
        "parent_id": "19769",
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
        "parent_id": "19769",
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
        "parent_id": "19769",
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
        "parent_id": "19769",
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
        "parent_id": "19769",
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
        "parent_id": "19769",
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
        "parent_id": "19769",
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
    "oid": "19783",
    "schema": "public",
    "reltuples": -1,
    "relkind": "r",
    "relam": "2",
    "relacl": null,
    "reltype": "19785",
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
        "parent_id": "19783",
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
        "parent_id": "19783",
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
        "parent_id": "19783",
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
        "parent_id": "19783",
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
        "parent_id": "19783",
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
        "parent_id": "19783",
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
        "parent_id": "19783",
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
        "parent_id": "19783",
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
	updateFuncByTableName["aggregated_detection"] = genericUpdateAggregatedDetection
	deleteFuncByTableName["aggregated_detection"] = genericDeleteAggregatedDetection
	deserializeFuncByTableName["aggregated_detection"] = DeserializeAggregatedDetection
	columnNamesByTableName["aggregated_detection"] = AggregatedDetectionColumns
	transformedColumnNamesByTableName["aggregated_detection"] = AggregatedDetectionTransformedColumns
	selectFuncByTableName["camera"] = genericSelectCameras
	insertFuncByTableName["camera"] = genericInsertCamera
	updateFuncByTableName["camera"] = genericUpdateCamera
	deleteFuncByTableName["camera"] = genericDeleteCamera
	deserializeFuncByTableName["camera"] = DeserializeCamera
	columnNamesByTableName["camera"] = CameraColumns
	transformedColumnNamesByTableName["camera"] = CameraTransformedColumns
	selectFuncByTableName["detection"] = genericSelectDetections
	insertFuncByTableName["detection"] = genericInsertDetection
	updateFuncByTableName["detection"] = genericUpdateDetection
	deleteFuncByTableName["detection"] = genericDeleteDetection
	deserializeFuncByTableName["detection"] = DeserializeDetection
	columnNamesByTableName["detection"] = DetectionColumns
	transformedColumnNamesByTableName["detection"] = DetectionTransformedColumns
	selectFuncByTableName["event"] = genericSelectEvents
	insertFuncByTableName["event"] = genericInsertEvent
	updateFuncByTableName["event"] = genericUpdateEvent
	deleteFuncByTableName["event"] = genericDeleteEvent
	deserializeFuncByTableName["event"] = DeserializeEvent
	columnNamesByTableName["event"] = EventColumns
	transformedColumnNamesByTableName["event"] = EventTransformedColumns
	selectFuncByTableName["geography_columns"] = genericSelectGeographyColumnsView
	insertFuncByTableName["geography_columns"] = genericInsertGeographyColumnView
	updateFuncByTableName["geography_columns"] = genericUpdateGeographyColumnView
	deleteFuncByTableName["geography_columns"] = genericDeleteGeographyColumnView
	deserializeFuncByTableName["geography_columns"] = DeserializeGeographyColumnView
	columnNamesByTableName["geography_columns"] = GeographyColumnViewColumns
	transformedColumnNamesByTableName["geography_columns"] = GeographyColumnViewTransformedColumns
	selectFuncByTableName["geometry_columns"] = genericSelectGeometryColumnsView
	insertFuncByTableName["geometry_columns"] = genericInsertGeometryColumnView
	updateFuncByTableName["geometry_columns"] = genericUpdateGeometryColumnView
	deleteFuncByTableName["geometry_columns"] = genericDeleteGeometryColumnView
	deserializeFuncByTableName["geometry_columns"] = DeserializeGeometryColumnView
	columnNamesByTableName["geometry_columns"] = GeometryColumnViewColumns
	transformedColumnNamesByTableName["geometry_columns"] = GeometryColumnViewTransformedColumns
	selectFuncByTableName["image"] = genericSelectImages
	insertFuncByTableName["image"] = genericInsertImage
	updateFuncByTableName["image"] = genericUpdateImage
	deleteFuncByTableName["image"] = genericDeleteImage
	deserializeFuncByTableName["image"] = DeserializeImage
	columnNamesByTableName["image"] = ImageColumns
	transformedColumnNamesByTableName["image"] = ImageTransformedColumns
	selectFuncByTableName["object"] = genericSelectObjects
	insertFuncByTableName["object"] = genericInsertObject
	updateFuncByTableName["object"] = genericUpdateObject
	deleteFuncByTableName["object"] = genericDeleteObject
	deserializeFuncByTableName["object"] = DeserializeObject
	columnNamesByTableName["object"] = ObjectColumns
	transformedColumnNamesByTableName["object"] = ObjectTransformedColumns
	selectFuncByTableName["spatial_ref_sys"] = genericSelectSpatialRefSys
	insertFuncByTableName["spatial_ref_sys"] = genericInsertSpatialRefSy
	updateFuncByTableName["spatial_ref_sys"] = genericUpdateSpatialRefSy
	deleteFuncByTableName["spatial_ref_sys"] = genericDeleteSpatialRefSy
	deserializeFuncByTableName["spatial_ref_sys"] = DeserializeSpatialRefSy
	columnNamesByTableName["spatial_ref_sys"] = SpatialRefSyColumns
	transformedColumnNamesByTableName["spatial_ref_sys"] = SpatialRefSyTransformedColumns
	selectFuncByTableName["video"] = genericSelectVideos
	insertFuncByTableName["video"] = genericInsertVideo
	updateFuncByTableName["video"] = genericUpdateVideo
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

func GetSelectFuncByTableName() map[string]types.SelectFunc {
	thisSelectFuncByTableName := make(map[string]types.SelectFunc)

	for tableName, selectFunc := range selectFuncByTableName {
		thisSelectFuncByTableName[tableName] = selectFunc
	}

	return thisSelectFuncByTableName
}

func GetInsertFuncByTableName() map[string]types.InsertFunc {
	thisInsertFuncByTableName := make(map[string]types.InsertFunc)

	for tableName, insertFunc := range insertFuncByTableName {
		thisInsertFuncByTableName[tableName] = insertFunc
	}

	return thisInsertFuncByTableName
}

func GetUpdateFuncByTableName() map[string]types.UpdateFunc {
	thisUpdateFuncByTableName := make(map[string]types.UpdateFunc)

	for tableName, updateFunc := range updateFuncByTableName {
		thisUpdateFuncByTableName[tableName] = updateFunc
	}

	return thisUpdateFuncByTableName
}

func GetDeleteFuncByTableName() map[string]types.DeleteFunc {
	thisDeleteFuncByTableName := make(map[string]types.DeleteFunc)

	for tableName, deleteFunc := range deleteFuncByTableName {
		thisDeleteFuncByTableName[tableName] = deleteFunc
	}

	return thisDeleteFuncByTableName
}
