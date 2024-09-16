package query

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/google/uuid"
	"github.com/initialed85/djangolang/pkg/helpers"
)

type CachedItems struct {
	Items      []map[string]any
	Count      int64
	TotalCount int64
	Limit      int64
	Offset     int64
}

var mu = new(sync.Mutex)
var cachedItemsByCacheKey = make(map[string]*CachedItems)
var cacheKeysByQueryID = make(map[uuid.UUID]*[]string)

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

func getCachedItems(cacheKey string) (*CachedItems, bool) {
	mu.Lock()
	defer mu.Unlock()

	cachedItems, ok := cachedItemsByCacheKey[cacheKey]
	if !ok {
		return nil, false
	}

	return cachedItems, true
}

func setCachedItems(queryID *uuid.UUID, cacheKey string, items *CachedItems) {
	if queryID == nil {
		return
	}

	mu.Lock()
	defer mu.Unlock()

	cacheKeys := cacheKeysByQueryID[*queryID]
	if cacheKeys == nil {
		cacheKeys = helpers.Ptr(make([]string, 0))
	}

	cachedItemsByCacheKey[cacheKey] = items
	*cacheKeys = append(*cacheKeys, cacheKey)
	cacheKeysByQueryID[*queryID] = cacheKeys
}

func clearCachedItems(queryID *uuid.UUID) {
	if queryID == nil {
		return
	}

	mu.Lock()
	defer mu.Unlock()

	possibleCacheKeys := cacheKeysByQueryID[*queryID]
	if possibleCacheKeys == nil {
		return
	}

	for _, cacheKey := range *possibleCacheKeys {
		delete(cachedItemsByCacheKey, cacheKey)
	}

	delete(cacheKeysByQueryID, *queryID)
}

type queryIDContextKey struct{}

var QueryIDContextKey = queryIDContextKey{}

func GetQueryID(ctx context.Context) *uuid.UUID {
	rawQueryID := ctx.Value(QueryIDContextKey)
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
		queryID := uuid.Must(uuid.NewRandom())
		ctx = context.WithValue(ctx, QueryIDContextKey, queryID)
		cleanup = func() {
			clearCachedItems(&queryID)
		}
	}

	return ctx, cleanup
}
