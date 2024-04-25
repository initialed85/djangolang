package some_db

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/initialed85/djangolang/pkg/types"

	"github.com/chanced/caps"
	_pluralize "github.com/gertd/go-pluralize"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/initialed85/djangolang/pkg/helpers"
	"github.com/jmoiron/sqlx"
	"golang.org/x/exp/maps"
)

const pingInterval = time.Second * 60

type Subscriber struct {
	ID         uuid.UUID
	TableNames map[string]struct{}
	Messages   chan []byte
}

var (
	pluralize = _pluralize.NewClient()
)

func init() {
	converter, ok := caps.DefaultConverter.(caps.StdConverter)
	if !ok {
		logger.Panicf("failed to cast %#+v to caps.StdConverter", caps.DefaultConverter)
	}

	converter.Set("Dob", "DOB")
	converter.Set("Cpo", "CPO")
	converter.Set("Mwh", "MWH")
	converter.Set("Kwh", "KWH")
	converter.Set("Wh", "WH")
	converter.Set("Json", "JSON")
	converter.Set("Jsonb", "JSONB")
	converter.Set("Mac", "MAC")
	converter.Set("Ip", "IP")
}

func GetSelectHandlerForTableName(tableName string, db *sqlx.DB) (types.SelectHandler, error) {
	selectFunc, ok := selectFuncByTableName[tableName]
	if !ok {
		return nil, fmt.Errorf("no selectFuncByTableName entry for tableName %v", tableName)
	}

	columns, ok := transformedColumnNamesByTableName[tableName]
	if !ok {
		return nil, fmt.Errorf("no columnNamesByTableName entry for tableName %v", tableName)
	}

	table, ok := tableByName[tableName]
	if !ok {
		return nil, fmt.Errorf("no tableByName entry for tableName %v", tableName)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		status := http.StatusOK
		var body []byte
		var err error

		defer func() {
			w.Header().Add("Content-Type", "application/json")

			w.WriteHeader(status)

			if err != nil {
				_, _ = w.Write([]byte(fmt.Sprintf("{\"error\": %#+v}", err.Error())))
				return
			}

			_, _ = w.Write(body)
		}()

		if r.Method != http.MethodGet {
			status = http.StatusMethodNotAllowed
			err = fmt.Errorf("%v; wanted %v, got %v", "StatusMethodNotAllowed", http.MethodGet, r.Method)
			return
		}

		limit := helpers.Ptr(50)
		rawLimit := strings.TrimSpace(r.URL.Query().Get("limit"))
		if len(rawLimit) > 0 {
			var parsedLimit int64
			parsedLimit, err = strconv.ParseInt(rawLimit, 10, 64)
			if err != nil {
				status = http.StatusBadRequest
				err = fmt.Errorf("failed to parse limit %#+v as int; err: %v", limit, err.Error())
				return
			}
			limit = helpers.Ptr(int(parsedLimit))
		}

		offset := helpers.Ptr(50)
		rawOffset := strings.TrimSpace(r.URL.Query().Get("offset"))
		if len(rawOffset) > 0 {
			var parsedOffset int64
			parsedOffset, err = strconv.ParseInt(rawOffset, 10, 64)
			if err != nil {
				status = http.StatusBadRequest
				err = fmt.Errorf("failed to parse offset %#+v as int; err: %v", offset, err.Error())
				return
			}
			offset = helpers.Ptr(int(parsedOffset))
		}

		descending := strings.HasPrefix(strings.ToLower(strings.TrimSpace(r.URL.Query().Get("order"))), "desc")

		orderBy := make([]string, 0)
		for k, vs := range r.URL.Query() {
			if k != "order_by" {
				continue
			}

			for _, v := range vs {
				_, ok := table.ColumnByName[v]
				if !ok {
					status = http.StatusBadRequest
					err = fmt.Errorf("bad order by; no table.ColumnByName entry for columnName %v", v)
					return
				}

				found := false

				for _, existingOrderBy := range orderBy {
					if existingOrderBy == v {
						found = true
						break
					}
				}

				if found {
					continue
				}

				orderBy = append(orderBy, v)
			}
		}

		var order *string
		if len(orderBy) > 0 {
			if descending {
				order = Descending(orderBy...)
			} else {
				order = Ascending(orderBy...)
			}
		}

		wheres := make([]string, 0)
		for k, vs := range r.URL.Query() {
			if k == "order" || k == "order_by" || k == "limit" || k == "offset" {
				continue
			}

			field := k
			matcher := "="

			parts := strings.Split(k, "__")
			if len(parts) > 1 {
				field = strings.Join(parts[0:len(parts)-1], "__")
				lastPart := strings.ToLower(parts[len(parts)-1])
				if lastPart == "eq" {
					matcher = "="
				} else if lastPart == "ne" {
					matcher = "!="
				} else if lastPart == "gt" {
					matcher = ">"
				} else if lastPart == "lt" {
					matcher = "<"
				} else if lastPart == "gte" {
					matcher = ">="
				} else if lastPart == "lte" {
					matcher = "<="
				} else if lastPart == "ilike" {
					matcher = "ILIKE"
				} else if lastPart == "in" {
					matcher = "IN"
				} else if lastPart == "nin" || lastPart == "not_in" {
					matcher = "NOT IN"
				}
			}

			innerWheres := make([]string, 0)

			for _, v := range vs {
				column, ok := table.ColumnByName[field]
				if !ok {
					status = http.StatusBadRequest
					err = fmt.Errorf("bad filter; no table.ColumnByName entry for columnName %v", field)
					return
				}

				if matcher != "IN" && matcher != "NOT IN" {
					if column.TypeTemplate == "string" || column.TypeTemplate == "uuid" {
						v = fmt.Sprintf("%#+v", v)
						innerWheres = append(
							innerWheres,
							fmt.Sprintf("%v %v '%v'", field, matcher, v[1:len(v)-1]),
						)
					} else {
						innerWheres = append(
							innerWheres,
							fmt.Sprintf("%v %v %v", field, matcher, v),
						)
					}
				} else {
					itemsBefore := strings.Split(v, ",")
					itemsAfter := make([]string, 0)

					for _, x := range itemsBefore {
						x = strings.TrimSpace(x)
						if x == "" {
							continue
						}

						if column.TypeTemplate == "string" || column.TypeTemplate == "uuid" {
							x = fmt.Sprintf("%#+v", x)
							itemsAfter = append(itemsAfter, fmt.Sprintf("'%v'", x[1:len(x)-1]))
						} else {
							itemsAfter = append(itemsAfter, x)
						}
					}

					innerWheres = append(
						innerWheres,
						fmt.Sprintf("%v %v (%v)", field, matcher, strings.Join(itemsAfter, ", ")),
					)
				}
			}

			wheres = append(wheres, fmt.Sprintf("(%v)", strings.Join(innerWheres, " OR ")))
		}

		var items []types.DjangolangObject
		items, err = selectFunc(r.Context(), db, columns, order, limit, offset, wheres...)
		if err != nil {
			status = http.StatusInternalServerError
			err = fmt.Errorf("query failed; err: %v", err.Error())
			return
		}

		body, err = json.Marshal(items)
		if err != nil {
			status = http.StatusInternalServerError
			err = fmt.Errorf("serialization failed; err: %v", err.Error())
		}
	}, nil
}

