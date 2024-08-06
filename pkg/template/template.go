package template

import (
	"bytes"
	"fmt"
	"go/format"
	"log"
	"regexp"
	"slices"
	"strings"

	"text/template"

	"github.com/chanced/caps"
	_pluralize "github.com/gertd/go-pluralize"
	"github.com/initialed85/djangolang/pkg/introspect"
	"github.com/initialed85/djangolang/pkg/types"
	"golang.org/x/exp/maps"

	"github.com/initialed85/djangolang/pkg/model_reference"
)

var (
	pluralize = _pluralize.NewClient()
)

func init() {
	converter, ok := caps.DefaultConverter.(caps.StdConverter)
	if !ok {
		panic(fmt.Sprintf("failed to cast %#+v to caps.StdConverter", caps.DefaultConverter))
	}

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

func Template(
	tableByName map[string]*introspect.Table,
	modulePath string,
	packageName string,
) (map[string]string, error) {
	templateDataByFileName := make(map[string]string)

	cmdMainFileData := model_reference.CmdMainFileData

	cmdMainFileData = strings.ReplaceAll(cmdMainFileData, "github.com/initialed85/djangolang/pkg/model_reference", fmt.Sprintf("%v/pkg/model_reference", modulePath))
	cmdMainFileData = strings.ReplaceAll(cmdMainFileData, "model_reference", packageName)

	templateDataByFileName["cmd/main.go"] = cmdMainFileData

	templateDataByFileName["0_meta.go"] = strings.ReplaceAll(
		model_reference.BaseFileData,
		"package model_reference",
		fmt.Sprintf("package %s", packageName),
	)

	templateDataByFileName["0_app.go"] = strings.ReplaceAll(
		model_reference.AppFileData,
		"package model_reference",
		fmt.Sprintf("package %s", packageName),
	)

	tableNames := maps.Keys(tableByName)
	slices.Sort(tableNames)

	for _, tableName := range tableNames {
		table := tableByName[tableName]

		// TODO: should probably be configurable
		if tableName == "schema_migrations" {
			continue
		}

		// TODO: add support for views
		if table.RelKind == "v" {
			continue
		}

		// TODO: add support for tables without primary keys
		if table.PrimaryKeyColumn == nil {
			log.Printf("warning: skipping table %s because it has no primary key", tableName)
			continue
		}

		// TODO: a factory or something
		intermediateData := model_reference.ReferenceFileData

		// TODO: a factory that doesn't re-parse every time
		parseTasks, err := Parse()
		if err != nil {
			return nil, err
		}

		getBaseVariables := func() map[string]string {

			baseVariables := map[string]string{
				"PackageName":  packageName,
				"ObjectName":   pluralize.Singular(caps.ToCamel(tableName)),
				"TableName":    tableName,
				"EndpointName": pluralize.Plural(caps.ToKebab(tableName)),
			}

			return baseVariables
		}

		getParseTask := func(name string) (ParseTask, error) {
			for _, parseTask := range parseTasks {
				if parseTask.Name == name {
					return parseTask, nil
				}
			}

			return ParseTask{}, fmt.Errorf("parse task %#+v unknown", name)
		}

		templateParseTask := func(parseTaskName string) error {
			parseTask, err := getParseTask(parseTaskName)
			if err != nil {
				return nil
			}

			replacedFragment := bytes.NewBufferString("")

			//
			// start
			//

			startTmpl, err := template.New(tableName).Option("missingkey=error").Parse(parseTask.ReplacedStartMatch)
			if err != nil {
				return fmt.Errorf("template.New (startTmpl) failed: %v", err)
			}

			startVariables := getBaseVariables()

			err = startTmpl.Execute(replacedFragment, startVariables)
			if err != nil {
				return fmt.Errorf("startTmpl.Execute failed: %v", err)
			}

			//
			// keep
			//

			keepVariables := getBaseVariables()

			keepVariables["PrimaryKeyColumnName"] = caps.ToCamel(table.PrimaryKeyColumn.Name)

			if parseTask.KeepIsPerColumn {
				for _, column := range table.Columns {
					if parseTask.KeepIsForPrimaryKeyOnly && !column.IsPrimaryKey {
						continue
					}

					if parseTask.KeepIsForNonPrimaryKeyOnly && column.IsPrimaryKey {
						continue
					}

					if parseTask.KeepIsForForeignKeysOnly && column.ForeignColumn == nil {
						continue
					}

					keepTmpl, err := template.New(tableName).Option("missingkey=error").Parse(parseTask.ReplacedKeepMatch)
					if err != nil {
						return fmt.Errorf("template.New (keepTmpl in parseTask.KeepIsPerColumn) failed: %v", err)
					}

					repeaterReplacedFragment := bytes.NewBufferString("")

					keepVariables["StructField"] = caps.ToCamel(column.Name)

					typeTemplate := column.TypeTemplate
					if !column.NotNull && !strings.HasPrefix(column.TypeTemplate, "*") && column.TypeTemplate != "any" {
						typeTemplate = fmt.Sprintf("*%v", typeTemplate)
					}

					keepVariables["TypeTemplate"] = typeTemplate

					keepVariables["TypeTemplateWithoutPointer"] = strings.TrimLeft(typeTemplate, "*")

					keepVariables["ColumnName"] = column.Name

					if column.ForeignColumn != nil {
						keepVariables["SelectFunc"] = fmt.Sprintf(
							"Select%v",
							caps.ToCamel(pluralize.Singular(column.ForeignColumn.TableName)),
						)
					}

					wrapColumnNameWithTypeCastInGeoJSON := false

					dataType := column.DataType
					dataType = strings.Trim(strings.Split(dataType, "(")[0], `"`)

					theType, err := types.GetTypeMetaForDataType(strings.TrimLeft(dataType, "*"))
					if err != nil {
						return err
					}

					keepVariables["ParseFunc"] = theType.ParseFuncTemplate

					keepVariables["IsZeroFunc"] = theType.IsZeroFuncTemplate

					keepVariables["FormatFunc"] = theType.FormatFuncTemplate

					structFieldAssignmentRef := ""
					if !column.NotNull {
						structFieldAssignmentRef = "&"
					}

					keepVariables["StructFieldAssignmentRef"] = structFieldAssignmentRef

					if wrapColumnNameWithTypeCastInGeoJSON {
						keepVariables["ColumnNameWithTypeCast"] = fmt.Sprintf(
							`ST_AsGeoJSON("%v"::geometry)::jsonb AS %v`,
							column.Name,
							column.Name,
						)
					} else {
						keepVariables["ColumnNameWithTypeCast"] = fmt.Sprintf(
							`"%v" AS %v`,
							column.Name,
							column.Name,
						)
					}

					if column.ForeignColumn != nil && parseTask.Name == "SelectLoadForeignObjects" {
						keepVariables["ForeignPrimaryKeyColumnVariable"] = fmt.Sprintf("%sTablePrimaryKeyColumn", pluralize.Singular(caps.ToCamel(column.ForeignTable.Name)))
					}

					err = keepTmpl.Execute(repeaterReplacedFragment, keepVariables)
					if err != nil {
						return err
					}

					if column.ForeignColumn != nil &&
						(parseTask.Name == "StructDefinition" ||
							parseTask.Name == "ReloadSetFields") {
						keepVariables["StructField"] += "Object"
						keepVariables["TypeTemplate"] = fmt.Sprintf("*%v", caps.ToCamel(pluralize.Singular(column.ForeignColumn.TableName)))
						keepVariables["ColumnName"] = fmt.Sprintf("%v_object", column.Name)

						err = keepTmpl.Execute(repeaterReplacedFragment, keepVariables)
						if err != nil {
							return fmt.Errorf("keepTmpl.Execute (in parseTask.KeepIsPerColumn) failed: %v", err)
						}
					}

					replacedFragment.Write(repeaterReplacedFragment.Bytes())
				}
			} else {
				keepVariables["DeleteSoftDelete"] = "/* soft-delete not applicable */"

				if parseTaskName == "DeleteSoftDelete" && slices.Contains(maps.Keys(table.ColumnByName), "deleted_at") {
					keepVariables["DeleteSoftDelete"] = parseTask.KeepMatch
				}

				keepTmpl, err := template.New(tableName).Option("missingkey=error").Parse(parseTask.ReplacedKeepMatch)
				if err != nil {
					return fmt.Errorf("template.New (keepTmpl in !parseTask.KeepIsPerColumn) failed: %v", err)
				}

				err = keepTmpl.Execute(replacedFragment, keepVariables)
				if err != nil {
					return fmt.Errorf("keepTmpl.Execute (in !parseTask.KeepIsPerColumn) failed: %v", err)
				}
			}

			//
			// end
			//

			endTmpl, err := template.New(tableName).Option("missingkey=error").Parse(parseTask.ReplacedEndMatch)
			if err != nil {
				return err
			}

			endVariables := getBaseVariables()

			err = endTmpl.Execute(replacedFragment, endVariables)
			if err != nil {
				return err
			}

			//
			// merge
			//

			parseTask.ReplacedFragment = replacedFragment.String()

			intermediateData = strings.Replace(
				intermediateData,
				parseTask.Fragment,
				parseTask.ReplacedFragment,
				1,
			)

			return nil
		}

		for _, parseTask := range parseTasks {
			err = templateParseTask(parseTask.Name)
			if err != nil {
				return nil, err
			}
		}

		baseTokenizeTasks := []TokenizeTask{
			{
				Find:    regexp.MustCompile(fmt.Sprintf(`package %v`, model_reference.ReferencePackageName)),
				Replace: "package {{ .PackageName }}",
			},
			{
				Find:    regexp.MustCompile(fmt.Sprintf(`%v`, model_reference.ReferenceTableName)),
				Replace: "{{ .TableName }}",
			},
			{
				Find:    regexp.MustCompile(fmt.Sprintf(`%v`, model_reference.ReferenceObjectName)),
				Replace: "{{ .ObjectName }}",
			},
			{
				Find:    regexp.MustCompile(fmt.Sprintf(`%v`, model_reference.ReferenceEndpointName)),
				Replace: "{{ .EndpointName }}",
			},
		}

		for _, tokenizeTask := range baseTokenizeTasks {
			intermediateData = tokenizeTask.Find.ReplaceAllString(intermediateData, tokenizeTask.Replace)
		}

		replacedIntermediateData := bytes.NewBufferString("")

		tmpl, err := template.New(tableName).Option("missingkey=error").Parse(intermediateData)
		if err != nil {
			return nil, err
		}

		err = tmpl.Execute(replacedIntermediateData, getBaseVariables())
		if err != nil {
			return nil, err
		}

		intermediateData = replacedIntermediateData.String()

		expr := regexp.MustCompile(`(?m)\s*//\s*.*$`)
		intermediateData = expr.ReplaceAllString(intermediateData, "")

		formatted, err := format.Source([]byte(intermediateData))
		if err != nil {
			for i, line := range strings.Split(intermediateData, "\n") {
				fmt.Printf("%v:\t %v\n", i, line)
			}
			log.Panicf("failed to format: %v", err)
		}
		intermediateData = string(formatted)

		templateDataByFileName[fmt.Sprintf("%v.go", tableName)] = intermediateData
	}

	return templateDataByFileName, nil
}
