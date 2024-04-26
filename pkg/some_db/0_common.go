package some_db

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
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
  "location_history": {
    "tablename": "location_history",
    "oid": "20411",
    "schema": "public",
    "reltuples": -1,
    "relkind": "r",
    "relam": "2",
    "relacl": null,
    "reltype": "20413",
    "relowner": "10",
    "relhasindex": true,
    "columns": [
      {
        "column": "id",
        "datatype": "uuid",
        "table": "location_history",
        "pos": 1,
        "typeid": "2950",
        "typelen": 16,
        "typemod": -1,
        "notnull": true,
        "hasdefault": true,
        "hasmissing": false,
        "ispkey": true,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20411",
        "zero_type": "00000000-0000-0000-0000-000000000000",
        "type_template": "uuid.UUID"
      },
      {
        "column": "created_at",
        "datatype": "timestamp with time zone",
        "table": "location_history",
        "pos": 2,
        "typeid": "1184",
        "typelen": 8,
        "typemod": -1,
        "notnull": true,
        "hasdefault": true,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20411",
        "zero_type": "0001-01-01T00:00:00Z",
        "type_template": "time.Time"
      },
      {
        "column": "updated_at",
        "datatype": "timestamp with time zone",
        "table": "location_history",
        "pos": 3,
        "typeid": "1184",
        "typelen": 8,
        "typemod": -1,
        "notnull": true,
        "hasdefault": true,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20411",
        "zero_type": "0001-01-01T00:00:00Z",
        "type_template": "time.Time"
      },
      {
        "column": "deleted_at",
        "datatype": "timestamp with time zone",
        "table": "location_history",
        "pos": 4,
        "typeid": "1184",
        "typelen": 8,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20411",
        "zero_type": "0001-01-01T00:00:00Z",
        "type_template": "*time.Time"
      },
      {
        "column": "timestamp",
        "datatype": "timestamp with time zone",
        "table": "location_history",
        "pos": 5,
        "typeid": "1184",
        "typelen": 8,
        "typemod": -1,
        "notnull": true,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20411",
        "zero_type": "0001-01-01T00:00:00Z",
        "type_template": "time.Time"
      },
      {
        "column": "point",
        "datatype": "point",
        "table": "location_history",
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
        "parent_id": "20411",
        "zero_type": {},
        "type_template": "*geom.Point"
      },
      {
        "column": "polygon",
        "datatype": "polygon",
        "table": "location_history",
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
        "parent_id": "20411",
        "zero_type": {},
        "type_template": "*geom.Polygon"
      },
      {
        "column": "parent_physical_thing_id",
        "datatype": "uuid",
        "table": "location_history",
        "pos": 8,
        "typeid": "2950",
        "typelen": 16,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": "physical_things",
        "fcolumn": "id",
        "parent_id": "20411",
        "zero_type": "00000000-0000-0000-0000-000000000000",
        "type_template": "*uuid.UUID"
      }
    ]
  },
  "logical_things": {
    "tablename": "logical_things",
    "oid": "20427",
    "schema": "public",
    "reltuples": -1,
    "relkind": "r",
    "relam": "2",
    "relacl": null,
    "reltype": "20429",
    "relowner": "10",
    "relhasindex": true,
    "columns": [
      {
        "column": "id",
        "datatype": "uuid",
        "table": "logical_things",
        "pos": 1,
        "typeid": "2950",
        "typelen": 16,
        "typemod": -1,
        "notnull": true,
        "hasdefault": true,
        "hasmissing": false,
        "ispkey": true,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20427",
        "zero_type": "00000000-0000-0000-0000-000000000000",
        "type_template": "uuid.UUID"
      },
      {
        "column": "created_at",
        "datatype": "timestamp with time zone",
        "table": "logical_things",
        "pos": 2,
        "typeid": "1184",
        "typelen": 8,
        "typemod": -1,
        "notnull": true,
        "hasdefault": true,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20427",
        "zero_type": "0001-01-01T00:00:00Z",
        "type_template": "time.Time"
      },
      {
        "column": "updated_at",
        "datatype": "timestamp with time zone",
        "table": "logical_things",
        "pos": 3,
        "typeid": "1184",
        "typelen": 8,
        "typemod": -1,
        "notnull": true,
        "hasdefault": true,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20427",
        "zero_type": "0001-01-01T00:00:00Z",
        "type_template": "time.Time"
      },
      {
        "column": "deleted_at",
        "datatype": "timestamp with time zone",
        "table": "logical_things",
        "pos": 4,
        "typeid": "1184",
        "typelen": 8,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20427",
        "zero_type": "0001-01-01T00:00:00Z",
        "type_template": "*time.Time"
      },
      {
        "column": "external_id",
        "datatype": "text",
        "table": "logical_things",
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
        "parent_id": "20427",
        "zero_type": "",
        "type_template": "*string"
      },
      {
        "column": "name",
        "datatype": "text",
        "table": "logical_things",
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
        "parent_id": "20427",
        "zero_type": "",
        "type_template": "string"
      },
      {
        "column": "type",
        "datatype": "text",
        "table": "logical_things",
        "pos": 7,
        "typeid": "25",
        "typelen": -1,
        "typemod": -1,
        "notnull": true,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20427",
        "zero_type": "",
        "type_template": "string"
      },
      {
        "column": "tags",
        "datatype": "text[]",
        "table": "logical_things",
        "pos": 8,
        "typeid": "1009",
        "typelen": -1,
        "typemod": -1,
        "notnull": true,
        "hasdefault": true,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20427",
        "zero_type": [],
        "type_template": "pq.StringArray"
      },
      {
        "column": "metadata",
        "datatype": "hstore",
        "table": "logical_things",
        "pos": 9,
        "typeid": "20265",
        "typelen": -1,
        "typemod": -1,
        "notnull": true,
        "hasdefault": true,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20427",
        "zero_type": {
          "Map": null
        },
        "type_template": "pg_types.Hstore"
      },
      {
        "column": "raw_data",
        "datatype": "jsonb",
        "table": "logical_things",
        "pos": 10,
        "typeid": "3802",
        "typelen": -1,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20427",
        "zero_type": null,
        "type_template": "any"
      },
      {
        "column": "parent_physical_thing_id",
        "datatype": "uuid",
        "table": "logical_things",
        "pos": 11,
        "typeid": "2950",
        "typelen": 16,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": "physical_things",
        "fcolumn": "id",
        "parent_id": "20427",
        "zero_type": "00000000-0000-0000-0000-000000000000",
        "type_template": "*uuid.UUID"
      },
      {
        "column": "parent_logical_thing_id",
        "datatype": "uuid",
        "table": "logical_things",
        "pos": 12,
        "typeid": "2950",
        "typelen": 16,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": "logical_things",
        "fcolumn": "id",
        "parent_id": "20427",
        "zero_type": "00000000-0000-0000-0000-000000000000",
        "type_template": "*uuid.UUID"
      }
    ]
  },
  "physical_things": {
    "tablename": "physical_things",
    "oid": "20392",
    "schema": "public",
    "reltuples": -1,
    "relkind": "r",
    "relam": "2",
    "relacl": null,
    "reltype": "20394",
    "relowner": "10",
    "relhasindex": true,
    "columns": [
      {
        "column": "id",
        "datatype": "uuid",
        "table": "physical_things",
        "pos": 1,
        "typeid": "2950",
        "typelen": 16,
        "typemod": -1,
        "notnull": true,
        "hasdefault": true,
        "hasmissing": false,
        "ispkey": true,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20392",
        "zero_type": "00000000-0000-0000-0000-000000000000",
        "type_template": "uuid.UUID"
      },
      {
        "column": "created_at",
        "datatype": "timestamp with time zone",
        "table": "physical_things",
        "pos": 2,
        "typeid": "1184",
        "typelen": 8,
        "typemod": -1,
        "notnull": true,
        "hasdefault": true,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20392",
        "zero_type": "0001-01-01T00:00:00Z",
        "type_template": "time.Time"
      },
      {
        "column": "updated_at",
        "datatype": "timestamp with time zone",
        "table": "physical_things",
        "pos": 3,
        "typeid": "1184",
        "typelen": 8,
        "typemod": -1,
        "notnull": true,
        "hasdefault": true,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20392",
        "zero_type": "0001-01-01T00:00:00Z",
        "type_template": "time.Time"
      },
      {
        "column": "deleted_at",
        "datatype": "timestamp with time zone",
        "table": "physical_things",
        "pos": 4,
        "typeid": "1184",
        "typelen": 8,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20392",
        "zero_type": "0001-01-01T00:00:00Z",
        "type_template": "*time.Time"
      },
      {
        "column": "external_id",
        "datatype": "text",
        "table": "physical_things",
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
        "parent_id": "20392",
        "zero_type": "",
        "type_template": "*string"
      },
      {
        "column": "name",
        "datatype": "text",
        "table": "physical_things",
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
        "parent_id": "20392",
        "zero_type": "",
        "type_template": "string"
      },
      {
        "column": "type",
        "datatype": "text",
        "table": "physical_things",
        "pos": 7,
        "typeid": "25",
        "typelen": -1,
        "typemod": -1,
        "notnull": true,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20392",
        "zero_type": "",
        "type_template": "string"
      },
      {
        "column": "tags",
        "datatype": "text[]",
        "table": "physical_things",
        "pos": 8,
        "typeid": "1009",
        "typelen": -1,
        "typemod": -1,
        "notnull": true,
        "hasdefault": true,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20392",
        "zero_type": [],
        "type_template": "pq.StringArray"
      },
      {
        "column": "metadata",
        "datatype": "hstore",
        "table": "physical_things",
        "pos": 9,
        "typeid": "20265",
        "typelen": -1,
        "typemod": -1,
        "notnull": true,
        "hasdefault": true,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20392",
        "zero_type": {
          "Map": null
        },
        "type_template": "pg_types.Hstore"
      },
      {
        "column": "raw_data",
        "datatype": "jsonb",
        "table": "physical_things",
        "pos": 10,
        "typeid": "3802",
        "typelen": -1,
        "typemod": -1,
        "notnull": false,
        "hasdefault": false,
        "hasmissing": false,
        "ispkey": false,
        "ftable": null,
        "fcolumn": null,
        "parent_id": "20392",
        "zero_type": null,
        "type_template": "any"
      }
    ]
  },
  "raster_columns": {
    "tablename": "raster_columns",
    "oid": "20239",
    "schema": "public",
    "reltuples": -1,
    "relkind": "v",
    "relam": "0",
    "relacl": [
      "postgres=arwdDxt/postgres",
      "=r/postgres"
    ],
    "reltype": "20241",
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
        "parent_id": "20239",
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
        "parent_id": "20239",
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
        "parent_id": "20239",
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
        "parent_id": "20239",
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
        "parent_id": "20239",
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
        "parent_id": "20239",
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
        "parent_id": "20239",
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
        "parent_id": "20239",
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
        "parent_id": "20239",
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
        "parent_id": "20239",
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
        "parent_id": "20239",
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
        "parent_id": "20239",
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
        "parent_id": "20239",
        "zero_type": [],
        "type_template": "*pq.StringArray"
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
        "parent_id": "20239",
        "zero_type": [],
        "type_template": "*pq.Float64Array"
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
        "parent_id": "20239",
        "zero_type": [],
        "type_template": "*pq.BoolArray"
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
        "parent_id": "20239",
        "zero_type": {
          "type": ""
        },
        "type_template": "*geojson.Geometry"
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
        "parent_id": "20239",
        "zero_type": false,
        "type_template": "*bool"
      }
    ]
  },
  "raster_overviews": {
    "tablename": "raster_overviews",
    "oid": "20248",
    "schema": "public",
    "reltuples": -1,
    "relkind": "v",
    "relam": "0",
    "relacl": [
      "postgres=arwdDxt/postgres",
      "=r/postgres"
    ],
    "reltype": "20250",
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
        "parent_id": "20248",
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
        "parent_id": "20248",
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
        "parent_id": "20248",
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
        "parent_id": "20248",
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
        "parent_id": "20248",
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
        "parent_id": "20248",
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
        "parent_id": "20248",
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
        "parent_id": "20248",
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
        "parent_id": "20248",
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
  }
}
`)

var TableNames = []string{"geography_columns", "geometry_columns", "location_history", "logical_things", "physical_things", "raster_columns", "raster_overviews", "spatial_ref_sys"}

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
	selectFuncByTableName["location_history"] = genericSelectLocationHistories
	insertFuncByTableName["location_history"] = genericInsertLocationHistory
	updateFuncByTableName["location_history"] = genericUpdateLocationHistory
	deleteFuncByTableName["location_history"] = genericDeleteLocationHistory
	deserializeFuncByTableName["location_history"] = DeserializeLocationHistory
	columnNamesByTableName["location_history"] = LocationHistoryColumns
	transformedColumnNamesByTableName["location_history"] = LocationHistoryTransformedColumns
	selectFuncByTableName["logical_things"] = genericSelectLogicalThings
	insertFuncByTableName["logical_things"] = genericInsertLogicalThing
	updateFuncByTableName["logical_things"] = genericUpdateLogicalThing
	deleteFuncByTableName["logical_things"] = genericDeleteLogicalThing
	deserializeFuncByTableName["logical_things"] = DeserializeLogicalThing
	columnNamesByTableName["logical_things"] = LogicalThingColumns
	transformedColumnNamesByTableName["logical_things"] = LogicalThingTransformedColumns
	selectFuncByTableName["physical_things"] = genericSelectPhysicalThings
	insertFuncByTableName["physical_things"] = genericInsertPhysicalThing
	updateFuncByTableName["physical_things"] = genericUpdatePhysicalThing
	deleteFuncByTableName["physical_things"] = genericDeletePhysicalThing
	deserializeFuncByTableName["physical_things"] = DeserializePhysicalThing
	columnNamesByTableName["physical_things"] = PhysicalThingColumns
	transformedColumnNamesByTableName["physical_things"] = PhysicalThingTransformedColumns
	selectFuncByTableName["raster_columns"] = genericSelectRasterColumnsView
	insertFuncByTableName["raster_columns"] = genericInsertRasterColumnView
	updateFuncByTableName["raster_columns"] = genericUpdateRasterColumnView
	deleteFuncByTableName["raster_columns"] = genericDeleteRasterColumnView
	deserializeFuncByTableName["raster_columns"] = DeserializeRasterColumnView
	columnNamesByTableName["raster_columns"] = RasterColumnViewColumns
	transformedColumnNamesByTableName["raster_columns"] = RasterColumnViewTransformedColumns
	selectFuncByTableName["raster_overviews"] = genericSelectRasterOverviewsView
	insertFuncByTableName["raster_overviews"] = genericInsertRasterOverviewView
	updateFuncByTableName["raster_overviews"] = genericUpdateRasterOverviewView
	deleteFuncByTableName["raster_overviews"] = genericDeleteRasterOverviewView
	deserializeFuncByTableName["raster_overviews"] = DeserializeRasterOverviewView
	columnNamesByTableName["raster_overviews"] = RasterOverviewViewColumns
	transformedColumnNamesByTableName["raster_overviews"] = RasterOverviewViewTransformedColumns
	selectFuncByTableName["spatial_ref_sys"] = genericSelectSpatialRefSys
	insertFuncByTableName["spatial_ref_sys"] = genericInsertSpatialRefSy
	updateFuncByTableName["spatial_ref_sys"] = genericUpdateSpatialRefSy
	deleteFuncByTableName["spatial_ref_sys"] = genericDeleteSpatialRefSy
	deserializeFuncByTableName["spatial_ref_sys"] = DeserializeSpatialRefSy
	columnNamesByTableName["spatial_ref_sys"] = SpatialRefSyColumns
	transformedColumnNamesByTableName["spatial_ref_sys"] = SpatialRefSyTransformedColumns
}

func SetDebug(desiredDebug bool) {
	mu.Lock()
	defer mu.Unlock()
	actualDebug = desiredDebug
	logger.Printf("runtime SetDebug() called, debugging enabled: %v", actualDebug)
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
