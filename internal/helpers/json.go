package helpers

import (
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func UnsafeJSONPrettyFormat(v any) string {
	b, _ := json.MarshalIndent(v, "", "  ")

	return string(b)
}

func UnsafeJSONFormat(v any) string {
	b, _ := json.Marshal(v)

	return string(b)
}
