package schema

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v2"

	_ "embed"
)

//go:embed table.sql.tmpl
var tableTemplate string

//go:embed table_without_soft_delete.sql.tmpl
var tableWithSoftDeleteTemplate string

//go:embed unique.sql.tmpl
var uniqueTemplate string

//go:embed unique_without_soft_delete.sql.tmpl
var uniqueWithoutSoftDeleteTemplate string

//go:embed one_to_one.sql.tmpl
var oneToOneTemplate string

//go:embed one_to_many.sql.tmpl
var oneToManyTemplate string

// TODO: implement this- commenting for now to avoid unused variable warning
// //go:embed cascade_delete.sql.tmpl
// var cascadeDeleteTemplate string

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

	createdTableNames := make(map[string]struct{})

	for _, object := range schema.Objects {
		for _, claim := range object.Claims {
			claim = strings.TrimSpace(claim)

			claimedFor := ""
			if len(claim) > 0 {
				claimedFor = fmt.Sprintf("%s_", claim)
			}

			object.Properties = append(object.Properties, []string{
				fmt.Sprintf("%sclaimed_until timestamptz NOT NULL DEFAULT to_timestamp(0)", claimedFor),
			}...)
		}
	}

	additionalObjects := make([]*Object, 0)
	additionalRelationships := make([]*Relationship, 0)

	// TODO: use additionalObjects + additionalRelationships when we get a feature for it

	schema.Objects = append(schema.Objects, additionalObjects...)
	schema.Relationships = append(schema.Relationships, additionalRelationships...)

	for _, object := range schema.Objects {
		tableSQL := tableTemplate
		if object.WithoutSoftDelete {
			tableSQL = tableWithSoftDeleteTemplate
		}

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
			if !object.WithoutSoftDelete {
				tableSQL = strings.ReplaceAll(
					tableSQL,
					"deleted_at timestamptz NULL DEFAULT NULL,",
					"deleted_at timestamptz NULL DEFAULT NULL",
				)
			} else {
				tableSQL = strings.ReplaceAll(
					tableSQL,
					"id uuid PRIMARY KEY NOT NULL UNIQUE DEFAULT gen_random_uuid (),",
					"id uuid PRIMARY KEY NOT NULL UNIQUE DEFAULT gen_random_uuid ()",
				)
			}
		}

		tableSQL = strings.ReplaceAll(tableSQL, "        $$fields$$\n", fieldsSQL)

		sql += "--\n"
		sql += fmt.Sprintf("-- %s\n", object.Name)
		sql += "--\n\n"
		sql += fmt.Sprintf("%s\n\n", strings.TrimSpace(tableSQL))

		createdTableNames[object.Name] = struct{}{}
	}

	for _, relationship := range schema.Relationships {
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
		} else if relationship.Type == ManyToOne {
			relationshipSQL = oneToManyTemplate

			name := ""
			if relationship.Name != nil {
				name = fmt.Sprintf("%s_", *relationship.Name)
			}

			relationshipSQL = strings.ReplaceAll(relationshipSQL, "$$table_1$$", relationship.Source)
			relationshipSQL = strings.ReplaceAll(relationshipSQL, "$$column_1$$", "id")
			relationshipSQL = strings.ReplaceAll(relationshipSQL, "$$table_2$$", relationship.Destination)
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
		} else if relationship.Type == ManyToMany {
			relationshipTableName := fmt.Sprintf("m2m_%s_related_%s", relationship.Source, relationship.Destination)
			if relationship.Name != nil {
				relationshipTableName = *relationship.Name
			}

			_, tableAlreadyExists := createdTableNames[relationshipTableName]

			relationshipTableSQL := ""

			if !tableAlreadyExists {
				relationshipTableSQL = tableTemplate

				if !tableAlreadyExists {
					if dropFirst {
						relationshipTableSQL = "DROP TABLE IF EXISTS \"$$schema$$\".\"$$table_1$$\" CASCADE;\n\n" + tableTemplate
					}
				}

				relationshipTableSQL = strings.ReplaceAll(relationshipTableSQL, "$$table_1$$", relationshipTableName)

				relationshipTableSQL = strings.ReplaceAll(
					relationshipTableSQL,
					"deleted_at timestamptz NULL DEFAULT NULL,",
					"deleted_at timestamptz NULL DEFAULT NULL",
				)

				relationshipTableSQL = strings.ReplaceAll(relationshipTableSQL, "        $$fields$$\n", "")
			}

			relationshipSQL1 := oneToManyTemplate

			relationshipSQL1 = strings.ReplaceAll(relationshipSQL1, "$$table_1$$", relationshipTableName)
			relationshipSQL1 = strings.ReplaceAll(relationshipSQL1, "$$column_1$$", "id")
			relationshipSQL1 = strings.ReplaceAll(relationshipSQL1, "$$table_2$$", relationship.Destination)
			relationshipSQL1 = strings.ReplaceAll(relationshipSQL1, "$$column_2$$", "id")
			relationshipSQL1 = strings.ReplaceAll(relationshipSQL1, "$$name$$", "")

			relationshipSQL2 := oneToManyTemplate

			relationshipSQL2 = strings.ReplaceAll(relationshipSQL2, "$$table_1$$", relationshipTableName)
			relationshipSQL2 = strings.ReplaceAll(relationshipSQL2, "$$column_1$$", "id")
			relationshipSQL2 = strings.ReplaceAll(relationshipSQL2, "$$table_2$$", relationship.Source)
			relationshipSQL2 = strings.ReplaceAll(relationshipSQL2, "$$column_2$$", "id")
			relationshipSQL2 = strings.ReplaceAll(relationshipSQL2, "$$name$$", "")

			if !tableAlreadyExists {
				relationshipSQL += fmt.Sprintf("%s\n\n", relationshipTableSQL)
			}

			relationshipSQL += fmt.Sprintf("%s\n\n", relationshipSQL1)
			relationshipSQL += fmt.Sprintf("%s\n\n", relationshipSQL2)
		} else {
			return "", fmt.Errorf("unknown relationship type %#+v", relationship.Type)
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
			if object.WithoutSoftDelete {
				thisUniqueSQL = uniqueWithoutSoftDeleteTemplate
			}

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
