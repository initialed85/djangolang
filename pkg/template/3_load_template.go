package template

import "strings"

var loadTemplate = strings.TrimSpace(`
	idsFor%v := make([]string, 0)
	for _, id := range maps.Keys(%v) {
		b, err := json.Marshal(id)
		if err != nil {
			return nil, err
		}

		s := strings.ReplaceAll(string(b), "\"", "'")

		idsFor%v = append(idsFor%v, s)
	}

	if len(idsFor%v) > 0 {
		rowsFor%v, err := Select%v(
			ctx,
			db,
			%vTransformedColumns,
			nil,
			nil,
			nil,
			fmt.Sprintf("%v IN (%%v)", strings.Join(idsFor%v, ", ")),
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
