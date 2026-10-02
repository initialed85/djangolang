package query

import (
	"context"
	"fmt"
	"testing"

	"github.com/initialed85/djangolang/pkg/helpers"
	"github.com/stretchr/testify/require"
)

func TestHandlePath(t *testing.T) {
	t.Run("SimpleWithDefaultMaxDepth", func(t *testing.T) {
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

	t.Run("ComplexWithDefaultMaxDepth", func(t *testing.T) {
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

	t.Run("ObjectIdentityAndCollectionMarkers", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// Object-specific paths remain distinct, while collection selectors use
		// {<nil>} as a cycle marker.
		ctx = WithMaxDepth(ctx, helpers.Ptr(0))

		var ok bool
		ctx, ok = HandleQueryPathGraphCycles(ctx, "meme{<nil>}", true)
		require.True(t, ok)
		ctx, ok = HandleQueryPathGraphCycles(ctx, "__ReferencedBy__meme_tag{root-id}", true)
		require.True(t, ok)
		ctx, ok = HandleQueryPathGraphCycles(ctx, "meme_tag{<nil>}", false)
		require.True(t, ok)

		// Both sibling tag objects load, because the object identities differ
		// and the non-incrementing collection selector is not retained.
		ctx, ok = HandleQueryPathGraphCycles(ctx, "tag{tag-a}", true)
		require.True(t, ok)
		ctx, ok = HandleQueryPathGraphCycles(ctx, "tag{<nil>}", false)
		require.True(t, ok)
		ctx, ok = HandleQueryPathGraphCycles(ctx, "tag{tag-b}", true)
		require.True(t, ok)
		ctx, ok = HandleQueryPathGraphCycles(ctx, "tag{<nil>}", false)
		require.True(t, ok)

		// An exact object identity remains a cycle, while the root meme object
		// is stopped by its repeated {<nil>} collection marker.
		ctx, ok = HandleQueryPathGraphCycles(ctx, "tag{tag-a}", true)
		require.False(t, ok)
		ctx, ok = HandleQueryPathGraphCycles(ctx, "meme{root-id}", true)
		require.True(t, ok)
		ctx, ok = HandleQueryPathGraphCycles(ctx, "meme{<nil>}", false)
		require.False(t, ok)
	})

	t.Run("ComplexWithUnlimitedDepth", func(t *testing.T) {
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
