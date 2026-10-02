package template

import (
	"strings"
	"testing"

	"github.com/initialed85/djangolang/pkg/model_reference"
	"github.com/stretchr/testify/require"
)

func TestFieldUpdateInReference(t *testing.T) {
	// Verify FieldUpdate methods exist in reference files
	fileData := model_reference.ReferenceFileData

	require.True(t, strings.Contains(fileData, "UpdateField"), "Reference should contain UpdateField")
	require.True(t, strings.Contains(fileData, "UpdateFields"), "Reference should contain UpdateFields")
}

func TestFieldUpdateMethodStructure(t *testing.T) {
	// Verify UpdateFields in reference has correct structure: values [colValA, colValB, ..., primary key]
	fileData := model_reference.ReferenceFileData

	// Find the UpdateFields method body for LogicalThing
	updateFieldsIdx := strings.Index(fileData, "func (m *LogicalThing) UpdateFields")
	require.NotEqual(t, -1, updateFieldsIdx, "UpdateFields method should exist in LogicalThing")

	// Extract method body (simplified check)
	updateFieldsBody := fileData[updateFieldsIdx:]

	// Verify the method uses map iteration and single m.ID append
	require.True(t, strings.Contains(updateFieldsBody, "range fields"), "UpdateFields should range over fields")

	// Check that values are built properly. The reference model resolves its
	// primary key through the model interface so generated models can use a
	// non-ID primary key as well.
	require.True(t, strings.Contains(updateFieldsBody, "values = append(values, m.GetPrimaryKeyValue())"), "UpdateFields should append the primary key value to values")
}

func TestFieldUpdateCaseGeneration(t *testing.T) {
	// Verify that case statements use concrete identifiers (not template syntax)
	fileData := model_reference.ReferenceFileData

	// Check that LogicalThing case statements exist with concrete table names
	require.True(t, strings.Contains(fileData, "LogicalThingTableIDColumn"), "Should have concrete table column identifiers")
	require.True(t, strings.Contains(fileData, "LogicalThingTableCreatedAtColumn"), "Should have concrete time column identifiers")
}
