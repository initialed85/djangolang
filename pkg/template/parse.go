package template

import (
	"fmt"
	"regexp"
	"slices"
	"strings"

	"github.com/initialed85/djangolang/pkg/model_reference"
)

// TokenizeTask tells the parser how to tokenize
type TokenizeTask struct {
	Find    *regexp.Regexp
	Replace string
}

// ParseTask tells the parser how to parse, *Expr must have no capture groups (meaning keep the whole match) or
// exactly one capture group (meaning keep only this); you'll probably want the single-capture-group option for
// repeating fields (e.g. structs)
//
// NOTE: There's a big gotcha in that TokenizeTask.Replace must only be able to execute successfully once, if this
// is not the case, then TokenizeTask.Replace will get applied each time that TokenizeTask.Find succeeds (and
// you'll end up with heaps of duplicated TokenizeTask.Replace instances)
type ParseTask struct {
	Name          string
	StartExpr     *regexp.Regexp
	KeepExpr      *regexp.Regexp
	EndExpr       *regexp.Regexp
	TokenizeTasks []TokenizeTask

	KeepIsPerColumn            bool
	KeepIsForPrimaryKeyOnly    bool
	KeepIsForNonPrimaryKeyOnly bool
	KeepIsForForeignKeysOnly   bool
	KeepIsForSoftDeletableOnly bool
	KeepIsForReferencedByOnly  bool
	KeepIsForClaimOnly         bool
	StartMatch                 string
	KeepMatch                  string
	EndMatch                   string
	Fragment                   string
	ReplacedStartMatch         string
	ReplacedKeepMatch          string
	ReplacedEndMatch           string
	ReplacedFragment           string
	StartVariableNameSet       map[string]struct{}
	KeepVariableNameSet        map[string]struct{}
	EndVariableNameSet         map[string]struct{}
}

