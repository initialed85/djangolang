package template

import "strings"

var loadTemplate = strings.TrimSpace(`
	idsFor%v := make([]any, 0)
	for _, id := range maps.Keys(%v) {
		idsFor%v = append(idsFor%v, id)
	}

	if len(idsFor%v) > 0 {
		rowsFor%v, err := Select%v(
			ctx,
			db,
			%vTransformedColumns,
			nil,
			nil,
			nil,
			types.Clause("%v IN $1", idsFor%v),
		)
		if err != nil {
			return nil, err
		}

		for _, row := range rowsFor%v {
			%v[%v] = row
		}

		for _, item := range items {
			item.%vObject = %v[%v]
		}
	}
`)
