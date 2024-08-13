package types

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/netip"
	"reflect"
	"strings"
	"time"

	"github.com/aymericbeaumet/go-tsvector"
	"github.com/cridenour/go-postgis"
	"github.com/google/uuid"
	"github.com/initialed85/djangolang/pkg/helpers"
	_pgtype "github.com/jackc/pgtype"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/lib/pq"
	"github.com/lib/pq/hstore"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/ewkb"
)

func typeOf(v any) string {
	if v == nil {
		return "nil"
	}

	return reflect.TypeOf(v).String()
}

func ParseNotImplemented(v any) (any, error) {
	return nil, fmt.Errorf("parse not implemented for %#+v", v)
}

func FormatNotImplemented(v any) (any, error) {
	return nil, fmt.Errorf("format not implemented for %#+v", v)
}

func IsZeroNotImplemented(v any) bool {
	return true
}

func GetOpenAPISchemaNotImplemented() *Schema {
	return &Schema{}
}

func GetOpenAPISchemaUUID() *Schema {
	return &Schema{
		Type:   TypeOfString,
		Format: FormatOfUUID,
	}
}

func ParseUUID(v any) (any, error) {
	switch v1 := v.(type) {

	case string:
		v2, err := uuid.Parse(v1)
		if err != nil {
			return uuid.UUID{}, fmt.Errorf("%#+v (%v) could not be parsed with uuid.Parse for ParseUUID; err: %v", v, typeOf(v), err)
		}

		return v2, nil

	case []byte:
		v2, err := uuid.Parse(string(v1))
		if err != nil {
			return uuid.UUID{}, fmt.Errorf("%#+v (%v) could not be parsed with uuid.Parse for ParseUUID; err: %v", v, typeOf(v), err)
		}

		return v2, nil

	case [16]byte:
		v2, err := uuid.FromBytes(v1[:])
		if err != nil {
			return uuid.UUID{}, fmt.Errorf("%#+v (%v) could not be parsed with uuid.ParseBytes for ParseUUID; err: %v", v, typeOf(v), err)
		}

		return v2, nil
	}

	return uuid.UUID{}, fmt.Errorf("%#+v (%v) could not be identified for ParseUUID", v, typeOf(v))
}

func FormatUUID(v any) (any, error) {
	v1, ok := v.(*uuid.UUID)
	if ok {
		if v1 == nil {
			return nil, nil
		}

		return *v1, nil
	}

	v2, ok := v.(uuid.UUID)
	if !ok {
		return nil, fmt.Errorf("%#+v (%v) could not be cast to uuid.UUID for FormatUUID", v, typeOf(v))
	}

	return v2, nil
}

func IsZeroUUID(v any) bool {
	if v == nil {
		return true
	}

	v1, ok := v.(*uuid.UUID)
	if ok {
		if v1 == nil {
			return true
		}

		v = *v1
	}

	v2, ok := v.(uuid.UUID)
	if !ok {
		return true
	}

	return v2 == uuid.Nil
}

func GetOpenAPISchemaTime() *Schema {
	return &Schema{
		Type:   TypeOfString,
		Format: FormatOfDateTime,
	}
}

func ParseTime(v any) (any, error) {
	switch v1 := v.(type) {
	case string:
		v2, err := time.Parse(time.RFC3339Nano, v1)
		if err != nil {
			v2, err = time.Parse(time.RFC3339, v1)
			if err != nil {
				v2, err = time.Parse("2006-01-02T15:04:05Z", v1)
				if err != nil {
					return uuid.UUID{}, fmt.Errorf("%#+v (%v) could not be parsed with time.Parse for ParseTime; err: %v", v, typeOf(v), err)
				}
			}
		}

		return v2, nil
	case time.Time:
		return v1, nil
	}

	return time.Time{}, fmt.Errorf("%#+v (%v) could not be identified ParseTime", v, typeOf(v))
}

func FormatTime(v any) (any, error) {
	v1, ok := v.(*time.Time)
	if ok {
		if v1 == nil {
			return nil, nil
		}

		return *v1, nil
	}

	v2, ok := v.(time.Time)
	if !ok {
		return nil, fmt.Errorf("%#+v (%v) could not be cast to time.Time for FormatTime", v, typeOf(v))
	}

	return v2, nil
}

func IsZeroTime(v any) bool {
	if v == nil {
		return true
	}

	v1, ok := v.(*time.Time)
	if ok {
		if v1 == nil {
			return true
		}

		v = *v1
	}

	v2, ok := v.(time.Time)
	if !ok {
		return true
	}

	return v2.IsZero()
}

func GetOpenAPISchemaDuration() *Schema {
	return &Schema{
		Type:   TypeOfInteger,
		Format: FormatOfInt64,
	}
}

