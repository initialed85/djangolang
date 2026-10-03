# Task 1: Explicit Selection Trees & Model-Owned Query Metadata

## Status
**Phase 1 — Metadata Types & Generation** (Complete ✅)
**Phase 2 — OpenAPI Migration** (Complete ✅)
- `pkg/openapi/openapi.go`: Added `NewFromQueryModels()`, `modelQueryParams()`, `introspectedRelationshipParams()`, `introspectedFilterParams()`, `dataTypeToSchema()`, `filterOpToMatcher()`
- `pkg/model_reference/0_meta.go`: Updated `GetOpenAPI()` to call `NewFromQueryModels()`
- All OpenAPI tests pass — spec is identical
**Phase 3 — Selection Tree Parser** (Complete ✅ — tree built & stored in context)
- `pkg/server/helpers.go`: Added `buildSelectionTreeFromParams()` — parses query params, builds `SelectionNode` tree with selected edges, stores in context via `SelectionTreeKey`
- `pkg/query/helpers.go`: Added `GetSelectionNode()` and `WithSelection()` helpers
- `pkg/query/types.go`: Added `RelationKey`, `SelectArguments`, `SelectionNode`, `SelectionTreeKey`
- Note: Generated relation loaders still use legacy `ShouldLoad` path. Phase 4 (template changes) needed to wire up selection tree to generated code.
**Phase 4 — Generated Relation Loaders** (Not started)
**Phase 5 — Integration & Cleanup** (Not started)

### Phase 1 Completed
- `pkg/query/types.go`: Added `QueryModel`, `ModelQueryInfo`, `RelationshipDefinition`, `FilterDefinition`, `FilterOperator`, `RelationKey`, `SelectArguments`, `SelectionNode`, `SelectionTreeKey`
- `pkg/template/template.go`: Added `generateGetQueryInfo()` function that produces the method from introspected table data
- `pkg/model_reference/logical_things.go`: Added `// <get-query-info>` fragment placeholder
- All generated models (both `model_generated` and `model_generated_from_schema`) implement `QueryModel`
- All tests pass

---

## Problem Statement

The current system uses `depth` as a wildcard for nested data loading. It's a blunt instrument — you get everything up to depth N or nothing. Query parsing is limited to one-hop filters on the root object. OpenAPI generation introspects the schema separately from the models themselves, creating a parallel interpretation of the same data.

We need:
1. **Explicit selection trees** — GraphQL-style parameter paths that control exactly what gets loaded, filtered, and ordered at each edge.
2. **Model-owned query metadata** — Models describe their own public query contract (relationships, filters), replacing OpenAPI's parallel schema introspection.

---

## Current Architecture

### Key Files & Their Roles

| File | Role |
|------|------|
| `pkg/server/helpers.go` | `GetSelectManyArguments()` — parses query params into `SelectManyArguments` (where/order/limit/offset). Stores nested args in `query.ArgumentKey` context. |
| `pkg/server/types.go` | `HTTPHandler` / `HTTPHandlerSummary` — carries `BuiltinModelObject` and `BuiltinTable`. `WithPrimaryKey`, `WithReload`, `WithInsert` interfaces. |
| `pkg/query/types.go` | Context keys: `DepthKey`, `PathKey`, `ArgumentKey`, `LoadKey`, `ShouldSkipKey`. `DepthValue`, `PathValue`, `FieldIdentifier`. |
| `pkg/query/helpers.go` | `WithMaxDepth()`, `WithPathValue()`, `HandleQueryPathGraphCycles()` — depth tracking + identity-aware cycle detection. `ShouldLoad()`, `WithLoad()` — load markers. |
| `pkg/openapi/openapi.go` | `NewFromIntrospectedSchema()` — walks `HTTPHandlerSummary.BuiltinTable` and `BuiltinIntrospectedModelObject` to build OpenAPI spec. Hardcoded `matchers` (eq, ne, gt, etc.). |
| `pkg/model_reference/0_meta.go` | `register()` — registers tables, objects, newFromItem fns, router fns. `GetOpenAPI()` — delegates to `openapi.NewFromIntrospectedSchema()`. |
| `pkg/model_generated/logical_things.go` | Generated models. Relation loading in `SelectLogicalThings()` handles FK loads with `HandleQueryPathGraphCycles()` + `ShouldLoad()`. Hardcoded FK constraints like `fmt.Sprintf("%v = $1", PhysicalThingTablePrimaryKeyColumn)`. |
| `pkg/introspect/types.go` | `Table` (Columns, ColumnByName, ForeignTables, ReferencedByColumns), `Column` (ForeignTable, ForeignColumn, ParentTable, ParentID). |
| `pkg/types/types.go` | `TypeMeta[T]` — per-data-type metadata (parse/format/schema functions). Used by query parsing and template generation. |

