package query

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/google/uuid"
)

var mu = new(sync.Mutex)
var cachedItemsByCacheKey = make(map[string][]map[string]any)
var cacheKeysByQueryID = make(map[uuid.UUID][]string)

func getCacheKey(
	queryID *uuid.UUID,
	columns []string,
	table string,
	where string,
	orderBy *string,
	limit *int,
	offset *int,
	values ...any,
) (string, error) {
	if queryID == nil {
		return "", nil
	}

	b, err := json.Marshal(
		map[string]any{
			"queryID": queryID,
			"columns": columns,
			"table":   table,
			"where":   where,
			"orderBy": orderBy,
			"limit":   limit,
			"offset":  offset,
			"values":  values,
		},
	)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func getCachedItems(cacheKey string) ([]map[string]any, bool) {
	mu.Lock()
	defer mu.Unlock()

	cachedItems, ok := cachedItemsByCacheKey[cacheKey]
	if !ok {
		return nil, false
	}

	return cachedItems, true
}

func setCachedItems(queryID *uuid.UUID, cacheKey string, items []map[string]any) {
	if queryID == nil {
		return
	}

	mu.Lock()
	defer mu.Unlock()

	cacheKeys := cacheKeysByQueryID[*queryID]
	if cacheKeys == nil {
		cacheKeys = make([]string, 0)
	}

	cachedItemsByCacheKey[cacheKey] = items
	cacheKeys = append(cacheKeys, cacheKey)
	cacheKeysByQueryID[*queryID] = cacheKeys
}

func clearCachedItems(queryID *uuid.UUID) {
	if queryID == nil {
		return
	}

	mu.Lock()
	defer mu.Unlock()

	for _, cacheKey := range cacheKeysByQueryID[*queryID] {
		delete(cachedItemsByCacheKey, cacheKey)
	}

	delete(cacheKeysByQueryID, *queryID)
}

var queryIdContextKey = struct{}{}

func GetQueryID(ctx context.Context) *uuid.UUID {
	rawQueryID := ctx.Value(queryIdContextKey)
	queryID, ok := rawQueryID.(uuid.UUID)
	if !ok {
		return nil
	}

	return &queryID
}

var dummyCleanup = func() {}

func WithQueryID(ctx context.Context) (context.Context, func()) {
	cleanup := dummyCleanup

	queryID := GetQueryID(ctx)
	if queryID == nil {
		ctx = context.WithValue(ctx, queryIdContextKey, uuid.Must(uuid.NewRandom()))
		cleanup = func() {
			clearCachedItems(queryID)
		}
	}

	return ctx, cleanup
}
