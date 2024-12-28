package schema

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"
	"runtime"
	"testing"

	_ "embed"

	"github.com/initialed85/djangolang/internal/hack"
	"github.com/initialed85/djangolang/pkg/config"
	"github.com/initialed85/djangolang/pkg/introspect"
	"github.com/initialed85/djangolang/pkg/template"
	"github.com/stretchr/testify/require"
)

//go:embed test.yaml
var testYAML []byte

func TestSchema(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := config.GetDBFromEnvironment(ctx)
	if err != nil {
		require.NoError(t, err)
	}
	defer func() {
		db.Close()
	}()

	tx, err := db.Begin(ctx)
	require.NoError(t, err)

	defer func() {
		_ = tx.Rollback(ctx)
	}()

	schema, err := Parse(testYAML)
	require.NoError(t, err)
	log.Printf("schema: %s", hack.UnsafeJSONPrettyFormat(schema))

	fmt.Printf("\n\n")

	sql, err := Dump(schema, "test", true)
	require.NoError(t, err)
	log.Printf("sql: %s", string(sql))

	_ = os.WriteFile("/tmp/dump.sql", []byte(sql), 0o777)

	_, err = tx.Exec(ctx, string(sql))
	require.NoError(t, err)

	tableByName, err := introspect.Introspect(ctx, tx, "test")
	require.NoError(t, err)
	log.Printf("tableByName: %s", hack.UnsafeJSONPrettyFormat(tableByName))

	templateDataByFileName, err := template.Template(tableByName, "github.com/initialed85/djangolang", "model_generated_from_schema", "test")
	require.NoError(t, err)
	require.NotNil(t, templateDataByFileName)
	require.Len(t, templateDataByFileName, 13)

	_, filePath, _, ok := runtime.Caller(0)
	require.True(t, ok)

	dirPath, _ := path.Split(filePath)
	dirPath = path.Join(dirPath, "../", "model_generated_from_schema")
	_ = os.RemoveAll(dirPath)
	err = os.MkdirAll(dirPath, 0o777)
	require.NoError(t, err)
	err = os.MkdirAll(path.Join(dirPath, "cmd"), 0o777)
	require.NoError(t, err)

	for fileName, templateData := range templateDataByFileName {
		log.Printf("schema writing %s", path.Join(dirPath, fileName))
		err = os.WriteFile(path.Join(dirPath, fileName), []byte(templateData), 0o777)
		require.NoError(t, err)
	}

	err = tx.Commit(ctx)
	require.NoError(t, err)
}
