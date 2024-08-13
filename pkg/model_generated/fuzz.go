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

type Fuzz struct {
	ID       uuid.UUID           `json:"id"`
	Column1  *time.Time          `json:"column1"`
	Column2  *time.Time          `json:"column2"`
	Column3  any                 `json:"column3"`
	Column4  any                 `json:"column4"`
	Column5  *[]string           `json:"column5"`
	Column6  *[]string           `json:"column6"`
	Column7  *string             `json:"column7"`
	Column8  *string             `json:"column8"`
	Column9  *pq.Int64Array      `json:"column9"`
	Column10 *pq.Int64Array      `json:"column10"`
	Column11 *pq.Int64Array      `json:"column11"`
	Column12 *int64              `json:"column12"`
	Column13 *int64              `json:"column13"`
	Column14 *int64              `json:"column14"`
	Column15 *pq.Float64Array    `json:"column15"`
	Column16 *pq.Float64Array    `json:"column16"`
	Column17 *pq.Float64Array    `json:"column17"`
	Column18 *pq.Float64Array    `json:"column18"`
	Column19 *float64            `json:"column19"`
	Column20 *float64            `json:"column20"`
	Column21 *float64            `json:"column21"`
	Column22 *float64            `json:"column22"`
	Column23 *pq.BoolArray       `json:"column23"`
	Column24 *bool               `json:"column24"`
	Column25 *map[string][]int   `json:"column25"`
	Column26 *uuid.UUID          `json:"column26"`
	Column27 *map[string]*string `json:"column27"`
	Column28 *pgtype.Point       `json:"column28"`
	Column29 *pgtype.Polygon     `json:"column29"`
	Column30 *postgis.PointZ     `json:"column30"`
	Column31 *postgis.PointZ     `json:"column31"`
	Column32 *netip.Prefix       `json:"column32"`
	Column33 *[]byte             `json:"column33"`
}

var FuzzTable = "fuzz"

var (
	FuzzTableIDColumn       = "id"
	FuzzTableColumn1Column  = "column1"
	FuzzTableColumn2Column  = "column2"
	FuzzTableColumn3Column  = "column3"
	FuzzTableColumn4Column  = "column4"
	FuzzTableColumn5Column  = "column5"
	FuzzTableColumn6Column  = "column6"
	FuzzTableColumn7Column  = "column7"
	FuzzTableColumn8Column  = "column8"
	FuzzTableColumn9Column  = "column9"
	FuzzTableColumn10Column = "column10"
	FuzzTableColumn11Column = "column11"
	FuzzTableColumn12Column = "column12"
	FuzzTableColumn13Column = "column13"
	FuzzTableColumn14Column = "column14"
	FuzzTableColumn15Column = "column15"
	FuzzTableColumn16Column = "column16"
	FuzzTableColumn17Column = "column17"
	FuzzTableColumn18Column = "column18"
	FuzzTableColumn19Column = "column19"
	FuzzTableColumn20Column = "column20"
	FuzzTableColumn21Column = "column21"
	FuzzTableColumn22Column = "column22"
	FuzzTableColumn23Column = "column23"
	FuzzTableColumn24Column = "column24"
	FuzzTableColumn25Column = "column25"
	FuzzTableColumn26Column = "column26"
	FuzzTableColumn27Column = "column27"
	FuzzTableColumn28Column = "column28"
	FuzzTableColumn29Column = "column29"
	FuzzTableColumn30Column = "column30"
	FuzzTableColumn31Column = "column31"
	FuzzTableColumn32Column = "column32"
	FuzzTableColumn33Column = "column33"
)

var (
	FuzzTableIDColumnWithTypeCast       = fmt.Sprintf(`"id" AS id`)
	FuzzTableColumn1ColumnWithTypeCast  = fmt.Sprintf(`"column1" AS column1`)
	FuzzTableColumn2ColumnWithTypeCast  = fmt.Sprintf(`"column2" AS column2`)
	FuzzTableColumn3ColumnWithTypeCast  = fmt.Sprintf(`"column3" AS column3`)
	FuzzTableColumn4ColumnWithTypeCast  = fmt.Sprintf(`"column4" AS column4`)
	FuzzTableColumn5ColumnWithTypeCast  = fmt.Sprintf(`"column5" AS column5`)
	FuzzTableColumn6ColumnWithTypeCast  = fmt.Sprintf(`"column6" AS column6`)
	FuzzTableColumn7ColumnWithTypeCast  = fmt.Sprintf(`"column7" AS column7`)
	FuzzTableColumn8ColumnWithTypeCast  = fmt.Sprintf(`"column8" AS column8`)
	FuzzTableColumn9ColumnWithTypeCast  = fmt.Sprintf(`"column9" AS column9`)
	FuzzTableColumn10ColumnWithTypeCast = fmt.Sprintf(`"column10" AS column10`)
	FuzzTableColumn11ColumnWithTypeCast = fmt.Sprintf(`"column11" AS column11`)
	FuzzTableColumn12ColumnWithTypeCast = fmt.Sprintf(`"column12" AS column12`)
	FuzzTableColumn13ColumnWithTypeCast = fmt.Sprintf(`"column13" AS column13`)
	FuzzTableColumn14ColumnWithTypeCast = fmt.Sprintf(`"column14" AS column14`)
	FuzzTableColumn15ColumnWithTypeCast = fmt.Sprintf(`"column15" AS column15`)
	FuzzTableColumn16ColumnWithTypeCast = fmt.Sprintf(`"column16" AS column16`)
	FuzzTableColumn17ColumnWithTypeCast = fmt.Sprintf(`"column17" AS column17`)
	FuzzTableColumn18ColumnWithTypeCast = fmt.Sprintf(`"column18" AS column18`)
	FuzzTableColumn19ColumnWithTypeCast = fmt.Sprintf(`"column19" AS column19`)
	FuzzTableColumn20ColumnWithTypeCast = fmt.Sprintf(`"column20" AS column20`)
	FuzzTableColumn21ColumnWithTypeCast = fmt.Sprintf(`"column21" AS column21`)
	FuzzTableColumn22ColumnWithTypeCast = fmt.Sprintf(`"column22" AS column22`)
	FuzzTableColumn23ColumnWithTypeCast = fmt.Sprintf(`"column23" AS column23`)
	FuzzTableColumn24ColumnWithTypeCast = fmt.Sprintf(`"column24" AS column24`)
	FuzzTableColumn25ColumnWithTypeCast = fmt.Sprintf(`"column25" AS column25`)
	FuzzTableColumn26ColumnWithTypeCast = fmt.Sprintf(`"column26" AS column26`)
	FuzzTableColumn27ColumnWithTypeCast = fmt.Sprintf(`"column27" AS column27`)
	FuzzTableColumn28ColumnWithTypeCast = fmt.Sprintf(`"column28" AS column28`)
	FuzzTableColumn29ColumnWithTypeCast = fmt.Sprintf(`"column29" AS column29`)
	FuzzTableColumn30ColumnWithTypeCast = fmt.Sprintf(`"column30" AS column30`)
	FuzzTableColumn31ColumnWithTypeCast = fmt.Sprintf(`"column31" AS column31`)
	FuzzTableColumn32ColumnWithTypeCast = fmt.Sprintf(`"column32" AS column32`)
	FuzzTableColumn33ColumnWithTypeCast = fmt.Sprintf(`"column33" AS column33`)
)