func getParseTasks() []ParseTask {
	parseTasks := []ParseTask{
		{
			Name:      "NamespaceID",
			StartExpr: regexp.MustCompile(`(?ms)^[ |\t]*var LogicalThingTableNamespaceID int32 = `),
			KeepExpr:  regexp.MustCompile(`(?msU)(\d+) // LogicalThingTableNamespaceID$`),
			EndExpr:   regexp.MustCompile(`(?msU)// LogicalThingTableNamespaceID$\n`),
			TokenizeTasks: []TokenizeTask{
				{
					Find:    regexp.MustCompile(`1337`),
					Replace: "{{ .NamespaceID }}",
				},
			},
			KeepIsPerColumn:            false,
			KeepIsForPrimaryKeyOnly:    false,
			KeepIsForNonPrimaryKeyOnly: false,
			KeepIsForForeignKeysOnly:   false,
			KeepIsForSoftDeletableOnly: false,
			KeepIsForReferencedByOnly:  false,
		},

		{
			Name:      "StructDefinition",
			StartExpr: regexp.MustCompile(`(?ms)^[ |\t]*type LogicalThing struct \{$\n`),
			KeepExpr:  regexp.MustCompile(`(?msU)^\s*ID\s+uuid.UUID\s+\x60json:"id"\x60\s*$\n`),
			EndExpr:   regexp.MustCompile(`(?msU)^}$\n`),
			TokenizeTasks: []TokenizeTask{
				{
					Find:    regexp.MustCompile(`ID`),
					Replace: "{{ .StructField }}",
				},
				{
					Find:    regexp.MustCompile(`uuid\.UUID`),
					Replace: "{{ .TypeTemplate }}",
				},
				{
					Find:    regexp.MustCompile(`id`),
					Replace: "{{ .ColumnName }}",
				},
			},
			KeepIsPerColumn:            true,
			KeepIsForPrimaryKeyOnly:    false,
			KeepIsForNonPrimaryKeyOnly: false,
			KeepIsForForeignKeysOnly:   false,
			KeepIsForSoftDeletableOnly: false,
			KeepIsForReferencedByOnly:  false,
		},

		{
			Name:      "ColumnVariables",
			StartExpr: regexp.MustCompile(`(?ms)^[ |\t]*var \( // ColumnVariables$\n`),
			KeepExpr:  regexp.MustCompile(`(?msU)^\s*LogicalThingTableIDColumn\s*=\s*"id"$\n`),
			EndExpr:   regexp.MustCompile(`(?msU)^\)$\n`),
			TokenizeTasks: []TokenizeTask{
				{
					Find:    regexp.MustCompile(`LogicalThingTableIDColumn`),
					Replace: "LogicalThingTable{{ .StructField }}Column",
				},
				{
					Find:    regexp.MustCompile(`"id"`),
					Replace: `"{{ .ColumnName }}"`,
				},
			},
			KeepIsPerColumn:            true,
			KeepIsForPrimaryKeyOnly:    false,
			KeepIsForNonPrimaryKeyOnly: false,
			KeepIsForForeignKeysOnly:   false,
			KeepIsForSoftDeletableOnly: false,
			KeepIsForReferencedByOnly:  false,
		},

		{
			Name:      "ColumnVariablesWithTypeCasts",
			StartExpr: regexp.MustCompile(`(?ms)^[ |\t]*var \( // ColumnVariablesWithTypeCasts$\n`),
			KeepExpr:  regexp.MustCompile(`(?msU)^\s*LogicalThingTableIDColumnWithTypeCast\s*=\s*.*\x60"id" AS id\x60$\n`),
			EndExpr:   regexp.MustCompile(`(?msU)^\)$\n`),
			TokenizeTasks: []TokenizeTask{
				{
					Find:    regexp.MustCompile(`LogicalThingTableIDColumn`),
					Replace: "LogicalThingTable{{ .StructField }}Column",
				},
				{
					Find:    regexp.MustCompile(`"id" AS id`),
					Replace: `{{ .ColumnNameWithTypeCast }}`,
				},
			},
			KeepIsPerColumn:            true,
			KeepIsForPrimaryKeyOnly:    false,
			KeepIsForNonPrimaryKeyOnly: false,
			KeepIsForForeignKeysOnly:   false,
			KeepIsForSoftDeletableOnly: false,
			KeepIsForReferencedByOnly:  false,
		},

		{
			Name:      "ColumnSlice",
			StartExpr: regexp.MustCompile(`(?ms)^[ |\t]*var LogicalThingTableColumns = \[\]string\{$\n`),
			KeepExpr:  regexp.MustCompile(`(?msU)^\s*LogicalThingTableIDColumn,\s*$\n`),
			EndExpr:   regexp.MustCompile(`(?msU)^}$\n`),
			TokenizeTasks: []TokenizeTask{
				{
					Find:    regexp.MustCompile(`LogicalThingTableIDColumn`),
					Replace: "LogicalThingTable{{ .StructField }}Column",
				},
			},
			KeepIsPerColumn:            true,
			KeepIsForPrimaryKeyOnly:    false,
			KeepIsForNonPrimaryKeyOnly: false,
			KeepIsForForeignKeysOnly:   false,
			KeepIsForSoftDeletableOnly: false,
			KeepIsForReferencedByOnly:  false,
		},

		{
			Name:      "ColumnWithTypeCastSlice",
			StartExpr: regexp.MustCompile(`(?ms)^[ |\t]*var LogicalThingTableColumnsWithTypeCasts = \[\]string\{$\n`),
			KeepExpr:  regexp.MustCompile(`(?msU)^\s*LogicalThingTableIDColumnWithTypeCast,\s*$\n`),
			EndExpr:   regexp.MustCompile(`(?msU)^}$\n`),
			TokenizeTasks: []TokenizeTask{
				{
					Find:    regexp.MustCompile(`LogicalThingTableIDColumn`),
					Replace: "LogicalThingTable{{ .StructField }}Column",
				},
			},
			KeepIsPerColumn:            true,
			KeepIsForPrimaryKeyOnly:    false,
			KeepIsForNonPrimaryKeyOnly: false,
			KeepIsForForeignKeysOnly:   false,
			KeepIsForSoftDeletableOnly: false,
			KeepIsForReferencedByOnly:  false,
		},

		{
			Name:      "PrimaryKeyColumn",
			StartExpr: regexp.MustCompile(`(?ms)^[ |\t]*var \( // PrimaryKeyColumn$\n`),
			KeepExpr:  regexp.MustCompile(`(?msU)^\s*LogicalThingTablePrimaryKeyColumn\s*=\s*LogicalThingTableIDColumn$\n`),
			EndExpr:   regexp.MustCompile(`(?msU)^\)$\n`),
			TokenizeTasks: []TokenizeTask{
				{
					Find:    regexp.MustCompile(`ID`),
					Replace: "{{ .PrimaryKeyColumnName }}",
				},
			},
			KeepIsPerColumn:            false,
			KeepIsForPrimaryKeyOnly:    true,
			KeepIsForNonPrimaryKeyOnly: false,
			KeepIsForForeignKeysOnly:   false,
			KeepIsForSoftDeletableOnly: false,
			KeepIsForReferencedByOnly:  false,
		},

		{
			Name:      "PrimaryKeyGetter",
			StartExpr: regexp.MustCompile(`(?ms)^[ |\t]*func \(m \*LogicalThing\) GetPrimaryKeyValue\(\) any \{$\n`),
			KeepExpr:  regexp.MustCompile(`(?msU)^\s*return m\.ID\s*$\n`),
			EndExpr:   regexp.MustCompile(`(?msU)^}$\n`),
			TokenizeTasks: []TokenizeTask{
				{
					Find:    regexp.MustCompile(`m \*LogicalThing`),
					Replace: "m *{{ .ObjectName }}",
				},
				{
					Find:    regexp.MustCompile(`m\.ID`),
					Replace: "m.{{ .PrimaryKeyColumnName }}",
				},
			},
			KeepIsPerColumn:            false,
			KeepIsForPrimaryKeyOnly:    false,
			KeepIsForNonPrimaryKeyOnly: false,
			KeepIsForForeignKeysOnly:   false,
			KeepIsForSoftDeletableOnly: false,
			KeepIsForReferencedByOnly:  false,
		},

		{
			Name:      "FromItemTypeSwitch",
			StartExpr: regexp.MustCompile(`(?ms)^[ |\t]*switch k {$\n`),
			KeepExpr:  regexp.MustCompile(`(?ms)^[ |\t]*case "id":$\n.*m\.ID = temp2$\n\n`),
			EndExpr:   regexp.MustCompile(`(?msU)^\s*}$\n^\s*}\s*return nil$\n`),
			TokenizeTasks: []TokenizeTask{
				{
					Find:    regexp.MustCompile(`id`),
					Replace: "{{ .ColumnName }}",
				},
				{
					Find:    regexp.MustCompile(`m\.ID`),
					Replace: "m.{{ .StructField }}",
				},
				{
					Find:    regexp.MustCompile(`types\.ParseUUID\(v\)`),
					Replace: "{{ .ParseFunc }}",
				},
				{
					Find:    regexp.MustCompile(`temp2, ok := temp1\.\(uuid\.UUID\)`),
					Replace: "{{ .TypeTemplateWithoutPointer }}",
				},
				{
					Find:    regexp.MustCompile(` = temp2`),
					Replace: "= {{ .StructFieldAssignmentRef }}temp2",
				},
			},
			KeepIsPerColumn:            true,
			KeepIsForPrimaryKeyOnly:    false,
			KeepIsForNonPrimaryKeyOnly: false,
			KeepIsForForeignKeysOnly:   false,
			KeepIsForSoftDeletableOnly: false,
			KeepIsForReferencedByOnly:  false,
		},

		{
			Name:      "ReloadSetFields",
			StartExpr: regexp.MustCompile(`(?ms)^[ |\t]*// <reload-set-fields>$\n`),
			KeepExpr:  regexp.MustCompile(`(?msU)\n^\s*m\.ID = o\.ID$`),
			EndExpr:   regexp.MustCompile(`(?msU)^\s*// </reload-set-fields>$\n`),
			TokenizeTasks: []TokenizeTask{
				{
					Find:    regexp.MustCompile(`m\.ID`),
					Replace: "m.{{ .StructField }}",
				},
				{
					Find:    regexp.MustCompile(`o\.ID`),
					Replace: "o.{{ .StructField }}",
				},
			},
			KeepIsPerColumn:            true,
			KeepIsForPrimaryKeyOnly:    false,
			KeepIsForNonPrimaryKeyOnly: false,
			KeepIsForForeignKeysOnly:   false,
			KeepIsForSoftDeletableOnly: false,
			KeepIsForReferencedByOnly:  false,
		},

		{
			Name:      "SelectLoadForeignObjects",
			StartExpr: regexp.MustCompile(`(?ms)^[ |\t]*// <select-load-foreign-objects>$\n`),
			KeepExpr:  regexp.MustCompile(`(?msU)^\s*// <select-load-foreign-object>(\s+.*\n)\s+// </select-load-foreign-object>$\n`),
			EndExpr:   regexp.MustCompile(`(?msU)^\s*// </select-load-foreign-objects>$\n`),
			TokenizeTasks: []TokenizeTask{
				{
					Find:    regexp.MustCompile(`object.ID`),
					Replace: "object.{{ .PrimaryKeyStructField }}",
				},
				{
					Find:    regexp.MustCompile(`object.ParentPhysicalThingID`),
					Replace: "object.{{ .StructField }}",
				},
				{
					Find:    regexp.MustCompile(`object.ParentPhysicalThingIDObject`),
					Replace: "object.{{ .StructField }}Object",
				},
				{
					Find:    regexp.MustCompile(`SelectPhysicalThing`),
					Replace: "{{ .SelectFunc }}",
				},
				{
					Find:    regexp.MustCompile(`LogicalThing.ParentPhysicalThingIDObject`),
					Replace: "{{ .ObjectName }}.{{ .StructField }}Object",
				},
				{
					Find:    regexp.MustCompile(`PhysicalThingTablePrimaryKeyColumn`),
					Replace: "{{ .ForeignPrimaryKeyColumnVariable }}",
				},
				{
					Find:    regexp.MustCompile(`LogicalThingTableParentLogicalThingIDColumn`),
					Replace: "{{ .ObjectName }}Table{{ .StructField }}Column",
				},
				{
					Find:    regexp.MustCompile(`PhysicalThingTable`),
					Replace: "{{ .ForeignObjectName }}Table",
				},
				{
					Find:    regexp.MustCompile(`types\.IsZeroUUID`),
					Replace: "{{ .IsZeroFunc }}",
				},
			},
			KeepIsPerColumn:            true,
			KeepIsForPrimaryKeyOnly:    false,
			KeepIsForNonPrimaryKeyOnly: false,
			KeepIsForForeignKeysOnly:   true,
			KeepIsForSoftDeletableOnly: false,
			KeepIsForReferencedByOnly:  false,
		},

		{
			Name:      "SelectLoadReferencedByObjects",
			StartExpr: regexp.MustCompile(`(?ms)^[ |\t]*// <select-load-referenced-by-objects>$\n`),
			KeepExpr:  regexp.MustCompile(`(?msU)^\s*// <select-load-referenced-by-object>(\s+.*\n)\s+// </select-load-referenced-by-object>$\n`),
			EndExpr:   regexp.MustCompile(`(?msU)^\s*// </select-load-referenced-by-objects>$\n`),
			TokenizeTasks: []TokenizeTask{
				{
					Find:    regexp.MustCompile(`object.ID`),
					Replace: "object.{{ .PrimaryKeyStructField }}",
				},
				{
					Find:    regexp.MustCompile(`object.ReferencedByLogicalThingParentLogicalThingIDObjects`),
					Replace: "object.{{ .StructField }}",
				},
				{
					Find:    regexp.MustCompile(`SelectLogicalThings`),
					Replace: "{{ .SelectFunc }}",
				},
				{
					Find:    regexp.MustCompile(`LogicalThingTable, LogicalThingTableParentLogicalThingIDColumn, LogicalThingTable, object\.ID`),
					Replace: "{{ .ForeignPrimaryKeyTableVariable }}, {{ .ForeignPrimaryKeyColumnVariable }}, {{ .ObjectName }}Table, object.{{ .PrimaryKeyStructField }}",
				},
				{
					Find:    regexp.MustCompile(`LogicalThingTablePrimaryKeyColumn`),
					Replace: "{{ .ForeignPrimaryKeyColumnVariable }}",
				},
				{
					Find:    regexp.MustCompile(`LogicalThingTableColumnLookup`),
					Replace: "{{ .ForeignPrimaryKeyTableVariable }}ColumnLookup",
				},
				{
					Find:    regexp.MustCompile(`LogicalThingTableParentLogicalThingIDColumn`),
					Replace: "{{ .ForeignPrimaryKeyColumnVariable }}",
				},
				{
					Find:    regexp.MustCompile(`LogicalThingTable`),
					Replace: "{{ .ForeignPrimaryKeyTableVariable }}",
				},
			},
			KeepIsPerColumn:            false,
			KeepIsForPrimaryKeyOnly:    false,
			KeepIsForNonPrimaryKeyOnly: false,
			KeepIsForForeignKeysOnly:   false,
			KeepIsForSoftDeletableOnly: false,
			KeepIsForReferencedByOnly:  true,
		},

		{
			Name:      "InsertSetFieldPrimaryKey",
			StartExpr: regexp.MustCompile(`(?ms)^[ |\t]*// <insert-set-fields-primary-key>$\n`),
			KeepExpr:  regexp.MustCompile(`(?ms)^[ |\t]*// <insert-set-field-primary-key>$(.*)^[ |\t]*// </insert-set-field-primary-key>$\n`),
			EndExpr:   regexp.MustCompile(`(?msU)^[ |\t]*// </insert-set-fields-primary-key>$\n`),
			TokenizeTasks: []TokenizeTask{
				{
					Find:    regexp.MustCompile(`types\.IsZeroUUID`),
					Replace: "{{ .IsZeroFunc }}",
				},
				{
					Find:    regexp.MustCompile(`types\.FormatUUID`),
					Replace: "{{ .FormatFunc }}",
				},
				{
					Find:    regexp.MustCompile(`m\.ID`),
					Replace: "m.{{ .StructField }}",
				},
				{
					Find:    regexp.MustCompile(`LogicalThingTableIDColumn`),
					Replace: "LogicalThingTable{{ .StructField }}Column",
				},
			},
			KeepIsPerColumn:            true,
			KeepIsForPrimaryKeyOnly:    true,
			KeepIsForNonPrimaryKeyOnly: false,
			KeepIsForForeignKeysOnly:   false,
			KeepIsForSoftDeletableOnly: false,
			KeepIsForReferencedByOnly:  false,
		},

		{
			Name:      "InsertSetField",
			StartExpr: regexp.MustCompile(`(?ms)^[ |\t]*// <insert-set-fields>$\n`),
			KeepExpr:  regexp.MustCompile(`(?ms)^[ |\t]*// <insert-set-field>$(.*)^[ |\t]*// </insert-set-field>$\n`),
			EndExpr:   regexp.MustCompile(`(?msU)^[ |\t]*// </insert-set-fields>$\n`),
			TokenizeTasks: []TokenizeTask{
				{
					Find:    regexp.MustCompile(`types\.IsZeroTime`),
					Replace: "{{ .IsZeroFunc }}",
				},
				{
					Find:    regexp.MustCompile(`types\.FormatTime`),
					Replace: "{{ .FormatFunc }}",
				},
				{
					Find:    regexp.MustCompile(`m\.CreatedAt`),
					Replace: "m.{{ .StructField }}",
				},
				{
					Find:    regexp.MustCompile(`LogicalThingTableCreatedAtColumn`),
					Replace: "LogicalThingTable{{ .StructField }}Column",
				},
			},
			KeepIsPerColumn:            true,
			KeepIsForPrimaryKeyOnly:    false,
			KeepIsForNonPrimaryKeyOnly: true,
			KeepIsForForeignKeysOnly:   false,
			KeepIsForSoftDeletableOnly: false,
			KeepIsForReferencedByOnly:  false,
		},

		{
			Name:      "InsertSetPrimaryKey",
			StartExpr: regexp.MustCompile(`(?ms)^[ |\t]*// <insert-set-primary-key>$\n`),
			KeepExpr:  regexp.MustCompile(`(?ms)^[ |\t]*(.*)$\n`),
			EndExpr:   regexp.MustCompile(`(?ms)^[ |\t]*// </insert-set-primary-key>$\n`),
			TokenizeTasks: []TokenizeTask{
				{
					Find:    regexp.MustCompile(`types\.ParseUUID\(v\)`),
					Replace: "{{ .ParseFunc }}",
				},
				{
					Find:    regexp.MustCompile(`m\.ID`),
					Replace: "m.{{ .StructField }}",
				},
				{
					Find:    regexp.MustCompile(`LogicalThingTablePrimaryKeyColumn`),
					Replace: "LogicalThingTable{{ .StructField }}Column",
				},
				{
					Find:    regexp.MustCompile(`uuid.UUID`),
					Replace: "{{ .TypeTemplate }}",
				},
			},
			KeepIsPerColumn:            true,
			KeepIsForPrimaryKeyOnly:    true,
			KeepIsForNonPrimaryKeyOnly: false,
			KeepIsForForeignKeysOnly:   false,
			KeepIsForSoftDeletableOnly: false,
			KeepIsForReferencedByOnly:  false,
		},

		{
			Name:      "UpdateSetField",
			StartExpr: regexp.MustCompile(`(?ms)^[ |\t]*// <update-set-fields>$\n`),
			KeepExpr:  regexp.MustCompile(`(?ms)^[ |\t]*// <update-set-field>$(.*)^[ |\t]*// </update-set-field>$\n`),
			EndExpr:   regexp.MustCompile(`(?msU)^[ |\t]*// </update-set-fields>$\n`),
			TokenizeTasks: []TokenizeTask{
				{
					Find:    regexp.MustCompile(`types\.IsZeroTime`),
					Replace: "{{ .IsZeroFunc }}",
				},
				{
					Find:    regexp.MustCompile(`types\.FormatTime`),
					Replace: "{{ .FormatFunc }}",
				},
				{
					Find:    regexp.MustCompile(`m\.CreatedAt`),
					Replace: "m.{{ .StructField }}",
				},
				{
					Find:    regexp.MustCompile(`LogicalThingTableCreatedAtColumn`),
					Replace: "LogicalThingTable{{ .StructField }}Column",
				},
			},
			KeepIsPerColumn:            true,
			KeepIsForPrimaryKeyOnly:    false,
			KeepIsForNonPrimaryKeyOnly: true,
			KeepIsForForeignKeysOnly:   false,
			KeepIsForSoftDeletableOnly: false,
			KeepIsForReferencedByOnly:  false,
		},

		{
			Name:      "UpdateSetPrimaryKey",
			StartExpr: regexp.MustCompile(`(?ms)^[ |\t]*// <update-set-primary-key>$\n`),
			KeepExpr:  regexp.MustCompile(`(?ms)^[ |\t]*(.*)$\n`),
			EndExpr:   regexp.MustCompile(`(?ms)^[ |\t]*// </update-set-primary-key>$\n`),
			TokenizeTasks: []TokenizeTask{
				{
					Find:    regexp.MustCompile(`types\.FormatUUID`),
					Replace: "{{ .FormatFunc }}",
				},
				{
					Find:    regexp.MustCompile(`m\.ID`),
					Replace: "m.{{ .StructField }}",
				},
				{
					Find:    regexp.MustCompile(`LogicalThingTablePrimaryKeyColumn`),
					Replace: "LogicalThingTable{{ .StructField }}Column",
				},
			},
			KeepIsPerColumn:            true,
			KeepIsForPrimaryKeyOnly:    true,
			KeepIsForNonPrimaryKeyOnly: false,
			KeepIsForForeignKeysOnly:   false,
			KeepIsForSoftDeletableOnly: false,
			KeepIsForReferencedByOnly:  false,
		},

		{
			Name:      "DeleteSetPrimaryKey",
			StartExpr: regexp.MustCompile(`(?ms)^[ |\t]*// <delete-set-primary-key>$\n`),
			KeepExpr:  regexp.MustCompile(`(?ms)^[ |\t]*(.*)$\n`),
			EndExpr:   regexp.MustCompile(`(?ms)^[ |\t]*// </delete-set-primary-key>$\n`),
			TokenizeTasks: []TokenizeTask{
				{
					Find:    regexp.MustCompile(`types\.FormatUUID`),
					Replace: "{{ .FormatFunc }}",
				},
				{
					Find:    regexp.MustCompile(`m\.ID`),
					Replace: "m.{{ .StructField }}",
				},
				{
					Find:    regexp.MustCompile(`LogicalThingTablePrimaryKeyColumn`),
					Replace: "LogicalThingTable{{ .StructField }}Column",
				},
			},
			KeepIsPerColumn:            true,
			KeepIsForPrimaryKeyOnly:    true,
			KeepIsForNonPrimaryKeyOnly: false,
			KeepIsForForeignKeysOnly:   false,
			KeepIsForSoftDeletableOnly: false,
			KeepIsForReferencedByOnly:  false,
		},

		{
			Name:      "DeleteSoftDelete",
			StartExpr: regexp.MustCompile(`(?ms)^[ |\t]*// <delete-soft-delete>$\n`),
			KeepExpr:  regexp.MustCompile(`(?ms)^[ |\t]*(.*)$\n`),
			EndExpr:   regexp.MustCompile(`(?ms)^[ |\t]*// </delete-soft-delete>$\n`),
			TokenizeTasks: []TokenizeTask{
				{
					Find:    regexp.MustCompile(`(?ms)^[ |\t]*hardDelete.*\t\t}\n\t}`),
					Replace: "{{ .DeleteSoftDelete }}",
				},
			},
			KeepIsPerColumn:            false,
			KeepIsForPrimaryKeyOnly:    false,
			KeepIsForNonPrimaryKeyOnly: false,
			KeepIsForForeignKeysOnly:   false,
			KeepIsForSoftDeletableOnly: true,
			KeepIsForReferencedByOnly:  false,
		},

		{
			Name:      "ClaimRequest",
			StartExpr: regexp.MustCompile(`(?ms)^*// <claim-request>$\n`),
			KeepExpr:  regexp.MustCompile(`(?ms)^*(.*)$\n`),
			EndExpr:   regexp.MustCompile(`(?ms)^*// </claim-request>$\n`),
			TokenizeTasks: []TokenizeTask{
				{
					Find:    regexp.MustCompile(`LogicalThingClaimRequest`),
					Replace: "{{ .ObjectName }}{{ .ClaimPrefixPascalCase }}ClaimRequest",
				},
			},
			KeepIsPerColumn:            false,
			KeepIsForPrimaryKeyOnly:    false,
			KeepIsForNonPrimaryKeyOnly: false,
			KeepIsForForeignKeysOnly:   false,
			KeepIsForSoftDeletableOnly: false,
			KeepIsForReferencedByOnly:  false,
			KeepIsForClaimOnly:         true,
		},

		{
			Name:      "ClaimMethod",
			StartExpr: regexp.MustCompile(`(?ms)^*// <claim-method>$\n`),
			KeepExpr:  regexp.MustCompile(`(?ms)^*(.*)$\n`),
			EndExpr:   regexp.MustCompile(`(?ms)^*// </claim-method>$\n`),
			TokenizeTasks: []TokenizeTask{
				{
					Find:    regexp.MustCompile(`\) Claim\(`),
					Replace: ") {{ .ClaimPrefixPascalCase }}Claim(",
				},
				{
					Find:    regexp.MustCompile(`claimed_until`),
					Replace: "{{ .ClaimPrefixSnakeCase }}claimed_until",
				},
				{
					Find:    regexp.MustCompile(`m\.ClaimedUntil`),
					Replace: "m.{{ .ClaimPrefixPascalCase }}ClaimedUntil",
				},
			},
			KeepIsPerColumn:            false,
			KeepIsForPrimaryKeyOnly:    false,
			KeepIsForNonPrimaryKeyOnly: false,
			KeepIsForForeignKeysOnly:   false,
			KeepIsForSoftDeletableOnly: false,
			KeepIsForReferencedByOnly:  false,
			KeepIsForClaimOnly:         true,
		},

		{
			Name:      "ClaimFunc",
			StartExpr: regexp.MustCompile(`(?ms)^*// <claim-func>$\n`),
			KeepExpr:  regexp.MustCompile(`(?ms)^*(.*)$\n`),
			EndExpr:   regexp.MustCompile(`(?ms)^*// </claim-func>$\n`),
			TokenizeTasks: []TokenizeTask{
				{
					Find:    regexp.MustCompile(`func ClaimLogicalThing\(`),
					Replace: "func {{ .ClaimPrefixPascalCase }}Claim{{ .ObjectName }}(",
				},
				{
					Find:    regexp.MustCompile(`claimed_until`),
					Replace: "{{ .ClaimPrefixSnakeCase }}claimed_until",
				},
				{
					Find:    regexp.MustCompile(`SelectLogicalThings\(`),
					Replace: "Select{{ .ObjectNamePlural }}(",
				},
				{
					Find:    regexp.MustCompile(`m\.ClaimedUntil`),
					Replace: "m.{{ .ClaimPrefixPascalCase }}ClaimedUntil",
				},
			},
			KeepIsPerColumn:            false,
			KeepIsForPrimaryKeyOnly:    false,
			KeepIsForNonPrimaryKeyOnly: false,
			KeepIsForForeignKeysOnly:   false,
			KeepIsForSoftDeletableOnly: false,
			KeepIsForReferencedByOnly:  false,
			KeepIsForClaimOnly:         true,
		},

		{
			Name:      "ClaimHandlers",
			StartExpr: regexp.MustCompile(`(?ms)^*// <claim-handlers>$\n`),
			KeepExpr:  regexp.MustCompile(`(?ms)^*(.*)$\n`),
			EndExpr:   regexp.MustCompile(`(?ms)^*// </claim-handlers>$\n`),
			TokenizeTasks: []TokenizeTask{
				{
					Find:    regexp.MustCompile(`HandlerForClaim`),
					Replace: "HandlerFor{{ .ClaimPrefixPascalCase }}Claim",
				},
				{
					Find:    regexp.MustCompile(`req LogicalThingClaimRequest`),
					Replace: "req {{ .ObjectName }}{{ .ClaimPrefixPascalCase }}ClaimRequest",
				},
				{
					Find:    regexp.MustCompile(`claim-logical-thing`),
					Replace: "{{ .ClaimPrefixKebabCase }}claim-{{ .EndpointNameSingular }}",
				},
				{
					Find:    regexp.MustCompile(`/logical-things/{primaryKey}/claim`),
					Replace: "/{{ .EndpointName }}/{primaryKey}/{{ .ClaimPrefixKebabCase }}claim",
				},
				{
					Find:    regexp.MustCompile(`ClaimLogicalThing\(`),
					Replace: "{{ .ClaimPrefixPascalCase }}Claim{{ .ObjectName }}(",
				},
				{
					Find:    regexp.MustCompile(`object\.Claim\(`),
					Replace: "object.{{ .ClaimPrefixPascalCase }}Claim(",
				},
			},
			KeepIsPerColumn:            false,
			KeepIsForPrimaryKeyOnly:    false,
			KeepIsForNonPrimaryKeyOnly: false,
			KeepIsForForeignKeysOnly:   false,
			KeepIsForSoftDeletableOnly: false,
			KeepIsForReferencedByOnly:  false,
			KeepIsForClaimOnly:         true,
		},
	}

	for i, parseTask := range parseTasks {
		slices.SortFunc(parseTask.TokenizeTasks, func(a, b TokenizeTask) int {
			if len(a.Find.String()) < len(b.Find.String()) {
				return 1
			} else if len(a.Find.String()) > len(b.Find.String()) {
				return -1
			} else {
				return 0
			}
		})

		parseTask.StartVariableNameSet = make(map[string]struct{})
		parseTask.KeepVariableNameSet = make(map[string]struct{})
		parseTask.EndVariableNameSet = make(map[string]struct{})

		parseTasks[i] = parseTask
	}

	return parseTasks
}

