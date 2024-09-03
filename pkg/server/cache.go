package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gomodule/redigo/redis"
)

func GetRequestHash(tableName string, wheres []string, orderBy string, limit int, offset int, depth int, values []any, primaryKey any) (string, error) {
	params := make(map[string]any)
	params["__order_by"] = orderBy
	params["__limit"] = limit
	params["__offset"] = offset
	params["__depth"] = depth

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
		return "", fmt.Errorf("failed to build request hash for params: %#+v; %v", params, err)
	}

	if primaryKey == nil {
		primaryKey = ""
	}

	return fmt.Sprintf("%v:%v:%v", tableName, primaryKey, string(b)), nil
}

func GetCachedResponseAsJSON(requestHash string, redisConn redis.Conn) ([]byte, bool, error) {
	if redisConn == nil {
		return nil, false, nil
	}

	cachedResponse, err := redis.String(redisConn.Do("GET", requestHash))
	if err != nil && !errors.Is(err, redis.ErrNil) {
		return nil, false, err
	}

	if !errors.Is(err, redis.ErrNil) {
		return []byte(cachedResponse), true, nil
	}

	return nil, false, nil
}

func AttemptCachedResponse(requestHash string, redisConn redis.Conn, w http.ResponseWriter) (bool, error) {
	if redisConn == nil {
		w.Header().Add("X-Djangolang-Cache-Status", "disabled")
		return false, nil
	}

	cachedResponseAsJSON, ok, err := GetCachedResponseAsJSON(requestHash, redisConn)
	if err != nil {
		w.Header().Add("X-Djangolang-Cache-Status", "error")
		HandleErrorResponse(w, http.StatusInternalServerError, err)
		return false, err
	}

	if !ok {
		w.Header().Add("X-Djangolang-Cache-Status", "miss")
		return false, nil
	}

	w.Header().Add("X-Djangolang-Cache-Status", "hit")
	WriteResponse(w, http.StatusOK, cachedResponseAsJSON)

	return false, nil
}

func StoreCachedResponse(requestHash string, redisConn redis.Conn, responseAsJSON []byte) error {
	if redisConn == nil {
		return nil
	}

	_, err := redisConn.Do("SET", requestHash, string(responseAsJSON))
	if err != nil {
		return fmt.Errorf("failed to set value for cache item %v = %v; %v", requestHash, string(responseAsJSON), err)
	}

	_, err = redisConn.Do("EXPIRE", requestHash, "86400")
	if err != nil {
		return fmt.Errorf("failed to set expiry for cache item %v = %v failed; %v", requestHash, string(responseAsJSON), err)
	}

	return nil
}