var FuzzTableColumns = []string{
	FuzzTableIDColumn,
	FuzzTableColumn1Column,
	FuzzTableColumn2Column,
	FuzzTableColumn3Column,
	FuzzTableColumn4Column,
	FuzzTableColumn5Column,
	FuzzTableColumn6Column,
	FuzzTableColumn7Column,
	FuzzTableColumn8Column,
	FuzzTableColumn9Column,
	FuzzTableColumn10Column,
	FuzzTableColumn11Column,
	FuzzTableColumn12Column,
	FuzzTableColumn13Column,
	FuzzTableColumn14Column,
	FuzzTableColumn15Column,
	FuzzTableColumn16Column,
	FuzzTableColumn17Column,
	FuzzTableColumn18Column,
	FuzzTableColumn19Column,
	FuzzTableColumn20Column,
	FuzzTableColumn21Column,
	FuzzTableColumn22Column,
	FuzzTableColumn23Column,
	FuzzTableColumn24Column,
	FuzzTableColumn25Column,
	FuzzTableColumn26Column,
	FuzzTableColumn27Column,
	FuzzTableColumn28Column,
	FuzzTableColumn29Column,
	FuzzTableColumn30Column,
	FuzzTableColumn31Column,
	FuzzTableColumn32Column,
	FuzzTableColumn33Column,
}

var FuzzTableColumnsWithTypeCasts = []string{
	FuzzTableIDColumnWithTypeCast,
	FuzzTableColumn1ColumnWithTypeCast,
	FuzzTableColumn2ColumnWithTypeCast,
	FuzzTableColumn3ColumnWithTypeCast,
	FuzzTableColumn4ColumnWithTypeCast,
	FuzzTableColumn5ColumnWithTypeCast,
	FuzzTableColumn6ColumnWithTypeCast,
	FuzzTableColumn7ColumnWithTypeCast,
	FuzzTableColumn8ColumnWithTypeCast,
	FuzzTableColumn9ColumnWithTypeCast,
	FuzzTableColumn10ColumnWithTypeCast,
	FuzzTableColumn11ColumnWithTypeCast,
	FuzzTableColumn12ColumnWithTypeCast,
	FuzzTableColumn13ColumnWithTypeCast,
	FuzzTableColumn14ColumnWithTypeCast,
	FuzzTableColumn15ColumnWithTypeCast,
	FuzzTableColumn16ColumnWithTypeCast,
	FuzzTableColumn17ColumnWithTypeCast,
	FuzzTableColumn18ColumnWithTypeCast,
	FuzzTableColumn19ColumnWithTypeCast,
	FuzzTableColumn20ColumnWithTypeCast,
	FuzzTableColumn21ColumnWithTypeCast,
	FuzzTableColumn22ColumnWithTypeCast,
	FuzzTableColumn23ColumnWithTypeCast,
	FuzzTableColumn24ColumnWithTypeCast,
	FuzzTableColumn25ColumnWithTypeCast,
	FuzzTableColumn26ColumnWithTypeCast,
	FuzzTableColumn27ColumnWithTypeCast,
	FuzzTableColumn28ColumnWithTypeCast,
	FuzzTableColumn29ColumnWithTypeCast,
	FuzzTableColumn30ColumnWithTypeCast,
	FuzzTableColumn31ColumnWithTypeCast,
	FuzzTableColumn32ColumnWithTypeCast,
	FuzzTableColumn33ColumnWithTypeCast,
}

var FuzzTableColumnLookup = map[string]*introspect.Column{
	FuzzTableIDColumn:       new(introspect.Column),
	FuzzTableColumn1Column:  new(introspect.Column),
	FuzzTableColumn2Column:  new(introspect.Column),
	FuzzTableColumn3Column:  new(introspect.Column),
	FuzzTableColumn4Column:  new(introspect.Column),
	FuzzTableColumn5Column:  new(introspect.Column),
	FuzzTableColumn6Column:  new(introspect.Column),
	FuzzTableColumn7Column:  new(introspect.Column),
	FuzzTableColumn8Column:  new(introspect.Column),
	FuzzTableColumn9Column:  new(introspect.Column),
	FuzzTableColumn10Column: new(introspect.Column),
	FuzzTableColumn11Column: new(introspect.Column),
	FuzzTableColumn12Column: new(introspect.Column),
	FuzzTableColumn13Column: new(introspect.Column),
	FuzzTableColumn14Column: new(introspect.Column),
	FuzzTableColumn15Column: new(introspect.Column),
	FuzzTableColumn16Column: new(introspect.Column),
	FuzzTableColumn17Column: new(introspect.Column),
	FuzzTableColumn18Column: new(introspect.Column),
	FuzzTableColumn19Column: new(introspect.Column),
	FuzzTableColumn20Column: new(introspect.Column),
	FuzzTableColumn21Column: new(introspect.Column),
	FuzzTableColumn22Column: new(introspect.Column),
	FuzzTableColumn23Column: new(introspect.Column),
	FuzzTableColumn24Column: new(introspect.Column),
	FuzzTableColumn25Column: new(introspect.Column),
	FuzzTableColumn26Column: new(introspect.Column),
	FuzzTableColumn27Column: new(introspect.Column),
	FuzzTableColumn28Column: new(introspect.Column),
	FuzzTableColumn29Column: new(introspect.Column),
	FuzzTableColumn30Column: new(introspect.Column),
	FuzzTableColumn31Column: new(introspect.Column),
	FuzzTableColumn32Column: new(introspect.Column),
	FuzzTableColumn33Column: new(introspect.Column),
}

