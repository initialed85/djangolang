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