func GetInsertHandlerForTableName(tableName string, db *sqlx.DB) (types.InsertHandler, error) {
	insertFunc, ok := insertFuncByTableName[tableName]
	if !ok {
		_, ok = selectFuncByTableName[tableName]
		if ok {
			return nil, nil
		}

		return nil, fmt.Errorf("no insertFuncByTableName entry for tableName %v", tableName)
	}

	deserializeFunc, ok := deserializeFuncByTableName[tableName]
	if !ok {
		return nil, fmt.Errorf("no deserializeFuncByTableName entry for tableName %v", tableName)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		status := http.StatusCreated
		var body []byte
		var err error

		defer func() {
			w.Header().Add("Content-Type", "application/json")

			w.WriteHeader(status)

			if err != nil {
				_, _ = w.Write([]byte(fmt.Sprintf("{\"error\": %#+v}", err.Error())))
				return
			}

			_, _ = w.Write(body)
		}()

		if r.Method != http.MethodPost {
			status = http.StatusMethodNotAllowed
			err = fmt.Errorf("%v; wanted %v, got %v", "StatusMethodNotAllowed", http.MethodPost, r.Method)
			return
		}

		b, err := io.ReadAll(r.Body)
		if err != nil {
			status = http.StatusBadRequest
			err = fmt.Errorf("read failed; err: %v", err.Error())
			return
		}

		object, err := deserializeFunc(b)
		if err != nil {
			status = http.StatusBadRequest
			err = fmt.Errorf("deserialization of %v failed; err: %v", string(b), err.Error())
			return
		}

		item, err := insertFunc(r.Context(), db, object)
		if err != nil {
			status = http.StatusInternalServerError
			if strings.Contains(err.Error(), "value violates unique constraint") {
				status = http.StatusConflict
			}

			err = fmt.Errorf("insert query of %#+v failed; err: %v", object, err.Error())
			return
		}

		body, err = json.Marshal(item)
		if err != nil {
			status = http.StatusInternalServerError
			err = fmt.Errorf("serialization failed; err: %v", err.Error())
			return
		}
	}, nil
}