### How Query Parsing Currently Works

```
HTTP Request
    → server.HTTPHandler.ServeHTTP()
        → GetSelectManyArguments(ctx, queryParams, table, nil, nil)
            → For each ReferencedByColumn:
                → Recursively call GetSelectManyArguments(ctx, queryParams, foreignTable, parentTable, column)
            → Parse __load params → query.WithLoad(ctx, tableName)
            → Parse column__eq/etc params → build wheres + values
            → Store nested args in context: context.WithValue(ctx, query.ArgumentKey, argumentsByFieldKey)
            → Return SelectManyArguments{Where, OrderBy, Limit, Offset, Values, Ctx}
        → Call generated handler (handleGetLogicalThings)
            → SelectLogicalThings(ctx, tx, where, orderBy, limit, offset, values...)
                → For each object: load FK relations via HandleQueryPathGraphCycles + ShouldLoad
```

### How Relation Loading Works (Generated Code)

In `SelectLogicalThings()` (generated), for each object:

```go
if !types.IsZeroUUID(object.ParentPhysicalThingID) {
    ctx, ok := query.HandleQueryPathGraphCycles(ctx, fmt.Sprintf("%s{%v}", PhysicalThingTable, object.ParentPhysicalThingID), true)
    shouldLoad := query.ShouldLoad(ctx, PhysicalThingTable)
    if ok || shouldLoad {
        object.ParentPhysicalThingIDObject, _, _, _, _, err = SelectPhysicalThing(
            ctx, tx,
            fmt.Sprintf("%v = $1", PhysicalThingTablePrimaryKeyColumn),  // HARDCODED FK constraint
            object.ParentPhysicalThingID,
        )
    }
}
```

Key observations:
- FK constraints are hardcoded strings like `fmt.Sprintf("%v = $1", PhysicalThingTablePrimaryKeyColumn)`
- No child-specific `where`/`order`/`limit` — uses parent's arguments
- `ShouldLoad` checks `ShouldLoad(ctx, PhysicalThingTable)` which looks for `"physical_things__load"` in the load map
- `HandleQueryPathGraphCycles` prevents cycles using identity-aware path tracking (pushed as `e337e01257b7ca320ebaadc64ae6337b6cef95db`)

### How OpenAPI Generation Works

`openapi.NewFromIntrospectedSchema(httpHandlerSummaries)`:
1. Walks `HTTPHandlerSummary.BuiltinTable.Columns` for FK relationship params (`{table}__load`)
2. Walks `HTTPHandlerSummary.BuiltinTable.ReferencedByColumns` for reverse FK params (`referenced_by_{table}__load`)
3. Walks `BuiltinIntrospectedModelObject.StructFields` for filter params (`{column}__eq`, `{column}__gt`, etc.)
4. Uses hardcoded `matchers` slice and `descriptionByMatcher` map
5. Ignores `BuiltinModelObject` entirely — doesn't use any model-owned metadata

### Context Key Chain

```
query.DepthKey       → DepthValue{ID, MaxDepth, CurrentDepth}
query.PathKey        → PathValue{ID, VisitedTableNames}
query.ArgumentKey    → map[FieldIdentifier]*SelectManyArguments  (nested per-relation args)
query.LoadKey        → map[string]struct{}  (e.g., {"physical_things__load": {}})
query.ShouldSkipKey  → *T  (generic skip mechanism)
```

---

## Proposed Architecture

### Part A: Model-Owned Query Metadata

#### New Interface

```go
type QueryModel interface {
    GetQueryInfo() query.ModelQueryInfo
}
```

#### New Types (in `pkg/query/types.go`)

