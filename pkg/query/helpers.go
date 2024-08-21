package query

import (
	"context"
	"log"
	"strings"
)

type depthKey struct{}

var DepthKey = depthKey{}

type DepthValue struct {
	MaxDepth     int
	currentDepth int
}

type PathKey struct {
	TableName string
}

type PathValue struct {
	VisitedTableNames []string
}

func WithMaxDepth(ctx context.Context, maxDepth *int) context.Context {
	actualMaxDepth := -1
	if maxDepth != nil {
		actualMaxDepth = *maxDepth
	}

	var depthValue DepthValue
	rawDepthValue := ctx.Value(DepthKey)
	if rawDepthValue == nil {
		depthValue = DepthValue{
			MaxDepth:     actualMaxDepth,
			currentDepth: -1,
		}
	} else {
		var castOk bool
		depthValue, castOk = rawDepthValue.(DepthValue)
		if !castOk {
			log.Panicf("expected context key %v to contain a DepthValue but it had %#+v", DepthKey, rawDepthValue)
		}
	}

	ctx = context.WithValue(ctx, DepthKey, depthValue)

	return ctx
}

func HandleQueryPathGraphCycles(ctx context.Context, tableName string, maxVisitCounts ...int) (context.Context, bool) {
	var depthValue DepthValue
	rawDepthValue := ctx.Value(DepthKey)
	if rawDepthValue != nil {
		var castOk bool
		depthValue, castOk = rawDepthValue.(DepthValue)
		if !castOk {
			log.Panicf("expected context key %v to contain a DepthValue but it had %#+v", DepthKey, rawDepthValue)
		}
	}

	if depthValue.MaxDepth != 0 {
		// TODO: this is a bit gross and implicit but it will do for now
		if !strings.HasPrefix(tableName, "__ReferencedBy__") {
			depthValue.currentDepth++
		}

		ctx = context.WithValue(ctx, DepthKey, depthValue)

		if depthValue.currentDepth > depthValue.MaxDepth {
			return ctx, false
		}

		return ctx, true
	}

	maxVisitCount := 1
	if len(maxVisitCounts) > 0 {
		maxVisitCount = maxVisitCounts[0]
	}

	if maxVisitCount == 0 {
		return ctx, false
	}

	pathKey := PathKey{TableName: tableName}
	var pathValue PathValue
	rawPathValue := ctx.Value(pathKey)
	if rawPathValue == nil {
		pathValue = PathValue{
			VisitedTableNames: []string{},
		}
	} else {
		var castOk bool
		pathValue, castOk = rawPathValue.(PathValue)
		if !castOk {
			log.Panicf("expected context key %v to contain a PathValue but it had %#+v", pathKey, rawPathValue)
		}
	}

	visitCount := 0
	for _, visitedTableName := range pathValue.VisitedTableNames {
		if visitedTableName != tableName {
			continue
		}

		if visitCount >= maxVisitCount {
			return ctx, false
		}

		visitCount++
	}

	pathValue.VisitedTableNames = append(pathValue.VisitedTableNames, tableName)
	ctx = context.WithValue(ctx, pathKey, pathValue)

	return ctx, true
}
