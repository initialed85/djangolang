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
	// Verify UpdateFields in reference has correct structure: values [colValA, colValB, ..., m.ID]
	fileData := model_reference.ReferenceFileData

	// Find the UpdateFields method body for LogicalThing
	updateFieldsIdx := strings.Index(fileData, "func (m *LogicalThing) UpdateFields")
	require.NotEqual(t, -1, updateFieldsIdx, "UpdateFields method should exist in LogicalThing")

	// Extract method body (simplified check)
	updateFieldsBody := fileData[updateFieldsIdx:]
	
	// Verify the method uses map iteration and single m.ID append
	require.True(t, strings.Contains(updateFieldsBody, "range fields"), "UpdateFields should range over fields")
	
	// Check that values are built properly
	require.True(t, strings.Contains(updateFieldsBody, "values = append(values, m.ID)"), "UpdateFields should append m.ID to values")
}

func TestFieldUpdateCaseGeneration(t *testing.T) {
	// Verify that case statements use concrete identifiers (not template syntax)
	fileData := model_reference.ReferenceFileData

	// Check that LogicalThing case statements exist with concrete table names
	require.True(t, strings.Contains(fileData, "LogicalThingTableIDColumn"), "Should have concrete table column identifiers")
	require.True(t, strings.Contains(fileData, "LogicalThingTableCreatedAtColumn"), "Should have concrete time column identifiers")
}
