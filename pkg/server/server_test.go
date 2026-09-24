package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/initialed85/djangolang/internal/hack"
	"github.com/initialed85/djangolang/pkg/helpers"
	"github.com/initialed85/djangolang/pkg/introspect"
	"github.com/stretchr/testify/require"
)

type SomePathParams struct {
	PatchID uuid.UUID `json:"patch_id"`
}

type SomeQueryParams struct {
	SomeInt          int      `json:"some_int"`
	SomeString       string   `json:"some_string"`
	SomeOptionalBool *bool    `json:"some_optional_bool"`
	SomeIntArray     []int    `json:"some_int_array"`
	SomeStringArray  []string `json:"some_string_array"`
	SomeBoolArray    []bool   `json:"some_bool_array"`
}

type Cabbage struct {
	What  map[string]bool     `json:"what"`
	Is    []map[int]time.Time `json:"is"`
	This  float64             `json:"this"`
	Thing []float64           `json:"thing"`
	Buddy *int                `json:"buddy"`
}

type SomeRequest struct {
	Timestamp time.Time `json:"timestamp"`
	Cabbages  []Cabbage
}

type SomeResponse struct {
	Timestamp        time.Time `json:"timestamp"`
	Cabbages         []Cabbage `json:"cabbages"`
	FavouriteCabbage *Cabbage  `json:"favourite_cabbage"`
}

func TestBrowserCacheMaxAge(t *testing.T) {
	now := time.Date(2026, 9, 24, 12, 0, 0, 0, time.UTC)
	table := &introspect.Table{
		ColumnByName: map[string]*introspect.Column{
			"seen_at": {Name: "seen_at", TypeTemplate: "time.Time"},
			"name":    {Name: "name", TypeTemplate: "string"},
		},
	}

	tests := []struct {
		name      string
		method    string
		status    int
		query     map[string]any
		expectAge time.Duration
		expectSet bool
	}{
		{name: "no timestamp filter", method: http.MethodGet, status: http.StatusOK, query: map[string]any{"name__eq": "camera"}, expectAge: time.Second, expectSet: true},
		{name: "no filters", method: http.MethodGet, status: http.StatusOK, expectAge: time.Second, expectSet: true},
		{name: "gt always debounces", method: http.MethodGet, status: http.StatusOK, query: map[string]any{"seen_at__gt": now.Add(-24 * time.Hour)}, expectAge: time.Second, expectSet: true},
		{name: "gte always debounces", method: http.MethodGet, status: http.StatusOK, query: map[string]any{"seen_at__gte": now.Add(-24 * time.Hour)}, expectAge: time.Second, expectSet: true},
		{name: "timestamp within one hour", method: http.MethodGet, status: http.StatusOK, query: map[string]any{"seen_at__lt": now.Add(-30 * time.Minute)}, expectAge: time.Minute, expectSet: true},
		{name: "historic timestamp", method: http.MethodGet, status: http.StatusOK, query: map[string]any{"seen_at__eq": now.Add(-2 * time.Hour)}, expectAge: time.Hour, expectSet: true},
		{name: "unknown timestamp value debounces", method: http.MethodGet, status: http.StatusOK, query: map[string]any{"seen_at__lt": "not-a-time"}, expectAge: time.Second, expectSet: true},
		{name: "non-GET is not cached", method: http.MethodPost, status: http.StatusOK, expectSet: false},
		{name: "error response is not cached", method: http.MethodGet, status: http.StatusInternalServerError, expectSet: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actualAge, actualSet := browserCacheMaxAge(test.method, test.status, test.query, table, now)
			require.Equal(t, test.expectSet, actualSet)
			if test.expectSet {
				require.Equal(t, test.expectAge, actualAge)
			}
		})
	}
}

