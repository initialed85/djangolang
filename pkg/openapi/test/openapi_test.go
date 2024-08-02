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
		// model_reference.PhysicalThing{},
		// model_reference.LogicalThing{},
		model_generated.Fuzz{},
	})
	require.NoError(t, err)

	fmt.Printf("\n\n%v\n\n", o.String())
}
