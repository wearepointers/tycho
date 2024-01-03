package utils

import (
	"encoding/json"
	"fmt"
	"strings"
	"unicode"
)

func Unmarshal[T any](s string) (*T, error) {
	out := new(T)
	if err := json.Unmarshal([]byte(s), out); err != nil {
		return nil, err
	}

	return out, nil
}

func Join(elems []string, sep string, prefix string, suffix string) string {
	switch len(elems) {
	case 0:
		return ""
	case 1:
		return fmt.Sprint(prefix, elems[0], suffix)
	}
	n := len(sep) * (len(elems) - 1)
	for i := 0; i < len(elems); i++ {
		n += len(elems[i])
	}

	var b strings.Builder
	b.Grow(n)
	b.WriteString(fmt.Sprint(prefix, elems[0], suffix))
	for _, s := range elems[1:] {
		b.WriteString(sep)
		b.WriteString(fmt.Sprint(prefix, s, suffix))
	}
	return b.String()
}

func SnakeToPascal(s string) string {
	var result []rune

	for i, r := range s {
		if i == 0 {
			result = append(result, unicode.ToUpper(r))
		} else if r == '_' {
			continue
		} else if i > 0 && s[i-1] == '_' {
			result = append(result, unicode.ToUpper(r))
		} else {
			result = append(result, r)
		}
	}

	return string(result)
}
