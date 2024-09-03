package openapi

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/initialed85/djangolang/pkg/model_generated"
	"github.com/initialed85/djangolang/pkg/openapi"
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

func TestOpenAPI(t *testing.T) {
	o, err := openapi.NewFromIntrospectedSchema(
		[]any{
			model_generated.NotNullFuzz{},
			model_generated.PhysicalThing{},
			model_generated.LogicalThing{},
			model_generated.LocationHistory{},
		},
		[]openapi.CustomHTTPHandlerSummary{
			{
				PathParams:  SomePathParams{},
				QueryParams: SomeQueryParams{},
				Request:     SomeRequest{},
				Response:    SomeResponse{},
				Method:      http.MethodPost,
				Path:        "/api/add-cabbages/{patch_id}",
				Status:      http.StatusCreated,
			},
		},
	)
	require.NoError(t, err)

	fmt.Printf("\n\n%v\n\n", o.String())

	// require.Fail(t, "failing at user request")
}