var variableExpression = regexp.MustCompile(`{{ .\w+ }}`)

func init() {
	_ = getParseTasks() // just to check the expressions (because MustCompile)
}

func Parse() ([]ParseTask, error) {
	fileData := model_reference.ReferenceFileData

	baseTokenizeTasks := []TokenizeTask{
		{
			Find:    regexp.MustCompile(fmt.Sprintf(`%v`, model_reference.ReferenceTableName)),
			Replace: "{{ .TableName }}",
		},
		{
			Find:    regexp.MustCompile(fmt.Sprintf(`%v`, model_reference.ReferenceObjectName)),
			Replace: "{{ .ObjectName }}",
		},
	}

	parseTasks := getParseTasks()

	for i, parseTask := range parseTasks {
		startMatches := parseTask.StartExpr.FindStringSubmatch(fileData)
		if len(startMatches) > 1 {
			parseTask.StartMatch = startMatches[1] // prefer subgroup match
		} else if len(startMatches) == 1 {
			parseTask.StartMatch = startMatches[0] // fall back to full match
		}
		if strings.TrimSpace(parseTask.StartMatch) == "" {
			return nil, fmt.Errorf("failed to find a match for StartExpr: %v", parseTask.StartExpr)
		}

		startIndexes := parseTask.StartExpr.FindStringIndex(fileData)
		if len(startIndexes) != 2 {
			return nil, fmt.Errorf("failed to find exactly two indexes for StartExpr: %v", parseTask.StartExpr)
		}
		parseTask.Fragment = fileData[startIndexes[0]:]

		endMatches := parseTask.EndExpr.FindStringSubmatch(parseTask.Fragment)
		if len(endMatches) > 1 {
			parseTask.EndMatch = endMatches[len(endMatches)-1] // prefer subgroup match
		} else if len(endMatches) == 1 {
			parseTask.EndMatch = endMatches[0] // fall back to full match
		}
		if strings.TrimSpace(parseTask.EndMatch) == "" {
			return nil, fmt.Errorf("failed to find a match for EndExpr: %v", parseTask.EndExpr)
		}

		endIndexes := parseTask.EndExpr.FindStringIndex(parseTask.Fragment)
		if len(endIndexes) != 2 {
			return nil, fmt.Errorf("failed to find exactly two indexes for StartExpr: %v", parseTask.StartExpr)
		}
		parseTask.Fragment = parseTask.Fragment[:endIndexes[1]]

		keepMatches := parseTask.KeepExpr.FindStringSubmatch(parseTask.Fragment)
		if len(keepMatches) > 1 {
			parseTask.KeepMatch = keepMatches[1] // prefer subgroup match
		} else if len(keepMatches) == 1 {
			parseTask.KeepMatch = keepMatches[0] // fall back to full match
		}
		if strings.TrimSpace(parseTask.KeepMatch) == "" {
			return nil, fmt.Errorf("failed to find a match for KeepExpr: %v in:\n%v", parseTask.KeepExpr, parseTask.Fragment)
		}

		parseTask.ReplacedStartMatch = parseTask.StartMatch
		parseTask.ReplacedKeepMatch = parseTask.KeepMatch
		parseTask.ReplacedEndMatch = parseTask.EndMatch

		for _, tokenizeTask := range parseTask.TokenizeTasks {
			parseTask.ReplacedStartMatch = tokenizeTask.Find.ReplaceAllString(parseTask.ReplacedStartMatch, tokenizeTask.Replace)
			parseTask.ReplacedKeepMatch = tokenizeTask.Find.ReplaceAllString(parseTask.ReplacedKeepMatch, tokenizeTask.Replace)
			parseTask.ReplacedEndMatch = tokenizeTask.Find.ReplaceAllString(parseTask.ReplacedEndMatch, tokenizeTask.Replace)
		}

		for _, tokenizeTask := range baseTokenizeTasks {
			parseTask.ReplacedStartMatch = tokenizeTask.Find.ReplaceAllString(parseTask.ReplacedStartMatch, tokenizeTask.Replace)
			parseTask.ReplacedKeepMatch = tokenizeTask.Find.ReplaceAllString(parseTask.ReplacedKeepMatch, tokenizeTask.Replace)
			parseTask.ReplacedEndMatch = tokenizeTask.Find.ReplaceAllString(parseTask.ReplacedEndMatch, tokenizeTask.Replace)
		}

		for _, match := range variableExpression.FindAllString(parseTask.ReplacedStartMatch, -1) {
			parseTask.StartVariableNameSet[match] = struct{}{}
		}

		for _, match := range variableExpression.FindAllString(parseTask.ReplacedKeepMatch, -1) {
			parseTask.KeepVariableNameSet[match] = struct{}{}
		}

		for _, match := range variableExpression.FindAllString(parseTask.ReplacedEndMatch, -1) {
			parseTask.EndVariableNameSet[match] = struct{}{}
		}

		// fmt.Printf("----------------------------------------------------------------\n")
		// fmt.Printf("StartExpr: %v\n", parseTask.StartExpr)
		// fmt.Printf("KeepExpr: %v\n", parseTask.KeepExpr)
		// fmt.Printf("EndExpr: %v\n", parseTask.EndExpr)
		// fmt.Printf("TokenizeTasks:%v\n", parseTask.TokenizeTasks)
		// fmt.Printf("start1: %v\n", parseTask.StartMatch)
		// fmt.Printf("keep1 : %v\n", parseTask.KeepMatch)
		// fmt.Printf("end1  : %v\n", parseTask.EndMatch)
		// fmt.Printf("start2: %v\n", parseTask.ReplacedStartMatch)
		// fmt.Printf("keep2 : %v\n", parseTask.ReplacedKeepMatch)
		// fmt.Printf("end2  : %v\n", parseTask.ReplacedEndMatch)
		// fmt.Printf("# %v.Fragment\n%v\n", parseTask.Name, parseTask.Fragment)
		// fmt.Printf("# %v.ReplacedFragment\n%v\n", parseTask.Name, parseTask.ReplacedFragment)
		// fmt.Printf("# %v.Variables\n%v | %v | %v\n",
		// 	parseTask.Name,
		// 	maps.Keys(parseTask.StartVariableNameSet),
		// 	maps.Keys(parseTask.KeepVariableNameSet),
		// 	maps.Keys(parseTask.EndVariableNameSet),
		// )

		parseTasks[i] = parseTask
	}

	return parseTasks, nil
}
