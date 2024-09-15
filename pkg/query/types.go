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