```go
type RelationDirection string
const (
    DirectionOutbound  RelationDirection = "outbound"   // FK from this table
    DirectionInbound   RelationDirection = "inbound"    // referenced by this table
)

type RelationshipCardinality string
const (
    CardinalityOne   RelationCardinality = "one"
    CardinalityMany  RelationCardinality = "many"
)

type RelationshipDefinition struct {
    Name          string                // e.g., "parent_physical_thing_id", "referenced_by_logical_thing_parent_logical_thing_id"
    Direction     RelationDirection
    TargetTable   string
    LocalColumn   string
    TargetColumn  string
    Cardinality   RelationCardinality
}

type FilterOperator string
const (
    OpEq, OpNe, OpGt, OpGte, OpLt, OpLte, OpIn, OpNotIn
    OpContains, OpNotContains
    OpIsNull, OpIsNotNull
    OpIsFalse, OpIsTrue
    OpLike, OpNotLike, OpILike, OpNotILike
)

type FilterDefinition struct {
    Name       string
    DataType   string  // e.g., "text", "uuid", "boolean"
    Nullable   bool
    Operators  []FilterOperator
}

type ModelQueryInfo struct {
    TableName        string
    PrimaryKeyColumn string
    Relationships    []RelationshipDefinition
    Filters          []FilterDefinition
}
```

#### Generated Implementation (example for LogicalThing)

```go
func (LogicalThing) GetQueryInfo() query.ModelQueryInfo {
    return query.ModelQueryInfo{
        TableName:        "logical_things",
        PrimaryKeyColumn: "id",
        Relationships: []query.RelationshipDefinition{
            {
                Name:        "parent_physical_thing_id",
                Direction:   query.DirectionOutbound,
                TargetTable: "physical_things",
                LocalColumn: "parent_physical_thing_id",
                TargetColumn: "id",
                Cardinality: query.CardinalityOne,
            },
            {
                Name:        "parent_logical_thing_id",
                Direction:   query.DirectionOutbound,
                TargetTable: "logical_things",
                LocalColumn: "parent_logical_thing_id",
                TargetColumn: "id",
                Cardinality: query.CardinalityOne,
            },
            {
                Name:        "referenced_by_logical_thing_parent_logical_thing_id",
                Direction:   query.DirectionInbound,
                TargetTable: "logical_things",
                LocalColumn: "parent_logical_thing_id",
                TargetColumn: "id",
                Cardinality: query.CardinalityMany,
            },
        },
        Filters: []query.FilterDefinition{
            {Name: "id", DataType: "uuid", Nullable: false, Operators: [...]},
            {Name: "name", DataType: "text", Nullable: false, Operators: [...]},
            {Name: "type", DataType: "text", Nullable: false, Operators: [...]},
            {Name: "created_at", DataType: "timestamp with time zone", Nullable: false, Operators: [...]},
            {Name: "updated_at", DataType: "timestamp with time zone", Nullable: false, Operators: [...]},
            {Name: "deleted_at", DataType: "timestamp with time zone", Nullable: true, Operators: [...]},
            // ... all columns
        },
    }
}
```

### Part B: Recursive Selection Tree

#### New Types (in `pkg/query/types.go`)

```go
type RelationKey struct {
    Direction   RelationDirection
    TableName   string
    LocalColumn string
}

type SelectArguments struct {
    Where    []string
    Values   []any
    OrderBy  *string
    Limit    *int
    Offset   *int
}

type SelectionNode struct {
    Selected bool
    Filters  SelectArguments
    Children map[RelationKey]*SelectionNode
}
```

#### New Context Keys

```go
type selectionTreeKey struct{}
var SelectionTreeKey = selectionTreeKey{}
```

#### Desugaring Ancestor Loads

Query param: `referenced_by_meme_tag__tag__load=`

Should imply:
- `referenced_by_meme_tag__load=` (ancestor)
- `referenced_by_meme_tag__tag__load=` (target)

The parser creates all missing ancestors automatically.

#### Query Param Grammar

```
{relation_path}__load=                    → select this edge
{relation_path}.{column}__eq=value        → filter child collection on this edge
{relation_path}.{column}__gt=value        → filter child collection on this edge
{relation_path}__order_by={column}__desc  → order child collection
{relation_path}__limit=N                  → limit child collection
{relation_path}__offset=N                 → offset child collection
```

Where `relation_path` is a `.`-separated chain of relation names:
- `parent_physical_thing_id` (outbound FK)
- `referenced_by_logical_thing_parent_logical_thing_id` (inbound FK)

