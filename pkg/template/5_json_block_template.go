package template

import "strings"

var jsonBlockTemplate = strings.Trim(`
        var temp1%v any
		var temp2%v []any

		if item.%v != nil {
			err = json.Unmarshal(item.%v.([]byte), &temp1%v)
			if err != nil {
				err = json.Unmarshal(item.%v.([]byte), &temp2%v)
				if err != nil {
					item.%v = nil
				} else {
					item.%v = &temp2%v
				}
			} else {
				item.%v = &temp1%v
			}
		}

		item.%v = &temp1%v
`, "\n") + "\n\n"
