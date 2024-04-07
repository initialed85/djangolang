package template

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/chanced/caps"
	_pluralize "github.com/gertd/go-pluralize"
	_helpers "github.com/initialed85/djangolang/internal/helpers"
	"github.com/initialed85/djangolang/pkg/helpers"
	"github.com/initialed85/djangolang/pkg/introspect"
	"golang.org/x/exp/maps"
)

var (
	logger    = helpers.GetLogger("generate")
	pluralize = _pluralize.NewClient()
)

func init() {
	converter, _ := caps.DefaultConverter.(caps.StdConverter)
	converter.Set("Dob", "DOB")
	converter.Set("Cpo", "CPO")
	converter.Set("Mwh", "MWH")
	converter.Set("Kwh", "KWH")
	converter.Set("Wh", "WH")
	converter.Set("Json", "JSON")
	converter.Set("Jsonb", "JSONB")
	converter.Set("Mac", "MAC")
	converter.Set("Ip", "IP")
}

func run(ctx context.Context, finalOutputPath string, tempOutputPath string) error {
	parts := strings.Split(strings.TrimRight(finalOutputPath, "/"), "/")

	packageName := parts[len(parts)-1]

	logger.Printf("preparing %v...", tempOutputPath)
	_ = os.RemoveAll(tempOutputPath)
	err := os.MkdirAll(tempOutputPath, 0o777)
	if err != nil {
		return err
	}
	logger.Printf("done.")

	conn, err := helpers.GetDBFromEnvironment(ctx)
	if err != nil {
		return err
	}
	defer func() {
		_ = conn.Close()
	}()

	schema := helpers.GetSchema()

	tableByName, err := introspect.Introspect(ctx, conn, schema)
	if err != nil {
		return err
	}

	tableNames := maps.Keys(tableByName)
	slices.Sort(tableNames)

	fileData := ""

	initData := ""

	for _, tableName := range tableNames {
		tableFileData := ""

		table := tableByName[tableName]

		structNameSingular := pluralize.Singular(caps.ToCamel(table.Name))
		structNamePlural := pluralize.Plural(caps.ToCamel(table.Name))

		if table.RelKind == "v" {
			structNameSingular += "View"
			structNamePlural += "View"
		}

		initData += fmt.Sprintf(
			"    selectFuncByTableName[%#+v] = genericSelect%v\n",
			table.Name,
			structNamePlural,
		)

		initData += fmt.Sprintf(
			"    insertFuncByTableName[%#+v] = genericInsert%v\n",
			table.Name,
			structNameSingular,
		)

		initData += fmt.Sprintf(
			"    updateFuncByTableName[%#+v] = genericUpdate%v\n",
			table.Name,
			structNameSingular,
		)

		initData += fmt.Sprintf(
			"    deleteFuncByTableName[%#+v] = genericDelete%v\n",
			table.Name,
			structNameSingular,
		)

		initData += fmt.Sprintf(
			"    deserializeFuncByTableName[%#+v] = Deserialize%v\n",
			table.Name,
			structNameSingular,
		)

		initData += fmt.Sprintf(
			"    columnNamesByTableName[%#+v] = %vColumns\n",
			table.Name,
			structNameSingular,
		)

		initData += fmt.Sprintf(
			"    transformedColumnNamesByTableName[%#+v] = %vTransformedColumns\n",
			table.Name,
			structNameSingular,
		)

		nameData := fmt.Sprintf(
			"var %vTable = \"%v\"\n",
			structNameSingular,
			table.Name,
		)

		structData := fmt.Sprintf(
			"type %v struct {\n",
			structNameSingular,
		)

		postprocessData := ""

		columnNames := make([]string, 0)
		insertColumnNames := make([]string, 0)
		transformedColumnNames := make([]string, 0)

		foreignObjectIDMapDeclarations := ""
		foreignObjectIDMapPopulations := ""
		foreignObjectLoading := ""

		primaryKeyColumnVariable := ""
		primaryKeyStructField := ""
		primaryKeyStructType := ""

		for _, column := range table.Columns {
			if column.ZeroType == nil && column.TypeTemplate != "any" {
				continue
			}

			fieldName := caps.ToCamel(column.Name)

			columnNames = append(columnNames, column.Name)

			if !column.IsPrimaryKey {
				insertColumnNames = append(insertColumnNames, column.Name)
			}

			transformedColumnName := column.Name

			if column.DataType == "point" || column.DataType == "polygon" || column.DataType == "geometry" {
				transformedColumnName = fmt.Sprintf("ST_AsGeoJSON(%v::geometry)::jsonb AS %v", column.Name, column.Name)
			} else if column.DataType == "tsvector" {
				transformedColumnName = fmt.Sprintf("array_to_json(tsvector_to_array(%v))::jsonb AS %v", column.Name, column.Name)
			} else if column.DataType == "interval" {
				transformedColumnName = fmt.Sprintf("extract(microseconds FROM %v)::numeric * 1000 AS %v", column.Name, column.Name)
			}

			transformedColumnNames = append(transformedColumnNames, transformedColumnName)

			nameData += fmt.Sprintf(
				"var %vTable%vColumn = \"%v\"\n",
				structNameSingular,
				fieldName,
				column.Name,
			)

			if column.IsPrimaryKey {
				primaryKeyColumnVariable = fmt.Sprintf(
					"%vTable%vColumn",
					structNameSingular,
					fieldName,
				)

				primaryKeyStructField = fieldName

				primaryKeyStructType = column.TypeTemplate
			}

			structData += fmt.Sprintf(
				"    %v %v `json:\"%v\" db:\"%v\"`\n",
				fieldName,
				column.TypeTemplate,
				column.Name,
				column.Name,
			)

			if column.ForeignColumn != nil {
				foreignObjectNameSingular := pluralize.Singular(caps.ToCamel(column.ForeignTable.Name))
				foreignObjectNamePlural := pluralize.Plural(caps.ToCamel(column.ForeignTable.Name))

				structData += fmt.Sprintf(
					"    %vObject *%v `json:\"%v_object,omitempty\"`\n",
					fieldName,
					foreignObjectNameSingular,
					column.Name,
				)

				foreignObjectMapName := fmt.Sprintf(
					"%vBy%v",
					caps.ToLowerCamel(foreignObjectNameSingular),
					fieldName,
				)

				foreignObjectIDMapDeclarations += fmt.Sprintf(
					"%v := make(map[%v]*%v, 0)\n    ",
					foreignObjectMapName,
					strings.TrimLeft(column.TypeTemplate, "*"),
					foreignObjectNameSingular,
				)

				if strings.HasPrefix(column.TypeTemplate, "*") {
					foreignObjectIDMapPopulations += fmt.Sprintf(
						"%v[helpers.Deref(item.%v)] = nil\n        ",
						foreignObjectMapName,
						fieldName,
					)
				} else {
					foreignObjectIDMapPopulations += fmt.Sprintf(
						"%v[item.%v] = nil\n        ",
						foreignObjectMapName,
						fieldName,
					)
				}

				possiblePointerFieldName := ""
				if strings.HasPrefix(column.TypeTemplate, "*") {
					possiblePointerFieldName = fmt.Sprintf("helpers.Deref(item.%v)", fieldName)
				} else {
					possiblePointerFieldName = fmt.Sprintf("item.%v", fieldName)
				}

				columnRef := ""
				if strings.HasPrefix(column.ForeignColumn.TypeTemplate, "*") {
					columnRef = fmt.Sprintf("helpers.Deref(row.%v)", caps.ToCamel(column.ForeignColumn.Name))
				} else {
					columnRef = fmt.Sprintf("row.%v", caps.ToCamel(column.ForeignColumn.Name))
				}

				foreignObjectLoading += fmt.Sprintf(
					loadTemplate,
					fieldName,
					foreignObjectMapName,
					fieldName,
					fieldName,
					fieldName,
					fieldName,
					foreignObjectNamePlural,
					foreignObjectNameSingular,
					column.ForeignColumn.Name,
					fieldName,
					fieldName,
					foreignObjectMapName,
					columnRef,
					fieldName,
					foreignObjectMapName,
					possiblePointerFieldName,
				) + "\n    \n    "
			}

			if column.ZeroType == nil && (column.TypeTemplate == "any" || column.TypeTemplate == "*any") {
				postprocessData += fmt.Sprintf(
					jsonBlockTemplate,
					fieldName,
					fieldName,
					fieldName,
					fieldName,
					fieldName,
					fieldName,
					fieldName,
					fieldName,
					fieldName,
					fieldName,
					fieldName,
					fieldName,
					fieldName,
					fieldName,
				)
			}
		}

		hasPrimaryKey := primaryKeyColumnVariable != "" && primaryKeyStructField != ""

		if postprocessData == "" {
			postprocessData = "        // this is where any post-scan processing would appear (if required)"
		}

		nameData += fmt.Sprintf("var %vColumns = %#+v\n", structNameSingular, columnNames)
		nameData += fmt.Sprintf("var %vTransformedColumns = %#+v\n", structNameSingular, transformedColumnNames)
		nameData += fmt.Sprintf("var %vInsertColumns = %#+v\n", structNameSingular, insertColumnNames)
		nameData += "\n"

		tableFileData += nameData

		if foreignObjectIDMapDeclarations == "" {
			foreignObjectIDMapDeclarations = "        // this is where any foreign object ID map declarations would appear (if required)"
		}

		if foreignObjectIDMapPopulations == "" {
			foreignObjectIDMapPopulations = "        // this is where foreign object ID map populations would appear (if required)"
		}

		if foreignObjectLoading == "" {
			foreignObjectLoading = "        // this is where any foreign object loading would appear (if required)"
		}

		selectData := fmt.Sprintf(
			selectFuncTemplate,
			structNamePlural,
			structNameSingular,
			structNameSingular,
			structNameSingular,
			table.Name,
			foreignObjectIDMapDeclarations,
			structNameSingular,
			structNameSingular,
			postprocessData,
			foreignObjectIDMapPopulations,
			foreignObjectLoading,
		)

		selectData += fmt.Sprintf(
			genericSelectFuncTemplate,
			structNamePlural,
			structNamePlural,
		)

		selectData += fmt.Sprintf(
			deserializeFuncTemplate,
			structNameSingular,
			structNameSingular,
		)

		tableFileData += selectData

		structData += "}\n\n"

		receiver := strings.ToLower(structNameSingular)[0:1]

		if hasPrimaryKey {
			structData += fmt.Sprintf(
				insertFuncTemplate,
				receiver,
				structNameSingular,
				structNameSingular,
				tableName,
				structNameSingular,
				receiver,
				receiver,
			)
		} else {
			structData += fmt.Sprintf(
				insertFuncTemplateNotImplemented,
				receiver,
				structNameSingular,
			)
		}

		structData += fmt.Sprintf(
			genericInsertFuncTemplate,
			structNameSingular,
		)

		if hasPrimaryKey {
			structData += fmt.Sprintf(
				primaryKeyFuncTemplate,
				receiver,
				structNameSingular,
				receiver,
				primaryKeyStructField,
				receiver,
				structNameSingular,
				receiver,
				primaryKeyStructField,
				primaryKeyStructType,
			)
		} else {
			structData += fmt.Sprintf(
				primaryKeyFuncTemplateNotImplemented,
				receiver,
				structNameSingular,
				receiver,
				structNameSingular,
			)
		}

		if hasPrimaryKey {
			structData += fmt.Sprintf(
				updateFuncTemplate,
				receiver,
				structNameSingular,
				structNameSingular,
				receiver,
				primaryKeyStructField,
				structNameSingular,
				receiver,
				receiver,
			)
		} else {
			structData += fmt.Sprintf(
				updateFuncTemplateNotImplemented,
				receiver,
				structNameSingular,
			)
		}

		structData += fmt.Sprintf(
			genericUpdateFuncTemplate,
			structNameSingular,
		)

		if hasPrimaryKey {
			structData += fmt.Sprintf(
				deleteFuncTemplate,
				receiver,
				structNameSingular,
				tableName,
				table.PrimaryKeyColumn.Name,
				receiver,
				primaryKeyStructField,
				receiver,
			)
		} else {
			structData += fmt.Sprintf(
				deleteFuncTemplateNotImplemented,
				receiver,
				structNameSingular,
			)
		}

		structData += fmt.Sprintf(
			genericDeleteFuncTemplate,
			structNameSingular,
		)

		tableFileData += structData

		tableFileData = headerTemplate + tableFileData

		fullOutputPath := filepath.Join(tempOutputPath, fmt.Sprintf("table_%v.go", table.Name))

		tableFileData = strings.ReplaceAll(tableFileData, "package templates", fmt.Sprintf("package %v", packageName))

		err = os.WriteFile(
			fullOutputPath,
			[]byte(tableFileData),
			0o777,
		)
		if err != nil {
			return err
		}
	}

	header := fmt.Sprintf(
		fileDataHeaderTemplate,
		packageName,
		fmt.Sprintf("`\n%v\n`", _helpers.UnsafeJSONPrettyFormat(tableByName)),
		fmt.Sprintf("var TableNames = %#+v\n\n", tableNames),
		strings.TrimSpace(initData),
	)

	fileData = header + fileData

	fullOutputPath := filepath.Join(tempOutputPath, "0_common.go")

	fileData = strings.ReplaceAll(fileData, "package templates", fmt.Sprintf("package %v", packageName))

	err = os.WriteFile(
		fullOutputPath,
		[]byte(fileData),
		0o777,
	)
	if err != nil {
		return err
	}

	fullOutputPath = filepath.Join(tempOutputPath, "0_server.go")

	fileData = strings.ReplaceAll(serverTemplate, "package templates", fmt.Sprintf("package %v", packageName))

	err = os.WriteFile(
		fullOutputPath,
		[]byte(fileData),
		0o777,
	)
	if err != nil {
		return err
	}

	fullOutputPath = filepath.Join(tempOutputPath, "0_stream.go")

	fileData = strings.ReplaceAll(streamTemplate, "package templates", fmt.Sprintf("package %v", packageName))

	err = os.WriteFile(
		fullOutputPath,
		[]byte(fileData),
		0o777,
	)
	if err != nil {
		return err
	}

	command1Ctx, command1Cancel := context.WithTimeout(ctx, time.Second*10)
	defer command1Cancel()
	out, err := exec.CommandContext(command1Ctx, "goimports", "-w", tempOutputPath).CombinedOutput()
	if err != nil {
		logger.Printf("error: goimports failed; err: %v; output follows:\n\n%v", err, string(out))
		return err
	}

	command2Ctx, command2Cancel := context.WithTimeout(ctx, time.Second*10)
	defer command2Cancel()
	out, err = exec.CommandContext(command2Ctx, "gofmt", tempOutputPath).CombinedOutput()
	if err != nil {
		logger.Printf("error: gofmt failed; err: %v; output follows:\n\n%v", err, string(out))
		return err
	}

	command3Ctx, command3Cancel := context.WithTimeout(ctx, time.Second*10)
	defer command3Cancel()
	out, err = exec.CommandContext(command3Ctx, "go", "vet", tempOutputPath).CombinedOutput()
	if err != nil {
		logger.Printf("error: go vet failed; err: %v; output follows:\n\n%v", err, string(out))
		return err
	}

	logger.Printf("moving %v to %v...", tempOutputPath, finalOutputPath)
	_ = os.RemoveAll(finalOutputPath)
	err = os.Rename(tempOutputPath, finalOutputPath)
	if err != nil {
		return err
	}
	_ = os.RemoveAll(tempOutputPath)
	logger.Printf("done.")

	return nil
}

func Run(ctx context.Context) error {
	relativeOutputPath := ""
	if len(os.Args) > 2 {
		relativeOutputPath = os.Args[2]
	}

	if relativeOutputPath == "" {
		return fmt.Errorf("first argument must be output path for folder to put generated templates in")
	}

	logger.Printf("resolving %v...", relativeOutputPath)
	finalOutputPath, err := filepath.Abs(relativeOutputPath)
	if err != nil {
		return err
	}
	logger.Printf("done.")

	tempOutputPath := fmt.Sprintf("%v_temp", finalOutputPath)

	err = run(ctx, finalOutputPath, tempOutputPath)
	if err != nil {
		logger.Printf("templating failed: %v\n\ncheck out %v to see where we got up to", err, tempOutputPath)
	}

	return nil
}
