package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
	"github.com/initialed85/djangolang/pkg/stream"
	"github.com/initialed85/structmeta/pkg/introspect"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/exp/maps"
)

type HTTPMiddleware func(http.Handler) http.Handler

type ObjectMiddleware func()

type GetRouterFn func(*pgxpool.Pool, *redis.Pool, []HTTPMiddleware, []ObjectMiddleware, WaitForChange) chi.Router

type WaitForChange func(context.Context, []stream.Action, string, uint32) (*Change, error)

type EmptyPathParams struct{}

type EmptyQueryParams struct{}

type EmptyRequest struct{}

type EmptyResponse struct{}

var introspectedEmptyPathParams *introspect.Object

var introspectedEmptyQueryParams *introspect.Object

var introspectedEmptyRequest *introspect.Object

var introspectedEmptyResponse *introspect.Object

func init() {
	var err error

	introspectedEmptyPathParams, err = introspect.Introspect(EmptyPathParams{})
	if err != nil {
		panic(err)
	}

	introspectedEmptyQueryParams, err = introspect.Introspect(EmptyQueryParams{})
	if err != nil {
		panic(err)
	}

	introspectedEmptyRequest, err = introspect.Introspect(EmptyRequest{})
	if err != nil {
		panic(err)
	}

	introspectedEmptyResponse, err = introspect.Introspect(EmptyResponse{})
	if err != nil {
		panic(err)
	}
}

type CustomHTTPHandler[T any, S any, Q any, R any] struct {
	Method                                        string
	Path                                          string
	PathParams                                    T
	QueryParams                                   S
	Request                                       Q
	Response                                      R
	Status                                        int
	Handle                                        func(context.Context, T, S, Q, any) (*R, error)
	AllPathParamKeys                              map[string]struct{}
	RequiredPathParamKeys                         map[string]struct{}
	AllQueryParamKeys                             map[string]struct{}
	RequiredQueryParamKeys                        map[string]struct{}
	PathParamsIntrospectedStructFieldObjectByKey  map[string]*introspect.StructFieldObject
	QueryParamsIntrospectedStructFieldObjectByKey map[string]*introspect.StructFieldObject
	RequestIntrospectedObject                     *introspect.Object
	ResponseIntrospectedObject                    *introspect.Object
	PathParamsIsEmpty                             bool
	QueryParamsIsEmpty                            bool
	RequestIsEmpty                                bool
	ResponseIsEmpty                               bool
}

func GetCustomHTTPHandler[T any, S any, Q any, R any](method string, path string, status int, handle func(context.Context, T, S, Q, any) (*R, error)) (*CustomHTTPHandler[T, S, Q, R], error) {
	s := CustomHTTPHandler[T, S, Q, R]{
		Method:                 method,
		Path:                   path,
		Handle:                 handle,
		Status:                 status,
		AllPathParamKeys:       make(map[string]struct{}, 0),
		RequiredPathParamKeys:  make(map[string]struct{}, 0),
		AllQueryParamKeys:      make(map[string]struct{}, 0),
		RequiredQueryParamKeys: make(map[string]struct{}, 0),
		PathParamsIntrospectedStructFieldObjectByKey:  make(map[string]*introspect.StructFieldObject),
		QueryParamsIntrospectedStructFieldObjectByKey: make(map[string]*introspect.StructFieldObject),
		RequestIntrospectedObject:                     nil,
		ResponseIntrospectedObject:                    nil,
	}

	//
	// handle path params config
	//

	pathParams := *new(T)
	pathParamsIntrospectedObject, err := introspect.Introspect(pathParams)
	if err != nil {
		return nil, err
	}

	if pathParamsIntrospectedObject != introspectedEmptyPathParams {
		for _, structFieldObject := range pathParamsIntrospectedObject.StructFields {
			key := structFieldObject.Tag.Get("json")
			if key == "" {
				key = structFieldObject.Field
			}

			s.AllPathParamKeys[key] = struct{}{}

			if structFieldObject.PointerValue == nil {
				s.RequiredPathParamKeys[key] = struct{}{}
			}

			s.PathParamsIntrospectedStructFieldObjectByKey[key] = structFieldObject
		}
	} else {
		s.PathParamsIsEmpty = true
	}

	//
	// handle query params config
	//

	queryParams := *new(S)
	queryParamsIntrospectedObject, err := introspect.Introspect(queryParams)
	if err != nil {
		return nil, err
	}

	if queryParamsIntrospectedObject != introspectedEmptyQueryParams {
		for _, structFieldObject := range queryParamsIntrospectedObject.StructFields {
			key := structFieldObject.Tag.Get("json")
			if key == "" {
				key = structFieldObject.Field
			}

			s.AllQueryParamKeys[key] = struct{}{}

			if structFieldObject.PointerValue == nil {
				s.RequiredQueryParamKeys[key] = struct{}{}
			}

			s.QueryParamsIntrospectedStructFieldObjectByKey[key] = structFieldObject
		}
	} else {
		s.QueryParamsIsEmpty = true
	}

	if queryParamsIntrospectedObject.MapKey != nil && queryParamsIntrospectedObject.MapValue != nil {
		s.AllQueryParamKeys = nil
		s.RequiredPathParamKeys = nil
	}

	s.RequestIntrospectedObject, err = introspect.Introspect(*new(Q))
	if err != nil {
		return nil, err
	}

	if s.RequestIntrospectedObject == introspectedEmptyRequest {
		s.RequestIntrospectedObject = nil
		s.RequestIsEmpty = true
	}

	s.ResponseIntrospectedObject, err = introspect.Introspect(*new(R))
	if err != nil {
		return nil, err
	}

	if s.ResponseIntrospectedObject == introspectedEmptyResponse {
		s.ResponseIntrospectedObject = nil
		s.ResponseIsEmpty = true
	}

	return &s, nil
}

