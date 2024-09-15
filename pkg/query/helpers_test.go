package query

import (
	"context"
	"fmt"
	"testing"

	"github.com/initialed85/djangolang/pkg/helpers"
	"github.com/stretchr/testify/require"
)

func TestHandlePath(t *testing.T) {
	t.Run("SimpleWithMaxVisitCountOfDefault", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		var ok bool

		tableNames := []string{
			"logical_things",
			"logical_things",
			"logical_things",
		}

		path := make([]string, 0)

		for i, tableName := range tableNames {
			ctx, ok = HandleQueryPathGraphCycles(ctx, tableName, true)
			if !ok {
				path = append(path, fmt.Sprintf("%d: %s = no", i, tableName))
			} else {
				path = append(path, fmt.Sprintf("%d: %s = yes", i, tableName))
			}
		}

		require.Equal(
			t,
			[]string([]string{
				"0: logical_things = yes",
				"1: logical_things = no",
				"2: logical_things = no",
			}),
			path,
		)
	})

	t.Run("ComplexWithMaxVisitCountOfDefault", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		var ok bool

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

		path := make([]string, 0)

		for i, tableName := range tableNames {
			ctx, ok = HandleQueryPathGraphCycles(ctx, tableName, true)
			if !ok {
				path = append(path, fmt.Sprintf("%d: %s = no", i, tableName))
			} else {
				path = append(path, fmt.Sprintf("%d: %s = yes", i, tableName))
			}
		}

		require.Equal(
			t,
			[]string{
				"0: logical_things = yes",
				"1: physical_things = no",
				"2: location_history = no",
				"3: physical_things = no",
				"4: logical_things = no",
				"5: physical_things = no",
				"6: location_history = no",
				"7: physical_things = no",
			},
			path,
		)
	})

	t.Run("ComplexWithMaxVisitCountOfSmart", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		ctx = WithMaxDepth(ctx, helpers.Ptr(0))

		var ok bool

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

		path := make([]string, 0)

		for i, tableName := range tableNames {
			ctx, ok = HandleQueryPathGraphCycles(ctx, tableName, true)
			if !ok {
				path = append(path, fmt.Sprintf("%d: %s = no", i, tableName))
			} else {
				path = append(path, fmt.Sprintf("%d: %s = yes", i, tableName))
			}
		}

		require.Equal(
			t,
			[]string{
				"0: logical_things = yes",
				"1: physical_things = yes",
				"2: location_history = yes",
				"3: physical_things = no",
				"4: logical_things = no",
				"5: physical_things = no",
				"6: location_history = no",
				"7: physical_things = no",
			},
			path,
		)
	})
}
