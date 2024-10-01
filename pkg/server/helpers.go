package server

import (
	"context"
	"encoding/json"
	"fmt"
	_log "log"
	"net/http"
	"strings"

	"github.com/initialed85/djangolang/pkg/config"
	"github.com/initialed85/djangolang/pkg/helpers"
	"github.com/initialed85/djangolang/pkg/introspect"
	"github.com/initialed85/djangolang/pkg/query"
	"github.com/initialed85/djangolang/pkg/types"
)

var log = helpers.GetLogger("server")

func ThisLogger() *_log.Logger {
	return log
}

var cors = config.CORS()

var unknownErrorResponse Response[any] = Response[any]{
	Status:  http.StatusInternalServerError,
	Success: false,
	Error:   []string{fmt.Errorf("unknown error (HTTP %d)", http.StatusInternalServerError).Error()},
	Objects: nil,
}
var unknownErrorResponseJSON []byte

func init() {
	b, err := json.Marshal(unknownErrorResponse)
	if err != nil {
		panic(err)
	}

	unknownErrorResponseJSON = b
}

type Response[T any] struct {
	Status     int      `json:"status"`
	Success    bool     `json:"success"`
	Error      []string `json:"error,omitempty"`
	Objects    []*T     `json:"objects,omitempty"`
	Count      int64    `json:"count"`
	TotalCount int64    `json:"total_count"`
	Limit      int64    `json:"limit"`
	Offset     int64    `json:"offset"`
}

func GetResponse(status int, err error, objects []*any, prettyFormats ...bool) (int, Response[any], []byte, error) {
	prettyFormat := len(prettyFormats) > 1 && prettyFormats[0]

	if status >= http.StatusBadRequest && err == nil {
		err = fmt.Errorf("unspecified error (HTTP %d)", status)
	}

	errorMessage := []string{}
	if err != nil {
		errorMessage = strings.Split(err.Error(), "; ")
	}

	response := Response[any]{
		Status:  status,
		Success: err == nil,
		Error:   errorMessage,
		Objects: objects,
	}

	var b []byte

	if status != http.StatusNoContent {
		if prettyFormat {
			b, err = json.MarshalIndent(response, "", "    ")
		} else {
			b, err = json.Marshal(response)
		}
	}

	if err != nil {
		response = Response[any]{
			Status:  http.StatusInternalServerError,
			Success: false,
			Error: []string{
				fmt.Sprintf("failed to marshal to JSON; status: %d", status),
				fmt.Sprintf("%v", err),
				fmt.Sprintf("objects: %#+v", objects),
			},
			Objects: objects,
		}

		if prettyFormat {
			b, err = json.MarshalIndent(response, "", "    ")
			if err != nil {
				return http.StatusInternalServerError, unknownErrorResponse, unknownErrorResponseJSON, err
			}
		} else {
			b, err = json.Marshal(response)
			if err != nil {
				return http.StatusInternalServerError, unknownErrorResponse, unknownErrorResponseJSON, err
			}
		}

		return http.StatusInternalServerError, response, b, err
	}

	return status, response, b, nil
}

func WriteResponse(w http.ResponseWriter, status int, b []byte) {
	w.Header().Add("Access-Control-Allow-Origin", cors)
	w.Header().Add("Content-Type", "application/json")

	w.WriteHeader(status)
	_, err := w.Write(b)
	if err != nil {
		log.Printf("warning: failed WriteResponse for status: %d: b: %s: %s", status, string(b), err.Error())
		return
	}
}

func HandleErrorResponse(w http.ResponseWriter, status int, err error) {
	status, _, b, _ := GetResponse(status, err, nil, true)
	WriteResponse(w, status, b)
}

func HandleObjectsResponse(w http.ResponseWriter, status int, objects []*any) []byte {
	status, _, b, _ := GetResponse(status, nil, objects)
	WriteResponse(w, status, b)
	return b
}

