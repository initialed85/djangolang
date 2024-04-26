package test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/initialed85/djangolang/pkg/helpers"
	"github.com/initialed85/djangolang/pkg/some_db"
	"github.com/initialed85/djangolang/pkg/types"
	"github.com/stretchr/testify/require"
)

func TestEndpoints(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := helpers.GetDBFromEnvironment(ctx)
	if err != nil {
		require.NoError(t, err)
	}
	defer func() {
		_ = db.Close()
	}()

	port, err := helpers.GetPort()
	require.NoError(t, err)

	changes := make(chan types.Change, 1024)
	go func() {
		err = some_db.RunServer(ctx, changes)
		require.NoError(t, err)
	}()
	runtime.Gosched()
	time.Sleep(time.Second * 1)

	mu := new(sync.Mutex)
	lastChangeByTableName := make(map[string]types.Change)
	lastWebSocketChangeByTableName := make(map[string]types.Change)

	dialer := websocket.Dialer{
		HandshakeTimeout:  time.Second * 10,
		Subprotocols:      []string{"djangolang"},
		EnableCompression: true,
	}

	conn, _, err := dialer.Dial(
		fmt.Sprintf("ws://localhost:%v/__subscribe", port),
		nil,
	)
	require.NoError(t, err)
	defer func() {
		conn.Close()
	}()

	rootWg := new(sync.WaitGroup)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			_, b, err := conn.ReadMessage()
			if err != nil {
				log.Printf("warning: conn.readMessage() returned err: %v", err)
				return
			}

			var change types.Change
			err = json.Unmarshal(b, &change)
			require.NoError(t, err)

			mu.Lock()
			lastWebSocketChangeByTableName[change.TableName] = change
			mu.Unlock()
		}
	}()

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

	rootWg.Add(1)
	t.Run("TestPostGetPutDelete", func(t *testing.T) {
		defer rootWg.Done()

		httpClient := http.Client{
			Timeout: time.Second * 5,
		}

		for i := 0; i < 50; i++ {
			physicalThings := make([]*some_db.PhysicalThing, 0)

			physicalThing := &some_db.PhysicalThing{
				Name:     fmt.Sprintf("SomePhysicalThing-%v", i),
				Type:     "http://some-url.org",
				Tags:     make([]string, 0),
				Metadata: make(map[string]sql.NullString),
			}
			physicalThingJSON, err := json.Marshal(physicalThing)
			require.NoError(t, err)

			otherPhysicalThing := &some_db.PhysicalThing{
				Name:     fmt.Sprintf("OtherPhysicalThing-%v", i),
				Type:     "http://some-url.org",
				Tags:     make([]string, 0),
				Metadata: make(map[string]sql.NullString),
			}
			otherPhysicalThingJSON, err := json.Marshal(otherPhysicalThing)
			require.NoError(t, err)

			resp, err := httpClient.Get(fmt.Sprintf("http://localhost:%v/physical-things", port))
			require.NoError(t, err)
			b, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			_ = resp.Body.Close()
			require.Equal(t, http.StatusOK, resp.StatusCode, string(b))
			err = json.Unmarshal(b, &physicalThings)
			require.NoError(t, err)
			require.Len(t, physicalThings, 0)

			resp, err = httpClient.Post(
				fmt.Sprintf("http://localhost:%v/physical-things", port),
				"application/json",
				bytes.NewBuffer(physicalThingJSON),
			)
			require.NoError(t, err)
			b, err = io.ReadAll(resp.Body)
			require.NoError(t, err)
			_ = resp.Body.Close()
			require.Equal(t, http.StatusCreated, resp.StatusCode, string(b))
			err = json.Unmarshal(b, &physicalThing)
			require.NoError(t, err)
			require.NotZero(t, physicalThing.ID)
			require.Equal(t, fmt.Sprintf("SomePhysicalThing-%v", i), physicalThing.Name)
			require.Equal(t, "http://some-url.org", physicalThing.Type)
			require.Eventually(t, func() bool {
				mu.Lock()
				change, ok := lastChangeByTableName[some_db.PhysicalThingTable]
				mu.Unlock()

				if !ok {
					return false
				}

				if change.Action != types.INSERT {
					return false
				}

				object := change.Object.(*some_db.PhysicalThing)

				return object.ID == physicalThing.ID
			}, time.Second*1, time.Millisecond*1)
			require.Eventually(t, func() bool {
				mu.Lock()
				change, ok := lastWebSocketChangeByTableName[some_db.PhysicalThingTable]
				mu.Unlock()

				if !ok {
					return false
				}

				if change.Action != types.INSERT {
					return false
				}

				object := change.Object.(map[string]any)

				return object["id"].(string) == physicalThing.ID.String()
			}, time.Second*1, time.Millisecond*1)

			resp, err = httpClient.Post(
				fmt.Sprintf("http://localhost:%v/physical-things", port),
				"application/json",
				bytes.NewBuffer(physicalThingJSON),
			)
			require.NoError(t, err)
			b, err = io.ReadAll(resp.Body)
			require.NoError(t, err)
			_ = resp.Body.Close()
			require.Equal(t, http.StatusConflict, resp.StatusCode, string(b))
			err = json.Unmarshal(b, &physicalThing)
			require.NoError(t, err)
			require.NotZero(t, physicalThing.ID)
			require.Equal(t, fmt.Sprintf("SomePhysicalThing-%v", i), physicalThing.Name)
			require.Equal(t, "http://some-url.org", physicalThing.Type)

			req, err := http.NewRequest(
				http.MethodPut,
				fmt.Sprintf("http://localhost:%v/physical-things/%v", port, physicalThing.ID.String()),
				bytes.NewBuffer(otherPhysicalThingJSON),
			)
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			resp, err = httpClient.Do(req)
			require.NoError(t, err)
			b, err = io.ReadAll(resp.Body)
			require.NoError(t, err)
			_ = resp.Body.Close()
			require.Equal(t, http.StatusOK, resp.StatusCode, string(b))
			err = json.Unmarshal(b, &physicalThing)
			require.NoError(t, err)
			require.Equal(t, fmt.Sprintf("OtherPhysicalThing-%v", i), physicalThing.Name)
			require.Eventually(t, func() bool {
				mu.Lock()
				change, ok := lastChangeByTableName[some_db.PhysicalThingTable]
				mu.Unlock()

				if !ok {
					return false
				}

				if change.Action != types.UPDATE {
					return false
				}

				object := change.Object.(*some_db.PhysicalThing)

				return object.ID == physicalThing.ID
			}, time.Second*1, time.Millisecond*1)
			require.Eventually(t, func() bool {
				mu.Lock()
				change, ok := lastWebSocketChangeByTableName[some_db.PhysicalThingTable]
				mu.Unlock()

				if !ok {
					return false
				}

				if change.Action != types.UPDATE {
					return false
				}

				object := change.Object.(map[string]any)

				return object["id"].(string) == physicalThing.ID.String()
			}, time.Second*1, time.Millisecond*1)

			resp, err = httpClient.Post(
				fmt.Sprintf("http://localhost:%v/physical-things", port),
				"application/json",
				bytes.NewBuffer(otherPhysicalThingJSON),
			)
			require.NoError(t, err)
			b, err = io.ReadAll(resp.Body)
			require.NoError(t, err)
			_ = resp.Body.Close()
			require.Equal(t, http.StatusConflict, resp.StatusCode, string(b))
			err = json.Unmarshal(b, &physicalThing)
			require.NoError(t, err)
			require.NotZero(t, physicalThing.ID)
			require.Equal(t, fmt.Sprintf("OtherPhysicalThing-%v", i), physicalThing.Name)
			require.Equal(t, "http://some-url.org", physicalThing.Type)

			req, err = http.NewRequest(
				http.MethodDelete,
				fmt.Sprintf("http://localhost:%v/physical-things/%v", port, physicalThing.ID.String()),
				nil,
			)
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			resp, err = httpClient.Do(req)
			require.NoError(t, err)
			b, err = io.ReadAll(resp.Body)
			require.NoError(t, err)
			_ = resp.Body.Close()
			require.Equal(t, http.StatusNoContent, resp.StatusCode, string(b))
			err = json.Unmarshal(b, &struct{}{})
			require.Error(t, err)
			require.Eventually(t, func() bool {
				mu.Lock()
				change, ok := lastChangeByTableName[some_db.PhysicalThingTable]
				mu.Unlock()

				if !ok {
					return false
				}

				if change.Action != types.DELETE {
					return false
				}

				return change.PrimaryKeyValue.(uuid.UUID) == physicalThing.ID
			}, time.Second*1, time.Millisecond*1)
			require.Eventually(t, func() bool {
				mu.Lock()
				change, ok := lastWebSocketChangeByTableName[some_db.PhysicalThingTable]
				mu.Unlock()

				if !ok {
					return false
				}

				if change.Action != types.DELETE {
					return false
				}

				return change.PrimaryKeyValue.(string) == physicalThing.ID.String()
			}, time.Second*1, time.Millisecond*1)

			resp, err = httpClient.Get(fmt.Sprintf("http://localhost:%v/physical-things", port))
			require.NoError(t, err)
			b, err = io.ReadAll(resp.Body)
			require.NoError(t, err)
			_ = resp.Body.Close()
			require.Equal(t, resp.StatusCode, http.StatusOK, string(b))
			err = json.Unmarshal(b, &physicalThings)
			require.NoError(t, err)
			require.Len(t, physicalThings, 0)

			req, err = http.NewRequest(
				http.MethodDelete,
				fmt.Sprintf("http://localhost:%v/physical-things/%v", port, physicalThing.ID.String()),
				bytes.NewBuffer(physicalThingJSON),
			)
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			resp, err = httpClient.Do(req)
			require.NoError(t, err)
			b, err = io.ReadAll(resp.Body)
			require.NoError(t, err)
			_ = resp.Body.Close()
			require.Equal(t, http.StatusNotFound, resp.StatusCode, string(b))
		}
	})

	rootWg.Add(1)
	t.Run("TestPerformance", func(t *testing.T) {
		defer rootWg.Done()

		wg := new(sync.WaitGroup)

		limiter := make(chan bool, 50)
		for i := 0; i < 50; i++ {
			limiter <- true
		}

		//
		// post
		//

		for i := 0; i < 1000; i++ {
			wg.Add(1)

			go func(i int) {
				defer wg.Done()

				<-limiter
				defer func() {
					limiter <- true
				}()

				httpClient := http.Client{
					Timeout: time.Second * 5,
				}

				physicalThing := &some_db.PhysicalThing{
					Name:     fmt.Sprintf("SomePhysicalThing-%v", i+1+1000),
					Type:     "http://some-url.org",
					Tags:     make([]string, 0),
					Metadata: make(map[string]sql.NullString),
				}

				b, err := json.Marshal(physicalThing)
				require.NoError(t, err)

				resp, err := httpClient.Post(
					fmt.Sprintf("http://localhost:%v/physical-things", port),
					"application/json",
					bytes.NewBuffer(b),
				)
				require.NoError(t, err)

				b, err = io.ReadAll(resp.Body)
				require.NoError(t, err)
				defer func() {
					_ = resp.Body.Close()
				}()
				require.Equal(t, http.StatusCreated, resp.StatusCode, string(b))
			}(i)
		}

		wg.Wait()

		//
		// get
		//

		httpClient := http.Client{
			Timeout: time.Second * 5,
		}

		resp, err := httpClient.Get(fmt.Sprintf("http://localhost:%v/physical-things?limit=2000&order_by=id&order=asc", port))
		require.NoError(t, err)

		b, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		defer func() {
			_ = resp.Body.Close()
		}()
		require.Equal(t, http.StatusOK, resp.StatusCode, string(b))

		physicalThings := make([]*some_db.PhysicalThing, 0)
		err = json.Unmarshal(b, &physicalThings)
		require.NoError(t, err)
		require.Equal(t, 1000, len(physicalThings))

		//
		// patch
		//

		for i, physicalThing := range physicalThings {
			wg.Add(1)

			go func(i int, physicalThing *some_db.PhysicalThing) {
				defer wg.Done()

				<-limiter
				defer func() {
					limiter <- true
				}()

				httpClient := http.Client{
					Timeout: time.Second * 5,
				}

				physicalThing.Name = fmt.Sprintf("OtherPhysicalThing-%v", i+1)

				b, err := json.Marshal(physicalThing)
				require.NoError(t, err)

				req, err := http.NewRequest(
					http.MethodPut,
					fmt.Sprintf("http://localhost:%v/physical-things/%v", port, physicalThing.ID.String()),
					bytes.NewBuffer(b),
				)
				require.NoError(t, err)
				req.Header.Set("Content-Type", "application/json")

				resp, err := httpClient.Do(req)
				require.NoError(t, err)

				b, err = io.ReadAll(resp.Body)
				require.NoError(t, err)
				defer func() {
					_ = resp.Body.Close()
				}()
				require.Equal(t, http.StatusOK, resp.StatusCode, string(b))

				err = json.Unmarshal(b, &physicalThing)
				require.NoError(t, err)
				require.Equal(t, fmt.Sprintf("OtherPhysicalThing-%v", i+1), physicalThing.Name)
			}(i, physicalThing)
		}

		wg.Wait()

		//
		// delete
		//

		for _, physicalThing := range physicalThings {
			wg.Add(1)

			go func(physicalThing *some_db.PhysicalThing) {
				defer wg.Done()

				<-limiter
				defer func() {
					limiter <- true
				}()

				httpClient := http.Client{
					Timeout: time.Second * 5,
				}

				req, err := http.NewRequest(
					http.MethodDelete,
					fmt.Sprintf("http://localhost:%v/physical-things/%v", port, physicalThing.ID.String()),
					nil,
				)
				require.NoError(t, err)

				resp, err = httpClient.Do(req)
				require.NoError(t, err)

				require.Equal(t, http.StatusNoContent, resp.StatusCode, string(b))
			}(physicalThing)
		}

		wg.Wait()
	})

	rootWg.Wait()
}
