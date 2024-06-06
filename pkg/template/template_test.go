package template

import (
	"context"
	"os"
	"path"
	"runtime"
	"testing"

	"github.com/initialed85/djangolang/pkg/helpers"
	"github.com/initialed85/djangolang/pkg/introspect"
	"github.com/stretchr/testify/require"
)

func TestTemplate(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := helpers.GetDBFromEnvironment(ctx)
	if err != nil {
		require.NoError(t, err)
	}
	defer func() {
		_ = db.Close()
	}()

	schema := helpers.GetSchema()

	tableByName, err := introspect.Introspect(ctx, db, schema)
	require.NoError(t, err)

	parseTasks, err := Parse()
	require.NoError(t, err)
	require.NotNil(t, parseTasks)

	templateDataByFileName, err := Template(tableByName, "model_generated")
	require.NoError(t, err)
	require.NotNil(t, templateDataByFileName)
	require.Len(t, templateDataByFileName, 4)

	_, filePath, _, ok := runtime.Caller(0)
	require.True(t, ok)

	dirPath, _ := path.Split(filePath)
	dirPath = path.Join(dirPath, "../", "model_generated")
	_ = os.RemoveAll(dirPath)
	err = os.MkdirAll(dirPath, 0o777)
	require.NoError(t, err)

	for fileName, templateData := range templateDataByFileName {
		_ = os.WriteFile(path.Join(dirPath, fileName), []byte(templateData), 0o777)
	}
}