func GetUpdateHandlerForTableName(tableName string, db *sqlx.DB) (types.UpdateHandler, error) {
	updateFunc, ok := updateFuncByTableName[tableName]
	if !ok {
		_, ok = selectFuncByTableName[tableName]
		if ok {
			return nil, nil
		}

		return nil, fmt.Errorf("no updateFuncByTableName entry for tableName %v", tableName)
	}

	selectFunc, ok := selectFuncByTableName[tableName]
	if !ok {
		return nil, fmt.Errorf("no selectFuncByTableName entry for tableName %v", tableName)
	}

	table, ok := tableByName[tableName]
	if !ok {
		return nil, fmt.Errorf("no tableByName entry for tableName %v", tableName)
	}

	deserializeFunc, ok := deserializeFuncByTableName[tableName]
	if !ok {
		return nil, fmt.Errorf("no deserializeFuncByTableName entry for tableName %v", tableName)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		status := http.StatusOK
		var body []byte
		var err error

		defer func() {
			w.Header().Add("Content-Type", "application/json")

			w.WriteHeader(status)

			if err != nil {
				_, _ = w.Write([]byte(fmt.Sprintf("{\"error\": %#+v}", err.Error())))
				return
			}

			_, _ = w.Write(body)
		}()

		if r.Method != http.MethodPut {
			status = http.StatusMethodNotAllowed
			err = fmt.Errorf("%v; wanted %v, got %v", "StatusMethodNotAllowed", http.MethodPut, r.Method)
			return
		}

		primaryKeyValue := chi.URLParam(r, "primaryKeyValue")

		var safePrimaryKeyValue any

		err = json.Unmarshal([]byte(primaryKeyValue), &safePrimaryKeyValue)
		if err != nil {
			status = http.StatusBadRequest
			err = fmt.Errorf(
				"parameters failed; err: %v",
				fmt.Errorf("failed to prepare %#+v: %v", primaryKeyValue, err).Error(),
			)
		}

		objects, err := selectFunc(
			r.Context(),
			db,
			nil,
			nil,
			nil,
			nil,
			fmt.Sprintf("%v = %v", table.PrimaryKeyColumn.Name, safePrimaryKeyValue),
		)
		if err != nil {
			err = fmt.Errorf("select for update failed; err: %v", err.Error())
			status = http.StatusInternalServerError
			return
		}

		if len(objects) < 1 {
			err = fmt.Errorf(
				"select for delete failed; err: %v",
				fmt.Errorf("failed to find object for primary key value: %#+v", safePrimaryKeyValue).Error(),
			)
			status = http.StatusNotFound
			return
		}

		if len(objects) > 1 {
			err = fmt.Errorf(
				"select for delete failed; err: %v",
				fmt.Errorf("found too many objects for primary key value: %#+v", safePrimaryKeyValue).Error(),
			)
			status = http.StatusInternalServerError
			return
		}

		var existingObject types.DjangolangObject = objects[0]

		b, err := io.ReadAll(r.Body)
		if err != nil {
			status = http.StatusBadRequest
			err = fmt.Errorf("read failed; err: %v", err.Error())
			return
		}

		updatedObject, err := deserializeFunc(b)
		if err != nil {
			status = http.StatusBadRequest
			err = fmt.Errorf("deserialization of %v failed; err: %v", string(b), err.Error())
			return
		}

		primaryKey, err := existingObject.GetPrimaryKey()
		if err != nil {
			status = http.StatusBadRequest
			err = fmt.Errorf("get primary key for %v failed; err: %v", string(b), err.Error())
			return
		}

		err = updatedObject.SetPrimaryKey(primaryKey)
		if err != nil {
			status = http.StatusBadRequest
			err = fmt.Errorf("set primary key for %v failed; err: %v", string(b), err.Error())
			return
		}

		item, err := updateFunc(r.Context(), db, updatedObject)
		if err != nil {
			status = http.StatusInternalServerError
			if strings.Contains(err.Error(), "value violates unique constraint") {
				status = http.StatusConflict
			}

			err = fmt.Errorf("update query of %#+v failed; err: %v", updatedObject, err.Error())
			return
		}

		body, err = json.Marshal(item)
		if err != nil {
			status = http.StatusInternalServerError
			err = fmt.Errorf("serialization failed; err: %v", err.Error())
			return
		}
	}, nil
}