var (
	FuzzTablePrimaryKeyColumn = FuzzTableIDColumn
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

func (m *Fuzz) GetPrimaryKeyColumn() string {
	return FuzzTablePrimaryKeyColumn
}

func (m *Fuzz) GetPrimaryKeyValue() any {
	return m.ID
}

func (m *Fuzz) FromItem(item map[string]any) error {
	if item == nil {
		return fmt.Errorf(
			"item unexpectedly nil during FuzzFromItem",
		)
	}

	if len(item) == 0 {
		return fmt.Errorf(
			"item unexpectedly empty during FuzzFromItem",
		)
	}

	wrapError := func(k string, v any, err error) error {
		return fmt.Errorf("%v: %#+v; error: %v", k, v, err)
	}

	for k, v := range item {
		_, ok := FuzzTableColumnLookup[k]
		if !ok {
			return fmt.Errorf(
				"item contained unexpected key %#+v during FuzzFromItem; item: %#+v",
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

		case "column1":
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

			m.Column1 = &temp2

		case "column2":
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

			m.Column2 = &temp2

		case "column3":
			if v == nil {
				continue
			}

			temp1, err := types.ParseJSON(v)
			if err != nil {
				return wrapError(k, v, err)
			}

			temp2, ok := temp1.(any)
			if !ok {
				if temp1 != nil {
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to any", temp1))
				}
			}

			m.Column3 = &temp2

		case "column4":
			if v == nil {
				continue
			}

			temp1, err := types.ParseJSON(v)
			if err != nil {
				return wrapError(k, v, err)
			}

			temp2, ok := temp1.(any)
			if !ok {
				if temp1 != nil {
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to any", temp1))
				}
			}

			m.Column4 = &temp2

		case "column5":
			if v == nil {
				continue
			}

			temp1, err := types.ParseStringArray(v)
			if err != nil {
				return wrapError(k, v, err)
			}

			temp2, ok := temp1.([]string)
			if !ok {
				if temp1 != nil {
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to []string", temp1))
				}
			}

			m.Column5 = &temp2

		case "column6":
			if v == nil {
				continue
			}

			temp1, err := types.ParseStringArray(v)
			if err != nil {
				return wrapError(k, v, err)
			}

			temp2, ok := temp1.([]string)
			if !ok {
				if temp1 != nil {
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to []string", temp1))
				}
			}

			m.Column6 = &temp2

		case "column7":
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

			m.Column7 = &temp2

		case "column8":
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

			m.Column8 = &temp2

		case "column9":
			if v == nil {
				continue
			}

			temp1, err := types.ParseNotImplemented(v)
			if err != nil {
				return wrapError(k, v, err)
			}

			temp2, ok := temp1.(pq.Int64Array)
			if !ok {
				if temp1 != nil {
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to pq.Int64Array", temp1))
				}
			}

			m.Column9 = &temp2

		case "column10":
			if v == nil {
				continue
			}

			temp1, err := types.ParseNotImplemented(v)
			if err != nil {
				return wrapError(k, v, err)
			}

			temp2, ok := temp1.(pq.Int64Array)
			if !ok {
				if temp1 != nil {
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to pq.Int64Array", temp1))
				}
			}

			m.Column10 = &temp2

		case "column11":
			if v == nil {
				continue
			}

			temp1, err := types.ParseNotImplemented(v)
			if err != nil {
				return wrapError(k, v, err)
			}

			temp2, ok := temp1.(pq.Int64Array)
			if !ok {
				if temp1 != nil {
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to pq.Int64Array", temp1))
				}
			}

			m.Column11 = &temp2

		case "column12":
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
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to int64", temp1))
				}
			}

			m.Column12 = &temp2

		case "column13":
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
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to int64", temp1))
				}
			}

			m.Column13 = &temp2

		case "column14":
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
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to int64", temp1))
				}
			}

			m.Column14 = &temp2

		case "column15":
			if v == nil {
				continue
			}

			temp1, err := types.ParseNotImplemented(v)
			if err != nil {
				return wrapError(k, v, err)
			}

			temp2, ok := temp1.(pq.Float64Array)
			if !ok {
				if temp1 != nil {
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to pq.Float64Array", temp1))
				}
			}

			m.Column15 = &temp2

		case "column16":
			if v == nil {
				continue
			}

			temp1, err := types.ParseNotImplemented(v)
			if err != nil {
				return wrapError(k, v, err)
			}

			temp2, ok := temp1.(pq.Float64Array)
			if !ok {
				if temp1 != nil {
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to pq.Float64Array", temp1))
				}
			}

			m.Column16 = &temp2

		case "column17":
			if v == nil {
				continue
			}

			temp1, err := types.ParseNotImplemented(v)
			if err != nil {
				return wrapError(k, v, err)
			}

			temp2, ok := temp1.(pq.Float64Array)
			if !ok {
				if temp1 != nil {
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to pq.Float64Array", temp1))
				}
			}

			m.Column17 = &temp2

		case "column18":
			if v == nil {
				continue
			}

			temp1, err := types.ParseNotImplemented(v)
			if err != nil {
				return wrapError(k, v, err)
			}

			temp2, ok := temp1.(pq.Float64Array)
			if !ok {
				if temp1 != nil {
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to pq.Float64Array", temp1))
				}
			}

			m.Column18 = &temp2

		case "column19":
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

			m.Column19 = &temp2

		case "column20":
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

			m.Column20 = &temp2

		case "column21":
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

			m.Column21 = &temp2

		case "column22":
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

			m.Column22 = &temp2

		case "column23":
			if v == nil {
				continue
			}

			temp1, err := types.ParseNotImplemented(v)
			if err != nil {
				return wrapError(k, v, err)
			}

			temp2, ok := temp1.(pq.BoolArray)
			if !ok {
				if temp1 != nil {
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to pq.BoolArray", temp1))
				}
			}

			m.Column23 = &temp2

		case "column24":
			if v == nil {
				continue
			}

			temp1, err := types.ParseBool(v)
			if err != nil {
				return wrapError(k, v, err)
			}

			temp2, ok := temp1.(bool)
			if !ok {
				if temp1 != nil {
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to bool", temp1))
				}
			}

			m.Column24 = &temp2

		case "column25":
			if v == nil {
				continue
			}

			temp1, err := types.ParseTSVector(v)
			if err != nil {
				return wrapError(k, v, err)
			}

			temp2, ok := temp1.(map[string][]int)
			if !ok {
				if temp1 != nil {
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to map[string][]int", temp1))
				}
			}

			m.Column25 = &temp2

		case "column26":
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

			m.Column26 = &temp2

		case "column27":
			if v == nil {
				continue
			}

			temp1, err := types.ParseHstore(v)
			if err != nil {
				return wrapError(k, v, err)
			}

			temp2, ok := temp1.(map[string]*string)
			if !ok {
				if temp1 != nil {
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to map[string]*string", temp1))
				}
			}

			m.Column27 = &temp2

		case "column28":
			if v == nil {
				continue
			}

			temp1, err := types.ParsePoint(v)
			if err != nil {
				return wrapError(k, v, err)
			}

			temp2, ok := temp1.(pgtype.Point)
			if !ok {
				if temp1 != nil {
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to pgtype.Point", temp1))
				}
			}

			m.Column28 = &temp2

		case "column29":
			if v == nil {
				continue
			}

			temp1, err := types.ParsePolygon(v)
			if err != nil {
				return wrapError(k, v, err)
			}

			temp2, ok := temp1.(pgtype.Polygon)
			if !ok {
				if temp1 != nil {
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to pgtype.Polygon", temp1))
				}
			}

			m.Column29 = &temp2

		case "column30":
			if v == nil {
				continue
			}

			temp1, err := types.ParseGeometry(v)
			if err != nil {
				return wrapError(k, v, err)
			}

			temp2, ok := temp1.(postgis.PointZ)
			if !ok {
				if temp1 != nil {
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to postgis.PointZ", temp1))
				}
			}

			m.Column30 = &temp2

		case "column31":
			if v == nil {
				continue
			}

			temp1, err := types.ParseGeometry(v)
			if err != nil {
				return wrapError(k, v, err)
			}

			temp2, ok := temp1.(postgis.PointZ)
			if !ok {
				if temp1 != nil {
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to postgis.PointZ", temp1))
				}
			}

			m.Column31 = &temp2

		case "column32":
			if v == nil {
				continue
			}

			temp1, err := types.ParseInet(v)
			if err != nil {
				return wrapError(k, v, err)
			}

			temp2, ok := temp1.(netip.Prefix)
			if !ok {
				if temp1 != nil {
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to netip.Prefix", temp1))
				}
			}

			m.Column32 = &temp2

		case "column33":
			if v == nil {
				continue
			}

			temp1, err := types.ParseBytes(v)
			if err != nil {
				return wrapError(k, v, err)
			}

			temp2, ok := temp1.([]byte)
			if !ok {
				if temp1 != nil {
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to []byte", temp1))
				}
			}

			m.Column33 = &temp2

		}
	}

	return nil
}

