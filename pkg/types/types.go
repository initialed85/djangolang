package types

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/google/uuid"
	"github.com/initialed85/djangolang/pkg/helpers"
	"github.com/lib/pq"
	"github.com/lib/pq/hstore"
	"github.com/twpayne/go-geom"
)

func ParsePtr[T any](ParseFunc func(any) (T, error), v any) (*T, error) {
	if v == nil {
		return nil, nil
	}

	v1, err := ParseFunc(v)
	if err != nil {
		return nil, err
	}

	return &v1, nil
}

func ParseUUID(v any) (uuid.UUID, error) {
	if v == nil {
		return uuid.UUID{}, nil
	}

	switch v1 := v.(type) {
	case []byte:
		v2, err := uuid.Parse(string(v1))
		if err != nil {
			return uuid.UUID{}, fmt.Errorf("%#+v (%v) could not be parsed with uuid.Parse for ParseUUID; err: %v", v, reflect.TypeOf(v).String(), err)
		}

		return v2, nil
	case [16]byte:
		v2, err := uuid.FromBytes(v1[:])
		if err != nil {
			return uuid.UUID{}, fmt.Errorf("%#+v (%v) could not be parsed with uuid.ParseBytes for ParseUUID; err: %v", v, reflect.TypeOf(v).String(), err)
		}

		return v2, nil
	}

	return uuid.UUID{}, fmt.Errorf("%#+v (%v) could not be identified for ParseTime", v, reflect.TypeOf(v).String())
}

func ParseTime(v any) (time.Time, error) {
	if v == nil {
		return time.Time{}, nil
	}

	v1, ok := v.(time.Time)
	if !ok {
		return time.Time{}, fmt.Errorf("%#+v (%v) could not be cast to time.Time for ParseTime", v, reflect.TypeOf(v).String())
	}

	return v1, nil
}

func ParseString(v any) (string, error) {
	if v == nil {
		return "", nil
	}

	v1, ok := v.(string)
	if !ok {
		return "", fmt.Errorf("%#+v (%v) could not be cast to string for ParseString", v, reflect.TypeOf(v).String())
	}

	return v1, nil
}

func ParseStringArray(v any) ([]string, error) {
	if v == nil {
		return nil, nil
	}

	v1, ok := v.([]byte)
	if !ok {
		temp1, ok := v.([]any)
		if !ok {
			return nil, fmt.Errorf("%#+v (%v) could not be cast to []byte for ParseStringArray", v, reflect.TypeOf(v).String())
		}

		temp2 := make([]string, 0)
		for _, v := range temp1 {
			s, ok := v.(string)
			if !ok {
				return nil, fmt.Errorf("%#+v (%v) could not be cast to []string for ParseStringArray", v, reflect.TypeOf(v).String())
			}

			temp2 = append(temp2, s)
		}

		return temp2, nil
	}

	v2 := pq.StringArray{}
	err := v2.Scan(v1)
	if err != nil {
		return nil, fmt.Errorf("%#+v (%v) could not be parsed with pq.StringArray.Scan for ParseStringArray; err: %v", v1, reflect.TypeOf(v).String(), err)
	}

	v3 := []string(v2)

	return v3, nil
}

func ParseHstore(v any) (map[string]*string, error) {
	if v == nil {
		return nil, nil
	}

	v1, ok := v.([]byte)
	if !ok {
		temp, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("%#+v (%v) could not be cast to string for ParseHstore", v, reflect.TypeOf(v).String())
		}

		v1 = []byte(temp)
	}

	v2 := hstore.Hstore{}
	err := v2.Scan(v1)
	if err != nil {
		return nil, fmt.Errorf("%#+v (%v) could not be parsed with pq.StringArray.Scan for ParseHstore; err: %v", v1, reflect.TypeOf(v).String(), err)
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

func ParseJSON(v any) (any, error) {
	if v == nil {
		return nil, nil
	}

	var v1 any

	switch v := v.(type) {
	case []byte:
		err := json.Unmarshal(v, &v1)
		if err != nil {
			return nil, fmt.Errorf("%#+v (%v) could not be parsed with json.Unmarshal for ParseJSON; err: %v", v1, reflect.TypeOf(v).String(), err)
		}

		return v1, nil
	case (map[string]any), []map[string]any, []any, any:
		return v, nil
	default:
	}

	return nil, fmt.Errorf("%#+v could not be identified ParseJSON", v)
}

func ParseInt(v any) (int64, error) {
	if v == nil {
		return 0, nil
	}

	v1, ok := v.(int64)
	if !ok {
		return 0, fmt.Errorf("%#+v (%v) could not be cast to int64 for ParseInt", v, reflect.TypeOf(v).String())
	}

	return v1, nil
}

func ParseFloat(v any) (float64, error) {
	if v == nil {
		return 0, nil
	}

	v1, ok := v.(float64)
	if !ok {
		return 0, fmt.Errorf("%#+v (%v) could not be cast to float64 for ParseFloat", v, reflect.TypeOf(v).String())
	}

	return v1, nil
}

func ParseBool(v any) (bool, error) {
	if v == nil {
		return false, nil
	}

	v1, ok := v.(bool)
	if !ok {
		return false, fmt.Errorf("%#+v (%v) could not be cast to bool for ParseBool", v, reflect.TypeOf(v).String())
	}

	return v1, nil
}

func ParsePoint(v any) (geom.Point, error) {
	if v == nil {
		return geom.Point{}, nil
	}

	v1, ok := v.([]byte)
	if !ok {
		return geom.Point{}, fmt.Errorf("%#+v (%v) could not be cast to []byte for ParsePoint", v, reflect.TypeOf(v).String())
	}

	_ = v1

	v2 := &geom.Point{}

	return *v2, nil
}

func ParsePolygon(v any) (geom.Polygon, error) {
	if v == nil {
		return geom.Polygon{}, nil
	}

	v1, ok := v.([]byte)
	if !ok {
		return geom.Polygon{}, fmt.Errorf("%#+v (%v) could not be cast to []byte for ParsePolygon", v, reflect.TypeOf(v).String())
	}

	_ = v1

	v2 := &geom.Polygon{}

	return *v2, nil
}
