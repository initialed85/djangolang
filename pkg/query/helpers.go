package query

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/initialed85/djangolang/pkg/config"
	"github.com/initialed85/djangolang/pkg/helpers"
	"github.com/initialed85/djangolang/pkg/introspect"
)

func GetCurrentDepthValue(ctx context.Context) *DepthValue {
	rawDepthValue := ctx.Value(DepthKey)
	if rawDepthValue == nil {
		return nil
	}

	depthValue, ok := rawDepthValue.(DepthValue)
	if !ok {
		log.Panicf("assertion failed: expected context key %v to contain a DepthValue but it had %#+v", DepthKey, rawDepthValue)
	}

	if depthValue.ID == "" {
		log.Panicf("assertion failed: looks like we've got an unitialized DepthValue; this should never happen")
	}

	return &depthValue
}

func GetCurrentPathValue(ctx context.Context) *PathValue {
	rawPathValue := ctx.Value(PathKey)
	if rawPathValue == nil {
		return nil
	}

	pathValue, ok := rawPathValue.(PathValue)
	if !ok {
		log.Panicf("assertion failed: expected context key %v to contain a PathValue but it had %#+v", PathKey, rawPathValue)
	}

	if pathValue.ID == "" {
		log.Panicf("assertion failed: looks like we've got an unitialized PathValue; this should never happen")
	}

	return &pathValue
}

func WithMaxDepth(ctx context.Context, maxDepth *int, increments ...bool) context.Context {
	if config.Debug() {
		log.Printf("entered WithMaxDepth; %#+v", ctx.Value(DepthKey))

		defer func() {
			log.Printf("exited WithMaxDepth; %#+v", ctx.Value(DepthKey))
		}()
	}

	actualMaxDepth := 1
	if maxDepth != nil {
		actualMaxDepth = *maxDepth
	}

	var depthValue DepthValue

	possibleDepthValue := GetCurrentDepthValue(ctx)
	if possibleDepthValue == nil {
		depthValue = DepthValue{
			ID:           uuid.Must(uuid.NewRandom()).String(),
			MaxDepth:     actualMaxDepth,
			CurrentDepth: 0,
		}
	} else {
		depthValue = *possibleDepthValue
	}

	if len(increments) > 0 && increments[0] {
		depthValue.CurrentDepth++
	}

	ctx = context.WithValue(ctx, DepthKey, depthValue)

	return ctx
}

func WithPathValue(ctx context.Context, tableName string, increments ...bool) context.Context {
	if config.Debug() {
		log.Printf("entered WithPathValue; %#+v", ctx.Value(PathKey))

		defer func() {
			log.Printf("exited WithPathValue; %#+v", ctx.Value(PathKey))
		}()
	}

	var pathValue PathValue

	possiblePathValue := GetCurrentPathValue(ctx)
	if possiblePathValue == nil {
		pathValue = PathValue{
			ID:                uuid.Must(uuid.NewRandom()).String(),
			VisitedTableNames: make([]string, 0),
		}
	} else {
		pathValue = *possiblePathValue
	}

	if len(increments) > 0 && increments[0] {
		pathValue.VisitedTableNames = append(pathValue.VisitedTableNames, tableName)
	}

	ctx = context.WithValue(ctx, PathKey, pathValue)

	return ctx
}

func HandleQueryPathGraphCycles(ctx context.Context, tableName string, increments ...bool) (context.Context, bool) {
	if config.Debug() {
		log.Printf("entered HandleQueryPathGraphCycles for %s (%#+v)", tableName, increments)
	}

	ctx = WithMaxDepth(ctx, helpers.Ptr(1), increments...)
	possibleDepthValue := GetCurrentDepthValue(ctx)
	if possibleDepthValue == nil {
		log.Panicf("assertion failed: DepthValue unexpectedly nil; this should never happen")
	}
	depthValue := *possibleDepthValue

	if depthValue.MaxDepth != 0 && depthValue.CurrentDepth > depthValue.MaxDepth {
		if config.Debug() {
			log.Printf("exited HandleQueryPathGraphCycles for %s (%#+v) after triggering DepthValue", tableName, increments)
		}
		return ctx, false
	}

	ctx = WithPathValue(ctx, tableName, increments...)
	possiblePathValue := GetCurrentPathValue(ctx)
	if possiblePathValue == nil {
		log.Panicf("assertion failed: PathValue unexpectedly nil; this should never happen")
	}
	pathValue := *possiblePathValue

	maxVisitCount := depthValue.MaxDepth
	if maxVisitCount == 0 {
		maxVisitCount = 1
	}

	visitCount := 0

	for _, visitedTableName := range pathValue.VisitedTableNames {
		if visitedTableName != tableName {
			continue
		}

		visitCount++

		if visitCount > maxVisitCount {
			if config.Debug() {
				log.Printf("exited HandleQueryPathGraphCycles for %s (%#+v) after triggering PathValue", tableName, increments)
			}

			return ctx, false
		}
	}

	if config.Debug() {
		log.Printf("exited HandleQueryPathGraphCycles for %s (%#+v)", tableName, increments)
	}

	return ctx, true
}