### Part C: Relation-Aware Generated Loaders

#### Generated Code Concept

```go
// In SelectLogicalThings()
relationKey := query.RelationKey{
    Direction:   query.DirectionOutbound,
    TableName:   PhysicalThingTable,
    LocalColumn: LogicalThingTableParentPhysicalThingIDColumn,
}

childNode := query.GetSelectionNode(ctx, relationKey)

if childNode != nil && childNode.Selected {
    ctx, cycleOk := query.HandleQueryPathGraphCycles(ctx, pathStr, true)
    if cycleOk {
        // Combine FK constraint with child filters
        where := []string{
            fmt.Sprintf("%s = $$??", PhysicalThingTablePrimaryKeyColumn),
        }
        values := []any{object.ParentPhysicalThingID}
        
        // Append child-specific filters
        where = append(where, childNode.Filters.Where...)
        values = append(values, childNode.Filters.Values...)
        
        childCtx := query.WithSelection(ctx, childNode)
        
        object.ParentPhysicalThingIDObject, _, _, _, _, err = SelectPhysicalThing(
            childCtx, tx,
            strings.Join(where, "\n    AND "),
            childNode.Filters.OrderBy,
            childNode.Filters.Limit,
            childNode.Filters.Offset,
            values...,
        )
    }
}
```

Key points:
- FK constraint and user filters share the `$$??` placeholder system
- Each relationship has its own `where`/`order`/`limit`/`offset`
- Child receives the child selection subtree via context
- Cycle protection remains unconditional (never bypassed)

### Part D: OpenAPI Consumes Model Metadata

```go
// Old:
openApi, err = openapi.NewFromIntrospectedSchema(httpHandlerSummaries)
    → openapi walks BuiltinTable.Columns, BuiltinTable.ReferencedByColumns,
      BuiltinIntrospectedModelObject.StructFields

// New:
openApi, err = openapi.NewFromQueryModels(httpHandlerSummaries)
    → openapi walks BuiltinModelObject.(QueryModel).GetQueryInfo()
    → Falls back to BuiltinTable introspection for unknown models
```

The introspected schema still generates the metadata, but OpenAPI consumes the model's declared public query contract.

---

## Implementation Order

### Phase 1: Metadata Types & Generation

1. Define `ModelQueryInfo`, `RelationshipDefinition`, `FilterDefinition`, `RelationKey`, `SelectArguments`, `SelectionNode` in `pkg/query/types.go`
2. Update template in `pkg/template/template.go` to generate `GetQueryInfo()` for every model
3. Test: ensure every generated model implements `QueryModel`

### Phase 2: OpenAPI Migration

4. Create `openapi.NewFromQueryModels()` that consumes `QueryModel` metadata
5. Keep `NewFromIntrospectedSchema()` as fallback for backward compat
6. Update `model_reference/0_meta.go` to use new function
7. Test: verify OpenAPI spec is identical (or better) than before

### Phase 3: Selection Tree Parser

8. Add `SelectionTreeKey` context key
9. Extend `GetSelectManyArguments()` to build `SelectionNode` tree from query params
10. Implement ancestor desugaring (e.g., `a__b__load` implies `a__load`)
11. Handle nested `where`/`order`/`limit`/`offset` on relation edges
12. Test: verify tree structure matches query params

### Phase 4: Generated Relation Loaders

13. Update template to generate relation loaders that:
    - Look up child `SelectionNode` from context
    - Combine FK constraint with child filters using `$$??` placeholders
    - Pass child `where`/`order`/`limit`/`offset` to child `Select*` calls
    - Pass child subtree via context
14. Test: verify relation loading works with nested filters/orders/limits

### Phase 5: Integration & Cleanup

15. Make `depth` a compatibility feature (no explicit tree → use depth behavior)
16. Add server-side max expansion depth safety ceiling
17. Update `model_reference` handlers to use new flow
18. Remove unused `ArgumentKey` context mechanism (or deprecate)
19. Update tests

---

## Key Design Decisions

### 1. Relation Identity

Cannot just use target table name — multiple FKs may point to the same table. Must include:
- Direction (outbound/inbound)
- Target table
- Local column (the FK column on the source table)

Example:
```
meme_tag → meme (outbound, local: meme_id)
meme_tag → tag  (outbound, local: tag_id)
meme  → meme (inbound, local: meme_id via meme_tag)
```