func (s *CustomHTTPHandler[T, S, Q, R]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rc := chi.RouteContext(r.Context())

	//
	// handle path params
	//

	unrecognizedPathParams := make([]string, 0)
	unparseablePathParams := make([]string, 0)
	rawPathParams := make(map[string]any)
	if rc != nil {
		for i, k := range rc.URLParams.Keys {
			// not sure what this is- not relevant though that's for sure
			if k == "*" {
				continue
			}

			_, ok := s.AllPathParamKeys[k]
			if !ok {
				unrecognizedPathParams = append(unrecognizedPathParams, k)
				continue
			}

			var v any

			err := json.Unmarshal([]byte(rc.URLParams.Values[i]), &v)
			if err != nil {
				err = json.Unmarshal([]byte(fmt.Sprintf("\"%s\"", rc.URLParams.Values[i])), &v)
				if err != nil {
					unparseablePathParams = append(unparseablePathParams, fmt.Sprintf("%s=%v", k, rc.URLParams.Values[i]))
					continue
				}
			}

			rawPathParams[k] = v
		}
	}

	b, err := json.Marshal(rawPathParams)
	if err != nil {
		HandleErrorResponse(
			w,
			http.StatusInternalServerError,
			fmt.Errorf("failed to convert rawPathParams %#+v to JSON; %v", rawPathParams, err),
		)
		return
	}

	pathParams := *new(T)
	err = json.Unmarshal(b, &pathParams)
	if err != nil {
		HandleErrorResponse(
			w,
			http.StatusInternalServerError,
			fmt.Errorf("failed to convert rawPathParams %v from JSON to %#+v; %v", string(b), pathParams, err),
		)
		return
	}

	if len(unrecognizedPathParams) > 0 {
		HandleErrorResponse(
			w,
			http.StatusBadRequest,
			fmt.Errorf("unrecognized path params: %s; wanted at most %s", strings.Join(unrecognizedPathParams, ", "), strings.Join(maps.Keys(s.AllPathParamKeys), ", ")),
		)
		return
	}

	if len(unparseablePathParams) > 0 {
		HandleErrorResponse(
			w,
			http.StatusBadRequest,
			fmt.Errorf("unparseable path params: %s", strings.Join(unparseablePathParams, ", ")),
		)
		return
	}

	//
	// handle query params
	//

	unrecognizedQueryParams := make([]string, 0)
	unparseableQueryParams := make([]string, 0)
	rawQueryParams := make(map[string]any)
	if r.URL != nil {
		for k, vs := range r.URL.Query() {
			if s.AllQueryParamKeys != nil {
				_, ok := s.AllQueryParamKeys[k]
				if !ok {
					unrecognizedQueryParams = append(unrecognizedQueryParams, k)
					continue
				}
			}

			for _, rawV := range vs {
				timestamp, err := time.Parse(time.RFC3339Nano, strings.ReplaceAll(rawV, " ", "+"))
				if err == nil {
					rawQueryParams[k] = timestamp
					continue
				}

				var v any

				err = json.Unmarshal([]byte(rawV), &v)
				if err != nil {
					err = json.Unmarshal([]byte(fmt.Sprintf("\"%s\"", rawV)), &v)
					if err != nil {
						unparseableQueryParams = append(unparseableQueryParams, fmt.Sprintf("%s=%v", k, rawV))
						continue
					}
				}

				rawQueryParams[k] = v
			}
		}
	}

	b, err = json.Marshal(rawQueryParams)
	if err != nil {
		HandleErrorResponse(
			w,
			http.StatusInternalServerError,
			fmt.Errorf("failed to convert rawQueryParams %#+v to JSON; %v", rawQueryParams, err),
		)
		return
	}

	queryParams := *new(S)
	err = json.Unmarshal(b, &queryParams)
	if err != nil {
		HandleErrorResponse(
			w,
			http.StatusInternalServerError,
			fmt.Errorf("failed to convert rawQueryParams %v from JSON %v; %#+v", b, queryParams, err),
		)
		return
	}

	if len(unrecognizedQueryParams) > 0 {
		HandleErrorResponse(
			w,
			http.StatusBadRequest,
			fmt.Errorf("unrecognized query params: %s", strings.Join(unrecognizedQueryParams, ", ")),
		)
		return
	}

	if len(unparseableQueryParams) > 0 {
		HandleErrorResponse(
			w,
			http.StatusBadRequest,
			fmt.Errorf("unparseable query params: %s", strings.Join(unparseableQueryParams, ", ")),
		)
		return
	}

	req := *new(Q)
	var rawReq any

	if s.RequestIntrospectedObject != nil {
		// TODO: do we need such strict validation? we're gonna treat it like JSON anyway, and it'll parse or it won't
		// contentType := strings.ToLower(strings.TrimSpace(r.Header.Get("Content-Type")))
		// if contentType != "application/json" {
		// 	HandleErrorResponse(
		// 		w,
		// 		http.StatusBadRequest,
		// 		fmt.Errorf("invalid Content-Type header: %s", fmt.Errorf("wanted %#+v, got %#+v", "application/json", contentType).Error()),
		// 	)
		// 	return
		// }

		reqBody, err := io.ReadAll(r.Body)
		if err != nil {
			HandleErrorResponse(
				w,
				http.StatusInternalServerError,
				fmt.Errorf("failed to read reqBody: %s", err.Error()),
			)
			return
		}

		err = json.Unmarshal(reqBody, &req)
		if err != nil {
			HandleErrorResponse(
				w,
				http.StatusInternalServerError,
				fmt.Errorf("failed to handle reqBody %s as JSON; %v", string(reqBody), err),
			)
			return
		}

		err = json.Unmarshal(reqBody, &rawReq)
		if err != nil {
			HandleErrorResponse(
				w,
				http.StatusInternalServerError,
				fmt.Errorf("failed to handle reqBody %s as JSON; %v", string(reqBody), err),
			)
			return
		}
	}

	res, err := s.Handle(r.Context(), pathParams, queryParams, req, rawReq)
	if err != nil {
		HandleErrorResponse(
			w,
			http.StatusBadRequest,
			fmt.Errorf("failed to invoke handler; %s", err.Error()),
		)
		return
	}

	b = []byte{}

	if s.ResponseIntrospectedObject != nil {
		b, err = json.Marshal(res)
		if err != nil {
			HandleErrorResponse(
				w,
				http.StatusInternalServerError,
				fmt.Errorf("failed to convert res %#+v to JSON; %v", res, err),
			)
			return
		}
	}

	w.WriteHeader(s.Status)

	if s.ResponseIntrospectedObject != nil {
		w.Header().Add("Content-Type", "application/json")
		_, _ = w.Write(b)
	}
}

type WithReload interface {
	Reload(context.Context, pgx.Tx, ...bool) error
}

type WithPrimaryKey interface {
	GetPrimaryKeyColumn() string
	GetPrimaryKeyValue() any
}

type WithInsert interface {
	Insert(context.Context, pgx.Tx) error
}

type Change struct {
	Timestamp time.Time      `json:"timestamp"`
	ID        uuid.UUID      `json:"id"`
	Action    stream.Action  `json:"action"`
	TableName string         `json:"table_name"`
	Item      map[string]any `json:"-"`
	Object    any            `json:"object"`
	Xid       uint32         `json:"xid"`
}

func (c *Change) String() string {
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

	return strings.TrimSpace(fmt.Sprintf(
		"(%s / %d) @ %s; %s %s %s",
		c.ID, c.Xid, c.Timestamp, c.Action, c.TableName, primaryKeySummary,
	))
}
