package query

import (
	"context"
	"strings"
	"fmt"

	"github.com/google/uuid"
	"github.com/initialed85/djangolang/pkg/config"
	"github.com/initialed85/djangolang/pkg/helpers"
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

// PathAction describes what HandleQueryPathGraphCycles decided.
// Continue means proceed with the FK load. Skip means don't load.
type PathAction int

const (
	PathActionContinue PathAction = iota // no cycle, no depth limit — proceed
	PathActionDepthExceeded              // depth limit hit — skip, but selectable by shouldLoad
	PathActionCycleDetected              // cycle detected — always skip, never bypassable
)

// HandleQueryPathGraphCycles checks depth limits and cycle detection for a
// given table in the current query path. It returns the (possibly updated)
// context and a PathAction describing why loading may be skipped.
func HandleQueryPathGraphCycles(ctx context.Context, tableName string, increments ...bool) (context.Context, PathAction) {
	if config.Debug() {
		log.Printf("entered HandleQueryPathGraphCycles for %s (%#+v)", tableName, increments)
	}

	ctx = WithMaxDepth(ctx, helpers.Ptr(1), increments...)
	possibleDepthValue := GetCurrentDepthValue(ctx)
	if possibleDepthValue == nil {
		log.Panicf("assertion failed: DepthValue unexpectedly nil; this should never happen")
	}
	depthValue := *possibleDepthValue

	depthExceeded := depthValue.MaxDepth != 0 && depthValue.CurrentDepth > depthValue.MaxDepth

	possiblePathValue := GetCurrentPathValue(ctx)
	cycleDetected := false
	if possiblePathValue != nil {
		for _, visitedPath := range possiblePathValue.VisitedTableNames {
			// Path entries retain their object identity. This means sibling rows
			// such as tag{id-a} and tag{id-b} remain loadable, while an exact
			// object repeat or a repeated collection marker such as tag{<nil>}
			// is a cycle.
			// Extract table name from visited path
			visitedTableName := visitedPath
			if strings.HasPrefix(visitedTableName, "__ReferencedBy__") {
				visitedTableName = strings.TrimPrefix(visitedTableName, "__ReferencedBy__")
			}
			if idx := strings.Index(visitedTableName, "{"); idx != -1 {
				visitedTableName = visitedTableName[:idx]
			}
			// Extract table name from target path
			targetTableName := tableName
			if strings.HasPrefix(targetTableName, "__ReferencedBy__") {
				targetTableName = strings.TrimPrefix(targetTableName, "__ReferencedBy__")
			}
			if idx := strings.Index(targetTableName, "{"); idx != -1 {
				targetTableName = targetTableName[:idx]
			}
			if visitedTableName != targetTableName {
				continue
			}
			// For collection markers ({<nil>}), check if it's the SAME marker
			// If both are collection markers and table names match, it's a cycle
			visitedNil := strings.Contains(visitedPath, "{<nil>}")
			checkNil := strings.Contains(tableName, "{<nil>}")
			if visitedNil && checkNil {
				cycleDetected = true
				break
			}
			// For non-marker entries, check exact match
			if visitedPath == tableName {
				cycleDetected = true
				break
			}
		}
	}

	if cycleDetected {
		if config.Debug() {
			log.Printf("exited HandleQueryPathGraphCycles for %s (%#+v) after triggering PathValue", tableName, increments)
		}
		ctx = WithPathValue(ctx, tableName, increments...)
		possiblePathValue = GetCurrentPathValue(ctx)
		if possiblePathValue == nil {
			log.Panicf("assertion failed: PathValue unexpectedly nil; this should never happen")
		}
		return ctx, PathActionCycleDetected
	}

	if depthExceeded {
		if config.Debug() {
			log.Printf("exited HandleQueryPathGraphCycles for %s (%#+v) after triggering DepthValue", tableName, increments)
		}
		ctx = WithPathValue(ctx, tableName, increments...)
		possiblePathValue = GetCurrentPathValue(ctx)
		if possiblePathValue == nil {
			log.Panicf("assertion failed: PathValue unexpectedly nil; this should never happen")
		}
		return ctx, PathActionDepthExceeded
	}

	ctx = WithPathValue(ctx, tableName, increments...)
	possiblePathValue = GetCurrentPathValue(ctx)
	if possiblePathValue == nil {
		log.Panicf("assertion failed: PathValue unexpectedly nil; this should never happen")
	}

	if config.Debug() {
		log.Printf("exited HandleQueryPathGraphCycles for %s (%#+v)", tableName, increments)
	}

	return ctx, PathActionContinue
}

// DepthExceededFromContext returns true if the current depth exceeds the max
// depth (without modifying context). Used by generated code after
// HandleQueryPathGraphCycles to distinguish depth from cycle.
func DepthExceededFromContext(ctx context.Context) bool {
	pv := GetCurrentDepthValue(ctx)
	if pv == nil {
		return false
	}
	return pv.MaxDepth != 0 && pv.CurrentDepth > pv.MaxDepth
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

func WithLoad(ctx context.Context, tableName string) context.Context {
	loadValueLookup := make(map[string]struct{}, 0)

	if config.Debug() {
		log.Printf("entered WithLoad for tableName: %s", tableName)
		defer func() {
			log.Printf("exited WithLoad for tableName: %s; loadValueLookup: %#+v", tableName, loadValueLookup)
		}()
	}

	if tableName == "" {
		return ctx
	}

	rawLoadValueLookup := ctx.Value(LoadKey)
	possibleLoadValueLookup, ok := rawLoadValueLookup.(map[string]struct{})
	if ok {
		loadValueLookup = possibleLoadValueLookup
	} else {
		loadValueLookup = make(map[string]struct{}, 0)
	}

	loadValue := fmt.Sprintf("%s__load", tableName)

	loadValueLookup[loadValue] = struct{}{}

	ctx = context.WithValue(ctx, LoadKey, loadValueLookup)

	return ctx
}

func ShouldLoad(ctx context.Context, tableName string) bool {
	shouldLoad := false
	var loadValueLookup map[string]struct{}

	if config.Debug() {
		log.Printf("entered ShouldLoad for tableName: %s", tableName)
		defer func() {
			log.Printf("exited ShouldLoad for tableName: %s; loadValueLookup: %#+v, shouldLoad: %v", tableName, loadValueLookup, shouldLoad)
		}()
	}

	if tableName == "" {
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

	_, shouldLoad = loadValueLookup[fmt.Sprintf("%s__load", tableName)]

	return shouldLoad
}

func WithSkip[T any](ctx context.Context, useInstead *T) context.Context {
	return context.WithValue(ctx, ShouldSkipKey, useInstead)
}

func WithoutSkip(ctx context.Context) context.Context {
	return context.WithValue(ctx, ShouldSkipKey, nil)
}

func ShouldSkip[T any](ctx context.Context) (*T, bool) {
	rawUseInstead := ctx.Value(ShouldSkipKey)
	if rawUseInstead == nil {
		return nil, false
	}

	useInstead, ok := rawUseInstead.(*T)
	if !ok {
		return nil, false
	}

	if useInstead == nil {
		return nil, false
	}

	return useInstead, true
}

// ---------------------------------------------------------------------------
// Selection tree helpers (Phase 3 — explicit selection trees)
// ---------------------------------------------------------------------------

// GetSelectionNode looks up the SelectionNode for the given RelationKey
// from the context. Returns nil if no tree is present or the node doesn't
// exist (meaning this edge was not explicitly selected).
// Walks the tree recursively to find nested children.
func GetSelectionNode(ctx context.Context, key RelationKey) *SelectionNode {
	raw := ctx.Value(SelectionTreeKey)
	if raw == nil {
		return nil
	}
	tree, ok := raw.(*SelectionNode)
	if !ok || tree == nil {
		return nil
	}
	if config.Debug() {
		log.Printf("GetSelectionNode: key=%v, root.Children=%d", key, len(tree.Children))
	}
	// First check at root level — match by TableName and LocalColumn, ignoring Direction
	for childKey, child := range tree.Children {
		if childKey.TableName == key.TableName && childKey.LocalColumn == key.LocalColumn {
			return child
		}
	}
	// Then walk nested children
	for _, child := range tree.Children {
		if config.Debug() {
			log.Printf("GetSelectionNode: checking nested child")
		}
		if found := findNodeInTree(child, key); found != nil {
			if config.Debug() {
				log.Printf("GetSelectionNode: found in nested tree")
			}
			return found
		}
	}
	if config.Debug() {
		log.Printf("GetSelectionNode: returning nil")
	}
	return nil
}

// findNodeInTree walks the tree recursively to find a matching RelationKey.
func findNodeInTree(node *SelectionNode, key RelationKey) *SelectionNode {
	if node == nil {
		return nil
	}
	if config.Debug() {
		log.Printf("findNodeInTree: key=%v, node.Children=%d", key, len(node.Children))
	}
	for childKey, child := range node.Children {
		if config.Debug() {
			log.Printf("findNodeInTree: checking child=%v", childKey)
		}
		// Match by TableName and LocalColumn, ignoring Direction
		if childKey.TableName == key.TableName && childKey.LocalColumn == key.LocalColumn {
			return child
		}
	}
	// No direct match; recurse into children
	for _, child := range node.Children {
		if found := findNodeInTree(child, key); found != nil {
			return found
		}
	}
	return nil
}

// WithSelection stores a SelectionNode tree in the context.
func WithSelection(ctx context.Context, tree *SelectionNode) context.Context {
	return context.WithValue(ctx, SelectionTreeKey, tree)
}

// ShouldLoadFromSelectionTree checks if the given table name appears as a
// selected edge anywhere in the selection tree stored in the context.
func ShouldLoadFromSelectionTree(ctx context.Context, tableName string) bool {
	raw := ctx.Value(SelectionTreeKey)
	if raw == nil {
		return false
	}
	tree, ok := raw.(*SelectionNode)
	if !ok || tree == nil {
		return false
	}
	return findNodeByTableName(tree, tableName)
}

// findNodeByTableName walks the tree recursively to find a node whose
// RelationKey TableName matches.
func findNodeByTableName(node *SelectionNode, tableName string) bool {
	if node == nil {
		return false
	}
	for key, child := range node.Children {
		if key.TableName == tableName && child.Selected {
			return true
		}
		if findNodeByTableName(child, tableName) {
			return true
		}
	}
	return false
}

// HasSelectionTree returns true if a selection tree is stored in the context.
// When true, FK loading should respect explicit selection rather than loading by default.
func HasSelectionTree(ctx context.Context) bool {
	raw := ctx.Value(SelectionTreeKey)
	return raw != nil
}
