package template

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	parseTasks, err := Parse()
	require.NoError(t, err)
	require.NotNil(t, parseTasks)
	require.NotEmpty(t, parseTasks)
}