func (m *Fuzz) Reload(
	ctx context.Context,
	tx *sqlx.Tx,
	includeDeleteds ...bool,
) error {
	extraWhere := ""
	if len(includeDeleteds) > 0 && includeDeleteds[0] {
		if slices.Contains(FuzzTableColumns, "deleted_at") {
			extraWhere = "\n    AND (deleted_at IS null OR deleted_at IS NOT null)"
		}
	}

	t, err := SelectFuzz(
		ctx,
		tx,
		fmt.Sprintf("%v = $1%v", m.GetPrimaryKeyColumn(), extraWhere),
		m.GetPrimaryKeyValue(),
	)
	if err != nil {
		return err
	}

	m.ID = t.ID
	m.Column1 = t.Column1
	m.Column2 = t.Column2
	m.Column3 = t.Column3
	m.Column4 = t.Column4
	m.Column5 = t.Column5
	m.Column6 = t.Column6
	m.Column7 = t.Column7
	m.Column8 = t.Column8
	m.Column9 = t.Column9
	m.Column10 = t.Column10
	m.Column11 = t.Column11
	m.Column12 = t.Column12
	m.Column13 = t.Column13
	m.Column14 = t.Column14
	m.Column15 = t.Column15
	m.Column16 = t.Column16
	m.Column17 = t.Column17
	m.Column18 = t.Column18
	m.Column19 = t.Column19
	m.Column20 = t.Column20
	m.Column21 = t.Column21
	m.Column22 = t.Column22
	m.Column23 = t.Column23
	m.Column24 = t.Column24
	m.Column25 = t.Column25
	m.Column26 = t.Column26
	m.Column27 = t.Column27
	m.Column28 = t.Column28
	m.Column29 = t.Column29
	m.Column30 = t.Column30
	m.Column31 = t.Column31
	m.Column32 = t.Column32
	m.Column33 = t.Column33

	return nil
}

