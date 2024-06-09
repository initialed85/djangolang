package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/initialed85/djangolang/pkg/helpers"
	"github.com/initialed85/djangolang/pkg/introspect"
	"github.com/initialed85/djangolang/pkg/stream"
	"github.com/jmoiron/sqlx"
	"golang.org/x/exp/maps"
)

const handshakeTimeout = time.Second * 10

type WithReload interface {
	Reload(context.Context, *sqlx.Tx) error
}

type WithPrimaryKey interface {
	GetPrimaryKeyColumn() string
	GetPrimaryKeyValue() any
}

type Change struct {
	ID        uuid.UUID      `json:"id"`
	Action    stream.Action  `json:"action"`
	TableName string         `json:"table_name"`
	Item      map[string]any `json:"-"`
	Object    any            `json:"object"`
}

func (c *Change) String() string {
	b, _ := json.Marshal(c.Object)

	primaryKeySummary := ""
	if c.Object != nil {
		object, ok := c.Object.(WithPrimaryKey)
		if ok {
			primaryKeySummary = fmt.Sprintf(
				"(%s = %s) ",
				object.GetPrimaryKeyColumn(),
				object.GetPrimaryKeyValue(),
			)
		}
	}

	return fmt.Sprintf(
		"(%s) %s %s %s%s",
		c.ID, c.Action, c.TableName, primaryKeySummary, string(b),
	)
}

