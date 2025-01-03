package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"runtime"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/initialed85/djangolang/pkg/config"
	"github.com/initialed85/djangolang/pkg/introspect"
	"github.com/initialed85/djangolang/pkg/stream"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/exp/maps"
)

var reloadChangeObjects = config.ReloadChangeObjects()

type PathValue struct {
	VisitedTableNames []string
}

type Waiter struct {
	Action    stream.Action
	TableName string
	Xid       uint32
}

const handshakeTimeout = time.Second * 10

func GetDefaultHTTPMiddlewares(extraHTTPMiddlewares ...HTTPMiddleware) []HTTPMiddleware {
	middleware.DefaultLogger = middleware.RequestLogger(
		&middleware.DefaultLogFormatter{
			Logger: log,
		},
	)

	httpMiddlewares := make([]HTTPMiddleware, 0)

	httpMiddlewares = append(httpMiddlewares, middleware.RequestID)
	httpMiddlewares = append(httpMiddlewares, middleware.RealIP)
	httpMiddlewares = append(httpMiddlewares, middleware.StripSlashes)
	httpMiddlewares = append(httpMiddlewares, middleware.Logger)
	httpMiddlewares = append(httpMiddlewares, middleware.Recoverer)

	httpMiddlewares = append(httpMiddlewares, extraHTTPMiddlewares...)

	return httpMiddlewares
}

