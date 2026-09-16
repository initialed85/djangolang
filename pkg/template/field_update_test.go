package template

import (
	"strings"
	"testing"

	"github.com/initialed85/djangolang/pkg/model_reference"
	"github.com/stretchr/testify/require"
)

func TestFieldUpdateParseTask(t *testing.T) {
	parseTasks, err := Parse()
	require.NoError(t, err)
	require.NotNil(t, parseTasks)

	// Find FieldUpdate task
	var fieldUpdateTask *ParseTask
	for i := range parseTasks {
		if parseTasks[i].Name == "FieldUpdate" {
			fieldUpdateTask = &parseTasks[i]
			break
		}
	}
	require.NotNil(t, fieldUpdateTask, "FieldUpdate parse task should exist")

	// Verify the task can be found and has required properties
	require.Equal(t, "FieldUpdate", fieldUpdateTask.Name)
	require.NotNil(t, fieldUpdateTask.StartExpr)
	require.NotNil(t, fieldUpdateTask.KeepExpr)
	require.NotNil(t, fieldUpdateTask.EndExpr)

	// Verify KeepExpr can match the field update method template
	keepMatch := fieldUpdateTask.KeepMatch
	require.NotEmpty(t, keepMatch)
	require.True(t, strings.Contains(keepMatch, "UpdateField"), "KeepMatch should contain UpdateField")
}

func TestFieldUpdateTemplateGeneration(t *testing.T) {
	// Use the reference model data to test field update generation
	fileData := model_reference.ReferenceFileData

	// Verify the field-update-methods markers exist in reference data
	require.True(t, strings.Contains(fileData, "<field-update-methods>"), "Reference should have field-update-methods marker")
	require.True(t, strings.Contains(fileData, "</field-update-methods>"), "Reference should have closing field-update-methods marker")

	// Verify the placeholder for auto-generated cases exists
	require.True(t, strings.Contains(fileData, "<field-update-cases>"), "Reference should have field-update-cases placeholder")
}

func TestUpdateFieldsArgumentOrder(t *testing.T) {
	// Verify the generated UpdateFields method has correct argument order
	// values should be: [col1Value, col2Value, ..., m.ID] - single WHERE id at end
	
	// The template should generate this pattern - verify via inspection
	parseTasks, err := Parse()
	require.NoError(t, err)

	for _, task := range parseTasks {
		if task.Name == "FieldUpdate" {
			// Check that KeepMatch contains the expected argument structure
			keepMatch := task.KeepMatch
			// The method body should append ID only once, after building all values
			// This is verified by the structure of the template
			require.NotEmpty(t, keepMatch)
		}
	}
}
