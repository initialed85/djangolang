package introspect

import (
	"context"
	"testing"

	"github.com/initialed85/djangolang/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestIntrospect(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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

	primaryKeyCount := 0
	for _, column := range tableByName["no_primary_things"].Columns {
		if column.IsPrimaryKey {
			primaryKeyCount++
		}
	}
	require.Equal(t, 0, primaryKeyCount)
	require.Nil(t, tableByName["no_primary_things"].PrimaryKeyColumn)

	primaryKeyCount = 0
	for _, column := range tableByName["two_primary_things"].Columns {
		if column.IsPrimaryKey {
			primaryKeyCount++
		}
	}
	require.Equal(t, 2, primaryKeyCount)
	require.NotNil(t, tableByName["two_primary_things"].PrimaryKeyColumns)
}
