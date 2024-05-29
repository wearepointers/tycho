package sql

import (
	"fmt"

	"github.com/wearepointers/tycho/utils"
)

func WhereComposite[T any](op Operator, cols []T, f func(c T) string) string {
	var arguments = make([]string, len(cols))
	var columns = make([]string, len(cols))
	for i, c := range cols {
		arguments[i] = "?"
		columns[i] = f(c)
	}

	return Query([]string{Expr(Group(columns...)), op.String(), Expr(Group(arguments...))}...)
}

func Where(s string, op Operator, v ...string) string {
	j := []string{s, op.String()}
	if len(v) > 0 {
		j = append(j, v...)
	}

	return Query(j...)
}

func WhereIn(s string, v any) (string, any) {
	isSlice, val := utils.IsSlice[any](v)
	if !isSlice {
		return "", nil
	}

	count := len(val)
	placeholders := make([]string, count)
	for i := 0; i < count; i++ {
		placeholders[i] = "?"
	}

	return Where(s, In, fmt.Sprintf("(%s)", Group(placeholders...))), v
}

func WhereNotIn(s string, v any) (string, any) {
	isSlice, val := utils.IsSlice[any](v)
	if !isSlice {
		return "", nil
	}

	count := len(val)
	placeholders := make([]string, count)
	for i := 0; i < count; i++ {
		placeholders[i] = "?"
	}

	return Where(s, NotIn, fmt.Sprintf("(%s)", Group(placeholders...))), v
}

func WhereLike(s string, v any) (string, any) {
	return Where(s, Like, "?"), fmt.Sprint("%", v, "%")
}

func WhereNotLike(s string, v any) (string, any) {
	return Where(s, NotLike, "?"), fmt.Sprint("%", v, "%")
}

func WhereStartsWith(s string, v any) (string, any) {
	return Where(s, Like, "?"), fmt.Sprint(v, "%")
}

func WhereEndsWith(s string, v any) (string, any) {
	return Where(s, Like, "?"), fmt.Sprint("%", v)
}