func RunServer(
	outerCtx context.Context,
	outerChanges chan Change,
	addr string,
	newFromItem func(string, map[string]any) (any, error),
	mutateRouterFn MutateRouterFn,
	db *pgxpool.Pool,
	redisPool *redis.Pool,
	httpMiddlewares []HTTPMiddleware,
	objectMiddlewares []ObjectMiddleware,
	addCustomHandlers func(chi.Router) error,
	tableByName introspect.TableByName,
	nodeNames ...string,
) error {
	ctx, cancel := context.WithCancel(outerCtx)
	defer cancel()

	if len(httpMiddlewares) == 0 {
		httpMiddlewares = GetDefaultHTTPMiddlewares()
	}

	changesByWaiterMu := new(sync.Mutex)
	changesByWaiter := make(map[Waiter]chan Change)

	var waitForChange WaitForChange = func(ctx context.Context, actions []stream.Action, tableName string, xid uint32) (*Change, error) {
		waiters := make([]Waiter, 0)
		for _, action := range actions {
			waiter := Waiter{
				Action:    action,
				TableName: tableName,
				Xid:       xid,
			}

			waiters = append(waiters, waiter)
		}

		changes := make(chan Change, 1)
		changesByWaiterMu.Lock()
		for _, waiter := range waiters {
			changesByWaiter[waiter] = changes
		}
		changesByWaiterMu.Unlock()

		defer func() {
			changesByWaiterMu.Lock()
			for _, waiter := range waiters {
				delete(changesByWaiter, waiter)
			}
			changesByWaiterMu.Unlock()
		}()

		select {
		case <-ctx.Done():
			break
		case change := <-changes:
			return &change, nil
		}

		return nil, fmt.Errorf("context canceled while waiting for change")
	}

	actualRouter := chi.NewRouter()

	for _, m := range httpMiddlewares {
		actualRouter.Use(m)
	}

	if mutateRouterFn != nil {
		mutateRouterFn(actualRouter, db, redisPool, objectMiddlewares, waitForChange)
	}

	changes := make(chan stream.Change, 1024)

	mu := new(sync.Mutex)
	outgoingMessagesBySubscriberIDByTableName := make(map[string]map[uuid.UUID]chan []byte)

	// this goroutine drains the changes from the change stream and publishes them out to the WebSocket clients
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case change := <-changes:
				object, err := newFromItem(change.TableName, change.Item)
				if err != nil {
					log.Printf("warning: failed to convert item to object for %s; %v", change.String(), err)
					return
				}

				// optional, disabled by default for performance and consistency
				if reloadChangeObjects {
					// for INSERT / UPDATE / SOFT_DELETE the row should still be there (so we can likely do a reload to get the
					// nested objects etc)
					if change.Action != stream.DELETE && change.Action != stream.TRUNCATE {
						func() {
							logErr := func(err error) {
								log.Printf("warning: failed to reload object for %s (will send out as-is); %v", change.String(), err)
							}

							tx, err := db.Begin(ctx)
							if err != nil {
								logErr(err)
								return
							}

							defer func() {
								_ = tx.Rollback(ctx)
							}()

							possibleObject, ok := object.(WithReload)
							if !ok {
								logErr(err)
								return
							}

							err = possibleObject.Reload(ctx, tx, true)
							if err != nil {
								logErr(err)
								return
							}

							object = possibleObject
						}()
					}
				}

				objectChange := Change{
					Timestamp: change.Timestamp,
					ID:        change.ID,
					Action:    change.Action,
					TableName: change.TableName,
					Item:      change.Item,
					Object:    object,
					Xid:       change.Xid,
				}

				waiter := Waiter{
					Action:    change.Action,
					TableName: change.TableName,
					Xid:       change.Xid,
				}

				changesByWaiterMu.Lock()
				changesForWaiter := changesByWaiter[waiter]
				changesByWaiterMu.Unlock()

				if changesForWaiter != nil {
					select {
					case changesForWaiter <- objectChange:
					default:
						log.Printf(
							"warning: attempt to write %v to %#+v would block (broken waiter / duplicate change); will cancel waiter...",
							waiter, objectChange.String(),
						)

						changesByWaiterMu.Lock()
						close(changesForWaiter)
						delete(changesByWaiter, waiter)
						changesByWaiterMu.Unlock()
					}
				}

				// TODO: unbounded goroutine use could blow out if we're dealing with a lot of changes, should
				//   probably be a limited number of workers or something like that
				// this provides a way to clone the change stream to an externally injected channel; note that
				// we don't wait around if the channel is blocked (it's on the reader to keep that channel happy)
				go func() {
					if outerChanges != nil {
						select {
						case outerChanges <- objectChange:
						default:
							log.Printf("warning: failed push %#+v to outerChanges", objectChange)
						}
					}
				}()
				runtime.Gosched()

				// TODO: unbounded goroutine use could blow out if we're dealing with a lot of changes, should
				//   probably be a limited number of workers or something like that
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
						log.Printf("warning: failed to marshal %#+v to JSON; %v", change, err)
						return
					}

					for _, outgoingMessages := range allOutgoingMessages {
						select {
						case outgoingMessages <- b:
						default:
						}
					}
				}()
				runtime.Gosched()
			}
		}
	}()
	runtime.Gosched()

	var err error

	// this goroutine runs the handler for the CDC stream
	go func() {
		defer cancel()

		// this is a blocking call
		err = stream.Run(ctx, changes, tableByName, nodeNames...)
		if err != nil {
			log.Printf("stream.Run failed: %v", err)
			return
		}
	}()
	runtime.Gosched()

	upgrader := websocket.Upgrader{
		HandshakeTimeout: handshakeTimeout,
		Subprotocols:     []string{"djangolang"},
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		EnableCompression: true,
	}

	actualRouter.Get("/__stream", func(w http.ResponseWriter, r *http.Request) {
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
			HandleErrorResponse(
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
			HandleErrorResponse(
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
			HandleErrorResponse(
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
		runtime.Gosched()

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

		// this goroutine cleans up any subscriptions on disconnect of the WebSocket client
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
		runtime.Gosched()

		// this goroutine handles reads from the WebSocket client
		go func() {
			defer connCancel()

			for {
				select {
				case <-ctx.Done():
					return
				default:
				}

				// we don't actually expect the WebSocket client to send any messages just yet
				_, _, err := conn.ReadMessage()
				if err != nil {
					_, _, outgoingMessage, _ := GetResponse(http.StatusBadRequest, fmt.Errorf("read failed: %v", err), nil)
					_ = conn.WriteControl(websocket.CloseAbnormalClosure, outgoingMessage, time.Now().Add(time.Second*1))
					return
				}
			}
		}()
		runtime.Gosched()

		// this goroutine publishes changes for the subscriptions applicable to the WebSocket client
		go func() {
			defer connCancel()

			for {
				select {
				case <-ctx.Done():
					return
				case b := <-outgoingMessages:
					err := conn.WriteMessage(websocket.BinaryMessage, b)
					if err != nil {
						_, _, outgoingMessage, _ := GetResponse(http.StatusBadRequest, fmt.Errorf("write failed: %v", err), nil)
						_ = conn.WriteControl(websocket.CloseAbnormalClosure, outgoingMessage, time.Now().Add(time.Second*1))
						return
					}
				}
			}
		}()
		runtime.Gosched()
	})

	if addCustomHandlers != nil {
		actualRouter.Route("/custom", func(r chi.Router) {
			err = addCustomHandlers(r)
			if err != nil {
				err = fmt.Errorf("failed to add custom handlers; %v", err)
				return
			}
		})
	}

	if err != nil {
		return err
	}

	apiRoot := config.APIRoot()
	finalRouter := chi.NewRouter()
	finalRouter.Mount(fmt.Sprintf("/%s", strings.Trim(apiRoot, "/")), actualRouter)

	var gatherRoutes func(string, []chi.Route)

	patterns := make(map[string]struct{})

	gatherRoutes = func(pattern string, routes []chi.Route) {
		for _, route := range routes {
			if route.SubRoutes != nil {
				gatherRoutes(
					strings.TrimRight(pattern+route.Pattern, "/*"),
					route.SubRoutes.Routes(),
				)
			}

			patterns[strings.TrimRight(pattern+route.Pattern, "/*")] = struct{}{}
		}
	}

	gatherRoutes("", finalRouter.Routes())

	routes := maps.Keys(patterns)
	slices.Sort(routes)
	slices.Reverse(routes)

	for _, route := range routes {
		log.Printf("registered: %s", route)
	}

	httpServer := &http.Server{
		Addr:    addr,
		Handler: finalRouter,
		BaseContext: func(l net.Listener) context.Context {
			return ctx
		},
	}

	// this goroutine cleans up the HTTP server on shutdown
	go func() {
		<-ctx.Done()

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), time.Second*1)
		defer shutdownCancel()

		_ = httpServer.Shutdown(shutdownCtx)
		_ = httpServer.Close()
	}()
	runtime.Gosched()

	// this is a blocking call
	err = httpServer.ListenAndServe()
	if err != nil {
		log.Printf("httpServer.ListenAndServe failed: %v", err)
		return err
	}

	return nil
}
