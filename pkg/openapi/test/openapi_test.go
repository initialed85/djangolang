package openapi

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/initialed85/djangolang/pkg/config"
	"github.com/initialed85/djangolang/pkg/helpers"
	"github.com/initialed85/djangolang/pkg/model_generated"
	"github.com/initialed85/djangolang/pkg/openapi"
	"github.com/initialed85/djangolang/pkg/query"
	"github.com/initialed85/djangolang/pkg/server"
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
	Timestamp        time.Time `json:"timestamp"`
	Cabbages         []Cabbage
	TakeThisWhatever any `json:"take_this_whatever"`
}

type SomeResponse struct {
	Timestamp        time.Time `json:"timestamp"`
	Cabbages         []Cabbage `json:"cabbages"`
	FavouriteCabbage *Cabbage  `json:"favourite_cabbage"`
	HereIsAWhatever  any       `json:"here_is_a_whatever"`
}

type ClaimRequest struct {
	ClaimDurationSeconds float64 `json:"claim_duration_seconds"`
}

func TestOpenAPI(t *testing.T) {
	t.Run("Contrived", func(t *testing.T) {
		o, err := openapi.NewFromIntrospectedSchema(
			[]server.HTTPHandlerSummary{
				{
					PathParams:  SomePathParams{},
					QueryParams: SomeQueryParams{},
					Request:     SomeRequest{},
					Response:    SomeResponse{},
					Method:      http.MethodPost,
					Path:        "/add-cabbages/{patch_id}",
					Status:      http.StatusCreated,
				},
				{
					PathParams:  SomePathParams{},
					QueryParams: SomeQueryParams{},
					Request:     SomeRequest{},
					Response:    &SomeResponse{},
					Method:      http.MethodPost,
					Path:        "/do-something-optional-with-cabbages/{patch_id}",
					Status:      http.StatusOK,
				},
				{
					PathParams:  SomePathParams{},
					QueryParams: server.EmptyQueryParams{},
					Request:     server.EmptyRequest{},
					Response:    server.EmptyResponse{},
					Method:      http.MethodDelete,
					Path:        "/remove-cabbages/{patch_id}",
					Status:      http.StatusNoContent,
				},
				{
					PathParams:  server.EmptyPathParams{},
					QueryParams: server.EmptyQueryParams{},
					Request:     server.EmptyRequest{},
					Response:    server.EmptyResponse{},
					Method:      http.MethodGet,
					Path:        "/do-nothing",
					Status:      http.StatusOK,
				},
				{
					PathParams:  server.EmptyPathParams{},
					QueryParams: server.EmptyQueryParams{},
					Request:     server.EmptyRequest{},
					Response:    server.EmptyResponse{},
					Method:      http.MethodPost,
					Path:        "/do-nothing",
					Status:      http.StatusOK,
				},
			},
		)
		require.NoError(t, err)

		fmt.Printf("\n\n%v\n\n", o.String())
	})

	t.Run("Camry", func(t *testing.T) {
		claimVideoForObjectDetectorHandler, err := model_generated.GetHTTPHandler(
			http.MethodPatch,
			"/custom/claim-video-for-object-detector",
			http.StatusOK,
			func(
				ctx context.Context,
				pathParams server.EmptyPathParams,
				queryParams server.EmptyQueryParams,
				req ClaimRequest,
				rawReq any,
			) (*model_generated.Video, error) {
				now := time.Now().UTC()

				claimUntil := now.Add(time.Second * time.Duration(req.ClaimDurationSeconds))

				if claimUntil.Sub(now) <= 0 {
					return nil, fmt.Errorf("claim_duration_seconds too short; must result in a claim that expires in the future")
				}

				video := &model_generated.Video{}

				return video, nil
			},
		)
		require.NoError(t, err)

		o, err := openapi.NewFromIntrospectedSchema(
			[]server.HTTPHandlerSummary{
				claimVideoForObjectDetectorHandler.Summarize(),
			},
		)
		require.NoError(t, err)

		fmt.Printf("\n\n%v\n\n", o.String())
	})

	t.Run("ModelGenerated", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		ctx = query.WithMaxDepth(ctx, helpers.Ptr(0))

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

		_ = model_generated.GetRouter(db, redisPool, nil, nil, nil)

		httpHandlerSummaries, err := model_generated.GetHTTPHandlerSummaries()
		require.NoError(t, err)

		o, err := openapi.NewFromIntrospectedSchema(httpHandlerSummaries)
		require.NoError(t, err)
		require.NotNil(t, o)

		fmt.Printf("\n\n%v\n\n", o.String())
	})
}