type SelectManyArguments struct {
	TableName       string
	ParentTableName *string
	ParentFieldName *string
	Wheres          []string
	Where           string
	OrderBy         *string
	Limit           *int
	Offset          *int
	Values          []any
	RequestHash     string
	Ctx             context.Context
}

type SelectOneArguments struct {
	TableName       string
	ParentTableName *string
	ParentFieldName *string
	RequestHash     string
	Wheres          []string
	Where           string
	Values          []any
	Ctx             context.Context
}

type LoadArguments struct {
	Ctx context.Context
}

func GetSelectManyArguments(ctx context.Context, queryParams map[string]any, table *introspect.Table, parentTable *introspect.Table, parentColumn *introspect.Column) (*SelectManyArguments, error) {
	if parentTable == nil || parentColumn == nil {
		argumentsByFieldKey := make(map[query.FieldIdentifier]*SelectManyArguments)
		for _, column := range table.ReferencedByColumns {
			if column.ForeignTable == nil || column.ForeignColumn == nil {
				continue
			}

			theseArguments, err := GetSelectManyArguments(ctx, queryParams, column.ForeignTable, column.ParentTable, column)
			if err != nil {
				return nil, err
			}

			argumentsByFieldKey[query.FieldIdentifier{
				TableName:  column.TableName,
				ColumnName: column.Name,
			}] = theseArguments
		}

		ctx = context.WithValue(ctx, query.ArgumentKey, argumentsByFieldKey)
	}

	insaneOrderParams := make([]string, 0)
	hadInsaneOrderParams := false

	unrecognizedParams := make([]string, 0)
	hadUnrecognizedParams := false

	unparseableParams := make([]string, 0)
	hadUnparseableParams := false

	var orderByDirection *string
	orderBys := make([]string, 0)

	var load string

	expectedParts := 2
	if parentTable != nil || parentColumn != nil {
		expectedParts = 3
	}

	values := make([]any, 0)
	wheres := make([]string, 0)

	// outer:
	for rawKey, rawValue := range queryParams {
		if rawKey == "limit" || rawKey == "offset" || rawKey == "depth" {
			continue
		}

		parts := strings.Split(rawKey, "__")

		if expectedParts == 3 && len(parts) == 2 {
			continue
		}

		if expectedParts == 2 && len(parts) == 3 {
			continue
		}

		isUnrecognized := len(parts) != expectedParts

		if expectedParts == 3 {
			parts = parts[1:]
		}

		comparison := ""
		isSliceComparison := false
		isValuelessComparison := false
		isLikeComparison := false

		var column *introspect.Column

		if !isUnrecognized {
			if parts[1] == "load" {
				if strings.HasPrefix(parts[0], "referenced_by_") {
					tableName := parts[0][len("referenced_by_"):]

					isUnrecognized = true
					for _, possibleColumn := range table.ReferencedByColumns {
						possibleTable := possibleColumn.ParentTable

						if possibleTable.Name == tableName {
							isUnrecognized = false
							ctx = query.WithLoad(ctx, fmt.Sprintf("referenced_by_%s", tableName))
							break
						}
					}
				} else {
					tableName := parts[0]

					isUnrecognized = true
					for _, possibleColumn := range table.Columns {
						if possibleColumn.ForeignTable == nil {
							continue
						}

						if possibleColumn.ForeignTable.Name != tableName {
							continue
						}

						isUnrecognized = false

						ctx = query.WithLoad(ctx, tableName)
						break
					}
				}

				if !isUnrecognized {
					load = parts[0]
					continue
				}
			}

			if !isUnrecognized {
				column = table.ColumnByName[parts[0]]

				if column == nil {
					if expectedParts == 3 {
						continue
					}

					isUnrecognized = true
				}

				if !isUnrecognized {
					switch parts[1] {

					case "eq":
						comparison = "="

					case "ne":
						comparison = "!="

					case "gt":
						comparison = ">"

					case "gte":
						comparison = ">="

					case "lt":
						comparison = "<"

					case "lte":
						comparison = "<="

					case "in":
						comparison = "IN"
						isSliceComparison = true

					case "notin":
						comparison = "NOT IN"
						isSliceComparison = true

					case "isnull":
						if column.NotNull {
							isUnrecognized = true
						}

						comparison = "IS null"
						isValuelessComparison = true

					case "isnotnull":
						if column.NotNull {
							isUnrecognized = true
						}

						comparison = "IS NOT null"
						isValuelessComparison = true

					case "isfalse":
						if !(column.TypeTemplate == "bool" || column.TypeTemplate == "*bool") {
							isUnrecognized = true
						}

						comparison = "IS false"
						isValuelessComparison = true

					case "istrue":
						if !(column.TypeTemplate == "bool" || column.TypeTemplate == "*bool") {
							isUnrecognized = true
						}

						comparison = "IS true"
						isValuelessComparison = true

					case "like":
						if !(column.TypeTemplate == "string" || column.TypeTemplate == "*string") {
							isUnrecognized = true
						}
						comparison = "LIKE"
						isLikeComparison = true

					case "notlike":
						if !(column.TypeTemplate == "string" || column.TypeTemplate == "*string") {
							isUnrecognized = true
						}
						comparison = "NOT LIKE"
						isLikeComparison = true

					case "ilike":
						if !(column.TypeTemplate == "string" || column.TypeTemplate == "*string") {
							isUnrecognized = true
						}
						comparison = "ILIKE"
						isLikeComparison = true

					case "notilike":
						if !(column.TypeTemplate == "string" || column.TypeTemplate == "*string") {
							isUnrecognized = true
						}
						comparison = "NOT ILIKE"
						isLikeComparison = true

					case "desc":
						if orderByDirection != nil && *orderByDirection != "DESC" {
							hadInsaneOrderParams = true
							insaneOrderParams = append(insaneOrderParams, rawKey)
							continue
						}

						orderByDirection = helpers.Ptr("DESC")
						orderBys = append(orderBys, parts[0])
						continue

					case "asc":
						if orderByDirection != nil && *orderByDirection != "ASC" {
							hadInsaneOrderParams = true
							insaneOrderParams = append(insaneOrderParams, rawKey)
							continue
						}

						orderByDirection = helpers.Ptr("ASC")
						orderBys = append(orderBys, parts[0])
						continue

					default:
						isUnrecognized = true
					}
				}
			}
		}

		if isValuelessComparison {
			wheres = append(wheres, fmt.Sprintf("%s %s", parts[0], comparison))
			continue
		}

		if isUnrecognized {
			unrecognizedParams = append(unrecognizedParams, fmt.Sprintf("%s=%s", rawKey, rawValue))
			hadUnrecognizedParams = true
			continue
		}

		if hadUnrecognizedParams {
			continue
		}

		if column == nil {
			log.Panicf("assertion error: made it to query param value handling piece without working out which column applies for %s", fmt.Sprintf("%s=%s", rawKey, rawValue))
		}

		typeMeta, err := types.GetTypeMetaForDataType(column.DataType)
		if err != nil {
			log.Panicf("assertion error: made it to query param value handling with a column that has no TypeMeta %s; %s", fmt.Sprintf("%s=%s", rawKey, rawValue), err.Error())
		}

		value, err := typeMeta.ParseFunc(rawValue)
		if err != nil {
			unparseableParams = append(unparseableParams, fmt.Sprintf("%s=%s (%s)", rawKey, rawValue, err.Error()))
			hadUnparseableParams = true
		}

		if !hadUnparseableParams {
			if isSliceComparison {
				wheres = append(wheres, fmt.Sprintf("%s %s ($$??)", parts[0], comparison))
				values = append(values, value)
			} else if isLikeComparison {
				value, ok := value.(string)
				if ok {
					value = fmt.Sprintf("%#+v", value)
					value = value[1 : len(value)-1]
				}

				wheres = append(wheres, fmt.Sprintf("%s %s E'%%%s%%'", parts[0], comparison, value))
				values = append(values, value)
			} else {
				wheres = append(wheres, fmt.Sprintf("%s %s $$??", parts[0], comparison))
				values = append(values, value)
			}
		}
	}

	if hadUnrecognizedParams {
		return nil, fmt.Errorf("unrecognized params %s", strings.Join(unrecognizedParams, ", "))
	}

	if hadUnparseableParams {
		return nil, fmt.Errorf("unparseable params %s", strings.Join(unparseableParams, ", "))
	}

	if hadInsaneOrderParams {
		return nil, fmt.Errorf("insane order params (e.g. conflicting asc / desc) %s", strings.Join(insaneOrderParams, ", "))
	}

	limit := 50
	rawLimit, ok := queryParams["limit"]
	if ok {
		rawLimit, ok := rawLimit.(float64)
		if ok {
			limit = int(rawLimit)
		}
	}

	offset := 0
	rawOffset, ok := queryParams["offset"]
	if ok {
		rawOffset, ok := rawOffset.(float64)
		if ok {
			offset = int(rawOffset)
		}
	}

	depth := 1
	rawDepth, ok := queryParams["depth"]
	if ok {
		rawDepth, ok := rawDepth.(float64)
		if ok {
			depth = int(rawDepth)

			ctx = query.WithMaxDepth(ctx, &depth)
		}
	}

	hashableOrderBy := ""
	var orderBy *string
	if len(orderBys) > 0 {
		hashableOrderBy = strings.Join(orderBys, ", ")
		if len(orderBys) > 1 {
			hashableOrderBy = fmt.Sprintf("(%v)", hashableOrderBy)
		}
		hashableOrderBy = fmt.Sprintf("%v %v", hashableOrderBy, *orderByDirection)
		orderBy = &hashableOrderBy
	}

	requestHash, err := GetRequestHash(table.Name, wheres, hashableOrderBy, limit, offset, depth, values, nil, load)
	if err != nil {
		return nil, err
	}

	where := strings.Join(wheres, "\n    AND ")

	var parentTableName *string
	if parentTable != nil {
		parentTableName = &parentTable.Name
	}

	var parentFieldName *string
	if parentTable != nil {
		parentFieldName = &parentColumn.Name
	}

	a := SelectManyArguments{
		TableName:       table.Name,
		ParentTableName: parentTableName,
		ParentFieldName: parentFieldName,
		Wheres:          wheres,
		Where:           where,
		OrderBy:         orderBy,
		Limit:           &limit,
		Offset:          &offset,
		Values:          values,
		RequestHash:     requestHash,
		Ctx:             ctx,
	}

	return &a, nil
}