func TestBrowserCacheHeadersOnGET(t *testing.T) {
	handler, err := GetHTTPHandler(
		http.MethodGet,
		"/browser-cache",
		http.StatusOK,
		func(context.Context, EmptyPathParams, EmptyQueryParams, EmptyRequest, any) (*SomeResponse, error) {
			return &SomeResponse{Timestamp: time.Date(2026, 9, 24, 12, 0, 0, 0, time.UTC)}, nil
		},
	)
	require.NoError(t, err)

	router := chi.NewRouter()
	router.Get(handler.FullPath, handler.ServeHTTP)

	request := httptest.NewRequest(http.MethodGet, "/browser-cache", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	require.Equal(t, http.StatusOK, response.Code)
	require.Equal(t, "max-age=1", response.Header().Get("Cache-Control"))
	require.Empty(t, response.Header().Get("ETag"))
}

func TestServer(t *testing.T) {
	t.Run("GetCustomHTTPHandler", func(t *testing.T) {
		customHTTPHandler, err := GetHTTPHandler(
			http.MethodGet,
			"/api/add-cabbages/{patch_id}",
			http.StatusOK,
			func(ctx context.Context, pathParams SomePathParams, queryParams SomeQueryParams, req SomeRequest, rawReq any) (*SomeResponse, error) {
				res := &SomeResponse{
					Timestamp: time.Now(),
					Cabbages: []Cabbage{{
						What:  map[string]bool{"hello": true, "world": true},
						Is:    []map[int]time.Time{{1: time.Now()}},
						This:  69,
						Thing: []float64{42.0},
						Buddy: nil,
					},
					},
					FavouriteCabbage: &Cabbage{
						What:  map[string]bool{"goodbye": true, "world": true},
						Is:    []map[int]time.Time{{2: time.Now()}},
						This:  1337,
						Thing: []float64{8008135.0},
						Buddy: helpers.Ptr(42),
					},
				}

				return res, nil
			},
		)
		require.NoError(t, err)

		router := chi.NewRouter()
		router.Get(customHTTPHandler.FullPath, customHTTPHandler.ServeHTTP)

		patchID := uuid.Must(uuid.NewRandom())

		requestBody := SomeRequest{
			Timestamp: time.Now(),
			Cabbages: []Cabbage{
				{
					What:  map[string]bool{"hello": true, "world": true},
					Is:    []map[int]time.Time{{1: time.Now()}},
					This:  69,
					Thing: []float64{42.0},
					Buddy: nil,
				},
			},
		}

		rawRequestBody, err := json.Marshal(requestBody)
		require.NoError(t, err)

		r := httptest.NewRequest(
			http.MethodGet,
			fmt.Sprintf("http://localhost/api/add-cabbages/%s?some_int=1&some_string=\"hello+world\"&some_int_array=[1,2,3]&some_string_array=[\"hello+world\",+\"goodbye+world\"]&some_bool_array=[true,false]", patchID),
			bytes.NewReader(rawRequestBody),
		)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, r)

		b, err := io.ReadAll(w.Body)
		require.NoError(t, err)
		require.Equal(t, customHTTPHandler.Status, w.Result().StatusCode, string(b))
		log.Printf("b: %s", string(b))
	})

	t.Run("GetCustomHTTPHandlerEmptyRequestAndEmptyResponse", func(t *testing.T) {
		customHTTPHandler, err := GetHTTPHandler(
			http.MethodGet,
			"/api/add-cabbages/{patch_id}",
			http.StatusOK,
			func(ctx context.Context, pathParams SomePathParams, queryParams SomeQueryParams, req EmptyRequest, rawReq any) (*EmptyResponse, error) {
				log.Printf("pathParams: %s\n\n", hack.UnsafeJSONPrettyFormat(pathParams))
				log.Printf("queryParams: %s\n\n", hack.UnsafeJSONPrettyFormat(queryParams))

				return nil, nil
			},
		)
		require.NoError(t, err)

		router := chi.NewRouter()
		router.Get(customHTTPHandler.FullPath, customHTTPHandler.ServeHTTP)

		patchID := uuid.Must(uuid.NewRandom())

		r := httptest.NewRequest(
			http.MethodGet,
			fmt.Sprintf("http://localhost/api/add-cabbages/%s?some_int=1&some_string=\"hello+world\"&some_int_array=[1,2,3]&some_string_array=[\"hello+world\",+\"goodbye+world\"]&some_bool_array=[true,false]", patchID),
			nil,
		)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, r)

		b, err := io.ReadAll(w.Body)
		require.NoError(t, err)
		require.Equal(t, customHTTPHandler.Status, w.Result().StatusCode, string(b))
		log.Printf("b: %s", string(b))
	})

	t.Run("GetCustomHTTPHandlerMapForDynamicQueryParams", func(t *testing.T) {
		customHTTPHandler, err := GetHTTPHandler(
			http.MethodGet,
			"/api/add-cabbages/{patch_id}",
			http.StatusOK,
			func(ctx context.Context, pathParams SomePathParams, queryParams map[string]any, req SomeRequest, rawReq any) (*SomeResponse, error) {
				log.Printf("pathParams: %s\n\n", hack.UnsafeJSONPrettyFormat(pathParams))
				log.Printf("queryParams: %s\n\n", hack.UnsafeJSONPrettyFormat(queryParams))
				log.Printf("req (%s): %s\n\n", reflect.TypeOf(req).String(), hack.UnsafeJSONPrettyFormat(req))
				log.Printf("rawReq (%s): %s\n\n", reflect.TypeOf(rawReq).String(), hack.UnsafeJSONPrettyFormat(rawReq))

				res := &SomeResponse{
					Timestamp: time.Now(),
					Cabbages: []Cabbage{{
						What:  map[string]bool{"hello": true, "world": true},
						Is:    []map[int]time.Time{{1: time.Now()}},
						This:  69,
						Thing: []float64{42.0},
						Buddy: nil,
					},
					},
					FavouriteCabbage: &Cabbage{
						What:  map[string]bool{"goodbye": true, "world": true},
						Is:    []map[int]time.Time{{2: time.Now()}},
						This:  1337,
						Thing: []float64{8008135.0},
						Buddy: helpers.Ptr(42),
					},
				}

				log.Printf("res: %s\n\n", hack.UnsafeJSONPrettyFormat(res))

				return res, nil
			},
		)
		require.NoError(t, err)

		router := chi.NewRouter()
		router.Get(customHTTPHandler.FullPath, customHTTPHandler.ServeHTTP)

		patchID := uuid.Must(uuid.NewRandom())

		requestBody := SomeRequest{
			Timestamp: time.Now(),
			Cabbages: []Cabbage{
				{
					What:  map[string]bool{"hello": true, "world": true},
					Is:    []map[int]time.Time{{1: time.Now()}},
					This:  69,
					Thing: []float64{42.0},
					Buddy: nil,
				},
			},
		}

		rawRequestBody, err := json.Marshal(requestBody)
		require.NoError(t, err)

		r := httptest.NewRequest(
			http.MethodGet,
			fmt.Sprintf("http://localhost/api/add-cabbages/%s?some_int=1&some_string=\"hello+world\"&some_int_array=[1,2,3]&some_string_array=[\"hello+world\",+\"goodbye+world\"]&some_bool_array=[true,false]", patchID),
			bytes.NewReader(rawRequestBody),
		)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, r)

		b, err := io.ReadAll(w.Body)
		require.NoError(t, err)
		require.Equal(t, customHTTPHandler.Status, w.Result().StatusCode, string(b))
		log.Printf("b: %s", string(b))
	})
}