func GetDeleteHandlerForTableName(tableName string, db *sqlx.DB) (types.DeleteHandler, error) {
	deleteFunc, ok := deleteFuncByTableName[tableName]
	if !ok {
		_, ok = selectFuncByTableName[tableName]
		if ok {
			return nil, nil
		}

		return nil, fmt.Errorf("no deleteFuncByTableName entry for tableName %v", tableName)
	}

	selectFunc, ok := selectFuncByTableName[tableName]
	if !ok {
		return nil, fmt.Errorf("no selectFuncByTableName entry for tableName %v", tableName)
	}

	table, ok := tableByName[tableName]
	if !ok {
		return nil, fmt.Errorf("no tableByName entry for tableName %v", tableName)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		status := http.StatusNoContent
		var body []byte
		var err error

		defer func() {
			w.Header().Add("Content-Type", "application/json")

			w.WriteHeader(status)

			if err != nil {
				_, _ = w.Write([]byte(fmt.Sprintf("{\"error\": %#+v}", err.Error())))
				return
			}

			_, _ = w.Write(body)
		}()

		if r.Method != http.MethodDelete {
			status = http.StatusMethodNotAllowed
			err = fmt.Errorf("%v; wanted %v, got %v", "StatusMethodNotAllowed", http.MethodDelete, r.Method)
			return
		}

		primaryKeyValue := chi.URLParam(r, "primaryKeyValue")

		var safePrimaryKeyValue any

		err = json.Unmarshal([]byte(primaryKeyValue), &safePrimaryKeyValue)
		if err != nil {
			status = http.StatusBadRequest
			err = fmt.Errorf(
				"parameters failed; err: %v",
				fmt.Errorf("failed to prepare %#+v: %v", primaryKeyValue, err).Error(),
			)
		}

		objects, err := selectFunc(
			r.Context(),
			db,
			nil,
			nil,
			nil,
			nil,
			fmt.Sprintf("%v = %v", table.PrimaryKeyColumn.Name, safePrimaryKeyValue),
		)
		if err != nil {
			err = fmt.Errorf("select for delete failed; err: %v", err.Error())
			status = http.StatusInternalServerError
			return
		}

		if len(objects) < 1 {
			err = fmt.Errorf(
				"select for delete failed; err: %v",
				fmt.Errorf("failed to find object for primary key value: %#+v", safePrimaryKeyValue).Error(),
			)
			status = http.StatusNotFound
			return
		}

		if len(objects) > 1 {
			err = fmt.Errorf(
				"select for delete failed; err: %v",
				fmt.Errorf("found too many objects for primary key value: %#+v", safePrimaryKeyValue).Error(),
			)
			status = http.StatusInternalServerError
			return
		}

		var object types.DjangolangObject = objects[0]

		err = deleteFunc(r.Context(), db, object)
		if err != nil {
			status = http.StatusInternalServerError
			if strings.Contains(err.Error(), "value violates unique constraint") {
				status = http.StatusConflict
			}

			err = fmt.Errorf("query failed; err: %v", err.Error())
			return
		}

		status = http.StatusNoContent
		body = []byte("")
	}, nil
}