func ParseDuration(v any) (any, error) {
	switch v1 := v.(type) {
	case string:
		v2 := pgtype.Interval{}
		err := v2.Scan(v1)
		if err != nil {
			return nil, fmt.Errorf("%#+v (%v) could not be parsed with pgtype.Interval.Scan for ParseDuration; err: %v", v, typeOf(v), err)
		}

		duration := time.Microsecond * time.Duration(v2.Microseconds)

		return time.Duration(duration), nil
	case []byte:
		v2 := pgtype.Interval{}
		err := v2.Scan(string(v1))
		if err != nil {
			return nil, fmt.Errorf("%#+v (%v) could not be parsed with pgtype.Interval.Scan for ParseDuration; err: %v", v, typeOf(v), err)
		}

		duration := time.Microsecond * time.Duration(v2.Microseconds)

		return time.Duration(duration), nil
	case time.Duration:
		return v1, nil
	case float64:
		return time.Duration(v1) * time.Nanosecond, nil
	case int64:
		return time.Duration(v1) * time.Nanosecond, nil
	case _pgtype.Interval:
		return (time.Microsecond * time.Duration(v1.Microseconds)), nil
	case pgtype.Interval:
		return (time.Microsecond * time.Duration(v1.Microseconds)), nil
	}

	return time.Duration(0), fmt.Errorf("%#+v (%v) could not be identified ParseDuration", v, typeOf(v))
}

func FormatDuration(v any) (any, error) {
	var v1 time.Duration

	v2, ok := v.(*time.Duration)
	if ok {
		if v2 == nil {
			return nil, nil
		}

		v1 = *v2
	} else {
		v1, ok = v.(time.Duration)
		if !ok {
			return nil, fmt.Errorf("%#+v (%v) could not be cast to time.Duration for FormatDuration", v, typeOf(v))
		}
	}

	months := int32(v1 / (time.Hour * 24 * 30))
	v1 -= time.Hour * 24 * 30 * time.Duration(months)

	days := int32(v1 / (time.Hour * 24))
	v1 -= time.Hour * 24 * time.Duration(months)

	microseconds := v1.Microseconds()
	v1 -= (time.Microsecond * time.Duration(microseconds)) // not needed, kept to balance the math

	v3 := pgtype.Interval{
		Microseconds: microseconds,
		Days:         days,
		Months:       months,
		Valid:        true,
	}

	return v3, nil
}

func IsZeroDuration(v any) bool {
	if v == nil {
		return true
	}

	v1, ok := v.(*time.Duration)
	if ok {
		if v1 == nil {
			return true
		}

		v = *v1
	}

	v2, ok := v.(time.Duration)
	if !ok {
		return true
	}

	return v2 == time.Duration(0)
}

func GetOpenAPISchemaString() *Schema {
	return &Schema{
		Type: TypeOfString,
	}
}

func ParseString(v any) (any, error) {
	switch v1 := v.(type) {
	case string:
		return v1, nil
	}

	return "", fmt.Errorf("%#+v (%v) could not be identified for ParseString", v, typeOf(v))
}

func FormatString(v any) (any, error) {
	v1, ok := v.(*string)
	if ok {
		if v1 == nil {
			return nil, nil
		}

		return *v1, nil
	}

	v2, ok := v.(string)
	if !ok {
		return nil, fmt.Errorf("%#+v (%v) could not be cast to string for FormatString", v, typeOf(v))
	}

	return v2, nil
}

func IsZeroString(v any) bool {
	if v == nil {
		return true
	}

	v1, ok := v.(*string)
	if ok {
		if v1 == nil {
			return true
		}

		v = *v1
	}

	v2, ok := v.(string)
	if !ok {
		return true
	}

	return v2 == ""
}

func GetOpenAPISchemaStringArray() *Schema {
	return &Schema{
		Type: TypeOfArray,
		Items: &Schema{
			Type: TypeOfString,
		},
	}
}

func ParseStringArray(v any) (any, error) {
	switch v1 := v.(type) {
	case []byte:
		v2 := pq.StringArray{}
		err := v2.Scan(v1)
		if err != nil {
			return nil, fmt.Errorf("%#+v (%v) could not be parsed with pq.StringArray.Scan for ParseStringArray; err: %v", v1, typeOf(v), err)
		}

		v3 := []string(v2)
		return v3, nil
	case []any:
		temp2 := make([]string, 0)
		for _, v := range v1 {
			s, ok := v.(string)
			if !ok {
				return nil, fmt.Errorf("%#+v (%v) could not be cast to []string for ParseStringArray", v, typeOf(v))
			}

			temp2 = append(temp2, s)
		}

		return temp2, nil
	}

	return nil, fmt.Errorf("%#+v (%v) could not be identified for ParseStringArray", v, typeOf(v))
}

