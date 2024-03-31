package template

import "strings"

var deserializeFuncTemplate = strings.TrimSpace(`
func Deserialize%v(b []byte) (DjangolangObject, error) {
	var object %v

	err := json.Unmarshal(b, &object)
	if err != nil {
		return nil, err
	}

	return &object, nil
}
`) + "\n\n"
