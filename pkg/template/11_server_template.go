package template

import "strings"

var serverTemplate = strings.TrimSpace(`
package templates

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/chanced/caps"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/initialed85/djangolang/pkg/helpers"
	"github.com/jmoiron/sqlx"
	"golang.org/x/exp/maps"
)

func GetSelectHandlerForTableName(tableName string, db *sqlx.DB) (SelectHandler, error) {
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

		var items []DjangolangObject
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

func GetInsertHandlerForTableName(tableName string, db *sqlx.DB) (InsertHandler, error) {
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

func GetUpdateHandlerForTableName(tableName string, db *sqlx.DB) (UpdateHandler, error) {
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

		if r.Method != http.MethodPatch {
			status = http.StatusMethodNotAllowed
			err = fmt.Errorf("%v; wanted %v, got %v", "StatusMethodNotAllowed", http.MethodPatch, r.Method)
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

		var existingObject DjangolangObject = objects[0]

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

func GetDeleteHandlerForTableName(tableName string, db *sqlx.DB) (DeleteHandler, error) {
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

		var object DjangolangObject = objects[0]

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

func GetSelectHandlerByEndpointName(db *sqlx.DB) (map[string]SelectHandler, error) {
	thisSelectHandlerByEndpointName := make(map[string]SelectHandler)

	for _, tableName := range TableNames {
		endpointName := caps.ToKebab[string](tableName)

		selectHandler, err := GetSelectHandlerForTableName(tableName, db)
		if err != nil {
			return nil, err
		}

		thisSelectHandlerByEndpointName[endpointName] = selectHandler
	}

	return thisSelectHandlerByEndpointName, nil
}

func GetInsertHandlerByEndpointName(db *sqlx.DB) (map[string]InsertHandler, error) {
	thisInsertHandlerByEndpointName := make(map[string]InsertHandler)

	for _, tableName := range TableNames {
		endpointName := caps.ToKebab[string](tableName)

		insertHandler, err := GetInsertHandlerForTableName(tableName, db)
		if err != nil {
			return nil, err
		}

		thisInsertHandlerByEndpointName[endpointName] = insertHandler
	}

	return thisInsertHandlerByEndpointName, nil
}

func GetUpdateHandlerByEndpointName(db *sqlx.DB) (map[string]UpdateHandler, error) {
	thisUpdateHandlerByEndpointName := make(map[string]UpdateHandler)

	for _, tableName := range TableNames {
		endpointName := caps.ToKebab[string](tableName)

		updateHandler, err := GetUpdateHandlerForTableName(tableName, db)
		if err != nil {
			return nil, err
		}

		thisUpdateHandlerByEndpointName[endpointName] = updateHandler
	}

	return thisUpdateHandlerByEndpointName, nil
}

func GetDeleteHandlerByEndpointName(db *sqlx.DB) (map[string]DeleteHandler, error) {
	thisDeleteHandlerByEndpointName := make(map[string]DeleteHandler)

	for _, tableName := range TableNames {
		endpointName := caps.ToKebab[string](tableName)

		deleteHandler, err := GetDeleteHandlerForTableName(tableName, db)
		if err != nil {
			return nil, err
		}

		thisDeleteHandlerByEndpointName[endpointName] = deleteHandler
	}

	return thisDeleteHandlerByEndpointName, nil
}

func RunServer(ctx context.Context) error {
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

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

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
				logger.Printf("registering UPDATE for update: %v", pathName)
				r.Patch("/", updateHandler)

				logger.Printf("registering DELETE for delete: %v", pathName)
				r.Delete("/", deleteHandler)
			})
		}
	}

	errs := make(chan error, 2)

	if os.Getenv("DJANGOLANG_SKIP_STREAM") != "1" {
		go func() {
			logger.Printf("starting stream...")

			errs <- runStream(ctx)
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
`) + "\n"