func FormatStringArray(v any) (any, error) {
	format := func(x []string) pq.StringArray {
		return pq.StringArray(x)
	}

	v1, ok := v.(*[]string)
	if ok {
		if v1 == nil {
			return nil, nil
		}

		return format(*v1), nil
	}

	v2, ok := v.([]string)
	if !ok {
		return nil, fmt.Errorf("%#+v (%v) could not be cast to []string for FormatStringArray", v, typeOf(v))
	}

	return format(v2), nil
}

func IsZeroStringArray(v any) bool {
	if v == nil {
		return true
	}

	v1, ok := v.([]string)
	if !ok {
		return true
	}

	return v1 == nil
}

func GetOpenAPISchemaHstore() *Schema {
	return &Schema{
		Type: TypeOfObject,
		AdditionalProperties: &Schema{
			Type:     TypeOfString,
			Nullable: true,
		},
	}
}

func ParseHstore(v any) (any, error) {
	v0, ok := v.(map[string]interface{})
	if ok {
		v1 := make(map[string]*string)
		for v0k, v0v := range v0 {
			if v0v == nil {
				v1[v0k] = nil
				continue
			}

			vs, ok := v0v.(string)
			if !ok {
				return nil, fmt.Errorf("%#+v (%v) could not be cast to string for ParseHstore", v0v, reflect.TypeOf(v0v).String())
			}

			v1[v0k] = &vs
		}

		return v1, nil
	}

	v1, ok := v.([]byte)
	if !ok {
		temp, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("%#+v (%v) could not be cast to string for ParseHstore", v, typeOf(v))
		}

		v1 = []byte(temp)
	}

	v2 := hstore.Hstore{}
	err := v2.Scan(v1)
	if err != nil {
		return nil, fmt.Errorf("%#+v (%v) could not be parsed with pq.StringArray.Scan for ParseHstore; err: %v", v1, typeOf(v), err)
	}

	v3 := make(map[string]*string)

	for k, v := range v2.Map {
		if v.Valid {
			v3[k] = helpers.Ptr(v.String)
		} else {
			v3[k] = helpers.Nil("")
		}
	}

	return v3, nil
}

func FormatHstore(v any) (any, error) {
	format := func(x map[string]*string) hstore.Hstore {
		y := hstore.Hstore{
			Map: make(map[string]sql.NullString),
		}

		for k, v := range x {
			z := sql.NullString{}

			if v != nil {
				z.String = *v
				z.Valid = true
			}

			y.Map[k] = z
		}

		return y
	}

	v1, ok := v.(*map[string]*string)
	if ok {
		if v1 == nil {
			return nil, nil
		}

		return format(*v1), nil
	}

	v2, ok := v.(map[string]*string)
	if !ok {
		return nil, fmt.Errorf("%#+v (%v) could not be cast to map[string]*string for FormatHstore", v, typeOf(v))
	}

	return format(v2), nil
}

func IsZeroHstore(v any) bool {
	if v == nil {
		return true
	}

	v1, ok := v.(map[string]*string)
	if !ok {
		return true
	}

	return v1 == nil
}

func GetOpenAPISchemaJSON() *Schema {
	return &Schema{
		Type:     TypeOfObject,
		Nullable: true,
	}
}

func ParseJSON(v any) (any, error) {
	switch v := v.(type) {
	case []byte:
		var v1 any
		err := json.Unmarshal(v, &v1)
		if err != nil {
			return nil, fmt.Errorf("%#+v (%v) could not be parsed with json.Unmarshal for ParseJSON; err: %v", v1, typeOf(v), err)
		}

		return v1, nil
	case (map[string]any), []map[string]any, []any, any:
		return v, nil
	default:
	}

	return nil, fmt.Errorf("%#+v could not be identified ParseJSON", v)
}

func FormatJSON(v any) (any, error) {
	if v == nil {
		return nil, nil
	}

	v1, ok := v.(*any)
	if ok {
		if v1 == nil {
			return nil, nil
		}

		return json.Marshal(*v1)
	}

	v2, ok := v.(any)
	if !ok {
		return nil, fmt.Errorf("%#+v (%v) could not be cast to any for FormatJSON", v, typeOf(v))
	}

	return json.Marshal(v2)
}

func IsZeroJSON(v any) bool {
	if v == nil {
		return true
	}

	return v == nil
}

func GetOpenAPISchemaInt() *Schema {
	return &Schema{
		Type:   TypeOfInteger,
		Format: FormatOfInt64,
	}
}

