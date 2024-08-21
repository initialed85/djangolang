package query

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHandlePath(t *testing.T) {
	t.Run("SimpleWithMaxVisitCountOfDefault1", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		ctx = context.WithValue(ctx, DepthKey, DepthValue{MaxDepth: 0})

		tableNames := []string{
			"logical_things",
			"logical_things",
			"logical_things",
		}

		var depth int
		var tableName string
		var ok bool

		for depth, tableName = range tableNames {
			ctx, ok = HandleQueryPathGraphCycles(ctx, tableName)
			if !ok {
				break
			}
		}

		require.Equal(t, "logical_things", tableName)
		require.Equal(t, 2, depth)
	})

	t.Run("SimpleWithMaxVisitCountOf2", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		ctx = context.WithValue(ctx, DepthKey, DepthValue{MaxDepth: 0})

		tableNames := []string{
			"logical_things",
			"logical_things",
			"logical_things",
		}

		var depth int
		var tableName string
		var ok bool

		for depth, tableName = range tableNames {
			ctx, ok = HandleQueryPathGraphCycles(ctx, tableName, 2)
			if !ok {
				break
			}
		}

		require.Equal(t, "logical_things", tableName)
		require.Equal(t, 2, depth)
	})

	t.Run("SimpleWithMaxVisitCountOf0", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		ctx = context.WithValue(ctx, DepthKey, DepthValue{MaxDepth: 0})

		tableNames := []string{
			"logical_things",
			"logical_things",
			"logical_things",
		}

		var depth int
		var tableName string
		var ok bool

		for depth, tableName = range tableNames {
			ctx, ok = HandleQueryPathGraphCycles(ctx, tableName, 0)
			if !ok {
				break
			}
		}

		require.Equal(t, "logical_things", tableName)
		require.Equal(t, 0, depth)
	})

	t.Run("ComplexWithMaxVisitCountOfDefault1", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		ctx = context.WithValue(ctx, DepthKey, DepthValue{MaxDepth: 0})

		tableNames := []string{
			"logical_things",
			"physical_things",
			"location_history",
			"physical_things",
			"logical_things",
			"physical_things",
			"location_history",
			"physical_things",
		}

		var depth int
		var tableName string
		var ok bool

		for depth, tableName = range tableNames {
			ctx, ok = HandleQueryPathGraphCycles(ctx, tableName)
			if !ok {
				break
			}
		}

		require.Equal(t, "physical_things", tableName)
		require.Equal(t, 5, depth)
	})

	t.Run("ComplexWithMaxVisitCountOf2", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		ctx = context.WithValue(ctx, DepthKey, DepthValue{MaxDepth: 0})

		tableNames := []string{
			"logical_things",
			"physical_things",
			"location_history",
			"physical_things",
			"logical_things",
			"physical_things",
			"location_history",
			"physical_things",
		}

		var depth int
		var tableName string
		var ok bool

		for depth, tableName = range tableNames {
			ctx, ok = HandleQueryPathGraphCycles(ctx, tableName, 2)
			if !ok {
				break
			}
		}

		require.Equal(t, "physical_things", tableName)
		require.Equal(t, 7, depth)
	})
}
