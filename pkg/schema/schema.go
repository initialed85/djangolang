package schema

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v2"

	_ "embed"
)

//go:embed table.sql.tmpl
var tableTemplate string

//go:embed unique.sql.tmpl
var uniqueTemplate string

//go:embed one_to_one.sql.tmpl
var oneToOneTemplate string

//go:embed one_to_many.sql.tmpl
var oneToManyTemplate string

func Parse(b []byte) (*Schema, error) {
	var schema Schema

	err := yaml.Unmarshal(b, &schema)
	if err != nil {
		return nil, err
	}

	return &schema, nil
}

func Dump(schema *Schema, schemaName string, dropFirsts ...bool) (string, error) {
	if schema == nil {
		return "", fmt.Errorf("schema unexpectedly nil")
	}

	dropFirst := len(dropFirsts) > 0 && dropFirsts[0]

	sql := ""

	sql += "CREATE SCHEMA IF NOT EXISTS $$schema$$;\n\n"

	sql += "SET LOCAL search_path = $$schema$$;\n\n"

	for _, object := range schema.Objects {
		tableSQL := tableTemplate

		if dropFirst {
			tableSQL = "DROP TABLE IF EXISTS \"$$schema$$\".\"$$table_1$$\" CASCADE;\n\n" + tableSQL
		}

		tableSQL = strings.ReplaceAll(tableSQL, "$$table_1$$", object.Name)

		var fieldsSQL string
		for i, property := range object.Properties {
			suffix := ","
			if i == len(object.Properties)-1 {
				suffix = ""
			}

			fieldsSQL += fmt.Sprintf("        %s%s\n", property, suffix)
		}

		if len(object.Properties) == 0 {
			tableSQL = strings.ReplaceAll(
				tableSQL,
				"deleted_at timestamptz NULL DEFAULT NULL,",
				"deleted_at timestamptz NULL DEFAULT NULL",
			)
		}

		tableSQL = strings.ReplaceAll(tableSQL, "        $$fields$$\n", fieldsSQL)

		sql += "--\n"
		sql += fmt.Sprintf("-- %s\n", object.Name)
		sql += "--\n\n"
		sql += fmt.Sprintf("%s\n\n", strings.TrimSpace(tableSQL))
	}

	for _, relationship := range schema.Relationships {
		if relationship.Type == ManyToOne {
			relationship = &Relationship{
				Name:        relationship.Name,
				Source:      relationship.Destination,
				Destination: relationship.Source,
				Type:        relationship.Type,
			}
		}

		var relationshipSQL string

		if relationship.Type == OneToMany {
			relationshipSQL = oneToManyTemplate

			name := ""
			if relationship.Name != nil {
				name = fmt.Sprintf("%s_", *relationship.Name)
			}

			relationshipSQL = strings.ReplaceAll(relationshipSQL, "$$table_1$$", relationship.Destination)
			relationshipSQL = strings.ReplaceAll(relationshipSQL, "$$column_1$$", "id")
			relationshipSQL = strings.ReplaceAll(relationshipSQL, "$$table_2$$", relationship.Source)
			relationshipSQL = strings.ReplaceAll(relationshipSQL, "$$column_2$$", "id")
			relationshipSQL = strings.ReplaceAll(relationshipSQL, "$$name$$", name)

			if relationship.Optional {
				relationshipSQL = strings.ReplaceAll(relationshipSQL, "NOT NULL", "NULL")
			}
		} else if relationship.Type == OneToOne {
			relationshipSQL = oneToOneTemplate

			name := ""
			if relationship.Name != nil {
				name = fmt.Sprintf("%s_", *relationship.Name)
			}

			relationshipSQL = strings.ReplaceAll(relationshipSQL, "$$table_1$$", relationship.Destination)
			relationshipSQL = strings.ReplaceAll(relationshipSQL, "$$column_1$$", "id")
			relationshipSQL = strings.ReplaceAll(relationshipSQL, "$$table_2$$", relationship.Source)
			relationshipSQL = strings.ReplaceAll(relationshipSQL, "$$column_2$$", "id")
			relationshipSQL = strings.ReplaceAll(relationshipSQL, "$$name$$", name)

			if relationship.Optional {
				relationshipSQL = strings.ReplaceAll(relationshipSQL, "NOT NULL", "NULL")
			}
		} else {
			infix := "relates_to"
			if relationship.Name != nil {
				infix = *relationship.Name
			}

			relationshipTableName := fmt.Sprintf("m2m_%s_%s_%s", relationship.Source, infix, relationship.Destination)

			relationshipTableSQL := tableTemplate

			if dropFirst {
				relationshipTableSQL = "DROP TABLE IF EXISTS \"$$schema$$\".\"$$table_1$$\" CASCADE;\n\n" + tableTemplate
			}

			relationshipTableSQL = strings.ReplaceAll(relationshipTableSQL, "$$table_1$$", relationshipTableName)

			relationshipTableSQL = strings.ReplaceAll(
				relationshipTableSQL,
				"deleted_at timestamptz NULL DEFAULT NULL,",
				"deleted_at timestamptz NULL DEFAULT NULL",
			)

			relationshipTableSQL = strings.ReplaceAll(relationshipTableSQL, "        $$fields$$\n", "")

			relationshipSQL1 := oneToManyTemplate

			sourceColumn1 := fmt.Sprintf("%s_id", relationship.Destination)
			if relationship.Name != nil {
				sourceColumn1 = fmt.Sprintf("%s_%s", strings.TrimRight(*relationship.Name, "_"), sourceColumn1)
			}

			relationshipSQL1 = strings.ReplaceAll(relationshipSQL1, "$$table_1$$", relationship.Source)
			relationshipSQL1 = strings.ReplaceAll(relationshipSQL1, "$$column_1$$", sourceColumn1)
			relationshipSQL1 = strings.ReplaceAll(relationshipSQL1, "$$table_2$$", relationship.Destination)
			relationshipSQL1 = strings.ReplaceAll(relationshipSQL1, "$$column_2$$", "id")
			relationshipSQL1 = strings.ReplaceAll(relationshipSQL1, "$$name$$", "")

			relationshipSQL2 := oneToManyTemplate

			sourceColumn2 := fmt.Sprintf("%s_id", relationship.Source)
			if relationship.Name != nil {
				sourceColumn2 = fmt.Sprintf("%s_%s", strings.TrimRight(*relationship.Name, "_"), sourceColumn2)
			}

			relationshipSQL2 = strings.ReplaceAll(relationshipSQL2, "$$table_1$$", relationship.Destination)
			relationshipSQL2 = strings.ReplaceAll(relationshipSQL2, "$$column_1$$", sourceColumn2)
			relationshipSQL2 = strings.ReplaceAll(relationshipSQL2, "$$table_2$$", relationship.Source)
			relationshipSQL2 = strings.ReplaceAll(relationshipSQL2, "$$column_2$$", "id")
			relationshipSQL2 = strings.ReplaceAll(relationshipSQL2, "$$name$$", "")

			relationshipSQL += fmt.Sprintf("%s\n\n", relationshipTableSQL)
			relationshipSQL += fmt.Sprintf("%s\n\n", relationshipSQL1)
			relationshipSQL += fmt.Sprintf("%s\n\n", relationshipSQL2)
		}

		sql += "--\n"
		sql += fmt.Sprintf("-- %s -> %s (%s)\n", relationship.Source, relationship.Destination, relationship.Type)
		sql += "--\n\n"
		sql += fmt.Sprintf("%s\n\n", strings.TrimSpace(relationshipSQL))
	}

	for _, object := range schema.Objects {
		if len(object.UniqueOn) == 0 {
			continue
		}

		uniqueSQL := ""

		for _, uniqueOns := range object.UniqueOn {
			thisUniqueSQL := uniqueTemplate

			thisUniqueSQL = strings.ReplaceAll(thisUniqueSQL, "$$table_1$$", object.Name)
			thisUniqueSQL = strings.ReplaceAll(
				thisUniqueSQL,
				"$$name$$",
				strings.Join(uniqueOns, "_"),
			)
			thisUniqueSQL = strings.ReplaceAll(
				thisUniqueSQL,
				"$$columns$$",
				strings.Join(uniqueOns, ", "),
			)

			uniqueSQL += "--\n"
			uniqueSQL += fmt.Sprintf("-- %s unique on (%s)\n", object.Name, strings.Join(uniqueOns, ", "))
			uniqueSQL += "--\n\n"

			uniqueSQL += fmt.Sprintf("%s\n\n", thisUniqueSQL)
		}

		sql += uniqueSQL
	}

	sql = strings.ReplaceAll(sql, "$$schema$$", schemaName)

	return sql, nil
}