func ParseInt(v any) (any, error) {
	if v != nil {
		switch v1 := v.(type) {

		case *int64:
			if v1 == nil {
				return nil, nil
			}

			return int64(*v1), nil
		case *int32:
			if v1 == nil {
				return nil, nil
			}

			return int64(*v1), nil
		case *int16:
			if v1 == nil {
				return nil, nil
			}

			return int64(*v1), nil
		case *int8:
			if v1 == nil {
				return nil, nil
			}

			return int64(*v1), nil
		case *uint64:
			if v1 == nil {
				return nil, nil
			}

			return int64(*v1), nil
		case *uint32:
			if v1 == nil {
				return nil, nil
			}

			return int64(*v1), nil
		case *uint16:
			if v1 == nil {
				return nil, nil
			}

			return int64(*v1), nil
		case *uint8:
			if v1 == nil {
				return nil, nil
			}

			return int64(*v1), nil
		case *float64:
			if v1 == nil {
				return nil, nil
			}

			return int64(*v1), nil
		case *float32:
			if v1 == nil {
				return nil, nil
			}

			return int64(*v1), nil

		case int64:
			return int64(v1), nil
		case int32:
			return int64(v1), nil
		case int16:
			return int64(v1), nil
		case int8:
			return int64(v1), nil
		case uint64:
			return int64(v1), nil
		case uint32:
			return int64(v1), nil
		case uint16:
			return int64(v1), nil
		case uint8:
			return int64(v1), nil
		case float64:
			return int64(v1), nil
		case float32:
			return int64(v1), nil
		}
	}

	return 0, fmt.Errorf("%#+v (%v) could not be identified for ParseInt", v, typeOf(v))
}

func FormatInt(v any) (any, error) {
	if v != nil {
		switch v1 := v.(type) {
		case *int64:
			if v1 == nil {
				return nil, nil
			}
			return int64(*v1), nil
		case *int32:
			if v1 == nil {
				return nil, nil
			}
			return int64(*v1), nil
		case *int16:
			if v1 == nil {
				return nil, nil
			}
			return int64(*v1), nil
		case *int8:
			if v1 == nil {
				return nil, nil
			}
			return int64(*v1), nil
		case *uint64:
			if v1 == nil {
				return nil, nil
			}
			return int64(*v1), nil
		case *uint32:
			if v1 == nil {
				return nil, nil
			}
			return int64(*v1), nil
		case *uint16:
			if v1 == nil {
				return nil, nil
			}
			return int64(*v1), nil
		case *uint8:
			if v1 == nil {
				return nil, nil
			}
			return int64(*v1), nil
		case *float64:
			if v1 == nil {
				return nil, nil
			}
			return int64(*v1), nil
		case *float32:
			if v1 == nil {
				return nil, nil
			}
			return int64(*v1), nil

		case int64:
			return int64(v1), nil
		case int32:
			return int64(v1), nil
		case int16:
			return int64(v1), nil
		case int8:
			return int64(v1), nil
		case uint64:
			return int64(v1), nil
		case uint32:
			return int64(v1), nil
		case uint16:
			return int64(v1), nil
		case uint8:
			return int64(v1), nil
		case float64:
			return int64(v1), nil
		case float32:
			return int64(v1), nil
		}
	}

	return nil, fmt.Errorf("%#+v (%v) could not be cast to int64 for FormatInt", v, typeOf(v))
}

func IsZeroInt(v any) bool {
	if v != nil {
		switch v1 := v.(type) {
		case *int64:
			if v1 == nil {
				break
			}

			return int64(*v1) == 0
		case *int32:
			if v1 == nil {
				break
			}

			return int64(*v1) == 0
		case *int16:
			if v1 == nil {
				break
			}

			return int64(*v1) == 0
		case *int8:
			if v1 == nil {
				break
			}

			return int64(*v1) == 0
		case *uint64:
			if v1 == nil {
				break
			}

			return int64(*v1) == 0
		case *uint32:
			if v1 == nil {
				break
			}

			return int64(*v1) == 0
		case *uint16:
			if v1 == nil {
				break
			}

			return int64(*v1) == 0
		case *uint8:
			if v1 == nil {
				break
			}

			return int64(*v1) == 0
		case *float64:
			if v1 == nil {
				break
			}

			return int64(*v1) == 0
		case *float32:
			if v1 == nil {
				break
			}

			return int64(*v1) == 0

		case int64:
			return int64(v1) == 0
		case int32:
			return int64(v1) == 0
		case int16:
			return int64(v1) == 0
		case int8:
			return int64(v1) == 0
		case uint64:
			return int64(v1) == 0
		case uint32:
			return int64(v1) == 0
		case uint16:
			return int64(v1) == 0
		case uint8:
			return int64(v1) == 0
		case float64:
			return int64(v1) == 0
		case float32:
			return int64(v1) == 0
		}
	}

	return true
}

func GetOpenAPISchemaFloat() *Schema {
	return &Schema{
		Type:   TypeOfNumber,
		Format: FormatOfDouble,
	}
}

