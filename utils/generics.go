package utils

import (
	"encoding/json"
)

func Unmarshal[T any](s string) (*T, error) {
	out := new(T)
	if err := json.Unmarshal([]byte(s), out); err != nil {
		return nil, err
	}

	return out, nil
}

func StructToJSON[T any](v T) string {
	b, _ := json.MarshalIndent(v, "", "  ")
	return string(b)
}

func OmitemptyMap(v map[string]any) map[string]any {
	for key, value := range v {
		if value == nil {
			delete(v, key)
		}

		isString, s := IsString(value)
		if isString && s == "" {
			delete(v, key)
		}

		isBool, b := IsBool(value)
		if isBool && !b {
			delete(v, key)
		}
	}

	return v
}