func GetSelectHandlerByEndpointName(db *sqlx.DB) (map[string]types.SelectHandler, error) {
	thisSelectHandlerByEndpointName := make(map[string]types.SelectHandler)

	for _, tableName := range TableNames {
		selectHandler, err := GetSelectHandlerForTableName(tableName, db)
		if err != nil {
			return nil, err
		}

		thisSelectHandlerByEndpointName[caps.ToKebab(pluralize.Singular(tableName))] = selectHandler
		thisSelectHandlerByEndpointName[caps.ToKebab(pluralize.Plural(tableName))] = selectHandler
	}

	return thisSelectHandlerByEndpointName, nil
}

func GetInsertHandlerByEndpointName(db *sqlx.DB) (map[string]types.InsertHandler, error) {
	thisInsertHandlerByEndpointName := make(map[string]types.InsertHandler)

	for _, tableName := range TableNames {
		insertHandler, err := GetInsertHandlerForTableName(tableName, db)
		if err != nil {
			return nil, err
		}

		thisInsertHandlerByEndpointName[caps.ToKebab(pluralize.Singular(tableName))] = insertHandler
		thisInsertHandlerByEndpointName[caps.ToKebab(pluralize.Plural(tableName))] = insertHandler
	}

	return thisInsertHandlerByEndpointName, nil
}

func GetUpdateHandlerByEndpointName(db *sqlx.DB) (map[string]types.UpdateHandler, error) {
	thisUpdateHandlerByEndpointName := make(map[string]types.UpdateHandler)

	for _, tableName := range TableNames {
		updateHandler, err := GetUpdateHandlerForTableName(tableName, db)
		if err != nil {
			return nil, err
		}

		thisUpdateHandlerByEndpointName[caps.ToKebab(pluralize.Singular(tableName))] = updateHandler
		thisUpdateHandlerByEndpointName[caps.ToKebab(pluralize.Plural(tableName))] = updateHandler
	}

	return thisUpdateHandlerByEndpointName, nil
}

func GetDeleteHandlerByEndpointName(db *sqlx.DB) (map[string]types.DeleteHandler, error) {
	thisDeleteHandlerByEndpointName := make(map[string]types.DeleteHandler)

	for _, tableName := range TableNames {
		deleteHandler, err := GetDeleteHandlerForTableName(tableName, db)
		if err != nil {
			return nil, err
		}

		thisDeleteHandlerByEndpointName[caps.ToKebab(pluralize.Singular(tableName))] = deleteHandler
		thisDeleteHandlerByEndpointName[caps.ToKebab(pluralize.Plural(tableName))] = deleteHandler
	}

	return thisDeleteHandlerByEndpointName, nil
}