func ParseFloat(v any) (any, error) {
	if v != nil {
		switch v1 := v.(type) {
		case *int64:
			if v1 == nil {
				return nil, nil
			}

			return float64(*v1), nil
		case *int32:
			if v1 == nil {
				return nil, nil
			}

			return float64(*v1), nil
		case *int16:
			if v1 == nil {
				return nil, nil
			}

			return float64(*v1), nil
		case *int8:
			if v1 == nil {
				return nil, nil
			}

			return float64(*v1), nil
		case *uint64:
			if v1 == nil {
				return nil, nil
			}

			return float64(*v1), nil
		case *uint32:
			if v1 == nil {
				return nil, nil
			}

			return float64(*v1), nil
		case *uint16:
			if v1 == nil {
				return nil, nil
			}

			return float64(*v1), nil
		case *uint8:
			if v1 == nil {
				return nil, nil
			}

			return float64(*v1), nil
		case *float64:
			if v1 == nil {
				return nil, nil
			}

			return float64(*v1), nil
		case *float32:
			if v1 == nil {
				return nil, nil
			}

			return float64(*v1), nil

		case int64:
			return float64(v1), nil
		case int32:
			return float64(v1), nil
		case int16:
			return float64(v1), nil
		case int8:
			return float64(v1), nil
		case uint64:
			return float64(v1), nil
		case uint32:
			return float64(v1), nil
		case uint16:
			return float64(v1), nil
		case uint8:
			return float64(v1), nil
		case float64:
			return float64(v1), nil
		case float32:
			return float64(v1), nil
		}
	}

	return 0, fmt.Errorf("%#+v (%v) could not be identified for ParseFloat", v, typeOf(v))
}

func FormatFloat(v any) (any, error) {
	if v != nil {
		switch v1 := v.(type) {
		case *int64:
			if v1 == nil {
				return nil, nil
			}

			return float64(*v1), nil
		case *int32:
			if v1 == nil {
				return nil, nil
			}

			return float64(*v1), nil
		case *int16:
			if v1 == nil {
				return nil, nil
			}

			return float64(*v1), nil
		case *int8:
			if v1 == nil {
				return nil, nil
			}

			return float64(*v1), nil
		case *uint64:
			if v1 == nil {
				return nil, nil
			}

			return float64(*v1), nil
		case *uint32:
			if v1 == nil {
				return nil, nil
			}

			return float64(*v1), nil
		case *uint16:
			if v1 == nil {
				return nil, nil
			}

			return float64(*v1), nil
		case *uint8:
			if v1 == nil {
				return nil, nil
			}

			return float64(*v1), nil
		case *float64:
			if v1 == nil {
				return nil, nil
			}

			return float64(*v1), nil
		case *float32:
			if v1 == nil {
				return nil, nil
			}

			return float64(*v1), nil

		case int64:
			return float64(v1), nil
		case int32:
			return float64(v1), nil
		case int16:
			return float64(v1), nil
		case int8:
			return float64(v1), nil
		case uint64:
			return float64(v1), nil
		case uint32:
			return float64(v1), nil
		case uint16:
			return float64(v1), nil
		case uint8:
			return float64(v1), nil
		case float64:
			return float64(v1), nil
		case float32:
			return float64(v1), nil
		}
	}

	return nil, fmt.Errorf("%#+v (%v) could not be cast to int64 for FormatFloat", v, typeOf(v))
}

func IsZeroFloat(v any) bool {
	if v != nil {
		switch v1 := v.(type) {

		case *int64:
			if v1 == nil {
				break
			}

			return float64(*v1) == 0.0
		case *int32:
			if v1 == nil {
				break
			}

			return float64(*v1) == 0.0
		case *int16:
			if v1 == nil {
				break
			}

			return float64(*v1) == 0.0
		case *int8:
			if v1 == nil {
				break
			}

			return float64(*v1) == 0.0
		case *uint64:
			if v1 == nil {
				break
			}

			return float64(*v1) == 0.0
		case *uint32:
			if v1 == nil {
				break
			}

			return float64(*v1) == 0.0
		case *uint16:
			if v1 == nil {
				break
			}

			return float64(*v1) == 0.0
		case *uint8:
			if v1 == nil {
				break
			}

			return float64(*v1) == 0.0
		case *float64:
			if v1 == nil {
				break
			}

			return float64(*v1) == 0.0
		case *float32:
			if v1 == nil {
				break
			}

			return float64(*v1) == 0.0

		case int64:
			return float64(v1) == 0.0
		case int32:
			return float64(v1) == 0.0
		case int16:
			return float64(v1) == 0.0
		case int8:
			return float64(v1) == 0.0
		case uint64:
			return float64(v1) == 0.0
		case uint32:
			return float64(v1) == 0.0
		case uint16:
			return float64(v1) == 0.0
		case uint8:
			return float64(v1) == 0.0
		case float64:
			return float64(v1) == 0.0
		case float32:
			return float64(v1) == 0.0
		}
	}

	return true
}

