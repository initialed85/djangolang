package server

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/initialed85/djangolang/pkg/introspect"
)

const (
	browserCacheDebounce = time.Second
	browserCacheRecent   = time.Minute
	browserCacheHistoric = time.Hour
)

// browserCacheMaxAge applies the browser debounce heuristic to successful GETs.
// A bound within an hour of now is relatively live; older timestamp windows are
// considered historical. Greater-than bounds are kept very short because new
// rows can continuously enter their result set.
func browserCacheMaxAge(method string, status int, queryParams map[string]any, table *introspect.Table, now time.Time) (time.Duration, bool) {
	if method != http.MethodGet || status < http.StatusOK || status >= http.StatusMultipleChoices {
		return 0, false
	}

	hasTimestampFilter := false
	hasRecentTimestamp := false

	for key, value := range queryParams {
		parts := strings.Split(key, "__")
		if len(parts) != 2 || table == nil {
			continue
		}

		column := table.ColumnByName[parts[0]]
		if column == nil || column.TypeTemplate != "time.Time" {
			continue
		}

		switch parts[1] {
		case "gt", "gte":
			return browserCacheDebounce, true
		case "eq", "ne", "lt", "lte", "in", "notin":
			hasTimestampFilter = true
			timestamps, ok := timestampValues(value)
			if !ok {
				// If we cannot classify a timestamp filter's value, prefer a
				// short debounce rather than allowing a long-lived stale result.
				return browserCacheDebounce, true
			}
			for _, timestamp := range timestamps {
				if !timestamp.Before(now.Add(-time.Hour)) && !timestamp.After(now.Add(time.Hour)) {
					hasRecentTimestamp = true
				}
			}
		case "isnull", "isnotnull":
			// These filters have no timestamp value to judge, and their result
			// may change as rows are inserted or updated.
			return browserCacheDebounce, true
		}
	}

	if !hasTimestampFilter {
		return browserCacheDebounce, true
	}
	if hasRecentTimestamp {
		return browserCacheRecent, true
	}
	return browserCacheHistoric, true
}

func timestampValues(value any) ([]time.Time, bool) {
	switch typedValue := value.(type) {
	case time.Time:
		return []time.Time{typedValue}, true
	case string:
		rawValues := strings.Split(typedValue, ",")
		timestamps := make([]time.Time, 0, len(rawValues))
		for _, rawValue := range rawValues {
			rawValue = strings.Trim(strings.TrimSpace(rawValue), `"`)
			rawValue = strings.ReplaceAll(rawValue, " ", "+")
			timestamp, err := time.Parse(time.RFC3339Nano, rawValue)
			if err != nil {
				timestamp, err = time.Parse(time.RFC3339, rawValue)
			}
			if err != nil {
				return nil, false
			}
			timestamps = append(timestamps, timestamp)
		}
		return timestamps, len(timestamps) > 0
	case []any:
		timestamps := make([]time.Time, 0, len(typedValue))
		for _, item := range typedValue {
			values, ok := timestampValues(item)
			if !ok {
				return nil, false
			}
			timestamps = append(timestamps, values...)
		}
		return timestamps, len(timestamps) > 0
	default:
		return nil, false
	}
}

func setBrowserCacheHeaders(w http.ResponseWriter, method string, status int, queryParams map[string]any, table *introspect.Table, now time.Time) {
	maxAge, ok := browserCacheMaxAge(method, status, queryParams, table, now)
	if !ok {
		return
	}
	w.Header().Set("Cache-Control", fmt.Sprintf("max-age=%d", int(maxAge/time.Second)))
}