func RunServer(ctx context.Context, outerChanges chan types.Change) error {
	port, err := helpers.GetPort()
	if err != nil {
		return err
	}

	db, err := helpers.GetDBFromEnvironment(ctx)
	if err != nil {
		return err
	}
	defer func() {
		_ = db.Close()
	}()

	r := chi.NewRouter()

	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.StripSlashes)
	r.Use(middleware.AllowContentType("application/json"))

	selectHandlerByEndpointName, err := GetSelectHandlerByEndpointName(db)
	if err != nil {
		return err
	}

	insertHandlerByEndpointName, err := GetInsertHandlerByEndpointName(db)
	if err != nil {
		return err
	}

	updateHandlerByEndpointName, err := GetUpdateHandlerByEndpointName(db)
	if err != nil {
		return err
	}

	deleteHandlerByEndpointName, err := GetDeleteHandlerByEndpointName(db)
	if err != nil {
		return err
	}

	endpointNames := maps.Keys(selectHandlerByEndpointName)
	slices.Sort(endpointNames)

	for _, endpointName := range endpointNames {
		selectHandler := selectHandlerByEndpointName[endpointName]
		if selectHandler != nil {
			pathName := fmt.Sprintf("/%v", endpointName)
			logger.Printf("registering GET for select: %v", pathName)
			r.Get(pathName, selectHandler)
		}

		insertHandler := insertHandlerByEndpointName[endpointName]
		if insertHandler != nil {
			pathName := fmt.Sprintf("/%v", endpointName)
			logger.Printf("registering POST for insert: %v", pathName)
			r.Post(pathName, insertHandler)
		}

		updateHandler := updateHandlerByEndpointName[endpointName]
		deleteHandler := deleteHandlerByEndpointName[endpointName]
		if updateHandler != nil && deleteHandler != nil {
			pathName := fmt.Sprintf("/%v/{primaryKeyValue}", endpointName)

			r.Route(pathName, func(r chi.Router) {
				logger.Printf("registering PUT for update: %v", pathName)
				r.Put("/", updateHandler)

				logger.Printf("registering DELETE for delete: %v", pathName)
				r.Delete("/", deleteHandler)
			})
		}
	}

	upgrader := websocket.Upgrader{
		HandshakeTimeout: timeout,
		Subprotocols:     []string{"djangolang"},
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		EnableCompression: true,
	}

	subscriberByID := make(map[uuid.UUID]*Subscriber)

	r.Get("/__subscribe", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			logger.Printf("warning: failed to upgrade websocket; err: %v", err)
			return
		}

		subscriber := &Subscriber{
			ID:         uuid.Must(uuid.NewRandom()),
			TableNames: make(map[string]struct{}),
			Messages:   make(chan []byte, 1024),
		}

		mu.Lock()
		subscriberByID[subscriber.ID] = subscriber
		mu.Unlock()

		connCtx, connCancel := context.WithCancel(ctx)
		go func() {
			<-connCtx.Done()
			conn.Close()

			mu.Lock()
			delete(subscriberByID, subscriber.ID)
			mu.Unlock()

			logger.Printf("subscriber %v left", subscriber.ID)
		}()

		logger.Printf("subscriber %v joined", subscriber.ID)

		conn.SetCloseHandler(func(code int, text string) error {
			logger.Printf("warning: websocket closed; code: %v, text: %#+v", code, text)
			connCancel()
			return nil
		})

		writeMu := new(sync.Mutex)

		// ping loop
		go func() {
			defer func() {
				connCancel()
			}()

			t := time.NewTimer(pingInterval)
			defer t.Stop()

			for {
				select {
				case <-connCtx.Done():
					return
				case <-t.C:
					writeMu.Lock()
					err = conn.WriteControl(websocket.PingMessage, nil, time.Now().Add(timeout))
					writeMu.Unlock()
					if err != nil {
						logger.Printf("error: websocket ping failed; err: %v", err)
						return
					}
				}
			}
		}()

		// read loop
		go func() {
			defer func() {
				connCancel()
			}()

			for {
				select {
				case <-connCtx.Done():
					return
				default:
				}

				messageType, b, err := conn.ReadMessage()
				if err != nil {
					logger.Printf("error: websocket read failed; err: %v", err)
					return
				}

				switch messageType {
				case websocket.TextMessage:
					logger.Printf("websocket message received: %#+v", string(b))
				default:
					continue
				}
			}
		}()

		// write loop
		go func() {
			defer func() {
				connCancel()
			}()

			for {
				select {
				case <-connCtx.Done():
					return
				case message := <-subscriber.Messages:
					writeMu.Lock()
					err = conn.WriteMessage(websocket.TextMessage, message)
					if err != nil {
						logger.Printf("error: websocket write failed; err: %v", err)
						return
					}
					writeMu.Unlock()
				}
			}
		}()

		runtime.Gosched()
	})

	errs := make(chan error, 2)

	if os.Getenv("DJANGOLANG_SKIP_STREAM") != "1" {
		innerChanges := make(chan types.Change, 1024)

		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				case change := <-innerChanges:

					message, err := json.Marshal(change)
					if err != nil {
						logger.Printf("warning: failed to marshal %#+v to JSON; err: %v", change, err)
						continue
					}

					mu.Lock()
					subscriptions := maps.Values(subscriberByID)
					mu.Unlock()

					for _, subscription := range subscriptions {
						// TODO: filter by table names

						select {
						case subscription.Messages <- message:
						default:
						}
					}

					select {
					case outerChanges <- change:
					default:
					}
				}
			}
		}()

		go func() {
			logger.Printf("starting stream...")

			errs <- runStream(ctx, innerChanges)
		}()
	}

	go func() {
		logger.Printf("starting server...")
		errs <- http.ListenAndServe(
			fmt.Sprintf(":%v", port),
			r,
		)
	}()

	select {
	case <-ctx.Done():
		break
	case err = <-errs:
		return fmt.Errorf("stream unexpectedly exited; err: %v", err)
	}

	return nil
}
