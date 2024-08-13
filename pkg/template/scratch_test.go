package template

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

type TestThing struct {
	Name string
}

func ScratchTest(t *testing.T) {
	t.Run("StructEquality", func(t *testing.T) {
		testThing1 := TestThing{Name: "TestThing1"}
		testThing2 := TestThing{Name: "TestThing2"}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		ctx = context.WithValue(ctx, testThing1, testThing1)
		ctx = context.WithValue(ctx, testThing2, testThing2)

		require.NotEqual(t, testThing1, testThing2)

		retrievedTestThing1 := ctx.Value(testThing1)
		require.Equal(t, testThing1, retrievedTestThing1)

		retrievedTestThing2 := ctx.Value(testThing2)
		require.Equal(t, testThing2, retrievedTestThing2)

		require.NotEqual(t, retrievedTestThing1, retrievedTestThing2)
	})
}
