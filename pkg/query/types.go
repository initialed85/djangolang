package query

type depthKey struct{}

var DepthKey = depthKey{}

type DepthValue struct {
	ID           string
	MaxDepth     int
	CurrentDepth int
}

type pathKey struct{}

var PathKey = pathKey{}

type PathValue struct {
	ID                string
	VisitedTableNames []string
}

type FieldIdentifier struct {
	TableName  string
	ColumnName string
}

type argumentKey struct{}

var ArgumentKey = argumentKey{}

type loadKey struct{}

var LoadKey = loadKey{}

type shouldSkipKey struct{}

var ShouldSkipKey = shouldSkipKey{}

// ---------------------------------------------------------------------------
// Model-owned query metadata (Phase 1 of explicit selection trees)
// ---------------------------------------------------------------------------

type RelationDirection string

const (
	RelationDirectionOutbound RelationDirection = "outbound" // FK column on this table
	RelationDirectionInbound  RelationDirection = "inbound"  // referenced by this table
)

type RelationCardinality string

const (
	RelationCardinalityOne   RelationCardinality = "one"
	RelationCardinalityMany  RelationCardinality = "many"
)

// RelationshipDefinition describes a single foreign-key relationship
// that a model exposes for loading and filtering.
type RelationshipDefinition struct {
	Name          string               // e.g. "parent_physical_thing_id", "referenced_by_logical_thing_parent_logical_thing_id"
	Direction     RelationDirection    // outbound or inbound
	TargetTable   string               // the table being pointed to
	LocalColumn   string               // the FK column on this model's table
	TargetColumn  string               // the referenced column (usually "id")
	Cardinality   RelationCardinality  // one or many
}

// FilterOperator represents a supported SQL comparison operator.
type FilterOperator string

const (
	FilterOpEq          FilterOperator = "eq"
	FilterOpNe          FilterOperator = "ne"
	FilterOpGt          FilterOperator = "gt"
	FilterOpGte         FilterOperator = "gte"
	FilterOpLt          FilterOperator = "lt"
	FilterOpLte         FilterOperator = "lte"
	FilterOpIn          FilterOperator = "in"
	FilterOpNotIn       FilterOperator = "notin"
	FilterOpContains    FilterOperator = "contains"
	FilterOpNotContains FilterOperator = "notcontains"
	FilterOpIsNull      FilterOperator = "isnull"
	FilterOpIsNotNull   FilterOperator = "isnotnull"
	FilterOpIsFalse     FilterOperator = "isfalse"
	FilterOpIsTrue      FilterOperator = "istrue"
	FilterOpLike        FilterOperator = "like"
	FilterOpNotLike     FilterOperator = "notlike"
	FilterOpILike       FilterOperator = "ilike"
	FilterOpNotILike    FilterOperator = "notilike"
)

// FilterDefinition describes a single column that can be used as a filter.
type FilterDefinition struct {
	Name     string           // JSON tag / column name
	DataType string           // PostgreSQL data type, e.g. "text", "uuid", "boolean"
	Nullable bool             // whether the column allows NULL
	Operators []FilterOperator // supported operators for this column type
}

// ModelQueryInfo is the public query contract for a single model/table.
type ModelQueryInfo struct {
	TableName        string
	PrimaryKeyColumn string
	Relationships    []RelationshipDefinition
	Filters          []FilterDefinition
}

// QueryModel is implemented by every generated model to expose its
// query metadata (relationships, filters) without needing a database
// connection or introspection.
type QueryModel interface {
	GetQueryInfo() ModelQueryInfo
}

// ---------------------------------------------------------------------------
// Selection tree types (for explicit selection trees — Phases 3+)
// ---------------------------------------------------------------------------

// RelationKey uniquely identifies a relationship edge for selection tree lookups.
type RelationKey struct {
	Direction   RelationDirection
	TableName   string
	LocalColumn string
}

// SelectArguments holds the query parameters that apply to a specific
// relation edge (child collection filters, ordering, pagination).
type SelectArguments struct {
	Where    []string
	Values   []any
	OrderBy  *string
	Limit    *int
	Offset   *int
}

// SelectionNode represents one node in the explicit selection tree.
type SelectionNode struct {
	Selected bool
	Filters  SelectArguments
	Children map[RelationKey]*SelectionNode
}

type selectionTreeKey struct{}

var SelectionTreeKey = selectionTreeKey{}
