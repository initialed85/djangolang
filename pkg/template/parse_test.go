package template

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	parseTasks, err := Parse()
	require.NoError(t, err)
	require.NotNil(t, parseTasks)
	require.Len(t, parseTasks, 9)

	for _, parseTask := range parseTasks {
		require.True(
			t,
			parseTask.StartMatch != parseTask.ReplacedStartMatch ||
				parseTask.KeepMatch != parseTask.ReplacedKeepMatch ||
				parseTask.EndMatch != parseTask.ReplacedEndMatch,
			fmt.Sprintf("%v had no changes", parseTask.Name),
		)

		require.Greater(
			t,
			len(parseTask.StartVariableNameSet)+
				len(parseTask.KeepVariableNameSet)+
				len(parseTask.EndVariableNameSet),
			0,
			fmt.Sprintf("%v had no variables", parseTask.Name),
		)
	}
}
