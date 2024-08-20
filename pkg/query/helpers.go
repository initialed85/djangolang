package query

import (
	"context"
	"log"
)

type shallowKey struct{}

var ShallowKey = shallowKey{}

type PathKey struct {
	TableName string
}

type PathValue struct {
	VisitedTableNames []string
}

func HandleQueryPathGraphCycles(ctx context.Context, tableName string, maxVisitCounts ...int) (context.Context, bool) {
	maxVisitCount := 1
	if len(maxVisitCounts) > 0 {
		maxVisitCount = maxVisitCounts[0]
	}

	rawShallow := ctx.Value(ShallowKey)
	shallow, ok := rawShallow.(bool)
	if ok && shallow {
		maxVisitCount = 0
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
			log.Panicf("expected context key %v to contain a *PathValue but it had %#+v", pathKey, rawPathValue)
		}
	}

	visitCount := 0
	for _, visitedTableName := range pathValue.VisitedTableNames {
		if visitedTableName == tableName {
			if visitCount >= maxVisitCount {
				return ctx, false
			}

			visitCount++
		}
	}

	pathValue.VisitedTableNames = append(pathValue.VisitedTableNames, tableName)
	ctx = context.WithValue(ctx, pathKey, pathValue)

	return ctx, true
}
