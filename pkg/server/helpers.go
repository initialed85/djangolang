package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/initialed85/djangolang/pkg/helpers"
	"github.com/initialed85/djangolang/pkg/introspect"
	"github.com/initialed85/djangolang/pkg/query"
)

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

func GetSelectManyArgumentsFromRequest(r *http.Request, table *introspect.Table, parentTable *introspect.Table, parentColumn *introspect.Column) (*SelectManyArguments, error) {
	ctx := r.Context()

	if parentTable == nil || parentColumn == nil {
		argumentsByFieldKey := make(map[query.FieldIdentifier]*SelectManyArguments)
		for _, column := range table.ReferencedByColumns {
			if column.ForeignTable == nil || column.ForeignColumn == nil {
				continue
			}

			theseArguments, err := GetSelectManyArgumentsFromRequest(r, column.ForeignTable, column.ParentTable, column)
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

	expectedParts := 2
	if parentTable != nil || parentColumn != nil {
		expectedParts = 3
	}

	values := make([]any, 0)
	wheres := make([]string, 0)
	for rawKey, rawValues := range r.URL.Query() {
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
		isNullComparison := false
		IsLikeComparison := false

		if !isUnrecognized {
			column := table.ColumnByName[parts[0]]
			if column == nil {
				if expectedParts == 3 {
					continue
				}
				isUnrecognized = true
			} else {
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
					comparison = "IS null"
					isNullComparison = true
				case "isnotnull":
					comparison = "IS NOT null"
					isNullComparison = true
				case "isfalse":
					comparison = "IS false"
					isNullComparison = true
				case "istrue":
					comparison = "IS true"
					isNullComparison = true
				case "like":
					comparison = "LIKE"
					IsLikeComparison = true
				case "notlike":
					comparison = "NOT LIKE"
					IsLikeComparison = true
				case "ilike":
					comparison = "ILIKE"
					IsLikeComparison = true
				case "notilike":
					comparison = "NOT ILIKE"
					IsLikeComparison = true
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

		if isNullComparison {
			wheres = append(wheres, fmt.Sprintf("%s %s", parts[0], comparison))
			continue
		}

		for _, rawValue := range rawValues {
			if isUnrecognized {
				unrecognizedParams = append(unrecognizedParams, fmt.Sprintf("%s=%s", rawKey, rawValue))
				hadUnrecognizedParams = true
				continue
			}

			if hadUnrecognizedParams {
				continue
			}

			attempts := make([]string, 0)

			if !IsLikeComparison {
				attempts = append(attempts, rawValue)
			}

			if isSliceComparison {
				attempts = append(attempts, fmt.Sprintf("[%s]", rawValue))

				vs := make([]string, 0)
				for _, v := range strings.Split(rawValue, ",") {
					vs = append(vs, fmt.Sprintf("\"%s\"", v))
				}

				attempts = append(attempts, fmt.Sprintf("[%s]", strings.Join(vs, ",")))
			}

			if IsLikeComparison {
				attempts = append(attempts, fmt.Sprintf("\"%%%s%%\"", rawValue))
			} else {
				attempts = append(attempts, fmt.Sprintf("\"%s\"", rawValue))
			}

			var err error

			for _, attempt := range attempts {
				var value any

				value, err = time.Parse(time.RFC3339Nano, strings.ReplaceAll(attempt, " ", "+"))
				if err != nil {
					value, err = time.Parse(time.RFC3339, strings.ReplaceAll(attempt, " ", "+"))
					if err != nil {
						err = json.Unmarshal([]byte(attempt), &value)
					}
				}

				if err == nil {
					if isSliceComparison {
						sliceValues, ok := value.([]any)
						if !ok {
							err = fmt.Errorf("failed to cast %#+v to []string", value)
							break
						}

						values = append(values, sliceValues...)

						sliceWheres := make([]string, 0)
						for range values {
							sliceWheres = append(sliceWheres, "$$??")
						}

						wheres = append(wheres, fmt.Sprintf("%s %s (%s)", parts[0], comparison, strings.Join(sliceWheres, ", ")))
					} else {
						values = append(values, value)
						wheres = append(wheres, fmt.Sprintf("%s %s $$??", parts[0], comparison))
					}

					break
				}
			}

			if err != nil {
				unparseableParams = append(unparseableParams, fmt.Sprintf("%s=%s", rawKey, rawValue))
				hadUnparseableParams = true
				continue
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
	rawLimit := r.URL.Query().Get("limit")
	if rawLimit != "" {
		possibleLimit, err := strconv.ParseInt(rawLimit, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse param limit=%s as int; %v", rawLimit, err)
		}

		limit = int(possibleLimit)
	}

	offset := 0
	rawOffset := r.URL.Query().Get("offset")
	if rawOffset != "" {
		possibleOffset, err := strconv.ParseInt(rawOffset, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse param offset=%s as int; %v", rawOffset, err)
		}

		offset = int(possibleOffset)
	}

	depth := 1
	rawDepth := r.URL.Query().Get("depth")
	if rawDepth != "" {
		possibleDepth, err := strconv.ParseInt(rawDepth, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse param depth=%s as int; %v", rawDepth, err)
		}

		depth = int(possibleDepth)

		ctx = query.WithMaxDepth(ctx, &depth)
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

	requestHash, err := helpers.GetRequestHash(table.Name, wheres, hashableOrderBy, limit, offset, depth, values, nil)
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

func GetSelectOneArgumentsFromRequest(r *http.Request, table *introspect.Table, primaryKey any, parentTable *introspect.Table, parentColumn *introspect.Column) (*SelectOneArguments, error) {
	ctx := r.Context()

	wheres := []string{fmt.Sprintf("%s = $$??", table.PrimaryKeyColumn.Name)}
	values := []any{primaryKey}

	unrecognizedParams := make([]string, 0)
	hadUnrecognizedParams := false

	for rawKey, rawValues := range r.URL.Query() {
		if rawKey == "depth" {
			continue
		}

		isUnrecognized := true

		for _, rawValue := range rawValues {
			if isUnrecognized {
				unrecognizedParams = append(unrecognizedParams, fmt.Sprintf("%s=%s", rawKey, rawValue))
				hadUnrecognizedParams = true
				continue
			}

			if hadUnrecognizedParams {
				continue
			}
		}
	}

	if hadUnrecognizedParams {
		return nil, fmt.Errorf("unrecognized params %s", strings.Join(unrecognizedParams, ", "))
	}

	depth := 1
	rawDepth := r.URL.Query().Get("depth")
	if rawDepth != "" {
		possibleDepth, err := strconv.ParseInt(rawDepth, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse param depth=%s as int; %v", rawDepth, err)
		}

		depth = int(possibleDepth)

		ctx = query.WithMaxDepth(ctx, &depth)
	}

	requestHash, err := helpers.GetRequestHash(table.Name, wheres, "", 2, 0, depth, values, primaryKey)
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

func GetLoadArgumentsFromRequest(r *http.Request) (*LoadArguments, error) {
	ctx := r.Context()

	unrecognizedParams := make([]string, 0)
	hadUnrecognizedParams := false

	for rawKey, rawValues := range r.URL.Query() {
		if rawKey == "depth" {
			continue
		}

		isUnrecognized := true

		for _, rawValue := range rawValues {
			if isUnrecognized {
				unrecognizedParams = append(unrecognizedParams, fmt.Sprintf("%s=%s", rawKey, rawValue))
				hadUnrecognizedParams = true
				continue
			}

			if hadUnrecognizedParams {
				continue
			}
		}
	}

	if hadUnrecognizedParams {
		return nil, fmt.Errorf("unrecognized params %s", strings.Join(unrecognizedParams, ", "))
	}

	depth := 1
	rawDepth := r.URL.Query().Get("depth")
	if rawDepth != "" {
		possibleDepth, err := strconv.ParseInt(rawDepth, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse param depth=%s as int; %v", rawDepth, err)
		}

		depth = int(possibleDepth)

		ctx = query.WithMaxDepth(ctx, &depth)
	}

	a := LoadArguments{
		Ctx: ctx,
	}

	return &a, nil
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

	expectedParts := 2
	if parentTable != nil || parentColumn != nil {
		expectedParts = 3
	}

	values := make([]any, 0)
	wheres := make([]string, 0)

outer:
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

		if !isUnrecognized {
			column := table.ColumnByName[parts[0]]
			if column == nil {
				if expectedParts == 3 {
					continue
				}
				isUnrecognized = true
			} else {
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
					comparison = "IS null"
					isValuelessComparison = true
				case "isnotnull":
					comparison = "IS NOT null"
					isValuelessComparison = true
				case "isfalse":
					comparison = "IS false"
					isValuelessComparison = true
				case "istrue":
					comparison = "IS true"
					isValuelessComparison = true
				case "like":
					comparison = "LIKE"
					isLikeComparison = true
				case "notlike":
					comparison = "NOT LIKE"
					isLikeComparison = true
				case "ilike":
					comparison = "ILIKE"
					isLikeComparison = true
				case "notilike":
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

		rawValues := make([]any, 0)

		if isSliceComparison {
			rawValueAsString, ok := rawValue.(string)
			if ok {
				for _, thisRawValueAsString := range strings.Split(rawValueAsString, ",") {
					rawValues = append(rawValues, thisRawValueAsString)
				}
			}
		} else {
			rawValues = append(rawValues, rawValue)
		}

		values := make([]string, 0)

		for _, thisRawValue := range rawValues {
			var valueAsString string

			value, err := json.Marshal(thisRawValue)
			if err != nil {
				unparseableParams = append(unparseableParams, fmt.Sprintf("%s=%s", rawKey, rawValue))
				continue outer
			}

			valueAsString = string(value)

			if strings.HasPrefix(valueAsString, `"`) && strings.HasSuffix(valueAsString, `"`) {
				valueAsString = fmt.Sprintf("%#+v", valueAsString[1:len(valueAsString)-1])

				valueAsString = valueAsString[1 : len(valueAsString)-1]

				if isLikeComparison {
					valueAsString = `%` + valueAsString + `%`
				}

				valueAsString = fmt.Sprintf(`E'%s'`, valueAsString)
			}

			values = append(values, valueAsString)
		}

		valuesAsString := strings.Join(values, ", ")

		if isSliceComparison {
			valuesAsString = fmt.Sprintf("(%s)", valuesAsString)
		}

		wheres = append(wheres, fmt.Sprintf("%s %s %s", parts[0], comparison, valuesAsString))
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
		rawLimit, ok := rawLimit.(string)
		if ok {
			possibleLimit, err := strconv.ParseInt(rawLimit, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("failed to parse param limit=%s as int; %v", rawLimit, err)
			}

			limit = int(possibleLimit)
		}
	}

	offset := 0
	rawOffset, ok := queryParams["offset"]
	if ok {
		rawOffset, ok := rawOffset.(string)
		if ok {
			possibleOffset, err := strconv.ParseInt(rawOffset, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("failed to parse param offset=%s as int; %v", rawOffset, err)
			}

			offset = int(possibleOffset)
		}
	}

	depth := 1
	rawDepth, ok := queryParams["depth"]
	if ok {
		rawDepth, ok := rawDepth.(string)
		if ok {
			possibleDepth, err := strconv.ParseInt(rawDepth, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("failed to parse param depth=%s as int; %v", rawDepth, err)
			}

			depth = int(possibleDepth)

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

	requestHash, err := helpers.GetRequestHash(table.Name, wheres, hashableOrderBy, limit, offset, depth, values, nil)
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

	requestHash, err := helpers.GetRequestHash(table.Name, wheres, "", 2, 0, depth, values, primaryKey)
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
