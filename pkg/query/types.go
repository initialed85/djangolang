package query

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

type FieldIdentifier struct {
	TableName  string
	ColumnName string
}

type argumentKey struct{}

var ArgumentKey = argumentKey{}