func GetOpenAPISchemaBool() *Schema {
	return &Schema{
		Type: TypeOfBoolean,
	}
}

func ParseBool(v any) (any, error) {
	switch v1 := v.(type) {
	case bool:
		return v1, nil
	}

	return 0, fmt.Errorf("%#+v (%v) could not be identified for ParseBool", v, typeOf(v))
}

func FormatBool(v any) (any, error) {
	v1, ok := v.(*bool)
	if ok {
		if v1 == nil {
			return nil, nil
		}

		return *v1, nil
	}

	v2, ok := v.(bool)
	if !ok {
		return nil, fmt.Errorf("%#+v (%v) could not be cast to bool for FormatBool", v, typeOf(v))
	}

	return v2, nil
}

func IsZeroBool(v any) bool {
	if v == nil {
		return true
	}

	v1, ok := v.(*bool)
	if ok {
		if v1 == nil {
			return true
		}

		v = *v1
	}

	v2, ok := v.(bool)
	if !ok {
		return true
	}

	return !v2
}

func GetOpenAPISchemaTSVector() *Schema {
	return &Schema{
		Type: TypeOfObject,
		AdditionalProperties: &Schema{
			Type: TypeOfArray,
			Items: &Schema{
				Type:   TypeOfInteger,
				Format: FormatOfInt32,
			},
		},
	}
}

func ParseTSVector(v any) (any, error) {
	switch v1 := v.(type) {
	case tsvector.TSVector:
		return v1, nil
	case string:
		v2 := tsvector.TSVector{}
		err := v2.Scan([]byte(v1))
		if err != nil {
			return nil, fmt.Errorf("%#+v (%v) could not be parsed with tsvector.TSVector.Scan for ParseTSVector; err: %v", v1, typeOf(v), err)
		}

		return v2.Lexemes(), nil
	case []byte:
		v2 := tsvector.TSVector{}
		err := v2.Scan(v1)
		if err != nil {
			return nil, fmt.Errorf("%#+v (%v) could not be parsed with tsvector.TSVector.Scan for ParseTSVector; err: %v", v1, typeOf(v), err)
		}

		return v2.Lexemes(), nil
	}

	return nil, fmt.Errorf("%#+v (%v) could not be identified for ParseTSVector", v, typeOf(v))
}

func FormatTSVector(v any) (any, error) {
	v1, ok := v.(*tsvector.TSVector)
	if ok {
		if v1 == nil {
			return nil, nil
		}

		return *v1, nil
	}

	v2, ok := v.(tsvector.TSVector)
	if !ok {
		return nil, fmt.Errorf("%#+v (%v) could not be cast to tsvector.TSVector for FormatTSVector", v, typeOf(v))
	}

	return v2, nil
}

func IsZeroTSVector(v any) bool {
	if v == nil {
		return true
	}

	v1, ok := v.(*bool)
	if ok {
		if v1 == nil {
			return true
		}

		v = *v1
	}

	v2, ok := v.(bool)
	if !ok {
		return true
	}

	return !v2
}

func GetOpenAPISchemaPoint() *Schema {
	return &Schema{
		Type: TypeOfObject,
		Properties: map[string]*Schema{
			"P": {
				Type: TypeOfObject,
				Properties: map[string]*Schema{
					"X": {
						Type:   TypeOfNumber,
						Format: FormatOfDouble,
					},
					"Y": {
						Type:   TypeOfNumber,
						Format: FormatOfDouble,
					},
				},
			},
		},
	}
}

func ParsePoint(v any) (any, error) {
	switch v1 := v.(type) {
	case map[string]any:
		rawRawP, ok := v1["P"]
		if !ok {
			return pgtype.Vec2{}, fmt.Errorf("%#+v (%v) could not be parsed as pgtype.Point for ParsePoint; err: %v", v1, typeOf(v), fmt.Errorf("missing 'P' key"))
		}

		rawP, ok := rawRawP.(map[string]any)
		if !ok {
			return pgtype.Vec2{}, fmt.Errorf("%#+v (%v) could not be parsed as pgtype.Point for ParsePoint; err: %v", v1, typeOf(v), fmt.Errorf("bad type for 'P' key"))
		}

		rawX, ok := rawP["X"]
		if !ok {
			return pgtype.Vec2{}, fmt.Errorf("%#+v (%v) could not be parsed as pgtype.Point for ParsePoint; err: %v", v1, typeOf(v), fmt.Errorf("missing 'X' key"))
		}

		x, ok := rawX.(float64)
		if !ok {
			return pgtype.Vec2{}, fmt.Errorf("%#+v (%v) could not be parsed as pgtype.Point for ParsePoint; err: %v", v1, typeOf(v), fmt.Errorf("bad type for 'X' key"))
		}

		rawY, ok := rawP["Y"]
		if !ok {
			return pgtype.Vec2{}, fmt.Errorf("%#+v (%v) could not be parsed as pgtype.Point for ParsePoint; err: %v", v1, typeOf(v), fmt.Errorf("missing 'Y' key"))
		}

		y, ok := rawY.(float64)
		if !ok {
			return pgtype.Vec2{}, fmt.Errorf("%#+v (%v) could not be parsed as pgtype.Point for ParsePoint; err: %v", v1, typeOf(v), fmt.Errorf("bad type for 'Y' key"))
		}

		p := pgtype.Vec2{
			X: x,
			Y: y,
		}

		return pgtype.Point{
			P:     p,
			Valid: true,
		}, nil
	case []byte:
		v2 := pgtype.Point{}
		err := v2.UnmarshalJSON(v1)
		if err != nil {
			return pgtype.Vec2{}, fmt.Errorf("%#+v (%v) could not be parsed with pgtype.Point.UnmarshalJSON for ParsePoint; err: %v", v1, typeOf(v), err)
		}

		return v2, nil
	case pgtype.Point:
		return v1, nil
	}

	return pgtype.Point{}, fmt.Errorf("%#+v (%v) could not be identified for ParsePoint", v, typeOf(v))
}

