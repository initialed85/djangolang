package template

import "strings"

var primaryKeyFuncTemplate = strings.TrimSpace(`
func (%v *%v) GetPrimaryKey() (any, error) {
	return %v.%v, nil
}

func (%v *%v) SetPrimaryKey(value any) error {
	%v.%v = value.(%v)

	return nil
}
`) + "\n\n"

var primaryKeyFuncTemplateNotImplemented = strings.TrimSpace(`
func (%v *%v) GetPrimaryKey() (any, error) {
	return nil, fmt.Errorf("not implemented (table has no primary key)")
}

func (%v *%v) SetPrimaryKey(value any) error {
	return fmt.Errorf("not implemented (table has no primary key)")
}
`) + "\n\n"