func GetSelectOneArguments(ctx context.Context, rawDepth *int, table *introspect.Table, primaryKey any, parentTable *introspect.Table, parentColumn *introspect.Column) (*SelectOneArguments, error) {
	wheres := []string{fmt.Sprintf("%s = $$??", table.PrimaryKeyColumn.Name)}
	values := []any{primaryKey}

	depth := 1
	if rawDepth != nil {
		depth = *rawDepth

		ctx = query.WithMaxDepth(ctx, &depth)
	}

	requestHash, err := GetRequestHash(table.Name, wheres, "", 2, 0, depth, values, primaryKey, "")
	if err != nil {
		return nil, err
	}

	where := strings.Join(wheres, "\n    AND ")

	var parentTableName *string
	if parentTable != nil {
		parentTableName = &parentTable.Name
	}

	var parentFieldName *string
	if parentTable != nil {
		parentFieldName = &parentColumn.Name
	}

	a := SelectOneArguments{
		TableName:       table.Name,
		ParentTableName: parentTableName,
		ParentFieldName: parentFieldName,
		Wheres:          wheres,
		Where:           where,
		Values:          values,
		RequestHash:     requestHash,
		Ctx:             ctx,
	}

	return &a, nil
}

func GetLoadArguments(ctx context.Context, rawDepth *int) (*LoadArguments, error) {
	depth := 1
	if rawDepth != nil {
		depth = *rawDepth

		ctx = query.WithMaxDepth(ctx, &depth)
	}

	a := LoadArguments{
		Ctx: ctx,
	}

	return &a, nil
}