func RunServer(
	ctx context.Context,
	outerChanges chan Change,
	addr string,
	newFromItem func(string, map[string]any) (any, error),
	getRouterFn func(*sqlx.DB, ...func(http.Handler) http.Handler) chi.Router,
	db *sqlx.DB,
	middlewares ...func(http.Handler) http.Handler,
) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	if len(middlewares) == 0 {
		middlewares = append(middlewares, middleware.Recoverer)
		middlewares = append(middlewares, middleware.RequestID)
		middlewares = append(middlewares, middleware.Logger)
		middlewares = append(middlewares, middleware.RealIP)
	}

	r := getRouterFn(db, middlewares...)

	schema := helpers.GetSchema()

	tableByName, err := introspect.Introspect(ctx, db, schema)
	if err != nil {
		return err
	}

	changes := make(chan stream.Change, 1024)

	mu := new(sync.Mutex)
	outgoingMessagesBySubscriberIDByTableName := make(map[string]map[uuid.UUID]chan []byte)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case change := <-changes:
				func() {
					object, err := newFromItem(change.TableName, change.Item)
					if err != nil {
						log.Printf("warning: failed to convert item to object for %s: %v", change.String(), err)
						return
					}

					// if change.Action != stream.DELETE && change.Action != stream.TRUNCATE {
					// 	func() {
					// 		logErr := func(err error) {
					// 			log.Printf("warning: failed to reload object for %s (will send out as-is): %v", change.String(), err)
					// 		}

					// 		tx, err := db.Beginx()
					// 		if err != nil {
					// 			logErr(err)
					// 			return
					// 		}

					// 		defer func() {
					// 			_ = tx.Rollback()
					// 		}()

					// 		possibleObject, ok := object.(WithReload)
					// 		if !ok {
					// 			logErr(err)
					// 			return
					// 		}

					// 		err = possibleObject.Reload(ctx, tx)
					// 		if err != nil {
					// 			logErr(err)
					// 			return
					// 		}

					// 		object = possibleObject
					// 	}()
					// }

					objectChange := Change{
						ID:        change.ID,
						Action:    change.Action,
						TableName: change.TableName,
						Item:      change.Item,
						Object:    object,
					}

					if outerChanges != nil {
						select {
						case outerChanges <- objectChange:
						default:
						}
					}

					go func() {
						var allOutgoingMessages []chan []byte

						mu.Lock()
						outgoingMessagesBySubscriberID := outgoingMessagesBySubscriberIDByTableName[change.TableName]
						if outgoingMessagesBySubscriberID != nil {
							allOutgoingMessages = maps.Values(outgoingMessagesBySubscriberID)
						}
						mu.Unlock()

						if allOutgoingMessages == nil {
							return
						}

						b, err := json.Marshal(objectChange)
						if err != nil {
							log.Printf("warning: failed to marshal %#+v to JSON: %v", change, err)
							return
						}

						for _, outgoingMessages := range allOutgoingMessages {
							select {
							case outgoingMessages <- b:
							default:
							}
						}
					}()
				}()
			}
		}
	}()

	go func() {
		defer cancel()

		err = stream.Run(ctx, changes, tableByName)
		if err != nil {
			log.Printf("stream.Run failed: %v", err)
			return
		}
	}()

	upgrader := websocket.Upgrader{
		HandshakeTimeout: handshakeTimeout,
		Subprotocols:     []string{"djangolang"},
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		EnableCompression: true,
	}

	r.Get("/__stream", func(w http.ResponseWriter, r *http.Request) {
		unrecognizedParams := make([]string, 0)
		for k, vs := range r.URL.Query() {
			if k == "include" || k == "exclude" {
				continue
			}

			for _, v := range vs {
				unrecognizedParams = append(unrecognizedParams, fmt.Sprintf("%s=%s", k, v))
			}
		}

		if len(unrecognizedParams) > 0 {
			helpers.HandleErrorResponse(
				w,
				http.StatusInternalServerError,
				fmt.Errorf("unrecognized params %s", strings.Join(unrecognizedParams, ", ")),
			)
			return
		}

		unknownTableNames := make([]string, 0)

		includeTableNames := make([]string, 0)
		for _, tableName := range r.URL.Query()["include"] {
			_, ok := tableByName[tableName]
			if !ok {
				unknownTableNames = append(unknownTableNames, fmt.Sprintf("%s=%s", "include", tableName))
				continue
			}

			includeTableNames = append(includeTableNames, tableName)
		}

		excludeTableNames := make([]string, 0)
		for _, tableName := range r.URL.Query()["exclude"] {
			_, ok := tableByName[tableName]
			if !ok {
				unknownTableNames = append(unknownTableNames, fmt.Sprintf("%s=%s", "exclude", tableName))
				continue
			}

			excludeTableNames = append(excludeTableNames, tableName)
		}

		if len(unknownTableNames) > 0 {
			helpers.HandleErrorResponse(
				w,
				http.StatusInternalServerError,
				fmt.Errorf("unknown table names %s", strings.Join(unknownTableNames, ", ")),
			)
			return
		}

		tableNames := make([]string, 0)

		if len(includeTableNames) == 0 {
			includeTableNames = maps.Keys(tableByName)
		}

		for _, includeTableName := range includeTableNames {
			include := true

			for _, excludeTableName := range excludeTableNames {
				if includeTableName == excludeTableName {
					include = false
					break
				}
			}

			if !include {
				continue
			}

			tableNames = append(tableNames, includeTableName)
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			helpers.HandleErrorResponse(
				w,
				http.StatusInternalServerError,
				fmt.Errorf("failed to upgrade WebSocket: %s", err),
			)
			return
		}

		connCtx, connCancel := context.WithCancel(ctx)
		go func() {
			defer connCancel()
			<-connCtx.Done()
			_ = conn.Close()
		}()

		conn.SetCloseHandler(func(code int, text string) error {
			connCancel()
			return nil
		})

		subscriberID := uuid.New()
		outgoingMessages := make(chan []byte, 128)

		mu.Lock()
		for _, tableName := range tableNames {
			outgoingMessagesBySubscriberID := outgoingMessagesBySubscriberIDByTableName[tableName]
			if outgoingMessagesBySubscriberID == nil {
				outgoingMessagesBySubscriberID = make(map[uuid.UUID]chan []byte)
			}

			outgoingMessagesBySubscriberID[subscriberID] = outgoingMessages
			outgoingMessagesBySubscriberIDByTableName[tableName] = outgoingMessagesBySubscriberID
		}
		mu.Unlock()

		go func() {
			<-connCtx.Done()
			mu.Lock()
			for _, tableName := range tableNames {
				outgoingMessagesBySubscriberID := outgoingMessagesBySubscriberIDByTableName[tableName]
				if outgoingMessagesBySubscriberID == nil {
					continue
				}

				delete(outgoingMessagesBySubscriberID, subscriberID)
				outgoingMessagesBySubscriberIDByTableName[tableName] = outgoingMessagesBySubscriberID
			}
			mu.Unlock()
		}()

		go func() {
			defer connCancel()

			for {
				select {
				case <-ctx.Done():
					return
				default:
				}

				_, _, err := conn.ReadMessage()
				if err != nil {
					_, _, outgoingMessage, _ := helpers.GetResponse(http.StatusBadRequest, fmt.Errorf("read failed: %v", err), nil)
					_ = conn.WriteControl(websocket.CloseAbnormalClosure, outgoingMessage, time.Now().Add(time.Second*1))
					return
				}
			}
		}()

		go func() {
			defer connCancel()

			for {
				select {
				case <-ctx.Done():
					return
				case b := <-outgoingMessages:
					err = conn.WriteMessage(websocket.BinaryMessage, b)
					if err != nil {
						_, _, outgoingMessage, _ := helpers.GetResponse(http.StatusBadRequest, fmt.Errorf("write failed: %v", err), nil)
						_ = conn.WriteControl(websocket.CloseAbnormalClosure, outgoingMessage, time.Now().Add(time.Second*1))
						return
					}
				}
			}
		}()
	})

	httpServer := &http.Server{
		Addr:    addr,
		Handler: r,
		BaseContext: func(l net.Listener) context.Context {
			return ctx
		},
	}

	go func() {
		<-ctx.Done()

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), time.Second*1)
		defer shutdownCancel()

		_ = httpServer.Shutdown(shutdownCtx)
		_ = httpServer.Close()
	}()

	err = httpServer.ListenAndServe()
	if err != nil {
		log.Printf("httpServer.ListenAndServe failed: %v", err)
		return err
	}

	return nil
}