func GetPaginationDetails(count int64, totalCount int64, rawLimit *int, rawOffset *int) (int64, int64, int64, int64) {
	limit := int64(0)
	if rawLimit != nil {
		limit = int64(*rawLimit)
	}

	offset := int64(0)
	if rawOffset != nil {
		offset = int64(*rawOffset)
	}

	if limit <= 0 {
		return count, totalCount, 1, 1
	}

	totalPages := totalCount / limit
	if totalPages <= 0 {
		totalPages = 1
	}

	cursor := offset + count
	if cursor <= 0 {
		return count, totalCount, 1, totalPages
	}

	page := (cursor / limit) + 1

	return count, totalCount, page, totalPages
}

func WithLoad(ctx context.Context, column *introspect.Column, referencedByColumn *introspect.Column) context.Context {
	loadValueLookup := make(map[string]struct{}, 0)

	if config.Debug() {
		columnSummary := "nil"
		if column != nil && column.ForeignColumn != nil {
			columnSummary = fmt.Sprintf("%s.%s -> %s.%s", column.TableName, column.Name, column.ForeignColumn.TableName, column.ForeignColumn.Name)
		}

		referencedByColumnSummary := "nil"
		if referencedByColumn != nil && referencedByColumn.ForeignColumn != nil {
			referencedByColumnSummary = fmt.Sprintf("%s.%s <- %s.%s", referencedByColumn.ForeignColumn.TableName, referencedByColumn.ForeignColumn.Name, referencedByColumn.TableName, referencedByColumn.Name)
		}

		log.Printf("entered WithLoad for column: %s, referencedByColumn: %s", columnSummary, referencedByColumnSummary)
		defer func() {
			log.Printf("exited WithLoad for column: %s, referencedByColumn: %s; loadValueLookup: %#+v", columnSummary, referencedByColumnSummary, loadValueLookup)
		}()
	}

	if (column == nil || column.ForeignColumn == nil) && (referencedByColumn == nil) {
		return ctx
	}

	rawLoadValueLookup := ctx.Value(LoadKey)
	if rawLoadValueLookup != nil {
		possibleLoadValueLookup, ok := rawLoadValueLookup.(map[string]struct{})
		if ok {
			loadValueLookup = possibleLoadValueLookup
		} else {
			loadValueLookup = make(map[string]struct{}, 0)
		}
	}

	loadValue := fmt.Sprintf("%s__load", column.ForeignColumn.TableName)
	if referencedByColumn != nil {
		loadValue = fmt.Sprintf("referenced_by_%s__load", referencedByColumn.TableName)
	}

	loadValueLookup[loadValue] = struct{}{}

	return context.WithValue(ctx, LoadKey, loadValueLookup)
}

func ShouldLoad(ctx context.Context, column *introspect.Column, referencedByColumn *introspect.Column) bool {
	shouldLoad := false
	var loadValueLookup map[string]struct{}

	if config.Debug() {
		columnSummary := "nil"
		if column != nil && column.ForeignColumn != nil {
			columnSummary = fmt.Sprintf("%s.%s -> %s.%s", column.TableName, column.Name, column.ForeignColumn.TableName, column.ForeignColumn.Name)
		}

		referencedByColumnSummary := "nil"
		if referencedByColumn != nil && referencedByColumn.ForeignColumn != nil {
			referencedByColumnSummary = fmt.Sprintf("%s.%s <- %s.%s", referencedByColumn.ForeignColumn.TableName, referencedByColumn.ForeignColumn.Name, referencedByColumn.TableName, referencedByColumn.Name)
		}

		log.Printf("entered ShouldLoad for column: %s, referencedByColumn: %s", columnSummary, referencedByColumnSummary)
		defer func() {
			log.Printf("exited ShouldLoad for column: %s, referencedByColumn: %s; loadValueLookup: %#+v, shouldLoad: %v", columnSummary, referencedByColumnSummary, loadValueLookup, shouldLoad)
		}()
	}

	if (column == nil || column.ForeignColumn == nil) && (referencedByColumn == nil) {
		return false
	}

	rawLoadValueLookup := ctx.Value(LoadKey)
	if rawLoadValueLookup != nil {
		possibleLoadValueLookup, ok := rawLoadValueLookup.(map[string]struct{})
		if ok {
			loadValueLookup = possibleLoadValueLookup
		} else {
			loadValueLookup = make(map[string]struct{}, 0)
		}
	}

	if column != nil {
		_, shouldLoad = loadValueLookup[fmt.Sprintf("%s__load", column.ForeignColumn.TableName)]
	} else if referencedByColumn != nil {
		_, shouldLoad = loadValueLookup[fmt.Sprintf("referenced_by_%s__load", referencedByColumn.TableName)]
	}

	return shouldLoad
}
