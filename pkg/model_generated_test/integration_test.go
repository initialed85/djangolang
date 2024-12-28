package model_generated_test

import (
	"context"
	"log"
	"net/http"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/initialed85/djangolang/pkg/config"
	"github.com/initialed85/djangolang/pkg/model_generated"
	"github.com/initialed85/djangolang/pkg/model_generated_from_schema"
	"github.com/initialed85/djangolang/pkg/server"
	"github.com/stretchr/testify/require"
)

func TestIntegration(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := config.GetDBFromEnvironment(ctx)
	if err != nil {
		require.NoError(t, err)
	}
	defer func() {
		db.Close()
	}()

	redisPool, err := config.GetRedisFromEnvironment()
	if err != nil {
		require.NoError(t, err)
	}
	defer func() {
		redisPool.Close()
	}()

	redisConn := redisPool.Get()
	defer func() {
		_ = redisConn.Close()
	}()

	httpClient := &HTTPClient{
		httpClient: &http.Client{
			Timeout: time.Second * 10,
		},
	}

	changes := make(chan server.Change, 1024)
	mu := new(sync.Mutex)
	lastChangeByTableName := make(map[string]server.Change)

	addCustomHandlers := func(router chi.Router) error {
		collectPrimaryKeysHandler, err := server.GetHTTPHandler(
			http.MethodGet,
			"/collect-mr-primaries",
			http.StatusOK,
			func(
				ctx context.Context,
				pathParams server.EmptyPathParams,
				queryParams server.EmptyQueryParams,
				req server.EmptyRequest,
				rawReq any,
			) (*CollectMrPrimariesResponse, error) {
				tx, err := db.Begin(ctx)
				if err != nil {
					return nil, err
				}

				defer func() {
					_ = tx.Rollback(ctx)
				}()

				collectMrPrimariesResponse := CollectMrPrimariesResponse{
					MrPrimaries: []int{},
				}

				rows, err := db.Query(ctx, "SELECT array_agg(mr_primary) FROM not_null_fuzz;")
				if err != nil {
					return nil, err
				}

				for rows.Next() {
					err = rows.Scan(&collectMrPrimariesResponse.MrPrimaries)
					if err != nil {
						return nil, err
					}
				}

				err = rows.Err()
				if err != nil {
					return nil, err
				}

				err = tx.Commit(ctx)
				if err != nil {
					return nil, err
				}

				return &collectMrPrimariesResponse, nil
			},
		)
		if err != nil {
			return err
		}

		router.Get(collectPrimaryKeysHandler.FullPath, collectPrimaryKeysHandler.ServeHTTP)

		return nil
	}

	go func() {
		err := model_generated_from_schema.RunServer(ctx, nil, "127.0.0.1:4040", db, redisPool, nil, nil, nil, "model_generated_from_schema_test")
		if err != nil {
			log.Printf("stream.Run failed: %v", err)
		}
	}()
	runtime.Gosched()
	time.Sleep(time.Millisecond * 100)

	go func() {
		err := model_generated.RunServer(ctx, changes, "127.0.0.1:5050", db, redisPool, nil, nil, addCustomHandlers, "model_generated_test")
		if err != nil {
			log.Printf("stream.Run failed: %v", err)
		}
	}()
	runtime.Gosched()
	time.Sleep(time.Millisecond * 100)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case change := <-changes:
				mu.Lock()
				lastChangeByTableName[change.TableName] = change
				mu.Unlock()
			}
		}
	}()
	runtime.Gosched()
	time.Sleep(time.Millisecond * 100)

	getLastChangeForTableName := func(tableName string) *server.Change {
		mu.Lock()
		change, ok := lastChangeByTableName[tableName]
		mu.Unlock()
		if !ok {
			return nil
		}

		return &change
	}

	require.Eventually(
		t,
		func() bool {
			resp, err := httpClient.Get("http://localhost:5050/logical-things")
			if err != nil {
				return false
			}

			if resp.StatusCode != http.StatusOK {
				return false
			}

			return true
		},
		time.Second*10,
		time.Millisecond*100,
	)

	testIntegration(t, ctx, db, redisConn, mu, lastChangeByTableName, httpClient, getLastChangeForTableName)
	testLocationHistory(t, ctx, db, redisConn, mu, lastChangeByTableName, httpClient, getLastChangeForTableName)
	testLogicalThings(t, ctx, db, redisConn, mu, lastChangeByTableName, httpClient, getLastChangeForTableName)
	testNotNullFuzz(t, ctx, db, redisConn, mu, lastChangeByTableName, httpClient, getLastChangeForTableName)
	testIntegrationOther(t, ctx, db, redisConn, mu, lastChangeByTableName, httpClient, getLastChangeForTableName)
}