func FormatPoint(v any) (any, error) {
	switch v1 := v.(type) {
	case *pgtype.Point:
		return *v1, nil
	case pgtype.Point:
		return v1, nil
	case *pgtype.Vec2:
		return pgtype.Point{
			P:     *v1,
			Valid: true,
		}, nil
	case pgtype.Vec2:
		return pgtype.Point{
			P:     v1,
			Valid: true,
		}, nil
	}

	return nil, fmt.Errorf("%#+v (%v) could not be cast to pgtype.Vec2 for FormatPoint", v, typeOf(v))
}

func IsZeroPoint(v any) bool {
	if v == nil {
		return true
	}

	v1, ok := v.(*pgtype.Point)
	if ok {
		if v1 == nil {
			return true
		}

		return v1.P == pgtype.Vec2{}
	}

	v2, ok := v.(pgtype.Point)
	if ok {
		return v2.P == pgtype.Vec2{}
	}

	v3, ok := v.(*pgtype.Vec2)
	if ok {
		if v3 == nil {
			return true
		}

		return *v3 == pgtype.Vec2{}
	}

	v4, ok := v.(pgtype.Vec2)
	if !ok {
		return v4 == pgtype.Vec2{}
	}

	return false
}

func GetOpenAPISchemaPolygon() *Schema {
	return &Schema{
		Type:  TypeOfArray,
		Items: GetOpenAPISchemaPoint(),
	}
}

func ParsePolygon(v any) (any, error) {
	switch v1 := v.(type) {
	case []any:
		v2 := pgtype.Polygon{
			P:     []pgtype.Vec2{},
			Valid: false,
		}

		for _, item := range v1 {
			rawPoint, err := ParsePoint(item)
			if err != nil {
				return pgtype.Polygon{}, fmt.Errorf("%#+v (%v) could not be parsed by item for ParsePolygon; err: %v", v1, typeOf(v), err)
			}

			point, ok := rawPoint.(pgtype.Point)
			if !ok {
				return pgtype.Polygon{}, fmt.Errorf("%#+v (%v) could not be parsed by item for ParsePolygon; err: %v", v1, typeOf(v), fmt.Errorf("failed cast to pgtype.Point"))
			}

			v2.P = append(v2.P, point.P)
		}

		v2.Valid = len(v2.P) > 0

		return v2, nil
	case []byte:
		v2 := _pgtype.Polygon{}
		err := v2.Scan(v1)
		if err != nil {
			return pgtype.Polygon{}, fmt.Errorf("%#+v (%v) could not be parsed with _pgtype.Polygon.Scan for ParsePolygon; err: %v", v1, typeOf(v), err)
		}

		v3 := pgtype.Polygon{
			P:     make([]pgtype.Vec2, 0),
			Valid: v2.Status == _pgtype.Present,
		}

		for _, p := range v2.P {
			v3.P = append(v3.P, pgtype.Vec2{X: p.X, Y: p.Y})
		}

		return v3, nil
	case pgtype.Polygon:
		return v1, nil
	}

	return nil, fmt.Errorf("%#+v (%v) could not be identified for ParsePolygon", v, typeOf(v))
}

func FormatPolygon(v any) (any, error) {
	switch v1 := v.(type) {
	case *pgtype.Polygon:
		return *v1, nil
	case pgtype.Polygon:
		return v, nil
	case *[]pgtype.Vec2:
		return pgtype.Polygon{
			P:     *v1,
			Valid: true,
		}, nil
	case []pgtype.Vec2:
		return pgtype.Polygon{
			P:     v1,
			Valid: true,
		}, nil

	}

	return nil, fmt.Errorf("%#+v (%v) could not be cast to []pgtype.Vec2 for FormatPolygon", v, typeOf(v))
}

