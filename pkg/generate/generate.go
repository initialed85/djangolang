package generate

import (
	"context"
	"os"
	"path/filepath"
	"time"

	"github.com/initialed85/djangolang/internal/helpers"
	"github.com/initialed85/djangolang/pkg/introspect"
)

func Run(ctx context.Context) error {
	relativeOutputPath := "./pkg/model"
	if len(os.Args) < 3 {
		relativeOutputPath = os.Args[2]
	}

	outputPath, err := filepath.Abs(relativeOutputPath)
	if err != nil {
		return err
	}

	// _ = os.RemoveAll(outputPath)
	// err = os.MkdirAll(outputPath, 0o600)
	// if err != nil {
	// 	return err
	// }

	conn, err := helpers.GetDBConnFromEnvironment(ctx)
	if err != nil {
		return err
	}
	defer func() {
		closeCtx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		_ = conn.Close(closeCtx)
	}()

	schema := helpers.GetSchema()

	tableByName, err := introspect.Introspect(ctx, conn, schema)
	if err != nil {
		return err
	}

	for _, table := range tableByName {
		for _, column := range table.ColumnByName {
			_ = column
		}
	}

	return nil
}
