package helpers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gomodule/redigo/redis"
)

type key struct {
	Name string
}

func GetRedisURL() string {
	redisURL := GetEnvironmentVariable("REDIS_URL")
	if redisURL == "" {
		log.Printf("REDIS_URL env var empty or unset; caching will be disabled")
	}

	return redisURL
}

func GetRequestHash(tableName string, wheres []string, orderBy string, limit int, offset int, values []any, primaryKey any) (string, error) {
	params := make(map[string]any)
	params["__order_by"] = orderBy
	params["__limit"] = limit
	params["__offset"] = offset

	i := 0
	for _, v := range wheres {
		if !strings.Contains(v, "$$??") {
			params[v] = nil
			continue
		}

		params[v] = values[i]

		i++
	}

	b, err := json.Marshal(params)
	if err != nil {
		return "", fmt.Errorf("failed to build request hash for params: %#+v: %v", params, err)
	}

	if primaryKey == nil {
		primaryKey = ""
	}

	return fmt.Sprintf("%v:%v:%v", tableName, primaryKey, string(b)), nil
}

func AttemptCachedResponse(requestHash string, redisConn redis.Conn, w http.ResponseWriter) (bool, error) {
	if redisConn == nil {
		w.Header().Add("X-Djangolang-Cache-Status", "disabled")
		return false, nil
	}

	cachedObjectsAsStringOfJSON, err := redis.String(redisConn.Do("GET", requestHash))
	if err != nil && !errors.Is(err, redis.ErrNil) {
		w.Header().Add("X-Djangolang-Cache-Status", "error")
		HandleErrorResponse(w, http.StatusInternalServerError, err)
		return false, err
	}

	if !errors.Is(err, redis.ErrNil) {
		w.Header().Add("X-Djangolang-Cache-Status", "hit")
		WriteResponse(w, http.StatusOK, []byte(cachedObjectsAsStringOfJSON))
		return true, nil
	}

	w.Header().Add("X-Djangolang-Cache-Status", "miss")
	return false, nil
}

func StoreCachedResponse(requestHash string, redisConn redis.Conn, returnedObjectsAsJSON string) error {
	if redisConn == nil {
		return nil
	}

	_, err := redisConn.Do("SET", requestHash, string(returnedObjectsAsJSON))
	if err != nil {
		return fmt.Errorf("failed to set value for cache item %v = %v: %v", requestHash, string(returnedObjectsAsJSON), err)
	}

	_, err = redisConn.Do("EXPIRE", requestHash, "86400")
	if err != nil {
		return fmt.Errorf("failed to set expiry for cache item %v = %v failed: %v", requestHash, string(returnedObjectsAsJSON), err)
	}

	return nil
}

func GetQueryHash(columns []string, table string, where string, orderBy *string, limit *int, offset *int, values ...any,
) (string, error) {
	seed := map[string]any{
		"columns": columns,
		"table":   table,
		"where":   where,
		"orderBy": orderBy,
		"limit":   limit,
		"offset":  offset,
		"values":  values,
	}

	queryHash, err := json.Marshal(seed)
	if err != nil {
		return "", fmt.Errorf("failed to get query hash for seed %#+v: %v", seed, err)
	}

	return string(queryHash), nil
}

func AttemptCachedQuery(ctx context.Context, queryHash string) (bool, context.Context, []map[string]any, error) {
	rawCachedItems := ctx.Value(key{Name: queryHash})
	if rawCachedItems != nil {
		cachedItems, ok := rawCachedItems.([]map[string]any)
		if !ok {
			return false, ctx, nil, fmt.Errorf("failed to cast %#+v to %#+v", rawCachedItems, nil)
		}

		return true, ctx, cachedItems, nil
	}

	return false, ctx, nil, nil
}

func StoreCachedItems(ctx context.Context, queryHash string, items []map[string]any) context.Context {
	return context.WithValue(ctx, key{Name: queryHash}, items)
}
