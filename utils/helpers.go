package utils

import "encoding/json"

func Unmarshal[T any](s string) (*T, error) {
	out := new(T)
	if err := json.Unmarshal([]byte(s), out); err != nil {
		return nil, err
	}

	return out, nil
}
