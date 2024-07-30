package helpers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gomodule/redigo/redis"
)

func GetRedisURL() string {
	redisURL := strings.TrimSpace(os.Getenv("REDIS_URL"))

	return redisURL
}

func GetRequestHash(tableName string, wheres []string, limit int, offset int, values []any, primaryKey any) (string, error) {
	params := make(map[string]any)
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

	_, err = redisConn.Do("EXPIRE", requestHash, 86400)
	if err != nil {
		return fmt.Errorf("failed to set expiry for cache item %v = %v failed: %v", requestHash, string(returnedObjectsAsJSON), err)
	}

	return nil
}