Public names: `referenced_by_meme_tag`, `referenced_by_tag`, etc.

### 2. Filtering Scope

Two distinct operations:
- **Filter loaded child collection**: `?referenced_by_meme_tag__tag__name__eq=cloud` — only populate `tag_id_object` when tag matches
- **Filter root collection by child existence**: `?exists[referenced_by_meme_tag][tag][name__eq]=cloud` — only return memes that have at least one matching tag

First implementation focuses on child collection filtering. Parent filtering via `EXISTS`/joins is deferred.

### 3. Depth as Wildcard

- No explicit selection tree → preserve current `depth` behavior
- Explicit selection tree present → load only selected paths
- `depth` remains as a compatibility/default expansion mechanism
- Cycle protection remains unconditional

### 4. OpenAPI Guardrails

OpenAPI cannot recursively emit infinite parameters for self-referential models. Walker should:
- Generate useful parameters for one or two levels
- Add a generic documented path parameter for arbitrary deeper paths
- Example: emit `referenced_by_meme_tag__load`, `referenced_by_meme_tag__tag__load`, `referenced_by_meme_tag__tag__name__eq`, then document that the same grammar can continue recursively

### 5. Placeholder Convention

Continue using `$$??` for parameterized values in generated code. This avoids hardcoding `$1`, `$2`, etc. in templates. The runtime fills in correct positional parameters.

---

## Related Commits & Files

- **Cycle detection fix**: `e337e01257b7ca320ebaadc64ae6337b6cef95db` on `origin/main`
  - Modified: `pkg/query/helpers.go`, `pkg/query/helpers_test.go`
  - Uses identity-aware path tracking with `{<nil>}` collection markers

- **Field update refactoring**: `ParseTask` per-column snippets
  - Modified: `pkg/template/parse.go`, `pkg/model_reference/logical_things.go`

- **Generated models**: `pkg/model_generated/logical_things.go`
  - Contains `SelectLogicalThings()`, `SelectLogicalThing()`, `Reload()`, `Insert()`, `Update()`, `UpdateField()`, `UpdateFields()`
  - Relation loading in `SelectLogicalThings()` at ~line 1476

- **Template**: `pkg/template/template.go`
  - Generates all model code from introspected schema

- **Model reference**: `pkg/model_reference/logical_things.go`
  - Contains custom handlers (e.g., `ClaimLogicalThing`)
  - Calls `server.GetSelectManyArguments()`

---

## Testing Strategy

### Unit Tests
- `query`: Selection tree parsing, ancestor desugaring, relation key equality
- `template`: Verify `GetQueryInfo()` generates correct metadata for all types
- `openapi`: Verify spec matches expected output

### Integration Tests
- `model_generated_test/`: Existing tests + new tests for nested selection
- Query parsing: Verify `?referenced_by_meme_tag__tag__name__in=cloud,devops` produces correct tree and SQL
- Relation loading: Verify child-specific `where`/`order`/`limit` applied correctly

### Regression Tests
- `depth=N` still works (no selection tree → depth behavior)
- `?tag__load=` still works (simple load without nested args)
- Cycle detection still prevents `meme → meme_tag → meme`

---

## Risks & Mitigations

1. **Generated code bloat**: Adding `GetQueryInfo()` adds ~50 lines per model. Mitigation: static data, no runtime cost.
2. **Template complexity**: Relation-aware loaders are harder to generate. Mitigation: Phase 4 after metadata is stable.
3. **OpenAPI spec changes**: New parameters may break existing clients. Mitigation: additive changes only, keep `depth` working.
4. **SQL injection**: User filters combined with FK constraints. Mitigation: `$$??` placeholders, type-checked parsing.
5. **Query complexity**: Deep nested filters could produce complex SQL. Mitigation: server-side max depth, query timeout.

---

## Environment Notes

- Run tests with: `./run.sh test-ci`
- Start environment: `./run.sh env up`
- Stop environment: `./run.sh env down`
- Cross-project port collision with Camry (PostgreSQL 5432, Redis 6379, Swagger 7071)
- Generated files under `pkg/model_generated/` and `pkg/model_generated_from_schema/` — never hand-edit
- If `git restore` fails due to root ownership: `sudo chown -R $(id -u):$(id -g) pkg/model_generated pkg/model_generated_from_schema`
