package pg_types

// this is lifted directly from https://github.com/lib/pq and adjusted so that rather
// than being a struct containing a map, it IS a map

import (
	"database/sql"
	"database/sql/driver"
	"strings"
)

// Hstore is a wrapper for transferring Hstore values back and forth easily.
type Hstore map[string]sql.NullString

// escapes and quotes hstore keys/values
// s should be a sql.NullString or string
func hQuote(s interface{}) string {
	var str string
	switch v := s.(type) {
	case sql.NullString:
		if !v.Valid {
			return "NULL"
		}
		str = v.String
	case string:
		str = v
	default:
		panic("not a string or sql.NullString")
	}

	str = strings.Replace(str, "\\", "\\\\", -1)
	return `"` + strings.Replace(str, "\"", "\\\"", -1) + `"`
}

// Scan implements the Scanner interface.
//
// Note h.Map is reallocated before the scan to clear existing values. If the
// hstore column's database value is NULL, then h.Map is set to nil instead.
func (h Hstore) Scan(value interface{}) error {
	if value == nil {
		h = nil
		return nil
	}
	h = make(Hstore)
	var b byte
	pair := [][]byte{{}, {}}
	pi := 0
	inQuote := false
	didQuote := false
	sawSlash := false
	bindex := 0
	for bindex, b = range value.([]byte) {
		if sawSlash {
			pair[pi] = append(pair[pi], b)
			sawSlash = false
			continue
		}

		switch b {
		case '\\':
			sawSlash = true
			continue
		case '"':
			inQuote = !inQuote
			if !didQuote {
				didQuote = true
			}
			continue
		default:
			if !inQuote {
				switch b {
				case ' ', '\t', '\n', '\r':
					continue
				case '=':
					continue
				case '>':
					pi = 1
					didQuote = false
					continue
				case ',':
					s := string(pair[1])
					if !didQuote && len(s) == 4 && strings.ToLower(s) == "null" {
						h[string(pair[0])] = sql.NullString{String: "", Valid: false}
					} else {
						h[string(pair[0])] = sql.NullString{String: string(pair[1]), Valid: true}
					}
					pair[0] = []byte{}
					pair[1] = []byte{}
					pi = 0
					continue
				}
			}
		}
		pair[pi] = append(pair[pi], b)
	}
	if bindex > 0 {
		s := string(pair[1])
		if !didQuote && len(s) == 4 && strings.ToLower(s) == "null" {
			h[string(pair[0])] = sql.NullString{String: "", Valid: false}
		} else {
			h[string(pair[0])] = sql.NullString{String: string(pair[1]), Valid: true}
		}
	}
	return nil
}

// Value implements the driver Valuer interface. Note if h.Map is nil, the
// database column value will be set to NULL.
func (h Hstore) Value() (driver.Value, error) {
	if h == nil {
		return nil, nil
	}
	parts := []string{}
	for key, val := range h {
		thispart := hQuote(key) + "=>" + hQuote(val)
		parts = append(parts, thispart)
	}
	return []byte(strings.Join(parts, ",")), nil
}
