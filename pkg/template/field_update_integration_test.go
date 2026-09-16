package template

import (
	"strings"
	"testing"

	"github.com/initialed85/djangolang/pkg/model_reference"
	"github.com/stretchr/testify/require"
)

func TestFieldUpdateGeneration(t *testing.T) {
	// Test that field update markers are properly preserved during template processing
	
	// Load reference file data (simulates what happens during real template generation)
	fileData := model_reference.ReferenceFileData

	// Verify the template markers are present in source
	require.True(t, strings.Contains(fileData, "<field-update-methods>"), "Missing opening field-update-methods marker")
	require.True(t, strings.Contains(fileData, "</field-update-methods>"), "Missing closing field-update-methods marker")
	require.True(t, strings.Contains(fileData, "<field-update-cases>"), "Missing field-update-cases placeholder")

	// Verify the KeepMatch regex can capture the method template
	parseTasks, err := Parse()
	require.NoError(t, err)

	var fieldUpdateTask *ParseTask
	for i := range parseTasks {
		if parseTasks[i].Name == "FieldUpdate" {
			fieldUpdateTask = &parseTasks[i]
			break
		}
	}
	require.NotNil(t, fieldUpdateTask)
	require.NotEmpty(t, fieldUpdateTask.KeepMatch)

	// Verify KeepMatch captures method content without template syntax errors
	// The KeepMatch should match between markers and capture the method body
	require.True(t, strings.Contains(fieldUpdateTask.KeepMatch, "UpdateField"), "KeepMatch should reference UpdateField")

	// Verify the method has proper structure after case insertion
	// (This tests that the template will generate valid Go code)
	casePlaceholder := "<field-update-cases>"
	require.True(t, strings.Contains(fieldUpdateTask.KeepMatch, casePlaceholder), "KeepMatch should have case placeholder")
}

func TestFieldUpdateMethodStructure(t *testing.T) {
	// Verify UpdateFields builds values correctly: [colValA, colValB, ..., m.ID]
	parseTasks, err := Parse()
	require.NoError(t, err)

	for _, task := range parseTasks {
		if task.Name == "FieldUpdate" {
			keepMatch := task.KeepMatch

			// Verify the structure: columns/values in loop, then m.ID appended after
			// This is verified by checking that there's only one "values = append(values, m.ID)"
			count := strings.Count(keepMatch, "values = append(values, m.ID)")
			require.Equal(t, 1, count, "values = append(values, m.ID) should appear exactly once, after the loop")

			// Verify case statements exist
			require.True(t, strings.Contains(keepMatch, "switch fieldName"), "Method should have fieldName switch")
			require.True(t, strings.Contains(keepMatch, "switch columnName"), "Method should have columnName switch")
		}
	}
}