func (m *Fuzz) Insert(
	ctx context.Context,
	tx *sqlx.Tx,
	setPrimaryKey bool,
	setZeroValues bool,
) error {
	columns := make([]string, 0)
	values := make([]any, 0)

	if setPrimaryKey && (setZeroValues || !types.IsZeroUUID(m.ID)) {
		columns = append(columns, FuzzTableIDColumn)

		v, err := types.FormatUUID(m.ID)
		if err != nil {
			return fmt.Errorf("failed to handle m.ID: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.Column1) {
		columns = append(columns, FuzzTableColumn1Column)

		v, err := types.FormatTime(m.Column1)
		if err != nil {
			return fmt.Errorf("failed to handle m.Column1: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.Column2) {
		columns = append(columns, FuzzTableColumn2Column)

		v, err := types.FormatTime(m.Column2)
		if err != nil {
			return fmt.Errorf("failed to handle m.Column2: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroJSON(m.Column3) {
		columns = append(columns, FuzzTableColumn3Column)

		v, err := types.FormatJSON(m.Column3)
		if err != nil {
			return fmt.Errorf("failed to handle m.Column3: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroJSON(m.Column4) {
		columns = append(columns, FuzzTableColumn4Column)

		v, err := types.FormatJSON(m.Column4)
		if err != nil {
			return fmt.Errorf("failed to handle m.Column4: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroStringArray(m.Column5) {
		columns = append(columns, FuzzTableColumn5Column)

		v, err := types.FormatStringArray(m.Column5)
		if err != nil {
			return fmt.Errorf("failed to handle m.Column5: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroStringArray(m.Column6) {
		columns = append(columns, FuzzTableColumn6Column)

		v, err := types.FormatStringArray(m.Column6)
		if err != nil {
			return fmt.Errorf("failed to handle m.Column6: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroString(m.Column7) {
		columns = append(columns, FuzzTableColumn7Column)

		v, err := types.FormatString(m.Column7)
		if err != nil {
			return fmt.Errorf("failed to handle m.Column7: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroString(m.Column8) {
		columns = append(columns, FuzzTableColumn8Column)

		v, err := types.FormatString(m.Column8)
		if err != nil {
			return fmt.Errorf("failed to handle m.Column8: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroNotImplemented(m.Column9) {
		columns = append(columns, FuzzTableColumn9Column)

		v, err := types.FormatNotImplemented(m.Column9)
		if err != nil {
			return fmt.Errorf("failed to handle m.Column9: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroNotImplemented(m.Column10) {
		columns = append(columns, FuzzTableColumn10Column)

		v, err := types.FormatNotImplemented(m.Column10)
		if err != nil {
			return fmt.Errorf("failed to handle m.Column10: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroNotImplemented(m.Column11) {
		columns = append(columns, FuzzTableColumn11Column)

		v, err := types.FormatNotImplemented(m.Column11)
		if err != nil {
			return fmt.Errorf("failed to handle m.Column11: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroInt(m.Column12) {
		columns = append(columns, FuzzTableColumn12Column)

		v, err := types.FormatInt(m.Column12)
		if err != nil {
			return fmt.Errorf("failed to handle m.Column12: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroInt(m.Column13) {
		columns = append(columns, FuzzTableColumn13Column)

		v, err := types.FormatInt(m.Column13)
		if err != nil {
			return fmt.Errorf("failed to handle m.Column13: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroInt(m.Column14) {
		columns = append(columns, FuzzTableColumn14Column)

		v, err := types.FormatInt(m.Column14)
		if err != nil {
			return fmt.Errorf("failed to handle m.Column14: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroNotImplemented(m.Column15) {
		columns = append(columns, FuzzTableColumn15Column)

		v, err := types.FormatNotImplemented(m.Column15)
		if err != nil {
			return fmt.Errorf("failed to handle m.Column15: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroNotImplemented(m.Column16) {
		columns = append(columns, FuzzTableColumn16Column)

		v, err := types.FormatNotImplemented(m.Column16)
		if err != nil {
			return fmt.Errorf("failed to handle m.Column16: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroNotImplemented(m.Column17) {
		columns = append(columns, FuzzTableColumn17Column)

		v, err := types.FormatNotImplemented(m.Column17)
		if err != nil {
			return fmt.Errorf("failed to handle m.Column17: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroNotImplemented(m.Column18) {
		columns = append(columns, FuzzTableColumn18Column)

		v, err := types.FormatNotImplemented(m.Column18)
		if err != nil {
			return fmt.Errorf("failed to handle m.Column18: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroFloat(m.Column19) {
		columns = append(columns, FuzzTableColumn19Column)

		v, err := types.FormatFloat(m.Column19)
		if err != nil {
			return fmt.Errorf("failed to handle m.Column19: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroFloat(m.Column20) {
		columns = append(columns, FuzzTableColumn20Column)

		v, err := types.FormatFloat(m.Column20)
		if err != nil {
			return fmt.Errorf("failed to handle m.Column20: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroFloat(m.Column21) {
		columns = append(columns, FuzzTableColumn21Column)

		v, err := types.FormatFloat(m.Column21)
		if err != nil {
			return fmt.Errorf("failed to handle m.Column21: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroFloat(m.Column22) {
		columns = append(columns, FuzzTableColumn22Column)

		v, err := types.FormatFloat(m.Column22)
		if err != nil {
			return fmt.Errorf("failed to handle m.Column22: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroNotImplemented(m.Column23) {
		columns = append(columns, FuzzTableColumn23Column)

		v, err := types.FormatNotImplemented(m.Column23)
		if err != nil {
			return fmt.Errorf("failed to handle m.Column23: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroBool(m.Column24) {
		columns = append(columns, FuzzTableColumn24Column)

		v, err := types.FormatBool(m.Column24)
		if err != nil {
			return fmt.Errorf("failed to handle m.Column24: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTSVector(m.Column25) {
		columns = append(columns, FuzzTableColumn25Column)

		v, err := types.FormatTSVector(m.Column25)
		if err != nil {
			return fmt.Errorf("failed to handle m.Column25: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroUUID(m.Column26) {
		columns = append(columns, FuzzTableColumn26Column)

		v, err := types.FormatUUID(m.Column26)
		if err != nil {
			return fmt.Errorf("failed to handle m.Column26: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroHstore(m.Column27) {
		columns = append(columns, FuzzTableColumn27Column)

		v, err := types.FormatHstore(m.Column27)
		if err != nil {
			return fmt.Errorf("failed to handle m.Column27: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroPoint(m.Column28) {
		columns = append(columns, FuzzTableColumn28Column)

		v, err := types.FormatPoint(m.Column28)
		if err != nil {
			return fmt.Errorf("failed to handle m.Column28: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroPolygon(m.Column29) {
		columns = append(columns, FuzzTableColumn29Column)

		v, err := types.FormatPolygon(m.Column29)
		if err != nil {
			return fmt.Errorf("failed to handle m.Column29: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroGeometry(m.Column30) {
		columns = append(columns, FuzzTableColumn30Column)

		v, err := types.FormatGeometry(m.Column30)
		if err != nil {
			return fmt.Errorf("failed to handle m.Column30: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroGeometry(m.Column31) {
		columns = append(columns, FuzzTableColumn31Column)

		v, err := types.FormatGeometry(m.Column31)
		if err != nil {
			return fmt.Errorf("failed to handle m.Column31: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroInet(m.Column32) {
		columns = append(columns, FuzzTableColumn32Column)

		v, err := types.FormatInet(m.Column32)
		if err != nil {
			return fmt.Errorf("failed to handle m.Column32: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroBytes(m.Column33) {
		columns = append(columns, FuzzTableColumn33Column)

		v, err := types.FormatBytes(m.Column33)
		if err != nil {
			return fmt.Errorf("failed to handle m.Column33: %v", err)
		}

		values = append(values, v)
	}

	item, err := query.Insert(
		ctx,
		tx,
		FuzzTable,
		columns,
		nil,
		false,
		false,
		FuzzTableColumns,
		values...,
	)
	if err != nil {
		return fmt.Errorf("failed to insert %#+v: %v", m, err)
	}
	v := item[FuzzTableIDColumn]

	if v == nil {
		return fmt.Errorf("failed to find %v in %#+v", FuzzTableIDColumn, item)
	}

	wrapError := func(err error) error {
		return fmt.Errorf(
			"failed to treat %v: %#+v as uuid.UUID: %v",
			FuzzTableIDColumn,
			item[FuzzTableIDColumn],
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

func (m *Fuzz) Update(
	ctx context.Context,
	tx *sqlx.Tx,
	setZeroValues bool,
	forceSetValuesForFields ...string,
) error {
	columns := make([]string, 0)
	values := make([]any, 0)

	if setZeroValues || !types.IsZeroTime(m.Column1) || slices.Contains(forceSetValuesForFields, FuzzTableColumn1Column) {
		columns = append(columns, FuzzTableColumn1Column)

		v, err := types.FormatTime(m.Column1)
		if err != nil {
			return fmt.Errorf("failed to handle m.Column1: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.Column2) || slices.Contains(forceSetValuesForFields, FuzzTableColumn2Column) {
		columns = append(columns, FuzzTableColumn2Column)

		v, err := types.FormatTime(m.Column2)
		if err != nil {
			return fmt.Errorf("failed to handle m.Column2: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroJSON(m.Column3) || slices.Contains(forceSetValuesForFields, FuzzTableColumn3Column) {
		columns = append(columns, FuzzTableColumn3Column)

		v, err := types.FormatJSON(m.Column3)
		if err != nil {
			return fmt.Errorf("failed to handle m.Column3: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroJSON(m.Column4) || slices.Contains(forceSetValuesForFields, FuzzTableColumn4Column) {
		columns = append(columns, FuzzTableColumn4Column)

		v, err := types.FormatJSON(m.Column4)
		if err != nil {
			return fmt.Errorf("failed to handle m.Column4: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroStringArray(m.Column5) || slices.Contains(forceSetValuesForFields, FuzzTableColumn5Column) {
		columns = append(columns, FuzzTableColumn5Column)

		v, err := types.FormatStringArray(m.Column5)
		if err != nil {
			return fmt.Errorf("failed to handle m.Column5: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroStringArray(m.Column6) || slices.Contains(forceSetValuesForFields, FuzzTableColumn6Column) {
		columns = append(columns, FuzzTableColumn6Column)

		v, err := types.FormatStringArray(m.Column6)
		if err != nil {
			return fmt.Errorf("failed to handle m.Column6: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroString(m.Column7) || slices.Contains(forceSetValuesForFields, FuzzTableColumn7Column) {
		columns = append(columns, FuzzTableColumn7Column)

		v, err := types.FormatString(m.Column7)
		if err != nil {
			return fmt.Errorf("failed to handle m.Column7: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroString(m.Column8) || slices.Contains(forceSetValuesForFields, FuzzTableColumn8Column) {
		columns = append(columns, FuzzTableColumn8Column)

		v, err := types.FormatString(m.Column8)
		if err != nil {
			return fmt.Errorf("failed to handle m.Column8: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroNotImplemented(m.Column9) || slices.Contains(forceSetValuesForFields, FuzzTableColumn9Column) {
		columns = append(columns, FuzzTableColumn9Column)

		v, err := types.FormatNotImplemented(m.Column9)
		if err != nil {
			return fmt.Errorf("failed to handle m.Column9: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroNotImplemented(m.Column10) || slices.Contains(forceSetValuesForFields, FuzzTableColumn10Column) {
		columns = append(columns, FuzzTableColumn10Column)

		v, err := types.FormatNotImplemented(m.Column10)
		if err != nil {
			return fmt.Errorf("failed to handle m.Column10: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroNotImplemented(m.Column11) || slices.Contains(forceSetValuesForFields, FuzzTableColumn11Column) {
		columns = append(columns, FuzzTableColumn11Column)

		v, err := types.FormatNotImplemented(m.Column11)
		if err != nil {
			return fmt.Errorf("failed to handle m.Column11: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroInt(m.Column12) || slices.Contains(forceSetValuesForFields, FuzzTableColumn12Column) {
		columns = append(columns, FuzzTableColumn12Column)

		v, err := types.FormatInt(m.Column12)
		if err != nil {
			return fmt.Errorf("failed to handle m.Column12: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroInt(m.Column13) || slices.Contains(forceSetValuesForFields, FuzzTableColumn13Column) {
		columns = append(columns, FuzzTableColumn13Column)

		v, err := types.FormatInt(m.Column13)
		if err != nil {
			return fmt.Errorf("failed to handle m.Column13: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroInt(m.Column14) || slices.Contains(forceSetValuesForFields, FuzzTableColumn14Column) {
		columns = append(columns, FuzzTableColumn14Column)

		v, err := types.FormatInt(m.Column14)
		if err != nil {
			return fmt.Errorf("failed to handle m.Column14: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroNotImplemented(m.Column15) || slices.Contains(forceSetValuesForFields, FuzzTableColumn15Column) {
		columns = append(columns, FuzzTableColumn15Column)

		v, err := types.FormatNotImplemented(m.Column15)
		if err != nil {
			return fmt.Errorf("failed to handle m.Column15: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroNotImplemented(m.Column16) || slices.Contains(forceSetValuesForFields, FuzzTableColumn16Column) {
		columns = append(columns, FuzzTableColumn16Column)

		v, err := types.FormatNotImplemented(m.Column16)
		if err != nil {
			return fmt.Errorf("failed to handle m.Column16: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroNotImplemented(m.Column17) || slices.Contains(forceSetValuesForFields, FuzzTableColumn17Column) {
		columns = append(columns, FuzzTableColumn17Column)

		v, err := types.FormatNotImplemented(m.Column17)
		if err != nil {
			return fmt.Errorf("failed to handle m.Column17: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroNotImplemented(m.Column18) || slices.Contains(forceSetValuesForFields, FuzzTableColumn18Column) {
		columns = append(columns, FuzzTableColumn18Column)

		v, err := types.FormatNotImplemented(m.Column18)
		if err != nil {
			return fmt.Errorf("failed to handle m.Column18: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroFloat(m.Column19) || slices.Contains(forceSetValuesForFields, FuzzTableColumn19Column) {
		columns = append(columns, FuzzTableColumn19Column)

		v, err := types.FormatFloat(m.Column19)
		if err != nil {
			return fmt.Errorf("failed to handle m.Column19: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroFloat(m.Column20) || slices.Contains(forceSetValuesForFields, FuzzTableColumn20Column) {
		columns = append(columns, FuzzTableColumn20Column)

		v, err := types.FormatFloat(m.Column20)
		if err != nil {
			return fmt.Errorf("failed to handle m.Column20: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroFloat(m.Column21) || slices.Contains(forceSetValuesForFields, FuzzTableColumn21Column) {
		columns = append(columns, FuzzTableColumn21Column)

		v, err := types.FormatFloat(m.Column21)
		if err != nil {
			return fmt.Errorf("failed to handle m.Column21: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroFloat(m.Column22) || slices.Contains(forceSetValuesForFields, FuzzTableColumn22Column) {
		columns = append(columns, FuzzTableColumn22Column)

		v, err := types.FormatFloat(m.Column22)
		if err != nil {
			return fmt.Errorf("failed to handle m.Column22: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroNotImplemented(m.Column23) || slices.Contains(forceSetValuesForFields, FuzzTableColumn23Column) {
		columns = append(columns, FuzzTableColumn23Column)

		v, err := types.FormatNotImplemented(m.Column23)
		if err != nil {
			return fmt.Errorf("failed to handle m.Column23: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroBool(m.Column24) || slices.Contains(forceSetValuesForFields, FuzzTableColumn24Column) {
		columns = append(columns, FuzzTableColumn24Column)

		v, err := types.FormatBool(m.Column24)
		if err != nil {
			return fmt.Errorf("failed to handle m.Column24: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTSVector(m.Column25) || slices.Contains(forceSetValuesForFields, FuzzTableColumn25Column) {
		columns = append(columns, FuzzTableColumn25Column)

		v, err := types.FormatTSVector(m.Column25)
		if err != nil {
			return fmt.Errorf("failed to handle m.Column25: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroUUID(m.Column26) || slices.Contains(forceSetValuesForFields, FuzzTableColumn26Column) {
		columns = append(columns, FuzzTableColumn26Column)

		v, err := types.FormatUUID(m.Column26)
		if err != nil {
			return fmt.Errorf("failed to handle m.Column26: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroHstore(m.Column27) || slices.Contains(forceSetValuesForFields, FuzzTableColumn27Column) {
		columns = append(columns, FuzzTableColumn27Column)

		v, err := types.FormatHstore(m.Column27)
		if err != nil {
			return fmt.Errorf("failed to handle m.Column27: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroPoint(m.Column28) || slices.Contains(forceSetValuesForFields, FuzzTableColumn28Column) {
		columns = append(columns, FuzzTableColumn28Column)

		v, err := types.FormatPoint(m.Column28)
		if err != nil {
			return fmt.Errorf("failed to handle m.Column28: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroPolygon(m.Column29) || slices.Contains(forceSetValuesForFields, FuzzTableColumn29Column) {
		columns = append(columns, FuzzTableColumn29Column)

		v, err := types.FormatPolygon(m.Column29)
		if err != nil {
			return fmt.Errorf("failed to handle m.Column29: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroGeometry(m.Column30) || slices.Contains(forceSetValuesForFields, FuzzTableColumn30Column) {
		columns = append(columns, FuzzTableColumn30Column)

		v, err := types.FormatGeometry(m.Column30)
		if err != nil {
			return fmt.Errorf("failed to handle m.Column30: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroGeometry(m.Column31) || slices.Contains(forceSetValuesForFields, FuzzTableColumn31Column) {
		columns = append(columns, FuzzTableColumn31Column)

		v, err := types.FormatGeometry(m.Column31)
		if err != nil {
			return fmt.Errorf("failed to handle m.Column31: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroInet(m.Column32) || slices.Contains(forceSetValuesForFields, FuzzTableColumn32Column) {
		columns = append(columns, FuzzTableColumn32Column)

		v, err := types.FormatInet(m.Column32)
		if err != nil {
			return fmt.Errorf("failed to handle m.Column32: %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroBytes(m.Column33) || slices.Contains(forceSetValuesForFields, FuzzTableColumn33Column) {
		columns = append(columns, FuzzTableColumn33Column)

		v, err := types.FormatBytes(m.Column33)
		if err != nil {
			return fmt.Errorf("failed to handle m.Column33: %v", err)
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
		FuzzTable,
		columns,
		fmt.Sprintf("%v = $$??", FuzzTableIDColumn),
		FuzzTableColumns,
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

func (m *Fuzz) Delete(
	ctx context.Context,
	tx *sqlx.Tx,
	hardDeletes ...bool,
) error {
	/* soft-delete not applicable */

	values := make([]any, 0)
	v, err := types.FormatUUID(m.ID)
	if err != nil {
		return fmt.Errorf("failed to handle m.ID: %v", err)
	}

	values = append(values, v)

	err = query.Delete(
		ctx,
		tx,
		FuzzTable,
		fmt.Sprintf("%v = $$??", FuzzTableIDColumn),
		values...,
	)
	if err != nil {
		return fmt.Errorf("failed to delete %#+v: %v", m, err)
	}

	_ = m.Reload(ctx, tx, true)

	return nil
}

func SelectFuzzes(
	ctx context.Context,
	tx *sqlx.Tx,
	where string,
	orderBy *string,
	limit *int,
	offset *int,
	values ...any,
) ([]*Fuzz, error) {
	if slices.Contains(FuzzTableColumns, "deleted_at") {
		if !strings.Contains(where, "deleted_at") {
			if where != "" {
				where += "\n    AND "
			}

			where += "deleted_at IS null"
		}
	}

	items, err := query.Select(
		ctx,
		tx,
		FuzzTableColumnsWithTypeCasts,
		FuzzTable,
		where,
		orderBy,
		limit,
		offset,
		values...,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to call SelectFuzzs; err: %v", err)
	}

	objects := make([]*Fuzz, 0)

	for _, item := range items {
		object := &Fuzz{}

		err = object.FromItem(item)
		if err != nil {
			return nil, err
		}

		objects = append(objects, object)
	}

	return objects, nil
}

func SelectFuzz(
	ctx context.Context,
	tx *sqlx.Tx,
	where string,
	values ...any,
) (*Fuzz, error) {
	objects, err := SelectFuzzes(
		ctx,
		tx,
		where,
		nil,
		helpers.Ptr(2),
		helpers.Ptr(0),
		values...,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to call SelectFuzz; err: %v", err)
	}

	if len(objects) > 1 {
		return nil, fmt.Errorf("attempt to call SelectFuzz returned more than 1 row")
	}

	if len(objects) < 1 {
		return nil, sql.ErrNoRows
	}

	object := objects[0]

	return object, nil
}

func handleGetFuzzes(w http.ResponseWriter, r *http.Request, db *sqlx.DB, redisPool *redis.Pool, objectMiddlewares []server.ObjectMiddleware) {
	ctx := r.Context()

	insaneOrderParams := make([]string, 0)
	hadInsaneOrderParams := false

	unrecognizedParams := make([]string, 0)
	hadUnrecognizedParams := false

	unparseableParams := make([]string, 0)
	hadUnparseableParams := false

	var orderByDirection *string
	orderBys := make([]string, 0)

	includes := make([]string, 0)

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
			column := FuzzTableColumnLookup[parts[0]]
			if column == nil {
				if parts[0] != "load" {
					isUnrecognized = true
				}
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
				case "load":
					includes = append(includes, parts[0])
					_ = includes

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

	requestHash, err := helpers.GetRequestHash(FuzzTable, wheres, hashableOrderBy, limit, offset, values, nil)
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

	tx, err := db.BeginTxx(r.Context(), nil)
	if err != nil {
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	defer func() {
		_ = tx.Rollback()
	}()

	where := strings.Join(wheres, "\n    AND ")

	objects, err := SelectFuzzes(ctx, tx, where, orderBy, &limit, &offset, values...)
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

func handleGetFuzz(w http.ResponseWriter, r *http.Request, db *sqlx.DB, redisPool *redis.Pool, objectMiddlewares []server.ObjectMiddleware, primaryKey string) {
	ctx := r.Context()

	wheres := []string{fmt.Sprintf("%s = $$??", FuzzTablePrimaryKeyColumn)}
	values := []any{primaryKey}

	requestHash, err := helpers.GetRequestHash(FuzzTable, wheres, "", 2, 0, values, primaryKey)
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

	tx, err := db.BeginTxx(r.Context(), nil)
	if err != nil {
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	defer func() {
		_ = tx.Rollback()
	}()

	where := strings.Join(wheres, "\n    AND ")

	object, err := SelectFuzz(ctx, tx, where, values...)
	if err != nil {
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	err = tx.Commit()
	if err != nil {
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	returnedObjectsAsJSON := helpers.HandleObjectsResponse(w, http.StatusOK, []*Fuzz{object})

	err = helpers.StoreCachedResponse(requestHash, redisConn, string(returnedObjectsAsJSON))
	if err != nil {
		log.Printf("warning: %v", err)
	}
}

func handlePostFuzzs(w http.ResponseWriter, r *http.Request, db *sqlx.DB, redisPool *redis.Pool, objectMiddlewares []server.ObjectMiddleware, waitForChange server.WaitForChange) {
	_ = redisPool

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

	objects := make([]*Fuzz, 0)
	for _, item := range allItems {
		object := &Fuzz{}
		err = object.FromItem(item)
		if err != nil {
			err = fmt.Errorf("failed to interpret %#+v as Fuzz in item form: %v", item, err)
			helpers.HandleErrorResponse(w, http.StatusBadRequest, err)
			return
		}

		objects = append(objects, object)
	}

	tx, err := db.BeginTxx(r.Context(), nil)
	if err != nil {
		err = fmt.Errorf("failed to begin DB transaction: %v", err)
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	defer func() {
		_ = tx.Rollback()
	}()

	xid, err := query.GetXid(r.Context(), tx)
	if err != nil {
		err = fmt.Errorf("failed to get xid: %v", err)
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}
	_ = xid

	for i, object := range objects {
		err = object.Insert(r.Context(), tx, false, false)
		if err != nil {
			err = fmt.Errorf("failed to insert %#+v: %v", object, err)
			helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
			return
		}

		objects[i] = object
	}

	errs := make(chan error, 1)
	go func() {
		_, err = waitForChange(r.Context(), []stream.Action{stream.INSERT}, FuzzTable, xid)
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

func handlePutFuzz(w http.ResponseWriter, r *http.Request, db *sqlx.DB, redisPool *redis.Pool, objectMiddlewares []server.ObjectMiddleware, waitForChange server.WaitForChange, primaryKey string) {
	_ = redisPool

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

	item[FuzzTablePrimaryKeyColumn] = primaryKey

	object := &Fuzz{}
	err = object.FromItem(item)
	if err != nil {
		err = fmt.Errorf("failed to interpret %#+v as Fuzz in item form: %v", item, err)
		helpers.HandleErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	tx, err := db.BeginTxx(r.Context(), nil)
	if err != nil {
		err = fmt.Errorf("failed to begin DB transaction: %v", err)
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	defer func() {
		_ = tx.Rollback()
	}()

	xid, err := query.GetXid(r.Context(), tx)
	if err != nil {
		err = fmt.Errorf("failed to get xid: %v", err)
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}
	_ = xid

	err = object.Update(r.Context(), tx, true)
	if err != nil {
		err = fmt.Errorf("failed to update %#+v: %v", object, err)
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	errs := make(chan error, 1)
	go func() {
		_, err = waitForChange(r.Context(), []stream.Action{stream.UPDATE, stream.SOFT_DELETE, stream.SOFT_RESTORE, stream.SOFT_UPDATE}, FuzzTable, xid)
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

	helpers.HandleObjectsResponse(w, http.StatusOK, []*Fuzz{object})
}

func handlePatchFuzz(w http.ResponseWriter, r *http.Request, db *sqlx.DB, redisPool *redis.Pool, objectMiddlewares []server.ObjectMiddleware, waitForChange server.WaitForChange, primaryKey string) {
	_ = redisPool

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
		if !slices.Contains(FuzzTableColumns, possibleField) {
			continue
		}

		forceSetValuesForFields = append(forceSetValuesForFields, possibleField)
	}

	item[FuzzTablePrimaryKeyColumn] = primaryKey

	object := &Fuzz{}
	err = object.FromItem(item)
	if err != nil {
		err = fmt.Errorf("failed to interpret %#+v as Fuzz in item form: %v", item, err)
		helpers.HandleErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	tx, err := db.BeginTxx(r.Context(), nil)
	if err != nil {
		err = fmt.Errorf("failed to begin DB transaction: %v", err)
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	defer func() {
		_ = tx.Rollback()
	}()

	xid, err := query.GetXid(r.Context(), tx)
	if err != nil {
		err = fmt.Errorf("failed to get xid: %v", err)
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}
	_ = xid

	err = object.Update(r.Context(), tx, false, forceSetValuesForFields...)
	if err != nil {
		err = fmt.Errorf("failed to update %#+v: %v", object, err)
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	errs := make(chan error, 1)
	go func() {
		_, err = waitForChange(r.Context(), []stream.Action{stream.UPDATE, stream.SOFT_DELETE, stream.SOFT_RESTORE, stream.SOFT_UPDATE}, FuzzTable, xid)
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

	helpers.HandleObjectsResponse(w, http.StatusOK, []*Fuzz{object})
}

func handleDeleteFuzz(w http.ResponseWriter, r *http.Request, db *sqlx.DB, redisPool *redis.Pool, objectMiddlewares []server.ObjectMiddleware, waitForChange server.WaitForChange, primaryKey string) {
	_ = redisPool

	var item = make(map[string]any)

	item[FuzzTablePrimaryKeyColumn] = primaryKey

	object := &Fuzz{}
	err := object.FromItem(item)
	if err != nil {
		err = fmt.Errorf("failed to interpret %#+v as Fuzz in item form: %v", item, err)
		helpers.HandleErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	tx, err := db.BeginTxx(r.Context(), nil)
	if err != nil {
		err = fmt.Errorf("failed to begin DB transaction: %v", err)
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	defer func() {
		_ = tx.Rollback()
	}()

	xid, err := query.GetXid(r.Context(), tx)
	if err != nil {
		err = fmt.Errorf("failed to get xid: %v", err)
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}
	_ = xid

	err = object.Delete(r.Context(), tx)
	if err != nil {
		err = fmt.Errorf("failed to delete %#+v: %v", object, err)
		helpers.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	errs := make(chan error, 1)
	go func() {
		_, err = waitForChange(r.Context(), []stream.Action{stream.DELETE, stream.SOFT_DELETE}, FuzzTable, xid)
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

func GetFuzzRouter(db *sqlx.DB, redisPool *redis.Pool, httpMiddlewares []server.HTTPMiddleware, objectMiddlewares []server.ObjectMiddleware, waitForChange server.WaitForChange) chi.Router {
	r := chi.NewRouter()

	for _, m := range httpMiddlewares {
		r.Use(m)
	}

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("handleGetFuzzes()")
		handleGetFuzzes(w, r, db, redisPool, objectMiddlewares)
	})

	r.Get("/{primaryKey}", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("handleGetFuzz(%v)", chi.URLParam(r, "primaryKey"))
		handleGetFuzz(w, r, db, redisPool, objectMiddlewares, chi.URLParam(r, "primaryKey"))
	})

	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		handlePostFuzzs(w, r, db, redisPool, objectMiddlewares, waitForChange)
	})

	r.Put("/{primaryKey}", func(w http.ResponseWriter, r *http.Request) {
		handlePutFuzz(w, r, db, redisPool, objectMiddlewares, waitForChange, chi.URLParam(r, "primaryKey"))
	})

	r.Patch("/{primaryKey}", func(w http.ResponseWriter, r *http.Request) {
		handlePatchFuzz(w, r, db, redisPool, objectMiddlewares, waitForChange, chi.URLParam(r, "primaryKey"))
	})

	r.Delete("/{primaryKey}", func(w http.ResponseWriter, r *http.Request) {
		handleDeleteFuzz(w, r, db, redisPool, objectMiddlewares, waitForChange, chi.URLParam(r, "primaryKey"))
	})

	return r
}

func NewFuzzFromItem(item map[string]any) (any, error) {
	object := &Fuzz{}

	err := object.FromItem(item)
	if err != nil {
		return nil, err
	}

	return object, nil
}

func init() {
	register(
		FuzzTable,
		Fuzz{},
		NewFuzzFromItem,
		"/fuzzes",
		GetFuzzRouter,
	)
}
