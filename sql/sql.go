package sql

import (
	"fmt"
	"strings"
)

func Clause(op LogicalOperator, mods ...string) string {
	return strings.Join(mods, op.String())
}

func Expr(s string) string {
	return fmt.Sprint("(", s, ")")
}

func Args(sql, prefix string) string {
	count := strings.Count(sql, "?")

	for i := 1; i <= count; i++ {
		sql = strings.Replace(sql, "?", fmt.Sprintf("%v%d", prefix, i), 1)
	}

	return sql
}
