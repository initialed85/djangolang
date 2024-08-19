package openapi

import (
	"fmt"
	"testing"

	"github.com/initialed85/djangolang/pkg/model_generated"
	"github.com/initialed85/djangolang/pkg/openapi"
	"github.com/stretchr/testify/require"
)

func TestOpenAPI(t *testing.T) {
	o, err := openapi.NewFromIntrospectedSchema([]any{
		model_generated.NotNullFuzz{},
		model_generated.PhysicalThing{},
		model_generated.LogicalThing{},
		model_generated.LocationHistory{},
	})
	require.NoError(t, err)

	fmt.Printf("\n\n%v\n\n", o.String())
}
