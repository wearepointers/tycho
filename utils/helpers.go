package utils

import (
	"encoding/json"
	"reflect"
)

func Unmarshal[T any](s string) (*T, error) {
	out := new(T)
	if err := json.Unmarshal([]byte(s), out); err != nil {
		return nil, err
	}

	return out, nil
}

func IsSlice(v any) bool {
	return reflect.TypeOf(v).Kind() == reflect.Slice
}