func IsZeroPolygon(v any) bool {
	if v == nil {
		return true
	}

	v1, ok := v.(*pgtype.Polygon)
	if ok {
		if v1 == nil {
			return true
		}

		v = *v1
	}

	v2, ok := v.(pgtype.Polygon)
	if !ok {
		return true
	}

	return v2.P == nil
}

func GetOpenAPISchemaGeometry() *Schema {
	return &Schema{
		Type: TypeOfObject,
		Properties: map[string]*Schema{
			"X": {
				Type:   TypeOfNumber,
				Format: FormatOfDouble,
			},
			"Y": {
				Type:   TypeOfNumber,
				Format: FormatOfDouble,
			},
			"Z": {
				Type:   TypeOfNumber,
				Format: FormatOfDouble,
			},
		},
	}
}

func ParseGeometry(v any) (any, error) {
	switch v1 := v.(type) {
	case string:
		v2 := postgis.PointZ{}
		err := v2.Scan([]byte(v1))
		if err != nil {
			return nil, fmt.Errorf("%#+v (%v) could not be parsed with ewkbhex.Decode for ParseGeometry; err: %v", v1, typeOf(v), err)
		}

		return v2, nil
	}

	return nil, fmt.Errorf("%#+v (%v) could not be identified for ParseGeometry", v, typeOf(v))
}

func FormatGeometry(v any) (any, error) {
	v1, ok := v.(*ewkb.Point)
	if ok {
		if v1 == nil {
			return nil, nil
		}

		return *v1, nil
	}

	v2, ok := v.(ewkb.Point)
	if !ok {
		return nil, fmt.Errorf("%#+v (%v) could not be cast to ewkb.Point for FormatGeometry", v, typeOf(v))
	}

	return v2, nil
}

func IsZeroGeometry(v any) bool {
	if v == nil {
		return true
	}

	v1, ok := v.(*geom.T)
	if ok {
		if v1 == nil {
			return true
		}

		v = *v1
	}

	v2, ok := v.(geom.T)
	if !ok {
		return true
	}

	return v2 == nil || v2.Empty()
}

func GetOpenAPISchemaInet() *Schema {
	return &Schema{
		Type:   TypeOfString,
		Format: FormatOfIPv4,
	}
}

func ParseInet(v any) (any, error) {
	switch v1 := v.(type) {
	case netip.Prefix:
		return v1, nil
	case []byte:
		v2 := string(v1)
		if !strings.Contains(v2, "/") {
			v2 += "/32"
		}

		v3, err := netip.ParsePrefix(v2)
		if err != nil {
			return netip.Prefix{}, fmt.Errorf("%#+v (%v) could not be parsed with netip.ParsePrefix for ParseInet; err: %v", v3, typeOf(v), err)
		}

		return v3, nil
	}

	return netip.Prefix{}, fmt.Errorf("%#+v (%v) could not be identified for ParseInet", v, typeOf(v))
}

func FormatInet(v any) (any, error) {
	v1, ok := v.(*netip.Prefix)
	if ok {
		if v1 == nil {
			return nil, nil
		}

		return *v1, nil
	}

	v2, ok := v.(netip.Prefix)
	if !ok {
		return nil, fmt.Errorf("%#+v (%v) could not be cast to netip.Prefix for FormatInet", v, typeOf(v))
	}

	return v2, nil
}

func IsZeroInet(v any) bool {
	if v == nil {
		return true
	}

	v1, ok := v.(*netip.Prefix)
	if ok {
		if v1 == nil {
			return true
		}

		v = *v1
	}

	v2, ok := v.(netip.Prefix)
	if !ok {
		return true
	}

	return v2.IsValid()
}

func GetOpenAPISchemaBytes() *Schema {
	return &Schema{
		Type:   TypeOfString,
		Format: FormatOfByte,
	}
}

func ParseBytes(v any) (any, error) {
	switch v1 := v.(type) {
	case []byte:
		return v1, nil
	}

	return 0, fmt.Errorf("%#+v (%v) could not be identified for ParseBytes", v, typeOf(v))
}

func FormatBytes(v any) (any, error) {
	v1, ok := v.(*[]byte)
	if ok {
		if v1 == nil {
			return nil, nil
		}

		return *v1, nil
	}

	v2, ok := v.([]byte)
	if !ok {
		return nil, fmt.Errorf("%#+v (%v) could not be cast to []byte for FormatBytes", v, typeOf(v))
	}

	return v2, nil
}

func IsZeroBytes(v any) bool {
	if v == nil {
		return true
	}

	v1, ok := v.([]byte)
	if !ok {
		return true
	}

	return len(v1) == 0
}
