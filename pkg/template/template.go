package template

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/format"
	_log "log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"text/template"

	"github.com/chanced/caps"
	_pluralize "github.com/gertd/go-pluralize"
	"github.com/initialed85/djangolang/pkg/helpers"
	"github.com/initialed85/djangolang/pkg/introspect"
	"github.com/initialed85/djangolang/pkg/types"
	"golang.org/x/exp/maps"

	"github.com/initialed85/djangolang/pkg/model_reference"
)

var log = helpers.GetLogger("template")

func ThisLogger() *_log.Logger {
	return log
}

var (
	pluralize = _pluralize.NewClient()
)

func init() {
	converter, ok := caps.DefaultConverter.(caps.StdConverter)
	if !ok {
		panic(fmt.Sprintf("failed to cast %#+v to caps.StdConverter", caps.DefaultConverter))
	}

	items := []string{
		"Dob",
		"Cpo",
		"Mwh",
		"Kwh",
		"Wh",
		"Json",
		"Jsonb",
		"Mac",
		"Ip",
		"M2m",
	}

	for _, item := range items {
		itemUpper := strings.ToUpper(item)
		converter.Set(item, itemUpper)
		converter.Set(strings.ToLower(item), itemUpper)
	}

	caps.DefaultConverter = converter
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

	tableByNameAsJSON, err := json.MarshalIndent(tableByName, "", "  ")
	if err != nil {
		return nil, err
	}

	templateDataByFileName["0_meta.go"] = strings.ReplaceAll(
		templateDataByFileName["0_meta.go"],
		"var tableByNameAsJSON = []byte(`{}`)",
		fmt.Sprintf("var tableByNameAsJSON = []byte(`%s`)", string(tableByNameAsJSON)),
	)

	templateDataByFileName["0_app.go"] = strings.ReplaceAll(
		model_reference.AppFileData,
		"package model_reference",
		fmt.Sprintf("package %s", packageName),
	)

	templateDataByFileName["0_app.go"] = strings.ReplaceAll(
		templateDataByFileName["0_app.go"],
		`helpers.GetLogger("model_reference")`,
		fmt.Sprintf("helpers.GetLogger(\"%s\")", packageName),
	)

	tableNames := maps.Keys(tableByName)
	slices.Sort(tableNames)

	for i, tableName := range tableNames {
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
				"PackageName":          packageName,
				"ObjectName":           pluralize.Singular(caps.ToCamel(tableName)),
				"ObjectNamePlural":     pluralize.Plural(caps.ToCamel(tableName)),
				"TableName":            tableName,
				"EndpointName":         pluralize.Plural(caps.ToKebab(tableName)),
				"EndpointNameSingular": pluralize.Singular(caps.ToKebab(tableName)),
				"NamespaceID":          fmt.Sprintf("1337 + %d", i+1),
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
			keepVariables["PrimaryKeyStructField"] = caps.ToCamel(table.PrimaryKeyColumn.Name)

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

					if typeTemplate != "any" {
						keepVariables["TypeTemplateWithoutPointer"] = fmt.Sprintf(`temp2, ok := temp1.(%s)`, strings.TrimLeft(typeTemplate, "*"))
					} else {
						keepVariables["TypeTemplateWithoutPointer"] = "temp2, ok := temp1, true"
					}

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

					if column.ForeignColumn != nil && (parseTask.Name == "SelectLoadForeignObjects" || parseTask.Name == "SelectLoadReferencedByObjects") {
						keepVariables["ForeignPrimaryKeyColumnVariable"] = fmt.Sprintf("%sTablePrimaryKeyColumn", pluralize.Singular(caps.ToCamel(column.ForeignTable.Name)))
						keepVariables["ForeignPrimaryKeyTableVariable"] = fmt.Sprintf("%sTable", pluralize.Singular(caps.ToCamel(column.TableName)))
						keepVariables["ForeignObjectName"] = pluralize.Singular(caps.ToCamel(column.ForeignTable.Name))
					}

					keepVariables["NotNull"] = fmt.Sprintf("%v", column.NotNull)
					keepVariables["HasDefault"] = fmt.Sprintf("%v", column.HasDefault)

					err = keepTmpl.Execute(repeaterReplacedFragment, keepVariables)
					if err != nil {
						return err
					}

					if column.ForeignColumn != nil {
						if parseTask.Name == "StructDefinition" || parseTask.Name == "ReloadSetFields" {
							keepVariables["StructField"] += "Object"
							keepVariables["TypeTemplate"] = fmt.Sprintf("*%v", caps.ToCamel(pluralize.Singular(column.ForeignColumn.TableName)))
							keepVariables["ColumnName"] = fmt.Sprintf("%v_object", column.Name)

							err = keepTmpl.Execute(repeaterReplacedFragment, keepVariables)
							if err != nil {
								return fmt.Errorf("keepTmpl.Execute (in parseTask.KeepIsPerColumn) failed: %v", err)
							}
						}
					}

					replacedFragment.Write(repeaterReplacedFragment.Bytes())
				}

				if parseTask.Name == "StructDefinition" || parseTask.Name == "ReloadSetFields" {
					for _, foreignColumn := range table.ReferencedByColumns {
						repeaterReplacedFragment := bytes.NewBufferString("")

						keepVariables["StructField"] = fmt.Sprintf("ReferencedBy%s%sObjects", caps.ToCamel(pluralize.Singular(foreignColumn.TableName)), caps.ToCamel(foreignColumn.Name))
						keepVariables["TypeTemplate"] = fmt.Sprintf("[]*%v", caps.ToCamel(pluralize.Singular(foreignColumn.TableName)))
						keepVariables["ColumnName"] = fmt.Sprintf("referenced_by_%s_%s_objects", pluralize.Singular(foreignColumn.TableName), foreignColumn.Name)

						keepTmpl, err := template.New(tableName).Option("missingkey=error").Parse(parseTask.ReplacedKeepMatch)
						if err != nil {
							return fmt.Errorf("template.New (keepTmpl in parseTask.ReferencedByTables) failed: %v", err)
						}

						err = keepTmpl.Execute(repeaterReplacedFragment, keepVariables)
						if err != nil {
							return fmt.Errorf("keepTmpl.Execute (in parseTask.ReferencedByTables) failed: %v", err)
						}

						replacedFragment.Write(repeaterReplacedFragment.Bytes())
					}
				}
			} else {
				keepVariables["DeleteSoftDelete"] = "/* soft-delete not applicable */"

				if parseTaskName == "DeleteSoftDelete" && slices.Contains(maps.Keys(table.ColumnByName), "deleted_at") {
					keepVariables["DeleteSoftDelete"] = parseTask.KeepMatch
				}

				if parseTaskName == "SelectLoadReferencedByObjects" {
					for _, foreignColumn := range table.ReferencedByColumns {
						repeaterReplacedFragment := bytes.NewBufferString("")

						keepVariables["StructField"] = fmt.Sprintf("ReferencedBy%s%sObjects", caps.ToCamel(pluralize.Singular(foreignColumn.TableName)), caps.ToCamel(foreignColumn.Name))
						keepVariables["TypeTemplate"] = fmt.Sprintf("[]*%v", caps.ToCamel(pluralize.Singular(foreignColumn.TableName)))
						keepVariables["ColumnName"] = fmt.Sprintf("referenced_by_%s_%s_objects", pluralize.Singular(foreignColumn.TableName), foreignColumn.Name)
						keepVariables["SelectFunc"] = fmt.Sprintf(
							"Select%v",
							caps.ToCamel(pluralize.Plural(foreignColumn.TableName)),
						)
						keepVariables["ForeignPrimaryKeyColumnVariable"] = fmt.Sprintf(
							"%sTable%sColumn",
							caps.ToCamel(pluralize.Singular(foreignColumn.TableName)),
							caps.ToCamel(pluralize.Singular(foreignColumn.Name)),
						)
						keepVariables["ForeignPrimaryKeyTableVariable"] = fmt.Sprintf("%sTable", pluralize.Singular(caps.ToCamel(foreignColumn.TableName)))
						keepVariables["ForeignObjectName"] = pluralize.Singular(caps.ToCamel(foreignColumn.TableName))

						keepTmpl, err := template.New(tableName).Option("missingkey=error").Parse(parseTask.ReplacedKeepMatch)
						if err != nil {
							return fmt.Errorf("template.New (keepTmpl in parseTask.SelectLoadReferencedByObjects) failed: %v", err)
						}

						err = keepTmpl.Execute(repeaterReplacedFragment, keepVariables)
						if err != nil {
							return fmt.Errorf("keepTmpl.Execute (in parseTask.SelectLoadReferencedByObjects) failed: %v", err)
						}

						if foreignColumn.TableName == table.Name {
							replacedFragment.WriteString("\n/*")
						}

						replacedFragment.Write(repeaterReplacedFragment.Bytes())

						if foreignColumn.TableName == table.Name {
							replacedFragment.WriteString("*/\n")
						}
					}
				} else {
					keepTmpl, err := template.New(tableName).Option("missingkey=error").Parse(parseTask.ReplacedKeepMatch)
					if err != nil {
						return fmt.Errorf("template.New (keepTmpl in !parseTask.KeepIsPerColumn) failed: %v", err)
					}

					err = keepTmpl.Execute(replacedFragment, keepVariables)
					if err != nil {
						return fmt.Errorf("keepTmpl.Execute (in !parseTask.KeepIsPerColumn) failed: %v", err)
					}
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
			// safe ones
			{
				Find:    regexp.MustCompile(fmt.Sprintf(`package %v`, model_reference.ReferencePackageName)),
				Replace: "package {{ .PackageName }}",
			},
			{
				Find:    regexp.MustCompile(fmt.Sprintf(`"%v"`, model_reference.ReferenceTableName)),
				Replace: `"{{ .TableName }}"`,
			},
			{
				Find:    regexp.MustCompile(fmt.Sprintf(`%v`, model_reference.ReferenceEndpointName)),
				Replace: "{{ .EndpointName }}",
			},
			{
				Find:    regexp.MustCompile(fmt.Sprintf(`claim-%v`, model_reference.ReferenceEndpointNameSingular)),
				Replace: "claim-{{ .EndpointNameSingular }}",
			},
			{
				Find:    regexp.MustCompile(fmt.Sprintf(`%vIntrospectedTable`, model_reference.ReferenceObjectName)),
				Replace: "{{ .ObjectName }}IntrospectedTable",
			},
			{
				Find:    regexp.MustCompile(fmt.Sprintf(`tableByName\[%vTable\]`, model_reference.ReferenceObjectName)),
				Replace: "tableByName[{{ .ObjectName }}Table]",
			},
			{
				Find:    regexp.MustCompile(fmt.Sprintf(`%vManyPathParams`, model_reference.ReferenceObjectName)),
				Replace: "{{ .ObjectName }}ManyPathParams",
			},
			{
				Find:    regexp.MustCompile(fmt.Sprintf(`%vOnePathParams`, model_reference.ReferenceObjectName)),
				Replace: "{{ .ObjectName }}OnePathParams",
			},
			{
				Find:    regexp.MustCompile(fmt.Sprintf(`%vLoadQueryParams`, model_reference.ReferenceObjectName)),
				Replace: "{{ .ObjectName }}LoadQueryParams",
			},
			{
				Find:    regexp.MustCompile(`PrimaryKey uuid\.UUID`),
				Replace: fmt.Sprintf("PrimaryKey %s", table.PrimaryKeyColumn.TypeTemplate),
			},
			{
				Find:    regexp.MustCompile(`primaryKey uuid\.UUID`),
				Replace: fmt.Sprintf("primaryKey %s", table.PrimaryKeyColumn.TypeTemplate),
			},
			{
				Find:    regexp.MustCompile(`item map\[string\]any\) \(\[\]\*LogicalThing`),
				Replace: fmt.Sprintf("item map[string]any) ([]*%s", model_reference.ReferenceObjectName),
			},
			{
				Find:    regexp.MustCompile(`\[\]\*LogicalThing{object}`),
				Replace: "[]*{{ .ObjectName }}{object}",
			},
			{
				Find:    regexp.MustCompile(`object\.ID = pathParams.PrimaryKey`),
				Replace: fmt.Sprintf("object.%s = pathParams.PrimaryKey", caps.ToCamel(table.PrimaryKeyColumn.Name)),
			},
			{
				Find:    regexp.MustCompile(`object \*LogicalThing`),
				Replace: "object *{{ .ObjectName }}",
			},
			{
				Find:    regexp.MustCompile(`req LogicalThing`),
				Replace: "req {{ .ObjectName }}",
			},
			{
				Find:    regexp.MustCompile(`server\.Response\[LogicalThing\]`),
				Replace: "server.Response[{{ .ObjectName }}]",
			},
			{
				Find:    regexp.MustCompile(`objects \[\]\*LogicalThing`),
				Replace: "objects []*{{ .ObjectName }}",
			},
			{
				Find:    regexp.MustCompile(`req \[\]\*LogicalThing`),
				Replace: "req []*{{ .ObjectName }}",
			},
			{
				Find:    regexp.MustCompile(`var cachedObjects \[\]\*LogicalThing`),
				Replace: "var cachedObjects []*{{ .ObjectName }}",
			},

			// plurals first
			{
				Find:    regexp.MustCompile(fmt.Sprintf(`func Select%v`, model_reference.ReferenceObjectNamePlural)),
				Replace: "func Select{{ .ObjectNamePlural }}",
			},
			{
				Find:    regexp.MustCompile(fmt.Sprintf(`entered Select%v`, model_reference.ReferenceObjectNamePlural)),
				Replace: "entered Select{{ .ObjectNamePlural }}",
			},
			{
				Find:    regexp.MustCompile(fmt.Sprintf(`exited Select%v`, model_reference.ReferenceObjectNamePlural)),
				Replace: "exited Select{{ .ObjectNamePlural }}",
			},
			{
				Find:    regexp.MustCompile(`recursion limit reached for __this_function__`),
				Replace: "recursion limit reached for Select{{ .ObjectNamePlural }}",
			},
			{
				Find:    regexp.MustCompile(`loading __this_function__`),
				Replace: "loading Select{{ .ObjectNamePlural }}",
			},
			{
				Find:    regexp.MustCompile(`loaded __this_function__`),
				Replace: "loaded Select{{ .ObjectNamePlural }}",
			},
			{
				Find:    regexp.MustCompile(fmt.Sprintf(`objects, count, totalCount, page, totalPages, err := Select%v`, model_reference.ReferenceObjectNamePlural)),
				Replace: "objects, count, totalCount, page, totalPages, err := Select{{ .ObjectNamePlural }}",
			},
			{
				Find:    regexp.MustCompile(fmt.Sprintf(`objects, _, _, _, _, err := Select%v`, model_reference.ReferenceObjectNamePlural)),
				Replace: "objects, _, _, _, _, err := Select{{ .ObjectNamePlural }}",
			},
			{
				Find:    regexp.MustCompile(fmt.Sprintf(`handleGet%v`, model_reference.ReferenceObjectNamePlural)),
				Replace: "handleGet{{ .ObjectNamePlural }}",
			},
			{
				Find:    regexp.MustCompile(fmt.Sprintf(`ms, _, _, _, _, err := Select%v\(`, model_reference.ReferenceObjectNamePlural)),
				Replace: "ms, _, _, _, _, err := Select{{ .ObjectNamePlural }}(",
			},

			// singulars last
			{
				Find:    regexp.MustCompile(fmt.Sprintf(`Claim%v`, model_reference.ReferenceObjectName)),
				Replace: "Claim{{ .ObjectName }}",
			},
			{
				Find:    regexp.MustCompile(fmt.Sprintf(`object, count, totalCount, page, totalPages, err := Select%v`, model_reference.ReferenceObjectName)),
				Replace: "object, count, totalCount, page, totalPages, err := Select{{ .ObjectName }}",
			},
			{
				Find:    regexp.MustCompile(fmt.Sprintf(`x, _, _, _, _, err := Select%v\(`, model_reference.ReferenceObjectName)),
				Replace: "x, _, _, _, _, err := Select{{ .ObjectName }}(",
			},
			{
				Find:    regexp.MustCompile(fmt.Sprintf(`Objects, _, _, _, _, err = Select%v\(`, model_reference.ReferenceObjectName)),
				Replace: "Objects, _, _, _, _, err = Select{{ .ObjectName }}(",
			},
			{
				Find:    regexp.MustCompile(fmt.Sprintf(`_, _, _, _, err = Select%v\(`, model_reference.ReferenceObjectName)),
				Replace: "_, _, _, _, err = Select{{ .ObjectName }}(",
			},
			{
				Find:    regexp.MustCompile(fmt.Sprintf(`%vClaimRequest`, model_reference.ReferenceObjectName)),
				Replace: "{{ .ObjectName }}ClaimRequest",
			},
			{
				Find:    regexp.MustCompile(fmt.Sprintf(`object, count, totalCount, _, _, err = Select%v`, model_reference.ReferenceObjectName)),
				Replace: "object, count, totalCount, _, _, err = Select{{ .ObjectName }}",
			},
			{
				Find:    regexp.MustCompile(fmt.Sprintf(`o, _, _, _, _, err := Select%v`, model_reference.ReferenceObjectName)),
				Replace: "o, _, _, _, _, err := Select{{ .ObjectName }}",
			},
			{
				Find:    regexp.MustCompile(fmt.Sprintf(`var %vTable`, model_reference.ReferenceObjectName)),
				Replace: "var {{ .ObjectName }}Table",
			},
			{
				Find:    regexp.MustCompile(fmt.Sprintf(`%vTablePrimaryKeyColumn`, model_reference.ReferenceObjectName)),
				Replace: "{{ .ObjectName }}TablePrimaryKeyColumn",
			},
			{
				Find:    regexp.MustCompile(fmt.Sprintf(`%vTableColumnLookup`, model_reference.ReferenceObjectName)),
				Replace: "{{ .ObjectName }}TableColumnLookup",
			},
			{
				Find:    regexp.MustCompile(fmt.Sprintf(`%vTableColumns`, model_reference.ReferenceObjectName)),
				Replace: "{{ .ObjectName }}TableColumns",
			},
			{
				Find:    regexp.MustCompile(fmt.Sprintf(`%vTable,`, model_reference.ReferenceObjectName)),
				Replace: "{{ .ObjectName }}Table,",
			},
			{
				Find:    regexp.MustCompile(fmt.Sprintf(`m \*%v`, model_reference.ReferenceObjectName)),
				Replace: "m *{{ .ObjectName }}",
			},
			{
				Find:    regexp.MustCompile(fmt.Sprintf(`func Select%v`, model_reference.ReferenceObjectName)),
				Replace: "func Select{{ .ObjectName }}",
			},
			{
				Find:    regexp.MustCompile(fmt.Sprintf(`\(\[\]\*%v`, model_reference.ReferenceObjectName)),
				Replace: "([]*{{ .ObjectName }}",
			},
			{
				Find:    regexp.MustCompile(fmt.Sprintf(`\(context\.Context, \[\]\*%v`, model_reference.ReferenceObjectName)),
				Replace: "(context.Context, []*{{ .ObjectName }}",
			},
			{
				Find:    regexp.MustCompile(fmt.Sprintf(`&%v`, model_reference.ReferenceObjectName)),
				Replace: "&{{ .ObjectName }}",
			},
			{
				Find:    regexp.MustCompile(fmt.Sprintf(`\(\*%v`, model_reference.ReferenceObjectName)),
				Replace: "(*{{ .ObjectName }}",
			},
			{
				Find:    regexp.MustCompile(fmt.Sprintf(`\(context\.Context, \*%v`, model_reference.ReferenceObjectName)),
				Replace: "(context.Context, *{{ .ObjectName }}",
			},
			{
				Find:    regexp.MustCompile(fmt.Sprintf(`, \[\]\*%v`, model_reference.ReferenceObjectName)),
				Replace: ", []*{{ .ObjectName }}",
			},
			{
				Find:    regexp.MustCompile(fmt.Sprintf(`handleGet%v`, model_reference.ReferenceObjectName)),
				Replace: "handleGet{{ .ObjectName }}",
			},
			{
				Find:    regexp.MustCompile(fmt.Sprintf(`handlePost%v`, model_reference.ReferenceObjectName)),
				Replace: "handlePost{{ .ObjectName }}",
			},
			{
				Find:    regexp.MustCompile(fmt.Sprintf(`handlePut%v`, model_reference.ReferenceObjectName)),
				Replace: "handlePut{{ .ObjectName }}",
			},
			{
				Find:    regexp.MustCompile(fmt.Sprintf(`handlePatch%v`, model_reference.ReferenceObjectName)),
				Replace: "handlePatch{{ .ObjectName }}",
			},
			{
				Find:    regexp.MustCompile(fmt.Sprintf(`handleDelete%v`, model_reference.ReferenceObjectName)),
				Replace: "handleDelete{{ .ObjectName }}",
			},
			{
				Find:    regexp.MustCompile(fmt.Sprintf(`func Get%v`, model_reference.ReferenceObjectName)),
				Replace: "func Get{{ .ObjectName }}",
			},
			{
				Find:    regexp.MustCompile(fmt.Sprintf(`New%v`, model_reference.ReferenceObjectName)),
				Replace: "New{{ .ObjectName }}",
			},
			{
				Find:    regexp.MustCompile(fmt.Sprintf(`during (\w*)%v`, model_reference.ReferenceObjectName)),
				Replace: "during $1{{ .ObjectName }}",
			},
			{
				Find:    regexp.MustCompile(fmt.Sprintf(`call (\w*)%v`, model_reference.ReferenceObjectName)),
				Replace: "call $1{{ .ObjectName }}",
			},
			{
				Find:    regexp.MustCompile(fmt.Sprintf(`as %v`, model_reference.ReferenceObjectName)),
				Replace: "as {{ .ObjectName }}",
			},
			{
				Find:    regexp.MustCompile(fmt.Sprintf(`%v{},`, model_reference.ReferenceObjectName)),
				Replace: "{{ .ObjectName }}{},",
			},
			{
				Find:    regexp.MustCompile(fmt.Sprintf(`Get%v`, model_reference.ReferenceObjectName)),
				Replace: "Get{{ .ObjectName }}",
			},
			{
				Find:    regexp.MustCompile(fmt.Sprintf(`thatCtx, ok := query\.HandleQueryPathGraphCycles\(ctx, %vTable\)`, model_reference.ReferenceObjectName)),
				Replace: "thatCtx, ok := query.HandleQueryPathGraphCycles(ctx, {{ .ObjectName }}Table)",
			},
			{
				Find:    regexp.MustCompile(fmt.Sprintf(`%vTableNamespaceID`, model_reference.ReferenceObjectName)),
				Replace: `{{ .ObjectName }}TableNamespaceID`,
			},
			{
				Find:    regexp.MustCompile(fmt.Sprintf(`skipping Select%v`, model_reference.ReferenceObjectNamePlural)),
				Replace: "skipping Select{{ .ObjectName }}",
			},
			{
				Find:    regexp.MustCompile(fmt.Sprintf(`ShouldLoad\(ctx, %vTable\)`, model_reference.ReferenceObjectName)),
				Replace: "ShouldLoad(ctx, {{ .ObjectName }}Table)",
			},
			{
				Find:    regexp.MustCompile(fmt.Sprintf(`ShouldLoad\(ctx, fmt.Sprintf\("referenced_by_%%s", %vTable\)\)`, model_reference.ReferenceObjectName)),
				Replace: `ShouldLoad(ctx, fmt.Sprintf("referenced_by_%s", {{ .ObjectName }}Table))`,
			},
			{
				Find:    regexp.MustCompile(fmt.Sprintf(`MutateRouterFor%v`, model_reference.ReferenceObjectName)),
				Replace: "MutateRouterFor{{ .ObjectName }}",
			},
		}

		columnNames := maps.Keys(table.ColumnByName)

		if !(slices.Contains(columnNames, "claimed_until") && slices.Contains(columnNames, "claimed_by")) {
			baseTokenizeTasks = append(baseTokenizeTasks,
				TokenizeTask{
					Find:    regexp.MustCompile(`m\.ClaimedUntil = &until`),
					Replace: `/* m.ClaimedUntil = &until */`,
				},
				TokenizeTask{
					Find:    regexp.MustCompile(`m\.ClaimedBy = &by`),
					Replace: `/* m.ClaimedBy = &by */`,
				},
			)
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

		templateDataByFileName[fmt.Sprintf("%v.go", tableName)] = intermediateData
	}

	temp, err := os.CreateTemp("", "djangolang")
	if err != nil {
		return nil, err
	}

	tempFolder, _ := filepath.Split(temp.Name())
	tempFolder = filepath.Join(tempFolder, "djangolang")
	err = os.MkdirAll(tempFolder, 0o777)
	if err != nil {
		return nil, err
	}

	fileNameByTempFile := make(map[string]string)

	for fileName, templateData := range templateDataByFileName {
		tempFile := filepath.Join(tempFolder, fileName)

		tempFolder, _ := filepath.Split(tempFile)
		err = os.MkdirAll(tempFolder, 0o777)
		if err != nil {
			return nil, err
		}

		err = os.WriteFile(tempFile, []byte(templateData), 0o777)
		if err != nil {
			return nil, err
		}

		fileNameByTempFile[fileName] = tempFile
	}

	cmd := exec.Command("goimports", "-w", tempFolder)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("%v: %v", err, string(out))
	}

	for fileName, tempFile := range fileNameByTempFile {
		unusedImportsRemoved, err := os.ReadFile(tempFile)
		if err != nil {
			return nil, err
		}

		formatted, err := format.Source([]byte(unusedImportsRemoved))

		if err != nil {
			for i, line := range strings.Split(string(unusedImportsRemoved), "\n") {
				fmt.Printf("%v:\t %v\n", i, line)
			}
			log.Panicf("failed to format: %v", err)
		}

		templateDataByFileName[fileName] = string(formatted)
	}

	return templateDataByFileName, nil
}
