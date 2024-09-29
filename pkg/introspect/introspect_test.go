package introspect

import (
	"context"
	"testing"

	"github.com/initialed85/djangolang/pkg/config"
	"github.com/initialed85/djangolang/pkg/helpers"
	"github.com/initialed85/djangolang/pkg/query"
	"github.com/stretchr/testify/require"
)

func TestIntrospect(t *testing.T) {
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

	schema := config.GetSchema()

	tx, err := db.Begin(ctx)
	if err != nil {
		require.NoError(t, err)
	}

	defer func() {
		_ = tx.Rollback(ctx)
	}()

	originalTableByName, err := Introspect(ctx, tx, schema)
	require.NoError(t, err)
	require.NotNil(t, originalTableByName["logical_things"].ForeignTables)
	require.NotNil(t, originalTableByName["logical_things"].ForeignTables[0])

	tableByNameAsJSON, err := json.MarshalIndent(originalTableByName, "", "  ")
	require.NoError(t, err)

	var tableByName TableByName
	err = json.Unmarshal(tableByNameAsJSON, &tableByName)
	require.NoError(t, err)
	require.NotNil(t, tableByName["logical_things"].ForeignTables)
	require.NotNil(t, tableByName["logical_things"].ForeignTables[0])
}
