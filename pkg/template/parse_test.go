package template

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	parseTasks, err := Parse()
	require.NoError(t, err)
	require.NotNil(t, parseTasks)

	// TODO: handy for debugging
	// for _, parseTask := range parseTasks {
	// 	t.Logf("Name: %#+v", parseTask.Name)
	// 	t.Logf("StartExpr: %#+v", parseTask.StartExpr)
	// 	t.Logf("KeepExpr: %#+v", parseTask.KeepExpr)
	// 	t.Logf("EndExpr: %#+v", parseTask.EndExpr)
	// 	t.Logf("TokenizeTasks: %#+v", parseTask.TokenizeTasks)
	// 	t.Logf("KeepIsPerColumn: %#+v", parseTask.KeepIsPerColumn)
	// 	t.Logf("KeepIsForPrimaryKeyOnly: %#+v", parseTask.KeepIsForPrimaryKeyOnly)
	// 	t.Logf("KeepIsForNonPrimaryKeyOnly: %#+v", parseTask.KeepIsForNonPrimaryKeyOnly)
	// 	t.Logf("KeepIsForForeignKeysOnly: %#+v", parseTask.KeepIsForForeignKeysOnly)
	// 	t.Logf("ReplaceText: %#+v", parseTask.ReplaceText)
	// 	t.Logf("StartMatch: %#+v", parseTask.StartMatch)
	// 	t.Logf("KeepMatch: %#+v", parseTask.KeepMatch)
	// 	t.Logf("EndMatch: %#+v", parseTask.EndMatch)
	// 	t.Logf("Fragment: %#+v", parseTask.Fragment)
	// 	t.Logf("ReplacedStartMatch: %#+v", parseTask.ReplacedStartMatch)
	// 	t.Logf("ReplacedKeepMatch: %#+v", parseTask.ReplacedKeepMatch)
	// 	t.Logf("ReplacedEndMatch: %#+v", parseTask.ReplacedEndMatch)
	// 	t.Logf("ReplacedFragment: %#+v", parseTask.ReplacedFragment)
	// 	t.Logf("StartVariableNameSet: %#+v", parseTask.StartVariableNameSet)
	// 	t.Logf("KeepVariableNameSet: %#+v", parseTask.KeepVariableNameSet)
	// 	t.Logf("EndVariableNameSet: %#+v", parseTask.EndVariableNameSet)
	// 	t.Logf("")
	// }
}
